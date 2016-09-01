package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var G_GlobalVariables TGlobalVariables

type TActivityLst []TActivityData

type TActivityData struct {
	ActivityID   int   //! 唯一活动ID
	activityType int   //! 活动所用类型模板
	beginTime    int64 //! 活动开启时间
	endTime      int64 //! 活动结束时间
	award        int   //! 当前活动使用奖励版本
	Status       int   //! 状态: 1:有效活动，0:无效活动。
	VersionCode  int32 //! 活动刷新版本号
	ResetCode    int32 //! 活动迭代版本号
}

type TGroupPurchaseInfo struct {
	ItemID  int //! 道具ID
	SaleNum int //! 销售数量
}

type TSevenDayBuyInfo struct {
	ActivityID int
	LimitBuy   [7]int
}

type TGlobalVariables struct {
	ID               int32                `bson:"_id"`
	NormalMoneyPoor  int                  //! 奖金池
	ExcitedMoneyPoor int                  //! 豪华奖金池
	GroupPurchaseLst []TGroupPurchaseInfo //! 团购货物列表
	SevenDayLimit    []TSevenDayBuyInfo   //! 七日活动已购买限购的人数列表
	LimitSaleNum     int                  //! 限时特惠道具购买人次

	ActivityLst TActivityLst //! 活动列表

	SvrAwardIncID int // 自增ID
	SvrAwardList  []TAwardData
}

func (self *TGlobalVariables) Init() {
	if self.DB_LoadGlobalVariables() == false {
		//! 未找到数据则初始化

		self.ID = 1
		self.NormalMoneyPoor = 0
		self.ExcitedMoneyPoor = 0

		//! 初始化七天活动限购
		self.InitSevenDayBuyLst()

		//! 获取今日开启活动
		openDay := GetOpenServerDay()
		for _, v := range gamedata.GT_ActivityLst {
			if v.ID == 0 {
				gamelog.Error("TGlobalVariables::Init Error Invalid ActivityID:%d", v.ID)
				continue
			}

			if v.ActivityType == gamedata.Activity_Seven {
				seven := TSevenDayBuyInfo{}
				seven.ActivityID = v.ID
				self.SevenDayLimit = append(self.SevenDayLimit, seven)
				//go self.DB_AddSevenDayBuyInfo(seven)
			}

			var activity TActivityData
			activity.ActivityID = v.ID
			activity.activityType = v.ActivityType
			activity.award = v.AwardType
			activity.beginTime, activity.endTime = gamedata.GetActivityEndTime(v.ID, openDay)
			activity.VersionCode = 0
			activity.ResetCode = 0
			activity.Status = v.Status
			self.ActivityLst = append(self.ActivityLst, activity)
		}

		mongodb.InsertToDB(appconfig.GameDbName, "GlobalVariables", self)
	}

	//! 检测新加活动select
	CheckActivityAdd()

	//! 计算活动状态
	self.CalcActivityTime()

	//! 赋值活动奖励模板
	self.SetActivityAwardType()

}

func (self *TGlobalVariables) CorrectionTime() {
	//! 特殊矫正活动时间
	for i, v := range self.ActivityLst {
		activityInfo := gamedata.GetActivityInfo(v.ActivityID)
		if activityInfo.ServerType == 2 && activityInfo.ActivityType == gamedata.Activity_Sign {
			//! 公服期签到为永久存在
			self.ActivityLst[i].beginTime = 0
			self.ActivityLst[i].endTime = 0
		}
	}
}

//! 获取活动奖励
func (self *TGlobalVariables) GetActivityAwardType(activityID int) int {
	//! 根据活动模板获取对应ID
	for i := 0; i < len(self.ActivityLst); i++ {
		if self.ActivityLst[i].ActivityID == activityID {
			return self.ActivityLst[i].award
		}
	}

	gamelog.Error("GetActivityAwardType Error : Invalid activityID:%d", activityID)
	return 0
}

//! 判断活动是否开启
func (self *TGlobalVariables) IsActivityOpen(activityID int) bool {
	now := time.Now().Unix()
	length := len(self.ActivityLst)
	for i := 0; i < length; i++ {
		data := &self.ActivityLst[i]
		if activityID == data.ActivityID &&
			data.Status == 1 &&
			((data.beginTime <= now && now <= data.endTime) || (data.beginTime == 0 && data.endTime == 0)) {
			return true
		}
	}

	return false
}

//! 判断当前是否为活动时间
//! 返回: 是否在活动操作期(有的专门设置了领奖期)    结束倒计时
func (self *TGlobalVariables) IsActivityTime(activityID int) (bool, int) {
	var endCountDown int = -1
	now := time.Now()
	for _, v := range G_GlobalVariables.ActivityLst {
		activityInfo := gamedata.GetActivityInfo(v.ActivityID)
		if G_GlobalVariables.IsActivityOpen(v.ActivityID) == true && activityID == v.ActivityID {
			if v.beginTime == 0 && v.endTime == 0 {
				return true, 0 //! 永久开启
			}

			endCountDown = int(v.endTime) - activityInfo.AwardTime*24*60*60

			// gamelog.Info("EndTime: %v    endCountDown: %v", v.endTime, endCountDown)
			break
		}
	}

	if int64(endCountDown) <= now.Unix() {
		return false, 0
	}

	return true, endCountDown
}

//! 判断活动是否为有效活动
func (self *TGlobalVariables) IsActivityValid(activityID int) bool {
	for _, v := range self.ActivityLst {
		if (activityID == v.ActivityID) && (v.Status == 1) {
			return true
		}
	}

	return false
}

func (self *TGlobalVariables) SetActivityAwardType() {
	for i, v := range self.ActivityLst {
		activityInfo := gamedata.GetActivityInfo(v.ActivityID)
		self.ActivityLst[i].award = activityInfo.AwardType
		self.ActivityLst[i].activityType = activityInfo.ActivityType
	}
}

//! 计算开启时间与关闭时间
func (self *TGlobalVariables) CalcActivityTime() {
	openDay := GetOpenServerDay()
	for i, v := range self.ActivityLst {
		if v.ActivityID == 0 {
			gamelog.Error("CalcActivityTime Error Invalid ActivityID:%d", v.ActivityID)
			continue
		}
		self.ActivityLst[i].beginTime, self.ActivityLst[i].endTime = gamedata.GetActivityEndTime(v.ActivityID, openDay)

		// gamelog.Info("ActivityID: %v, BeginTime: %v  EndTime: %v", self.ActivityLst[i].ActivityID,
		// 	self.ActivityLst[i].beginTime,
		// 	self.ActivityLst[i].endTime)

	}
}

func (self *TGlobalVariables) GetGroupPurchaseItemInfo(itemID int) (*TGroupPurchaseInfo, int) {
	length := len(self.GroupPurchaseLst)
	for i := 0; i < length; i++ {
		if self.GroupPurchaseLst[i].ItemID == itemID {
			return &self.GroupPurchaseLst[i], i
		}
	}

	//! 不存在物品信息, 初始化
	var newRecord TGroupPurchaseInfo
	newRecord.ItemID = itemID
	newRecord.SaleNum = 0
	self.GroupPurchaseLst = append(self.GroupPurchaseLst, newRecord)
	go self.DB_AddNewGroupPurchaseRecord(&newRecord)

	return &self.GroupPurchaseLst[length], length
}

func (self *TGlobalVariables) GetSevenDayLimit(activityID int) *TSevenDayBuyInfo {
	for i := 0; i < len(G_GlobalVariables.SevenDayLimit); i++ {
		if G_GlobalVariables.SevenDayLimit[i].ActivityID == activityID {
			return &G_GlobalVariables.SevenDayLimit[i]
		}
	}

	return nil
}

func (self *TGlobalVariables) AddSevenDayLimit(activityID int, index int) {
	for i := 0; i < len(G_GlobalVariables.SevenDayLimit); i++ {
		if G_GlobalVariables.SevenDayLimit[i].ActivityID == activityID {
			G_GlobalVariables.SevenDayLimit[i].LimitBuy[index] += 1
			go self.DB_SaveSevenDayLimit(index)
			break
		}
	}
}

func (self *TGlobalVariables) InitSevenDayBuyLst() bool {
	sevenDayLst := []TActivityModule{}
	s := mongodb.GetDBSession()
	defer s.Close()

	index := 0
	for i := 0; i < len(G_GlobalVariables.ActivityLst); i++ {
		if G_GlobalVariables.ActivityLst[i].activityType == gamedata.Activity_Seven {
			if G_GlobalVariables.IsActivityOpen(G_GlobalVariables.ActivityLst[i].ActivityID) == true {
				filedName := fmt.Sprintf("sevenday.%d.buylst", index)
				err := s.DB(appconfig.GameDbName).C("PlayerActivity").Find(bson.M{filedName: bson.M{"$exists": true}}).All(&sevenDayLst)
				if err != nil {
					if err != mgo.ErrNotFound {
						gamelog.Error("Init DB Error!!!")
						return false
					}
				}

				for _, v := range sevenDayLst {
					for _, n := range v.SevenDay[index].BuyLst {
						self.SevenDayLimit[index].LimitBuy[n-1] += 1
					}
				}
			}
			index++
		}
	}

	return true
}

func (self *TGlobalVariables) AddGroupPurchaseRecord(itemID int, saleNum int) int {
	//! 获取团购记录
	recordInfo, index := self.GetGroupPurchaseItemInfo(itemID)

	//! 添加购买数
	recordInfo.SaleNum += saleNum
	go self.DB_UpdateGroupPurchaseSaleNum(index)

	return recordInfo.SaleNum
}

//! 添加团购记录
func (self *TGlobalVariables) DB_AddNewGroupPurchaseRecord(record *TGroupPurchaseInfo) {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$push": bson.M{"grouppurchaselst": *record}})
}

//! 更新团购商品购买记录
func (self *TGlobalVariables) DB_UpdateGroupPurchaseSaleNum(index int) {
	filedName := fmt.Sprintf("grouppurchaselst.%d.salenum", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
		filedName: self.GroupPurchaseLst[index].SaleNum}})
}

func (self *TGlobalVariables) DB_CleanMoneyPoor() {
	if self.NormalMoneyPoor != 0 || self.ExcitedMoneyPoor != 0 {
		self.NormalMoneyPoor = 0
		self.ExcitedMoneyPoor = 0
		mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
			"normalmoneypoor":  self.NormalMoneyPoor,
			"excitedmoneypoor": self.ExcitedMoneyPoor}})
	}
}

func (self *TGlobalVariables) DB_CleanGroupPurchase() {
	if len(self.GroupPurchaseLst) != 0 {
		self.GroupPurchaseLst = []TGroupPurchaseInfo{}

		mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
			"grouppurchaselst": self.GroupPurchaseLst}})
	}
}

func (self *TGlobalVariables) DB_SaveMoneyPoor() {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
		"normalmoneypoor":  self.NormalMoneyPoor,
		"excitedmoneypoor": self.ExcitedMoneyPoor}})
}

func (self *TGlobalVariables) DB_LoadGlobalVariables() bool {
	if mongodb.Find(appconfig.GameDbName, "GlobalVariables", "_id", 1, self) != 0 {
		return false
	}
	return true
}

func (self *TGlobalVariables) DB_SaveGlobalVariables() {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": self})
}

func (self *TGlobalVariables) DB_AddNewActivity(activity TActivityData) {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$push": bson.M{"activitylst": activity}})
}

func (self *TGlobalVariables) DB_UpdateActivityLst() {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
		"activitylst": self.ActivityLst}})
}

func (self *TGlobalVariables) DB_AddSevenDayBuyInfo(seven TSevenDayBuyInfo) {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$push": bson.M{"sevendaylimit": seven}})
}

func (self *TGlobalVariables) DB_CleanSevenDayInfo(activityID int) {
	index := -1
	for i := 0; i < len(self.SevenDayLimit); i++ {
		if self.SevenDayLimit[i].ActivityID == activityID {
			self.SevenDayLimit[i].LimitBuy = [7]int{0, 0, 0, 0, 0, 0, 0}
			index = i
			break
		}
	}

	if index < 0 {
		return
	}

	filedName := fmt.Sprintf("sevendaylimit.%d.limitbuy", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{filedName: self.SevenDayLimit[index].LimitBuy}})
}

func (self *TGlobalVariables) DB_SaveSevenDayLimit(index int) {
	filedName := fmt.Sprintf("sevendaylimit.%d", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{filedName: self.SevenDayLimit[index]}})
}

func (self *TGlobalVariables) DB_UpdateActivityInfo(index int) {
	filedName := fmt.Sprintf("activitylst.%d.status", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
		filedName: self.ActivityLst[index].Status}})
}

// GM调用的增删全服奖励接口
func (self *TGlobalVariables) AddSvrAward(pAwardData *TAwardData) {
	self.SvrAwardIncID += 1
	pAwardData.ID = self.SvrAwardIncID
	pAwardData.Time = time.Now().Unix()
	self.SvrAwardList = append(self.SvrAwardList, *pAwardData)
	self.DB_AddAward(pAwardData)
	self.DB_SaveIncrementID()
}
func (self *TGlobalVariables) DelSvrAward(id int) {
	for i, v := range self.SvrAwardList {
		if v.ID == id {
			self.SvrAwardList = append(self.SvrAwardList[:i], self.SvrAwardList[i+1:]...)
			self.DB_DelAward(id)
			break
		}
	}
}
func (self *TGlobalVariables) DB_AddAward(pAwardData *TAwardData) {
	mongodb.UpdateToDB(appconfig.GameDbName, "SvrAwardCenter", bson.M{"_id": 0}, bson.M{"$push": bson.M{"svrawardlist": *pAwardData}})
}
func (self *TGlobalVariables) DB_DelAward(id int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "SvrAwardCenter", bson.M{"_id": 0}, bson.M{"$pull": bson.M{"svrawardlist": bson.M{"id": id}}})
}
func (self *TGlobalVariables) DB_SaveIncrementID() {
	mongodb.UpdateToDB(appconfig.GameDbName, "SvrAwardCenter", bson.M{"_id": 0}, bson.M{"$set": bson.M{"svrawardincid": self.SvrAwardIncID}})
}

func (self *TGlobalVariables) DB_UpdateLimitSaleNum() {
	mongodb.UpdateToDB(appconfig.GameDbName, "GlobalVariables", bson.M{"_id": 1}, bson.M{"$set": bson.M{
		"limitsalenum": self.LimitSaleNum}})
}

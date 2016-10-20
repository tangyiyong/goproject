package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"time"
	"utility"
)

var G_GlobalVariables TGlobalVariables

type TActivityData struct {
	ActivityID  int32 //! 唯一活动ID
	VersionCode int32 //! 活动刷新版本号
	ResetCode   int32 //! 活动迭代版本号

	activityType int //! 活动所用类型模板
	award        int //! 当前活动使用奖励版本
	Status       int //! 状态: 1:有效活动，0:无效活动。

	beginTime  int32 //! 活动开启时间
	actEndTime int32 //! 活动操作结束时间
	endTime    int32 //! 领奖结束时间
}

type TGroupPurchaseInfo struct {
	ItemID  int //! 道具ID
	SaleNum int //! 销售数量
}

type TSevenDayBuyInfo struct {
	ActivityID int32
	LimitBuy   [7]int
}

type TGlobalVariables struct {
	ID               int32                `bson:"_id"`
	NormalMoneyPoor  int                  //! 奖金池
	ExcitedMoneyPoor int                  //! 豪华奖金池
	GroupPurchaseLst []TGroupPurchaseInfo //! 团购货物列表
	SevenDayLimit    []TSevenDayBuyInfo   //! 七日活动已购买限购的人数列表
	LimitSaleNum     int                  //! 限时特惠道具购买人次

	ActivityLst   []TActivityData //! 活动列表
	SvrAwardIncID int             // 自增ID
	SvrAwardList  []TAwardData
}

func (self *TGlobalVariables) Init() {
	if self.DB_LoadGlobalVariables() == false {
		self.ID = 1
		self.NormalMoneyPoor = 0
		self.ExcitedMoneyPoor = 0

		//! 初始化七天活动限购
		self.InitSevenDayBuyLst()

		for _, v := range gamedata.GT_ActivityLst {
			if v.ID == 0 {
				gamelog.Error("TGlobalVariables::Init Error Invalid ActivityID:%d", v.ID)
				continue
			}

			if v.ActType == gamedata.Activity_Seven {
				seven := TSevenDayBuyInfo{}
				seven.ActivityID = v.ID
				self.SevenDayLimit = append(self.SevenDayLimit, seven)
			}

			var activity TActivityData
			activity.ActivityID = v.ID
			activity.activityType = v.ActType
			activity.Status = v.Status
			self.ActivityLst = append(self.ActivityLst, activity)
		}

		mongodb.InsertToDB("GlobalVariables", self)
	}

	self.CheckActivityNew()

	self.UpdateActivity()
}

func (self *TGlobalVariables) CheckActivityNew() {
	for _, v := range gamedata.GT_ActivityLst {
		if v.ID == 0 {
			gamelog.Error("CheckActivityAdd Error: Invalid ActivityID:%d", v.ID)
			continue
		}

		isExist := false
		for _, n := range G_GlobalVariables.ActivityLst {
			if n.ActivityID == v.ID {
				isExist = true
				break
			}
		}

		if isExist == true {
			continue
		}

		if v.ActType == gamedata.Activity_Seven {
			seven := TSevenDayBuyInfo{}
			seven.ActivityID = v.ID
			G_GlobalVariables.SevenDayLimit = append(G_GlobalVariables.SevenDayLimit, seven)
			G_GlobalVariables.DB_AddSevenDayBuyInfo(seven)
		}

		var activity TActivityData
		activity.ActivityID = v.ID
		activity.Status = v.Status
		G_GlobalVariables.ActivityLst = append(G_GlobalVariables.ActivityLst, activity)
		G_GlobalVariables.DB_AddNewActivity(activity)
	}
}

func (self *TGlobalVariables) UpdateActivity() bool {
	openday := GetOpenServerDay()
	for i := 0; i < len(self.ActivityLst); i++ {
		pActInfo := gamedata.GetActivityInfo(self.ActivityLst[i].ActivityID)
		if pActInfo == nil {
			gamelog.Error("UpdateActivity Error : Invalid activityID:%d", self.ActivityLst[i].ActivityID)
			return false
		}

		self.ActivityLst[i].activityType = pActInfo.ActType
		self.ActivityLst[i].award = pActInfo.AwardType
		self.ActivityLst[i].beginTime, self.ActivityLst[i].actEndTime, self.ActivityLst[i].endTime = CalcActivityTime(self.ActivityLst[i].ActivityID, openday)
	}

	return true
}

//! 获取活动奖励
func (self *TGlobalVariables) GetActivityData(activityID int32) *TActivityData {
	//! 根据活动模板获取对应ID
	for i := 0; i < len(self.ActivityLst); i++ {
		if self.ActivityLst[i].ActivityID == activityID {
			return &self.ActivityLst[i]
		}
	}

	gamelog.Error("GetActivityData Error : Invalid activityID:%d", activityID)
	return nil
}

//! 获取活动奖励
func (self *TGlobalVariables) GetActivityAwardType(activityID int32) int {
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
func (self *TGlobalVariables) IsActivityOpen(activityID int32) bool {
	nowTime := utility.GetCurTime()
	openday := GetOpenServerDay()

	var pActData *TActivityData = nil
	for i := 0; i < len(self.ActivityLst); i++ {
		if activityID == self.ActivityLst[i].ActivityID {
			pActData = &self.ActivityLst[i]
			break
		}
	}

	if pActData == nil {
		gamelog.Error3("IsActivityOpen Error : Invalid activityID:%d", activityID)
		return false
	}

	if pActData.Status != 1 {
		return false
	}

	pActivityInfo := gamedata.GetActivityInfo(activityID)
	if pActivityInfo == nil {
		gamelog.Error("IsActivityOpen Error: pActivityInfo:nil")
		return false
	}

	if pActivityInfo.TimeType == gamedata.Time_NewSvr {
		if openday > 30 {
			return false
		}

	} else if pActivityInfo.TimeType == gamedata.Time_PublicSvr {
		if openday <= 30 {
			return false
		}

	} else if pActivityInfo.TimeType != gamedata.Time_AllSvr {
		gamelog.Error("IsActivityOpen Error: Invalid TimeType:%d", pActivityInfo.TimeType)
		return false
	}

	if pActivityInfo.CycleType == gamedata.CyCle_All {
		return true
	}

	if pActData.beginTime <= nowTime && nowTime < pActData.endTime {
		return true
	}

	return false
}

//! 判断当前是否为活动时间
//! 返回: 是否在活动操作期(有的专门设置了领奖期)    结束倒计时
func (self *TGlobalVariables) IsActivityTime(activityID int32) (bOk bool) {
	bOk = false
	nowTime := utility.GetCurTime()
	openday := GetOpenServerDay()
	var pActData *TActivityData = nil
	for i := 0; i < len(self.ActivityLst); i++ {
		if activityID == self.ActivityLst[i].ActivityID {
			pActData = &self.ActivityLst[i]
			break
		}
	}

	if pActData == nil {
		gamelog.Error("IsActivityTime Error : Invalid activityID:%d", activityID)
		return
	}

	if pActData.Status != 1 {
		return
	}

	pActivityInfo := gamedata.GetActivityInfo(activityID)
	if pActivityInfo == nil {
		gamelog.Error("IsActivityTime Error: pActivityInfo:nil")
		return
	}

	if pActivityInfo.TimeType == gamedata.Time_NewSvr {
		if openday > 30 {
			return
		}

	} else if pActivityInfo.TimeType == gamedata.Time_PublicSvr {
		if openday <= 30 {
			return
		}

	} else if pActivityInfo.TimeType != gamedata.Time_AllSvr {
		gamelog.Error("IsActivityTime Error: Invalid TimeType:%d", pActivityInfo.TimeType)
		return
	}

	if pActivityInfo.CycleType != gamedata.CyCle_All {
		bOk = true
		return
	}

	if pActData.beginTime > nowTime || nowTime > pActData.actEndTime {
		return
	}

	bOk = true
	return
}

//! 计算开启时间与关闭时间
func CalcActivityTime(activityID int32, openDay int) (beginTime int32, actEndTime int32, endTime int32) {
	pActInfo := gamedata.GetActivityInfo(activityID)
	if pActInfo == nil {
		gamelog.Error("CalcActivityTime Error : Invalid activityid:%d", activityID)
		return
	}

	if pActInfo.CycleType == gamedata.CyCle_All {
		return
	}

	nowTime := time.Now()
	if pActInfo.CycleType == gamedata.Cycle_Month { //! 按照月计算
		beginDate := time.Date(nowTime.Year(), nowTime.Month(), pActInfo.BeginTime, 0, 0, 0, 0, nowTime.Location())
		endDate := time.Date(nowTime.Year(), nowTime.Month(), pActInfo.EndTime, 23, 59, 59, 59, nowTime.Location())
		if endDate.Unix() <= nowTime.Unix() {
			beginDate = beginDate.AddDate(0, 1, 0)
			endDate = endDate.AddDate(0, 1, 0)
		}

		beginTime = int32(beginDate.Unix())
		actEndTime = int32(endDate.Unix()) - int32(pActInfo.AwardTime*86400)
		endTime = int32(endDate.Unix())
	} else if pActInfo.CycleType == gamedata.Cycle_Week { //! 按照周计算
		weekDay := int(nowTime.Weekday())
		if weekDay == 0 { //! 特殊处理周末
			weekDay = 7
		}

		beginDate := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location())
		beginDate = beginDate.AddDate(0, 0, pActInfo.BeginTime-weekDay)
		endDate := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 59, nowTime.Location())
		endDate = endDate.AddDate(0, 0, pActInfo.EndTime-weekDay)

		if endDate.Unix() <= nowTime.Unix() {
			beginDate = beginDate.AddDate(0, 0, 7)
			endDate = endDate.AddDate(0, 0, 7)
		}

		endTime = int32(endDate.Unix())
		actEndTime = int32(endDate.Unix()) - int32(pActInfo.AwardTime*86400)
		beginTime = int32(beginDate.Unix())
	} else if pActInfo.CycleType == gamedata.Cycle_OpenDay {
		beginDate := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location())
		beginDate = beginDate.AddDate(0, 0, -1*openDay)
		beginDate = beginDate.AddDate(0, 0, pActInfo.BeginTime)
		endDate := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 23, 59, 59, 59, nowTime.Location())
		endDate = endDate.AddDate(0, 0, pActInfo.EndTime-openDay)

		if endDate.Unix() <= nowTime.Unix() {
			beginTime = 0xFFFFFFF
			actEndTime = 0xFFFFFFF
			endTime = 0xFFFFFFF
		} else {
			beginTime = int32(beginDate.Unix())
			actEndTime = int32(endDate.Unix()) - int32(pActInfo.AwardTime*86400)
			endTime = int32(endDate.Unix())
		}

	} else if pActInfo.CycleType == gamedata.Cycle_FixDay {
		day := pActInfo.BeginTime % 100
		month := (pActInfo.BeginTime - day) / 100
		if day < 1 || day > 31 || month < 1 || month > 12 {
			gamelog.Error("CalcActivityTime Error : Invalid Activity BeginTime: %d", pActInfo.BeginTime)
			return
		}

		beginDate := time.Date(nowTime.Year(), time.Month(month), day, 0, 0, 0, 0, nowTime.Location())

		day = pActInfo.EndTime % 100
		month = (pActInfo.EndTime - day) / 100
		if day < 1 || day > 31 || month < 1 || month > 12 {
			gamelog.Error("CalcActivityTime Error :  Invalid Activity EndTime: %d", pActInfo.EndTime)
			return
		}

		endDate := time.Date(nowTime.Year(), time.Month(month), day, 23, 59, 59, 59, nowTime.Location())

		if endDate.Unix() <= nowTime.Unix() {
			beginTime = 0xFFFFFFF
			actEndTime = 0xFFFFFFF
			endTime = 0xFFFFFFF
		} else {
			beginTime = int32(beginDate.Unix())
			actEndTime = int32(endDate.Unix()) - int32(pActInfo.AwardTime*86400)
			endTime = int32(endDate.Unix())
		}
	}

	return
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
	self.DB_AddNewGroupPurchaseRecord(&newRecord)

	return &self.GroupPurchaseLst[length], length
}

func (self *TGlobalVariables) GetSevenDayLimit(activityID int32) *TSevenDayBuyInfo {
	for i := 0; i < len(G_GlobalVariables.SevenDayLimit); i++ {
		if G_GlobalVariables.SevenDayLimit[i].ActivityID == activityID {
			return &G_GlobalVariables.SevenDayLimit[i]
		}
	}

	return nil
}

func (self *TGlobalVariables) AddSevenDayLimit(activityID int32, index int) {
	for i := 0; i < len(G_GlobalVariables.SevenDayLimit); i++ {
		if G_GlobalVariables.SevenDayLimit[i].ActivityID == activityID {
			G_GlobalVariables.SevenDayLimit[i].LimitBuy[index] += 1
			self.DB_SaveSevenDayLimit(index)
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
	self.DB_UpdateGroupPurchaseSaleNum(index)

	return recordInfo.SaleNum
}

//! 添加团购记录
func (self *TGlobalVariables) DB_AddNewGroupPurchaseRecord(record *TGroupPurchaseInfo) {
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$push": bson.M{"grouppurchaselst": *record}})
}

//! 更新团购商品购买记录
func (self *TGlobalVariables) DB_UpdateGroupPurchaseSaleNum(index int) {
	filedName := fmt.Sprintf("grouppurchaselst.%d.salenum", index)
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
		filedName: self.GroupPurchaseLst[index].SaleNum}})
}

func (self *TGlobalVariables) DB_CleanMoneyPoor() {
	if self.NormalMoneyPoor != 0 || self.ExcitedMoneyPoor != 0 {
		self.NormalMoneyPoor = 0
		self.ExcitedMoneyPoor = 0
		mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
			"normalmoneypoor":  self.NormalMoneyPoor,
			"excitedmoneypoor": self.ExcitedMoneyPoor}})
	}
}

func (self *TGlobalVariables) DB_CleanGroupPurchase() {
	if len(self.GroupPurchaseLst) != 0 {
		self.GroupPurchaseLst = []TGroupPurchaseInfo{}

		mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
			"grouppurchaselst": self.GroupPurchaseLst}})
	}
}

func (self *TGlobalVariables) DB_SaveMoneyPoor() {
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
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
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": self})
}

func (self *TGlobalVariables) DB_AddNewActivity(activity TActivityData) {
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$push": bson.M{"activitylst": activity}})
}

func (self *TGlobalVariables) DB_UpdateActivityLst() {
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
		"activitylst": self.ActivityLst}})
}

func (self *TGlobalVariables) DB_AddSevenDayBuyInfo(seven TSevenDayBuyInfo) {
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$push": bson.M{"sevendaylimit": seven}})
}

func (self *TGlobalVariables) DB_CleanSevenDayInfo(activityID int32) {
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
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{filedName: self.SevenDayLimit[index].LimitBuy}})
}

func (self *TGlobalVariables) DB_SaveSevenDayLimit(index int) {
	filedName := fmt.Sprintf("sevendaylimit.%d", index)
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{filedName: self.SevenDayLimit[index]}})
}

func (self *TGlobalVariables) DB_UpdateActivityStatus(index int) {
	filedName := fmt.Sprintf("activitylst.%d.status", index)
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
		filedName: self.ActivityLst[index].Status}})
}

// GM调用的增删全服奖励接口
func (self *TGlobalVariables) AddSvrAward(pAwardData *TAwardData) {
	self.SvrAwardIncID += 1
	pAwardData.ID = self.SvrAwardIncID
	pAwardData.Time = utility.GetCurTime()
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
	mongodb.UpdateToDB("SvrAwardCenter", &bson.M{"_id": 0}, &bson.M{"$push": bson.M{"svrawardlist": *pAwardData}})
}
func (self *TGlobalVariables) DB_DelAward(id int) {
	mongodb.UpdateToDB("SvrAwardCenter", &bson.M{"_id": 0}, &bson.M{"$pull": bson.M{"svrawardlist": bson.M{"id": id}}})
}
func (self *TGlobalVariables) DB_SaveIncrementID() {
	mongodb.UpdateToDB("SvrAwardCenter", &bson.M{"_id": 0}, &bson.M{"$set": bson.M{"svrawardincid": self.SvrAwardIncID}})
}

func (self *TGlobalVariables) DB_UpdateLimitSaleNum() {
	mongodb.UpdateToDB("GlobalVariables", &bson.M{"_id": 1}, &bson.M{"$set": bson.M{
		"limitsalenum": self.LimitSaleNum}})
}

package mainlogic

import (
	"appconfig"
	"fmt"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"utility"
)

type TLimitSaleInfo struct {
	ID     int  //! 商品ID
	Status bool //! 是否购买 true -> 已购买  false -> 未购买
}

//! 限时优惠活动
type TActivityLimitSale struct {
	ActivityID int //! 活动ID

	Score       int              //! 当前积分
	ItemLst     []TLimitSaleInfo //! 当天优惠物品
	RefreshMark bool             //! 刷新标记
	AwardMark   Mark             //! 全民奖励领取标记
	WeekReset   int              //! 全民奖励刷新周

	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityLimitSale) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLimitSale) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr

	self.Score = 0
	self.ItemLst = []TLimitSaleInfo{}
	self.AwardMark = 0
	self.WeekReset = utility.GetCurDay()

	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityLimitSale) Refresh(versionCode int) {
	//! 刷新贩售物品
	self.RefreshItem()
	self.RefreshMark = true
	self.VersionCode = versionCode

	//! 如果积分满100分, 则清空
	if self.Score >= 100 {
		self.Score = 0
	}

	if utility.IsSameWeek(self.WeekReset) == false {
		//! 刷新全民奖励
		self.AwardMark = 0
		self.WeekReset = utility.GetCurDay()
		G_GlobalVariables.LimitSaleNum = 0
		go G_GlobalVariables.DB_UpdateLimitSaleNum()
	}

	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityLimitSale) End(versionCode int, resetCode int) {

	self.ResetCode = resetCode
	self.VersionCode = versionCode

	go self.DB_Reset()
}

func (self *TActivityLimitSale) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityLimitSale) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityLimitSale) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.Score >= 100 {
		//! 积分满时显示红点
		self.RefreshMark = false
		go self.DB_SaveRefreshMark()
		return true
	}

	if self.RefreshMark == true {
		self.RefreshMark = false
		go self.DB_SaveRefreshMark()
		return true
	}

	return false
}

func (self *TActivityLimitSale) RefreshItem() {
	if len(self.ItemLst) != 0 {
		self.ItemLst = []TLimitSaleInfo{}
	}

	itemIDLst := gamedata.RandLimitSaleItem()
	for i := 0; i < len(itemIDLst); i++ {
		var item TLimitSaleInfo
		item.ID = itemIDLst[i]
		item.Status = false
		self.ItemLst = append(self.ItemLst, item)
	}
}

func (self *TActivityLimitSale) DB_UpdateScore() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"limitsale.score": self.Score}})
}

func (self *TActivityLimitSale) DB_UpdateStatus(index int) {
	filedName := fmt.Sprintf("limitsale.itemlst.%d.status", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.ItemLst[index].Status}})
}

func (self *TActivityLimitSale) DB_SaveRefreshMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"limitsale.refreshmark": self.RefreshMark}})
}

func (self *TActivityLimitSale) DB_UpdateAwardMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"limitsale.awardmark": self.AwardMark}})
}

func (self *TActivityLimitSale) DB_Refresh() bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"limitsale.versioncode": self.VersionCode,
		"limitsale.refreshmark": self.RefreshMark,
		"limitsale.score":       self.Score,
		"limitsale.awardmark":   self.AwardMark,
		"limitsale.weekreset":   self.WeekReset,
		"limitsale.itemlst":     self.ItemLst}})
}

func (self *TActivityLimitSale) DB_Reset() bool {
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"limitsale.versioncode": self.VersionCode,
		"limitsale.resetcode":   self.ResetCode,
		"limitsale.refreshmark": self.RefreshMark,
		"limitsale.score":       self.Score,
		"limitsale.awardmark":   self.AwardMark,
		"limitsale.weekreset":   self.WeekReset,
		"limitsale.itemlst":     self.ItemLst}})
}

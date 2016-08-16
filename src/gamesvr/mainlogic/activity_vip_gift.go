package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"utility"
)

//! VIP每周礼包
type TVipWeekItem struct {
	ID       int
	Award    int
	BuyTimes int
}

//! VIP礼包活动
type TActivityVipGift struct {
	ActivityID int //! 活动ID

	IsRecvWelfare bool           //! 是否领取日常礼包
	WeekGift      []TVipWeekItem //! VIP每周礼包
	ResetWeek     int            //! VIP每周福利刷新

	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityVipGift) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityVipGift) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.VersionCode = vercode
	self.ResetCode = resetcode
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	//! 设置每日礼包下次刷新时间
	self.ResetWeek = utility.GetCurDay()
	self.IsRecvWelfare = false

	//! 刷新VIP每周礼包
	self.RefreshWeekGift(false)
}

//! 刷新数据
func (self *TActivityVipGift) Refresh(versionCode int) {
	gamelog.Info("TActivityVipGift Refresh")
	self.CheckWeekGiftRefresh()
	self.VersionCode = versionCode
	self.IsRecvWelfare = false

	//! 更新至数据库
	self.DB_Refresh()

}

//! 活动结束
func (self *TActivityVipGift) End(versionCode int, resetCode int) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.IsRecvWelfare = false
	self.WeekGift = []TVipWeekItem{}
	self.ResetWeek = 0

	go self.DB_Reset()
}

func (self *TActivityVipGift) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityVipGift) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityVipGift) RedTip() bool {
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.IsRecvWelfare == true {
		return false
	}

	return true
}

func (self *TActivityVipGift) RefreshWeekGift(isSave bool) {
	if len(self.WeekGift) > 0 {
		self.WeekGift = []TVipWeekItem{}
	}

	itemLst := gamedata.GetVipWeekItem(self.activityModule.ownplayer.GetLevel())

	for _, v := range itemLst {
		var item TVipWeekItem
		item.ID = v.ID
		item.Award = v.Award
		item.BuyTimes = 0
		self.WeekGift = append(self.WeekGift, item)
	}

	if isSave == true {
		go self.DB_UpdateWeekGiftToDatabase()
	}
}

func (self *TActivityVipGift) CheckWeekGiftRefresh() {
	if utility.IsSameWeek(self.ResetWeek) {
		return
	}

	//! 重置刷新时间
	self.ResetWeek = utility.GetCurDay()

	//! 刷新每周礼包信息
	self.RefreshWeekGift(true)
}

//! 重置活动
func (self *TActivityVipGift) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"vipgift.activityid":    self.ActivityID,
		"vipgift.weekgift":      self.WeekGift,
		"vipgift.versioncode":   self.VersionCode,
		"vipgift.resetcode":     self.ResetCode,
		"vipgift.isrecvwelfare": self.IsRecvWelfare,
		"vipgift.resetweek":     self.ResetWeek}})
	return true
}

//! 更新每周礼包信息
func (self *TActivityVipGift) DB_UpdateWeekGiftToDatabase() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"vipgift.weekgift":  self.WeekGift,
		"vipgift.resetweek": self.ResetWeek}})
}

//! 更新VIP日常福利领取时间到数据库
func (self *TActivityVipGift) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"vipgift.isrecvwelfare": self.IsRecvWelfare,
		"vipgift.versioncode":   self.VersionCode}})

	return true
}

func (self *TActivityVipGift) DB_SaveDailyResetTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"vipgift.isrecvwelfare": self.IsRecvWelfare}})
}

//! 更新购买次数
func (self *TActivityVipGift) DB_UpdateBuyTimes(id int, times int) {
	index := -1
	for i, v := range self.WeekGift {
		if v.ID == id {
			index = i
		}
	}

	if index < 0 {
		gamelog.Error("DB_UpdateBuyTimes Fail: Not find week gift id: %d", id)
		return
	}

	filedName := fmt.Sprintf("vipgift.weekgift.%d.buytimes", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: times}})
}

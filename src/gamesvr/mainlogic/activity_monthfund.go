package mainlogic

import (
	"appconfig"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 月基金
type TActivityMonthFund struct {
	ActivityID int //! 活动ID

	Day       int  //! 基金领取天数
	AwardMark Mark //! 基金领取标记

	VersionCode int32 //! 版本号
	ResetCode   int32 //! 迭代号

	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityMonthFund) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityMonthFund) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr

	self.Day = 0

	self.activityModule.activityPtrs[activityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityMonthFund) Refresh(versionCode int32) {
	if self.Day != 0 {
		self.Day -= int(versionCode - self.VersionCode)
	}

	self.VersionCode = versionCode
	go self.DB_Refresh()
}

func (self *TActivityMonthFund) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	//! 补发未领取奖励
	self.AwardRetroactive()

	self.AwardMark = 0
	self.Day = 0

	go self.DB_Reset()
}

func (self *TActivityMonthFund) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityMonthFund) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityMonthFund) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	//! 未购买月基金
	if self.Day == 0 {
		return false
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	awardCount := gamedata.GetMonthFundAwardCount(awardType)

	if self.AwardMark.Get(uint(awardCount-self.Day+1)) == false { //! 今日奖励未领取
		return true
	}

	return false
}

func (self *TActivityMonthFund) SetMonthFund(rmb int) {
	if self.Day != 0 {
		//! 已经激活月基金
		return
	}

	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return //! 活动尚未开启
	}

	if len(self.activityModule.MonthCard.CardDays) < 3 {
		return
	}

	if self.activityModule.MonthCard.CardDays[1] == 0 ||
		self.activityModule.MonthCard.CardDays[2] == 0 {
		//! 必须激活双月卡
		return
	}

	chargeInfo := gamedata.GetChargeItem(gamedata.MonthFundCostMoneyNum)
	if rmb == chargeInfo.RenMinBi {
		awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
		awardCount := gamedata.GetMonthFundAwardCount(awardType)
		self.Day += awardCount
		self.AwardMark = 0
		go self.DB_MonthFund()
	}
}

func (self *TActivityMonthFund) AwardRetroactive() {
	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	awardCount := gamedata.GetMonthFundAwardCount(awardType)

	awardLst := []gamedata.ST_ItemData{}
	for i := 1; i <= awardCount; i++ {
		if self.AwardMark.Get(uint(i)) == false {
			//! 加入补发奖励名单
			fundInfo := gamedata.GetMonthFundAward(awardType, i)
			awardLst = append(awardLst, gamedata.ST_ItemData{fundInfo.ItemID, fundInfo.ItemNum})
		}
	}

	//! 发送补偿邮件
	SendAwardMail(self.activityModule.PlayerID, Text_MonthFund, awardLst, []string{})
}

func (self *TActivityMonthFund) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthfund.activityid":  self.ActivityID,
		"monthfund.day":         self.Day,
		"monthfund.awardmark":   self.AwardMark,
		"monthfund.resetcode":   self.ResetCode,
		"monthfund.versioncode": self.VersionCode}})
}

func (self *TActivityMonthFund) DB_Refresh() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthfund.day":         self.Day,
		"monthfund.versioncode": self.VersionCode}})
}

func (self *TActivityMonthFund) DB_MonthFund() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthfund.day":       self.Day,
		"monthfund.awardmark": self.AwardMark}})
}

func (self *TActivityMonthFund) DB_UpdateAwardMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"monthfund.awardmark": self.AwardMark}})
}

package mainlogic

import (
	"appconfig"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 周周盈
type TActivityWeekAward struct {
	ActivityID int //! 活动ID

	LoginDay    int  //! 登录天数
	RechargeNum int  //! 充值数目
	AwardMark   Mark //! 奖励标记 位运算

	VersionCode int32 //! 版本号
	ResetCode   int32 //! 迭代号

	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityWeekAward) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityWeekAward) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr

	self.LoginDay = 1
	self.RechargeNum = 0
	self.AwardMark = 0

	self.activityModule.activityPtrs[activityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityWeekAward) Refresh(versionCode int32) {
	//! 刷新签到标记
	self.LoginDay += int(versionCode - self.VersionCode)
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 结束
func (self *TActivityWeekAward) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.LoginDay = 0
	self.RechargeNum = 0
	self.AwardMark = 0
	go self.DB_Reset()
}

func (self *TActivityWeekAward) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityWeekAward) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityWeekAward) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	awardLst := gamedata.GetWeekAwardInfoLst(awardType)
	if self.LoginDay == 1 && self.AwardMark.Get(uint(1)) == false {
		return true //! 免费奖励未领取
	}

	for i := 1; i <= 7; i++ {
		awardInfo := awardLst[i]
		if self.AwardMark.Get(uint(i+1)) == false &&
			self.LoginDay >= i &&
			self.RechargeNum >= awardInfo.RechargeNum {
			return true
		}
	}

	return false
}

func (self *TActivityWeekAward) AddRechargeNum(rechargeNum int) {
	isEnd, _ := G_GlobalVariables.IsActivityTime(self.ActivityID)
	if isEnd == false {
		return
	}

	self.RechargeNum += rechargeNum
	go self.DB_UpdateRechargeNum()
}

func (self *TActivityWeekAward) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"weekaward.loginday":    self.LoginDay,
		"weekaward.versioncode": self.VersionCode}})
	return true
}

func (self *TActivityWeekAward) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"weekaward.activityid":  self.ActivityID,
		"weekaward.loginday":    self.LoginDay,
		"weekaward.AwardMark":   self.AwardMark,
		"weekaward.rechargenum": self.RechargeNum,
		"weekaward.versioncode": self.VersionCode,
		"weekaward.resetcode":   self.ResetCode}})
	return true
}

func (self *TActivityWeekAward) DB_UpdateRechargeNum() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"weekaward.rechargenum": self.RechargeNum}})
}

func (self *TActivityWeekAward) DB_UpdateAwardMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"weekaward.awardmark": self.AwardMark}})
}

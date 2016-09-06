package mainlogic

import (
	"gamesvr/gamedata"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TActivityMoneyGod struct {
	ActivityID      int              //! 活动ID
	CurrentTimes    int              //! 当前剩余领取次数
	CumulativeTimes int              //! 当前累积次数
	TotalMoney      int              //! 累积银币
	NextTime        int64            //! 下次迎财神时间
	VersionCode     int32            //! 版本号
	ResetCode       int32            //! 迭代号
	activityModule  *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityMoneyGod) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityMoneyGod) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.CurrentTimes = 3
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityMoneyGod) Refresh(versionCode int32) {
	if self.CurrentTimes != 0 {
		//! 迎财神中断,奖金池清空
		self.TotalMoney = 0
		self.CumulativeTimes = 0
	}

	//! 迎财神次数重置
	self.CurrentTimes = 3
	self.NextTime = 0
	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityMoneyGod) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.DB_Reset()
}

func (self *TActivityMoneyGod) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityMoneyGod) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityMoneyGod) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	self.CheckMoneyGod()
	now := time.Now().Unix()
	//! 可迎财神
	if now >= self.NextTime && self.CurrentTimes > 0 {
		return true
	}

	//! 可领取累积银币
	activityInfo := gamedata.GetActivityInfo(self.ActivityID)
	moneyInfo := gamedata.GetMoneyGoldInfo(activityInfo.AwardType)
	if self.CumulativeTimes >= moneyInfo.AwardTimes {
		return true
	}

	return false
}

//! 迎财神时间检测
func (self *TActivityMoneyGod) CheckMoneyGod() {
	now := time.Now().Unix()
	if now < self.NextTime || self.CurrentTimes == 0 || self.NextTime == 0 {
		return
	}

	self.NextTime = 0
	self.DB_UpdateNextTime()
}

func (self *TActivityMoneyGod) DB_UpdateNextTime() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"moneygod.nexttime": self.NextTime}})
}

func (self *TActivityMoneyGod) DB_Refresh() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"moneygod.currenttimes":    self.CurrentTimes,
		"moneygod.nexttime":        self.NextTime,
		"moneygod.totalmoney":      self.TotalMoney,
		"moneygod.cumulativetimes": self.CumulativeTimes,
		"moneygod.versioncode":     self.VersionCode}})
}

func (self *TActivityMoneyGod) DB_UpdateCumulativeTimes() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"moneygod.totalmoney":      self.TotalMoney,
		"moneygod.cumulativetimes": self.CumulativeTimes}})
}

func (self *TActivityMoneyGod) DB_Reset() {
	GameSvrUpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"moneygod.activityid":  self.ActivityID,
		"moneygod.versioncode": self.VersionCode,
		"moneygod.resetcode":   self.ResetCode}})
}

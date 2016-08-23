package mainlogic

import (
	"appconfig"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

type TActivityFirstRecharge struct {
	ActivityID         int              //! 活动ID
	FirstRechargeAward int              //! 首充标记 0->不能领取 1->可以领取 2->已领取并开启次充奖励 3->次充奖励可领取 4->已领取次充奖励
	VersionCode        int32            //! 版本号
	ResetCode          int32            //! 迭代号
	activityModule     *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityFirstRecharge) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityFirstRecharge) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.FirstRechargeAward = 0
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityFirstRecharge) Refresh(versionCode int32) {
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityFirstRecharge) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	go self.DB_Reset()
}

func (self *TActivityFirstRecharge) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityFirstRecharge) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityFirstRecharge) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.FirstRechargeAward == 1 || self.FirstRechargeAward == 3 {
		return true
	}
	return false
}

//!	首充/次充检测
func (self *TActivityFirstRecharge) CheckRecharge(rmb int) {
	if self.FirstRechargeAward == 0 { //! 首充
		self.FirstRechargeAward = 1
	} else if self.FirstRechargeAward == 2 && rmb == gamedata.NextAwardNeedRecharge { //! 次充
		self.FirstRechargeAward = 3
	}

	go self.DB_SetFirstRechargeMark()

}

//! 更新
func (self *TActivityFirstRecharge) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"firstrecharge.versioncode": self.VersionCode}})
	return true
}

//! 重置
func (self *TActivityFirstRecharge) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"firstrecharge.activityid":  self.ActivityID,
		"firstrecharge.resetcode":   self.ResetCode,
		"firstrecharge.versioncode": self.VersionCode}})
	return true
}

func (self *TActivityFirstRecharge) DB_SetFirstRechargeMark() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"firstrecharge.firstrechargeaward": self.FirstRechargeAward}})
}

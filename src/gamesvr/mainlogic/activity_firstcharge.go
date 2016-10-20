package mainlogic

import (
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

type TActivityFirstCharge struct {
	ActivityID  int32            //! 活动ID
	FirstAward  int              //! 首充标记 0->不能领取 1->可以领取 2->已领取
	NextAward   int              //! 次充标记 0->不能领取 1->可以领取 2->已领取
	VersionCode int32            //! 版本号
	ResetCode   int32            //! 迭代号
	modulePtr   *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityFirstCharge) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityFirstCharge) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.FirstAward = 0
	self.NextAward = 0
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityFirstCharge) Refresh(versionCode int32) {
	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityFirstCharge) End(versionCode int32, resetCode int32) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.DB_Reset()
}

func (self *TActivityFirstCharge) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityFirstCharge) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityFirstCharge) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.FirstAward == 1 || self.NextAward == 1 {
		return true
	}
	return false
}

//!	首充/次充检测
func (self *TActivityFirstCharge) CheckRecharge(rmb int) {
	if self.FirstAward == 0 { //! 首充
		self.FirstAward = 1
	} else if rmb >= gamedata.NextAwardNeedRecharge && self.NextAward == 0 { //! 次充
		self.NextAward = 1
	}

	self.DB_SetFirstRechargeMark()

}

//! 更新
func (self *TActivityFirstCharge) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"firstcharge.versioncode": self.VersionCode}})
}

//! 重置
func (self *TActivityFirstCharge) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"firstcharge.activityid":  self.ActivityID,
		"firstcharge.resetcode":   self.ResetCode,
		"firstcharge.versioncode": self.VersionCode}})
}

func (self *TActivityFirstCharge) DB_SetFirstRechargeMark() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"firstcharge.firstaward": self.FirstAward,
		"firstcharge.nextaward":  self.NextAward}})
}

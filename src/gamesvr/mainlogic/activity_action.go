package mainlogic

import (
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"time"
)

//! 领取体力
type TActivityAction struct {
	RecvAction  BitsType         //! 领取体力标记
	ActivityID  int32            //! 活动ID
	VersionCode int32            //! 版本号
	ResetCode   int32            //! 迭代号
	modulePtr   *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityAction) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityAction) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
	self.RecvAction = 0
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityAction) Refresh(versionCode int32) {
	//! 重置体力领取标记
	self.RecvAction = 0
	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityAction) End(versionCode int32, resetCode int32) {
	self.RecvAction = 0
	self.ResetCode = resetCode
	self.VersionCode = versionCode
	self.DB_Reset()
}

func (self *TActivityAction) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityAction) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityAction) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	pActivityAction := gamedata.GetActivityAction(awardType)
	if pActivityAction == nil {
		return false
	}

	now := time.Now()
	sec := now.Hour()*3600 + now.Minute()*60 + now.Second()

	index := -1
	for i := 0; i < 4; i++ {
		if sec >= pActivityAction.Time_Begin[i] && sec <= pActivityAction.Time_End[i] {
			index = i
			break
		}
	}

	if index < 0 {
		return false
	}

	//! 判断当前时间段是否已领取
	if self.RecvAction.Get(index+1) == true {
		return false
	}

	return true
}

func (self *TActivityAction) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"receiveaction.activityid":  self.ActivityID,
		"receiveaction.versioncode": self.VersionCode,
		"receiveaction.recvaction":  self.RecvAction,
		"receiveaction.resetcode":   self.ResetCode}})
}

func (self *TActivityAction) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		"receiveaction.recvaction":  self.RecvAction,
		"receiveaction.versioncode": self.VersionCode}})
}

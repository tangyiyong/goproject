package mainlogic

import (
	"appconfig"
	"gamesvr/gamedata"
	"mongodb"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//! 领取体力
type TActivityReceiveAction struct {
	ActivityID int  //! 活动ID
	RecvAction Mark //! 领取体力标记

	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityReceiveAction) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityReceiveAction) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.RecvAction = 0
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityReceiveAction) Refresh(versionCode int) {
	//! 重置体力领取标记
	self.RecvAction = 0
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityReceiveAction) End(versionCode int, resetCode int) {
	self.RecvAction = 0
	self.ResetCode = resetCode
	self.VersionCode = versionCode
	go self.DB_Reset()
}

func (self *TActivityReceiveAction) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityReceiveAction) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityReceiveAction) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	activityInfo := gamedata.GetActivityInfo(self.ActivityID)
	index := gamedata.IsRecvActionTime(activityInfo.AwardType)
	if index < 0 {
		return false
	}

	//! 判断当前时间段是否已领取
	if self.RecvAction.Get(uint(index)) == true {
		return false
	}

	return true
}

func (self *TActivityReceiveAction) GetNextActionAwardTime() int {
	now := time.Now()
	nowSec := now.Hour()*3600 + now.Minute()*60 + now.Second()
	var nextTime int
	var endTime int

	isExist := false
	for _, v := range gamedata.GT_RecvActionLst {
		if nowSec < v.Time_Begin {
			nextTime = v.Time_Begin
			isExist = true
			break
		}
		endTime = v.Time_End
	}

	if isExist == false {
		return gamedata.GT_RecvActionLst[0].Time_Begin + (endTime - nowSec)
	}

	return (nextTime - nowSec)
}

func (self *TActivityReceiveAction) DB_Reset() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"receiveaction.activityid":  self.ActivityID,
		"receiveaction.versioncode": self.VersionCode,
		"receiveaction.recvaction":  self.RecvAction,
		"receiveaction.resetcode":   self.ResetCode}})
	return true
}

func (self *TActivityReceiveAction) DB_Refresh() bool {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		"receiveaction.recvaction":  self.RecvAction,
		"receiveaction.versioncode": self.VersionCode}})
	return true
}

package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 登录送礼活动
type TActivityLogin struct {
	ActivityID int //! 活动ID

	LoginDay   int  //! 登录天数
	LoginAward Mark //! 登录奖励

	VersionCode    int32            //! 版本号
	ResetCode      int32            //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityLogin) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLogin) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.LoginDay = 0
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityLogin) Refresh(versionCode int32) {
	//! 累计登陆
	self.VersionCode = versionCode

	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityLogin) End(versionCode int32, resetCode int32) {
	self.LoginDay = 0
	self.LoginAward = 0

	self.VersionCode = versionCode

	go self.DB_Reset()
}

func (self *TActivityLogin) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityLogin) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityLogin) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	return true
}

func (self *TActivityLogin) AddLoginDay(index int) {
	self.LoginDay++
	go self.DB_AddLoginDay(index)
}

func (self *TActivityLogin) DB_Reset() {
	index := -1
	for i, v := range self.activityModule.Login {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Login DB_Reset fail. self.ActivityID: %d", self.ActivityID)
		return
	}

	filedName1 := fmt.Sprintf("login.%d.loginday", index)
	filedName2 := fmt.Sprintf("login.%d.loginaward", index)
	filedName4 := fmt.Sprintf("login.%d.versioncode", index)
	filedName5 := fmt.Sprintf("login.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName1: self.LoginDay,
		filedName2: self.LoginAward,
		filedName4: self.VersionCode,
		filedName5: self.ResetCode}})
}

func (self *TActivityLogin) DB_AddLoginDay(index int) {

	filedName := fmt.Sprintf("login.%d.loginday", index)
	filedName3 := fmt.Sprintf("login.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.LoginDay,
		filedName3: self.VersionCode}})

}

func (self *TActivityLogin) DB_Refresh() {
	index := -1
	for i, v := range self.activityModule.Login {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Login DB_Refresh fail. self.ActivityID: %d", self.ActivityID)
		return
	}

	filedName := fmt.Sprintf("login.%d.loginday", index)
	filedName3 := fmt.Sprintf("login.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.LoginDay,
		filedName3: self.VersionCode}})
}

func (self *TActivityLogin) DB_UpdateLoginAward(activityIndex int) {
	filedName := fmt.Sprintf("login.%d.loginaward", activityIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.LoginAward}})
}

package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 登录送礼活动
type TActivityLogin struct {
	ActivityID int //! 活动ID

	LoginDay   int  //! 登录天数
	LoginAward Mark //! 登录奖励
	IsLogin    bool //! 是否登录

	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityLogin) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLogin) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.LoginDay = 0
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityLogin) Refresh(versionCode int) {
	//! 累计登陆
	self.IsLogin = false

	self.VersionCode = versionCode

	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityLogin) End(versionCode int, resetCode int) {
	self.LoginDay = 0
	self.IsLogin = false
	self.LoginAward = 0

	self.VersionCode = versionCode

	go self.DB_Reset()
}

func (self *TActivityLogin) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityLogin) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityLogin) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	if self.IsLogin == true {
		return false
	}

	return true
}

func (self *TActivityLogin) AddLoginDay(index int) {
	if self.IsLogin == true {
		return
	}
	self.LoginDay++
	self.IsLogin = true
	go self.DB_AddLoginDay(index)
}

func (self *TActivityLogin) DB_Reset() bool {
	index := -1
	for i, v := range self.activityModule.Login {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Login DB_Reset fail. self.ActivityID: %d", self.ActivityID)
		return false
	}

	filedName1 := fmt.Sprintf("login.%d.loginday", index)
	filedName2 := fmt.Sprintf("login.%d.loginaward", index)
	filedName3 := fmt.Sprintf("login.%d.islogin", index)
	filedName4 := fmt.Sprintf("login.%d.versioncode", index)
	filedName5 := fmt.Sprintf("login.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName1: self.LoginDay,
		filedName2: self.LoginAward,
		filedName3: self.IsLogin,
		filedName4: self.VersionCode,
		filedName5: self.ResetCode}})
	return true
}

func (self *TActivityLogin) DB_AddLoginDay(index int) {

	filedName := fmt.Sprintf("login.%d.loginday", index)
	filedName2 := fmt.Sprintf("login.%d.islogin", index)
	filedName3 := fmt.Sprintf("login.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.LoginDay,
		filedName2: self.IsLogin,
		filedName3: self.VersionCode}})

}

func (self *TActivityLogin) DB_Refresh() bool {

	index := -1
	for i, v := range self.activityModule.Login {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Login DB_Refresh fail. self.ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("login.%d.loginday", index)
	filedName2 := fmt.Sprintf("login.%d.islogin", index)
	filedName3 := fmt.Sprintf("login.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.LoginDay,
		filedName2: self.IsLogin,
		filedName3: self.VersionCode}})

	return true
}

func (self *TActivityLogin) DB_UpdateLoginAward(activityIndex int) {
	filedName := fmt.Sprintf("login.%d.loginaward", activityIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.LoginAward}})
}

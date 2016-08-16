package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 充值活动数据
type TActivityRecharge struct {
	ActivityID       int              //! 活动ID
	RechargeValue    int              //! 活动期间累积充值数额
	RechargeAwardLst Mark             //! 累积充值领取标记 (索引)
	VersionCode      int              //! 版本号
	ResetCode        int              //! 迭代号
	activityModule   *TActivityModule //! 指针
}

//! 赋值基础数据
func (self *TActivityRecharge) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityRecharge) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivityRecharge) Refresh(versionCode int) {
	//! 累积充值不会刷新
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityRecharge) End(versionCode int, resetCode int) {
	self.RechargeValue = 0
	self.RechargeAwardLst = 0
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	go self.DB_Reset()
}

func (self *TActivityRecharge) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivityRecharge) GetResetV() int {
	return self.ResetCode
}

func (self *TActivityRecharge) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	awardLst := gamedata.GetRechargeInfo(awardType)

	recvLst := IntLst{}
	for i, v := range awardLst {
		if self.RechargeValue > v.Recharge {
			recvLst.Add(i + 1)
		}
	}

	for _, v := range recvLst {
		if self.RechargeAwardLst.Get(uint(v)) == false {
			return true
		}
	}

	return false
}

func (self *TActivityRecharge) DB_Refresh() bool {
	index := -1
	for i, v := range self.activityModule.Recharge {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Recharge DB_Refresh fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("recharge.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.VersionCode}})
	return true
}

func (self *TActivityRecharge) DB_Reset() bool {
	index := -1
	for i, v := range self.activityModule.Recharge {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Recharge DB_Reset fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName1 := fmt.Sprintf("recharge.%d.rechargevalue", index)
	filedName2 := fmt.Sprintf("recharge.%d.rechargeawardlst", index)
	filedName3 := fmt.Sprintf("recharge.%d.versioncode", index)
	filedName4 := fmt.Sprintf("recharge.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName1: self.RechargeValue,
		filedName2: self.RechargeAwardLst,
		filedName3: self.VersionCode,
		filedName4: self.ResetCode}})
	return true
}

func (self *TActivityRecharge) DB_UpdateRechargeMark(index int, mark int) {
	filedName := fmt.Sprintf("recharge.%d.rechargeawardlst", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: mark}})
}

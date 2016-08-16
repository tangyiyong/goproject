package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

type TSingleRechargeRecord struct {
	Money  int //! 单笔充值人民币
	Status int
}

type SingleRechargeRecordLst []TSingleRechargeRecord

func (self SingleRechargeRecordLst) Len() int {
	return len(self)
}

func (self SingleRechargeRecordLst) Less(i int, j int) bool {
	if (self)[i].Money < (self)[j].Money {
		return true
	}
	return false
}

func (self SingleRechargeRecordLst) Swap(i int, j int) {
	temp := (self)[i]
	(self)[i] = (self)[j]
	(self)[j] = temp
}

type TActivityRechargeInfo struct {
	Index int //! 领取索引
	Times int //! 领取次数
}

//! 充值活动数据
type TActivitySingleRecharge struct {
	ActivityID     int                     //! 活动ID
	RechargeRecord SingleRechargeRecordLst //! 活动期间单笔充值记录
	SingleAwardLst []TActivityRechargeInfo //! 单充奖励领取记录
	VersionCode    int                     //! 版本号
	ResetCode      int                     //! 迭代号
	activityModule *TActivityModule        //! 指针
}

//! 赋值基础数据
func (self *TActivitySingleRecharge) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivitySingleRecharge) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode
}

//! 刷新数据
func (self *TActivitySingleRecharge) Refresh(versionCode int) {
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivitySingleRecharge) End(versionCode int, resetCode int) {
	self.RechargeRecord = []TSingleRechargeRecord{}
	self.SingleAwardLst = []TActivityRechargeInfo{}
	self.VersionCode = versionCode
	self.ResetCode = resetCode
	go self.DB_Reset()
}

func (self *TActivitySingleRecharge) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivitySingleRecharge) GetResetV() int {
	return self.ResetCode
}

func (self *TActivitySingleRecharge) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	awardType := G_GlobalVariables.GetActivityAwardType(self.ActivityID)
	singleRechargeTaskLst := gamedata.GetRechargeInfo(awardType)

	for i, v := range singleRechargeTaskLst {
		isHaveTimes := true
		for _, n := range self.SingleAwardLst {
			if n.Index == i+1 {
				if n.Times >= v.Times {
					isHaveTimes = false
					break
				}
			}
		}

		if isHaveTimes == false {
			continue
		}

		for _, n := range self.RechargeRecord {
			if n.Money >= v.Recharge && n.Status == 0 {
				return true
			}
		}
	}

	return false
}

func (self *TActivitySingleRecharge) GetSingleRechargeAwardTimes(index int) (*TActivityRechargeInfo, int) {

	for i, n := range self.SingleAwardLst {
		if n.Index == index {
			return &self.SingleAwardLst[i], i
		}
	}
	return nil, 0
}

func (self *TActivitySingleRecharge) DB_Refresh() bool {
	index := -1
	for i, v := range self.activityModule.SingleRecharge {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("SingleRecharge DB_Refresh fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("singlerecharge.%d.versioncode", index)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.VersionCode}})
	return true
}

func (self *TActivitySingleRecharge) DB_Reset() bool {
	index := -1
	for i, v := range self.activityModule.SingleRecharge {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("SingleRecharge DB_Reset fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName1 := fmt.Sprintf("singlerecharge.%d.rechargerecord", index)
	filedName2 := fmt.Sprintf("singlerecharge.%d.singleawardlst", index)
	filedName3 := fmt.Sprintf("singlerecharge.%d.versioncode", index)
	filedName4 := fmt.Sprintf("singlerecharge.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName1: self.RechargeRecord,
		filedName2: self.SingleAwardLst,
		filedName3: self.VersionCode,
		filedName4: self.ResetCode}})
	return true
}

func (self *TActivitySingleRecharge) DB_AddSingleRecharge(index int, info TActivityRechargeInfo) {
	filedName := fmt.Sprintf("singlerecharge.%d.singleawardlst", index)
	mongodb.AddToArray(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, filedName, info)
}

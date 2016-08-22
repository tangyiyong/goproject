package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 七日活动表结构
type TActivitySevenDay struct {
	ActivityID int         //! 活动ID
	TaskList   []TTaskInfo //! 任务列表
	BuyLst     IntLst      //! 已购买限购商品列表

	VersionCode    int              //! 版本号
	ResetCode      int              //! 迭代号
	activityModule *TActivityModule //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivitySevenDay) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivitySevenDay) Init(activityID int, mPtr *TActivityModule, vercode int, resetcode int) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self

	self.TaskList = []TTaskInfo{}
	self.BuyLst = []int{}

	self.VersionCode = vercode
	self.ResetCode = resetcode

	awardType := G_GlobalVariables.GetActivityAwardType(activityID)
	taskLst := gamedata.GetSevenTaskInfoFromAwardType(awardType)
	for _, v := range taskLst {
		var info TTaskInfo
		if v.TaskID == 0 {
			continue
		}
		info.TaskID = v.TaskID
		info.TaskStatus = 0
		info.TaskCount = 0
		info.TaskType = v.TaskType
		self.TaskList = append(self.TaskList, info)
	}

}

//! 刷新数据
func (self *TActivitySevenDay) Refresh(versionCode int) {
	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivitySevenDay) End(versionCode int, resetCode int) {
	self.VersionCode = versionCode
	self.ResetCode = resetCode

	self.TaskList = []TTaskInfo{}
	self.BuyLst = []int{}

	go self.DB_Reset()
}

func (self *TActivitySevenDay) GetRefreshV() int {
	return self.VersionCode
}

func (self *TActivitySevenDay) GetResetV() int {
	return self.ResetCode
}

func (self *TActivitySevenDay) RedTip() bool {
	return false
}

func (self *TActivitySevenDay) DB_Refresh() bool {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_Refresh fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("sevenday.%d.versioncode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.VersionCode}})
	return true
}

func (self *TActivitySevenDay) DB_Reset() bool {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_Reset fail. ActivityID: %d", self.ActivityID)
		return false
	}

	filedName1 := fmt.Sprintf("sevenday.%d.tasklist", index)
	filedName2 := fmt.Sprintf("sevenday.%d.buylst", index)
	filedName3 := fmt.Sprintf("sevenday.%d.resetcode", index)
	filedName4 := fmt.Sprintf("sevenday.%d.versioncode", index)

	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName1: self.TaskList,
		filedName2: self.BuyLst,
		filedName3: self.ResetCode,
		filedName4: self.VersionCode}})
	return true
}

//! 设置玩家任务进度
func (self *TActivitySevenDay) DB_UpdatePlayerSevenTask(taskID int, count int, status int) bool {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_UpdatePlayerSevenTask fail: Not find activityID: %d ", self.ActivityID)
		return false
	}

	indexTask := -1
	for i, v := range self.TaskList {
		if v.TaskID == taskID {
			indexTask = i
			break
		}
	}

	if indexTask < 0 {
		gamelog.Error("Sevenday DB_UpdatePlayerSevenTaskStatus fail: Not find activityID: %d  taskID: %d", self.ActivityID, taskID)
		return false
	}

	filedName := fmt.Sprintf("sevenday.%d.tasklist.%d.taskstatus", index, indexTask)
	filedName2 := fmt.Sprintf("sevenday.%d.tasklist.%d.taskcount", index, indexTask)
	return mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  status,
		filedName2: count}})
}

//! 设置玩家七日活动限购购买标记
func (self *TActivitySevenDay) DB_AddPlayerSevenTaskMark(ID int) {
	index := -1
	for i, v := range self.activityModule.SevenDay {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("Sevenday DB_AddPlayerSevenTaskMark fail")
		return
	}

	filedName1 := fmt.Sprintf("sevenday.%d.buylst", index)
	mongodb.AddToArray(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, filedName1, ID)
}

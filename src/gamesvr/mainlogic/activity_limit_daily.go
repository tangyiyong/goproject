package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 限时日常任务类型
type TLimitDailyTask struct {
	TaskType int //! 任务类型
	Count    int //! 当前次数
	Need     int //! 需要次数
	Status   int //! 状态: 0->未完成 1->已完成 2->已领取
	Award    int //! 奖励
	IsSelect int //! 是否为多选一类型
}

//! 限时日常
type TActivityLimitDaily struct {
	ActivityID     int               //! 活动ID
	TaskLst        []TLimitDailyTask //! 任务链
	VersionCode    int32             //! 版本号
	ResetCode      int32             //! 迭代号
	activityModule *TActivityModule  //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityLimitDaily) SetModulePtr(mPtr *TActivityModule) {
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLimitDaily) Init(activityID int, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.activityModule = mPtr
	self.activityModule.activityPtrs[self.ActivityID] = self
	awardType := G_GlobalVariables.GetActivityAwardType(activityID)
	taskLst := gamedata.GetActivityLimitDaily(awardType)
	self.VersionCode = vercode
	self.ResetCode = resetcode

	for _, n := range taskLst {

		var task TLimitDailyTask
		task.Count = 0
		task.Need = n.Count
		task.TaskType = n.TaskType
		task.Status = 0
		task.Award = n.Award
		task.IsSelect = n.IsSelect

		self.TaskLst = append(self.TaskLst, task)
	}
}

//! 刷新数据
func (self *TActivityLimitDaily) Refresh(versionCode int32) {
	//! 清空限时任务
	for j, _ := range self.TaskLst {
		if self.TaskLst[j].TaskType != gamedata.TASK_RECHARGE {
			self.TaskLst[j].Count = 0
			self.TaskLst[j].Status = 0
		}
	}

	self.VersionCode = versionCode
	go self.DB_Refresh()
}

//! 活动结束
func (self *TActivityLimitDaily) End(versionCode int32, resetCode int32) {
	for j, _ := range self.TaskLst {
		self.TaskLst[j].Count = 0
		self.TaskLst[j].Status = 0
	}

	self.VersionCode = versionCode
	self.ResetCode = resetCode
	go self.DB_Reset()
}

func (self *TActivityLimitDaily) IsAllComplete() bool {
	for _, v := range self.TaskLst {
		if v.TaskType == gamedata.TASK_COMPLETE_ALL_TASK {
			continue
		}
		if v.Status == 0 {
			return false
		}
	}

	return true
}

func (self *TActivityLimitDaily) GetRefreshV() int32 {
	return self.VersionCode
}

func (self *TActivityLimitDaily) GetResetV() int32 {
	return self.ResetCode
}

func (self *TActivityLimitDaily) RedTip() bool {
	//! 活动未开启, 不亮起红点
	if G_GlobalVariables.IsActivityOpen(self.ActivityID) == false {
		return false
	}

	length := len(self.TaskLst)
	for i := 0; i < length; i++ {
		if self.TaskLst[i].Status == 1 {
			return true
		}
	}
	return false
}

func (self *TActivityLimitDaily) DB_Refresh() bool {
	index := -1
	for i, v := range self.activityModule.LimitDaily {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("LimitDaily DB_Refresh fail. self.ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("limitdaily.%d.tasklst", index)
	filedName2 := fmt.Sprintf("limitdaily.%d.versioncode", index)
	filedName3 := fmt.Sprintf("limitdaily.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.TaskLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
	return true
}

func (self *TActivityLimitDaily) DB_Reset() bool {
	index := -1
	for i, v := range self.activityModule.LimitDaily {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("LimitDaily DB_Reset fail. self.ActivityID: %d", self.ActivityID)
		return false
	}

	filedName := fmt.Sprintf("limitdaily.%d.tasklst", index)
	filedName2 := fmt.Sprintf("limitdaily.%d.versioncode", index)
	filedName3 := fmt.Sprintf("limitdaily.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.TaskLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
	return true
}

func (self *TActivityLimitDaily) DB_SaveTask() {
	index := -1
	for i, v := range self.activityModule.LimitDaily {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("LimitDaily DB_SaveTask fail")
		return
	}

	filedName := fmt.Sprintf("limitdaily.%d.tasklst", index)
	filedName2 := fmt.Sprintf("limitdaily.%d.versioncode", index)
	filedName3 := fmt.Sprintf("limitdaily.%d.resetcode", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName:  self.TaskLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
}

func (self *TActivityLimitDaily) DB_UpdateTaskStatus(activityIndex int, taskIndex int) {
	filedName := fmt.Sprintf("limitdaily.%d.tasklst.%d.status", activityIndex, taskIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerActivity", bson.M{"_id": self.activityModule.PlayerID}, bson.M{"$set": bson.M{
		filedName: self.TaskLst[taskIndex].Status}})
}

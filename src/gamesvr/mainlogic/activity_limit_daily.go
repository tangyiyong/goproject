package mainlogic

import (
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
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
	ActivityID  int32             //! 活动ID
	TaskLst     []TLimitDailyTask //! 任务链
	VersionCode int32             //! 版本号
	ResetCode   int32             //! 迭代号
	modulePtr   *TActivityModule  //! 活动模块指针
}

//! 赋值基础数据
func (self *TActivityLimitDaily) SetModulePtr(mPtr *TActivityModule) {
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
}

//! 创建初始化
func (self *TActivityLimitDaily) Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32) {
	delete(mPtr.activityPtrs, self.ActivityID)
	self.ActivityID = activityID
	self.modulePtr = mPtr
	self.modulePtr.activityPtrs[self.ActivityID] = self
	self.VersionCode = vercode
	self.ResetCode = resetcode

	awardType := G_GlobalVariables.GetActivityAwardType(activityID)
	pTaskList := gamedata.GetActivityLimitDaily(awardType)
	self.TaskLst = make([]TLimitDailyTask, len(pTaskList))
	for i, n := range pTaskList {
		self.TaskLst[i].Count = 0
		self.TaskLst[i].Need = n.Count
		self.TaskLst[i].TaskType = n.TaskType
		self.TaskLst[i].Status = 0
		self.TaskLst[i].Award = n.Award
		self.TaskLst[i].IsSelect = n.IsSelect
	}
}

//! 刷新数据
func (self *TActivityLimitDaily) Refresh(versionCode int32) {
	//! 清空限时任务
	for j := 0; j < len(self.TaskLst); j++ {
		if self.TaskLst[j].TaskType != gamedata.TASK_RECHARGE {
			self.TaskLst[j].Count = 0
			self.TaskLst[j].Status = 0
		}
	}

	self.VersionCode = versionCode
	self.DB_Refresh()
}

//! 活动结束
func (self *TActivityLimitDaily) End(versionCode int32, resetCode int32) {
	for j, _ := range self.TaskLst {
		self.TaskLst[j].Count = 0
		self.TaskLst[j].Status = 0
	}

	self.VersionCode = versionCode
	self.ResetCode = resetCode
	self.DB_Reset()
}

func (self *TActivityLimitDaily) IsAllComplete() bool {
	for j := 0; j < len(self.TaskLst); j++ {
		if self.TaskLst[j].TaskType == gamedata.TASK_COMPLETE_ALL_TASK {
			continue
		}
		if self.TaskLst[j].Status == 0 {
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

	for i := 0; i < len(self.TaskLst); i++ {
		if self.TaskLst[i].Status == 1 {
			return true
		}
	}
	return false
}

func (self *TActivityLimitDaily) DB_Refresh() {
	index := -1
	for i, v := range self.modulePtr.LimitDaily {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("LimitDaily DB_Refresh fail. self.ActivityID: %d", self.ActivityID)
		return
	}

	filedName := fmt.Sprintf("limitdaily.%d.tasklst", index)
	filedName2 := fmt.Sprintf("limitdaily.%d.versioncode", index)
	filedName3 := fmt.Sprintf("limitdaily.%d.resetcode", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  self.TaskLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
}

func (self *TActivityLimitDaily) DB_Reset() {
	index := -1
	for i, v := range self.modulePtr.LimitDaily {
		if v.ActivityID == self.ActivityID {
			index = i
			break
		}
	}

	if index < 0 {
		gamelog.Error("LimitDaily DB_Reset fail. self.ActivityID: %d", self.ActivityID)
		return
	}

	filedName := fmt.Sprintf("limitdaily.%d.tasklst", index)
	filedName2 := fmt.Sprintf("limitdaily.%d.versioncode", index)
	filedName3 := fmt.Sprintf("limitdaily.%d.resetcode", index)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName:  self.TaskLst,
		filedName2: self.VersionCode,
		filedName3: self.ResetCode}})
}

func (self *TActivityLimitDaily) DB_UpdateTaskStatus(activityIndex int, taskIndex int) {
	filedName := fmt.Sprintf("limitdaily.%d.tasklst.%d.status", activityIndex, taskIndex)
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.modulePtr.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.TaskLst[taskIndex].Status}})
}

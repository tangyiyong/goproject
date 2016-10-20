package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 设置玩家任务进度
func (taskmodule *TTaskMoudle) DB_UpdateTask(nIndex int) {
	filedName := fmt.Sprintf("tasklist.%d", nIndex)
	mongodb.UpdateToDB("PlayerTask", &bson.M{"_id": taskmodule.PlayerID}, &bson.M{"$set": bson.M{filedName: taskmodule.TaskList[nIndex]}})
}

//! 设置玩家成就进度
func (taskmodule *TTaskMoudle) DB_UpdateAchieve(nIndex int) {
	filedName := fmt.Sprintf("achievelist.%d", nIndex)
	mongodb.UpdateToDB("PlayerTask", &bson.M{"_id": taskmodule.PlayerID}, &bson.M{"$set": bson.M{filedName: taskmodule.AchieveList[nIndex]}})
}

//! 增加玩家成就达成列表
func (taskmodule *TTaskMoudle) DB_AddAchieveID(achievementID int) {
	mongodb.UpdateToDB("PlayerTask", &bson.M{"_id": taskmodule.PlayerID}, &bson.M{"$push": bson.M{"achieveids": achievementID}})
}

//! 设置玩家任务积分
func (taskmodule *TTaskMoudle) DB_UpdateTaskScore(score int) {
	mongodb.UpdateToDB("PlayerTask", &bson.M{"_id": taskmodule.PlayerID}, &bson.M{"$set": bson.M{"taskscore": score}})
}

//! 设置玩家任务积分宝箱领取状态
func (taskmodule *TTaskMoudle) DB_UpdateTaskScoreAwardStatus() {
	mongodb.UpdateToDB("PlayerTask", &bson.M{"_id": taskmodule.PlayerID}, &bson.M{"$set": bson.M{"scoreawardstatus": taskmodule.ScoreAwardStatus}})
}

//! 日常任务信息存储数据库
func (self *TTaskMoudle) DB_UpdateDailyTaskInfo() {
	mongodb.UpdateToDB("PlayerTask", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"taskscore":        self.TaskScore,
		"scoreawardstatus": self.ScoreAwardStatus,
		"scoreawardid":     self.ScoreAwardID,
		"tasklist":         self.TaskList,
		"resetday":         self.ResetDay}})
}

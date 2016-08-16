package gamedata

import (
	"gamelog"
)

type ST_TaskAwardInfo struct {
	TaskAwardID int //! 任务奖励唯一标识
	NeedScore   int //! 需要积分
	Award       int //! 奖励物品ID
	MinLevel    int //! 奖励等级
	MaxLevel    int //! 奖励等级
}

var GT_Task_Award_List []ST_TaskAwardInfo = nil

func InitTaskAwardParser(total int) bool {
	GT_Task_Award_List = make([]ST_TaskAwardInfo, total+1)

	return true
}

func ParseTaskAwardRecord(rs *RecordSet) {
	awardID := rs.GetFieldInt("id")
	GT_Task_Award_List[awardID].TaskAwardID = awardID
	GT_Task_Award_List[awardID].NeedScore = rs.GetFieldInt("needscore")
	GT_Task_Award_List[awardID].Award = rs.GetFieldInt("award")
	GT_Task_Award_List[awardID].MinLevel = rs.GetFieldInt("minlevel")
	GT_Task_Award_List[awardID].MaxLevel = rs.GetFieldInt("maxlevel")
}

//! 返回积分奖励信息
func GetTaskScoreAwardData(awardID int) *ST_TaskAwardInfo {
	if awardID > len(GT_Task_Award_List) || awardID <= 0 {
		gamelog.Error("GetTaskScoreAwardData Error: invalid awardID %d", awardID)
		return nil
	}
	return &GT_Task_Award_List[awardID]
}

//! 根据等级返回对应日常任务宝箱
func GetTaskScoreAwardID(level int) (awardLst []int) {
	for _, v := range GT_Task_Award_List {
		if level >= v.MinLevel && level < v.MaxLevel {
			awardLst = append(awardLst, v.TaskAwardID)
		}
	}

	return awardLst
}

package gamedata

import (
	"gamelog"
)

type ST_CopyBase struct {
	CopyID         int    //! 副本ID
	Type           int    //! 副本类型
	Name           string //! 副本名字
	ActionType     int    //! 消耗行动力类型
	ActionValue    int    //! 消耗行动力值
	MoneyID        int    //! 货币ID
	MoneyNum       int    //! 货币数量
	FirstAward     int    //! 首胜奖励
	AwardID        int    //! 普通奖励
	Experience     int    //! 副本经验
	MaxBattleTimes int    //! 最大挑战次数
}

const (
	COPY_TYPE_Main      = 1  //主线
	COPY_TYPE_Elite     = 2  //精英
	COPY_TYPE_Famous    = 3  //名将
	COPY_TYPE_Daily     = 4  //日常
	COPY_TYPE_GuaJi     = 5  //挂机
	COPY_TYPE_SGWS      = 6  //三国无双
	COPY_TYPE_SGWSJY    = 7  //三国无双精英
	COPY_TYPE_Territory = 8  //领地征讨
	COPY_TYPE_Rebel     = 9  //叛军
	COPY_TYPE_Guild     = 10 //公会副本
	COPY_TYPE_Wander    = 11 //云游副本

)

var (
	GT_CopyBaseList map[int]*ST_CopyBase //主线副本
)

func InitCopyParser(total int) bool {
	GT_CopyBaseList = make(map[int]*ST_CopyBase) //副本表
	return true
}

func ParseCopyRecord(rs *RecordSet) {
	CopyID := CheckAtoi(rs.Values[0], 0)
	copybase := new(ST_CopyBase)
	copybase.CopyID = CopyID
	copybase.Type = rs.GetFieldInt("type")
	copybase.Name = rs.Values[2]
	copybase.ActionType = rs.GetFieldInt("action_type")
	copybase.ActionValue = rs.GetFieldInt("action_cost")
	copybase.MoneyID = rs.GetFieldInt("money_id")
	copybase.MoneyNum = rs.GetFieldInt("money_num")
	copybase.FirstAward = rs.GetFieldInt("firstaward")
	copybase.AwardID = rs.GetFieldInt("awardid")
	copybase.Experience = rs.GetFieldInt("experience")
	copybase.MaxBattleTimes = rs.GetFieldInt("maxbattletimes")
	GT_CopyBaseList[CopyID] = copybase
}

func GetCopyBaseInfo(copyid int) *ST_CopyBase {
	pCopyBase, ok := GT_CopyBaseList[copyid]
	if pCopyBase == nil || !ok {
		gamelog.Error("GetCopyBaseInfo Error: Invalid copyid :%d", copyid)
		return nil
	}
	return pCopyBase
}

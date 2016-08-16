package gamedata

//行动力配制表
//配制了所有行动类型和属性

import (
	"gamelog"
)

type ST_ActionInfo struct {
	ActionID int //行动力ID
	UnitTime int //恢复1个单位需要的时间(秒)
	Max      int //最大值
}

var GT_Action_List []ST_ActionInfo = nil

func InitActionParser(total int) bool {
	GT_Action_List = make([]ST_ActionInfo, total+1)
	return true
}

func ParseActionRecord(rs *RecordSet) {
	actionID := rs.GetFieldInt("id")
	GT_Action_List[actionID].ActionID = actionID
	GT_Action_List[actionID].UnitTime = rs.GetFieldInt("unittime")
	GT_Action_List[actionID].Max = rs.GetFieldInt("max")
}

func GetActionCount() int {
	return len(GT_Action_List) - 1
}

func GetActionInfo(actionID int) *ST_ActionInfo {
	if actionID >= len(GT_Action_List) || actionID == 0 {
		gamelog.Error("GetActionInfo Error: invalid actionID :%d", actionID)
		return nil
	}

	return &GT_Action_List[actionID]
}

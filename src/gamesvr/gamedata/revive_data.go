package gamedata

//阵营战复活配制表

import (
	"gamelog"
)

type ST_ReviveInfo struct {
	ID           int //水晶ID
	HpRatio      int //复血的比例
	CostMoneyID  int //花费货币ID
	CostMoneyNum int //花费货币数
	BuffTime     int //buff持续时间
	IncRatio     int //加强属性比例
	Stay         int //是否原地复活
}

var GT_Revive_List []ST_ReviveInfo = nil

func InitReviveParser(total int) bool {
	GT_Revive_List = make([]ST_ReviveInfo, total+1)
	return true
}

func ParseReviveRecord(rs *RecordSet) {
	ID := rs.GetFieldInt("id")
	GT_Revive_List[ID].ID = ID
	GT_Revive_List[ID].HpRatio = rs.GetFieldInt("hp_ratio")
	GT_Revive_List[ID].IncRatio = rs.GetFieldInt("inc_ratio")
	GT_Revive_List[ID].BuffTime = rs.GetFieldInt("buff_time")
	GT_Revive_List[ID].Stay = rs.GetFieldInt("stay")
	GT_Revive_List[ID].CostMoneyID = rs.GetFieldInt("cost_money_id")
	GT_Revive_List[ID].CostMoneyNum = rs.GetFieldInt("cost_money_num")
}

func GetReviveInfo(id int) *ST_ReviveInfo {
	if id >= len(GT_Revive_List) || id == 0 {
		gamelog.Error("GetReviveInfo Error: invalid id :%d", id)
		return nil
	}

	return &GT_Revive_List[id]
}

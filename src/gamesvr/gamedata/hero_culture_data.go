package gamedata

import (
	"gamelog"
)

type ST_CultureMaxItem struct {
	AttackType int    //攻击类型
	MaxRation  [5]int //五个属性最大值的等级因子
}

var (
	GT_CultureMaxList []ST_CultureMaxItem //培养化费表
)

func InitCultureMaxParser(total int) bool {
	GT_CultureMaxList = make([]ST_CultureMaxItem, total+1)
	return true
}

func ParseCultureMaxRecord(rs *RecordSet) {
	attacktype := rs.GetFieldInt("attacktype")
	GT_CultureMaxList[attacktype].AttackType = attacktype
	GT_CultureMaxList[attacktype].MaxRation[0] = rs.GetFieldInt("p1_max")
	GT_CultureMaxList[attacktype].MaxRation[1] = rs.GetFieldInt("p2_max")
	GT_CultureMaxList[attacktype].MaxRation[2] = rs.GetFieldInt("p3_max")
	GT_CultureMaxList[attacktype].MaxRation[3] = rs.GetFieldInt("p4_max")
	GT_CultureMaxList[attacktype].MaxRation[4] = rs.GetFieldInt("p5_max")
}

func GetCultureMaxInfo(attacktype int) *ST_CultureMaxItem {
	if attacktype >= len(GT_CultureMaxList) || attacktype <= 0 {
		gamelog.Error("GetCultureMaxInfo Error : Invalid attacktype  %d", attacktype)
		return nil
	}

	return &GT_CultureMaxList[attacktype]
}

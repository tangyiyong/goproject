package gamedata

import (
	"gamelog"
)

type ST_EquipStrengCost struct {
	Level    int     //等级
	MoneyID  int     //货币ID
	MoneyNum [10]int //货币数量
}

var (
	GT_EquipStrengCostList []ST_EquipStrengCost = nil
)

func InitEquipStrengthCostParser(total int) bool {
	GT_EquipStrengCostList = make([]ST_EquipStrengCost, total+1)

	return true
}

func ParseEquipStrengthCostRecord(rs *RecordSet) {
	Level := rs.GetFieldInt("level")
	GT_EquipStrengCostList[Level].Level = Level
	GT_EquipStrengCostList[Level].MoneyID = CheckAtoi(rs.Values[1], 1)
	GT_EquipStrengCostList[Level].MoneyNum[0] = CheckAtoi(rs.Values[2], 2)
	GT_EquipStrengCostList[Level].MoneyNum[1] = CheckAtoi(rs.Values[3], 3)
	GT_EquipStrengCostList[Level].MoneyNum[2] = CheckAtoi(rs.Values[4], 4)
	GT_EquipStrengCostList[Level].MoneyNum[3] = CheckAtoi(rs.Values[5], 5)
	GT_EquipStrengCostList[Level].MoneyNum[4] = CheckAtoi(rs.Values[6], 6)
	GT_EquipStrengCostList[Level].MoneyNum[5] = CheckAtoi(rs.Values[7], 7)
	GT_EquipStrengCostList[Level].MoneyNum[6] = CheckAtoi(rs.Values[8], 8)
	GT_EquipStrengCostList[Level].MoneyNum[7] = CheckAtoi(rs.Values[9], 9)
	GT_EquipStrengCostList[Level].MoneyNum[8] = CheckAtoi(rs.Values[10], 10)
	GT_EquipStrengCostList[Level].MoneyNum[9] = CheckAtoi(rs.Values[11], 11)

	return
}

func GetEquipStrengthCostInfo(level int) *ST_EquipStrengCost {
	if level >= len(GT_EquipStrengCostList) || level <= 0 {
		gamelog.Error("GetEquipStrengCostInfo Error: invalid level %d", level)
		return nil
	}

	return &GT_EquipStrengCostList[level]
}

package gamedata

import (
	"gamelog"
)

type ST_EquipRefineCost struct {
	Level    int     //等级
	NeedExp  [10]int //依据品质
	TotalExp [10]int
}

var (
	GT_EquipRefineCostList []ST_EquipRefineCost = nil
)

func InitEquipRefineCostParser(total int) bool {
	GT_EquipRefineCostList = make([]ST_EquipRefineCost, total)

	return true
}

func ParseEquipRefineCostRecord(rs *RecordSet) {
	level := CheckAtoi(rs.Values[0], 0)
	GT_EquipRefineCostList[level].Level = level
	GT_EquipRefineCostList[level].NeedExp[0] = CheckAtoi(rs.Values[1], 1)
	GT_EquipRefineCostList[level].NeedExp[1] = CheckAtoi(rs.Values[2], 2)
	GT_EquipRefineCostList[level].NeedExp[2] = CheckAtoi(rs.Values[3], 3)
	GT_EquipRefineCostList[level].NeedExp[3] = CheckAtoi(rs.Values[4], 4)
	GT_EquipRefineCostList[level].NeedExp[4] = CheckAtoi(rs.Values[5], 5)
	GT_EquipRefineCostList[level].NeedExp[5] = CheckAtoi(rs.Values[6], 6)
	GT_EquipRefineCostList[level].NeedExp[6] = CheckAtoi(rs.Values[7], 7)
	GT_EquipRefineCostList[level].NeedExp[7] = CheckAtoi(rs.Values[8], 8)
	GT_EquipRefineCostList[level].NeedExp[8] = CheckAtoi(rs.Values[9], 9)
	GT_EquipRefineCostList[level].NeedExp[9] = CheckAtoi(rs.Values[10], 10)

	GT_EquipRefineCostList[level].TotalExp[0] = CheckAtoi(rs.Values[11], 11)
	GT_EquipRefineCostList[level].TotalExp[1] = CheckAtoi(rs.Values[12], 12)
	GT_EquipRefineCostList[level].TotalExp[2] = CheckAtoi(rs.Values[13], 13)
	GT_EquipRefineCostList[level].TotalExp[3] = CheckAtoi(rs.Values[14], 14)
	GT_EquipRefineCostList[level].TotalExp[4] = CheckAtoi(rs.Values[15], 15)
	GT_EquipRefineCostList[level].TotalExp[5] = CheckAtoi(rs.Values[16], 16)
	GT_EquipRefineCostList[level].TotalExp[6] = CheckAtoi(rs.Values[17], 17)
	GT_EquipRefineCostList[level].TotalExp[7] = CheckAtoi(rs.Values[18], 18)
	GT_EquipRefineCostList[level].TotalExp[8] = CheckAtoi(rs.Values[19], 19)
	GT_EquipRefineCostList[level].TotalExp[9] = CheckAtoi(rs.Values[20], 20)
	return
}

func GetEquipRefineCostInfo(level int) *ST_EquipRefineCost {
	if level >= len(GT_EquipRefineCostList) {
		gamelog.Error("GetEquipRefineCostInfo Error : Invalid level %d", level)
		return nil
	}

	return &GT_EquipRefineCostList[level]
}

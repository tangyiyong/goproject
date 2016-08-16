package gamedata

import (
	"gamelog"
)

type ST_GemStrengthCost struct {
	Level    int     //等级
	MoneyID  int     //货币ID
	MoneyNum int     //*经验数
	NeedExp  [10]int //需要的经验数
	TotalExp [10]int
}

var (
	GT_GemStrengthCostList []ST_GemStrengthCost = nil
)

func InitGemStrengthCostParser(total int) bool {

	GT_GemStrengthCostList = make([]ST_GemStrengthCost, total+1)

	return true
}

func ParseGemStrengthCostRecord(rs *RecordSet) {
	level := rs.GetFieldInt("level")

	GT_GemStrengthCostList[level].Level = level
	GT_GemStrengthCostList[level].MoneyID = rs.GetFieldInt("money_id")
	GT_GemStrengthCostList[level].MoneyNum = rs.GetFieldInt("money_num")
	GT_GemStrengthCostList[level].NeedExp[0] = CheckAtoi(rs.Values[3], 3)
	GT_GemStrengthCostList[level].NeedExp[1] = CheckAtoi(rs.Values[4], 4)
	GT_GemStrengthCostList[level].NeedExp[2] = CheckAtoi(rs.Values[5], 5)
	GT_GemStrengthCostList[level].NeedExp[3] = CheckAtoi(rs.Values[6], 6)
	GT_GemStrengthCostList[level].NeedExp[4] = CheckAtoi(rs.Values[7], 7)
	GT_GemStrengthCostList[level].NeedExp[5] = CheckAtoi(rs.Values[8], 8)
	GT_GemStrengthCostList[level].NeedExp[6] = CheckAtoi(rs.Values[9], 9)
	GT_GemStrengthCostList[level].NeedExp[7] = CheckAtoi(rs.Values[10], 10)
	GT_GemStrengthCostList[level].NeedExp[8] = CheckAtoi(rs.Values[11], 11)
	GT_GemStrengthCostList[level].NeedExp[9] = CheckAtoi(rs.Values[12], 12)

	GT_GemStrengthCostList[level].TotalExp[0] = CheckAtoi(rs.Values[13], 13)
	GT_GemStrengthCostList[level].TotalExp[1] = CheckAtoi(rs.Values[14], 14)
	GT_GemStrengthCostList[level].TotalExp[2] = CheckAtoi(rs.Values[15], 15)
	GT_GemStrengthCostList[level].TotalExp[3] = CheckAtoi(rs.Values[16], 16)
	GT_GemStrengthCostList[level].TotalExp[4] = CheckAtoi(rs.Values[17], 17)
	GT_GemStrengthCostList[level].TotalExp[5] = CheckAtoi(rs.Values[18], 18)
	GT_GemStrengthCostList[level].TotalExp[6] = CheckAtoi(rs.Values[19], 19)
	GT_GemStrengthCostList[level].TotalExp[7] = CheckAtoi(rs.Values[20], 20)
	GT_GemStrengthCostList[level].TotalExp[8] = CheckAtoi(rs.Values[21], 21)
	GT_GemStrengthCostList[level].TotalExp[9] = CheckAtoi(rs.Values[22], 22)

	return
}

func GetGemStrengthCostInfo(level int) *ST_GemStrengthCost {
	if level >= len(GT_GemStrengthCostList) || level <= 0 {
		gamelog.Error("GetGemStrengthCostInfo Error: invalid level %d", level)
		return nil
	}

	return &GT_GemStrengthCostList[level]
}

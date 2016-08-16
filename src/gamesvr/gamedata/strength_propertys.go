package gamedata

import (
	"gamelog"
)

type ST_StrengthInfo struct {
	Quality     int       //物品品质
	PropertyInc [6][2]int //不同部位的属性等级增加值
}

var (
	GT_StrengthList []ST_StrengthInfo = nil
)

func InitStrengthParser(total int) bool {
	GT_StrengthList = make([]ST_StrengthInfo, total+1)

	return true
}

func ParseStrengthRecord(rs *RecordSet) {
	Quality := rs.GetFieldInt("quality")
	GT_StrengthList[Quality].Quality = Quality
	GT_StrengthList[Quality].PropertyInc[0][0] = CheckAtoi(rs.Values[1], 1)
	GT_StrengthList[Quality].PropertyInc[1][0] = CheckAtoi(rs.Values[2], 2)
	GT_StrengthList[Quality].PropertyInc[2][0] = CheckAtoi(rs.Values[3], 3)
	GT_StrengthList[Quality].PropertyInc[3][0] = CheckAtoi(rs.Values[4], 4)
	GT_StrengthList[Quality].PropertyInc[4] = ParseTo2IntSlice(rs.Values[5])
	GT_StrengthList[Quality].PropertyInc[5] = ParseTo2IntSlice(rs.Values[6])

	return
}

func GetStrengthInfo(quality int) *ST_StrengthInfo {
	if quality > len(GT_StrengthList) || quality <= 0 {
		gamelog.Error("GetStrengthInfo Error: invalid quality %d", quality)
		return nil
	}

	return &GT_StrengthList[quality]
}

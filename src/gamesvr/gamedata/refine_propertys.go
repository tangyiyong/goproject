package gamedata

import (
	"gamelog"
)

type ST_RefineInfo struct {
	Quality     int       //物品品质
	PropertyInc [6][2]int //装备的位置属性改变
}

var (
	GT_RefineList []ST_RefineInfo = nil
)

func InitRefineParser(total int) bool {
	GT_RefineList = make([]ST_RefineInfo, total+1)

	return true
}

//解析精炼记录
func ParseRefineRecord(rs *RecordSet) {
	Quality := CheckAtoi(rs.Values[0], 0)
	GT_RefineList[Quality].Quality = Quality

	GT_RefineList[Quality].PropertyInc[0] = ParseTo2IntSlice(rs.Values[1])
	GT_RefineList[Quality].PropertyInc[1] = ParseTo2IntSlice(rs.Values[2])
	GT_RefineList[Quality].PropertyInc[2] = ParseTo2IntSlice(rs.Values[3])
	GT_RefineList[Quality].PropertyInc[3] = ParseTo2IntSlice(rs.Values[4])
	GT_RefineList[Quality].PropertyInc[4] = ParseTo2IntSlice(rs.Values[5])
	GT_RefineList[Quality].PropertyInc[5] = ParseTo2IntSlice(rs.Values[6])

	return
}

//获取精炼信息
func GetRefineInfo(quality int) *ST_RefineInfo {
	if quality >= len(GT_RefineList) || quality <= 0 {
		gamelog.Error("GetRefineInfo Error: invalid quality %d", quality)
		return nil
	}
	return &GT_RefineList[quality]
}

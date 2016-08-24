package gamedata

import (
	"gamelog"
)

const (
	Sanguo_Add_Attr     = 1 //! 全队增加属性
	Sanguo_Give_Item    = 2 //! 给予道具
	Sanguo_Main_Hero_Up = 3 //! 主角品质提升
)

//! 三国志
type ST_SanGuoZhiInfo struct {
	ID       int  //! 唯一标识
	Index    int  //! 星宿
	Type     int  //! 提升种类
	AttrID   int  //! 属性ID
	Value    int8 //! 提升数量
	CostType int  //! 花费道具ID
	CostNum  int  //! 花费道具数量
}

var GT_SanGuoZhiList []ST_SanGuoZhiInfo

//! 初始化
func InitSanGuoZhiParser(total int) bool {
	GT_SanGuoZhiList = make([]ST_SanGuoZhiInfo, total+1)
	return true
}

//! 分析CSV
func ParseSanGuoZhiRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)

	GT_SanGuoZhiList[id].ID = id
	GT_SanGuoZhiList[id].Index = rs.GetFieldInt("index")
	GT_SanGuoZhiList[id].Type = rs.GetFieldInt("type")
	GT_SanGuoZhiList[id].AttrID = rs.GetFieldInt("attrid")
	GT_SanGuoZhiList[id].Value = int8(rs.GetFieldInt("value"))
	GT_SanGuoZhiList[id].CostType = rs.GetFieldInt("costtype")
	GT_SanGuoZhiList[id].CostNum = rs.GetFieldInt("costnum")
}

func GetSanGuoZhiInfo(id int) *ST_SanGuoZhiInfo {
	if id >= len(GT_SanGuoZhiList) || id <= 0 {
		gamelog.Error("GetSanGuoZhiInfo Error: invalid id %d", id)
		return nil
	}

	return &GT_SanGuoZhiList[id]
}

func IsStarEnd(starID int) bool {
	if starID == len(GT_SanGuoZhiList)-1 {
		return true
	}
	return false
}

//! 获取星宿数量
func GetSanGoIndexNumber() int {
	indexLst := []int{}

	for _, v := range GT_SanGuoZhiList {
		isExist := false
		for _, b := range indexLst {
			if b == v.Index {
				isExist = true
				break
			}
		}

		if isExist == false {
			indexLst = append(indexLst, v.Index)
		}
	}

	return len(indexLst)
}

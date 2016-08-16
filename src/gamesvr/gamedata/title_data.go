package gamedata

import (
	"fmt"
	"gamelog"
)

type ST_TitleProperty struct {
	PropertyID int
	IsPercent  bool
	Value      int
}

//! 称号静态表
type ST_Title struct {
	ID         int                 //! 称号ID
	Property   [3]ST_TitleProperty //! 属性
	IsAll      bool                //! 是否全体有效
	Time       int
	CostItemID int //! 消耗道具ID
}

var GT_TitleLst []ST_Title

func InitTitleParser(total int) bool {
	GT_TitleLst = make([]ST_Title, total+1)
	return true
}

func ParseTitleRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)

	for i := 1; i <= 3; i++ {
		filedName := fmt.Sprintf("property%d", i)
		GT_TitleLst[id].Property[i-1].PropertyID = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("is_percent%d", i)
		GT_TitleLst[id].Property[i-1].IsPercent = (rs.GetFieldInt(filedName) == 1)

		filedName = fmt.Sprintf("value%d", i)
		GT_TitleLst[id].Property[i-1].Value = rs.GetFieldInt(filedName)
	}

	GT_TitleLst[id].ID = id
	GT_TitleLst[id].IsAll = (rs.GetFieldInt("is_all") == 1)
	GT_TitleLst[id].Time = rs.GetFieldInt("time")
	GT_TitleLst[id].CostItemID = rs.GetFieldInt("cost_item")
}

//! 根据ID获取称号信息
func GetTitleInfo(id int) *ST_Title {
	if id > len(GT_TitleLst)-1 {
		gamelog.Error("GetTitleInfo Error: invalid id %d", id)
		return nil
	}

	return &GT_TitleLst[id]
}

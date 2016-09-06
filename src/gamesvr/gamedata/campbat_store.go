package gamedata

import (
	"gamelog"
)

//! 阵营战商店
type ST_CampBatStoreItem struct {
	ID           int //条目ID
	Type         int //类型
	ItemID       int
	ItemNum      int
	CostMoneyID  int
	CostMoneyNum int
	NeedScore    int
	MaxBuyTime   int
	NeedLevel    int
}

var (
	GT_CampBatStore_List []ST_CampBatStoreItem = nil
)

//! 初始化阵营战分析器
func InitCampBatStoreParser(total int) bool {
	GT_CampBatStore_List = make([]ST_CampBatStoreItem, total+1)
	return true
}

func ParseCampBatStoreRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_CampBatStore_List[id].ID = id
	GT_CampBatStore_List[id].Type = rs.GetFieldInt("type")
	GT_CampBatStore_List[id].ItemID = rs.GetFieldInt("itemid")
	GT_CampBatStore_List[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_CampBatStore_List[id].CostMoneyID = rs.GetFieldInt("cost_moneyid")
	GT_CampBatStore_List[id].CostMoneyNum = rs.GetFieldInt("cost_moneynum")
	GT_CampBatStore_List[id].NeedScore = rs.GetFieldInt("needscore")
	GT_CampBatStore_List[id].MaxBuyTime = rs.GetFieldInt("buytimes")
	GT_CampBatStore_List[id].NeedLevel = rs.GetFieldInt("needlevel")
}

func GetCampBatStoreItem(id int) *ST_CampBatStoreItem {
	if id >= len(GT_CampBatStore_List) || id == 0 {
		gamelog.Error("GetCampBatStoreItem Error: invalid id :%d", id)
		return nil
	}

	return &GT_CampBatStore_List[id]
}

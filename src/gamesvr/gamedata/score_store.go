package gamedata

import (
	"gamelog"
)

//! 积分赛段位奖励
type ST_ScoreStoreItem struct {
	ID           int //条目ID
	ItemID       int
	ItemNum      int
	CostMoneyID  int
	CostMoneyNum int
	CostItemID   int
	CostItemNum  int
	MaxBuyTime   int
}

var (
	GT_ScoreStore_List []ST_ScoreStoreItem = nil
)

//! 初始化积分赛分析器
func InitScoreStoreParser(total int) bool {
	GT_ScoreStore_List = make([]ST_ScoreStoreItem, total+1)
	return true
}

func ParseScoreStoreRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_ScoreStore_List[id].ID = id
	GT_ScoreStore_List[id].ItemID = rs.GetFieldInt("itemid")
	GT_ScoreStore_List[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_ScoreStore_List[id].CostMoneyID = rs.GetFieldInt("cost_moneyid")
	GT_ScoreStore_List[id].CostMoneyNum = rs.GetFieldInt("cost_moneynum")
	GT_ScoreStore_List[id].CostItemID = rs.GetFieldInt("cost_item")
	GT_ScoreStore_List[id].CostItemNum = rs.GetFieldInt("cost_itemnum")
	GT_ScoreStore_List[id].MaxBuyTime = rs.GetFieldInt("buytimes")
}

func GetScoreStoreItem(id int) *ST_ScoreStoreItem {
	if id >= len(GT_ScoreStore_List) || id == 0 {
		gamelog.Error("GetScoreStoreItem Error: invalid id :%d", id)
		return nil
	}

	return &GT_ScoreStore_List[id]
}

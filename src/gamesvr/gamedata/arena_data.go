package gamedata

import (
	"gamelog"
)

//! 排行榜基本信息
type ST_ArenaInfo struct {
	CopyID             int //! 副本基本信息
	VictoryMoneyID     int //! 胜利奖励货币ID
	VictoryMoneyNum    int //! 胜利奖励货币数量
	FailedMoneyID      int //! 失败奖励货币ID
	FailedMoneyNum     int //! 失败奖励货币数量
	DailyAwardNeedRank int //! 每日排名奖励的排名数
}

//! 排行榜排名奖励
type ST_ArenaRankInfo struct {
	ID          int //! 唯一标识
	RankLowest  int //! 排行奖励区间
	RankHighest int //! 排行奖励区间
	RankAward   int //! 排名奖励 值取award表
}

//! 排行榜商店
type ST_ArenaStoreInfo struct {
	ID          int //! 唯一标识
	Type        int //! 0->商品 1->奖励
	ItemID      int //! 货物ID
	ItemNum     int //! 货物数量
	MoneyID     int //! 货币ID
	MoneyNum    int //! 货币数量
	CostItemID  int //! 需要道具ID
	CostItemNum int //! 需要道具数量
	NeedRank    int //! 需求排名
	NeedLevel   int //! 需求等级
}

//! 竞技场排名上升奖励货币表
type ST_ArenaMoneyAwardInfo struct {
	RankHigh int
	RankLow  int
	MoneyID  int
	MoneyNum int
}

var (
	Arena_Base         ST_ArenaInfo
	GT_ArenaRank_List  []ST_ArenaRankInfo
	GT_ArenaStore_List []ST_ArenaStoreInfo
	GT_ArenaMoney_List []ST_ArenaMoneyAwardInfo
)

func InitArenaParser(total int) bool {
	return true
}

func InitArenaRankParser(total int) bool {
	GT_ArenaRank_List = make([]ST_ArenaRankInfo, total+1)
	return true
}

func InitArenaStoreParser(total int) bool {
	GT_ArenaStore_List = make([]ST_ArenaStoreInfo, total+1)
	return true
}

func InitArenaMoneyParser(total int) bool {
	//GT_ArenaMoney_List = make([]ST_ArenaMoneyAwardInfo, total+1)
	return true
}

func ParseArenaMoneyRecord(rs *RecordSet) {
	var info ST_ArenaMoneyAwardInfo
	info.RankHigh = rs.GetFieldInt("money_award_rank_high")
	info.RankLow = rs.GetFieldInt("money_award_rank_low")
	info.MoneyID = rs.GetFieldInt("award_money_id")
	info.MoneyNum = rs.GetFieldInt("award_money_num")

	GT_ArenaMoney_List = append(GT_ArenaMoney_List, info)
}

func GetArenaMoneyAward(oldRank int, newRank int) (int, int) {
	oldSection := 0
	newSection := 0
	moneyNum := 0

	isHaveAward := false
	for i, v := range GT_ArenaMoney_List {
		if v.RankLow <= oldRank && oldRank <= v.RankHigh {
			oldSection = i
		}

		if v.RankLow <= newRank && newRank <= v.RankHigh {
			newSection = i
			isHaveAward = true
		}
	}

	if isHaveAward == false {
		return 0, 0
	}

	max := oldRank
	if oldSection == 0 {
		oldSection = len(GT_ArenaMoney_List) - 1
		max = GT_ArenaMoney_List[oldSection].RankHigh
	}

	for i := oldSection; i >= newSection; i-- {

		min := GT_ArenaMoney_List[i].RankLow
		if min < newRank {
			min = newRank
		}

		money := (max - GT_ArenaMoney_List[i].RankLow)
		if money < 0 {
			money *= -1
		}
		moneyNum += (money * GT_ArenaMoney_List[i].MoneyNum)
		max = GT_ArenaMoney_List[i].RankLow

	}

	return GT_ArenaMoney_List[0].MoneyID, moneyNum
}

func ParseArenaRecord(rs *RecordSet) {
	Arena_Base.CopyID = rs.GetFieldInt("copy_id")
	Arena_Base.VictoryMoneyID = rs.GetFieldInt("victory_money_id")
	Arena_Base.VictoryMoneyNum = rs.GetFieldInt("victory_money_num")
	Arena_Base.FailedMoneyID = rs.GetFieldInt("fail_money_id")
	Arena_Base.FailedMoneyNum = rs.GetFieldInt("fail_money_num")
	Arena_Base.DailyAwardNeedRank = rs.GetFieldInt("daily_award_need_rank")
}

func ParseArenaRankRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_ArenaRank_List[id].ID = id
	GT_ArenaRank_List[id].RankLowest = rs.GetFieldInt("ranklowest")
	GT_ArenaRank_List[id].RankHighest = rs.GetFieldInt("rankhighest")
	GT_ArenaRank_List[id].RankAward = rs.GetFieldInt("rankaward")
}

func ParseArenaStoreRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_ArenaStore_List[id].ID = id
	GT_ArenaStore_List[id].Type = rs.GetFieldInt("type")
	GT_ArenaStore_List[id].ItemID = rs.GetFieldInt("itemid")
	GT_ArenaStore_List[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_ArenaStore_List[id].MoneyID = rs.GetFieldInt("cost_moneyid")
	GT_ArenaStore_List[id].MoneyNum = rs.GetFieldInt("cost_moneynum")
	GT_ArenaStore_List[id].CostItemID = rs.GetFieldInt("cost_itemid")
	GT_ArenaStore_List[id].CostItemNum = rs.GetFieldInt("cost_itemnum")
	GT_ArenaStore_List[id].NeedRank = rs.GetFieldInt("needrank")
	GT_ArenaStore_List[id].NeedLevel = rs.GetFieldInt("needlevel")
}

func GetArenaConfig() *ST_ArenaInfo {
	return &Arena_Base
}

func GetArenaStoreItem(id int) *ST_ArenaStoreInfo {
	if id >= len(GT_ArenaStore_List) || id <= 0 {
		gamelog.Error("GetArenaStoreItem fail. ID:%d", id)
		return nil
	}

	return &GT_ArenaStore_List[id]
}

func GetArenaRankAward(rank int) int {
	for i, v := range GT_ArenaRank_List {
		if rank >= v.RankLowest && rank <= v.RankHighest {
			return GT_ArenaRank_List[i].RankAward
		}
	}

	return 0
}

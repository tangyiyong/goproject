package gamedata

import (
	"gamelog"
)

type ST_GemRefineCost struct {
	Level      int //等级
	MoneyID    int //货币
	MoneyNum   int //需要货币值
	GemNum     int //同名宝物数量
	ItemID     int //物品的ID
	ItemNum    int //物品数量
	TotalMoney int
	TotalGem   int
	TotalItem  int
}

var (
	GT_GemRefineCostList []ST_GemRefineCost = nil
)

func InitGemRefineCostParser(total int) bool {
	GT_GemRefineCostList = make([]ST_GemRefineCost, total)

	return true
}

//解析精炼记录
func ParseGemRefineCostRecord(rs *RecordSet) {
	Level := rs.GetFieldInt("level")
	GT_GemRefineCostList[Level].MoneyID = rs.GetFieldInt("money_id")
	GT_GemRefineCostList[Level].MoneyNum = rs.GetFieldInt("money_num")
	GT_GemRefineCostList[Level].GemNum = rs.GetFieldInt("gem_num")
	GT_GemRefineCostList[Level].ItemID = rs.GetFieldInt("item_id")
	GT_GemRefineCostList[Level].ItemNum = rs.GetFieldInt("item_num")
	GT_GemRefineCostList[Level].TotalItem = rs.GetFieldInt("total_item")
	GT_GemRefineCostList[Level].TotalGem = rs.GetFieldInt("total_gem")
	GT_GemRefineCostList[Level].TotalMoney = rs.GetFieldInt("total_money")

	return
}

//获取精炼信息
func GetGemRefineCostInfo(level int) *ST_GemRefineCost {
	if level >= len(GT_GemRefineCostList) || level < 0 {
		gamelog.Error("GetGemRefineCostInfo Error: invalid level %d", level)
		return nil
	}

	return &GT_GemRefineCostList[level]
}

package gamedata

import (
	"gamelog"
)

type ST_BreakPropertys [5]int

type ST_BreakLevelInfo struct {
	Level         int //等级
	MoneyID       int //货币ID
	MoneyNum      int //货币数
	HeroNum       int //英雄数
	ItemID        int //道具ID
	ItemNum       int //物品数
	HostItemNum   int //主角需要的物品数
	NeedLevel     int //需求的主角等级
	IncPercent    int //加成百分比
	TotalMoneyNum int //花费货币
	TotalHeroNum  int //花费英雄
	TotalItemNum  int //花费道具
}

var (
	GT_HeroBreak_List []ST_BreakLevelInfo = nil
)

func InitHeroBreakParser(total int) bool {
	GT_HeroBreak_List = make([]ST_BreakLevelInfo, total+1)

	return true
}

func ParseHeroBreakRecord(rs *RecordSet) {
	Level := rs.GetFieldInt("level")
	GT_HeroBreak_List[Level].MoneyID = rs.GetFieldInt("money_id")
	GT_HeroBreak_List[Level].ItemID = rs.GetFieldInt("item_id")
	GT_HeroBreak_List[Level].MoneyNum = rs.GetFieldInt("money_num")
	GT_HeroBreak_List[Level].HeroNum = rs.GetFieldInt("hero_num")
	GT_HeroBreak_List[Level].ItemNum = rs.GetFieldInt("item_num")
	GT_HeroBreak_List[Level].HostItemNum = rs.GetFieldInt("host_num")
	GT_HeroBreak_List[Level].IncPercent = rs.GetFieldInt("inc_percent")
	GT_HeroBreak_List[Level].NeedLevel = rs.GetFieldInt("require_level")
	GT_HeroBreak_List[Level].TotalMoneyNum = rs.GetFieldInt("total_money_num")
	GT_HeroBreak_List[Level].TotalHeroNum = rs.GetFieldInt("total_hero_num")
	GT_HeroBreak_List[Level].TotalItemNum = rs.GetFieldInt("total_item_num")

	return
}

func GetHeroBreakInfo(level int8) *ST_BreakLevelInfo {
	if int(level) >= len(GT_HeroBreak_List) || level < 0 {
		gamelog.Error("GetHeroBreakInfo Error : Invalid level :%d", level)
		return nil
	}

	return &GT_HeroBreak_List[level]
}

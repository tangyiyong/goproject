package gamedata

import (
	"gamelog"
)

type ST_DestinyItem struct {
	Level        int //等级
	CostItemID   int //消耗物品ID
	OneTimeCost  int //单次消耗量
	PropertyInc  int //当前等级的属性增加值
	UpgradeRatio int //天命成功概率
	Return       int //分解返还道具数目
}

var (
	GT_HeroDestinyItem_List []ST_DestinyItem = nil
)

func InitHeroDestinyParser(total int) bool {
	GT_HeroDestinyItem_List = make([]ST_DestinyItem, total+1)
	return true
}

func ParseHeroDestinyRecord(rs *RecordSet) {
	Level := rs.GetFieldInt("level")
	GT_HeroDestinyItem_List[Level].Level = Level
	GT_HeroDestinyItem_List[Level].CostItemID = rs.GetFieldInt("item_id")
	GT_HeroDestinyItem_List[Level].OneTimeCost = rs.GetFieldInt("one_time_cost")
	GT_HeroDestinyItem_List[Level].PropertyInc = rs.GetFieldInt("inc_percent")
	GT_HeroDestinyItem_List[Level].UpgradeRatio = rs.GetFieldInt("upgrade_ratio")
	GT_HeroDestinyItem_List[Level].Return = rs.GetFieldInt("return")
	return
}

func GetHeroDestinyInfo(level int) *ST_DestinyItem {
	if level >= len(GT_HeroDestinyItem_List) || level < 0 {
		gamelog.Error("GetHeroDestinyInfo Error : Invalid level :%d", level)
		return nil
	}

	return &GT_HeroDestinyItem_List[level]
}

package gamedata

import (
	"gamelog"
)

type ST_GodItem struct {
	Level        int    //等级
	Propertys    [5]int //属性表
	NeedID       int    //需要ID
	NeedType     int    //需求道具类型
	NeedNum      int    //需要数量
	NeedMoneyID  int    //需要的货币ID
	NeedMoneyNum int    //需要的货数量
	NeedItemID   int    //需要的道具ID
	NeedItemNum  int    //需要的道具数量
	TotalSouls   int    //需要的灵魂
	TotalPiece   int    //需要的碎片
	TotalItem    int    //需要消耗的道具
	TotalMoney   int    //需要消耗的货币
}

var (
	GT_HeroGodItem_List []ST_GodItem = nil
)

func InitHeroGodParser(total int) bool {
	GT_HeroGodItem_List = make([]ST_GodItem, total+1)
	return true
}

func ParseHeroGodRecord(rs *RecordSet) {
	Level := rs.GetFieldInt("god_level")

	GT_HeroGodItem_List[Level].Level = Level

	GT_HeroGodItem_List[Level].Propertys[0] = rs.GetFieldInt("property_1")
	GT_HeroGodItem_List[Level].Propertys[1] = rs.GetFieldInt("property_2")
	GT_HeroGodItem_List[Level].Propertys[2] = rs.GetFieldInt("property_3")
	GT_HeroGodItem_List[Level].Propertys[3] = rs.GetFieldInt("property_4")
	GT_HeroGodItem_List[Level].Propertys[4] = rs.GetFieldInt("property_5")
	GT_HeroGodItem_List[Level].NeedID = rs.GetFieldInt("need_id")
	GT_HeroGodItem_List[Level].NeedNum = rs.GetFieldInt("need_num")
	GT_HeroGodItem_List[Level].NeedType = rs.GetFieldInt("need_type")
	GT_HeroGodItem_List[Level].NeedMoneyID = rs.GetFieldInt("need_money_id")
	GT_HeroGodItem_List[Level].NeedMoneyNum = rs.GetFieldInt("need_money_num")
	GT_HeroGodItem_List[Level].NeedItemID = rs.GetFieldInt("need_item_id")
	GT_HeroGodItem_List[Level].NeedItemNum = rs.GetFieldInt("need_item_num")
	GT_HeroGodItem_List[Level].TotalSouls = rs.GetFieldInt("total_souls")
	GT_HeroGodItem_List[Level].TotalPiece = rs.GetFieldInt("total_piece")
	GT_HeroGodItem_List[Level].TotalItem = rs.GetFieldInt("total_item")
	GT_HeroGodItem_List[Level].TotalMoney = rs.GetFieldInt("total_money")

	return
}

func GetHeroGodInfo(level int) *ST_GodItem {
	if level >= len(GT_HeroGodItem_List) || level < 0 {
		gamelog.Error("GetHeroGodInfo Error : Invalid level :%d", level)
		return nil
	}

	return &GT_HeroGodItem_List[level]
}

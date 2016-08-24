package gamedata

import (
	"gamelog"
)

type ST_HeroInfo struct {
	HeroID         int    //英雄ID
	ItemID         int    //对应的道具ID
	Camp           int8   //阵营
	Quality        int8   //品质
	AttackType     int    //攻击类型   1 :物攻  2:法攻
	SellID         int    //出售货币ID
	SellPrice      int    //出售价格
	DecomposeID    int    //分解的货币ID
	DecomposePrice int    //分解的货币价格
	PieceID        int    //碎片道具ID
	PieceNum       int    //合成碎片数
	HeroExp        int    //经验
	Setup          int    //是否可以上阵
	CanSell        bool   //是否可出售
	ModelID        int    //模型ID
	BasePropertys  [5]int //基础属性
	Skills         [5]int //拥有技能
}

var (
	GT_Hero_List []ST_HeroInfo = nil
)

func InitHeroParser(total int) bool {
	GT_Hero_List = make([]ST_HeroInfo, 1500)
	return true
}

func ParseHeroRecord(rs *RecordSet) {
	HeroID := rs.GetFieldInt("id")
	if (HeroID <= 0) || (HeroID >= len(GT_Hero_List)) {
		gamelog.Error("ParseHeroRecord Error: Invalid heroid :%d", HeroID)
		return
	}
	GT_Hero_List[HeroID].HeroID = HeroID
	GT_Hero_List[HeroID].ItemID = rs.GetFieldInt("itemid")
	GT_Hero_List[HeroID].Camp = int8(rs.GetFieldInt("camp"))
	GT_Hero_List[HeroID].Quality = int8(rs.GetFieldInt("quality"))
	GT_Hero_List[HeroID].AttackType = rs.GetFieldInt("attacktype")
	GT_Hero_List[HeroID].SellID = rs.GetFieldInt("sell_money_id_1")
	GT_Hero_List[HeroID].SellPrice = rs.GetFieldInt("sell_money_num_1")
	GT_Hero_List[HeroID].DecomposeID = rs.GetFieldInt("sell_money_id_2")
	GT_Hero_List[HeroID].DecomposePrice = rs.GetFieldInt("sell_money_num_2")
	GT_Hero_List[HeroID].CanSell = rs.GetFieldInt("sell") == 1
	GT_Hero_List[HeroID].PieceNum = rs.GetFieldInt("piece_num")
	GT_Hero_List[HeroID].ModelID = rs.GetFieldInt("modelid")
	GT_Hero_List[HeroID].Setup = rs.GetFieldInt("setup")
	GT_Hero_List[HeroID].PieceID = rs.GetFieldInt("chip_id")
	GT_Hero_List[HeroID].HeroExp = rs.GetFieldInt("experience")
	GT_Hero_List[HeroID].BasePropertys[0] = rs.GetFieldInt("p1")
	GT_Hero_List[HeroID].BasePropertys[1] = rs.GetFieldInt("p2")
	GT_Hero_List[HeroID].BasePropertys[2] = rs.GetFieldInt("p3")
	GT_Hero_List[HeroID].BasePropertys[3] = rs.GetFieldInt("p4")
	GT_Hero_List[HeroID].BasePropertys[4] = rs.GetFieldInt("p5")
	GT_Hero_List[HeroID].Skills[0] = rs.GetFieldInt("skill_1")
	GT_Hero_List[HeroID].Skills[1] = rs.GetFieldInt("skill_2")
	GT_Hero_List[HeroID].Skills[2] = rs.GetFieldInt("skill_3")
	GT_Hero_List[HeroID].Skills[3] = rs.GetFieldInt("skill_4")
	GT_Hero_List[HeroID].Skills[4] = rs.GetFieldInt("skill_5")

	if GT_Hero_List[HeroID].HeroExp <= 0 {
		panic("field experience should not be zero!")
	}

	if GT_Hero_List[HeroID].CanSell {
		if (GT_Hero_List[HeroID].SellID == 0) || GT_Hero_List[HeroID].SellPrice == 0 {
			panic("field sell_money_id_1 and sell_money_num_1 should not be zero!")
		}
	}

	return
}

func GetHeroInfo(heroid int) *ST_HeroInfo {
	if heroid >= len(GT_Hero_List) || heroid <= 0 {
		gamelog.Error("GetHeroInfo Error: invalid heroid :%d", heroid)
		return nil
	}

	if GT_Hero_List[heroid].HeroID != heroid {
		gamelog.Error("GetHeroInfo Error: invalid heroid2 :%d", heroid)
		return nil
	}

	return &GT_Hero_List[heroid]
}

func GetCampHeroCount() []int {
	campHeroCount := []int{0, 0, 0, 0}
	for _, v := range GT_Hero_List {
		if v.Camp == 1 {
			campHeroCount[0]++
		} else if v.Camp == 2 {
			campHeroCount[1]++
		} else if v.Camp == 3 {
			campHeroCount[2]++
		} else if v.Camp == 4 {
			campHeroCount[3]++
		}
	}

	return campHeroCount
}

func GetHeroQuality(heroid int) int8 {
	if heroid >= len(GT_Hero_List) || heroid <= 0 {
		gamelog.Error("GetHeroInfo Error: invalid heroid :%d", heroid)
		return 0
	}

	return GT_Hero_List[heroid].Quality
}

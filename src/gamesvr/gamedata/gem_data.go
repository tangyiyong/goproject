package gamedata

import (
	"gamelog"
)

type ST_GemInfo struct {
	GemID             int    //宝物ID
	ItemID            int    //对应的道具ID
	Quality           int    //装备的品质
	SellID            int    //出售货币ID
	SellPrice         int    //出售价格
	CanSell           bool   //是否可出售
	Position          int    //装备的位置
	Setup             int    //是否可以上阵
	Experience        int    //宝物带有的经验
	PieceIDs          []int  //宝物碎片列表
	BasePropertys     [2]int //基础属性值
	StrengthPropertys [2]int //强化的两个属性
	RefinePropertys   [2]int //精炼的两个属性
}

var (
	GT_GemList []ST_GemInfo = nil
)

func InitGemParser(total int) bool {

	GT_GemList = make([]ST_GemInfo, total+1)

	return true
}

func ParseGemRecord(rs *RecordSet) {
	GemID := rs.GetFieldInt("id")
	if (GemID <= 0) || (GemID >= len(GT_GemList)) {
		gamelog.Error("ParseGemRecord Error: Invalid GemID :%d", GemID)
		return
	}

	GT_GemList[GemID].GemID = GemID
	GT_GemList[GemID].Quality = rs.GetFieldInt("quality")
	GT_GemList[GemID].SellID = rs.GetFieldInt("sell_money_id")
	GT_GemList[GemID].SellPrice = rs.GetFieldInt("sell_money_num")
	GT_GemList[GemID].CanSell = rs.GetFieldInt("sell") == 1
	GT_GemList[GemID].Position = rs.GetFieldInt("pos")
	GT_GemList[GemID].Setup = rs.GetFieldInt("setup")
	GT_GemList[GemID].Experience = rs.GetFieldInt("exp")
	GT_GemList[GemID].PieceIDs = ParseToIntSlice(rs.GetFieldString("chip_id"))
	GT_GemList[GemID].BasePropertys = ParseTo2IntSlice(rs.GetFieldString("base_propertys"))
	GT_GemList[GemID].StrengthPropertys = ParseTo2IntSlice(rs.GetFieldString("strenth_propertys"))
	GT_GemList[GemID].RefinePropertys = ParseTo2IntSlice(rs.GetFieldString("refine_propertys"))
	GT_GemList[GemID].ItemID = rs.GetFieldInt("itemid")

	if (GT_GemList[GemID].SellID == 0) || GT_GemList[GemID].SellPrice == 0 {
		panic("field sell_money_id and sell_money_num should not be zero!")
	}

	if (GT_GemList[GemID].Position < 5) || (GT_GemList[GemID].Position > 6) {
		panic("field pos is wrong value!")
	}

	if (GT_GemList[GemID].StrengthPropertys[0] < 1) || GT_GemList[GemID].StrengthPropertys[0] > 11 {
		if GT_GemList[GemID].StrengthPropertys[0] != 20 && (GT_GemList[GemID].StrengthPropertys[0] != 21) {
			panic("field strenth_propertys is wrong value!" + string(GT_GemList[GemID].StrengthPropertys[0]))
		}
	}

	if GT_GemList[GemID].StrengthPropertys[1] < 1 || GT_GemList[GemID].StrengthPropertys[1] > 11 {
		if GT_GemList[GemID].StrengthPropertys[1] != 20 && GT_GemList[GemID].StrengthPropertys[1] != 21 {
			panic("field strenth_propertys is wrong value!" + string(GT_GemList[GemID].StrengthPropertys[1]))
		}
	}

	if GT_GemList[GemID].RefinePropertys[0] < 1 || GT_GemList[GemID].RefinePropertys[0] > 11 {
		if GT_GemList[GemID].RefinePropertys[0] != 20 && (GT_GemList[GemID].RefinePropertys[0] != 21) {
			panic("field refine_propertys is wrong value!" + string(GT_GemList[GemID].RefinePropertys[0]))
		}
	}

	if GT_GemList[GemID].RefinePropertys[1] < 1 || GT_GemList[GemID].RefinePropertys[1] > 11 {
		if GT_GemList[GemID].RefinePropertys[1] != 20 && (GT_GemList[GemID].RefinePropertys[1] != 21) {
			panic("field refine_propertys is wrong value!" + string(GT_GemList[GemID].RefinePropertys[1]))
		}
	}

	return
}

func GetGemInfo(gemid int) *ST_GemInfo {
	if gemid >= len(GT_GemList) || gemid <= 0 {
		gamelog.Error("GetGemInfo Error : equipid gemid %d", gemid)
		return nil
	}

	return &GT_GemList[gemid]
}

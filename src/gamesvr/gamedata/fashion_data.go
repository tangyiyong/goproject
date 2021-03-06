package gamedata

//时装配制表

import (
	"gamelog"
)

type ST_FashionInfo struct {
	ID       int32 //时装ID
	Quality  int   //品质
	PieceID  int   //装备碎片ID
	PieceNum int   //合成碎片数
}

type ST_FashionMapInfo struct {
	ID         int //宠物图鉴ID
	FashionIds []int32
	Buffs      [3]ST_PropertyBuff //宠物的Buff集
}

type ST_FashionLevel struct {
	ID             int    //时装ID
	Level          int    //等级
	CostItemID     int    //道具ID
	CostItemNum    int    //道具数量
	CostItemSum    int    //道具数量
	CostMoneyID    int    //货币ID
	CostMoneyNum   int    //货币数量
	CostMoneySum   int    //货币总量
	PropertyValues [5]int //五个属性加成
	ExtraID        int    //强化额外属性ID
	ExtraPro       int    //强化额外属性百分比
}

var (
	GT_Fashion_List         []ST_FashionInfo    = nil //时装基础信息
	GT_FashionMap_List      []ST_FashionMapInfo = nil //时装图鉴信息
	GT_FashionStrength_List [10][200]ST_FashionLevel
)

func InitFashionParser(total int) bool {
	GT_Fashion_List = make([]ST_FashionInfo, total+1)
	return true
}

func ParseFashionRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_Fashion_List[id].ID = int32(id)
	GT_Fashion_List[id].Quality = rs.GetFieldInt("quality")
	GT_Fashion_List[id].PieceID = rs.GetFieldInt("chip_id")
	GT_Fashion_List[id].PieceNum = rs.GetFieldInt("chip_num")
}

func GetFashionInfo(id int32) *ST_FashionInfo {
	if int(id) >= len(GT_Fashion_List) || id == 0 {
		gamelog.Error("GetFashionInfo Error: invalid id :%d", id)
		return nil
	}

	return &GT_Fashion_List[id]
}

func InitFashionMapParser(total int) bool {
	GT_FashionMap_List = make([]ST_FashionMapInfo, total+1)
	return true
}

func ParseFashionMapRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_FashionMap_List[id].ID = id

	fashionid := rs.GetFieldInt("fashion_id1")
	GT_FashionMap_List[id].FashionIds = append(GT_FashionMap_List[id].FashionIds, int32(fashionid))

	fashionid = rs.GetFieldInt("fashion_id2")
	if fashionid > 0 {
		GT_FashionMap_List[id].FashionIds = append(GT_FashionMap_List[id].FashionIds, int32(fashionid))
	}

	fashionid = rs.GetFieldInt("fashion_id3")
	if fashionid > 0 {
		GT_FashionMap_List[id].FashionIds = append(GT_FashionMap_List[id].FashionIds, int32(fashionid))
	}

	GT_FashionMap_List[id].Buffs[0].PropertyID = rs.GetFieldInt("property1")
	GT_FashionMap_List[id].Buffs[0].Value = rs.GetFieldInt("value1")
	GT_FashionMap_List[id].Buffs[0].IsPercent = rs.GetFieldInt("is_percent1") == 1
	GT_FashionMap_List[id].Buffs[1].PropertyID = rs.GetFieldInt("property2")
	GT_FashionMap_List[id].Buffs[1].Value = rs.GetFieldInt("value2")
	GT_FashionMap_List[id].Buffs[1].IsPercent = rs.GetFieldInt("is_percent2") == 1
	GT_FashionMap_List[id].Buffs[2].PropertyID = rs.GetFieldInt("property3")
	GT_FashionMap_List[id].Buffs[2].Value = rs.GetFieldInt("value3")
	GT_FashionMap_List[id].Buffs[2].IsPercent = rs.GetFieldInt("is_percent3") == 1
}

func InitFashionStrengthParser(total int) bool {
	return true
}

func ParseFashionStrengthRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	level := rs.GetFieldInt("level")
	GT_FashionStrength_List[id][level].ID = id
	GT_FashionStrength_List[id][level].Level = level
	GT_FashionStrength_List[id][level].CostItemID = rs.GetFieldInt("cost_item_id")
	GT_FashionStrength_List[id][level].CostItemNum = rs.GetFieldInt("cost_item_num")
	GT_FashionStrength_List[id][level].CostItemSum = rs.GetFieldInt("cost_item_sum")
	GT_FashionStrength_List[id][level].CostMoneyID = rs.GetFieldInt("cost_money_id")
	GT_FashionStrength_List[id][level].CostMoneyNum = rs.GetFieldInt("cost_money_num")
	GT_FashionStrength_List[id][level].CostMoneySum = rs.GetFieldInt("cost_money_sum")
	GT_FashionStrength_List[id][level].PropertyValues[0] = rs.GetFieldInt("value_1")
	GT_FashionStrength_List[id][level].PropertyValues[1] = rs.GetFieldInt("value_2")
	GT_FashionStrength_List[id][level].PropertyValues[2] = rs.GetFieldInt("value_3")
	GT_FashionStrength_List[id][level].PropertyValues[3] = rs.GetFieldInt("value_4")
	GT_FashionStrength_List[id][level].PropertyValues[4] = rs.GetFieldInt("value_5")
	GT_FashionStrength_List[id][level].ExtraID = rs.GetFieldInt("extra_id")
	GT_FashionStrength_List[id][level].ExtraPro = rs.GetFieldInt("extra_value")

}

func GetFashionLevelInfo(id int32, level int32) *ST_FashionLevel {
	return &GT_FashionStrength_List[id][level]
}

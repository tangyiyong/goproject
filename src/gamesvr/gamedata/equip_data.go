package gamedata

import (
	"gamelog"
)

type ST_EquipInfo struct {
	EquipID          int    //装备ID
	ItemID           int    //对应的道具ID
	Quality          int    //装备的品质
	Position         int    //装备的位置
	SuitID           int    //套装ID
	PieceID          int    //装备碎片ID
	PieceNum         int    //合成碎片数
	SellID           [2]int //出售货币ID
	SellPrice        [2]int //出售价格
	CanSell          bool   //是否可出售
	BaseProperty     int    //基础属性值
	StrengthProperty int    //强化的1个属性
	RefinePropertys  [2]int //精炼的两个属性
}

var (
	GT_EquipList []ST_EquipInfo = nil
)

func InitEquipParser(total int) bool {
	GT_EquipList = make([]ST_EquipInfo, total+1)
	return true
}

func ParseEquipRecord(rs *RecordSet) {
	EquipID := CheckAtoi(rs.Values[0], 0)
	if (EquipID <= 0) || (EquipID >= len(GT_EquipList)) {
		gamelog.Error("ParseEquipRecord Error: Invalid EquipID :%d", EquipID)
		return
	}

	GT_EquipList[EquipID].EquipID = EquipID
	GT_EquipList[EquipID].ItemID = rs.GetFieldInt("itemid")
	GT_EquipList[EquipID].Quality = rs.GetFieldInt("quality")
	GT_EquipList[EquipID].Position = rs.GetFieldInt("pos")
	GT_EquipList[EquipID].SuitID = rs.GetFieldInt("suit_id")
	GT_EquipList[EquipID].PieceNum = rs.GetFieldInt("piece_num")
	GT_EquipList[EquipID].SellID[0] = rs.GetFieldInt("sell_money_id_1")
	GT_EquipList[EquipID].SellPrice[0] = rs.GetFieldInt("sell_money_num_1")
	GT_EquipList[EquipID].SellID[1] = rs.GetFieldInt("sell_money_id_2")
	GT_EquipList[EquipID].SellPrice[1] = rs.GetFieldInt("sell_money_num_2")
	GT_EquipList[EquipID].CanSell = rs.GetFieldInt("sell") == 1
	GT_EquipList[EquipID].PieceID = rs.GetFieldInt("chip_id")
	GT_EquipList[EquipID].BaseProperty = rs.GetFieldInt("base_propertys")
	GT_EquipList[EquipID].StrengthProperty = rs.GetFieldInt("strenth_propertys")
	GT_EquipList[EquipID].RefinePropertys = ParseTo2IntSlice(rs.GetFieldString("refine_propertys"))

	if (GT_EquipList[EquipID].Position < 1) || (GT_EquipList[EquipID].Position > 4) {
		panic("field pos is wrong value!")
	}

	if GT_EquipList[EquipID].StrengthProperty <= 0 {
		panic("field strenth_propertys is wrong value!")
	}

	return
}

func GetEquipmentInfo(equipid int) *ST_EquipInfo {
	if equipid >= len(GT_EquipList) || equipid <= 0 {
		gamelog.Error("GetEquipmentInfo Error : equipid itemid %d", equipid)
		return nil
	}

	return &GT_EquipList[equipid]
}

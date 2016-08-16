package gamedata

//宠物配制表
//配制了所有物类型和属性

import (
	"gamelog"
)

type ST_PetInfo struct {
	PetID     int    //宠物ID
	ItemID    int    //对应的道具ID
	Quality   int    //品质
	CanSell   bool   //是否可以售卖
	SellID    int    //出售货币ID
	SellPrice int    //出售价格
	PieceID   int    //碎片道具ID
	PieceNum  int    //合成碎片数
	BasePtys  [5]int //基础属性
}

var GT_Pet_List []ST_PetInfo = nil

func InitPetParser(total int) bool {
	GT_Pet_List = make([]ST_PetInfo, total+1)
	return true
}

func ParsePetRecord(rs *RecordSet) {
	petID := rs.GetFieldInt("id")
	if (petID <= 0) || (petID >= len(GT_Pet_List)) {
		gamelog.Error("ParsePetRecord Error: Invalid petID :%d", petID)
		return
	}

	GT_Pet_List[petID].PetID = petID
	GT_Pet_List[petID].ItemID = rs.GetFieldInt("itemid")
	GT_Pet_List[petID].Quality = rs.GetFieldInt("quality")
	GT_Pet_List[petID].PieceNum = rs.GetFieldInt("piece_num")
	GT_Pet_List[petID].PieceID = rs.GetFieldInt("chip_id")
	GT_Pet_List[petID].BasePtys[0] = rs.GetFieldInt("p1")
	GT_Pet_List[petID].BasePtys[1] = rs.GetFieldInt("p2")
	GT_Pet_List[petID].BasePtys[2] = rs.GetFieldInt("p3")
	GT_Pet_List[petID].BasePtys[3] = rs.GetFieldInt("p4")
	GT_Pet_List[petID].BasePtys[4] = rs.GetFieldInt("p5")
	GT_Pet_List[petID].CanSell = rs.GetFieldInt("sell") == 1
	GT_Pet_List[petID].SellID = rs.GetFieldInt("sell_money_id_1")
	GT_Pet_List[petID].SellPrice = rs.GetFieldInt("sell_money_num_1")

}

func GetPetInfo(petID int) *ST_PetInfo {
	if petID >= len(GT_Pet_List) || petID == 0 {
		gamelog.Error("GetPetInfo Error: invalid petID :%d", petID)
		return nil
	}

	return &GT_Pet_List[petID]
}

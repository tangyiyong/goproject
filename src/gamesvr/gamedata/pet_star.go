package gamedata

//宠物配制表
//配制了所有宠物升星

import (
	"gamelog"
)

type ST_PetStarInfo struct {
	Quality       int //宠物品质
	Star          int //宠物星级
	NeedLevel     int //升级所需要的等级
	MoneyID       int //需要的货币ID
	MoneyNum      int //需要的货币系数*经验值
	NeedItemID    int
	NeedItemNum   int
	PieceNum      int
	Inc_percent   int
	PropertyTrans int
	TotalItemNum  int
	TotalPiece    int
	TotalMoney    int
}

var (
	GT_PetStar_List [8][6]ST_PetStarInfo // = nil
)

func InitPetStarParser(total int) bool {
	//GT_PetStar_List = make([]ST_PetStarInfo, total+1)
	return true
}

func ParsePetStarRecord(rs *RecordSet) {
	quality := rs.GetFieldInt("quality")
	star := rs.GetFieldInt("star")

	GT_PetStar_List[quality][star].Quality = quality
	GT_PetStar_List[quality][star].Star = star
	GT_PetStar_List[quality][star].NeedLevel = rs.GetFieldInt("need_level")
	GT_PetStar_List[quality][star].MoneyID = rs.GetFieldInt("money_id")
	GT_PetStar_List[quality][star].MoneyNum = rs.GetFieldInt("money_num")
	GT_PetStar_List[quality][star].NeedItemID = rs.GetFieldInt("need_item_id")
	GT_PetStar_List[quality][star].NeedItemNum = rs.GetFieldInt("need_item_num")
	GT_PetStar_List[quality][star].PieceNum = rs.GetFieldInt("piece_num")
	GT_PetStar_List[quality][star].PropertyTrans = rs.GetFieldInt("property_trans")
	GT_PetStar_List[quality][star].TotalItemNum = rs.GetFieldInt("total_item_num")
	GT_PetStar_List[quality][star].TotalPiece = rs.GetFieldInt("total_piece_num")
	GT_PetStar_List[quality][star].TotalMoney = rs.GetFieldInt("total_money")
	GT_PetStar_List[quality][star].Inc_percent = rs.GetFieldInt("inc_percent")
}

func GetPetStarInfo(quality int, star int) *ST_PetStarInfo {
	if quality <= 0 || quality >= 8 {
		gamelog.Error("GetPetStarInfo Error: Invalid Quality :%d", quality)
		return nil
	}

	if star < 0 || star >= 6 {
		gamelog.Error("GetPetStarInfo Error: Invalid star :%d", star)
		return nil
	}

	return &GT_PetStar_List[quality][star]
}

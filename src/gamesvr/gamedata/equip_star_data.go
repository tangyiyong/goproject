package gamedata

import (
	"gamelog"
)

type ST_EquipStarInfo struct {
	Quality       int    //装备品质
	Position      int    //装备位置
	PropertyID    int    //属性ID
	TotalProperty int    //星级总属性
	ExtraProperty int    //星级额外属性
	AddProperty   int    //一次增加属性
	NeedExp       int    //升级需要经验
	AddExp        int    //增加经验
	AddLuck       int    //增加幸运值
	StarLvl       int    //星级
	MoneyID       [2]int //货币
	MoneyNum      [2]int //货币数量
	PieceNum      int    //碎片数量
	Luck          [4]int //幸运值
	Ratio         [4]int //幸运率
	BaseProperty  int    //当前所有星级属性之和
}

var (
	GT_EquipStarList [10][5][6]ST_EquipStarInfo // = nil
	MinStarQuality   int
)

func InitEquipStarParser(total int) bool {
	MinStarQuality = 100
	return true
}

func ParseEquipStarRecord(rs *RecordSet) {
	quality := rs.GetFieldInt("quality")
	pos := rs.GetFieldInt("position")
	starlvl := rs.GetFieldInt("star_lvl")

	if quality < MinStarQuality {
		MinStarQuality = quality
	}

	GT_EquipStarList[quality][pos][starlvl].Quality = quality
	GT_EquipStarList[quality][pos][starlvl].Position = pos
	GT_EquipStarList[quality][pos][starlvl].StarLvl = starlvl
	GT_EquipStarList[quality][pos][starlvl].BaseProperty = rs.GetFieldInt("base_property")
	GT_EquipStarList[quality][pos][starlvl].PropertyID = rs.GetFieldInt("propertyid")
	GT_EquipStarList[quality][pos][starlvl].TotalProperty = rs.GetFieldInt("total_property")
	GT_EquipStarList[quality][pos][starlvl].ExtraProperty = rs.GetFieldInt("extra_property")
	GT_EquipStarList[quality][pos][starlvl].AddProperty = rs.GetFieldInt("add_property")
	GT_EquipStarList[quality][pos][starlvl].NeedExp = rs.GetFieldInt("need_exp")
	GT_EquipStarList[quality][pos][starlvl].AddLuck = rs.GetFieldInt("add_luck")
	GT_EquipStarList[quality][pos][starlvl].AddExp = rs.GetFieldInt("add_exp")
	GT_EquipStarList[quality][pos][starlvl].MoneyID[0] = rs.GetFieldInt("money_id_1")
	GT_EquipStarList[quality][pos][starlvl].MoneyID[1] = rs.GetFieldInt("money_id_2")
	GT_EquipStarList[quality][pos][starlvl].MoneyNum[0] = rs.GetFieldInt("money_num_1")
	GT_EquipStarList[quality][pos][starlvl].MoneyNum[1] = rs.GetFieldInt("money_num_2")
	GT_EquipStarList[quality][pos][starlvl].PieceNum = rs.GetFieldInt("chip_num")

	GT_EquipStarList[quality][pos][starlvl].Luck[0] = rs.GetFieldInt("luck_num1")
	GT_EquipStarList[quality][pos][starlvl].Luck[1] = rs.GetFieldInt("luck_num2")
	GT_EquipStarList[quality][pos][starlvl].Luck[2] = rs.GetFieldInt("luck_num3")
	GT_EquipStarList[quality][pos][starlvl].Luck[3] = rs.GetFieldInt("luck_num4")
	GT_EquipStarList[quality][pos][starlvl].Ratio[0] = rs.GetFieldInt("success_rate1")
	GT_EquipStarList[quality][pos][starlvl].Ratio[1] = rs.GetFieldInt("success_rate2")
	GT_EquipStarList[quality][pos][starlvl].Ratio[2] = rs.GetFieldInt("success_rate3")
	GT_EquipStarList[quality][pos][starlvl].Ratio[3] = rs.GetFieldInt("success_rate4")

	return
}

func GetEquipStarInfo(quality int, pos int, star int) *ST_EquipStarInfo {
	if quality < MinStarQuality || quality > 10 {
		gamelog.Error("GetEquipStarInfo  Error: Invalid quality :%d, minQuality:%d", quality, MinStarQuality)
		return nil
	}
	if pos <= 0 || pos > 4 {
		gamelog.Error("GetEquipStarInfo  Error: Invalid Pos :%d", pos)
		return nil
	}

	if star < 0 || star > 5 {
		gamelog.Error("GetEquipStarInfo  Error: Invalid star :%d", star)
		return nil
	}

	return &GT_EquipStarList[quality][pos][star]
}

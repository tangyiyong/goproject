package gamedata

//宠物配制表
//配制了所有宠物的神炼属性

import (
	"gamelog"
)

type Pet_God struct {
	Level     int //宠物等级
	NeedExp   int //升级需要经验
	TotalExp  int
	Propertys [3]ST_PropertyBuff //属性值
}

type ST_PetGodInfo struct {
	ID     int         //宠物ID
	Levels [31]Pet_God //等级属性
}

var (
	GT_PetGod_List   []ST_PetGodInfo = nil
	GT_Max_God_Level                 = 0
)

func InitPetGodParser(total int) bool {
	GT_PetGod_List = make([]ST_PetGodInfo, total+1)
	return true
}

func ParsePetGodRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_PetGod_List[id].ID = id

	godlvl := rs.GetFieldInt("god_lvl")
	GT_PetGod_List[id].Levels[godlvl].Level = godlvl
	GT_PetGod_List[id].Levels[godlvl].NeedExp = rs.GetFieldInt("god_exp")
	GT_PetGod_List[id].Levels[godlvl].Propertys[0].PropertyID = rs.GetFieldInt("property1")
	GT_PetGod_List[id].Levels[godlvl].Propertys[0].Value = rs.GetFieldInt("value1")
	GT_PetGod_List[id].Levels[godlvl].Propertys[0].IsPercent = rs.GetFieldInt("is_percent1") == 1
	GT_PetGod_List[id].Levels[godlvl].Propertys[1].PropertyID = rs.GetFieldInt("property2")
	GT_PetGod_List[id].Levels[godlvl].Propertys[1].Value = rs.GetFieldInt("value2")
	GT_PetGod_List[id].Levels[godlvl].Propertys[1].IsPercent = rs.GetFieldInt("is_percent2") == 1
	GT_PetGod_List[id].Levels[godlvl].Propertys[2].PropertyID = rs.GetFieldInt("property3")
	GT_PetGod_List[id].Levels[godlvl].Propertys[2].Value = rs.GetFieldInt("value3")
	GT_PetGod_List[id].Levels[godlvl].Propertys[2].IsPercent = rs.GetFieldInt("is_percent3") == 1
	GT_PetGod_List[id].Levels[godlvl].TotalExp = rs.GetFieldInt("total_exp")

	if godlvl > GT_Max_God_Level {
		GT_Max_God_Level = godlvl
	}
}

func GetPetGodInfo(petID int, level int) *Pet_God {
	if petID >= len(GT_PetGod_List) || petID <= 0 {
		gamelog.Error("GetPetInfo Error: invalid petID :%d", petID)
		return nil
	}

	if level < 0 || level >= len(GT_PetGod_List[petID].Levels) {
		gamelog.Error("GetPetInfo Error: invalid level :%d", level)
		return nil
	}

	return &GT_PetGod_List[petID].Levels[level]
}

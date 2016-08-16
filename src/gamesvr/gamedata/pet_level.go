package gamedata

//宠物配制表
//配制了所有宠物等级

import (
	"gamelog"
)

type Pet_Level struct {
	Level      int    //宠物等级
	NeedExp    int    //升级所需要的经验
	MoneyID    int    //需要的货币ID
	MoneyNum   int    //需要的货币系数*经验值
	Propertys  [5]int //属性值
	TotalExp   int
	TotalMoney int
}

type ST_PetLevelInfo struct {
	ID     int            //宠物ID
	Levels [101]Pet_Level //等级属性
}

var (
	GT_PetLevel_List []ST_PetLevelInfo = nil
)

func InitPetLevelParser(total int) bool {
	GT_PetLevel_List = make([]ST_PetLevelInfo, total+1)
	return true
}

func ParsePetLevelRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	GT_PetLevel_List[id].ID = id

	Level := rs.GetFieldInt("level")
	GT_PetLevel_List[id].Levels[Level].Level = Level
	GT_PetLevel_List[id].Levels[Level].NeedExp = rs.GetFieldInt("needexp")
	GT_PetLevel_List[id].Levels[Level].MoneyID = rs.GetFieldInt("money_id")
	GT_PetLevel_List[id].Levels[Level].MoneyNum = rs.GetFieldInt("money_num")
	GT_PetLevel_List[id].Levels[Level].Propertys[0] = rs.GetFieldInt("p1")
	GT_PetLevel_List[id].Levels[Level].Propertys[1] = rs.GetFieldInt("p2")
	GT_PetLevel_List[id].Levels[Level].Propertys[2] = rs.GetFieldInt("p3")
	GT_PetLevel_List[id].Levels[Level].Propertys[3] = rs.GetFieldInt("p4")
	GT_PetLevel_List[id].Levels[Level].Propertys[4] = rs.GetFieldInt("p5")
	GT_PetLevel_List[id].Levels[Level].TotalMoney = rs.GetFieldInt("total_money")
	GT_PetLevel_List[id].Levels[Level].TotalExp = rs.GetFieldInt("total_exp")
}

func GetPetLevelInfo(id int, level int) *Pet_Level {
	if id >= len(GT_PetLevel_List) || id <= 0 {
		gamelog.Error("GetPetLevelInfo Error : Invalid id :%d", id)
		return nil
	}

	if level >= len(GT_PetLevel_List[id].Levels) || level <= 0 {
		gamelog.Error("GetPetLevelInfo Error : Invalid level :%d", level)
		return nil
	}

	return &GT_PetLevel_List[id].Levels[level]
}

package gamedata

//行动力配制表
//配制了所有行动类型和属性

import (
	"fmt"
	"gamelog"
)

type ST_GuildSkill struct {
	ID         int
	PropertyID int
}

var GT_GuildSkillID_Map map[int]int

func InitGuildSkillParser(total int) bool {
	GT_GuildSkillID_Map = make(map[int]int)
	return true
}

func ParseGuildSkillRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	propertyID := rs.GetFieldInt("propertyid")
	GT_GuildSkillID_Map[id] = propertyID
}

func GetGuildSkillPropertyID(id int) int {
	id, ok := GT_GuildSkillID_Map[id]
	if ok == false {
		gamelog.Error("GetGuildSkillPropertyID Fail, Invalid id %d", id)
		return -1
	}

	return id
}

type ST_GuildSkillLimit struct {
	Level       int
	PropertyLst [9]int
}

var GT_GuildSkillLimit_List []ST_GuildSkillLimit = nil

func InitGuildSkillLimitParser(total int) bool {
	GT_GuildSkillLimit_List = make([]ST_GuildSkillLimit, total+5)
	return true
}

func ParseGuildSkillLimitRecord(rs *RecordSet) {
	level := rs.GetFieldInt("level")
	GT_GuildSkillLimit_List[level].Level = level

	for i := 1; i <= 9; i++ {
		filedName := fmt.Sprintf("skill_%d", i)
		GT_GuildSkillLimit_List[level].PropertyLst[i-1] = rs.GetFieldInt(filedName)
	}
}

func GetGuildSkillLimit(level int, id int) int {
	if level > len(GT_GuildSkillLimit_List)-1 {
		gamelog.Error("GetGuildSkillLimit Error: invalid level :%d", level)
		return 0
	}

	return GT_GuildSkillLimit_List[level].PropertyLst[id-1]
}

type ST_GuildSkillMax struct {
	Level   int
	Skill   [9]int
	NeedNum [9]int
	NeedExp [9]int
}

var GT_GuildSkillMax_List []ST_GuildSkillMax

func InitGuildSkillMaxParser(total int) bool {
	GT_GuildSkillMax_List = make([]ST_GuildSkillMax, total+1)
	return true
}

func ParseGuildSkillMaxRecord(rs *RecordSet) {
	level := CheckAtoi(rs.Values[0], 0)
	GT_GuildSkillMax_List[level].Level = level

	for i := 1; i <= 9; i++ {
		filedName := fmt.Sprintf("skill_%d", i)
		GT_GuildSkillMax_List[level].Skill[i-1] = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("neednum_%d", i)
		GT_GuildSkillMax_List[level].NeedNum[i-1] = rs.GetFieldInt(filedName)

		filedName = fmt.Sprintf("needexp_%d", i)
		GT_GuildSkillMax_List[level].NeedExp[i-1] = rs.GetFieldInt(filedName)
	}
}

func GetGuildSkillValue(level int, id int) (value int) {
	value = GT_GuildSkillMax_List[level].Skill[id-1]

	return value
}

func GetGuildSkillExpValue(level int8) int {
	if int(level) >= len(GT_GuildSkillMax_List) {
		gamelog.Error("GetGuildSkillExpValue Error: invalid level :%d", level)
		return 0
	}

	return GT_GuildSkillMax_List[level].Skill[8]
}

func GetGuildSkillNeedExp(level int, id int) int {
	if level >= len(GT_GuildSkillMax_List) {
		gamelog.Error("GetGuildSkillNeedExp Error: invalid level :%d", level)
		return 0
	}

	return GT_GuildSkillMax_List[level].NeedExp[id-1]
}

func GetGuildSkillNeedMoney(level int, id int) (int, int) {
	if level >= len(GT_GuildSkillMax_List) {
		gamelog.Error("GetGuildSkillNeedMoney Error: invalid level :%d", level)
		return 0, 0
	}

	return GuildSKillStudyNeedMoneyID, GT_GuildSkillMax_List[level].NeedNum[id-1]
}

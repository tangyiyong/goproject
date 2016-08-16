package gamedata

//行动力配制表
//配制了所有行动类型和属性

import (
	"gamelog"
)

type ST_GuildSkillInfo struct {
	Level     int     //技能等级
	Property  [11]int //属性值
	NeedExp   [11]int //需求经验
	NeedMoney [11]int //需求货币

	ExpInc     int //经验增加值
	NeedExpNum int
	NeedExpExp int
}

var GT_GuildSkill_List []ST_GuildSkillInfo = nil

func InitGuildSkillParser(total int) bool {
	GT_GuildSkill_List = make([]ST_GuildSkillInfo, total+1)
	return true
}

func ParseGuildSkillRecord(rs *RecordSet) {
	level := rs.GetFieldInt("level")
	GT_GuildSkill_List[level].Level = level
	GT_GuildSkill_List[level].Property[0] = rs.GetFieldInt("property_1")
	GT_GuildSkill_List[level].NeedExp[0] = rs.GetFieldInt("needexp_1")
	GT_GuildSkill_List[level].NeedMoney[0] = rs.GetFieldInt("neednum_1")

	GT_GuildSkill_List[level].Property[1] = rs.GetFieldInt("property_20")
	GT_GuildSkill_List[level].NeedExp[1] = rs.GetFieldInt("needexp_20")
	GT_GuildSkill_List[level].NeedMoney[1] = rs.GetFieldInt("neednum_20")

	GT_GuildSkill_List[level].Property[2] = rs.GetFieldInt("property_3")
	GT_GuildSkill_List[level].NeedExp[2] = rs.GetFieldInt("needexp_3")
	GT_GuildSkill_List[level].NeedMoney[2] = rs.GetFieldInt("neednum_3")

	GT_GuildSkill_List[level].Property[3] = rs.GetFieldInt("property_20")
	GT_GuildSkill_List[level].NeedExp[3] = rs.GetFieldInt("needexp_20")
	GT_GuildSkill_List[level].NeedMoney[3] = rs.GetFieldInt("neednum_20")

	GT_GuildSkill_List[level].Property[4] = rs.GetFieldInt("property_5")
	GT_GuildSkill_List[level].NeedExp[4] = rs.GetFieldInt("needexp_5")
	GT_GuildSkill_List[level].NeedMoney[4] = rs.GetFieldInt("neednum_5")

	GT_GuildSkill_List[level].Property[7] = rs.GetFieldInt("property_8")
	GT_GuildSkill_List[level].NeedExp[7] = rs.GetFieldInt("needexp_8")
	GT_GuildSkill_List[level].NeedMoney[7] = rs.GetFieldInt("neednum_8")

	GT_GuildSkill_List[level].Property[8] = rs.GetFieldInt("property_9")
	GT_GuildSkill_List[level].NeedExp[8] = rs.GetFieldInt("needexp_9")
	GT_GuildSkill_List[level].NeedMoney[8] = rs.GetFieldInt("neednum_9")

	GT_GuildSkill_List[level].Property[9] = rs.GetFieldInt("property_10")
	GT_GuildSkill_List[level].NeedExp[9] = rs.GetFieldInt("needexp_10")
	GT_GuildSkill_List[level].NeedMoney[9] = rs.GetFieldInt("neednum_10")

	GT_GuildSkill_List[level].Property[10] = rs.GetFieldInt("property_11")
	GT_GuildSkill_List[level].NeedExp[10] = rs.GetFieldInt("needexp_11")
	GT_GuildSkill_List[level].NeedMoney[10] = rs.GetFieldInt("neednum_11")

	GT_GuildSkill_List[level].ExpInc = rs.GetFieldInt("exp_inc")
	GT_GuildSkill_List[level].NeedExpNum = rs.GetFieldInt("neednum_exp")
	GT_GuildSkill_List[level].NeedExpExp = rs.GetFieldInt("needexp_exp")
}

func GetGuildSkillValue(level int, propertyid int) int {
	if propertyid == AttackPropertyID {
		propertyid = AttackPhysicID
	}

	return GT_GuildSkill_List[level].Property[propertyid]
}

func GetGuildSkillExpValue(level int) int {
	if level >= len(GT_GuildSkill_List) {
		gamelog.Error("GetGuildSkillExpValue Error: invalid level :%d", level)
		return 0
	}

	return GT_GuildSkill_List[level].ExpInc
}

func GetGuildSkillNeedExp(level int, propertyid int) int {
	if level >= len(GT_GuildSkill_List) {
		gamelog.Error("GetGuildSkillNeedExp Error: invalid level :%d", level)
		return 0
	}

	if propertyid == AttackPropertyID {
		propertyid = AttackPhysicID
	}

	return GT_GuildSkill_List[level].NeedExp[propertyid-1]
}

func GetGuildSkillExpNeedExp(level int) int {
	if level >= len(GT_GuildSkill_List) {
		gamelog.Error("GetGuildSkillExpNeedExp Error: invalid level :%d", level)
		return 0
	}

	return GT_GuildSkill_List[level].NeedExpExp
}

func GetGuildSkillExpNeedMoney(level int) (int, int) {
	if level >= len(GT_GuildSkill_List) {
		gamelog.Error("GetGuildSkillExpNeedMoney Error: invalid level :%d", level)
		return 0, 0
	}

	return GuildSKillStudyNeedMoneyID, GT_GuildSkill_List[level].NeedExpNum
}

func GetGuildSkillNeedMoney(level int, propertyid int) (int, int) {
	if level >= len(GT_GuildSkill_List) {
		gamelog.Error("GetGuildSkillNeedMoney Error: invalid level :%d", level)
		return 0, 0
	}

	if propertyid == AttackPropertyID {
		propertyid = AttackPhysicID
	}

	return GuildSKillStudyNeedMoneyID, GT_GuildSkill_List[level].NeedMoney[propertyid]
}

type ST_GuildSkillLimit struct {
	Level       int
	PropertyLst [11]int
	ExpLevel    int
}

var GT_GuildSkillLimit_List []ST_GuildSkillLimit = nil

func InitGuildSkillLimitParser(total int) bool {
	GT_GuildSkillLimit_List = make([]ST_GuildSkillLimit, total+1+4)
	return true
}

func ParseGuildSkillLimitRecord(rs *RecordSet) {
	level := rs.GetFieldInt("level")
	GT_GuildSkillLimit_List[level].Level = level
	GT_GuildSkillLimit_List[level].PropertyLst[0] = rs.GetFieldInt("property_1_level")
	GT_GuildSkillLimit_List[level].PropertyLst[1] = rs.GetFieldInt("property_20_level")
	GT_GuildSkillLimit_List[level].PropertyLst[3] = rs.GetFieldInt("property_20_level")
	GT_GuildSkillLimit_List[level].PropertyLst[2] = rs.GetFieldInt("property_3_level")
	GT_GuildSkillLimit_List[level].PropertyLst[4] = rs.GetFieldInt("property_5_level")
	GT_GuildSkillLimit_List[level].PropertyLst[7] = rs.GetFieldInt("property_8_level")
	GT_GuildSkillLimit_List[level].PropertyLst[8] = rs.GetFieldInt("property_9_level")
	GT_GuildSkillLimit_List[level].PropertyLst[9] = rs.GetFieldInt("property_10_level")
	GT_GuildSkillLimit_List[level].PropertyLst[10] = rs.GetFieldInt("property_11_level")
	GT_GuildSkillLimit_List[level].ExpLevel = rs.GetFieldInt("property_exp_level")
}

func GetGuildSkillLimit(level int, property int) int {
	if level > len(GT_GuildSkillLimit_List)-1 {
		gamelog.Error("GetGuildSkillLimit Error: invalid level :%d", level)
		return 0
	}

	if property == AttackPropertyID {
		property = AttackPhysicID
	}

	return GT_GuildSkillLimit_List[level].PropertyLst[property-1]
}

func GetGuildExpIncSKillLimit(level int) int {
	return GT_GuildSkillLimit_List[level].ExpLevel
}

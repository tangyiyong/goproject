package gamedata

import (
	"gamelog"
)

type ST_BreakTalent struct {
	HeroID  int     //英雄ID
	Talents [15]int //天赋集
}

var (
	GT_HeroBreakTalent_List []ST_BreakTalent = nil
)

func InitHeroBreakTalentParser(total int) bool {
	GT_HeroBreakTalent_List = make([]ST_BreakTalent, 1500)

	return true
}

func ParseHeroBreakTalentRecord(rs *RecordSet) {
	heroid := rs.GetFieldInt("hero_id")
	GT_HeroBreakTalent_List[heroid].HeroID = heroid
	GT_HeroBreakTalent_List[heroid].Talents[0] = CheckAtoi(rs.Values[2], 2)
	GT_HeroBreakTalent_List[heroid].Talents[1] = CheckAtoi(rs.Values[3], 3)
	GT_HeroBreakTalent_List[heroid].Talents[2] = CheckAtoi(rs.Values[4], 4)
	GT_HeroBreakTalent_List[heroid].Talents[3] = CheckAtoi(rs.Values[5], 5)
	GT_HeroBreakTalent_List[heroid].Talents[4] = CheckAtoi(rs.Values[6], 6)
	GT_HeroBreakTalent_List[heroid].Talents[5] = CheckAtoi(rs.Values[7], 7)
	GT_HeroBreakTalent_List[heroid].Talents[6] = CheckAtoi(rs.Values[8], 8)
	GT_HeroBreakTalent_List[heroid].Talents[7] = CheckAtoi(rs.Values[9], 9)
	GT_HeroBreakTalent_List[heroid].Talents[8] = CheckAtoi(rs.Values[10], 10)
	GT_HeroBreakTalent_List[heroid].Talents[9] = CheckAtoi(rs.Values[11], 11)
	GT_HeroBreakTalent_List[heroid].Talents[10] = CheckAtoi(rs.Values[12], 12)
	GT_HeroBreakTalent_List[heroid].Talents[11] = CheckAtoi(rs.Values[13], 13)
	GT_HeroBreakTalent_List[heroid].Talents[12] = CheckAtoi(rs.Values[14], 14)
	GT_HeroBreakTalent_List[heroid].Talents[13] = CheckAtoi(rs.Values[15], 15)
	GT_HeroBreakTalent_List[heroid].Talents[14] = CheckAtoi(rs.Values[16], 16)

	return
}

func GetHeroBreakTalentInfo(heroid int) *ST_BreakTalent {
	if heroid >= len(GT_HeroBreakTalent_List) || heroid <= 0 {
		gamelog.Error("GetHeroBreakTalentInfo Error : Invalid heroid :%d", heroid)
		return nil
	}

	if heroid != GT_HeroBreakTalent_List[heroid].HeroID {
		gamelog.Error("GetHeroBreakTalentInfo Error : Invalid heroid2 :%d", heroid)
		return nil
	}

	return &GT_HeroBreakTalent_List[heroid]
}

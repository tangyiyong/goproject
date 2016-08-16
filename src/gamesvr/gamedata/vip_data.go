package gamedata

import (
	"gamelog"
)

//! VIP信息表
type ST_VipLevelExpInfo struct {
	Level    int //! VIP等级
	Exp      int //! 等级需求VIP经验
	VipAward int //! VIP每日福利 值取Award表
}

var GT_Vip_List []ST_VipLevelExpInfo = nil

func InitVipParser(total int) bool {
	GT_Vip_List = make([]ST_VipLevelExpInfo, total+1)
	return true
}

func ParseVipRecord(rs *RecordSet) {
	Level := rs.GetFieldInt("level")
	GT_Vip_List[Level].Level = Level
	GT_Vip_List[Level].Exp = rs.GetFieldInt("exp")
	GT_Vip_List[Level].VipAward = rs.GetFieldInt("vipaward")
	return
}

func GetVipInfo(level int) *ST_VipLevelExpInfo {
	if level >= len(GT_Vip_List) || level < 0 {
		gamelog.Error("GetVipInfo Error: invalid level %d", level)
		return nil
	}

	return &GT_Vip_List[level]
}

func CalcVipLevelByExp(exp int, oldlevel int) int {
	if len(GT_Vip_List) == 0 {
		gamelog.Error("CalcVipLevelByExp fail. Vip list is nil")
		return oldlevel
	}

	pVipInfo := GetVipInfo(oldlevel)
	if pVipInfo == nil {
		gamelog.Error("CalcVipLevelByExp fail. Vip info is nil")
		return oldlevel
	}

	if pVipInfo.Exp > exp {
		return oldlevel
	}

	i := oldlevel
	for {
		i += 1
		pVipInfo = GetVipInfo(i)
		if pVipInfo == nil || pVipInfo.Exp == 0 || pVipInfo.Level == 0 {
			return i - 1
		}

		if pVipInfo.Exp > exp {
			return i
		}
	}

	gamelog.Error("CalcVipLevelByExp fail.return 0")
	return 0
}

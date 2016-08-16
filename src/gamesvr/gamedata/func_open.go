package gamedata

import (
	"gamelog"
)

type ST_FuncOpen struct {
	FuncID    int //! 功能ID
	OpenLevel int //! 开放等级
	VipLevel  int //! VIP提前开放等级
	Logic     int //! 逻辑关系 1->优先VIP 2->同时满足
}

var GT_FuncOpen_List []ST_FuncOpen = nil

func InitFuncOpenParser(total int) bool {
	GT_FuncOpen_List = make([]ST_FuncOpen, total+1)
	return true
}

func ParseFuncOpenRecord(rs *RecordSet) {
	funcID := CheckAtoi(rs.Values[0], 0)
	GT_FuncOpen_List[funcID].FuncID = funcID
	GT_FuncOpen_List[funcID].OpenLevel = rs.GetFieldInt("level")
	GT_FuncOpen_List[funcID].VipLevel = rs.GetFieldInt("viplevel")
	GT_FuncOpen_List[funcID].Logic = rs.GetFieldInt("logic")
}

func GetFuncOpenInfo(funcid int) *ST_FuncOpen {
	if funcid >= len(GT_FuncOpen_List) || funcid <= 0 {
		gamelog.Error("GetFuncOpenInfo Error: invalid funcid %d", funcid)
		return nil
	}

	return &GT_FuncOpen_List[funcid]
}

//! 检测是否满足功能开启条件
func IsFuncOpen(funcid int, level int, viplevel int) bool {
	pFuncOpen := GetFuncOpenInfo(funcid)
	if pFuncOpen == nil {
		gamelog.Error("IsFuncOpen Error: invalid funcid %d", funcid)
		return false
	}

	if pFuncOpen.Logic == 1 {
		if level >= pFuncOpen.OpenLevel || viplevel >= pFuncOpen.VipLevel {
			return true
		}
	} else if pFuncOpen.Logic == 2 {
		if level >= pFuncOpen.OpenLevel && viplevel >= pFuncOpen.VipLevel {
			return true
		}
	} else if pFuncOpen.Logic == 3 {
		if viplevel >= pFuncOpen.VipLevel {
			return true
		}
	} else if pFuncOpen.Logic == 4 {
		if level >= pFuncOpen.OpenLevel {
			return true
		}
	} else {
		gamelog.Error("IsFuncOpen Error: invalid Logic %d", pFuncOpen.Logic)
	}

	return false
}

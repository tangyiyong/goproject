package gamedata

import (
	"fmt"
	"gamelog"
)

//! VIP特权表
const (
	VIPLevelLimit = 12
)

type ST_FuncVipInfo struct {
	FuncID   int     //! 功能ID
	VipValue [13]int //! 对应VIP数据
}

var GT_FuncVipList []ST_FuncVipInfo = nil

func InitVipPrivilegeParser(total int) bool {
	GT_FuncVipList = make([]ST_FuncVipInfo, 200)
	return true
}

func ParseVipPrivilegeRecord(rs *RecordSet) {
	funcid := rs.GetFieldInt("func_id")
	GT_FuncVipList[funcid].FuncID = funcid

	//! 加入特权
	for i := 0; i <= VIPLevelLimit; i++ {
		fieldName := fmt.Sprintf("vip%d", i)
		GT_FuncVipList[funcid].VipValue[i] = rs.GetFieldInt(fieldName)
	}
}

//! 查询特权
func GetFuncVipValue(funcID int, vipLevel int) int {
	if vipLevel > VIPLevelLimit || vipLevel < 0 {
		gamelog.Error("GetFuncVipValue Error : Invalid vip level: %d", vipLevel)
		return 0
	}

	if funcID <= FUNC_BEGIN_ID || funcID >= FUNC_END_ID {
		gamelog.Error("GetFuncVipValue Error : Invalid funcID: %d", funcID)
		return 0
	}

	return GT_FuncVipList[funcID].VipValue[vipLevel]
}

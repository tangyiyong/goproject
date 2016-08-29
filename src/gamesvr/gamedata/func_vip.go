package gamedata

import (
	"gamelog"
)

//! VIP特权表
const (
	VIPLevelLimit = int8(12)
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

	GT_FuncVipList[funcid].VipValue[0] = rs.GetFieldInt("vip0")
	GT_FuncVipList[funcid].VipValue[1] = rs.GetFieldInt("vip1")
	GT_FuncVipList[funcid].VipValue[2] = rs.GetFieldInt("vip2")
	GT_FuncVipList[funcid].VipValue[3] = rs.GetFieldInt("vip3")
	GT_FuncVipList[funcid].VipValue[4] = rs.GetFieldInt("vip4")
	GT_FuncVipList[funcid].VipValue[5] = rs.GetFieldInt("vip5")
	GT_FuncVipList[funcid].VipValue[6] = rs.GetFieldInt("vip6")
	GT_FuncVipList[funcid].VipValue[7] = rs.GetFieldInt("vip7")
	GT_FuncVipList[funcid].VipValue[8] = rs.GetFieldInt("vip8")
	GT_FuncVipList[funcid].VipValue[9] = rs.GetFieldInt("vip9")
	GT_FuncVipList[funcid].VipValue[10] = rs.GetFieldInt("vip10")
	GT_FuncVipList[funcid].VipValue[11] = rs.GetFieldInt("vip11")
	GT_FuncVipList[funcid].VipValue[12] = rs.GetFieldInt("vip12")
}

//! 查询特权
func GetFuncVipValue(funcID int, vipLevel int8) int {
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

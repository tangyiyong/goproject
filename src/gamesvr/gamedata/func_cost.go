package gamedata

import (
	"gamelog"
)

type TimesCostMap map[int]int

var GT_FuncCostMap []TimesCostMap

//! 初始化
func InitFuncCostParser(total int) bool {
	GT_FuncCostMap = make([]TimesCostMap, total+1)
	return true
}

//! 解析CSV
func ParseFuncCostRecord(rs *RecordSet) {

	times := CheckAtoi(rs.Values[0], 0)
	if GT_FuncCostMap[times] == nil {
		GT_FuncCostMap[times] = make(TimesCostMap)
	}

	for k, v := range rs.colmap {
		if v == 0 {
			continue
		}

		funcid := CheckAtoiName(k[6:], k)
		GT_FuncCostMap[times][funcid] = CheckAtoiName(rs.Values[v], k)
	}
}

//! 获取重置花费
func GetFuncTimeCost(funcID int, times int) int {
	if times >= len(GT_FuncCostMap) || times <= 0 {
		gamelog.Error("GetResetCost Error : Invalid times :%d", times)
		return 0
	}

	nCost, ok := GT_FuncCostMap[times][funcID]
	if !ok {
		gamelog.Error("GetResetCost Error : Invalid funcID :%d", funcID)
		return 0
	}

	if nCost < 0 {
		for i := times - 1; i > 0; i-- {
			nCost, _ = GT_FuncCostMap[times][funcID]
			if nCost >= 0 {
				break
			}
		}
	}

	return nCost
}

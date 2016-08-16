package gamedata

import (
	"gamelog"
)

type ST_MoneyInfo struct {
	MoneyID int //货币ID
	Max     int //上限
}

var GT_Money_List []ST_MoneyInfo = nil

func InitMoneyParser(total int) bool {
	GT_Money_List = make([]ST_MoneyInfo, total+1)
	return true
}

func ParseMoneyRecord(rs *RecordSet) {
	moneyID := CheckAtoi(rs.Values[0], 0)

	GT_Money_List[moneyID].MoneyID = moneyID
	GT_Money_List[moneyID].Max = CheckAtoi(rs.Values[3], 3)
}

func GetMoneyInfo(moneyID int) *ST_MoneyInfo {
	if moneyID >= len(GT_Money_List) || moneyID <= 0 {
		gamelog.Error("GetMoneyInfo Error: invalid moneyid :%d", moneyID)
		return nil
	}

	return &GT_Money_List[moneyID]
}

func GetMoneyCount() int {
	return len(GT_Money_List) - 1
}

func GetMoneyMaxValue(moneyID int) int {
	if moneyID >= len(GT_Money_List) || moneyID <= 0 {
		gamelog.Error("GetMoneyMaxValue Error: invalid moneyid :%d", moneyID)
		return 0
	}

	return GT_Money_List[moneyID].Max
}

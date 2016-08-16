package gamedata

import (
	"gamelog"
)

const (
	OpenFund_Level  = 1
	OpenFund_BuyNum = 2
)

//! 开服基金奖励表
type ST_OpenFundAward struct {
	ID    int
	Count int
	Award int
}

var GT_OpenFundLst [2][]ST_OpenFundAward

func InitOpenFundParser(total int) bool {
	return true
}

func ParseOpenFundRecord(rs *RecordSet) {
	fundType := rs.GetFieldInt("type")

	var openFundAward ST_OpenFundAward
	openFundAward.ID = rs.GetFieldInt("id")
	openFundAward.Award = rs.GetFieldInt("award")
	openFundAward.Count = rs.GetFieldInt("count")

	GT_OpenFundLst[fundType-1] = append(GT_OpenFundLst[fundType-1], openFundAward)
}

func GetOpenFundInfo(fundType int, id int) *ST_OpenFundAward {
	if fundType > OpenFund_BuyNum || fundType < OpenFund_Level {
		gamelog.Info("GetOpenFundInfo error: invalid fundtype: %d", fundType)
		return nil
	}

	if id > len(GT_OpenFundLst[fundType-1]) || id <= 0 {
		gamelog.Error("GetOpenFundInfo error: not found type: %d id: %d", fundType, id)
		return nil
	}

	return &GT_OpenFundLst[fundType-1][id-1]
}

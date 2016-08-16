package gamedata

import (
	"gamelog"
)

//！ 卡牌兑换物品表
type TCMExchangeItemCsv struct {
	ID         int
	AwardType  byte
	Items      []ST_ItemData
	NeedCards  []ST_ItemData
	DailyTimes uint16
}

var G_CMExchangeItemCsv []TCMExchangeItemCsv

func InitCMExchangeItemParser(total int) bool {
	G_CMExchangeItemCsv = make([]TCMExchangeItemCsv, total+1)
	return true
}
func ParseCMExchangeItemRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_CMExchangeItemCsv[id]
	data.ID = id
	data.AwardType = byte(rs.GetFieldInt("award_type"))
	data.Items = rs.GetFieldItems("item")
	data.NeedCards = rs.GetFieldItems("need_card")
	data.DailyTimes = uint16(rs.GetFieldInt("daily_times"))
}
func GetCMExchangeItemCsvInfo(exchangeID int) *TCMExchangeItemCsv {
	if exchangeID <= 0 || exchangeID >= len(G_CMExchangeItemCsv) {
		gamelog.Error("GetCMExchangeItemCsvInfo Error: Invalid ID:%d", exchangeID)
		return nil
	}
	return &G_CMExchangeItemCsv[exchangeID]
}

//！ 卡牌表
type TCardCsv struct {
	ID        int
	PointBuy  int
	PointSell int
}

var G_CardCsv []TCardCsv

func InitCardCsvParser(total int) bool {
	G_CardCsv = make([]TCardCsv, total+1)
	return true
}
func ParseCardCsvRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_CardCsv[id]
	data.ID = id
	data.PointBuy = rs.GetFieldInt("point_buy")
	data.PointSell = rs.GetFieldInt("point_sell")
}
func GetCardCsvInfo(cardID int) *TCardCsv {
	if cardID <= 0 || cardID >= len(G_CardCsv) {
		gamelog.Error("GetCardCsvInfo Error: Invalid ID:%d", cardID)
		return nil
	}
	return &G_CardCsv[cardID]
}

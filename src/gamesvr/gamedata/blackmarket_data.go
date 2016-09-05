package gamedata

import (
	"gamelog"
)

//! 黑市商店表
type ST_BlackMarket struct {
	ID           int
	ItemID       int
	ItemNum      int
	CostMoneyID  int
	CostMoneyNum int
	Level_Min    int
	Level_Max    int
	Recommend    int
}

var GT_BlackMarketLst []ST_BlackMarket

func InitBlackMarketParser(total int) bool {
	GT_BlackMarketLst = make([]ST_BlackMarket, total)
	return true
}

func ParseBlackMarketRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_BlackMarketLst[id-1].ID = id
	GT_BlackMarketLst[id-1].ItemID = rs.GetFieldInt("itemid")
	GT_BlackMarketLst[id-1].ItemNum = rs.GetFieldInt("itemnum")
	GT_BlackMarketLst[id-1].CostMoneyID = rs.GetFieldInt("costmoneyid")
	GT_BlackMarketLst[id-1].CostMoneyNum = rs.GetFieldInt("costmoneynum")
	GT_BlackMarketLst[id-1].Level_Min = rs.GetFieldInt("level_min")
	GT_BlackMarketLst[id-1].Level_Max = rs.GetFieldInt("level_max")
	GT_BlackMarketLst[id-1].Recommend = rs.GetFieldInt("recommend")
}

func BlackMarketRandGoods(level int) []int {
	var randLst []ST_BlackMarket
	for _, v := range GT_BlackMarketLst {
		if level >= v.Level_Min && level <= v.Level_Max {
			randLst = append(randLst, v)
		}
	}

	//! 随机六个商品
	goodsLst := []int{}
	length := len(randLst)
	for i := 0; i < 6; i++ {

		id := randLst[r.Intn(length)].ID

		isExist := false
		for _, v := range goodsLst {
			if v == id {
				isExist = true
				break
			}
		}

		if isExist == true {
			i--
			continue
		}

		goodsLst = append(goodsLst, randLst[r.Intn(len(randLst))].ID)
	}

	return goodsLst
}

//! 获取黑市商品信息
func GetBlackMarketGoodsInfo(id int) *ST_BlackMarket {
	if id > len(GT_BlackMarketLst) {
		gamelog.Error("GetBlackMarketGoodsInfo Error: invalid id %v", id)
		return nil
	}

	return &GT_BlackMarketLst[id-1]
}

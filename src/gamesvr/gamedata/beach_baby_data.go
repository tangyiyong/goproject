package gamedata

import (
	"gamelog"
	"math/rand"
)

type TBeachBabyGoodsCsv struct {
	ID         int
	AwardType  int
	ItemID     int
	ItemNum    int
	Weight     int
	isSelected bool // 做随机用的填充变量
}

var G_BeachBabyGoodsCsv []TBeachBabyGoodsCsv
var G_BeachBabyGoods_Type map[int][]*TBeachBabyGoodsCsv // [活动AwardType] = 商品列表

func GetBeachBabyGoodsCsv(id int) *TBeachBabyGoodsCsv {
	if id <= 0 || id >= len(G_BeachBabyGoodsCsv) {
		gamelog.Error("GetBeachBabyGoodsCsv Error: Invalid ID:%d", id)
		return nil
	}
	return &G_BeachBabyGoodsCsv[id]
}

func CreateBeachBabyGoodsTypeMap() {
	if G_BeachBabyGoods_Type == nil {
		G_BeachBabyGoods_Type = make(map[int][]*TBeachBabyGoodsCsv)

		for i := 1; i < len(G_BeachBabyGoodsCsv); i++ {
			data := &G_BeachBabyGoodsCsv[i]
			G_BeachBabyGoods_Type[data.AwardType] = append(G_BeachBabyGoods_Type[data.AwardType], data)
		}
	}
}
func RandSelect_BeachBabyGoods(activityID int, selectCnt int) (ret []int) {
	CreateBeachBabyGoodsTypeMap()

	csv := GetActivityInfo(activityID)
	if csv == nil {
		gamelog.Error("RandSelect_BeachBabyGoods GetActivityInfo() Error: ActivityID:%d", activityID)
		return nil
	}
	goodsList := G_BeachBabyGoods_Type[csv.AwardType]
	total, length := 0, len(goodsList)
	for i := 0; i < length; i++ {
		goodsList[i].isSelected = false
		total += goodsList[i].Weight
	}

	if selectCnt > length {
		gamelog.Error("RandSelect_BeachBabyGoods Error: Goods not enough!!! ActID:%d, AwardType:%d, length:%d, selectCnt:%d", activityID, csv.AwardType, length, selectCnt)
		return nil
	}

	for j := 0; j < selectCnt; j++ {
		rand := rand.Intn(total)
		// gamelog.Info("RandSelect_BeachBabyGoods ---- rand:%d, idx:%d", rand, j+1)
		for i := 0; i < length; i++ {
			goods := goodsList[i] // 此处已经是指针了
			if goods.isSelected {
				// gamelog.Info("RandSelect_BeachBabyGoods ---- continue ID:%d", goods.ID)
				continue
			}
			if rand < goods.Weight {
				// gamelog.Info("RandSelect_BeachBabyGoods ---- got ID:%d  W:%d", goods.ID, goods.Weight)
				ret = append(ret, goods.ID)
				goods.isSelected = true
				total -= goods.Weight
				break
			} else {
				rand -= goods.Weight
			}
		}
	}
	return ret
}

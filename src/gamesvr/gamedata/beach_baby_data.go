package gamedata

import (
	"gamelog"
	"math/rand"
)

type TBeachBabyItemInfo struct {
	ID         int
	AwardType  int
	ItemID     int
	ItemNum    int
	Weight     int
	isSelected bool // 做随机用的填充变量
}

var G_BeachBabyItem_List []TBeachBabyItemInfo
var G_BeachBabyGoods_Type map[int][]*TBeachBabyItemInfo // [活动AwardType] = 商品列表

func InitBeachBabyParser(total int) bool {
	G_BeachBabyItem_List = make([]TBeachBabyItemInfo, total+1)
	return true
}

func ParseBeachBabyRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	G_BeachBabyItem_List[id].ID = id
	G_BeachBabyItem_List[id].AwardType = rs.GetFieldInt("award_type")
	G_BeachBabyItem_List[id].ItemID = rs.GetFieldInt("itemid")
	G_BeachBabyItem_List[id].ItemNum = rs.GetFieldInt("itemnum")
	G_BeachBabyItem_List[id].Weight = rs.GetFieldInt("weight")
}

func GetBeachBabyItemInfo(id int) *TBeachBabyItemInfo {
	if id <= 0 || id >= len(G_BeachBabyItem_List) {
		gamelog.Error("GetBeachBabyItemInfo Error: Invalid ID:%d", id)
		return nil
	}
	return &G_BeachBabyItem_List[id]
}

func CreateBeachBabyGoodsTypeMap() {
	if G_BeachBabyGoods_Type == nil {
		G_BeachBabyGoods_Type = make(map[int][]*TBeachBabyItemInfo)

		for i := 1; i < len(G_BeachBabyItem_List); i++ {
			data := &G_BeachBabyItem_List[i]
			G_BeachBabyGoods_Type[data.AwardType] = append(G_BeachBabyGoods_Type[data.AwardType], data)
		}
	}
}
func RandSelect_BeachBabyGoods(awardType int, selectCnt int) (ret []int) {
	CreateBeachBabyGoodsTypeMap()

	goodsList := G_BeachBabyGoods_Type[awardType]
	total, length := 0, len(goodsList)
	for i := 0; i < length; i++ {
		goodsList[i].isSelected = false
		total += goodsList[i].Weight
	}

	if selectCnt > length {
		gamelog.Error("RandSelect_BeachBabyGoods Error: Goods not enough!!! AwardType:%d, length:%d, selectCnt:%d", awardType, length, selectCnt)
		return nil
	}

	for j := 0; j < selectCnt; j++ {
		rand := rand.Intn(total)

		for i := 0; i < length; i++ {
			goods := goodsList[i] // 此处已经是指针了
			if goods.isSelected {

				continue
			}
			if rand < goods.Weight {
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

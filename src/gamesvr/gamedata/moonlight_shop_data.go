package gamedata

import (
	"gamelog"
	"math/rand"
)

//！ 月光币兑换表
type TMoonlightShopExchangeCsv struct {
	ID         int
	GetToken   int
	CostType   int
	CostNum    int
	DailyTimes byte
}

var G_MoonlightShopExchangeCsv []TMoonlightShopExchangeCsv

func InitMoonlightShopExchangeCsv(total int) bool {
	G_MoonlightShopExchangeCsv = make([]TMoonlightShopExchangeCsv, total+1)
	return true
}
func ParseMoonlightShopExchangeCsv(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_MoonlightShopExchangeCsv[id]
	data.ID = id
	data.GetToken = rs.GetFieldInt("get_token")
	data.CostType = rs.GetFieldInt("cost_type")
	data.CostNum = rs.GetFieldInt("cost_num")
	data.DailyTimes = byte(rs.GetFieldInt("daily_times"))
}
func GetMoonlightShopExchangeCsv(id int) *TMoonlightShopExchangeCsv {
	if id <= 0 || id >= len(G_MoonlightShopExchangeCsv) {
		gamelog.Error("GetMoonlightShopExchangeCsv Error: Invalid ID:%d", id)
		return nil
	}
	return &G_MoonlightShopExchangeCsv[id]
}

//！ 商品表
type TMoonlightGoodsCsv struct {
	ID          int
	AwardType   int
	ItemID      int
	ItemNum     int
	Price       int
	MinDiscount byte // 百分比
	MaxDiscount byte
	DailyTimes  byte
	Weight      int
	isSelected  bool // 做随机用的填充变量
}

var G_MoonlightGoodsCsv []TMoonlightGoodsCsv
var G_MoonlightGoods_Type map[int][]*TMoonlightGoodsCsv // [活动AwardType] = 商品列表

func InitMoonlightGoodsCsv(total int) bool {
	G_MoonlightGoodsCsv = make([]TMoonlightGoodsCsv, total+1)
	G_MoonlightGoods_Type = make(map[int][]*TMoonlightGoodsCsv)
	return true
}
func ParseMoonlightGoodsCsv(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_MoonlightGoodsCsv[id]
	data.ID = id
	data.AwardType = rs.GetFieldInt("award_type")
	data.ItemID = rs.GetFieldInt("itemid")
	data.ItemNum = rs.GetFieldInt("itemnum")
	data.Price = rs.GetFieldInt("price")
	data.MinDiscount = byte(rs.GetFieldInt("discount_min"))
	data.MaxDiscount = byte(rs.GetFieldInt("discount_max"))
	data.DailyTimes = byte(rs.GetFieldInt("daily_times"))
	data.Weight = rs.GetFieldInt("weight")

	G_MoonlightGoods_Type[data.AwardType] = append(G_MoonlightGoods_Type[data.AwardType], data)
}
func GetMoonlightGoodsCsv(id int) *TMoonlightGoodsCsv {
	if id <= 0 || id >= len(G_MoonlightGoodsCsv) {
		gamelog.Error("GetMoonlightGoodsCsv Error: Invalid ID:%d", id)
		return nil
	}
	return &G_MoonlightGoodsCsv[id]
}
func RandSelect_MoonlightGoods(activityID int, selectCnt int) (ret []int) {
	csv := GetActivityInfo(activityID)
	if csv == nil {
		gamelog.Error("RandSelect_MoonlightGoods GetActivityInfo() Error: ActivityID:%d", activityID)
		return nil
	}
	goodsList := G_MoonlightGoods_Type[csv.AwardType]
	total, length := 0, len(goodsList)
	for i := 0; i < length; i++ {
		goodsList[i].isSelected = false
		total += goodsList[i].Weight
	}

	if selectCnt > length {
		gamelog.Error("RandSelect_MoonlightGoods Error: Goods not enough!!! AwardType:%d, length:%d, selectCnt:%d", csv.AwardType, length, selectCnt)
		return nil
	}

	for j := 0; j < selectCnt; j++ {
		rand := rand.Intn(total)
		for i := 0; i < selectCnt; i++ {
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

//! 积分奖励表
type TMoonlightAwardCsv struct {
	ID        int
	AwardType int
	ItemID    int
	ItemNum   int
	NeedScore int
}

var G_MoonlightAwardCsv []TMoonlightAwardCsv

func InitMoonlightShopAwardCsv(total int) bool {
	total++
	if total >= 64 {
		gamelog.Error("MoonlightShopAward Init Error: length must less then 64") // 使用位标记，记录领奖情况
	}

	G_MoonlightAwardCsv = make([]TMoonlightAwardCsv, total)
	return true
}
func ParseMoonlightShopAwardCsv(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_MoonlightAwardCsv[id]
	data.ID = id
	data.AwardType = rs.GetFieldInt("award_type")
	data.ItemID = rs.GetFieldInt("itemid")
	data.ItemNum = rs.GetFieldInt("itemnum")
	data.NeedScore = rs.GetFieldInt("need_score")
}
func GetMoonlightShopAwardCsv(id int) *TMoonlightAwardCsv {
	if id <= 0 || id >= len(G_MoonlightAwardCsv) {
		gamelog.Error("GetMoonlightShopAwardCsv Error: Invalid ID:%d", id)
		return nil
	}
	return &G_MoonlightAwardCsv[id]
}

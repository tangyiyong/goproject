package gamedata

import (
	"gamelog"
	"math/rand"
)

//！ 月光币兑换表
type ST_MoonShopExchangeInfo struct {
	ID          int
	ItemID      int
	ItemNum     int
	CostMoneyID int
	CostNum     int
	DailyTimes  byte
}

var G_MoonShopExchg_List []ST_MoonShopExchangeInfo
var MoonShop_Money_ID = 0

func InitMoonShopExchgParser(total int) bool {
	G_MoonShopExchg_List = make([]ST_MoonShopExchangeInfo, total+1)
	return true
}
func ParseMoontShopExchgRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	G_MoonShopExchg_List[id].ID = id
	G_MoonShopExchg_List[id].ItemID = rs.GetFieldInt("item_id")
	G_MoonShopExchg_List[id].ItemNum = rs.GetFieldInt("item_num")
	G_MoonShopExchg_List[id].CostMoneyID = rs.GetFieldInt("cost_money_id")
	G_MoonShopExchg_List[id].CostNum = rs.GetFieldInt("cost_num")
	G_MoonShopExchg_List[id].DailyTimes = byte(rs.GetFieldInt("daily_times"))
	MoonShop_Money_ID = G_MoonShopExchg_List[id].ItemID
}
func GetMoonShopExchgInfo(id int) *ST_MoonShopExchangeInfo {
	if id <= 0 || id >= len(G_MoonShopExchg_List) {
		gamelog.Error("GetMoonShopExchgInfo Error: Invalid ID:%d", id)
		return nil
	}
	return &G_MoonShopExchg_List[id]
}

//！ 商品表
type ST_MoonGoodsInfo struct {
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

var G_MoonGoods_List []ST_MoonGoodsInfo
var G_MoonlightGoods_Type map[int][]*ST_MoonGoodsInfo // [活动AwardType] = 商品列表

func InitMoonGoodsParser(total int) bool {
	G_MoonGoods_List = make([]ST_MoonGoodsInfo, total+1)
	G_MoonlightGoods_Type = make(map[int][]*ST_MoonGoodsInfo)
	return true
}
func ParseMoonGoodsRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_MoonGoods_List[id]
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
func GetMoonGoodsInfo(id int) *ST_MoonGoodsInfo {
	if id <= 0 || id >= len(G_MoonGoods_List) {
		gamelog.Error("v Error: Invalid ID:%d", id)
		return nil
	}
	return &G_MoonGoods_List[id]
}
func RandSelect_MoonlightGoods(activityID int32, selectCnt int) (ret []int) {
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
type ST_MoonAwardInfo struct {
	ID        int
	AwardType int
	ItemID    int
	ItemNum   int
	NeedScore int
}

var G_MoonAward_List []ST_MoonAwardInfo

func InitMoonShopAwardParser(total int) bool {
	G_MoonAward_List = make([]ST_MoonAwardInfo, total+1)
	return true
}
func ParseMoonShopAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	data := &G_MoonAward_List[id]
	data.ID = id
	data.AwardType = rs.GetFieldInt("award_type")
	data.ItemID = rs.GetFieldInt("itemid")
	data.ItemNum = rs.GetFieldInt("itemnum")
	data.NeedScore = rs.GetFieldInt("need_score")
}
func GetMoonShopAwardInfo(id int) *ST_MoonAwardInfo {
	if id <= 0 || id >= len(G_MoonAward_List) {
		gamelog.Error("GetMoonShopAwardInfo Error: Invalid ID:%d", id)
		return nil
	}
	return &G_MoonAward_List[id]
}

package gamedata

import (
	"gamelog"
)

const (
	HuntTreasureMap_Money = 1
	HuntTreasureMap_Item  = 2
	HuntTreasureMap_Event = 3
	HuntTreasureMap_Store = 4
	HuntTreasureMap_Move  = 5
)

//! 巡回探宝地图
type ST_HuntTreasureMap struct {
	ID        int //! 唯一ID
	AwardType int //! 活动ID
	Type      int //! 格子类型
	Award     int //! 格子奖励 若无奖励则为0
}

var GT_HuntTreasureMapLst map[int][]ST_HuntTreasureMap

func InitHuntTreasureMapParser(total int) bool {
	GT_HuntTreasureMapLst = make(map[int][]ST_HuntTreasureMap)
	return true
}

func ParseHuntTreasureMapRecord(rs *RecordSet) {
	var mapNode ST_HuntTreasureMap
	mapNode.ID = CheckAtoi(rs.Values[0], 0)
	mapNode.AwardType = rs.GetFieldInt("award_type")
	mapNode.Award = rs.GetFieldInt("award")
	mapNode.Type = rs.GetFieldInt("type")
	GT_HuntTreasureMapLst[mapNode.AwardType] = append(GT_HuntTreasureMapLst[mapNode.AwardType], mapNode)
}

func GetHuntTreasureMapCount(awardType int) int {
	return len(GT_HuntTreasureMapLst[awardType])
}

func GetHuntTreasureMap(id int, awardType int) *ST_HuntTreasureMap {
	if id > len(GT_HuntTreasureMapLst[awardType]) || id <= 0 {
		gamelog.Info("GetHuntTreasureMap Error: Invalid ID %d", id)
		return nil
	}

	return &GT_HuntTreasureMapLst[awardType][id-1]
}

//! 巡回探宝商店
type ST_HuntTreasureStore struct {
	ID        int //! 唯一ID
	AwardType int //! 活动奖励类型
	ItemID    int //! 道具ID
	ItemNum   int //! 道具数量
	MoneyID   int //! 货币ID
	MoneyNum  int //! 货币数量
	Score     int //! 购买增加积分
	Weight    int //! 权重
}

var GT_HuntTreasureStoreLst map[int][]ST_HuntTreasureStore

func InitHuntTreasureStoreParser(total int) bool {
	GT_HuntTreasureStoreLst = make(map[int][]ST_HuntTreasureStore)
	return true
}

func ParseHuntTreasureStoreRecord(rs *RecordSet) {
	var store ST_HuntTreasureStore

	id := CheckAtoi(rs.Values[0], 0)
	store.ID = id
	store.AwardType = rs.GetFieldInt("award_type")
	store.ItemID = rs.GetFieldInt("itemid")
	store.ItemNum = rs.GetFieldInt("itemnum")
	store.MoneyID = rs.GetFieldInt("moneyid")
	store.MoneyNum = rs.GetFieldInt("moneynum")
	store.Score = rs.GetFieldInt("score")
	store.Weight = rs.GetFieldInt("weight")
	GT_HuntTreasureStoreLst[store.AwardType] = append(GT_HuntTreasureStoreLst[store.AwardType], store)
}

func RandHuntTreasureStoreItem(num int, awardType int) (itemLst []int) {
	totalWeight := 0
	for _, v := range GT_HuntTreasureStoreLst[awardType] {
		totalWeight += v.Weight
	}

	for i := 0; i < num; i++ {
		randWeight := r.Intn(totalWeight)
		curWeight := 0

		for _, v := range GT_HuntTreasureStoreLst[awardType] {
			if randWeight >= curWeight && randWeight < v.Weight+curWeight {
				itemLst = append(itemLst, v.ID)
				break
			}
			curWeight += v.Weight
		}
	}

	return itemLst
}

func GetHuntTreasureStoreItem(id int, awardType int) *ST_HuntTreasureStore {
	if id > len(GT_HuntTreasureStoreLst[awardType]) || id <= 0 {
		gamelog.Info("GetHuntTreasureStoreItem Error: Invalid ID %d", id)
		return nil
	}

	return &GT_HuntTreasureStoreLst[awardType][id-1]
}

//! 巡回奖励
type ST_HuntTreasureAward struct {
	ID        int
	AwardType int //! 活动奖励
	NeedTurn  int
	Award     int
}

var GT_HuntTreasureAwardLst map[int][]ST_HuntTreasureAward

func InitHuntTreasureAwardParser(total int) bool {
	GT_HuntTreasureAwardLst = make(map[int][]ST_HuntTreasureAward)
	return true
}

func ParseHuntTreasureAwardRecord(rs *RecordSet) {
	var award ST_HuntTreasureAward
	id := CheckAtoi(rs.Values[0], 1)
	award.ID = id
	award.AwardType = rs.GetFieldInt("award_type")
	award.NeedTurn = rs.GetFieldInt("need_turn")
	award.Award = rs.GetFieldInt("award")
	GT_HuntTreasureAwardLst[award.AwardType] = append(GT_HuntTreasureAwardLst[award.AwardType], award)
}

func GetHuntTreasureAwardCount(awardType int) int {
	return len(GT_HuntTreasureAwardLst[awardType])
}

func GetHuntTreasureAward(id int, awardType int) *ST_HuntTreasureAward {
	if id > len(GT_HuntTreasureAwardLst[awardType]) || id <= 0 {
		gamelog.Info("GetHuntTreasureAward Error: Invalid ID %d", id)
		return nil
	}

	return &GT_HuntTreasureAwardLst[awardType][id-1]
}

type ST_OperationalRankAward struct {
	ID                   int
	ActivityType         int //! 活动类型
	AwardType            int //! 活动奖励模板
	Rank_Min             int //! 排名取值区间
	Rank_Max             int
	TodayNormalRankAward int //! 今日普通榜奖励
	TodayEliteRankAward  int //! 今日精英榜奖励
	TotalNormalRankAward int //! 总计普通榜奖励
	TotalEliteRankAward  int //! 总计精英榜奖励

}

var GT_OperationalRankAwardLst []ST_OperationalRankAward

func InitOperationalRankAwardParser(total int) bool {
	GT_OperationalRankAwardLst = make([]ST_OperationalRankAward, total+1)
	return true
}

func ParseOperationalRankAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_OperationalRankAwardLst[id].ID = id
	GT_OperationalRankAwardLst[id].ActivityType = rs.GetFieldInt("activity_type")
	GT_OperationalRankAwardLst[id].AwardType = rs.GetFieldInt("award_type")
	GT_OperationalRankAwardLst[id].Rank_Min = rs.GetFieldInt("rank_min")
	GT_OperationalRankAwardLst[id].Rank_Max = rs.GetFieldInt("rank_max")
	GT_OperationalRankAwardLst[id].TodayNormalRankAward = rs.GetFieldInt("today_normal_award")
	GT_OperationalRankAwardLst[id].TodayEliteRankAward = rs.GetFieldInt("today_elite_award")
	GT_OperationalRankAwardLst[id].TotalNormalRankAward = rs.GetFieldInt("total_normal_award")
	GT_OperationalRankAwardLst[id].TotalEliteRankAward = rs.GetFieldInt("total_elite_award")

}

func GetOperationalRankAwardFromRank(activityType int, awardType int, rank int) *ST_OperationalRankAward {
	for i, v := range GT_OperationalRankAwardLst {
		if v.ID == 0 || v.AwardType != awardType || v.ActivityType != activityType {
			continue
		}

		if v.Rank_Min <= rank && rank <= v.Rank_Max {
			return &GT_OperationalRankAwardLst[i]
		}
	}

	return nil
}

//! 幸运轮盘库
type ST_LuckyWheel struct {
	ID        int //! 唯一标识
	AwardType int //! 活动ID
	Day       int //! 天数
	Type      int //! 1->普通 2->豪华
	ItemID    int //! 物品ID
	ItemNum   int //! 物品数量
	IsSpecial int //! 是否为特殊奖励
	Weight    int //! 权重
}

var GT_LuckyWheelLst []ST_LuckyWheel

func InitLuckyWheelParser(total int) bool {
	GT_LuckyWheelLst = make([]ST_LuckyWheel, total+1)
	return true
}

func ParseLuckyWheelRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_LuckyWheelLst[id].ID = id
	GT_LuckyWheelLst[id].AwardType = rs.GetFieldInt("award_type")
	GT_LuckyWheelLst[id].Type = rs.GetFieldInt("type")
	GT_LuckyWheelLst[id].ItemID = rs.GetFieldInt("item_id")
	GT_LuckyWheelLst[id].ItemNum = rs.GetFieldInt("item_num")
	GT_LuckyWheelLst[id].IsSpecial = rs.GetFieldInt("is_special")
	GT_LuckyWheelLst[id].Weight = rs.GetFieldInt("weight")
}

func GetLuckyWheelItemFromDay(awardType int, wheelType int) []int {
	idLst := []int{}
	for _, v := range GT_LuckyWheelLst {
		if wheelType == v.Type && v.AwardType == awardType {
			idLst = append(idLst, v.ID)
		}
	}

	return idLst
}

func GetLuckyWheelItemFromID(id int) *ST_LuckyWheel {
	if id > len(GT_LuckyWheelLst)-1 || id <= 0 {
		gamelog.Info("GetLuckyWheelItemFromID Error: Invalid ID %d", id)
		return nil
	}

	return &GT_LuckyWheelLst[id]
}

//! 团购库
type ST_GroupPurchase struct {
	AwardType   int //! 活动ID
	ItemID      int //! 道具ID
	ItemNum     int //! 道具数量
	Discount    int //! 折扣
	NeedSaleNum int //! 需要销售数量
	UseItemMax  int //! 使用折扣券上限
	MoneyNum    int //! 货币数量
	BuyTimes    int //! 每日可购买次数
}

var GT_GroupPurchaseLst map[int][]ST_GroupPurchase

func InitGroupPurchaseParser(total int) bool {
	GT_GroupPurchaseLst = make(map[int][]ST_GroupPurchase)
	return true
}

func ParseGroupPurchaseRecord(rs *RecordSet) {
	var groupPurchase ST_GroupPurchase
	groupPurchase.AwardType = rs.GetFieldInt("award_type")
	groupPurchase.ItemID = rs.GetFieldInt("item_id")
	groupPurchase.ItemNum = rs.GetFieldInt("item_num")
	groupPurchase.Discount = rs.GetFieldInt("discount")
	groupPurchase.NeedSaleNum = rs.GetFieldInt("need_num")
	groupPurchase.UseItemMax = rs.GetFieldInt("use_max")
	groupPurchase.MoneyNum = rs.GetFieldInt("money_num")
	groupPurchase.BuyTimes = rs.GetFieldInt("buy_times")
	GT_GroupPurchaseLst[groupPurchase.AwardType] = append(GT_GroupPurchaseLst[groupPurchase.AwardType], groupPurchase)
}

func GetGroupPurchaseItemInfo(itemID int, awardType int) *ST_GroupPurchase {
	_, ok := GT_GroupPurchaseLst[awardType]
	if ok == false {
		gamelog.Error("GetGroupPurchaseItemInfo Error: Can't find activity %d", awardType)
		return nil
	}

	for i, _ := range GT_GroupPurchaseLst[awardType] {
		if GT_GroupPurchaseLst[awardType][i].ItemID == itemID {
			return &GT_GroupPurchaseLst[awardType][i]
		}
	}
	return nil
}

//! 根据销量计算当前价格
func GetGroupPurchaseItemInfoFromSale(itemID int, awardType int, saleNum int) *ST_GroupPurchase {
	_, ok := GT_GroupPurchaseLst[awardType]
	if ok == false {
		gamelog.Error("GetGroupPurchaseItemInfoFromSale Error: Can't find activity %d", awardType)
		return nil
	}

	index := -1
	for i, _ := range GT_GroupPurchaseLst[awardType] {
		if GT_GroupPurchaseLst[awardType][i].ItemID != itemID {
			continue
		}

		if GT_GroupPurchaseLst[awardType][i].NeedSaleNum <= saleNum {
			index = i
		}
	}

	if index < 0 {
		gamelog.Error("GetGroupPurchaseItemInfoFromSale Error: Can't find item %d", itemID)
		return nil
	}

	return &GT_GroupPurchaseLst[awardType][index]
}

//! 团购积分奖励表
type ST_GroupPurchaseScoreAward struct {
	ID        int //! 唯一标识
	AwardType int //! 活动ID
	ItemID    int //! 奖励道具ID
	ItemNum   int //! 奖励道具数量
	NeedScore int //! 需求积分
}

var GT_GroupPurchaseScoreLst []ST_GroupPurchaseScoreAward

func InitGroupPurchaseScoreParser(total int) bool {
	GT_GroupPurchaseScoreLst = make([]ST_GroupPurchaseScoreAward, total+1)
	return true
}

func ParseGroupPurchaseScoreRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)

	GT_GroupPurchaseScoreLst[id].ID = rs.GetFieldInt("id")
	GT_GroupPurchaseScoreLst[id].AwardType = rs.GetFieldInt("award_type")
	GT_GroupPurchaseScoreLst[id].ItemID = rs.GetFieldInt("itemid")
	GT_GroupPurchaseScoreLst[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_GroupPurchaseScoreLst[id].NeedScore = rs.GetFieldInt("need_score")
}

func GetGroupPurchaseScoreAward(id int) *ST_GroupPurchaseScoreAward {
	if id > len(GT_GroupPurchaseScoreLst)-1 || id < 0 {
		gamelog.Error("GetGroupPurchaseScoreAward Error: Invalid id %d", id)
		return nil
	}
	return &GT_GroupPurchaseScoreLst[id]
}

func GetGroupPurchaseScoreAwardCount() int {
	return len(GT_GroupPurchaseScoreLst) - 1
}

//! 欢庆佳节商店表
type ST_FestivalSale struct {
	ID        int //! 唯一标识
	AwardType int //! 活动奖励模板
	ItemID    int //! 商品ID
	ItemNum   int //! 商品数量
	MoneyID   int //! 货币ID
	MoneyNum  int //! 货币数量
	BuyTimes  int //! 货币
}

var GT_FestivalBuyLst []ST_FestivalSale

func InitFestivalSaleParser(total int) bool {
	GT_FestivalBuyLst = make([]ST_FestivalSale, total+1)
	return true
}

func ParseFestivalSaleRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_FestivalBuyLst[id].ID = rs.GetFieldInt("id")
	GT_FestivalBuyLst[id].ItemID = rs.GetFieldInt("item_id")
	GT_FestivalBuyLst[id].ItemNum = rs.GetFieldInt("item_num")
	GT_FestivalBuyLst[id].AwardType = rs.GetFieldInt("award_type")
	GT_FestivalBuyLst[id].MoneyID = rs.GetFieldInt("money_id")
	GT_FestivalBuyLst[id].MoneyNum = rs.GetFieldInt("money_num")
	GT_FestivalBuyLst[id].BuyTimes = rs.GetFieldInt("buy_times")
}

func GetFestivalItemInfo(id int) *ST_FestivalSale {
	if id < 0 || id > len(GT_FestivalBuyLst)-1 {
		gamelog.Error("GetFestivalItemInfo Error: Invalid id %d", id)
		return nil
	}

	return &GT_FestivalBuyLst[id]
}

//! 欢庆佳节任务库
type ST_FestivalTask struct {
	ID        int //! 唯一标识
	AwardType int //! 活动奖励模板
	TaskType  int //! 任务类型
	Need      int //! 达标数额
	Award     int //! 奖励
}

var GT_FestivalTaskLst map[int][]ST_FestivalTask

func InitFestivalTaskParser(total int) bool {
	GT_FestivalTaskLst = make(map[int][]ST_FestivalTask)
	return true
}

func ParseFestivalTaskRecord(rs *RecordSet) {

	var task ST_FestivalTask
	task.ID = rs.GetFieldInt("id")
	task.AwardType = rs.GetFieldInt("award_type")
	task.TaskType = rs.GetFieldInt("task_type")
	task.Need = rs.GetFieldInt("need")
	task.Award = rs.GetFieldInt("award")
	GT_FestivalTaskLst[task.AwardType] = append(GT_FestivalTaskLst[task.AwardType], task)
}

func GetFestivalTaskFromType(awardType int) []ST_FestivalTask {
	taskLst, ok := GT_FestivalTaskLst[awardType]
	if ok == false {
		gamelog.Error("GetFestivalTaskFromType Error: Not exist %d", awardType)
		return []ST_FestivalTask{}
	}

	return taskLst
}

func GetFestivalTaskInfo(awardType int, taskID int) *ST_FestivalTask {
	length := len(GT_FestivalTaskLst[awardType])
	for i := 0; i < length; i++ {
		if GT_FestivalTaskLst[awardType][i].ID == taskID {
			return &GT_FestivalTaskLst[awardType][i]
		}
	}
	return nil
}

//! 欢庆佳节兑换表
type ST_FestivalExchange struct {
	ID            int //! 唯一标识
	AwardType     int //! 活动奖励模板
	NeedItemID    int //! 需求道具
	NeedItemNum   int //! 需求道具数量
	Award         int //! 兑换奖励
	ExchangeTimes int //! 兑换次数
}

var GT_FestivalExchangeLst []ST_FestivalExchange

func InitFestivalExchangeParser(total int) bool {
	GT_FestivalExchangeLst = make([]ST_FestivalExchange, total+1)
	return true
}

func ParseFestivalExchangeRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_FestivalExchangeLst[id].ID = rs.GetFieldInt("id")
	GT_FestivalExchangeLst[id].AwardType = rs.GetFieldInt("award_type")
	GT_FestivalExchangeLst[id].NeedItemID = rs.GetFieldInt("need_item_id")
	GT_FestivalExchangeLst[id].NeedItemNum = rs.GetFieldInt("need_item_num")
	GT_FestivalExchangeLst[id].Award = rs.GetFieldInt("award")
	GT_FestivalExchangeLst[id].ExchangeTimes = rs.GetFieldInt("exchange_times")
}

func GetExchangeInfoFromID(id int) *ST_FestivalExchange {
	if id <= 0 && id > len(GT_FestivalExchangeLst)-1 {
		gamelog.Error("GetExchangeInfoFromID Error: Invalid id %d", id)
		return nil
	}

	return &GT_FestivalExchangeLst[id]
}

func GetExchangeInfoLst(awardType int) (exchangeLst []ST_FestivalExchange) {
	length := len(GT_FestivalExchangeLst)
	for i := 0; i < length; i++ {
		if GT_FestivalExchangeLst[i].AwardType == awardType && GT_FestivalExchangeLst[i].ID != 0 {
			exchangeLst = append(exchangeLst, GT_FestivalExchangeLst[i])
		}
	}

	return exchangeLst
}

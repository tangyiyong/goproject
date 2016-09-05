package msg

//! 玩家查询巡回探宝状态
//! 消息: /query_hunt_treasure
type MSG_QueryHuntTreasure_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryHuntTreasure_Ack struct {
	RetCode      int
	CurrentPos   int  //! 当前所在格ID
	HuntTurns    int  //! 巡回轮数
	Score        int  //! 总积分
	EndCountDown int  //! 活动结束倒计时
	TodayRank    int  //! 今日排名变化
	TotalRank    int  //! 累计排名变化
	AwardType    int  //! 返回奖励模板
	FreeTimes    int  //! 今日免费次数
	IsHaveStore  bool //! 是否存在商店
}

//! 玩家开始掷骰
//! 消息: /start_hunt
type Msg_StartHuntTreasure_Req struct {
	PlayerID        int32
	SessionKey      string
	IsUseLucklyDice int //! 是否使用幸运骰子
	IsStartTenTimes int //! 投掷十次
	Steps           int //! 幸运骰子决定点数
}

type Msg_StartHuntTreasure_Ack struct {
	RetCode       int
	CurrentPos    int            //! 移动后所在格ID
	MoveTypeLst   []int          //! 移动到过的格子
	ExMove        []int          //! 倘若出现移动事件,此参数代表移动格子
	RandomScore   int            //! 随机点数
	Score         int            //! 总积分
	TodayRank     int            //! 今日排名变化
	TotalRank     int            //! 累计排名变化
	HuntTurn      int            //! 巡回轮数
	AwardItem     []MSG_ItemData //! 获取奖励
	CostItem      MSG_ItemData   //! 消耗物品
	CostFreeTimes int            //! 消耗免费次数
}

//! 玩家查询巡回奖励领取情况
//! 消息: /query_hunt_award
type MSG_QueryHuntTurn_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryHuntTurn_Ack struct {
	RetCode   int
	AwardMask int //! 进行位运算 1 0 表示已领取或者未领取
}

//! 玩家领取巡回奖励
//! 消息: /get_hunt_award
type MSG_GetHuntTurnAward_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_GetHuntTurnAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 获取运营活动排行榜
//! 消息: /get_activity_rank
type MSG_GetActivityRank_Req struct {
	PlayerID   int32
	SessionKey string
	Type       int //! 按照Activity_Type取值
}

type MSG_GetActivityRank_Ack struct {
	RetCode          int
	TodayRankLst     []MSG_OperationalActivityRank
	YesterdayRankLst []MSG_OperationalActivityRank
	TotalRankLst     []MSG_OperationalActivityRank

	ScoreLst [3]int //! 玩家自身分数  0->昨天 1->今天 2->公共的
	RankLst  [3]int //! 玩家自身排名

	EliteScore           int  //! 进入精英榜最低要求积分
	IsRecvTodayRankAward bool //! 是否领取今日奖励
	IsRecvTotalRankAward bool //! 是否领取总排行奖励
}

//! 玩家领取活动排行榜奖励
//! 消息: /get_activity_rank_award
type MSG_GetActivityRankAward_Req struct {
	PlayerID     int32
	SessionKey   string
	ActivityType int //! 按照Activity_Type取值
	AwardType    int //! 1->昨日榜 2->总榜
}

type MSG_GetActivityRankAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 玩家查询排行榜
//! 消息: /get_hunt_rank
type MSG_GetHuntTreasureRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_OperationalActivityRank struct {
	PlayerID   int32
	PlayerName string
	HeroID     int
	Quality    int8
	Level      int
	Score      int
}

type MSG_GetHuntTreasureRank_Ack struct {
	RetCode          int
	TodayRankLst     []MSG_OperationalActivityRank
	YesterdayRankLst []MSG_OperationalActivityRank
	TotalRankLst     []MSG_OperationalActivityRank

	ScoreLst [3]int //! 玩家自身分数  0->昨天 1->今天 2->公共的
	RankLst  [3]int //! 玩家自身排名

	EliteScore           int  //! 进入精英榜最低要求积分
	IsRecvTodayRankAward bool //! 是否领取今日奖励
	IsRecvTotalRankAward bool //! 是否领取总排行奖励
}

//! 玩家领取排行榜奖励
//! 消息: /get_hunt_rank_award
type MSG_GetHuntRankAward_Req struct {
	PlayerID   int32
	SessionKey string
	RankType   int //! 1->昨日排行奖励 2->累积排行奖励
}

type MSG_GetHuntRankAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 玩家查询巡回商店
//! 消息: /query_hunt_store
type MSG_QueryHuntStore_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_HuntStoreItem struct {
	ID       int
	ItemID   int
	ItemNum  int
	MoneyID  int
	MoneyNum int
	Score    int  //! 购买物品所获积分
	IsBuy    bool //! 是否已经购买
}

type MSG_QueryHuntStore_Ack struct {
	RetCode  int
	GoodsLst []MSG_HuntStoreItem
}

//! 玩家购买巡回商店物品
//! 消息: /buy_hunt_store
type MSG_BuyHuntStoreItem_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_BuyHuntStoreItem_Ack struct {
	RetCode int
}

//! 清除商店物品
//! 消息: /clean_hunt_store
type MSG_CleanHuntStore_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_CleanHuntStore_Ack struct {
	RetCode int
}

//! 玩家查询转盘信息
//! 消息: /query_lucky_wheel
type MSG_QueryLuckyWheel_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_LuckyWheelAward struct {
	ItemID    int
	ItemNum   int
	IsSpecial int //! 1 表示是  0 表示不是
}

type MSG_QueryLuckyWheel_Ack struct {
	RetCode int

	NormalMoneyPoor  int //! 普通奖金池
	ExcitedMoneyPoor int //! 豪华奖金池

	NormalAwardLst  []MSG_LuckyWheelAward //! 普通轮盘奖励
	ExcitedAwardLst []MSG_LuckyWheelAward //! 豪华轮盘奖励

	NormalFreeTimes  int //! 普通轮盘免费次数
	ExcitedFreeTimes int //! 豪华轮盘免费次数

	CostMoneyID  [4]int
	CostMoneyNum [4]int

	Score        int //! 今日积分
	EndCountDown int //! 活动结束倒计时
	TodayRank    int //! 今日排名变化
	TotalRank    int //! 累计排名变化
}

//! 玩家转动轮盘
//! 消息: /rotating_wheel
type MSG_RotatingWheel_Req struct {
	PlayerID        int32
	SessionKey      string
	IsExcited       int //! 1 -> 普通 2-> 豪华
	IsStartTenTimes int //! 是否十连转 1 -> 是  0 -> 否
}

type MSG_RotatingWheel_Ack struct {
	RetCode    int
	AwardItem  []MSG_ItemData
	AwardIndex int

	Score int //! 积分

	NormalMoneyPoor  int //! 普通奖金池
	ExcitedMoneyPoor int //! 豪华奖金池

	TodayRank int //! 今日排名变化
	TotalRank int //! 累计排名变化

	CostItem      MSG_ItemData
	CostFreeTimes int

	MoneyNum int //! 货币数量
}

//! 玩家查询排行榜
//! 消息: /query_lucky_wheel_rank
type MSG_QueryLuckyWheelRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryLuckyWheelRank_Ack struct {
	RetCode          int
	TodayRankLst     []MSG_OperationalActivityRank
	YesterdayRankLst []MSG_OperationalActivityRank
	TotalRankLst     []MSG_OperationalActivityRank

	ScoreLst [3]int //! 玩家自身分数
	RankLst  [3]int //! 玩家自身排名

	EliteScore           int  //! 进入精英榜最低要求积分
	IsRecvTodayRankAward bool //! 是否领取今日奖励
	IsRecvTotalRankAward bool //! 是否领取总排行奖励
}

//! 玩家领取排行榜奖励
//! 消息: /get_wheel_rank_award
type MSG_GetLuckyWheelRankAward_Req struct {
	PlayerID   int32
	SessionKey string
	RankType   int //! 1->昨日排行奖励 2->累积排行奖励
}

type MSG_GetLuckyWheelRankAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 玩家请求团购信息
//! 消息: /get_group_purchase_info
type MSG_GetGroupPurchaseInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GroupPurchase struct {
	ItemID    int //! 物品ID
	SaleNum   int //! 销售数量
	CanBuyNum int //! 还可购买次数
}

type MSG_GetGroupPurchaseInfo_Ack struct {
	RetCode        int
	AwardType      int
	ItemInfo       []MSG_GroupPurchase
	Score          int   //! 总积分
	EndTime        int64 //! 结束
	AwardTime      int64 //! 领奖
	ScoreAwardMark []int
}

//! 玩家请求团购
//! 消息: /buy_group_purchase
type MSG_BuyGroupPurchase_Req struct {
	PlayerID   int32
	SessionKey string
	ItemID     int
}

type MSG_BuyGroupPurchase_Ack struct {
	RetCode      int
	CostItemID   int //! 使用团购券ID
	CostItemNum  int //! 使用团购券数量
	CostMoneyID  int //! 使用货币ID
	CostMoneyNum int //! 使用货币数量
	ItemID       int //! 获取道具ID
	ItemNum      int //! 获取道具数量
	SaleNum      int //! 销售数量
	Score        int //! 总积分
}

//! 玩家请求查询积分奖励
//! 消息: /get_group_purchase_score
type MSG_GetGroupPurchaseScore_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGroupPurchaseScore_Ack struct {
	RetCode int
}

//! 玩家请求获取积分奖励
//! 消息: /get_group_score_award
type MSG_GetGroupScoreAward_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_GetGroupScoreAward_Ack struct {
	RetCode int
}

//! 查询欢庆佳节活动信息
//! 消息: /get_festival_task
type MSG_GetFestivalTask_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FestivalTask struct {
	ID       int
	CurCount int //! 当前完成次数
	Status   int //! 当前任务状态: 0->未完成 1->已完成 2->已领取
}

type MSG_FestivalSale struct {
	ID    int
	Times int
}

type MSG_FestivalExchange struct {
	ID          int //! ID
	NeedItemID  int //! 所需物品
	NeedItemNum int //!
	Award       int //! 兑换物品
	Times       int //! 剩余兑换次数
}

type MSG_GetFestivalTask_Ack struct {
	RetCode           int
	TaskLst           []MSG_FestivalTask
	BuyLst            []MSG_FestivalSale
	ExchangeRecordLst []MSG_FestivalExchange
}

//! 领取欢庆佳节活动任务奖励
//! 消息: /get_festival_task_award
type MSG_GetFestivalTaskAward_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_GetFestivalTaskAward_Ack struct {
	RetCode  int
	AwardLst []MSG_ItemData
}

//! 兑换欢庆佳节活动奖励
//! 消息: /exchange_festival_award
type MSG_ExchangeFestivalAward_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_ExchangeFestivalAward_Ack struct {
	RetCode  int
	AwardLst []MSG_ItemData
}

//! 欢庆佳节半折贩售
//! 消息: /festival_discount_sale
type MSG_BuyFestivalSale_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_BuyFestivalSale_Ack struct {
	RetCode      int
	ItemID       int
	ItemNum      int
	CostMoneyID  int
	CostMoneyNum int
}

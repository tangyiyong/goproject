package msg

//! 请求获取当天开启活动
//! 消息: /get_activity 或者 /get_activity_list
type MSG_GetActivity_Req struct {
	PlayerID   int
	SessionKey string
}

type TActivityDiscount struct {
	Index    int //! 索引
	BuyTimes int //! 剩余购买次数
}

type MSG_ActivityInfo struct {
	ID        int
	Icon      int   //! 活动图标
	Type      int   //! 活动套用模板
	AwardType int   //! 活动套用奖励
	RedTip    bool  //! 是否存在操作提示
	BeginTime int64 //! 开始时间
	EndTime   int64 //! 结束时间
	AwardTime int   //! 领奖时间
	IsInside  int   //! 是否在里面
}

type MSG_GetActivity_Ack struct {
	RetCode            int
	ActivityLst        []MSG_ActivityInfo
	RemoveActivityIcon []int //! 需删除的图标
}

//! 查询累计登录活动信息
//! 消息: /query_activity_login
type MSG_QueryActivity_Login_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 活动ID
}

type MSG_QueryActivity_Login_Ack struct {
	//! 登录活动信息
	RetCode    int
	ActivityID int //! 活动ID
	AwardType  int //! 活动奖励
	LoginDay   int //! 可领取天数
	AwardMark  int //! 进行位运算判断领取标记
}

//! 查询首冲活动信息
//! 消息: /query_activity_firstrecharge
type MSG_QueryActivity_FirstRecharge_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_QueryActivity_FirstRecharge_Ack struct {
	//! 首充活动状态
	RetCode             int
	FirstRechargeStatus int //! 0->不能领取 1->可以领取 2->已领取首充并开启次充奖励 3->次充奖励可领取 4->已领取次充奖励
}

//! 查询领取体力活动信息
//! 消息: /query_activity_action
type MSG_QueryActivity_Action_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_QueryActivity_Action_Ack struct {
	//! 领体力
	RetCode                 int
	RecvAction              int //! 进行位运算
	NextAwardTime           int
	RetroactiveCostMoneyID  int //! 补签消耗货币
	RetroactiveCostMoneyNum int
}

//! 查询充值回馈活动信息(累计充值)
//! 消息: /query_activity_totalrecharge
type MSG_QueryActivity_TotalRecharge_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 活动ID
}

type MSG_QueryActivity_TotalRecharge_Ack struct {
	//! 充值回馈
	RetCode     int
	ActivityID  int
	AwardType   int
	RechargeNum int //! 活动期间累积充值数额
	AwardMark   int //! 累积充值领取标记 (索引 按位运算)
}

//! 查询月卡天数
//! 消息: /query_monthcard_days
type MSG_QueryActivity_MonthCard_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_QueryActivity_MonthCard_Ack struct {
	RetCode int
	Days    []int  //! 若为0则是没有该月卡
	Status  []bool //! true为已经领取  false为未领取
}

//! 查询单笔充值活动信息
//! 消息: /query_activity_singlerecharge
type MSG_QueryActivity_SingleRecharge_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
}

type MSG_QueryActivity_SingleRecharge_Ack struct {
	//! 单笔充值情况
	RetCode           int
	ActivityID        int
	AwardType         int
	SingleRechargeLst []MSG_SingleRecharge
}

type MSG_SingleRecharge struct {
	Index  int //! 索引 从1开始
	Times  int //! 当前剩余次数
	Status int //! 0 不可领  1 可领取   2  已领取
}

//! 请求领取登录奖励
//! 消息: /get_login_award
type MSG_GetActivity_LoginAward_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 任务ID
	Index      int //! 奖励索引 从1开始
	Choice     int //! 若为三选一奖励,则该字段用于发送选择奖励索引
}

type MSG_GetActivity_LoginAward_Ack struct {
	AwardItem []MSG_ItemData
	AwardMark int
	RetCode   int
}

//! 请求领取首充奖励
//! 消息: /get_first_recharge
type MSG_GetActivity_FirstRecharge_Req struct {
	PlayerID     int
	SessionKey   string
	GetAwardType int //! 1-> 首充  2->次充
}

type MSG_GetActivity_FirstRecharge_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 请求领取体力
//! 消息: /get_activity_action
type MSG_GetActivity_Action_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetActivity_Action_Ack struct {
	RetCode       int
	AwardItem     []MSG_ItemData
	ActionValue   int   //! 行动力值
	ActionTime    int64 //! 行动力恢复起始时间
	NextAwardTime int   //! 下次领取倒计时
	Index         int
}

//! 领取体力补签协议
//! 消息: /get_action_retroactive
type MSG_GetAction_Retroactive_Req struct {
	PlayerID   int
	SessionKey string
	Index      int //! 从1开始  1-4分别代表四个时间段
}

type MSG_GetAction_Retroactive_Ack struct {
	RetCode       int
	AwardItem     []MSG_ItemData
	ActionValue   int   //! 行动力值
	ActionTime    int64 //! 行动力恢复起始时间
	NextAwardTime int   //! 下次领取倒计时
	CostItem      []MSG_ItemData
	Index         int
}

//! 查询迎财神活动信息
//! 消息: /query_activity_moneygod
type MSG_QueryActivity_MoneyGod_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_QueryActivity_MoneyGod_Ack struct {
	//! 迎财神
	RetCode         int
	CurrentTimes    int   //! 当前剩余领取次数
	TotalMoney      int   //! 累积银币
	NextTime        int64 //! 下次迎财神时间
	CumulativeTimes int   //! 累积连续领取次数
}

//! 玩家请求迎财神
//! 消息: /welcome_money_god
type MSG_WelcomeMoneyGod_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_WelcomeMoneyGod_Ack struct {
	RetCode    int
	TotalMoney int   //! 累积金钱
	NextTime   int64 //! 下次可迎财神时间
	ExAwardID  int   //! 若此ID不为零,则表示获取了额外的物品奖励
	ExAwardNum int
	MoneyID    int //! 获取的银币
	MoneyNum   int
}

//! 玩家请求迎财神活动中累积奖励
//! 消息: /get_money_god_award
type MSG_GetMoneyGodTotalAward_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetMoneyGodTotalAward_Ack struct {
	RetCode  int
	MoneyID  int
	MoneyNum int
}

//! 查询折扣贩售活动信息
//! 消息: /query_activity_discountsale
type MSG_QueryActivity_DisountSale_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 活动ID
}

type MSG_QueryActivity_DisountSale_Ack struct {
	//! 折扣贩售
	RetCode    int
	ActivityID int
	AwardType  int
	ShopLst    []TActivityDiscount //! 已购买物品信息
}

//! 玩家请求购买折扣贩售商品
//! 消息: /buy_discount_sale
type MSG_BuyDiscountItem_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 活动ID
	Index      int //! 商品索引,从1开始
	Choice     int //! 若为多选一奖励,则该字段用于发送选择奖励索引
	Count      int //! 购买数量
}

type MSG_BuyDiscountItem_Ack struct {
	RetCode    int
	ActivityID int
	MoneyID    int
	MoneyNum   int
	AwardItem  []MSG_ItemData
	Index      int
	BuyNum     int
}

//! 玩家请求领取充值回馈奖励
//! 消息: /get_recharge_award
type MSG_GetRechargeAward_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 活动ID  避免存在多个充值回馈活动
	Index      int
}

type MSG_GetRechargeAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
	AwardMark int
}

//! 玩家请求领取单充回馈奖励
//! 消息: /get_single_award
type MSG_GetSingleAward_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int //! 活动ID 避免存在多个单笔充值
	Index      int
}

type MSG_GetSingleAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 玩家获取签到信息
//! 消息:/get_sign
type MSG_GetSignInfo_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetSignInfo_Ack struct {
	RetCode        int
	SignDay        int  //! 当前签到天数
	IsSign         bool //! 普通签到状态 false -> 未签  true -> 已签
	SignPlusStatus int  //! 豪华签到状态 0-> 不可领取  1-> 可领取  2-> 已领取

	SignIndex     int            //! 普通签到奖励索引
	SignPlusAward []MSG_ItemData //! 豪华签到奖励
}

//! 玩家请求日常签到
//! 消息: /daily_sign
type MSG_DailySign_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_DailySign_Ack struct {
	RetCode   int
	ItemID    int
	ItemNum   int
	AwardType int //! 奖励模板
}

//! 玩家请求豪华签到
//! 消息: /sign_plus
type MSG_PlusSign_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_PlusSign_Ack struct {
	RetCode   int
	AwardInfo []MSG_ItemData
}

//! 请求查询购买基金状态
//! 消息: /get_fund_status
type MSG_GetOpenFundStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetOpenFundStatus_Ack struct {
	RetCode       int
	IsBuy         bool //! 玩家是否购买基金
	BuyNum        int  //! 购买基金人数
	FundLevelMark int  //! 等级奖励领取标记
	FundCountMark int  //! 购买基金人数奖励领取
	CostMoneyID   int  //! 购买基金花费
	CostMoneyNum  int  //! 购买基金花费
	ReceiveMoney  int  //! 剩余领取钻石
}

//! 请求购买开服基金
//! 消息: /buy_fund
type MSG_BuyFund_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_BuyFund_Ack struct {
	RetCode int
	BuyNum  int //! 购买基金人数
}

//! 请求领取基金奖励-全服奖励
//! 消息: /get_fund_all_award
type MSG_ReceiveFundAllAward_Req struct {
	PlayerID   int
	SessionKey string
	ID         int //! 奖励ID
}

type MSG_ReceiveFundAllAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
	Index     int
}

//! 请求领取基金奖励-等级返利
//! 消息: /get_func_level_award
type MSG_ReceiveFundLevelAward_Req struct {
	PlayerID   int
	SessionKey string
	ID         int
}

type MSG_ReceiveFundLevelAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
	Index     int
}

//! 获取限时日常任务信息
//! 消息: /get_limit_daily_task
type MSG_GetLimitDailyTask_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
}

//! 限时日常任务类型
type TLimitDailyTask struct {
	Index  int //! 任务类型
	Count  int //! 当前次数
	Status int //! 状态: 0->未完成 1->已完成 2->已领取
}

type MSG_GetLimitDailyTask_Ack struct {
	RetCode int
	TaskLst []TLimitDailyTask
}

//! 领取限时日常任务奖励
//! 消息: /get_limit_daily_award
type MSG_GetLimitDailyAward_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
	Index      int //! 从0开始
	Select     int //! 若is_select字段不为0, 则只能从奖励中选取其中之一, 下标从1开始
}

type MSG_GetLimitDailyAward_Ack struct {
	RetCode   int
	AwardItem []MSG_ItemData
}

//! 卡牌大师
type MSG_CardMaster_CardList_Req struct { // 消息：/act_card_master_card_list
	PlayerID   int
	SessionKey string
}
type MSG_CardMaster_CardList_Ack struct {
	RetCode       int
	FreeTimes     byte
	Score         int
	Point         int
	Cards         []MSG_ItemData
	ExchangeTimes []MSG_ItemData
}
type MSG_CardMaster_Draw_Req struct { // 消息：/act_card_master_draw
	PlayerID   int
	SessionKey string
	Type       byte // 1：普通抽、2：普通十连、3：高级抽、4：高级十连
}
type MSG_CardMaster_Draw_Ack struct {
	RetCode int
	Cards   []MSG_ItemData
}
type MSG_CardMaster_Card2Item_Req struct { // 消息：/act_card_master_card2item
	PlayerID   int
	SessionKey string
	ExchangeID int
}
type MSG_CardMaster_Card2Item_Ack struct {
	RetCode int
}
type MSG_CardMaster_Card2Point_Req struct { // 消息：/act_card_master_card2point
	PlayerID   int
	SessionKey string
	Cards      []MSG_ItemData
}
type MSG_CardMaster_Card2Point_Ack struct {
	RetCode int
	Point   int
}
type MSG_CardMaster_Point2Card_Req struct { // 消息：/act_card_master_point2card
	PlayerID   int
	SessionKey string
	Cards      []MSG_ItemData
}
type MSG_CardMaster_Point2Card_Ack struct {
	RetCode int
	Point   int
}

//! 月光集市
type MSG_MoonlightShop_GetInfo_Req struct { // 消息：/act_moonlight_shop_get_info
	PlayerID   int
	SessionKey string
}
type MSG_MoonlightShop_ExchangeToken_Req struct { // 消息：/act_moonlight_shop_exchangetoken
	PlayerID   int
	SessionKey string
	ExchangeID int
}
type MSG_MoonlightShop_ExchangeToken_Ack struct {
	RetCode int
}
type MSG_MoonlightShop_ReduceDiscount_Req struct { // 消息：/act_moonlight_shop_reducediscount
	PlayerID   int
	SessionKey string
	GoodsID    int
}
type MSG_MoonlightShop_ReduceDiscount_Ack struct {
	RetCode  int
	Discount byte
}
type MSG_MoonlightShop_RefreshShop_Buy_Req struct { // 消息：/act_moonlight_shop_refreshshop_buy
	PlayerID   int
	SessionKey string
}
type MSG_MoonlightShop_RefreshShop_Auto_Req struct { // 消息：/act_moonlight_shop_refreshshop_auto
	PlayerID   int
	SessionKey string
}
type MSG_MoonlightShop_BuyGoods_Req struct { // 消息：/act_moonlight_shop_buygoods
	PlayerID   int
	SessionKey string
	GoodsID    int
}
type MSG_MoonlightShop_BuyGoods_Ack struct {
	RetCode int
}
type MSG_MoonlightShop_GetScoreAward_Req struct { // 消息：/act_moonlight_shop_getscoreaward
	PlayerID   int
	SessionKey string
	AwardID    int
}
type MSG_MoonlightShop_GetScoreAward_Ack struct {
	RetCode int
}

//! 沙滩宝贝
type MSG_BeachBaby_Info_Req struct { // 消息：/act_beach_baby_info
	PlayerID   int
	SessionKey string
}
type MSG_BeachBaby_OpenGoods_Req struct { // 消息：/act_beach_baby_open_goods
	PlayerID   int
	SessionKey string
	Index      int
}
type MSG_BeachBaby_OpenGoods_Ack struct {
	RetCode   int
	Item      MSG_ItemData
	IsGetItem bool
}
type MSG_BeachBaby_OpenAllGoods_Req struct { // 消息：/act_beach_baby_open_all_goods
	PlayerID   int
	SessionKey string
}
type MSG_BeachBaby_Refresh_Auto_Req struct { // 消息：/act_beach_baby_refresh_auto
	PlayerID   int
	SessionKey string
}
type MSG_BeachBaby_Refresh_Buy_Req struct { // 消息：/act_beach_baby_refresh_buy
	PlayerID   int
	SessionKey string
}
type MSG_BeachBaby_GetFreeConch_Req struct { // 消息：/act_beach_baby_get_freeconch
	PlayerID   int
	SessionKey string
}
type MSG_BeachBaby_GetFreeConch_Ack struct {
	RetCode int
}
type MSG_BeachBaby_SelectGoodsID_Req struct { // 消息：/act_beach_baby_select_goods
	PlayerID   int
	SessionKey string
	IDs        []int
}
type MSG_BeachBaby_SelectGoodsID_Ack struct {
	RetCode int
}

//! 玩家请求VIP日常福利领取状态
//! 消息: /query_vip_welfare
type MSG_GetVipDailyWelfareStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_WeekGiftInfo struct {
	ID       int
	BuyTimes int
}

type MSG_GetVipDailyWelfareStatus_Ack struct {
	RetCode    int                //! 返回码
	SignStatus int                //! 0->不可领取 1->未领取 2->已领取
	GiftLst    []MSG_WeekGiftInfo //! 每周礼包信息
}

//! 玩家请求VIP日常福利
//! 消息: /get_vip_welfare
type MSG_GetVipDailyWelfare_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetVipDailyWelfare_Ack struct {
	RetCode   int //! 返回码
	AwardItem []MSG_ItemData
}

//! 玩家请求VIP每周礼包信息
//! 消息: /get_vip_week_gift
type MSG_GetVipWeekGiftInfo_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetVipWeekGiftInfo_Ack struct {
	RetCode int
}

//! 玩家请求购买VIP每周礼包
//! 消息: /buy_vip_week_gift
type MSG_BuyVipWeekGiftInfo_Req struct {
	PlayerID   int
	SessionKey string
	ID         int
	BuyTimes   int
}

type MSG_BuyVipWeekGiftInfo_Ack struct {
	RetCode   int
	MoneyID   int
	MoneyNum  int
	AwardItem []MSG_ItemData
	BuyTimes  int
	ID        int
}

//! 玩家请求周周盈状态
//! 消息: /get_week_award_status
type MSG_GetWeekAwardStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetWeekAwardStatus_Ack struct {
	RetCode     int //! 返回码
	LoginDay    int //! 登录天数
	RechargeNum int //! 充值数目
	AwardMark   int //! 奖励标记 包含已领取奖励ID
}

//! 玩家请求领取周周盈奖励
//! 消息: /get_week_award
type MSG_GetWeekAward_Req struct {
	PlayerID   int
	SessionKey string
	Index      int //! Index  从1开始
	Select     int //! 选择哪个奖励
}

type MSG_GetWeekAward_Ack struct {
	RetCode   int
	AwardMark int //! 奖励标记 包含已领取奖励ID
	AwardItem []MSG_ItemData
}

//! 玩家请求等级礼包信息
//! 消息: /get_level_gift_info
type MSG_GetLevelGiftInfo_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_LevelGiftInfo struct {
	ID       int
	BuyTimes int   //! 当前可购买次数
	DeadLine int64 //! 过期时间
}

type MSG_GetLevelGiftInfo_Ack struct {
	RetCode int
	GiftLst []MSG_LevelGiftInfo
}

//! 玩家请求购买等级礼包
//! 消息: /buy_level_gift
type MSG_BuyLevelGift_Req struct {
	PlayerID   int
	SessionKey string
	GiftID     int
}

type MSG_BuyLevelGift_Ack struct {
	RetCode      int
	AwardItem    []MSG_ItemData
	CostMoneyID  int
	CostMoneyNum int
	BuyTimes     int //! 剩余可购买次数
}

//! 玩家请求查询月基金状态
//! 消息: /get_monthfund_status
type MSG_GetMonthFundStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetMonthFundStatus_Ack struct {
	RetCode    int
	Day        int  //! 还剩领取天数
	IsReceived bool //! 今天是否已领取
	CountDown  int  //! 截止购买时间
	MoneyID    int  //! 购买花费货币ID
	MoneyNum   int  //! 购买需花费货币数量
}

//! 玩家请求领取月基金奖励
//! 消息: /receive_month_fund
type MSG_ReceiveMonthFund_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_ReceiveMonthFund_Ack struct {
	RetCode  int
	AwardLst []MSG_ItemData
}

//! 玩家请求等级礼包信息
//! 消息: /get_rank_gift_info
type MSG_GetRankGiftInfo_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_RankGiftInfo struct {
	ID       int
	BuyTimes int //! 当前可购买次数
}

type MSG_GetRankGiftInfo_Ack struct {
	RetCode int
	GiftLst []MSG_RankGiftInfo
	Rank    int //! 历史最高名次
}

//! 玩家请求购买等级礼包
//! 消息: /buy_rank_gift
type MSG_BuyRankGift_Req struct {
	PlayerID   int
	SessionKey string
	GiftID     int
}

type MSG_BuyRankGift_Ack struct {
	RetCode      int
	AwardItem    []MSG_ItemData
	CostMoneyID  int
	CostMoneyNum int
	BuyTimes     int //! 剩余可购买次数
}

//! 玩家查询限时特惠物品信息
//! 消息: /get_limit_sale_info
type MSG_GetLimitSaleInfo_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_LimitSaleItemInfo struct {
	ID     int
	Status bool
}

type MSG_GetLimitSaleInfo_Ack struct {
	RetCode   int
	Score     int //! 玩家积分
	ItemLst   []MSG_LimitSaleItemInfo
	AwardMark int
}

//! 玩家购买限时特惠
//! 消息: /buy_limit_sale_item
type MSG_BuyLimitSaleItem_Req struct {
	PlayerID   int
	SessionKey string
	Index      int //! 索引, 从1开始
}

type MSG_BuyLimitSaleItem_Ack struct {
	RetCode  int
	AwardLst []MSG_ItemData
	Score    int //! 玩家积分
}

//! 玩家请求领取全民奖励
//! 消息: /get_limitsale_all_award
type MSG_GetLimitSale_AllAward_Req struct {
	PlayerID   int
	SessionKey string
	ID         int
}

type MSG_GetLimitSale_AllAward_Ack struct {
	RetCode   int
	AwardMark int
	AwardItem []MSG_ItemData
}

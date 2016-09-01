package gamedata

import (
	"fmt"
	// "gamelog"
	"strings"
)

var (
	BeachBaby_Token_ItemID      int    // 沙滩贝壳
	BeachBaby_CostMoneyType     int    // 翻贝壳消耗的货币类型
	BeachBaby_GetFreeToken_Time []byte // 8|12|16|20 免费领取	len <= 8  用int8位标记的
	BeachBaby_GetFreeToken_Cnt  byte
	BeachBaby_OpenGoods_Cost    []byte // 次数-价格    len() == 16
	BeachBaby_Refresh_CD        byte   // 商品刷新间隔：分钟
	BeachBaby_Refresh_Cost      byte   // 购买直接刷新商品
	BeachBaby_SelectGoods_Cost  []byte // 自选所需钻石数，编号表示次数
	BeachBaby_TodayRank_Limit   int    // 单日榜 特殊奖门槛
	BeachBaby_TotalRank_Limit   int    // 累计榜 特殊奖门槛

	MoonlightShop_Token_ItemID      int //月光币
	MoonlightShop_Discount_Cost     []byte
	MoonlightShop_Discount_OneTiems []byte //一次打折的范围
	MoonlightShop_Shop_Refresh_CD   byte   //商品刷新间隔：分钟
	MoonlightShop_Shop_Refresh_Cost byte   //购买直接刷新商品
	MoonlightShop_BuyTimes_Max      int    //单日最大购买次数

	CultureItemID  int //培养道具ID
	CultureItemNum int //培养道具数量

	//属性
	AttackPropertyID  int //攻击属性ID
	AttackPhysicID    int //物攻属性ID
	AttackMagicID     int //魔攻属性ID
	DefencePropertyID int //防御属性ID
	DefencePhysicID   int //物防属性ID
	DefenceMagicID    int //魔防属性ID
	AllPropertyID     int //全属性ID

	//商店
	StoreFreeRefreshTimes          int //商店可免费刷新次数上限
	StoreFreeRefreshAddTime        int //商店刷新免费次数刷新时间
	HeroStoreRefreshNeedMoneyType  int //神将商店刷新需要的货币ID
	HeroStoreRefreshNeedMoneyNum   int //神将商店刷新需要货币数量
	AwakeStoreRefreshNeedMoneyType int //觉醒商店刷新需要的货币ID
	AwakeStoreRefreshNeedMoneyNum  int //觉醒商店刷新需要货币数量
	PetStoreRefreshNeedMoneyType   int //战宠商店刷新需要的货币ID
	PetStoreRefreshNeedMoneyNum    int //战宠商店刷新需要货币数量
	StoreRefreshNeedItem           int //商店刷新需要的道具
	StoreRefreshItemNum            int //商店刷新需要道具的个数

	//名将副本
	FamousCopyChallengeTimes int //名将副本每天固定挑战次数

	//! 各类功能使用货币ID
	MallItemMoneyID       int //道具商城使用货币ID
	MallGiftMoneyID       int //礼包商城使用货币ID
	EliteCopyResetMoneyID int //精英副本重置使用货币ID
	VipWeeklyGiftMoneyID  int //VIP每周礼包使用货币ID

	//首充礼包ID
	FirstRechargeAwardID  int //首充礼包ID
	NextRechargeAwardID   int //次充礼包ID
	NextAwardNeedRecharge int //次充礼包需求充值额度

	//竞技场
	ArenaRankAwardCalcTime int //竞技场结算时间(整点)

	//三国无双
	SangokuMusouEliteFreeTimes      int //三国无双精英挑战每天免费次数
	SangokuMusouResetCopyMoneyID    int //三国无双重置普通挑战货币ID
	SangokuMusouAddEliteCopyMoneyID int //三国无双增加精英挑战次数货币ID
	CritMultiple                    int //普通暴击倍数
	BigCritMultiple                 int //大暴击倍数
	LuckyCritMultiple               int //幸运暴击倍数
	CritPro                         int //普通暴击概率
	BigCritPro                      int //大暴击概率
	LuckyCritPro                    int //幸运暴击概率

	//领地攻讨
	SuppressRiotAwardItem       int //帮忙好友镇压暴动奖励
	SuppressRiotAwardItemNum    int //帮忙好友镇压暴动奖励数量
	SuppressRiotFriendAwardItem int //好友帮忙镇压暴动奖励
	SuppressRiotFriendAwardNum  int //好友帮忙镇压暴动奖励数量
	RiotPro                     int //暴动概率
	RiotTime                    int //暴动持续时间

	//武将召唤
	NormalSummonFreeTimes  int //普通武将免费召唤次数
	NormalSummonFreeCDTime int //普通武将免费召唤间隔时间
	SeniorSummonFreeCDTime int //高级武将免费抽取时间
	TenSummonDiscount      int //十连抽折扣
	SeniorSummonPoint      int //高级武将每次增加积分

	//围剿叛军
	FindRebelPro                   int //发现叛军概率
	LowerRebelPro                  int //低级叛军概率
	MiddleRebelPro                 int //中极叛军概率
	SeniorRebelPro                 int //高级叛军概率
	RebelEscapeTime                int //叛军逃跑时间
	RebelAchievements              int //围剿叛军战功ID
	AttackRebelActionID            int //围剿叛军所需行动力ID
	NormalAttackRebelNeedActionNum int //普通围剿叛军需求行动力数量
	SeniorAttackRebelNeedActionNum int //高级围剿叛军需求道具数量

	//挖矿系统
	MiningActionRecoverLimit int //挖矿行动力恢复上限
	MiningActionRecoverTime  int //挖矿行动力恢复时间
	MiningAttackBossAction   int //挖矿每次攻击Boss消耗行动力
	MiningCostActionNum      int //挖矿消耗行动力数量
	MiningCostActionID       int //挖矿消耗行动力ID
	MiningMapLength          int //矿洞地图边长
	MiningMapSlice           int //挖矿地图切分
	MiningResetMoneyID       int //挖矿重置所需货币ID
	MiningEnterPointX        int //挖矿入口点坐标X
	MiningEnterPointY        int //挖矿入口点坐标Y
	MiningBossValue          int //boss 积分

	//精英关卡
	EliteInvadeTime1 int //! 精英关卡入侵时间
	EliteInvadeNum1  int //! 精英关卡入侵叛军数量
	EliteInvadeTime2 int
	EliteInvadeNum2  int
	EliteInvadeTime3 int
	EliteInvadeNum3  int
	EliteInvadeTime4 int
	EliteInvadeNum4  int

	ChargeMoneyID int //充值货币ID
	VipExpMoneyID int //VIP经验的货币ID

	SevenActivityAwardDay int //! 七日活动领奖持续天数

	OpenFundPriceID  int //! 开服基金购买需求货币ID
	OpenFundPriceNum int //! 开服基金购买需求货币数量

	//! 围剿叛军战功换算比例
	RebelExploitPoint int

	//挂机
	HangUpInitGridNum    int //挂机背包初始格子数
	HangUpQuickFight     int //一次快速战斗的事件次数
	HangUpOpenGridNum    int //单次开放格子数
	HangUpBuyGridMoneyID int //开放格子数所使用的货币

	//! 工会系统
	GuildActionRecoverTime   int //! 公会副本行动力回复时间
	GuildCopyBattleTimeBegin int //! 公会副本开启时间
	GuildCopyBattleTimeEnd   int //! 工会副本结束时间
	CreateGuildMoneyID       int //! 公会需要货币ID
	CreateGuildMoneyNum      int //! 公会需要货币数量
	UpdateGuildNameMoneyID   int //! 修改名字所需货币ID
	UpdateGuildNameMoneyNum  int //! 修改名字所需货币数量
	GuildSacrificeCrit       int //! 军团祭天暴击概率

	//黑市
	BlackMarketRefreshTime []int //! 黑市刷新时间
	EnterVipLevel          int   //! 入口VIP等级需求
	BlackMarketPro         int   //! 黑市入口出现概率

	//公会技能学习需求货币ID
	GuildSKillStudyNeedMoneyID int

	//积分赛
	OneTimeFightScore     int //单场战斗得失分
	ScoreBuyTimeMoneyID   int //购买战斗次数货币
	ScoreCopyID           int //副本ID
	ScoreMoneyID          int //积分赛货币
	ScoreMoneyNum         int //单场战斗得失货币
	ScoreSeriesWinAwardID int //连胜奖励ID
	ScoreSeriesWinTimes   int //连胜次数

	//名人堂免费送花次数
	FameHallFreeTimes int

	//夺粮战
	FoodWarOpenTime            int   //! 攻击时间
	FoodWarEndTime             int   //! 结束时间
	FoodWarAttackTimes         int   //! 初始攻击次数
	FoodWarRevengeTimes        int   //! 初始复仇次数
	FoodWarFixedFood           int   //! 初始固定粮草
	FoodWarNonFixedFood        int   //! 初始流动粮草
	FoodWarOpenDay             []int //! 活动开启时间
	FoodWarTimeAddFood         int   //! 每小时增加固定粮草数量
	FoodWarRobBili             int   //! 抢夺粮草比例
	FoodWarCopyID              int   //! 夺粮战副本ID
	FoodWarVictoryMoneyID      int   //! 胜利货币ID
	FoodWarVictoryMoneyNum     int   //! 胜利货币数量
	FoodWarFailedMoneyID       int   //! 失败货币ID
	FoodWarFailedMoneyNum      int   //! 失败货币数量
	FoodWarBuyTimesNeedMoneyID int   //! 夺粮战购买次数需求货币ID

	//好友
	GiveActionID  int //赠送体力值ID
	GiveActionNum int //赠送体力值数量
	MaxRecvTime   int //最大接收次数

	//领取体力活动补签
	ActionActivityRetroactiveMoneyID  int //! 领取体力活动补签需求货币ID
	ActionActivityRetroactiveMoneyNum int //! 领取体力活动补签需求货币数量

	//将灵
	HeroSoulsStoreFixedItemID       int //! 将灵商店固定物品
	HeroSoulsStoreFixedItemMoneyID  int
	HeroSoulsStoreFixedItemMoneyNum int

	HeroSoulsStoreFixedItemID2       int
	HeroSoulsStoreFixedItemMoneyID2  int
	HeroSoulsStoreFixedItemMoneyNum2 int

	HeroSoulsStoreRefreshTime []int //! 将灵商店刷新时间

	HeroSoulsRefreshCostMoneyID    int //! 将灵刷新消耗货币ID
	HeroSoulsRefreshCostMoneyValue int //! 将灵刷新消耗货币系数 计算方法:  10+已击杀英灵数*20
	HeroSoulsRefreshGetMoneyID     int //! 将灵刷新获取货币ID(英魂)
	HeroSoulsRefreshGetMoneyValue  int //! 将灵刷新获取货币系数(英魂) 计算方法: 系数*品质

	HeroSoulsChallengeTimes  int //! 每日挑战将灵次数
	BuyChallengeTimesMoneyID int //! 购买将灵挑战次数使用货币ID

	//云游
	WanderInitTime    int //云游重置初始次数
	WanderBeginID     int //云游起始副本ID
	WanderEndID       int //云游结束副本ID
	WanderSingleBoxID int //云游单抽奖励ID
	WanderTenBoxID    int //云游十连抽奖励ID
	WanderDrawMoneyID int //抽宝符货币ID
	WanderDrawNum     int //单抽货币数
	WanderTenDrawNum  int //十连抽货币数
	WanderTenGiftID   int //十连抽必送道具ID
	WanderTenGiftNum  int //十连抽必送道具数量

	//! 分解
	HeroExpDecomposeItemID     int //! 分解英雄经验道具
	HeroGodDecomposeSoulsID    int //! 分解英雄化身将魂ID
	HeroGodDecomposeItemID     int //! 分解英雄化身物品ID
	EquipRefineDecomposeItemID int //! 分解装备精炼道具ID
	RebornCostMoneyID          int //! 重生消耗货币ID
	RebornCostMoneyNum         int //! 重生消耗货币数量
	GemStrengthDecomposeItemID int //! 重生宝物获取道具ID
	GemRefineDecomposeItemID   int //! 重生宝物道具获得道具ID
	PetExpDecomposeItemID      int //! 分解宠物经验道具
	PetGodDecomposeItemID      int //! 分解宠物神练道具
	PetDecomposeSoulsID        int //! 分解宠物兽魂货币ID

	//! 巡回探宝
	LuckyDiceItemID        int //! 幸运骰子道具ID
	HuntTicketItemID       int //! 巡回探宝游戏券道具ID
	HuntFreeTimes          int //! 免费次数
	HuntCostMoneyID        int //! 消耗货币ID
	HuntCostMoneyNum       int //! 消耗货币数量
	EliteHuntRankNeedScore int //! 巡回探宝精英榜需求分数

	//! 幸运转盘
	LuckyWheelCostItemID  int //! 幸运转盘需求花费道具ID
	NormalWheelFreeTimes  int //! 普通转盘每日免费次数
	ExcitedWheelFreeTimes int //! 高级转盘每日免费次数
	NormalWheelMoneyID    int //! 普通转盘花费货币ID
	NormalWheelMoneyNum   int //! 普通转盘花费货币Num
	ExcitedWheelMoneyID   int //! 高级转盘花费货币ID
	ExcitedWheelMoneyNum  int //! 高级转盘花费货币Num

	//! 团购
	GroupPurchaseCostItemID  int //! 团购券ID
	GroupPurchaseCostMoneyID int //! 团购统一使用货币ID

	//! 卡牌大师
	CardMaster_CostType          int  // 抽卡消耗的货币类型
	CardMaster_FreeTimes         byte // 每日免费次数
	CardMaster_RaffleTicket      int  // 抽奖券ItemID
	CardMaster_NormalCost        int  // 普通抽奖：钻石消耗
	CardMaster_NormalJiFen       int  // 普通抽奖：获得积分
	CardMaster_NormalAwardID     int  // 普通抽奖ID
	CardMaster_NormalAwardID_10  int
	CardMaster_NormalCost_10     int // 普通十连抽：钻石消耗
	CardMaster_SpecialCost       int // 高级抽奖：钻石消耗
	CardMaster_SpecialJiFen      int // 高级抽奖：获得积分
	CardMaster_SpecialAwardID    int // 高级抽奖ID
	CardMaster_SpecialAwardID_10 int
	CardMaster_SpecialCost_10    int // 高级十连抽：钻石消耗
	CardMaster_BigJoker_CardID   int
	CardMaster_TodayRank_Limit   int // 单日榜 特殊奖门槛
	CardMaster_TotalRank_Limit   int // 累计榜 特殊奖门槛

	//! 召唤
	NormalSummonAwardID int
	SeniorSummonAwardID int
	OrangeSummonAwardID int

	//阵营战配制
	CampBat_MoveTimes    int //每日可搬水晶次数
	Campbat_MaxMoveTime  int //每次搬水晶次的最大时间
	CampBat_SelCampAward int //选择推荐阵营的奖励
	CampBat_RoomMatchLvl int //高低等级房间分界线(线上的也属性低等级)
	CampBat_NtyKillNum   int //连杀走马灯人数
	CampBat_Chg_MoneyID  int //置换花费的货币ID
	CampBat_Chg_MoneyNum int //置换花费的货币数
	CampBat_KillHonorMax int //每日击杀荣誉上限
	Campbat_KillHonorOne int //击杀一个玩家获得的荣誉

	//月基金
	MonthFundCostMoneyID  int //! 月基金购买货币ID
	MonthFundCostMoneyNum int //! 月基金购买货币数量

	//时装
	FashionMeltingSum     int //!时装的熔炼总值
	FashionMeltingAwardID int //!时装熔炼的产出ID

	//竞技场
	ArenaBattleVictoryPercent int //! 竞技场挑战多次胜率

)

func InitOptionParser(total int) bool {

	return true
}

func ParseOptionRecord(rs *RecordSet) {
	switch rs.Values[0] {
	case "hero_store_free_refresh_times":
		{
			StoreFreeRefreshTimes = CheckAtoiName(rs.Values[2], "hero_store_free_refresh_times")
		}
	case "hero_store_free_refresh_add_time":
		{
			StoreFreeRefreshAddTime = CheckAtoiName(rs.Values[2], "hero_store_free_refresh_add_time")
		}
	case "hero_store_refresh_need_item":
		{
			StoreRefreshNeedItem = CheckAtoiName(rs.Values[2], "hero_store_refresh_need_item")
		}
	case "hero_store_refresh_need_money_type":
		{
			HeroStoreRefreshNeedMoneyType = CheckAtoiName(rs.Values[2], "hero_store_refresh_need_money_type")
		}
	case "hero_store_refresh_need_money_num":
		{
			HeroStoreRefreshNeedMoneyNum = CheckAtoiName(rs.Values[2], "hero_store_refresh_need_money_num")
		}
	case "awake_store_refresh_need_money_type":
		{
			AwakeStoreRefreshNeedMoneyType = CheckAtoiName(rs.Values[2], "awake_store_refresh_need_money_type")
		}
	case "awake_store_refresh_need_money_num":
		{
			AwakeStoreRefreshNeedMoneyNum = CheckAtoiName(rs.Values[2], "awake_store_refresh_need_money_num")
		}
	case "pet_store_refresh_need_money_type":
		{
			PetStoreRefreshNeedMoneyType = CheckAtoiName(rs.Values[2], "pet_store_refresh_need_money_type")
		}
	case "pet_store_refresh_need_money_num":
		{
			PetStoreRefreshNeedMoneyNum = CheckAtoiName(rs.Values[2], "pet_store_refresh_need_money_num")
		}
	case "hero_store_refresh_item_number":
		{
			StoreRefreshItemNum = CheckAtoiName(rs.Values[2], "hero_store_refresh_item_number")
		}
	case "famous_copy_challenge_times":
		{
			FamousCopyChallengeTimes = CheckAtoiName(rs.Values[2], "famous_copy_challenge_times")
		}
	case "property_attack":
		{
			AttackPropertyID = CheckAtoiName(rs.Values[2], "property_attack")
		}
	case "property_attack_physic":
		{
			AttackPhysicID = CheckAtoiName(rs.Values[2], "property_attack_physic")
		}
	case "property_attack_magic":
		{
			AttackMagicID = CheckAtoiName(rs.Values[2], "property_attack_magic")
		}
	case "property_defence":
		{
			DefencePropertyID = CheckAtoiName(rs.Values[2], "property_defence")
		}
	case "property_defence_physic":
		{
			DefencePhysicID = CheckAtoiName(rs.Values[2], "property_defence_physic")
		}
	case "property_defence_magic":
		{
			DefenceMagicID = CheckAtoiName(rs.Values[2], "property_defence_magic")
		}
	case "property_all":
		{
			AllPropertyID = CheckAtoiName(rs.Values[2], "property_all")
		}
	case "mall_item_money_id":
		{
			MallItemMoneyID = CheckAtoiName(rs.Values[2], "mall_item_money_id")
		}
	case "elite_copy_reset_money_id":
		{
			EliteCopyResetMoneyID = CheckAtoiName(rs.Values[2], "elite_copy_reset_money_id")
		}
	case "mall_gift_money_id":
		{
			MallGiftMoneyID = CheckAtoiName(rs.Values[2], "mall_gift_money_id")
		}
	case "vip_weekly_gift_money_id":
		{
			VipWeeklyGiftMoneyID = CheckAtoiName(rs.Values[2], "vip_weekly_gift_money_id")
		}
	case "first_recharge_gift_id":
		{
			FirstRechargeAwardID = CheckAtoiName(rs.Values[2], "first_recharge_gift_id")
		}
	case "next_recharge_gift_id":
		{
			NextRechargeAwardID = CheckAtoiName(rs.Values[2], "next_recharge_gift_id")
		}
	case "sanguo_wushuang_elite_challenge_times":
		{
			SangokuMusouEliteFreeTimes = CheckAtoiName(rs.Values[2], "sanguo_wushuang_elite_challenge_times")
		}
	case "sanguo_wushuang_reset_copy_moneyid":
		{
			SangokuMusouResetCopyMoneyID = CheckAtoiName(rs.Values[2], "sanguo_wushuang_reset_copy_moneyid")
		}
	case "sanguo_wushuang_add_elite_copy_moneyid":
		{
			SangokuMusouAddEliteCopyMoneyID = CheckAtoiName(rs.Values[2], "sanguo_wushuang_add_elite_copy_moneyid")
		}
	case "arena_rank_calc_time":
		{
			ArenaRankAwardCalcTime = CheckAtoiName(rs.Values[2], "arena_rank_calc_time")
		}
	case "summon_free_times_normal":
		{
			NormalSummonFreeTimes = CheckAtoiName(rs.Values[2], "summon_free_times_normal")
		}
	case "summon_free_cd_time_normal":
		{
			NormalSummonFreeCDTime = CheckAtoiName(rs.Values[2], "summon_free_cd_time_normal")
		}
	case "summon_free_cd_time_senior":
		{
			SeniorSummonFreeCDTime = CheckAtoiName(rs.Values[2], "summon_free_cd_time_senior")
		}
	case "summon_ten_times_discount":
		{
			TenSummonDiscount = CheckAtoiName(rs.Values[2], "summon_ten_times_discount")
		}
	case "summon_point":
		{
			SeniorSummonPoint = CheckAtoiName(rs.Values[2], "summon_point")
		}
	case "suppress_riot_award_item":
		{
			SuppressRiotAwardItem = CheckAtoiName(rs.Values[2], "suppress_riot_award_item")
		}
	case "suppress_riot_award_item_num":
		{
			SuppressRiotAwardItemNum = CheckAtoiName(rs.Values[2], "suppress_riot_award_item_num")
		}
	case "suppress_riot_friend_award_item":
		{
			SuppressRiotFriendAwardItem = CheckAtoiName(rs.Values[2], "suppress_riot_friend_award_item")
		}
	case "suppress_riot_friend_award_item_num":
		{
			SuppressRiotFriendAwardNum = CheckAtoiName(rs.Values[2], "suppress_riot_friend_award_item_num")
		}
	case "rebel_escape_time":
		{
			RebelEscapeTime = CheckAtoiName(rs.Values[2], "rebel_escape_time")
		}
	case "find_rebel_pro":
		{
			FindRebelPro = CheckAtoiName(rs.Values[2], "find_rebel_pro")
		}
	case "lower_rebel_pro":
		{
			LowerRebelPro = CheckAtoiName(rs.Values[2], "lower_rebel_pro")
		}
	case "middle_rebel_pro":
		{
			MiddleRebelPro = CheckAtoiName(rs.Values[2], "middle_rebel_pro")
		}
	case "senior_rebel_pro":
		{
			SeniorRebelPro = CheckAtoiName(rs.Values[2], "senior_rebel_pro")
		}
	case "rebel_achievements":
		{
			RebelAchievements = CheckAtoiName(rs.Values[2], "rebel_achievements")
		}
	case "attack_rebel_need_action_id":
		{
			AttackRebelActionID = CheckAtoiName(rs.Values[2], "attack_rebel_need_action_id")
		}
	case "normal_attack_rebel_action_num":
		{
			NormalAttackRebelNeedActionNum = CheckAtoiName(rs.Values[2], "normal_attack_rebel_action_num")
		}
	case "senior_attack_rebel_action_num":
		{
			SeniorAttackRebelNeedActionNum = CheckAtoiName(rs.Values[2], "senior_attack_rebel_action_num")
		}
	case "mining_map_length":
		{
			MiningMapLength = CheckAtoiName(rs.Values[2], "mining_map_length")
		}
	case "mining_map_slice":
		{
			MiningMapSlice = CheckAtoiName(rs.Values[2], "mining_map_slice")
		}
	case "mining_cost_action_id":
		{
			MiningCostActionID = CheckAtoiName(rs.Values[2], "mining_cost_action_id")
		}
	case "mining_cost_action_num":
		{
			MiningCostActionNum = CheckAtoiName(rs.Values[2], "mining_cost_action_num")
		}
	case "mining_attack_boss_action":
		{
			MiningAttackBossAction = CheckAtoiName(rs.Values[2], "mining_attack_boss_action")
		}
	case "mining_action_recover_time":
		{
			MiningActionRecoverTime = CheckAtoiName(rs.Values[2], "mining_action_recover_time")
		}
	case "mining_action_recover_limit":
		{
			MiningActionRecoverLimit = CheckAtoiName(rs.Values[2], "mining_action_recover_limit")
		}
	case "mining_reset_money_id":
		{
			MiningResetMoneyID = CheckAtoiName(rs.Values[2], "mining_reset_money_id")
		}
	case "elite_invade_time1":
		{
			EliteInvadeTime1 = CheckAtoiName(rs.Values[2], "elite_invade_time1")
		}
	case "elite_invade_num1":
		{
			EliteInvadeNum1 = CheckAtoiName(rs.Values[2], "elite_invade_num1")
		}
	case "elite_invade_time2":
		{
			EliteInvadeTime2 = CheckAtoiName(rs.Values[2], "elite_invade_time2")
		}
	case "elite_invade_num2":
		{
			EliteInvadeNum2 = CheckAtoiName(rs.Values[2], "elite_invade_num2")
		}
	case "elite_invade_time3":
		{
			EliteInvadeTime3 = CheckAtoiName(rs.Values[2], "elite_invade_time3")
		}
	case "elite_invade_num3":
		{
			EliteInvadeNum3 = CheckAtoiName(rs.Values[2], "elite_invade_num3")
		}
	case "elite_invade_time4":
		{
			EliteInvadeTime4 = CheckAtoiName(rs.Values[2], "elite_invade_time4")
		}
	case "elite_invade_num4":
		{
			EliteInvadeNum4 = CheckAtoiName(rs.Values[2], "elite_invade_num4")
		}
	case "riot_pro":
		{
			RiotPro = CheckAtoiName(rs.Values[2], "riot_pro")
		}
	case "riot_time":
		{
			RiotTime = CheckAtoiName(rs.Values[2], "riot_time")
		}
	case "crit_multiple":
		{
			CritMultiple = CheckAtoiName(rs.Values[2], "crit_multiple")
		}
	case "big_crit_multiple":
		{
			BigCritMultiple = CheckAtoiName(rs.Values[2], "big_crit_multiple")
		}
	case "lucky_crit_multiple":
		{
			LuckyCritMultiple = CheckAtoiName(rs.Values[2], "lucky_crit_multiple")
		}
	case "crit_pro":
		{
			CritPro = CheckAtoiName(rs.Values[2], "crit_pro")
		}
	case "big_crit_pro":
		{
			BigCritPro = CheckAtoiName(rs.Values[2], "big_crit_pro")
		}
	case "lucky_crit_pro":
		{
			LuckyCritPro = CheckAtoiName(rs.Values[2], "lucky_crit_pro")
		}
	case "charge_money_id":
		{
			ChargeMoneyID = CheckAtoiName(rs.Values[2], "charge_money_id")
		}
	case "vip_exp_money_id":
		{
			VipExpMoneyID = CheckAtoiName(rs.Values[2], "vip_exp_money_id")
		}
	case "rebel_bili":
		{
			RebelExploitPoint = CheckAtoiName(rs.Values[2], "rebel_bili")
			if RebelExploitPoint <= 0 {
				panic(fmt.Sprintf("[%s] can not be zero !!!", rs.Values[0]))
			}
		}
	case "seven_activity_award_day":
		{
			SevenActivityAwardDay = CheckAtoiName(rs.Values[2], "seven_activity_award_day")
		}
	case "open_fund_money_id":
		{
			OpenFundPriceID = CheckAtoiName(rs.Values[2], "open_fund_money_id")
		}
	case "open_fund_money_num":
		{
			OpenFundPriceNum = CheckAtoiName(rs.Values[2], "open_fund_money_num")
		}
	case "mining_enter_point_x":
		{
			MiningEnterPointX = CheckAtoiName(rs.Values[2], "mining_enter_point_x")
		}
	case "mining_enter_point_y":
		{
			MiningEnterPointY = CheckAtoiName(rs.Values[2], "mining_enter_point_y")
		}
	case "mining_monster_boss_value":
		{
			MiningBossValue = CheckAtoiName(rs.Values[2], "mining_monster_boss_value")
		}
	case "hero_culture_item_id":
		{
			CultureItemID = CheckAtoiName(rs.Values[2], "hero_culture_item_id")
		}
	case "hero_culture_item_num":
		{
			CultureItemNum = CheckAtoiName(rs.Values[2], "hero_culture_item_num")
		}
	case "hangup_init_grid_num":
		{
			HangUpInitGridNum = CheckAtoiName(rs.Values[2], "hangup_init_grid_num")
		}
	case "hangup_quick_times":
		{
			HangUpQuickFight = CheckAtoiName(rs.Values[2], "hangup_quick_times")
		}
	case "hangup_grid_num":
		{
			HangUpOpenGridNum = CheckAtoiName(rs.Values[2], "hangup_grid_num")
		}
	case "hangup_buy_grid_money":
		{
			HangUpBuyGridMoneyID = CheckAtoiName(rs.Values[2], "hangup_buy_grid_money")
		}
	case "guild_action_recover_time":
		{
			GuildActionRecoverTime = CheckAtoiName(rs.Values[2], "guild_action_recover_time")
		}
	case "guild_copy_battle_time_begin":
		{
			GuildCopyBattleTimeBegin = CheckAtoiName(rs.Values[2], "guild_copy_battle_time_begin")
		}
	case "guild_copy_battle_time_end":
		{
			GuildCopyBattleTimeEnd = CheckAtoiName(rs.Values[2], "guild_copy_battle_time_end")
		}
	case "guild_guild_sacrifice_crit":
		{
			GuildSacrificeCrit = CheckAtoiName(rs.Values[2], "guild_guild_sacrifice_crit")
		}
	case "create_guild_money_id":
		{
			CreateGuildMoneyID = CheckAtoiName(rs.Values[2], "create_guild_money_id")
		}
	case "create_guild_money_num":
		{
			CreateGuildMoneyNum = CheckAtoiName(rs.Values[2], "create_guild_money_num")
		}
	case "update_guild_name_money_id":
		{
			UpdateGuildNameMoneyID = CheckAtoiName(rs.Values[2], "update_guild_name_money_id")
		}
	case "update_guild_name_money_num":
		{
			UpdateGuildNameMoneyNum = CheckAtoiName(rs.Values[2], "update_guild_name_money_num")
		}
	case "black_market_refresh_time":
		{
			timeLst := strings.Split(rs.Values[2], "|")
			BlackMarketRefreshTime = make([]int, len(timeLst))
			for i, v := range timeLst {
				BlackMarketRefreshTime[i] = CheckAtoiName(v, "black_market_refresh_time")
			}
		}
	case "enter_vip_level":
		{
			EnterVipLevel = CheckAtoiName(rs.Values[2], "enter_vip_level")
		}
	case "black_market_pro":
		{
			BlackMarketPro = CheckAtoiName(rs.Values[2], "black_market_pro")
		}
	case "guild_skill_need_money_id":
		{
			GuildSKillStudyNeedMoneyID = CheckAtoiName(rs.Values[2], "guild_skill_need_money_id")
		}
	case "score_one_fight_score":
		{
			OneTimeFightScore = CheckAtoiName(rs.Values[2], "score_one_fight_score")
		}
	case "score_copy_id":
		{
			ScoreCopyID = CheckAtoiName(rs.Values[2], "score_copy_id")
		}
	case "score_money_id":
		{
			ScoreMoneyID = CheckAtoiName(rs.Values[2], "score_money_id")
		}
	case "score_money_num":
		{
			ScoreMoneyNum = CheckAtoiName(rs.Values[2], "score_money_num")
		}
	case "score_series_win_times":
		{
			ScoreSeriesWinTimes = CheckAtoiName(rs.Values[2], "score_series_win_times")
		}
	case "score_series_win_awardid":
		{
			ScoreSeriesWinAwardID = CheckAtoiName(rs.Values[2], "score_series_win_awardid")
		}
	case "score_buytimes_moneyid":
		{
			ScoreBuyTimeMoneyID = CheckAtoiName(rs.Values[2], "score_buytimes_moneyid")
		}

	case "fame_hall_free_times":
		{
			FameHallFreeTimes = CheckAtoiName(rs.Values[2], "fame_hall_free_times")
		}
	case "food_war_open_time":
		{
			FoodWarOpenTime = CheckAtoiName(rs.Values[2], "food_war_open_time")
		}
	case "food_war_end_time":
		{
			FoodWarEndTime = CheckAtoiName(rs.Values[2], "food_war_end_time")
		}
	case "food_war_open_day":
		{
			values := strings.Split(rs.Values[2], "|")
			for i, v := range values {
				FoodWarOpenDay = append(FoodWarOpenDay, CheckAtoi(v, i))
			}
		}

	case "food_war_attack_times":
		{
			FoodWarAttackTimes = CheckAtoiName(rs.Values[2], "food_war_attack_times")

		}
	case "food_war_revenge_times":
		{
			FoodWarRevengeTimes = CheckAtoiName(rs.Values[2], "food_war_revenge_times")
		}
	case "food_war_fixed_food":
		{
			FoodWarFixedFood = CheckAtoiName(rs.Values[2], "food_war_fixed_food")
		}
	case "food_war_nonfixed_food":
		{
			FoodWarNonFixedFood = CheckAtoiName(rs.Values[2], "food_war_nonfixed_food")
		}
	case "food_war_time_add_food":
		{
			FoodWarTimeAddFood = CheckAtoiName(rs.Values[2], "food_war_time_add_food")
		}
	case "food_war_rob_bili":
		{
			FoodWarRobBili = CheckAtoiName(rs.Values[2], "food_war_rob_bili")
		}
	case "food_war_copy_id":
		{
			FoodWarCopyID = CheckAtoiName(rs.Values[2], "food_war_copy_id")
		}
	case "food_war_victory_money_id":
		{
			FoodWarVictoryMoneyID = CheckAtoiName(rs.Values[2], "food_war_victory_money_id")
		}
	case "food_war_victory_money_num":
		{
			FoodWarVictoryMoneyNum = CheckAtoiName(rs.Values[2], "food_war_victory_money_num")
		}
	case "food_war_failed_money_id":
		{
			FoodWarFailedMoneyID = CheckAtoiName(rs.Values[2], "food_war_failed_money_id")
		}
	case "food_war_failed_money_num":
		{
			FoodWarFailedMoneyNum = CheckAtoiName(rs.Values[2], "food_war_failed_money_num")
		}
	case "food_war_buy_times_need_money_id":
		{
			FoodWarBuyTimesNeedMoneyID = CheckAtoiName(rs.Values[2], "food_war_buy_times_need_money_id")
		}
	case "give_action_id":
		{
			GiveActionID = CheckAtoiName(rs.Values[2], "give_action_id")
		}
	case "give_action_num":
		{
			GiveActionNum = CheckAtoiName(rs.Values[2], "give_action_num")
		}
	case "max_recv_time":
		{
			MaxRecvTime = CheckAtoiName(rs.Values[2], "max_recv_time")
		}
	case "action_retroactive_money_id":
		{
			ActionActivityRetroactiveMoneyID = CheckAtoiName(rs.Values[2], "action_retroactive_money_id")
		}
	case "action_retroactive_money_num":
		{
			ActionActivityRetroactiveMoneyNum = CheckAtoiName(rs.Values[2], "action_retroactive_money_num")
		}
	case "hero_souls_store_fixed_item_id":
		{
			HeroSoulsStoreFixedItemID = CheckAtoiName(rs.Values[2], "hero_souls_store_fixed_item_id")
		}
	case "hero_souls_store_fixed_item_money_id":
		{
			HeroSoulsStoreFixedItemMoneyID = CheckAtoiName(rs.Values[2], "hero_souls_store_fixed_item_money_id")
		}
	case "hero_souls_store_fixed_item_money_num":
		{
			HeroSoulsStoreFixedItemMoneyNum = CheckAtoiName(rs.Values[2], "hero_souls_store_fixed_item_money_num")
		}
	case "hero_souls_store_fixed_item_id2":
		{
			HeroSoulsStoreFixedItemID2 = CheckAtoiName(rs.Values[2], "hero_souls_store_fixed_item_id2")
		}
	case "hero_souls_store_fixed_item_money_id2":
		{
			HeroSoulsStoreFixedItemMoneyID2 = CheckAtoiName(rs.Values[2], "hero_souls_store_fixed_item_money_id2")
		}
	case "hero_souls_store_fixed_item_money_num2":
		{
			HeroSoulsStoreFixedItemMoneyNum2 = CheckAtoiName(rs.Values[2], "hero_souls_store_fixed_item_money_num2")
		}
	case "hero_souls_store_refresh_time":
		{
			values := strings.Split(rs.Values[2], "|")
			for i, v := range values {
				HeroSoulsStoreRefreshTime = append(HeroSoulsStoreRefreshTime, CheckAtoi(v, i))
			}
		}
	case "hero_souls_refresh_cost_money_id":
		{
			HeroSoulsRefreshCostMoneyID = CheckAtoiName(rs.Values[2], "hero_souls_refresh_cost_money_id")
		}
	case "hero_souls_refresh_cost_money_value":
		{
			HeroSoulsRefreshCostMoneyValue = CheckAtoiName(rs.Values[2], "hero_souls_refresh_cost_money_value")
		}
	case "hero_souls_refresh_get_money_id":
		{
			HeroSoulsRefreshGetMoneyID = CheckAtoiName(rs.Values[2], "hero_souls_refresh_get_money_id")
		}
	case "hero_souls_refresh_get_money_value":
		{
			HeroSoulsRefreshGetMoneyValue = CheckAtoiName(rs.Values[2], "hero_souls_refresh_get_money_value")
		}
	case "hero_souls_challenge_times":
		{
			HeroSoulsChallengeTimes = CheckAtoiName(rs.Values[2], "hero_souls_challenge_times")
		}
	case "buy_challenge_times_money_id":
		{
			BuyChallengeTimesMoneyID = CheckAtoiName(rs.Values[2], "buy_challenge_times_money_id")
		}
	case "wander_init_reset_time":
		{
			WanderInitTime = CheckAtoiName(rs.Values[2], "wander_init_reset_time")
		}
	case "wander_copy_begin_id":
		{
			WanderBeginID = CheckAtoiName(rs.Values[2], "wander_copy_begin_id")
		}
	case "wander_copy_end_id":
		{
			WanderEndID = CheckAtoiName(rs.Values[2], "wander_copy_end_id")
		}
	case "wander_box_single_id":
		{
			WanderSingleBoxID = CheckAtoiName(rs.Values[2], "wander_box_single_id")
		}
	case "wander_box_ten_id":
		{
			WanderTenBoxID = CheckAtoiName(rs.Values[2], "wander_box_ten_id")
		}
	case "wander_draw_money_id":
		{
			WanderDrawMoneyID = CheckAtoiName(rs.Values[2], "wander_draw_money_id")
		}
	case "wander_single_money_num":
		{
			WanderDrawNum = CheckAtoiName(rs.Values[2], "wander_single_money_num")
		}
	case "wander_ten_money_num":
		{
			WanderTenDrawNum = CheckAtoiName(rs.Values[2], "wander_ten_money_num")
		}
	case "wander_ten_gift_id":
		{
			WanderTenGiftID = CheckAtoiName(rs.Values[2], "wander_ten_gift_id")
		}
	case "wander_ten_gift_num":
		{
			WanderTenGiftNum = CheckAtoiName(rs.Values[2], "wander_ten_gift_num")
		}
	case "hero_exp_decompose_item_id":
		{
			HeroExpDecomposeItemID = CheckAtoiName(rs.Values[2], "hero_exp_decompose_item_id")
		}
	case "hero_god_decompose_souls_id":
		{
			HeroGodDecomposeSoulsID = CheckAtoiName(rs.Values[2], "hero_god_decompose_souls_id")
		}
	case "hero_god_decompose_item_id":
		{
			HeroGodDecomposeItemID = CheckAtoiName(rs.Values[2], "hero_god_decompose_item_id")
		}
	case "reborn_cost_money_id":
		{
			RebornCostMoneyID = CheckAtoiName(rs.Values[2], "reborn_cost_money_id")
		}
	case "reborn_cost_money_num":
		{
			RebornCostMoneyNum = CheckAtoiName(rs.Values[2], "reborn_cost_money_num")
		}
	case "equip_refine_decompose_item_id":
		{
			EquipRefineDecomposeItemID = CheckAtoiName(rs.Values[2], "equip_refine_decompose_item_id")
		}
	case "gem_strength_decompose_item_id":
		{
			GemStrengthDecomposeItemID = CheckAtoiName(rs.Values[2], "gem_strength_decompose_item_id")
		}
	case "gem_refine_decompose_item_id":
		{
			GemRefineDecomposeItemID = CheckAtoiName(rs.Values[2], "gem_refine_decompose_item_id")
		}
	case "pet_exp_decompose_item_id":
		{
			PetExpDecomposeItemID = CheckAtoiName(rs.Values[2], "pet_exp_decompose_item_id")
		}
	case "pet_god_decompose_item_id":
		{
			PetGodDecomposeItemID = CheckAtoiName(rs.Values[2], "pet_god_decompose_item_id")
		}
	case "pet_decompose_souls_id":
		{
			PetDecomposeSoulsID = CheckAtoiName(rs.Values[2], "pet_decompose_souls_id")
		}
	case "lucky_dice_item_id":
		{
			LuckyDiceItemID = CheckAtoiName(rs.Values[2], "lucky_dice_item_id")
		}
	case "hunt_ticket_item_id":
		{
			HuntTicketItemID = CheckAtoiName(rs.Values[2], "hunt_ticket_item_id")
		}
	case "hunt_free_time":
		{
			HuntFreeTimes = CheckAtoiName(rs.Values[2], "hunt_free_time")
		}
	case "hunt_cost_money_id":
		{
			HuntCostMoneyID = CheckAtoiName(rs.Values[2], "hunt_cost_money_id")
		}
	case "hunt_cost_money_num":
		{
			HuntCostMoneyNum = CheckAtoiName(rs.Values[2], "hunt_cost_money_num")
		}
	case "elite_hunt_rank_need_score":
		{
			EliteHuntRankNeedScore = CheckAtoiName(rs.Values[2], "elite_hunt_rank_need_score")
		}
	case "lucky_wheel_cost_item_id":
		{
			LuckyWheelCostItemID = CheckAtoiName(rs.Values[2], "lucky_wheel_cost_item_id")
		}
	case "normal_wheel_free_times":
		{
			NormalWheelFreeTimes = CheckAtoiName(rs.Values[2], "normal_wheel_free_times")
		}
	case "excited_wheel_free_times":
		{
			ExcitedWheelFreeTimes = CheckAtoiName(rs.Values[2], "excited_wheel_free_times")
		}
	case "normal_wheel_money_id":
		{
			NormalWheelMoneyID = CheckAtoiName(rs.Values[2], "normal_wheel_money_id")
		}
	case "normal_wheel_money_num":
		{
			NormalWheelMoneyNum = CheckAtoiName(rs.Values[2], "normal_wheel_money_num")
		}
	case "excited_wheel_money_id":
		{
			ExcitedWheelMoneyID = CheckAtoiName(rs.Values[2], "excited_wheel_money_id")
		}
	case "excited_wheel_money_num":
		{
			ExcitedWheelMoneyNum = CheckAtoiName(rs.Values[2], "excited_wheel_money_num")
		}
	case "act_card_master_cost_type":
		{
			CardMaster_CostType = CheckAtoiName(rs.Values[2], "act_card_master_cost_type")
		}
	case "act_card_master_free_times":
		{
			CardMaster_FreeTimes = byte(CheckAtoiName(rs.Values[2], "act_card_master_free_times"))
		}
	case "act_card_master_raffle_ticket":
		{
			CardMaster_RaffleTicket = CheckAtoiName(rs.Values[2], "act_card_master_raffle_ticket")
		}
	case "act_card_master_normal_cost":
		{
			CardMaster_NormalCost = CheckAtoiName(rs.Values[2], "act_card_master_normal_cost")
		}
	case "act_card_master_normal_cost_10":
		{
			CardMaster_NormalCost_10 = CheckAtoiName(rs.Values[2], "act_card_master_normal_cost_10")
		}
	case "act_card_master_normal_jifen":
		{
			CardMaster_NormalJiFen = CheckAtoiName(rs.Values[2], "act_card_master_normal_jifen")
		}
	case "act_card_master_normal_award_id":
		{
			CardMaster_NormalAwardID = CheckAtoiName(rs.Values[2], "act_card_master_normal_award_id")
		}
	case "act_card_master_normal_award_id_10":
		{
			CardMaster_NormalAwardID_10 = CheckAtoiName(rs.Values[2], "act_card_master_normal_award_id_10")
		}
	case "act_card_master_special_cost":
		{
			CardMaster_SpecialCost = CheckAtoiName(rs.Values[2], "act_card_master_special_cost")
		}
	case "act_card_master_special_cost_10":
		{
			CardMaster_SpecialCost_10 = CheckAtoiName(rs.Values[2], "act_card_master_special_cost_10")
		}
	case "act_card_master_special_jifen":
		{
			CardMaster_SpecialJiFen = CheckAtoiName(rs.Values[2], "act_card_master_special_jifen")
		}
	case "act_card_master_special_award_id":
		{
			CardMaster_SpecialAwardID = CheckAtoiName(rs.Values[2], "act_card_master_special_award_id")
		}
	case "act_card_master_special_award_id_10":
		{
			CardMaster_SpecialAwardID_10 = CheckAtoiName(rs.Values[2], "act_card_master_special_award_id_10")
		}
	case "act_card_master_big_joker_card_id":
		{
			CardMaster_BigJoker_CardID = CheckAtoiName(rs.Values[2], "act_card_master_big_joker_card_id")
		}
	case "act_card_master_today_limit_score":
		{
			CardMaster_TodayRank_Limit = CheckAtoiName(rs.Values[2], "act_card_master_today_limit_score")
		}
	case "act_card_master_total_limit_score":
		{
			CardMaster_TotalRank_Limit = CheckAtoiName(rs.Values[2], "act_card_master_total_limit_score")
		}
	case "moonlight_shop_token_item_id":
		{
			MoonlightShop_Token_ItemID = CheckAtoiName(rs.Values[2], "moonlight_shop_token_item_id")
		}
	case "moonlight_shop_discount_cost":
		{
			MoonlightShop_Discount_Cost = ParseStringToByteArray(rs.Values[2])
		}
	case "moonlight_shop_discount_one_tiems":
		{
			MoonlightShop_Discount_OneTiems = ParseStringToByteArray(rs.Values[2])
		}
	case "moonlight_shop_refresh_cd":
		{
			MoonlightShop_Shop_Refresh_CD = byte(CheckAtoiName(rs.Values[2], "moonlight_shop_refresh_cd"))
		}
	case "moonlight_shop_refresh_cost":
		{
			MoonlightShop_Shop_Refresh_Cost = byte(CheckAtoiName(rs.Values[2], "moonlight_shop_refresh_cost"))
		}
	case "moonlight_shop_buy_times_max":
		{
			MoonlightShop_BuyTimes_Max = CheckAtoiName(rs.Values[2], "moonlight_shop_buy_times_max")
		}
	case "beach_baby_token_item_id":
		{
			BeachBaby_Token_ItemID = CheckAtoiName(rs.Values[2], "beach_baby_token_item_id")
		}
	case "beach_baby_cost_money_type":
		{
			BeachBaby_CostMoneyType = CheckAtoiName(rs.Values[2], "beach_baby_cost_money_type")
		}
	case "beach_baby_get_free_token_time":
		{
			BeachBaby_GetFreeToken_Time = ParseStringToByteArray(rs.Values[2])
			if len(BeachBaby_GetFreeToken_Time) > 8 {
				panic(fmt.Sprintf("[%s] element cnt is up to 8 !!!", rs.Values[0]))
			}
		}
	case "beach_baby_get_free_token_cnt":
		{
			BeachBaby_GetFreeToken_Cnt = byte(CheckAtoiName(rs.Values[2], "beach_baby_get_free_token_cnt"))
		}
	case "beach_baby_open_goods_cost":
		{
			BeachBaby_OpenGoods_Cost = ParseStringToByteArray(rs.Values[2])
			if len(BeachBaby_OpenGoods_Cost) != 16 { // 沙滩宝贝商品数量
				panic(fmt.Sprintf("[%s] element cnt is not 16 !!!", rs.Values[0]))
			}
		}
	case "beach_baby_refresh_cd":
		{
			BeachBaby_Refresh_CD = byte(CheckAtoiName(rs.Values[2], "beach_baby_refresh_cd"))
		}
	case "beach_baby_refresh_cost":
		{
			BeachBaby_Refresh_Cost = byte(CheckAtoiName(rs.Values[2], "beach_baby_refresh_cost"))
		}
	case "beach_baby_select_goods_cost":
		{
			BeachBaby_SelectGoods_Cost = ParseStringToByteArray(rs.Values[2])
		}
	case "beach_baby_today_limit_score":
		{
			BeachBaby_TodayRank_Limit = CheckAtoiName(rs.Values[2], "beach_baby_today_limit_score")
		}
	case "beach_baby_total_limit_score":
		{
			BeachBaby_TotalRank_Limit = CheckAtoiName(rs.Values[2], "beach_baby_total_limit_score")
		}
	case "group_purchase_cost_item_id":
		{
			GroupPurchaseCostItemID = CheckAtoiName(rs.Values[2], "group_purchase_cost_item_id")
		}
	case "group_purchase_cost_money_id":
		{
			GroupPurchaseCostMoneyID = CheckAtoiName(rs.Values[2], "group_purchase_cost_money_id")
		}
	case "normal_summon_award_id":
		{
			NormalSummonAwardID = CheckAtoiName(rs.Values[2], "normal_summon_award_id")
		}
	case "senior_summon_award_id":
		{
			SeniorSummonAwardID = CheckAtoiName(rs.Values[2], "senior_summon_award_id")
		}
	case "orange_summon_award_id":
		{
			OrangeSummonAwardID = CheckAtoiName(rs.Values[2], "orange_summon_award_id")
		}
	case "campbat_move_crystal_times":
		{
			CampBat_MoveTimes = CheckAtoiName(rs.Values[2], "campbat_move_crystal_times")
		}
	case "campbat_select_camp_award":
		{
			CampBat_SelCampAward = CheckAtoiName(rs.Values[2], "campbat_select_camp_award")
		}
	case "campbat_room_match_level":
		{
			CampBat_RoomMatchLvl = CheckAtoiName(rs.Values[2], "campbat_room_match_level")
		}
	case "campbat_kill_num":
		{
			CampBat_NtyKillNum = CheckAtoiName(rs.Values[2], "campbat_kill_num")
		}
	case "campbat_max_move_time":
		{
			Campbat_MaxMoveTime = CheckAtoiName(rs.Values[2], "campbat_max_move_time")
		}
	case "campbat_chg_money_id":
		{
			CampBat_Chg_MoneyID = CheckAtoiName(rs.Values[2], "campbat_chg_money_id")
		}
	case "campbat_chg_money_num":
		{
			CampBat_Chg_MoneyNum = CheckAtoiName(rs.Values[2], "campbat_chg_money_num")
		}
	case "campbat_kill_honor_max":
		{
			CampBat_KillHonorMax = CheckAtoiName(rs.Values[2], "campbat_kill_honor_max")
		}
	case "campbat_kill_honor_one":
		{
			Campbat_KillHonorOne = CheckAtoiName(rs.Values[2], "campbat_kill_honor_one")
		}
	case "month_fund_money_id":
		{
			MonthFundCostMoneyID = CheckAtoiName(rs.Values[2], "month_fund_money_id")
		}
	case "month_fund_money_num":
		{
			MonthFundCostMoneyNum = CheckAtoiName(rs.Values[2], "month_fund_money_num")
		}
	case "next_recharge_need":
		{
			NextAwardNeedRecharge = CheckAtoiName(rs.Values[2], "next_recharge_need")
		}
	case "fashion_melting_sum":
		{
			FashionMeltingSum = CheckAtoiName(rs.Values[2], "fashion_melting_sum")
		}
	case "fashion_melting_awardid":
		{
			FashionMeltingAwardID = CheckAtoiName(rs.Values[2], "fashion_melting_awardid")
		}
	case "arena_battle_victory_percent":
		{
			ArenaBattleVictoryPercent = CheckAtoiName(rs.Values[2], "arena_battle_victory_percent")
		}
	default:
		{
			panic(fmt.Sprintf("[%s] not processed !!!", rs.Values[0]))
		}
	}
}

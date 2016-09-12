package gamedata

const (
	FUNC_BEGIN_ID                       = 0   //!功能ID开始标记
	FUNC_DOUBLE_POWER_UP                = 1   //! 强化打造两倍暴击
	FUNC_TRIPLE_POWER_UP                = 2   //! 强化打造三倍暴击
	FUNC_SUPPRESS_TERRITORY             = 3   //! 解决领地暴动次数
	FUNC_MAIN_COPY_RESET                = 4   //! 主线副本每日重置次数
	FUNC_HERO_BAG_CAPACITY              = 5   //! 武将背包额外增加容量
	FUNC_GEM_BAG_CAPACITY               = 6   //! 宝物背包额外增加容量
	FUNC_SANGUOWUSHUANG_RESET           = 7   //! 三国无双每日重置次数
	FUNC_HERO_STORE_RESET               = 8   //! 神将商店每日重置次数
	FUNC_AWAKEN_STORE_RESET             = 9   //! 觉醒商店每日重置次数
	FUNC_BATTLE_PET_STORE_RESET         = 10  //! 战宠商店每日重置次数
	FUNC_SERVER_BATTLE_BUY_TIMES        = 11  //! 跨服演武每日可购买挑战次数
	FUNC_BUY_SANGUOWUSHUANG_ELITE_TIMES = 12  //! 每日可购买三国无双精英挑战次数
	FUNC_FAMOUS_COPY_CHALLENGE_TIMES    = 13  //! 名将副本每日可攻略次数
	FUNC_RESET_EXPERIENCED_TIMES        = 14  //! 重置百战沙场次数
	FUNC_BUY_EXPEDITION_ORDER_TIMES     = 15  //! 每日可购买征讨令次数
	FUNC_BUY_PHYSICAL_STRENGTH_TIMES    = 16  //! 每日可购买体力丹次数
	FUNC_BUY_ENERGY_TIMES               = 17  //! 每日可购买精力丹次数
	FUNC_BUY_GOLD_DRAGON_TIMES          = 18  //! 每日可购买金龙宝宝次数
	FUNC_BUY_COIN_TIMES                 = 19  //! 每日可购买银两次数
	FUNC_BUY_ORANGE_EQUI_TIMES          = 20  //! 每日可购买橙色装备宝箱次数
	FUNC_BUY_ORANGE_TREASURE_TIMES      = 21  //! 每日可购买橙色宝物宝箱次数
	FUNC_GUILD_COPY_BUY_TIMES           = 22  //! 每日可购买公会副本挑战次数
	FUNC_FAMOUS_TRIALS_TIMES            = 23  //! 每日将灵,名将试炼可挑战次数
	FUNC_PLUTUS_TIMES                   = 24  //! 每日可招财次数
	FUNC_GLIDE_WORSHIP_CEREMONY_PLUS    = 25  //! 公会高级祭天功能
	FUNC_ROB_TREASURE_FIVE_TIMES        = 26  //! 夺宝五次功能
	FUNC_CULTURE_FIVE_TIMES             = 27  //! 培养五次功能
	FUNC_SANGUOWUSHUANG_SWEEP           = 28  //! 三国无双一键三星功能
	FUNC_TERRITORY_PATROL_INTERMEDIATE  = 29  //! 中级领地巡逻功能 20分钟收益
	FUNC_CULTURE_TEN_TIMES              = 30  //! 培养十次功能
	FUNC_ARENA_FIVE_CHALLENGE           = 31  //! 竞技场连战五次
	FUNC_REBEL_SIEGE_SKIP               = 32  //! 围剿叛军跳过功能
	FUNC_ROB_TREASURE_ONE_KEY           = 33  //! 一键夺宝功能
	FUNC_TERRITORY_PATROL_SENIOR        = 34  //! 高级领地巡逻功能
	FUNC_MINING                         = 35  //! 挖矿功能 预留
	FUNC_GUAJI                          = 36  //! 挂机功能 预留
	FUNC_BUY_REFRESH_ITEM               = 37  //! 每日可购买刷新令次数
	FUNC_FREE_WAR_ITEM_PLUS             = 38  //! 每日可购买免战牌(大)次数
	FUNC_FREE_WAR_ITEM                  = 39  //! 每日可购买免战牌(小)次数
	FUNC_BUY_FASHION_ESSENCE            = 40  //! 每日可购买时装精华次数
	FUNC_BUY_DESTINY_ITEM               = 41  //! 每日可购买天命石次数
	FUNC_BUY_GOLD_EXP_TREASURE          = 42  //! 每日可购买黄金经验宝物次数
	FUNC_BEST_REFINED_STONE             = 43  //! 每日可购买极品精炼石次数
	FUNC_TREASURE_REFINED_STONE         = 44  //! 每日可购买宝物精炼石次数
	FUNC_RED_EQUI_CASE                  = 45  //! 每日可购买红色装备箱子
	FUNC_VIP_GIFT_0                     = 46  //! VIP0超值礼包
	FUNC_VIP_GIFT_1                     = 47  //! VIP1超值礼包
	FUNC_VIP_GIFT_2                     = 48  //! VIP2超值礼包
	FUNC_VIP_GIFT_3                     = 49  //! VIP3超值礼包
	FUNC_VIP_GIFT_4                     = 50  //! VIP4超值礼包
	FUNC_VIP_GIFT_5                     = 51  //! VIP5超值礼包
	FUNC_VIP_GIFT_6                     = 52  //! VIP6超值礼包
	FUNC_VIP_GIFT_7                     = 53  //! VIP7超值礼包
	FUNC_VIP_GIFT_8                     = 54  //! VIP8超值礼包
	FUNC_VIP_GIFT_9                     = 55  //! VIP9超值礼包
	FUNC_VIP_GIFT_10                    = 56  //! VIP10超值礼包
	FUNC_VIP_GIFT_11                    = 57  //! VIP11超值礼包
	FUNC_VIP_GIFT_12                    = 58  //! VIP12超值礼包
	FUNC_POS_START                      = 59  //! 上阵阵位开启
	FUNC_POS_END                        = 65  //! 上阵阵位结束
	FUNC_HERO_LEVEL_UP                  = 66  //! 英雄升级
	FUNC_HERO_BREAK                     = 67  //! 角色突破
	FUNC_EQUI_STRENGTHEN                = 68  //! 装备强化
	FUNC_EQUI_STRENGTHEN_FIVE           = 69  //! 装备强化五次
	FUNC_HERO_CULTURE                   = 70  //! 英雄培养
	FUNC_ARENA                          = 71  //! 竞技场
	FUNC_HERO_STORE                     = 72  //! 神将商店
	FUNC_SANGUOZHI                      = 73  //! 三国志
	FUNC_ROB_GEM                        = 74  //! 夺宝
	FUNC_COPY_SWEEP                     = 75  //! 副本扫荡
	FUNC_SANGUOWUSHUANG                 = 76  //! 三国无双
	FUNC_DAILY_COPY                     = 77  //! 日常副本
	FUNC_FAMOUS_COPY                    = 78  //! 名将副本
	FUNC_EQUI_REFINED                   = 79  //! 装备精炼
	FUNC_MINING_OPEN                    = 80  //! 挖矿功能
	FUNC_TERRITORY                      = 81  //! 领地征讨
	FUNC_REBEL_SIEGE                    = 82  //! 围剿叛军
	FUNC_GUILD                          = 83  //! 公会
	FUNC_ELITE_COPY                     = 84  //! 精英副本
	FUNC_ELITE_ENEMY                    = 85  //! 精英外敌
	FUNC_HERO_AWKEN                     = 86  //! 英雄觉醒
	FUNC_AWAKEN_STORE                   = 87  //! 觉醒商店
	FUNC_HEROSOULS_TIMES                = 88  //! 每日将灵,名将试炼可刷新次数
	FUNC_EQUIP_BAG_CAPACITY             = 89  //! 装备背包容量数
	FUNC_DISCOUNT_SELL                  = 90  //! 折扣贩卖
	FUNC_BACK_POS_BEGIN                 = 91  //! 援军第一格开放
	FUNC_BACK_POS_END                   = 96  //! 援军第六格开放
	FUNC_BACK_CHEER                     = 97  //! 援军助威
	FUNC_BATTLE_PET                     = 98  //! 战宠
	FUNC_DESTINY                        = 99  //! 天命
	FUNC_GEM_REFINE                     = 100 //! 宝物精炼
	FUNC_SUMMON                         = 101 //! 商城召唤
	FUNC_HEROGOD                        = 102 //! 化神系统
	FUNC_MAIN_COPY                      = 103 //! 主线副本
	FUNC_ARENA_STORE                    = 104 //! 竞技场商店
	FUNC_SGWS_STORE                     = 105 //! 三国无双商店
	FUNC_REBEL_STORE                    = 106 //! 围剿叛军商店
	FUNC_CHARGE                         = 107 //! 充值
	FUNC_TEAM                           = 108 //! 阵容
	FUNC_MALL_ITEM                      = 109 //! 商城道具
	FUNC_MALL_GIFT                      = 110 //! 商城礼包
	FUNC_GLYPH_REFINE                   = 111 //! 雕文洗练
	FUNC_HANGUP_GRID_OPNE               = 112 //! 挂机系统格子开启花费
	FUNC_HANGUP_QUICKTIME               = 113 //! 挂机系统快速战斗次数
	FUNC_SCORE_FIGHT_TIME               = 114 //! 积分赛战斗次数
	FUNC_FOODWAR_ATTACK_TIMES           = 115 //! 夺粮战购买掠夺次数
	FUNC_FOODWAR_REVENGE_TIMES          = 116 //! 夺粮战购买复仇次数
	FUNC_FRIEND_NUM_LIMIT               = 117 //! 好友数量上限
	FUNC_EQUI_STAR                      = 118 //! 装备升星
	FUNC_PET_POS_BEGIN                  = 119 //! 战宠第一格开放
	FUNC_PET_POS_END                    = 124 //! 战宠第六格开放
	FUNC_PET_STORE                      = 125 //! 宠物商店
	FUNC_HEROSOULS_STORE                = 126 //! 英魂商店
	FUNC_BLACK_STORE                    = 127 //! 黑市商店
	FUNC_WANDER                         = 128 //! 云游系统
	FUNC_PET_REFINE                     = 129 //! 宠物神练
	FUNC_AWARDCENTER                    = 130 //! 奖励中心
	FUNC_MAIL                           = 131 //! 邮件
	FUNC_FRIEND                         = 132 //! 好友
	FUNC_FAMOUSHALL                     = 133 //! 名人堂
	FUNC_BAG                            = 134 //! 背包
	FUNC_CAMPBAT                        = 135 //! 阵营战
	FUNC_DAILYTASK                      = 136 //! 日常任务
	FUNC_PET_STAR                       = 137 //! 宠物升星
	FUNC_CHAT                           = 138 //! 聊天
	FUNC_SCORE_SYSTEM                   = 139 //! 积分赛系统
	FUNC_TOP_RACE                       = 140 //! 争霸赛
	FUNC_FASHION_COMPROMISES            = 141 //! 时装融合
	FUNC_HOLY_MELTING                   = 142 //! 圣物融合

	//功能ID结束标记
	FUNC_END_ID = 139 //!功能ID结束标记
)

package main

import (
	"gamesvr/mainlogic"
	"gamesvr/tcpclient"
	"msg"
	"net/http"
	"utility"
)

func RegHttpMsgHandler() {
	type THttpFuncInfo struct {
		url      string
		fun      func(http.ResponseWriter, *http.Request)
		needLock bool
	}
	var configSlice = []THttpFuncInfo{
		//! 任务的消息处理
		{"/get_tasks", mainlogic.Hand_GetAllTask, false},                //! 玩家请求全部的任务信息
		{"/receive_task", mainlogic.Hand_GetTaskAward, false},           //! 玩家请求领取完成任务奖励
		{"/receive_taskscore", mainlogic.Hand_GetTaskScoreAward, false}, //! 玩家请求领取任务积分宝箱奖励

		//玩家登录离开创建角色处理
		{"/user_login_game", mainlogic.Hand_PlayerLoginGame, false},   //玩家登录游戏服
		{"/user_enter_game", mainlogic.Hand_PlayerEnterGame, false},   //玩家进入游戏服
		{"/user_leave_game", mainlogic.Hand_PlayerLeaveGame, false},   //玩家离开游戏服
		{"/create_new_player", mainlogic.Hand_CreateNewPlayer, false}, //玩家创建角色
		{"/query_server_time", mainlogic.Hand_QueryServerTime, false}, //玩家查询服务器时间
		{"/get_login_data", mainlogic.Hand_GetLoginData, false},       //玩家加载登录信息

		//! 成就的消息处理
		{"/receive_achievement", mainlogic.Hand_GetAchievementAward, false}, //! 玩家请求领取成就奖励

		//! 签到的消息处理
		{"/get_sign", mainlogic.Hand_GetSignInfo, false}, //! 玩家请求获取签到信息
		{"/daily_sign", mainlogic.Hand_DailySign, false}, //! 玩家请求获取日常签到奖励
		{"/sign_plus", mainlogic.Hand_SignPlus, false},   //! 玩家请求获取豪华签到奖励

		//! 神将商店的消息处理
		{"/get_store", mainlogic.Hand_GetHeroStoreInfo, false},         //! 玩家请求商店信息
		{"/refresh_store", mainlogic.Hand_RefreshHeroStore, false},     //! 玩家请求刷新神将商店
		{"/store_buy", mainlogic.Hand_HeroStore_Buy, false},            //! 玩家请求购买神将商店商品
		{"/get_all_store_data", mainlogic.Hand_GetAllStoreInfo, false}, //! 玩家请求所有商店信息

		//! 商城消息处理
		{"/get_vip_gift", mainlogic.Hand_GetVipGiftInfo, false},              //! 玩家请求礼包信息
		{"/buy_vip_gift", mainlogic.Hand_BuyVipGift, false},                  //! 玩家请求购买礼包
		{"/buy_vip_week_gift", mainlogic.Hand_BuyVipWeekGift, false},         //! 玩家请求购买VIP每周礼包
		{"/get_goods_buy_times", mainlogic.Hand_GetMallGoodsBuyTimes, false}, //! 玩家请求道具商城商品购买次数
		{"/get_items_buy_times", mainlogic.Hand_GetMallGoodsBuyTimes, false}, //! 玩家请求道具商城商品购买次数(复制)
		{"/buy_goods", mainlogic.Hand_BuyMallGoods, false},                   //! 玩家请求购买道具商城商品
		{"/get_summon_status", mainlogic.Hand_GetSummonStatus, false},        //! 玩家请求召唤信息状态
		{"/get_summon", mainlogic.Hand_GetSummon, false},                     //! 玩家请求召唤武将
		{"/exchange_hero", mainlogic.Hand_ExchangeHero, false},               //! 玩家请求积分兑换英雄

		//! Vip的消息处理
		{"/query_vip_welfare", mainlogic.Hand_GetDailyVipStatus, false}, //! 玩家请求VIP日常福利领取状态
		{"/get_vip_welfare", mainlogic.Hand_DailyVipWelfare, false},     //! 玩家请求领取VIP日常福利

		//! 副本类消息
		{"/battle_check", mainlogic.Hand_BattleCheck, false},                       //! 检查挑战副本的条件
		{"/battle_result", mainlogic.Hand_BattleResult, false},                     //! 检查挑战副本的条件
		{"/sweep_copy", mainlogic.Hand_SweepCopy, false},                           //! 扫荡副本
		{"/get_main_star_award", mainlogic.Hand_GetMainStarAward, false},           //! 玩家请求主线关卡星级奖励
		{"/get_main_scene_award", mainlogic.Hand_GetMainSceneAward, false},         //! 玩家请求主线关卡场景奖励
		{"/get_rebel_find_info", mainlogic.Hand_GetRebelFindInfo, false},           //! 获取发现叛军简单信息
		{"/get_main_reset_times", mainlogic.Hand_GetMainResetTimes, false},         //! 玩家查询主线关卡重置次数
		{"/reset_main_battletimes", mainlogic.Hand_ResetMainBattleTimes, false},    //! 玩家重置主线关卡挑战次数
		{"/get_famous_detail", mainlogic.Hand_GetFamousCopyDetailInfo, false},      //! 获取名将副本详细信息
		{"/get_famous_award", mainlogic.Hand_GetFamousCopyAward, false},            //! 获取名将副本章节奖励
		{"/get_elite_star_award", mainlogic.Hand_GetEliteStarAward, false},         //! 获取精英关卡星级奖励
		{"/get_elite_scene_award", mainlogic.Hand_GetEliteSceneAward, false},       //! 获取精英关卡场景奖励
		{"/get_elite_reset_times", mainlogic.Hand_GetEliteResetTimes, false},       //! 获取精英关卡重置次数
		{"/reset_elite_battletimes", mainlogic.Hand_ResetEliteBattleTimes, false},  //! 重置精英关卡挑战
		{"/attack_elite_invade", mainlogic.Hand_AttackInvade, false},               //! 请求攻击精英关卡入侵
		{"/get_elite_invade_status", mainlogic.Hand_GetEliteCopyInvadeInfo, false}, //! 请求入侵消息
		{"/get_copy_data", mainlogic.Hand_GetCopyData, false},                      //! 获取副本数据

		//! 三国志消息
		{"/get_sanguozhi_info", mainlogic.Hand_SanGuoZhiInfo, false},             //! 玩家请求获取三国志信息
		{"/set_sanguozhi", mainlogic.Hand_SetSanGuoZhi, false},                   //! 玩家请求三国志命星
		{"/get_sanguozhi_attr", mainlogic.Hand_GetSanGuoStarAddAttribute, false}, //! 玩家请求获取三国志加属性信息

		//! 竞技场消息处理
		{"/arena_result", mainlogic.Hand_ChallengeArenaResult, false},                      //! 玩家反馈挑战竞技场结果
		{"/get_arena_info", mainlogic.Hand_GetArenaInfo, false},                            //! 玩家请求竞技场可挑战玩家信息
		{"/arena_check", mainlogic.Hand_ArenaCheck, false},                                 //! 玩家请求挑战排位检测
		{"/arena_store_buy_item", mainlogic.Hand_BuyArenaStoreItem, false},                 //! 玩家请求购买声望商店物品
		{"/arena_store_query_award", mainlogic.Hand_GetArenaStoreAleadyBuyAwardLst, false}, //! 玩家请求已购买的声望商店奖励ID列表
		{"/arena_battle", mainlogic.Hand_ArenaBattle, false},                               //! 多次挑战竞技场

		//! 夺宝消息处理
		{"/get_rob_list", mainlogic.Hand_GetRobList, false},            //! 玩家请求抢劫名单
		{"/refresh_rob_list", mainlogic.Hand_GetRobList, false},        //! 玩家请求刷新抢劫名单
		{"/rob_treasure", mainlogic.Hand_RobTreasure, false},           //! 玩家请求夺宝
		{"/get_free_war_time", mainlogic.Hand_GetFreeWarTime, false},   //! 玩家请求免战时间
		{"/get_rob_hero_info", mainlogic.Hand_GetRobHeroInfo, false},   //! 获取抢劫玩家武将信息
		{"/treasure_composed", mainlogic.Hand_TreasureComposed, false}, //! 玩家请求合成宝物
		{"/treasure_melting", mainlogic.Hand_TreasureMelting, false},   //! 玩家请求宝物熔炼

		//! 三国无双消息处理
		{"/get_sangokumusou_star", mainlogic.Hand_GetSangokuMusou_StarInfo, false},                   //! 获取星数信息
		{"/get_sangokumusou_copy", mainlogic.Hand_GetSangokuMuSou_CopyInfo, false},                   //! 获取闯关信息
		{"/get_sangokumusou_elite_copy", mainlogic.Hand_GetSangokuMusou_EliteCopy, false},            //! 获取精英挑战闯关信息
		{"/pass_sangokumusou", mainlogic.Hand_PassSangokuMusou_Copy, false},                          //! 通关三国无双回馈
		{"/pass_sgws_elite", mainlogic.Hand_PassSangokuMusou_EliteCopy, false},                       //! 通关三国无双精英挑战
		{"/sweep_sangoumusou", mainlogic.Hand_SangokuMusou_Sweep, false},                             //! 请求三星扫荡章节
		{"/get_sangokumusou_chapter_award", mainlogic.Hand_GetSangokuMusou_ChapterAward, false},      //! 请求章节奖励
		{"/get_sangokumusou_attr", mainlogic.Hand_GetSangokuMusou_ChapterAttr, false},                //! 请求随机三个属性奖励
		{"/set_sangokumusou_attr", mainlogic.Hand_SetSangokuMusou_ChapterAttr, false},                //! 选择属性奖励
		{"/get_sangokumusou_all_attr", mainlogic.Hand_GetSangokuMusou_Attr, false},                   //! 获取当前所有属性奖励
		{"/get_sangokumusou_treasure", mainlogic.Hand_GetSangokuMusou_Treasure, false},               //! 获取无双秘藏
		{"/buy_sangokumusou_treasure", mainlogic.Hand_BuySangokuMusou_Treasure, false},               //! 购买无双秘藏
		{"/reset_sangokumusou_copy", mainlogic.Hand_SangokuMusou_ResetCopy, false},                   //! 重置挑战
		{"/get_sangoukumusou_elite_add_times", mainlogic.Hand_GetSangokuMusou_AddEliteCopy, false},   //! 获取可增加精英挑战次数
		{"/add_sangoukumusou_elite_copy", mainlogic.Hand_SangokuMusou_AddEliteCopy, false},           //! 增加精英挑战次数
		{"/get_sangokumusou_store_aleady_buy", mainlogic.Hand_GetSangokuMusouStore_AleadyBuy, false}, //! 获取无双商店购买次数信息
		{"/buy_sangokumusou_store", mainlogic.Hand_GetSangokuMusou_StoreItem, false},                 //! 购买无双商店商品
		{"/get_sangokumusou_status", mainlogic.Hand_GetSanguowsStatus, false},                        //! 获取三国无双状态

		//! 领地征讨消息处理
		{"/get_territory_status", mainlogic.Hand_GetTerritoryStatus, false},              //! 玩家请求获取自身领地状态
		{"/challenge_territory", mainlogic.Hand_ChallengeTerritory, false},               //! 玩家回馈挑战领地成功结果
		{"/get_friend_territory_status", mainlogic.Hand_GetFriendTerritoryStatus, false}, //! 玩家请求获取好友领地状态
		{"/get_friend_territory_info", mainlogic.Hand_GetFriendTerritoryDetail, false},   //! 玩家请求获取好友领地详情
		{"/help_riot", mainlogic.Hand_SuppressRiot, false},                               //! 玩家请求帮助好友镇压暴动
		{"/get_territory_award", mainlogic.Hand_GetTerritoryAward, false},                //! 玩家请求收获巡逻领地奖励
		{"/patrol_territory", mainlogic.Hand_PatrolTerritory, false},                     //! 玩家请求领地放置武将巡逻
		{"/territory_skill_up", mainlogic.Hand_TerritorySkillLevelUp, false},             //! 玩家请求提升领地技能等级
		// ("/query_territory_award", mainlogic.Hand_GetTerritoryAwardLst)           //! 玩家请求查询领地巡逻奖励
		{"/query_territory_riot", mainlogic.Hand_GetTerritoryRiotInfo, false}, //! 玩家请求查询领地暴动信息

		//! 围剿叛军消息处理
		{"/get_rebel_info", mainlogic.Hand_GetRebelInfo, false},                    //! 玩家请求获取叛军
		{"/attack_rebel", mainlogic.Hand_AttackRebel, false},                       //! 玩家请求攻击叛军
		{"/share_rebel", mainlogic.Hand_ShareRebel, false},                         //! 玩家请求分享叛军
		{"/get_exploit_award_status", mainlogic.Hand_GetExploitAwardStatus, false}, //! 获取功勋奖励领奖状态
		{"/get_exploit_award", mainlogic.Hand_GetExploitAward, false},              //! 玩家请求领取功勋奖励
		{"/buy_rebel_store", mainlogic.Hand_BuyRebelStore, false},                  //! 玩家请求购买战功商店物品

		//! 活动消息处理
		{"/get_seven_activity", mainlogic.Hand_GetSevenActivityInfo, false},                       //! 玩家请求获取七日活动信息
		{"/get_seven_activity_award", mainlogic.Hand_GetSevenActivityAward, false},                //! 玩家获取七日活动奖励
		{"/buy_seven_activity_limit", mainlogic.Hand_BuySevenActivityLimit, false},                //! 玩家购买七日活动限购商品
		{"/get_activity", mainlogic.Hand_GetActivity, false},                                      //! 玩家请求当前开启活动
		{"/get_activity_list", mainlogic.Hand_GetActivity, false},                                 //! 玩家请求当前开启活动
		{"/get_open_server_day", mainlogic.Hand_GetServerOpenDay, false},                          //! 获取服务器开启天数
		{"/query_activity_login", mainlogic.Hand_QueryActivityLoginInfo, false},                   //! 查询累计登录活动信息
		{"/query_activity_firstrecharge", mainlogic.Hand_QueryActivityFirstRechargeInfo, false},   //! 查询首冲活动信息
		{"/query_monthcard_days", mainlogic.Hand_QueryActivityMonthCardDays, false},               //! 查询月卡剩余天数
		{"/query_activity_action", mainlogic.Hand_QueryActivityActionInfo, false},                 //! 查询领取体力活动信息
		{"/query_activity_moneygod", mainlogic.Hand_QueryActivityMoneyGodInfo, false},             //! 查询迎财神活动信息
		{"/query_activity_discountsale", mainlogic.Hand_QueryActivityDiscountSaleInfo, false},     //! 查询折扣贩售活动信息
		{"/query_activity_totalrecharge", mainlogic.Hand_QueryActivityTotalRechargeInfo, false},   //! 查询累计充值活动信息
		{"/query_activity_singlerecharge", mainlogic.Hand_QueryActivitySingleRechargeInfo, false}, //! 查询单笔充值活动信息
		{"/get_login_award", mainlogic.Hand_GetActivityLoginAward, false},                         //! 领取累计登录奖励
		{"/get_first_recharge", mainlogic.Hand_GetFirstRechargeAward, false},                      //! 领取首冲奖励
		{"/get_activity_action", mainlogic.Hand_ReceiveActivityAction, false},                     //! 领取体力奖励
		{"/welcome_money_god", mainlogic.Hand_WelcomeMoneyGold, false},                            //! 领取迎财神奖励
		{"/get_money_god_award", mainlogic.Hand_MoneyGoldAward, false},                            //! 领取迎财神累积奖励
		{"/buy_discount_sale", mainlogic.Hand_BuyDiscountSaleItem, false},                         //! 购买折扣贩售商品
		{"/get_recharge_award", mainlogic.Hand_GetRechargeAward, false},                           //! 领取累积充值奖励
		{"/get_action_retroactive", mainlogic.Hand_ActionRetroactive, false},                      //! 补签体力活动
		{"/get_single_award", mainlogic.Hand_GetSingleRechargeAward, false},                       //! 获取单充奖励
		{"/query_hunt_treasure", mainlogic.Hand_QueryHuntTreasure, false},                         //! 查询巡回探宝状态
		{"/start_hunt", mainlogic.Hand_StartHuntTreasure, false},                                  //! 开始巡回探宝掷骰
		{"/query_hunt_award", mainlogic.Hand_QueryHuntTurnsAward, false},                          //! 查询巡回探宝奖励领取情况
		{"/get_hunt_award", mainlogic.Hand_GetHuntTurnsAward, false},                              //! 获取巡回探宝奖励
		{"/query_hunt_store", mainlogic.Hand_QueryHuntTreasureStore, false},                       //! 玩家查询巡回商店
		{"/buy_hunt_store", mainlogic.Hand_BuyHuntTreasureStroreItem, false},                      //! 玩家购买巡回商店物品
		{"/query_lucky_wheel", mainlogic.Hand_QueryLuckyWheel, false},                             //! 查询幸运转盘
		{"/rotating_wheel", mainlogic.Hand_RotatingWheel, false},                                  //! 申请转动轮盘
		{"/get_group_purchase_info", mainlogic.Hand_GetGroupPurchaseInfo, false},                  //! 获取团购信息
		{"/buy_group_purchase", mainlogic.Hand_BuyGroupPurchaseItem, false},                       //! 玩家请求购买团购
		{"/get_group_score_award", mainlogic.Hand_GetGroupPurchaseScoreAward, false},              //! 玩家请求积分奖励
		{"/get_festival_task", mainlogic.Hand_GetFestivalTask, false},                             //! 获取欢庆佳节任务
		{"/get_festival_task_award", mainlogic.Hand_GetFestivalTaskAward, false},                  //! 玩家请求领取欢庆佳节任务奖励
		{"/exchange_festival_award", mainlogic.Hand_ExchangeFestivalAward, false},                 //! 玩家兑换奖励
		{"/festival_discount_sale", mainlogic.Hand_BuyFestivalSaleItem, false},                    //! 玩家请求购买欢庆佳节半价限购
		{"/clean_hunt_store", mainlogic.Hand_CleanHuntStore, false},                               //! 清除巡回商店数据
		{"/get_activity_rank", mainlogic.Hand_GetActivityRank, false},                             //! 获取活动排行榜
		{"/get_activity_rank_award", mainlogic.Hand_GetActivityRankAward, false},                  //! 获取活动排行榜奖励
		{"/get_week_award_status", mainlogic.Hand_GetWeekAwardStatus, false},                      //! 获取周周盈状态
		{"/get_week_award", mainlogic.Hand_GetWeekAward, false},                                   //! 获取周周盈奖励
		{"/get_level_gift_info", mainlogic.Hand_GetLevelGiftInfo, false},                          //! 获取等级礼包信息
		{"/buy_level_gift", mainlogic.Hand_BuyLevelGift, false},                                   //! 购买等级礼包
		{"/get_rank_gift_info", mainlogic.Hand_GetRankGiftInfo, false},                            //! 获取排名礼包信息
		{"/buy_rank_gift", mainlogic.Hand_BuyRankGift, false},                                     //! 购买排名礼包
		{"/get_monthfund_status", mainlogic.Hand_GetMonthFundStatus, false},                       //! 获取月基金状态
		{"/receive_month_fund", mainlogic.Hand_ReceiveMonthFund, false},                           //! 玩家领取月基金奖励
		{"/get_limit_daily_task", mainlogic.Hand_GetLimitDailyTaskInfo, false},                    //! 获取限时日常任务信息
		{"/get_limit_daily_award", mainlogic.Hand_GetLimitDailyTaskAward, false},                  //! 获取限时日常奖励
		{"/get_limit_sale_info", mainlogic.Hand_GetLimitSaleItemInfo, false},                      //! 查询限时特惠物品信息
		{"/buy_limit_sale_item", mainlogic.Hand_BuyLimitSaleItem, false},                          //! 购买限时特惠物品
		{"/get_limitsale_all_award", mainlogic.Hand_GetLimitSaleAllAward, false},                  //! 获取限时特惠全民奖励

		//! 挖矿消息处理
		{"/enter_mining", mainlogic.Hand_GetMiningInfo, false},                             //! 获取挖矿信息
		{"/mining_get_award", mainlogic.Hand_GetRandBossAward, false},                      //! 获取打败Boss后九个翻牌奖励
		{"/select_mining_award", mainlogic.Hand_SelectBossAward, false},                    //! 选择打败Boss后九个翻牌奖励
		{"/mining_guaji", mainlogic.Hand_MiningGuaji, false},                               //! 挖矿挂机
		{"/mining_guaji_time", mainlogic.Hand_GetMiningGuajiTime, false},                   //! 查询挖矿挂机倒计时
		{"/mining_guaji_award", mainlogic.Hand_GetMiningGuajiAward, false},                 //! 获取挂机收益
		{"/mining_dig", mainlogic.Hand_MiningDig, false},                                   //! 玩家请求获取挖地图某个点的信息
		{"/mining_element_stone", mainlogic.Hand_MiningElement_RefiningStone, false},       //! 挖矿元素-精炼石
		{"/mining_event_action_award", mainlogic.Hand_MiningEvent_ActionAward, false},      //! 挖矿事件-行动力奖励
		{"/mining_event_black_market", mainlogic.Hand_MiningEvent_BalckMarket, false},      //! 挖矿事件-黑市
		{"/mining_buy_black_market", mainlogic.Hand_MiningEvent_BuyBlackMarketItem, false}, //! 挖矿事件-黑市购买
		{"/mining_event_monster", mainlogic.Hand_MiningEvent_Monster, false},               //! 挖矿事件-怪物
		{"/mining_event_treasure", mainlogic.Hand_MiningEvent_Treasure, false},             //! 挖矿事件-宝箱
		{"/mining_event_box", mainlogic.Hand_MiningEvent_MagicBox, false},                  //! 挖矿事件-魔盒
		{"/mining_event_scan", mainlogic.Hand_MiningEvent_Scan, false},                     //! 挖矿事件-扫描
		{"/mining_event_question", mainlogic.Hand_MiningEvent_Question, false},             //! 挖矿事件-答题
		{"/mining_event_buff", mainlogic.Hand_MiningEvent_Buff, false},                     //! 挖矿事件-Buff
		{"/get_mining_status", mainlogic.Hand_GetMiningStatusCode, false},                  //! 获取挖矿状态码
		{"/mining_event_monster_info", mainlogic.Hand_MiningEvent_MonsterInfo, false},      //! 获取怪物信息

		//! 八卦镜
		{"/use_baguajing", mainlogic.Hand_UseBaGuaJing, false}, //! 玩家请求使用八卦镜

		//! 开服基金
		{"/get_fund_status", mainlogic.Hand_GetOpenFundStatus, false},          //! 请求查询购买基金状态
		{"/buy_fund", mainlogic.Hand_BuyOpenFund, false},                       //! 请求购买开服基金
		{"/get_fund_all_award", mainlogic.Hand_GetOpenFundAllAward, false},     //! 请求领取基金奖励-全服奖励
		{"/get_func_level_award", mainlogic.Hand_GetOpenFundLevelAward, false}, //! 请求领取基金奖励-等级返利

		//! 领奖中心
		{"/query_award_center", mainlogic.Hand_GetAwardCenterInfo, false},
		{"/get_award_center", mainlogic.Hand_RecvAwardCenter, false},
		{"/onekey_award_center", mainlogic.Hand_RecvAwardCenterAwardOneyKey, false},

		//! 公会协议
		{"/create_guild", mainlogic.Hand_CreateGuild, false},
		{"/get_guild", mainlogic.Hand_GetGuild, false},
		{"/get_guild_lst", mainlogic.Hand_GetMoreGuild, false},
		{"/enter_guild", mainlogic.Hand_EnterGuild, false},
		{"/get_apply_guild_list", mainlogic.Hand_GetApplyGuildList, false},
		{"/get_apply_guild_member_list", mainlogic.Hand_GetApplyGuildMemberList, false},
		{"/apply_through", mainlogic.Hand_ApplicationThrough, false},
		{"/leave_guild", mainlogic.Hand_ExitGuild, false},
		{"/get_sacrifice_status", mainlogic.Hand_GetSacrificeStatus, false},
		{"/guild_sacrifice", mainlogic.Hand_GuildSacrifice, false},
		{"/get_sacrifice_award", mainlogic.Hand_GetSacrificeAward, false},
		{"/query_guild_store", mainlogic.Hand_GetGuildStoreInfo, false},
		{"/buy_guild_store", mainlogic.Hand_BuyGuildItem, false},
		{"/attack_guild_copy", mainlogic.Hand_AttackGuildCopy, false},
		{"/get_guild_copy_status", mainlogic.Hand_GetGuildCopyStatus, false},
		{"/get_guild_copy_award", mainlogic.Hand_GetGuildCopyTreasure, false},
		{"/query_recv_copy_award", mainlogic.Hand_QueryGuildCopyTreasure, false},
		{"/query_guild_copy_rank", mainlogic.Hand_QueryGuildCopyRank, false},
		{"/query_guild_msg_board", mainlogic.Hand_QueryGuildMsgBoard, false},
		{"/remove_guild_msg_board", mainlogic.Hand_RemoveGuildMsgBoard, false},
		{"/write_guild_msg_board", mainlogic.Hand_UseGuildMsgBoard, false},
		{"/kick_member", mainlogic.Hand_KickGuildMember, false},
		{"/update_guild_name", mainlogic.Hand_UpdateGuildName, false},
		{"/update_guild_info", mainlogic.Hand_UpdateGuildInfo, false},
		{"/research_guild_skill", mainlogic.Hnad_ResearchGuildSkill, false},
		{"/study_guild_skill", mainlogic.Hand_StudyGuildSkill, false},
		{"/get_guild_skill", mainlogic.Hand_GetGuildSkillInfo, false},
		{"/get_guild_skill_limit", mainlogic.Hand_GetGuildSkillResearchInfo, false},
		{"/get_guild_chapter_status", mainlogic.Hand_GetGuildChapterRecvLst, false},
		{"/get_guild_chapter_award", mainlogic.Hand_GetGuildChapterAward, false},
		{"/get_guild_member_list", mainlogic.Hand_GetGuildMemberList, false},
		{"/get_player_info", mainlogic.Hand_GetPlayerInfo, false},
		{"/get_guild_log", mainlogic.Hand_GetGuildLog, false},
		{"/change_guild_role", mainlogic.Hand_ChangeGuildMemberPose, false},
		//("/guild_levelup", mainlogic.Hand_GuildLevelUp) //! 改变逻辑,暂时不用
		{"/get_guild_chapter_award_all", mainlogic.Hand_GetAllGuildChapterAward, false},
		{"/update_guild_backstatus", mainlogic.Hand_UpdateGuildChapterBackStatus, false},
		{"/search_guild", mainlogic.Hand_SearchGuild, false},
		{"/cancellation_guild_apply", mainlogic.Hand_CancellationGuildApply, false},
		{"/get_guild_status", mainlogic.Hand_GetGuildStatus, false},

		//! 黑市协议
		{"/get_black_market_info", mainlogic.Hand_GetBlackMarketInfo, false},
		{"/get_black_market_status", mainlogic.Hand_GetBlackMarketStatus, false},
		{"/buy_black_market", mainlogic.Hand_BuyBlackMarketGoods, false},

		//! 名人堂协议
		{"/send_flower", mainlogic.Hand_SendFlower, false},
		{"/get_charm", mainlogic.Hand_GetCharmValue, false},

		//! 称号协议
		{"/get_title", mainlogic.Hand_GetTitle, false},
		{"/activate_title", mainlogic.Hand_ActivateTitle, false},
		{"/equi_title", mainlogic.Hand_EquipTitle, false},

		//! 夺粮战
		{"/get_foodwar_challenger", mainlogic.Hand_FoodWar_GetChallenger, false},
		{"/get_foodwar_time", mainlogic.Hand_FoodWar_GetTime, false},
		{"/get_foodwar_revenge_status", mainlogic.Hand_FoodWar_RevengeStatus, false},
		{"/rob_food", mainlogic.Hand_RobFood, false},
		{"/get_food_rank", mainlogic.Hand_FoodWar_GetRank, false},
		{"/buy_food_times", mainlogic.Hand_FoodWar_BuyTimes, false},
		{"/recv_food_award", mainlogic.Hand_FoodWar_RecvAward, false},
		{"/query_food_award", mainlogic.Hand_FoodWar_QueryAward, false},
		{"/get_foodwar_status", mainlogic.Hand_FoodWar_GetStatus, false},
		{"/revenge_rob", mainlogic.Hand_FoodWar_Revenge, false},

		//! 英魂
		{"/activate_herosouls", mainlogic.Hand_ActivateHeroSouls, false},
		{"/query_herosouls_chapter", mainlogic.Hand_QueryChapterHeroSoulsDetail, false},
		{"/get_herosouls_lst", mainlogic.Hand_GetHeroSoulsLst, false},
		{"/refresh_herosouls", mainlogic.Hand_RefreshHeroSoulsLst, false},
		{"/challenge_herosouls", mainlogic.Hand_ChallengeHeroSouls, false},
		{"/buy_challenge_herosouls", mainlogic.Hand_BuyChallengeHeroSoulsTimes, false},
		{"/reset_herosouls_lst", mainlogic.Hand_ResetHeroSoulsLst, false},
		{"/query_herosouls_rank", mainlogic.Hand_QueryHeroSoulsRank, false},
		{"/query_herosouls_store", mainlogic.Hand_QueryHeroSoulsStoreInfo, false},
		{"/buy_herosouls", mainlogic.Hand_BuyHeroSoulsStoreItem, false},
		{"/query_herosouls_achievement", mainlogic.Hand_QuerySoulMapInfo, false},
		{"/activate_herosouls_achievement", mainlogic.Hand_ActivateheroSoulsAchievement, false},
		{"/query_herosouls_property", mainlogic.Hand_QueryHeroSoulsPerproty, false},

		//英雄消息处理
		{"/get_battle_data", mainlogic.Hand_GetBattleData, false},     //玩家请求上阵数据
		{"/upgrade_hero", mainlogic.Hand_UpgradeHero, false},          //升级英雄(非主角)
		{"/change_hero", mainlogic.Hand_ChangeHero, false},            //玩家更换英雄
		{"/change_back_hero", mainlogic.Hand_ChangeBackHero, false},   //玩家更换援军英雄
		{"/breakout_hero", mainlogic.Hand_BreakOutHero, false},        //玩家突破英雄
		{"/culture_hero", mainlogic.Hand_CultureHero, false},          //玩家培养英雄
		{"/compose_hero", mainlogic.Hand_ComposeHero, false},          //玩家天命英雄
		{"/upgod_hero", mainlogic.Hand_UpgodHero, false},              //玩家化神英雄
		{"/change_career", mainlogic.Hand_Change_Career, false},       //玩家更改职业
		{"/set_wake_item", mainlogic.Hand_SetWakeItem, false},         //玩家设置觉醒道具
		{"/up_wake_level", mainlogic.Hand_UpWakeLevel, false},         //玩家提升觉醒等级
		{"/compose_wake_item", mainlogic.Hand_ComposeWakeItem, false}, //玩家合成觉醒道具等级
		{"/query_destiny", mainlogic.Hand_QueryHeroDestiny, false},    //玩家查询天命状态
		{"/destiny_hero", mainlogic.Hand_DestinyHero, false},          //玩家天命英雄
		{"/upgrade_diaowen", mainlogic.Hand_UpgradeDiaoWen, false},    //玩家升品雕文
		{"/xilian_diaowen", mainlogic.Hand_XiLianDiaoWen, false},      //玩家洗炼雕文
		{"/xilian_tihuan", mainlogic.Hand_XiLianTiHuan, false},        //玩家洗炼替换雕文
		{"/upgrade_pet", mainlogic.Hand_UpgradePet, false},            //升级宠物
		{"/upstar_pet", mainlogic.Hand_UpstarPet, false},              //升星宠物
		{"/upgod_pet", mainlogic.Hand_UpgodPet, false},                //神炼宠物
		{"/change_pet", mainlogic.Hand_ChangePet, false},              //更换宠物
		{"/unset_pet", mainlogic.Hand_UnsetPet, false},                //下阵宠物
		{"/compose_pet", mainlogic.Hand_ComposePet, false},            //装备合成

		//装备
		{"/change_equip", mainlogic.Hand_ChangeEquip, false},         //装备更换
		{"/equip_strengthen", mainlogic.Hand_EquipStrengthen, false}, //装备强化
		{"/equip_refine", mainlogic.Hand_EquipRefine, false},         //装备精炼
		{"/equip_risestar", mainlogic.Hand_EquipRiseStar, false},     //装备升星
		{"/compose_equip", mainlogic.Hand_ComposeEquip, false},       //装备合成

		//宝物
		{"/change_gem", mainlogic.Hand_ChangeGem, false},         //宝物更换
		{"/gem_strengthen", mainlogic.Hand_GemStrengthen, false}, //宝物强化
		{"/gem_refine", mainlogic.Hand_GemRefine, false},         //宝物精炼

		//时装
		{"/fashion_set", mainlogic.Hand_FashionSet, false},           //时装装备
		{"/fashion_strength", mainlogic.Hand_FashionStrength, false}, //时装强化
		{"/fashion_recast", mainlogic.Hand_FashionRecast, false},     //时装重铸
		{"/fashion_compose", mainlogic.Hand_FashionCompose, false},   //时装合成
		{"/fashion_melting", mainlogic.Hand_FashionMelting, false},   //时装熔炼

		//背包
		{"/get_bag_data", mainlogic.Hand_GetBagData, false},              //请求背包数据
		{"/get_bag_heros", mainlogic.Hand_GetBagHeros, false},            //请求背包中的所有英雄
		{"/get_bag_equips", mainlogic.Hand_GetBagEquips, false},          //请求背包中的所有装备
		{"/get_bag_hero_piece", mainlogic.Hand_GetBagHerosPiece, false},  //请求背包中的所有英雄碎片
		{"/get_bag_equip_piece", mainlogic.Hand_GetBagEquipPiece, false}, //请求背包中的所有装备碎片
		{"/get_bag_gem_piece", mainlogic.Hand_GetBagGemPiece, false},     //请求背包中的所有宝物碎片
		{"/get_bag_gems", mainlogic.Hand_GetBagGems, false},              //请求背包中的所有的宝物
		{"/get_bag_items", mainlogic.Hand_GetBagItems, false},            //请求背包里的道具
		{"/get_bag_wake_items", mainlogic.Hand_GetBagWakeItems, false},   //请求背包里的觉醒道具
		{"/get_bag_pets", mainlogic.Hand_GetBagPets, false},              //请求背包中的所有宠物
		{"/get_bag_pet_piece", mainlogic.Hand_GetBagPetsPiece, false},    //请求背包中的所有宠物碎片
		{"/use_item", mainlogic.Hand_UseItem, false},                     //使用背包里的道具
		{"/sell_item", mainlogic.Hand_SellItem, false},                   //使用背包里的道具

		//角色信息
		{"/get_role_data", mainlogic.Hand_GetRoleData, false},            //玩家获取角色数据
		{"/levelup_notify", mainlogic.Hand_LevelUpNotify, false},         //玩家等级升级通知
		{"/change_role_name", mainlogic.Hand_ChangeRoleName, false},      //更改角色名字
		{"/get_new_wizard", mainlogic.Hand_GetNewWizard, false},          //读取新手向导
		{"/set_new_wizard", mainlogic.Hand_SetNewWizard, false},          //设置新手向导
		{"/get_collection_heros", mainlogic.Hand_GetCollectHeros, false}, //获取玩家收集过的英雄

		//回收
		{"/query_hero_decompose_cost", mainlogic.Hand_QueryHeroDecomposeCost, false},   //! 查询分解英雄材料
		{"/decompose_hero", mainlogic.Hand_DecomposeHero, false},                       //! 分解英雄
		{"/query_hero_relive_cost", mainlogic.Hand_QueryHeroRelive, false},             //! 查询重生英雄材料
		{"/relive_hero", mainlogic.Hand_ReliveHero, false},                             //! 重生英雄
		{"/query_equip_decompose_cost", mainlogic.Hand_QueryEquipDecomposeCost, false}, //! 查询分解装备材料
		{"/decompose_equip", mainlogic.Hand_DecomposeEquip, false},                     //! 分解装备
		{"/query_equip_relive_cost", mainlogic.Hand_QueryEquipRelive, false},           //! 查询重生装备材料
		{"/relive_equip", mainlogic.Hand_ReliveEquip, false},                           //! 重生装备
		{"/query_pet_decompose_cost", mainlogic.Hand_QueryDecomposePetCost, false},     //! 查询分解战宠材料
		{"/decompose_pet", mainlogic.Hand_DecomposePet, false},                         //! 分解战宠
		{"/query_pet_relive_cost", mainlogic.Hand_QueryPetRelive, false},               //! 查询重生战宠材料
		{"/relive_pet", mainlogic.Hand_RelivePet, false},                               //! 重生战宠
		{"/query_gem_relive_cost", mainlogic.Hand_QueryGemRelive, false},               //! 查询宝物重生材料
		{"/relive_gem", mainlogic.Hand_ReliveGem, false},                               //! 重生宝物

		//挂机
		{"/hangup_get_info", mainlogic.Hand_GetHangUpInfo, false}, //请求挂机信息
		{"/hangup_set_boss", mainlogic.Hand_SetBoss, false},       //设置挂机BOSS
		{"/hangup_quick_fight", mainlogic.Hand_QuickFight, false}, //快速战斗请求
		{"/hangup_add_grid", mainlogic.Hand_AddGrid, false},       //快速战斗请求
		{"/hangup_use_exp", mainlogic.Hand_UseExpItem, false},     //一键使用经验丹

		//以下为测试消息
		{"/test_get_money", mainlogic.Hand_TestGetMoney, false},
		{"/test_get_action", mainlogic.Hand_TestGetAction, false},
		{"/test_uplevel", mainlogic.Hand_TestUplevel, false},
		{"/test_uplevel_ten", mainlogic.Hand_TestUplevelTen, false},
		{"/test_get_bag_heros", mainlogic.Hand_GetBagHeros, false},            //请求背包中的所有英雄
		{"/test_get_bag_equips", mainlogic.Hand_GetBagEquips, false},          //请求背包中的所有装备
		{"/test_get_bag_hero_piece", mainlogic.Hand_GetBagHerosPiece, false},  //请求背包中的所有英雄碎片
		{"/test_get_bag_equip_piece", mainlogic.Hand_GetBagEquipPiece, false}, //请求背包中的所有装备碎片
		{"/test_get_bag_gem_piece", mainlogic.Hand_GetBagGemPiece, false},     //请求背包中的所有宝物碎片
		{"/test_get_bag_gems", mainlogic.Hand_GetBagGems, false},              //请求背包中的所有的宝物
		{"/test_get_bag_items", mainlogic.Hand_GetBagItems, false},            //请求背包里的道具
		{"/test_get_bag_wake_items", mainlogic.Hand_GetBagWakeItems, false},   //请求背包里的觉醒道具
		{"/test_add_vip", mainlogic.Hand_TestAddVip, false},                   //请求增加Vip
		{"/test_add_guild", mainlogic.Hand_TestAddGuildExp, false},            //请求增加公会经验
		{"/test_compress", mainlogic.Hand_TestCompress, false},                //测试压缩协议
		{"/test_heros_property", mainlogic.Hand_GetHerosProperty, false},      //测试获取玩家各项属性
		{"/test_add_item", mainlogic.Hand_TestAddItem, false},
		{"/test_charge_money", mainlogic.Hand_TestAddCharge, false}, //! 测试活动充值相关
		{"/test_pass_copy", mainlogic.Hand_TestPassCopy, false},     //! 测试直接通关副本

		//充值消息处理
		{"/get_charge_info", mainlogic.Hand_GetChargeInfo, false},       //! 玩家请求充值结果
		{"/get_charge_result", mainlogic.Hand_GetChargeResult, false},   //! 玩家请求充值结果
		{"/receive_month_card", mainlogic.Hand_ReceiveMonthCard, false}, //! 玩家请求领取月卡

		//邮件系统
		{"/receive_all_mails", mainlogic.Hand_ReceiveAllMails, false}, //! 玩家请求所有的邮件

		//排行榜
		{"/get_level_rank", mainlogic.Hand_GetLevelRank, false},            //! 玩家请求等级排行榜
		{"/get_fight_rank", mainlogic.Hand_GetFightRank, false},            //! 玩家请求战力排行榜
		{"/get_sanguows_rank", mainlogic.Hand_GetSanguowsRank, false},      //! 玩家请求无双排行榜
		{"/get_arena_rank", mainlogic.Hand_GetArenaRank, false},            //! 玩家请求竞技场排行榜
		{"/get_rebel_rank", mainlogic.Hand_GetRebelRank, false},            //! 玩家请求叛军排行榜
		{"/get_guild_level_rank", mainlogic.Hand_GetGuildLevelRank, false}, //! 玩家请求公会等级排行榜
		{"/get_guild_copy_rank", mainlogic.Hand_GetGuildCopyRank, false},   //! 玩家请求公会副本排行榜
		{"/get_score_rank", mainlogic.Hand_GetScoreRank, false},            //! 玩家请求积分赛排行榜
		{"/get_wander_rank", mainlogic.Hand_GetWanderRank, false},          //! 玩家请求云游戏排行榜
		{"/get_campbat_rank", mainlogic.Hand_GetCampBatRank, false},        //! 玩家请求阵营战排行榜

		//积分赛
		{"/get_score_data", mainlogic.Hand_GetScoreData, false}, //获取积分赛主界面信息
		{"/get_score_battle_check", mainlogic.Hand_GetScoreBattleCheck, false},
		{"/set_score_battle_result", mainlogic.Hand_SetScoreBattleResult, false},
		{"/get_score_time_award", mainlogic.Hand_GetScoreTimeAward, false},
		{"/recv_score_time_award", mainlogic.Hand_RecvScoreTimeAward, false},
		{"/buy_score_fight_time", mainlogic.Hand_BuyScoreFightTime, false},
		{"/buy_score_store_item", mainlogic.Hand_BuyScoreStoreItem, false},
		{"/recv_score_continue_award", mainlogic.Hand_RecvContinueWinAward, false},
		{"/get_score_report_req", mainlogic.Hand_GetScoreBattleReport, false},

		//跨服服来请求玩家数据
		{"/select_target_player", mainlogic.Hand_SelectTargetPlayer, false},
		{"/get_fight_target", mainlogic.Hand_GetFightTarget, false},

		//社交
		{"/get_all_friend", mainlogic.Hand_GetAllFriend, false},
		{"/get_online_friend", mainlogic.Hand_GetOnlineFriend, false},
		{"/add_friend_request", mainlogic.Hand_AddFriendReq, false},
		{"/del_friend_request", mainlogic.Hand_DelFriendReq, false},
		{"/process_friend_request", mainlogic.Hand_ProcessFriendReq, false},
		{"/get_apply_list", mainlogic.Hand_GetApplyList, false},
		{"/search_friend", mainlogic.Hand_SearchFriend, false},
		{"/recomand_friend", mainlogic.Hand_RecomandFriend, false},
		{"/give_action", mainlogic.Hand_GiveAction, false},
		{"/receive_action", mainlogic.Hand_ReceiveAction, false},

		//云游
		{"/wander_getinfo", mainlogic.Hand_WanderGetInfo, false},
		{"/wander_reset", mainlogic.Hand_WanderReset, false},
		{"/wander_sweep", mainlogic.Hand_WanderSweep, false},
		{"/wander_openbox", mainlogic.Hand_WanderOpenBox, false},
		{"/wander_check", mainlogic.Hand_WanderCheck, false},
		{"/wander_result", mainlogic.Hand_WanderResult, false},

		//卡牌大师
		{"/act_card_master_draw", mainlogic.Hand_CardMaster_Draw, false},
		{"/act_card_master_card_list", mainlogic.Hand_CardMaster_CardList, false},
		{"/act_card_master_card2item", mainlogic.Hand_CardMaster_Card2Item, false},
		{"/act_card_master_card2point", mainlogic.Hand_CardMaster_Card2Point, false},
		{"/act_card_master_point2card", mainlogic.Hand_CardMaster_Point2Card, false},

		//月光集市
		{"/act_moonlight_shop_get_info", mainlogic.Hand_MoonlightShop_GetInfo, false},
		{"/act_moonlight_shop_exchangetoken", mainlogic.Hand_MoonlightShop_ExchangeToken, false},
		{"/act_moonlight_shop_reducediscount", mainlogic.Hand_MoonlightShop_ReduceDiscount, false},
		{"/act_moonlight_shop_refreshshop_buy", mainlogic.Hand_MoonlightShop_RefreshShop_Buy, false},
		{"/act_moonlight_shop_refreshshop_auto", mainlogic.Hand_MoonlightShop_RefreshShop_Auto, false},
		{"/act_moonlight_shop_buygoods", mainlogic.Hand_MoonlightShop_BuyGoods, false},
		{"/act_moonlight_shop_getscoreaward", mainlogic.Hand_MoonlightShop_GetScoreAward, false},

		// 沙滩宝贝
		{"/act_beach_baby_info", mainlogic.Hand_BeachBaby_Info, false},
		{"/act_beach_baby_open_goods", mainlogic.Hand_BeachBaby_OpenGoods, false},
		{"/act_beach_baby_open_all_goods", mainlogic.Hand_BeachBaby_OpenAllGoods, false},
		{"/act_beach_baby_refresh_auto", mainlogic.Hand_BeachBaby_Refresh_Auto, false},
		{"/act_beach_baby_refresh_buy", mainlogic.Hand_BeachBaby_Refresh_Buy, false},
		{"/act_beach_baby_get_freeconch", mainlogic.Hand_BeachBaby_GetFreeConch, false},
		{"/act_beach_baby_select_goods", mainlogic.Hand_BeachBaby_SelectGoodsID, false},

		//阵营战
		{"/register_battle_svr", mainlogic.Hand_RegBattleSvr, false}, //注册阵营战服务器
		{"/get_recommandcamp", mainlogic.Hand_RecommandCamp, false},  //获取推荐的阵营战服务器
		{"/set_battlecamp", mainlogic.Hand_SetBattleCamp, false},
		{"/enter_campbattle", mainlogic.Hand_EnterCampBattle, false},
		{"/get_campbat_data", mainlogic.Hand_GetCampBatData, false},
		{"/get_campbat_store_state", mainlogic.Hand_GetCampbatStoreState, false},
		{"/buy_campbat_store_item", mainlogic.Hand_BuyCampbatStoreItem, false},

		//请求界面红点提示
		{"/get_mainui_tip", mainlogic.Hand_GetMainUITip, false},

		//以下全是GM通过后台操作的消息
		//★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★
		{"/update_gamedata", mainlogic.Hand_UpdateGameData, false},
		{"/add_svr_award", mainlogic.Hand_AddSvrAward, false},
		{"/del_svr_award", mainlogic.Hand_DelSvrAward, false},
		{"/send_award_to_player", mainlogic.Hand_SendAwardToPlayer, false},
		{"/get_server_info", mainlogic.Hand_GetServerInfo, false},
		//★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★

		// SDK
		{"/create_recharge_order", mainlogic.Handle_Create_Recharge_Order, false},
		{"/sdk_recharge_success", mainlogic.Handle_Recharge_Success, false},
	}

	max := len(configSlice)
	mainlogic.GMap_HandleMsg_Lock = make(map[string]bool, max)
	for i := 0; i < max; i++ {
		data := &configSlice[i]
		http.HandleFunc(data.url, data.fun)
		mainlogic.GMap_HandleMsg_Lock[data.url] = data.needLock
	}
}

func RegTcpMsgHandler() {
	tcpclient.HandleFunc(msg.MSG_CHECK_IN_ACK, func(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) { return }) //这个消息不用处理
	tcpclient.HandleFunc(msg.MSG_ONLINE_NOTIFY, mainlogic.Hand_OnlineNotify)                                            //玩家上下线通知
	tcpclient.HandleFunc(msg.MSG_DISCONNECT, mainlogic.Hand_DisConnect)
	tcpclient.HandleFunc(msg.MSG_CONNECT, mainlogic.Hand_Connect)

	//以下的消息来自阵营战服务器
	tcpclient.HandleFunc(msg.MSG_LOAD_CAMPBAT_REQ, mainlogic.Hand_LoadCampBatInfo)
	tcpclient.HandleFunc(msg.MSG_KILL_EVENT_REQ, mainlogic.Hand_KillEventReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_QUERY_REQ, mainlogic.Hand_PlayerQueryReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_CHANGE_REQ, mainlogic.Hand_PlayerChangeReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_REVIVE_REQ, mainlogic.Hand_PlayerReviveReq)

	tcpclient.HandleFunc(msg.MSG_START_CARRY_REQ, mainlogic.Hand_StartCarryReq)
	tcpclient.HandleFunc(msg.MSG_FINISH_CARRY_REQ, mainlogic.Hand_FinishCarryReq)

}

func RegConsoleCmdHandler() {
	utility.HandleFunc("setloglevel", mainlogic.HandCmd_SetLogLevel) //例如 setloglevel [1]
}

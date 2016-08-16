//	http.HandleFunc("/get_elite_chapter_detail_info", mainlogic.Hand_GetEliteDetailInfo) //! 获取精英关卡章节详细信息
package main

import (
	"gamesvr/mainlogic"
	"gamesvr/tcpclient"
	"msg"
	"net/http"
	"utility"
)

func RegHttpMsgHandler() {

	//! 任务的消息处理
	http.HandleFunc("/get_tasks", mainlogic.Hand_GetAllTask)                //! 玩家请求全部的任务信息
	http.HandleFunc("/receive_task", mainlogic.Hand_GetTaskAward)           //! 玩家请求领取完成任务奖励
	http.HandleFunc("/receive_taskscore", mainlogic.Hand_GetTaskScoreAward) //! 玩家请求领取任务积分宝箱奖励
	http.HandleFunc("/get_taskscore", mainlogic.Hand_GetTaskScoreInfo)      //! 玩家请求全部任务积分信息

	//玩家登录离开创建角色处理
	http.HandleFunc("/user_login_game", mainlogic.Hand_PlayerLoginGame)   //玩家登录游戏服
	http.HandleFunc("/user_enter_game", mainlogic.Hand_PlayerEnterGame)   //玩家进入游戏服
	http.HandleFunc("/user_leave_game", mainlogic.Hand_PlayerLeaveGame)   //玩家离开游戏服
	http.HandleFunc("/create_new_player", mainlogic.Hand_CreateNewPlayer) //玩家创建角色
	http.HandleFunc("/query_server_time", mainlogic.Hand_QueryServerTime) //玩家查询服务器时间
	http.HandleFunc("/get_login_data", mainlogic.Hand_GetLoginData)       //玩家加载登录信息

	//! 成就的消息处理
	http.HandleFunc("/get_achievement", mainlogic.Hand_GetAllAchievement)       //! 玩家请求全部成就信息
	http.HandleFunc("/receive_achievement", mainlogic.Hand_GetAchievementAward) //! 玩家请求领取成就奖励

	//! 签到的消息处理
	http.HandleFunc("/get_sign", mainlogic.Hand_GetSignInfo) //! 玩家请求获取签到信息
	http.HandleFunc("/daily_sign", mainlogic.Hand_DailySign) //! 玩家请求获取日常签到奖励
	http.HandleFunc("/sign_plus", mainlogic.Hand_SignPlus)   //! 玩家请求获取豪华签到奖励

	//! 神将商店的消息处理
	http.HandleFunc("/get_store", mainlogic.Hand_GetHeroStoreInfo)         //! 玩家请求商店信息
	http.HandleFunc("/refresh_store", mainlogic.Hand_RefreshHeroStore)     //! 玩家请求刷新神将商店
	http.HandleFunc("/store_buy", mainlogic.Hand_HeroStore_Buy)            //! 玩家请求购买神将商店商品
	http.HandleFunc("/get_all_store_data", mainlogic.Hand_GetAllStoreInfo) //! 玩家请求所有商店信息

	//! 商城消息处理
	http.HandleFunc("/get_vip_gift", mainlogic.Hand_GetVipGiftInfo)              //! 玩家请求礼包信息
	http.HandleFunc("/buy_vip_gift", mainlogic.Hand_BuyVipGift)                  //! 玩家请求购买礼包
	http.HandleFunc("/buy_vip_week_gift", mainlogic.Hand_BuyVipWeekGift)         //! 玩家请求购买VIP每周礼包
	http.HandleFunc("/get_goods_buy_times", mainlogic.Hand_GetMallGoodsBuyTimes) //! 玩家请求道具商城商品购买次数
	http.HandleFunc("/get_items_buy_times", mainlogic.Hand_GetMallGoodsBuyTimes) //! 玩家请求道具商城商品购买次数(复制)
	http.HandleFunc("/buy_goods", mainlogic.Hand_BuyMallGoods)                   //! 玩家请求购买道具商城商品
	http.HandleFunc("/get_summon_status", mainlogic.Hand_GetSummonStatus)        //! 玩家请求召唤信息状态
	http.HandleFunc("/get_summon", mainlogic.Hand_GetSummon)                     //! 玩家请求召唤武将
	http.HandleFunc("/exchange_hero", mainlogic.Hand_ExchangeHero)               //! 玩家请求积分兑换英雄

	//! Vip的消息处理
	http.HandleFunc("/query_vip_welfare", mainlogic.Hand_GetDailyVipStatus) //! 玩家请求VIP日常福利领取状态
	http.HandleFunc("/get_vip_welfare", mainlogic.Hand_DailyVipWelfare)     //! 玩家请求领取VIP日常福利

	//! 副本类消息
	http.HandleFunc("/battle_check", mainlogic.Hand_BattleCheck)                       //! 检查挑战副本的条件
	http.HandleFunc("/battle_result", mainlogic.Hand_BattleResult)                     //! 检查挑战副本的条件
	http.HandleFunc("/sweep_copy", mainlogic.Hand_SweepCopy)                           //! 扫荡副本
	http.HandleFunc("/get_main_star_award", mainlogic.Hand_GetMainStarAward)           //! 玩家请求主线关卡星级奖励
	http.HandleFunc("/get_main_scene_award", mainlogic.Hand_GetMainSceneAward)         //! 玩家请求主线关卡场景奖励
	http.HandleFunc("/get_rebel_find_info", mainlogic.Hand_GetRebelFindInfo)           //! 获取发现叛军简单信息
	http.HandleFunc("/get_main_reset_times", mainlogic.Hand_GetMainResetTimes)         //! 玩家查询主线关卡重置次数
	http.HandleFunc("/reset_main_battletimes", mainlogic.Hand_ResetMainBattleTimes)    //! 玩家重置主线关卡挑战次数
	http.HandleFunc("/get_famous_detail", mainlogic.Hand_GetFamousCopyDetailInfo)      //! 获取名将副本详细信息
	http.HandleFunc("/get_famous_award", mainlogic.Hand_GetFamousCopyAward)            //! 获取名将副本章节奖励
	http.HandleFunc("/get_elite_star_award", mainlogic.Hand_GetEliteStarAward)         //! 获取精英关卡星级奖励
	http.HandleFunc("/get_elite_scene_award", mainlogic.Hand_GetEliteSceneAward)       //! 获取精英关卡场景奖励
	http.HandleFunc("/get_elite_reset_times", mainlogic.Hand_GetEliteResetTimes)       //! 获取精英关卡重置次数
	http.HandleFunc("/reset_elite_battletimes", mainlogic.Hand_ResetEliteBattleTimes)  //! 重置精英关卡挑战
	http.HandleFunc("/attack_elite_invade", mainlogic.Hand_AttackInvade)               //! 请求攻击精英关卡入侵
	http.HandleFunc("/get_elite_invade_status", mainlogic.Hand_GetEliteCopyInvadeInfo) //! 请求入侵消息
	http.HandleFunc("/get_copy_data", mainlogic.Hand_GetCopyData)                      //! 获取副本数据

	//! 三国志消息
	http.HandleFunc("/get_sanguozhi_info", mainlogic.Hand_SanGuoZhiInfo)             //! 玩家请求获取三国志信息
	http.HandleFunc("/set_sanguozhi", mainlogic.Hand_SetSanGuoZhi)                   //! 玩家请求三国志命星
	http.HandleFunc("/get_sanguozhi_attr", mainlogic.Hand_GetSanGuoStarAddAttribute) //! 玩家请求获取三国志加属性信息

	//! 竞技场消息处理
	http.HandleFunc("/arena_result", mainlogic.Hand_ChallengeArenaResult)                      //! 玩家反馈挑战竞技场结果
	http.HandleFunc("/get_arena_info", mainlogic.Hand_GetArenaInfo)                            //! 玩家请求竞技场可挑战玩家信息
	http.HandleFunc("/arena_check", mainlogic.Hand_ArenaCheck)                                 //! 玩家请求挑战排位检测
	http.HandleFunc("/arena_store_buy_item", mainlogic.Hand_BuyArenaStoreItem)                 //! 玩家请求购买声望商店物品
	http.HandleFunc("/arena_store_query_award", mainlogic.Hand_GetArenaStoreAleadyBuyAwardLst) //! 玩家请求已购买的声望商店奖励ID列表

	//! 夺宝消息处理
	http.HandleFunc("/get_rob_list", mainlogic.Hand_GetRobList)            //! 玩家请求抢劫名单
	http.HandleFunc("/refresh_rob_list", mainlogic.Hand_GetRobList)        //! 玩家请求刷新抢劫名单
	http.HandleFunc("/rob_treasure", mainlogic.Hand_RobTreasure)           //! 玩家请求夺宝
	http.HandleFunc("/get_free_war_time", mainlogic.Hand_GetFreeWarTime)   //! 玩家请求免战时间
	http.HandleFunc("/get_rob_hero_info", mainlogic.Hand_GetRobHeroInfo)   //! 获取抢劫玩家武将信息
	http.HandleFunc("/treasure_composed", mainlogic.Hand_TreasureComposed) //! 玩家请求合成宝物
	http.HandleFunc("/treasure_melting", mainlogic.Hand_TreasureMelting)   //! 玩家请求宝物熔炼

	//! 三国无双消息处理
	http.HandleFunc("/get_sangokumusou_star", mainlogic.Hand_GetSangokuMusou_StarInfo)                   //! 获取星数信息
	http.HandleFunc("/get_sangokumusou_copy", mainlogic.Hand_GetSangokuMuSou_CopyInfo)                   //! 获取闯关信息
	http.HandleFunc("/get_sangokumusou_elite_copy", mainlogic.Hand_GetSangokuMusou_EliteCopy)            //! 获取精英挑战闯关信息
	http.HandleFunc("/pass_sangokumusou", mainlogic.Hand_PassSangokuMusou_Copy)                          //! 通关三国无双回馈
	http.HandleFunc("/pass_sgws_elite", mainlogic.Hand_PassSangokuMusou_EliteCopy)                       //! 通关三国无双精英挑战
	http.HandleFunc("/sweep_sangoumusou", mainlogic.Hand_SangokuMusou_Sweep)                             //! 请求三星扫荡章节
	http.HandleFunc("/get_sangokumusou_chapter_award", mainlogic.Hand_GetSangokuMusou_ChapterAward)      //! 请求章节奖励
	http.HandleFunc("/get_sangokumusou_attr", mainlogic.Hand_GetSangokuMusou_ChapterAttr)                //! 请求随机三个属性奖励
	http.HandleFunc("/set_sangokumusou_attr", mainlogic.Hand_SetSangokuMusou_ChapterAttr)                //! 选择属性奖励
	http.HandleFunc("/get_sangokumusou_all_attr", mainlogic.Hand_GetSangokuMusou_Attr)                   //! 获取当前所有属性奖励
	http.HandleFunc("/get_sangokumusou_treasure", mainlogic.Hand_GetSangokuMusou_Treasure)               //! 获取无双秘藏
	http.HandleFunc("/buy_sangokumusou_treasure", mainlogic.Hand_BuySangokuMusou_Treasure)               //! 购买无双秘藏
	http.HandleFunc("/reset_sangokumusou_copy", mainlogic.Hand_SangokuMusou_ResetCopy)                   //! 重置挑战
	http.HandleFunc("/get_sangoukumusou_elite_add_times", mainlogic.Hand_GetSangokuMusou_AddEliteCopy)   //! 获取可增加精英挑战次数
	http.HandleFunc("/add_sangoukumusou_elite_copy", mainlogic.Hand_SangokuMusou_AddEliteCopy)           //! 增加精英挑战次数
	http.HandleFunc("/get_sangokumusou_store_aleady_buy", mainlogic.Hand_GetSangokuMusouStore_AleadyBuy) //! 获取无双商店购买次数信息
	http.HandleFunc("/buy_sangokumusou_store", mainlogic.Hand_GetSangokuMusou_StoreItem)                 //! 购买无双商店商品
	http.HandleFunc("/get_sangokumusou_status", mainlogic.Hand_GetSanguowsStatus)                        //! 获取三国无双状态

	//! 领地征讨消息处理
	http.HandleFunc("/get_territory_status", mainlogic.Hand_GetTerritoryStatus)              //! 玩家请求获取自身领地状态
	http.HandleFunc("/challenge_territory", mainlogic.Hand_ChallengeTerritory)               //! 玩家回馈挑战领地成功结果
	http.HandleFunc("/get_friend_territory_status", mainlogic.Hand_GetFriendTerritoryStatus) //! 玩家请求获取好友领地状态
	http.HandleFunc("/get_friend_territory_info", mainlogic.Hand_GetFriendTerritoryDetail)   //! 玩家请求获取好友领地详情
	http.HandleFunc("/help_riot", mainlogic.Hand_SuppressRiot)                               //! 玩家请求帮助好友镇压暴动
	http.HandleFunc("/get_territory_award", mainlogic.Hand_GetTerritoryAward)                //! 玩家请求收获巡逻领地奖励
	http.HandleFunc("/patrol_territory", mainlogic.Hand_PatrolTerritory)                     //! 玩家请求领地放置武将巡逻
	http.HandleFunc("/territory_skill_up", mainlogic.Hand_TerritorySkillLevelUp)             //! 玩家请求提升领地技能等级
	// http.HandleFunc("/query_territory_award", mainlogic.Hand_GetTerritoryAwardLst)           //! 玩家请求查询领地巡逻奖励
	http.HandleFunc("/query_territory_riot", mainlogic.Hand_GetTerritoryRiotInfo) //! 玩家请求查询领地暴动信息

	//! 围剿叛军消息处理
	http.HandleFunc("/get_rebel_info", mainlogic.Hand_GetRebelInfo)                    //! 玩家请求获取叛军
	http.HandleFunc("/attack_rebel", mainlogic.Hand_AttackRebel)                       //! 玩家请求攻击叛军
	http.HandleFunc("/share_rebel", mainlogic.Hand_ShareRebel)                         //! 玩家请求分享叛军
	http.HandleFunc("/get_exploit_award_status", mainlogic.Hand_GetExploitAwardStatus) //! 获取功勋奖励领奖状态
	http.HandleFunc("/get_exploit_award", mainlogic.Hand_GetExploitAward)              //! 玩家请求领取功勋奖励
	http.HandleFunc("/buy_rebel_store", mainlogic.Hand_BuyRebelStore)                  //! 玩家请求购买战功商店物品

	//! 活动消息处理
	http.HandleFunc("/get_seven_activity", mainlogic.Hand_GetSevenActivityInfo)                       //! 玩家请求获取七日活动信息
	http.HandleFunc("/get_seven_activity_award", mainlogic.Hand_GetSevenActivityAward)                //! 玩家获取七日活动奖励
	http.HandleFunc("/get_seven_activity_limit_num", mainlogic.Hand_GetSevenActivityLimitInfo)        //! 玩家获取七日活动限购信息
	http.HandleFunc("/buy_seven_activity_limit", mainlogic.Hand_BuySevenActivityLimit)                //! 玩家购买七日活动限购商品
	http.HandleFunc("/get_activity", mainlogic.Hand_GetActivity)                                      //! 玩家请求当前开启活动
	http.HandleFunc("/get_activity_list", mainlogic.Hand_GetActivity)                                 //! 玩家请求当前开启活动
	http.HandleFunc("/get_open_server_day", mainlogic.Hand_GetServerOpenDay)                          //! 获取服务器开启天数
	http.HandleFunc("/query_activity_login", mainlogic.Hand_QueryActivityLoginInfo)                   //! 查询累计登录活动信息
	http.HandleFunc("/query_activity_firstrecharge", mainlogic.Hand_QueryActivityFirstRechargeInfo)   //! 查询首冲活动信息
	http.HandleFunc("/query_monthcard_days", mainlogic.Hand_QueryActivityMonthCardDays)               //! 查询月卡剩余天数
	http.HandleFunc("/query_activity_action", mainlogic.Hand_QueryActivityActionInfo)                 //! 查询领取体力活动信息
	http.HandleFunc("/query_activity_moneygod", mainlogic.Hand_QueryActivityMoneyGodInfo)             //! 查询迎财神活动信息
	http.HandleFunc("/query_activity_discountsale", mainlogic.Hand_QueryActivityDiscountSaleInfo)     //! 查询折扣贩售活动信息
	http.HandleFunc("/query_activity_totalrecharge", mainlogic.Hand_QueryActivityTotalRechargeInfo)   //! 查询累计充值活动信息
	http.HandleFunc("/query_activity_singlerecharge", mainlogic.Hand_QueryActivitySingleRechargeInfo) //! 查询单笔充值活动信息
	http.HandleFunc("/get_login_award", mainlogic.Hand_GetActivityLoginAward)                         //! 领取累计登录奖励
	http.HandleFunc("/get_first_recharge", mainlogic.Hand_GetFirstRechargeAward)                      //! 领取首冲奖励
	http.HandleFunc("/get_activity_action", mainlogic.Hand_ReceiveActivityAction)                     //! 领取体力奖励
	http.HandleFunc("/welcome_money_god", mainlogic.Hand_WelcomeMoneyGold)                            //! 领取迎财神奖励
	http.HandleFunc("/get_money_god_award", mainlogic.Hand_MoneyGoldAward)                            //! 领取迎财神累积奖励
	http.HandleFunc("/buy_discount_sale", mainlogic.Hand_BuyDiscountSaleItem)                         //! 购买折扣贩售商品
	http.HandleFunc("/get_recharge_award", mainlogic.Hand_GetRechargeAward)                           //! 领取累积充值奖励
	http.HandleFunc("/get_action_retroactive", mainlogic.Hand_ActionRetroactive)                      //! 补签体力活动
	http.HandleFunc("/get_single_award", mainlogic.Hand_GetSingleRechargeAward)                       //! 获取单充奖励
	http.HandleFunc("/query_hunt_treasure", mainlogic.Hand_QueryHuntTreasure)                         //! 查询巡回探宝状态
	http.HandleFunc("/start_hunt", mainlogic.Hand_StartHuntTreasure)                                  //! 开始巡回探宝掷骰
	http.HandleFunc("/query_hunt_award", mainlogic.Hand_QueryHuntTurnsAward)                          //! 查询巡回探宝奖励领取情况
	http.HandleFunc("/get_hunt_award", mainlogic.Hand_GetHuntTurnsAward)                              //! 获取巡回探宝奖励
	http.HandleFunc("/query_hunt_store", mainlogic.Hand_QueryHuntTreasureStore)                       //! 玩家查询巡回商店
	http.HandleFunc("/buy_hunt_store", mainlogic.Hand_BuyHuntTreasureStroreItem)                      //! 玩家购买巡回商店物品
	http.HandleFunc("/query_lucky_wheel", mainlogic.Hand_QueryLuckyWheel)                             //! 查询幸运转盘
	http.HandleFunc("/rotating_wheel", mainlogic.Hand_RotatingWheel)                                  //! 申请转动轮盘
	http.HandleFunc("/get_group_purchase_info", mainlogic.Hand_GetGroupPurchaseInfo)                  //! 获取团购信息
	http.HandleFunc("/buy_group_purchase", mainlogic.Hand_BuyGroupPurchaseItem)                       //! 玩家请求购买团购
	http.HandleFunc("/get_group_purchase_score", mainlogic.Hand_QueryGroupPurchaseScoreAward)         //! 玩家请求查询积分奖励
	http.HandleFunc("/get_group_score_award", mainlogic.Hand_GetGroupPurchaseScoreAward)              //! 玩家请求积分奖励
	http.HandleFunc("/get_festival_task", mainlogic.Hand_GetFestivalTask)                             //! 获取欢庆佳节任务
	http.HandleFunc("/get_festival_exchange", mainlogic.Hand_GetFestivalExchangeInfo)                 //! 获取欢庆佳节兑换信息
	http.HandleFunc("/get_festival_task_award", mainlogic.Hand_GetFestivalTaskAward)                  //! 玩家请求领取欢庆佳节任务奖励
	http.HandleFunc("/exchange_festival_award", mainlogic.Hand_ExchangeFestivalAward)                 //! 玩家兑换奖励
	http.HandleFunc("/clean_hunt_store", mainlogic.Hand_CleanHuntStore)                               //! 清除巡回商店数据
	http.HandleFunc("/get_diff_price", mainlogic.Hand_GetGroupPurchaseCost)                           //! 领取团购差价
	http.HandleFunc("/get_activity_rank", mainlogic.Hand_GetActivityRank)                             //! 获取活动排行榜
	http.HandleFunc("/get_activity_rank_award", mainlogic.Hand_GetActivityRankAward)                  //! 获取活动排行榜奖励
	http.HandleFunc("/get_week_award_status", mainlogic.Hand_GetWeekAwardStatus)                      //! 获取周周盈状态
	http.HandleFunc("/get_week_award", mainlogic.Hand_GetWeekAward)                                   //! 获取周周盈奖励
	http.HandleFunc("/get_level_gift_info", mainlogic.Hand_GetLevelGiftInfo)                          //! 获取等级礼包信息
	http.HandleFunc("/buy_level_gift", mainlogic.Hand_BuyLevelGift)                                   //! 购买等级礼包
	http.HandleFunc("/get_rank_gift_info", mainlogic.Hand_GetRankGiftInfo)                            //! 获取排名礼包信息
	http.HandleFunc("/buy_rank_gift", mainlogic.Hand_BuyRankGift)                                     //! 购买排名礼包
	http.HandleFunc("/get_monthfund_status", mainlogic.Hand_GetMonthFundStatus)                       //! 获取月基金状态
	http.HandleFunc("/receive_month_fund", mainlogic.Hand_ReceiveMonthFund)                           //! 玩家领取月基金奖励

	//! 挖矿消息处理
	http.HandleFunc("/enter_mining", mainlogic.Hand_GetMiningInfo)                             //! 获取挖矿信息
	http.HandleFunc("/mining_get_award", mainlogic.Hand_GetRandBossAward)                      //! 获取打败Boss后九个翻牌奖励
	http.HandleFunc("/select_mining_award", mainlogic.Hand_SelectBossAward)                    //! 选择打败Boss后九个翻牌奖励
	http.HandleFunc("/mining_guaji", mainlogic.Hand_MiningGuaji)                               //! 挖矿挂机
	http.HandleFunc("/mining_guaji_time", mainlogic.Hand_GetMiningGuajiTime)                   //! 查询挖矿挂机倒计时
	http.HandleFunc("/mining_guaji_award", mainlogic.Hand_GetMiningGuajiAward)                 //! 获取挂机收益
	http.HandleFunc("/mining_dig", mainlogic.Hand_MiningDig)                                   //! 玩家请求获取挖地图某个点的信息
	http.HandleFunc("/mining_element_stone", mainlogic.Hand_MiningElement_RefiningStone)       //! 挖矿元素-精炼石
	http.HandleFunc("/mining_event_action_award", mainlogic.Hand_MiningEvent_ActionAward)      //! 挖矿事件-行动力奖励
	http.HandleFunc("/mining_event_black_market", mainlogic.Hand_MiningEvent_BalckMarket)      //! 挖矿事件-黑市
	http.HandleFunc("/mining_buy_black_market", mainlogic.Hand_MiningEvent_BuyBlackMarketItem) //! 挖矿事件-黑市购买
	http.HandleFunc("/mining_event_monster", mainlogic.Hand_MiningEvent_Monster)               //! 挖矿事件-怪物
	http.HandleFunc("/mining_event_treasure", mainlogic.Hand_MiningEvent_Treasure)             //! 挖矿事件-宝箱
	http.HandleFunc("/mining_event_box", mainlogic.Hand_MiningEvent_MagicBox)                  //! 挖矿事件-魔盒
	http.HandleFunc("/mining_event_scan", mainlogic.Hand_MiningEvent_Scan)                     //! 挖矿事件-扫描
	http.HandleFunc("/mining_event_question", mainlogic.Hand_MiningEvent_Question)             //! 挖矿事件-答题
	http.HandleFunc("/mining_event_buff", mainlogic.Hand_MiningEvent_Buff)                     //! 挖矿事件-Buff
	http.HandleFunc("/get_mining_status", mainlogic.Hand_GetMiningStatusCode)                  //! 获取挖矿状态码
	http.HandleFunc("/mining_event_monster_info", mainlogic.Hand_MiningEvent_MonsterInfo)      //! 获取怪物信息

	//! 八卦镜
	http.HandleFunc("/use_baguajing", mainlogic.Hand_UseBaGuaJing) //! 玩家请求使用八卦镜

	//! 开服基金
	http.HandleFunc("/get_fund_status", mainlogic.Hand_GetOpenFundStatus)          //! 请求查询购买基金状态
	http.HandleFunc("/buy_fund", mainlogic.Hand_BuyOpenFund)                       //! 请求购买开服基金
	http.HandleFunc("/get_fund_all_award", mainlogic.Hand_GetOpenFundAllAward)     //! 请求领取基金奖励-全服奖励
	http.HandleFunc("/get_func_level_award", mainlogic.Hand_GetOpenFundLevelAward) //! 请求领取基金奖励-等级返利

	//! 领奖中心
	http.HandleFunc("/query_award_center", mainlogic.Hand_GetAwardCenterInfo)
	http.HandleFunc("/get_award_center", mainlogic.Hand_RecvAwardCenter)

	//! 公会协议
	http.HandleFunc("/create_guild", mainlogic.Hand_CreateGuild)
	http.HandleFunc("/get_guild", mainlogic.Hand_GetGuild)
	http.HandleFunc("/get_guild_lst", mainlogic.Hand_GetMoreGuild)
	http.HandleFunc("/enter_guild", mainlogic.Hand_EnterGuild)
	http.HandleFunc("/get_apply_guild_list", mainlogic.Hand_GetApplyGuildList)
	http.HandleFunc("/get_apply_guild_member_list", mainlogic.Hand_GetApplyGuildMemberList)
	http.HandleFunc("/apply_through", mainlogic.Hand_ApplicationThrough)
	http.HandleFunc("/leave_guild", mainlogic.Hand_ExitGuild)
	http.HandleFunc("/get_sacrifice_status", mainlogic.Hand_GetSacrificeStatus)
	http.HandleFunc("/guild_sacrifice", mainlogic.Hand_GuildSacrifice)
	http.HandleFunc("/get_sacrifice_award", mainlogic.Hand_GetSacrificeAward)
	http.HandleFunc("/query_guild_store", mainlogic.Hand_GetGuildStoreInfo)
	http.HandleFunc("/buy_guild_store", mainlogic.Hand_BuyGuildItem)
	http.HandleFunc("/attack_guild_copy", mainlogic.Hand_AttackGuildCopy)
	http.HandleFunc("/get_guild_copy_status", mainlogic.Hand_GetGuildCopyStatus)
	http.HandleFunc("/get_guild_copy_award", mainlogic.Hand_GetGuildCopyTreasure)
	http.HandleFunc("/query_recv_copy_award", mainlogic.Hand_QueryGuildCopyTreasure)
	http.HandleFunc("/query_guild_copy_rank", mainlogic.Hand_QueryGuildCopyRank)
	http.HandleFunc("/query_guild_msg_board", mainlogic.Hand_QueryGuildMsgBoard)
	http.HandleFunc("/remove_guild_msg_board", mainlogic.Hand_RemoveGuildMsgBoard)
	http.HandleFunc("/write_guild_msg_board", mainlogic.Hand_UseGuildMsgBoard)
	http.HandleFunc("/kick_member", mainlogic.Hand_KickGuildMember)
	http.HandleFunc("/update_guild_name", mainlogic.Hand_UpdateGuildName)
	http.HandleFunc("/update_guild_info", mainlogic.Hand_UpdateGuildInfo)
	http.HandleFunc("/research_guild_skill", mainlogic.Hnad_ResearchGuildSkill)
	http.HandleFunc("/study_guild_skill", mainlogic.Hand_StudyGuildSkill)
	http.HandleFunc("/get_guild_skill", mainlogic.Hand_GetGuildSkillInfo)
	http.HandleFunc("/get_guild_skill_limit", mainlogic.Hand_GetGuildSkillResearchInfo)
	http.HandleFunc("/get_guild_chapter_status", mainlogic.Hand_GetGuildChapterRecvLst)
	http.HandleFunc("/get_guild_chapter_award", mainlogic.Hand_GetGuildChapterAward)
	http.HandleFunc("/get_guild_member_list", mainlogic.Hand_GetGuildMemberList)
	http.HandleFunc("/get_player_info", mainlogic.Hand_GetPlayerInfo)
	http.HandleFunc("/get_guild_log", mainlogic.Hand_GetGuildLog)
	http.HandleFunc("/change_guild_role", mainlogic.Hand_ChangeGuildMemberPose)
	//http.HandleFunc("/guild_levelup", mainlogic.Hand_GuildLevelUp) //! 改变逻辑,暂时不用
	http.HandleFunc("/get_guild_chapter_award_all", mainlogic.Hand_GetAllGuildChapterAward)
	http.HandleFunc("/update_guild_backstatus", mainlogic.Hand_UpdateGuildChapterBackStatus)
	http.HandleFunc("/search_guild", mainlogic.Hand_SearchGuild)
	http.HandleFunc("/cancellation_guild_apply", mainlogic.Hand_CancellationGuildApply)
	http.HandleFunc("/get_guild_status", mainlogic.Hand_GetGuildStatus)

	//! 黑市协议
	http.HandleFunc("/get_black_market_info", mainlogic.Hand_GetBlackMarketInfo)
	http.HandleFunc("/get_black_market_status", mainlogic.Hand_GetBlackMarketStatus)
	http.HandleFunc("/buy_black_market", mainlogic.Hand_BuyBlackMarketGoods)

	//! 名人堂协议
	http.HandleFunc("/send_flower", mainlogic.Hand_SendFlower)
	http.HandleFunc("/get_charm", mainlogic.Hand_GetCharmValue)

	//! 称号协议
	http.HandleFunc("/get_title", mainlogic.Hand_GetTitle)
	http.HandleFunc("/activate_title", mainlogic.Hand_ActivateTitle)
	http.HandleFunc("/equi_title", mainlogic.Hand_EquipTitle)

	//! 夺粮战
	http.HandleFunc("/get_foodwar_challenger", mainlogic.Hand_FoodWar_GetChallenger)
	http.HandleFunc("/get_foodwar_time", mainlogic.Hand_FoodWar_GetTime)
	http.HandleFunc("/get_foodwar_revenge_status", mainlogic.Hand_FoodWar_RevengeStatus)
	http.HandleFunc("/rob_food", mainlogic.Hand_RobFood)
	http.HandleFunc("/get_food_rank", mainlogic.Hand_FoodWar_GetRank)
	http.HandleFunc("/buy_food_times", mainlogic.Hand_FoodWar_BuyTimes)
	http.HandleFunc("/recv_food_award", mainlogic.Hand_FoodWar_RecvAward)
	http.HandleFunc("/query_food_award", mainlogic.Hand_FoodWar_QueryAward)
	http.HandleFunc("/get_foodwar_status", mainlogic.Hand_FoodWar_GetStatus)
	http.HandleFunc("/revenge_rob", mainlogic.Hand_FoodWar_Revenge)

	//! 英魂
	http.HandleFunc("/activate_herosouls", mainlogic.Hand_ActivateHeroSouls)
	http.HandleFunc("/query_herosouls_chapter", mainlogic.Hand_QueryChapterHeroSoulsDetail)
	http.HandleFunc("/get_herosouls_lst", mainlogic.Hand_GetHeroSoulsLst)
	http.HandleFunc("/refresh_herosouls", mainlogic.Hand_RefreshHeroSoulsLst)
	http.HandleFunc("/challenge_herosouls", mainlogic.Hand_ChallengeHeroSouls)
	http.HandleFunc("/buy_challenge_herosouls", mainlogic.Hand_BuyChallengeHeroSoulsTimes)
	http.HandleFunc("/reset_herosouls_lst", mainlogic.Hand_ResetHeroSoulsLst)
	http.HandleFunc("/query_herosouls_rank", mainlogic.Hand_QueryHeroSoulsRank)
	http.HandleFunc("/query_herosouls_store", mainlogic.Hand_QueryHeroSoulsStoreInfo)
	http.HandleFunc("/buy_herosouls", mainlogic.Hand_BuyHeroSoulsStoreItem)
	http.HandleFunc("/query_herosouls_achievement", mainlogic.Hand_QuerySoulMapInfo)
	http.HandleFunc("/activate_herosouls_achievement", mainlogic.Hand_ActivateheroSoulsAchievement)
	http.HandleFunc("/query_herosouls_property", mainlogic.Hand_QueryHeroSoulsPerproty)

	//英雄消息处理
	http.HandleFunc("/get_battle_data", mainlogic.Hand_GetBattleData)     //玩家请求上阵数据
	http.HandleFunc("/upgrade_hero", mainlogic.Hand_UpgradeHero)          //升级英雄(非主角)
	http.HandleFunc("/change_hero", mainlogic.Hand_ChangeHero)            //玩家更换英雄
	http.HandleFunc("/change_back_hero", mainlogic.Hand_ChangeBackHero)   //玩家更换援军英雄
	http.HandleFunc("/breakout_hero", mainlogic.Hand_BreakOutHero)        //玩家突破英雄
	http.HandleFunc("/culture_hero", mainlogic.Hand_CultureHero)          //玩家培养英雄
	http.HandleFunc("/compose_hero", mainlogic.Hand_ComposeHero)          //玩家天命英雄
	http.HandleFunc("/upgod_hero", mainlogic.Hand_UpgodHero)              //玩家化神英雄
	http.HandleFunc("/change_career", mainlogic.Hand_Change_Career)       //玩家更改职业
	http.HandleFunc("/set_wake_item", mainlogic.Hand_SetWakeItem)         //玩家设置觉醒道具
	http.HandleFunc("/up_wake_level", mainlogic.Hand_UpWakeLevel)         //玩家提升觉醒等级
	http.HandleFunc("/compose_wake_item", mainlogic.Hand_ComposeWakeItem) //玩家合成觉醒道具等级
	http.HandleFunc("/query_destiny", mainlogic.Hand_QueryHeroDestiny)    //玩家查询天命状态
	http.HandleFunc("/destiny_hero", mainlogic.Hand_DestinyHero)          //玩家天命英雄
	http.HandleFunc("/upgrade_diaowen", mainlogic.Hand_UpgradeDiaoWen)    //玩家升品雕文
	http.HandleFunc("/xilian_diaowen", mainlogic.Hand_XiLianDiaoWen)      //玩家洗炼雕文
	http.HandleFunc("/xilian_tihuan", mainlogic.Hand_XiLianTiHuan)        //玩家洗炼替换雕文
	http.HandleFunc("/upgrade_pet", mainlogic.Hand_UpgradePet)            //升级宠物
	http.HandleFunc("/upstar_pet", mainlogic.Hand_UpstarPet)              //升星宠物
	http.HandleFunc("/upgod_pet", mainlogic.Hand_UpgodPet)                //神炼宠物
	http.HandleFunc("/change_pet", mainlogic.Hand_ChangePet)              //更换宠物
	http.HandleFunc("/unset_pet", mainlogic.Hand_UnsetPet)                //下阵宠物
	http.HandleFunc("/compose_pet", mainlogic.Hand_ComposePet)            //装备合成

	//装备
	http.HandleFunc("/change_equip", mainlogic.Hand_ChangeEquip)         //装备更换
	http.HandleFunc("/equip_strengthen", mainlogic.Hand_EquipStrengthen) //装备强化
	http.HandleFunc("/equip_refine", mainlogic.Hand_EquipRefine)         //装备精炼
	http.HandleFunc("/equip_risestar", mainlogic.Hand_EquipRiseStar)     //装备升星
	http.HandleFunc("/compose_equip", mainlogic.Hand_ComposeEquip)       //装备合成

	//宝物
	http.HandleFunc("/change_gem", mainlogic.Hand_ChangeGem)         //宝物更换
	http.HandleFunc("/gem_strengthen", mainlogic.Hand_GemStrengthen) //宝物强化
	http.HandleFunc("/gem_refine", mainlogic.Hand_GemRefine)         //宝物精炼

	//时装
	http.HandleFunc("/fashion_set", mainlogic.Hand_FashionSet)           //时装装备
	http.HandleFunc("/fashion_strength", mainlogic.Hand_FashionStrength) //时装强化
	http.HandleFunc("/fashion_recast", mainlogic.Hand_FashionRecast)     //时装重铸
	http.HandleFunc("/fashion_compose", mainlogic.Hand_FashionCompose)   //时装合成
	http.HandleFunc("/fashion_melting", mainlogic.Hand_FashionMelting)   //时装熔炼

	//背包
	http.HandleFunc("/get_bag_data", mainlogic.Hand_GetBagData)              //请求背包数据
	http.HandleFunc("/get_bag_heros", mainlogic.Hand_GetBagHeros)            //请求背包中的所有英雄
	http.HandleFunc("/get_bag_equips", mainlogic.Hand_GetBagEquips)          //请求背包中的所有装备
	http.HandleFunc("/get_bag_hero_piece", mainlogic.Hand_GetBagHerosPiece)  //请求背包中的所有英雄碎片
	http.HandleFunc("/get_bag_equip_piece", mainlogic.Hand_GetBagEquipPiece) //请求背包中的所有装备碎片
	http.HandleFunc("/get_bag_gem_piece", mainlogic.Hand_GetBagGemPiece)     //请求背包中的所有宝物碎片
	http.HandleFunc("/get_bag_gems", mainlogic.Hand_GetBagGems)              //请求背包中的所有的宝物
	http.HandleFunc("/get_bag_items", mainlogic.Hand_GetBagItems)            //请求背包里的道具
	http.HandleFunc("/get_bag_wake_items", mainlogic.Hand_GetBagWakeItems)   //请求背包里的觉醒道具
	http.HandleFunc("/get_bag_pets", mainlogic.Hand_GetBagPets)              //请求背包中的所有宠物
	http.HandleFunc("/get_bag_pet_piece", mainlogic.Hand_GetBagPetsPiece)    //请求背包中的所有宠物碎片
	http.HandleFunc("/use_item", mainlogic.Hand_UseItem)                     //使用背包里的道具
	http.HandleFunc("/sell_item", mainlogic.Hand_SellItem)                   //使用背包里的道具

	//角色信息
	http.HandleFunc("/get_role_data", mainlogic.Hand_GetRoleData)            //玩家获取角色数据
	http.HandleFunc("/levelup_notify", mainlogic.Hand_LevelUpNotify)         //玩家等级升级通知
	http.HandleFunc("/change_role_name", mainlogic.Hand_ChangeRoleName)      //更改角色名字
	http.HandleFunc("/get_new_wizard", mainlogic.Hand_GetNewWizard)          //读取新手向导
	http.HandleFunc("/set_new_wizard", mainlogic.Hand_SetNewWizard)          //设置新手向导
	http.HandleFunc("/get_collection_heros", mainlogic.Hand_GetCollectHeros) //获取玩家收集过的英雄

	//回收
	http.HandleFunc("/query_hero_decompose_cost", mainlogic.Hand_QueryHeroDecomposeCost) //! 查询分解英雄材料
	http.HandleFunc("/decompose_hero", mainlogic.Hand_DecomposeHero)                     //! 分解英雄
	http.HandleFunc("/query_hero_relive_cost", mainlogic.Hand_QueryHeroRelive)           //! 查询重生英雄材料
	http.HandleFunc("/relive_hero", mainlogic.Hand_ReliveHero)                           //! 重生英雄

	http.HandleFunc("/query_equip_decompose_cost", mainlogic.Hand_QueryEquipDecomposeCost) //! 查询分解装备材料
	http.HandleFunc("/decompose_equip", mainlogic.Hand_DecomposeEquip)                     //! 分解装备
	http.HandleFunc("/query_equip_relive_cost", mainlogic.Hand_QueryEquipRelive)           //! 查询重生装备材料
	http.HandleFunc("/relive_equip", mainlogic.Hand_ReliveEquip)                           //! 重生装备

	http.HandleFunc("/query_pet_decompose_cost", mainlogic.Hand_QueryDecomposePetCost) //! 查询分解战宠材料
	http.HandleFunc("/decompose_pet", mainlogic.Hand_DecomposePet)                     //! 分解战宠
	http.HandleFunc("/query_pet_relive_cost", mainlogic.Hand_QueryPetRelive)           //! 查询重生战宠材料
	http.HandleFunc("/relive_pet", mainlogic.Hand_RelivePet)                           //! 重生战宠

	http.HandleFunc("/query_gem_relive_cost", mainlogic.Hand_QueryGemRelive) //! 查询宝物重生材料
	http.HandleFunc("/relive_gem", mainlogic.Hand_ReliveGem)                 //! 重生宝物

	//挂机
	http.HandleFunc("/hangup_get_info", mainlogic.Hand_GetHangUpInfo) //请求挂机信息
	http.HandleFunc("/hangup_set_boss", mainlogic.Hand_SetBoss)       //设置挂机BOSS
	http.HandleFunc("/hangup_quick_fight", mainlogic.Hand_QuickFight) //快速战斗请求
	http.HandleFunc("/hangup_add_grid", mainlogic.Hand_AddGrid)       //快速战斗请求
	http.HandleFunc("/hangup_use_exp", mainlogic.Hand_UseExpItem)     //一键使用经验丹

	//以下为测试消息
	http.HandleFunc("/test_get_money", mainlogic.Hand_TestGetMoney)
	http.HandleFunc("/test_get_action", mainlogic.Hand_TestGetAction)
	http.HandleFunc("/test_uplevel", mainlogic.Hand_TestUplevel)
	http.HandleFunc("/test_uplevel_ten", mainlogic.Hand_TestUplevelTen)
	http.HandleFunc("/test_get_bag_heros", mainlogic.Hand_GetBagHeros)            //请求背包中的所有英雄
	http.HandleFunc("/test_get_bag_equips", mainlogic.Hand_GetBagEquips)          //请求背包中的所有装备
	http.HandleFunc("/test_get_bag_hero_piece", mainlogic.Hand_GetBagHerosPiece)  //请求背包中的所有英雄碎片
	http.HandleFunc("/test_get_bag_equip_piece", mainlogic.Hand_GetBagEquipPiece) //请求背包中的所有装备碎片
	http.HandleFunc("/test_get_bag_gem_piece", mainlogic.Hand_GetBagGemPiece)     //请求背包中的所有宝物碎片
	http.HandleFunc("/test_get_bag_gems", mainlogic.Hand_GetBagGems)              //请求背包中的所有的宝物
	http.HandleFunc("/test_get_bag_items", mainlogic.Hand_GetBagItems)            //请求背包里的道具
	http.HandleFunc("/test_get_bag_wake_items", mainlogic.Hand_GetBagWakeItems)   //请求背包里的觉醒道具
	http.HandleFunc("/test_add_vip", mainlogic.Hand_TestAddVip)                   //请求增加Vip
	http.HandleFunc("/test_add_guild", mainlogic.Hand_TestAddGuildExp)            //请求增加公会经验
	http.HandleFunc("/test_compress", mainlogic.Hand_TestCompress)                //测试压缩协议
	http.HandleFunc("/test_heros_property", mainlogic.Hand_GetHerosProperty)      //测试获取玩家各项属性
	http.HandleFunc("/test_add_item", mainlogic.Hand_TestAddItem)
	http.HandleFunc("/test_charge_money", mainlogic.Hand_TestAddCharge) //! 测试活动充值相关

	//充值消息处理
	http.HandleFunc("/get_charge_info", mainlogic.Hand_GetChargeInfo)       //! 玩家请求充值结果
	http.HandleFunc("/get_charge_result", mainlogic.Hand_GetChargeResult)   //! 玩家请求充值结果
	http.HandleFunc("/receive_month_card", mainlogic.Hand_ReceiveMonthCard) //! 玩家请求领取月卡

	//邮件系统
	http.HandleFunc("/receive_all_mails", mainlogic.Hand_ReceiveAllMails) //! 玩家请求所有的邮件

	//排行榜
	http.HandleFunc("/get_level_rank", mainlogic.Hand_GetLevelRank)            //! 玩家请求等级排行榜
	http.HandleFunc("/get_fight_rank", mainlogic.Hand_GetFightRank)            //! 玩家请求战力排行榜
	http.HandleFunc("/get_sanguows_rank", mainlogic.Hand_GetSanguowsRank)      //! 玩家请求无双排行榜
	http.HandleFunc("/get_arena_rank", mainlogic.Hand_GetArenaRank)            //! 玩家请求竞技场排行榜
	http.HandleFunc("/get_rebel_rank", mainlogic.Hand_GetRebelRank)            //! 玩家请求叛军排行榜
	http.HandleFunc("/get_guild_level_rank", mainlogic.Hand_GetGuildLevelRank) //! 玩家请求公会等级排行榜
	http.HandleFunc("/get_guild_copy_rank", mainlogic.Hand_GetGuildCopyRank)   //! 玩家请求公会副本排行榜
	http.HandleFunc("/get_score_rank", mainlogic.Hand_GetScoreRank)            //! 玩家请求积分赛排行榜
	http.HandleFunc("/get_wander_rank", mainlogic.Hand_GetWanderRank)          //! 玩家请求云游戏排行榜
	http.HandleFunc("/get_campbat_rank", mainlogic.Hand_GetCampBatRank)        //! 玩家请求阵营战排行榜

	//积分赛
	http.HandleFunc("/get_score_target", mainlogic.Hand_GetScoreTarget)
	http.HandleFunc("/get_score_battle_check", mainlogic.Hand_GetScoreBattleCheck)
	http.HandleFunc("/set_score_battle_result", mainlogic.Hand_SetScoreBattleResult)
	http.HandleFunc("/get_score_time_award", mainlogic.Hand_GetScoreTimeAward)
	http.HandleFunc("/rcv_score_time_award", mainlogic.Hand_RcvScoreTimeAward)
	http.HandleFunc("/buy_score_fight_time", mainlogic.Hand_BuyScoreFightTime)
	http.HandleFunc("/get_score_store_state", mainlogic.Hand_GetScoreStoreState)
	http.HandleFunc("/buy_score_store_item", mainlogic.Hand_BuyScoreStoreItem)

	//跨服服来请求玩家数据
	http.HandleFunc("/select_target_player", mainlogic.Hand_SelectTargetPlayer)
	http.HandleFunc("/get_fight_target", mainlogic.Hand_GetFightTarget)

	//社交
	http.HandleFunc("/get_all_friend", mainlogic.Hand_GetAllFriend)
	http.HandleFunc("/get_online_friend", mainlogic.Hand_GetOnlineFriend)
	http.HandleFunc("/add_friend_request", mainlogic.Hand_AddFriendReq)
	http.HandleFunc("/del_friend_request", mainlogic.Hand_DelFriendReq)
	http.HandleFunc("/process_friend_request", mainlogic.Hand_ProcessFriendReq)
	http.HandleFunc("/get_apply_list", mainlogic.Hand_GetApplyList)
	http.HandleFunc("/search_friend", mainlogic.Hand_SearchFriend)
	http.HandleFunc("/recomand_friend", mainlogic.Hand_RecomandFriend)
	http.HandleFunc("/give_action", mainlogic.Hand_GiveAction)
	http.HandleFunc("/receive_action", mainlogic.Hand_ReceiveAction)

	//云游
	http.HandleFunc("/wander_getinfo", mainlogic.Hand_WanderGetInfo)
	http.HandleFunc("/wander_reset", mainlogic.Hand_WanderReset)
	http.HandleFunc("/wander_sweep", mainlogic.Hand_WanderSweep)
	http.HandleFunc("/wander_openbox", mainlogic.Hand_WanderOpenBox)
	http.HandleFunc("/wander_check", mainlogic.Hand_WanderCheck)
	http.HandleFunc("/wander_result", mainlogic.Hand_WanderResult)

	//卡牌大师
	http.HandleFunc("/act_card_master_draw", mainlogic.Hand_CardMaster_Draw)
	http.HandleFunc("/act_card_master_card_list", mainlogic.Hand_CardMaster_CardList)
	http.HandleFunc("/act_card_master_card2item", mainlogic.Hand_CardMaster_Card2Item)
	http.HandleFunc("/act_card_master_card2point", mainlogic.Hand_CardMaster_Card2Point)
	http.HandleFunc("/act_card_master_point2card", mainlogic.Hand_CardMaster_Point2Card)

	//月光集市
	http.HandleFunc("/act_moonlight_shop_get_info", mainlogic.Hand_MoonlightShop_GetInfo)
	http.HandleFunc("/act_moonlight_shop_exchangetoken", mainlogic.Hand_MoonlightShop_ExchangeToken)
	http.HandleFunc("/act_moonlight_shop_reducediscount", mainlogic.Hand_MoonlightShop_ReduceDiscount)
	http.HandleFunc("/act_moonlight_shop_refreshshop_buy", mainlogic.Hand_MoonlightShop_RefreshShop_Buy)
	http.HandleFunc("/act_moonlight_shop_refreshshop_auto", mainlogic.Hand_MoonlightShop_RefreshShop_Auto)
	http.HandleFunc("/act_moonlight_shop_buygoods", mainlogic.Hand_MoonlightShop_BuyGoods)
	http.HandleFunc("/act_moonlight_shop_getscoreaward", mainlogic.Hand_MoonlightShop_GetScoreAward)

	// 沙滩宝贝
	http.HandleFunc("/act_beach_baby_info", mainlogic.Hand_BeachBaby_Info)
	http.HandleFunc("/act_beach_baby_open_goods", mainlogic.Hand_BeachBaby_OpenGoods)
	http.HandleFunc("/act_beach_baby_open_all_goods", mainlogic.Hand_BeachBaby_OpenAllGoods)
	http.HandleFunc("/act_beach_baby_refresh_auto", mainlogic.Hand_BeachBaby_Refresh_Auto)
	http.HandleFunc("/act_beach_baby_refresh_buy", mainlogic.Hand_BeachBaby_Refresh_Buy)
	http.HandleFunc("/act_beach_baby_get_freeconch", mainlogic.Hand_BeachBaby_GetFreeConch)
	http.HandleFunc("/act_beach_baby_select_goods", mainlogic.Hand_BeachBaby_SelectGoodsID)

	//阵营战
	http.HandleFunc("/register_battle_svr", mainlogic.Hand_RegBattleSvr) //注册阵营战服务器
	http.HandleFunc("/get_recommandcamp", mainlogic.Hand_RecommandCamp)  //获取推荐的阵营战服务器
	http.HandleFunc("/set_battlecamp", mainlogic.Hand_SetBattleCamp)
	http.HandleFunc("/enter_campbattle", mainlogic.Hand_EnterCampBattle)
	http.HandleFunc("/get_campbat_data", mainlogic.Hand_GetCampBatData)
	http.HandleFunc("/get_campbat_store_state", mainlogic.Hand_GetCampbatStoreState)
	http.HandleFunc("/buy_campbat_store_item", mainlogic.Hand_BuyCampbatStoreItem)

	//请求界面红点提示
	http.HandleFunc("/get_mainui_tip", mainlogic.Hand_GetMainUITip)

	//以下全是GM通过后台操作的消息
	//★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★
	http.HandleFunc("/update_gamedata", mainlogic.Hand_UpdateGameData)
	http.HandleFunc("/add_svr_award", mainlogic.Hand_AddSvrAward)
	http.HandleFunc("/del_svr_award", mainlogic.Hand_DelSvrAward)
	http.HandleFunc("/send_award_to_player", mainlogic.Hand_SendAwardToPlayer)
	http.HandleFunc("/server_state_info", mainlogic.Hand_ServerStateInfo)
	//★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★

}

func RegTcpMsgHandler() {
	tcpclient.HandleFunc(msg.MSG_CHECK_IN_ACK, func(pTcpConn *tcpclient.TCPConn, pdata []byte) { return }) //这个消息不用处理
	tcpclient.HandleFunc(msg.MSG_ONLINE_NOTIFY, mainlogic.Hand_OnlineNotify)                               //玩家上下线通知
	tcpclient.HandleFunc(msg.MSG_DISCONNECT, mainlogic.Hand_DisConnect)
	tcpclient.HandleFunc(msg.MSG_CONNECT, mainlogic.Hand_Connect)

	//以下的消息来自阵营战服务器
	tcpclient.HandleFunc(msg.MSG_LOAD_CAMPBAT_REQ, mainlogic.Hand_LoadCampBatInfo)
	tcpclient.HandleFunc(msg.MSG_KILL_EVENT_REQ, mainlogic.Hand_KillEventReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_QUERY_REQ, mainlogic.Hand_PlayerQueryReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_CHANGE_REQ, mainlogic.Hand_PlayerChangeReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_CARRY_REQ, mainlogic.Hand_PlayerCarryReq)
	tcpclient.HandleFunc(msg.MSG_PLAYER_REVIVE_REQ, mainlogic.Hand_PlayerReviveReq)

}

func RegConsoleCmdHandler() {
	utility.HandleFunc("setloglevel", mainlogic.HandCmd_SetLogLevel) //例如 setloglevel [1]
}

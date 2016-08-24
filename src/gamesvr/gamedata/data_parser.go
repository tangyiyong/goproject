package gamedata

//InitDataParser 初始化CSV文件解析器
var NullParser = TDataParser{nil, nil, nil}
var G_DataParserMap = map[string]TDataParser{
	//通用选项配置表
	"type_option": TDataParser{InitOptionParser, ParseOptionRecord, nil},

	//副本配置表
	"type_copy_base":   TDataParser{InitCopyParser, ParseCopyRecord, nil},
	"type_copy_main":   TDataParser{InitMainParser, ParseMainRecord, nil},
	"type_copy_elite":  TDataParser{InitEliteParser, ParseEliteRecord, nil},
	"type_copy_daily":  TDataParser{InitDailyParse, ParseDailyRecord, nil},
	"type_copy_famous": TDataParser{InitFamousParser, ParseFamousRecord, nil},

	//道具物品配置表
	"type_item": TDataParser{InitItemParser, ParseItemRecord, nil},

	//商城配置表
	"type_mall": TDataParser{InitMallParser, ParseMallRecord, nil},

	//VIP配置表
	"type_vip":      TDataParser{InitVipParser, ParseVipRecord, nil},
	"type_vip_week": TDataParser{InitVipWeekParser, ParseVipWeekRecord, nil},

	//任务成就配置表
	"type_task":                 TDataParser{InitTaskParser, ParseTaskRecord, nil},
	"type_achievement":          TDataParser{InitAchievementParser, ParseAchievementRecord, nil},
	"type_seven_activity":       TDataParser{InitSevenActivityParser, ParseSevenActivityRecord, nil},
	"type_seven_activity_store": TDataParser{InitSevenActivityStoreRecord, ParseSevenActivityStoreRecord, nil},
	"type_task_award":           TDataParser{InitTaskAwardParser, ParseTaskAwardRecord, nil},
	"type_tasktype":             TDataParser{InitTaskTypeParser, ParseTaskTypeRecord, nil},

	//奖励配置表
	"type_award": TDataParser{InitAwardParser, ParseAwardRecord, nil},

	//签到配置表
	"type_sign":      TDataParser{InitSignParser, ParseSignRecord, nil},
	"type_sign_plus": TDataParser{InitSignPlusParser, ParseSignPlusRecord, nil},

	//神将商店
	"type_store": TDataParser{InitStoreParse, ParseStoreRecord, nil},

	//三国志配置表
	"type_sanguozhi": TDataParser{InitSanGuoZhiParser, ParseSanGuoZhiRecord, nil},

	//招贤配置表
	"type_summon_config": TDataParser{InitSummonConfigParser, ParseSummonConfigRecord, nil},

	//竞技场配置表
	"type_arena":       TDataParser{InitArenaParser, ParseArenaRecord, nil},
	"type_arena_rank":  TDataParser{InitArenaRankParser, ParseArenaRankRecord, nil},
	"type_arena_store": TDataParser{InitArenaStoreParser, ParseArenaStoreRecord, nil},
	"type_arena_money": TDataParser{InitArenaMoneyParser, ParseArenaMoneyRecord, nil},

	//夺宝配置表
	"type_rob":              TDataParser{InitRobParser, ParseRobRecord, nil},
	"type_treasure_melting": TDataParser{InitTreasureMeltingParser, ParseTreasureMeltingRecord, nil},

	//三国无双配置表
	"type_sgws_chapter":       TDataParser{InitSangokuMusouChapter, ParseSangokuMusouChapterRecord, nil},
	"type_sgws_chapter_attr":  TDataParser{InitSangokuMusouAttrMarkupParser, ParseSangokuMusouAttrMarkupRecord, nil},
	"type_sgws_chapter_award": TDataParser{InitSangokuMusouChapterAwardParser, ParseSangokuMusouChapterAwardRecord, nil},
	"type_sgws_elite_copy":    TDataParser{InitSangokuMusouEliteCopyParser, ParseSangokuMusouEliteCopyRecord, nil},
	"type_sgws_store":         TDataParser{InitSangokuMusouStoreParser, ParseSangokuMusouStoreRecord, nil},
	"type_sgws_treasure":      TDataParser{InitSangokuMusouSaleParser, ParseSangokuMusouSaleRecord, nil},

	//领地攻伐表
	"type_territory":          TDataParser{InitTerritoryParser, ParseTerritoryRecord, nil},
	"type_territoryaward":     TDataParser{InitTerritoryAwardParser, ParseTerritoryAwardRecord, nil},
	"type_territoryskill":     TDataParser{InitTerritorySkillParser, ParseTerritorySkillRecord, nil},
	"type_territorypatrol":    TDataParser{InitTerritoryPatrolParser, ParseTerritoryPatrolRecord, nil},
	"type_territoryawardtype": TDataParser{InitTerritoryAwardTypeParser, ParseTerritoryAwardTypeRecord, nil},

	//叛军围剿表
	"type_rebel_siege":      TDataParser{InitRebelSiegeParser, ParseRebelSiegeRecord, nil},
	"type_rebel_action":     TDataParser{InitRebelActionAwardParser, ParseRebelActionAwardRecord, nil},
	"type_rebel_award":      TDataParser{InitExploitAwardParser, ParseExploitAwardRecord, nil},
	"type_rebel_store":      TDataParser{InitExploitStoreParser, ParseExploitStoreRecord, nil},
	"type_rebel_activity":   TDataParser{InitRebelActivityParser, ParseRebelActivityRecord, nil},
	"type_rebel_rank_award": TDataParser{InitRebelRankAwardParser, ParseRebelRankAwardRecord, nil},

	//挖矿表
	"type_mining_stone_num":    TDataParser{InitMiningStoneRandomParser, ParseMiningStoneRecord, nil},
	"type_mining_element":      TDataParser{InitMiningElementParser, ParseMiningElementRecord, nil},
	"type_mining_event":        TDataParser{InitMiningEventParser, ParserMiningEventRecord, nil},
	"type_mining_black_market": TDataParser{InitMiningEventBlackMarketParser, ParserMiningEventBlackMarketRecord, nil},
	"type_mining_question":     NullParser,
	"type_mining_buff":         TDataParser{InitMiningEventBuffParser, ParserMiningEventBuffRecord, nil},
	"type_mining_treasure":     TDataParser{InitMiningEventTreasureParser, ParserMiningEventTreasureRecord, nil},
	"type_mining_monster":      TDataParser{InitMiningEventMonsterPerser, ParseMiningEventMonsterRecord, nil},
	"type_mining_award":        TDataParser{InitMiningAwardParser, ParserMiningAwardRecord, nil},
	"type_mining_guaji":        TDataParser{InitMiningGuaJiParser, ParserMiningGuaJiRecord, nil},
	"type_mining_rand":         TDataParser{InitMiningRandParser, ParserMiningRandRecord, nil},

	//活动配置表
	"type_activity":                   TDataParser{InitActivityParser, ParseActivityRecord, nil},
	"type_activity_competition":       TDataParser{InitCompetitionParser, ParseCompetitionRecord, nil},
	"type_activity_action":            TDataParser{InitRecvActionParser, ParseRecvActionRecord, nil},
	"type_activity_discount":          TDataParser{InitDiscountSaleParser, ParseDiscountSaleRecord, nil},
	"type_activity_login":             TDataParser{InitActivityLoginParser, ParseActivityLoginRecord, nil},
	"type_activity_moneygod":          TDataParser{InitActivityMoneyParser, ParseActivityMoneyRecord, nil},
	"type_activity_recharge":          TDataParser{InitActivityRechargeParser, ParseActivityRechargeRecord, nil},
	"type_activity_limitdaily":        TDataParser{InitActivityLimitDailyParser, ParseActivityLimitDailyRecord, nil},
	"type_activity_hunt_map":          TDataParser{InitHuntTreasureMapParser, ParseHuntTreasureMapRecord, nil},
	"type_activity_hunt_store":        TDataParser{InitHuntTreasureStoreParser, ParseHuntTreasureStoreRecord, nil},
	"type_activity_hunt_turn":         TDataParser{InitHuntTreasureAwardParser, ParseHuntTreasureAwardRecord, nil},
	"type_activity_rank":              TDataParser{InitOperationalRankAwardParser, ParseOperationalRankAwardRecord, nil},
	"type_activity_lucky_wheel":       TDataParser{InitLuckyWheelParser, ParseLuckyWheelRecord, nil},
	"type_activity_group_purchase":    TDataParser{InitGroupPurchaseParser, ParseGroupPurchaseRecord, nil},
	"type_activity_group_score":       TDataParser{InitGroupPurchaseScoreParser, ParseGroupPurchaseScoreRecord, nil},
	"type_activity_festival_task":     TDataParser{InitFestivalTaskParser, ParseFestivalTaskRecord, nil},
	"type_activity_festival_exchange": TDataParser{InitFestivalExchangeParser, ParseFestivalExchangeRecord, nil},
	"type_activity_week_award":        TDataParser{InitActivityWeekAwardParser, ParseActivityWeekAwardRecord, nil},
	"type_activity_level_gift":        TDataParser{InitActivityLevelGiftParser, ParseActivityLevelGiftRecord, nil},
	"type_activity_month_fund":        TDataParser{InitActivityMonthFundParser, ParseActivityMonthFundRecord, nil},
	"type_activity_limitsale":         TDataParser{InitLimitSaleItemParser, ParseLimitSaleItemRecord, nil},
	"type_activity_limitsale_award":   TDataParser{InitLimitSaleAllAwardParser, ParseLimitSaleAllAwardRecord, nil},

	//公会配置表
	"type_guild_base":            TDataParser{InitGuildParser, ParseGuildBaseRecord, nil},
	"type_guild_copy":            TDataParser{InitGuildCopyParser, ParseGuildCopyRecord, nil},
	"type_guild_copy_award":      TDataParser{InitGuildCopyAwardParser, ParseGuildCopyAwardRecord, nil},
	"type_guild_role":            TDataParser{InitGuildRoleParser, ParseGuildRoleRecord, nil},
	"type_guild_sacrifice":       TDataParser{InitGuildSacrificeParser, ParseGuildSacrificeRecord, nil},
	"type_guild_sacrifice_award": TDataParser{InitGuildSacrificeAwardParser, ParseGuildSacrificeAwardRecord, nil},
	"type_guild_store":           TDataParser{InitGuildStoreParser, ParseGuildStoreRecord, nil},
	"type_guild_skill":           TDataParser{InitGuildSkillParser, ParseGuildSkillRecord, nil},
	"type_guild_skill_level":     TDataParser{InitGuildSkillLimitParser, ParseGuildSkillLimitRecord, nil},

	//装备配置表
	"type_equipment":           TDataParser{InitEquipParser, ParseEquipRecord, nil},
	"type_equip_star":          TDataParser{InitEquipStarParser, ParseEquipStarRecord, nil},
	"type_equip_refine_cost":   TDataParser{InitEquipRefineCostParser, ParseEquipRefineCostRecord, nil},
	"type_equip_strength_cost": TDataParser{InitEquipStrengthCostParser, ParseEquipStrengthCostRecord, nil},
	"type_equipsuit":           TDataParser{InitEquipSuitParser, ParseEquipSuitRecord, nil},
	"type_shenbin":             TDataParser{InitShenBinParser, ParseShenBinRecord, FinishShenBinParser},

	//八卦镜
	"type_baguajing": TDataParser{InitBaGuaJingParser, ParseBaGuaJingRecord, nil},

	//称号表
	"type_title": TDataParser{InitTitleParser, ParseTitleRecord, nil},

	//夺粮战
	"type_foodwar_rank":  TDataParser{InitFoodWarRankAwardParser, ParseFoodWarRankAwardRecord, nil},
	"type_foodwar_award": TDataParser{InitFoodWarAwardParser, ParseFoodWarAwardRecord, nil},

	//开服基金
	"type_open_fund": TDataParser{InitOpenFundParser, ParseOpenFundRecord, nil},

	//将灵表
	"type_herosouls_link":    TDataParser{InitHeroSoulsParser, ParseHeroSoulsRecord, nil},
	"type_herosouls_map":     TDataParser{InitSoulMapParser, ParseSoulMapRecord, nil},
	"type_herosouls_store":   TDataParser{InitHeroSoulsStoreParser, ParseHeroSoulsStoreRecrod, nil},
	"type_herosouls_trials":  TDataParser{InitHeroSoulsTrialsParser, ParseHeroSoulsTrialRecord, nil},
	"type_herosouls_chapter": TDataParser{InitHeroSoulsChapterParser, ParseHeroSoulsChapterRecord, nil},

	//宝物配置表
	"type_gem":               TDataParser{InitGemParser, ParseGemRecord, nil},
	"type_gem_refine_cost":   TDataParser{InitGemRefineCostParser, ParseGemRefineCostRecord, nil},
	"type_gem_strength_cost": TDataParser{InitGemStrengthCostParser, ParseGemStrengthCostRecord, nil},

	"type_refine":   TDataParser{InitRefineParser, ParseRefineRecord, nil},
	"type_strength": TDataParser{InitStrengthParser, ParseStrengthRecord, nil},

	//强化大师表
	"type_master": TDataParser{InitMasterParser, ParseMasterRecord, nil},

	//行动力货币配置表
	"type_action":        TDataParser{InitActionParser, ParseActionRecord, nil},
	"type_money":         TDataParser{InitMoneyParser, ParseMoneyRecord, nil},
	"type_property_type": TDataParser{InitPropertyParser, ParsePropertyRecord, nil},

	//功能配置表
	"type_func_open":     TDataParser{InitFuncOpenParser, ParseFuncOpenRecord, nil},
	"type_reset_cost":    TDataParser{InitFuncCostParser, ParseFuncCostRecord, nil},
	"type_vip_privilege": TDataParser{InitVipPrivilegeParser, ParseVipPrivilegeRecord, nil},

	//挂机表
	"type_hangup": TDataParser{InitHangUpParser, ParseHangUpRecord, nil},

	//英雄配制表
	"type_hero":              TDataParser{InitHeroParser, ParseHeroRecord, nil},
	"type_hero_level":        TDataParser{InitHeroLevelParser, ParseHeroLevelRecord, nil},
	"type_hero_relation":     TDataParser{InitHeroRelationParser, ParseHeroRelationRecord, nil},
	"type_hero_break":        TDataParser{InitHeroBreakParser, ParseHeroBreakRecord, nil},
	"type_hero_break_talent": TDataParser{InitHeroBreakTalentParser, ParseHeroBreakTalentRecord, nil},
	"type_hero_culture_max":  TDataParser{InitCultureMaxParser, ParseCultureMaxRecord, nil},
	"type_hero_destiny":      TDataParser{InitHeroDestinyParser, ParseHeroDestinyRecord, nil},
	"type_hero_talent":       TDataParser{InitTalentParser, ParseTalentRecord, nil},
	"type_relation":          TDataParser{InitHeroRelationBuffParser, ParseHeroRelationBuffRecord, nil},
	"type_hero_diaowen":      TDataParser{InitDiaoWenParser, ParseDiaoWenRecord, nil},
	"type_hero_xilian":       TDataParser{InitXiLianParser, ParseXiLianRecord, nil},
	"type_hero_god":          TDataParser{InitHeroGodParser, ParseHeroGodRecord, nil},
	"type_hero_friend":       TDataParser{InitHeroFriendParser, ParseHeroFriendRecord, nil},

	//黑市表
	"type_black_market": TDataParser{InitBlackMarketParser, ParseBlackMarketRecord, nil},

	//机器人表
	"type_robot": TDataParser{InitRobotParser, ParseRobotRecord, nil},

	//充值表
	"type_monthcard": TDataParser{InitMonthCardParser, ParseMonthCardRecord, nil},
	"type_charge":    TDataParser{InitChargeItemParser, ParseChargeItemRecord, nil},

	//觉醒表
	"type_hero_wake":         TDataParser{InitWakeLevelParser, ParseWakeLevelRecord, nil},
	"type_hero_wake_compose": TDataParser{InitWakeComposeParser, ParseWakeComposeRecord, nil},

	//积分赛
	"type_jifen_duan":  TDataParser{InitScoreDwParser, ParseScoreDwRecord, nil},
	"type_jifen_award": TDataParser{InitScoreAwardParser, ParseScoreAwardRecord, nil},
	"type_jifen_store": TDataParser{InitScoreStoreParser, ParseScoreStoreRecord, nil},

	//宠物表
	"type_pet":       TDataParser{InitPetParser, ParsePetRecord, nil},
	"type_pet_level": TDataParser{InitPetLevelParser, ParsePetLevelRecord, nil},
	"type_pet_god":   TDataParser{InitPetGodParser, ParsePetGodRecord, nil},
	"type_pet_star":  TDataParser{InitPetStarParser, ParsePetStarRecord, nil},
	"type_pet_map":   TDataParser{InitPetMapParser, ParsePetMapRecord, nil},

	//卡牌大师
	"type_activity_card_exchange": TDataParser{InitCMExchangeItemParser, ParseCMExchangeItemRecord, nil},
	"type_activity_card":          TDataParser{InitCardCsvParser, ParseCardCsvRecord, nil},

	//月光集市
	"type_activity_moonlight_exch":  TDataParser{InitMoonlightShopExchangeCsv, ParseMoonlightShopExchangeCsv, nil},
	"type_activity_moonlight_goods": TDataParser{InitMoonlightGoodsCsv, ParseMoonlightGoodsCsv, nil},
	"type_activity_moonlight_award": TDataParser{InitMoonlightShopAwardCsv, ParseMoonlightShopAwardCsv, nil},

	//阵营战
	"type_crystal":       TDataParser{InitCrystalParser, ParseCrystalRecord, nil},
	"type_revive":        TDataParser{InitReviveParser, ParseReviveRecord, nil},
	"type_campbat_rank":  TDataParser{InitCampBatRankParser, ParseCampBatRankRecord, nil},
	"type_campbat_store": TDataParser{InitCampBatStoreParser, ParseCampBatStoreRecord, nil},

	//时装
	"type_fashion":          TDataParser{InitFashionParser, ParseFashionRecord, nil},
	"type_fashion_map":      TDataParser{InitFashionMapParser, ParseFashionMapRecord, nil},
	"type_fashion_strength": TDataParser{InitFashionStrengthParser, ParseFashionStrengthRecord, nil},

	//等待被解析的表
	//以下为不需要服务器读的表
	"type_copy_daily_res":  NullParser,
	"type_quality":         NullParser,
	"type_copytype":        NullParser,
	"type_model":           NullParser,
	"type_item_type":       NullParser,
	"type_camp":            NullParser,
	"type_name":            NullParser,
	"type_item_usetype":    NullParser,
	"type_monster":         NullParser,
	"type_fashion_show":    NullParser,
	"type_role_show":       NullParser,
	"type_buff_show":       NullParser,
	"type_adventure":       NullParser,
	"type_goto":            NullParser,
	"type_getway":          NullParser,
	"type_skill":           NullParser,
	"type_activity_type":   NullParser,
	"type_skill_attribute": NullParser,
	"type_mail":            NullParser,
	"type_guild_log":       NullParser,
	"type_buff_client":     NullParser,
	"type_buff_structure":  NullParser,
	"type_sgws_difficult":  NullParser,
	"type_level_function":  NullParser,
}

/* 反射解析表结构
1、表数据格式：
		数  值：1234
		字符串：zhoumf
		数值对：(24|1)(11|1)...
		数  组：10|20|30...

2、首次出现的有效行(非注释的)，即为表头

3、行列注释：已"#"开头的行，没命名/前缀"(c)"的列    有些列仅client显示用的

4、使用方式如下：
		type TTestCsv struct { // 字段须与csv表格的顺序一致
			ID    int
			Des   string
			Item  []IntPair
			Card  []IntPair
			Array []string
		}
		var G_MapCsv = make(map[int]*TTestCsv)  // map结构读表，将【&G_MapCsv】注册进G_ReflectParserMap即可自动读取
		var G_SliceCsv []TTestCsv = nil 		// 数组结构读表，注册【&G_SliceCsv】
*/
var G_ReflectParserMap = map[string]interface{}{
	// "test_name": &G_MapCsv,
	// "test_name": &G_SliceCsv,
	"type_activity_beach_goods": &G_BeachBabyGoodsCsv,
}

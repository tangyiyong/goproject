package gamedata

//InitDataParser 初始化CSV文件解析器
func InitDataParser() {

	//通用选项配置表
	DataParserMap["type_option"] = TDataParser{InitOptionParser, ParseOptionRecord}

	//副本配置表
	DataParserMap["type_copy_base"] = TDataParser{InitCopyParser, ParseCopyRecord}
	DataParserMap["type_copy_main"] = TDataParser{InitMainParser, ParseMainRecord}
	DataParserMap["type_copy_elite"] = TDataParser{InitEliteParser, ParseEliteRecord}
	DataParserMap["type_copy_daily"] = TDataParser{InitDailyParse, ParseDailyRecord}
	DataParserMap["type_copy_famous"] = TDataParser{InitFamousParser, ParseFamousRecord}

	//道具物品配置表
	DataParserMap["type_item"] = TDataParser{InitItemParser, ParseItemRecord}

	//商城配置表
	DataParserMap["type_mall"] = TDataParser{InitMallParser, ParseMallRecord}

	//VIP配置表
	DataParserMap["type_vip"] = TDataParser{InitVipParser, ParseVipRecord}
	DataParserMap["type_vip_week"] = TDataParser{InitVipWeekParser, ParseVipWeekRecord}

	//任务成就配置表
	DataParserMap["type_task"] = TDataParser{InitTaskParser, ParseTaskRecord}
	DataParserMap["type_achievement"] = TDataParser{InitAchievementParser, ParseAchievementRecord}
	DataParserMap["type_seven_activity"] = TDataParser{InitSevenActivityParser, ParseSevenActivityRecord}
	DataParserMap["type_seven_activity_store"] = TDataParser{InitSevenActivityStoreRecord, ParseSevenActivityStoreRecord}
	DataParserMap["type_task_award"] = TDataParser{InitTaskAwardParser, ParseTaskAwardRecord}
	DataParserMap["type_tasktype"] = TDataParser{InitTaskTypeParser, ParseTaskTypeRecord}

	//奖励配置表
	DataParserMap["type_award"] = TDataParser{InitAwardParser, ParseAwardRecord}

	//签到配置表
	DataParserMap["type_sign"] = TDataParser{InitSignParser, ParseSignRecord}
	DataParserMap["type_sign_plus"] = TDataParser{InitSignPlusParser, ParseSignPlusRecord}

	//神将商店
	DataParserMap["type_store"] = TDataParser{InitStoreParse, ParseStoreRecord}

	//三国志配置表
	DataParserMap["type_sanguozhi"] = TDataParser{InitSanGuoZhiParser, ParseSanGuoZhiRecord}

	//招贤配置表
	DataParserMap["type_summon_config"] = TDataParser{InitSummonConfigParser, ParseSummonConfigRecord}

	//竞技场配置表
	DataParserMap["type_arena"] = TDataParser{InitArenaParser, ParseArenaRecord}
	DataParserMap["type_arena_rank"] = TDataParser{InitArenaRankParser, ParseArenaRankRecord}
	DataParserMap["type_arena_store"] = TDataParser{InitArenaStoreParser, ParseArenaStoreRecord}
	DataParserMap["type_arena_money"] = TDataParser{InitArenaMoneyParser, ParseArenaMoneyRecord}

	//夺宝配置表
	DataParserMap["type_rob"] = TDataParser{InitRobParser, ParseRobRecord}
	DataParserMap["type_treasure_melting"] = TDataParser{InitTreasureMeltingParser, ParseTreasureMeltingRecord}

	//三国无双配置表
	DataParserMap["type_sgws_chapter"] = TDataParser{InitSangokuMusouChapter, ParseSangokuMusouChapterRecord}
	DataParserMap["type_sgws_chapter_attr"] = TDataParser{InitSangokuMusouAttrMarkupParser, ParseSangokuMusouAttrMarkupRecord}
	DataParserMap["type_sgws_chapter_award"] = TDataParser{InitSangokuMusouChapterAwardParser, ParseSangokuMusouChapterAwardRecord}
	DataParserMap["type_sgws_elite_copy"] = TDataParser{InitSangokuMusouEliteCopyParser, ParseSangokuMusouEliteCopyRecord}
	DataParserMap["type_sgws_store"] = TDataParser{InitSangokuMusouStoreParser, ParseSangokuMusouStoreRecord}
	DataParserMap["type_sgws_treasure"] = TDataParser{InitSangokuMusouSaleParser, ParseSangokuMusouSaleRecord}

	//领地攻伐表
	DataParserMap["type_territory"] = TDataParser{InitTerritoryParser, ParseTerritoryRecord}
	DataParserMap["type_territoryaward"] = TDataParser{InitTerritoryAwardParser, ParseTerritoryAwardRecord}
	DataParserMap["type_territoryskill"] = TDataParser{InitTerritorySkillParser, ParseTerritorySkillRecord}
	DataParserMap["type_territorypatrol"] = TDataParser{InitTerritoryPatrolParser, ParseTerritoryPatrolRecord}
	DataParserMap["type_territoryawardtype"] = TDataParser{InitTerritoryAwardTypeParser, ParseTerritoryAwardTypeRecord}

	//叛军围剿表
	DataParserMap["type_rebel_siege"] = TDataParser{InitRebelSiegeParser, ParseRebelSiegeRecord}
	DataParserMap["type_rebel_action"] = TDataParser{InitRebelActionAwardParser, ParseRebelActionAwardRecord}
	DataParserMap["type_rebel_award"] = TDataParser{InitExploitAwardParser, ParseExploitAwardRecord}
	DataParserMap["type_rebel_store"] = TDataParser{InitExploitStoreParser, ParseExploitStoreRecord}
	DataParserMap["type_rebel_activity"] = TDataParser{InitRebelActivityParser, ParseRebelActivityRecord}
	DataParserMap["type_rebel_rank_award"] = TDataParser{InitRebelRankAwardParser, ParseRebelRankAwardRecord}

	//挖矿表
	DataParserMap["type_mining_stone_num"] = TDataParser{InitMiningStoneRandomParser, ParseMiningStoneRecord}
	DataParserMap["type_mining_element"] = TDataParser{InitMiningElementParser, ParseMiningElementRecord}
	DataParserMap["type_mining_event"] = TDataParser{InitMiningEventParser, ParserMiningEventRecord}
	DataParserMap["type_mining_black_market"] = TDataParser{InitMiningEventBlackMarketParser, ParserMiningEventBlackMarketRecord}
	DataParserMap["type_mining_question"] = NullParser
	DataParserMap["type_mining_buff"] = TDataParser{InitMiningEventBuffParser, ParserMiningEventBuffRecord}
	DataParserMap["type_mining_treasure"] = TDataParser{InitMiningEventTreasureParser, ParserMiningEventTreasureRecord}
	DataParserMap["type_mining_monster"] = TDataParser{InitMiningEventMonsterPerser, ParseMiningEventMonsterRecord}
	DataParserMap["type_mining_award"] = TDataParser{InitMiningAwardParser, ParserMiningAwardRecord}
	DataParserMap["type_mining_guaji"] = TDataParser{InitMiningGuaJiParser, ParserMiningGuaJiRecord}
	DataParserMap["type_mining_rand"] = TDataParser{InitMiningRandParser, ParserMiningRandRecord}

	//活动配置表
	DataParserMap["type_activity"] = TDataParser{InitActivityParser, ParseActivityRecord}
	DataParserMap["type_activity_competition"] = TDataParser{InitCompetitionParser, ParseCompetitionRecord}
	DataParserMap["type_activity_action"] = TDataParser{InitRecvActionParser, ParseRecvActionRecord}
	DataParserMap["type_activity_discount"] = TDataParser{InitDiscountSaleParser, ParseDiscountSaleRecord}
	DataParserMap["type_activity_login"] = TDataParser{InitActivityLoginParser, ParseActivityLoginRecord}
	DataParserMap["type_activity_moneygod"] = TDataParser{InitActivityMoneyParser, ParseActivityMoneyRecord}
	DataParserMap["type_activity_recharge"] = TDataParser{InitActivityRechargeParser, ParseActivityRechargeRecord}
	DataParserMap["type_activity_limitdaily"] = TDataParser{InitActivityLimitDailyParser, ParseActivityLimitDailyRecord}
	DataParserMap["type_activity_hunt_map"] = TDataParser{InitHuntTreasureMapParser, ParseHuntTreasureMapRecord}
	DataParserMap["type_activity_hunt_store"] = TDataParser{InitHuntTreasureStoreParser, ParseHuntTreasureStoreRecord}
	DataParserMap["type_activity_hunt_turn"] = TDataParser{InitHuntTreasureAwardParser, ParseHuntTreasureAwardRecord}
	DataParserMap["type_activity_rank"] = TDataParser{InitOperationalRankAwardParser, ParseOperationalRankAwardRecord}
	DataParserMap["type_activity_lucky_wheel"] = TDataParser{InitLuckyWheelParser, ParseLuckyWheelRecord}
	DataParserMap["type_activity_group_purchase"] = TDataParser{InitGroupPurchaseParser, ParseGroupPurchaseRecord}
	DataParserMap["type_activity_group_score"] = TDataParser{InitGroupPurchaseScoreParser, ParseGroupPurchaseScoreRecord}
	DataParserMap["type_activity_festival_task"] = TDataParser{InitFestivalTaskParser, ParseFestivalTaskRecord}
	DataParserMap["type_activity_festival_exchange"] = TDataParser{InitFestivalExchangeParser, ParseFestivalExchangeRecord}
	DataParserMap["type_activity_week_award"] = TDataParser{InitActivityWeekAwardParser, ParseActivityWeekAwardRecord}
	DataParserMap["type_activity_level_gift"] = TDataParser{InitActivityLevelGiftParser, ParseActivityLevelGiftRecord}
	DataParserMap["type_activity_month_fund"] = TDataParser{InitActivityMonthFundParser, ParseActivityMonthFundRecord}

	//公会配置表
	DataParserMap["type_guild_base"] = TDataParser{InitGuildParser, ParseGuildBaseRecord}
	DataParserMap["type_guild_copy"] = TDataParser{InitGuildCopyParser, ParseGuildCopyRecord}
	DataParserMap["type_guild_copy_award"] = TDataParser{InitGuildCopyAwardParser, ParseGuildCopyAwardRecord}
	DataParserMap["type_guild_role"] = TDataParser{InitGuildRoleParser, ParseGuildRoleRecord}
	DataParserMap["type_guild_sacrifice"] = TDataParser{InitGuildSacrificeParser, ParseGuildSacrificeRecord}
	DataParserMap["type_guild_sacrifice_award"] = TDataParser{InitGuildSacrificeAwardParser, ParseGuildSacrificeAwardRecord}
	DataParserMap["type_guild_store"] = TDataParser{InitGuildStoreParser, ParseGuildStoreRecord}
	DataParserMap["type_guild_skill"] = TDataParser{InitGuildSkillParser, ParseGuildSkillRecord}
	DataParserMap["type_guild_skill_level"] = TDataParser{InitGuildSkillLimitParser, ParseGuildSkillLimitRecord}

	//装备配置表
	DataParserMap["type_equipment"] = TDataParser{InitEquipParser, ParseEquipRecord}
	DataParserMap["type_equip_star"] = TDataParser{InitEquipStarParser, ParseEquipStarRecord}
	DataParserMap["type_equip_refine_cost"] = TDataParser{InitEquipRefineCostParser, ParseEquipRefineCostRecord}
	DataParserMap["type_equip_strength_cost"] = TDataParser{InitEquipStrengthCostParser, ParseEquipStrengthCostRecord}
	DataParserMap["type_equipsuit"] = TDataParser{InitEquipSuitParser, ParseEquipSuitRecord}

	//八卦镜
	DataParserMap["type_baguajing"] = TDataParser{InitBaGuaJingParser, ParseBaGuaJingRecord}

	//称号表
	DataParserMap["type_title"] = TDataParser{InitTitleParser, ParseTitleRecord}

	//夺粮战
	DataParserMap["type_foodwar_rank"] = TDataParser{InitFoodWarRankAwardParser, ParseFoodWarRankAwardRecord}
	DataParserMap["type_foodwar_award"] = TDataParser{InitFoodWarAwardParser, ParseFoodWarAwardRecord}

	//开服基金
	DataParserMap["type_open_fund"] = TDataParser{InitOpenFundParser, ParseOpenFundRecord}

	//将灵表
	DataParserMap["type_herosouls_link"] = TDataParser{InitHeroSoulsParser, ParseHeroSoulsRecord}
	DataParserMap["type_herosouls_map"] = TDataParser{InitSoulMapParser, ParseSoulMapRecord}
	DataParserMap["type_herosouls_store"] = TDataParser{InitHeroSoulsStoreParser, ParseHeroSoulsStoreRecrod}
	DataParserMap["type_herosouls_trials"] = TDataParser{InitHeroSoulsTrialsParser, ParseHeroSoulsTrialRecord}
	DataParserMap["type_herosouls_chapter"] = TDataParser{InitHeroSoulsChapterParser, ParseHeroSoulsChapterRecord}

	//宝物配置表
	DataParserMap["type_gem"] = TDataParser{InitGemParser, ParseGemRecord}
	DataParserMap["type_gem_refine_cost"] = TDataParser{InitGemRefineCostParser, ParseGemRefineCostRecord}
	DataParserMap["type_gem_strength_cost"] = TDataParser{InitGemStrengthCostParser, ParseGemStrengthCostRecord}

	DataParserMap["type_refine"] = TDataParser{InitRefineParser, ParseRefineRecord}
	DataParserMap["type_strength"] = TDataParser{InitStrengthParser, ParseStrengthRecord}

	//强化大师表
	DataParserMap["type_master"] = TDataParser{InitMasterParser, ParseMasterRecord}

	//行动力货币配置表
	DataParserMap["type_action"] = TDataParser{InitActionParser, ParseActionRecord}
	DataParserMap["type_money"] = TDataParser{InitMoneyParser, ParseMoneyRecord}
	DataParserMap["type_property_type"] = TDataParser{InitPropertyParser, ParsePropertyRecord}

	//功能配置表
	DataParserMap["type_func_open"] = TDataParser{InitFuncOpenParser, ParseFuncOpenRecord}
	DataParserMap["type_reset_cost"] = TDataParser{InitFuncCostParser, ParseFuncCostRecord}
	DataParserMap["type_vip_privilege"] = TDataParser{InitVipPrivilegeParser, ParseVipPrivilegeRecord}

	//挂机表
	DataParserMap["type_hangup"] = TDataParser{InitHangUpParser, ParseHangUpRecord}

	//英雄配制表
	DataParserMap["type_hero"] = TDataParser{InitHeroParser, ParseHeroRecord}
	DataParserMap["type_hero_level"] = TDataParser{InitHeroLevelParser, ParseHeroLevelRecord}
	DataParserMap["type_hero_relation"] = TDataParser{InitHeroRelationParser, ParseHeroRelationRecord}
	DataParserMap["type_hero_break"] = TDataParser{InitHeroBreakParser, ParseHeroBreakRecord}
	DataParserMap["type_hero_break_talent"] = TDataParser{InitHeroBreakTalentParser, ParseHeroBreakTalentRecord}
	DataParserMap["type_hero_culture_max"] = TDataParser{InitCultureMaxParser, ParseCultureMaxRecord}
	DataParserMap["type_hero_destiny"] = TDataParser{InitHeroDestinyParser, ParseHeroDestinyRecord}
	DataParserMap["type_hero_relation"] = TDataParser{InitHeroRelationParser, ParseHeroRelationRecord}
	DataParserMap["type_hero_talent"] = TDataParser{InitTalentParser, ParseTalentRecord}
	DataParserMap["type_relation"] = TDataParser{InitHeroRelationBuffParser, ParseHeroRelationBuffRecord}
	DataParserMap["type_hero_diaowen"] = TDataParser{InitDiaoWenParser, ParseDiaoWenRecord}
	DataParserMap["type_hero_xilian"] = TDataParser{InitXiLianParser, ParseXiLianRecord}
	DataParserMap["type_hero_god"] = TDataParser{InitHeroGodParser, ParseHeroGodRecord}
	DataParserMap["type_hero_friend"] = TDataParser{InitHeroFriendParser, ParseHeroFriendRecord}

	//黑市表
	DataParserMap["type_black_market"] = TDataParser{InitBlackMarketParser, ParseBlackMarketRecord}

	//机器人表
	DataParserMap["type_robot"] = TDataParser{InitRobotParser, ParseRobotRecord}

	//充值表
	DataParserMap["type_monthcard"] = TDataParser{InitMonthCardParser, ParseMonthCardRecord}
	DataParserMap["type_charge"] = TDataParser{InitChargeItemParser, ParseChargeItemRecord}

	//觉醒表
	DataParserMap["type_hero_wake"] = TDataParser{InitWakeLevelParser, ParseWakeLevelRecord}
	DataParserMap["type_hero_wake_compose"] = TDataParser{InitWakeComposeParser, ParseWakeComposeRecord}

	//积分赛
	DataParserMap["type_jifen_duan"] = TDataParser{InitScoreDwParser, ParseScoreDwRecord}
	DataParserMap["type_jifen_award"] = TDataParser{InitScoreAwardParser, ParseScoreAwardRecord}
	DataParserMap["type_jifen_store"] = TDataParser{InitScoreStoreParser, ParseScoreStoreRecord}

	//宠物表
	DataParserMap["type_pet"] = TDataParser{InitPetParser, ParsePetRecord}
	DataParserMap["type_pet_level"] = TDataParser{InitPetLevelParser, ParsePetLevelRecord}
	DataParserMap["type_pet_god"] = TDataParser{InitPetGodParser, ParsePetGodRecord}
	DataParserMap["type_pet_star"] = TDataParser{InitPetStarParser, ParsePetStarRecord}
	DataParserMap["type_pet_map"] = TDataParser{InitPetMapParser, ParsePetMapRecord}

	//卡牌大师
	DataParserMap["type_activity_card_exchange"] = TDataParser{InitCMExchangeItemParser, ParseCMExchangeItemRecord}
	DataParserMap["type_activity_card"] = TDataParser{InitCardCsvParser, ParseCardCsvRecord}

	//月光集市
	DataParserMap["type_activity_moonlight_exch"] = TDataParser{InitMoonlightShopExchangeCsv, ParseMoonlightShopExchangeCsv}
	DataParserMap["type_activity_moonlight_goods"] = TDataParser{InitMoonlightGoodsCsv, ParseMoonlightGoodsCsv}
	DataParserMap["type_activity_moonlight_award"] = TDataParser{InitMoonlightShopAwardCsv, ParseMoonlightShopAwardCsv}

	//阵营战
	DataParserMap["type_crystal"] = TDataParser{InitCrystalParser, ParseCrystalRecord}
	DataParserMap["type_revive"] = TDataParser{InitReviveParser, ParseReviveRecord}
	DataParserMap["type_campbat_rank"] = TDataParser{InitCampBatRankParser, ParseCampBatRankRecord}
	DataParserMap["type_campbat_store"] = TDataParser{InitCampBatStoreParser, ParseCampBatStoreRecord}

	//时装
	DataParserMap["type_fashion"] = TDataParser{InitFashionParser, ParseFashionRecord}
	DataParserMap["type_fashion_map"] = TDataParser{InitFashionMapParser, ParseFashionMapRecord}
	DataParserMap["type_fashion_strength"] = TDataParser{InitFashionStrengthParser, ParseFashionStrengthRecord}

	//等待被解析的表
	//以下为不需要服务器读的表
	DataParserMap["type_copy_daily_res"] = NullParser
	DataParserMap["type_quality"] = NullParser
	DataParserMap["type_copytype"] = NullParser
	DataParserMap["type_model"] = NullParser
	DataParserMap["type_item_type"] = NullParser
	DataParserMap["type_camp"] = NullParser
	DataParserMap["type_name"] = NullParser
	DataParserMap["type_item_usetype"] = NullParser
	DataParserMap["type_monster"] = NullParser
	DataParserMap["type_fashion_show"] = NullParser
	DataParserMap["type_role_show"] = NullParser
	DataParserMap["type_buff_show"] = NullParser
	DataParserMap["type_adventure"] = NullParser
	DataParserMap["type_goto"] = NullParser
	DataParserMap["type_getway"] = NullParser
	DataParserMap["type_skill"] = NullParser
	DataParserMap["type_activity_type"] = NullParser
	DataParserMap["type_skill_attribute"] = NullParser
	DataParserMap["type_mail"] = NullParser
	DataParserMap["type_guild_log"] = NullParser
	DataParserMap["type_buff_client"] = NullParser
	DataParserMap["type_buff_structure"] = NullParser
	DataParserMap["type_sgws_difficult"] = NullParser
	DataParserMap["type_level_function"] = NullParser
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
func InitReflectParser() {
	G_ReflectParserMap = map[string]interface{}{
		// "test_name": &G_MapCsv,
		// "test_name": &G_SliceCsv,
		"type_activity_beach_goods": &G_BeachBabyGoodsCsv,
	}
}

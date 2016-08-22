package gamedata

//InitDataParser 初始化CSV文件解析器
func InitDataParser() {

	//通用选项配置表
	DataParserMap["type_option"] = TDataParser{InitOptionParser, ParseOptionRecord, nil}

	//副本配置表
	DataParserMap["type_copy_base"] = TDataParser{InitCopyParser, ParseCopyRecord, nil}
	DataParserMap["type_copy_main"] = TDataParser{InitMainParser, ParseMainRecord, nil}
	DataParserMap["type_copy_elite"] = TDataParser{InitEliteParser, ParseEliteRecord, nil}
	DataParserMap["type_copy_daily"] = TDataParser{InitDailyParse, ParseDailyRecord, nil}
	DataParserMap["type_copy_famous"] = TDataParser{InitFamousParser, ParseFamousRecord, nil}

	//道具物品配置表
	DataParserMap["type_item"] = TDataParser{InitItemParser, ParseItemRecord, nil}

	//商城配置表
	DataParserMap["type_mall"] = TDataParser{InitMallParser, ParseMallRecord, nil}

	//VIP配置表
	DataParserMap["type_vip"] = TDataParser{InitVipParser, ParseVipRecord, nil}
	DataParserMap["type_vip_week"] = TDataParser{InitVipWeekParser, ParseVipWeekRecord, nil}

	//任务成就配置表
	DataParserMap["type_task"] = TDataParser{InitTaskParser, ParseTaskRecord, nil}
	DataParserMap["type_achievement"] = TDataParser{InitAchievementParser, ParseAchievementRecord, nil}
	DataParserMap["type_seven_activity"] = TDataParser{InitSevenActivityParser, ParseSevenActivityRecord, nil}
	DataParserMap["type_seven_activity_store"] = TDataParser{InitSevenActivityStoreRecord, ParseSevenActivityStoreRecord, nil}
	DataParserMap["type_task_award"] = TDataParser{InitTaskAwardParser, ParseTaskAwardRecord, nil}
	DataParserMap["type_tasktype"] = TDataParser{InitTaskTypeParser, ParseTaskTypeRecord, nil}

	//奖励配置表
	DataParserMap["type_award"] = TDataParser{InitAwardParser, ParseAwardRecord, nil}

	//签到配置表
	DataParserMap["type_sign"] = TDataParser{InitSignParser, ParseSignRecord, nil}
	DataParserMap["type_sign_plus"] = TDataParser{InitSignPlusParser, ParseSignPlusRecord, nil}

	//神将商店
	DataParserMap["type_store"] = TDataParser{InitStoreParse, ParseStoreRecord, nil}

	//三国志配置表
	DataParserMap["type_sanguozhi"] = TDataParser{InitSanGuoZhiParser, ParseSanGuoZhiRecord, nil}

	//招贤配置表
	DataParserMap["type_summon_config"] = TDataParser{InitSummonConfigParser, ParseSummonConfigRecord, nil}

	//竞技场配置表
	DataParserMap["type_arena"] = TDataParser{InitArenaParser, ParseArenaRecord, nil}
	DataParserMap["type_arena_rank"] = TDataParser{InitArenaRankParser, ParseArenaRankRecord, nil}
	DataParserMap["type_arena_store"] = TDataParser{InitArenaStoreParser, ParseArenaStoreRecord, nil}
	DataParserMap["type_arena_money"] = TDataParser{InitArenaMoneyParser, ParseArenaMoneyRecord, nil}

	//夺宝配置表
	DataParserMap["type_rob"] = TDataParser{InitRobParser, ParseRobRecord, nil}
	DataParserMap["type_treasure_melting"] = TDataParser{InitTreasureMeltingParser, ParseTreasureMeltingRecord, nil}

	//三国无双配置表
	DataParserMap["type_sgws_chapter"] = TDataParser{InitSangokuMusouChapter, ParseSangokuMusouChapterRecord, nil}
	DataParserMap["type_sgws_chapter_attr"] = TDataParser{InitSangokuMusouAttrMarkupParser, ParseSangokuMusouAttrMarkupRecord, nil}
	DataParserMap["type_sgws_chapter_award"] = TDataParser{InitSangokuMusouChapterAwardParser, ParseSangokuMusouChapterAwardRecord, nil}
	DataParserMap["type_sgws_elite_copy"] = TDataParser{InitSangokuMusouEliteCopyParser, ParseSangokuMusouEliteCopyRecord, nil}
	DataParserMap["type_sgws_store"] = TDataParser{InitSangokuMusouStoreParser, ParseSangokuMusouStoreRecord, nil}
	DataParserMap["type_sgws_treasure"] = TDataParser{InitSangokuMusouSaleParser, ParseSangokuMusouSaleRecord, nil}

	//领地攻伐表
	DataParserMap["type_territory"] = TDataParser{InitTerritoryParser, ParseTerritoryRecord, nil}
	DataParserMap["type_territoryaward"] = TDataParser{InitTerritoryAwardParser, ParseTerritoryAwardRecord, nil}
	DataParserMap["type_territoryskill"] = TDataParser{InitTerritorySkillParser, ParseTerritorySkillRecord, nil}
	DataParserMap["type_territorypatrol"] = TDataParser{InitTerritoryPatrolParser, ParseTerritoryPatrolRecord, nil}
	DataParserMap["type_territoryawardtype"] = TDataParser{InitTerritoryAwardTypeParser, ParseTerritoryAwardTypeRecord, nil}

	//叛军围剿表
	DataParserMap["type_rebel_siege"] = TDataParser{InitRebelSiegeParser, ParseRebelSiegeRecord, nil}
	DataParserMap["type_rebel_action"] = TDataParser{InitRebelActionAwardParser, ParseRebelActionAwardRecord, nil}
	DataParserMap["type_rebel_award"] = TDataParser{InitExploitAwardParser, ParseExploitAwardRecord, nil}
	DataParserMap["type_rebel_store"] = TDataParser{InitExploitStoreParser, ParseExploitStoreRecord, nil}
	DataParserMap["type_rebel_activity"] = TDataParser{InitRebelActivityParser, ParseRebelActivityRecord, nil}
	DataParserMap["type_rebel_rank_award"] = TDataParser{InitRebelRankAwardParser, ParseRebelRankAwardRecord, nil}

	//挖矿表
	DataParserMap["type_mining_stone_num"] = TDataParser{InitMiningStoneRandomParser, ParseMiningStoneRecord, nil}
	DataParserMap["type_mining_element"] = TDataParser{InitMiningElementParser, ParseMiningElementRecord, nil}
	DataParserMap["type_mining_event"] = TDataParser{InitMiningEventParser, ParserMiningEventRecord, nil}
	DataParserMap["type_mining_black_market"] = TDataParser{InitMiningEventBlackMarketParser, ParserMiningEventBlackMarketRecord, nil}
	DataParserMap["type_mining_question"] = NullParser
	DataParserMap["type_mining_buff"] = TDataParser{InitMiningEventBuffParser, ParserMiningEventBuffRecord, nil}
	DataParserMap["type_mining_treasure"] = TDataParser{InitMiningEventTreasureParser, ParserMiningEventTreasureRecord, nil}
	DataParserMap["type_mining_monster"] = TDataParser{InitMiningEventMonsterPerser, ParseMiningEventMonsterRecord, nil}
	DataParserMap["type_mining_award"] = TDataParser{InitMiningAwardParser, ParserMiningAwardRecord, nil}
	DataParserMap["type_mining_guaji"] = TDataParser{InitMiningGuaJiParser, ParserMiningGuaJiRecord, nil}
	DataParserMap["type_mining_rand"] = TDataParser{InitMiningRandParser, ParserMiningRandRecord, nil}

	//活动配置表
	DataParserMap["type_activity"] = TDataParser{InitActivityParser, ParseActivityRecord, nil}
	DataParserMap["type_activity_competition"] = TDataParser{InitCompetitionParser, ParseCompetitionRecord, nil}
	DataParserMap["type_activity_action"] = TDataParser{InitRecvActionParser, ParseRecvActionRecord, nil}
	DataParserMap["type_activity_discount"] = TDataParser{InitDiscountSaleParser, ParseDiscountSaleRecord, nil}
	DataParserMap["type_activity_login"] = TDataParser{InitActivityLoginParser, ParseActivityLoginRecord, nil}
	DataParserMap["type_activity_moneygod"] = TDataParser{InitActivityMoneyParser, ParseActivityMoneyRecord, nil}
	DataParserMap["type_activity_recharge"] = TDataParser{InitActivityRechargeParser, ParseActivityRechargeRecord, nil}
	DataParserMap["type_activity_limitdaily"] = TDataParser{InitActivityLimitDailyParser, ParseActivityLimitDailyRecord, nil}
	DataParserMap["type_activity_hunt_map"] = TDataParser{InitHuntTreasureMapParser, ParseHuntTreasureMapRecord, nil}
	DataParserMap["type_activity_hunt_store"] = TDataParser{InitHuntTreasureStoreParser, ParseHuntTreasureStoreRecord, nil}
	DataParserMap["type_activity_hunt_turn"] = TDataParser{InitHuntTreasureAwardParser, ParseHuntTreasureAwardRecord, nil}
	DataParserMap["type_activity_rank"] = TDataParser{InitOperationalRankAwardParser, ParseOperationalRankAwardRecord, nil}
	DataParserMap["type_activity_lucky_wheel"] = TDataParser{InitLuckyWheelParser, ParseLuckyWheelRecord, nil}
	DataParserMap["type_activity_group_purchase"] = TDataParser{InitGroupPurchaseParser, ParseGroupPurchaseRecord, nil}
	DataParserMap["type_activity_group_score"] = TDataParser{InitGroupPurchaseScoreParser, ParseGroupPurchaseScoreRecord, nil}
	DataParserMap["type_activity_festival_task"] = TDataParser{InitFestivalTaskParser, ParseFestivalTaskRecord, nil}
	DataParserMap["type_activity_festival_exchange"] = TDataParser{InitFestivalExchangeParser, ParseFestivalExchangeRecord, nil}
	DataParserMap["type_activity_week_award"] = TDataParser{InitActivityWeekAwardParser, ParseActivityWeekAwardRecord, nil}
	DataParserMap["type_activity_level_gift"] = TDataParser{InitActivityLevelGiftParser, ParseActivityLevelGiftRecord, nil}
	DataParserMap["type_activity_month_fund"] = TDataParser{InitActivityMonthFundParser, ParseActivityMonthFundRecord, nil}
	DataParserMap["type_activity_limitsale"] = TDataParser{InitLimitSaleItemParser, ParseLimitSaleItemRecord, nil}
	DataParserMap["type_activity_limitsale_award"] = TDataParser{InitLimitSaleAllAwardParser, ParseLimitSaleAllAwardRecord, nil}

	//公会配置表
	DataParserMap["type_guild_base"] = TDataParser{InitGuildParser, ParseGuildBaseRecord, nil}
	DataParserMap["type_guild_copy"] = TDataParser{InitGuildCopyParser, ParseGuildCopyRecord, nil}
	DataParserMap["type_guild_copy_award"] = TDataParser{InitGuildCopyAwardParser, ParseGuildCopyAwardRecord, nil}
	DataParserMap["type_guild_role"] = TDataParser{InitGuildRoleParser, ParseGuildRoleRecord, nil}
	DataParserMap["type_guild_sacrifice"] = TDataParser{InitGuildSacrificeParser, ParseGuildSacrificeRecord, nil}
	DataParserMap["type_guild_sacrifice_award"] = TDataParser{InitGuildSacrificeAwardParser, ParseGuildSacrificeAwardRecord, nil}
	DataParserMap["type_guild_store"] = TDataParser{InitGuildStoreParser, ParseGuildStoreRecord, nil}
	DataParserMap["type_guild_skill"] = TDataParser{InitGuildSkillParser, ParseGuildSkillRecord, nil}
	DataParserMap["type_guild_skill_level"] = TDataParser{InitGuildSkillLimitParser, ParseGuildSkillLimitRecord, nil}

	//装备配置表
	DataParserMap["type_equipment"] = TDataParser{InitEquipParser, ParseEquipRecord, nil}
	DataParserMap["type_equip_star"] = TDataParser{InitEquipStarParser, ParseEquipStarRecord, nil}
	DataParserMap["type_equip_refine_cost"] = TDataParser{InitEquipRefineCostParser, ParseEquipRefineCostRecord, nil}
	DataParserMap["type_equip_strength_cost"] = TDataParser{InitEquipStrengthCostParser, ParseEquipStrengthCostRecord, nil}
	DataParserMap["type_equipsuit"] = TDataParser{InitEquipSuitParser, ParseEquipSuitRecord, nil}
	DataParserMap["type_shenbin"] = TDataParser{InitShenBinParser, ParseShenBinRecord, FinishShenBinParser}

	//八卦镜
	DataParserMap["type_baguajing"] = TDataParser{InitBaGuaJingParser, ParseBaGuaJingRecord, nil}

	//称号表
	DataParserMap["type_title"] = TDataParser{InitTitleParser, ParseTitleRecord, nil}

	//夺粮战
	DataParserMap["type_foodwar_rank"] = TDataParser{InitFoodWarRankAwardParser, ParseFoodWarRankAwardRecord, nil}
	DataParserMap["type_foodwar_award"] = TDataParser{InitFoodWarAwardParser, ParseFoodWarAwardRecord, nil}

	//开服基金
	DataParserMap["type_open_fund"] = TDataParser{InitOpenFundParser, ParseOpenFundRecord, nil}

	//将灵表
	DataParserMap["type_herosouls_link"] = TDataParser{InitHeroSoulsParser, ParseHeroSoulsRecord, nil}
	DataParserMap["type_herosouls_map"] = TDataParser{InitSoulMapParser, ParseSoulMapRecord, nil}
	DataParserMap["type_herosouls_store"] = TDataParser{InitHeroSoulsStoreParser, ParseHeroSoulsStoreRecrod, nil}
	DataParserMap["type_herosouls_trials"] = TDataParser{InitHeroSoulsTrialsParser, ParseHeroSoulsTrialRecord, nil}
	DataParserMap["type_herosouls_chapter"] = TDataParser{InitHeroSoulsChapterParser, ParseHeroSoulsChapterRecord, nil}

	//宝物配置表
	DataParserMap["type_gem"] = TDataParser{InitGemParser, ParseGemRecord, nil}
	DataParserMap["type_gem_refine_cost"] = TDataParser{InitGemRefineCostParser, ParseGemRefineCostRecord, nil}
	DataParserMap["type_gem_strength_cost"] = TDataParser{InitGemStrengthCostParser, ParseGemStrengthCostRecord, nil}

	DataParserMap["type_refine"] = TDataParser{InitRefineParser, ParseRefineRecord, nil}
	DataParserMap["type_strength"] = TDataParser{InitStrengthParser, ParseStrengthRecord, nil}

	//强化大师表
	DataParserMap["type_master"] = TDataParser{InitMasterParser, ParseMasterRecord, nil}

	//行动力货币配置表
	DataParserMap["type_action"] = TDataParser{InitActionParser, ParseActionRecord, nil}
	DataParserMap["type_money"] = TDataParser{InitMoneyParser, ParseMoneyRecord, nil}
	DataParserMap["type_property_type"] = TDataParser{InitPropertyParser, ParsePropertyRecord, nil}

	//功能配置表
	DataParserMap["type_func_open"] = TDataParser{InitFuncOpenParser, ParseFuncOpenRecord, nil}
	DataParserMap["type_reset_cost"] = TDataParser{InitFuncCostParser, ParseFuncCostRecord, nil}
	DataParserMap["type_vip_privilege"] = TDataParser{InitVipPrivilegeParser, ParseVipPrivilegeRecord, nil}

	//挂机表
	DataParserMap["type_hangup"] = TDataParser{InitHangUpParser, ParseHangUpRecord, nil}

	//英雄配制表
	DataParserMap["type_hero"] = TDataParser{InitHeroParser, ParseHeroRecord, nil}
	DataParserMap["type_hero_level"] = TDataParser{InitHeroLevelParser, ParseHeroLevelRecord, nil}
	DataParserMap["type_hero_relation"] = TDataParser{InitHeroRelationParser, ParseHeroRelationRecord, nil}
	DataParserMap["type_hero_break"] = TDataParser{InitHeroBreakParser, ParseHeroBreakRecord, nil}
	DataParserMap["type_hero_break_talent"] = TDataParser{InitHeroBreakTalentParser, ParseHeroBreakTalentRecord, nil}
	DataParserMap["type_hero_culture_max"] = TDataParser{InitCultureMaxParser, ParseCultureMaxRecord, nil}
	DataParserMap["type_hero_destiny"] = TDataParser{InitHeroDestinyParser, ParseHeroDestinyRecord, nil}
	DataParserMap["type_hero_relation"] = TDataParser{InitHeroRelationParser, ParseHeroRelationRecord, nil}
	DataParserMap["type_hero_talent"] = TDataParser{InitTalentParser, ParseTalentRecord, nil}
	DataParserMap["type_relation"] = TDataParser{InitHeroRelationBuffParser, ParseHeroRelationBuffRecord, nil}
	DataParserMap["type_hero_diaowen"] = TDataParser{InitDiaoWenParser, ParseDiaoWenRecord, nil}
	DataParserMap["type_hero_xilian"] = TDataParser{InitXiLianParser, ParseXiLianRecord, nil}
	DataParserMap["type_hero_god"] = TDataParser{InitHeroGodParser, ParseHeroGodRecord, nil}
	DataParserMap["type_hero_friend"] = TDataParser{InitHeroFriendParser, ParseHeroFriendRecord, nil}

	//黑市表
	DataParserMap["type_black_market"] = TDataParser{InitBlackMarketParser, ParseBlackMarketRecord, nil}

	//机器人表
	DataParserMap["type_robot"] = TDataParser{InitRobotParser, ParseRobotRecord, nil}

	//充值表
	DataParserMap["type_monthcard"] = TDataParser{InitMonthCardParser, ParseMonthCardRecord, nil}
	DataParserMap["type_charge"] = TDataParser{InitChargeItemParser, ParseChargeItemRecord, nil}

	//觉醒表
	DataParserMap["type_hero_wake"] = TDataParser{InitWakeLevelParser, ParseWakeLevelRecord, nil}
	DataParserMap["type_hero_wake_compose"] = TDataParser{InitWakeComposeParser, ParseWakeComposeRecord, nil}

	//积分赛
	DataParserMap["type_jifen_duan"] = TDataParser{InitScoreDwParser, ParseScoreDwRecord, nil}
	DataParserMap["type_jifen_award"] = TDataParser{InitScoreAwardParser, ParseScoreAwardRecord, nil}
	DataParserMap["type_jifen_store"] = TDataParser{InitScoreStoreParser, ParseScoreStoreRecord, nil}

	//宠物表
	DataParserMap["type_pet"] = TDataParser{InitPetParser, ParsePetRecord, nil}
	DataParserMap["type_pet_level"] = TDataParser{InitPetLevelParser, ParsePetLevelRecord, nil}
	DataParserMap["type_pet_god"] = TDataParser{InitPetGodParser, ParsePetGodRecord, nil}
	DataParserMap["type_pet_star"] = TDataParser{InitPetStarParser, ParsePetStarRecord, nil}
	DataParserMap["type_pet_map"] = TDataParser{InitPetMapParser, ParsePetMapRecord, nil}

	//卡牌大师
	DataParserMap["type_activity_card_exchange"] = TDataParser{InitCMExchangeItemParser, ParseCMExchangeItemRecord, nil}
	DataParserMap["type_activity_card"] = TDataParser{InitCardCsvParser, ParseCardCsvRecord, nil}

	//月光集市
	DataParserMap["type_activity_moonlight_exch"] = TDataParser{InitMoonlightShopExchangeCsv, ParseMoonlightShopExchangeCsv, nil}
	DataParserMap["type_activity_moonlight_goods"] = TDataParser{InitMoonlightGoodsCsv, ParseMoonlightGoodsCsv, nil}
	DataParserMap["type_activity_moonlight_award"] = TDataParser{InitMoonlightShopAwardCsv, ParseMoonlightShopAwardCsv, nil}

	//阵营战
	DataParserMap["type_crystal"] = TDataParser{InitCrystalParser, ParseCrystalRecord, nil}
	DataParserMap["type_revive"] = TDataParser{InitReviveParser, ParseReviveRecord, nil}
	DataParserMap["type_campbat_rank"] = TDataParser{InitCampBatRankParser, ParseCampBatRankRecord, nil}
	DataParserMap["type_campbat_store"] = TDataParser{InitCampBatStoreParser, ParseCampBatStoreRecord, nil}

	//时装
	DataParserMap["type_fashion"] = TDataParser{InitFashionParser, ParseFashionRecord, nil}
	DataParserMap["type_fashion_map"] = TDataParser{InitFashionMapParser, ParseFashionMapRecord, nil}
	DataParserMap["type_fashion_strength"] = TDataParser{InitFashionStrengthParser, ParseFashionStrengthRecord, nil}

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

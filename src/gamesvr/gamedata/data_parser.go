package gamedata

//InitDataParser 初始化CSV文件解析器
var NullParser = TDataParser{nil, nil, nil}
var G_DataParserMap = map[string]TDataParser{
	//通用选项配置表
	"type_option": {InitOptionParser, ParseOptionRecord, nil},

	//副本配置表
	"type_copy_base":   {InitCopyParser, ParseCopyRecord, nil},
	"type_copy_main":   {InitMainParser, ParseMainRecord, nil},
	"type_copy_elite":  {InitEliteParser, ParseEliteRecord, nil},
	"type_copy_daily":  {InitDailyParse, ParseDailyRecord, nil},
	"type_copy_famous": {InitFamousParser, ParseFamousRecord, nil},

	//道具物品配置表
	"type_item": {InitItemParser, ParseItemRecord, nil},

	//商城配置表
	"type_mall": {InitMallParser, ParseMallRecord, nil},

	//VIP配置表
	"type_vip":      {InitVipParser, ParseVipRecord, nil},
	"type_vip_week": {InitVipWeekParser, ParseVipWeekRecord, nil},

	//任务成就配置表
	"type_task":                 {InitTaskParser, ParseTaskRecord, nil},
	"type_achievement":          {InitAchievementParser, ParseAchievementRecord, nil},
	"type_seven_activity":       {InitSevenActivityParser, ParseSevenActivityRecord, nil},
	"type_seven_activity_store": {InitSevenActivityStoreRecord, ParseSevenActivityStoreRecord, nil},
	"type_task_award":           {InitTaskAwardParser, ParseTaskAwardRecord, nil},
	"type_tasktype":             {InitTaskTypeParser, ParseTaskTypeRecord, nil},

	//奖励配置表
	"type_award": {InitAwardParser, ParseAwardRecord, nil},

	//签到配置表
	"type_sign":      {InitSignParser, ParseSignRecord, nil},
	"type_sign_plus": {InitSignPlusParser, ParseSignPlusRecord, nil},

	//神将商店
	"type_store": {InitStoreParse, ParseStoreRecord, nil},

	//三国志配置表
	"type_sanguozhi": {InitSanGuoZhiParser, ParseSanGuoZhiRecord, nil},

	//招贤配置表
	"type_summon_config": {InitSummonConfigParser, ParseSummonConfigRecord, nil},

	//竞技场配置表
	"type_arena":       {InitArenaParser, ParseArenaRecord, nil},
	"type_arena_rank":  {InitArenaRankParser, ParseArenaRankRecord, nil},
	"type_arena_store": {InitArenaStoreParser, ParseArenaStoreRecord, nil},
	"type_arena_money": {InitArenaMoneyParser, ParseArenaMoneyRecord, nil},

	//夺宝配置表
	"type_rob":              {InitRobParser, ParseRobRecord, nil},
	"type_treasure_melting": {InitTreasureMeltingParser, ParseTreasureMeltingRecord, nil},

	//三国无双配置表
	"type_sgws_chapter":       {InitSangokuMusouChapter, ParseSangokuMusouChapterRecord, nil},
	"type_sgws_chapter_attr":  {InitSangokuMusouAttrMarkupParser, ParseSangokuMusouAttrMarkupRecord, nil},
	"type_sgws_chapter_award": {InitSangokuMusouChapterAwardParser, ParseSangokuMusouChapterAwardRecord, nil},
	"type_sgws_elite_copy":    {InitSangokuMusouEliteCopyParser, ParseSangokuMusouEliteCopyRecord, nil},
	"type_sgws_store":         {InitSangokuMusouStoreParser, ParseSangokuMusouStoreRecord, nil},
	"type_sgws_treasure":      {InitSangokuMusouSaleParser, ParseSangokuMusouSaleRecord, nil},

	//领地攻伐表
	"type_territory":          {InitTerritoryParser, ParseTerritoryRecord, nil},
	"type_territoryaward":     {InitTerritoryAwardParser, ParseTerritoryAwardRecord, nil},
	"type_territoryskill":     {InitTerritorySkillParser, ParseTerritorySkillRecord, nil},
	"type_territorypatrol":    {InitTerritoryPatrolParser, ParseTerritoryPatrolRecord, nil},
	"type_territoryawardtype": {InitTerritoryAwardTypeParser, ParseTerritoryAwardTypeRecord, nil},

	//叛军围剿表
	"type_rebel_siege":      {InitRebelSiegeParser, ParseRebelSiegeRecord, nil},
	"type_rebel_action":     {InitRebelActionAwardParser, ParseRebelActionAwardRecord, nil},
	"type_rebel_award":      {InitExploitAwardParser, ParseExploitAwardRecord, nil},
	"type_rebel_store":      {InitExploitStoreParser, ParseExploitStoreRecord, nil},
	"type_rebel_activity":   {InitRebelActivityParser, ParseRebelActivityRecord, nil},
	"type_rebel_rank_award": {InitRebelRankAwardParser, ParseRebelRankAwardRecord, nil},

	//挖矿表
	"type_mining_stone_num":    {InitMiningStoneRandomParser, ParseMiningStoneRecord, nil},
	"type_mining_element":      {InitMiningElementParser, ParseMiningElementRecord, nil},
	"type_mining_event":        {InitMiningEventParser, ParserMiningEventRecord, nil},
	"type_mining_black_market": {InitMiningEventBlackMarketParser, ParserMiningEventBlackMarketRecord, nil},
	"type_mining_question":     NullParser,
	"type_mining_buff":         {InitMiningEventBuffParser, ParserMiningEventBuffRecord, nil},
	"type_mining_treasure":     {InitMiningEventTreasureParser, ParserMiningEventTreasureRecord, nil},
	"type_mining_monster":      {InitMiningEventMonsterPerser, ParseMiningEventMonsterRecord, nil},
	"type_mining_award":        {InitMiningAwardParser, ParserMiningAwardRecord, nil},
	"type_mining_guaji":        {InitMiningGuaJiParser, ParserMiningGuaJiRecord, nil},
	"type_mining_rand":         {InitMiningRandParser, ParserMiningRandRecord, nil},

	//活动配置表
	"type_activity":                   {InitActivityParser, ParseActivityRecord, nil},
	"type_activity_competition":       {InitCompetitionParser, ParseCompetitionRecord, nil},
	"type_activity_action":            {InitRecvActionParser, ParseRecvActionRecord, nil},
	"type_activity_discount":          {InitDiscountSaleParser, ParseDiscountSaleRecord, nil},
	"type_activity_login":             {InitActivityLoginParser, ParseActivityLoginRecord, nil},
	"type_activity_moneygod":          {InitActivityMoneyParser, ParseActivityMoneyRecord, nil},
	"type_activity_recharge":          {InitActivityRechargeParser, ParseActivityRechargeRecord, nil},
	"type_activity_limitdaily":        {InitActivityLimitDailyParser, ParseActivityLimitDailyRecord, nil},
	"type_activity_hunt_map":          {InitHuntTreasureMapParser, ParseHuntTreasureMapRecord, nil},
	"type_activity_hunt_store":        {InitHuntTreasureStoreParser, ParseHuntTreasureStoreRecord, nil},
	"type_activity_hunt_turn":         {InitHuntTreasureAwardParser, ParseHuntTreasureAwardRecord, nil},
	"type_activity_rank":              {InitOperationalRankAwardParser, ParseOperationalRankAwardRecord, nil},
	"type_activity_lucky_wheel":       {InitLuckyWheelParser, ParseLuckyWheelRecord, nil},
	"type_activity_group_purchase":    {InitGroupPurchaseParser, ParseGroupPurchaseRecord, nil},
	"type_activity_group_score":       {InitGroupPurchaseScoreParser, ParseGroupPurchaseScoreRecord, nil},
	"type_activity_festival_task":     {InitFestivalTaskParser, ParseFestivalTaskRecord, nil},
	"type_activity_festival_exchange": {InitFestivalExchangeParser, ParseFestivalExchangeRecord, nil},
	"type_activity_week_award":        {InitActivityWeekAwardParser, ParseActivityWeekAwardRecord, nil},
	"type_activity_level_gift":        {InitActivityLevelGiftParser, ParseActivityLevelGiftRecord, nil},
	"type_activity_month_fund":        {InitActivityMonthFundParser, ParseActivityMonthFundRecord, nil},
	"type_activity_limitsale":         {InitLimitSaleItemParser, ParseLimitSaleItemRecord, nil},
	"type_activity_limitsale_award":   {InitLimitSaleAllAwardParser, ParseLimitSaleAllAwardRecord, nil},

	//公会配置表
	"type_guild_base":            {InitGuildParser, ParseGuildBaseRecord, nil},
	"type_guild_copy":            {InitGuildCopyParser, ParseGuildCopyRecord, nil},
	"type_guild_copy_award":      {InitGuildCopyAwardParser, ParseGuildCopyAwardRecord, nil},
	"type_guild_role":            {InitGuildRoleParser, ParseGuildRoleRecord, nil},
	"type_guild_sacrifice":       {InitGuildSacrificeParser, ParseGuildSacrificeRecord, nil},
	"type_guild_sacrifice_award": {InitGuildSacrificeAwardParser, ParseGuildSacrificeAwardRecord, nil},
	"type_guild_store":           {InitGuildStoreParser, ParseGuildStoreRecord, nil},
	"type_guild_skill":           {InitGuildSkillParser, ParseGuildSkillRecord, nil},
	"type_guild_skill_max":       {InitGuildSkillLimitParser, ParseGuildSkillLimitRecord, nil},
	"type_guild_skill_level":     {InitGuildSkillMaxParser, ParseGuildSkillMaxRecord, nil},

	//装备配置表
	"type_equipment":           {InitEquipParser, ParseEquipRecord, nil},
	"type_equip_star":          {InitEquipStarParser, ParseEquipStarRecord, nil},
	"type_equip_refine_cost":   {InitEquipRefineCostParser, ParseEquipRefineCostRecord, nil},
	"type_equip_strength_cost": {InitEquipStrengthCostParser, ParseEquipStrengthCostRecord, nil},
	"type_equipsuit":           {InitEquipSuitParser, ParseEquipSuitRecord, nil},
	"type_shenbin":             {InitShenBinParser, ParseShenBinRecord, FinishShenBinParser},
	"type_shenbin_skill":       {InitShenBinSkillParser, ParseShenBinSkillRecord, nil},

	//八卦镜
	"type_baguajing": {InitBaGuaJingParser, ParseBaGuaJingRecord, nil},

	//称号表
	"type_title": {InitTitleParser, ParseTitleRecord, nil},

	//夺粮战
	"type_foodwar_rank":  {InitFoodWarRankAwardParser, ParseFoodWarRankAwardRecord, nil},
	"type_foodwar_award": {InitFoodWarAwardParser, ParseFoodWarAwardRecord, nil},

	//开服基金
	"type_open_fund": {InitOpenFundParser, ParseOpenFundRecord, nil},

	//将灵表
	"type_herosouls_link":    {InitHeroSoulsParser, ParseHeroSoulsRecord, nil},
	"type_herosouls_map":     {InitSoulMapParser, ParseSoulMapRecord, nil},
	"type_herosouls_store":   {InitHeroSoulsStoreParser, ParseHeroSoulsStoreRecrod, nil},
	"type_herosouls_trials":  {InitHeroSoulsTrialsParser, ParseHeroSoulsTrialRecord, nil},
	"type_herosouls_chapter": {InitHeroSoulsChapterParser, ParseHeroSoulsChapterRecord, nil},

	//宝物配置表
	"type_gem":               {InitGemParser, ParseGemRecord, nil},
	"type_gem_refine_cost":   {InitGemRefineCostParser, ParseGemRefineCostRecord, nil},
	"type_gem_strength_cost": {InitGemStrengthCostParser, ParseGemStrengthCostRecord, nil},

	"type_refine":   {InitRefineParser, ParseRefineRecord, nil},
	"type_strength": {InitStrengthParser, ParseStrengthRecord, nil},

	//强化大师表
	"type_master": {InitMasterParser, ParseMasterRecord, nil},

	//行动力货币配置表
	"type_action":        {InitActionParser, ParseActionRecord, nil},
	"type_money":         {InitMoneyParser, ParseMoneyRecord, nil},
	"type_property_type": {InitPropertyParser, ParsePropertyRecord, nil},

	//功能配置表
	"type_func_open":     {InitFuncOpenParser, ParseFuncOpenRecord, nil},
	"type_reset_cost":    {InitFuncCostParser, ParseFuncCostRecord, nil},
	"type_vip_privilege": {InitVipPrivilegeParser, ParseVipPrivilegeRecord, nil},

	//挂机表
	"type_hangup": {InitHangUpParser, ParseHangUpRecord, nil},

	//英雄配制表
	"type_hero":              {InitHeroParser, ParseHeroRecord, nil},
	"type_hero_level":        {InitHeroLevelParser, ParseHeroLevelRecord, nil},
	"type_hero_relation":     {InitHeroRelationParser, ParseHeroRelationRecord, nil},
	"type_hero_break":        {InitHeroBreakParser, ParseHeroBreakRecord, nil},
	"type_hero_break_talent": {InitHeroBreakTalentParser, ParseHeroBreakTalentRecord, nil},
	"type_hero_culture_max":  {InitCultureMaxParser, ParseCultureMaxRecord, nil},
	"type_hero_destiny":      {InitHeroDestinyParser, ParseHeroDestinyRecord, nil},
	"type_hero_talent":       {InitTalentParser, ParseTalentRecord, nil},
	"type_relation":          {InitHeroRelationBuffParser, ParseHeroRelationBuffRecord, nil},
	"type_hero_diaowen":      {InitDiaoWenParser, ParseDiaoWenRecord, nil},
	"type_hero_xilian":       {InitXiLianParser, ParseXiLianRecord, nil},
	"type_hero_god":          {InitHeroGodParser, ParseHeroGodRecord, nil},
	"type_hero_friend":       {InitHeroFriendParser, ParseHeroFriendRecord, nil},

	//黑市表
	"type_black_market": {InitBlackMarketParser, ParseBlackMarketRecord, nil},

	//机器人表
	"type_robot": {InitRobotParser, ParseRobotRecord, nil},

	//充值表
	"type_monthcard": {InitMonthCardParser, ParseMonthCardRecord, nil},
	"type_charge":    {InitChargeItemParser, ParseChargeItemRecord, nil},

	//觉醒表
	"type_hero_wake":         {InitWakeLevelParser, ParseWakeLevelRecord, nil},
	"type_hero_wake_compose": {InitWakeComposeParser, ParseWakeComposeRecord, nil},

	//积分赛
	"type_jifen_duan":  {InitScoreDwParser, ParseScoreDwRecord, nil},
	"type_jifen_award": {InitScoreAwardParser, ParseScoreAwardRecord, nil},
	"type_jifen_store": {InitScoreStoreParser, ParseScoreStoreRecord, nil},

	//宠物表
	"type_pet":       {InitPetParser, ParsePetRecord, nil},
	"type_pet_level": {InitPetLevelParser, ParsePetLevelRecord, nil},
	"type_pet_god":   {InitPetGodParser, ParsePetGodRecord, nil},
	"type_pet_star":  {InitPetStarParser, ParsePetStarRecord, nil},
	"type_pet_map":   {InitPetMapParser, ParsePetMapRecord, nil},

	//卡牌大师
	"type_activity_card_exchange": {InitCMExchangeItemParser, ParseCMExchangeItemRecord, nil},
	"type_activity_card":          {InitCardCsvParser, ParseCardCsvRecord, nil},

	//月光集市
	"type_activity_moonlight_exch":  {InitMoonlightShopExchangeCsv, ParseMoonlightShopExchangeCsv, nil},
	"type_activity_moonlight_goods": {InitMoonlightGoodsCsv, ParseMoonlightGoodsCsv, nil},
	"type_activity_moonlight_award": {InitMoonlightShopAwardCsv, ParseMoonlightShopAwardCsv, nil},

	//阵营战
	"type_crystal":       {InitCrystalParser, ParseCrystalRecord, nil},
	"type_revive":        {InitReviveParser, ParseReviveRecord, nil},
	"type_campbat_rank":  {InitCampBatRankParser, ParseCampBatRankRecord, nil},
	"type_campbat_store": {InitCampBatStoreParser, ParseCampBatStoreRecord, nil},

	//时装
	"type_fashion":          {InitFashionParser, ParseFashionRecord, nil},
	"type_fashion_map":      {InitFashionMapParser, ParseFashionMapRecord, nil},
	"type_fashion_strength": {InitFashionStrengthParser, ParseFashionStrengthRecord, nil},

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

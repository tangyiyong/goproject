package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"mongodb"
	"utility"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	G_LevelRanker        utility.TRanker //等级排行榜
	G_FightRanker        utility.TRanker //战力排行榜
	G_RebelExploitRanker utility.TRanker //叛军功勋排行榜
	G_RebelDamageRanker  utility.TRanker //叛军伤害排行榜
	G_SgwsStarRanker     utility.TRanker //三国无双排行榜
	G_GuildLevelRanker   utility.TRanker //公会等级排行榜
	G_GuildCopyRanker    utility.TRanker //公会副本排行榜
	G_FoodWarRanker      utility.TRanker //夺粮战排行榜
	G_WanderRanker       utility.TRanker //云游排行榜
	G_HeroSoulsRanker    utility.TRanker //将灵排行榜

	G_HuntTreasureTodayRanker     utility.TRanker //巡回探宝今日排行榜
	G_HuntTreasureYesterdayRanker utility.TRanker //巡回探宝昨日排行榜
	G_HuntTreasureTotalRanker     utility.TRanker //巡回探宝累计排行榜

	G_LuckyWheelTodayRanker     utility.TRanker //幸运转盘今日排行榜
	G_LuckyWheelYesterdayRanker utility.TRanker //幸运转盘昨日排行榜
	G_LuckyWheelTotalRanker     utility.TRanker //幸运转盘累计排行榜

	G_CardMasterTodayRanker     utility.TRanker //卡牌大师积分榜
	G_CardMasterYesterdayRanker utility.TRanker //昨日
	G_CardMasterTotalRanker     utility.TRanker //累计

	G_BeachBabyTodayRanker     utility.TRanker
	G_BeachBabyYesterdayRanker utility.TRanker
	G_BeachBabyTotalRanker     utility.TRanker

	//阵营战
	G_CampBat_TodayKill    utility.TRanker    //今日击杀
	G_CampBat_TodayDestroy utility.TRanker    //今日击杀
	G_CampBat_CampKill     [3]utility.TRanker //本阵营今日击杀排行榜
	G_CampBat_CampDestroy  [3]utility.TRanker //本阵营今日团灭排行榜
	G_CampBat_KillSum      utility.TRanker    //开服以来的总击杀排行榜
	G_CampBat_DestroySum   utility.TRanker    //开服以来的总团灭排行榜
)

const (
	RT_TodayKill    = 1 //今日击杀
	RT_TodayDestroy = 2 //今日击杀
	RT_CampKill     = 3 //本阵营今日击杀排行榜
	RT_CampDestroy  = 4 //本阵营今日团灭排行榜
	RT_KillSum      = 5 //开服以来的总击杀排行榜
	RT_DestroySum   = 6 //开服以来的总团灭排行榜
)

func InitRankMgr() {
	//等级排行榜
	InitLevelRanker()

	//战力排行榜
	InitFightRanker()

	//叛军排行榜
	InitRebelExploitRanker()
	InitRebelDamageRanker()

	//三国无双排行榜
	InitSgwsRanker()

	//公会等级排行榜
	InitGuildLevelRanker()

	//公会副本排行榜
	InitGuildCopyRanker()

	//夺粮战排行榜
	InitFoodWarRanker()

	//云游排行榜
	InitWanderRanker()

	//将灵排行榜
	InitHeroSoulsRanker()

	//巡回探宝排行榜
	InitHuntTreasureRanker()

	//幸运转盘
	InitLuckyWheelRanker()

	//卡牌大师
	InitCardMasterRanker()

	//沙滩宝贝
	InitBeachBabyRanker()

	//阵营战排行榜
	InitCampBattleRanker()
}

//公会副本排行榜
func InitGuildCopyRanker() bool {
	G_GuildCopyRanker.InitRanker(10, 50)

	s := mongodb.GetDBSession()
	defer s.Close()

	var guilds []TGuild
	err := s.DB(appconfig.GameDbName).C("Guild").Find(nil).Sort("-historypasschapter").Limit(50).All(&guilds)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitGuildCopyRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(guilds); i++ {
		G_GuildCopyRanker.SetRankItem(guilds[i].GuildID, int(guilds[i].HistoryPassChapter))
	}

	return true
}

//公会等级排行榜
func InitGuildLevelRanker() bool {
	G_GuildLevelRanker.InitRanker(10, 50)

	s := mongodb.GetDBSession()
	defer s.Close()

	var guilds []TGuild
	err := s.DB(appconfig.GameDbName).C("Guild").Find(nil).Sort("-level").Limit(50).All(&guilds)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitGuildLevelRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(guilds); i++ {
		G_GuildLevelRanker.SetRankItem(guilds[i].GuildID, guilds[i].Level)
	}

	return true
}

//等级排行榜
func InitLevelRanker() bool {
	G_LevelRanker.InitRanker(50, 200)

	s := mongodb.GetDBSession()
	defer s.Close()

	var simplevec []TSimpleInfo
	err := s.DB(appconfig.GameDbName).C("PlayerSimple").Find(nil).Sort("-level").Limit(200).All(&simplevec)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitLevelRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(simplevec); i++ {
		G_LevelRanker.SetRankItem(simplevec[i].PlayerID, simplevec[i].Level)
	}

	return true
}

//战力排行榜
func InitFightRanker() bool {
	G_FightRanker.InitRanker(50, 200)

	s := mongodb.GetDBSession()
	defer s.Close()

	var simplevec []TSimpleInfo
	err := s.DB(appconfig.GameDbName).C("PlayerSimple").Find(nil).Sort("-fightvalue").Limit(200).All(&simplevec)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitFightRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(simplevec); i++ {
		G_FightRanker.SetRankItem(simplevec[i].PlayerID, int(simplevec[i].FightValue))
	}

	return true
}

//叛军功勋排行榜
func InitRebelExploitRanker() bool {
	G_RebelExploitRanker.InitRanker(5, 200)
	rankLst := []TRebelModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerRebel", "exploit", -1, 200, &rankLst)

	for _, v := range rankLst {
		G_RebelExploitRanker.SetRankItem(v.PlayerID, v.Exploit)
	}

	return true
}

//叛军伤害排行榜
func InitRebelDamageRanker() bool {
	G_RebelDamageRanker.InitRanker(5, 200)
	rankLst := []TRebelModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerRebel", "damage", -1, 200, &rankLst)

	for _, v := range rankLst {
		G_RebelDamageRanker.SetRankItem(v.PlayerID, int(v.Damage))
	}

	return true
}

//三国无双排行榜
func InitSgwsRanker() bool {
	G_SgwsStarRanker.InitRanker(5, 200)
	playerLst := []TSangokuMusouModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerSangokuMusou", "historystar", -1, 200, &playerLst)
	for _, v := range playerLst {
		G_SgwsStarRanker.SetRankItem(v.PlayerID, v.HistoryStar)
	}

	return true
}

//夺粮战排行榜
func InitFoodWarRanker() bool {
	G_FoodWarRanker.InitRanker(5, 200)
	rankLst := []TFoodWarModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerFoodWar", "totalfood", -1, 200, &rankLst)
	for _, v := range rankLst {
		G_FoodWarRanker.SetRankItem(v.PlayerID, v.TotalFood)
	}

	return true
}

//云游排行榜
func InitWanderRanker() bool {
	G_WanderRanker.InitRanker(10, 50)

	s := mongodb.GetDBSession()
	defer s.Close()

	var wandervec []TWanderModule
	err := s.DB(appconfig.GameDbName).C("PlayerWander").Find(nil).Sort("-maxcopyid").Limit(50).All(&wandervec)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitLevelRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(wandervec); i++ {
		G_WanderRanker.SetRankItem(wandervec[i].PlayerID, wandervec[i].MaxCopyID)
	}

	return true
}

//英灵排行榜
func InitHeroSoulsRanker() bool {
	//! 初始化参数
	G_HeroSoulsRanker.InitRanker(20, 200)

	rankLst := []THeroSoulsModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerHeroSouls", "soulmapvalue", -1, 200, &rankLst)
	for _, v := range rankLst {
		G_HeroSoulsRanker.SetRankItem(v.PlayerID, v.SoulMapValue*10000+len(v.HeroSoulsLink))
	}

	return true
}

//巡回探宝排行榜
func InitHuntTreasureRanker() bool {
	//! 初始化参数
	G_HuntTreasureTotalRanker.InitRanker(20, 50)

	rankLst := []TActivityModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", "hunttreasure.score", -1, 50, &rankLst)
	for _, v := range rankLst {
		G_HuntTreasureTotalRanker.SetRankItem(v.PlayerID, v.HuntTreasure.Score)
	}

	G_HuntTreasureTodayRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}

	filedName := fmt.Sprintf("hunttreasure.todayscore.%d", utility.GetCurDayMod())
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", filedName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_HuntTreasureTodayRanker.SetRankItem(v.PlayerID, v.HuntTreasure.TodayScore[utility.GetCurDayMod()])
	}

	G_HuntTreasureYesterdayRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	filedName = fmt.Sprintf("hunttreasure.todayscore.%d", 1-utility.GetCurDayMod())
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", filedName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_HuntTreasureYesterdayRanker.SetRankItem(v.PlayerID, v.HuntTreasure.TodayScore[1-utility.GetCurDayMod()])
	}

	return true
}

//巡回探宝排行榜
func InitLuckyWheelRanker() bool {
	//! 初始化参数
	G_LuckyWheelTotalRanker.InitRanker(20, 50)

	indexToday := 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
	}

	rankLst := []TActivityModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", "luckywheel.score", -1, 50, &rankLst)
	for _, v := range rankLst {
		G_LuckyWheelTotalRanker.SetRankItem(v.PlayerID, v.LuckyWheel.TotalScore)
	}

	G_LuckyWheelTodayRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	filedName := fmt.Sprintf("luckywheel.today.%d", utility.GetCurDayMod())
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", filedName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_LuckyWheelTodayRanker.SetRankItem(v.PlayerID, v.LuckyWheel.TodayScore[utility.GetCurDayMod()])
	}

	G_LuckyWheelYesterdayRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	filedName = fmt.Sprintf("luckywheel.today.%d", 1-indexToday)
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", filedName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_LuckyWheelYesterdayRanker.SetRankItem(v.PlayerID, v.LuckyWheel.TodayScore[1-indexToday])
	}

	return true
}

func InitCardMasterRanker() bool { // 卡牌大师积分榜
	indexToday, indexYesterday := 0, 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
		indexYesterday = 0
	} else {
		indexToday = 0
		indexYesterday = 1
	}
	// 今日排行
	G_CardMasterTodayRanker.InitRanker(20, 50)
	rankLst := []TActivityModule{}
	FieldName := fmt.Sprintf("cardmaster.jifen.%d", indexToday)
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", FieldName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_CardMasterTodayRanker.SetRankItem(v.PlayerID, v.CardMaster.JiFen[indexToday])
	}
	// 昨日排行
	G_CardMasterYesterdayRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	FieldName = fmt.Sprintf("cardmaster.jifen.%d", indexYesterday)
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", FieldName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_CardMasterYesterdayRanker.SetRankItem(v.PlayerID, v.CardMaster.JiFen[indexYesterday])
	}
	// 累计排好
	G_CardMasterTotalRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", "cardmaster.totaljifen", -1, 50, &rankLst)
	for _, v := range rankLst {
		G_CardMasterTotalRanker.SetRankItem(v.PlayerID, v.CardMaster.TotalJiFen)
	}
	return true
}

func InitBeachBabyRanker() bool {
	indexToday, indexYesterday := 0, 0
	if utility.GetCurDayMod() == 1 {
		indexToday = 1
		indexYesterday = 0
	} else {
		indexToday = 0
		indexYesterday = 1
	}
	// 今日排行
	G_BeachBabyTodayRanker.InitRanker(20, 50)
	rankLst := []TActivityModule{}
	FieldName := fmt.Sprintf("beachbaby.score.%d", indexToday)
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", FieldName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_BeachBabyTodayRanker.SetRankItem(v.PlayerID, v.BeachBaby.Score[indexToday])
	}
	// 昨日排行
	G_BeachBabyYesterdayRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	FieldName = fmt.Sprintf("beachbaby.score.%d", indexYesterday)
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", FieldName, -1, 50, &rankLst)
	for _, v := range rankLst {
		G_BeachBabyYesterdayRanker.SetRankItem(v.PlayerID, v.BeachBaby.Score[indexYesterday])
	}
	// 累计排好
	G_BeachBabyTotalRanker.InitRanker(20, 50)
	rankLst = []TActivityModule{}
	mongodb.Find_Sort(appconfig.GameDbName, "PlayerActivity", "beachbaby.totalscore", -1, 50, &rankLst)
	for _, v := range rankLst {
		G_BeachBabyTotalRanker.SetRankItem(v.PlayerID, v.BeachBaby.TotalScore)
	}
	return true
}

//初始化阵营战排行榜
func InitCampBattleRanker() bool {
	G_CampBat_TodayKill.InitRanker(20, 50)
	G_CampBat_TodayDestroy.InitRanker(20, 50)
	G_CampBat_CampKill[0].InitRanker(20, 50)
	G_CampBat_CampKill[1].InitRanker(20, 50)
	G_CampBat_CampKill[2].InitRanker(20, 50)
	G_CampBat_CampDestroy[0].InitRanker(20, 50)
	G_CampBat_CampDestroy[1].InitRanker(20, 50)
	G_CampBat_CampDestroy[2].InitRanker(20, 50)
	G_CampBat_KillSum.InitRanker(20, 50)
	G_CampBat_DestroySum.InitRanker(20, 50)

	s := mongodb.GetDBSession()
	defer s.Close()

	var result []TCampBattleModule
	err := s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": bson.M{"$gt": 0}}).Sort("-kill").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_TodayKill.SetRankItem(result[i].PlayerID, result[i].Kill)
	}

	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": bson.M{"$gt": 0}}).Sort("-destroy").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_TodayDestroy.SetRankItem(result[i].PlayerID, result[i].Destroy)
	}
	//////////////////////////////////////////
	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": bson.M{"$gt": 0}}).Sort("-killsum").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_KillSum.SetRankItem(result[i].PlayerID, result[i].KillSum)
	}
	//////////////////////////////////////////
	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": bson.M{"$gt": 0}}).Sort("-destroysum").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_DestroySum.SetRankItem(result[i].PlayerID, result[i].KillSum)
	}
	//////////////////////////////////////////
	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": 1}).Sort("-kill").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_CampKill[0].SetRankItem(result[i].PlayerID, result[i].Kill)
	}

	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": 2}).Sort("-kill").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_CampKill[1].SetRankItem(result[i].PlayerID, result[i].Kill)
	}

	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": 3}).Sort("-kill").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_CampKill[2].SetRankItem(result[i].PlayerID, result[i].Kill)
	}

	//////////////////////////////////////////
	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": 1}).Sort("-destroy").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_CampDestroy[0].SetRankItem(result[i].PlayerID, result[i].Destroy)
	}

	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": 2}).Sort("-destroy").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_CampDestroy[1].SetRankItem(result[i].PlayerID, result[i].Destroy)
	}

	err = s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"battlecamp": 3}).Sort("-destroy").Limit(50).All(&result)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitCampBattleRanker DB Error!!!")
		return false
	}

	for i := 0; i < len(result); i++ {
		G_CampBat_CampDestroy[2].SetRankItem(result[i].PlayerID, result[i].Destroy)
	}

	return true
}

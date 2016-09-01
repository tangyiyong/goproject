package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	OneDay = 24 * 60 * 60
)

const (
	Timer_Func_Arena    = 1 //! 竞技场
	Timer_Func_Rebel    = 2 //! 叛军击杀
	Timer_Func_FoodWar  = 4 //! 夺粮战
	Timer_Func_Activity = 6 //! 活动刷新
	Timer_Func_NewDay   = 7 //! 新的一天
)

var G_Timer Timer

type DealFunc func(now int64) bool
type TimerFunc struct {
	FuncID    int
	ResetTime int64
	CDTime    int
	deal      DealFunc
	IsOpen    bool
}

type Timer struct {
	ID      int `bson:"_id"`
	FuncLst []TimerFunc
}

func GetTodayTime() int64 {
	now := time.Now()
	todayTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return todayTime.Unix()
}

func (self *Timer) Init() {
	self.ID = 1
	if self.LoadTimer() == false {
		mongodb.InsertToDB(appconfig.GameDbName, "Timer", self)
	}

	//! 加入各种模块功能
	//! 竞技场
	resetTime := GetTodayTime() + int64(gamedata.ArenaRankAwardCalcTime*60*60)
	self.AddTimeFunc(Timer_Func_Arena, resetTime, 24*60*60, self.ArenaFunc, true)

	//! 围剿叛军
	resetTime = GetTodayTime() + int64(24*60*60)
	self.AddTimeFunc(Timer_Func_Rebel, resetTime, 24*60*60, self.RebelFunc, true)

	//! 粮草战
	resetTime = GetTodayTime() + int64(gamedata.FoodWarEndTime)
	self.AddTimeFunc(Timer_Func_FoodWar, resetTime, 24*60*60, self.FoodWarFunc, true)

	//! 活动
	resetTime = GetTodayTime() + 24*60*60
	self.AddTimeFunc(Timer_Func_Activity, resetTime, 24*60*60, ActivityTimerFunc, true)

	resetTime = GetTodayTime() + 24*60*60
	self.AddTimeFunc(Timer_Func_NewDay, resetTime, 24*60*60, self.OnNewDayFunc, true)

	go self.OnTimer()
}

//! 结算排行奖励
func (self *Timer) ArenaFunc(now int64) bool {
	gamelog.Info("Timer: ArenaFunc")
	//! 获取竞技场排行榜配置
	arenaConfig := gamedata.GetArenaConfig()
	if arenaConfig == nil {
		gamelog.Error("GetArenaConfig Fail.")
		return true
	}

	//! 开始结算奖励
	for i, v := range G_Rank_List {
		if i+1 > arenaConfig.DailyAwardNeedRank {
			break
		}

		if v.IsRobot == true || v.PlayerID == 0 {
			continue
		}

		award := gamedata.GetArenaRankAward(i + 1)
		if award == 0 {
			continue
		}

		value := strconv.Itoa(i + 1)
		SendAwardMail(v.PlayerID, Text_Arean_Win, gamedata.GetItemsFromAwardID(award), []string{value})
	}

	return true
}

func (self *Timer) OnNewDayFunc(nowUnix int64) bool {
	gamelog.Info("Timer: OnNewDayFunc")
	//先发今日击杀
	for i := 0; i < len(G_CampBat_TodayKill.List); i++ {
		if G_CampBat_TodayKill.List[i].RankID <= 0 {
			break
		}

		pRankAward := gamedata.GetCampBatRank(i + 1)
		if pRankAward == nil {
			break
		}

		value := strconv.Itoa(i + 1)
		SendAwardMail(G_CampBat_TodayKill.List[i].RankID, TextCampBatTodayKill, gamedata.GetItemsFromAwardID(pRankAward.TotalKillID), []string{value})
	}

	for i := 0; i < len(G_CampBat_TodayDestroy.List); i++ {
		if G_CampBat_TodayDestroy.List[i].RankID <= 0 {
			break
		}
		pRankAward := gamedata.GetCampBatRank(i + 1)
		if pRankAward == nil {
			break
		}

		value := strconv.Itoa(i + 1)
		SendAwardMail(G_CampBat_TodayDestroy.List[i].RankID, TextCampBatTodayDestroy, gamedata.GetItemsFromAwardID(pRankAward.TotalDestroyID), []string{value})
	}

	G_CampBat_TodayKill.Clear()
	mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerCampBat", nil, bson.M{"$set": bson.M{"kill": 0}})

	G_CampBat_TodayDestroy.Clear()
	mongodb.UpdateToDBAll(appconfig.GameDbName, "PlayerCampBat", nil, bson.M{"$set": bson.M{"destroy": 0}})

	for i := 0; i < len(G_CampBat_CampKill[0].List); i++ {
		pRankAward := gamedata.GetCampBatRank(i + 1)
		if pRankAward == nil {
			break
		}

		if G_CampBat_CampKill[0].List[i].RankID > 0 {
			value := strconv.Itoa(i + 1)
			SendAwardMail(G_CampBat_CampKill[0].List[i].RankID, TextCampBatCampKill, gamedata.GetItemsFromAwardID(pRankAward.CampKillID), []string{value})
		}

		if G_CampBat_CampKill[1].List[i].RankID > 0 {
			value := strconv.Itoa(i + 1)
			SendAwardMail(G_CampBat_CampKill[1].List[i].RankID, TextCampBatCampKill, gamedata.GetItemsFromAwardID(pRankAward.CampKillID), []string{value})
		}

		if G_CampBat_CampKill[2].List[i].RankID > 0 {
			value := strconv.Itoa(i + 1)
			SendAwardMail(G_CampBat_CampKill[2].List[i].RankID, TextCampBatCampKill, gamedata.GetItemsFromAwardID(pRankAward.CampKillID), []string{value})
		}

	}

	G_CampBat_CampKill[0].Clear()
	G_CampBat_CampKill[1].Clear()
	G_CampBat_CampKill[2].Clear()

	for i := 0; i < len(G_CampBat_CampDestroy[0].List); i++ {
		pRankAward := gamedata.GetCampBatRank(i + 1)
		if pRankAward == nil {
			break
		}

		if G_CampBat_CampDestroy[0].List[i].RankID > 0 {
			value := strconv.Itoa(i + 1)
			SendAwardMail(G_CampBat_CampDestroy[0].List[i].RankID, TextCampBatCampDestroy, gamedata.GetItemsFromAwardID(pRankAward.CampDestroyID), []string{value})
		}

		if G_CampBat_CampDestroy[1].List[i].RankID > 0 {
			value := strconv.Itoa(i + 1)
			SendAwardMail(G_CampBat_CampDestroy[1].List[i].RankID, TextCampBatCampDestroy, gamedata.GetItemsFromAwardID(pRankAward.CampDestroyID), []string{value})
		}

		if G_CampBat_CampDestroy[2].List[i].RankID > 0 {
			value := strconv.Itoa(i + 1)
			SendAwardMail(G_CampBat_CampDestroy[2].List[i].RankID, TextCampBatCampDestroy, gamedata.GetItemsFromAwardID(pRankAward.CampDestroyID), []string{value})
		}

	}

	G_CampBat_CampDestroy[0].Clear()
	G_CampBat_CampDestroy[1].Clear()
	G_CampBat_CampDestroy[2].Clear()

	return true
}

//! 结算夺粮战排行奖励
func (self *Timer) FoodWarFunc(nowUnix int64) bool {
	gamelog.Info("Timer: FoodWarFunc")
	isOpen := false
	now := time.Now()
	nowSec := now.Hour()*3600 + now.Minute()*60 + now.Second()

	//! 获取星期几
	day := int(now.Weekday())
	if now.Weekday() == time.Sunday {
		day = 7
	}

	//! 判断开启天数
	for _, v := range gamedata.FoodWarOpenDay {
		if v == day {
			isOpen = true
		}
	}

	//! 判断开启时间
	if nowSec <= gamedata.FoodWarEndTime && nowSec >= gamedata.FoodWarOpenTime && isOpen == true {
		isOpen = true
	} else {
		isOpen = false
	}

	if isOpen == false {
		return true
	}

	var awardData TAwardData
	awardData.TextType = Text_FoodWar_Rank
	awardData.Time = time.Now().Unix()
	for i, v := range G_FoodWarRanker.List {
		if v.RankID <= 0 {
			break
		}
		award := gamedata.GetFoodWarRankAward(i + 1)
		if award == 0 {
			continue
		}

		awardData.ItemLst = gamedata.GetItemsFromAwardID(award)
		awardData.Value = []string{strconv.Itoa(i + 1)}
		SendAwardToPlayer(v.RankID, &awardData)
	}

	return true
}

//! 结算排行奖励
func (self *Timer) RebelFunc(now int64) bool {
	gamelog.Info("Timer: RebelFunc")
	//! 开始结算奖励
	var awardData TAwardData
	awardData.TextType = Text_Rebel_Damage
	awardData.Time = time.Now().Unix()
	for i, v := range G_RebelExploitRanker.List {
		if v.RankID == 0 {
			continue
		}

		award := gamedata.GetRebelRankAward(i+1, Rank_Exploit)
		awardData.ItemLst = gamedata.GetItemsFromAwardID(award)
		awardData.Value = []string{strconv.Itoa(i + 1)}
		SendAwardToPlayer(v.RankID, &awardData)
	}

	awardData.TextType = Text_Rebel_Exploit
	for i, v := range G_RebelDamageRanker.List {
		if v.RankID == 0 {
			continue
		}

		award := gamedata.GetRebelRankAward(i+1, Rank_Damage)
		awardData.ItemLst = gamedata.GetItemsFromAwardID(award)
		awardData.Value = []string{strconv.Itoa(i + 1)}
		SendAwardToPlayer(v.RankID, &awardData)
	}

	return true
}

func (self *Timer) AddTimeFunc(funcID int, resetTime int64, cdTime int, deal func(int64) bool, isOpen bool) {
	for i, v := range self.FuncLst {
		if v.FuncID == funcID {
			//! 已有对应功能
			self.FuncLst[i].deal = deal
			return
		}
	}

	funcData := TimerFunc{funcID, resetTime, cdTime, deal, isOpen}
	self.FuncLst = append(self.FuncLst, funcData)
	self.SaveTimer()
}

//! 计时器
func (self *Timer) OnTimer() {
	timer := time.NewTimer(time.Second)

	for {
		select {
		case <-timer.C:
			now := time.Now().Unix()
			self.OnTimerFunc(now)
			timer.Reset(time.Second)
		}
	}
}

//! 计时器处理函数
func (self *Timer) OnTimerFunc(now int64) {
	for i, v := range self.FuncLst {
		if now >= self.FuncLst[i].ResetTime && self.FuncLst[i].IsOpen == true {
			if self.FuncLst[i].deal == nil {
				gamelog.Error("OnTimerFunc Error: timerid:%d has not deal function", self.FuncLst[i].FuncID)
				continue
			}

			isLoop := self.FuncLst[i].deal(now)
			if isLoop == false {
				self.FuncLst[i].IsOpen = false
			}

			self.FuncLst[i].ResetTime += int64(v.CDTime)
			self.SaveTimer()
		}
	}
}

//! 存储到数据库
func (self *Timer) SaveTimer() {
	mongodb.UpdateToDB(appconfig.GameDbName, "Timer", bson.M{"_id": 1}, bson.M{"$set": self})
}

//! 读取到内存
func (self *Timer) LoadTimer() bool {
	if mongodb.Find(appconfig.GameDbName, "Timer", "_id", 1, self) != 0 {
		return false
	}
	return true
}

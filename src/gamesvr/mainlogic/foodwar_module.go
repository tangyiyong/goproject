package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TRevengeInfo struct {
	PlayerID int
	RobFood  int
}

//! 夺粮战
type TFoodWarModule struct {
	PlayerID int `bson:"_id"`

	FixedFood int //! 固定粮草
	TotalFood int //! 总计粮草

	AttackTimes  int //! 攻打次数
	RevengeTimes int //! 复仇次数

	BuyAttackTimes  int //! 已购买攻击次数
	BuyRevengeTimes int //! 已购买复仇次数

	RevengeLst []TRevengeInfo //! 复仇名单

	NextTime int64

	AwardRecvLst IntLst //! 粮草奖励领取记录

	ResetDay int

	ownplayer *TPlayer
}

func (self *TFoodWarModule) SetPlayerPtr(playerid int, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

func (self *TFoodWarModule) OnCreate(playerID int) {

	self.ResetDay = utility.GetCurDay()

	self.AttackTimes = gamedata.FoodWarAttackTimes
	self.RevengeTimes = gamedata.FoodWarRevengeTimes
	now := time.Now()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local).Unix()
	nextTime += 3600
	self.NextTime = nextTime
	self.ResetDay = utility.GetCurDay()

	self.TotalFood = gamedata.FoodWarFixedFood + gamedata.FoodWarNonFixedFood
	self.FixedFood = gamedata.FoodWarFixedFood
	self.RevengeLst = []TRevengeInfo{}
	self.BuyAttackTimes = 0
	self.BuyRevengeTimes = 0
	self.AwardRecvLst = IntLst{}

	//! 加入粮草排行榜
	G_FoodWarRanker.SetRankItem(self.PlayerID, self.TotalFood)

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerFoodWar", self)
}

func (self *TFoodWarModule) OnDestroy(playerID int) {

}

func (self *TFoodWarModule) OnPlayerOnline(playerID int) {

}

//! 玩家离开游戏
func (self *TFoodWarModule) OnPlayerOffline(playerID int) {

}

//! 读取玩家
func (self *TFoodWarModule) OnPlayerLoad(playerid int, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerFoodWar").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerFoodWar Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 判断是否活动开启时间
func (self *TFoodWarModule) IsActivityOpen() bool {
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

	return isOpen
}

//! 重置玩家信息
func (self *TFoodWarModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TFoodWarModule) OnNewDay(newday int) {

	self.AttackTimes = gamedata.FoodWarAttackTimes
	self.RevengeTimes = gamedata.FoodWarRevengeTimes
	now := time.Now()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local).Unix()
	nextTime += 3600
	self.NextTime = nextTime
	self.ResetDay = newday

	self.TotalFood = gamedata.FoodWarFixedFood + gamedata.FoodWarNonFixedFood
	self.FixedFood = gamedata.FoodWarFixedFood
	self.RevengeLst = []TRevengeInfo{}
	self.BuyAttackTimes = 0
	self.BuyRevengeTimes = 0
	self.AwardRecvLst = IntLst{}

	//! 加入粮草排行榜
	G_FoodWarRanker.SetRankItem(self.PlayerID, self.TotalFood)
	go self.DB_Reset()
}

//! 获取复仇信息
func (self *TFoodWarModule) GetRevengeInfo(targetPlayerID int) *TRevengeInfo {
	for i, v := range self.RevengeLst {
		if v.PlayerID == targetPlayerID {
			return &self.RevengeLst[i]
		}
	}

	return nil
}

//! 检测时间增长
func (self *TFoodWarModule) CheckTime() {
	//! 判断活动是否开启
	if self.IsActivityOpen() == false {
		return
	}

	now := time.Now()

	interval := now.Unix() - self.NextTime
	addTimes := 0
	if interval < 0 {
		//! 未到下次增加时间
		return
	} else {
		addTimes = 1
	}

	//! 计算累加次数
	addTimes += int(interval / 3600)

	self.AttackTimes += addTimes
	if self.AttackTimes >= gamedata.FoodWarAttackTimes {
		self.AttackTimes = gamedata.FoodWarAttackTimes
	}

	self.TotalFood += gamedata.FoodWarTimeAddFood * addTimes
	self.FixedFood += gamedata.FoodWarTimeAddFood * addTimes

	//! 更新粮草排行榜
	G_FoodWarRanker.SetRankItem(self.PlayerID, self.TotalFood)

	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local).Unix()
	nextTime += 3600
	self.NextTime = nextTime

	go self.DB_CheckTime()
}

//! 获取玩家粮草信息
func (self *TFoodWarModule) GetPlayerFoodInfo(playerID int) *TFoodWarModule {
	s := mongodb.GetDBSession()
	defer s.Close()

	var food TFoodWarModule

	err := s.DB(appconfig.GameDbName).C("PlayerFoodWar").Find(bson.M{"_id": playerID}).One(&food)
	if err != nil {
		gamelog.Error("PlayerFoodWar Load Error :%s， PlayerID: %d", err.Error(), playerID)
		return nil
	}

	return &food
}
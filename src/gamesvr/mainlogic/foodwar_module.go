package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"sync"
	"time"
	"utility"
)

type TRevengeInfo struct {
	PlayerID int32
	RobFood  int
}

//! 夺粮战
type TFoodWarModule struct {
	PlayerID     int32          `bson:"_id"`
	FixedFood    int            //! 固定粮草
	TotalFood    int            //! 总计粮草
	AttackTimes  int            //! 攻打次数
	RevengeTimes int            //! 复仇次数
	BuyTimes     int            //! 已购买攻击次数
	BuyRevTimes  int            //! 已购买复仇次数
	RevengeLst   []TRevengeInfo //! 复仇名单
	NextTime     uint32
	AwardRecvLst IntLst //! 粮草奖励领取记录
	ResetDay     uint32
	ownplayer    *TPlayer
}

func (self *TFoodWarModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TFoodWarModule) OnCreate(playerid int32) {

	self.ResetDay = utility.GetCurDay()

	self.AttackTimes = gamedata.FoodWarAttackTimes
	self.RevengeTimes = gamedata.FoodWarRevengeTimes
	now := time.Now()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local).Unix()
	nextTime += 3600
	self.NextTime = uint32(nextTime)
	self.ResetDay = utility.GetCurDay()

	self.TotalFood = gamedata.FoodWarFixedFood + gamedata.FoodWarNonFixedFood
	self.FixedFood = gamedata.FoodWarFixedFood
	self.RevengeLst = []TRevengeInfo{}
	self.BuyTimes = 0
	self.BuyRevTimes = 0
	self.AwardRecvLst = IntLst{}

	//! 加入粮草排行榜
	G_FoodWarRanker.SetRankItem(self.PlayerID, self.TotalFood)

	//! 插入数据库
	mongodb.InsertToDB("PlayerFoodWar", self)
}

func (self *TFoodWarModule) OnDestroy(playerid int32) {

}

func (self *TFoodWarModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TFoodWarModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TFoodWarModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerFoodWar").Find(&bson.M{"_id": playerid}).One(self)
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

func (self *TFoodWarModule) OnNewDay(newday uint32) {

	self.AttackTimes = gamedata.FoodWarAttackTimes
	self.RevengeTimes = gamedata.FoodWarRevengeTimes
	now := time.Now()
	nextTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local).Unix()
	nextTime += 3600
	self.NextTime = uint32(nextTime)
	self.ResetDay = newday

	self.TotalFood = gamedata.FoodWarFixedFood + gamedata.FoodWarNonFixedFood
	self.FixedFood = gamedata.FoodWarFixedFood
	self.RevengeLst = []TRevengeInfo{}
	self.BuyTimes = 0
	self.BuyRevTimes = 0
	self.AwardRecvLst = IntLst{}

	//! 加入粮草排行榜
	G_FoodWarRanker.SetRankItem(self.PlayerID, self.TotalFood)
	self.DB_Reset()
}

//! 获取复仇信息
func (self *TFoodWarModule) GetRevengeInfo(targetPlayerID int32) *TRevengeInfo {
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

	interval := uint32(now.Unix()) - self.NextTime
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
	self.NextTime = uint32(nextTime)

	self.DB_CheckTime()
}

//! 获取玩家粮草信息
func (self *TFoodWarModule) GetPlayerFoodInfo(playerid int32) *TFoodWarModule {
	s := mongodb.GetDBSession()
	defer s.Close()

	var food TFoodWarModule

	err := s.DB(appconfig.GameDbName).C("PlayerFoodWar").Find(&bson.M{"_id": playerid}).One(&food)
	if err != nil {
		gamelog.Error("PlayerFoodWar Load Error :%s， PlayerID: %d", err.Error(), playerid)
		return nil
	}

	return &food
}

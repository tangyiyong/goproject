package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TFameHallInfo struct {
	PlayerID   int
	HeroID     int
	CharmValue int
}

// 0 战力  1 等级
var G_FameHallLst [2][6]TFameHallInfo

//! 名人堂
type TFameHallModule struct {
	PlayerID int `bson:"_id"`

	CharmValue  int    //! 魅力值
	FreeTimes   int    //! 免费次数
	ResetDay    int    //! 重置天数
	SendFightID IntLst //! 已送花朵
	SendLevelID IntLst //! 已送花朵
	ownplayer   *TPlayer
}

func (self *TFameHallModule) SetPlayerPtr(playerid int, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

func (self *TFameHallModule) OnCreate(playerID int) {
	//! 初始化各类参数
	self.FreeTimes = gamedata.FameHallFreeTimes
	self.CharmValue = 0
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerFameHall", self)
}

func (self *TFameHallModule) OnDestroy(playerID int) {

}

func (self *TFameHallModule) OnPlayerOnline(playerID int) {

}

//! 玩家离开游戏
func (self *TFameHallModule) OnPlayerOffline(playerID int) {

}

//! 读取玩家
func (self *TFameHallModule) OnPlayerLoad(playerid int, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerFameHall").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerFameHall Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 检测重置
func (self *TFameHallModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TFameHallModule) OnNewDay(newday int) {
	//! 重置参数
	self.SendFightID = IntLst{}
	self.SendLevelID = IntLst{}
	self.ResetDay = newday
	self.FreeTimes = gamedata.FameHallFreeTimes

	go self.DB_Reset()
}

func (self *TFameHallModule) RedTip() bool {
	//! 免费次数
	if self.FreeTimes != 0 {
		return true
	}

	return false
}

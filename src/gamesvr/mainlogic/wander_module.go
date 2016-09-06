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

//! 云游模块
type TWanderModule struct {
	PlayerID   int32  `bson:"_id"`
	MaxCopyID  int    //通过的最大副本ID
	CurCopyID  int    //当前打到的副本ID
	CanBattle  int    //是否可以战斗
	LeftTime   int    //今日重置次数
	SglFreeDay uint32 //己免费日
	SingleFree bool   //单抽免费
	ResetDay   uint32
	ownplayer  *TPlayer
}

func (self *TWanderModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TWanderModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.ResetDay = utility.GetCurDay()
	self.SglFreeDay = self.ResetDay
	self.CurCopyID = 0
	self.LeftTime = gamedata.WanderInitTime
	self.MaxCopyID = 0
	self.CanBattle = 1
	self.SingleFree = true

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerWander", self)
}

func (self *TWanderModule) OnDestroy(playerid int32) {

}

func (self *TWanderModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TWanderModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TWanderModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerWander").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerWander Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TWanderModule) CheckReset() {
	curDay := utility.GetCurDay()

	if curDay != self.SglFreeDay {
		if time.Now().Hour() >= 9 {
			self.SglFreeDay = curDay
			self.SingleFree = true
			self.DB_ResetSingleFreeDay()
		}
	}

	if curDay == self.ResetDay {
		return
	}

	self.OnNewDay(curDay)
}

func (self *TWanderModule) OnNewDay(newday uint32) {
	self.ResetDay = newday
	self.LeftTime = gamedata.WanderInitTime
	self.DB_Reset()
}

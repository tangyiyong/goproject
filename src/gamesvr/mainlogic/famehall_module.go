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
	PlayerID   int32
	HeroID     int
	CharmValue int
}

// 0 战力  1 等级
var G_FameHallLst [2][6]TFameHallInfo

//! 名人堂
type TFameHallModule struct {
	PlayerID int32 `bson:"_id"`

	CharmValue  int      //! 魅力值
	FreeTimes   int      //! 免费次数
	ResetDay    uint32   //! 重置天数
	SendFightID Int32Lst //! 已送花朵
	SendLevelID Int32Lst //! 已送花朵
	ownplayer   *TPlayer
}

func (self *TFameHallModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TFameHallModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.FreeTimes = gamedata.FameHallFreeTimes
	self.CharmValue = 0
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerFameHall", self)
}

func (self *TFameHallModule) OnDestroy(playerid int32) {

}

func (self *TFameHallModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TFameHallModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TFameHallModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
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

//! 刷新名人堂
func (self *TFameHallModule) RefreshFameHallLstFunc() bool {
	if len(G_FameHallLst) != 0 {
		G_FameHallLst = [2][6]TFameHallInfo{}
	}

	if G_FightRanker.List.Len() < 6 || G_LevelRanker.List.Len() < 6 {
		gamelog.Error("RefreshFameHallLstFunc Error: G_FightRanker or G_LevelRanker length is not enough")
		return true
	}

	index := 0
	for i := 0; i < G_FightRanker.List.Len(); i++ {
		isExist := false
		for _, v := range G_FameHallLst[0] {
			if v.PlayerID == G_FightRanker.List[i].RankID {
				isExist = true
				break
			}
		}

		if isExist == false {
			G_FameHallLst[0][index].PlayerID = G_FightRanker.List[i].RankID
			index++
		}

		if index >= 6 {
			break
		}
	}

	index = 0
	for i := 0; i < G_LevelRanker.List.Len(); i++ {
		isExist := false
		for _, v := range G_FameHallLst[1] {
			if v.PlayerID == G_LevelRanker.List[i].RankID {
				isExist = true
				break
			}
		}

		if isExist == false {
			G_FameHallLst[1][index].PlayerID = G_LevelRanker.List[i].RankID
			index++
		}

		if index >= 6 {
			break
		}
	}

	for i, v := range G_FameHallLst {

		for n, m := range v {
			if m.PlayerID == 0 {
				continue
			}

			player := GetPlayerByID(m.PlayerID)
			if player == nil { //! 内存中不存在则读取数据库
				s := mongodb.GetDBSession()
				defer s.Close()
				var fameHall TFameHallModule
				err := s.DB(appconfig.GameDbName).C("PlayerFameHall").Find(bson.M{"_id": m.PlayerID}).One(&fameHall)
				if err != nil {
					gamelog.Error("FameHallRefresh Load Error :%s， PlayerID: %d", err.Error(), m.PlayerID)
					continue
				}
				G_FameHallLst[i][n].CharmValue = fameHall.CharmValue

				var heroModule THeroMoudle
				err = s.DB(appconfig.GameDbName).C("PlayerHero").Find(bson.M{"_id": m.PlayerID}).One(&heroModule)
				if err != nil {
					gamelog.Error("FameHallRefresh Load Error :%s， PlayerID: %d", err.Error(), m.PlayerID)
					continue
				}
				G_FameHallLst[i][n].HeroID = heroModule.CurHeros[0].ID
			} else {
				G_FameHallLst[i][n].CharmValue = player.FameHallModule.CharmValue
				G_FameHallLst[i][n].HeroID = player.HeroMoudle.CurHeros[0].ID
			}
		}

	}

	return true
}

//! 检测重置
func (self *TFameHallModule) CheckReset() {
	self.RefreshFameHallLstFunc()

	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TFameHallModule) OnNewDay(newday uint32) {
	//! 重置参数
	self.SendFightID = Int32Lst{}
	self.SendLevelID = Int32Lst{}
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

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"mongodb"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var G_GemPlayersIndex int

type RobPlayerInfo struct {
	PlayerID int             //! 玩家ID
	Name     string          //! 名字
	Level    int             //! 等级
	HeroID   [BATTLE_NUM]int //! 英雄ID
	IsRobot  int             //! 机器人标记
}

type TRobModule struct {
	PlayerID    int   `bson:"_id"`
	FreeWarTime int64 //! 免战时间
	ownplayer   *TPlayer
}

func (self *TRobModule) SetPlayerPtr(playerID int, player *TPlayer) {
	self.PlayerID = playerID
	self.ownplayer = player
}

func (self *TRobModule) OnCreate(playerID int) {
	//! 初始化信息
	self.FreeWarTime = 0

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerRob", self)
}

func (self *TRobModule) OnDestroy(playerID int) {

}

func (self *TRobModule) OnPlayerOnline(playerID int) {

}

func (self *TRobModule) OnPlayerOffline(playerID int) {

}

func (self *TRobModule) OnPlayerLoad(playerid int, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerRob").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerRob Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 使用物品免战牌接口
func (self *TRobModule) AddFreeWarTime(freeTime int) {
	if self.FreeWarTime == 0 {
		self.FreeWarTime = time.Now().Unix() + int64(freeTime)
	} else {
		//! 若已处于免战时间,则累加时间
		self.FreeWarTime = self.FreeWarTime + int64(freeTime)
	}

	go self.UpdateFreeWarTime()
}

//! 刷新免战时间
func (self *TRobModule) RefreshFreeWarTime() {
	now := time.Now().Unix()
	if now >= self.FreeWarTime { //! 免战时间结束
		self.FreeWarTime = 0
		go self.UpdateFreeWarTime()
	}
}

func (self *TRobModule) GetRobList(itemID int, exclude IntLst) (robPlayerLst []RobPlayerInfo) {
	//! 玩家取2个
	count := 0
	if G_GemPlayersIndex >= len(g_SelectPlayers) {
		G_GemPlayersIndex = 0
	}

	for i := G_GemPlayersIndex; i < len(g_SelectPlayers); i++ {
		//! 检查是否处于免战时间
		if time.Now().Unix() < g_SelectPlayers[i].RobModule.FreeWarTime {
			continue
		}
		if g_SelectPlayers[i].BagMoudle.GetGemPieceCount(itemID) > 0 {
			var info RobPlayerInfo
			info.Name = g_SelectPlayers[i].RoleMoudle.Name
			info.Level = g_SelectPlayers[i].GetLevel()
			info.PlayerID = g_SelectPlayers[i].GetPlayerID()
			info.IsRobot = 0
			for i, b := range g_SelectPlayers[i].HeroMoudle.CurHeros {
				info.HeroID[i] = b.HeroID
			}

			if exclude.IsExist(info.PlayerID) < 0 {
				robPlayerLst = append(robPlayerLst, info)
			}

			if len(robPlayerLst) == 2 {
				G_GemPlayersIndex = i
				break
			}
		}
		count++

		if count == 200 { //! 限制查找人数
			G_GemPlayersIndex = i
			break
		}
	}

	//! 机器人
	robotNum := 5 - len(robPlayerLst)
	for i := 0; i < robotNum; i++ {
		robot := gamedata.RandRobot()
		if robot == nil {
			gamelog.Error("Rand Robot Error: robot is nil")
			return
		}

		var info RobPlayerInfo
		info.Name = robot.Name
		for j := 0; j < BATTLE_NUM; j++ {
			info.HeroID[j] = robot.Heros[j].HeroID
		}

		info.Level = robot.Level
		info.PlayerID = robot.RobotID
		if info.PlayerID == 0 {
			gamelog.Error("Get Robot Error : %v", *robot)
		}

		info.IsRobot = 1
		robPlayerLst = append(robPlayerLst, info)
	}

	return robPlayerLst
}

//! 抢劫NPC
func (self *TRobModule) RobNPC(itemID int) bool {
	//! 获取宝物信息
	itemInfo := gamedata.GetItemInfo(itemID)
	if itemInfo == nil {
		gamelog.Error("GetItemInfo fail. ItemID: %d", itemID)
		return false
	}

	//! 获取配置信息
	robConfig := gamedata.GetRobConfig()
	if robConfig == nil {
		gamelog.Error("GetRobConfig fail.")
		return false
	}

	//! 初始化随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randValue := r.Intn(1000)
	rangevalue := 0
	for i, v := range robConfig.Quality {
		if v == itemInfo.Quality {
			rangevalue = robConfig.RobPro[i]
		}
	}

	if randValue < rangevalue {
		return true
	}

	return false
}

//! 抢劫玩家
func (self *TRobModule) RobPlayer(targetLevel int) bool {
	//! 获取配置信息
	robConfig := gamedata.GetRobConfig()
	if robConfig == nil {
		gamelog.Error("GetRobConfig fail.")
	}

	//! 初始化随机种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randValue := r.Intn(100)

	isLow := false
	levelDifference := self.ownplayer.GetLevel() - targetLevel
	if levelDifference < 0 {
		levelDifference *= -1
		isLow = true
	}

	if levelDifference <= robConfig.GeneralLevelDifference {
		if randValue < robConfig.PlayerGeneralPro {
			return true
		}
	}

	if isLow == true && levelDifference > robConfig.HighLevelDifference {
		if randValue < robConfig.PlayerHighPro {
			return true
		}
	}

	if isLow == false && levelDifference > robConfig.LowLevelDifference {
		if randValue < robConfig.PlayerLowPro {
			return true
		}
	}

	return false
}

func (self *TRobModule) UpdateFreeWarTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRob", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"freewartime": self.FreeWarTime}})
}

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

type TArenaRankData struct {
	PlayerID int32
	Rank     int
	IsRobot  bool
}

type TArenaModule struct {
	PlayerID    int32  `bson:"_id"` //! 唯一标识
	StoreAward  IntLst //! 商店已购买奖励ID
	CurrentRank int    //! 当前玩家排名
	HistoryRank int    //! 历史最高排名
	ownplayer   *TPlayer
}

func (self *TArenaModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TArenaModule) OnCreate(playerid int32) {
	//! 初始化信息
	self.CurrentRank = 5001
	self.HistoryRank = 5001

	//! 判断排名
	if self.CurrentRank <= 5000 {
		G_Rank_List[self.CurrentRank-1].PlayerID = self.PlayerID
		G_Rank_List[self.CurrentRank-1].IsRobot = false
	}

	//! 插入数据库
	mongodb.InsertToDB("PlayerArena", self)
}

func (self *TArenaModule) RedTip() bool {
	length := len(gamedata.GT_ArenaStore_List)
	for i := 0; i < length; i++ {
		if gamedata.GT_ArenaStore_List[i].Type == 2 &&
			self.HistoryRank >= gamedata.GT_ArenaStore_List[i].NeedRank &&
			self.ownplayer.GetLevel() >= gamedata.GT_ArenaStore_List[i].NeedLevel &&
			self.StoreAward.IsExist(gamedata.GT_ArenaStore_List[i].ID) == -1 {
			return true
		}
	}

	return false
}

//! 玩家销毁角色
func (self *TArenaModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TArenaModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TArenaModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TArenaModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerArena").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerArena Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 获取排名玩家信息
func (self *TArenaModule) GetRankPlayer(rank int) *TArenaRankInfo {
	for i, _ := range G_Rank_List {
		if i+1 == rank {
			return &G_Rank_List[i]
		}
	}
	return nil
}

//! 获取排名信息
func (self *TArenaModule) GetRankPlayerInfo() (playerInfo []TArenaRankData) {

	playerRankLst := make([]int, 10)
	if self.CurrentRank > 50 { //! 按照概率抽取
		// {32, 28, 24, 20, 16, 12, 8, 4, -4, -8}
		//  0   1   2   3   4   5   6  7   8   9
		//! 获取排行榜玩家ID
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		currentRank := self.CurrentRank

		if currentRank > 5000 {
			currentRank = 5000
		}

		for i := 0; i < 10; i++ {
			pro := r.Intn(3) + 1 + (28 - 4*i)
			if pro < 0 {
				pro *= -1
				if self.CurrentRank+currentRank*pro/100 <= 5000 {
					//! 不刷新出5000之外的名次
					playerRankLst[i] = currentRank + currentRank*pro/100
				} else {
					if currentRank-i <= 5000 && self.CurrentRank < 5000 {
						playerRankLst[i] = currentRank - (10 - i)
					}
				}
				continue
			}

			playerRankLst[i] = currentRank - currentRank*pro/100
			if playerRankLst[i] == currentRank {
				i -= 1
			}

		}
	} else if self.CurrentRank <= 50 && self.CurrentRank > 10 { //! 按照顺次抽取
		for i := 7; i >= 0; i-- {
			playerRankLst[7-i] = self.CurrentRank - i
		}

		playerRankLst[8] = self.CurrentRank + 1
		playerRankLst[9] = self.CurrentRank + 2
	} else if self.CurrentRank <= 10 { //! 仅显示前十玩家
		return playerInfo
	}

	//! 根据排名获取玩家信息
	for _, v := range playerRankLst {
		if v == 0 {
			continue
		}

		player := self.GetRankPlayer(v)
		if player == nil {
			gamelog.Error("GetRankPlayer error: invalid rank %d", v)
			continue
		}

		var topinfo TArenaRankData
		topinfo.PlayerID = player.PlayerID
		topinfo.Rank = v
		topinfo.IsRobot = player.IsRobot
		playerInfo = append(playerInfo, topinfo)
	}

	return playerInfo
}

//! 刷新可挑战名单
func (self *TArenaModule) RefreshChallangeLst() []TArenaRankData {

	challangeLst := []TArenaRankData{}

	//! 获取前十玩家
	for i := 0; i < 10; i++ {
		var topinfo TArenaRankData
		topinfo.PlayerID = G_Rank_List[i].PlayerID
		topinfo.Rank = i + 1
		topinfo.IsRobot = G_Rank_List[i].IsRobot
		challangeLst = append(challangeLst, topinfo)
	}

	if self.CurrentRank > 10 {
		//! 获取显示可挑战玩家
		challanger := self.GetRankPlayerInfo()
		challangeLst = append(challangeLst, challanger...)
	}

	if self.CurrentRank <= 10 {
		return challangeLst //! 玩家已到前十, 固定显示, 不需要插入
	}

	sortLst := []TArenaRankData{}

	//! 插入自己
	var myRank TArenaRankData
	myRank.PlayerID = self.PlayerID
	myRank.Rank = self.CurrentRank
	myRank.IsRobot = false

	isInside := false
	for _, v := range challangeLst {
		if self.CurrentRank < v.Rank && isInside == false {
			isInside = true
			sortLst = append(sortLst, myRank)
		}

		sortLst = append(sortLst, v)
	}

	if isInside == false {
		//! 不存在比自己小的时候, 插入最后
		sortLst = append(sortLst, myRank)
	}

	return sortLst
}

//! 返回玩家排名
func (self *TArenaModule) GetRankFromID(playerid int32) int {
	for i, v := range G_Rank_List {
		if v.PlayerID == playerid {
			return i + 1
		}
	}

	gamelog.Error("GetRankFromID fail. PlayerID:%d", playerid)
	return 0
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

//角色基本数据表结构
type TVipMoudle struct {
	PlayerID     int32  `bson:"_id"` //主键 玩家ID
	FirstCharges []bool //己完成的首充表

	ownplayer *TPlayer //父player指针
}

func (playerVip *TVipMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	playerVip.PlayerID = playerid
	playerVip.ownplayer = pPlayer
}

func (playerVip *TVipMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	playerVip.PlayerID = playerid

	count := gamedata.GetChargeItemCount()
	playerVip.FirstCharges = make([]bool, count)

	//创建数据库记录
	mongodb.InsertToDB(appconfig.GameDbName, "PlayerVip", playerVip)
}

//玩家对象销毁
func (playerVip *TVipMoudle) OnDestroy(playerid int32) {
	playerVip = nil
}

//玩家进入游戏
func (playerVip *TVipMoudle) OnPlayerOnline(playerid int32) {
	//
}

//OnPlayerOffline 玩家离开游戏
func (playerVip *TVipMoudle) OnPlayerOffline(playerid int32) {
	//
}

//! 预取玩家数据
func (playerVip *TVipMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerVip").Find(bson.M{"_id": playerid}).One(playerVip)
	if err != nil {
		gamelog.Error("PlayerVip Load Error: %s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	playerVip.PlayerID = playerid
}

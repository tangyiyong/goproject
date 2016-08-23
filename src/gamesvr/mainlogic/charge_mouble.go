package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

type TChargeMoudle struct {
	PlayerID    int32 `bson:"_id"` //主键 玩家ID
	ChargeTimes []int

	ownplayer *TPlayer //父player指针
}

//！ 活动框架代码
func (self *TChargeMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}
func (self *TChargeMoudle) OnCreate(playerid int32) {
	self.PlayerID = playerid
	count := gamedata.GetChargeItemCount()
	self.ChargeTimes = make([]int, count)

	//创建数据库记录
	mongodb.InsertToDB(appconfig.GameDbName, "PlayerCharge", self)
}
func (self *TChargeMoudle) OnDestroy(playerid int32) {
	self = nil
}
func (self *TChargeMoudle) OnPlayerOnline(playerid int32) {
	//
}
func (self *TChargeMoudle) OnPlayerOffline(playerid int32) {
	//
}
func (self *TChargeMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerCharge").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerCharge Load Error: %s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 数据操作代码
func (self *TChargeMoudle) AddChargeTimes(id int) int {
	if id <= 0 || id >= len(self.ChargeTimes) {
		gamelog.Error("AddChargeTimes Error : Invalid id :%d", id)
		return 0
	}
	self.ChargeTimes[id]++
	self.db_SaveChargeTimes(id)
	return self.ChargeTimes[id]
}

//! DB相关
func (self *TChargeMoudle) db_SaveChargeTimes(nIndex int) {
	FieldName := fmt.Sprintf("chargetimes.%d", nIndex)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCharge", bson.M{"_id": self.PlayerID},
		bson.M{"$set": bson.M{FieldName: self.ChargeTimes[nIndex]}})
}

//! 逻辑代码
func (self *TChargeMoudle) IsFirstCharge(id int) bool {
	if id <= 0 || id >= len(self.ChargeTimes) {
		return false
	}
	return self.ChargeTimes[id] == 0
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"msg"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

//! 三国志增益信息
type TSanGuoZhiData struct {
	ID       int
	GainType int
	Attr     int
	Value    int
}

type TSanGuoZhiDataLst []TSanGuoZhiData

type TSanGuoZhiModule struct {
	PlayerID  int32 `bson:"_id"`
	CurStarID int
	// SanGuoZhi TSanGuoZhiDataLst
	ownplayer *TPlayer
}

func (self *TSanGuoZhiModule) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

//! 玩家创建角色
func (self *TSanGuoZhiModule) OnCreate(playerid int32) {
	//! 初始化信息
	self.PlayerID = playerid

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerSanGuoZhi", self)
}

//! 玩家销毁角色
func (self *TSanGuoZhiModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TSanGuoZhiModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TSanGuoZhiModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TSanGuoZhiModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerSanGuoZhi").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerSanGuoZhi Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 检测命星材料是否足够
func (self *TSanGuoZhiModule) CheckItemEnough(starID int) (bool, int) {
	info := gamedata.GetSanGuoZhiInfo(starID)
	if info == nil {
		//! 无法获取该星信息
		return false, msg.RE_INVALID_PARAM
	}

	bEnough := self.ownplayer.BagMoudle.IsItemEnough(info.CostType, info.CostNum)
	if !bEnough {
		return false, msg.RE_SANGUOZHI_ITEM_NOT_ENOUGH
	}

	return true, msg.RE_SUCCESS
}

//! 将升星信息存储到数据库
func (self *TSanGuoZhiModule) SaveSanGuoZhiStar() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSanGuoZhi", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"curstarid": self.CurStarID}})
}

func (self *TSanGuoZhiModule) RedTip() bool {
	//! 判断升星材料是否足够
	isEngouth, _ := self.CheckItemEnough(self.CurStarID + 1)
	if isEngouth == true {
		return true
	}

	return false
}

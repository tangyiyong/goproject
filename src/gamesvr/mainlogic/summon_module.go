package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//! 普通召唤
type TNormalSummon struct {
	SummonCounts int   //! 今日召唤次数
	SummonTime   int64 //! 下次可免费召唤时间戳
	ResetTime    int64 //! 重置次数时间
}

//! 高级召唤
type TSeniorSummon struct {
	SummonPoint int   //! 高级召唤积分
	SummonTime  int64 //! 下次可免费召唤时间戳
	OrangeCount int   //! 十次送橙将
}

//! 商城召唤模块
type TSummonModule struct {
	PlayerID int32 `bson:"_id"`

	Normal TNormalSummon //! 普通召唤
	Senior TSeniorSummon //! 高级召唤

	IsFirst bool

	ownplayer *TPlayer
}

func (self *TSummonModule) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

//! 玩家创建角色
func (self *TSummonModule) OnCreate(playerid int32) {
	//! 初始化信息
	self.Normal.SummonTime = time.Now().Unix()
	self.Senior.SummonTime = time.Now().Unix()

	self.IsFirst = true

	self.Normal.ResetTime = GetTodayTime() + 24*60*60

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerSummon", self)
}

//! 玩家销毁角色
func (self *TSummonModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TSummonModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TSummonModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TSummonModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerSummon").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerSummon Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TSummonModule) RedTip() bool {
	now := time.Now().Unix()
	if now >= self.Senior.SummonTime {
		//! 免费高级召唤
		return true
	}

	if now >= self.Normal.SummonTime && self.Normal.SummonCounts < 3 {
		//! 免费普通召唤
		return true
	}

	summonConfig := gamedata.GetSummonConfig(gamedata.Summon_Normal)
	if self.ownplayer.BagMoudle.IsItemEnough(summonConfig.CostItemID, 1) == true {
		return true
	}

	summonConfig = gamedata.GetSummonConfig(gamedata.Summon_Senior)
	if self.ownplayer.BagMoudle.IsItemEnough(summonConfig.CostItemID, 1) == true {
		return true
	}

	return false
}

//! 更新召唤状态
func (self *TSummonModule) UpdateSummonStatus() {
	//! 获取当期时间戳
	now := time.Now().Unix()

	//! 判断重置
	if now > self.Normal.ResetTime {
		self.Normal.SummonCounts = 0
		self.Normal.SummonTime = now
		self.Normal.ResetTime = GetTodayTime() + 24*60*60
		go self.UpdateNormalSummon()
	}
}

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

const (
	CB_KILL    = 1
	CB_DESTROY = 2
)

var (
	G_CampKillNum [3]int //三个阵营今天的击杀人数
)

//! 活动模块
type TCampBattleModule struct {
	PlayerID        int32           `bson:"_id"`
	BattleCamp      int             //阵营战阵营
	Kill            int             //今日杀
	Destroy         int             //今日团
	KillSum         int             //总击杀
	DestroySum      int             //总团灭
	KillHonor       int             //今日击杀荣誉
	LeftTimes       int             //搬运水晶次数
	CrystalID       int             //搬运水晶的ID
	EndTime         int             //搬运结束时间,  超时就是搬运失败
	StoreBuyRecord  []TStoreBuyData //购买商店的次数
	AwardStoreIndex IntLst          //奖励商店的购买ID
	ResetDay        uint32          //重置天

	///////////////以下为临时数据
	enterCode int32    //阵营战的连接进入码
	ownplayer *TPlayer //玩家角色指针
}

func (self *TCampBattleModule) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

func (self *TCampBattleModule) OnCreate(playerid int32) {
	self.BattleCamp = 0
	self.ResetDay = utility.GetCurDay()
	self.CrystalID = 1
	self.LeftTimes = gamedata.CampBat_MoveTimes
	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerCampBat", self)
}

func (self *TCampBattleModule) OnDestroy(playerid int32) {

}

func (self *TCampBattleModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TCampBattleModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TCampBattleModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerCampBat").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerCampBat Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TCampBattleModule) CheckReset() {
	curDay := utility.GetCurDay()
	if curDay == self.ResetDay {
		return
	}

	self.OnNewDay(curDay)
}

func (self *TCampBattleModule) OnNewDay(newday uint32) {
	self.ResetDay = newday
	self.Kill = 0
	self.Destroy = 0
	self.KillHonor = 0
	self.LeftTimes = gamedata.CampBat_MoveTimes
	self.EndTime = 0
	self.DB_Reset()
}

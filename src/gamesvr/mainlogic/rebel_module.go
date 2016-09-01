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

const (
	Rank_Exploit = 1
	Rank_Damage  = 2
)

//! 围剿叛军模块
type TRebelModule struct {
	PlayerID int32 `bson:"_id"`

	RebelID         int    //! 当前叛军ID
	CurLife         int    //! 当前血量
	Level           int    //! 叛军等级
	EscapeTime      int64  //! 逃跑时间
	Damage          int    //! 单次伤害
	Exploit         int    //! 功勋
	ExploitAwardLst IntLst //! 功勋奖励领取标记
	ResetDay        uint32 //! 重置功勋奖励领取标记时间
	IsShare         bool   //! 是否分享

	ownplayer *TPlayer
}

func (self *TRebelModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TRebelModule) OnCreate(playerid int32) {
	//! 初始化信息

	//! 设置重置时间
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerRebel", self)
}

//! 玩家销毁角色
func (self *TRebelModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TRebelModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TRebelModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TRebelModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerRebel").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerRebel Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 重置功勋奖励
func (self *TRebelModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TRebelModule) OnNewDay(newday uint32) {
	self.Exploit = 0
	self.Damage = 0
	self.ExploitAwardLst = IntLst{}
	self.ResetDay = utility.GetCurDay()
	go self.DB_UpdateResetTime()
}

//! 是否发现叛军
func (self *TRebelModule) IsHaveRebel() bool {
	if self.CurLife > 0 && self.RebelID != 0 {
		return true
	}
	return false
}

//! 随机叛军
func (self *TRebelModule) RandRebel() {
	rebelInfo := gamedata.RandRebel(self.ownplayer.GetLevel())
	if rebelInfo == nil {
		gamelog.Error("RandRebel Error : cannot rand rebel!")
		return
	}

	//! 设置属性
	self.RebelID = rebelInfo.CopyID
	self.CurLife = rebelInfo.LifeValue * 10000 //! 需要根据等级加成
	self.Level += 1
	self.EscapeTime = time.Now().Unix() + int64(gamedata.RebelEscapeTime)
	self.IsShare = false

	//! 更新到数据库
	go self.DB_UpdateRebelInfo()
}

//! 检测逃跑时间
func (self *TRebelModule) CheckEscapeTime() {
	if time.Now().Unix() < self.EscapeTime {
		return
	}

	self.RebelID = 0
	self.EscapeTime = 0
	self.IsShare = false
	go self.DB_UpdateRebelInfo()
}

//! 获取当期叛军等级
func (self *TRebelModule) GetRebelLevel() int {
	return self.Level
}

//! 获取玩家叛军信息
func (self *TRebelModule) GetPlayerRebelPtr(playerid int32) (*TRebelModule, string) {
	var rebelModule *TRebelModule
	playerName := ""

	//! 尝试从内存读取玩家信息
	player := GetPlayerByID(playerid)
	if player == nil {
		rebelModule = new(TRebelModule)
		isSuccess := mongodb.Find(appconfig.GameDbName, "PlayerRebel", "_id", playerid, rebelModule)
		if isSuccess != 0 {
			return nil, ""
		}

		roleInfo := new(TRoleMoudle)
		isSuccess = mongodb.Find(appconfig.GameDbName, "PlayerRole", "_id", playerid, roleInfo)
		if isSuccess != 0 {
			return nil, ""
		}

		playerName = roleInfo.Name
	} else {
		rebelModule = &player.RebelModule
		playerName = player.RoleMoudle.Name
	}

	return rebelModule, playerName
}

//! 检测活动开启
func (self *TRebelModule) GetOpenActivity() int {
	now := time.Now()
	nowtime := now.Hour()*60*60 + now.Minute()*60 + now.Second()

	return gamedata.GetRebelOpenActivity(nowtime)
}

func (self *TRebelModule) RedTip() bool {
	if self.CurLife > 0 && self.RebelID != 0 {
		return true
	}

	//! 获取好友发现的叛军
	for _, v := range self.ownplayer.FriendMoudle.FriendList {
		rebelModulePtr, _ := self.GetPlayerRebelPtr(v.PlayerID)
		rebelModulePtr.CheckEscapeTime()
		if rebelModulePtr.RebelID != 0 && rebelModulePtr.IsShare == true {
			return true
		}
	}

	awardLst := gamedata.GetExploitAwardFromLevel(self.ownplayer.GetLevel())
	for _, v := range awardLst {
		if self.ExploitAwardLst.IsExist(v.ID) == -1 && self.Exploit >= v.NeedExploit {
			return true
		}
	}

	return false
}

//! 使用征讨令接口
func (self *TRebelModule) UseItem(itemID int, itemNum int) {
	//! 获取物品信息
	itemInfo := gamedata.GetItemInfo(itemID)
	if itemInfo == nil {
		gamelog.Error("GetItemInfo error: invalid itemID: %d", itemID)
		return
	}

	self.ownplayer.RoleMoudle.AddAction(gamedata.AttackRebelActionID, itemNum)
}

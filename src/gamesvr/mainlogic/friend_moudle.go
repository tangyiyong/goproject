package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TFriendInfo struct {
	PlayerID int  //玩家ID
	IsGive   bool //是否己赠送
	HasAct   bool //是否有未领取
}

type TFriendMoudle struct {
	PlayerID   int `bson:"_id"`
	FriendList []TFriendInfo
	ApplyList  IntLst   //玩家申请列表
	BlackList  IntLst   //黑名单列表
	RcvNum     int      //今天领取体力次数
	ResetDay   int      //刷新时间点
	ownplayer  *TPlayer //父player指针
}

func (self *TFriendMoudle) SetPlayerPtr(playerid int, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

//OnCreate 响应角色创建
func (self *TFriendMoudle) OnCreate(playerid int) {
	self.ResetDay = utility.GetCurDay()
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerFriend", self)
}

//OnDestroy player销毁
func (self *TFriendMoudle) OnDestroy(playerid int) {

}

//OnPlayerOnline player进入游戏
func (self *TFriendMoudle) OnPlayerOnline(playerid int) {

}

//OnPlayerOffline player 离开游戏
func (self *TFriendMoudle) OnPlayerOffline(playerid int) {

}

//OnLoad 从数据库中加载
func (self *TFriendMoudle) OnPlayerLoad(playerid int, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerFriend").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerFriend Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TFriendMoudle) CheckReset() {
	curDay := utility.GetCurDay()
	if curDay == self.ResetDay {
		return
	}

	self.OnNewDay(curDay)
}

func (self *TFriendMoudle) RedTip() bool {
	//! 好友申请
	if len(self.ApplyList) > 0 {
		return true
	}

	//! 检查领取体力
	for _, v := range self.FriendList {
		if v.HasAct == true {
			return true
		}
	}

	return false
}

func (self *TFriendMoudle) OnNewDay(newday int) {
	self.ResetDay = newday
	self.RcvNum = 0
	for i := 0; i < len(self.FriendList); i++ {
		self.FriendList[i].HasAct = false
		self.FriendList[i].IsGive = false
	}

	self.DB_UpdateFriend()
}

//获取好友信息
func (self *TFriendMoudle) GetFriendByID(id int) (*TFriendInfo, int) {
	nCount := len(self.FriendList)
	for i := 0; i < nCount; i++ {
		if self.FriendList[i].PlayerID == id {
			return &self.FriendList[i], i
		}
	}

	return nil, -1
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"msg"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TScorePlayer struct {
	PlayerID   int32  //玩家ID
	Name       string //角色名
	HeroID     int    //英雄ID
	SvrID      int32  //服务器ID
	SvrName    string //服务器名字
	FightValue int32  //战力
	Level      int    //等级
	Quality    int8   //品质
}

//角色基本数据表结构
type TScoreMoudle struct {
	PlayerID   int32             `bson:"_id"` //玩家ID
	FightTime  int               //今天战斗次数
	Score      int               //自己的积分
	ScoreEnemy [3]TScorePlayer   //积分目标
	RecvAward  []int             //己经领取得积分奖励ID
	BuyTime    int               //己购买战斗的次数
	BuyRecord  []msg.MSG_BuyData //购买商店的次数
	SeriesWin  uint32            //连胜次数+领奖状态高16位是领取状态，低16位
	ResetDay   uint32            //重置时间标线
	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	rank      int      //自己的排行(不需要存库)
	ownplayer *TPlayer //父player指针
}

func (self *TScoreMoudle) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TScoreMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	self.PlayerID = playerid
	self.FightTime = 0
	self.Score = 0
	self.RecvAward = make([]int, 0)
	self.ResetDay = utility.GetCurDay()

	//创建数据库记录
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerScore", self)
}

func (self *TScoreMoudle) CheckReset() {
	if utility.IsSameDay(self.ResetDay) {
		return
	}

	self.OnNewDay()
}

func (self *TScoreMoudle) OnNewDay() {
	self.FightTime = 0
	self.BuyTime = 0
	self.ResetDay = utility.GetCurDay()
	self.RecvAward = make([]int, 0)
	self.BuyRecord = []msg.MSG_BuyData{}
	self.SeriesWin = 0
	self.DB_OnNewDay()
}

//玩家对象销毁
func (self *TScoreMoudle) OnDestroy(playerid int32) {
	self = nil
}

//玩家进入游戏
func (self *TScoreMoudle) OnPlayerOnline(playerid int32) {
}

//OnPlayerOffline 玩家离开游戏
func (self *TScoreMoudle) OnPlayerOffline(playerid int32) {
}

//玩家离开游戏
func (self *TScoreMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	var bRet = true
	err := s.DB(appconfig.GameDbName).C("PlayerScore").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerScore Load Error :%s， PlayerID: %d", err.Error(), playerid)
		bRet = false
	}

	if wg != nil {
		wg.Done()
	}

	self.PlayerID = playerid
	return bRet
}

//从跨服服务器取三个目标

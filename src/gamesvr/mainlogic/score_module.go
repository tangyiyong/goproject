package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TScorePlayer struct {
	PlayerID   int32
	Name       string //角色名
	HeroID     int    //英雄ID
	SvrID      int
	SvrName    string
	FightValue int //战力
	Level      int //等级
}

//角色基本数据表结构
type TScoreMoudle struct {
	PlayerID        int32           `bson:"_id"` //玩家ID
	FightTime       int             //今天战斗次数
	Score           int             //自己的积分
	ScoreEnemy      [3]TScorePlayer //积分目标
	RecvAward       []int           //己经领取得积分奖励ID
	ResetDay        uint32          //重置时间标线
	BuyTime         int             //己购买战斗的次数
	StoreBuyRecord  []TStoreBuyData //购买商店的次数
	AwardStoreIndex Int32Lst        //奖励商店的购买ID
	WinTime         int             //连胜次数

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
	rank      int      //自己的排行(不需要存库)
	ownplayer *TPlayer //父player指针
}

func (score *TScoreMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	score.PlayerID = playerid
	score.ownplayer = pPlayer
}

func (score *TScoreMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	score.PlayerID = playerid
	score.FightTime = 0
	score.Score = 0
	score.RecvAward = make([]int, 0)
	score.ResetDay = utility.GetCurDay()

	//创建数据库记录
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerScore", score)
}

func (score *TScoreMoudle) CheckReset() {
	if utility.IsSameDay(score.ResetDay) {
		return
	}

	score.OnNewDay()
}

func (score *TScoreMoudle) OnNewDay() {
	score.FightTime = 0
	score.BuyTime = 0
	score.ResetDay = utility.GetCurDay()
	score.RecvAward = make([]int, 0)
	score.StoreBuyRecord = []TStoreBuyData{}
	score.DB_SaveShoppingInfo()
	score.DB_SaveBuyFightTime()
	score.DB_SaveScoreAndFightTime()
	score.DB_UpdateRecvAward()
}

//玩家对象销毁
func (score *TScoreMoudle) OnDestroy(playerid int32) {
	score = nil
}

//玩家进入游戏
func (score *TScoreMoudle) OnPlayerOnline(playerid int32) {
}

//OnPlayerOffline 玩家离开游戏
func (score *TScoreMoudle) OnPlayerOffline(playerid int32) {
}

//玩家离开游戏
func (score *TScoreMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	var bRet = true
	err := s.DB(appconfig.GameDbName).C("PlayerScore").Find(bson.M{"_id": playerid}).One(score)
	if err != nil {
		gamelog.Error("PlayerScore Load Error :%s， PlayerID: %d", err.Error(), playerid)
		bRet = false
	}

	if wg != nil {
		wg.Done()
	}

	score.PlayerID = playerid
	return bRet
}

//从跨服服务器取三个目标

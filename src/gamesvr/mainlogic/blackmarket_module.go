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

type TBlackMarketGoods struct {
	ID    int
	IsBuy bool
}

//! 黑市
type TBlackMarketModule struct {
	PlayerID int32 `bson:"_id"`

	GoodsLst    []TBlackMarketGoods //! 商品
	RefreshTime int64               //! 刷新时间

	IsOpen      bool   //! 是否开启
	OpenEndTime int64  //! 商店结束时间
	ResetDay    uint32 //! 隔天刷新

	ownplayer *TPlayer
}

func (self *TBlackMarketModule) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = pPlayer
}

func (self *TBlackMarketModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.ResetDay = utility.GetCurDay()
	self.IsOpen = false

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerBlackMarket", self)
}

func (self *TBlackMarketModule) RefreshGoods(isSave bool) {
	if len(self.GoodsLst) > 0 {
		self.GoodsLst = []TBlackMarketGoods{}
	}

	self.RefreshTime = self.GetNextRefreshTime()
	self.OpenEndTime = time.Now().Unix() + 60*60

	randGoodsIDLst := gamedata.BlackMarketRandGoods(self.ownplayer.GetLevel())

	for _, v := range randGoodsIDLst {
		var goods TBlackMarketGoods
		goods.ID = v
		goods.IsBuy = false
		self.GoodsLst = append(self.GoodsLst, goods)
	}

	self.IsOpen = true

	if isSave == true {
		go self.DB_SaveGoods()
	}
}

func (self *TBlackMarketModule) CheckReset() {
	if time.Now().Unix() > self.RefreshTime {
		self.RefreshGoods(true)
	}

	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TBlackMarketModule) OnNewDay(newday uint32) {
	self.IsOpen = false
	self.OpenEndTime = 0
	self.ResetDay = newday
	go self.DB_Reset()
}

func (self *TBlackMarketModule) GetNextRefreshTime() int64 {
	now := time.Now()
	nowsec := now.Hour()*60*60 + now.Minute()*60 + now.Second()

	refreshSec := 0
	for i := 0; i < len(gamedata.BlackMarketRefreshTime); i++ {
		if nowsec < gamedata.BlackMarketRefreshTime[i] {
			refreshSec = gamedata.BlackMarketRefreshTime[i]
			break
		}
	}

	//! 隔天刷新
	if refreshSec == 0 {
		hour := gamedata.BlackMarketRefreshTime[0] / 3600
		min := (gamedata.BlackMarketRefreshTime[0] - hour*3600) / 60
		sec := gamedata.BlackMarketRefreshTime[0] - hour*3600 - min*60

		now.AddDate(0, 0, 1)
		refreshTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
		return refreshTime.Unix()
	}

	hour := refreshSec / 3600
	min := (refreshSec - hour*3600) / 60
	sec := refreshSec - hour*3600 - min*60
	refreshTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
	return refreshTime.Unix()
}

func (self *TBlackMarketModule) OnDestroy(playerid int32) {

}

func (self *TBlackMarketModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TBlackMarketModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TBlackMarketModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerBlackMarket").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerBlackMarket Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TBlackMarketModule) DB_SaveGoods() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerBlackMarket", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"goodslst":    self.GoodsLst,
		"refreshtime": self.RefreshTime,
		"openendtime": self.OpenEndTime,
		"isopen":      self.IsOpen}})
}

func (self *TBlackMarketModule) DB_UpdateBuyMark(id int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerBlackMarket", bson.M{"_id": self.PlayerID, "goodslst.id": id}, bson.M{"$set": bson.M{"goodslst.$.isbuy": true}})
}

func (self *TBlackMarketModule) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerBlackMarket", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"openendtime": self.OpenEndTime,
		"isopen":      self.IsOpen}})
}

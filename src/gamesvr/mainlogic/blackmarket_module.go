package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"sync"
	"time"
	"utility"
)

type TBlackMarketGoods struct {
	ID    int
	IsBuy bool
}

//! 黑市
type TBlackMarketModule struct {
	PlayerID    int32               `bson:"_id"`
	GoodsLst    []TBlackMarketGoods //! 商品
	RefreshTime int32               //! 刷新时间
	IsOpen      bool                //! 是否开启
	BlackTime   int32               //! 商店结束时间
	ResetDay    uint32              //! 隔天刷新
	ownplayer   *TPlayer
}

func (self *TBlackMarketModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TBlackMarketModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.ResetDay = utility.GetCurDay()
	self.IsOpen = false

	//! 插入数据库
	mongodb.InsertToDB("PlayerBlackMarket", self)
}

func (self *TBlackMarketModule) RefreshGoods(isSave bool) {
	if len(self.GoodsLst) > 0 {
		self.GoodsLst = []TBlackMarketGoods{}
	}

	self.RefreshTime = self.GetNextRefreshTime()
	self.BlackTime = utility.GetCurTime() + 60*60

	randGoodsIDLst := gamedata.BlackMarketRandGoods(self.ownplayer.GetLevel())

	for _, v := range randGoodsIDLst {
		var goods TBlackMarketGoods
		goods.ID = v
		goods.IsBuy = false
		self.GoodsLst = append(self.GoodsLst, goods)
	}

	self.IsOpen = true

	if isSave == true {
		self.DB_SaveGoods()
	}
}

func (self *TBlackMarketModule) CheckReset() {
	if utility.GetCurTime() > self.RefreshTime {
		self.RefreshGoods(true)
	}

	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TBlackMarketModule) OnNewDay(newday uint32) {
	self.IsOpen = false
	self.BlackTime = 0
	self.ResetDay = newday
	self.DB_Reset()
}

func (self *TBlackMarketModule) GetNextRefreshTime() int32 {
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
		return int32(refreshTime.Unix())
	}

	hour := refreshSec / 3600
	min := (refreshSec - hour*3600) / 60
	sec := refreshSec - hour*3600 - min*60
	refreshTime := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, now.Location())
	return int32(refreshTime.Unix())
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

	err := s.DB(appconfig.GameDbName).C("PlayerBlackMarket").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerBlackMarket Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TBlackMarketModule) GetBlackMarketItemInfo(index int) *TBlackMarketGoods {
	if index >= len(self.GoodsLst) || index < 0 {
		gamelog.Error("GetBlackMarketItemInfo Error: Invalid index: %d", index)
		return nil
	}

	return &self.GoodsLst[index]
}

func (self *TBlackMarketModule) DB_SaveGoods() {
	mongodb.UpdateToDB("PlayerBlackMarket", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"goodslst":    self.GoodsLst,
		"refreshtime": self.RefreshTime,
		"blacktime":   self.BlackTime,
		"isopen":      self.IsOpen}})
}

func (self *TBlackMarketModule) DB_UpdateBuyMark(id int) {
	filedName := fmt.Sprintf("goodslst.%d.isbuy", id)
	mongodb.UpdateToDB("PlayerBlackMarket", &bson.M{"_id": self.PlayerID},
		&bson.M{"$set": bson.M{filedName: true}})
}

func (self *TBlackMarketModule) DB_Reset() {
	mongodb.UpdateToDB("PlayerBlackMarket", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"blacktime": self.BlackTime,
		"isopen":    self.IsOpen}})
}

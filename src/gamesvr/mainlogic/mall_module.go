package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

//! 购物
type TItemShoppingInfo struct {
	ItemID   int //! 物品ID
	ItemType int //! 购买物品类型
	BuyTimes int //! 购物次数
}

type TMallModule struct {
	PlayerID     int32 `bson:"_id"`
	ShoppingInfo []TItemShoppingInfo
	ResetDay     uint32 //! 重置时间
	ownplayer    *TPlayer
}

func (self *TMallModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TMallModule) OnCreate(playerid int32) {
	//! 初始化信息
	self.PlayerID = playerid

	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerMall", self)
}

//! 玩家销毁角色
func (self *TMallModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TMallModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TMallModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TMallModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerMall").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerMall Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 重置购买次数
func (self *TMallModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TMallModule) OnNewDay(newday uint32) {
	self.ResetDay = newday
	for i := 0; i < len(self.ShoppingInfo); i++ {
		//! 普通商品购买次数刷新
		if self.ShoppingInfo[i].ItemType == 0 {
			self.ShoppingInfo[i].BuyTimes = 0
		}
	}

	go self.UpdateResetShoppingInfo()
}

//! 获取用户已购买的VIP礼包
func (self *TMallModule) GetUserAleadyShoppingGift(goodstype int) (itemLst IntLst) {
	for _, v := range self.ShoppingInfo {
		if v.ItemType == goodstype {
			itemLst = append(itemLst, v.ItemID)
		}
	}

	return itemLst
}

//! 获取购买次数
func (self *TMallModule) GetItemShoppingInfo(id int) *TItemShoppingInfo {
	for i, v := range self.ShoppingInfo {
		if v.ItemID == id {
			return &self.ShoppingInfo[i]
		}
	}
	return nil
}

//! 增加购买次数
func (self *TMallModule) AddItemShoppingTimes(id int, times int) {

	for i := 0; i < len(self.ShoppingInfo); i++ {
		if self.ShoppingInfo[i].ItemID == id {
			self.ShoppingInfo[i].BuyTimes = self.ShoppingInfo[i].BuyTimes + times
			go self.UpdateShoppingInfo()
		}
	}

}

//! 数据库重置购买次数
func (self *TMallModule) UpdateResetShoppingInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerMall", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"shoppinginfo": self.ShoppingInfo,
		"resetday":     self.ResetDay}})
}

//! 存储购买物品信息
func (self *TMallModule) UpdateShoppingInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerMall", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"shoppinginfo": self.ShoppingInfo}})
}

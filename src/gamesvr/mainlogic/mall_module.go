package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"msg"
	"sync"
	"utility"
)

type TMallModule struct {
	PlayerID     int32             `bson:"_id"`
	NormalRecord []msg.MSG_BuyData //普通商品的次数
	VipBagRecord Int32Lst          //VIP礼包商店的次数
	ResetDay     uint32            //! 重置时间
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
	self.NormalRecord = []msg.MSG_BuyData{}
	self.VipBagRecord = Int32Lst{}

	//! 插入数据库
	mongodb.InsertToDB("PlayerMall", self)
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

	err := s.DB(appconfig.GameDbName).C("PlayerMall").Find(&bson.M{"_id": playerid}).One(self)
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
	for i := 0; i < len(self.NormalRecord); i++ {
		self.NormalRecord[i].Times = 0
	}

	self.UpdateNormalRecord()
}

//! 获取购买次数
func (self *TMallModule) GetItemBuyData(id int32) *msg.MSG_BuyData {
	for i := 0; i < len(self.NormalRecord); i++ {
		if self.NormalRecord[i].ID == id {
			return &self.NormalRecord[i]
		}
	}
	return nil
}

//! 数据库重置购买次数
func (self *TMallModule) UpdateNormalRecord() {
	mongodb.UpdateToDB("PlayerMall", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"normalrecord": self.NormalRecord,
		"resetday":     self.ResetDay}})
}

//! 数据库重置购买次数
func (self *TMallModule) UpdateVipBagRecord() {
	mongodb.UpdateToDB("PlayerMall", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"vipbagrecord": self.VipBagRecord}})
}

//! 数据库重置购买次数
func (self *TMallModule) UpdateNormalRecordAt(nIndex int) {
	filedName := fmt.Sprintf("normalrecord.%d", nIndex)
	mongodb.UpdateToDB("PlayerMall", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.NormalRecord[nIndex]}})
}


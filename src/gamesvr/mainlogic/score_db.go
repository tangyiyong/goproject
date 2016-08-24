package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//!
func (self *TScoreMoudle) DB_SaveScoreAndFightTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"fighttime": self.FightTime,
		"score":     self.Score,
		"wintime":   self.WinTime}})
}

//!
func (self *TScoreMoudle) DB_UpdateRecvAward() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"recvaward": self.RecvAward}})
}

//!保存购买次数
func (self *TScoreMoudle) DB_SaveBuyFightTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"buytime": self.BuyTime}})
}

//! 更新购买次数
func (self *TScoreMoudle) DB_UpdateStoreItemBuyTimes(id int, times int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID, "storebuyrecord.id": id}, bson.M{"$set": bson.M{
		"shoppinglst.$.times": times}})
}

//! 增加购买信息
func (self *TScoreMoudle) DB_AddStoreItemBuyInfo(info TStoreBuyData) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID}, "storebuyrecord", info)
}

//! 重置购买信息
func (self *TScoreMoudle) DB_SaveShoppingInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"storebuyrecord": self.StoreBuyRecord}})
}

//! 增加购买奖励信息
func (self *TScoreMoudle) DB_AddStoreAwardInfo(id int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerScore", bson.M{"_id": self.PlayerID}, "awardstoreindex", id)
}

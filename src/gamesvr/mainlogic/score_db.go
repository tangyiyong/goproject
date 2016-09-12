package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"msg"
)

//!
func (self *TScoreMoudle) DB_SaveScoreAndFightTime() {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"fighttime": self.FightTime,
		"score":     self.Score,
		"serieswin": self.SeriesWin}})
}

//!
func (self *TScoreMoudle) DB_UpdateRecvAward() {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"recvaward": self.RecvAward}})
}

//!保存购买次数
func (self *TScoreMoudle) DB_SaveBuyFightTime() {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buytime": self.BuyTime}})
}

//! 更新购买次数
func (self *TScoreMoudle) DB_UpdateStoreItemBuyTimes(index int, times int) {
	filedName := fmt.Sprintf("buyrecord.%d.times", index)
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: times}})
}

//! 增加购买信息
func (self *TScoreMoudle) DB_AddStoreItemBuyInfo(info msg.MSG_BuyData) {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"buyrecord": info}})
}

//! 重置购买信息
func (self *TScoreMoudle) DB_SaveShoppingInfo() {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buyrecord": self.BuyRecord}})
}

//! 增加购买奖励信息
func (self *TScoreMoudle) DB_AddStoreAwardInfo(id int) {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"awardstoreindex": id}})
}

//! 重置购买信息
func (self *TScoreMoudle) DB_OnNewDay() {
	mongodb.UpdateToDB("PlayerScore", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buytime":   self.BuyTime,
		"serieswin": self.SeriesWin,
		"fighttime": self.FightTime,
		"recvaward": self.RecvAward,
		"buyrecord": self.BuyRecord}})
}

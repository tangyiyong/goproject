package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! 每日重置
func (self *TFameHallModule) DB_Reset() {
	mongodb.UpdateToDB("PlayerFameHall", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"resetday":    self.ResetDay,
		"freetimes":   self.FreeTimes,
		"sendfightid": self.SendFightID,
		"sendlevelid": self.SendLevelID}})
}

//! 更新玩家免费赠送次数
func (self *TFameHallModule) DB_UpdateFreeTimes() {
	mongodb.UpdateToDB("PlayerFameHall", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"freetimes": self.FreeTimes}})
}

//! 更新玩家魅力值
func (self *TFameHallModule) DB_UpdateCharm() {
	mongodb.UpdateToDB("PlayerFameHall", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"charmvalue": self.CharmValue}})
}

//! 增加赠送玩家ID
func (self *TFameHallModule) DB_AddSendFightID(index int32) {
	mongodb.UpdateToDB("PlayerFameHall", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"sendfightid": index}})
}

func (self *TFameHallModule) DB_AddSendLevelID(index int32) {
	mongodb.UpdateToDB("PlayerFameHall", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"sendlevelid": index}})
}

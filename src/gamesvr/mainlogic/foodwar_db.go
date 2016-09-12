package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TFoodWarModule) DB_Reset() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"resetday":        self.ResetDay,
		"attacktimes":     self.AttackTimes,
		"revengetimes":    self.RevengeTimes,
		"nexttime":        self.NextTime,
		"totalfood":       self.TotalFood,
		"fixedfood":       self.FixedFood,
		"revengelst":      self.RevengeLst,
		"buyattacktimes":  self.BuyAttackTimes,
		"buyrevengetimes": self.BuyRevengeTimes,
		"awardrecvlst":    self.AwardRecvLst}})
}

func (self *TFoodWarModule) DB_AddAwardRecvRecord(id int) {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"awardrecvlst": id}})
}

func (self *TFoodWarModule) DB_CheckTime() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"attacktimes": self.AttackTimes,
		"totalfood":   self.TotalFood,
		"fixedfood":   self.FixedFood,
		"nexttime":    self.NextTime}})
}

func (self *TFoodWarModule) DB_SaveFood() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"totalfood": self.TotalFood,
		"fixedfood": self.FixedFood}})
}

func (self *TFoodWarModule) DB_SaveAttackTimes() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"attacktimes": self.AttackTimes}})
}

func (self *TFoodWarModule) DB_SaveRevengeTimes() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"revengetimes": self.RevengeTimes}})
}

func (self *TFoodWarModule) DB_AddRevengeLst(player TRevengeInfo) {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"revengelst": player}})
}

func (self *TFoodWarModule) DB_RemoveRevengeLst(player TRevengeInfo) {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"revengelst": player}})
}

func (self *TFoodWarModule) DB_SaveBuyAttackTimes() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buyattacktimes": self.BuyAttackTimes,
		"attacktimes":    self.AttackTimes}})
}

func (self *TFoodWarModule) DB_SaveBuyRevengeTimes() {
	mongodb.UpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"revengetimes":    self.RevengeTimes,
		"buyrevengetimes": self.BuyRevengeTimes}})
}

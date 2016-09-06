package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
)

func (self *TFoodWarModule) DB_Reset() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
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
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"awardrecvlst": id}})
}

func (self *TFoodWarModule) DB_CheckTime() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"attacktimes": self.AttackTimes,
		"totalfood":   self.TotalFood,
		"fixedfood":   self.FixedFood,
		"nexttime":    self.NextTime}})
}

func (self *TFoodWarModule) DB_SaveFood() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"totalfood": self.TotalFood,
		"fixedfood": self.FixedFood}})
}

func (self *TFoodWarModule) DB_SaveAttackTimes() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"attacktimes": self.AttackTimes}})
}

func (self *TFoodWarModule) DB_SaveRevengeTimes() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"revengetimes": self.RevengeTimes}})
}

func (self *TFoodWarModule) DB_AddRevengeLst(player TRevengeInfo) {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"revengelst": player}})
}

func (self *TFoodWarModule) DB_RemoveRevengeLst(player TRevengeInfo) {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"revengelst": player}})
}

func (self *TFoodWarModule) DB_SaveBuyAttackTimes() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buyattacktimes": self.BuyAttackTimes,
		"attacktimes":    self.AttackTimes}})
}

func (self *TFoodWarModule) DB_SaveBuyRevengeTimes() {
	GameSvrUpdateToDB("PlayerFoodWar", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"revengetimes":    self.RevengeTimes,
		"buyrevengetimes": self.BuyRevengeTimes}})
}

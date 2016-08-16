package mainlogic

import (
	//"fmt"
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TFoodWarModule) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
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
	mongodb.AddToArray(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, "awardrecvlst", id)
}

func (self *TFoodWarModule) DB_CheckTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"attacktimes": self.AttackTimes,
		"totalfood":   self.TotalFood,
		"fixedfood":   self.FixedFood,
		"nexttime":    self.NextTime}})
}

func (self *TFoodWarModule) DB_SaveFood() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"totalfood": self.TotalFood,
		"fixedfood": self.FixedFood}})
}

func (self *TFoodWarModule) DB_SaveAttackTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"attacktimes": self.AttackTimes}})
}

func (self *TFoodWarModule) DB_SaveRevengeTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"revengetimes": self.RevengeTimes}})
}

func (self *TFoodWarModule) DB_AddRevengeLst(player TRevengeInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, "revengelst", player)
}

func (self *TFoodWarModule) DB_RemoveRevengeLst(player TRevengeInfo) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, "revengelst", player)
}

func (self *TFoodWarModule) DB_SaveBuyAttackTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"buyattacktimes": self.BuyAttackTimes,
		"attacktimes":    self.AttackTimes}})
}

func (self *TFoodWarModule) DB_SaveBuyRevengeTimes() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerFoodWar", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"revengetimes":    self.RevengeTimes,
		"buyrevengetimes": self.BuyRevengeTimes}})
}

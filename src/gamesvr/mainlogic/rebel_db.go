package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 更新重置信息
func (self *TRebelModule) UpdateResetTime() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRebel", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"exploit":         self.Exploit,
		"exploitawardlst": self.ExploitAwardLst,
		"damage":          self.Damage,
		"resetday":        self.ResetDay}})
}

//! 更新领取标记
func (self *TRebelModule) UpdateExploitAward(id int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerRebel", bson.M{"_id": self.PlayerID}, "exploitawardlst", id)
}

//! 更新叛军信息
func (self *TRebelModule) UpdateRebelInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRebel", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"rebelid":    self.RebelID,
		"curlife":    self.CurLife,
		"level":      self.Level,
		"escapetime": self.EscapeTime,
		"isshare":    self.IsShare}})
}

//! 更新战功信息
func (self *TRebelModule) UpdateExploit() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerRebel", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"exploit": self.Exploit,
		"damage":  self.Damage}})
}

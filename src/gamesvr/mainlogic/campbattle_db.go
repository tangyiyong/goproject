package mainlogic

import (
	"appconfig"
	//"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//修改一个阵营战阵营
func (self *TCampBattleModule) DB_SaveBattleCamp() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"battlecamp": self.BattleCamp}})
}

func (self *TCampBattleModule) DB_SaveKillData() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"kill": self.Kill,
		"killsum":    self.KillSum,
		"killhonor":  self.KillHonor,
		"destroy":    self.Destroy,
		"destroysum": self.DestroySum}})
}

func (self *TCampBattleModule) DB_SaveMoveStautus() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"lefttimes": self.LeftTimes,
		"endtime":   self.EndTime,
		"crystalid": self.CrystalID}})
}

func (self *TCampBattleModule) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"lefttimes": self.LeftTimes,
		"endtime":   self.EndTime,
		"killhonor": self.KillHonor,
		"kill":      self.Kill,
		"resetday":  self.ResetDay,
		"destroy":   self.Destroy}})
}

//! 更新购买次数
func (self *TCampBattleModule) DB_UpdateStoreItemBuyTimes(id int, times int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID, "storebuyrecord.id": id}, bson.M{"$set": bson.M{
		"shoppinglst.$.times": times}})
}

//! 增加购买信息
func (self *TCampBattleModule) DB_AddStoreItemBuyInfo(info TStoreBuyData) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, "storebuyrecord", info)
}

//! 重置购买信息
func (self *TCampBattleModule) DB_SaveShoppingInfo() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"storebuyrecord": self.StoreBuyRecord}})
}

//! 增加购买奖励信息
func (self *TCampBattleModule) DB_AddStoreAwardInfo(id int) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerCampBat", bson.M{"_id": self.PlayerID}, "awardstoreindex", id)
}

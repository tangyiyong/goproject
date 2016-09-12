package mainlogic

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
	"msg"
)

//修改一个阵营战阵营
func (self *TCampBattleModule) DB_SaveBattleCamp() {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"battlecamp": self.BattleCamp}})
}

func (self *TCampBattleModule) DB_SaveKillData() {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"kill": self.Kill,
		"killsum":    self.KillSum,
		"killhonor":  self.KillHonor,
		"destroy":    self.Destroy,
		"destroysum": self.DestroySum}})
}

func (self *TCampBattleModule) DB_SaveMoveStautus() {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"lefttimes": self.LeftTimes,
		"endtime":   self.EndTime,
		"crystalid": self.CrystalID}})
}

func (self *TCampBattleModule) DB_Reset() {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"lefttimes": self.LeftTimes,
		"endtime":   self.EndTime,
		"killhonor": self.KillHonor,
		"kill":      self.Kill,
		"resetday":  self.ResetDay,
		"destroy":   self.Destroy}})
}

//! 更新购买次数
func (self *TCampBattleModule) DB_UpdateStoreItemBuyTimes(nindex int, times int) {
	filedName := fmt.Sprintf("buyrecord.%d", nindex)
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.BuyRecord[nindex]}})
}

//! 增加购买信息
func (self *TCampBattleModule) DB_AddStoreItemBuyInfo(info msg.MSG_BuyData) {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"buyrecord": info}})
}

//! 重置购买信息
func (self *TCampBattleModule) DB_SaveShoppingInfo() {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buyrecord": self.BuyRecord}})
}

//! 增加购买奖励信息
func (self *TCampBattleModule) DB_AddStoreAwardInfo(id int) {
	mongodb.UpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"awardstoreindex": id}})
}

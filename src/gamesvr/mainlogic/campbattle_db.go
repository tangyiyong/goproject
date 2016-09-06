package mainlogic

import (
	"fmt"
	"msg"

	"gopkg.in/mgo.v2/bson"
)

//修改一个阵营战阵营
func (self *TCampBattleModule) DB_SaveBattleCamp() {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"battlecamp": self.BattleCamp}})
}

func (self *TCampBattleModule) DB_SaveKillData() {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"kill": self.Kill,
		"killsum":    self.KillSum,
		"killhonor":  self.KillHonor,
		"destroy":    self.Destroy,
		"destroysum": self.DestroySum}})
}

func (self *TCampBattleModule) DB_SaveMoveStautus() {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"lefttimes": self.LeftTimes,
		"endtime":   self.EndTime,
		"crystalid": self.CrystalID}})
}

func (self *TCampBattleModule) DB_Reset() {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
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
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		filedName: self.BuyRecord[nindex]}})
}

//! 增加购买信息
func (self *TCampBattleModule) DB_AddStoreItemBuyInfo(info msg.MSG_BuyData) {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"buyrecord": info}})
}

//! 重置购买信息
func (self *TCampBattleModule) DB_SaveShoppingInfo() {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"buyrecord": self.BuyRecord}})
}

//! 增加购买奖励信息
func (self *TCampBattleModule) DB_AddStoreAwardInfo(id int) {
	GameSvrUpdateToDB("PlayerCampBat", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"awardstoreindex": id}})
}

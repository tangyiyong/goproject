package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

//! DB相关
func (self *TBeachBabyInfo) DB_SaveAllGoods() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"beachbaby.goods":    self.Goods,
		"beachbaby.autotime": self.AutoTime}})
}
func (self *TBeachBabyInfo) DB_SaveScore() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"beachbaby.score":      self.Score,
		"beachbaby.totalscore": self.TotalScore}})
}
func (self *TBeachBabyInfo) DB_SaveRankAwardFlag() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"beachbaby.rankaward": self.RankAward}})
}
func (self *TBeachBabyInfo) DB_Refresh() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"beachbaby.score":        self.Score,
		"beachbaby.totalscore":   self.TotalScore,
		"beachbaby.freeconchbit": self.FreeConch,
		"beachbaby.rankaward":    self.RankAward,
		"beachbaby.activityid":   self.ActivityID,
		"beachbaby.versioncode":  self.VersionCode,
		"beachbaby.resetcode":    self.ResetCode}})
}

func (self *TBeachBabyInfo) DB_Reset() {
	mongodb.UpdateToDB("PlayerActivity", &bson.M{"_id": self.activityModule.PlayerID}, &bson.M{"$set": bson.M{
		"beachbaby.score":        self.Score,
		"beachbaby.totalscore":   self.TotalScore,
		"beachbaby.freeconchbit": self.FreeConch,
		"beachbaby.rankaward":    self.RankAward,
		"beachbaby.activityid":   self.ActivityID,
		"beachbaby.versioncode":  self.VersionCode,
		"beachbaby.resetcode":    self.ResetCode}})
}

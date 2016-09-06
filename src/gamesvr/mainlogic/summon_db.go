package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
)

func (self *TSummonModule) DB_UpdateNormalSummon() {
	GameSvrUpdateToDB("PlayerSummon", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"normal": self.Normal}})
}

func (self *TSummonModule) DB_UpdateSeniorSummon() {
	GameSvrUpdateToDB("PlayerSummon", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"senior": self.Senior}})
}

func (self *TSummonModule) DB_UpdateFirstSummon() {
	GameSvrUpdateToDB("PlayerSummon", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"isfirst": self.IsFirst}})
}

package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TSummonModule) DB_UpdateNormalSummon() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSummon", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"normal": self.Normal}})
}

func (self *TSummonModule) DB_UpdateSeniorSummon() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSummon", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"senior": self.Senior}})
}

func (self *TSummonModule) DB_UpdateFirstSummon() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSummon", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"isfirst": self.IsFirst}})
}

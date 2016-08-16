package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TSummonModule) UpdateNormalSummon() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSummon", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"normal": self.Normal}})
}

func (self *TSummonModule) UpdateSeniorSummon() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSummon", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"senior": self.Senior}})
}

func (self *TSummonModule) UpdateFirstSummon() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSummon", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"isfirst": self.IsFirst}})
}

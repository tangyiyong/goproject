package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TWanderModule) DB_Reset() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerWander", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"resetday": self.ResetDay,
		"curcopyid":  self.CurCopyID,
		"canbattle":  self.CanBattle,
		"maxcopyid":  self.MaxCopyID,
		"singlefree": self.SingleFree,
		"lefttime":   self.LeftTime}})
}

func (self *TWanderModule) DB_ResetSingleFreeDay() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerWander", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{"sglfreeday": self.SglFreeDay,
		"singlefree": self.SingleFree}})
}

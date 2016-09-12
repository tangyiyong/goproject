package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TWanderModule) DB_Reset() {
	mongodb.UpdateToDB("PlayerWander", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"resetday": self.ResetDay,
		"curcopyid":  self.CurCopyID,
		"canbattle":  self.CanBattle,
		"maxcopyid":  self.MaxCopyID,
		"singlefree": self.SingleFree,
		"lefttime":   self.LeftTime}})
}

func (self *TWanderModule) DB_ResetSingleFreeDay() {
	mongodb.UpdateToDB("PlayerWander", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{"sglfreeday": self.SglFreeDay,
		"singlefree": self.SingleFree}})
}

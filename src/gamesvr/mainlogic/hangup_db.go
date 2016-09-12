package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (hang *THangUpMoudle) DB_SaveHangUpState() {
	mongodb.UpdateToDB("PlayerHang", &bson.M{"_id": hang.PlayerID}, &bson.M{"$set": bson.M{
		"curbossid": hang.CurBossID,
		"starttime": hang.StartTime,
		"expitems":  hang.ExpItems,
		"history":   hang.History,
		"gridnum":   hang.GridNum,
		"quicktime": hang.QuickTime}})
}

func (hang *THangUpMoudle) DB_ClearHangUpBag() {
	mongodb.UpdateToDB("PlayerHang", &bson.M{"_id": hang.PlayerID}, &bson.M{"$set": bson.M{"expitems": hang.ExpItems}})
}

func (hang *THangUpMoudle) DB_SaveQuickFightResult() {
	mongodb.UpdateToDB("PlayerHang", &bson.M{"_id": hang.PlayerID}, &bson.M{"$set": bson.M{
		"expitems":  hang.ExpItems,
		"history":   hang.History,
		"quicktime": hang.QuickTime}})
}

func (hang *THangUpMoudle) DB_SaveQuickFightTime() {
	mongodb.UpdateToDB("PlayerHang", &bson.M{"_id": hang.PlayerID}, &bson.M{"$set": bson.M{"quicktime": hang.QuickTime}})
}

func (hang *THangUpMoudle) DB_SaveGridNum() {
	mongodb.UpdateToDB("PlayerHang", &bson.M{"_id": hang.PlayerID}, &bson.M{"$set": bson.M{"gridnum": hang.GridNum}})
}

package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (simple *TSimpleInfoMgr) DB_SetFightValue(playerid int32, fightvalue int, level int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"fightvalue": fightvalue, "level": level}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetPlayerName(playerid int32, name string) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"name": name}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetLogoffTime(playerid int32, time int64) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"logofftime": time}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetLoginDay(playerid int32, day uint32) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"loginday": day}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetHeroID(playerid int32, heroid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"heroid": heroid}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetHeroQuality(playerid int32, quality int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"quality": quality}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetVipLevel(playerid int32, viplevel int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"viplevel": viplevel}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetBatCamp(playerid int32, camp int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"batcamp": camp}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetAwardCenterID(playerid int32, awardCenterID int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"awardcenterid": awardCenterID}})
}

func (simple *TSimpleInfoMgr) DB_SetGuildID(playerid int32, guildid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"guildid": guildid}})
}

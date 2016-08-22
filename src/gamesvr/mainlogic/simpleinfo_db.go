package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (simple *TSimpleInfoMgr) DB_SetFightValue(playerid int, fightvalue int, level int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"fightvalue": fightvalue, "level": level}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetPlayerName(playerid int, name string) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"name": name}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetLogoffTime(playerid int, time int64) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"logofftime": time}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetLoginDay(playerid int, day int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"loginday": day}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetHeroID(playerid int, heroid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"heroid": heroid}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetHeroQuality(playerid int, quality int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"quality": quality}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetVipLevel(playerid int, viplevel int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"viplevel": viplevel}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetBatCamp(playerid int, camp int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerid}, bson.M{"$set": bson.M{"batcamp": camp}})
	return
}

func (simple *TSimpleInfoMgr) DB_SetAwardCenterID(playerID int, awardCenterID int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerID}, bson.M{"$set": bson.M{"awardcenterid": awardCenterID}})
}

func (simple *TSimpleInfoMgr) DB_SetGuildID(playerID int, guildid int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerSimple", bson.M{"_id": playerID}, bson.M{"$set": bson.M{"guildid": guildid}})
}

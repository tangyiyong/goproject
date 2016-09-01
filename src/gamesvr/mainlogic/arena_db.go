package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TArenaModule) DB_UpdateRankToDatabase() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerArena", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"currentrank": self.CurrentRank,
		"historyrank": self.HistoryRank}})
}

func (self *TArenaModule) DB_UpdateStoreToDatabase() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerArena", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"storeaward": self.StoreAward}})
}

func (self *TArenaModule) DB_UpdateChallangeRank(playerid int32, rank int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerArena", bson.M{"_id": playerid}, bson.M{"$set": bson.M{
		"currentrank": rank}})
}

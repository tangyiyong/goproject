package mainlogic

import (
	"appconfig"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TArenaModule) UpdateRankToDatabase() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerArena", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"currentrank": self.CurrentRank,
		"historyrank": self.HistoryRank}})
}

func (self *TArenaModule) UpdateStoreToDatabase() {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerArena", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		"storeaward": self.StoreAward}})
}

func (self *TArenaModule) UpdateChallangeRank(playerid int32, rank int) {
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerArena", bson.M{"_id": playerid}, bson.M{"$set": bson.M{
		"currentrank": rank}})
}

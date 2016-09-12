package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func (self *TArenaModule) DB_UpdateRankToDatabase() {
	mongodb.UpdateToDB("PlayerArena", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"currentrank": self.CurrentRank,
		"historyrank": self.HistoryRank}})
}

func (self *TArenaModule) DB_UpdateStoreToDatabase() {
	mongodb.UpdateToDB("PlayerArena", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"storeaward": self.StoreAward}})
}

func (self *TArenaModule) DB_UpdateChallangeRank(playerid int32, rank int) {
	mongodb.UpdateToDB("PlayerArena", &bson.M{"_id": playerid}, &bson.M{"$set": bson.M{
		"currentrank": rank}})
}

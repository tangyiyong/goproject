package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
)

func (self *TArenaModule) DB_UpdateRankToDatabase() {
	GameSvrUpdateToDB("PlayerArena", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"currentrank": self.CurrentRank,
		"historyrank": self.HistoryRank}})
}

func (self *TArenaModule) DB_UpdateStoreToDatabase() {
	GameSvrUpdateToDB("PlayerArena", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{
		"storeaward": self.StoreAward}})
}

func (self *TArenaModule) DB_UpdateChallangeRank(playerid int32, rank int) {
	GameSvrUpdateToDB("PlayerArena", &bson.M{"_id": playerid}, &bson.M{"$set": bson.M{
		"currentrank": rank}})
}

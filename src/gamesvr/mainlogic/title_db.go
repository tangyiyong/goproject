package mainlogic

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

func (self *TTitleModule) DB_AddTitleInfo(info TitleInfo) {
	GameSvrUpdateToDB("PlayerTitle", &bson.M{"_id": self.PlayerID}, &bson.M{"$push": bson.M{"titlelst": info}})
}

func (self *TTitleModule) DB_RemoveTitleInfo(info *TitleInfo) {
	GameSvrUpdateToDB("PlayerTitle", &bson.M{"_id": self.PlayerID}, &bson.M{"$pull": bson.M{"titlelst": *info}})
}

func (self *TTitleModule) DB_UpdateTitleStatus(index int, status int) {
	fieldName := fmt.Sprint("titlelst.%d.status", index)
	GameSvrUpdateToDB("PlayerTitle", &bson.M{"_id": self.PlayerID}, &bson.M{"$set": bson.M{fieldName: status}})
}

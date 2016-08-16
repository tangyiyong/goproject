package mainlogic

import (
	"appconfig"
	"fmt"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

func (self *TTitleModule) DB_AddTitleInfo(info TitleInfo) {
	mongodb.AddToArray(appconfig.GameDbName, "PlayerTitle", bson.M{"_id": self.PlayerID}, "titlelst", info)
}

func (self *TTitleModule) DB_RemoveTitleInfo(info *TitleInfo) {
	mongodb.RemoveFromArray(appconfig.GameDbName, "PlayerTitle", bson.M{"_id": self.PlayerID}, "titlelst", *info)
}

func (self *TTitleModule) DB_UpdateTitleStatus(index int, status int) {
	fieldName := fmt.Sprint("titlelst.%d.status", index)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerTitle", bson.M{"_id": self.PlayerID}, bson.M{"$set": bson.M{
		fieldName: status}})
}

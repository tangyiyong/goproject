package mainlogic

import (
	"appconfig"
	"fmt"
	"gamelog"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
)

//! 更新VIP日常福利领取时间到数据库
func (playerVip *TVipMoudle) DB_SaveFisrtCharge(id int) {
	if id <= 0 {
		gamelog.Error("DB_SaveFisrtCharge Error: Invalid id :%d", id)
		return
	}
	fieldName := fmt.Sprintf("firstcharges.%d", id-1)
	mongodb.UpdateToDB(appconfig.GameDbName, "PlayerVip", bson.M{"_id": playerVip.PlayerID}, bson.M{"$set": bson.M{fieldName: playerVip.FirstCharges[id-1]}})
}

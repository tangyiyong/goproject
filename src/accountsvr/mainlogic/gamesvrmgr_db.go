package mainlogic

import (
	"appconfig"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func DB_UpdateSvrState(svrid int, svrflag int) {
	mongodb.UpdateToDB(appconfig.AccountDbName, "GameSvrList", bson.M{"_id": svrid}, bson.M{"$set": bson.M{"svrflag": svrflag}})
}

package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func DB_UpdateSvrInfo(svrid int32, svrinfo TGameServerInfo) {
	mongodb.UpdateToDB("GameSvrList", &bson.M{"_id": svrid}, &bson.M{"$set": svrinfo})
}

func DB_UpdateLastSvrID(accountid int32, svrid int32) {
	mongodb.UpdateToDB("Account", &bson.M{"_id": accountid}, &bson.M{"$set": bson.M{"lastsvrid": svrid}})
}

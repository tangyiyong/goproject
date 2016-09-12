package mainlogic

import (
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

func DB_UpdateSvrState(svrid int32, svrflag uint32) {
	mongodb.UpdateToDB("GameSvrList", &bson.M{"_id": svrid}, &bson.M{"$set": bson.M{"svrflag": svrflag}})
}

func DB_UpdateCountAndLastSvr(accountid int32, svrid int32) {
	mongodb.UpdateToDB("Account", &bson.M{"_id": accountid}, &bson.M{"$set": bson.M{"lastsvrid": svrid, "logincount": 1}})
}

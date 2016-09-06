package mainlogic

import (
	"appconfig"
	"gamelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mongodb"
)

var G_DbSession *mgo.Session = nil
var G_ColMap map[string]*mgo.Collection
var G_LastColName string = "x"

type TDB_Param struct {
	IsAll   bool    //是否更新全部记录
	ColName string  //集合名
	Search  *bson.M //条件
	Stuff   *bson.M //数据
}

var G_DB_ParamList chan *TDB_Param //参数队列

func GameSvrUpdateToDB(collection string, search *bson.M, stuff *bson.M) {
	//mongodb.UpdateToDB(appconfig.GameDbName, collection, search, stuff)

	//以下新的实现
	var param TDB_Param
	param.IsAll = false
	param.ColName = collection
	param.Search = search
	param.Stuff = stuff
	G_DB_ParamList <- &param
}

func GameSvrUpdateToDBAll(collection string, search *bson.M, stuff *bson.M) {
	//mongodb.UpdateToDBAll(appconfig.GameDbName, collection, search, stuff)

	//以下新的实现
	var param TDB_Param
	param.IsAll = true
	param.ColName = collection
	param.Search = search
	param.Stuff = stuff

	G_DB_ParamList <- &param
}

func InitDbProcesser() bool {
	G_DB_ParamList = make(chan *TDB_Param, 1024)
	G_ColMap = make(map[string]*mgo.Collection, 1)
	G_DbSession = mongodb.GetDBSession()
	if G_DbSession == nil {
		gamelog.Error("InitDbProcesser  Error : GetDBSession Failed!!!")
		return false
	}

	go DBProcess()

	return true
}

func DBProcess() {
	var collectoin *mgo.Collection = nil
	var err error
	var ok bool
	var maxcount int = 0
	for param := range G_DB_ParamList {
		if maxcount < len(G_DB_ParamList) {
			maxcount = len(G_DB_ParamList)
			gamelog.Error("DBProcess :%d wait to process", maxcount)
		}

		if param.ColName != G_LastColName {
			collectoin, ok = G_ColMap[param.ColName]
			if ok == false || collectoin == nil {
				collectoin = G_DbSession.DB(appconfig.GameDbName).C(param.ColName)
				G_ColMap[param.ColName] = collectoin
			}
		}

		if param.IsAll {
			_, err = collectoin.UpdateAll(param.Search, param.Stuff)
		} else {
			err = collectoin.Update(param.Search, param.Stuff)
		}

		if err != nil {
			gamelog.Error3("UpdateToDB Failed:Collection:[%s] search:[%v], stuff:[%v], Error:%v", param.ColName, param.Search, param.Stuff, err.Error())
		}
	}
}

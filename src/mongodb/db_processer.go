package mongodb

import (
	"gamelog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB_OP_UPDATE_ALL    = 1
	DB_OP_UPDATE_SINGLE = 2
	DB_OP_INSERT        = 3
	DB_OP_REMOVE        = 4
)

var G_DbSession *mgo.Session = nil
var G_ColMap map[string]*mgo.Collection
var G_LastColName string = "x"
var G_Db_Name string

type TDB_Param struct {
	OpType  int32       //数据操作类型
	ColName string      //集合名
	Search  *bson.M     //条件
	Stuff   *bson.M     //数据
	Data    interface{} //账号数据
}

var G_DB_ParamList chan *TDB_Param //参数队列

//type TDB_Collection struct {
//	Name       string          //集合名
//	ParamList  chan *TDB_Param //参数队列
//	Session    *mgo.Session
//	Collection *mgo.Collection
//}
//
//func (self *TDB_Collection) Init(name string) {
//	self.ParamList = make(chan *TDB_Param, 10240)
//	self.Session = GetDBSession()
//	self.Name = name
//	self.Collection = self.Session.DB().C(name)

//	go
//}
//
//var G_CollectionMapEx map[string]*TDB_Collection

func InsertToDB(collection string, pData interface{}) {
	//以下新的实现
	var param TDB_Param
	param.OpType = DB_OP_INSERT
	param.ColName = collection
	param.Data = pData
	G_DB_ParamList <- &param
}

func UpdateToDB(collection string, search *bson.M, stuff *bson.M) {
	//以下新的实现
	var param TDB_Param
	param.OpType = DB_OP_UPDATE_SINGLE
	param.ColName = collection
	param.Search = search
	param.Stuff = stuff
	G_DB_ParamList <- &param
}

func UpdateToDBAll(collection string, search *bson.M, stuff *bson.M) {
	//以下新的实现
	var param TDB_Param
	param.OpType = DB_OP_UPDATE_ALL
	param.ColName = collection
	param.Search = search
	param.Stuff = stuff

	G_DB_ParamList <- &param
}

func InitDbProcesser(dbname string) bool {
	G_DB_ParamList = make(chan *TDB_Param, 10240)
	G_ColMap = make(map[string]*mgo.Collection, 1)
	G_DbSession = GetDBSession()
	G_Db_Name = dbname
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
				collectoin = G_DbSession.DB(G_Db_Name).C(param.ColName)
				G_ColMap[param.ColName] = collectoin
			}
		}

		if param.OpType == DB_OP_INSERT {
			err = collectoin.Insert(param.Data)
		} else if param.OpType == DB_OP_UPDATE_ALL {
			//nil interface 不等于nil
			if param.Search == nil {
				_, err = collectoin.UpdateAll(nil, param.Stuff)
			} else {
				_, err = collectoin.UpdateAll(param.Search, param.Stuff)
			}
		} else if param.OpType == DB_OP_UPDATE_SINGLE {
			err = collectoin.Update(param.Search, param.Stuff)
		}

		if err != nil {
			gamelog.Error("DBProcess Failed:Collection:[%s] search:[%v], stuff:[%v], Error:%v", param.ColName, param.Search, param.Stuff, err.Error())
		}
	}
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gopkg.in/mgo.v2"
	"mongodb"
	"sync"
	"time"
)

type TGameServerInfo struct {
	SvrDomainID   int32 `bson:"_id"` //账号ID
	SvrDomainName string
	SvrFlag       int32
	svrOutAddr    string
	svrInnerAddr  string //内部地址
	isSvrOK       bool   //服务器是否正常
	updateTime    int64
}

var (
	G_ServerList [10000]TGameServerInfo
	ListLock     sync.Mutex
)

func InitGameSvrMgr() {
	s := mongodb.GetDBSession()
	defer s.Close()

	var tempList []TGameServerInfo
	err := s.DB(appconfig.AccountDbName).C("GameSvrList").Find(nil).Sort("+_id").All(&tempList)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitGameSvrMgr DB Error!!!")
		return
	}

	for i := 0; i < len(tempList); i++ {
		G_ServerList[tempList[i].SvrDomainID].SvrDomainID = tempList[i].SvrDomainID
		G_ServerList[tempList[i].SvrDomainID].SvrDomainName = tempList[i].SvrDomainName
		G_ServerList[tempList[i].SvrDomainID].SvrFlag = tempList[i].SvrFlag
	}

	go CheckGameStateRoutine()

	return
}

func CheckGameStateRoutine() {
	regtimer := time.Tick(10 * time.Second)
	for {
		ListLock.Lock()
		curtime := time.Now().Unix()
		for i := 0; i < 10000; i++ {
			if G_ServerList[i].SvrDomainID <= 0 {
				continue
			}

			if (curtime - G_ServerList[i].updateTime) > (70) {
				G_ServerList[i].isSvrOK = false
			}
		}
		ListLock.Unlock()
		<-regtimer
	}
}

func UpdateGameSvrInfo(domainid int32, svrname string, outaddr string, inaddr string) {
	if domainid <= 0 || domainid >= 10000 {
		gamelog.Error("UpdateGameSvrInfo Error : Invalid DomainID:%d", domainid)
		return
	}

	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[domainid].SvrDomainID == 0 {
		G_ServerList[domainid].SvrDomainID = domainid
		G_ServerList[domainid].SvrDomainName = svrname
		G_ServerList[domainid].svrInnerAddr = inaddr
		G_ServerList[domainid].svrOutAddr = outaddr
		G_ServerList[domainid].isSvrOK = true
		G_ServerList[domainid].SvrFlag = 1
		G_ServerList[domainid].updateTime = time.Now().Unix()
		mongodb.InsertToDB(appconfig.AccountDbName, "GameSvrList", &G_ServerList[domainid])
	} else {
		if G_ServerList[domainid].SvrDomainName != svrname {
			gamelog.Error("UpdateGameSvrInfo Error : %d has two domainname:%s, %s", domainid, svrname, G_ServerList[domainid].SvrDomainName)
		}
		G_ServerList[domainid].svrInnerAddr = inaddr
		G_ServerList[domainid].svrOutAddr = outaddr
		G_ServerList[domainid].isSvrOK = true
		G_ServerList[domainid].updateTime = time.Now().Unix()
	}

}

func GetGameSvrName(domainid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[domainid].SvrDomainID == 0 {
		gamelog.Error("GetGameSvrName Error Invalid domainid :%d", domainid)
		return ""
	}

	return G_ServerList[domainid].SvrDomainName
}

func GetGameSvrOutAddr(domainid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[domainid].SvrDomainID == 0 {
		gamelog.Error("GetGameSvrAddr Error Invalid domainid :%d", domainid)
		return ""
	}

	return G_ServerList[domainid].svrOutAddr
}

func GetGameSvrInAddr(domainid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()
	if G_ServerList[domainid].SvrDomainID == 0 {
		return ""
	}

	return G_ServerList[domainid].svrInnerAddr
}

func GetGameSvrInfo(domainid int32) (pInfo *TGameServerInfo) {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[domainid].SvrDomainID == 0 {
		gamelog.Error("GetGameSvrInfo Error Invalid domainid :%d", domainid)
		return nil
	}

	return &G_ServerList[domainid]
}

func RemoveGameSvrInfo(domainid int32) bool {
	ListLock.Lock()
	defer ListLock.Unlock()
	G_ServerList[domainid].SvrDomainID = 0
	return true
}

func GetRecommendSvrID() *TGameServerInfo {
	ListLock.Lock()
	defer ListLock.Unlock()

	if len(G_ServerList) <= 0 {
		gamelog.Error("GetRecommendSvrID Error Has No Server!!")
		return nil
	}

	for i := 0; i < 10000; i++ {
		if G_ServerList[i].SvrDomainID != 0 && G_ServerList[i].isSvrOK == true {
			return &G_ServerList[i]
		}
	}

	gamelog.Error("GetRecommendSvrID Error No Avaliable Server!!")
	return nil
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gopkg.in/mgo.v2"
	"mongodb"
	"sync"
	"time"
)

const (
	SFG_RECOMMAND = 0x00000001 //推荐服务器
	SFG_CREATE    = 0x00000002 //允许创建角色
	SFG_LOGIN     = 0x00000004 //允许登录角色
	SFG_VISIBLE   = 0x00000008 //服务器可见

	////组合标记
	SFG_ALL    = 0x0000000F
	SFG_NORMAL = 0x0000000E //可注册，可登录，可见，
)

type TGameServerInfo struct {
	SvrID   int32 `bson:"_id"` //账号ID
	SvrName string
	SvrFlag uint32 //游戏服标记

	//以下的变量不是存数据库
	svrOutAddr   string
	svrInnerAddr string //内部地址
	isSvrOK      bool   //服务器是否正常
	updateTime   int64  //更新时间
}

var (
	G_ServerList  [10000]TGameServerInfo
	G_RecommendID int32 //推荐的服务器ID
	ListLock      sync.Mutex
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
		G_ServerList[tempList[i].SvrID].SvrID = tempList[i].SvrID
		G_ServerList[tempList[i].SvrID].SvrName = tempList[i].SvrName
		G_ServerList[tempList[i].SvrID].SvrFlag = tempList[i].SvrFlag
	}

	//go CheckGameStateRoutine()

	return
}

func CheckGameStateRoutine() {
	regtimer := time.Tick(10 * time.Second)
	for {
		ListLock.Lock()
		curtime := time.Now().Unix()
		for i := 0; i < 10000; i++ {
			if G_ServerList[i].SvrID <= 0 {
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

func UpdateGameSvrInfo(svrid int32, svrname string, outaddr string, inaddr string) {
	if svrid <= 0 || svrid >= 10000 {
		gamelog.Error("UpdateGameSvrInfo Error : Invalid svrid:%d", svrid)
		return
	}

	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[svrid].SvrID == 0 {
		G_ServerList[svrid].SvrID = svrid
		G_ServerList[svrid].SvrName = svrname
		G_ServerList[svrid].svrInnerAddr = inaddr
		G_ServerList[svrid].svrOutAddr = outaddr
		G_ServerList[svrid].isSvrOK = true
		G_ServerList[svrid].SvrFlag = SFG_ALL
		G_ServerList[svrid].updateTime = time.Now().Unix()
		mongodb.InsertToDB("GameSvrList", &G_ServerList[svrid])
	} else {
		if G_ServerList[svrid].SvrName != svrname {
			gamelog.Error("UpdateGameSvrInfo Error : %d has two domainname:%s, %s", svrid, svrname, G_ServerList[svrid].SvrName)
		}
		G_ServerList[svrid].SvrName = svrname
		G_ServerList[svrid].svrInnerAddr = inaddr
		G_ServerList[svrid].svrOutAddr = outaddr
		G_ServerList[svrid].isSvrOK = true
		G_ServerList[svrid].updateTime = time.Now().Unix()
	}

}

func GetGameSvrName(svrid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrName Error Invalid svrid :%d", svrid)
		return ""
	}

	return G_ServerList[svrid].SvrName
}

func GetGameSvrOutAddr(svrid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrAddr Error Invalid svrid :%d", svrid)
		return ""
	}

	return G_ServerList[svrid].svrOutAddr
}

func GetGameSvrInAddr(svrid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()
	if G_ServerList[svrid].SvrID == 0 {
		return ""
	}

	return G_ServerList[svrid].svrInnerAddr
}

func GetGameSvrInfo(svrid int32) (pInfo *TGameServerInfo) {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrInfo Error Invalid svrid :%d", svrid)
		return nil
	}

	return &G_ServerList[svrid]
}

func RemoveGameSvrInfo(svrid int32) bool {
	ListLock.Lock()
	defer ListLock.Unlock()
	G_ServerList[svrid].SvrID = 0
	return true
}

func GetRecommendSvrID() *TGameServerInfo {
	ListLock.Lock()
	defer ListLock.Unlock()

	if G_RecommendID > 0 && G_ServerList[G_RecommendID].SvrID != 0 && G_ServerList[G_RecommendID].isSvrOK == true && (G_ServerList[G_RecommendID].SvrFlag&SFG_VISIBLE > 0) {
		return &G_ServerList[G_RecommendID]
	}

	for i := 9999; i > 0; i-- {
		if G_ServerList[i].SvrID != 0 && G_ServerList[i].isSvrOK == true && (G_ServerList[i].SvrFlag&SFG_VISIBLE) > 0 {
			return &G_ServerList[i]
		}
	}

	gamelog.Error("GetRecommendSvrID Error No Avaliable Server!!")
	return nil
}

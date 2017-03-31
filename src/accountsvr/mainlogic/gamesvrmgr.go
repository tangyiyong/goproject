package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"utility"

	"gopkg.in/mgo.v2"
)

const (
	SS_Ready     = 1 //未开放
	SS_NewSvr    = 2 //新服
	SS_Good      = 3 //流畅
	SS_Busy      = 4 //拥挤
	SS_Maintain  = 5 //维护
	SS_Full      = 6 //爆满
	SS_Recommend = 7 //推荐
	SS_Close     = 8 //关闭
)

type TGameServerInfo struct {
	SvrID        int32  `bson:"_id"` //账号ID
	SvrName      string //服务器名字
	ControlFlag  uint32 //控制标记
	SvrState     uint32 //显示标记
	SvrDefault   uint32 //是否默认
	SvrOutAddr   string //外部地址
	SvrInnerAddr string //内部地址
	SvrOpenTime  int32  //服务器开服时间

	//以下的变量不是存数据库
	isSvrOK    bool  //服务器是否正常
	updateTime int32 //更新时间
}

var (
	G_ServerList  [10000]TGameServerInfo
	G_RecommendID int32 //推荐的服务器ID
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
		G_ServerList[tempList[i].SvrID].SvrState = tempList[i].SvrState
		G_ServerList[tempList[i].SvrID].SvrInnerAddr = tempList[i].SvrInnerAddr
		G_ServerList[tempList[i].SvrID].SvrOutAddr = tempList[i].SvrOutAddr
		G_ServerList[tempList[i].SvrID].SvrDefault = tempList[i].SvrDefault
		G_ServerList[tempList[i].SvrID].SvrOpenTime = tempList[i].SvrOpenTime
		G_ServerList[tempList[i].SvrID].isSvrOK = false

		if G_ServerList[tempList[i].SvrID].SvrDefault == 1 {
			G_RecommendID = tempList[i].SvrID
		}
	}

	go CheckGameStateRoutine()

	return
}

/*
func CheckGameStateRoutine() {
	regtimer := time.Tick(10 * time.Second)
	for {
		ListLock.Lock()
		curtime := utility.GetCurTime()
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
*/

func UpdateGameSvrInfo(svrid int32, outaddr string, inaddr string, opentime int32) string {
	if svrid <= 0 || svrid >= 10000 {
		gamelog.Error("UpdateGameSvrInfo Error : Invalid svrid:%d", svrid)
		return "Error"
	}

	if G_ServerList[svrid].SvrID == 0 {
		G_ServerList[svrid].SvrID = svrid
		G_ServerList[svrid].SvrName = "UNKNOW"
		G_ServerList[svrid].SvrInnerAddr = inaddr
		G_ServerList[svrid].SvrOutAddr = outaddr
		G_ServerList[svrid].isSvrOK = true
		G_ServerList[svrid].SvrState = SS_Ready
		G_ServerList[svrid].SvrDefault = 0
		G_ServerList[svrid].SvrOpenTime = 0
		G_ServerList[svrid].updateTime = utility.GetCurTime()
		mongodb.InsertToDB("GameSvrList", &G_ServerList[svrid])
	} else {
		if G_ServerList[svrid].SvrOutAddr != outaddr {
			gamelog.Error("UpdateGameSvrInfo Error : **************** Server:%d has two addr:[%s]-[%s]*********", svrid, G_ServerList[svrid].SvrOutAddr, outaddr)
		}
		G_ServerList[svrid].SvrInnerAddr = inaddr
		G_ServerList[svrid].SvrOutAddr = outaddr
		G_ServerList[svrid].isSvrOK = true
		G_ServerList[svrid].SvrOpenTime = opentime
		G_ServerList[svrid].updateTime = utility.GetCurTime()
		DB_UpdateSvrInfo(svrid, G_ServerList[svrid])
	}

	return G_ServerList[svrid].SvrName
}

func GetGameSvrName(svrid int32) string {
	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrName Error Invalid svrid :%d", svrid)
		return ""
	}

	return G_ServerList[svrid].SvrName
}

func GetGameSvrOutAddr(svrid int32) string {
	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrAddr Error Invalid svrid :%d", svrid)
		return ""
	}

	return G_ServerList[svrid].SvrOutAddr
}

func GetGameSvrInAddr(svrid int32) string {
	if G_ServerList[svrid].SvrID == 0 {
		return ""
	}

	return G_ServerList[svrid].SvrInnerAddr
}

func GetGameSvrInfo(svrid int32) (pInfo *TGameServerInfo) {
	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrInfo Error Invalid svrid :%d", svrid)
		return nil
	}

	return &G_ServerList[svrid]
}

func RemoveGameSvrInfo(svrid int32) bool {
	G_ServerList[svrid].SvrID = 0
	return true
}

func GetRecommendSvrID() *TGameServerInfo {
	if G_RecommendID > 0 && G_ServerList[G_RecommendID].SvrID != 0 && G_ServerList[G_RecommendID].isSvrOK == true && (G_ServerList[G_RecommendID].SvrState > SS_Ready) {
		return &G_ServerList[G_RecommendID]
	}

	for i := 9999; i > 0; i-- {
		if G_ServerList[i].SvrID != 0 && G_ServerList[i].isSvrOK == true && G_ServerList[i].SvrState > SS_Ready {
			return &G_ServerList[i]
		}
	}

	gamelog.Error("GetRecommendSvrID Error No Avaliable Server!!")
	return nil
}

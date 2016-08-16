package mainlogic

import (
	"gamelog"
	"gamesvr/tcpclient"
	"sync"
)

type TBattleServerInfo struct {
	BatSvrID     int    //阵营战服务器的ID(就是端口号)
	SvrState     int    // 0: 不可用 1:可用
	SvrOutAddr   string //外部地址(带端口号)
	SvrInnerAddr string //内部地址(带端口号)
	PlayerNum    int    //玩家人数
	BatClient    tcpclient.TCPClient
}

var (
	G_ServerList = make(map[int]*TBattleServerInfo)
	ListLock     sync.Mutex
)

func GetRecommendSvrAddr() (ret string) {
	ListLock.Lock()
	defer ListLock.Unlock()
	if len(G_ServerList) <= 0 {
		gamelog.Error("GetRecommendSvrAddr Error : Not Avalible Battle Server!!!")
		ret = ""
		return
	}

	for _, v := range G_ServerList {
		if v.SvrState >= 1 {
			ret = v.SvrOutAddr
			return
		}
	}

	ret = ""
	gamelog.Error("GetRecommendSvrAddr Error : Not Avalible Battle Server!!!")
	return
}

func SetBattleSvrConnectOK(svrid int, conok bool) bool {
	ListLock.Lock()
	defer ListLock.Unlock()

	pBatSvrInfo, ok := G_ServerList[svrid]
	if !ok || pBatSvrInfo == nil {
		return false
	}

	if conok {
		pBatSvrInfo.SvrState = 1
	} else {
		pBatSvrInfo.SvrState = 0
	}

	return true
}

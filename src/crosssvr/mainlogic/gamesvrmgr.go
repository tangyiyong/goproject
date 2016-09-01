package mainlogic

import (
	"gamelog"
	"sync"
	"time"
)

type TGameServerInfo struct {
	SvrDomainID   int32
	SvrDomainName string
	UpdateTime    int64
	SvrOutAddr    string
	SvrInnerAddr  string //内部地址
	IsSvrOK       bool   //服务器是否正常
	SvrFightUrl   string //var reqUrl = "http://" + addr + "/get_fight_target"
}

var (
	G_ServerList   [10000]TGameServerInfo
	CurSelectIndex = 0
	ListLock       sync.Mutex
)

func InitGameSvrMgr() {
	go func() {
		regtimer := time.Tick(10 * time.Second)
		for {
			ListLock.Lock()
			curtime := time.Now().Unix()
			for i := 0; i < 10000; i++ {
				if G_ServerList[i].SvrDomainID <= 0 {
					continue
				}

				if (curtime - G_ServerList[i].UpdateTime) > (70) {
					G_ServerList[i].IsSvrOK = false
				}
			}
			ListLock.Unlock()
			<-regtimer
		}
	}()

	return
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
		G_ServerList[domainid].SvrInnerAddr = inaddr
		G_ServerList[domainid].SvrOutAddr = outaddr
		G_ServerList[domainid].IsSvrOK = true
		G_ServerList[domainid].UpdateTime = time.Now().Unix()
	} else {
		if G_ServerList[domainid].SvrDomainName != svrname {
			gamelog.Error("UpdateGameSvrInfo Error : %d has two domainname:%s, %s", domainid, svrname, G_ServerList[domainid].SvrDomainName)
		}

		G_ServerList[domainid].SvrInnerAddr = inaddr
		G_ServerList[domainid].SvrOutAddr = outaddr
		G_ServerList[domainid].IsSvrOK = true
		G_ServerList[domainid].UpdateTime = time.Now().Unix()
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

	return G_ServerList[domainid].SvrOutAddr
}

func GetGameSvrInAddr(domainid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()
	if G_ServerList[domainid].SvrDomainID == 0 {
		return ""
	}

	return G_ServerList[domainid].SvrInnerAddr
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

func GetSelectSvrAddr() string {
	ListLock.Lock()
	defer ListLock.Unlock()

	for _, v := range G_ServerList {
		if v.IsSvrOK == true {
			return v.SvrInnerAddr
		}
	}

	return ""
}

func GetGameSvrFightTarAddr(domainid int32) string {
	ListLock.Lock()
	defer ListLock.Unlock()
	if G_ServerList[domainid].SvrDomainID == 0 {
		return ""
	}

	if len(G_ServerList[domainid].SvrFightUrl) <= 0 {
		G_ServerList[domainid].SvrFightUrl = "http://" + G_ServerList[domainid].SvrInnerAddr + "/get_fight_target"
	}

	return G_ServerList[domainid].SvrFightUrl
}

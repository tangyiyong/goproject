package mainlogic

import (
	"gamelog"
	"sync"
	"time"
)

const (
	SVR_STATE_CLOSE  = 0
	SVR_STATE_NORMAL = 1
)

type GameServerInfo struct {
	SvrDomainID   int
	SvrDomainName string
	SvrState      int //
	UpdateTime    int64
	SvrOutAddr    string
	SvrInnerAddr  string //内部地址
	SvrFightUrl   string //var reqUrl = "http://" + addr + "/get_fight_target"
}

var (
	G_ServerList   = make(map[int]*GameServerInfo)
	CurSelectIndex = 0
	ListLock       sync.Mutex
)

func InitGameSvrMgr() {
	go func() {
		regtimer := time.Tick(10 * time.Second)
		for {
			ListLock.Lock()
			curtime := time.Now().Unix()
			for _, info := range G_ServerList {
				if (curtime - info.UpdateTime) > (70) {
					info.SvrState = 0
				}
			}
			ListLock.Unlock()
			<-regtimer
		}
	}()

	return
}

func UpdateGameSvrInfo(domainid int, doname string, outaddr string, inaddr string) {
	if domainid <= 0 {
		return
	}

	ListLock.Lock()
	defer ListLock.Unlock()

	pGameSvrInfo, ok := G_ServerList[domainid]
	if !ok || pGameSvrInfo == nil {
		var pInfo *GameServerInfo = new(GameServerInfo)
		pInfo.SvrDomainID = domainid
		pInfo.SvrDomainName = doname
		pInfo.SvrInnerAddr = inaddr
		pInfo.SvrOutAddr = outaddr
		pInfo.SvrState = 1
		pInfo.UpdateTime = time.Now().Unix()
		G_ServerList[domainid] = pInfo
		return
	}

	if pGameSvrInfo.SvrDomainName != doname {
		gamelog.Error("UpdateGameSvrInfo Error : %d has two domainname:%s, %s", domainid, doname, pGameSvrInfo.SvrDomainName)
	}

	pGameSvrInfo.UpdateTime = time.Now().Unix()
	pGameSvrInfo.SvrState = 1

}

func AddGameSvrInfo(pInfo *GameServerInfo) {
	ListLock.Lock()
	defer ListLock.Unlock()

	if pInfo == nil {
		gamelog.Error("AddGameSvrInfo Error pInof is nil")
		return
	}

	pGameSvrInfo, ok := G_ServerList[pInfo.SvrDomainID]
	if !ok || pGameSvrInfo == nil {

		G_ServerList[pInfo.SvrDomainID] = pInfo

		return
	}

	if pGameSvrInfo.SvrDomainName != pInfo.SvrDomainName {
		gamelog.Error("AddGameSvrInfo Error : %d has tow domainname:%s, %s", pInfo.SvrDomainID, pInfo.SvrDomainName, pGameSvrInfo.SvrDomainName)
	}

	pGameSvrInfo.UpdateTime = pInfo.UpdateTime
	pGameSvrInfo.SvrState = 1

	return
}

func GetGameSvrName(domainid int) string {
	ListLock.Lock()
	defer ListLock.Unlock()

	pGameSvrInfo, ok := G_ServerList[domainid]
	if !ok || pGameSvrInfo == nil {
		return ""
	}

	return pGameSvrInfo.SvrDomainName
}

func GetGameSvrInAddr(domainid int) string {
	ListLock.Lock()
	defer ListLock.Unlock()
	pGameSvrInfo, ok := G_ServerList[domainid]
	if !ok || pGameSvrInfo == nil {
		return ""
	}

	return pGameSvrInfo.SvrInnerAddr
}

func GetGameSvrFightTarAddr(domainid int) string {
	ListLock.Lock()
	defer ListLock.Unlock()
	pGameSvrInfo, ok := G_ServerList[domainid]
	if !ok || pGameSvrInfo == nil {
		return ""
	}

	if len(pGameSvrInfo.SvrFightUrl) <= 0 {
		pGameSvrInfo.SvrFightUrl = "http://" + pGameSvrInfo.SvrInnerAddr + "/get_fight_target"
	}

	return pGameSvrInfo.SvrFightUrl
}

func GetGameSvrInfo(domainid int) (pInfo *GameServerInfo) {
	ListLock.Lock()
	defer ListLock.Unlock()
	pGameSvrInfo, ok := G_ServerList[domainid]
	if !ok || pGameSvrInfo == nil {
		return nil
	}

	return pGameSvrInfo
}

func RemoveGameSvrInfo(domainid int) bool {
	ListLock.Lock()
	defer ListLock.Unlock()

	delete(G_ServerList, domainid)

	return true
}

func GetSelectSvrAddr() string {
	ListLock.Lock()
	defer ListLock.Unlock()

	for _, v := range G_ServerList {
		if v.SvrState == SVR_STATE_NORMAL {
			return v.SvrInnerAddr
		}
	}

	return ""
}

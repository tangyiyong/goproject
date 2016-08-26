package mainlogic

import (
	"gamelog"
	"sync"
	"tcpserver"
)

type TLogData struct {
	ServerID int32
}

var (
	G_SvrConns   map[int32]*tcpserver.TCPConn
	G_ConnsMutex sync.Mutex
)

func InitConnMgr() bool {
	G_SvrConns = make(map[int32]*tcpserver.TCPConn, 1)
	return true
}

func GetConnByID(playerid int32) *tcpserver.TCPConn {
	G_ConnsMutex.Lock()
	pConn, _ := G_SvrConns[playerid]
	G_ConnsMutex.Unlock()
	return pConn
}

func DelConnByID(playerid int32) {
	G_ConnsMutex.Lock()
	delete(G_SvrConns, playerid)
	G_ConnsMutex.Unlock()
	return
}

func AddConnByID(playerid int32, pTcpConn *tcpserver.TCPConn) {
	G_ConnsMutex.Lock()
	G_SvrConns[playerid] = pTcpConn
	G_ConnsMutex.Unlock()
	return
}

func AddTcpConn(serverid int32, name string, pTcpConn *tcpserver.TCPConn) {
	pTcpConn.Data = new(TLogData)
	pTcpConn.Data.(*TLogData).ServerID = serverid
	pTcpConn.Cleaned = false
	AddConnByID(serverid, pTcpConn)
	return
}

func CheckAndClean(serverid int32) {
	if serverid == 0 {
		gamelog.Error("CheckAndClean Error: Invalid serverid:0")
		return
	}
	G_ConnsMutex.Lock()
	defer G_ConnsMutex.Unlock()
	pOldConn, ok := G_SvrConns[serverid]
	if !ok {
		return
	}

	delete(G_SvrConns, serverid)
	pOldConn.Cleaned = true
	pOldConn.Close()
}

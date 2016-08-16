package msgprocess

import (
	"gamelog"
	"sync"
	"tcpserver"
)

type TChatData struct {
	GuildID  int
	PlayerID int
}

type TGuildConns struct {
	connMap map[int]*tcpserver.TCPConn
}

var (
	G_GuildConns  map[int]*TGuildConns
	G_NameConns   map[string]*tcpserver.TCPConn
	G_PlayerGuild map[int]int
	G_PlayerConns map[int]*tcpserver.TCPConn
	G_GameSvrConn *tcpserver.TCPConn
	G_ConnsMutex  sync.Mutex
)

func Init() bool {
	G_GuildConns = make(map[int]*TGuildConns, 1)
	G_NameConns = make(map[string]*tcpserver.TCPConn, 1)
	G_PlayerConns = make(map[int]*tcpserver.TCPConn, 1)
	G_PlayerGuild = make(map[int]int, 1)

	return true
}

func (guild *TGuildConns) Init(guildid int) {
	guild.connMap = make(map[int]*tcpserver.TCPConn, 30)
}

func GetConnByID(playerid int) *tcpserver.TCPConn {
	G_ConnsMutex.Lock()
	pConn, _ := G_PlayerConns[playerid]
	G_ConnsMutex.Unlock()
	return pConn
}

func AddTcpConn(playerid int, guildid int, name string, pTcpConn *tcpserver.TCPConn) {
	G_ConnsMutex.Lock()
	defer G_ConnsMutex.Unlock()

	pTcpConn.Data = new(TChatData)
	pTcpConn.Data.(*TChatData).GuildID = guildid
	pTcpConn.Data.(*TChatData).PlayerID = playerid
	pTcpConn.Cleaned = false

	tGuildConns, ok := G_GuildConns[guildid]
	if ok {
		tGuildConns.connMap[playerid] = pTcpConn
	} else {
		var pGuildConns = new(TGuildConns)
		pGuildConns.Init(guildid)
		pGuildConns.connMap[playerid] = pTcpConn
		G_GuildConns[guildid] = pGuildConns
		G_PlayerGuild[playerid] = guildid
	}

	G_NameConns[name] = pTcpConn
	G_PlayerConns[playerid] = pTcpConn
	return
}

func CheckAndClean(playerid int) {
	if playerid == 0 {
		gamelog.Error("CheckAndClean Error: Invalid PlayerID:0")
		return
	}
	G_ConnsMutex.Lock()
	defer G_ConnsMutex.Unlock()
	pOldConn, ok := G_PlayerConns[playerid]
	if !ok {
		return
	}

	GuildID, ok := G_PlayerGuild[playerid]
	if ok {
		oldGuildConns, ok := G_GuildConns[GuildID]
		if ok {
			delete(oldGuildConns.connMap, playerid)
		}
	}

	delete(G_PlayerGuild, playerid)
	delete(G_PlayerConns, playerid)
	pOldConn.Cleaned = true
	pOldConn.Close()
}

func ChangeConnGuild(playerid int, newguildid int) {
	G_ConnsMutex.Lock()
	defer G_ConnsMutex.Unlock()

	pTcpConn, ok := G_PlayerConns[playerid]
	if !ok || pTcpConn == nil {
		return
	}

	if pTcpConn.Data.(*TChatData).GuildID == newguildid {
		return
	}

	//首先从之前的位置清掉
	oldGuildConns, ok := G_GuildConns[pTcpConn.Data.(*TChatData).GuildID]
	if ok {
		delete(oldGuildConns.connMap, playerid)
	}

	tGuildConns, ok := G_GuildConns[newguildid]
	if ok {
		tGuildConns.connMap[playerid] = pTcpConn
	} else {
		var pGuildConns = new(TGuildConns)
		pGuildConns.Init(newguildid)
		pGuildConns.connMap[playerid] = pTcpConn
		G_GuildConns[newguildid] = pGuildConns
		G_PlayerGuild[playerid] = newguildid
	}

	return
}

func SendMessageByID(playerid int, msgid int16, msgdata []byte) bool {
	pConn := GetConnByID(playerid)
	if pConn == nil {
		gamelog.Error("SendMessageByID Invalid playerid : %d", playerid)
		return false
	}

	return pConn.WriteMsg(msgid, msgdata)
}

func SendMessageByName(playername string, msgid int16, msgdata []byte) bool {
	G_ConnsMutex.Lock()
	pConn, ok := G_NameConns[playername]
	if !ok {
		G_ConnsMutex.Unlock()
		gamelog.Error("SendMessageByName Invalid name : %s", playername)
		return false
	}
	G_ConnsMutex.Unlock()

	return pConn.WriteMsg(msgid, msgdata)
}

func SendMessageToGuild(guildid int, msgid int16, msgdata []byte, sendPlayerID int) bool {
	G_ConnsMutex.Lock()
	tGuildConns, ok := G_GuildConns[guildid]
	if !ok {
		G_ConnsMutex.Unlock()
		gamelog.Error("SendMessageToGuild: can not find the target guild!!!")
		return false
	}
	G_ConnsMutex.Unlock()
	for playerID, conn := range tGuildConns.connMap {
		if playerID == sendPlayerID {
			continue
		}
		conn.WriteMsg(msgid, msgdata)
	}

	return true
}

func SendMessageToWorld(msgid int16, msgdata []byte, sendPlayerID int) bool {
	G_ConnsMutex.Lock()
	for playerID, conn := range G_PlayerConns {
		if playerID == sendPlayerID {
			continue
		}
		conn.WriteMsg(msgid, msgdata)
	}
	G_ConnsMutex.Unlock()
	return true
}

func SendMessageToGameSvr(msgid int16, msgdata []byte) bool {
	if G_GameSvrConn == nil {
		gamelog.Error("SendMessageToGameSvr Error, G_GameSvrConn is nil!!")
		return false
	}

	G_GameSvrConn.WriteMsg(msgid, msgdata)

	return true
}

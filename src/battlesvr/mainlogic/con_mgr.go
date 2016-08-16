package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
	"strconv"
	"sync"
	"tcpserver"
	"time"
)

type TBattleData struct {
	PlayerID int
	RoomID   int
	PackNo   int //包序号
}

var (
	G_PlayerConns  map[int]*tcpserver.TCPConn
	G_ConnsMutex   sync.Mutex
	G_GameSvrConns *tcpserver.TCPConn = nil
)

func InitConMgr() bool {
	G_PlayerConns = make(map[int]*tcpserver.TCPConn, 1)
	return true
}

func GetConnByID(playerid int) *tcpserver.TCPConn {
	G_ConnsMutex.Lock()
	pConn, _ := G_PlayerConns[playerid]
	G_ConnsMutex.Unlock()
	return pConn
}

func DelConnByID(playerid int) {
	G_ConnsMutex.Lock()
	delete(G_PlayerConns, playerid)
	G_ConnsMutex.Unlock()
	return
}

func AddConnByID(playerid int, pTcpConn *tcpserver.TCPConn) {
	G_ConnsMutex.Lock()
	G_PlayerConns[playerid] = pTcpConn
	G_ConnsMutex.Unlock()
	return
}

func AddTcpConn(playerid int, roomid int, pTcpConn *tcpserver.TCPConn) {
	pData := new(TBattleData)
	pData.RoomID = roomid
	pData.PlayerID = playerid
	pTcpConn.Data = pData
	pTcpConn.Cleaned = false
	AddConnByID(playerid, pTcpConn)
	return
}

func CheckAndClean(playerid int) {
	if playerid == 0 {
		gamelog.Error("CheckAndClean Error: Invalid PlayerID:0")
		return
	}
	G_ConnsMutex.Lock()
	pConn := GetConnByID(playerid)
	if pConn != nil {
		DelConnByID(playerid)
		pBattleData := pConn.Data.(*TBattleData)
		G_RoomMgr.RemovePlayerFromRoom(pBattleData.RoomID, playerid)
		pConn.Cleaned = true
		pConn.Close()
		gamelog.Error("CheckAndClean Error: Clean the unclosed Connection:%d", playerid)
	}
}

//func SendMessageToPlayer(playerid int, msgid int16, msgdata []byte) bool {
//	G_ConnsMutex.Lock()
//	pConn, ok := G_PlayerConns[playerid]
//	if !ok {
//		G_ConnsMutex.Unlock()
//		gamelog.Error("SendMessageToPlayer Invalid playerid : %d", playerid)
//		return false
//	}
//	G_ConnsMutex.Unlock()

//	return pConn.WriteMsg(msgid, msgdata)
//}

func SendMessageToPlayer(playerid int, msgid int16, pmsg msg.TMsg) bool {
	var writer msg.PacketWriter
	writer.BeginWrite(msgid)
	pmsg.Write(&writer)
	writer.EndWrite()

	G_ConnsMutex.Lock()
	pConn, ok := G_PlayerConns[playerid]
	if !ok {
		G_ConnsMutex.Unlock()
		gamelog.Error("SendMessageToPlayer Invalid playerid : %d", playerid)
		return false
	}
	G_ConnsMutex.Unlock()

	return pConn.WriteMsgData(writer.GetDataPtr())
}

func SendMessageToRoom(playerid int, roomid int, msgid int16, pmsg msg.TMsg) bool {
	if roomid <= 0 {
		gamelog.Error("SendMessageToRoom Invalid roomid : %d ", roomid)
		return false
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Error("SendMessageToRoom Invalid roomid : %d", roomid)
		return false
	}

	var writer msg.PacketWriter
	writer.BeginWrite(msgid)
	pmsg.Write(&writer)
	writer.EndWrite()

	G_ConnsMutex.Lock()
	for i := 0; i < max_room_player; i++ {
		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID != playerid && pRoom.Players[i].PlayerID > 0 {
			pConn, ok := G_PlayerConns[pRoom.Players[i].PlayerID]
			if ok && pConn != nil {
				pConn.WriteMsgData(writer.GetDataPtr())
			}
		}
	}
	G_ConnsMutex.Unlock()

	return true
}

//func SendMessageToRoom(playerid int, roomid int, msgid int16, msgdata []byte) bool {
//	pRoom := G_RoomMgr.GetRoomByID(roomid)
//	if pRoom == nil {
//		gamelog.Error("SendMessageToRoom Invalid roomid : %d", roomid)
//		return false
//	}

//	G_ConnsMutex.Lock()
//	for i := 0; i < max_room_player; i++ {
//		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID != playerid && pRoom.Players[i].PlayerID > 0 {
//			pConn, ok := G_PlayerConns[pRoom.Players[i].PlayerID]
//			if ok && pConn != nil {
//				pConn.WriteMsg(msgid, msgdata)
//			}
//		}
//	}
//	G_ConnsMutex.Unlock()

//	return true
//}

func SendMessageToGameSvr(msgid int16, pmsg msg.TMsg) bool {
	if G_GameSvrConns == nil {
		gamelog.Error("SendMessageToGameSvr Error : G_GameSvrConns is Nil!!")
		return false
	}

	var writer msg.PacketWriter
	writer.BeginWrite(msgid)
	pmsg.Write(&writer)
	writer.EndWrite()
	G_GameSvrConns.WriteMsgData(writer.GetDataPtr())
	return true
}

func RegisterToGameSvr() {
	var registerReq msg.MSG_RegBattleSvr_Req
	registerReq.BatSvrID = appconfig.BattleSvrPort
	registerReq.ServerOuterAddr = appconfig.BattleSvrOuterIp + ":" + strconv.Itoa(appconfig.BattleSvrPort)
	registerReq.ServerInnerAddr = appconfig.BattleSvrInnerIp + ":" + strconv.Itoa(appconfig.BattleSvrPort)
	b, _ := json.Marshal(registerReq)

	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(appconfig.RegToGameSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to Game Server failed, err : %s !!!!", err.Error())
			time.Sleep(2 * time.Second)
			continue
		}

		response.Body.Close()
		return
	}
}

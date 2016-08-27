package mainlogic

import (
	"battlesvr/gamedata"
	"gamelog"
	"msg"
	"runtime"
	"runtime/debug"
	"tcpserver"
)

func Init() bool {
	//初始化连接管理器
	InitConMgr()

	//初始化房间管理器
	InitRoomMgr()

	//加载配制文件
	gamedata.LoadConfig()

	return true
}

func BatSvrMsgDispatcher(pTcpConn *tcpserver.TCPConn, MsgID int16, extra int16, pdata []byte) {
	if pTcpConn == nil {
		gamelog.Error("BatSvrMsgDispatcher Error: pTcpConn == nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				gamelog.Error("MsgID %d Error  %s", MsgID, debug.Stack())
			}
		}
	}()

	switch MsgID {
	case msg.MSG_ENTER_ROOM_REQ:
		G_RoomMgr.Hand_EnterRoom(pTcpConn, pdata)
	case msg.MSG_LOAD_CAMPBAT_ACK:
		G_RoomMgr.Hand_LoadCampBatAck(pTcpConn, pdata)
	case msg.MSG_CHECK_IN_REQ:
		G_RoomMgr.Hand_CheckInReq(pTcpConn, pdata)
	case msg.MSG_DISCONNECT:
		G_RoomMgr.Hand_Disconnect(pTcpConn, pdata)
	case msg.MSG_HEART_BEAT:
		G_RoomMgr.Hand_HeartBeat(pTcpConn, pdata)
	default:
		{
			pBatData := pTcpConn.Data.(*TBattleData)
			if pTcpConn.Data == nil || pBatData == nil {
				gamelog.Error("BatSvrMsgDispatcher Error: pBatData == nil ||pTcpConn.Data == nil ")
				return
			}

			var pRoom *TBattleRoom = nil
			if pBatData.RoomID == 0 || pBatData.PlayerID < 10000 {
				pRoom = G_RoomMgr.GetRoomByID(extra)
			} else {
				pRoom = G_RoomMgr.GetRoomByID(pBatData.RoomID)
			}

			if pRoom == nil {
				gamelog.Error("BatSvrMsgDispatcher Error: Invalid RoomID:%d, PlayerID:%d", pBatData.RoomID, pBatData.PlayerID)
				return
			}

			var msg TMessage
			msg.MsgID = MsgID
			msg.MsgData = pdata
			pRoom.MsgList <- msg
		}

	}

}

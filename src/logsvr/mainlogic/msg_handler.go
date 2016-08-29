package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"tcpserver"
)

func Hand_CheckInReq(pTcpConn *tcpserver.TCPConn, extra int16, pdata []byte) {
	var req msg.MSG_CheckIn_Req
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_CheckInReq : Unmarshal error!!!!")
		return
	}

	if req.PlayerID > 10000 || (req.GuildID != -1) {
		gamelog.Error("Hand_CheckInReq Invalid playerid:%d", req.PlayerID)
		pTcpConn.Close()
		return
	}

	CheckAndClean(req.PlayerID)
	gamelog.Info("message: Hand_CheckInReq id:%d, name:%s", req.PlayerID, req.PlayerName)
	AddTcpConn(req.PlayerID, req.PlayerName, pTcpConn)

	return
}

func Hand_OnLogData(pTcpConn *tcpserver.TCPConn, extra int16, pdata []byte) {
	G_LogDataMgr.Append(pdata)
	//如果不用写数据库，则消息不用解析出来
	//如果需要写数据库，则消息需要解析出来
	// var req msg.MSG_SvrLogData
	//if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
	//	gamelog.Error("Hand_OnLogData : Message Reader Error!!!!")
	//	return
	//}
}

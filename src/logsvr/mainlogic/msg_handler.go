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
	//CreateLogFile(pTcpConn.ConnID)
	return
}

func Hand_DisConnect(pTcpConn *tcpserver.TCPConn, extra int16, pdata []byte) {
	if pTcpConn == nil || pTcpConn.ConnID <= 0 {
		return
	}

	if pTcpConn.Cleaned == false {
		CheckAndClean(pTcpConn.ConnID)
	}

	return
}

func Hand_OnLogData(pTcpConn *tcpserver.TCPConn, extra int16, pdata []byte) {
	WriteSvrLog(pdata, pTcpConn.ConnID)
	//如果不用写数据库，则消息不用解析出来
	//如果需要写数据库，则消息需要解析出来
	// var req msg.MSG_SvrLogData
	//if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
	//	gamelog.Error("Hand_OnLogData : Message Reader Error!!!!")
	//	return
	//}
}

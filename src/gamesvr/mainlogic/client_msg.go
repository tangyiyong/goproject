package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/tcpclient"
	"msg"
	"time"
)

func Hand_Connect(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_Connect")
	SendCheckInMsg(pTcpConn)

	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	if pClient.ConType == tcpclient.CON_TYPE_LOGSVR {

	} else if pClient.ConType == tcpclient.CON_TYPE_CHAT {

	} else if pClient.ConType == tcpclient.CON_TYPE_BATSVR {
		SetBattleSvrConnectOK(pClient.SvrID, true)
	}

	return
}

func Hand_DisConnect(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_DisConnect")

	pClient := pTcpConn.Data.(*tcpclient.TCPClient)
	if pClient == nil {
		return
	}

	if pClient.ConType == tcpclient.CON_TYPE_LOGSVR {

	} else if pClient.ConType == tcpclient.CON_TYPE_CHAT {

	} else if pClient.ConType == tcpclient.CON_TYPE_BATSVR {
		SetBattleSvrConnectOK(pClient.SvrID, false)
	}

	return
}

func Hand_OnlineNotify(pTcpConn *tcpclient.TCPConn, extra int16, pdata []byte) {
	gamelog.Info("message: Hand_OnlineNotify")
	var req msg.MSG_OnlineNotify_Req
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_OnlineNotify : Unmarshal error!!!!")
		return
	}

	if req.PlayerID == 0 {
		gamelog.Error("Hand_OnlineNotify req.PlayerID == 0")
		return
	}

	pSimple := G_SimpleMgr.GetSimpleInfoByID(req.PlayerID)
	if pSimple == nil {
		gamelog.Error("Hand_OnlineNotify : Error pSimple is nil!!!!")
		return
	}

	pSimple.isOnline = req.Online

	if req.Online == false {
		pSimple.LogoffTime = time.Now().Unix()
		G_SimpleMgr.DB_SetLogoffTime(req.PlayerID, pSimple.LogoffTime)
	}

	return
}

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/tcpclient"
	"msg"
	"time"
)

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

package msgprocess

import (
	"encoding/json"
	"msg"
)

//向游戏服发送上下线通知
func SendOnlineNotify(playerid int32, online bool) bool {
	var req msg.MSG_OnlineNotify_Req
	req.PlayerID = playerid
	req.Online = online
	buff, _ := json.Marshal(req)
	return SendMessageToGameSvr(msg.MSG_ONLINE_NOTIFY, buff)
}

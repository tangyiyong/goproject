package mainlogic

import (
	"appconfig"
	"encoding/json"
	"gamesvr/tcpclient"
	"msg"
)

var G_ChatClient tcpclient.TCPClient

func SendMessageToChat(msgID int16, msgdata []byte) bool {
	return G_ChatClient.TcpConn.WriteMsg(msgID, msgdata)
}

//直接将消息发送到客户端
func SendMessageToClient(playerid int, msgID int16, msgdata []byte) bool {
	return G_ChatClient.TcpConn.WriteMsgContinue(playerid, msgID, msgdata)
}

func ConnectToChatSvr(addr string) {
	G_ChatClient.ConType = tcpclient.CON_TYPE_CHAT
	G_ChatClient.SvrID = 1 //聊天服忽略此值
	G_ChatClient.ConnectToSvr(addr, 10)
}

func SendCheckInMsg(pTcpConn *tcpclient.TCPConn) bool {
	var req msg.MSG_CheckIn_Req
	req.GuildID = -1
	req.PlayerID = appconfig.DomainID
	req.PlayerName = "gamesvr"
	buffer, _ := json.Marshal(&req)
	return pTcpConn.WriteMsg(msg.MSG_CHECK_IN_REQ, buffer)
}

func SendGuildChangeMsg(playerid int, guilid int) bool {
	var req msg.MSG_GuildNotify_Req
	req.PlayerID = playerid
	req.NewGuildID = guilid
	buffer, _ := json.Marshal(&req)

	return SendMessageToChat(msg.MSG_GUILD_NOTIFY, buffer)
}

package mainlogic

import (
	"appconfig"
	"encoding/json"
	"gamelog"
	"gamesvr/tcpclient"
	"msg"
	"strconv"
)

var G_ChatClient tcpclient.TCPClient
var G_LogClient tcpclient.TCPClient

//给聊天服发送消息
func SendMessageToChat(msgID int16, extra int16, msgdata []byte) bool {
	if G_ChatClient.TcpConn == nil {
		gamelog.Error("SendMessageToChat Error: G_ChatClient.TcpConn is nullptr!!!")
		return false
	}
	return G_ChatClient.TcpConn.WriteMsg(msgID, extra, msgdata)
}

//直接将消息发送到客户端
func SendMessageToClient(playerid int32, msgID int16, extra int16, msgdata []byte) bool {
	if G_ChatClient.TcpConn == nil {
		gamelog.Error("SendMessageToClient Error: G_ChatClient.TcpConn is nullptr!!!")
		return false
	}
	return G_ChatClient.TcpConn.WriteMsgContinue(playerid, msgID, extra, msgdata)
}

//游戏服连接其它的服务器
func ConnectToOtherSvr() bool {
	//连接聊天服务器
	ConnectToChatSvr(appconfig.ChatSvrInnerIp + ":" + strconv.Itoa(appconfig.ChatSvrPort))

	//连接日志服务器
	ConnectToLogSvr(appconfig.LogSvrOuterIp + ":" + strconv.Itoa(appconfig.LogSvrPort))
	return true
}

//连接聊天服务器
func ConnectToChatSvr(addr string) {
	G_ChatClient.ConType = tcpclient.CON_TYPE_CHAT
	G_ChatClient.SvrID = 1 //聊天服忽略此值
	G_ChatClient.ConnectToSvr(addr, 10)
}

//连接日志服务器
func ConnectToLogSvr(addr string) {
	G_LogClient.ConType = tcpclient.CON_TYPE_LOGSVR
	G_LogClient.SvrID = 1 //日志服忽略此值
	G_LogClient.ConnectToSvr(addr, 10)
}

//向服务器发送签到消息
func SendCheckInMsg(pTcpConn *tcpclient.TCPConn) bool {
	if pTcpConn == nil {
		gamelog.Error("SendCheckInMsg Error: pTcpConn is nullptr!!!")
		return false
	}
	var req msg.MSG_CheckIn_Req
	req.GuildID = -1
	req.PlayerID = int32(appconfig.GameSvrID)
	req.PlayerName = "gamesvr"
	buffer, _ := json.Marshal(&req)
	return pTcpConn.WriteMsg(msg.MSG_CHECK_IN_REQ, 0, buffer)
}

func SendGuildChangeMsg(playerid int32, guilid int32) bool {
	var req msg.MSG_GuildNotify_Req
	req.PlayerID = playerid
	req.NewGuildID = guilid
	buffer, _ := json.Marshal(&req)
	return SendMessageToChat(msg.MSG_GUILD_NOTIFY, 0, buffer)
}

func SendGameSvrNotify(playerid int32, funcid int) bool {
	var req msg.MSG_GameSvr_Notify
	req.FuncID = funcid
	buffer, _ := json.Marshal(&req)
	return SendMessageToClient(playerid, msg.MSG_GAME_SERVER_NOTIFY, 0, buffer)
}

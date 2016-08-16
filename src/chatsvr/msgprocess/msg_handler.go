package msgprocess

import (
	"encoding/binary"
	"encoding/json"
	"gamelog"
	"msg"
	"tcpserver"
)

func Hand_CheckInReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	var req msg.MSG_CheckIn_Req
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_CheckInReq : Unmarshal error!!!!")
		return
	}

	if req.PlayerID == 0 {
		gamelog.Error("Hand_CheckInReq req.PlayerID == 0")
		return
	}

	//收到的是服务器连接
	if (req.PlayerID < 10000) && (req.GuildID == -1) {
		G_GameSvrConn = pTcpConn
		return
	}

	CheckAndClean(req.PlayerID)
	gamelog.Info("message: Hand_CheckInReq id:%d, name:%s", req.PlayerID, req.PlayerName)
	AddTcpConn(req.PlayerID, req.GuildID, req.PlayerName, pTcpConn)
	var response msg.MSG_Chat_Ack
	response.RetCode = msg.RE_SUCCESS
	b, _ := json.Marshal(response)
	pTcpConn.WriteMsg(msg.MSG_CHECK_IN_ACK, b)

	if req.PlayerID >= 10000 {
		SendOnlineNotify(req.PlayerID, true)
	}

	return
}

func Hand_ChatMsgReq(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: Hand_ChatMsgReq")
	if pTcpConn.Data.(*TChatData).PlayerID == 0 {
		gamelog.Error("Hand_ChatMsgReq pTcpConn.PlayerID == 0!!!!")
		return
	}

	var req msg.MSG_Chat_Req
	var response msg.MSG_Chat_Ack
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_ChatMsgReq : Unmarshal error!!!!")
		return
	}

	gamelog.Info("recv: %v", req)

	var chatMsgNotify msg.MSG_Chat_Msg_Notify
	chatMsgNotify.SourceName = req.SourceName
	chatMsgNotify.SourcePlayerID = pTcpConn.Data.(*TChatData).PlayerID
	chatMsgNotify.TargetChannel = req.TargetChannel
	chatMsgNotify.TargetGuildID = req.TargetGuildID
	chatMsgNotify.MsgContent = req.MsgContent
	chatMsgNotify.HeroID = req.HeroID
	chatMsgNotify.Quality = req.Quality

	buff, _ := json.Marshal(chatMsgNotify)

	if req.TargetChannel == msg.MSG_CHANNEL_PLAYER {
		SendMessageByName(req.TargetName, msg.MSG_CHATMSG_NOTIFY, buff)
	} else if req.TargetChannel == msg.MSG_CHANNEL_GUILD {
		SendMessageToGuild(req.TargetGuildID, msg.MSG_CHATMSG_NOTIFY, buff, chatMsgNotify.SourcePlayerID)
	} else if req.TargetChannel == msg.MSG_CHANNEL_WORLD {
		SendMessageToWorld(msg.MSG_CHATMSG_NOTIFY, buff, chatMsgNotify.SourcePlayerID)
	}

	response.RetCode = msg.RE_SUCCESS
	b, _ := json.Marshal(response)
	pTcpConn.WriteMsg(msg.MSG_CHATMSG_ACK, b)
	gamelog.Info("Return: %s", b)
	return
}

func Hand_DisConnect(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	if pTcpConn.Data == nil || pTcpConn.Data.(*TChatData).PlayerID <= 0 {
		return
	}
	SendOnlineNotify(pTcpConn.Data.(*TChatData).PlayerID, false)

	if pTcpConn.Cleaned == false {
		CheckAndClean(pTcpConn.Data.(*TChatData).PlayerID)
	}

	return
}

func Hand_Game_To_Client(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	if len(pdata) < 6 {
		gamelog.Error("Hand_Game_To_Client : message data errror!!!!")
		return
	}
	playerid := int(binary.LittleEndian.Uint64(pdata[:4]))
	msgid := int16(binary.LittleEndian.Uint16(pdata[4:6]))

	if playerid == 0 {
		SendMessageToWorld(msgid, pdata[6:], 0)
	} else {
		SendMessageByID(playerid, msgid, pdata[6:])
	}

	return
}

func Hand_GuildChange_Notify(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	var req msg.MSG_GuildNotify_Req
	if json.Unmarshal(pdata, &req) != nil {
		gamelog.Error("Hand_GuildChange_Notify : Unmarshal error!!!!")
		return
	}

	ChangeConnGuild(req.PlayerID, req.NewGuildID)

	return
}

func Hand_HeartBeat(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_HEART_BEAT")
	var req msg.MSG_HeartBeat_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_HeartBeat : Message Reader Error!!!!")
		return
	}

	return
}
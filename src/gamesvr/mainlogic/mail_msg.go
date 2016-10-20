package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
	"utility"
)

//收到所有的邮件请求
//消息:/receive_all_mails
type MSG_ReceiveAllMails_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_ReceiveAllMails_Ack struct {
	RetCode  int         //返回码
	MailList []TMailInfo //邮件列表
}

//! 玩家请求邮件信息
func Hand_ReceiveAllMails(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req MSG_ReceiveAllMails_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ReceiveAllMails unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_ReceiveAllMails_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.MailList = player.MailMoudle.MailList
	response.RetCode = msg.RE_SUCCESS
	player.MailMoudle.MailList = []TMailInfo{}
	player.MailMoudle.DB_ClearAllMails()
}

//! 发送私人邮件
func Hand_SendPrivateMail(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SendPrivateMail_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SendPrivateMail unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SendPrivateMail_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	var mail TMailInfo
	mail.TextType = Text_Private_mail
	mail.MailTime = utility.GetCurTime()
	mail.MailParams = make([]string, 2)
	mail.MailParams[0] = req.Title
	mail.MailParams[1] = req.Content
	SendMailToPlayer(req.TargetID, &mail)
	response.RetCode = msg.RE_SUCCESS

}

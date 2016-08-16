package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

func Hand_GetTitle(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTitle_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTitle Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetTitle_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.TitleModule.CheckTitleDeadLine()

	for _, v := range player.TitleModule.TitleLst {
		var title msg.TitleInfo
		title.TitleID = v.TitleID
		title.EndTime = v.EndTime
		title.Status = v.Status
		response.TitleLst = append(response.TitleLst, title)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求激活称号
func Hand_ActivateTitle(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ActivateTitle_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTitle Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_ActivateTitle_Ack

	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.TitleModule.CheckTitleDeadLine()

	isExist := false
	for _, v := range player.TitleModule.TitleLst {
		if v.TitleID == req.TitleID {
			isExist = true
		}
	}

	if isExist == false {
		gamelog.Error("Hand_GetTitle Error: Player have not this title id: %d", req.TitleID)
		response.RetCode = msg.RE_NOT_HAVE_TITLE
		return
	}

	titleData := gamedata.GetTitleInfo(req.TitleID)
	if player.BagMoudle.IsItemEnough(titleData.CostItemID, 1) == false {
		gamelog.Error("Hand_ActivateTitle Error: Item not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		return
	}

	player.BagMoudle.RemoveNormalItem(titleData.CostItemID, 1)
	player.TitleModule.AddTitle(req.TitleID)

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求佩戴称号
func Hand_EquipTitle(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_EquiTitle_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTitle Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_EquiTitle_Ack

	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.TitleModule.CheckTitleDeadLine()

	isExist := false
	for _, v := range player.TitleModule.TitleLst {
		if v.TitleID == req.TitleID {
			isExist = true
		}
	}

	if isExist == false {
		gamelog.Error("Hand_GetTitle Error: Player have not this title id: %d", req.TitleID)
		response.RetCode = msg.RE_NOT_HAVE_TITLE
		return
	}

	player.TitleModule.EquiTitle(req.TitleID)

	response.RetCode = msg.RE_SUCCESS
}

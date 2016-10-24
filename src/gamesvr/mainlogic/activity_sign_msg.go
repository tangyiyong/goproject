package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 获取签到状态
func Hand_GetSignData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSignData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSignData Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetSignData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.Sign.ActivityID) == false {
		gamelog.Error("Hand_GetSignData Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	response.IsSign = player.ActivityModule.Sign.IsSign
	if response.IsSign == true {
		//! 已签到
		response.SignDay = player.ActivityModule.Sign.SignDay
	} else {
		//! 未签到
		response.SignDay = player.ActivityModule.Sign.SignDay + 1
	}

	signAwardCount := gamedata.GetSignAwardCount()
	if response.SignDay > signAwardCount {
		response.SignDay = signAwardCount
	}

	signAward := gamedata.GetSignData(response.SignDay)

	response.SignIndex = signAward.Type

	length := len(player.ActivityModule.Sign.SignPlusAward)
	for i := 0; i < length; i++ {
		var item msg.MSG_ItemData
		item.ID = player.ActivityModule.Sign.SignPlusAward[i].ItemID
		item.Num = player.ActivityModule.Sign.SignPlusAward[i].ItemNum
		response.SignPlusAward = append(response.SignPlusAward, item)
	}

	response.SignPlusStatus = player.ActivityModule.Sign.SignPlusStatus

	response.RetCode = msg.RE_SUCCESS
}

//! 进行日常签到
func Hand_DailySign(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_DailySign_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_DailySign Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_DailySign_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	player.ActivityModule.CheckReset()

	//! 检测时间
	if player.ActivityModule.Sign.IsSign == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		gamelog.Error("Hand_DailySign error: Aleady Sign. playerID: %v", player.playerid)
		return
	}

	//! 签到
	ret := false
	ret, response.ItemID, response.ItemNum = player.ActivityModule.Sign.Sign()
	if ret == false {
		gamelog.Error("Hand_DailySign error: Sign fail. playerID: %v", player.playerid)
		return
	}

	response.RetCode = msg.RE_SUCCESS
	response.AwardType = G_GlobalVariables.GetActivityAwardType(player.ActivityModule.Sign.ActivityID)
}

//! 领取豪华签到
func Hand_SignPlus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_PlusSign_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SignPlus Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_PlusSign_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	player.ActivityModule.CheckReset()

	//! 检测豪华签到状态
	ret := (player.ActivityModule.Sign.SignPlusStatus == SignPlus_Can_Receive)
	if ret == false {
		gamelog.Error("Hand_SignPlus error: Users do not recharge or don't enough.")
		response.RetCode = msg.RE_CAN_NOT_SIGN_PLUS
		return
	}

	//! 检测豪华签到时间
	if player.ActivityModule.Sign.IsSignPlus == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		gamelog.Error("Can't sign plus. playerID: %v", player.playerid)
		return
	}

	//! 进行豪华签到
	awardLst := player.ActivityModule.Sign.SignPlus()

	for _, v := range awardLst {
		award := msg.MSG_ItemData{v.ItemID, v.ItemNum}
		response.AwardInfo = append(response.AwardInfo, award)
	}
	response.RetCode = msg.RE_SUCCESS
}

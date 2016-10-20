package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求领取月卡
func Hand_ReceiveMonthCard(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_ReceiveMonthCard_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ReceiveMonthCard Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ReceiveMonthCard_Ack
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

	player.ActivityModule.CheckReset()

	pMonthCard := gamedata.GetChargeItem(req.CardID)
	if pMonthCard == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ReceiveMonthCard Error : Invalid Cardid :%d", req.CardID)
		return
	}

	if player.ActivityModule.MonthCard.CardDays[pMonthCard.ID-1] <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ReceiveMonthCard Error : can receive month card days: %v  req.ID: %d", player.ActivityModule.MonthCard.CardDays[pMonthCard.ID-1], req.CardID)
		return
	}

	if player.ActivityModule.MonthCard.CardStatus[pMonthCard.ID-1] == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		gamelog.Error("Hand_ReceiveMonthCard Error : can receive month card  status: %v  req.ID: %d", player.ActivityModule.MonthCard.CardStatus[pMonthCard.ID], req.CardID)
		return
	}

	player.ActivityModule.MonthCard.CardStatus[pMonthCard.ID-1] = true
	player.ActivityModule.MonthCard.DB_UpdateCardStatus()

	player.RoleMoudle.AddMoney(gamedata.ChargeMoneyID, pMonthCard.ExtraAward)

	//! 返回领取状态
	response.RetCode = msg.RE_SUCCESS
	response.CardID = req.CardID

	response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{gamedata.ChargeMoneyID, pMonthCard.ExtraAward})
}

//! 查询月卡天数
func Hand_QueryActivityMonthCardDays(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryActivity_MonthCard_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryActivityTotalRechargeInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QueryActivity_MonthCard_Ack
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

	for i := 0; i < 2; i++ {
		response.Days = append(response.Days, player.ActivityModule.MonthCard.CardDays[i])
		response.Status = append(response.Status, player.ActivityModule.MonthCard.CardStatus[i])
	}

	response.RetCode = msg.RE_SUCCESS
}

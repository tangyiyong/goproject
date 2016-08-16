package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 查询月基金状态
func Hand_GetMonthFundStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetMonthFundStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMonthFundStatus Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetMonthFundStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.MonthFund.ActivityID) == false {
		gamelog.Error("Hand_GetMonthFundStatus Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	_, countDown := G_GlobalVariables.IsActivityTime(player.ActivityModule.MonthFund.ActivityID)

	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.MonthFund.ActivityID)
	awardCount := gamedata.GetMonthFundAwardCount(awardType)

	response.CountDown = countDown
	response.MoneyID, response.MoneyNum = gamedata.MonthFundCostMoneyID, gamedata.MonthFundCostMoneyNum
	response.Day = player.ActivityModule.MonthFund.Day
	response.IsReceived = player.ActivityModule.MonthFund.AwardMark.Get(uint(awardCount - response.Day + 1))
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求月基金领取奖励
func Hand_ReceiveMonthFund(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ReceiveMonthFund_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ReceiveMonthFund Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_ReceiveMonthFund_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.ActivityModule.CheckReset()

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.MonthFund.ActivityID) == false {
		gamelog.Error("Hand_ReceiveMonthFund Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.MonthFund.ActivityID)
	awardCount := gamedata.GetMonthFundAwardCount(awardType)

	day := player.ActivityModule.MonthFund.Day
	award := gamedata.GetMonthFundAward(awardType, awardCount-day+1)

	if player.ActivityModule.MonthFund.AwardMark.Get(uint(awardCount-day+1)) == true {
		gamelog.Error("Hand_ReceiveMonthFund Error: Aleady receive award")
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	player.ActivityModule.MonthFund.AwardMark.Set(uint(awardCount - day + 1))
	go player.ActivityModule.MonthFund.DB_UpdateAwardMark()

	player.BagMoudle.AddAwardItem(award.ItemID, award.ItemNum)
	response.AwardLst = append(response.AwardLst, msg.MSG_ItemData{award.ItemID, award.ItemNum})

	response.RetCode = msg.RE_SUCCESS
}

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家获取欢庆佳节任务信息
func Hand_GetFestivalTask(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFestivalTask_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFestivalTask Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetFestivalTask_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.Festival.ActivityID) == false {
		gamelog.Error("Hand_GetFestivalTask Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	activityInfo := gamedata.GetActivityInfo(player.ActivityModule.Festival.ActivityID)

	taskLst := player.ActivityModule.Festival.TaskLst
	length := len(taskLst)
	for i := 0; i < length; i++ {
		var task msg.MSG_FestivalTask
		task.ID = taskLst[i].ID

		taskInfo := gamedata.GetFestivalTaskInfo(activityInfo.AwardType, task.ID)
		task.Need = taskInfo.Need
		task.CurCount = taskLst[i].Count
		task.Page = taskInfo.Page
		task.PageName = taskInfo.PageName
		task.TaskType = taskLst[i].TaskType
		task.Goto = taskInfo.Goto
		task.Award = taskInfo.Award
		task.Desc = taskInfo.Desc
		task.Status = taskLst[i].Status

		response.TaskLst = append(response.TaskLst, task)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家获取欢庆佳节活动兑换信息
func Hand_GetFestivalExchangeInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFestivalExchange_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFestivalExchangeInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetFestivalExchange_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.Festival.ActivityID) == false {
		gamelog.Error("Hand_GetFestivalExchangeInfo Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	activityInfo := gamedata.GetActivityInfo(player.ActivityModule.Festival.ActivityID)
	exchangeLst := gamedata.GetExchangeInfoLst(activityInfo.AwardType)

	length := len(exchangeLst)
	for i := 0; i < length; i++ {
		var exchange msg.MSG_FestivalExchange
		exchange.ID = exchangeLst[i].ID
		exchange.Award = exchangeLst[i].Award
		exchange.NeedItemID1 = exchangeLst[i].NeedItemID1
		exchange.NeedItemID2 = exchangeLst[i].NeedItemID2
		exchange.NeedItemNum1 = exchangeLst[i].NeedItemNum1
		exchange.NeedItemNum2 = exchangeLst[i].NeedItemNum2
		exchange.Times = exchangeLst[i].ExchangeTimes
		response.ExchangeRecordLst = append(response.ExchangeRecordLst, exchange)
	}

	//! 扣除已兑换次数
	length = len(response.ExchangeRecordLst)
	for _, v := range player.ActivityModule.Festival.ExchangeLst {
		for i := 0; i < length; i++ {
			if response.ExchangeRecordLst[i].ID == v.ID {
				response.ExchangeRecordLst[i].Times -= v.Times
			}
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求领取欢庆佳节任务奖励
func Hand_GetFestivalTaskAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetFestivalTaskAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetFestivalTaskAward Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetFestivalTaskAward_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.Festival.ActivityID) == false {
		gamelog.Error("Hand_GetFestivalTaskAward Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 获取任务信息
	taskInfo, index := player.ActivityModule.Festival.GetTaskInfo(req.ID)
	if taskInfo == nil {
		gamelog.Error("Hand_GetFestivalTaskAward Error: Invalid taskID %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if taskInfo.Status != 1 {
		gamelog.Error("Hand_GetFestivalTaskAward Error: Aleady receive this task award")
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	taskInfo.Status = 2
	go player.ActivityModule.Festival.DB_UpdateTaskStatus(index)

	//! 发放奖励
	awardLst := gamedata.GetItemsFromAwardID(taskInfo.Award)
	response.AwardLst = []msg.MSG_ItemData{}
	for _, v := range awardLst {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.AwardLst = append(response.AwardLst, item)
	}

	player.BagMoudle.AddAwardItems(awardLst)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家兑换奖励
func Hand_ExchangeFestivalAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ExchangeFestivalAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ExchangeFestivalAward Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_ExchangeFestivalAward_Ack
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

	exchangeInfo := gamedata.GetExchangeInfoFromID(req.ID)
	if exchangeInfo == nil {
		gamelog.Error("Hand_ExchangeFestivalAward Error: Invalid Exchange id %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	exchange, index := player.ActivityModule.Festival.GetExchangeInfo(req.ID)
	if exchange.Times+1 > exchangeInfo.ExchangeTimes {
		//! 超出兑换次数
		gamelog.Error("Hand_ExchangeFestivalAward Error: Over exchange times")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	if exchangeInfo.NeedItemID1 != 0 {
		if player.BagMoudle.IsItemEnough(exchangeInfo.NeedItemID1, exchangeInfo.NeedItemNum1) == false {
			gamelog.Error("Hand_ExchangeFestivalAward Error: Item not enough %d", exchangeInfo.NeedItemID1)
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}

		if exchangeInfo.NeedItemID2 != 0 {
			if player.BagMoudle.IsItemEnough(exchangeInfo.NeedItemID2, exchangeInfo.NeedItemNum2) == false {
				gamelog.Error("Hand_ExchangeFestivalAward Error: Item not enough %d", exchangeInfo.NeedItemID1)
				response.RetCode = msg.RE_NOT_ENOUGH_ITEM
				return
			}
			player.BagMoudle.RemoveNormalItem(exchangeInfo.NeedItemID2, exchangeInfo.NeedItemNum2)
		}
	}

	player.BagMoudle.RemoveNormalItem(exchangeInfo.NeedItemID1, exchangeInfo.NeedItemNum1)

	//! 增加兑换次数
	exchange.Times += 1
	go player.ActivityModule.Festival.DB_UpdateExchangeTimes(index, exchange.Times)

	//! 给予兑换奖励
	awardLst := gamedata.GetItemsFromAwardID(exchangeInfo.Award)
	response.AwardLst = []msg.MSG_ItemData{}
	for _, v := range awardLst {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.AwardLst = append(response.AwardLst, item)
	}

	response.RetCode = msg.RE_SUCCESS
}

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

	taskLst := player.ActivityModule.Festival.TaskLst
	length := len(taskLst)
	for i := 0; i < length; i++ {
		var task msg.MSG_FestivalTask
		task.ID = taskLst[i].ID
		task.CurCount = taskLst[i].Count
		task.Status = taskLst[i].Status

		response.TaskLst = append(response.TaskLst, task)
	}

	activityInfo := gamedata.GetActivityInfo(player.ActivityModule.Festival.ActivityID)
	exchangeLst := gamedata.GetExchangeInfoLst(activityInfo.AwardType)

	length = len(exchangeLst)
	for i := 0; i < length; i++ {
		var exchange msg.MSG_FestivalExchange
		exchange.ID = exchangeLst[i].ID
		exchange.Award = exchangeLst[i].Award
		exchange.NeedItemID = exchangeLst[i].NeedItemID
		exchange.NeedItemNum = exchangeLst[i].NeedItemNum
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

	//! 购买记录
	for _, v := range player.ActivityModule.Festival.BuyLst {
		var saleInfo msg.MSG_FestivalSale
		saleInfo.ID = v.ID
		saleInfo.Times = v.BuyTimes
		response.BuyLst = append(response.BuyLst, saleInfo)
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
	player.ActivityModule.Festival.DB_UpdateTaskStatus(index)

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

	if player.BagMoudle.IsItemEnough(exchangeInfo.NeedItemID, exchangeInfo.NeedItemNum) == false {
		gamelog.Error("Hand_ExchangeFestivalAward Error: Item not enough %d", exchangeInfo.NeedItemID)
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		return
	}

	player.BagMoudle.RemoveNormalItem(exchangeInfo.NeedItemID, exchangeInfo.NeedItemNum)

	//! 增加兑换次数
	exchange.Times += 1
	player.ActivityModule.Festival.DB_UpdateExchangeTimes(index, exchange.Times)

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

//! 玩家请求领取欢庆佳节任务奖励
func Hand_BuyFestivalSaleItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyFestivalSale_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyFestivalSaleItem Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyFestivalSale_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.Festival.ActivityID) == false {
		gamelog.Error("Hand_BuyFestivalSaleItem Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 获取商品信息
	itemInfo := gamedata.GetFestivalItemInfo(req.ID)
	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.Festival.ActivityID)
	if itemInfo == nil || awardType != itemInfo.AwardType {
		gamelog.Error("Hand_BuyFestivalSaleItem Error: Invalid Param  ID %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.MoneyID, itemInfo.MoneyNum) == false {
		gamelog.Error("Hand_BuyFestivalSaleItem Error: Money not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 检测次数是否足够
	saleInfo := player.ActivityModule.Festival.GetFestivalSaleInfo(req.ID)
	if saleInfo.BuyTimes >= itemInfo.BuyTimes {
		gamelog.Error("Hand_BuyFestivalSaleItem Error: Buy times not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(itemInfo.MoneyID, itemInfo.MoneyNum)
	response.CostMoneyID, response.CostMoneyNum = itemInfo.MoneyID, itemInfo.MoneyNum

	//! 发放商品
	player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum)
	response.ItemID, response.ItemNum = itemInfo.ItemID, itemInfo.ItemNum

	//! 增加次数
	saleInfo.BuyTimes++
	player.ActivityModule.Festival.DB_UpdateBuyTimes(saleInfo)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

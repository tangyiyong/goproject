package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家查询限时特惠商品信息
func Hand_GetLimitSaleItemInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetLimitSaleInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetLimitSaleItemInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetLimitSaleInfo_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.LimitSale.ActivityID) == false {
		gamelog.Error("Hand_GetLimitSaleItemInfo Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	//! 返回数据
	response.Score = player.ActivityModule.LimitSale.Score

	for i := 0; i < len(player.ActivityModule.LimitSale.ItemLst); i++ {
		item := player.ActivityModule.LimitSale.ItemLst[i]
		response.ItemLst = append(response.ItemLst, msg.MSG_LimitSaleItemInfo{item.ID, item.Status})
	}

	response.AwardMark = int(player.ActivityModule.LimitSale.AwardMark)
	response.SaleNum = G_GlobalVariables.LimitSaleNum
	response.RetCode = msg.RE_SUCCESS
	response.DiscountChargeID = player.ActivityModule.LimitSale.DiscountChargeID
}

//! 玩家购买限时特惠商品
func Hand_BuyLimitSaleItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyLimitSaleItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyLimitSaleItem Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyLimitSaleItem_Ack
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

	//! 检测参数合法性
	if req.Index > len(player.ActivityModule.LimitSale.ItemLst) ||
		req.Index <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_BuyLimitSaleItem Fail: Invalid Index %d ", req.Index)
		return
	}

	item := &player.ActivityModule.LimitSale.ItemLst[req.Index-1]
	if item.Status == true {
		gamelog.Error("Hand_BuyLimitSaleItem Fail: Aleady buy Index %d", req.Index)
		response.RetCode = msg.RE_ALEADY_BUY
		return
	}

	itemInfo := gamedata.GetLimitSaleItemInfo(item.ID)

	//! 判断货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.MoneyID, itemInfo.MoneyNum) == false {
		gamelog.Error("Hand_BuylimitSaleItem Fail: Not enough money")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(itemInfo.MoneyID, itemInfo.MoneyNum)

	//! 增加积分
	player.ActivityModule.LimitSale.Score += itemInfo.Score
	if player.ActivityModule.LimitSale.Score >= 100 {
		if player.ActivityModule.LimitSale.DiscountChargeID == 0 {
			player.ActivityModule.LimitSale.RandDiscountCharge()
		}

		player.ActivityModule.LimitSale.Score = 100
	}
	player.ActivityModule.LimitSale.DB_UpdateScore()

	//! 奖励物品
	player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum)
	response.AwardLst = append(response.AwardLst, msg.MSG_ItemData{itemInfo.ItemID, itemInfo.ItemNum})

	//! 改变状态
	item.Status = true
	player.ActivityModule.LimitSale.DB_UpdateStatus(req.Index - 1)

	//! 增加购买人次
	G_GlobalVariables.LimitSaleNum += 1
	G_GlobalVariables.DB_UpdateLimitSaleNum()

	response.DiscountChargeID = player.ActivityModule.LimitSale.DiscountChargeID
	response.Score = player.ActivityModule.LimitSale.Score
	response.RetCode = msg.RE_SUCCESS
}

//! 领取全民奖励
func Hand_GetLimitSaleAllAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetLimitSale_AllAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetLimitSaleAllAward Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetLimitSale_AllAward_Ack
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

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	player.ActivityModule.CheckReset()

	if player.ActivityModule.LimitSale.AwardMark.Get(req.ID) == true {
		gamelog.Error("Hand_GetLimitSaleAllAward Error: Aleady received. ID: %d", req.ID)
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	awardInfo := gamedata.GetLimitSaleAllAwardInfo(req.ID)
	if awardInfo == nil {
		gamelog.Error("Hand_GetLimitSaleAllAward Error: Invalid Param ID: %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if awardInfo.NeedNum > G_GlobalVariables.LimitSaleNum {
		gamelog.Error("Hand_GetLimitSaleAllAward Error: Not enough buy num")
		response.RetCode = msg.RE_NOT_ENOUGH_NUMBER
		return
	}

	awardLst := gamedata.GetItemsFromAwardID(awardInfo.Award)
	player.BagMoudle.AddAwardItems(awardLst)
	for _, v := range awardLst {
		var awardItem msg.MSG_ItemData
		awardItem.ID = v.ItemID
		awardItem.Num = v.ItemNum
		response.AwardLst = append(response.AwardLst, awardItem)
	}

	player.ActivityModule.LimitSale.AwardMark.Set(req.ID)
	player.ActivityModule.LimitSale.DB_UpdateAwardMark()

	response.AwardMark = int(player.ActivityModule.LimitSale.AwardMark)
	response.RetCode = msg.RE_SUCCESS
}

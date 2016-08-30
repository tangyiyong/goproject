package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 获取团购信息
func Hand_GetGroupPurchaseInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message:%s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetGroupPurchaseInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetGroupPurchaseInfo Error: Unmarshal fail")
		return
	}

	var response msg.MSG_GetGroupPurchaseInfo_Ack
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

	player.ActivityModule.CheckReset()

	response.Score = player.ActivityModule.GroupPurchase.Score

	for _, v := range player.ActivityModule.GroupPurchase.PurchaseCostLst {
		var itemInfo msg.MSG_GroupPurchase
		itemInfo.ItemID = v.ItemID

		groupItemData := gamedata.GetGroupPurchaseItemInfo(v.ItemID, player.ActivityModule.GroupPurchase.ActivityID)
		itemInfo.CanBuyNum = groupItemData.BuyTimes - v.Times

		groupItemInfo, _ := G_GlobalVariables.GetGroupPurchaseItemInfo(v.ItemID)
		itemInfo.SaleNum = groupItemInfo.SaleNum
		response.ItemInfo = append(response.ItemInfo, itemInfo)
	}

	for _, v := range G_GlobalVariables.ActivityLst {
		if v.ActivityID == player.ActivityModule.GroupPurchase.ActivityID {
			activityInfo := gamedata.GetActivityInfo(v.ActivityID)
			response.EndTime = v.endTime - int64(activityInfo.AwardTime*60*60*24)
			response.AwardTime = v.endTime
		}
	}

	response.AwardType = G_GlobalVariables.GetActivityAwardType(player.ActivityModule.GroupPurchase.ActivityID)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买团购
func Hand_BuyGroupPurchaseItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message:%s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_BuyGroupPurchase_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyGroupPurchaseItem Error: Unmarshal fail")
		return
	}

	var response msg.MSG_BuyGroupPurchase_Ack
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

	player.ActivityModule.CheckReset()

	//! 检查参数合法性
	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.GroupPurchase.ActivityID)
	itemData := gamedata.GetGroupPurchaseItemInfo(req.ItemID, awardType)

	//! 检测团购次数是否足够
	costInfo, costIndex := player.ActivityModule.GroupPurchase.GetGroupItemInfo(req.ItemID)
	canBuyNum := itemData.BuyTimes - costInfo.Times

	if canBuyNum <= 0 {
		gamelog.Error("Hand_BuyGroupPurchaseItem Error: Buy times not enough")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 获取售价信息
	saleInfo, _ := G_GlobalVariables.GetGroupPurchaseItemInfo(req.ItemID)
	itemSaleInfo := gamedata.GetGroupPurchaseItemInfoFromSale(req.ItemID, awardType, saleInfo.SaleNum)

	//! 检测玩家货币是否足够
	useItemNum := 0
	useMoneyNum := itemSaleInfo.MoneyNum
	if itemSaleInfo.MoneyNum > player.RoleMoudle.Moneys[gamedata.GroupPurchaseCostMoneyID-1] {
		//! 使用团购券来湊
		useItemNum = itemSaleInfo.MoneyNum - player.RoleMoudle.Moneys[gamedata.GroupPurchaseCostMoneyID-1]
		if useItemNum > itemData.UseItemMax {
			gamelog.Error("Hand_BuyGroupPurchaseItem Error: Money not enough")
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}

		useMoneyNum -= useItemNum
	}

	if useItemNum != 0 {
		//! 扣除团购券
		player.BagMoudle.RemoveNormalItem(gamedata.GroupPurchaseCostItemID, useItemNum)
		response.CostItemID = gamedata.GroupPurchaseCostItemID
		response.CostItemNum = useItemNum
	}

	player.RoleMoudle.CostMoney(gamedata.GroupPurchaseCostMoneyID, useMoneyNum)
	response.CostMoneyID = gamedata.GroupPurchaseCostMoneyID
	response.CostMoneyNum = useMoneyNum

	//! 发放物品
	player.BagMoudle.AddAwardItem(itemData.ItemID, itemData.ItemNum)

	//! 增加积分
	response.Score = useItemNum + useMoneyNum*10
	player.ActivityModule.GroupPurchase.Score += response.Score
	response.Score = player.ActivityModule.GroupPurchase.Score
	go player.ActivityModule.GroupPurchase.DB_SaveScore()

	//! 增加购买记录
	costInfo.MoneyNum += useMoneyNum
	costInfo.Times += 1
	go player.ActivityModule.GroupPurchase.DB_UpdatePurchaseCostInfo(costIndex)

	//! 增加总购买记录
	G_GlobalVariables.AddGroupPurchaseRecord(req.ItemID, 1)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询积分奖励
func Hand_QueryGroupPurchaseScoreAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message:%s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetGroupPurchaseScore_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetGroupPurchaseScoreAward Error: Unmarshal fail")
		return
	}

	var response msg.MSG_GetGroupPurchaseScore_Ack
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

	player.ActivityModule.CheckReset()

	response.ScoreAwardMark = player.ActivityModule.GroupPurchase.ScoreAwardMark
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求积分奖励
func Hand_GetGroupPurchaseScoreAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message:%s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetGroupScoreAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetGroupPurchaseScoreAward Error: Unmarshal fail")
		return
	}

	var response msg.MSG_GetGroupScoreAward_Ack
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

	player.ActivityModule.CheckReset()

	//! 检查玩家是否已经领取
	if player.ActivityModule.GroupPurchase.ScoreAwardMark.IsExist(req.ID) > 0 {
		gamelog.Error("Hand_GetGroupPurchaseScoreAward Error: Player aleady received")
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查积分是否足够
	scoreInfo := gamedata.GetGroupPurchaseScoreAward(req.ID)
	if scoreInfo == nil {
		gamelog.Error("Hand_GetGroupPurchaseScoreAward Error: Invalid Param")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if player.ActivityModule.GroupPurchase.Score < scoreInfo.NeedScore {
		gamelog.Error("Hand_GetGroupPurchaseScoreAward Error: Score not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_SCORE
		return
	}

	//! 领取奖励
	player.BagMoudle.AddAwardItem(scoreInfo.ItemID, scoreInfo.ItemNum)

	//! 增加记录
	player.ActivityModule.GroupPurchase.ScoreAwardMark.Add(req.ID)
	player.ActivityModule.GroupPurchase.DB_AddScoreAward(req.ID)
}

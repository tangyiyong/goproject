package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 获取开服基金状态
func Hand_GetOpenFundStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetOpenFundStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetOpenFundStatus : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetOpenFundStatus_Ack
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

	response.BuyNum = G_BuyFundNum
	response.IsBuy = player.ActivityModule.OpenFund.IsBuyFund
	response.FundCountMark = int(player.ActivityModule.OpenFund.FundCountMark)
	response.FundLevelMark = int(player.ActivityModule.OpenFund.FundLevelMark)

	response.CostMoneyID = gamedata.OpenFundPriceID
	response.CostMoneyNum = gamedata.OpenFundPriceNum

	for _, v := range gamedata.GT_OpenFundLst[gamedata.OpenFund_Level] {
		if player.ActivityModule.OpenFund.FundLevelMark.Get(v.ID) != true {
			awardLst := gamedata.GetItemsFromAwardID(v.Award)
			for _, m := range awardLst {
				response.ReceiveMoney += m.ItemNum
			}
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 购买开服基金
func Hand_BuyOpenFund(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyFund_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyOpenFund : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyFund_Ack
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

	//! 检查是否已经购买
	if player.ActivityModule.OpenFund.IsBuyFund == true {
		gamelog.Error("Hand_BuyOpenFund error: repeat purchase. PlayerID: %v", req.PlayerID)
		response.RetCode = msg.RE_ALEADY_BUY
		return
	}

	//! 检查玩家货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.OpenFundPriceID, gamedata.OpenFundPriceNum) == false {
		gamelog.Error("Hand_BuyOpenFund error: money is not enough. PlayerID: %v", req.PlayerID)
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除货币,更改标记
	player.RoleMoudle.CostMoney(gamedata.OpenFundPriceID, gamedata.OpenFundPriceNum)

	player.ActivityModule.OpenFund.IsBuyFund = true
	player.ActivityModule.OpenFund.UpdateBuyFundMark()

	//! 购买基金人数+1
	G_BuyFundNum += 1

	response.FundCountMark = int(player.ActivityModule.OpenFund.FundCountMark)
	response.FundLevelMark = int(player.ActivityModule.OpenFund.FundLevelMark)

	for _, v := range gamedata.GT_OpenFundLst[gamedata.OpenFund_Level] {
		if player.ActivityModule.OpenFund.FundLevelMark.Get(v.ID) != true {
			awardLst := gamedata.GetItemsFromAwardID(v.Award)
			for _, m := range awardLst {
				response.ReceiveMoney += m.ItemNum
			}
		}
	}

	response.BuyNum = G_BuyFundNum
	response.RetCode = msg.RE_SUCCESS
}

//! 领取开服基金全民奖励
func Hand_GetOpenFundAllAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ReceiveFundAllAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetOpenFundAllAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_ReceiveFundAllAward_Ack
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

	//! 检测是否已够买基金
	if player.ActivityModule.OpenFund.IsBuyFund == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取领取奖励类型
	awardInfo := gamedata.GetOpenFundInfo(gamedata.OpenFund_BuyNum, req.ID)
	if awardInfo == nil {
		gamelog.Error("Hand_GetOpenFundAllAward error: GetOpenFundInfo nil")
		return
	}

	//! 判断人数是否达标
	if G_BuyFundNum < awardInfo.Count {
		gamelog.Error("Hand_GetOpenFundAllAward error: BuyNum: %d  NeedNum: %d  ID: %d", G_BuyFundNum, awardInfo.Count, req.ID)
		response.RetCode = msg.RE_NOT_ENOUGH_NUMBER
		return
	}

	//! 判断是否已经领取
	if player.ActivityModule.OpenFund.FundCountMark.Get(req.ID) == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 领取奖励
	awarditems := gamedata.GetItemsFromAwardID(awardInfo.Award)
	player.BagMoudle.AddAwardItems(awarditems)

	for _, v := range awarditems {
		award := msg.MSG_ItemData{v.ItemID, v.ItemNum}
		response.AwardItem = append(response.AwardItem, award)
	}

	//! 改变标记
	player.ActivityModule.OpenFund.FundCountMark.Set(req.ID)
	player.ActivityModule.OpenFund.UpdateFundCountMark()

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.Index = req.ID
	response.AwardMark = int(player.ActivityModule.OpenFund.FundCountMark)
}

//! 领取开服基金等级奖励
func Hand_GetOpenFundLevelAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ReceiveFundLevelAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetOpenFundLevelAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_ReceiveFundLevelAward_Ack
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

	//! 检测是否已够买基金
	if player.ActivityModule.OpenFund.IsBuyFund == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取领取奖励类型
	awardInfo := gamedata.GetOpenFundInfo(gamedata.OpenFund_Level, req.ID)
	if awardInfo == nil {
		gamelog.Error("Hand_GetOpenFundLevelAward error: GetOpenFundInfo nil")
		return
	}

	//! 判断等级是否达标
	if player.GetLevel() < awardInfo.Count {
		gamelog.Error("Hand_GetOpenFundLevelAward error: Level: %d  NeedLevel: %d", player.GetLevel(), awardInfo.Count)
		response.RetCode = msg.RE_NOT_ENOUGH_NUMBER
		return
	}

	//! 判断是否已经领取
	if player.ActivityModule.OpenFund.FundLevelMark.Get(req.ID) == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 领取奖励
	awarditems := gamedata.GetItemsFromAwardID(awardInfo.Award)
	player.BagMoudle.AddAwardItems(awarditems)

	for _, v := range awarditems {
		award := msg.MSG_ItemData{v.ItemID, v.ItemNum}
		response.AwardItem = append(response.AwardItem, award)
	}

	//! 改变标记
	player.ActivityModule.OpenFund.FundLevelMark.Set(req.ID)
	player.ActivityModule.OpenFund.UpdateFundLevelMark()

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.Index = req.ID
	response.AwardMark = int(player.ActivityModule.OpenFund.FundLevelMark)
}

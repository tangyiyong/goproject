package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求等级礼包信息
func Hand_GetRankGiftInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetRankGiftInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetRankGiftInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetRankGiftInfo_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.RankGift.ActivityID) == false {
		gamelog.Error("Hand_GetRankGiftInfo Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	giftLst := player.ActivityModule.RankGift.GiftLst
	length := len(player.ActivityModule.RankGift.GiftLst)
	for i := 0; i < length; i++ {
		var gift msg.MSG_RankGiftInfo
		gift.ID = giftLst[i].GiftID
		gift.BuyTimes = giftLst[i].BuyTimes
		response.GiftLst = append(response.GiftLst, gift)
	}

	//! 因为查看了购买列表, 将新商品标记置位false
	player.ActivityModule.RankGift.IsHaveNewItem = false
	player.ActivityModule.RankGift.DB_UpdateNewItemMark()

	response.Rank = player.ArenaModule.HistoryRank
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买等级礼包
func Hand_BuyRankGift(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyRankGift_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyRankGift Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuyRankGift_Ack
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

	if G_GlobalVariables.IsActivityOpen(player.ActivityModule.RankGift.ActivityID) == false {
		gamelog.Error("Hand_BuyRankGift Error: Activity not open")
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(player.ActivityModule.RankGift.ActivityID)

	//! 获取等级礼包信息
	rankGift := player.ActivityModule.RankGift.GetRankGiftInfo(req.GiftID)
	if rankGift == nil {
		gamelog.Error("Hand_BuyRankGift Error: Not find gift info. id: %d", req.GiftID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if rankGift.BuyTimes <= 0 {
		gamelog.Error("Hand_BuyRankGift Error: BuyTimes not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	rankGiftInfo := gamedata.GetLevelGiftInfo(awardType, player.ActivityModule.RankGift.ActivityID)
	if rankGiftInfo == nil {
		gamelog.Error("Hand_BuyRankGift Error: GetRankGiftInfo nil")
		response.RetCode = msg.RE_UNKNOWN_ERR
		return
	}

	//! 检测货币是否足够
	if rankGiftInfo.MoneyID != 0 {
		if player.RoleMoudle.CheckMoneyEnough(rankGiftInfo.MoneyID, rankGiftInfo.MoneyNum) == false {
			gamelog.Error("Hand_BuyRankGift Error: Not enough money")
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}

		//! 扣除货币
		player.RoleMoudle.CostMoney(rankGiftInfo.MoneyID, rankGiftInfo.MoneyNum)
		response.CostMoneyID = rankGiftInfo.MoneyID
		response.CostMoneyNum = rankGiftInfo.MoneyNum
	}

	rankGift.BuyTimes -= 1
	player.ActivityModule.RankGift.DB_UpdateBuyTimes(rankGift.GiftID, rankGift.BuyTimes)
	response.BuyTimes = rankGift.BuyTimes

	//! 给予商品
	awardLst := gamedata.GetItemsFromAwardID(rankGiftInfo.Award)
	player.BagMoudle.AddAwardItems(awardLst)

	for _, v := range awardLst {
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})

	}

	response.RetCode = msg.RE_SUCCESS
}

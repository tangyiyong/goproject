package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"gamesvr/tcpclient"
	"msg"
	"net/http"
	"time"
	"utility"
)

func Hand_RegBattleSvr(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RegBattleSvr_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetBattleCamp Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_RegBattleSvr_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var pInfo *TBattleServerInfo = new(TBattleServerInfo)
	pInfo.BatSvrID = req.BatSvrID
	pInfo.SvrInnerAddr = req.ServerInnerAddr
	pInfo.SvrOutAddr = req.ServerOuterAddr
	pInfo.SvrState = 0
	ListLock.Lock()
	G_ServerList[req.BatSvrID] = pInfo
	ListLock.Unlock()
	pInfo.BatClient.ConType = tcpclient.CON_TYPE_BATSVR
	pInfo.BatClient.SvrID = req.BatSvrID
	pInfo.BatClient.ConnectToSvr(pInfo.SvrInnerAddr, 10)
}

func Hand_SetBattleCamp(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SetBattleCamp_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetBattleCamp Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_SetBattleCamp_Ack
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

	player.CamBattleModule.CheckReset()

	if false == gamedata.IsFuncOpen(gamedata.FUNC_CAMPBAT, player.GetLevel(), player.GetVipLevel()) {
		gamelog.Error("Hand_SetBattleCamp Error: Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	if req.BattleCamp <= 0 || req.BattleCamp > 3 {
		gamelog.Error("Hand_SetBattleCamp Error: Invalid BatCamp:%d", req.BattleCamp)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if (req.BattleCamp + player.CamBattleModule.BattleCamp) == 0 {
		//玩家就是选择的我们推存的阵营
		var award TAwardData
		award.TextType = Text_Recommand_Camp
		award.ItemLst = gamedata.GetItemsFromAwardIDEx(gamedata.CampBat_SelCampAward)
		award.Time = time.Now().Unix()
		SendAwardToPlayer(player.playerid, &award)
	}

	player.CamBattleModule.BattleCamp = req.BattleCamp
	player.CamBattleModule.DB_SaveBattleCamp()
	G_SimpleMgr.Set_BatCamp(req.PlayerID, req.BattleCamp)
	G_CampBat_CampKill[req.BattleCamp-1].SetRankItem(req.PlayerID, 0)
	G_CampBat_CampDestroy[req.BattleCamp-1].SetRankItem(req.PlayerID, 0)
	response.RetCode = msg.RE_SUCCESS
}

func Hand_RecommandCamp(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetRecommandCamp_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RecommandCamp Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetRecommandCamp_Ack
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

	player.CamBattleModule.CheckReset()

	if false == gamedata.IsFuncOpen(gamedata.FUNC_CAMPBAT, player.GetLevel(), player.GetVipLevel()) {
		gamelog.Error("Hand_RecommandCamp Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.BattleCamp = 1
	player.CamBattleModule.BattleCamp = 0 - response.BattleCamp
	response.RetCode = msg.RE_SUCCESS
}

//! 进入阵营战
func Hand_EnterCampBattle(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_EnterCampBattle_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_EnterCampBattle Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_EnterCampBattle_Ack
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

	player.CamBattleModule.CheckReset()

	if false == gamedata.IsFuncOpen(gamedata.FUNC_CAMPBAT, player.GetLevel(), player.GetVipLevel()) {
		gamelog.Error("Hand_EnterCampBattle Function not open")
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	player.CamBattleModule.enterCode = int32(utility.Rand())
	response.BattleSvrAddr = GetRecommendSvrAddr()
	gamelog.Error("Hand_EnterCampBattle Genarate EnterCode:%d", player.CamBattleModule.enterCode)
	response.EnterCode = player.CamBattleModule.enterCode
	response.RetCode = msg.RE_SUCCESS
}

//! 获取阵营战数据
func Hand_GetCampBatData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetCampBatData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QueryCampBatInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetCampBatData_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = GetPlayerByID(req.PlayerID)
	if player == nil {
		response.RetCode = msg.RE_INVALID_PLAYERID
		gamelog.Error("Hand_QueryCampBatInfo Error: Invalid PlayerID", req.PlayerID)
		return
	}

	player.CamBattleModule.CheckReset()

	response.KillNum = player.CamBattleModule.Kill
	response.LeftTimes = player.CamBattleModule.LeftTimes
	response.MyRank = G_CampBat_TodayKill.GetRankIndex(player.playerid, response.KillNum)
	response.CampKill = G_CampKillNum

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求阵营战商店的状态
//! 消息: /get_campbat_store_state
func Hand_GetCampbatStoreState(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetCampbatStoreState_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetCampbatStoreState Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetCampbatStoreState_Ack
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

	response.AwardIndex = append(response.AwardIndex, player.CamBattleModule.AwardStoreIndex...)

	for _, v := range player.CamBattleModule.StoreBuyRecord {
		var itemInfo msg.MSG_StoreBuyData
		itemInfo.ID = v.ID
		itemInfo.Times = v.Times
		response.ItemLst = append(response.ItemLst, itemInfo)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买积分商店道具
//! 消息: /buy_campbat_store_item
func Hand_BuyCampbatStoreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BuyCampbatStoreItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyScoreStoreItem Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyCampbatStoreItem_Ack
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

	//! 获取购买物品信息
	itemInfo := gamedata.GetScoreStoreItem(req.StoreItemID)
	if itemInfo == nil {
		gamelog.Error("Hand_BuyScoreStoreItem Error: GetScoreStoreItem nil ID: %d ", req.StoreItemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断购买等级
	if player.GetLevel() < itemInfo.NeedLevel {
		gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough level")
		response.RetCode = msg.RE_NOT_ENOUGH_LEVEL
		return
	}

	//! 根据类型判断积分
	if itemInfo.Type == 2 && itemInfo.NeedScore > player.ScoreMoudle.Score {
		gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough Score")
		response.RetCode = msg.RE_NOT_ENOUGH_SCORE
		return
	}

	//! 判断货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum) == false {
		gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough money")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 判断道具是否足够
	if itemInfo.CostItemID != 0 {
		if player.BagMoudle.IsItemEnough(itemInfo.CostItemID, itemInfo.CostItemNum*req.BuyNum) == false {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough Item")
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}
	}

	//! 检测购买次数是否足够
	if itemInfo.Type == 1 {
		//! 普通商品
		isExist := false
		for i, v := range player.CamBattleModule.StoreBuyRecord {
			if v.ID == req.StoreItemID {
				isExist = true
				if v.Times+req.BuyNum > itemInfo.MaxBuyTime {
					gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
					response.RetCode = msg.RE_NOT_ENOUGH_TIMES
					return
				}

				player.CamBattleModule.StoreBuyRecord[i].Times += req.BuyNum
				go player.CamBattleModule.DB_UpdateStoreItemBuyTimes(i, player.CamBattleModule.StoreBuyRecord[i].Times)
			}
		}

		if isExist == false {
			//! 首次购买
			if req.BuyNum > itemInfo.MaxBuyTime {
				gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
				response.RetCode = msg.RE_NOT_ENOUGH_TIMES
				return
			}

			var itemData TStoreBuyData
			itemData.ID = req.StoreItemID
			itemData.Times = req.BuyNum
			player.ScoreMoudle.StoreBuyRecord = append(player.CamBattleModule.StoreBuyRecord, itemData)
			go player.CamBattleModule.DB_AddStoreItemBuyInfo(itemData)
		}

		//! 扣除货币
		player.RoleMoudle.CostMoney(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum)

		if itemInfo.CostItemID != 0 {
			player.BagMoudle.RemoveNormalItem(itemInfo.CostItemID, itemInfo.CostItemNum*req.BuyNum)
		}

		//! 发放物品
		player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum*req.BuyNum)

	} else if itemInfo.Type == 2 {
		//! 奖励
		if player.CamBattleModule.AwardStoreIndex.IsExist(req.StoreItemID) >= 0 {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
			response.RetCode = msg.RE_NOT_ENOUGH_TIMES
			return
		}

		player.RoleMoudle.CostMoney(itemInfo.CostMoneyID, itemInfo.CostMoneyNum)
		player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum)

		player.CamBattleModule.AwardStoreIndex.Add(req.StoreItemID)
		go player.CamBattleModule.DB_AddStoreAwardInfo(req.StoreItemID)
	}

	response.RetCode = msg.RE_SUCCESS
}

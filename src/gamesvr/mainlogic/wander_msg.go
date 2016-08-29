package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 获取今日活动
func Hand_WanderReset(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_WanderReset_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WanderReset Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WanderReset_Ack
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

	if player.WanderMoudle.LeftTime <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderReset Error : not enough reset time!!!")
		return
	}

	player.WanderMoudle.LeftTime -= 1
	player.WanderMoudle.CurCopyID = 0
	player.WanderMoudle.CanBattle = 1

	response.RetCode = msg.RE_SUCCESS
	response.CurCopyID = 0
	response.LeftTime = player.WanderMoudle.LeftTime
	response.CanBattle = 1

	player.WanderMoudle.DB_Reset()
}

//! 获取今日活动
func Hand_WanderOpenBox(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_WanderOpenBox_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WanderOpenBox Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WanderOpenBox_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Error("Hand_WanderOpenBox Error: %s", string(b))
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if req.DrawType == 1 { //单抽
		if player.WanderMoudle.SingleFree == false {
			if false == player.RoleMoudle.CheckMoneyEnough(gamedata.WanderDrawMoneyID, gamedata.WanderDrawNum) {
				response.RetCode = msg.RE_NOT_ENOUGH_MONEY
				gamelog.Error("Hand_WanderOpenBox Error : Not Enough money!!")
				return
			}

			player.RoleMoudle.CostMoney(gamedata.WanderDrawMoneyID, gamedata.WanderDrawNum)
		}

		awardLst := gamedata.GetItemsFromAwardID(gamedata.WanderSingleBoxID)
		player.BagMoudle.AddAwardItems(awardLst)
		for _, v := range awardLst {
			var item msg.MSG_ItemData
			item.ID = v.ItemID
			item.Num = v.ItemNum
			response.ItemLst = append(response.ItemLst, item)
		}

		player.WanderMoudle.SingleFree = true

	} else if req.DrawType == 2 { //十连抽
		if false == player.RoleMoudle.CheckMoneyEnough(gamedata.WanderDrawMoneyID, gamedata.WanderTenDrawNum) {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			gamelog.Error("Hand_WanderOpenBox Error : Not Enough money!!")
			return
		}

		awardLst := gamedata.GetItemsFromAwardIDEx(gamedata.WanderTenBoxID)
		player.BagMoudle.AddAwardItems(awardLst)
		for _, v := range awardLst {
			var item msg.MSG_ItemData
			item.ID = v.ItemID
			item.Num = v.ItemNum
			response.ItemLst = append(response.ItemLst, item)
		}

		response.ItemLst = append(response.ItemLst, msg.MSG_ItemData{ID: gamedata.WanderTenGiftID, Num: gamedata.WanderTenGiftNum})
		player.BagMoudle.AddAwardItem(gamedata.WanderTenGiftID, gamedata.WanderTenGiftNum)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 获取今日活动
func Hand_WanderSweep(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_WanderSweep_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WanderSweep Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WanderSweep_Ack
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

	if player.WanderMoudle.CanBattle == 0 {
		gamelog.Error("Hand_WanderSweep Error: you can't battle unless you reset")
		return
	}

	if req.TargetCopyID > player.WanderMoudle.MaxCopyID {
		gamelog.Error("Hand_WanderSweep Error: cant sweep to max copyid:%d", req.TargetCopyID)
		return
	}

	if req.TargetCopyID != player.WanderMoudle.CurCopyID+1 {
		gamelog.Error("Hand_WanderSweep Error: Invalid target copyid:%d", req.TargetCopyID)
		return
	}

	dropItem := gamedata.GetItemsFromAwardID(req.TargetCopyID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	player.BagMoudle.AddAwardItems(dropItem)
	player.WanderMoudle.CurCopyID = req.TargetCopyID
	response.CurCopyID = player.WanderMoudle.CurCopyID
	response.RetCode = msg.RE_SUCCESS
	player.WanderMoudle.DB_Reset()

	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_WANDER_SWEEP, 1)
}

//! 获取今日活动
func Hand_WanderGetInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_WanderGetInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WanderGetInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WanderGetInfo_Ack
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

	if false == gamedata.IsFuncOpen(gamedata.FUNC_WANDER, player.GetLevel(), player.GetVipLevel()) {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		gamelog.Error("Hand_WanderGetInfo Error : Func Wander not open!!!")
		return
	}

	response.RetCode = msg.RE_SUCCESS
	response.MaxCopyID = player.WanderMoudle.MaxCopyID
	response.CurCopyID = player.WanderMoudle.CurCopyID
	response.LeftTime = player.WanderMoudle.LeftTime
	response.CanBattle = player.WanderMoudle.CanBattle
}

//! 获取今日活动
func Hand_WanderCheck(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_WanderCheck_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WanderCheck Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WanderCheck_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.TargetCopyID <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderCheck Error: Invalid target copyid :%d", req.TargetCopyID)
		return
	}

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if player.WanderMoudle.CanBattle == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderCheck Error: you can't battle unless you reset")
		return
	}

	if req.TargetCopyID <= player.WanderMoudle.CurCopyID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderCheck Error: Invalid target copyid :%d, maxid:%d", req.TargetCopyID, player.WanderMoudle.CurCopyID)
		return
	}

	if player.WanderMoudle.CurCopyID == 0 {
		if req.TargetCopyID != gamedata.WanderBeginID {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_WanderCheck Error: Invalid target copyid :%d, beginid:%d, curid:%d", req.TargetCopyID, gamedata.WanderBeginID, player.WanderMoudle.CurCopyID)
			return
		}
	} else if req.TargetCopyID != player.WanderMoudle.CurCopyID+1 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderCheck Error: Invalid target copyid :%d, beginid:%d, curid:%d", req.TargetCopyID, gamedata.WanderBeginID, player.WanderMoudle.CurCopyID)
		return
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 获取今日活动
func Hand_WanderResult(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_WanderResult_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_WanderResult Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_WanderResult_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.TargetCopyID <= 0 {
		gamelog.Error("Hand_WanderResult Error: Invalid target copyid :%d", req.TargetCopyID)
		return
	}

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if player.WanderMoudle.CanBattle == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderResult Error: you can't battle unless you reset")
		return
	}

	if player.WanderMoudle.MaxCopyID == 0 {
		if req.TargetCopyID != gamedata.WanderBeginID {
			gamelog.Error("Hand_WanderResult Error: Invalid target copyid :%d", req.TargetCopyID)
			return
		}
	} else if req.TargetCopyID != player.WanderMoudle.MaxCopyID+1 {
		gamelog.Error("Hand_WanderResult Error: Invalid target copyid :%d", req.TargetCopyID)
		return
	}

	if req.Win == 0 {
		response.RetCode = msg.RE_SUCCESS
		player.WanderMoudle.CanBattle = 0
		response.CurCopyID = player.WanderMoudle.CurCopyID
		return
	}

	pCopyInfo := gamedata.GetCopyBaseInfo(req.TargetCopyID)
	if pCopyInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_WanderResult Error : Invalid copy id:%d", req.TargetCopyID)
		return
	}

	dropItem := gamedata.GetItemsFromAwardID(pCopyInfo.AwardID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}
	player.BagMoudle.AddAwardItems(dropItem)
	player.WanderMoudle.MaxCopyID = req.TargetCopyID
	player.WanderMoudle.CurCopyID = req.TargetCopyID
	G_WanderRanker.SetRankItem(req.PlayerID, response.CurCopyID)
	response.CurCopyID = player.WanderMoudle.CurCopyID
	response.RetCode = msg.RE_SUCCESS
	player.WanderMoudle.DB_Reset()
}

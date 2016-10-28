package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//请求挂机信息
func Hand_GetHangUpInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetHangUp_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetHangUpInfo : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetHangUp_Ack
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

	if player.HangMoudle.CurBossID > 0 {
		if player.HangMoudle.ReceiveHangUpProduce() == true {
			player.HangMoudle.DB_SaveHangUpState()
		}
	}

	response.RetCode = msg.RE_SUCCESS
	response.GridNum = player.HangMoudle.GridNum
	response.ExpItems = player.HangMoudle.ExpItems
	response.LeftQuick = gamedata.GetFuncVipValue(gamedata.FUNC_HANGUP_QUICKTIME, player.GetVipLevel()) - player.HangMoudle.QuickTime
	response.History = player.HangMoudle.History
	response.CurBossID = player.HangMoudle.CurBossID
}

//设置挂机BOSS
func Hand_SetBoss(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SetHangUpBoss_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetBoss : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_SetHangUpBoss_Ack
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

	pHangUpInfo := gamedata.GetHangUpInfo(req.BossID)
	if pHangUpInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_SetBoss : Invalid BossID:%d", req.BossID)
		return
	}

	if player.GetLevel() < pHangUpInfo.Level {
		response.RetCode = msg.RE_NOT_ENOUGH_HERO_LEVEL
		gamelog.Error("Hand_SetBoss : Not Enough Level:%d", pHangUpInfo.Level)
		return
	}
	if player.HangMoudle.CurBossID != 0 {
		player.HangMoudle.ReceiveHangUpProduce()
	}
	player.HangMoudle.StartTime = utility.GetCurTime()
	player.HangMoudle.CurBossID = req.BossID
	player.HangMoudle.DB_SaveHangUpState()
	response.CurBossID = req.BossID
	response.RetCode = msg.RE_SUCCESS
}

//快速战斗请求
func Hand_QuickFight(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_QuickFight_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_QuickFight : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_QuickFight_Ack
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

	if player.HangMoudle.CurBossID <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_QuickFight : Invalid CurBossID:%d", player.HangMoudle.CurBossID)
		return
	}

	pHangUpInfo := gamedata.GetHangUpInfo(player.HangMoudle.CurBossID)
	if pHangUpInfo == nil {
		gamelog.Error("Hand_QuickFight : Invalid BossID:%d", player.HangMoudle.CurBossID)
		return
	}

	maxtime := gamedata.GetFuncVipValue(gamedata.FUNC_HANGUP_QUICKTIME, player.GetVipLevel())
	if player.HangMoudle.QuickTime >= maxtime {
		response.RetCode = msg.RE_NOT_ENOUGH_REFRESH_TIMES
		gamelog.Error("Hand_QuickFight : Time Limit:%d", maxtime)
		return
	}

	for i := 0; i < gamedata.HangUpQuickFight; i++ {
		if utility.Rand() < player.HangMoudle.CalcHangUpRatio(player.GetFightValue(), pHangUpInfo.FightValue) {
			for j := 0; j < pHangUpInfo.ProduceNum; j++ {
				if len(player.HangMoudle.ExpItems) < player.HangMoudle.GridNum {
					player.HangMoudle.ExpItems = append(player.HangMoudle.ExpItems, pHangUpInfo.ProduceID)
				}
			}
			player.HangMoudle.History = append(player.HangMoudle.History, msg.THisHang{player.HangMoudle.CurBossID, pHangUpInfo.ProduceID,
				pHangUpInfo.ProduceNum, utility.GetCurTime()})
		} else {
			player.HangMoudle.History = append(player.HangMoudle.History, msg.THisHang{player.HangMoudle.CurBossID, pHangUpInfo.ProduceID,
				0, utility.GetCurTime()})
		}
	}
	player.HangMoudle.QuickTime += 1

	player.HangMoudle.DB_SaveQuickFightResult()
	response.QuickTime = player.HangMoudle.QuickTime
	response.RetCode = msg.RE_SUCCESS
	response.History = player.HangMoudle.History
	response.ExpItems = player.HangMoudle.ExpItems
}

//一键使用经验丹
func Hand_UseExpItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_UseExpItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_UseExpItem : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_UseExpItem_Ack
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

	for j := 0; j < len(player.HangMoudle.ExpItems); j++ {
		pItemInfo := gamedata.GetItemInfo(player.HangMoudle.ExpItems[j])
		if pItemInfo == nil {
			gamelog.Error("Hand_UseExpItem : Invalid Experience Item ID: %d", player.HangMoudle.ExpItems[j])
		} else {
			response.CurExp += pItemInfo.Data1
		}
	}

	player.HeroMoudle.AddMainHeroExp(response.CurExp)

	player.HangMoudle.ExpItems = make([]int, 0)

	player.HangMoudle.DB_ClearHangUpBag()

	response.RetCode = msg.RE_SUCCESS

	return
}

func Hand_AddGrid(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_AddGridNum_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_AddGrid : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_AddGridNum_Ack
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

	CostMoney := gamedata.GetFuncTimeCost(gamedata.FUNC_HANGUP_GRID_OPNE, player.HangMoudle.AddGridTime+1)
	if false == player.RoleMoudle.CheckMoneyEnough(gamedata.HangUpBuyGridMoneyID, CostMoney) {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_AddGrid : Not Enough Money")
		return
	}

	player.HangMoudle.GridNum += gamedata.HangUpOpenGridNum
	player.RoleMoudle.CostMoney(gamedata.HangUpBuyGridMoneyID, CostMoney)
	player.HangMoudle.DB_SaveGridNum()
	response.RetCode = msg.RE_SUCCESS
	response.GridNum = player.HangMoudle.GridNum
}

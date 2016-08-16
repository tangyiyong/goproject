package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

//! 玩家请求送花
func Hand_SendFlower(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SendFlower_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetGuild Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_SendFlower_Ack
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

	player.FameHallModule.CheckReset()

	//! 检测是否还有免费次数
	if player.FameHallModule.FreeTimes <= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 更新目标魅力值
	if req.SendIndex >= 6 {
		gamelog.Error("Hand_SendFlower Error: Invalid Index %d", req.SendIndex)
	}

	playerID := G_FameHallLst[req.SendType][req.SendIndex].PlayerID

	//! 检测是否为已送目标
	if req.SendType == 0 {
		if player.FameHallModule.SendFightID.IsExist(req.SendIndex) >= 0 {
			gamelog.Error("Hand_SendFlower Error: Aleady send flower index: %d", req.SendIndex)
			response.RetCode = msg.RE_AlEADY_SEND
			return
		}

		player.FameHallModule.SendFightID = append(player.FameHallModule.SendFightID, req.SendIndex)
		go player.FameHallModule.DB_AddSendFightID(req.SendIndex)
	} else {
		if player.FameHallModule.SendLevelID.IsExist(req.SendIndex) >= 0 {
			gamelog.Error("Hand_SendFlower Error: Aleady send flower index: %d", req.SendIndex)
			response.RetCode = msg.RE_AlEADY_SEND
			return
		}

		player.FameHallModule.SendLevelID = append(player.FameHallModule.SendLevelID, req.SendIndex)
		go player.FameHallModule.DB_AddSendLevelID(req.SendIndex)
	}

	targetPlayer := GetPlayerByID(playerID)
	if targetPlayer == nil {
		targetPlayer = LoadPlayerFromDB(playerID)
		if targetPlayer == nil {
			gamelog.Error("Hand_SendFlower Error: Not find player")
			response.RetCode = msg.RE_INVALID_PLAYERID
			return
		}
	}

	targetPlayer.FameHallModule.CharmValue += 1
	go targetPlayer.FameHallModule.DB_UpdateCharm()

	for i, j := range G_FameHallLst {
		for n, m := range j {
			if m.PlayerID == playerID {
				G_FameHallLst[i][n].CharmValue += 1
			}
		}
	}

	//! 减去自己送花次数
	player.FameHallModule.FreeTimes -= 1
	go targetPlayer.FameHallModule.DB_UpdateFreeTimes()

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS

	for i := 0; i < len(G_FameHallLst[req.SendType]); i++ {
		response.CharmValue = append(response.CharmValue, G_FameHallLst[req.SendType][i].CharmValue)
	}
}

//! 玩家请求查询魅力值
func Hand_GetCharmValue(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetCharm_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetCharmValue Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_GetCharm_Ack
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

	if req.RankType > 1 || req.RankType < 0 {
		gamelog.Error("Hand_GetCharmValue Error: Invalid Type %d", req.RankType)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	player.FameHallModule.CheckReset()

	for i := 0; i < len(G_FameHallLst[req.RankType]); i++ {
		if G_FameHallLst[req.RankType][i].PlayerID == 0 {
			continue
		}

		simpleInfo := G_SimpleMgr.GetSimpleInfoByID(G_FameHallLst[req.RankType][i].PlayerID)
		response.FightValue = append(response.FightValue, simpleInfo.FightValue)
		response.CharmValue = append(response.CharmValue, G_FameHallLst[req.RankType][i].CharmValue)
		response.Level = append(response.Level, simpleInfo.Level)
		response.Name = append(response.Name, simpleInfo.Name)
		response.PlayerID = append(response.PlayerID, G_FameHallLst[req.RankType][i].PlayerID)
		response.HeroID = response.HeroID
		if len(response.CharmValue) >= 6 {
			break
		}
	}

	response.Times = player.FameHallModule.FreeTimes

	if req.RankType == 0 {
		response.SendID = append(response.SendID, player.FameHallModule.SendFightID...)
	} else {
		response.SendID = append(response.SendID, player.FameHallModule.SendLevelID...)
	}

	response.RetCode = msg.RE_SUCCESS
}

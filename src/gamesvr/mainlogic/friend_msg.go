package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"strings"
)

func Hand_GetAllFriend(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetAllFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAllFriend unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetAllFriend_Ack
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

	player.FriendMoudle.CheckReset()

	response.FriendLst = make([]msg.MSG_FriendInfo, len(player.FriendMoudle.FriendList))
	response.RetCode = msg.RE_SUCCESS
	response.RcvNum = player.FriendMoudle.RcvNum
	for i := 0; i < len(player.FriendMoudle.FriendList); i++ {
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(player.FriendMoudle.FriendList[i].PlayerID)
		if pSimpleInfo != nil {
			response.FriendLst[i].PlayerID = pSimpleInfo.PlayerID
			response.FriendLst[i].FightValue = pSimpleInfo.FightValue
			response.FriendLst[i].HeroID = pSimpleInfo.HeroID
			response.FriendLst[i].Quality = pSimpleInfo.Quality
			response.FriendLst[i].OffTime = pSimpleInfo.LogoffTime
			response.FriendLst[i].GuildName = GetGuildName(pSimpleInfo.GuildID)
			response.FriendLst[i].Name = pSimpleInfo.Name
			response.FriendLst[i].Level = pSimpleInfo.Level
			response.FriendLst[i].IsGive = 0
			response.FriendLst[i].HasAct = 0

			if pSimpleInfo.isOnline == true {
				response.FriendLst[i].OffTime = 0
			}

			if player.FriendMoudle.FriendList[i].IsGive == true {
				response.FriendLst[i].IsGive = 1
			}
			if player.FriendMoudle.FriendList[i].HasAct == true {
				response.FriendLst[i].HasAct = 1
			}
		}
	}

}

func Hand_GetOnlineFriend(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetOnlineFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetOnlineFriend unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetOnlineFriend_Ack
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

	player.FriendMoudle.CheckReset()

	response.RetCode = msg.RE_SUCCESS
	for i := 0; i < len(player.FriendMoudle.FriendList); i++ {
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(player.FriendMoudle.FriendList[i].PlayerID)
		if pSimpleInfo != nil && pSimpleInfo.isOnline == true {
			response.OnlineLst = append(response.OnlineLst, msg.MSG_OnlineInfo{})
			response.OnlineLst[len(response.OnlineLst)-1].PlayerID = pSimpleInfo.PlayerID
			response.OnlineLst[len(response.OnlineLst)-1].Name = pSimpleInfo.Name
			response.OnlineLst[len(response.OnlineLst)-1].FightValue = pSimpleInfo.FightValue
			response.OnlineLst[len(response.OnlineLst)-1].HeroID = pSimpleInfo.HeroID
			response.OnlineLst[len(response.OnlineLst)-1].Quality = pSimpleInfo.Quality
			response.OnlineLst[len(response.OnlineLst)-1].GuildName = GetGuildName(pSimpleInfo.GuildID)
			response.OnlineLst[len(response.OnlineLst)-1].Level = pSimpleInfo.Level
		}
	}

}

func Hand_DelFriendReq(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_DelFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_DelFriendReq unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_DelFriend_Ack
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

	//删除的好友要存在
	pTarget := GetPlayerByID(req.TargetID)
	if pTarget != nil {

	}

	return
}

func Hand_AddFriendReq(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_AddFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_AddFriendReq unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_AddFriend_Ack
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

	//不能加自己好友
	if req.PlayerID == req.TargetID {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_AddFriendReq Error : can't add myself as friend!")
		return
	}

	//己到好友上限不能加好友
	nCount := len(player.FriendMoudle.FriendList)
	if nCount >= gamedata.GetFuncVipValue(gamedata.FUNC_FRIEND_NUM_LIMIT, player.GetVipLevel()) {
		response.RetCode = msg.RE_REACH_FRIEND_NUM_LIMIT
		gamelog.Error("Hand_AddFriendReq Error : alreay reach friend limit, cant add friend.")
		return
	}

	//不能加己是好友的人为好友
	for i := 0; i < nCount; i++ {
		if req.TargetID == player.FriendMoudle.FriendList[i].PlayerID {
			response.RetCode = msg.RE_ALREADY_FRIEND
			gamelog.Error("Hand_AddFriendReq Error : is already friend!!.")
			return
		}
	}

	response.RetCode = msg.RE_SUCCESS
	pTargetPlayer := GetPlayerByID(req.TargetID)
	if pTargetPlayer != nil {
		if pTargetPlayer.FriendMoudle.ApplyList.IsExist(req.PlayerID) < 0 {
			pTargetPlayer.FriendMoudle.ApplyList = append(pTargetPlayer.FriendMoudle.ApplyList, req.PlayerID)
		} else {
			return
		}
	}

	DB_AddFriendAppList(req.TargetID, req.PlayerID)

	return
}

func SelectFriendTarget(player *TPlayer, value int) bool {
	nCount := len(player.FriendMoudle.FriendList)
	if nCount >= gamedata.GetFuncVipValue(gamedata.FUNC_FRIEND_NUM_LIMIT, player.GetVipLevel()) {
		return false
	}

	return true
}

func Hand_RecomandFriend(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RecomandFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RecomandFriend unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RecomandFriend_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if len(g_Players) <= 1 {
		response.RetCode = msg.RE_SUCCESS
		return
	}

	for {
		pTarget := GetSelectPlayer(SelectFriendTarget, 1000)
		if pTarget == nil {
			gamelog.Error("Hand_RecomandFriend Error : GetSelectPlayerDef return nil!!!")
			break
		}

		if pTarget.playerid == req.PlayerID {
			continue
		}

		if len(response.FriendLst) >= 4 {
			break
		}

		var bFriend bool = false
		for _, v := range player.FriendMoudle.FriendList {
			if v.PlayerID == pTarget.playerid {
				bFriend = true
				break
			}
		}

		if bFriend == true {
			continue
		}

		var bFind bool = false
		for _, v := range response.FriendLst {
			if v.PlayerID == pTarget.playerid {
				bFind = true
				break
			}
		}

		if bFind == true {
			break
		}

		response.FriendLst = append(response.FriendLst, msg.MSG_FriendInfo{})
		nCount := len(response.FriendLst)
		if pTarget.pSimpleInfo != nil {
			response.FriendLst[nCount-1].PlayerID = pTarget.pSimpleInfo.PlayerID
			response.FriendLst[nCount-1].FightValue = pTarget.pSimpleInfo.FightValue
			response.FriendLst[nCount-1].HeroID = pTarget.pSimpleInfo.HeroID
			response.FriendLst[nCount-1].Quality = pTarget.pSimpleInfo.Quality
			response.FriendLst[nCount-1].OffTime = pTarget.pSimpleInfo.LogoffTime
			response.FriendLst[nCount-1].GuildName = GetGuildName(pTarget.pSimpleInfo.GuildID)
			response.FriendLst[nCount-1].Name = pTarget.pSimpleInfo.Name
			response.FriendLst[nCount-1].Level = pTarget.pSimpleInfo.Level
			if pTarget.pSimpleInfo.isOnline == true {
				response.FriendLst[nCount-1].OffTime = 0
			}
		}
	}

	response.RetCode = msg.RE_SUCCESS
	return
}

func Hand_GiveAction(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GiveAction_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GiveAction unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GiveAction_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if req.TargetID == 0 {

	} else {
		pFriendInfo, nIndex := player.FriendMoudle.GetFriendByID(req.TargetID)
		if nIndex < 0 {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GiveAction Error: Invalid TargetID :%d", req.TargetID)
			return
		}

		if pFriendInfo.IsGive == true {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_GiveAction Error: Already Give Action")
			return
		}

		pFriendInfo.IsGive = true
		DB_UpdateIsGive(player.playerid, nIndex, true)

		pTarget := GetPlayerByID(req.TargetID)
		if pTarget != nil {
			ptFriendInfo, ntIndex := pTarget.FriendMoudle.GetFriendByID(player.playerid)
			if ntIndex >= 0 {
				ptFriendInfo.HasAct = true
				DB_UpdateHasAct(req.TargetID, ntIndex, true)
			} else {

				gamelog.Error("Hand_GiveAction Error: Tartet Player cant find this friend t:%d, this:%d", req.TargetID, player.playerid)
			}
		} else {
			//DB_UpdateHasAct(req.TargetID, nIndex, true)
		}
	}

	//! 增送精力次数
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SEND_ACTION, 1)

	response.RetCode = msg.RE_SUCCESS
}

func Hand_ReceiveAction(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_ReceiveAction_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ReceiveAction unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ReceiveAction_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if player.FriendMoudle.RcvNum >= gamedata.MaxRecvTime {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ReceiveAction Error : Cant Receive More time cur:%d!!", player.FriendMoudle.RcvNum)
		return
	}

	nIndex := -1
	if req.TargetID != 0 {
		for i := 0; i < len(player.FriendMoudle.FriendList); i++ {
			if player.FriendMoudle.FriendList[i].PlayerID == req.TargetID {
				nIndex = i
				break
			}
		}
	}

	if req.TargetID == 0 {
		actNum := 0
		for i := 0; i < len(player.FriendMoudle.FriendList); i++ {
			if player.FriendMoudle.FriendList[i].HasAct == true && player.FriendMoudle.RcvNum < gamedata.MaxRecvTime {
				player.FriendMoudle.FriendList[nIndex].HasAct = false
				DB_UpdateHasAct(req.PlayerID, nIndex, false)
				actNum += gamedata.GiveActionNum
				player.FriendMoudle.RcvNum += 1
			}
		}
		player.RoleMoudle.AddAction(gamedata.GiveActionID, actNum)
		DB_UpdateRcvNum(req.PlayerID, player.FriendMoudle.RcvNum)
	} else {
		if nIndex < 0 || nIndex >= len(player.FriendMoudle.FriendList) {
			response.RetCode = msg.RE_INVALID_PARAM
			gamelog.Error("Hand_ReceiveAction Error: Invalid TargetID:%d", req.TargetID)
			return
		}

		if player.FriendMoudle.FriendList[nIndex].HasAct == false {
			response.RetCode = msg.RE_ALREADY_RECEIVED
			gamelog.Error("Hand_ReceiveAction Error: Alreay Received:%d", req.TargetID)
			return
		}

		player.FriendMoudle.FriendList[nIndex].HasAct = false
		player.FriendMoudle.RcvNum += 1
		DB_UpdateHasAct(req.PlayerID, nIndex, false)
		player.RoleMoudle.AddAction(gamedata.GiveActionID, gamedata.GiveActionNum)
		DB_UpdateRcvNum(req.PlayerID, player.FriendMoudle.RcvNum)
	}
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(gamedata.GiveActionID)
	response.RetCode = msg.RE_SUCCESS
	response.FriendLst = make([]msg.MSG_FriendInfo, len(player.FriendMoudle.FriendList))
	for i := 0; i < len(player.FriendMoudle.FriendList); i++ {
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(player.FriendMoudle.FriendList[i].PlayerID)
		if pSimpleInfo != nil {
			response.FriendLst[i].PlayerID = pSimpleInfo.PlayerID
			response.FriendLst[i].FightValue = pSimpleInfo.FightValue
			response.FriendLst[i].HeroID = pSimpleInfo.HeroID
			response.FriendLst[i].Quality = pSimpleInfo.Quality
			response.FriendLst[i].OffTime = pSimpleInfo.LogoffTime
			response.FriendLst[i].GuildName = GetGuildName(pSimpleInfo.GuildID)
			response.FriendLst[i].Name = pSimpleInfo.Name
			response.FriendLst[i].Level = pSimpleInfo.Level
			response.FriendLst[i].IsGive = 0
			response.FriendLst[i].HasAct = 0
			if pSimpleInfo.isOnline == true {
				response.FriendLst[i].OffTime = 0
			}
			if player.FriendMoudle.FriendList[i].IsGive == true {
				response.FriendLst[i].IsGive = 1
			}
			if player.FriendMoudle.FriendList[i].HasAct == true {
				response.FriendLst[i].HasAct = 1
			}
		}
	}
}

func Hand_SearchFriend(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SearchFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SearchFriend unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SearchFriend_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.RetCode = msg.RE_SUCCESS
	for _, pTarget := range g_Players {
		if pTarget.playerid == req.PlayerID {
			continue
		}
		if true == strings.Contains(pTarget.RoleMoudle.Name, req.Name) {
			if pTarget.pSimpleInfo != nil {
				response.FriendLst = append(response.FriendLst, msg.MSG_FriendInfo{})
				response.FriendLst[len(response.FriendLst)-1].PlayerID = pTarget.pSimpleInfo.PlayerID
				response.FriendLst[len(response.FriendLst)-1].FightValue = pTarget.pSimpleInfo.FightValue
				response.FriendLst[len(response.FriendLst)-1].HeroID = pTarget.pSimpleInfo.HeroID
				response.FriendLst[len(response.FriendLst)-1].Quality = pTarget.pSimpleInfo.Quality
				response.FriendLst[len(response.FriendLst)-1].OffTime = pTarget.pSimpleInfo.LogoffTime
				response.FriendLst[len(response.FriendLst)-1].GuildName = GetGuildName(pTarget.pSimpleInfo.GuildID)
				response.FriendLst[len(response.FriendLst)-1].Name = pTarget.pSimpleInfo.Name
				response.FriendLst[len(response.FriendLst)-1].Level = pTarget.pSimpleInfo.Level
				if pTarget.pSimpleInfo.isOnline == true {
					response.FriendLst[len(response.FriendLst)-1].OffTime = 0
				}
			}
		}
	}
}

func Hand_GetApplyList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetApplyList_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetApplyList unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetApplyList_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.FriendLst = make([]msg.MSG_FriendInfo, len(player.FriendMoudle.ApplyList))
	response.RetCode = msg.RE_SUCCESS
	for i := 0; i < len(player.FriendMoudle.ApplyList); i++ {
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(player.FriendMoudle.ApplyList[i])
		if pSimpleInfo != nil {
			response.FriendLst[i].PlayerID = pSimpleInfo.PlayerID
			response.FriendLst[i].Name = pSimpleInfo.Name
			response.FriendLst[i].FightValue = pSimpleInfo.FightValue
			response.FriendLst[i].HeroID = pSimpleInfo.HeroID
			response.FriendLst[i].Quality = pSimpleInfo.Quality
			response.FriendLst[i].OffTime = pSimpleInfo.LogoffTime
			response.FriendLst[i].GuildName = GetGuildName(pSimpleInfo.GuildID)
			response.FriendLst[i].Level = pSimpleInfo.Level
		}
	}
}

func Hand_ProcessFriendReq(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_ProcessFriend_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ProcessFriendReq unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ProcessFriend_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if req.TargetID == 0 { //全部处理
		if req.IsAgree == 1 {
			for _, v := range player.FriendMoudle.ApplyList {

				pFriend, _ := player.FriendMoudle.GetFriendByID(v)
				if pFriend != nil {
					continue
				}

				player.FriendMoudle.FriendList = append(player.FriendMoudle.FriendList, TFriendInfo{v, false, false})
				DB_AddFriend(player.playerid, &player.FriendMoudle.FriendList[len(player.FriendMoudle.FriendList)-1])
				pTarget := GetPlayerByID(v)
				if pTarget != nil {
					pTarget.FriendMoudle.FriendList = append(pTarget.FriendMoudle.FriendList, TFriendInfo{req.PlayerID, false, false})
				}

				DB_AddFriend(v, &TFriendInfo{req.PlayerID, false, false})
				pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(v)
				if pSimpleInfo != nil {
					response.FriendLst = append(response.FriendLst, msg.MSG_FriendInfo{})
					response.FriendLst[len(response.FriendLst)-1].PlayerID = pSimpleInfo.PlayerID
					response.FriendLst[len(response.FriendLst)-1].Name = pSimpleInfo.Name
					response.FriendLst[len(response.FriendLst)-1].FightValue = pSimpleInfo.FightValue
					response.FriendLst[len(response.FriendLst)-1].HeroID = pSimpleInfo.HeroID
					response.FriendLst[len(response.FriendLst)-1].Quality = pSimpleInfo.Quality
					response.FriendLst[len(response.FriendLst)-1].OffTime = pSimpleInfo.LogoffTime
					response.FriendLst[len(response.FriendLst)-1].GuildName = GetGuildName(pSimpleInfo.GuildID)
					response.FriendLst[len(response.FriendLst)-1].Level = pSimpleInfo.Level
					if pSimpleInfo.isOnline == true {
						response.FriendLst[len(response.FriendLst)-1].OffTime = 0
					}
				}
			}
		}

		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FRIEND_NUM, len(player.FriendMoudle.ApplyList))

		player.FriendMoudle.ApplyList = make([]int32, 0)
		player.FriendMoudle.DB_ClearAppList()

	} else { //单独处理
		var nIndex int = player.FriendMoudle.ApplyList.IsExist(req.TargetID)
		if nIndex < 0 || nIndex >= len(player.FriendMoudle.ApplyList) {
			gamelog.Error("Hand_ProcessFriendReq Error : Invalid TargetID:%d", req.TargetID)
			return
		}

		pFriend, _ := player.FriendMoudle.GetFriendByID(req.TargetID)
		if pFriend != nil {
			return
		}

		if req.IsAgree == 1 {
			player.FriendMoudle.FriendList = append(player.FriendMoudle.FriendList, TFriendInfo{req.TargetID, false, false})
			DB_AddFriend(player.playerid, &player.FriendMoudle.FriendList[len(player.FriendMoudle.FriendList)-1])
			pTarget := GetPlayerByID(req.TargetID)
			if pTarget != nil {
				pTarget.FriendMoudle.FriendList = append(pTarget.FriendMoudle.FriendList, TFriendInfo{req.PlayerID, false, false})
			}
			DB_AddFriend(req.TargetID, &TFriendInfo{req.PlayerID, false, false})

			pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(req.TargetID)
			if pSimpleInfo != nil {
				response.FriendLst = append(response.FriendLst, msg.MSG_FriendInfo{})
				response.FriendLst[len(response.FriendLst)-1].PlayerID = pSimpleInfo.PlayerID
				response.FriendLst[len(response.FriendLst)-1].Name = pSimpleInfo.Name
				response.FriendLst[len(response.FriendLst)-1].FightValue = pSimpleInfo.FightValue
				response.FriendLst[len(response.FriendLst)-1].HeroID = pSimpleInfo.HeroID
				response.FriendLst[len(response.FriendLst)-1].Quality = pSimpleInfo.Quality
				response.FriendLst[len(response.FriendLst)-1].OffTime = pSimpleInfo.LogoffTime
				response.FriendLst[len(response.FriendLst)-1].GuildName = GetGuildName(pSimpleInfo.GuildID)
				response.FriendLst[len(response.FriendLst)-1].Level = pSimpleInfo.Level
				if pSimpleInfo.isOnline == true {
					response.FriendLst[len(response.FriendLst)-1].OffTime = 0
				}
			}
			player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FRIEND_NUM, 1)
		}

		player.FriendMoudle.ApplyList = append(player.FriendMoudle.ApplyList[:nIndex], player.FriendMoudle.ApplyList[nIndex+1:]...)
		DB_RemoveFriendAppList(player.playerid, req.TargetID)
	}

	response.RetCode = msg.RE_SUCCESS
	response.RcvNum = player.FriendMoudle.RcvNum

	return
}

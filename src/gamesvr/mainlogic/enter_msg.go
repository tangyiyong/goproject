package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"gamesvr/sessionmgr"
	"mongodb"
	"msg"
	"net/http"
	"utility"
)

func Hand_LoginGame(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_LoginGameSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_PlayerLoginGame : Unmarshal error!!!!")
		return
	}

	if false == sessionmgr.CheckLoginTime(req.AccountID) {
		gamelog.Error("CheckLoginTime Error , Repeate Login!!!")
		return
	}

	var response msg.MSG_LoginGameSvr_Ack
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var playerid = GetPlayerIDByAccountID(req.AccountID)

	bcheck := CheckUserIsLogin(req.AccountID, playerid, req.LoginKey)
	if !bcheck {
		response.RetCode = msg.RE_INVALID_LOGINKEY
		gamelog.Error("Invalid Login key!!!!!")
		return
	}
	response.SessionKey = sessionmgr.NewSessionKey()
	sessionmgr.AddSessionKey(req.AccountID, response.SessionKey)
	response.PlayerID = playerid
	response.RetCode = msg.RE_SUCCESS
}

func Hand_EnterGame(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_EnterGameSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_PlayerEnterGame : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_EnterGameSvr_Ack
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if !sessionmgr.CheckSessionKey(req.PlayerID, req.SessionKey) {
		response.RetCode = msg.RE_INVALID_SESSIONKEY
		return
	}

	var player *TPlayer = GetPlayerByID(req.PlayerID)
	if player == nil {
		//如果内存中没有，查看数据库中是否存在
		player = LoadPlayerFromDB(req.PlayerID)
		if player == nil {
			gamelog.Error("Hand_PlayerEnterGame Error : Invalid Playerid :%d", req.PlayerID)
			response.RetCode = msg.RE_INVALID_PLAYERID
			return
		}
	}

	//! 还原XorCode
	var xorCode uint
	for i := 0; i < 4; i++ {
		index := 8 * (3 - i)
		left := uint(G_XorCode[i])
		xorCode = uint(xorCode) + left<<uint(index)
	}
	response.XorCode = int(xorCode)
	player.OnPlayerOnline(req.PlayerID)
	response.SvrTime = utility.GetCurTime()
	response.ChatSvrAddr = appconfig.ChatSvrAddr
	response.PlayerName = player.RoleMoudle.Name
	response.FightValue = G_SimpleMgr.Get_FightValue(req.PlayerID)
	response.RetCode = msg.RE_SUCCESS

	if player.pSimpleInfo == nil {
		gamelog.Error("Hand_PlayerEnterGame Error : player.pSimpleInfo == nil")
	} else {
		response.GuildID = player.pSimpleInfo.GuildID
	}

	gamelog.Info("message: user_enter_game : %s", response.PlayerName)

	if player.pSimpleInfo.LoginDay == 0 {
		SendLogNotify(req.PlayerID, EVENT_LOGIN_GAME, 0, player.GetLevel(), player.GetVipLevel(), 1, 0)
	} else {
		SendLogNotify(req.PlayerID, EVENT_LOGIN_GAME, 0, player.GetLevel(), player.GetVipLevel(), 0, 0)
	}

	//! 玩家登陆
	if player.IsTodayLogin() == false {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_USER_LOGIN, 1)
		player.ActivityModule.AddLoginDay()
	}

	G_SimpleMgr.Set_LoginDay(req.PlayerID, utility.GetCurDay())
	G_SimpleMgr.DB_SetLoginIp(req.PlayerID, r.Host)
}

func Hand_LeaveGame(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_LeaveGameSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_PlayerLeaveGame : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_LeaveGameSvr_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	sessionmgr.DeleteSessionKey(req.PlayerID)
	var player *TPlayer = GetPlayerByID(req.PlayerID)
	if player == nil {
		gamelog.Error("Hand_PlayerLeaveGame: cannot find the player info!!!!")
		response.RetCode = msg.RE_INVALID_PLAYERID
		return
	}

	player.OnPlayerOffline(req.PlayerID)
	response.RetCode = msg.RE_SUCCESS
	return
}

//
func Hand_CreatePlayer(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_CreateNewPlayerReq
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_CreateNewPlayer : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_CreateNewPlayerAck
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if (req.AccountID <= 0) || (req.HeroID <= 0) {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if !sessionmgr.CheckSessionKey(req.AccountID, req.SessionKey) {
		response.RetCode = msg.RE_INVALID_SESSIONKEY
		return
	}

	response.RetCode = CheckPlayerName(req.PlayerName)
	if response.RetCode != msg.RE_SUCCESS {
		return
	}

	var pSimpleInfo *TSimpleInfo = nil
	response.RetCode, pSimpleInfo = CheckCreatePlayer(req.AccountID, req.PlayerName)
	if response.RetCode != msg.RE_SUCCESS {
		return
	}

	player, _ := CreatePlayer(pSimpleInfo.PlayerID, req.PlayerName, req.HeroID)
	if player == nil {
		gamelog.Error("Create Player Failed player == nil !!!")
		response.RetCode = msg.RE_UNKNOWN_ERR
		return
	}

	player.pSimpleInfo = pSimpleInfo
	player.OnCreate(pSimpleInfo.PlayerID)
	pSimpleInfo.Quality = player.HeroMoudle.CurHeros[0].Quality
	pSimpleInfo.Level = player.HeroMoudle.CurHeros[0].Level
	pSimpleInfo.Name = req.PlayerName
	pSimpleInfo.HeroID = req.HeroID
	pSimpleInfo.FightValue = player.HeroMoudle.CalcFightValue(nil)
	pSimpleInfo.ChannelID = req.ChannelID
	response.PlayerID = pSimpleInfo.PlayerID
	response.RetCode = msg.RE_SUCCESS
	G_LevelRanker.SetRankItem(pSimpleInfo.PlayerID, pSimpleInfo.Level)
	G_FightRanker.SetRankItemEx(pSimpleInfo.PlayerID, 0, int(pSimpleInfo.FightValue))
	mongodb.InsertToDB("PlayerSimple", pSimpleInfo)
	IncRegisterCnt()
	return
}

func Hand_QueryServerTime(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_QueryServerTime_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryServerTime : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_QueryServerTime_Ack
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

	response.RetCode = msg.RE_SUCCESS
	response.SvrTime = utility.GetCurTime()
}

func CheckCreatePlayer(accountid int32, name string) (int, *TSimpleInfo) {
	G_SimpleMgr.SimpleLock.Lock()
	defer G_SimpleMgr.SimpleLock.Unlock()

	_, ok := G_SimpleMgr.SimpleList[accountid]
	if ok {
		gamelog.Error("CheckCreatePlayer Error : Repeat Create Role !!!!")
		return msg.RE_ALEADY_HAVE_ROLE, nil
	}

	_, ok = G_SimpleMgr.NameIDMap[name]
	if ok {
		gamelog.Error("CheckCreatePlayer Error : Name :%s has been used !!!!!", name)
		return msg.RE_ROLE_NAME_EXIST, nil
	}

	pInfo := new(TSimpleInfo)
	pInfo.AccountID = accountid
	pInfo.PlayerID = accountid
	pInfo.isOnline = false
	pInfo.Name = name

	G_SimpleMgr.SimpleList[pInfo.PlayerID] = pInfo
	G_SimpleMgr.NameIDMap[pInfo.Name] = pInfo.PlayerID

	return msg.RE_SUCCESS, pInfo
}

func CheckPlayerName(name string) int {
	if len(name) <= 0 {
		return msg.RE_INVALID_NAME
	}

	if len(name) > 20 {
		return msg.RE_INVALID_NAME
	}

	return msg.RE_SUCCESS
}

func CheckUserIsLogin(accountid int32, playerid int32, loginkey string) bool {
	var verifyuserReq msg.MSG_VerifyUserLogin_Req
	verifyuserReq.AccountID = accountid
	verifyuserReq.PlayerID = playerid
	verifyuserReq.LoginKey = loginkey
	verifyuserReq.SvrID = GetCurServerID()
	verifyuserBufferReq, _ := json.Marshal(&verifyuserReq)
	resp, err := http.Post(appconfig.VerifyUserLoginUrl, "text/HTML", bytes.NewReader(verifyuserBufferReq))
	if err != nil {
		gamelog.Error("CheckUserIsLogin post failed error : %s", err.Error())
		return false
	}

	verifyuserBufferAck := make([]byte, resp.ContentLength)
	resp.Body.Read(verifyuserBufferAck)
	resp.Body.Close()

	var verifyuserReqAck msg.MSG_VerifyUserLogin_Ack
	json.Unmarshal(verifyuserBufferAck, &verifyuserReqAck)
	if verifyuserReqAck.RetCode == msg.RE_SUCCESS {
		return true
	}

	return false
}

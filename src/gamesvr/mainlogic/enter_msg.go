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
	"time"
	"utility"
)

func Hand_PlayerLoginGame(w http.ResponseWriter, r *http.Request) {
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
	bcheck := true //CheckUserIsLogin(req.AccountID, "", req.LoginKey)
	if !bcheck {
		response.RetCode = msg.RE_INVALID_LOGINKEY
		gamelog.Error("Invalid Login key!!!!!")
		return
	}
	response.SessionKey = sessionmgr.NewSessionKey()
	sessionmgr.AddSessionKey(req.AccountID, response.SessionKey)
	response.PlayerID = GetPlayerIDByAccountID(req.AccountID)
	response.RetCode = msg.RE_SUCCESS
}

func Hand_GetLoginData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_EnterGameSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetLoginData : Unmarshal error!!!!")
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

	var pPlayer *TPlayer = GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		gamelog.Error("Hand_GetLoginData Error : Invalid Playerid :%d", req.PlayerID)
		return
	}

}

func Hand_PlayerEnterGame(w http.ResponseWriter, r *http.Request) {
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

	var pPlayer *TPlayer = GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		//如果内存中没有，查看数据库中是否存在
		pPlayer = LoadPlayerFromDB(req.PlayerID)
		if pPlayer == nil {
			gamelog.Error("Hand_PlayerEnterGame Error : Invalid Playerid :%d", req.PlayerID)
			return
		}
	}

	pPlayer.OnPlayerOnline(req.PlayerID)
	response.SvrTime = time.Now().Unix()

	response.ChatSvrAddr = appconfig.ChatSvrAddr
	response.PlayerName = pPlayer.RoleMoudle.Name
	response.FightValue = G_SimpleMgr.Get_FightValue(req.PlayerID)
	response.RetCode = msg.RE_SUCCESS

	if pPlayer.pSimpleInfo == nil {
		gamelog.Error("Hand_PlayerEnterGame Error : pPlayer.pSimpleInfo == nil")
	} else {
		response.GuildID = pPlayer.pSimpleInfo.GuildID
	}

	gamelog.Info("message: user_enter_game : %s", response.PlayerName)

	//! 玩家登陆
	if pPlayer.IsTodayLogin() == false {
		pPlayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_USER_LOGIN, 1)
		pPlayer.ActivityModule.AddLoginDay()
	}

	G_SimpleMgr.Set_LoginDay(req.PlayerID, utility.GetCurDay())
}

func Hand_PlayerLeaveGame(w http.ResponseWriter, r *http.Request) {
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
	var pPlayer *TPlayer = GetPlayerByID(req.PlayerID)
	if pPlayer == nil {
		gamelog.Error("Hand_PlayerLeaveGame: cannot find the player info!!!!")
		response.RetCode = msg.RE_INVALID_PLAYERID
		return
	}

	pPlayer.OnPlayerOffline(req.PlayerID)

	response.RetCode = msg.RE_SUCCESS
	return
}

//
func Hand_CreateNewPlayer(w http.ResponseWriter, r *http.Request) {
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

	pPlayer, _ := CreatePlayer(pSimpleInfo.PlayerID, req.PlayerName, req.HeroID)
	if pPlayer == nil {
		gamelog.Error("Create Player Failed pPlayer == nil !!!")
		response.RetCode = msg.RE_UNKNOWN_ERR
		return
	}

	pPlayer.pSimpleInfo = pSimpleInfo
	pPlayer.OnCreate(pSimpleInfo.PlayerID)
	pSimpleInfo.Quality = pPlayer.HeroMoudle.CurHeros[0].Quality
	pSimpleInfo.Level = pPlayer.HeroMoudle.CurHeros[0].Level
	pSimpleInfo.Name = req.PlayerName
	pSimpleInfo.HeroID = req.HeroID
	pSimpleInfo.FightValue = pPlayer.HeroMoudle.CalcFightValue(nil)
	response.PlayerID = pSimpleInfo.PlayerID
	response.RetCode = msg.RE_SUCCESS
	G_LevelRanker.SetRankItem(pSimpleInfo.PlayerID, pSimpleInfo.Level)
	G_FightRanker.SetRankItem(pSimpleInfo.PlayerID, pSimpleInfo.FightValue)

	if false == mongodb.InsertToDB(appconfig.GameDbName, "PlayerSimple", pSimpleInfo) {
		gamelog.Error("Hand_CreateNewPlayer Error: Insert to PlayserSimple Failed!!!")
		return
	}

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

	if !sessionmgr.CheckSessionKey(req.PlayerID, req.SessionKey) {
		response.RetCode = msg.RE_INVALID_SESSIONKEY
		return
	}

	response.RetCode = msg.RE_SUCCESS
	response.SvrTime = time.Now().Unix()
}

func CheckCreatePlayer(accountid int, name string) (int, *TSimpleInfo) {
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

	G_SimpleMgr.SimpleList[accountid] = pInfo
	G_SimpleMgr.NameIDMap[pInfo.Name] = pInfo.PlayerID

	return msg.RE_SUCCESS, pInfo
}

func CheckPlayerName(name string) int {
	if len(name) <= 0 {
		return msg.RE_INVALID_NAME
	}

	return msg.RE_SUCCESS
}

func CheckUserIsLogin(accountid int, accountname string, loginkey string) bool {
	var verifyuserReq msg.MSG_VerifyUserLogin_Req
	verifyuserReq.AccountID = accountid
	verifyuserReq.AccountName = accountname
	verifyuserReq.LoginKey = loginkey
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

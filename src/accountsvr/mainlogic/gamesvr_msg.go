package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"mongodb"
	"msg"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//处理游戏服力器的注册请求
func Handle_RegisterGameSvr(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	var buffer []byte

	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_RegToAccountSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_RegisterGameSvr : Unmarshal error!!!!")
		return
	}
	UpdateGameSvrInfo(req.SvrID, req.SvrName, req.SvrOuterAddr, req.SvrInnerAddr)
	var response msg.MSG_RegToAccountSvr_Ack
	response.RetCode = msg.RE_SUCCESS
	b, _ := json.Marshal(&response)
	w.Write(b)
}

//获取游戏公告
func Handle_GetGamePublic(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	var buffer []byte

	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_RegToAccountSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_RegisterGameSvr : Unmarshal error!!!!")
		return
	}
	UpdateGameSvrInfo(req.SvrID, req.SvrName, req.SvrOuterAddr, req.SvrInnerAddr)
	var response msg.MSG_RegToAccountSvr_Ack
	response.RetCode = msg.RE_SUCCESS
	b, _ := json.Marshal(&response)
	w.Write(b)
}

func Handle_SetGameSvrState(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_SetGameSvrState_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_SetGameSvrState : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_SetGameSvrState_Ack
	response.RetCode = msg.RE_INVALID_PARAM
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.SvrID <= 0 || req.SvrID >= 10000 {
		gamelog.Error("Handle_SetGameSvrState Error: Invalid SvrID:%d", req.SvrID)
		return
	}

	G_ServerList[req.SvrID].SvrState = req.SvrState
	G_ServerList[req.SvrID].SvrDefault = req.SvrDefault
	DB_UpdateSvrInfo(req.SvrID, G_ServerList[req.SvrID])
	response.RetCode = msg.RE_SUCCESS

	if G_ServerList[req.SvrID].SvrDefault == 1 {
		if G_RecommendID > 0 {
			G_ServerList[G_RecommendID].SvrDefault = 0
			DB_UpdateSvrInfo(G_RecommendID, G_ServerList[G_RecommendID])
		}
		G_RecommendID = req.SvrID
	} else if G_RecommendID == req.SvrID {
		G_RecommendID = 0
	}

	return
}

func Handle_GmServerList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetServerList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GetServerList : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetServerList_Ack
	response.RetCode = msg.RE_INVALID_PARAM
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	nCount := len(G_ServerList)
	response.SvrList = make([]msg.ServerNode, 0, 10)
	for i := 0; i < nCount; i++ {
		if G_ServerList[i].SvrID != 0 {
			response.SvrList = append(response.SvrList, msg.ServerNode{G_ServerList[i].SvrID,
				G_ServerList[i].SvrName,
				G_ServerList[i].SvrState,
				G_ServerList[i].SvrDefault,
				G_ServerList[i].SvrOutAddr})
		}
	}

	response.RetCode = msg.RE_SUCCESS
	return
}

func Handle_GmLogin(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GmLogin_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GmLogin : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GmLogin_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//检查是否具有GM操作权限
	if false == appconfig.CheckGmRight(req.SessionID, req.SessionKey, r.RemoteAddr[:strings.IndexRune(r.RemoteAddr, ':')]) {
		gamelog.Error("Handle_GmLogin Error Invalid Gm request!!!")
		response.RetCode = msg.RE_INVALID_NAME
		return
	}

	response.RetCode = msg.RE_SUCCESS
	return
}

func QueryAccountID(name string, svrid int32) int32 {
	var req msg.MSG_QueryAccountID_Req
	req.Name = name
	b, _ := json.Marshal(req)
	requrl := "http://" + GetGameSvrOutAddr(svrid) + "/query_account_id"
	http.DefaultClient.Timeout = 2 * time.Second
	httpret, err := http.Post(requrl, "text/HTML", bytes.NewReader(b))
	if err != nil {
		gamelog.Error("QueryAccountID Error:  err : %s !!!!", err.Error())
		return 0
	}

	buffer := make([]byte, httpret.ContentLength)
	httpret.Body.Read(buffer)
	httpret.Body.Close()

	var ack msg.MSG_QueryAccountID_Ack
	err = json.Unmarshal(buffer, &ack)
	if err != nil {
		gamelog.Error("QueryAccountID Error: Error: %s", err.Error())
		return 0
	}

	return ack.AccountID
}

func Handle_GmEnableAccount(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GmEnableAccount_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GmEnableAccount : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GmEnableAccount_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.PlayerID == 0 {
		//查询角色的账号ID
		req.PlayerID = QueryAccountID(req.RoleName, req.SvrID)
	}

	if req.PlayerID <= 0 {
		gamelog.Error("Handle_GmEnableAccount : Cant find the player with name:%s, id:%d!!!!", req.RoleName, req.SvrID)
		return
	}

	pAccount, _ := G_AccountMgr.GetAccountByID(req.PlayerID)
	if pAccount == nil {
		gamelog.Error("Handle_GmEnableAccount : Invalid AccountID:%d!!!!", req.PlayerID)
		response.RetCode = msg.RE_INVALID_ACCOUNTID
		return
	}

	pAccount.Enable = req.Enable
	mongodb.UpdateToDB("Account", &bson.M{"_id": req.PlayerID}, &bson.M{"$set": bson.M{"enable": pAccount.Enable}})
	response.RetCode = msg.RE_SUCCESS
	return
}

func CheckGameStateRoutine() {
	var req msg.MSG_QueryAccountID_Req
	req.Name = "123"
	b, _ := json.Marshal(req)

	for i := 0; i < 10000; i++ {
		if G_ServerList[i].SvrID <= 0 {
			continue
		}

		if G_ServerList[i].isSvrOK == true {
			continue
		}

		requrl := "http://" + GetGameSvrOutAddr(G_ServerList[i].SvrID) + "/query_account_id"
		http.DefaultClient.Timeout = 1 * time.Second
		response, err := http.Post(requrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("CheckGameStateRoutine Error , err : %s !!!!", err.Error())
			continue
		} else {
			G_ServerList[i].isSvrOK = true
		}

		response.Body.Close()
	}
}

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
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

func Handle_SetGamesvrFlag(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_SetGameSvrFlag_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_SetGamesvrFlag : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_SetGameSvrFlag_Ack
	response.RetCode = msg.RE_SUCCESS
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	if req.SvrID <= 0 || req.SvrID >= 10000 {
		gamelog.Error("Handle_SetGamesvrFlag Error: Invalid SvrID:%d", req.SvrID)
		return
	}

	G_ServerList[req.SvrID].SvrFlag = req.Flag
	DB_UpdateSvrState(req.SvrID, req.Flag)
	response.RetCode = msg.RE_SUCCESS
	return
}

func Handle_GetServerList(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetServerList_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_GetServerList : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetServerList_Ack
	response.RetCode = msg.RE_SUCCESS
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
				G_ServerList[i].SvrFlag,
				G_ServerList[i].svrOutAddr})
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

	response.RetCode = msg.RE_SUCCESS
	return
}

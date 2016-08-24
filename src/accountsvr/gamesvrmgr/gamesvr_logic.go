package gamesvrmgr

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

//处理游戏服力器的注册请求
func Handle_RegisterGameSvr(w http.ResponseWriter, r *http.Request) {
	var buffer []byte
	buffer = make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_RegisterGameSvr_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Handle_RegisterGameSvr : Unmarshal error!!!!")
		return
	}

	UpdateGameSvrInfo(req.ServerDomainID, req.ServerDomainName, req.ServerOuterAddr, req.ServerInnerAddr)

	var response msg.MSG_RegisterGameSvr_Ack
	response.RetCode = msg.RE_SUCCESS

	b, _ := json.Marshal(&response)
	w.Write(b)
}

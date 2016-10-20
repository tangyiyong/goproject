/***********************************************************************
* @ 游戏服msg的处理
* @ brief
    1、gamesvr先通知SDK进程，建立新充值订单

    2、第三方充值信息到达后，验证是否为有效订单

* @ author zhoumf
* @ date 2016-8-18
***********************************************************************/
package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

func Hand_RegGamesvrAddr(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RegToSdkSvr_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("HandSvr_GamesvrAddr unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.SDKMsg_GamesvrAddr_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	UpdateGameSvrInfo(req.SvrID, req.SvrName, req.SvrOuterAddr, req.SvrInnerAddr)
	response.RetCode = msg.RE_SUCCESS
}

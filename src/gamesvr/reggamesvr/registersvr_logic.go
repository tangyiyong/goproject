package reggamesvr

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
	"strconv"
	"time"
)

func RegisterToSvr() {
	//注册到账号服
	go RegisterToAccountSvr()

	//注册到跨服服
	go RegisterToCrossSvr()

	//注册到SDK服
	//go RegisterToSdkSvr()
}

//注册到账号服务器
func RegisterToAccountSvr() {
	var registerReq msg.MSG_RegToAccountSvr_Req
	registerReq.SvrID = int32(appconfig.GameSvrID)
	registerReq.SvrOpenTime = int32(appconfig.GameOpenSvrTime)
	registerReq.SvrOuterAddr = appconfig.GameSvrOuterIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	registerReq.SvrInnerAddr = appconfig.GameSvrInnerIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	b, _ := json.Marshal(registerReq)

	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(appconfig.RegToAccountSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to Account Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		buffer := make([]byte, response.ContentLength)
		response.Body.Read(buffer)
		response.Body.Close()
		var ack msg.MSG_RegToAccountSvr_Ack
		err = json.Unmarshal(buffer, &ack)
		appconfig.GameSvrName = ack.SvrName
		response.Body.Close()
		return
	}

}

//注册到跨服服务器
func RegisterToCrossSvr() {
	var registerReq msg.MSG_RegToCrossSvr_Req
	registerReq.SvrID = int32(appconfig.GameSvrID)
	registerReq.SvrName = appconfig.GameSvrName
	registerReq.SvrOuterAddr = appconfig.GameSvrOuterIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	registerReq.SvrInnerAddr = appconfig.GameSvrInnerIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	b, _ := json.Marshal(registerReq)

	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(appconfig.RegToCrossSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to Account Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		response.Body.Close()
		return
	}
}

func RegisterToSdkSvr() {
	var req msg.MSG_RegToSdkSvr_Req
	req.SvrID = int32(appconfig.GameSvrID)
	req.SvrInnerAddr = appconfig.GameSvrInnerIp
	req.SvrOuterAddr = appconfig.GameSvrOuterIp
	req.SvrName = appconfig.GameSvrName
	b, _ := json.Marshal(req)
	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(appconfig.RegToSdkSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to SDK Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		response.Body.Close()
		return
	}
}

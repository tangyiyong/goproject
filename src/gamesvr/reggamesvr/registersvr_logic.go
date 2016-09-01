package reggamesvr

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"fmt"
	"gamelog"
	"msg"
	"net/http"
	"strconv"
	"time"
)

func RegisterToSvr() {
	//注册到账号服
	go RegisterToAccountRoutine()

	//注册到跨服服
	go RegisterToCrossRoutine()

	//注册到SDK服
	//go RegisterToSdkSvr()
}

//注册到账号服务器
func RegisterToAccountRoutine() {
	var registerReq msg.MSG_RegToAccountSvr_Req
	registerReq.ServerDomainID = int32(appconfig.DomainID)
	registerReq.ServerDomainName = appconfig.DomainName
	registerReq.ServerOuterAddr = appconfig.GameSvrOuterIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	registerReq.ServerInnerAddr = appconfig.GameSvrInnerIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	b, _ := json.Marshal(registerReq)

	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(appconfig.RegToAccountSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to Account Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		response.Body.Close()
		time.Sleep(60 * time.Second)
	}

}

//注册到跨服服务器
func RegisterToCrossRoutine() {
	var registerReq msg.MSG_RegToCrossSvr_Req
	registerReq.ServerDomainID = int32(appconfig.DomainID)
	registerReq.ServerDomainName = appconfig.DomainName
	registerReq.ServerOuterAddr = appconfig.GameSvrOuterIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	registerReq.ServerInnerAddr = appconfig.GameSvrInnerIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
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
		time.Sleep(60 * time.Second)
	}
}

func RegisterToSdkSvr() {
	PorstUrl := fmt.Sprintf("http://%s:%d/reg_gamesvr_addr", appconfig.SdkSvrInnerIp, appconfig.SdkSvrPort)
	var req msg.SDKMsg_GamesvrAddr_Req
	req.GamesvrID = appconfig.DomainID
	req.Url = fmt.Sprintf("http://%s:%d/", appconfig.GameSvrInnerIp, appconfig.GameSvrPort)
	b, _ := json.Marshal(req)
	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(PorstUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to SDK Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		response.Body.Close()
		return
	}
}

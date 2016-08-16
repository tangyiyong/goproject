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

func RegisterToAccountSvr() {
	go RegisterRoutine()
}

//注册到账号服务器
func RegisterRoutine() {
	var registerReq msg.MSG_RegisterGameSvr_Req
	registerReq.ServerDomainID = appconfig.DomainID
	registerReq.ServerDomainName = appconfig.DomainName
	registerReq.ServerOuterAddr = appconfig.GameSvrOuterIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	registerReq.ServerInnerAddr = appconfig.GameSvrInnerIp + ":" + strconv.Itoa(appconfig.GameSvrPort)
	b, _ := json.Marshal(registerReq)
	bConnectOK := false

	for {
		http.DefaultClient.Timeout = 2 * time.Second
		response, err := http.Post(appconfig.RegToAccountSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to Account Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		} else {
			bConnectOK = true
		}
		response.Body.Close()

		response, err = http.Post(appconfig.RegToCrossSvrUrl, "text/HTML", bytes.NewReader(b))
		if err != nil {
			gamelog.Error("Register to Cross Server failed, err : %s !!!!", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		response.Body.Close()

		if bConnectOK == true {
			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
	}

}

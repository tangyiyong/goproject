/***********************************************************************
* @ 内部游戏服的地址列表
* @ brief
    1、SDK进程负责接收所有第三方消息，验证后转发至对应的gamesvr

    2、须预先加载服务器列表

* @ author zhoumf
* @ date 2016-8-16
***********************************************************************/
package sdklogic

import (
	"bytes"
	"encoding/json"
	"gamelog"
	"net/http"
	"strconv"
	"utility"
)

var (
	SvrAddr_PATH = utility.GetCurrPath() + "svr_addr.csv"

	SvrID_Addr = make(map[int]string)
)

func LoadSvrAddrList() {
	records, err := utility.LoadCsv(SvrAddr_PATH)
	if err != nil {
		gamelog.Error("LoadSvrAddrList : %s", err.Error())
		return
	}

	for i := 0; i < len(records); i++ {
		id, _ := strconv.Atoi(records[i][0])
		SvrID_Addr[id] = records[i][1]
	}
}
func RegisterGamesvrAddr(svrID int, url string) {
	SvrID_Addr[svrID] = url
}

// strKey = "sdk_recharge_info"
func RelayToGamesvr(svrId int, strKey string, pMsg interface{}) {
	url, ok := SvrID_Addr[svrId]
	if ok {
		url += strKey
		data, _ := json.Marshal(pMsg)
		if _, err := postServerReq(url, data); err != nil {
			gamelog.Error("RelayToGamesvr--PostServerReq: svrId(%d) %s", svrId, err.Error())
		}
	} else {
		gamelog.Error("RelayToGamesvr: svrId(%d)", svrId)
	}
}
func postServerReq(url string, buf []byte) ([]byte, error) {
	resp, err := http.Post(url, "text/HTML", bytes.NewReader(buf))
	backBuf := make([]byte, resp.ContentLength)
	resp.Body.Read(backBuf)
	resp.Body.Close()
	return backBuf, err
}

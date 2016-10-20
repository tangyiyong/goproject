package mainlogic

import (
	"bytes"
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
	"time"
)

//! 玩家请求积分赛目标信息
//cross_query_score_target
func Hand_GetScoreTarget(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_CrossQueryScoreTarget_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreTarget : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_CrossQueryScoreTarget_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	for i := 0; i < 3; i++ {
		svrUrl := GetSelectSvrAddr()
		if len(svrUrl) <= 0 {
			gamelog.Error("Hand_GetScoreTarget Error Invalid SvrUrl!!")
			return
		}
		response.TargetList[i], _ = GetScoreTargetItem(svrUrl)
	}

}

func GetScoreTargetItem(addr string) (msg.MSG_Target, bool) {
	var reqUrl = "http://" + addr + "/select_target_player"
	var GameSelectPlayerReq msg.MSG_GameSelectPlayer_Req
	b, _ := json.Marshal(GameSelectPlayerReq)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(reqUrl, "text/HTML", bytes.NewReader(b))
	if err != nil || httpret == nil {
		gamelog.Error("GetScoreTarget failed, err : %s !!!!", err.Error())
		return msg.MSG_Target{}, false
	}

	buffer := make([]byte, httpret.ContentLength)
	httpret.Body.Read(buffer)
	httpret.Body.Close()

	var GameSelectPlayerAck msg.MSG_GameSelectPlayer_Ack
	err = json.Unmarshal(buffer, &GameSelectPlayerAck)
	if err != nil {
		gamelog.Error("GetScoreTarget  Unmarshal fail, Error: %s", err.Error())
		return msg.MSG_Target{}, false
	}

	return GameSelectPlayerAck.Target, true
}

func Handle_GetFightTarget(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetFightTarget_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Handle_GetFightTarget : Unmarshal fail, Error: %s", err.Error())
		return
	}

	addr := GetGameSvrFightTarAddr(req.SvrID)
	if len(addr) <= 0 {
		gamelog.Error("Handle_GetFightTarget : Cant get the addr of svr %d", req.SvrID)
		return
	}

	b, _ := json.Marshal(req)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(addr, "text/HTML", bytes.NewReader(b))
	if err != nil || httpret == nil {
		gamelog.Error("Handle_GetFightTarget failed, err : %s !!!!", err.Error())
		return
	}

	buffer = make([]byte, httpret.ContentLength)
	httpret.Body.Read(buffer)
	httpret.Body.Close()
	w.Write(buffer)

	return
}

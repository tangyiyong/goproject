package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"msg"
	"net/http"
	"sync"
)

var TestLock sync.Mutex

func PostServerReq(url string, body io.Reader) ([]byte, error) {
	TestLock.Lock()
	//t1 := time.Now().UnixNano()
	resp, err := http.Post(url, "text/HTML", body)
	if resp == nil {
		return nil, err
	}
	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()
	//fmt.Println("t:", time.Now().UnixNano()-t1)
	TestLock.Unlock()

	return buffer, err
}
func PostMsg(url string, msg interface{}) []byte {
	b, _ := json.Marshal(msg)
	buf, err := PostServerReq(url, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return buf
}

func TestGetScoreTarget(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/get_score_target"
	var req msg.MSG_GetScoreTarget_Req
	req.SessionKey = sessionkey
	req.PlayerID = playerid
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ack msg.MSG_GetScoreTarget_Ack
	json.Unmarshal(buffer, &ack)
	if ack.RetCode != 0 {
		return
	}

	fmt.Println("%v", ack)

	return
}

func TestGetScoreRank(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/get_score_rank"
	var req msg.MSG_GetScoreRank_Req
	req.SessionKey = sessionkey
	req.PlayerID = playerid
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ack msg.MSG_GetScoreRank_Ack
	json.Unmarshal(buffer, &ack)
	if ack.RetCode != 0 {
		return
	}

	fmt.Println("%v", ack)

	return
}

func TestGetMoney(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/test_get_money"
	var req msg.MSG_GetTestMoney_Req
	req.SessionKey = sessionkey
	req.PlayerID = playerid
	b, _ := json.Marshal(req)
	_, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func TestGetAction(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/test_get_action"
	var req msg.MSG_GetTestAction_Req
	req.SessionKey = sessionkey
	req.PlayerID = playerid
	b, _ := json.Marshal(req)
	_, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func TestAddSvrAward(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/add_svr_award"
	var req msg.MSG_SvrAward_Add_Req
	req.Value = []string{"渣渣"}
	req.ItemLst = []msg.MSG_ItemData{{1, 1}}
	PostMsg(reqUrl, req)
}
func TestDelSvrAward(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/del_svr_award"
	var req msg.MSG_SvrAward_Del_Req
	req.ID = 1
	PostMsg(reqUrl, req)
}
func TestAwardCenterQuery(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/query_award_center"
	var req msg.MSG_AwardCenter_Query_Req
	req.PlayerID = playerid
	req.SessionKey = sessionkey
	PostMsg(reqUrl, req)
}
func TestAwardCenterGet(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/get_award_center"
	var req msg.MSG_AwardCenter_Get_Req
	req.PlayerID = playerid
	req.SessionKey = sessionkey
	req.AwardID = 3
	PostMsg(reqUrl, req)
}

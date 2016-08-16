// main
package main

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"fmt"
	"gamelog"
	"io"
	"msg"
	"net/http"
	"strconv"
	"sync"
	"time"
	"utility"
)

var (
	registerurl string = "http://127.0.0.1:8081/register"
	loginurl    string = "http://127.0.0.1:8081/login"
	serverlist  string = "http://127.0.0.1:8081:8081/serverlist"
	Password    string = "123456"
)

var TestLock sync.Mutex

func TestServerList() {

}

func TestGameSvr() {
	appconfig.LoadConfig()
	gamelog.InitLogger("httptest", true)

	for i := 3; i < 20; i++ {
		time.Sleep(time.Millisecond)
		go TestUser(i)
	}

	utility.StartConsoleWait()
}

func TestUser(i int) {
	Name := "kkkk" + strconv.Itoa(i)
	//根据账号名和密码，注册账号

	if !userRegister(Name, Name) {
		gamelog.Error("注山账号失败!!")
		return
	}

	//登录账号服务器
	accountid, loginkey, bRet := userLogin(Name, Name)
	if !bRet {
		gamelog.Error("登录账号服务器失败!!")
		return
	}

	//用账号服返回的loginkey 去登录游戏服
	retCode, playerid, sessionkey := userLoginGame(accountid, loginkey)
	if retCode != 0 {
		gamelog.Error("登录游戏服务器失败!!")
		return
	}

	if (playerid <= 0) && (retCode == 0) {
		//如果没有角色，就创建角色
		playerid, bRet = userCreatePlayer(accountid, sessionkey)
		if !bRet {
			gamelog.Error("创建角色失败!!!!")
			return
		}
	}

	//用游戏服返回的Sessionkey 去登录游戏服

	if !userEnterGame(playerid, sessionkey) {
		gamelog.Error("进入游戏服失败!!!!")
		return
	}
	/*
		TestGetAction(playerid, sessionkey)
		TestGetMoney(playerid, sessionkey)
		TestUpLevel(playerid, sessionkey)

		TestGetScoreTarget(playerid, sessionkey)
		TestGetScoreRank(playerid, sessionkey)
	*/

	//buyItem(playerid, sessionkey)

	TestAddSvrAward(playerid, sessionkey)
	TestAwardCenterQuery(playerid, sessionkey)
	TestAwardCenterGet(playerid, sessionkey)
	TestDelSvrAward(playerid, sessionkey)
}

func userRegister(name string, password string) bool {
	var req msg.MSG_RegAccount_Req
	req.Name = name
	req.Password = password
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(registerurl, bytes.NewReader(b))
	if err != nil {
		return false
	}
	var ack msg.MSG_RegAccount_Ack
	json.Unmarshal(buffer, &ack)

	return true
}

func userLogin(name string, password string) (int, string, bool) {
	var req msg.MSG_Login_Req
	req.Name = name
	req.Password = password
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(loginurl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return 0, "", false
	}
	var ack msg.MSG_Login_Ack
	err = json.Unmarshal(buffer, &ack)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ack.RetCode != 0 {
		return 0, "", false
	}

	return ack.AccountID, ack.LoginKey, true
}

func userLoginGame(accountid int, loginkey string) (int, int, string) {
	reqUrl := "http://127.0.0.1:8082/user_login_game"
	var req msg.MSG_LoginGameSvr_Req
	req.AccountID = accountid
	req.LoginKey = loginkey
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return msg.RE_UNKNOWN_ERR, 0, ""
	}

	var ack msg.MSG_LoginGameSvr_Ack
	json.Unmarshal(buffer, &ack)
	return ack.RetCode, ack.PlayerID, ack.SessionKey
}

func userCreatePlayer(accountid int, loginkey string) (int, bool) {
	reqUrl := "http://127.0.0.1:8082/create_new_player"
	var req msg.MSG_CreateNewPlayerReq
	req.AccountID = accountid
	req.SessionKey = loginkey
	req.PlayerName = "name" + strconv.Itoa(accountid)
	req.HeroID = 3
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return 0, false
	}
	var ack msg.MSG_CreateNewPlayerAck
	json.Unmarshal(buffer, &ack)

	if ack.RetCode != 0 {
		fmt.Println("userCreatePlayer failed: ", ack)
	}

	return ack.PlayerID, true
}

func userEnterGame(playerid int, sessionkey string) bool {
	reqUrl := "http://127.0.0.1:8082/user_enter_game"
	var req msg.MSG_EnterGameSvr_Req
	req.PlayerID = playerid
	req.SessionKey = sessionkey
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ack msg.MSG_EnterGameSvr_Ack
	json.Unmarshal(buffer, &ack)
	if ack.RetCode != 0 {
		return false
	}

	return true
}

func buyItem(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/buy_goods"
	var req msg.MSG_BuyGoods_Req
	req.ID = 2
	req.Num = 1
	req.SessionKey = sessionkey
	req.PlayerID = playerid
	b, _ := json.Marshal(req)

	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var ack msg.MSG_EnterGameSvr_Ack

	json.Unmarshal(buffer, &ack)
	if ack.RetCode != 0 {
		return
	}

	return
}

func PostServerReq(url string, body io.Reader) ([]byte, error) {
	//TestLock.Lock()
	//t1 := time.Now().UnixNano()
	resp, err := http.Post(url, "text/HTML", body)
	buffer := make([]byte, resp.ContentLength)
	resp.Body.Read(buffer)
	resp.Body.Close()
	//fmt.Println("t:", time.Now().UnixNano()-t1)
	//TestLock.Unlock()

	return buffer, err
}
func PostMsg(url string, msg interface{}) {
	b, _ := json.Marshal(msg)
	_, err := PostServerReq(url, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
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

func TestUpLevel(playerid int, sessionkey string) {
	reqUrl := "http://127.0.0.1:8082/test_uplevel"
	var req msg.MSG_TestUpLevel_Req
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

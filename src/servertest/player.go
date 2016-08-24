// main
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamelog"
	"gamesvr/tcpclient"
	"msg"
	"strconv"
	"time"
)

var (
	registerurl string = "http://127.0.0.1:8081/register"
	loginurl    string = "http://127.0.0.1:8081/login"
	serverlist  string = "http://127.0.0.1:8081:8081/serverlist"
	Password    string = "123456"
)

type TPlayer struct {
	AccountName string
	Password    string
	AccountID   int32

	PlayerName string
	PlayerID   int32

	LoginKey   string
	SessoinKey string

	BattleSvrAddr string //战场服务器IP地址
	BatCamp       int
	BatClient     tcpclient.TCPClient
	EnterCode     int32 //进入码

	Heros   [6]msg.MSG_HeroObj
	IsEnter bool
}

var G_PlayerMgr map[string]*TPlayer

func InitPlayerMgr() {
	G_PlayerMgr = make(map[string]*TPlayer, 1)
}

func CreatePlayer(index int) *TPlayer {
	var pPlayer = new(TPlayer)
	pPlayer.AccountName = "acc" + strconv.Itoa(index)
	pPlayer.Password = "123"
	pPlayer.PlayerName = "name" + strconv.Itoa(index)
	G_PlayerMgr[pPlayer.AccountName] = pPlayer
	return pPlayer
}

func StartTest() {
	for _, v := range G_PlayerMgr {
		go v.TestRoutine()
	}
}

func (self *TPlayer) TestRoutine() {
	//根据账号名和密码，注册账号

	if !self.userRegister() {
		gamelog.Error("注山账号失败!!")
		return
	}

	//登录账号服务器
	bRet := self.userLogin()
	if !bRet {
		gamelog.Error("登录账号服务器失败!!")
		return
	}

	//用账号服返回的loginkey 去登录游戏服
	bRet = self.userLoginGame()
	if !bRet {
		gamelog.Error("登录游戏服务器失败!!")
		return
	}

	if self.PlayerID <= 0 {
		//如果没有角色，就创建角色
		bRet = self.userCreatePlayer()
		if !bRet {
			gamelog.Error("创建角色失败!!!!")
			return
		}
	}

	//用游戏服返回的Sessionkey 去登录游戏服
	if !self.userEnterGame() {
		gamelog.Error("进入游戏服失败!!!!")
		return
	}

	// self.TestUpLevel()
	// self.TestUpLevel()
	// self.TestUpLevel()
	// self.TestUpLevel()
	self.Create_Recharge_Order_2Gamesvr()
	self.Recharge_Syccess_2SDK()

	if !self.userSetBatCamp() {
		gamelog.Error("设置阵营失败!!!!")
		return
	}

	return

	if !self.userEnterBattle() {
		gamelog.Error("进入阵营失败!!!!")
		return
	}

	//if !self.userStartCarry() {
	//	gamelog.Error("进入阵营失败!!!!")
	//	return
	//}

	for {
		if self.IsEnter == true {
			self.userMove()
		}

		time.Sleep(200 * time.Millisecond)
	}

}

func (self *TPlayer) userRegister() bool {
	var req msg.MSG_RegAccount_Req
	req.Name = self.AccountName
	req.Password = self.Password
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(registerurl, bytes.NewReader(b))
	if err != nil {
		return false
	}
	var ack msg.MSG_RegAccount_Ack
	json.Unmarshal(buffer, &ack)

	return true
}

func (self *TPlayer) userLogin() bool {
	var req msg.MSG_Login_Req
	req.Name = self.AccountName
	req.Password = self.Password
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(loginurl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ack msg.MSG_Login_Ack
	err = json.Unmarshal(buffer, &ack)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ack.RetCode != 0 {
		return false
	}

	self.AccountID = ack.AccountID
	self.LoginKey = ack.LoginKey

	return true
}

func (self *TPlayer) userLoginGame() bool {
	reqUrl := "http://127.0.0.1:8082/user_login_game"
	var req msg.MSG_LoginGameSvr_Req
	req.AccountID = self.AccountID
	req.LoginKey = self.LoginKey
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	var ack msg.MSG_LoginGameSvr_Ack
	json.Unmarshal(buffer, &ack)

	if ack.RetCode != msg.RE_SUCCESS {
		return false
	}

	self.PlayerID = ack.PlayerID
	self.SessoinKey = ack.SessionKey
	return true
}

func (self *TPlayer) userCreatePlayer() bool {
	reqUrl := "http://127.0.0.1:8082/create_new_player"
	var req msg.MSG_CreateNewPlayerReq
	req.AccountID = self.AccountID
	req.SessionKey = self.SessoinKey
	req.PlayerName = self.PlayerName
	req.HeroID = 3
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ack msg.MSG_CreateNewPlayerAck
	json.Unmarshal(buffer, &ack)

	if ack.RetCode != 0 {
		fmt.Println("userCreatePlayer failed: ", ack)
		return false
	}

	return true
}

func (self *TPlayer) userEnterGame() bool {
	reqUrl := "http://127.0.0.1:8082/user_enter_game"
	var req msg.MSG_EnterGameSvr_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ack msg.MSG_EnterGameSvr_Ack
	json.Unmarshal(buffer, &ack)
	if ack.RetCode != 0 {
		fmt.Println("userEnterGame failed: ", ack)
		return false
	}

	return true
}

func (self *TPlayer) TestUpLevel() {
	reqUrl := "http://127.0.0.1:8082/test_uplevel_ten"
	var req msg.MSG_TestUpLevelTen_Req
	req.SessionKey = self.SessoinKey
	req.PlayerID = self.PlayerID
	b, _ := json.Marshal(req)
	_, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func (self *TPlayer) buyItem() {
	reqUrl := "http://127.0.0.1:8082/buy_goods"
	var req msg.MSG_BuyGoods_Req
	req.ID = 2
	req.Num = 1
	req.SessionKey = self.SessoinKey
	req.PlayerID = self.PlayerID
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

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
	//"time"
)

type TPlayer struct {
	AccountName string
	Password    string
	AccountID   int32

	GameSvrAddr string //游戏服地址
	ChatSvrAddr string //聊天服务器地址
	ChatClient  tcpclient.TCPClient
	PlayerName  string
	PlayerID    int32

	LoginKey   string
	SessoinKey string

	BattleSvrAddr string //战场服务器IP地址
	BatCamp       int8
	BatClient     tcpclient.TCPClient

	EnterCode int32 //进入码

	Heros   [6]msg.MSG_HeroObj
	IsEnter bool

	PackNo int32 //消息编号
}

var G_PlayerMgr map[string]*TPlayer

func InitPlayerMgr() {
	G_PlayerMgr = make(map[string]*TPlayer, 1)
}

func CreatePlayer(index int) *TPlayer {
	var player = new(TPlayer)
	player.AccountName = "acc" + strconv.Itoa(index)
	player.Password = "123"
	player.PlayerName = "name" + strconv.Itoa(index)
	G_PlayerMgr[player.AccountName] = player
	return player
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

	self.GetBagData()
	self.GetCopyData()
	self.GetBattleData()
	self.GetActivitylist()

	//if !self.userSetBatCamp() {
	//	gamelog.Error("设置阵营失败!!!!")
	//	return
	//}

	//if !self.userEnterBattle() {
	//	gamelog.Error("进入阵营失败!!!!")
	//	return
	//}

	//if !self.userStartCarry() {
	//	gamelog.Error("进入阵营失败!!!!")
	//	return
	//}

	//for {
	//	if self.IsEnter == true {
	//		self.userMove()
	//	}

	//	time.Sleep(200 * time.Millisecond)
	//}

}

func (self *TPlayer) userRegister() bool {
	var req msg.MSG_RegAccount_Req
	req.Name = self.AccountName
	req.Password = self.Password
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(accountip+"/register", bytes.NewReader(b))
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
	buffer, err := PostServerReq(accountip+"/login", bytes.NewReader(b))
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
	self.GameSvrAddr = "http://" + ack.LastSvrAddr

	return true
}

func (self *TPlayer) userLoginGame() bool {
	var req msg.MSG_LoginGameSvr_Req
	req.AccountID = self.AccountID
	req.LoginKey = self.LoginKey
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(self.GameSvrAddr+"/user_login_game", bytes.NewReader(b))
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
	var req msg.MSG_CreateNewPlayerReq
	req.AccountID = self.AccountID
	req.SessionKey = self.SessoinKey
	req.PlayerName = self.PlayerName
	req.HeroID = 3
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(self.GameSvrAddr+"/create_new_player", bytes.NewReader(b))
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

	self.PlayerID = ack.PlayerID

	return true
}

func (self *TPlayer) userEnterGame() bool {
	var req msg.MSG_EnterGameSvr_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(self.GameSvrAddr+"/user_enter_game", bytes.NewReader(b))
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

	self.ChatSvrAddr = ack.ChatSvrAddr

	self.ChatClient.ConType = tcpclient.CON_TYPE_CHAT
	self.ChatClient.ExtraData = self
	self.ChatClient.ConnectToSvr(self.ChatSvrAddr, 1)

	return true
}

func (self *TPlayer) userCheckIn() bool {
	var req msg.MSG_CheckIn_Req
	req.PlayerID = self.PlayerID
	req.GuildID = 0
	req.PlayerName = self.PlayerName
	b, _ := json.Marshal(req)

	if self == nil {
		fmt.Println("userCheckIn failed self == nil:")
		return false
	}

	if self.ChatClient.TcpConn == nil {
		fmt.Println("userCheckIn failed TcpConn == nil:")
		return false
	}

	self.ChatClient.TcpConn.WriteMsg(msg.MSG_CHECK_IN_REQ, 0, b)

	return true
}

func (self *TPlayer) TestUpLevel() {
	var req msg.MSG_TestUpLevelTen_Req
	req.SessionKey = self.SessoinKey
	req.PlayerID = self.PlayerID
	b, _ := json.Marshal(req)
	_, err := PostServerReq(self.GameSvrAddr+"/test_uplevel_ten", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}

func (self *TPlayer) buyItem() {
	var req msg.MSG_BuyGoods_Req
	req.ID = 2
	req.Num = 1
	req.SessionKey = self.SessoinKey
	req.PlayerID = self.PlayerID
	b, _ := json.Marshal(req)

	buffer, err := PostServerReq(self.GameSvrAddr+"/buy_goods", bytes.NewReader(b))
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

//消息:/
type MSG_GetData_Req struct {
	PlayerID   int32
	SessionKey string
}

func (self *TPlayer) GetBagData() bool {
	var req MSG_GetData_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	b, _ := json.Marshal(req)
	_, err := PostServerReq(self.GameSvrAddr+"/get_bag_data", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (self *TPlayer) GetGiftCode() bool {
	var req msg.MSG_RecvGiftCode_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	req.GiftCode = "afdsfasf"
	b, _ := json.Marshal(req)
	_, err := PostServerReq(self.GameSvrAddr+"/recv_gift_code", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (self *TPlayer) GetCopyData() bool {
	var req MSG_GetData_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	b, _ := json.Marshal(req)
	_, err := PostServerReq(self.GameSvrAddr+"/get_copy_data", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (self *TPlayer) GetBattleData() bool {
	var req MSG_GetData_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	b, _ := json.Marshal(req)
	_, err := PostServerReq(self.GameSvrAddr+"/get_battle_data", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func (self *TPlayer) GetActivitylist() bool {
	var req MSG_GetData_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	b, _ := json.Marshal(req)
	_, err := PostServerReq(self.GameSvrAddr+"/get_activity_list", bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

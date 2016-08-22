// main
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamelog"
	"gamesvr/tcpclient"
	"msg"
	"utility"
)

func (self *TPlayer) userSetBatCamp() bool {
	reqUrl := "http://127.0.0.1:8082/set_battlecamp"
	var req msg.MSG_SetBattleCamp_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	self.BatCamp = self.PlayerID%3 + 1
	req.BattleCamp = self.BatCamp
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ack msg.MSG_SetBattleCamp_Ack
	json.Unmarshal(buffer, &ack)

	if ack.RetCode != 0 {
		fmt.Println("userSetBatCamp failed: ", ack)
		return false
	}

	return true
}

func (self *TPlayer) userEnterBattle() bool {
	reqUrl := "http://127.0.0.1:8082/enter_campbattle"
	var req msg.MSG_EnterCampBattle_Req
	req.PlayerID = self.PlayerID
	req.SessionKey = self.SessoinKey
	req.BattleCamp = self.BatCamp
	b, _ := json.Marshal(req)
	buffer, err := PostServerReq(reqUrl, bytes.NewReader(b))
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	var ack msg.MSG_EnterCampBattle_Ack
	json.Unmarshal(buffer, &ack)
	if ack.RetCode != 0 {
		fmt.Println("userSetBatCamp failed: ", ack)
		return false
	}

	self.BattleSvrAddr = ack.BattleSvrAddr
	self.EnterCode = ack.EnterCode
	self.BatClient.ConType = tcpclient.CON_TYPE_BATSVR
	self.BatClient.ExtraData = self
	self.BatClient.ConnectToSvr(self.BattleSvrAddr, 1)
	return true
}

func (self *TPlayer) userStartCarry() bool {
	if self.PlayerID < 10000 {
		gamelog.Error("userStartCarry Error : Invalid playerid:%d", self.PlayerID)
	}

	var req msg.MSG_StartCarry_Req
	req.PlayerID = self.PlayerID

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_START_CARRY_REQ)
	req.Write(&writer)
	writer.EndWrite()

	if self == nil {
		fmt.Println("userStartCarry failed self == nil:")
		return false
	}

	if self.BatClient.TcpConn == nil {
		fmt.Println("userStartCarry failed TcpConn == nil:")
		return false
	}

	self.BatClient.TcpConn.WriteMsgData(writer.GetDataPtr())
	return true
}

func (self *TPlayer) userEnterRoom() bool {
	if self.PlayerID < 10000 {
		gamelog.Error("userEnterRoom Error : Invalid playerid:%d", self.PlayerID)
	}

	var req msg.MSG_EnterRoom_Req
	req.PlayerID = self.PlayerID
	req.EnterCode = int(self.EnterCode)
	req.MsgNo = 1

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_ENTER_ROOM_REQ)
	req.Write(&writer)
	writer.EndWrite()

	if self == nil {
		fmt.Println("userEnterRoom failed self == nil:")
		return false
	}

	if self.BatClient.TcpConn == nil {
		fmt.Println("userEnterRoom failed TcpConn == nil:")
		return false
	}

	self.BatClient.TcpConn.WriteMsgData(writer.GetDataPtr())
	return true
}

func (self *TPlayer) userMove() bool {

	self.Heros[0].Position[0] = self.Heros[0].Position[0] + float32(utility.Rand()%10-5)
	self.Heros[0].Position[2] = self.Heros[0].Position[2] + float32(utility.Rand()%10-5)

	var req msg.MSG_Move_Req
	req.MsgNo = 1
	req.MoveEvents_Cnt = 1
	req.MoveEvents = append(req.MoveEvents, msg.MSG_Move_Item{self.Heros[0].ObjectID, self.Heros[0].Position})

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_MOVE_STATE)
	req.Write(&writer)
	writer.EndWrite()

	if self == nil {
		fmt.Println("userMove failed self == nil:")
		return false
	}

	if self.BatClient.TcpConn == nil {
		fmt.Println("userMove failed TcpConn == nil:")
		return false
	}

	if self.BatClient.TcpConn.IsConnected() == true {
		fmt.Println("userMove not connected:")
		return false
	}

	self.BatClient.TcpConn.WriteMsgData(writer.GetDataPtr())
	fmt.Println("userMove ----:")
	return true
}

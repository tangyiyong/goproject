package mainlogic

import (
	"battlesvr/gamedata"
	"gamelog"
	"msg"
	"time"
)

func (self *TBattleRoom) SendHeroStateToRoom(roomid int16, objectid int32) {
	if roomid <= 0 {
		gamelog.Error("SendHeroStateToRoom Error: Invalid RoomID %d", roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Error("SendHeroStateToRoom Error: Invalid RoomID2 %d", roomid)
		return
	}

	var ackHeroState msg.MSG_HeroState_Nty
	for i := 0; i < len(pRoom.Players); i++ {
		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID > 0 {
			for j := 0; j < 6; j++ {
				ackHeroState.Heros = append(ackHeroState.Heros, msg.MSG_HeroItem{ObjectID: pRoom.Players[i].HeroObj[j].ObjectID, CurHp: pRoom.Players[i].HeroObj[j].CurHp})
			}
		}
	}
	ackHeroState.Heros_Cnt = int32(len(ackHeroState.Heros))
	SendMessageToRoom(0, roomid, msg.MSG_HERO_STATE, &ackHeroState)
	return
}

func (self *TBattleRoom) Hand_LeaveRoom(pdata []byte) {
	gamelog.Info("message: MSG_LEAVE_ROOM_REQ")
	var req msg.MSG_LeaveRoom_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_LeaveRoom : Message Reader Error!!!!")
		return
	}

	pTcpConn := GetConnByID(req.PlayerID)
	if pTcpConn == nil || pTcpConn.Data == nil {
		gamelog.Error("Hand_LeaveRoom : Invalid Playerid:%d", req.PlayerID)
		return
	}

	var response msg.MSG_LeaveRoom_Notify
	response.ObjectIDs = self.GetPlayerHeros(req.PlayerID)
	SendMessageToRoom(req.PlayerID, self.RoomID, msg.MSG_LEAVE_ROOM_NTY, &response)
	pTcpConn.Cleaned = true
	pTcpConn.Close()
	DelConnByID(req.PlayerID)
	G_RoomMgr.RemovePlayerFromRoom(self.RoomID, req.PlayerID)
	return
}

func (self *TBattleRoom) Hand_LeaveByDisconnect(pdata []byte) {
	gamelog.Info("message: MSG_LEAVE_BY_DISCONNT")
	var playerid = int32(pdata[3])<<24 | int32(pdata[3])<<16 | int32(pdata[1])<<8 | int32(pdata[0])
	var response msg.MSG_LeaveRoom_Notify
	response.ObjectIDs = self.GetPlayerHeros(playerid)
	SendMessageToRoom(playerid, self.RoomID, msg.MSG_LEAVE_ROOM_NTY, &response)
	G_RoomMgr.RemovePlayerFromRoom(self.RoomID, playerid)
	return
}

func (self *TBattleRoom) Hand_MoveState(pdata []byte) {
	gamelog.Info("message: MSG_MOVE_STATE")
	var req msg.MSG_Move_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_MoveState : Message Reader Error!!!!")
		return
	}

	if req.MoveEvents_Cnt <= 0 {
		gamelog.Error("Hand_MoveState Error: Invalid MoveEvents_Cnt:%d!!!!", req.MoveEvents_Cnt)
		return
	}

	SendMessageToRoom(req.PlayerID, self.RoomID, msg.MSG_MOVE_STATE, &req)
	pRoom := G_RoomMgr.GetRoomByID(self.RoomID)
	if pRoom == nil {
		gamelog.Error("Hand_MoveState : Invalid RoomID:%d!!!!", self.RoomID)
		return
	}

	for i := 0; i < len(req.MoveEvents); i++ {
		pHeroObject := pRoom.GetHeroObject(req.MoveEvents[i].S_ID)
		if pHeroObject == nil {
			gamelog.Error("Hand_MoveState Error: Invalid Hero ObjectID:%d", req.MoveEvents[i].S_ID)
			continue
		}

		pHeroObject.Position = req.MoveEvents[i].Position
	}

	return
}

func (self *TBattleRoom) Hand_BuffState(pdata []byte) {
	gamelog.Info("message: MSG_BUFF_STATE")
}

func (self *TBattleRoom) Hand_PlayerQueryReq(pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_QUERY_REQ")
	G_GameSvrConns.WriteMsg(msg.MSG_PLAYER_QUERY_REQ, self.RoomID, pdata)
	return
}

func (self *TBattleRoom) Hand_PlayerQueryAck(pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_QUERY_ACK")

	var req msg.MSG_PlayerQuery_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerQueryAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_PlayerQueryAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	pConn.WriteMsg(msg.MSG_PLAYER_QUERY_ACK, 0, pdata)

	return
}

func (self *TBattleRoom) Hand_StartCarryReq(pdata []byte) {
	gamelog.Info("message: MSG_START_CARRY_REQ")
	var req msg.MSG_StartCarry_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_StartCarryReq : Message Reader Error!!!!")
		return
	}

	pBattleObj := self.GetBattleByPID(req.PlayerID)
	if pBattleObj == nil {
		gamelog.Error("Hand_StartCarryReq Error: Invalid playerid:%d", req.PlayerID)
		return
	}

	if pBattleObj.MoveEndTime > 0 {
		gamelog.Error("Hand_StartCarryReq Error: Has Already Carry a Ctystal:%d", req.PlayerID)
		return
	}

	pRect := &gamedata.GetSceneInfo().Camps[pBattleObj.BatCamp-1].MoveBegin
	if !pBattleObj.IsTeamIn(pRect) {
		gamelog.Error("Hand_StartCarryReq Error: Not In The Start Carry Rect:%d", req.PlayerID)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_START_CARRY_REQ, self.RoomID, pdata)

	return
}

func (self *TBattleRoom) Hand_FinishCarryReq(pdata []byte) {
	gamelog.Info("message: MSG_FINISH_CARRY_REQ")
	var req msg.MSG_FinishCarry_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_FinishCarryReq Error: Message Reader Error!!!!")
		return
	}

	pBattleObj := self.GetBattleByPID(req.PlayerID)
	if pBattleObj == nil {
		gamelog.Error("Hand_FinishCarryReq Error: Invalid playerid:%d", req.PlayerID)
		return
	}

	if pBattleObj.MoveEndTime <= 0 {
		gamelog.Error("Hand_FinishCarryReq Error: Has not start:%d", req.PlayerID)
		return
	}

	if int32(time.Now().Unix()) > pBattleObj.MoveEndTime {
		gamelog.Error("Hand_FinishCarryReq Error: Too late:%d", req.PlayerID)
		return
	}

	pRect := &gamedata.GetSceneInfo().Camps[pBattleObj.BatCamp-1].MoveEnd
	if !pBattleObj.IsTeamIn(pRect) {
		gamelog.Error("Hand_FinishCarryReq Error: Not In The Start Carry Rect:%d", req.PlayerID)
		return
	}

	G_GameSvrConns.WriteMsg(msg.MSG_FINISH_CARRY_REQ, self.RoomID, pdata)

	return
}

func (self *TBattleRoom) Hand_StartCarryAck(pdata []byte) {
	gamelog.Info("message: MSG_START_CARRY_ACK")
	var req msg.MSG_StartCarry_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_StartCarryAck Error: Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_StartCarryAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	playerid := pConn.Data.(*TBattleData).PlayerID
	roomid := pConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_StartCarryAck Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_StartCarryAck Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pBattleObj := pRoom.GetBattleByPID(playerid)
	if pBattleObj == nil {
		gamelog.Error("Hand_StartCarryAck Error: Invalid playerid:%d", playerid)
		return
	}

	if req.RetCode == msg.RE_SUCCESS {
		pBattleObj.MoveEndTime = req.EndTime
	}

	pConn.WriteMsg(msg.MSG_START_CARRY_ACK, 0, pdata)

	return
}

func (self *TBattleRoom) Hand_FinishCarryAck(pdata []byte) {
	gamelog.Info("message: MSG_FINISH_CARRY_ACK")
	var req msg.MSG_FinishCarry_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_FinishCarryAck Error: Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_FinishCarryAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	if req.RetCode == msg.RE_SUCCESS {
		playerid := pConn.Data.(*TBattleData).PlayerID
		roomid := pConn.Data.(*TBattleData).RoomID
		if playerid <= 0 || roomid <= 0 {
			gamelog.Info("Hand_FinishCarryAck Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
			return
		}

		pRoom := G_RoomMgr.GetRoomByID(roomid)
		if pRoom == nil {
			gamelog.Info("Hand_FinishCarryAck Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
			return
		}

		pBattleObj := pRoom.GetBattleByPID(playerid)
		if pBattleObj == nil {
			gamelog.Error("Hand_FinishCarryAck Error: Invalid playerid:%d", playerid)
			return
		}

		pBattleObj.MoveEndTime = 0
	}

	pConn.WriteMsg(msg.MSG_FINISH_CARRY_ACK, 0, pdata)

	return
}

func (self *TBattleRoom) Hand_PlayerChangeReq(pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_CHANGE_REQ")
	G_GameSvrConns.WriteMsg(msg.MSG_PLAYER_CHANGE_REQ, self.RoomID, pdata)
	return
}

func (self *TBattleRoom) Hand_PlayerChangeAck(pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_CHANGE_ACK")
	var req msg.MSG_PlayerChange_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerChangeAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_PlayerChangeAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	pConn.WriteMsg(msg.MSG_PLAYER_CHANGE_ACK, 0, pdata)

	return
}

func (self *TBattleRoom) Hand_PlayerReviveReq(pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_REVIVE_REQ")
	G_GameSvrConns.WriteMsg(msg.MSG_PLAYER_REVIVE_REQ, self.RoomID, pdata)
	return
}

func (self *TBattleRoom) Hand_PlayerReviveAck(pdata []byte) {
	gamelog.Info("message: MSG_PLAYER_REVIVE_ACK")
	var req msg.MSG_ServerRevive_Ack
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerReviveAck : Message Reader Error!!!!")
		return
	}

	pConn := GetConnByID(req.PlayerID)
	if pConn == nil {
		gamelog.Error("Hand_PlayerReviveAck : Invalid PlayerID:%d!!!!", req.PlayerID)
		return
	}

	playerid := pConn.Data.(*TBattleData).PlayerID
	roomid := pConn.Data.(*TBattleData).RoomID
	if playerid <= 0 || roomid <= 0 {
		gamelog.Info("Hand_PlayerReviveReq Error: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(roomid)
	if pRoom == nil {
		gamelog.Info("Hand_PlayerChangeReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	pBattleObj := pRoom.GetBattleByPID(playerid)
	if pBattleObj == nil {
		gamelog.Info("Hand_PlayerChangeReq Error2: Invalid PlayerID :%d and roomid :%d", playerid, roomid)
		return
	}

	if req.ReviveOpt == 2 {
		pBattleObj.ReviveTime[0] += 1
	} else if req.ReviveOpt == 3 {
		pBattleObj.ReviveTime[1] += 1
	} else if req.ReviveOpt == 4 {
	}

	if req.Stay != 1 { //表示要安全区复活

	}

	var response msg.MSG_PlayerRevive_Ack
	response.MoneyID = req.MoneyID
	response.MoneyNum = req.MoneyNum
	response.BatCamp = pBattleObj.BatCamp
	response.RetCode = req.RetCode

	for i := 0; i < 6; i++ {
		if pBattleObj.HeroObj[i].HeroID <= 0 {
			break
		}
		response.Heros = append(response.Heros, msg.MSG_HeroObj{pBattleObj.HeroObj[i].HeroID,
			pBattleObj.HeroObj[i].ObjectID,
			pBattleObj.HeroObj[i].CurHp,
			pBattleObj.HeroObj[i].Position})
	}
	response.Heros_Cnt = int32(len(response.Heros))

	var writer msg.PacketWriter
	writer.BeginWrite(msg.MSG_PLAYER_REVIVE_ACK, 0)
	response.Write(&writer)
	writer.EndWrite()
	pConn.WriteMsgData(writer.GetDataPtr())

	return
}

func (self *TBattleRoom) Hand_PlayerChatReq(pdata []byte) {
	gamelog.Info("message: MSG_CAMPBAT_CHAT_REQ")
	var req msg.MSG_CmapBatChat_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_PlayerChatReq : Message Reader Error!!!!")
		return
	}
	SendMessageToRoom(req.PlayerID, self.RoomID, msg.MSG_CAMPBAT_CHAT_REQ, &req)
	return
}

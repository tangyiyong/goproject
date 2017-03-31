package mainlogic

import (
	"gamelog"
	"msg"
)

const (
	room_player_num = 15 //一个房间的最大人数
	room_camp_num   = 3  //一个房间的最大阵营数
	camp_player_num = 5  //一个阵营的最大人数
)

type TMessage struct {
	MsgID   int16
	MsgData []byte
}

type TBattleRoom struct {
	RoomID  int16                        //房间ID
	Players [room_player_num]*TBattleObj //三个阵营的人员
	CampNum [room_camp_num]int32         //各个阵营人数
	LastTick int32                       //上一次tick时间
	MsgList chan TMessage                //消息队列
}

func (self *TBattleRoom) Init(id int16) bool {
	self.RoomID = id
	self.MsgList = make(chan TMessage, 100)
	self.LastTick = utility.GetCurTime();
	go self.MsgProcess()
	return true
}

func (self *TBattleRoom) MsgProcess() {
	for msgItem := range self.MsgList {
		switch msgItem.MsgID {
		case msg.MSG_MOVE_STATE:
			self.Hand_MoveState(msgItem.MsgData)
		case msg.MSG_SKILL_STATE:
			self.Hand_SkillState(msgItem.MsgData)
		case msg.MSG_BUFF_STATE:
			self.Hand_BuffState(msgItem.MsgData)
		case msg.MSG_LEAVE_ROOM_REQ:
			self.Hand_LeaveRoom(msgItem.MsgData)
		case msg.MSG_PLAYER_QUERY_REQ:
			self.Hand_PlayerQueryReq(msgItem.MsgData)
		case msg.MSG_PLAYER_QUERY_ACK:
			self.Hand_PlayerQueryAck(msgItem.MsgData)
		case msg.MSG_PLAYER_CHANGE_REQ:
			self.Hand_PlayerChangeReq(msgItem.MsgData)
		case msg.MSG_PLAYER_CHANGE_ACK:
			self.Hand_PlayerChangeAck(msgItem.MsgData)
		case msg.MSG_PLAYER_REVIVE_REQ:
			self.Hand_PlayerReviveReq(msgItem.MsgData)
		case msg.MSG_PLAYER_REVIVE_ACK:
			self.Hand_PlayerReviveAck(msgItem.MsgData)
		case msg.MSG_CAMPBAT_CHAT_REQ:
			self.Hand_PlayerChatReq(msgItem.MsgData)
		case msg.MSG_START_CARRY_REQ:
			self.Hand_StartCarryReq(msgItem.MsgData)
		case msg.MSG_FINISH_CARRY_REQ:
			self.Hand_FinishCarryReq(msgItem.MsgData)
		case msg.MSG_START_CARRY_ACK:
			self.Hand_StartCarryAck(msgItem.MsgData)
		case msg.MSG_FINISH_CARRY_ACK:
			self.Hand_FinishCarryAck(msgItem.MsgData)
		case msg.MSG_LEAVE_BY_DISCONNT:
			self.Hand_LeaveByDisconnect(msgItem.MsgData)
		}
	}
	
	if((utility.GetCurTime()-self.LastTick) <25)
	{
		return ;
	}
	
	self.LastTick = utility.GetCurTime();
	
	self.Update();
}

func (self *TBattleRoom) Update() bool {
	
	
	return true;
}

//由于客户端原因，index_player需要从2开始计
func (self *TBattleRoom) AddPlayer(pBattleObj *TBattleObj) bool {
	if pBattleObj == nil || pBattleObj.PlayerID <= 0 || pBattleObj.BatCamp <= 0 || pBattleObj.BatCamp > room_camp_num {
		gamelog.Error("AddPlayer Error Invalid Parameter playerid:%d, batcamp:%d!!!", pBattleObj.PlayerID, pBattleObj.BatCamp)
		return false
	}

	var i int32 = 0
	for ; i < room_player_num; i++ {
		if self.Players[i] == nil {
			self.Players[i] = pBattleObj
			for j := int32(0); j < 6; j++ {
				pBattleObj.HeroObj[j].ObjectID = (i+2)<<16 | j
			}
			break
		}
	}

	if i == room_player_num {
		gamelog.Error("AddPlayer Error No space for new player!")
		return false
	}

	self.CampNum[pBattleObj.BatCamp-1] += 1
	return true
}

func (self *TBattleRoom) RemovePlayer(playerid int32) bool {
	if playerid <= 0 {
		gamelog.Error("RemovePlayer Error Invalid Parameter!!!")
		return false
	}

	var i = 0
	for ; i < room_player_num; i++ {
		if self.Players[i] == nil {
			continue
		}

		if self.Players[i].PlayerID == playerid {
			self.CampNum[self.Players[i].BatCamp-1] -= 1
			self.Players[i] = nil
			break
		}
	}

	return true
}

func (self *TBattleRoom) GetHeroObject(objectid int32) *THeroObj {
	idx_player := (objectid >> 16) - 2
	idx_hero := objectid & 0x00ff

	if idx_player >= room_player_num || idx_player < 0 {
		gamelog.Error("GetHeroObject Error Objectid:%d, Invalid idx_player:%d", objectid, idx_player)
		return nil
	}

	if self.Players[idx_player] == nil {
		gamelog.Error("GetHeroObject Error Objectid:%d, Invalid idx_player:%d", objectid, idx_player)
		return nil
	}

	if idx_hero >= 6 || idx_hero < 0 {
		gamelog.Error("GetHeroObject Error Objectid:%d, Invalid idx_hero:%d", objectid, idx_hero)
		return nil
	}

	return &self.Players[idx_player].HeroObj[idx_hero]
}

func (self *TBattleRoom) GetBattleByPID(playerid int32) *TBattleObj {
	for i := 0; i < room_player_num; i++ {
		if self.Players[i] != nil && self.Players[i].PlayerID == playerid {
			return self.Players[i]
		}
	}

	return nil
}
func (self *TBattleRoom) GetBattleByOID(objectid int32) *TBattleObj {
	idx_player := (objectid >> 16) - 2
	return self.Players[idx_player]
}

func (self *TBattleRoom) GetPlayerHeros(playerid int32) (ret [6]int32) {
	if playerid <= 0 {
		gamelog.Error("GetPlayerHeros Error : playerid:%d", playerid)
		return
	}

	for i := 0; i < room_player_num; i++ {
		if self.Players[i] != nil && self.Players[i].PlayerID == playerid {
			for j := 0; j < 6; j++ {
				ret[j] = self.Players[i].HeroObj[j].ObjectID
			}
		}
	}

	return
}

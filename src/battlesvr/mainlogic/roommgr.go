package mainlogic

import (
	"gamelog"
	"sync"
	//"tcpserver"
)

const (
	Room_Type_Low  = 1 //低等级房间
	Room_Type_High = 2 //高等级房间

	LowRoom_StartID  = int16(1)     //低等级房间起始ID
	HighRooM_StartID = int16(10000) //高等级房间起始ID
)

type TRoomMgr struct {
	sync.Mutex
	//两个等级房间都创建1000间， 每个等级都可以容纳15000人，基本够用。
	LowRooms  [1000]TBattleRoom //低等级房间
	HighRooms [1000]TBattleRoom //高等级房间
}

var (
	G_RoomMgr TRoomMgr
)

func InitRoomMgr() bool {
	for i := int16(0); i < 1000; i++ {
		G_RoomMgr.LowRooms[i].Init(i + 1)
		G_RoomMgr.HighRooms[i].Init(i + HighRooM_StartID)
	}
	return true
}

func (self *TRoomMgr) GetRoomByID(roomid int16) *TBattleRoom {
	if roomid <= 0 {
		gamelog.Error("GetRoomByID Error : Invalid roomid :%d", roomid)
		return nil
	}

	if roomid < HighRooM_StartID {
		return &self.LowRooms[roomid-1]
	} else {
		return &self.HighRooms[roomid-HighRooM_StartID]
	}

	return nil
}

func (self *TRoomMgr) AddPlayerToRoom(roomtype int32, batcamp int8, pBattleObj *TBattleObj) int16 {
	self.Lock()
	defer self.Unlock()

	if roomtype == Room_Type_Low {
		for i := 0; i < len(self.LowRooms); i++ {
			if self.LowRooms[i].CampNum[batcamp-1] >= camp_player_num {
				continue
			} else {
				self.LowRooms[i].AddPlayer(pBattleObj)
				return self.LowRooms[i].RoomID
			}
		}
	} else if roomtype == Room_Type_High {
		for i := 0; i < len(self.HighRooms); i++ {
			if self.HighRooms[i].CampNum[batcamp-1] >= camp_player_num {
				continue
			} else {
				self.HighRooms[i].AddPlayer(pBattleObj)
				return self.HighRooms[i].RoomID
			}
		}
	} else {
		gamelog.Error("AddPlayerToRoom Error : Invalid RoomType :%d", roomtype)
	}

	return 0
}

func (self *TRoomMgr) RemovePlayerFromRoom(roomid int16, playerid int32) bool {
	if roomid <= 0 || playerid <= 0 {
		gamelog.Error("GetPlayerHeroIDs Error : Invalid roomid :%d and playerid:%d", roomid, playerid)
		return false
	}

	self.Lock()
	defer self.Unlock()

	var pRoom *TBattleRoom = nil
	if roomid < HighRooM_StartID {
		pRoom = &self.LowRooms[roomid-1]
	} else {
		pRoom = &self.HighRooms[roomid-HighRooM_StartID]
	}

	pRoom.RemovePlayer(playerid)

	return true
}

package mainlogic

import (
	"gamelog"
	"sync"
	//"tcpserver"
)

const (
	Room_Type_Low  = 1 //低等级房间
	Room_Type_High = 2 //高等级房间

	LowRoom_StartID  = 1     //低等级房间起始ID
	HighRooM_StartID = 10000 //高等级房间起始ID
)

type TRoomMgr struct {
	sync.Mutex
	LowRooms  []TBattleRoom //低等级房间
	HighRooms []TBattleRoom //高等级房间
}

var (
	G_RoomMgr TRoomMgr
)

func InitRoomMgr() bool {
	//两个等级房间都创建1000间， 每个等级都可以容纳15000人，基本够用。
	G_RoomMgr.LowRooms = make([]TBattleRoom, 1000, 1000)
	G_RoomMgr.HighRooms = make([]TBattleRoom, 1000, 1000)
	for i := 0; i < 1000; i++ {
		G_RoomMgr.LowRooms[i].Init(i+1, Room_Type_Low)
		G_RoomMgr.HighRooms[i].Init(i+HighRooM_StartID, Room_Type_High)
	}
	return true
}

func (mgr *TRoomMgr) GetRoomByID(roomid int) *TBattleRoom {
	if roomid <= 0 {
		gamelog.Error("GetRoomByID Error : Invalid roomid :%d", roomid)
		return nil
	}

	if roomid < HighRooM_StartID {
		return &mgr.LowRooms[roomid-1]
	} else {
		return &mgr.HighRooms[roomid-HighRooM_StartID]
	}

	return nil
}

func (mgr *TRoomMgr) GetPlayerHeroIDs(roomid int, playerid int) (ret [6]int) {
	mgr.Lock()
	defer mgr.Unlock()

	var pRoom *TBattleRoom = nil
	if roomid < HighRooM_StartID {
		pRoom = &mgr.LowRooms[roomid-1]
	} else {
		pRoom = &mgr.HighRooms[roomid-HighRooM_StartID]
	}

	for i := 0; i < max_room_player; i++ {
		if pRoom.Players[i] != nil && pRoom.Players[i].PlayerID == playerid {
			for j := 0; j < 6; j++ {
				ret[j] = pRoom.Players[i].HeroObj[j].ObjectID
			}
		}
	}

	return
}

func (mgr *TRoomMgr) AddPlayerToRoom(roomtype int, batcamp int, pBattleObj *TBattleObj) int {
	mgr.Lock()
	defer mgr.Unlock()

	if roomtype == Room_Type_Low {
		for i := 0; i < len(mgr.LowRooms); i++ {
			if mgr.LowRooms[i].CampNum[batcamp-1] >= max_one_camp_player {
				continue
			} else {
				mgr.LowRooms[i].AddPlayer(pBattleObj)
				return mgr.LowRooms[i].RoomID
			}
		}
	} else if roomtype == Room_Type_High {
		for i := 0; i < len(mgr.HighRooms); i++ {
			if mgr.HighRooms[i].CampNum[batcamp-1] >= max_one_camp_player {
				continue
			} else {
				mgr.HighRooms[i].AddPlayer(pBattleObj)
				return mgr.HighRooms[i].RoomID
			}
		}
	} else {
		gamelog.Error("AddPlayerToRoom Error : Invalid RoomType :%d", roomtype)
	}

	return 0
}

func (mgr *TRoomMgr) RemovePlayerFromRoom(roomid int, playerid int) bool {
	mgr.Lock()
	defer mgr.Unlock()

	var pRoom *TBattleRoom = nil
	if roomid < HighRooM_StartID {
		pRoom = &mgr.LowRooms[roomid-1]
	} else {
		pRoom = &mgr.HighRooms[roomid-HighRooM_StartID]
	}

	pRoom.RemovePlayer(playerid)

	return true
}

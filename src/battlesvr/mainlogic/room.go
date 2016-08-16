package mainlogic

import (
	"gamelog"
)

const (
	max_room_player     = 15 //一个房间的最大人数
	max_room_camp       = 3  //一个房间的最大阵营数
	max_one_camp_player = 5  //一个阵营的最大人数
)

type TBattleRoom struct {
	RoomID   int                          //房间ID
	RoomType int                          //等级类型
	Players  [max_room_player]*TBattleObj //三个阵营的人员
	CampNum  [max_room_camp]int           //各个阵营人数
}

func (room *TBattleRoom) Init(id int, roomtype int) bool {
	room.RoomID = id
	room.RoomType = roomtype
	return true
}

//由于客户端原因，index_player需要从2开始计

func (room *TBattleRoom) AddPlayer(pBattleObj *TBattleObj) bool {
	if pBattleObj == nil || pBattleObj.PlayerID <= 0 || pBattleObj.BatCamp <= 0 || pBattleObj.BatCamp > max_room_camp {
		gamelog.Error("AddPlayer Error Invalid Parameter playerid:%d, batcamp:%d!!!", pBattleObj.PlayerID, pBattleObj.BatCamp)
		return false
	}

	var i = 0
	for ; i < max_room_player; i++ {
		if room.Players[i] == nil {
			room.Players[i] = pBattleObj
			for j := 0; j < 6; j++ {
				pBattleObj.HeroObj[j].ObjectID = (i+2)<<16 | j
			}
			break
		}
	}

	if i == max_room_player {
		gamelog.Error("AddPlayer Error No space for new player!")
		return false
	}

	room.CampNum[pBattleObj.BatCamp-1] += 1
	return true
}

func (room *TBattleRoom) RemovePlayer(playerid int) bool {
	if playerid <= 0 {
		gamelog.Error("RemovePlayer Error Invalid Parameter!!!")
		return false
	}

	var i = 0
	for ; i < max_room_player; i++ {
		if room.Players[i] == nil {
			continue
		}

		if room.Players[i].PlayerID == playerid {
			room.CampNum[room.Players[i].BatCamp-1] -= 1
			room.Players[i] = nil
			break
		}
	}

	return true
}

func (room *TBattleRoom) GetHeroObject(objectid int) *THeroObj {
	idx_player := (objectid>>16) - 2
	idx_hero := objectid & 0x00ff

	if idx_player >= max_room_player || idx_player < 0 {
		gamelog.Error("GetHeroObject Error Objectid:%d, Invalid idx_player:%d", objectid, idx_player)
		return nil
	}

	if room.Players[idx_player] == nil {
		gamelog.Error("GetHeroObject Error Objectid:%d, Invalid idx_player:%d", objectid, idx_player)
		return nil
	}

	if idx_hero >= 6 || idx_hero < 0 {
		gamelog.Error("GetHeroObject Error Objectid:%d, Invalid idx_hero:%d", objectid, idx_hero)
		return nil
	}

	return &room.Players[idx_player].HeroObj[idx_hero]
}

func (room *TBattleRoom) GetBattleByPID(playerid int) *TBattleObj {
	for i := 0; i < max_room_player; i++ {
		if room.Players[i] != nil && room.Players[i].PlayerID == playerid {
			return room.Players[i]
		}
	}

	return nil
}
func (room *TBattleRoom) GetBattleByOID(objectid int) *TBattleObj {
	idx_player := (objectid>>16) - 2
	return room.Players[idx_player]
}

package mainlogic

import (
	"appconfig"
	//"fmt"
	"gamelog"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

var (
	G_Player_Mutex sync.Mutex
	G_Players      map[int32]*TPlayer //玩家集

	G_SelectPlayers  []*TPlayer //用来选择用的玩家表
	G_CurSelectIndex int        //当前选择索引

)

func GetPlayerByID(playerid int32) *TPlayer {
	G_Player_Mutex.Lock()
	defer G_Player_Mutex.Unlock()
	info, ok := G_Players[playerid]
	if ok {
		return info
	}

	return nil
}

func CreatePlayer(playerid int32, name string, heroid int) (*TPlayer, bool) {
	G_Player_Mutex.Lock()
	_, ok := G_Players[playerid]
	if ok {
		G_Player_Mutex.Unlock()
		gamelog.Error("Create Player Failed Error : playerid : %d exist!!!", playerid)
		return nil, false
	}

	player := new(TPlayer)
	G_Players[playerid] = player
	G_SelectPlayers = append(G_SelectPlayers, player)
	player.InitModules(playerid)
	player.SetPlayerName(name)
	player.SetMainHeroID(heroid)
	G_Player_Mutex.Unlock()

	return player, true
}

func LoadPlayerFromDB(playerid int32) *TPlayer {
	if playerid <= 0 {
		gamelog.Error("LoadPlayerFromDB Error : Invalid playerid :%d", playerid)
		return nil
	}

	G_Player_Mutex.Lock()
	player := new(TPlayer)
	G_Players[playerid] = player
	G_SelectPlayers = append(G_SelectPlayers, player)
	G_Player_Mutex.Unlock()
	player.OnPlayerLoad(playerid)
	player.pSimpleInfo = G_SimpleMgr.GetSimpleInfoByID(playerid)

	return player
}

func DestroyPlayer(playerid int32) bool {
	G_Player_Mutex.Lock()
	defer G_Player_Mutex.Unlock()

	player, ok := G_Players[playerid]
	if ok {
		delete(G_Players, playerid)
		player.OnDestroy(playerid)
	}

	return true
}

type TReslutID struct {
	ID int32 `bson:"_id"` //主键 玩家ID
}

//将一些有价值的玩家预先加载到服务器中
func PreLoadPlayers() {
	s := mongodb.GetDBSession()
	defer s.Close()

	query := s.DB(appconfig.GameDbName).C("PlayerSimple").Find(nil).Select(&bson.M{"_id": 1}).Sort("-logofftime").Limit(3000)
	iter := query.Iter()

	//fmt.Printf("PreLoadPlayers:%10d", 1)
	var result TReslutID
	for iter.Next(&result) {
		if result.ID < 10000 {
			gamelog.Error("PreLoadPlayers Error: Invalid PlayerID:%d", result.ID)
			continue
		}

		//fmt.Printf("\b\b\b\b\b\b\b\b")
		//fmt.Printf("%8d", result.ID)

		player := new(TPlayer)
		G_Players[result.ID] = player
		G_SelectPlayers = append(G_SelectPlayers, player)
		player.OnPlayerLoadSync(result.ID)
		player.pSimpleInfo = G_SimpleMgr.GetSimpleInfoByID(result.ID)
	}
	//fmt.Printf("\b\b\b\b\b\b\b\b")
	//fmt.Printf("Successed!!\n")
}

func GetSelectPlayer(selectfunc func(p *TPlayer, value int) bool, selectvalue int) *TPlayer {
	nTotal := len(G_SelectPlayers)
	if nTotal <= 0 {
		return nil
	}
	if nTotal <= G_CurSelectIndex {
		for i := 0; i < nTotal; i++ {
			if true == selectfunc(G_SelectPlayers[i], selectvalue) {
				G_CurSelectIndex = i + 1
				return G_SelectPlayers[i]
			}
		}
		G_CurSelectIndex = 0
	} else {
		for i := G_CurSelectIndex; i < nTotal; i++ {
			if true == selectfunc(G_SelectPlayers[i], selectvalue) {
				G_CurSelectIndex = i + 1
				return G_SelectPlayers[i]
			}
		}

		for i := 0; i < G_CurSelectIndex; i++ {
			if true == selectfunc(G_SelectPlayers[i], selectvalue) {
				G_CurSelectIndex = i + 1
				return G_SelectPlayers[i]
			}
		}
	}

	return nil
}

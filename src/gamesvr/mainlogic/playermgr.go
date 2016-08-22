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
	mMutex        sync.Mutex
	g_Players     map[int]*TPlayer //玩家集
	g_OnlineCount int              //在线玩家数量

	g_CurSelectIndex int        //当前选择索引
	g_SelectPlayers  []*TPlayer //用来选择用的玩家表

)

func GetPlayerByID(playerid int) *TPlayer {
	mMutex.Lock()
	defer mMutex.Unlock()
	info, ok := g_Players[playerid]
	if ok {
		return info
	}

	return nil
}

func CreatePlayer(playerid int, name string, heroid int) (*TPlayer, bool) {
	mMutex.Lock()
	_, ok := g_Players[playerid]
	if ok {
		mMutex.Unlock()
		gamelog.Error("Create Player Failed Error : playerid : %d exist!!!")
		return nil, false
	}

	pPlayer := new(TPlayer)
	g_Players[playerid] = pPlayer
	g_SelectPlayers = append(g_SelectPlayers, pPlayer)
	pPlayer.InitModules(playerid)
	pPlayer.SetPlayerName(name)
	pPlayer.SetMainHeroID(heroid)
	mMutex.Unlock()

	return pPlayer, true
}

func LoadPlayerFromDB(playerid int) *TPlayer {
	if playerid <= 0 {
		gamelog.Error("LoadPlayerFromDB Error : Invalid playerid :%d", playerid)
		return nil
	}

	mMutex.Lock()
	pPlayer := new(TPlayer)
	g_Players[playerid] = pPlayer
	g_SelectPlayers = append(g_SelectPlayers, pPlayer)
	mMutex.Unlock()
	pPlayer.OnPlayerLoad(playerid)
	pPlayer.pSimpleInfo = G_SimpleMgr.GetSimpleInfoByID(playerid)

	return pPlayer
}

func GetOnlineCount() int {
	mMutex.Lock()
	defer mMutex.Unlock()

	return g_OnlineCount
}

func DestroyPlayer(playerid int) bool {
	mMutex.Lock()
	defer mMutex.Unlock()

	pPlayer, ok := g_Players[playerid]
	if ok {
		delete(g_Players, playerid)
		pPlayer.OnDestroy(playerid)
	}

	return true
}

type TReslutID struct {
	ID int "_id"
}

//将一些有价值的玩家预先加载到服务器中
func PreLoadPlayers() {
	s := mongodb.GetDBSession()
	defer s.Close()

	query := s.DB(appconfig.GameDbName).C("PlayerRole").Find(nil).Select(bson.M{"_id": 1}).Sort("-Level").Limit(10000)
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

		pPlayer := new(TPlayer)
		g_Players[result.ID] = pPlayer
		g_SelectPlayers = append(g_SelectPlayers, pPlayer)
		pPlayer.OnPlayerLoadSync(result.ID)
		pPlayer.pSimpleInfo = G_SimpleMgr.GetSimpleInfoByID(result.ID)
	}
	//fmt.Printf("\b\b\b\b\b\b\b\b")
	//fmt.Printf("Successed!!\n")
}

func GetSelectPlayer(selectfunc func(p *TPlayer, value int) bool, selectvalue int) *TPlayer {
	nTotal := len(g_SelectPlayers)
	if nTotal <= 0 {
		return nil
	}
	if nTotal <= g_CurSelectIndex {
		for i := 0; i < nTotal; i++ {
			if true == selectfunc(g_SelectPlayers[i], selectvalue) {
				g_CurSelectIndex = i + 1
				return g_SelectPlayers[i]
			}
		}
		g_CurSelectIndex = 0
	} else {
		for i := g_CurSelectIndex; i < nTotal; i++ {
			if true == selectfunc(g_SelectPlayers[i], selectvalue) {
				g_CurSelectIndex = i + 1
				return g_SelectPlayers[i]
			}
		}

		for i := 0; i < g_CurSelectIndex; i++ {
			if true == selectfunc(g_SelectPlayers[i], selectvalue) {
				g_CurSelectIndex = i + 1
				return g_SelectPlayers[i]
			}
		}
	}

	return nil
}

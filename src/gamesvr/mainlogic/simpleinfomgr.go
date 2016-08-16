package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2"
)

type TSimpleInfo struct {
	PlayerID      int    `bson:"_id"` //玩家ID
	AccountID     int    //账号ID
	GuildID       int    //公会ID
	HeroID        int    //英雄ID
	Quality       int    //主角品质
	Name          string //玩家名字
	Level         int    //玩家等级
	VipLevel      int    //玩家的VIP等级
	FightValue    int    //玩家的战力
	LogoffTime    int64  //离线时间
	AwardCenterID int    //奖励中心ID
	isOnline      bool   //是否在线
	BatCamp       int    //阵营战阵营
}

type TSimpleInfoMgr struct {
	SimpleList map[int]*TSimpleInfo //角色离线信息map
	NameIDMap  map[string]int       //名字到ID的map
	SimpleLock sync.Mutex
}

var G_SimpleMgr TSimpleInfoMgr

func (mgr *TSimpleInfoMgr) Init() bool {
	mgr.SimpleList = make(map[int]*TSimpleInfo, 1000)
	mgr.NameIDMap = make(map[string]int, 1000)
	s := mongodb.GetDBSession()
	defer s.Close()

	var simplevec []TSimpleInfo
	err := s.DB(appconfig.GameDbName).C("PlayerSimple").Find(nil).Sort("+_id").All(&simplevec)
	if err != nil {
		if err != mgo.ErrNotFound {
			gamelog.Error("Init SimpleInfo Mgr DB Error!!!")
			return false
		}
	}

	nCount := len(simplevec)
	if nCount <= 0 {
		return true
	}

	var i int
	for i = 0; i < nCount; i++ {
		mgr.SimpleList[simplevec[i].AccountID] = &simplevec[i]
		mgr.NameIDMap[simplevec[i].Name] = simplevec[i].PlayerID
	}

	return true
}

func GetPlayerIDByAccountID(accountid int) int {
	G_SimpleMgr.SimpleLock.Lock()
	defer G_SimpleMgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[accountid]
	if ok && pInfo != nil {
		return pInfo.PlayerID
	}

	return 0
}

func (mgr *TSimpleInfoMgr) GetPlayerIDByName(name string) int {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	playerid, ok := mgr.NameIDMap[name]
	if ok {
		return playerid
	}

	return 0
}

func (mgr *TSimpleInfoMgr) GetPlayerLogoffTime(playerid int) int64 {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	if ok && pInfo != nil {
		return pInfo.LogoffTime
	}

	return 0
}

func (mgr *TSimpleInfoMgr) GetSimpleInfoByID(playerid int) *TSimpleInfo {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	if ok && pInfo != nil {
		return pInfo
	}

	gamelog.Error3("GetSimpleInfoByID Error , Invalid playerid:%d", playerid)
	return nil
}

func (mgr *TSimpleInfoMgr) GetPlayerAwardCenterID(playerID int) int {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[playerID]
	if ok && pInfo != nil {
		return pInfo.AwardCenterID
	}

	return 0
}

func (mgr *TSimpleInfoMgr) RemoveSimpleInfo(playerid int) bool {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()
	delete(mgr.SimpleList, playerid)
	return true
}

func (mgr *TSimpleInfoMgr) Get_FightValue(playerid int) int {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		return pInfo.FightValue
	}
	return 0
}

func (mgr *TSimpleInfoMgr) Set_FightValue(playerid int, fightvalue int, level int) bool {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		if fightvalue != pInfo.FightValue || level != pInfo.Level {
			pInfo.FightValue = fightvalue
			pInfo.Level = level
			G_SimpleMgr.DB_SetFightValue(playerid, fightvalue, pInfo.Level)

			return true
		}

		return false
	}

	gamelog.Error("Set_FightValue Error Cant find the player:%d", playerid)
	return false
}

func (mgr *TSimpleInfoMgr) Set_PlayerName(playerid int, name string) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.Name = name
		G_SimpleMgr.DB_SetPlayerName(playerid, name)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_LogoffTime(playerid int, time int64) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.LogoffTime = time
		G_SimpleMgr.DB_SetLogoffTime(playerid, time)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_HeroID(playerid int, heroid int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.HeroID = heroid
		G_SimpleMgr.DB_SetHeroID(playerid, heroid)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_HeroQuality(playerid int, quality int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.Quality = quality
		G_SimpleMgr.DB_SetHeroQuality(playerid, quality)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_VipLevel(playerid int, viplevel int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.VipLevel = viplevel
		G_SimpleMgr.DB_SetVipLevel(playerid, viplevel)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_AwardCenterID(playerID int, awardCenterID int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerID]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.AwardCenterID = awardCenterID
		G_SimpleMgr.DB_SetAwardCenterID(playerID, awardCenterID)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_BatCamp(playerid int, camp int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		pInfo.BatCamp = camp
		G_SimpleMgr.DB_SetBatCamp(playerid, camp)
		return
	}
	return
}

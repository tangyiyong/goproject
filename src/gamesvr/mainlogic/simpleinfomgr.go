package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2"
)

type TSimpleInfo struct {
	PlayerID      int32  `bson:"_id"` //玩家ID
	AccountID     int32  //账号ID
	GuildID       int    //公会ID
	HeroID        int    //英雄ID
	Quality       int8   //主角品质
	Name          string //玩家名字
	Level         int    //玩家等级
	VipLevel      int    //玩家的VIP等级
	FightValue    int    //玩家的战力
	LogoffTime    int64  //离线时间
	AwardCenterID int    //奖励中心ID
	BatCamp       int8   //阵营战阵营
	LoginDay      uint32 //登录日期
	isOnline      bool   //是否在线
}

type TSimpleInfoMgr struct {
	SimpleList map[int32]*TSimpleInfo //角色离线信息map
	NameIDMap  map[string]int32       //名字到ID的map
	SimpleLock sync.Mutex
}

var G_SimpleMgr TSimpleInfoMgr

func (mgr *TSimpleInfoMgr) Init() bool {
	mgr.SimpleList = make(map[int32]*TSimpleInfo, 1000)
	mgr.NameIDMap = make(map[string]int32, 1000)
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

	for i := 0; i < nCount; i++ {
		mgr.SimpleList[simplevec[i].PlayerID] = &simplevec[i]
		mgr.NameIDMap[simplevec[i].Name] = simplevec[i].PlayerID
	}

	return true
}

func GetPlayerIDByAccountID(accountid int32) int32 {
	G_SimpleMgr.SimpleLock.Lock()
	defer G_SimpleMgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[accountid]
	if ok && pInfo != nil {
		return pInfo.PlayerID
	}

	return 0
}

func (mgr *TSimpleInfoMgr) GetPlayerIDByName(name string) int32 {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	playerid, ok := mgr.NameIDMap[name]
	if ok {
		return playerid
	}

	return 0
}

func (mgr *TSimpleInfoMgr) GetPlayerLogoffTime(playerid int32) int64 {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	if ok && pInfo != nil {
		return pInfo.LogoffTime
	}

	return 0
}

func (mgr *TSimpleInfoMgr) GetSimpleInfoByID(playerid int32) *TSimpleInfo {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	if ok && pInfo != nil {
		return pInfo
	}

	gamelog.Error3("GetSimpleInfoByID Error , Invalid playerid:%d", playerid)
	return nil
}

func (mgr *TSimpleInfoMgr) GetPlayerAwardCenterID(playerid int32) int {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()

	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	if ok && pInfo != nil {
		return pInfo.AwardCenterID
	}

	return 0
}

func (mgr *TSimpleInfoMgr) RemoveSimpleInfo(playerid int32) bool {
	mgr.SimpleLock.Lock()
	defer mgr.SimpleLock.Unlock()
	delete(mgr.SimpleList, playerid)
	return true
}

func (mgr *TSimpleInfoMgr) Get_FightValue(playerid int32) int {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		return pInfo.FightValue
	}
	return 0
}

func (mgr *TSimpleInfoMgr) Set_FightValue(playerid int32, fightvalue int, level int) bool {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && (pInfo.FightValue != fightvalue || pInfo.Level != level) {
		pInfo.FightValue = fightvalue
		pInfo.Level = level
		G_SimpleMgr.DB_SetFightValue(playerid, fightvalue, pInfo.Level)

		return true
	}

	return false
}

func (mgr *TSimpleInfoMgr) Set_PlayerName(playerid int32, name string) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.Name != name {
		pInfo.Name = name
		G_SimpleMgr.DB_SetPlayerName(playerid, name)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_LogoffTime(playerid int32, time int64) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.LogoffTime != time {
		pInfo.LogoffTime = time
		G_SimpleMgr.DB_SetLogoffTime(playerid, time)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_HeroID(playerid int32, heroid int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.HeroID != heroid {
		pInfo.HeroID = heroid
		G_SimpleMgr.DB_SetHeroID(playerid, heroid)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_HeroQuality(playerid int32, quality int8) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.Quality != quality {
		pInfo.Quality = quality
		G_SimpleMgr.DB_SetHeroQuality(playerid, quality)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_VipLevel(playerid int32, viplevel int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.VipLevel != viplevel {
		pInfo.VipLevel = viplevel
		G_SimpleMgr.DB_SetVipLevel(playerid, viplevel)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_AwardCenterID(playerid int32, awardCenterID int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.AwardCenterID != awardCenterID {
		pInfo.AwardCenterID = awardCenterID
		G_SimpleMgr.DB_SetAwardCenterID(playerid, awardCenterID)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_BatCamp(playerid int32, camp int8) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.BatCamp != camp {
		pInfo.BatCamp = camp
		G_SimpleMgr.DB_SetBatCamp(playerid, camp)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_LoginDay(playerid int32, day uint32) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.LoginDay != day {
		pInfo.LoginDay = day
		G_SimpleMgr.DB_SetLoginDay(playerid, day)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Set_GuildID(playerid int32, guildid int) {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil && pInfo.GuildID != guildid {
		pInfo.GuildID = guildid
		G_SimpleMgr.DB_SetGuildID(playerid, guildid)
		return
	}
	return
}

func (mgr *TSimpleInfoMgr) Get_GuildID(playerid int32) int {
	mgr.SimpleLock.Lock()
	pInfo, ok := G_SimpleMgr.SimpleList[playerid]
	mgr.SimpleLock.Unlock()
	if ok && pInfo != nil {
		return pInfo.GuildID
	}
	return 0
}

package sessionmgr

import (
	"fmt"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

var (
	SessionKeyMap map[int32]string = make(map[int32]string, 1024)
	LoginTimeMap  map[int32]int32  = make(map[int32]int32, 1024)
	SessionMutex  sync.Mutex
)

func AddSessionKey(playerid int32, sessionkey string) {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()

	SessionKeyMap[playerid] = sessionkey
	LoginTimeMap[playerid] = utility.GetCurTime()
}

func NewSessionKey() string {
	return bson.NewObjectId().Hex()
}

func CheckLoginTime(playerid int32) bool {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()
	logintime, ok := LoginTimeMap[playerid]
	if !ok {
		return true
	}

	if (utility.GetCurTime() - logintime) < 5 {
		return false
	}

	return true
}

func CheckSessionKey(playerid int32, sessionkey string) bool {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()
	key, ok := SessionKeyMap[playerid]
	if !ok {
		return false
	}

	if key == sessionkey {
		return true
	}

	fmt.Println("session key is not valid !!!")
	return false
}

func DeleteSessionKey(playerid int32) bool {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()
	_, ok := SessionKeyMap[playerid]
	if !ok {
		return true
	}

	delete(SessionKeyMap, playerid)

	return true
}

package sessionmgr

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	SessionKeyMap map[int32]string = make(map[int32]string, 1024)
	LoginTimeMap  map[int32]int64  = make(map[int32]int64, 1024)
	SessionMutex  sync.Mutex
)

func AddSessionKey(playerid int32, sessionkey string) {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()

	SessionKeyMap[playerid] = sessionkey
	LoginTimeMap[playerid] = time.Now().Unix()
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

	if (time.Now().Unix() - logintime) < 5 {
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

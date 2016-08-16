package sessionmgr

import (
	"fmt"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	SessionKeyMap map[int]string = make(map[int]string, 1024)
	LoginTimeMap  map[int]int64  = make(map[int]int64, 1024)
	SessionMutex  sync.Mutex
)

func AddSessionKey(playerid int, sessionkey string) {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()

	SessionKeyMap[playerid] = sessionkey
	LoginTimeMap[playerid] = time.Now().Unix()
}

func NewSessionKey() string {
	return bson.NewObjectId().Hex()
}

func CheckLoginTime(playerid int) bool {
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

func CheckSessionKey(playerid int, sessionkey string) bool {
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

func DeleteSessionKey(playerid int) bool {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()
	_, ok := SessionKeyMap[playerid]
	if !ok {
		return true
	}

	delete(SessionKeyMap, playerid)

	return true
}

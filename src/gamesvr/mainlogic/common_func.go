package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/sessionmgr"
	"msg"
	"time"
)

//! 常规检查
func GetPlayerAndCheck(playerID int, sessionKey string, url string) (*TPlayer, int) {
	//! 检查SessionKey
	if !sessionmgr.CheckSessionKey(playerID, sessionKey) {
		gamelog.Error("message %s error : Invalid Session Key!", url)
		return nil, msg.RE_INVALID_SESSIONKEY
	}

	//! 获取玩家信息
	var pPlayer *TPlayer = GetPlayerByID(playerID)
	if pPlayer == nil {
		gamelog.Error("message %s error : Invalid Player ID :%d!", url, playerID)
		return nil, msg.RE_INVALID_PLAYERID
	}

	return pPlayer, msg.RE_UNKNOWN_ERR
}

//! 获取开服天数
func GetOpenServerDay() int {
	now := time.Now().Unix()

	if now < appconfig.GameOpenSvrTime {
		gamelog.Error("GetOpenServerDay Error : Invalid Open Server Time")
		return 1
	}

	day := (now-appconfig.GameOpenSvrTime)/86400 + 1

	if day <= 0 {
		day = 1
	}

	return int(day)
}

//获取当前服务器ID
func GetCurServerID() int {
	return appconfig.DomainID
}

//获取当前服务器名称
func GetCurServerName() string {
	return appconfig.DomainName
}

//! 自定义类型
type IntLst []int

func (self *IntLst) IsExist(value int) int {
	for i := 0; i < len(*self); i++ {
		if value == (*self)[i] {
			return i
		}
	}
	return -1
}

func (self *IntLst) Add(value int) {
	*self = append(*self, value)
}

func (self IntLst) Len() int {
	return len(self)
}

func (self IntLst) Less(i int, j int) bool {
	return (self)[i] < (self)[j]
}

func (self IntLst) Swap(i int, j int) {
	temp := (self)[i]
	(self)[i] = (self)[j]
	(self)[j] = temp
}

type Int64Lst []int64

func (self *Int64Lst) IsExist(value int64) int {
	for i := 0; i < len(*self); i++ {
		if value == (*self)[i] {
			return i
		}
	}
	return -1
}

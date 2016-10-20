package mainlogic

import (
	"sync"
)

var (
	G_StatisticMutex sync.Mutex
	G_OnlineCnt      int //在线玩家数量
	G_MaxOnlineCnt   int //最高在线数据
	G_RegisterCnt    int //总注册玩家
)

//增加在线人数
func IncOnlineCnt() {
	G_StatisticMutex.Lock()
	G_OnlineCnt = G_OnlineCnt + 1
	if G_MaxOnlineCnt < G_OnlineCnt {
		G_MaxOnlineCnt = G_OnlineCnt
	}
	G_StatisticMutex.Unlock()
}

//减少在线人数
func DecOnlineCnt() {
	G_StatisticMutex.Lock()
	G_OnlineCnt = G_OnlineCnt - 1
	G_StatisticMutex.Unlock()
}

//增加注册人数
func IncRegisterCnt() {
	G_StatisticMutex.Lock()
	G_RegisterCnt = G_RegisterCnt + 1
	G_StatisticMutex.Unlock()
}

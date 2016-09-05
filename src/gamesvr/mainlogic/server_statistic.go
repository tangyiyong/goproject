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

func IncOnlineCnt() {
	G_StatisticMutex.Lock()
	G_OnlineCnt = G_OnlineCnt + 1
	G_StatisticMutex.Unlock()
}

func DecOnlineCnt() {
	G_StatisticMutex.Lock()
	G_OnlineCnt = G_OnlineCnt - 1
	G_StatisticMutex.Unlock()
}

func IncRegisterCnt() {
	G_StatisticMutex.Lock()
	G_RegisterCnt = G_RegisterCnt + 1
	G_StatisticMutex.Unlock()
}

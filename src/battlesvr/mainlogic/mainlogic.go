package mainlogic

import (
	"battlesvr/gamedata"
)

func Init() bool {
	//初始化连接管理器
	InitConMgr()

	//初始化房间管理器
	InitRoomMgr()

	gamedata.LoadConfig()

	return true
}

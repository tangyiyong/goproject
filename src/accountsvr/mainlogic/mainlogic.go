package mainlogic

import (
	"appconfig"
	"mongodb"
)

func Init() {
	//初始化数据库处理器
	mongodb.InitDbProcesser(appconfig.AccountDbName)

	//初始化游戏服管理对象
	InitGameSvrMgr()

	//初始化账号管理器
	InitAccountMgr()

	//初始化礼包码管理器
	InitGiftCodeMgr()
}

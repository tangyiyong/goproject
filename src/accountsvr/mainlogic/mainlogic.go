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

	//初始化黑白名单管理器
	InitNetMgr()

	//初始化封禁账户管理器
	InitDisableMgr()
}

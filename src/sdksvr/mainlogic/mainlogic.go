package mainlogic

import (
	"appconfig"
)

func Init() {

	//连接数据库
	G_DbConn.Open(appconfig.SdkDataSource)

	//初始化游戏服管理对象
	InitGameSvrMgr()

}

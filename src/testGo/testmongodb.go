// main
package main

import (
	"gamelog"
	//"gamesvr/mainlogic"
	"mongodb"

	//"strings"
	//"utility"
	"appconfig"
)

var ()

func TestMongoDB() {
	//初始化日志系统
	//加载配制文件
	appconfig.LoadConfig()
	gamelog.InitLogger("test", 0)
	mongodb.Init("localhost:27017")
}

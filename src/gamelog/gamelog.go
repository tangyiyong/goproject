package gamelog

import (
	"os"
	"time"
	"utility"
)

var (
	g_logDir = utility.GetCurrPath() + "log\\"
)

func InitLogger(name string, bScreen bool) {
	var err error = nil
	if !utility.IsDirExists(g_logDir) {
		err = os.MkdirAll(g_logDir, os.ModePerm)
	}
	if err != nil {
		panic("InitLogger error : " + err.Error())
		return
	}

	timeStr := time.Now().Format("20060102_150405")
	logFileName := g_logDir + name + "_" + timeStr + ".log"

	InitDebugLog(logFileName, bScreen)
}

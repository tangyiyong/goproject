package mainlogic

import (
	"appconfig"
	"gamelog"
)

func InitLogMgr() {

}

type TLog interface {
	Start() bool
	WriteLog(pdata []byte)
	Flush()
	Close()
	SetFlushCnt(cnt int)
}

var (
	G_SvrLogMgr = make([]TLog, 1000) // svrid不会超过一定数目
)

const (
	FT_BINIARY = 1 //二进制文件类型
	FT_MYSQL   = 2 //MySql文件类型
)

func CreateLogFile(svrid int32) {

	var logfile TLog = nil
	switch appconfig.LogFileType {
	case FT_BINIARY:
		logfile = CreateBinaryFile(appconfig.LogFileName, svrid)
	case FT_MYSQL:
		logfile = CreateMysqlFile(appconfig.LogFileName, svrid)
	}

	if logfile == nil {
		gamelog.Error("CreateLogFile Failed type:%d!!!", appconfig.LogFileType)
		return
	}

	if false == logfile.Start() {
		gamelog.Error("CreateLogFile Error, Log Start Failed!!!")
		return
	}

	logfile.SetFlushCnt(appconfig.LogSvrFlushCnt)

	G_SvrLogMgr[svrid] = logfile

	gamelog.Error("CreateLogFile Successed!!!")
	return
}
func WriteSvrLog(pdata []byte, svrid int32) {
	if G_SvrLogMgr[svrid] == nil {
		gamelog.Error("WriteSvrLog Error: Invalid svrid : %d", svrid)
		return
	}

	G_SvrLogMgr[svrid].WriteLog(pdata)
}

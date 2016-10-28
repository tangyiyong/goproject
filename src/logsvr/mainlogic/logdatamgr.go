package mainlogic

import (
	"appconfig"
	"gamelog"
	"time"
)

var (
	G_LogFile TLog = nil
	G_LogChan chan []byte
)

func InitLogMgr() bool {
	G_LogFile = new(TMysqlLog)
	if G_LogFile == nil {
		gamelog.Error("InitLogMgr Error: new TMysqlLog Failed!!")
		return false
	}

	G_LogFile.SetFlushCnt(appconfig.LogSvrFlushCnt)
	if false == G_LogFile.Start(appconfig.LogFileName, 0) {
		gamelog.Error("InitLogMgr Error, Log Start Failed!!!")
		return false
	}

	G_LogChan = make(chan []byte, 10240)
	go LogRoutine()
	go TimerRoutine()
	return true
}

func LogRoutine() {
	for logdata := range G_LogChan {
		G_LogFile.WriteLog(logdata)
	}
}

func TimerRoutine() {
	regtimer := time.Tick(60 * time.Second)
	for {
		G_LogFile.WriteLog(nil)
		<-regtimer
	}
}

type TLog interface {
	Start(file string, svrid int32) bool
	WriteLog(pdata []byte)
	Flush()
	Close()
	SetFlushCnt(cnt int)
}

func WriteSvrLog(pdata []byte, svrid int32) {
	G_LogChan <- pdata
}

//var (
//	G_SvrLogMgr = make([]TLog, 1000) // svrid不会超过一定数目
//)

//const (
//	FT_BINIARY = 1 //二进制文件类型
//	FT_MYSQL   = 2 //MySql文件类型
//)

//func CreateLogFile(svrid int32) {
//	var logfile TLog = nil
//	switch appconfig.LogFileType {
//	case FT_BINIARY:
//		logfile = new(TBinaryLog)
//	case FT_MYSQL:
//		logfile = new(TMysqlLog)
//	}

//	if logfile == nil {
//		gamelog.Error("CreateLogFile Failed type:%d!!!", appconfig.LogFileType)
//		return
//	}

//	logfile.SetFlushCnt(appconfig.LogSvrFlushCnt)
//	if false == logfile.Start(appconfig.LogFileName, svrid) {
//		gamelog.Error("CreateLogFile Error, Log Start Failed!!!")
//		return
//	}

//	G_SvrLogMgr[svrid] = logfile
//	gamelog.Error("CreateLogFile Successed!!!")
//	return
//}

//func WriteSvrLog(pdata []byte, svrid int32) {
//	if G_SvrLogMgr[svrid] == nil {
//		gamelog.Error("WriteSvrLog Error: Invalid svrid : %d", svrid)
//		return
//	}

//	G_SvrLogMgr[svrid].WriteLog(pdata)
//}

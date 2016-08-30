package mainlogic

import (
	"appconfig"
	"gamelog"
)

type TLog interface {
	Start() bool
	WriteLog(pdata []byte)
	Flush()
	Close()
}

const (
	FT_BINIARY = 1 //二进制文件类型
	FT_MYSQL   = 2 //MySql文件类型
)

var G_LogFile TLog

func InitLogMgr() bool {
	switch appconfig.LogFileType {
	case FT_BINIARY:
		{ //二进制文件
			G_LogFile = CreateBinaryFile(appconfig.LogFileName)
		}
	case FT_MYSQL:
		{ //mysql数据库
			G_LogFile = CreateMysqlFile(appconfig.LogFileName)
		}
	}

	if G_LogFile == nil {
		return false
	}

	if false == G_LogFile.Start() {
		gamelog.Error("InitLogMgr Error")
		return false
	}

	return true
}

func AppendLog(pdata []byte) {
	G_LogFile.WriteLog(pdata)
}

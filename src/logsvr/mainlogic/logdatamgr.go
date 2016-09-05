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

func NewOneFile(svrid int32) (file TLog) {
	switch appconfig.LogFileType {
	case FT_BINIARY:
		file = CreateBinaryFile(appconfig.LogFileName, svrid)
	case FT_MYSQL:
		file = CreateMysqlFile(appconfig.LogFileName, svrid)
	}

	if file == nil || false == file.Start() {
		gamelog.Error("NewOneFile Error, Start Log Failed!!!")
		return
	}

	file.SetFlushCnt(appconfig.LogSvrFlushCnt)

	G_SvrLogMgr[svrid] = file
	return
}
func WriteSvrLog(pdata []byte, svrid int32) {
	if v := G_SvrLogMgr[svrid]; v != nil {
		v.WriteLog(pdata)
	}
	gamelog.Error("WriteSvrLog Error: svrid(%d)", svrid)
}

package mainlogic

var (
	G_LogDataMgr *AsyncLog

	g_binaryLog *TBinaryLog
	g_mysqlLog  *TMysqlLog
)

func InitLogMgr() {
	G_LogDataMgr = NewAsyncLog(1024, _doWriteBinaryLog)

	g_binaryLog = NewBinaryLog("logsvr")
	g_mysqlLog = NewMysqlLog()
	if g_binaryLog == nil || g_mysqlLog == nil {
		panic("logsvr InitMgr fail!")
		return
	}
}

func _doWriteBinaryLog(data1, data2 [][]byte) {
	g_binaryLog.Write(data1, data2)
}
func _doWriteMysqlLog(data1, data2 [][]byte) {
	g_mysqlLog.Write(data1, data2)
}

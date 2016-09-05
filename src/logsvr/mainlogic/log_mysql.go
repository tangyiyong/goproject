package mainlogic

import (
	"database/sql"
	"fmt"
	"gamelog"
	"msg"
	_ "mysql"
	"time"
)

type TMysqlLog struct {
	db       *sql.DB
	tx       *sql.Tx
	query    string
	writeCnt int
	flushCnt int
}

func (self *TMysqlLog) Start() bool {
	self.tx, _ = self.db.Begin()
	return true
}

func (self *TMysqlLog) WriteLog(pdata []byte) {
	stmt, _ := self.tx.Prepare(self.query)
	var req msg.MSG_SvrLogData
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("MysqlLog::WriteLog : Message Reader Error!!!!")
		return
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	stmt.Exec(req.EventID, req.PlayerID, req.SvrID, timeStr, req.Param[0], req.Param[1], req.Param[2], req.Param[3])
	stmt.Close()

	self.writeCnt++
	if self.writeCnt >= self.flushCnt {
		self.Flush()
	}
}

func (self *TMysqlLog) Close() {
	self.tx.Commit()
	self.db.Close()
	self.writeCnt = 0
}

func (self *TMysqlLog) Flush() {
	self.tx.Commit()
	self.tx, _ = self.db.Begin()
}

func (self *TMysqlLog) SetFlushCnt(cnt int) {
	self.flushCnt = cnt
	if cnt <= 0 {
		self.flushCnt = 100
	}
}

func CreateMysqlFile(filename string, svrid int32) *TMysqlLog {
	var log TMysqlLog
	var err error = nil
	log.db, err = sql.Open("mysql", filename)
	if err != nil {
		gamelog.Error("Create MysqlLog Error : %s", err.Error())
		return nil
	}
	log.query = fmt.Sprintf("INSERT log_%d SET EventID=?,SrcID=?,TargetID=?,Time=?,Param1=?,Param2=?,Param3=?,Param4=?", svrid)

	return &log
}

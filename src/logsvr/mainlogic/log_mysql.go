package mainlogic

import (
	"database/sql"
	"fmt"
	"gamelog"
	"msg"
	_ "mysql"
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
	stmt.Exec(req.SvrID, req.EventID, req.PlayerID, req.Time, req.Param[0], req.Param[1], req.Param[2], req.Param[3])
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

func CreateMysqlFile(filename string, svrid int32) TLog {
	var log TMysqlLog
	var err error = nil
	log.db, err = sql.Open("mysql", filename)
	if err != nil {
		gamelog.Error("CreateMysqlFile Error : %s", err.Error())
		return nil
	}

	err = log.db.Ping()
	if err != nil {
		gamelog.Error("CreateMysqlFile Error : db.ping : %s", err.Error())
		return nil
	}

	log.query = fmt.Sprintf("INSERT log_%d SET SvrID=?,EventID=?,PlayerID=?,Time=?,Param1=?,Param2=?,Param3=?,Param4=?", svrid)

	return &log
}

package mainlogic

import (
	"database/sql"
	"fmt"
	"gamelog"
	"msg"
	_ "mysql"
	"time"
)

// create table logsvr
// (
// 	ID int unsigned not null auto_increment primary key,
// 	EventID int not null,
// 	SrcID int not null,
// 	TargetID int,
// 	Time datetime,
// 	Param1 int,
// 	Param2 int,
// 	Param3 int,
// 	Param4 int
// );
const (
	g_user     = "root"
	g_password = ""
	g_addr     = "localhost:3306"
	g_dbname   = "mysql"
	g_table    = "logsvr"
)

type TMysqlLog struct {
	db *sql.DB
}

func NewMysqlLog() *TMysqlLog {
	var err error = nil
	log := new(TMysqlLog)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=30s&strict=true", g_user, g_password, g_addr, g_dbname)
	log.db, err = sql.Open("mysql", dsn)
	if err != nil {
		gamelog.Error("NewBinaryLog : %s", err.Error())
		return nil
	}
	return log
}
func (self *TMysqlLog) Close() {
	self.db.Close()
}
func (self *TMysqlLog) Write(data1, data2 [][]byte) {
	query := fmt.Sprintf("INSERT %s SET EventID=?,SrcID=?,TargetID=?,Time=?,Param1=?,Param2=?,Param3=?,Param4=?", g_table)

	//1、开启事务
	tx, _ := self.db.Begin()
	stmt, _ := tx.Prepare(query)

	//2、编辑数据
	for _, v := range data1 {
		self._transaction(stmt, v)
	}
	for _, v := range data2 {
		self._transaction(stmt, v)
	}

	//3、提交事务
	tx.Commit()
}
func (self *TMysqlLog) _transaction(stmt *sql.Stmt, pdata []byte) {
	// 将buf解析为结构体
	var req msg.MSG_SvrLogData
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("MysqlLog::InsertDB : Message Reader Error!!!!")
		return
	}
	// timeStr := time.Unix(req.Time, 0).Format("2006-01-02 15:04:05")
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	stmt.Exec(req.EventID, req.SrcID, req.TargetID, timeStr, req.Param[0], req.Param[1], req.Param[2], req.Param[3])
}

func (self *TMysqlLog) InsertDB(pdata *msg.MSG_SvrLogData) {
	// 直接插入数据库
	query := fmt.Sprintf("INSERT %s SET EventID=?,SrcID=?,TargetID=?,Time=?,Param1=?,Param2=?,Param3=?,Param4=?", g_table)
	stmt, err := self.db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		gamelog.Error("MysqlLog::Prepare : %s", err.Error())
		return
	}
	// timeStr := time.Unix(pdata.Time, 0).Format("2006-01-02 15:04:05")
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(pdata.EventID, pdata.SrcID, pdata.TargetID, timeStr, pdata.Param[0], pdata.Param[1], pdata.Param[2], pdata.Param[3])
	if err != nil {
		gamelog.Error("MysqlLog::Exec : %s", err.Error())
		return
	}
}

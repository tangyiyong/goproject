package mainlogic

import (
	"fmt"
	"gamelog"
	"msg"
	_ "mysql"
	"testing"
	"time"
)

func Test_MysqlLog_Insert(t *testing.T) {
	log := NewMysqlLog()

	req := msg.MSG_SvrLogData{}
	req.EventID = 111111
	req.SrcID = 2222
	req.TargetID = 3333

	// 插入数据库
	query := fmt.Sprintf("INSERT %s SET EventID=?,SrcID=?,TargetID=?,Time=?,Param1=?,Param2=?,Param3=?,Param4=?", g_table)
	stmt, err := log.db.Prepare(query)
	// timeStr := time.Unix(req.Time, 0).Format("2006-01-02 15:04:05")
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	defer stmt.Close()
	if err != nil {
		gamelog.Error("MysqlLog::Prepare : %s", err.Error())
		return
	}
	for i := 0; i < 1000; i++ {

		_, err = stmt.Exec(req.EventID, req.SrcID, req.TargetID, timeStr, req.Param[0], req.Param[1], req.Param[2], req.Param[3])
		if err != nil {
			gamelog.Error("MysqlLog::Exec : %s", err.Error())
			return
		}
	}
}
func Test_MysqlLog_Affair(t *testing.T) {
	log := NewMysqlLog()

	req := msg.MSG_SvrLogData{}
	req.EventID = 111111
	req.SrcID = 2222
	req.TargetID = 3333

	query := fmt.Sprintf("INSERT %s SET EventID=?,SrcID=?,TargetID=?,Time=?,Param1=?,Param2=?,Param3=?,Param4=?", g_table)
	// timeStr := time.Unix(req.Time, 0).Format("2006-01-02 15:04:05")
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	//开启事务
	Tx, _ := log.db.Begin()
	stm, _ := Tx.Prepare(query)
	for i := 0; i < 1000; i++ {
		stm.Exec(req.EventID, req.SrcID, req.TargetID, timeStr, req.Param[0], req.Param[1], req.Param[2], req.Param[3])
	}
	Tx.Commit()
}

package mainlogic

import (
	"database/sql"
	"gamelog"
	_ "mysql"
)

type TMysql struct {
	db      *sql.DB
	datasrc string
}

var G_DbConn TMysql

func (self *TMysql) Open(datasrc string) bool {
	var err error = nil
	self.datasrc = datasrc
	self.db, err = sql.Open("mysql", self.datasrc)
	if err != nil {
		gamelog.Error("TMysql Open file Error : %s", err.Error())
		return false
	}

	err = self.db.Ping()
	if err != nil {
		gamelog.Error("TMysql ping Error : %s", err.Error())
		return false
	}

	return true
}

func (self *TMysql) Exec(sql string, args ...interface{}) sql.Result {
	stmt, err := self.db.Prepare(sql)
	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
	}

	return result
}

func (self *TMysql) Close() {
	self.db.Close()
}

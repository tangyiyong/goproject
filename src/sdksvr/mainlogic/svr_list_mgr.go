package mainlogic

import (
	"database/sql"
	"gamelog"
)

type TGameServerInfo struct {
	SvrID        int32  `bson:"_id"` //账号ID
	SvrName      string //服务器名字
	SvrOutAddr   string //外部地址
	SvrInnerAddr string //内部地址
}

var (
	G_ServerList  [10000]TGameServerInfo
	G_RecommendID int32 //推荐的服务器ID
)

func InitGameSvrMgr() {
	rows, err := G_DbConn.db.Query("select * from gamesvrlist limit 10000")
	if err != nil && err != sql.ErrNoRows {
		//! 创建表
		sql := `CREATE TABLE if not exists gamesvrlist(
			id int not null primary key,
			name varchar(32),
			outaddr varchar(32),
			inneraddr varchar(32));`

		G_DbConn.Exec(sql)
	} else {

		for rows.Next() {
			var svrID int32
			var name string
			var outAddr string
			var innerAddr string

			err = rows.Scan(&svrID, &name, &outAddr, &innerAddr)
			if err != nil {
				panic(err)
			}

			G_ServerList[svrID].SvrID = svrID
			G_ServerList[svrID].SvrName = name
			G_ServerList[svrID].SvrOutAddr = outAddr
			G_ServerList[svrID].SvrInnerAddr = innerAddr
		}
	}

	return
}

func UpdateGameSvrInfo(svrid int32, svrname string, outaddr string, inaddr string) {
	if svrid <= 0 || svrid >= 10000 {
		gamelog.Error("UpdateGameSvrInfo Error : Invalid svrid:%d", svrid)
		return
	}

	if G_ServerList[svrid].SvrID == 0 {
		G_ServerList[svrid].SvrID = svrid
		G_ServerList[svrid].SvrName = svrname
		G_ServerList[svrid].SvrInnerAddr = inaddr
		G_ServerList[svrid].SvrOutAddr = outaddr

		sql := `INSERT INTO gamesvrlist
				(id,
				name,
				outaddr,
				inneraddr)
				VALUES
				(?,?,?,?);`

		G_DbConn.Exec(sql, svrid, svrname, outaddr, inaddr)
	} else {
		G_ServerList[svrid].SvrName = svrname
		G_ServerList[svrid].SvrInnerAddr = inaddr
		G_ServerList[svrid].SvrOutAddr = outaddr

		sql := `UPDATE gamesvrlist
				SET
				name = ?,
				outaddr = ?,
				inneraddr = ?
				WHERE id = ?;
		`
		G_DbConn.Exec(sql, svrname, outaddr, inaddr, svrid)
	}

}

func GetGameSvrOutAddr(svrid int32) string {
	if G_ServerList[svrid].SvrID == 0 {
		gamelog.Error("GetGameSvrAddr Error Invalid svrid :%d", svrid)
		return ""
	}

	return G_ServerList[svrid].SvrOutAddr
}

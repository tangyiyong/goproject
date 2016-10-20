package appconfig

import (
	"fmt"
	"strings"
)

type TGM_User struct {
	SessionID  string //管理员的ID
	SessionKey string //管理员的Key
	IpAddr     string //管理员的登录IP限制，为"0", 表示管理员可以来自任意的IP, 否则管理员只能来自指定的IP
}

var G_Users map[string]*TGM_User

func ParseGmUser(uinfo string) {
	slice := strings.Split(uinfo, ",")
	if len(slice) < 3 {
		panic("ParseConfigValue Invalid Gm User Data")
	}

	if G_Users == nil {
		G_Users = make(map[string]*TGM_User)
	}

	var pGmUser *TGM_User = new(TGM_User)
	pGmUser.SessionID = strings.TrimSpace(slice[0])
	pGmUser.SessionKey = strings.TrimSpace(slice[1])
	pGmUser.IpAddr = strings.TrimSpace(slice[2])

	G_Users[pGmUser.SessionID] = pGmUser
}

func CheckGmRight(id string, key string, ip string) bool {
	//GM ID是否存在
	pUser, ok := G_Users[id]
	if !ok || pUser == nil {
		fmt.Println("CheckGmRight Error: ***ID", id)
		return false
	}

	//SessionKey　是否对得上
	if key != pUser.SessionKey {
		fmt.Println("CheckGmRight Error: ***req.Key", key, "local.key", pUser.SessionKey)
		return false
	}

	//GM来源的IP是否符合要求
	if pUser.IpAddr != "0" && pUser.IpAddr != ip {
		fmt.Println("CheckGmRight Error: ***req.ip", ip, "local.key", pUser.IpAddr)
		return false
	}

	return true
}

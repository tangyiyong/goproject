package msg

//登录请求
//消息:/login
type MSG_Login_Req struct {
	Name     string //账户名
	Password string //密码
}

type MSG_Login_Ack struct {
	RetCode     int    //返回码 0:成功 1: 账号不存在 2: 密码不正确
	AccountID   int    //账号ID
	LoginKey    string //登录key
	LastSvrID   int    //上次登录SvrID
	LastSvrName string //上次登录svrName
	LastSvrAddr string //上次登录svr address
}

//注册账号请求
//消息:/register
type MSG_RegAccount_Req struct {
	Name     string //账户名
	Password string //密码
}

type MSG_RegAccount_Ack struct {
	RetCode int //返回码 0:成功 1: 无效的账号名称 2: 无效的密码 3:账号名己存在
}

//请求服务器列表
//消息:/serverlist
type MSG_ServerList_Req struct {
	AccountID   int    //账号ID
	AccountName string //账户名
	LoginKey    string //登录key
}

type ServerNode struct {
	SvrDomainID   int
	SvrDomainName string
	SvrState      int
	UpdateTime    int64
	SvrOutAddr    string
}

type MSG_ServerList_Ack struct {
	RetCode int
	SvrList []ServerNode //服务器结点表
}

//验证玩家登录请求
//消息:/verifyuserlogin
type MSG_VerifyUserLogin_Req struct {
	AccountID   int    //账号ID
	AccountName string //账户名
	LoginKey    string //登录key
	DomainID    int    //服务器ID
}

type MSG_VerifyUserLogin_Ack struct {
	RetCode int //返回码 0:成功 1: 失败
}

//游客玩家注册
//消息:/tourist_register
type MSG_TourRegAccount_Req struct {
	Name     string //账户名
	Password string //密码
}

type MSG_TourRegAccount_Ack struct {
	RetCode  int    //返回码
	Name     string //账户名
	Password string //密码
}

//邦定游客账号
//消息:/bind_tourist
type MSG_BindTourist_Req struct {
	Name        string //账户名
	Password    string //密码
	NewName     string //新账号名
	NewPassword string //新密码
}

type MSG_BindTourist_Ack struct {
	RetCode  int    //返回码
	Name     string //新账号名
	Password string //新密码
}

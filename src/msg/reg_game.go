package msg

//定义游戏服向账号服注册的消息

//游戏服向账号服务器注册消息
type MSG_RegToAccountSvr_Req struct {
	ServerDomainID   int //
	ServerDomainName string
	ServerOuterAddr  string
	ServerInnerAddr  string
}

//游戏服向账号服务器注册的返回消息
type MSG_RegToAccountSvr_Ack struct {
	RetCode int
}

//游戏服向账号服务器注册消息
type MSG_RegToCrossSvr_Req struct {
	ServerDomainID   int //
	ServerDomainName string
	ServerOuterAddr  string
	ServerInnerAddr  string
}

//游戏服向账号服务器注册的返回消息
type MSG_RegToCrossSvr_Ack struct {
	RetCode int
}

package msg

//! 玩家查询称号状态
//! 消息: /get_title
type MSG_GetTitle_Req struct {
	PlayerID   int32
	SessionKey string
}

type TitleInfo struct {
	TitleID int   //! 拥有称号ID
	EndTime int64 //! 结束时间
	Status  int   //! 0->未激活 1->已激活 2->已佩戴
}

type MSG_GetTitle_Ack struct {
	RetCode  int
	TitleLst []TitleInfo
}

//! 玩家请求激活称号
//! 消息: /activate_title
type MSG_ActivateTitle_Req struct {
	PlayerID   int32
	SessionKey string
	TitleID    int
}

type MSG_ActivateTitle_Ack struct {
	RetCode int
}

//! 玩家请求装备称号
//! 消息: /equi_title
type MSG_EquiTitle_Req struct {
	PlayerID   int32
	SessionKey string
	TitleID    int
}

type MSG_EquiTitle_Ack struct {
	RetCode int
}

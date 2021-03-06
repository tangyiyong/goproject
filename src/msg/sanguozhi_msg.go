package msg

//! 获取三国志信息
//! 消息: /get_sanguozhi_info
type MSG_GetSanGuoZhiInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSanGuoZhiInfo_Ack struct {
	RetCode   int
	CurOpenID int32 //! 当前开启ID (之前所有星都已命星)
}

//! 命星
//! 消息: /set_sanguozhi
type MSG_SetSanGuoZhi_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_SetSanGuoZhi_Ack struct {
	RetCode    int
	Quality    int8         //! 主角品质
	AwardItem  MSG_ItemData //! 物品
	FightValue int32        //! 战力
}

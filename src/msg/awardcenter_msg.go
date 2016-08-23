package msg

//! 玩家请求查询奖励中心信息
//! 消息: /query_award_center
type MSG_AwardCenter_Query_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_AwardCenter_Data struct {
	ID       int
	TextType int            //! 类型
	Value    []string       //! 参数
	ItemLst  []MSG_ItemData //! 奖励内容
	Time     int64          //! 发放奖励时间戳
}

type MSG_AwardCenter_Query_Ack struct {
	RetCode  int
	AwardLst []MSG_AwardCenter_Data
}

//! 玩家请求领取奖励中心奖励
//! 消息: /get_award_center
type MSG_AwardCenter_Get_Req struct {
	PlayerID   int32
	SessionKey string
	AwardID    int
}

type MSG_AwardCenter_Get_Ack struct {
	RetCode int
}

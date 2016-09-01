package msg

//! 玩家请求送花
//! 消息: /send_flower
type MSG_SendFlower_Req struct {
	PlayerID   int32
	SessionKey string
	SendIndex  int32 //! 0-5 对应 1-6名
	SendType   int   //! 0是战力排行榜 1是等级排行榜
}

type MSG_SendFlower_Ack struct {
	RetCode    int
	RankType   int //! 0是战力排行榜 1是等级排行榜
	CharmValue []int
}

//! 玩家请求查询前六名魅力值
//! 消息: /get_charm
type MSG_GetCharm_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_CharmPlayerInfo struct {
	HeroID     int
	CharmValue int
	Name       string
	FightValue int
	PlayerID   int32
	Level      int
}

type MSG_GetCharm_Ack struct {
	RetCode      int
	FightRankLst []MSG_CharmPlayerInfo
	LevelRankLst []MSG_CharmPlayerInfo

	Times       int //! 次数
	FightSendID []int32
	LevelSendID []int32
}

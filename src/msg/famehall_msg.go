package msg

//! 玩家请求送花
//! 消息: /send_flower
type MSG_SendFlower_Req struct {
	PlayerID   int
	SessionKey string
	SendIndex  int //! 0-5 对应 1-6名
	SendType   int //! 0是战力排行榜 1是等级排行榜
}

type MSG_SendFlower_Ack struct {
	RetCode    int
	RankType   int //! 0是战力排行榜 1是等级排行榜
	CharmValue []int
}

//! 玩家请求查询前六名魅力值
//! 消息: /get_charm
type MSG_GetCharm_Req struct {
	PlayerID   int
	SessionKey string
	RankType   int //! 0-> 战力排行榜 1->等级排行榜
}

type MSG_GetCharm_Ack struct {
	RetCode    int
	HeroID     []int
	CharmValue []int
	Name       []string
	FightValue []int
	PlayerID   []int
	Level      []int

	Times  int //! 次数
	SendID []int
}

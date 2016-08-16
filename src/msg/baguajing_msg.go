package msg

//! 玩家使用八卦镜
//! 消息: /use_baguajing
type MSG_UseBaguajing_Req struct {
	PlayerID   int
	SessionKey string
	BagPos     int //! 背包英雄索引
	HeroID     int //! 兑换英雄ID
}

type MSG_UseBaguajing_Ack struct {
	RetCode int
}

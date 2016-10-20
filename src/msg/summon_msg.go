package msg

//! 玩家请求查询召唤刷新状态
//! 消息: /get_summon_status
type MSG_GetSummonStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

//! 普通召唤
type MSG_NormalSummon struct {
	SummonCounts int   //! 今日还可免费召唤次数
	SummonTime   int32 //! 下次可免费召唤时间戳
}

//! 高级召唤
type MSG_SeniorSummon struct {
	Point       int   //! 当前召唤积分
	SummonTime  int32 //! 下次可免费召唤时间戳
	OrangeCount int   //! 还有N次可以得橙将
}

type MSG_GetSummonStatus_Ack struct {
	RetCode      int
	NormalSummon MSG_NormalSummon
	SeniorSummon MSG_SeniorSummon
	Discount     int //! 元宝十连抽折扣
}

//! 玩家请求召唤
/*
const (
	Summon_Normal   = 1 //! 普通招贤
	Summon_Senior   = 2 //! 高级招贤
)
*/
//! 消息: /get_summon
type MSG_GetSummon_Req struct {
	PlayerID   int32
	SessionKey string
	SummonType int //! 召唤类型
	NumberType int //! 0-> 单抽  1-> 十连抽
}

type MSG_GetSummon_Ack struct {
	RetCode      int
	HeroID       []int //! 召唤英雄ID
	NormalSummon MSG_NormalSummon
	SeniorSummon MSG_SeniorSummon
	IsFree       bool //! 是否免费
}

//! 玩家请求积分兑换指定英雄
//! 消息: /exchange_hero
type MSG_ExchangeHero_Req struct {
	PlayerID   int32
	SessionKey string
	HeroID     int //! 兑换英雄ID
}

type MSG_ExchangeHero_Ack struct {
	RetCode int
}

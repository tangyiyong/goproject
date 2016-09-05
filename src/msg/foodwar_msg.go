package msg

//! 玩家请求夺粮战挑战玩家列表
//! 消息: /get_foodwar_challenger
type MSG_FoodWar_GetChallenger_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FoodWar_Challenger struct {
	PlayerID   int32
	PlayerName string
	HeroID     int
	Quality    int8
	CanRobFood int
	TotalFood  int
	FightValue int32
	Level      int
}

type MSG_FoodWar_GetChallenger_Ack struct {
	RetCode       int
	ChallengerLst []MSG_FoodWar_Challenger
}

//! 玩家请求查询夺粮战状态
//! 消息: /get_foodwar_status
type MSG_FoodWar_GetStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FoodWar_GetStatus_Ack struct {
	RetCode        int
	Rank           int
	RobFood        int
	TotalFood      int
	ChallengerLst  []MSG_FoodWar_Challenger
	AttackTimes    int   //! 剩余抢夺次数
	BuyAttackTimes int   //! 已购买抢夺次数
	RecoverTime    int64 //! 下次恢复时间
	AttackTimesMax int   //! 上限
	FixFood        int   //! 固定粮草
	RecoverFood    int   //! 每小时增加固定粮草数
}

//! 玩家请求查询复仇状态
//! 消息: /get_foodwar_revenge_status
type MSG_FoodWar_RevengeStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FoodWar_RevengeStatus_Ack struct {
	RetCode         int                      //! 返回码
	RevengeTimes    int                      //! 剩余复仇次数
	BuyRevengeTimes int                      //! 已购买复仇次数
	RevengeLst      []MSG_FoodWar_Challenger //! 复仇名单
	RecoverTime     int64                    //! 下次恢复时间
}

//! 玩家请求查询次数以及恢复时间
//! 消息: /get_foodwar_time
type MSG_FoodWar_GetFoodWarTime_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FoodWar_GetFoodWarTime_Ack struct {
	RetCode         int
	AttackTimes     int   //! 剩余抢夺次数
	RevengeTimes    int   //! 剩余复仇次数
	BuyAttackTimes  int   //! 已购买抢夺次数
	BuyRevengeTimes int   //! 已购买复仇次数
	RecoverTime     int64 //! 下次恢复时间
}

//! 玩家请求掠夺粮草
//! 消息: /rob_food
type MSG_FoodWar_RobFood_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int32
	IsWin          int //! 1->胜利 0->失败
}

type MSG_FoodWar_RobFood_Ack struct {
	RetCode       int
	RobFood       int //! 掠夺粮草
	TotalFood     int //! 总粮草
	MoneyID       int
	MoneyNum      int
	Rank          int //! 名次
	ChallengerLst []MSG_FoodWar_Challenger
}

//! 玩家请求复仇
//! 消息: /revenge_rob
type MSG_FoodWar_RevengeRob_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int32
	IsWin          int //! 1->胜利 0->失败
}

type MSG_FoodWar_RevengeRob_Ack struct {
	RetCode   int
	RobFood   int //! 掠夺粮草
	FixFood   int //! 固定粮草
	TotalFood int //! 总粮草
	MoneyID   int
	MoneyNum  int
	Rank      int //! 名次
}

//! 玩家请求查询排行榜信息
//! 消息: /get_food_rank
type MSG_FoodWar_GetRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FoodWar_GetRank_Ack struct {
	RetCode int
	RankLst []MSG_FoodWar_Challenger
}

//! 请求购买次数
//! 消息: /buy_food_times
type MSG_FoodWar_BuyTimes_Req struct {
	PlayerID   int32
	SessionKey string
	TimesType  int //! 1->掠夺次数 2->复仇次数
	Times      int //! 购买次数
}

type MSG_FoodWar_BuyTimes_Ack struct {
	RetCode      int
	CostMoneyID  int
	CostMoneyNum int
}

//! 请求领取粮草奖励
//! 消息: /recv_food_award
type MSG_FoodWar_GetFoodAward_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_FoodWar_GetFoodAward_Ack struct {
	RetCode int
	Award   int
}

//! 请求查询粮草奖励领取状况
//! 消息: /query_food_award
type MSG_FoodWar_QueryFoodAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FoodWar_QueryFoodAward_Ack struct {
	RetCode  int
	AwardLst []int
}

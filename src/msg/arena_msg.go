package msg

//! 玩家竞技场挑战结果
//! 消息: /arena_result
type MSG_ArenaResult_Req struct {
	PlayerID   int32
	SessionKey string
	Rank       int //! 挑战的玩家排名
	IsVictory  int //! 是否胜利
}

type MSG_ArenaResult_Ack struct {
	RetCode     int
	IsVictory   int
	DropItem    []MSG_ItemData //! 翻牌奖励 第一个为获取的奖励 其余两个是翻牌未获取的奖励
	ExtraAward  MSG_ItemData   //! 额外元宝奖励
	HistoryRank int            //! 历史排名
	SelfRank    int            //! 自己的排名
}

type MSG_ArenaPlayerInfo struct {
	PlayerID   int32
	Rank       int
	Level      int
	Name       string
	HeroID     int
	Quality    int8
	FightValue int
}

//! 玩家竞技场挑战次数
//! 消息: /arena_battle
type MSG_ArenaBattle_Req struct {
	PlayerID   int32
	SessionKey string
	Rank       int //! 挑战的玩家名次
	IsUseItem  int //! 使用道具 1-> 使用 0-> 不使用
}

type MSG_ArenaBattle_Ack struct {
	RetCode   int
	IsVictory bool
	ItemID    int
	ItemNum   int
	Exp       int
	Money     int //! 银币
	Money2    int //! 声望
}

//! 玩家请求竞技场信息
//! 消息: /get_arena_info
type MSG_GetArenaInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetArenaInfo_Ack struct {
	RetCode     int
	PlayerLst   []MSG_ArenaPlayerInfo
	SelfRank    int   //! 自己的排名
	HistoryRank int   //! 历史排名
	IDLst       []int //! 已购买物品列表
}

//! 玩家请求声望商店购买商品
//! 消息: /arena_store_buy_item
type MSG_GetArenaStoreItem_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int //! 表中唯一标识
	Num        int //! 购买次数
}

type MSG_GetArenaStoreItem_Ack struct {
	RetCode int
}

//! 玩家请求已购买的声望商店奖励ID列表
//! 消息: /arena_store_query_award
type MSG_QueryArenaStoreAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryArenaStoreAward_Ack struct {
	RetCode int
	IDLst   []int
}

//! 玩家竞技场挑战检测
//! 消息:/arena_check
type MSG_ArenaCheck_Req struct {
	PlayerID   int32
	SessionKey string
	Rank       int //! 挑战的玩家排名
}

type MSG_ArenaCheck_Ack struct {
	RetCode    int
	TargetType int            // 1 : 玩家数据 2 : 机器人数据
	Name       string         //目标的名字
	PlayerData MSG_PlayerData //! 玩家武将信息
}

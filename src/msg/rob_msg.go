package msg

//! 玩家请求抢劫名单
//! 消息: /get_rob_list
type MSG_GetRobList_Req struct {
	PlayerID   int32
	SessionKey string
	TreasureID int //! 宝物碎片ID
}

type MSG_RobPlayerInfo struct {
	PlayerID int32  //! 玩家ID
	Name     string //! 名字
	Level    int    //! 等级
	HeroID   [6]int //! 英雄ID
	IsRobot  int    //! 机器人标记
}

type MSG_GetRobList_Ack struct {
	RetCode int
	Lst     []MSG_RobPlayerInfo
}

//! 玩家请求换批对手
//! 消息: /refresh_rob_list
type MSG_RefreshRobList_Req struct {
	PlayerID   int32
	SessionKey string
	TreasureID int     //! 宝物碎片ID
	CurRobLst  []int32 //! 现存表中的抢劫名单玩家ID
}

type MSG_RefreshRobList_Ack struct {
	RetCode int
	Lst     []MSG_RobPlayerInfo
}

//! 玩家请求抢夺
//! 消息: /rob_treasure
type MSG_RobTreasure_Req struct {
	PlayerID    int32
	SessionKey  string
	TreasureID  int   //! 宝物碎片ID
	RobPlayerID int32 //! 抢劫玩家ID
	IsRobot     int   //! 是否为机器人 0->玩家 1->机器人
}

type MSG_RobTreasure_Ack struct {
	RetCode     int
	RobSuccess  bool //! 返回抢劫是否成功
	MoneyID     int  //! 获取货币
	MoneyNum    int
	Exp         int             //! 获取经验
	DropItem    [3]MSG_ItemData //! 掉落物品
	FreeWarTime int64
}

//! 玩家请求免战时间
//! 消息: /get_free_war_time
type MSG_FreeWarTime_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FreeWarTime_Ack struct {
	RetCode     int
	FreeWarTime int64
}

//! 玩家请求宝物合成
//! 消息: /treasure_composed
type MSG_TreasureComposed_Req struct {
	PlayerID   int32
	SessionKey string
	GemID      int //! 合成的宝物ID
	Num        int //! 合成宝物个数
}

type MSG_TreasureComposed_Ack struct {
	RetCode int
	GemID   int //! 需要合成的宝物ID
	Num     int //! 需求合成宝物个数
}

//! 玩家请求熔炼
//! 消息: /treasure_melting
type MSG_TreasureMelting_Req struct {
	PlayerID      int32
	SessionKey    string
	GemPos        int
	TargetPieceID int
}

type MSG_TreasureMelting_Ack struct {
	RetCode int
}

//! 玩家请求抢劫玩家详细信息
//! 消息: /get_rob_hero_info
type MSG_GetRobPlayerInfo_Req struct {
	PlayerID    int32
	SessionKey  string
	RobPlayerID int32 //! 抢劫玩家ID
	IsRobot     int   //! 是否为机器人
	TreasureID  int   //! 宝物碎片ID
}

type MSG_GetRobPlayerInfo_Ack struct {
	RetCode    int
	PlayerData MSG_PlayerData //! 玩家武将信息
}

//! 一键合成协议
//! 消息: /one_key_composed
type MSG_OneKeyRob_Req struct {
	PlayerID   int32
	SessionKey string
	GemID      int //! 宝物ID
	IsUseItem  int //! 0->未使用  1->已使用
}

type MSG_OneKeyRob_Ack struct {
	RetCode    int
	PieceID    int  //! 碎片ID
	RobSuccess bool //! 返回抢劫是否成功
	MoneyID    int  //! 获取货币
	MoneyNum   int  //! 货币数量
	Exp        int  //! 获取经验
	ItemID     int  //! 掉落物品ID
	ItemNum    int  //! 掉落物品数量
}

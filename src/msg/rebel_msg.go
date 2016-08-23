package msg

//! 玩家请求获取叛军信息
//! 消息: /get_rebel_info
type MSG_GetRebelInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_RebelInfo struct {
	PlayerID   int32  //! 玩家ID
	RebelID    int    //! 叛军ID
	Level      int    //! 叛军等级
	CurLife    int    //! 叛军当前血量
	FindName   string //! 发现者名称
	EscapeTime int64  //! 剩余逃走时间
	IsShare    bool   //! 是否分享 0->不分享 1->分享
}

type MSG_GetRebelInfo_Ack struct {
	RetCode     int
	InfoLst     []MSG_RebelInfo
	Exploit     int
	TopDamage   int
	DamageRank  int //! 伤害排行榜 若为0则未进榜
	ExploitRank int //! 功勋排行榜 若为0则未进榜
}

//! 玩家请求攻击叛军
//! 消息: /attack_rebel
type MSG_Attack_Rebel_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int32
	Damage         int //! 伤害
	AttackType     int //! 1->普通攻击 2->全力一击
}

type MSG_Attack_Rebel_Ack struct {
	RetCode     int
	IsKill      int //! 是否击杀叛军  0->未击杀 1->已击杀
	DamageRank  int //! 伤害排行榜 若为0则未进榜
	ExploitRank int //! 功勋排行榜 若为0则未进榜
}

//! 玩家请求分享叛军
//! 消息: /share_rebel
type MSG_ShareRebel_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_ShareRebel_Ack struct {
	RetCode int //! 返回码
}

//! 玩家请求功勋奖励领取状态
//! 消息: /get_exploit_award_status
type MSG_GetExploitAwardStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetExploitAwardStatus_Ack struct {
	RetCode int
	RecvLst []int //! 已领取的功勋奖励ID
}

//! 玩家请求领取功勋奖励
//! 消息: /get_exploit_award
type MSG_GetExploitAward_Req struct {
	PlayerID       int32
	SessionKey     string
	ExploitAwardID int //! 功勋奖励ID
}

type MSG_GetExploitAward_Ack struct {
	RetCode int
}

//! 玩家请求购买战功商店物品
//! 消息: /buy_rebel_store
type MSG_BuyRebelStore_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int //! 商品ID 取静态表
	Num        int //! 购买个数
}

type MSG_BuyRebelStore_Ack struct {
	RetCode int
}

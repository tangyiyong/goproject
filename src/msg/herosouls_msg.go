package msg

//! 激活将灵链接
//! 消息: /activate_herosouls
type MSG_ActivateHeroSouls_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int //! 需要激活的链接ID
}

type MSG_ActivateHeroSouls_Ack struct {
	RetCode       int
	UnLockChapter int //! 解锁下一章节
}

//! 详细查询将灵
//! 消息: /query_herosouls_chapter
type MSG_QueryHeroSoulsChapter_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_HeroSoulsLink struct {
	ID    int //! 将灵链接ID
	Level int //! 将灵等级
}

type MSG_QueryHeroSoulsChapter_Ack struct {
	RetCode       int
	HeroSouls     []MSG_HeroSoulsLink
	UnLockChapter int
}

//! 获取将灵轮盘
//! 消息: /get_herosouls_lst
type MSG_GetHeroSoulsLst_Req struct {
	PlayerID   int32
	SessionKey string
}

type THeroSouls struct {
	ID      int  //! 唯一ID
	HeroID  int  //! 将灵ID
	IsExist bool //! true 存在 false 不存在
}

type MSG_GetHeroSoulsLst_Ack struct {
	RetCode           int
	HeroSoulsLst      []THeroSouls //! 可挑战将灵
	TargetIndex       int          //! 指针指向索引
	BuyChallengeTimes int          //! 已购买挑战次数
	ChallengeTimes    int          //! 剩余挑战次数
	RefreshMoneyNum   int          //! 刷新
	CountDown         int          //! 倒计时
}

//! 刷新将灵轮盘指针
//! 消息: /refresh_herosouls
type MSG_RefreshHeroSoulsLst_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_RefreshHeroSoulsLst_Ack struct {
	RetCode      int
	TargetIndex  int
	CostMoneyID  int //! 花费的货币ID
	CostMoneyNum int //! 花费的货币数量
	GetMoneyID   int //! 获取英魂ID
	GetMoneyNum  int //! 获取英魂数
}

//! 挑战将灵
//! 消息: /challenge_herosouls
type MSG_ChallengeHeroSouls_Req struct {
	PlayerID   int32
	SessionKey string
	//英雄核查数据
	HeroCkD []MSG_HeroCheckData
}

type MSG_ChallengeHeroSouls_Ack struct {
	RetCode int
}

//! 购买挑战英灵次数
//! 消息: /buy_challenge_herosouls
type MSG_BuyChallengeHeroSoulsTimes_Req struct {
	PlayerID   int32
	SessionKey string
	Times      int //! 购买次数
}

type MSG_BuyChallengeHeroSoulsTimes_Ack struct {
	RetCode      int
	CostMoneyID  int
	CostMoneyNum int
}

//! 重置轮盘
//! 消息: /reset_herosouls_lst
type MSG_ResetHeroSoulsLst_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_ResetHeroSoulsLst_Ack struct {
	RetCode      int
	HeroSoulsLst []THeroSouls //! 可挑战将灵
	TargetIndex  int          //! 指针指向索引
}

//! 查询将灵排行榜
//! 消息: /query_herosouls_rank
type MSG_QueryHeroSoulsRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type THeroSoulsRank struct {
	HeroID     int
	Name       string
	Quality    int8
	SoulsValue int //! 阵图值
	SoulsCount int //! 阵图个数
}

type MSG_QueryHeroSoulsRank_Ack struct {
	RetCode        int
	RankLst        []THeroSoulsRank
	SelfRank       int
	SelfSoulsValue int //! 阵图值
	SelfSoulsCount int //! 阵图个数
}

//! 查询将灵商店信息
//! 消息: /query_herosouls_store
type MSG_QueryHeroSoulsStore_Req struct {
	PlayerID   int32
	SessionKey string
}

type THeroSoulsStore struct {
	ItemID   int  //! 商品ID
	IsBuy    bool //! 是否已经购买
	MoneyID  int  //! 货币ID
	MoneyNum int  //! 货币数量
}

type MSG_QueryHeroSoulsStore_Ack struct {
	RetCode   int
	CountDown int               //! 倒计时
	GoodsLst  []THeroSoulsStore //! 商品列表
}

//! 购买指定将灵
//! 消息: /buy_herosouls
type MSG_BuyHeroSouls_Req struct {
	PlayerID   int32
	SessionKey string
	ItemID     int
}

type MSG_BuyHeroSouls_Ack struct {
	RetCode      int
	CostMoneyID  int
	CostMoneyNum int
}

//! 查询英灵阵图成就
//! 消息: /query_herosouls_achievement
type MSG_QueryHeroSoulsAchievement_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryHeroSoulsAchievement_Ack struct {
	RetCode      int
	SoulMapValue int
	Achievement  int //! 当前成就
}

//! 激活下一阵图成就
//! 消息: /activate_herosouls_achievement
type MSG_ActivateHeroSoulsAchievement_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_ActivateHeroSoulsAchievement_Ack struct {
	RetCode int
}

//! 查询属性加成汇总
//! 消息: /query_herosouls_property
type MSG_QueryHeroSoulsProperty_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryHeroSoulsProperty_Ack struct {
	RetCode                int
	PropertyIntLst         [11]int //! 加成实际数值
	PropertyPercentLst     [11]int //! 加成百分比
	CampPropertyKillLst    [4]int  //! 对阵营加伤百分比
	CampPropertyDefenceLst [4]int  //! 对阵营减伤百分比
}

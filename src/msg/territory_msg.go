package msg

//! 玩家请求当前领地状态
//! 消息: /get_territory_status
type MSG_GetTerritoryStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_TerritoryInfo struct {
	ID              int                     //! 领地ID
	HeroID          int                     //! 英雄ID
	AwardItem       []MSG_ItemData          //! 获得奖励
	PatrolBeginTime int64                   //! 巡逻开启时间
	PatrolType      int                     //! 巡逻类型
	PatrolEndTime   int64                   //! 巡逻结束时间
	SkillLevel      int                     //! 领地技能等级
	RiotInfo        []MSG_TerritoryRiotData //! 暴动信息
}

type MSG_GetTerritoryStatus_Ack struct {
	RetCode           int
	TerritoryLst      []MSG_TerritoryInfo
	SuppressRiotTimes int //! 当前镇压暴动次数
	TotalPatrolTime   int //! 总计巡逻时间(小时)
	RiotTime          int //! 暴动持续时间
}

//! 玩家回馈挑战领地结果
//! 消息: /challenge_territory
type MSG_ChallengeTerritory_Req struct {
	PlayerID    int
	SessionKey  string
	TerritoryID int
}

type MSG_ChallengeTerritory_Ack struct {
	RetCode int
}

//! 玩家请求好友领地状态
//! 消息: /get_friend_territory_status
type MSG_GetFriendTerritoryStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_FriendTerritoryStatus struct {
	PlayerID      int                 //! 好友ID
	Level         int                 //! 好友等级
	Quality       int                 //! 好友品质
	TerritoryLst  []MSG_TerritoryInfo //! 领地信息
	LastLoginTime int64               //! 上次登录时间
}

type MSG_GetFriendTerritoryStatus_Ack struct {
	RetCode    int
	FriendInfo []MSG_FriendTerritoryStatus
}

//! 玩家请求查看好友领地详情
//! 消息: /get_friend_territory_info
type MSG_GetFriendTerritoryInfo_Req struct {
	PlayerID    int
	SessionKey  string
	FriendID    int //! 朋友ID
	TerritoryID int //! 领地ID
}

//! 暴动信息
type MSG_TerritoryRiotData struct {
	BeginTime  int64  //! 开始时间
	DealTime   int64  //! 处理时间
	HelperName string //! 帮忙处理好友姓名
}

type MSG_GetFriendTerritoryInfo_Ack struct {
	RetCode   int
	HeroID    int
	AwardInfo []MSG_ItemData          //! 领地当前奖励信息
	RiotInfo  []MSG_TerritoryRiotData //! 暴动信息
}

//! 玩家请求帮忙好友镇压暴动
//! 消息: /help_riot
type MSG_HelpRiot_Req struct {
	PlayerID          int
	SessionKey        string
	TargetID          int //! 好友ID
	TargetTerritoryID int //! 好友领地ID
}

type MSG_HelpRiot_Ack struct {
	RetCode int
	ItemID  int //! 奖励物品ID
	ItemNum int //! 奖励物品数量
}

//! 玩家请求收获领地奖励
//! 消息: /get_territory_award
type MSG_GetTerritoryAward_Req struct {
	PlayerID    int
	SessionKey  string
	TerritoryID int
}

type MSG_GetTerritoryAward_Ack struct {
	RetCode int
}

//! 玩家置放武将到领地巡逻
//! 消息: /patrol_territory
type MSG_PatrolTerritory_Req struct {
	PlayerID    int
	SessionKey  string
	HeroID      int
	TerritoryID int
	PatrolType  int //! 巡逻类型
	AwardType   int //! 奖励类型
}

type MSG_PatrolTerritory_Ack struct {
	RetCode         int
	AwardItem       []MSG_ItemData          //! 获得奖励
	PatrolBeginTime int64                   //! 巡逻开启时间
	RiotInfo        []MSG_TerritoryRiotData //! 暴动信息
	ActionValue     int                     //! 行动力值
	ActionTime      int64                   //! 行动力恢复起始时间
}

//! 玩家请求升级领地技能
//! 消息: /territory_skill_up
type MSG_TerritorySkillUp_Req struct {
	PlayerID    int
	SessionKey  string
	TerritoryID int
}

type MSG_TerritorySkillUp_Ack struct {
	RetCode int
}

//! 玩家请求查询领地暴动状态
//! 消息: /query_territory_riot
type MSG_GetTerritoryRiot_Req struct { //! 需要每次请求
	PlayerID    int
	SessionKey  string
	TerritoryID int
}

type MSG_GetTerritoryRiot_Ack struct {
	RetCode  int
	RiotInfo []MSG_TerritoryRiotData //! 暴动信息
}

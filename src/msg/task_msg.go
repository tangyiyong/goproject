package msg

//玩家请求日常任务数据
//消息:/get_tasks
type MSG_GetTasks_Req struct {
	PlayerID   int
	SessionKey string
}

//角色表结构
type TTaskInfo struct {
	TaskID     int //! 任务ID
	TaskStatus int //! 任务状态 0-> 未完成 1-> 已完成  2-> 已领取
	TaskCount  int //! 任务次数
}

type MSG_GetTasks_Ack struct {
	RetCode int
	Tasks   []TTaskInfo
}

//! 玩家请求日常任务完成奖励
//! 消息: /receive_task
type MSG_GetTaskAward_Req struct {
	PlayerID   int
	SessionKey string
	TaskID     int
}

type MSG_GetTaskAward_Ack struct {
	RetCode   int
	TaskScore int //! 现在任务积分
	ItemLst   []MSG_ItemData
}

//! 玩家请求任务积分宝箱奖励
//! 消息: /receive_taskscore
type MSG_GetTaskScoreAward_Req struct {
	PlayerID     int
	SessionKey   string
	ScoreAwardID int //! 请求领取积分宝箱ID
}

type MSG_GetTaskScoreAward_Ack struct {
	RetCode       int
	ScoreAwardLst []MSG_ScoreAward //! 现在积分宝箱领取状态
	ItemLst       []MSG_ItemData
}

//! 玩家请求任务积分信息
//! 消息: /get_taskscore
type MSG_GetTaskScores_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_ScoreAward struct {
	ScoreAwardID int
	Status       bool
}

type MSG_GetTaskScores_Ack struct {
	RetCode       int
	TaskScore     int              //! 当前任务积分
	ScoreAwardLst []MSG_ScoreAward //! 积分宝箱ID
}

//! 获取当前成就任务
//! 消息: /get_achievement
type MSG_GetAchievementAll_Req struct {
	PlayerID   int
	SessionKey string
}

//! 角色成就表结构
type TAchievementInfo struct {
	ID         int //! 成就ID
	TaskStatus int //! 成就达成状态 0-> 未完成 1-> 已完成  2-> 已领取
	TaskCount  int //! 成就达成次数
}

type MSG_GetAchievementAll_Ack struct {
	RetCode int
	List    []TAchievementInfo
}

//! 请求成就奖励
//! 消息: /receive_achievement
type MSG_GetAchievementAward_Req struct {
	PlayerID      int
	SessionKey    string
	AchievementID int
}

type MSG_GetAchievementAward_Ack struct {
	RetCode    int
	NewAchieve TAchievementInfo //! 替换已完成的成就
}

//! 玩家请求开服天数
//! 消息: /get_open_server_day
type MSG_GetOpenServerDay_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetOpenServerDay_Ack struct {
	RetCode int
	OpenDay int
}

//! 玩家请求七日活动进度信息
//! 消息: /get_seven_activity
type MSG_GetSevenActivity_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
}

type MSG_GetSevenActivity_Ack struct {
	RetCode          int
	SevenActivityLst []TTaskInfo //! 只发送进度不为0的任务信息,结构参考上面TTaskInfo结构
}

//! 玩家请求领取七日活动奖励
//! 消息: /get_seven_activity_award
type MSG_GetSevenActivityAward_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
	TaskID     int
	ItemID     int //! 对应三选一奖励中,物品奖励ID
}

type MSG_GetSevenActivityAward_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

//! 玩家请求半价限购剩余件数与自己已购买信息
//! 消息: /get_seven_activity_limit_num
type MSG_GetSevenActivityLimitInfo_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
}

type MSG_GetSevenActivityLimitInfo_Ack struct {
	RetCode   int
	LimitInfo [7]int
	BuyLst    []int
}

//! 玩家请求购买半价限购
//! 消息: /buy_seven_activity_limit
type MSG_BuySevenActivityLimitItem_Req struct {
	PlayerID   int
	SessionKey string
	ActivityID int
	OpenDay    int //! 购买第几天的限购商品
}

type MSG_BuySevenActivityLimitItem_Ack struct {
	RetCode  int
	BuyTimes int //! 当前物品已购买次数
}

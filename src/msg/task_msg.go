package msg

//玩家请求日常任务数据
//消息:/get_tasks
type MSG_GetTaskData_Req struct {
	PlayerID   int32
	SessionKey string
}

//角色表结构
type MSG_TaskInfo struct {
	ID     int //! 任务ID
	Status int //! 任务状态 0-> 未完成 1-> 已完成  2-> 已领取
	Count  int //! 任务次数
}

type MSG_ScoreAward struct {
	ScoreAwardID int
	Status       bool
}

type MSG_GetTaskData_Ack struct {
	RetCode       int
	Tasks         []MSG_TaskInfo   //! 日常任务积分信息
	TaskScore     int              //! 当前任务积分
	ScoreAwardLst []MSG_ScoreAward //! 积分宝箱ID
	List          []MSG_TaskInfo   //! 成就信息
}

//! 玩家请求日常任务完成奖励
//! 消息: /receive_task
type MSG_RecvTaskAward_Req struct {
	PlayerID   int32
	SessionKey string
	TaskID     int
}

type MSG_RecvTaskAward_Ack struct {
	RetCode   int
	TaskScore int //! 现在任务积分
	ItemLst   []MSG_ItemData
}

//! 玩家请求任务积分宝箱奖励
//! 消息: /receive_taskscore
type MSG_RecvTaskScoreAward_Req struct {
	PlayerID     int32
	SessionKey   string
	ScoreAwardID int //! 请求领取积分宝箱ID
}

type MSG_RecvTaskScoreAward_Ack struct {
	RetCode       int
	ScoreAwardLst []MSG_ScoreAward //! 现在积分宝箱领取状态
	ItemLst       []MSG_ItemData
}

//! 请求成就奖励
//! 消息: /receive_achievement
type MSG_RecvAchievementAward_Req struct {
	PlayerID      int32
	SessionKey    string
	AchievementID int
}

type MSG_RecvAchievementAward_Ack struct {
	RetCode    int
	NewAchieve MSG_TaskInfo //! 替换已完成的成就
}

//! 玩家请求开服天数
//! 消息: /get_open_server_day
type MSG_GetOpenServerDay_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetOpenServerDay_Ack struct {
	RetCode int
	OpenDay int
}

//! 玩家请求七日活动进度信息
//! 消息: /get_seven_activity
type MSG_GetSevenActivity_Req struct {
	PlayerID   int32
	SessionKey string
	ActivityID int32
}

type MSG_GetSevenActivity_Ack struct {
	RetCode          int
	ActivityID       int32
	OpenDay          int            //! 当天为第几天
	SevenActivityLst []MSG_TaskInfo //! 只发送进度不为0的任务信息,结构参考上面TTaskInfo结构
	LimitInfo        [7]int
	BuyLst           []int
}

//! 玩家请求领取七日活动奖励
//! 消息: /get_seven_activity_award
type MSG_GetSevenActivityAward_Req struct {
	PlayerID   int32
	SessionKey string
	ActivityID int32
	TaskID     int
	ItemID     int //! 对应三选一奖励中,物品奖励ID
}

type MSG_GetSevenActivityAward_Ack struct {
	RetCode    int
	ActivityID int32
	ItemLst    []MSG_ItemData
}

//! 玩家请求购买半价限购
//! 消息: /buy_seven_activity_limit
type MSG_BuySevenActivityLimitItem_Req struct {
	PlayerID   int32
	SessionKey string
	ActivityID int32
	OpenDay    int //! 购买第几天的限购商品
}

type MSG_BuySevenActivityLimitItem_Ack struct {
	RetCode    int
	ActivityID int32
	BuyTimes   int //! 当前物品已购买次数
}

package msg

type MSG_Target struct {
	PlayerID   int32  //角色ID
	HeroID     int    //英雄ID
	Name       string //角色名
	FightValue int    //战力
	Level      int    //等级
	SvrID      int    //服务器ID
	SvrName    string //服务器名
	Quality    int8   //品质
}

//获取积分赛主界面信息
//get_score_data
type MSG_GetScoreData_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetScoreData_Ack struct {
	RetCode   int           //返回值
	Rank      int           //排名
	Score     int           //积分
	FightTime int           //今天战斗次数
	WinTime   int           //连胜次数
	IsRecv    int           //是否领取连胜次数奖励
	BuyTime   int           //己购买战斗次数
	Targets   []MSG_Target  //三个战斗目标
	ItemLst   []MSG_BuyData //购买物品次数
}

type MSG_ScoreRankInfo struct {
	PlayerID   int32  //! 角色ID
	HeroID     int    //! 英雄ID
	Quality    int8   //! 品质
	Name       string //! 角色名字
	FightValue int    //! 战力值
	Score      int    //！积分值
	SvrID      int    //! 服务器ID
	SvrName    string //! 服务器名字
}

//! 玩家请求积分赛排行榜
//! 消息: /get_score_rank
type MSG_GetScoreRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetScoreRank_Ack struct {
	RetCode       int                 //返回码
	ScoreRankList []MSG_ScoreRankInfo //功勋榜
	MyRank        int                 //自己的排名
	MyScore       int                 //自己的积分
}

//请求发送积分赛结果
//set_score_battle_result
type MSG_SetScoreBattleResult_Req struct {
	PlayerID    int32
	SessionKey  string
	TargetIndex int //目标玩家索引
	WinBattle   int //是否占胜对手 1: 胜, 0: 失败
}

type MSG_SetScoreBattleResult_Ack struct {
	RetCode int
	Targets []MSG_Target //三个战斗目标
	Rank    int          //最新排名
	Score   int          //最新积分
}

//! 玩家购买积分赛战斗次数
//! 消息: /buy_score_fight_time
type MSG_BuyScoreTime_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_BuyScoreTime_Ack struct {
	RetCode     int   //!返回码
	ActionID    int   //!行动力ID
	ActionValue int   //!行动力值
	ActionTime  int64 //!行动力恢复起始时间
}

//! 玩家请求积分赛战斗次数奖励
//! 消息: /get_score_time_award
type MSG_GetScoreTimeAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetScoreTimeAward_Ack struct {
	RetCode   int   //返回码
	FightTime int   //己战斗次数
	Awards    []int //己获取的奖励ID
}

//! 玩家收取积分赛战斗次数奖励
//! 消息: /recv_score_time_award
type MSG_RecvScoreTimeAward_Req struct {
	PlayerID    int32
	SessionKey  string
	TimeAwardID int //次数奖励的ID
}

type MSG_RecvScoreTimeAward_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //奖励的物品
}

//! 玩家领取连赢次数奖励
//! 消息: /recv_score_continue_award
type MSG_RecvScoreContinueAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_RecvScoreContinueAward_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //奖励的物品
}

//! 玩家请求购买积分商店道具
//! 消息: /buy_score_store_item
type MSG_BuyScoreStoreItem_Req struct {
	PlayerID    int32
	SessionKey  string
	StoreItemID int32 //商店道具ID
	BuyNum      int   //购买数量
}

type MSG_BuyScoreStoreItem_Ack struct {
	RetCode int //返回码
}

//! 游戏服向积分服请求积分排行榜
//! 消息: /cross_query_score_rank
type MSG_CrossQueryScoreRank_Req struct {
}

type MSG_CrossQueryScoreRank_Ack struct {
	RetCode       int                 //返回码
	MyRank        int                 //本人的排行榜
	ScoreRankList []MSG_ScoreRankInfo //功勋榜
}

//! 游戏服向积分服请求目标玩家排行榜
//! 消息: /cross_query_score_target
type MSG_CrossQueryScoreTarget_Req struct {
	PlayerID   int32  //角色ID
	HeroID     int    //英雄ID
	SvrID      int    //服务器ID
	SvrName    string //服务器名
	Score      int    //当前的积分
	Level      int    //等级
	FightValue int    //战力
	Quality    int8   //战力
	PlayerName string //角色名

}

type MSG_CrossQueryScoreTarget_Ack struct {
	RetCode    int           //返回码
	NewRank    int           //新的排名
	TargetList [3]MSG_Target //目标列表
}

//! 跨服务器向游戏服请求选择可战玩家信息
//! 消息: /game_selct_player
type MSG_GameSelectPlayer_Req struct {
}

type MSG_GameSelectPlayer_Ack struct {
	RetCode int        //返回码
	Target  MSG_Target //目标玩家
}

//! 向跨服服务器请求战斗目标数据
type MSG_GetFightTarget_Req struct {
	SvrID    int   //服务器ID
	PlayerID int32 //角色ID
}

//! 积分赛挑战检测
//! 消息:/get_score_battle_check
type MSG_GetScoreBattleCheck_Req struct {
	PlayerID    int32
	SessionKey  string
	TargetIndex int //! 挑战的玩家索引
}

type MSG_GetScoreBattleCheck_Ack struct {
	RetCode    int            //返回码
	PlayerData MSG_PlayerData //目标玩家
}

type MSG_GetFightTarget_Ack struct {
	RetCode    int            //返回码
	PlayerData MSG_PlayerData //目标玩家
}

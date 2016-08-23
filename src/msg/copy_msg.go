package msg

//! 获取副本数据
//消息: /get_copy_data
type MSG_GetCopyData_Req struct {
	PlayerID   int32
	SessionKey string
}

type TFamousInfo struct {
	CurCopyID   int //! 名将副本
	CurChapter  int //! 名将当前章节
	BattleTimes int //! 名将总共挑战次数
}

type MSG_GetCopyData_Ack struct {
	RetCode        int
	CopyMainInfo   TMainInfo       //! 主线
	CopyEliteInfo  TEliteCopyData  //! 精英
	CopyFamousInfo TFamousInfo     //! 名将
	CopyDailyInfo  []MSG_DailyCopy //! 日常副本
}

//玩家副本战斗检查
//消息:/battle_check
type MSG_BattleCheck_Req struct {
	PlayerID   int32
	SessionKey string
	CopyType   int //副本类型
	CopyID     int //副本ID
	Chapter    int //副本章节
}

type MSG_BattleCheck_Ack struct {
	RetCode int
}

//玩家副本战斗结果
//消息:/battle_result
type MSG_BattleResult_Req struct {
	PlayerID   int32
	SessionKey string
	CopyType   int //副本类型
	CopyID     int //副本ID
	Chapter    int //副本章节
	StarNum    int //战斗星数
}

type MSG_ItemData struct {
	ID  int //! 掉落物品ID
	Num int //! 掉落物品数量
}

type MSG_BattleResult_Ack struct {
	RetCode     int
	ItemLst     []MSG_ItemData
	FirstItem   []MSG_ItemData //! 首胜奖励 目前只对应名将副本
	IsFindRebel bool           //! 是否发现叛军 目前只对应主线副本 发现后发送get_rebel_find_info获取叛军信息
	OpenEndTime int64          //! 是否发现黑市
	Exp         int            //! 获取经验
	ActionValue int            //! 行动力值
	ActionTime  int64          //! 行动力恢复起始时间
}

//! 请求叛军信息
//! 消息: /get_rebel_find_info
type MSG_GetRebelFindInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetRebelFindInfo_Ack struct {
	RetCode int //! 返回码
	RebelID int //! 叛军ID
	Level   int //! 叛军等级
}

//玩家挑战挂机BOSS
//消息:/challenge_guaji_boss
type MSG_ChallenGuaJi_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
	Chapter    int    //章节
	CopyID     int    //副本ID
}

type MSG_ChallenGuaJi_Ack struct {
	RetCode    int   //! 返回码
	CurCopyID  int   //! 已通过关卡ID
	CurChapter int   //! 当前的章节
	NextTime   int64 //! 挑战Boss时间
	MoneyTime  int64 //! 货币挂机时间
	ExpTime    int64 //! 经验挂机时间
}

//! 玩家请求查询主线副本信息
//! 消息: /get_main_chapter_info
type MSG_GetMainChapterInfo_Req struct {
	PlayerID   int32    //! 玩家ID
	SessionKey string //! Session Key
}

type TMainCopy struct {
	CopyID      int //! 副本ID
	BattleTimes int //! 战斗次数
	ResetCount  int //! 当天刷新次数
	StarNum     int //! 星数
}

type TMainChapter struct {
	Chapter    int
	StarAward  [3]bool //! 6 12 15星 额外星级宝箱
	SceneAward [3]bool //! 1 3 5关卡 场景宝箱
}

type TMainInfo struct {
	CurCopyID  int //! 当前副本ID
	CurChapter int //! 当前章节ID

	CopyInfo []TMainCopy
	Chapter  []TMainChapter
}

type MSG_GetMainChapterInfo_Ack struct {
	RetCode int
	Info    TMainInfo
}

//! 玩家请求查询精英副本信息
//! 消息: /get_elite_chapter_info
type MSG_GetEliteChapterInfo_Req struct {
	PlayerID   int64
	SessionKey string
}

type TEliteCopy struct {
	CopyID      int //! 副本ID
	BattleTimes int //! 战斗次数
	ResetCount  int //! 当天刷新次数
	StarNum     int //! 星数
}

type TEliteChapter struct {
	Chapter    int
	StarAward  [3]bool //! 6 12 15星 额外星级宝箱
	SceneAward bool    //! 场景宝箱领取标记
	IsInvade   bool    //! 该章节是否遭遇入侵
}

type TEliteCopyData struct {
	CurCopyID  int             //! 当前副本ID
	CurChapter int             //! 当前章节ID
	Chapter    []TEliteChapter //! 章节信息
	CopyInfo   []TEliteCopy
}

type MSG_GetEliteChapterInfo_Ack struct {
	RetCode int
	Info    TEliteCopyData
}

//! 玩家请求入侵信息
//! 消息: /get_elite_invade_status
type MSG_GetEliteInvadeStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetEliteInvadeStatus_Ack struct {
	RetCode  int
	InvadeID []int
}

//! 玩家请求获取主线关卡星级奖励
//! 消息: /get_main_star_award
type MSG_GetMainStarAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
	StarAward  int
}

type MSG_GetMainStarAward_Ack struct {
	RetCode int //! 返回码
}

//! 玩家请求获取精英关卡星级奖励
//! 消息: /get_elite_star_award
type MSG_GetEliteStarAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
	StarAward  int
}

type MSG_GetEliteStarAward_Ack struct {
	RetCode int //! 返回码
}

//! 玩家请求获取主线关卡场景奖励
//! 消息: /get_main_scene_award
type MSG_GetMainSceneAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
	SceneAward int
}

type MSG_GetMainSceneAward_Ack struct {
	RetCode int
}

//! 玩家请求获取精英关卡场景奖励
//! 消息: /get_elite_scene_award
type MSG_GetEliteSceneAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
}

type MSG_GetEliteSceneAward_Ack struct {
	RetCode int
}

//! 玩家查询重置主线副本次数
//! 消息: /get_main_reset_times
type MSG_GetMainRefreshTimes_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
	CopyID     int
}

type MSG_GetMainRefreshTimes_Ack struct {
	RetCode      int
	RefreshTimes int //! 还能重置次数
}

//! 玩家查询重置精英副本次数
//! 消息: /get_elite_reset_times
type MSG_GetEliteRefreshTimes_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
	CopyID     int
}

type MSG_GetEliteRefreshTimes_Ack struct {
	RetCode      int
	RefreshTimes int //! 还能重置次数
}

//! 玩家请求重置主线副本挑战次数
//! 消息: /reset_main_battletimes
type MSG_ResetMainBattleTimes_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int //! 章节
	CopyID     int //! 副本ID
}

type MSG_ResetMainBattleTimes_Ack struct {
	RetCode  int
	MoneyID  int
	MoneyNum int
}

//! 玩家请求重置精英副本挑战次数
//! 消息: /reset_elite_battletimes
type MSG_ResetEliteBattleTimes_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int //! 章节
	CopyID     int //! 副本ID
}

type MSG_ResetEliteBattleTimes_Ack struct {
	RetCode  int
	MoneyID  int
	MoneyNum int
}

//! 玩家请求攻击入侵
//! 消息: /attack_elite_invade
type MSG_AttackEliteInvade_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
}

type MSG_AttackEliteInvade_Ack struct {
	RetCode   int
	Exp       int
	DropItems []MSG_ItemData
}

//! 玩家请求获取日常副本信息
//! 消息: /get_daily_copy
type MSG_GetDailyCopyInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_DailyCopy struct {
	ResType     int
	IsChallenge bool //! 是否挑战过了
}

type MSG_GetDailyCopyInfo_Ack struct {
	RetCode  int
	CopyInfo []MSG_DailyCopy
}

//! 玩家请求获取名将副本章节信息
//! 消息: /get_famous_chapter
type MSG_GetFamousCopyChapterInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetFamousCopyChapterInfo_Ack struct {
	RetCode     int
	CurCopyID   int
	CurChapter  int //! 当前章节
	BattleTimes int //! 总共挑战次数
}

//! 玩家请求获取名将副本详细信息
//! 消息: /get_famous_detail
type MSG_GetFamousCopyDetailInfo_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int //! 章节
}

type MSG_FamousCopyDetailInfo struct {
	CopyID      int
	BattleTimes int //! 每一个关卡的挑战次数
}

type MSG_GetFamousCopyDetailInfo_Ack struct {
	RetCode int
	CopyLst []MSG_FamousCopyDetailInfo
}

//! 玩家请求获取名将副本章节宝箱
//! 消息: /get_famous_award
type MSG_GetFamousCopyAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
}

type MSG_GetFamousCopyAward_Ack struct {
	RetCode int
}

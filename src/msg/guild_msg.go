package msg

//! 玩家请求公会状态
//! 消息: /get_guild_status
type MSG_GetGuildStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildStatus_Ack struct {
	RetCode int

	//! 副本
	ActionTimes        int                     //! 还能攻打次数
	NextRecoverTime    int64                   //! 下次恢复行动力时间戳
	CampLife           []MSG_CampLife          //! 当前攻打副本四个阵营血量
	IsBack             bool                    //! 是否回退章节
	CopyTreasure       []MSG_GuildCopyTreasure //! 章节宝藏领取情况
	PassChapter        int                     //! 当前攻打章节
	HistoryPassChapter int                     //! 历史通关最高章节
	AwardChapter       []MSG_PassAwardChapter  //! 已通关可领取奖励章节
	IsRecvCopyAward    []MSG_RecvCopyMark      //! 已领取的奖励信息

	//! 商店
	BuyLst []MSG_GuildGoods //! 已购买列表信息

	//! 祭天
	SacrificeStatus   int //! 祭天状态  0未祭天 1 2 3对应祭天模式
	SacrificeSchedule int //! 祭天进度
	SacrificeNum      int //! 当前祭天人数
	RecvLst           [4]int

	//! 技能
	SkillLst []TGuildSkill
}

//! 玩家请求创建公会
//! 消息: /create_guild
type MSG_CreateNewGuild_Req struct {
	PlayerID   int32
	SessionKey string
	Name       string
	Icon       int
}

type MSG_CreateNewGuild_Ack struct {
	RetCode  int
	NewGuild MSG_GuildInfo
}

//! 玩家查询公会状态
//! 消息: /get_guild
type MSG_GetGuildInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GuildInfo struct {
	GuildID   int    //! 军团ID
	Name      string //! 军团名字
	BossID    int32  //! 会长ID
	BossName  string //! 会长名字
	MemberNum int    //! 成员数量
	Icon      int    //! 军团Icon
	Notice    string //! 军团公告
	Level     int    //! 军团等级
	CurExp    int    //! 军团经验
}

type MSG_GetGuildInfo_Ack struct {
	RetCode     int
	IsHaveGuild bool            //! 是否有公会
	GuildLst    []MSG_GuildInfo //! 公会列表 如果有公会,则为公会信息, 若没公会,则为前五名公会数据
	CopyEndTime int64           //! 公会结束时间
}

//! 查看更多公会列表
//! 消息: /get_guild_lst
type MSG_GetGuildLst_Req struct {
	PlayerID   int32
	SessionKey string
	Index      int //! 从此值往后五个公会
}

type MSG_GetGuildLst_Ack struct {
	RetCode  int
	GuildLst []MSG_GuildInfo
}

//! 请求加入公会
//! 消息: /enter_guild
type MSG_EnterGuild_Req struct {
	PlayerID   int32
	SessionKey string
	GuildID    int
}

type MSG_EnterGuild_Ack struct {
	RetCode int
}

//! 请求查询已申请公会列表
//! 消息: /get_apply_guild_list
type MSG_GetApplyGuildList_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetApplyGuildList_Ack struct {
	RetCode  int
	GuildLst []MSG_GuildInfo
}

//! 撤销申请公会
//! /cancellation_guild_apply
type MSG_CancellationGuildApply_Req struct {
	PlayerID   int32
	SessionKey string
	GuildID    int
}

type MSG_CancellationGuildApply_Ack struct {
	RetCode int
}

//! 请求搜索公会
//! 消息: /search_guild
type MSG_SearchGuild_Req struct {
	PlayerID   int32
	SessionKey string
	GuildName  string
}

type MSG_SearchGuild_Ack struct {
	RetCode  int
	GuildLst []MSG_GuildInfo
}

//! 请求查询申请加入公会成员列表
//! 消息: /get_apply_guild_member_list
type MSG_GetApplyGuildMemberList_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_MemberInfo struct {
	PlayerID     int32
	Name         string
	Quality      int
	Level        int
	Role         int
	FightValue   int
	Contribution int
	OfflineTime  int64
	IsOnline     bool
}

type MSG_GetApplyGuildMemberList_Ack struct {
	RetCode       int
	MemberInfoLst []MSG_MemberInfo
}

//! 请求查询公会成员列表
//! 消息: /get_guild_member_list
type MSG_GetGuildMemberList_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildMemberList_Ack struct {
	RetCode   int
	MemberLst []MSG_MemberInfo
}

//! 接受成员入帮
//! 消息: /apply_through
type MSG_ApplyThrough_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int32
}

type MSG_ApplyThrough_Ack struct {
	RetCode int
}

//! 请求退出公会
//! 消息: /leave_guild
type MSG_LeaveGuild_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_LeaveGuild_Ack struct {
	RetCode int
}

//! 玩家查询公会祭天状态
//! 消息: /get_sacrifice_status
type MSG_GetSacrificeStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSacrificeStatus_Ack struct {
	RetCode           int
	SacrificeStatus   int //! 祭天状态  0未祭天 1 2 3对应祭天模式
	SacrificeSchedule int //! 祭天进度
	SacrificeNum      int //! 当前祭天人数
	RecvLst           [4]int
}

//! 玩家请求公会祭天
//! 消息: /guild_sacrifice
type MSG_GuildSacrifice_Req struct {
	PlayerID    int32
	SessionKey  string
	SacrificeID int //! 祭祀方式
}

type MSG_GuildSacrifice_Ack struct {
	RetCode           int
	MoneyID           int //! 获取货币
	MoneyNum          int
	CurExp            int
	GuildLevel        int
	SacrificeSchedule int //! 祭天进度
	SacrificeNum      int //! 当前祭天人数
}

//! 玩家请求领取祭天奖励
//! 消息: /get_sacrifice_award
type MSG_GetSacrificeAward_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_GetSacrificeAward_Ack struct {
	RetCode int
}

//! 玩家请求查询祭天奖励领取状态
//! 消息: /sacrifice_award_status
type MSG_GetSacrificeAwardStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSacrificeAwardStatus_Ack struct {
	RetCode int
}

//! 玩家请求查询公会商店商品信息
//! 消息: /query_guild_store
type MSG_QueryGuildStoreStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GuildGoods struct {
	ID    int //! 商品ID
	Times int //! 剩余次数/已购买次数
}

type MSG_QueryGuildStoreStatus_Ack struct {
	RetCode int
	BuyLst  []MSG_GuildGoods //! 已购买列表信息
}

//! 玩家请求购买军团商店物品
//! 消息: /buy_guild_store
type MSG_BuyGuildStoreItem_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
	Num        int
}

type MSG_BuyGuildStoreItem_Ack struct {
	RetCode int
}

//! 玩家请求攻击公会副本
//! 消息: /attack_guild_copy
type MSG_AttackGuildCopy_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
	CopyID     int   //! 攻击某个阵营 0 1 2 3
	Damage     int64 //! 造成伤害
}

type MSG_AttackGuildCopy_Ack struct {
	RetCode      int
	CampLife     []MSG_CampLife
	AwardChapter []MSG_PassAwardChapter
	IsPass       bool
	GuildLevel   int //! 当前公会等级
	CurExp       int //! 当前公会经验
}

//! 玩家查询工会副本状态
//! 消息: /get_guild_copy_status
type MSG_GetGuildCopyStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GuildCopyTreasure struct {
	CopyID     int
	Index      int
	AwardID    int
	PlayerName string
}

//! 查询玩家领取副本奖励情况
//! 消息: /query_recv_copy_award
type MSG_QueryGuildCopyTreasure_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
}

type MSG_QueryGuildCopyTreasure_Ack struct {
	RetCode      int
	CopyTreasure []MSG_GuildCopyTreasure
}

//! 通关阵营记录
type MSG_PassAwardChapter struct {
	PassChapter int
	CopyID      int
	PassTime    int64
	PlayerName  string //! 击杀者姓名
}

type MSG_CampLife struct {
	CopyID int
	Life   int64
}

type MSG_RecvCopyMark struct {
	Chapter int
	CopyID  int
}

type MSG_GetGuildCopyStatus_Ack struct {
	RetCode            int
	ActionTimes        int                     //! 还能攻打次数
	NextRecoverTime    int64                   //! 下次恢复行动力时间戳
	CampLife           []MSG_CampLife          //! 当前攻打副本四个阵营血量
	IsBack             bool                    //! 是否回退章节
	CopyTreasure       []MSG_GuildCopyTreasure //! 章节宝藏领取情况
	PassChapter        int                     //! 当前攻打章节
	HistoryPassChapter int                     //! 历史通关最高章节
	AwardChapter       []MSG_PassAwardChapter  //! 已通关可领取奖励章节
	IsRecvCopyAward    []MSG_RecvCopyMark      //! 已领取的奖励信息
}

//! 玩家请求领取副本通关奖励
//! 消息: /get_guild_copy_award
type MSG_GetGuildCopyAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int //! 章节
	CopyID     int //! 领取某个阵营的奖励
	ID         int //! 领取多少号箱子
}

type MSG_GetGuildCopyAward_Ack struct {
	RetCode int
	ItemID  int //! 奖励物品ID
	ItemNum int //! 奖励物品数量
	AwardID int
}

//! 玩家请求领取章节通关奖励
//! 消息: /get_guild_chapter_award
type MSG_GetGuildChapterAward_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
}

type MSG_GetGuildChapterAward_Ack struct {
	RetCode int
}

//! 玩家请求一键领取章节通关奖励
//! 消息: /get_guild_chapter_award_all
type MSG_GetGuildChapterAwardAll_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildChapterAwardAll_Ack struct {
	RetCode     int
	Award       []MSG_ItemData
	RecvChapter []int
}

//! 玩家请求查询章节通关领取状态
//! 消息: /get_guild_chapter_status
type MSG_GetGuildChapterAwardStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildChapterAwardStatus_Ack struct {
	RetCode int
	RecvLst []int
}

//! 玩家请求修改帮派信息
//! 消息: /update_guild_info
type MSG_UpdateGuildInfo_Req struct {
	PlayerID    int32
	SessionKey  string
	Icon        int    //! 图标
	Notice      string //! 公告
	Declaration string //! 宣言
}

type MSG_UpdateGuildInfo_Ack struct {
	RetCode int
}

//! 玩家请求修改帮派名字
//! 消息: /update_guild_name
type MSG_UpdateGuildName_Req struct {
	PlayerID   int32
	SessionKey string
	Name       string
}

type MSG_UpdateGuildName_Ack struct {
	RetCode int
}

//! 玩家请求设置公会副本回退状态
//! 消息: /update_guild_backstatus
type MSG_UpdateGuildBackStatus_Req struct {
	PlayerID   int32
	SessionKey string
	IsBack     int //! 0->不回退 1->回退到上一章
}

type MSG_UpdateGuildBackStatus_Ack struct {
	RetCode int
}

//! 玩家请求踢出帮派成员
//! 消息: /kick_member
type MSG_KickGuildMember_Req struct {
	PlayerID     int32
	SessionKey   string
	KickPlayerID int32
}

type MSG_KickGuildMember_Ack struct {
	RetCode int
}

//! 公会留言板留言
//! 消息: /write_guild_msg_board
type MSG_WriteGuildMsgBoard_Req struct {
	PlayerID   int32
	SessionKey string
	Message    string
}

type MSG_WriteGuildMsgBoard_Ack struct {
	RetCode int
}

//! 删除公会留言
//! 消息: /remove_guild_msg_board
type MSG_RemoveGuildMsgBoard_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int
	TargetTime     int64
}

type MSG_RemoveGuildMsgBoard_Ack struct {
	RetCode int
}

//! 查询公会留言板
//! 消息: /query_guild_msg_board
type MSG_QueryGuildMsgBoard_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GuildBoard struct {
	PlayerID   int32
	PlayerName string
	Message    string
	Time       int64
}

type MSG_QueryGuildMsgBoard_Ack struct {
	RetCode int
	MsgLst  []MSG_GuildBoard
}

//! 查询公会副本排行榜信息
//! 消息: /query_guild_copy_rank
type MSG_QueryGuildCopyRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GuildCopyRank struct {
	PlayerID    int32
	PlayerName  string
	Damage      int64
	BattleTimes int
}

type MSG_QueryGuildCopyRank_Ack struct {
	RetCode int
	RankLst []MSG_GuildCopyRank
}

//! 请求研究公会技能
//! 消息: /research_guild_skill
type MSG_ResearchGuildSkill_Req struct {
	PlayerID   int32
	SessionKey string
	SkillID    int
}

type MSG_ResearchGuildSkill_Ack struct {
	RetCode int
}

//! 请求学习公会技能
//! 消息: /study_guild_skill
type MSG_StudyGuildSkill_Req struct {
	PlayerID   int32
	SessionKey string
	SkillID    int
}

type MSG_StudyGuildSkill_Ack struct {
	RetCode int
}

type TGuildSkill struct {
	SkillID int
	Level   int
}

//! 请求公会技能等级
//! 消息: /get_guild_skill
type MSG_GetGuildSkillInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildSkillInfo_Ack struct {
	RetCode  int
	SkillLst []TGuildSkill
}

//! 请求公会技能研发等级
//! 消息: /get_guild_skill_limit
type MSG_GetGuildSkillResearch_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildSkillResearch_Ack struct {
	RetCode  int
	SkillLst []TGuildSkill
}

//! 查询公会状态
//! 消息: /get_guild_log
type MSG_GetGuildLog_Req struct {
	PlayerID   int32
	SessionKey string
}

type GuildEvent struct {
	PlayerID   int32
	PlayerName string
	Type       int //! 祭天->类型
	Value      int //! 祭天->经验  升级->等级  职位->新职位
	Action     int //!
	Time       int64
}

type MSG_GetGuildLog_Ack struct {
	RetCode int
	LogLst  []GuildEvent
}

//! 修改公会职位
//! 消息: /change_guild_role
type MSG_ChangeGuildRole_Req struct {
	PlayerID       int32
	SessionKey     string
	TargetPlayerID int32
	Role           int
}

type MSG_ChangeGuildRole_Ack struct {
	RetCode int
}

//! 升级工会
//! 消息: /guild_levelup
type MSG_GuildLevelUp_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GuildLevelUp_Ack struct {
	RetCode int
}

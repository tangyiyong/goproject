package msg

type MSG_FriendInfo struct {
	PlayerID   int32
	Name       string //好友的名字
	HeroID     int    //英雄ID
	Quality    int8   //英雄品质
	GuildName  string //军团名字
	FightValue int32  //战力
	Level      int    //等级
	OffTime    int32  //离线时间 >0表示离线，==0表示在线
	IsGive     int    //0:表示未赠送, 1:表示己赠送
	HasAct     int    //0:表没有未领取， 1:表示有未领取
}

//! 玩家请求竞技场名次排行榜信息
//! 消息: /get_all_friend
type MSG_GetAllFriend_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetAllFriend_Ack struct {
	RetCode   int
	RcvNum    int              //己领个数
	FriendLst []MSG_FriendInfo //好友列表
}

//! 玩家请求添加好友
//! 消息: /add_friend_request
type MSG_AddFriend_Req struct {
	PlayerID   int32
	SessionKey string
	TargetID   int32  //目标玩家ID
	TargetName string //目标玩家名
}

type MSG_AddFriend_Ack struct {
	RetCode int
}

//! 玩家请求删除好友
//! 消息: /del_friend_request
type MSG_DelFriend_Req struct {
	PlayerID   int32
	SessionKey string
	TargetID   int32 //目标玩家ID
}

type MSG_DelFriend_Ack struct {
	RetCode int
}

//! 玩家回应请求添加好友
//! 消息: /process_friend_request
type MSG_ProcessFriend_Req struct {
	PlayerID   int32
	SessionKey string
	TargetID   int32 //目标玩家ID TargetID: 0 表示处理全部, 否则表示处理一个
	IsAgree    int   // 0: 表示拒绝 ，1 表示同意
}

type MSG_ProcessFriend_Ack struct {
	RetCode   int
	RcvNum    int              //己领个数
	FriendLst []MSG_FriendInfo //新的好友列表, 只有批量同意的情况下，才有效
}

//! 玩家请求申请添加好友列表
//! 消息: /get_apply_list
type MSG_GetApplyList_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetApplyList_Ack struct {
	RetCode   int
	FriendLst []MSG_FriendInfo //好友申请列表
}

//! 玩家增送行动力
//! 消息: /give_action
type MSG_GiveAction_Req struct {
	PlayerID   int32
	SessionKey string
	TargetID   int32 //目标玩家
}

type MSG_GiveAction_Ack struct {
	RetCode   int
	RcvNum    int              //己领个数
	FriendLst []MSG_FriendInfo //好友列表
}

//! 玩家收取行动力
//! 消息: /receive_action
type MSG_ReceiveAction_Req struct {
	PlayerID   int32
	SessionKey string
	TargetID   int32 //目标玩家， 0:表示全部收取， 有值: 表示只收取一个人
}

type MSG_ReceiveAction_Ack struct {
	RetCode     int
	ActionValue int              //! 行动力值
	ActionTime  int32            //! 行动力恢复起始时间
	RcvNum      int              //己领个数
	FriendLst   []MSG_FriendInfo //好友列表
}

//! 玩家搜索好友
//! 消息: /search_friend
type MSG_SearchFriend_Req struct {
	PlayerID   int32
	SessionKey string
	Name       string //搜索名字
}

type MSG_SearchFriend_Ack struct {
	RetCode   int
	FriendLst []MSG_FriendInfo //好友列表
}

//! 玩家请求推荐好友
//! 消息: /recomand_friend
type MSG_RecomandFriend_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_RecomandFriend_Ack struct {
	RetCode   int
	FriendLst []MSG_FriendInfo //好友列表
}

type MSG_OnlineInfo struct {
	PlayerID   int32
	Name       string //好友的名字
	HeroID     int    //英雄ID
	Quality    int8   //英雄品质
	GuildName  string //军团名字
	FightValue int32  //战力
	Level      int    //等级
}

//! 玩家请求竞技场名次排行榜信息
//! 消息: /get_online_friend
type MSG_GetOnlineFriend_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetOnlineFriend_Ack struct {
	RetCode   int
	OnlineLst []MSG_OnlineInfo //好友列表
}

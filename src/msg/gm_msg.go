package msg

//! GM请求更新配制文件
//! 消息: /update_gamedata
type MSG_UpdateGameData_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	TbName     string //表名
}

type MSG_UpdateGameData_Ack struct {
	RetCode int //返回码
}

//! GM增发全服奖励
//! 消息: /add_svr_award
type MSG_SvrAward_Add_Req struct {
	SessionID  string         //GM SessionID
	SessionKey string         //GM SessionKey
	Value      []string       //! 参数
	ItemLst    []MSG_ItemData //! 奖励内容
}

type MSG_SvrAward_Add_Ack struct {
	RetCode int
}

//! GM删除全服奖励
//! 消息: /del_svr_award
type MSG_SvrAward_Del_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	ID         int
}

type MSG_SvrAward_Del_Ack struct {
	RetCode int
}

//! GM发个人奖励
//! 消息: /send_award_to_player
type MSG_Send_Award_Player_Req struct {
	SessionID  string         //GM SessionID
	SessionKey string         //GM SessionKey
	TargetID   int32          //目标玩家
	Value      string         //参数
	ItemLst    []MSG_ItemData //奖励内容
}

type MSG_Send_Award_Player_Ack struct {
	RetCode int
}

//! 查看当前服务器状态
//! 消息: /get_server_info
type MSG_GetServerInfo_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey

}

type MSG_GetServerInfo_Ack struct {
	SvrID        int32  //当前的服务器ID
	SvrName      string //当前服务器名字
	OnlineCnt    int    //在线人数
	MaxOnlineCnt int    //总人数
	RegisterCnt  int    //总注册人数
}

//验证玩家登录请求
//消息:/set_gamesvr_flag
type MSG_SetGameSvrFlag_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
	SvrID      int32  //服务器ID
	Flag       uint32 //服务器标记
}

type MSG_SetGameSvrFlag_Ack struct {
	RetCode int //返回码 0:成功 1: 失
}

//请求服务器列表
//消息:/get_server_list
type MSG_GetServerList_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
}

type MSG_GetServerList_Ack struct {
	RetCode int
	SvrList []ServerNode //服务器结点表
}

//gm用户登录
//消息:/gm_login
type MSG_GmLogin_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey
}

type MSG_GmLogin_Ack struct {
	RetCode int
}

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
	TargetID   int            //目标玩家
	Value      string         //参数
	ItemLst    []MSG_ItemData //奖励内容
}

type MSG_Send_Award_Player_Ack struct {
	RetCode int
}

//! 查看当前服务器状态
//! 消息: /server_state_info
type MSG_ServerStateInfo_Req struct {
	SessionID  string //GM SessionID
	SessionKey string //GM SessionKey

}

type MSG_ServerStateInfo_Ack struct {
	SvrName     string //当前服务器名字
	OnlineCount int    //在线人数
	TotalCount  int    //总人数
	MemAlloc    uint64
	MemInuse    uint64
	MenObjNum   uint64
}

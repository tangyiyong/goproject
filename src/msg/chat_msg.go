package msg

const (
	MSG_CHANNEL_INVALID = 0 //无效的频道
	MSG_CHANNEL_PLAYER  = 1 //私聊频道
	MSG_CHANNEL_GUILD   = 2 //公会频道
	MSG_CHANNEL_WORLD   = 3 //世界频道
)

//MSG_CHECK_IN_REQ
type MSG_CheckIn_Req struct {
	PlayerName string
	PlayerID   int32
	GuildID    int
}

//MSG_CHATMSG_REQ
type MSG_Chat_Req struct {
	SourceName     string
	TargetChannel  int
	TargetGuildID  int
	TargetName     string
	TargetPlayerID int
	MsgContent     string
	HeroID         int
	Quality        int
}

//MSG_CHATMSG_ACK
type MSG_Chat_Ack struct {
	RetCode int
}

//聊天服向客户端发的消息
//MSG_CHATMSG_NOTIFY
type MSG_Chat_Msg_Notify struct {
	SourcePlayerID int32
	SourceName     string
	TargetChannel  int
	TargetGuildID  int
	MsgContent     string
	HeroID         int
	Quality        int
}

type MSG_CheckIn_Ack struct {
	PlayerName string
	PlayerID   int32
}

//Game server 发送聊天服的通知消息
type MSG_GameSvr_Nofity struct {
	PlayerName string
	PlayerID   int32
}

//玩家公会变化通知消息
//MSG_GUILD_NOTIFY
type MSG_GuildNotify_Req struct {
	PlayerID   int32
	NewGuildID int //新的公会ID
}

//玩家上下线变化通知消息
//MSG_ONLINE_NOTIFY
type MSG_OnlineNotify_Req struct {
	PlayerID int32
	Online   bool //true 上线， false 下线
}

//跑马灯通知消息
//MSG_HORSELAME_NOTIFY
type MSG_HorseLame_Notify struct {
	TextType int      //文本索引
	Params   []string //参数
	Camps    []int    //阵营战专用
}

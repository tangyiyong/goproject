package msg

const (
	MSG_BEGIN      = 0
	MSG_CONNECT    = 1 //连接成功
	MSG_DISCONNECT = 2 //断开连接

	//客户端发给聊天服的消息
	MSG_CHECK_IN_REQ = 3 //进入聊天服
	MSG_CHECK_IN_ACK = 4 //进入聊天服回复
	MSG_CHATMSG_REQ  = 5 //聊天聊天请求

	MSG_GAME_TO_CLIENT = 6 //游戏服到客户端的中转消息

	//聊天服发给客户端的消息
	MSG_CHATMSG_ACK    = 7 //客户端聊天的回复， 一般有错误才回复， 目标不存在，被禁言， 发送错误.....
	MSG_CHATMSG_NOTIFY = 8 //聊天服发给客户端的聊天消息， 个人，公会，世界

	//游戏服发给聊天服的消息
	MSG_GUILD_NOTIFY = 9 //玩家公会变化通知消息

	//聊天服发给游戏服的消息
	MSG_ONLINE_NOTIFY = 10 //玩家上下线通知

	//跑马灯通知消息
	MSG_HORSELAME_NOTIFY = 11
	//心跳消息
	MSG_HEART_BEAT = 12 //心跳消息

	MSG_ENTER_ROOM_REQ    = 13 //角色进入阵营房间(Client To BatSvr)
	MSG_ENTER_ROOM_ACK    = 14 //角色进入阵营的回复(回复给玩家自己)(BatSvr To Client)
	MSG_ENTER_ROOM_NTY    = 15 //收到其它玩家进入战场消息(BatSvr To Client)
	MSG_LEAVE_ROOM_REQ    = 16 //角色离开阵营房间
	MSG_LEAVE_ROOM_ACK    = 17 //角色离开阵营房间回复(BatSvr To Client)
	MSG_LEAVE_ROOM_NTY    = 18 //收到其它玩家离开战场消息(BatSvr To Client)
	MSG_MOVE_STATE        = 19 //角色移动信息
	MSG_SKILL_STATE       = 20 //角色技能信息
	MSG_BUFF_STATE        = 21 //角色BUFF信息
	MSG_HERO_STATE        = 22 //角色英雄状态信息
	MSG_START_CARRY_REQ   = 23 //开始搬运水晶请求
	MSG_START_CARRY_ACK   = 24 //开始搬运水晶回复
	MSG_FINISH_CARRY_REQ  = 25 //玩家完成搬运水晶请求
	MSG_FINISH_CARRY_ACK  = 26 //玩家完成搬运水晶回复
	MSG_NEW_SKILL_NTY     = 27 //服务器更新技能通知
	MSG_REVIVE_NTY        = 28 //玩家复活的通知消息
	MSG_PLAYER_QUERY_REQ  = 29 //玩家查询水晶品质
	MSG_PLAYER_QUERY_ACK  = 30 //玩家查询水晶品质回复
	MSG_PLAYER_CHANGE_REQ = 31 //玩家请求置换水晶
	MSG_PLAYER_CHANGE_ACK = 32 //玩家请求置换水晶回复
	MSG_PLAYER_REVIVE_REQ = 33 //玩家请求复活
	MSG_PLAYER_REVIVE_ACK = 34 //玩家请求复活回复
	MSG_KILL_EVENT_REQ    = 35 //玩家杀伤事件请求
	MSG_KILL_EVENT_ACK    = 36 //玩家杀伤事件回复
	MSG_LOAD_CAMPBAT_REQ  = 37 //请求加载战斗数据
	MSG_LOAD_CAMPBAT_ACK  = 38 //请求加载战斗数据回复
	MSG_ALL_DIE_NTY       = 39 //玩家英雄全部死亡通知
	MSG_CAMPBAT_CHAT_REQ  = 40 //阵营战聊天请求
	MSG_CAMPBAT_CHAT_ACK  = 41 //阵营战聊天回复

	MSG_SVR_LOGDATA = 99 //GameSvr 发到 LogSvr的日志数据

	//
	MSG_END = 100
)

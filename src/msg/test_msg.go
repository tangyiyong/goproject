package msg

//玩家获取货币
//消息:/test_get_money
type MSG_GetTestMoney_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_GetTestMoney_Ack struct {
	RetCode int     //返回码
	Moneys  [10]int //货币表
}

//玩家获取货币
//消息:/test_get_action
type MSG_GetTestAction_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_GetTestAction_Ack struct {
	RetCode int   //返回码
	Actions []int //行动力
}

//玩家升级主角等级
//消息:/test_uplevel
type MSG_TestUpLevel_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_TestUpLevel_Ack struct {
	RetCode  int //返回码
	RetLevel int //返回等级
}

//玩家升级10级
//消息:/test_uplevel_ten
type MSG_TestUpLevelTen_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_TestUpLevelTen_Ack struct {
	RetCode  int //返回码
	RetLevel int //返回等级
}

//玩家增加VIP经验
//消息:/test_add_vip
type MSG_TestAddVip_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_TestAddVip_Ack struct {
	RetCode  int //返回码
	VipLevel int //返回当前等级
	VipExp   int //返回当前的VIP经验
}

//玩家增加工会经验
//消息:/test_add_guild
type MSG_TestAddGuild_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_TestAddGuild_Ack struct {
	RetCode    int //返回码
	GuildLevel int //返回当前等级
	GuildExp   int //返回当前的经验
}

//测试玩家的各项属性
//消息:/test_heros_property
type MSG_TestHerosProperty_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
}

type MSG_TestObject struct {
	HeroID          int
	PropertyValue   [11]int
	PropertyPercent [11]int
	CampDef         [5]int
	CampKill        [5]int
}

type MSG_TestHerosProperty_Ack struct {
	RetCode    int               //返回码
	BattleCamp int               //角色阵营
	PlayerID   int32               //角色ID
	Level      int               //主角等级
	Heros      [6]MSG_TestObject //英雄对象
}

// 添加物品
type MSG_TestAddItem_Req struct { // 消息：/test_add_item
	PlayerID   int32
	SessionKey string
	ItemID     int
	AddNum     int
}
type MSG_TestAddItem_Ack struct {
	RetCode int
	Count   int
}

//玩家测试充值相关活动
//消息:/test_charge_money
type MSG_ChargeTestMoney_Req struct {
	PlayerID   int32    //玩家ID
	SessionKey string //Sessionkey
	RMB        int    //! 充值人民币
	ChargeID   int    //充值ID
}

type MSG_ChargeTestMoney_Ack struct {
	RetCode  int //返回码
	VIPExp   int
	VIPLevel int
}

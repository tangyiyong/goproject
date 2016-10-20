package msg

type Target_Hero struct {
	HeroID  int //英雄ID
	PosType int //位置类型， 1：上阵， 2 援军 3:背包
	HeroPos int //英雄索引
}

type Cost_Hero struct {
	HeroID  int //英雄ID
	HeroPos int //英雄索引
}

//玩家升级英雄
//消息:/upgrade_hero
type MSG_UpgradeHero_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄
	CostHeros  []Cost_Hero //英雄列表
}

type MSG_UpgradeHero_Ack struct {
	RetCode    int   //返回码
	HeroID     int   //英雄ID
	NewLevel   int   //英雄等级
	NewExp     int   //英雄经验
	CostMoney  int   //花费货币值
	FightValue int32 //战力
}

//玩家更换英雄
//消息:/change_hero
type MSG_ChangeHero_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetID   int    //上阵英雄ID
	TargetPos  int    //上阵位置
	SourcePos  int    //背包位置
	SourceID   int    //背包英雄ID
}

type MSG_ChangeHero_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

//玩家更换援军英雄
//消息:/change_back_hero
type MSG_ChangeBackHero_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetID   int    //上阵英雄ID
	TargetPos  int    //上阵位置
	SourcePos  int    //背包位置
	SourceID   int    //背包英雄ID
}

type MSG_ChangeBackHero_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

//玩家取下援军英雄
//消息:/unset_backhero
type MSG_UnsetBackHero_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	HeroID     int    //目标英雄ID
	HeroPos    int    //目标位置
}

type MSG_UnsetBackHero_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

//玩家英雄突破
//消息:/breakout_hero
type MSG_BreakOut_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄位置
	CostHeros  []Cost_Hero //英雄列表
}

type MSG_BreakOut_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力值
	NewLevel   int8  //新的突破等级
	CostItems  int   //消耗的材料数
	CostMoney  int   //消耗的货币数
}

//玩家培养英雄
//消息:/culture_hero
type MSG_CultureHero_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄位置
	Times      int         //培养次数
}

type MSG_CultureHero_Ack struct {
	RetCode    int    //返回码
	Cultures   [5]int //培养新值
	FightValue int32  //战力
	CostItems  int    //消耗道具
}

//玩家天命英雄
//消息:/destiny_hero
type MSG_DestinyHero_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄位置
}

type MSG_DestinyHero_Ack struct {
	RetCode         int    //返回码
	NewDestinyState uint32 //新的天命状态
	FightValue      int32  //新的战力
	CostItemNum     int    //消耗数量
}

//玩家查询天命信息
//消息:/query_destiny
type MSG_QueryDestinyState_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄位置
}

type MSG_QueryDestinyState_Ack struct {
	RetCode         int    //返回码
	NewDestinyState uint32 //新的天命状态
	FightValue      int32  //新的战力
}

//玩家合成英雄
//消息:/compose_hero
type MSG_ComposeHero_Req struct {
	PlayerID    int32  //玩家ID
	SessionKey  string //Sessionkey
	HeroPieceID int    //英雄碎片ID
}

type MSG_ComposeHero_Ack struct {
	RetCode int //返回码
	HeroID  int //合成的英雄ID
}

//客户端通知服务器玩家升级
//消息:/levelup_notify
type MSG_LevelUpNotify_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_LevelUpNotify_Ack struct {
	RetCode    int   //返回码
	Level      int   //新的等级
	CurSvrTime int32 //当前的服务器时间
	CurExp     int   //当前的经验值
	FightValue int32 //战力
}

//玩家设置觉醒道具
//消息:/set_wake_item
type MSG_SetWakeItem_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄
	Pos        int         //目标位置
	ID         int         //觉醒道具ID

}

type MSG_SetWakeItem_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
	Pos        int   //目标位置
	ID         int   //觉醒道具ID
}

//玩家觉醒等级
//消息:/up_wake_level
type MSG_UpWakeLevel_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄
	SourcePos  int         //消耗英雄的位置
}

type MSG_UpWakeLevel_Ack struct {
	RetCode    int   //返回码
	WakeLevel  int   //新的觉醒等级
	FightValue int32 //新的战力
}

//玩家合成觉醒道具
//消息:/compose_wake_item
type MSG_ComposeWakeItem_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	ItemID     int    //觉醒道具ID
}

type MSG_ComposeWakeItem_Ack struct {
	RetCode int //返回码
	ItemID  int //觉醒道具ID
}

//玩家升品雕文
//消息:/upgrade_diaowen
type MSG_UpgradeDiaoWen_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄
	DiaoWenID  int         //雕文ID
}

type MSG_UpgradeDiaoWen_Ack struct {
	RetCode        int   //返回码
	DiaoWenID      int   //雕文ID
	DiaoWenQuality int32 //雕文品质
	FightValue     int32 //新的战力
}

//玩家洗炼雕文
//消息:/xilian_diaowen
type MSG_XiLianDiaoWen_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄
	DiaoWenID  int         //雕文ID
	LockIndex  [4]int      //锁定的功能索引
}

type MSG_XiLianDiaoWen_Ack struct {
	RetCode      int      //返回码
	DiaoWenID    int      //雕文ID
	RandValue    [5]int32 //四个随机值 生命, 攻击
	CostMoneyID  int      //消耗的货币ID
	CostMoneyNum int      //消耗的货币值
}

//玩家洗炼替换雕文
//消息:/xilian_tihuan
type MSG_XiLianTiHuan_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	TargetHero Target_Hero //目标英雄
	DiaoWenID  int         //雕文ID
}

type MSG_XiLianTiHuan_Ack struct {
	RetCode       int      //返回码
	DiaoWenID     int      //雕文ID
	FightValue    int32    //新的战力
	PropertyValue [5]int32 //五个属性值
}

//玩家升级宠物
//消息:/upgrade_pet
type MSG_UpgradePet_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:护佑  2:背包中
	PosIndex   int    //位置索引
	PetID      int    //装备ID
	ItemID     []int  //道具ID
	ItemNum    []int  //道具数量
}

type MSG_UpgradePet_Ack struct {
	RetCode    int   //返回码
	NewLevel   int   //新等级
	NewExp     int   //新的经验
	CostMoney  int   //消耗的货币
	FightValue int32 //战力
}

//玩家升星宠物
//消息:/upstar_pet
type MSG_UpstarPet_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:护佑  2:背包中
	PosIndex   int    //位置索引
	PetID      int    //装备ID
}

type MSG_UpstarPet_Ack struct {
	RetCode    int   //返回码
	NewStar    int   //新的星级
	CostMoney  int   //消耗的货币
	FightValue int32 //战力
}

//玩家神炼宠物
//消息:/upgod_pet
type MSG_UpgodPet_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:护佑  2:背包中
	PosIndex   int    //位置索引
	PetID      int    //装备ID
	ItemID     int    //道具ID
}

type MSG_UpgodPet_Ack struct {
	RetCode    int   //返回码
	Exp        int   //神炼经验
	Level      int   //神炼等级
	FightValue int32 //战力
}

//玩家更换宠物
//消息:/change_pet
type MSG_ChangePet_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetID   int    //上阵宠物ID
	TargetPos  int    //上阵位置
	SourcePos  int    //背包位置
	SourceID   int    //背包宠物ID
}

type MSG_ChangePet_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

//玩家下阵宠物
//消息:/unset_pet
type MSG_UnsetPet_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetID   int    //上阵宠物ID
	TargetPos  int    //上阵位置
}

type MSG_UnsetPet_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

//玩家合成宠物
//消息:/compose_pet
type MSG_ComposePet_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PetPieceID int    //宠物碎片ID
}

type MSG_ComposePet_Ack struct {
	RetCode int //返回码
	PetID   int //合成的宠物ID
}

//玩家化神英雄
//消息:/upgod_hero
type MSG_UpgodHero_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:上阵， 2:援军， 3:背包
	PosIndex   int    //位置索引
	HeroID     int    //英雄ID
}

type MSG_UpgodHero_Ack struct {
	RetCode    int   //返回码
	GodLevel   int   //化神等级
	Quality    int8  //品质
	FightValue int32 //战力
}

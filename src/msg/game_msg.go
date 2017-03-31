package msg

//玩家登录游戏服务器
//消息:/user_login_game
type MSG_LoginGameSvr_Req struct {
	AccountID int32  //账号ID
	LoginKey  string //登录key
}

//玩家登录游戏服务器返回
type MSG_LoginGameSvr_Ack struct {
	RetCode    int    // 0 表示登录成功   1 无效的登录key ,玩家需要从新登录
	PlayerID   int32  //玩家角色ID
	SessionKey string //SessionKey
}

//玩家创建角色请求
//消息:/create_new_player
type MSG_CreateNewPlayerReq struct {
	AccountID  int32  //账号ID
	SessionKey string //玩家角色ID
	PlayerName string //玩家角色名
	HeroID     int    //英雄ID
	ChannelID  int32  //渠道ID
}

type MSG_CreateNewPlayerAck struct {
	RetCode  int // 0 表示创建成功   1 表示无效的SessionKey, 2 表示角色不合法,3 表示创建角色失败
	PlayerID int32
}

//玩家进入游戏服
//消息:/user_enter_game
type MSG_EnterGameSvr_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_EnterGameSvr_Ack struct {
	RetCode     int    //返回码
	XorCode     int    //异或字节
	GuildID     int32  //工会ID
	SvrTime     int32  //服务器时间
	ChatSvrAddr string //聊天服的地址
	PlayerName  string //玩家角色名
	FightValue  int32  //战力
}

//玩家获取角色数据
//消息:/get_role_data
type MSG_GetRoleData_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetRoleData_Ack struct {
	RetCode    int     //返回码
	VipLevel   int8    //VIP等级
	VipExp     int     //VIP经验
	BatCamp    int8    //阵营战阵营
	Actions    []int   //行动力表
	ActionTime []int32 //行动力时间
	Moneys     []int   //货币表
	NewWizard  string  //新手向导信息
}

//玩家离开游戏服
//消息:/user_leave_game
type MSG_LeaveGameSvr_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_LeaveGameSvr_Ack struct {
	RetCode int
}

type MSG_QueryServerTime_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_QueryServerTime_Ack struct {
	RetCode int
	SvrTime int32 //服务器时间
}

//玩家更换装备
//消息:/change_equip
type MSG_ChangeEquip_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetPos  int    //目标位置索引
	TargetID   int    //目标位置装备ID
	SourcePos  int    //背包位置索引
	SourceID   int    //背包位置装备ID
}

type MSG_ChangeEquip_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

//玩家强化装备
//消息:/equip_strengthen
type MSG_EquipStrengthen_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:上阵  2:背包中
	PosIndex   int    //位置索引
	EquipID    int    //装备ID
	Times      int    //强化次数
}

type MSG_EquipStrengthen_Ack struct {
	RetCode    int   //返回码
	BaoJi      int   //1 : 暴击
	NewLevel   int   //新的强化等级
	FightValue int32 //新的战力
	CostMoney  int   //花掉的钱
}

//玩家精炼装备
//消息:/equip_refine
type MSG_EquipRefine_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:上阵  3:背包中
	PosIndex   int    //位置索引
	EquipID    int    //装备ID
	ItemID     int    //道具ID
	ItemNum    int    //道具个数
}

type MSG_EquipRefine_Ack struct {
	RetCode    int   //返回码
	Exp        int   //新的经验
	Level      int   //新的精炼等级
	PosType    int   //位置类型 1:上阵  3:背包中
	PosIndex   int   //位置索引
	ItemID     int   //道具ID
	FightValue int32 //新的战力
}

//玩家升星装备
//消息:/equip_risestar
type MSG_EquipRiseStar_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	PosType    int    //位置类型 1:上阵  3:背包中
	PosIndex   int    //位置索引
	EquipID    int    //装备ID
	CondIndex  int    //升星条件索引
}

type MSG_EquipRiseStar_Ack struct {
	RetCode    int   //返回码
	Exp        int   //新的经验
	Luck       int   //幸运值
	Level      int   //新的升星等级
	FightValue int32 //新的战力
}

//玩家合成装备
//消息:/compose_equip
type MSG_ComposeEquip_Req struct {
	PlayerID     int32  //玩家ID
	SessionKey   string //Sessionkey
	EquipPieceID int    //装备碎片ID
}

type MSG_ComposeEquip_Ack struct {
	RetCode int //返回码
	EquipID int //合成的装备ID
}

//使用物品
//消息:/use_item
type MSG_UseItem_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	ItemID     int    //物品ID
	ItemNum    int    //物品数量
	Index      int    //选择索引, 选择性的背包
}

type MSG_UseItem_Ack struct {
	RetCode int            //返回码
	Items   []MSG_ItemData //物器列表
}

type MSG_SellItem struct {
	Pos int //位置
	ID  int //ID
}

//使用物品
//消息:/sell_item
type MSG_SellItem_Req struct {
	PlayerID   int32          //玩家ID
	SessionKey string         //Sessionkey
	ItemType   int            //物品类型
	Items      []MSG_SellItem //物器列表
}

type MSG_SellItem_Ack struct {
	RetCode  int //返回码
	MoneyID  int //货币ID
	MoneyNum int //货币值
}

//玩家更换宝物
//消息:/change_gem
type MSG_ChangeGem_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetPos  int    //目标位置索引
	TargetID   int    //目标位置装备ID
	SourcePos  int    //背包位置索引
	SourceID   int    //背包位置装备ID
}

type MSG_ChangeGem_Ack struct {
	RetCode    int   //返回码
	FightValue int32 //战力
}

type Cost_Gem struct {
	GemID  int //宝物ID
	GemPos int //宝物索引
}

//玩家强化宝物
//消息:/gem_strengthen
type MSG_GemStrengthen_Req struct {
	PlayerID   int32      //玩家ID
	SessionKey string     //Sessionkey
	GemPosType int        //位置类型 1:上阵  3:背包中
	GemIndex   int        //位置索引
	GemID      int        //宝物ID
	CostGems   []Cost_Gem //宝物列表
}

type MSG_GemStrengthen_Ack struct {
	RetCode      int   //返回码
	Level        int   //新的强化等级
	Exp          int   //经验
	NewPos       int   //新位置
	FightValue   int32 //战力
	CostMoneyID  int   //消耗货币ID
	CostMoneyNum int   //消耗货币数
}

//玩家精炼宝物
//消息:/gem_refine
type MSG_GemRefine_Req struct {
	PlayerID   int32      //玩家ID
	SessionKey string     //Sessionkey
	GemPosType int        //位置类型 1:上阵  2:援军中 3:背包中
	GemIndex   int        //位置索引
	GemID      int        //装备ID
	CostGems   []Cost_Gem //宝物列表
}

type MSG_GemRefine_Ack struct {
	RetCode      int   //返回码
	Level        int   //新的精炼等级
	FightValue   int32 //战力
	CostMoneyID  int   //消耗货币ID
	CostMoneyNum int   //消耗货币数
}

//玩家更改职业
//消息:/change_career
type MSG_ChangeCareer_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_ChangeCareer_Ack struct {
	RetCode   int //返回码
	NewHeroID int //新的英雄ID
}

//玩家查询英雄培养与天命消耗
//消息: /query_hero_decompose_cost
type MSG_QueryHeroDecomposeCost_Req struct {
	PlayerID   int32
	SessionKey string
	CostHeros  []Cost_Hero
}

type MSG_QueryHeroDecomposeCost_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

//玩家分解英雄
//消息:/decompose_hero
type MSG_DecomposeHero_Req struct {
	PlayerID   int32       //玩家ID
	SessionKey string      //Sessionkey
	CostHeros  []Cost_Hero //分解英雄列表
}

type MSG_DecomposeHero_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //分解所获物品
}

//! 玩家查询分解宠物所得
//! 消息: /query_pet_decompose_cost
type MSG_QueryPetDecomposeCost_Req struct {
	PlayerID   int32
	SessionKey string
	PetID      int
	PetPos     int
}

type MSG_QueryPetDecomposeCost_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

//! 玩家分解宠物
//! 消息: /decompose_pet
type MSG_DecomposePet_Req struct {
	PlayerID   int32
	SessionKey string
	PetID      int
	PetPos     int
}

type MSG_DecomposePet_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

//! 玩家重生宠物
//! 消息: /relive_pet
type MSG_RelivePet_Req struct {
	PlayerID   int32
	SessionKey string
	PetID      int
	PetPos     int
}

type MSG_RelivePet_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

type Cost_Equip struct {
	EquipID  int //英雄ID
	EquipPos int //英雄索引
}

//! 玩家查询分解装备所得
//! 消息: /query_equip_decompose_cost
type MSG_QueryEquipDecomposeCost_Req struct {
	PlayerID   int32
	SessionKey string
	CostEquips []Cost_Equip
}

type MSG_QueryEquipDecomposeCost_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData //分解所获物品
}

//! 玩家查询重生战宠所得
//! 消息: /query_pet_relive_cost
type MSG_QueryPetReliveCost_Req struct {
	PlayerID   int32
	SessionKey string
	PetID      int
	PetPos     int
}

type MSG_QueryPetReliveCost_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

//! 玩家查询重生英雄所得
//! 消息: /query_hero_relive_cost
type MSG_QueryHeroReliveCost_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	HeroID     int    //英雄ID
	HeroPos    int    //英雄位置
}

type MSG_QueryHeroReliveCost_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //重生所获物品
}

//! 玩家查询重生装备所得
//! 消息: /query_equip_relive_cost
type MSG_QueryEquipReliveCost_Req struct {
	PlayerID     int32
	SessionKey   string
	CostEquipID  int
	CostEquipNum int
}

type MSG_QueryEquipReliveCost_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData //分解所获物品
}

//! 玩家查询重生宝物所得
//! 消息: /query_gem_relive_cost
type MSG_QueryGemReliveCost_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	GemID      int    //宝物ID
	GemPos     int    //宝物位置
}

type MSG_QueryGemReliveCost_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //重生所获物品
}

//玩家分解装备
//消息:/decompose_equip
type MSG_DecomposeEquip_Req struct {
	PlayerID   int32        //玩家ID
	SessionKey string       //Sessionkey
	CostEquips []Cost_Equip //分解英雄列表
}

type MSG_DecomposeEquip_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //分解所获物品
}

//玩家重生英雄
//消息:/relive_hero
type MSG_ReliveHero_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	HeroID     int    //英雄ID
	HeroPos    int    //英雄位置
}

type MSG_ReliveHero_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //重生所获物品
}

//玩家重生装备
//消息:/relive_equip
type MSG_ReliveEquip_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	EquipID    int    //装备ID
	EquipPos   int    //装备位置
}

type MSG_ReliveEquip_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //分解所获物品
}

//玩家重生宝物
//消息:/relive_gem
type MSG_ReliveGem_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	GemID      int    //宝物ID
	GemPos     int    //宝物位置
}

type MSG_ReliveGem_Ack struct {
	RetCode int            //返回码
	ItemLst []MSG_ItemData //分解所获物品
}

//请求挂机信息
//消息:/hangup_get_info
type MSG_GetHangUp_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type THisHang struct {
	BossID  int //
	ItemID  int
	ItemNum int
	Time    int32
}

type MSG_GetHangUp_Ack struct {
	RetCode   int        //返回码
	CurBossID int        //当前的BossID
	GridNum   int        //格子数
	ExpItems  []int      //经验丹数
	LeftQuick int        //剩下快速战斗次数
	History   []THisHang //历史记录
}

//设置新的挂机Boss
//消息:/hangup_set_boss
type MSG_SetHangUpBoss_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	BossID     int    //
}

type MSG_SetHangUpBoss_Ack struct {
	RetCode   int //返回码
	CurBossID int //当前的BossID
}

//一键使用经验丹
//消息:/hangup_use_exp
type MSG_UseExpItem_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_UseExpItem_Ack struct {
	RetCode int //返回码
	CurExp  int //增加的经验值
}

//增加格子数
//消息:/hangup_add_grid
type MSG_AddGridNum_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_AddGridNum_Ack struct {
	RetCode int //返回码
	GridNum int //新格子数
}

//快速战斗请法语
//消息:/hangup_quick_fight
type MSG_QuickFight_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_QuickFight_Ack struct {
	RetCode   int        //返回码
	ExpItems  []int      //经验丹数
	History   []THisHang //历史记录
	QuickTime int        //己使用的快速战斗次数
}

type MSG_PlayerInfo struct {
	PlayerID   int32
	HeroID     int    //英雄ID
	Name       string //名字
	GuildName  string //工会名
	GuildIcon  int    //工会图标
	Level      int
	FightValue int32
	Quality    int8
	Value      int
	Camp       int8
}

//请求等级排行榜
//消息:/get_level_rank
type MSG_GetLevelRank_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_GetLevelRank_Ack struct {
	RetCode   int              //返回码
	PlayerLst []MSG_PlayerInfo //玩家信息列表
	MyRank    int              //自己的排名
}

//请求等级排行榜
//消息:/get_fight_rank
type MSG_GetFightRank_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_GetFightRank_Ack struct {
	RetCode   int              //返回码
	PlayerLst []MSG_PlayerInfo //玩家信息列表
	MyRank    int              //自己的排名
}

//! 请求三国无双全服排行榜
//! 消息: /get_sanguows_rank
type MSG_GetSanguows_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_SanguowsInfo struct {
	Name       string //角色名字
	HeroID     int    //英雄ID
	Star       int    //星数
	FightValue int32  //战力
	Quality    int8
}

type MSG_GetSanguows_Ack struct {
	RetCode   int
	PlayerLst []MSG_SanguowsInfo
	MyRank    int
	MyStar    int //自己的星数
}

type MSG_ArenaInfo struct {
	PlayerID   int32
	Name       string //角色名字
	HeroID     int    //英雄ID
	FightValue int32  //战力
	Quality    int8
	Level      int //等级
}

//! 玩家请求竞技场名次排行榜信息
//! 消息: /get_arena_rank
type MSG_GetArenaRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetArenaRank_Ack struct {
	RetCode   int
	PlayerLst []MSG_ArenaInfo
	MyRank    int
}

type MSG_RebelRankInfo struct {
	PlayerID   int32  //! 角色ID
	HeroID     int    //! 英雄ID
	Name       string //! 角色名字
	FightValue int32  //! 战力值
	Level      int    //! 角色等级
	Value      int    //! 伤害值/功勋值
	Quality    int8
}

//! 玩家请求查询排行榜
//! 消息: /get_rebel_rank
type MSG_GetRebelRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetRebelRank_Ack struct {
	RetCode        int                 //返回码
	ExploitRankLst []MSG_RebelRankInfo //功勋榜
	DamageRankLst  []MSG_RebelRankInfo //伤害榜
	MyExploitRank  int                 //自怀的功勋排名
	MyDamageRank   int                 //自己的伤害排名
}

type MSG_GuildRankInfo struct {
	GuildID     int32  //公会ID
	GuildName   string //公公名
	Icon        int    //公会图标
	Level       int    //公会等级
	CurNum      int
	MaxNum      int
	Name        string //团长名
	CopyChapter int32  //副本章节
}

//! 玩家请求公会等级排行榜
//! 消息: /get_guild_level_rank
type MSG_GetGuildLevelRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildLevelRank_Ack struct {
	RetCode   int                 //返回码
	GuildList []MSG_GuildRankInfo //功勋榜
}

//! 玩家请求公会副本排行榜
//! 消息: /get_guild_copy_rank
type MSG_GetGuildCopyRank_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetGuildCopyRank_Ack struct {
	RetCode   int                 //返回码
	GuildList []MSG_GuildRankInfo //功勋榜
}

//! 玩家请求己收积过的英雄列表
//! 消息: /get_collection_heros
type MSG_GetCollectionHeros_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetCollectionHeros_Ack struct {
	RetCode int     //返回码
	Heros   []int16 //英雄表
}

//! 玩家请求重置云游信息
//! 消息: /wander_getinfo
type MSG_WanderGetInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_WanderGetInfo_Ack struct {
	RetCode   int //返回码
	MaxCopyID int //当前的云游的最大副本ID
	CurCopyID int //当前的云游的当前战斗副本ID
	CanBattle int //是否可以战斗0:不可以， 1:可以。
	LeftTime  int //剩余重置次数
}

//! 玩家请求重置云游战斗次数
//! 消息: /wander_reset
type MSG_WanderReset_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_WanderReset_Ack struct {
	RetCode   int //返回码
	CurCopyID int //当前云游副本ID
	LeftTime  int //剩余重置次数
	CanBattle int //是否可以战斗0:不可以， 1:可以。
}

//! 玩家请求开云游宝箱
//! 消息: /wander_openbox
type MSG_WanderOpenBox_Req struct {
	PlayerID   int32
	SessionKey string
	DrawType   int // 1 : 单抽， 2 : 十连抽
}

type MSG_WanderOpenBox_Ack struct {
	RetCode int //返回码
	ItemLst []MSG_ItemData
}

//! 玩家请求扫荡云游战斗次数
//! 消息: /wander_sweep
type MSG_WanderSweep_Req struct {
	PlayerID     int32
	SessionKey   string
	TargetCopyID int //目标副本ID
}

type MSG_WanderSweep_Ack struct {
	RetCode   int //返回码
	CurCopyID int //当前战胜的副本ID
	ItemLst   []MSG_ItemData
}

//请求云游排行榜
//消息:/get_wander_rank
type MSG_GetWanderRank_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_GetWanderRank_Ack struct {
	RetCode   int              //返回码
	PlayerLst []MSG_PlayerInfo //玩家信息列表
	MyRank    int              //自己的排名
}

//请求阵营战排行榜
//消息:/get_campbat_rank
type MSG_GetCampBatRank_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	RankType   int    //排榜类型
}

type MSG_GetCampBatRank_Ack struct {
	RetCode   int              //返回码
	PlayerLst []MSG_PlayerInfo //玩家信息列表
	MyRank    int              //自己的排名
}

//! 玩家请求云游战斗检查
//! 消息: /wander_check
type MSG_WanderCheck_Req struct {
	PlayerID     int32
	SessionKey   string
	TargetCopyID int //
}

type MSG_WanderCheck_Ack struct {
	RetCode int //返回码
}

//! 玩家请求云游战斗结果
//! 消息: /wander_result
type MSG_WanderResult_Req struct {
	PlayerID     int32
	SessionKey   string
	TargetCopyID int //副本ID
	Win          int //战斗结果 0: 失败 ,1 胜利
	//英雄核查数据
	HeroCkD []MSG_HeroCheckData
}

type MSG_WanderResult_Ack struct {
	RetCode   int            //返回码
	CurCopyID int            //当前战胜的副本ID
	ItemLst   []MSG_ItemData //奖励列表
}

//! 请求全部的红点提示
//! 消息: /get_mainui_tip
type MSG_GetMainUITip_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetMainUITip_Ack struct {
	RetCode int   //返回码
	FuncID  []int //功能ID
}

//游戏服向账号服务器注册消息
type MSG_RegBattleSvr_Req struct {
	BatSvrID         int // 就是port号
	ServerDomainName string
	ServerOuterAddr  string
	ServerInnerAddr  string
}

//游戏服向账号服务器注册的返回消息
type MSG_RegBattleSvr_Ack struct {
	RetCode int
}

//! 玩家请求设置时装
//! 消息: /fashion_set
type MSG_FashionSet_Req struct {
	PlayerID   int32
	SessionKey string
	FashionID  int32 //时装ID
}

type MSG_FashionSet_Ack struct {
	RetCode    int
	FightValue int32 //战力
}

//! 玩家请求强化时装
//! 消息: /fashion_strength
type MSG_FashionStrength_Req struct {
	PlayerID   int32
	SessionKey string
	FashionID  int32 //时装ID
}

type MSG_FashionStrength_Ack struct {
	RetCode    int
	FID        int32 //时装ID
	FLevel     int32 //时装等级
	FightValue int32 //战力
}

//! 玩家请求重铸时装
//! 消息: /fashion_recast
type MSG_FashionRecast_Req struct {
	PlayerID   int32
	SessionKey string
	FashionID  int32 //时装ID
}

type MSG_FashionRecast_Ack struct {
	RetCode    int
	FID        int32 //时装ID
	FLevel     int32 //时装等级
	MoneyID    int   //钱ID
	MoneyNum   int   //钱数
	CostID     int   //道具ID
	CostNum    int   //道具数
	FightValue int32 //战力值
}

//! 玩家请求合成时装
//! 消息: /fashion_compose
type MSG_FashionCompose_Req struct {
	PlayerID   int32
	SessionKey string
	FashionID  int32 //时装ID
}

type MSG_FashionCompose_Ack struct {
	RetCode   int
	FashionID int32 //合成后的时装ID
}

//! 时装熔炼值
//! 消息: /fashion_melt_value
type MSG_FashionMeltValue_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FashionMeltValue_Ack struct {
	RetCode int
	Value   int32 //值
}

//! 时装熔炼
//! 消息: /fashion_melting
type MSG_FashionMelting_Req struct {
	PlayerID   int32
	SessionKey string
	PieceID    int //碎片ID
	PieceNum   int //碎片个数
}

type MSG_FashionMelting_Ack struct {
	RetCode int
	Value   int32 //值
}

//! 时装熔炼奖励
//! 消息: /fashion_melt_award
type MSG_FashionMeltAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_FashionMeltAward_Ack struct {
	RetCode   int
	FashionID int32 //时装ID
	PieceID   int   //碎片ID
	PieceNum  int   // 碎片数量
}

//! 游戏服下发的主动通知
type MSG_GameSvr_Notify struct {
	FuncID int //需要加红点的应用
}

// //! 好友系统主动通知
// type MSG_GameSvr_Nofity_Friend struct {
// 	Action         int    //! 1->被申请添加好友  2->被同意添加好友
// 	TargetPlayerID int32  //! 目标ID
// 	TargetName     string //! 名称
// 	HeroID         int32  //! 主英雄ID
// 	Quality        int32  //! 品质
// }

// //! 公会系统主动通知
// type MSG_GameSvr_Nofity_Gulid struct {
// 	Action int //! 1->同意入会 2->有人申请公会
// }

//发送私人邮件
//消息:/send_private_mail
type MSG_SendPrivateMail_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
	TargetID   int32  //目标玩家ID
	Title      string //标题
	Content    string //内容
}

type MSG_SendPrivateMail_Ack struct {
	RetCode int //返回码
}

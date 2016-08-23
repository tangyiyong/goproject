package mainlogic

//玩家加载登录数据
//消息:/get_login_data
type MSG_GetLoginData_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetLoginData_Ack struct {
	RetCode int //返回码
}

//玩家请求上阵数据
//消息:/get_battle_data
type MSG_GetBattleData_Req struct {
	PlayerID   int32
	SessionKey string
}

//玩家请求上阵英雄数据回复
type MSG_GetBattleData_Ack struct {
	RetCode   int            //返回码
	CurHeros  [6]THeroData   //六个上阵英雄
	BackHeros [6]THeroData   //六个援军英雄
	Equips    [24]TEquipData //上阵装备信息
	Gems      [12]TGemData   //上阵宝物信息
	Pets      [6]TPetData    //上阵宠物
	Title     int            //称号ID
}

///////////////////////////////////////////////
//玩家请求背包数据
//消息:/get_bag_data
type MSG_GetBagData_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagData_Ack struct {
	RetCode     int
	Normals     []TItemData    //普通道具背包
	Heros       []THeroData    //英雄
	Equips      []TEquipData   //装备
	HeroPieces  []TItemData    //英雄碎片
	EquipPieces []TItemData    //装备碎片
	Gems        []TGemData     //宝物
	GemPieces   []TItemData    //宝物碎片
	Pets        []TPetData     //宠物
	PetPieces   []TItemData    //宠物碎片
	WakeItems   []TItemData    //觉醒道具
	HeroSouls   []TItemData    //英魂道具
	Fashions    []TFashionData //时装
	FasPieces   []TItemData    //时装碎片
}

///////////////////////////////////////////////
//玩家请求背包英雄数据
//消息:/get_bag_heros
type MSG_GetBagHeros_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagHeros_Ack struct {
	RetCode int
	Heros   []THeroData
}

//玩家请求背包装备数据
//消息:/get_bag_equips
type MSG_GetBagEquip_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagEquip_Ack struct {
	RetCode int
	Equips  []TEquipData
}

//玩家请求英雄碎片数据
//消息:/get_bag_hero_piece
type MSG_GetBagHerosPiece_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagHerosPiece_Ack struct {
	RetCode    int
	HeroPieces []TItemData
}

//玩家请求背包装备碎片
//消息:/get_bag_equip_piece
type MSG_GetBagEquipPiece_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagEquipPiece_Ack struct {
	RetCode     int
	EquipPieces []TItemData
}

//玩家请求背包宝石数据
//消息:/get_bag_gems
type MSG_GetBagGems_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagGems_Ack struct {
	RetCode int
	Gems    []TGemData
}

//请求背包里的道具
//消息:/get_bag_items
type MSG_GetBagItems_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagItems_Ack struct {
	RetCode int
	Items   []TItemData
}

//请求背包里的觉醒道具
//消息:/get_bag_wake_items
type MSG_GetBagWakeItems_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagWakeItems_Ack struct {
	RetCode int
	Items   []TItemData
}

//玩家请求背包宝物碎片
//消息:/get_bag_gem_piece
type MSG_GetBagGemPiece_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagGemPiece_Ack struct {
	RetCode   int
	GemPieces []TItemData
}

///////////////////////////////////////////////
//玩家请求背包宠物数据
//消息:/get_bag_pets
type MSG_GetBagPets_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagPets_Ack struct {
	RetCode int
	Pets    []TPetData
}

//玩家请求宠物碎片数据
//消息:/get_bag_pet_piece
type MSG_GetBagPetsPiece_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetBagPetsPiece_Ack struct {
	RetCode   int
	PetPieces []TItemData
}

package msg

//! 玩家获取商店商品信息
//! 消息: /get_all_store_data
type MSG_GetAllStoreData_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_StoreData struct {
	GoodsInfoLst   []MSG_StoreItem
	FreeCount      int   //! 当前可免费刷新次数
	FreeCountLimit int   //! 当天免费刷新上限
	RefreshCount   int   //! 当天还剩多少次数可刷新
	FreeRefeshTime int32 //! 倒计时
}

type MSG_GetAllStoreData_Ack struct {
	RetCode int
	Hero    MSG_StoreData //! 英雄商店
	Awake   MSG_StoreData //! 觉醒商店
	Pet     MSG_StoreData //! 战宠商店
}

//! 玩家获取神将商店商品信息
//! 消息:/get_store
type MSG_GetStoreData_Req struct {
	PlayerID   int32
	SessionKey string
	StoreType  int
}

type MSG_StoreItem struct {
	ID     int //! 唯一标识
	Status int //! 0->未购买 1->已购买  不得重复购买
}

type MSG_GetStoreData_Ack struct {
	RetCode        int //! 返回码
	GoodsInfoLst   []MSG_StoreItem
	FreeCount      int   //! 当前可免费刷新次数
	FreeCountLimit int   //! 当天免费刷新上限
	RefreshCount   int   //! 当天还剩多少次数可刷新
	FreeRefeshTime int32 //! 倒计时
}

//! 刷新神将商店
//! 消息: /refresh_store
type MSG_RefreshStore_Req struct {
	PlayerID   int32
	SessionKey string
	StoreType  int
}

type MSG_RefreshStore_Ack struct {
	RetCode        int   //! 返回码
	FreeCount      int   //! 当前可免费刷新次数
	FreeCountLimit int   //! 当天免费刷新上限
	RefreshCount   int   //! 当天还剩多少次数可刷新
	FreeRefeshTime int32 //! 倒计时
	CostType       int   //! 1->免费次数 2->道具 3->货币
	CostNum        int   //! 扣除次数
	GoodsInfoLst   []MSG_StoreItem
}

//! 购买神将商店商品
//! 消息: /store_buy
type MSG_StoreBuyItem_Req struct {
	PlayerID   int32
	SessionKey string
	StoreType  int
	Index      int
}

type MSG_StoreBuyItem_Ack struct {
	RetCode int //! 返回码
}

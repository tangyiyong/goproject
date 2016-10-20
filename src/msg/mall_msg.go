package msg

//! 玩家请求获取VIP礼包信息
//! 消息: /get_vip_gift
type MSG_GetVipGifts_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetVipGifts_Ack struct {
	RetCode int     //! 返回码
	ID      []int32 //! 显示礼包ID
}

//! 玩家请求购买VIP礼包
//! 消息: /buy_vip_gift
type MSG_BuyVipGift_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int32 //! ID
}

type MSG_BuyVipGift_Ack struct {
	RetCode  int     //! 返回码
	ID       []int32 //! 显示礼包ID
	MoneyID  int
	MoneyNum int
	ItemID   int32
	ItemNum  int32
}

//! 玩家请求查询道具商品购买次数
//! 消息: /get_goods_buy_times
type MSG_GetMallBuyTimes_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetMallBuyTimes_Ack struct {
	RetCode     int
	BuyTimesLst []MSG_BuyData
}

//! 玩家请求购买普通商品
//! 消息: /buy_goods
type MSG_BuyGoods_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int32 //! 物品
	Num        int   //! 数量
}

type MSG_BuyGoods_Ack struct {
	RetCode  int //! 返回码
	BuyTimes MSG_BuyData
	MoneyID  int
	MoneyNum int
}

//! 玩家请求物品剩余信息
//! 消息: /get_one_buy_times
type MSG_GetBuyTimes_Req struct {
	PlayerID   int32
	SessionKey string
	ItemID     int32
}

type MSG_GetBuyTimes_Ack struct {
	RetCode  int   //! 返回码
	FuncID   int   //! 功能ID
	ItemID   int32 //! 物品ID
	BuyTimes int   //! 剩余次数  返回负数则表示不限购, 可无限购买
}

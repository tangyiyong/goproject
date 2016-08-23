package msg

//! 玩家请求获取VIP礼包信息
//! 消息: /get_vip_gift
type MSG_GetVipGifts_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetVipGifts_Ack struct {
	RetCode int   //! 返回码
	ID      []int //! 显示礼包ID
}

//! 玩家请求购买VIP礼包
//! 消息: /buy_vip_gift
type MSG_BuyVipGift_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int //! ID
}

type MSG_BuyVipGift_Ack struct {
	RetCode  int   //! 返回码
	ID       []int //! 显示礼包ID
	MoneyID  int
	MoneyNum int
}

//! 玩家请求查询道具商品购买次数
//! 消息: /get_goods_buy_times
type MSG_GetGoodsBuyTimes_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GoodsBuyTimesInfo struct {
	ID    int
	Times int //! 还能购买次数
}

type MSG_GetGoodsBuyTimes_Ack struct {
	RetCode     int
	BuyTimesLst []MSG_GoodsBuyTimesInfo
}

//! 玩家请求购买普通商品
//! 消息: /buy_goods
type MSG_BuyGoods_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int //! 物品
	Num        int //! 数量
}

type MSG_BuyGoods_Ack struct {
	RetCode  int //! 返回码
	BuyTimes MSG_GoodsBuyTimesInfo
	MoneyID  int
	MoneyNum int
}

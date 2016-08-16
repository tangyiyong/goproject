package msg

//! 玩家查询黑市信息
//! 消息: /get_black_market_info
type MSG_GetBlackMarket_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_BlackMarketGoods struct {
	ID           int
	ItemID       int
	ItemNum      int
	CostMoneyID  int
	CostMoneyNum int
	IsBuy        bool
	Recommend    int
}

type MSG_GetBlackMarket_Ack struct {
	RetCode     int
	GoodsLst    []MSG_BlackMarketGoods //! 商品
	OpenEndTime int64                  //! 结束时间戳
	RefreshTime int64                  //! 刷新时间戳
}

//! 玩家查询黑市状态
//! 消息: /get_black_market_status
type MSG_GetBlackMarketStatus_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetBlackMarketStatus_Ack struct {
	RetCode     int
	OpenEndTime int64
}

//! 玩家购买黑市商品
//! 消息: /buy_black_market
type MSG_BuyBlackMarket_Req struct {
	PlayerID   int
	SessionKey string
	ID         int
}

type MSG_BuyBlackMarket_Ack struct {
	RetCode int
}

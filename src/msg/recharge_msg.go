package msg

//http.HandleFunc("/get_charge_result", mainlogic.Hand_GetFirstRechargeGift) //! 玩家请求充值结果
//http.HandleFunc("/charge_money", mainlogic.Hand_GetFirstRechargeStatus)    //! 玩家请求充值

//! 玩家请求充值信息
//get_charge_info
type MSG_GetChargeInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetChargeInfo_Ack struct {
	RetCode          int
	ChargeTimes      []int  //首充状态
	CardDays         [2]int //月卡剩余天数
	ActivityChargeID int    //优惠充值ID（0无）
}

//! 玩家请求充值结果
//! 消息: /get_charge_result
type MSG_GetChargeResult_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetChargeResult_Ack struct {
	RetCode  int  //返回码
	VipLevel int8 //当前的VIP等级
	VipExp   int  //当前的Vip经验
	MoneyNum int  //元宝数
}

//! 玩家领取月卡
//! 消息: /receive_month_card
type MSG_ReceiveMonthCard_Req struct {
	PlayerID   int32
	SessionKey string
	CardID     int //月卡ID
}

type MSG_ReceiveMonthCard_Ack struct {
	RetCode   int            //返回码
	AwardItem []MSG_ItemData //! 奖励
	CardID    int            //! 月卡ID
}

//! 玩家领取激活码
//! 消息: /recv_gift_code
type MSG_RecvGiftCode_Req struct {
	PlayerID   int32
	SessionKey string
	GiftCode   string //激活码
}

type MSG_RecvGiftCode_Ack struct {
	RetCode   int            //返回码
	AwardItem []MSG_ItemData //! 奖励
}

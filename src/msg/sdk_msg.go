package msg

// 消息：
//recharge_notify
type Msg_Recharge_Notify struct {
	PlayerID int32
	ChargeID int //充值表ID
	RMB      int //第三方发来的充值数
}

package msg

//请求推荐阵营
//消息:/get_recommandcamp
type MSG_GetRecommandCamp_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetRecommandCamp_Ack struct {
	RetCode    int //
	BattleCamp int //请求推荐的阵营
}

//请求设置阵营
//消息:/set_battlecamp
type MSG_SetBattleCamp_Req struct {
	PlayerID   int
	SessionKey string
	BattleCamp int //请求的阵营
}

type MSG_SetBattleCamp_Ack struct {
	RetCode int //
}

//请求阵营战主界面数据
//消息:/get_campbat_data
type MSG_GetCampBatData_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetCampBatData_Ack struct {
	RetCode   int    //返回码
	KillNum   int    //我的击杀数
	MyRank    int    //我的排名
	LeftTimes int    //剩余搬动次数
	CampKill  [3]int //排行数据
}

//请求进入阵营战
//消息:/enter_campbattle
type MSG_EnterCampBattle_Req struct {
	PlayerID   int
	SessionKey string
	BattleCamp int //请求的阵营
}

type MSG_EnterCampBattle_Ack struct {
	RetCode       int    //
	EnterCode     int    //进入阵营战协议码
	BattleSvrAddr string //阵营战服务器
}

//! 玩家请求积分商店的状态
//! 消息: /get_campbat_store_state
type MSG_GetCampbatStoreState_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetCampbatStoreState_Ack struct {
	RetCode    int                //返回码
	ItemLst    []MSG_StoreBuyData //购买物品次数
	AwardIndex []int              //奖励商店的索引
}

//! 玩家请求购买积分商店道具
//! 消息: /buy_campbat_store_item
type MSG_BuyCampbatStoreItem_Req struct {
	PlayerID    int
	SessionKey  string
	StoreItemID int //商店道具ID
	BuyNum      int //购买数量
}

type MSG_BuyCampbatStoreItem_Ack struct {
	RetCode int //返回码
}

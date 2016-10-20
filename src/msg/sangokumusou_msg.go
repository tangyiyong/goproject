package msg

//! 获取三国无双状态
//! 消息: /get_sangokumusou_status
type MSG_GetSangokuMusouStatus_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusouStatus_Ack struct {
	RetCode int

	CurStar       int  //! 当前星数
	HistoryStar   int  //! 历史最高星数
	CanUseStar    int  //! 当前可以使用星数
	IsEnd         int  //! 是否已经结束
	IsBuyTreasure bool //! 是否购买无双迷藏

	PassCopyID   int    //! 当前通关关卡
	IsRecvAward  int    //! 是否领取上章奖励
	IsSelectBuff int    //! 是否选择上章Buff
	BattleTimes  int    //! 关卡已挑战次数
	CopyLst      [3]int //! 关卡星数信息
	AttrLst      []MSG_SangokuMusou_Attr2
	ItemLst      []MSG_BuyData //! 已购买物品信息
}

//! 获取三国无双星数信息
//! 消息:/get_sangokumusou_star
type MSG_GetSangokuMusouStarInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusouStarInfo_Ack struct {
	RetCode     int
	CurStar     int //! 当前星数
	HistoryStar int //! 历史最高星数
	CanUseStar  int //! 当前可以使用星数
	IsEnd       int //! 是否已经结束
}

//! 获取三国无双闯关信息
//! 消息: /get_sangokumusou_copy
type MSG_GetSangokuMusouCopyInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusouCopyInfo_Ack struct {
	RetCode      int
	PassCopyID   int    //! 当前通关关卡
	IsRecvAward  int    //! 是否领取上章奖励
	IsSelectBuff int    //! 是否选择上章Buff
	BattleTimes  int    //! 关卡已挑战次数
	CopyLst      [3]int //! 关卡星数信息
}

//! 获取三国无双精英挑战闯关信息
//! 消息: /get_sangokumusou_elite_copy
type MSG_GetSangokuMusouEliteCopyInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusouEliteCopyInfo_Ack struct {
	RetCode         int
	PassEliteCopyID int //! 当前通关精英关卡
	BattleTimes     int //! 精英关卡已挑战次数
	HistoryCopyID   int //! 历史最高通关关卡ID
}

//! 通关三国无双
//! 消息: /pass_sangokumusou
type MSG_PassSangokuMusouCopy_Req struct {
	PlayerID   int32
	SessionKey string
	CopyType   int //! 0->普通副本 1->精英副本
	CopyID     int
	StarNum    int //! 获得星数
	IsVictory  int //! 是否胜利
}

type MSG_SangokuMusouDropItem struct {
	ItemID   int
	ItemNum  int
	CritType int //! 0->没暴击  1->暴击  2->大暴击  3->幸运暴击
}

type MSG_PassSangokuMusouCopy_Ack struct {
	RetCode  int
	DropItem []MSG_SangokuMusouDropItem
}

//! 通关三国无双精英挑战
//! 消息: /pass_sgws_elite
type MSG_PassSangokuMusouEliteCopy_Req struct {
	PlayerID   int32
	SessionKey string
	CopyID     int
}

type MSG_PassSangokuMusouEliteCopy_Ack struct {
	RetCode        int
	IsFirstVictory int //! 是否为首胜
}

//! 请求扫荡该章节
//! 消息: /sweep_sangoumusou
type MSG_PassSangokuMusouCopy_sweep_Req struct {
	PlayerID   int32
	SessionKey string
	Chapter    int
}

type MSG_PassSangokuMusouCopy_sweep_Ack struct {
	RetCode  int
	DropItem [][]MSG_SangokuMusouDropItem
}

//! 请求章节奖励
//! 消息: /get_sangokumusou_chapter_award
type MSG_GetSangokuMusouChapterAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusouChapterAward_Ack struct {
	RetCode  int
	AwardLst []MSG_ItemData
}

//! 请求随机章节属性奖励
//! 消息: /get_sangokumusou_attr
type MSG_GetSangokuMusouAttrInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_SangokuMusou_Attr1 struct {
	ID       int
	AttrID   int //! 属性
	Value    int //! 加成值
	CostStar int //! 消耗星数
}

type MSG_GetSangokuMusouAttrInfo_Ack struct {
	RetCode int
	AttrLst []MSG_SangokuMusou_Attr1
}

//! 请求玩家已选择所有属性奖励
//! 消息: /get_sangokumusou_all_attr
type MSG_GetSangokuMusouAllAttrInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_SangokuMusou_Attr2 struct {
	AttrID int //! 属性
	Value  int //! 加成值
}

type MSG_GetSangokuMusouAllAttrInfo_Ack struct {
	RetCode int
	AttrLst []MSG_SangokuMusou_Attr2
}

//! 选择章节属性奖励
//! 消息: /set_sangokumusou_attr
type MSG_SetSangokuMusouAttrInfo_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
}

type MSG_SetSangokuMusouAttrInfo_Ack struct {
	RetCode int
	AttrLst []MSG_SangokuMusou_Attr2
}

//! 请求无双秘藏
//! 消息: /get_sangokumusou_treasure
type MSG_GetSangokuMusouTreasureInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusouTreasureInfo_Ack struct {
	RetCode    int
	TreasureID int
}

//! 请求购买无双秘藏
//! 消息: /buy_sangokumusou_treasure
type MSG_BuySangokuMusouTreasure_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_BuySangokuMusouTreasure_Ack struct {
	RetCode int
}

//! 玩家请求重置普通挑战
//! 消息: /reset_sangokumusou_copy
type MSG_SangokuMusou_Reset_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_SangokuMusou_Reset_Ack struct {
	RetCode    int
	ResetTimes int
	MoneyID    int
	MoneyNum   int
}

//! 玩家查询精英挑战可增加次数
//! 消息: /get_sangoukumusou_elite_add_times
type MSG_GetSangokuMusou_Add_BattleTimes_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetSangokuMusou_Add_BattleTimes_Ack struct {
	RetCode    int
	ResetTimes int
}

//! 玩家请求增加精英挑战次数
//! 消息: /add_sangoukumusou_elite_copy
type MSG_SangokuMusou_Add_BattleTimes_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_SangokuMusou_Add_BattleTimes_Ack struct {
	RetCode int
}

//! 请求三国无双商店查询已购买物品信息
//! 消息: /get_sangokumusou_store_aleady_buy
type MSG_GetSangokuMusouStoreAleadyBuy_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_BuyData struct {
	ID    int32 //! 物品ID
	Times int   //! 购买次数
}

type MSG_GetSangokuMusouStoreAleadyBuy_Ack struct {
	RetCode int
	ItemLst []MSG_BuyData
}

//! 请求购买三国无双商店物品
//! 消息: /buy_sangokumusou_store
type MSG_BuySangokuMusouStoreItem_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int
	Num        int
}

type MSG_BuySangokuMusouStoreItem_Ack struct {
	RetCode int
	Times   int
}

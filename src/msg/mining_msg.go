package msg

//! 玩家查询当前挖矿信息
//! 消息: /enter_mining
type MSG_GetMiningInfo_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_MiningBuff struct {
	BuffType int
	Value    int
	Times    int
}

type MSG_GetMiningInfo_Ack struct {
	RetCode int
	MapData MSG_MiningMapData //! 地图信息
}

//! 玩家请求当前矿洞状态码
//! 消息: /get_mining_status
type MSG_GetMiningStatus_Req struct {
	PlayerID   int32
	SessionKey string
	StatusCode int
}

type MSG_GetMiningStatus_Ack struct {
	RetCode     int
	IsVerified  bool           //! 验证是否通过
	Buff        MSG_MiningBuff //! Buff列表
	LastPos     MSG_Pos        //! 最后一次玩家操作坐标
	Point       int            //! 玩家当前积分
	StatusCode  int            //! 状态码
	GuajiType   int            //! 挂机类型
	GuajiTime   int32          //! 距离挂机结算剩余时间
	ResetTimes  int            //! 重置次数
	GuajiStatus bool           //! 0->未挂机  1->挂机中
}

//! 玩家请求获取某个坐标点信息
//! 消息: /mining_dig
type MSG_Pos struct {
	X int32
	Y int32
}

type MSG_MiningDig_Req struct {
	PlayerID   int32
	SessionKey string
	Pos        []MSG_Pos //! 坐标 左上角为原点
}

type MSG_MiningDigData struct {
	IsDig       bool  //! 是否已经被挖掘
	Element     int32 //! 元素
	Monster     int   //! 元素类型为事件且事件类型为怪物,则Monster=怪物ID
	MonsterLife int   //! 怪物当前剩余血量
}

type MSG_MiningDig_Ack struct {
	RetCode    int                 //! 返回码
	MapData    []MSG_MiningDigData //! 地图信息
	StatusCode int                 //! 状态码
}

//! 玩家请求事件对应处理-行动力奖励
//! 消息: /mining_event_action_award
type MSG_MiningEvent_ActionAward_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningEvent_ActionAward_Ack struct {
	RetCode      int
	AddActionID  int                 //! 行动力ID
	AddActionNum int                 //! 增加行动力值
	StatusCode   int                 //! 状态码
	VisualPos    []MSG_MiningDigData //! 新可视区域
	Point        int
	ActionValue  int   //! 行动力值
	ActionTime   int32 //! 行动力恢复起始时间
}

//! 挖矿元素-精炼石
//! 消息: /mining_element_stone
type MSG_MiningElement_GetStone_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningElement_GetStone_Ack struct {
	RetCode     int
	ItemID      int
	ItemNum     int
	Point       int
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

//! 挖矿事件-黑市
//! 消息: /mining_event_black_market
type MSG_MiningEvent_GetBlackMarket_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningEvent_GetBlackMarket_Ack struct {
	RetCode     int
	GoodsLst    []int               //! 黑市商品
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

//! 购买黑市商品
//! 消息: /mining_buy_black_market
type MSG_MiningEvent_BuyBlackMarket_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int32
}

type MSG_MiningEvent_BuyBlackMarket_Ack struct {
	RetCode     int
	Point       int   //! 玩家当前积分
	ActionValue int   //! 行动力值
	ActionTime  int32 //! 行动力恢复起始时间
}

//! 挖矿事件-查看怪物
//! 消息: /mining_event_monster_info
type MSG_MiningEvent_Monster_Info_Req struct {
	PlayerID   int32
	SessionKey string
	Pos        MSG_Pos
}

type MSG_MiningEvent_Monster_Info_Ack struct {
	RetCode     int
	Level       int
	Life        int
	TotalLife   int
	CopyID      int
	MonsterType int32
}

//! 挖矿事件-怪物
//! 消息: /mining_event_monster
type MSG_MiningEvent_Monster_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
	Damage     int //! 对怪物造成伤害
	//英雄核查数据
	HeroCkD []MSG_HeroCheckData
}

type MSG_MiningEvent_Monster_Ack struct {
	RetCode     int
	IsKill      bool                //! 是否击杀
	DropItem    []MSG_ItemData      //! 掉落物品
	Point       int                 //! 当前积分
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

//! 挖矿事件-宝箱
//! 消息: /mining_event_treasure
type MSG_MiningEvent_Treasure_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningEvent_Treasure_Ack struct {
	RetCode     int
	AwardItem   []MSG_ItemData
	Point       int
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

//! 挖矿事件-魔盒
//! 消息: /mining_event_box
type MSG_MiningEvent_Box_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningEvent_Box_Ack struct {
	RetCode     int
	RandPoint   int                 //! 抽取的分数
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

//! 挖矿事件-扫描
//! 消息: /mining_event_scan
type MSG_MiningEvent_Scan_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningEvent_Scan_Ack struct {
	RetCode     int
	StatusCode  int //! 状态码
	Point       int
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

type MSG_MiningMonster struct {
	Index int32
	ID    int
	Life  int
}

type MSG_MiningMapData struct {
	DigStatus   [60]string
	MonsterInfo []MSG_MiningMonster
	Element     []int32 //! 前16位为index  后16位为值
}

//! 挖矿事件-答题
//! 消息: /mining_event_question
type MSG_MiningEvent_Question_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
	AddPoint   int //! 答对题目后增加的积分
}

type MSG_MiningEvent_Question_Ack struct {
	RetCode     int
	Point       int                 //! 玩家当前积分
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	ActionValue int                 //! 行动力值
	ActionTime  int32               //! 行动力恢复起始时间
}

//! 挖矿事件-Buff
//! 消息: /mining_event_buff
type MSG_MiningEvent_Buff_Req struct {
	PlayerID   int32
	SessionKey string
	PlayerPos  MSG_Pos
}

type MSG_MiningEvent_Buff_Ack struct {
	RetCode     int
	Buff        MSG_MiningBuff
	StatusCode  int                 //! 状态码
	VisualPos   []MSG_MiningDigData //! 新可视区域
	Point       int
	ActionValue int   //! 行动力值
	ActionTime  int32 //! 行动力恢复起始时间
}

//! 请求随机九种打完Boss翻牌奖励
//! 消息: /mining_get_award
type MSG_MiningGetAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_MiningAward struct {
	ID      int
	ItemID  int
	ItemNum int
	Status  bool //! 是否翻牌
}

type MSG_MiningGetAward_Ack struct {
	RetCode  int
	AwardLst []MSG_MiningAward
}

//! 玩家请求选择打完Boss后翻牌奖励
//! 消息: /select_mining_award
type MSG_MiningSelectAward_Req struct {
	PlayerID   int32
	SessionKey string
	SelectID   int
}

type MSG_MiningSelectAward_Ack struct {
	RetCode int
	Point   int  //! 玩家当前积分
	IsEnd   bool //! 是否结束
}

//! 玩家请求挂机
//! 消息: /mining_guaji
type MSG_MiningGuaji_Req struct {
	PlayerID   int32
	SessionKey string
	ID         int //! 挂机类型ID
}

type MSG_MiningGuaji_Ack struct {
	RetCode       int
	GuajiCalcTime int32 //! 距离挂机结算剩余时间
	ActionValue   int   //! 行动力值
	ActionTime    int32 //! 行动力恢复起始时间
}

//! 玩家查询挂机倒计时
//! 消息: /mining_guaji_time
type MSG_MiningGuajiTime_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_MiningGuajiTime_Ack struct {
	RetCode     int
	GuajiTime   int32 //! 距离挂机结算剩余时间
	GuajiStatus bool  //! 0->未挂机  1->挂机中  2->挂机结束未领取
}

//! 玩家请求领取挂机结算奖励
//! 消息: /mining_guaji_award
type MSG_GetMiningGuajiAward_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_GetMiningGuajiAward_Ack struct {
	RetCode int
	ItemLst []MSG_ItemData
}

//! 玩家请求重置矿洞信息
//! 消息: /mining_map_reset
type MSG_MiningMapReset_Req struct {
	PlayerID   int32
	SessionKey string
}

type MSG_MiningMapReset_Ack struct {
	RetCode int
}

package gamedata

//道具类型
const (
	TYPE_NORMAL        = 1  //普通道具
	TYPE_HERO          = 2  //英雄
	TYPE_HERO_PIECE    = 3  //英雄碎片
	TYPE_EQUIPMENT     = 4  //装备
	TYPE_EQUIP_PIECE   = 5  //装备碎片
	TYPE_GEM           = 6  //宝石
	TYPE_GEM_PIECE     = 7  //宝石碎片
	TYPE_WAKE          = 8  //觉醒道具
	TYPE_MONEY         = 9  //货币
	TYPE_ACTION        = 10 //行动力
	TYPE_PET           = 11 //宠物
	TYPE_PET_PIECE     = 12 //宠物碎片
	TYPE_HEROSOUL      = 13 //将灵
	TYPE_FASHION       = 14 //时装
	TYPE_FASHION_PIECE = 15 //时装碎片
)

//道具子类型
const (
	SUB_TYPE_MONEY        = 1  //货币
	SUB_TYPE_ACTION       = 2  //行动力
	SUB_TYPE_EQUIP_REFINE = 3  //装备精炼道具
	SUB_TYPE_GEM_STRENGTH = 4  //宝物强化道具
	SUB_TYPE_GEM_REFINE   = 5  //宝物精炼道具
	SUB_TYPE_GIFT_BAG     = 6  //礼包道具
	SUB_TYPE_FREE_WAR     = 7  //免战道具
	SUB_TYPE_PET_UPLVL    = 8  //宠物升级道具
	SUB_TYPE_PET_GOD      = 9  //宠物神炼道具
	SUB_TYPE_CHARGE       = 10 //充值额度道具
)

//角色攻击类型
const (
	ATYPE_PHYSI   = 1 //物攻类型
	ATYPE_MAGIC   = 2 //魔攻类型
	ATYPE_DEFFEND = 3 //辅助类型
	ATYPE_ASSIT   = 4 //防御类型
)

type ST_Range struct {
	Value [2]int //两个值
}

type ST_PropertyBuff struct {
	PropertyID int  //属性ID
	Value      int  //属性值
	IsPercent  bool //是否百分比
}

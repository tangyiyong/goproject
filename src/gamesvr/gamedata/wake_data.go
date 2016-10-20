package gamedata

import (
	"gamelog"
)

type ST_WakeLevel struct {
	Level            int    //充值ID
	NeedLevel        int    //需求等级
	NeedItem         [4]int //需要的道具
	NeedMoneyID      int
	NeedMoneyNum     int
	NeedHeroNum      int
	NeedWakeID       int     //觉醒丹ID
	NeedWakeNum      int     //觉醒丹数
	HostWakeNum      int     //主角觉醒丹数
	ExtraProperty    int     //额外属性ID
	ExtraValue       int     //额外属性值
	PropertyValues   [11]int //属性数值
	PropertyPercents [11]int //属性增加百分比
	bCalculate       bool    //是否己计算
}

type ComposeItem struct {
	ItemID  int
	ItemNum int
}

type ST_WakeCompose struct {
	ID       int            //道具ID
	MoneyID  int            //货币ID
	MoneyNum int            //货币值
	Items    [4]ComposeItem //合成需要的道具
}

var (
	GT_WakeLevelList   []ST_WakeLevel          //觉醒等级
	GT_WakeComposeList map[int]*ST_WakeCompose //觉配合成
	MaxWakeLevel       int                     //最大觉醒等级
)

func InitWakeLevelParser(total int) bool {
	GT_WakeLevelList = make([]ST_WakeLevel, total+1)
	MaxWakeLevel = 0
	return true
}

func InitWakeComposeParser(total int) bool {
	GT_WakeComposeList = make(map[int]*ST_WakeCompose, 1)
	return true
}

func ParseWakeLevelRecord(rs *RecordSet) {
	level := rs.GetFieldInt("level")
	GT_WakeLevelList[level].Level = level
	GT_WakeLevelList[level].NeedItem[0] = rs.GetFieldInt("need_item_1")
	GT_WakeLevelList[level].NeedItem[1] = rs.GetFieldInt("need_item_2")
	GT_WakeLevelList[level].NeedItem[2] = rs.GetFieldInt("need_item_3")
	GT_WakeLevelList[level].NeedItem[3] = rs.GetFieldInt("need_item_4")
	GT_WakeLevelList[level].NeedMoneyID = rs.GetFieldInt("money_id")
	GT_WakeLevelList[level].NeedMoneyNum = rs.GetFieldInt("money_num")
	GT_WakeLevelList[level].NeedHeroNum = rs.GetFieldInt("same_hero_num")
	GT_WakeLevelList[level].ExtraProperty = rs.GetFieldInt("extra_property")
	GT_WakeLevelList[level].ExtraValue = rs.GetFieldInt("extra_property_value")
	GT_WakeLevelList[level].NeedLevel = rs.GetFieldInt("need_level")
	GT_WakeLevelList[level].NeedWakeNum = rs.GetFieldInt("need_wake_num")
	GT_WakeLevelList[level].HostWakeNum = rs.GetFieldInt("host_wake_num")
	GT_WakeLevelList[level].NeedWakeID = rs.GetFieldInt("need_wake_id")
	GT_WakeLevelList[level].bCalculate = false

	if MaxWakeLevel < level {
		MaxWakeLevel = level
	}
}

func ParseWakeComposeRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	compose := new(ST_WakeCompose)
	compose.ID = id

	compose.Items[0].ItemID = rs.GetFieldInt("item_1")
	compose.Items[0].ItemNum = rs.GetFieldInt("item_1_num")
	compose.Items[1].ItemID = rs.GetFieldInt("item_2")
	compose.Items[1].ItemNum = rs.GetFieldInt("item_2_num")
	compose.Items[2].ItemID = rs.GetFieldInt("item_3")
	compose.Items[2].ItemNum = rs.GetFieldInt("item_3_num")
	compose.Items[3].ItemID = rs.GetFieldInt("item_4")
	compose.Items[3].ItemNum = rs.GetFieldInt("item_4_num")
	compose.MoneyID = rs.GetFieldInt("money_id")
	compose.MoneyNum = rs.GetFieldInt("money_num")
	GT_WakeComposeList[id] = compose
}

func FinishWakeLevelParser() bool {
	for level := 0; level <= MaxWakeLevel; level++ {
		pWakeLevel := &GT_WakeLevelList[level]
		if level == 0 {
			for i := 0; i < len(pWakeLevel.NeedItem); i++ {
				itemid := pWakeLevel.NeedItem[i]
				if itemid != 0 {
					pItemInfo := GetItemInfo(itemid)
					if pItemInfo == nil {
						return false
					}
					pWakeLevel.PropertyValues[0] += pItemInfo.Propertys[0]
					pWakeLevel.PropertyValues[1] += pItemInfo.Propertys[1]
					pWakeLevel.PropertyValues[2] += pItemInfo.Propertys[2]
					pWakeLevel.PropertyValues[3] += pItemInfo.Propertys[1]
					pWakeLevel.PropertyValues[4] += pItemInfo.Propertys[2]
				}
			}

			if pWakeLevel.ExtraProperty == AllPropertyID {
				pWakeLevel.PropertyPercents[0] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[1] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[2] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[3] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[4] += pWakeLevel.ExtraValue
			} else if pWakeLevel.ExtraProperty == AttackPropertyID {
				pWakeLevel.PropertyValues[AttackMagicID-1] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyValues[AttackPhysicID-1] += pWakeLevel.ExtraValue
			} else if pWakeLevel.ExtraProperty == DefencePropertyID {
				pWakeLevel.PropertyValues[DefenceMagicID-1] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyValues[DefencePhysicID-1] += pWakeLevel.ExtraValue
			} else if pWakeLevel.ExtraProperty > 0 && pWakeLevel.ExtraProperty <= 11 {
				pWakeLevel.PropertyValues[pWakeLevel.ExtraProperty-1] += pWakeLevel.ExtraValue
			}
		} else {
			pWakeLevel.PropertyValues = GT_WakeLevelList[level-1].PropertyValues
			pWakeLevel.PropertyPercents = GT_WakeLevelList[level-1].PropertyPercents
			for i := 0; i < len(pWakeLevel.NeedItem); i++ {
				itemid := pWakeLevel.NeedItem[i]
				if itemid != 0 {
					pItemInfo := GetItemInfo(itemid)
					pWakeLevel.PropertyValues[0] += pItemInfo.Propertys[0]
					pWakeLevel.PropertyValues[1] += pItemInfo.Propertys[1]
					pWakeLevel.PropertyValues[2] += pItemInfo.Propertys[2]
					pWakeLevel.PropertyValues[3] += pItemInfo.Propertys[1]
					pWakeLevel.PropertyValues[4] += pItemInfo.Propertys[2]
				}
			}

			if pWakeLevel.ExtraProperty == AllPropertyID {
				pWakeLevel.PropertyPercents[0] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[1] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[2] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[3] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyPercents[4] += pWakeLevel.ExtraValue
			} else if pWakeLevel.ExtraProperty == AttackPropertyID {
				pWakeLevel.PropertyValues[AttackMagicID-1] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyValues[AttackPhysicID-1] += pWakeLevel.ExtraValue
			} else if pWakeLevel.ExtraProperty == DefencePropertyID {
				pWakeLevel.PropertyValues[DefenceMagicID-1] += pWakeLevel.ExtraValue
				pWakeLevel.PropertyValues[DefencePhysicID-1] += pWakeLevel.ExtraValue
			} else if pWakeLevel.ExtraProperty > 0 && pWakeLevel.ExtraProperty <= 11 {
				pWakeLevel.PropertyValues[pWakeLevel.ExtraProperty-1] += pWakeLevel.ExtraValue
			}
		}
	}
	return true
}

func GetWakeLevelItem(level int) *ST_WakeLevel {
	if level >= len(GT_WakeLevelList) {
		gamelog.Error("GetWakeLevelItem Error : Invalid level  %d", level)
		return nil
	}

	return &GT_WakeLevelList[level]
}

func GetWakeComposeInfo(itemid int) *ST_WakeCompose {
	pWakeCompose, ok := GT_WakeComposeList[itemid]
	if pWakeCompose == nil || !ok {
		gamelog.Error("GetWakeComposeInfo Error: Invalid itemid :%d", itemid)
		return nil
	}

	return pWakeCompose
}

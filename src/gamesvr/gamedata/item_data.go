package gamedata

import (
	"gamelog"
)

type ST_ItemInfo struct {
	ItemID    int   //物品ID
	Type      int   //物品类型
	SubType   int   //物品子类型
	Quality   int   //物品的品质
	SellID    int   //出售货币ID
	SellPrice int   //出售价格
	UseType   int   //使用类型
	Data1     int   //参数1
	Data2     int   //参数2
	Propertys []int //觉醒道具的属性
}

var (
	GT_ItemList map[int]*ST_ItemInfo
)

func InitItemParser(total int) bool {

	GT_ItemList = make(map[int]*ST_ItemInfo)

	return true
}

func ParseItemRecord(rs *RecordSet) {
	ItemID := rs.GetFieldInt("id")
	iteminfo := new(ST_ItemInfo)

	iteminfo.ItemID = ItemID
	iteminfo.Type = rs.GetFieldInt("type")
	iteminfo.SubType = rs.GetFieldInt("sub_type")
	iteminfo.SellID = rs.GetFieldInt("sell_money_id")
	iteminfo.SellPrice = rs.GetFieldInt("sell_money_num")
	iteminfo.Quality = rs.GetFieldInt("quality")
	iteminfo.UseType = rs.GetFieldInt("usetype")
	iteminfo.Data2 = rs.GetFieldInt("data2")
	if iteminfo.Type != TYPE_WAKE {
		iteminfo.Data1 = rs.GetFieldInt("data1")
	} else {
		iteminfo.Propertys = ParseToIntSlice(rs.GetFieldString("data1"))
		if len(iteminfo.Propertys) < 3 {
			panic("field data1 is not a valid property string!")
		}
	}

	if iteminfo.Type == TYPE_MONEY {
		if iteminfo.Data1 <= 0 {
			panic("field data1 is not a valid money id!")
		}
	}

	if iteminfo.Type == TYPE_GEM {
		if iteminfo.Data1 <= 0 {
			panic("field data1 is not a valid gem id!")
		}
	}

	if iteminfo.Type == TYPE_HERO {
		if iteminfo.Data1 <= 0 {
			panic("field data1 is not a valid hero id!")
		}
	}

	if iteminfo.Type == TYPE_EQUIPMENT {
		if iteminfo.Data1 <= 0 {
			panic("field data1 is not a valid equip id!")
		}
	}

	if iteminfo.Type == TYPE_PET {
		if iteminfo.Data1 <= 0 {
			panic("field data1 is not a valid equip id!")
		}
	}

	if iteminfo.Type == TYPE_GEM_PIECE {
		if iteminfo.Data1 <= 0 {
			panic("field data1 is not a valid gem id!")
		}

		if iteminfo.Data2 <= 0 {
			panic("field data2 is not a valid piece index!")
		}
	}

	GT_ItemList[ItemID] = iteminfo

	return
}

func GetItemInfo(itemid int) *ST_ItemInfo {
	pItemInfo, ok := GT_ItemList[itemid]
	if pItemInfo == nil || !ok {
		gamelog.Error("GetItemInfo Error: Invalid itemid :%d", itemid)
		return nil
	}

	return pItemInfo
}

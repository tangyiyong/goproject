package gamedata

import (
	"gamelog"
)

const (
	StoreType_Hero  = 1
	StoreType_Awake = 2
	StoreType_Pet   = 3
)

//! 商店信息表
type ST_StoreInfo struct {
	ID        int //! 唯一标识
	Type      int //! 商店类型: 1->神将商店  2->觉醒商店  3->战宠商店
	ItemID    int //! 商品ID
	ItemNum   int //! 商品数量
	MoneyID   int //! 购买需求货币
	MoneyNum  int //! 价格
	NeedLevel int //! 需要等级
}

var GT_Store_List [4][]ST_StoreInfo

//! 初始化
func InitStoreParse(total int) bool {
	return true
}

//! 分析CSV
func ParseStoreRecord(rs *RecordSet) {
	storeType := rs.GetFieldInt("type")

	var info ST_StoreInfo
	id := rs.GetFieldInt("id")
	info.ID = id
	info.Type = storeType
	info.ItemID = rs.GetFieldInt("itemid")
	info.ItemNum = rs.GetFieldInt("itemnum")
	info.MoneyID = rs.GetFieldInt("moneytype")
	info.MoneyNum = rs.GetFieldInt("price")
	info.NeedLevel = rs.GetFieldInt("needlevel")

	GT_Store_List[storeType] = append(GT_Store_List[storeType], info)
}

//! 随机道具
func RandomStoreItem(number int, level int, itemtype int) (itemLst []ST_StoreInfo) {
	if itemtype > StoreType_Pet || itemtype < StoreType_Hero {
		gamelog.Error("RandomStoreItem error: invalid itemtype: %d", itemtype)
		return
	}

	itemCount := len(GT_Store_List[itemtype])
	if itemCount <= 0 {
		gamelog.Error("RandomStoreItem error: No Item data itemtype: %d", itemtype)
		return itemLst
	}
	for i := 0; i < number; i++ {
		id := r.Intn(itemCount)
		item := GT_Store_List[itemtype][id]
		itemLst = append(itemLst, item)
	}

	return itemLst
}

package gamedata

import (
	"gamelog"
)

const (
	Mall_Normal  = 0
	Mall_VipGift = 1
)

//! 商城信息表
type ST_MallItemInfo struct {
	ItemID  int32 //! 道具ID
	Type    int   //! 0-> 普通商品  1-> 礼包商品
	FuncID  int   //! 功能ID
	ItemNum int   //! 道具数量
	Value   int   //! 总共价值
}

var GT_MallItemLst []ST_MallItemInfo

//! 初始化
func InitMallParser(total int) bool {
	return true
}

//! 解析CSV
func ParseMallRecord(rs *RecordSet) {
	var mallItemInfo ST_MallItemInfo
	mallItemInfo.Type = rs.GetFieldInt("type")
	mallItemInfo.ItemID = int32(rs.GetFieldInt("itemid"))
	mallItemInfo.ItemNum = rs.GetFieldInt("itemnum")
	mallItemInfo.Value = rs.GetFieldInt("value")
	mallItemInfo.FuncID = rs.GetFieldInt("func_id")
	GT_MallItemLst = append(GT_MallItemLst, mallItemInfo)
}

func GetMallItemInfo(id int32) *ST_MallItemInfo {
	for i, v := range GT_MallItemLst {
		if v.ItemID == id {
			return &GT_MallItemLst[i]
		}
	}

	gamelog.Error("GetMallItemInfo Fail. id: %d", id)
	return nil
}

func GetVipGiftBagIDs() []int32 {
	item := []int32{}
	for _, v := range GT_MallItemLst {
		if v.Type == Mall_VipGift {
			item = append(item, int32(v.ItemID))
		}
	}
	return item
}

//! 根据物品获取功能ID
func GetFuncID(itemID int32) int {
	for _, v := range GT_MallItemLst {
		if v.ItemID == itemID {
			return v.FuncID
		}
	}
	return 0
}

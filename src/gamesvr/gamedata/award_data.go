package gamedata

import (
	"gamelog"
	"strings"
	"utility"
)

//奖励配制表

//获取物品数据
type ST_ItemData struct {
	ItemID  int
	ItemNum int
}
type IntPair struct {
	ID  int
	Cnt int
}

//掉落项
type ST_DropItem struct {
	ItemID  int    //物品ID
	ItemNum [2]int //物品数量
	Ratio   [2]int //概率范围
}

//奖励项
type ST_AwardItem struct {
	AwardID    int           //奖励ID
	FixItems   []ST_DropItem //必掉物品
	RatioCount int           //概率掉落个数
	RatioItems []ST_DropItem //机率掉落物品
}

var (
	GT_AwardList map[int]*ST_AwardItem = nil
)

func InitAwardParser(total int) bool {

	GT_AwardList = make(map[int]*ST_AwardItem) //奖励表

	return true
}

func ParseAwardRecord(rs *RecordSet) {
	awardID := CheckAtoi(rs.Values[0], 0)
	pAwardItem := new(ST_AwardItem)

	pAwardItem.AwardID = awardID
	pAwardItem.FixItems = nil
	pAwardItem.RatioItems = nil
	var bOk = false
	//解析固定掉落
	if rs.Values[2] != "NULL" {
		sFix := strings.Trim(rs.Values[2], "()")
		slice := strings.Split(sFix, ")(")
		itemCount := len(slice)
		pAwardItem.FixItems = make([]ST_DropItem, itemCount)
		for i, V := range slice {
			pAwardItem.FixItems[i], bOk = ParseToDropItem(V)
			if bOk == false {
				panic("field fix_item is wrong :s" + V)
			}
		}
	}

	pAwardItem.RatioCount = rs.GetFieldInt("ratio_num")

	if rs.Values[4] != "NULL" {
		var RatioBegin = 0
		sFix := strings.Trim(rs.Values[4], "()")
		slice := strings.Split(sFix, ")(")
		pAwardItem.RatioItems = make([]ST_DropItem, len(slice))
		for i, V := range slice {
			pAwardItem.RatioItems[i], bOk = ParseToDropItem(V)
			if bOk == false {
				panic("field ratio_item is wrong :" + V)
			}
			pAwardItem.RatioItems[i].Ratio[0] = RatioBegin + 1
			RatioBegin += pAwardItem.RatioItems[i].Ratio[1]
			pAwardItem.RatioItems[i].Ratio[1] = RatioBegin
		}
	}

	if len(pAwardItem.FixItems) <= 0 && len(pAwardItem.RatioItems) <= 0 {
		panic("field fix_item and ratio_item is wrong!")
	}

	GT_AwardList[awardID] = pAwardItem

}

func ParseToDropItem(drop string) (ST_DropItem, bool) {
	var item ST_DropItem
	pv := strings.Split(drop, "|")
	if len(pv) < 3 {
		return item, false
	}

	item.ItemID = CheckAtoi(pv[0], 10)
	numv := strings.Split(pv[1], "&")
	item.ItemNum[0] = CheckAtoi(numv[0], 11)
	item.ItemNum[1] = CheckAtoi(numv[1], 12)
	item.Ratio[1] = CheckAtoi(pv[2], 13)

	return item, true
}

func GetAwardItemByIndex(awardid int, index int) ST_ItemData {
	var item ST_ItemData
	pAwardItem, ok := GT_AwardList[awardid]
	if pAwardItem == nil || !ok {
		gamelog.Error("GetItemByIndex Error: Invalid awardid :%d", awardid)
		return item
	}

	if index >= len(pAwardItem.FixItems) {
		gamelog.Error("GetItemByIndex Error: Invalid index :%d", index)
		return item
	}

	item.ItemID = pAwardItem.FixItems[index].ItemID
	item.ItemNum = pAwardItem.FixItems[index].ItemNum[0]
	return item
}

func GetItemsFromAwardIDEx(awardid int) []ST_ItemData {
	pAwardItem, ok := GT_AwardList[awardid]
	if pAwardItem == nil || !ok {
		gamelog.Error("GetItemsFromAwardID Error: Invalid awardid :%d", awardid)
		return nil
	}

	var ret []ST_ItemData
	var item ST_ItemData
	if pAwardItem.FixItems != nil {
		for _, v := range pAwardItem.FixItems {
			item.ItemID = v.ItemID
			if v.ItemNum[0] == v.ItemNum[1] {
				item.ItemNum = v.ItemNum[0]
			} else {
				item.ItemNum = v.ItemNum[0] + utility.Rand()%(v.ItemNum[1]-v.ItemNum[0]+1)
			}

			if item.ItemNum > 0 {
				ret = append(ret, item)
			}

		}
	}

	if pAwardItem.RatioItems != nil {
		for i := 0; i < pAwardItem.RatioCount; i++ {
			randvalue := utility.Rand()
			for _, v := range pAwardItem.RatioItems {
				if (randvalue >= v.Ratio[0]) && (randvalue <= v.Ratio[1]) {
					item.ItemID = v.ItemID
					if v.ItemNum[1] == v.ItemNum[0] {
						item.ItemNum = v.ItemNum[0]
					} else {
						item.ItemNum = v.ItemNum[0] + utility.Rand()%(v.ItemNum[1]-v.ItemNum[0]+1)
					}

					if item.ItemNum > 0 {
						ret = append(ret, item)
					}
				}
			}
		}
	}
	return ret
}

//同一个掉落项只能掉一次，就是进行了去重
func GetItemsFromAwardID(awardid int) []ST_ItemData {
	pAwardItem, ok := GT_AwardList[awardid]
	if pAwardItem == nil || !ok {
		gamelog.Error("GetItemsFromAwardID Error: Invalid awardid :%d", awardid)
		return nil
	}

	var ret []ST_ItemData
	var item ST_ItemData
	if pAwardItem.FixItems != nil {
		for _, v := range pAwardItem.FixItems {
			item.ItemID = v.ItemID
			if v.ItemNum[0] == v.ItemNum[1] {
				item.ItemNum = v.ItemNum[0]
			} else {
				item.ItemNum = v.ItemNum[0] + utility.Rand()%(v.ItemNum[1]-v.ItemNum[0]+1)
			}

			if item.ItemNum > 0 {
				ret = append(ret, item)
			}

		}
	}

	if pAwardItem.RatioItems != nil {
		for i := 0; i < pAwardItem.RatioCount; i++ {
			randvalue := utility.Rand()
			for _, v := range pAwardItem.RatioItems {
				if (randvalue >= v.Ratio[0]) && (randvalue <= v.Ratio[1]) {
					item.ItemID = v.ItemID
					if v.ItemNum[1] == v.ItemNum[0] {
						item.ItemNum = v.ItemNum[0]
					} else {
						item.ItemNum = v.ItemNum[0] + utility.Rand()%(v.ItemNum[1]-v.ItemNum[0]+1)
					}

					if item.ItemNum > 0 {
						has := false
						for _, t := range ret {
							if t.ItemID == item.ItemID {
								has = true
							}
						}
						if has == false {
							ret = append(ret, item)
						}
					}
				}
			}
		}
	}
	return ret
}

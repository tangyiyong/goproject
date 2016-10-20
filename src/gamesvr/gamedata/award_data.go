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
	Ratio   int    //概率范围
}

//奖励项
type ST_AwardItem struct {
	AwardID    int           //奖励ID
	FixItems   []ST_DropItem //必掉物品
	RatioCount int           //概率掉落个数
	Distinct   int           //是否需要去重
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
		var RatioBegin = 1
		var tempvalue = 0
		sFix := strings.Trim(rs.Values[4], "()")
		slice := strings.Split(sFix, ")(")
		pAwardItem.RatioItems = make([]ST_DropItem, len(slice)+1)
		for i, V := range slice {
			pAwardItem.RatioItems[i], bOk = ParseToDropItem(V)
			if bOk == false {
				panic("field ratio_item is wrong" + V)
			}
			tempvalue = pAwardItem.RatioItems[i].Ratio
			pAwardItem.RatioItems[i].Ratio = RatioBegin
			RatioBegin += tempvalue
		}

		pAwardItem.RatioItems[len(slice)].ItemID = 0
		pAwardItem.RatioItems[len(slice)].Ratio = 10000
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
	item.Ratio = CheckAtoi(pv[2], 13)

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

func GetItemsFromAwardID(awardid int) []ST_ItemData {
	pAwardItem, ok := GT_AwardList[awardid]
	if pAwardItem == nil || !ok {
		gamelog.Error("GetItemsFromAwardID Error: Invalid awardid :%d", awardid)
		return nil
	}

	var ret []ST_ItemData
	var item ST_ItemData
	if pAwardItem.FixItems != nil {
		for i := 0; i < len(pAwardItem.FixItems); i++ {
			item.ItemID = pAwardItem.FixItems[i].ItemID
			if pAwardItem.FixItems[i].ItemNum[0] == pAwardItem.FixItems[i].ItemNum[1] {
				item.ItemNum = pAwardItem.FixItems[i].ItemNum[0]
			} else {
				item.ItemNum = pAwardItem.FixItems[i].ItemNum[0] +
					utility.Rand()%(pAwardItem.FixItems[i].ItemNum[1]-pAwardItem.FixItems[i].ItemNum[0]+1)
			}

			if item.ItemNum > 0 {
				ret = append(ret, item)
			}
		}
	}

	if pAwardItem.RatioItems != nil {
		for cycle := 0; cycle < pAwardItem.RatioCount; cycle++ {
			randvalue := utility.Rand()
			for i := 0; i < (len(pAwardItem.RatioItems) - 1); i++ {
				if (randvalue >= pAwardItem.RatioItems[i].Ratio) && (randvalue < pAwardItem.RatioItems[i+1].Ratio) {
					item.ItemID = pAwardItem.RatioItems[i].ItemID
					if pAwardItem.RatioItems[i].ItemNum[1] == pAwardItem.RatioItems[i].ItemNum[0] {
						item.ItemNum = pAwardItem.RatioItems[i].ItemNum[0]
					} else {
						item.ItemNum = pAwardItem.RatioItems[i].ItemNum[0] + utility.Rand()%(pAwardItem.RatioItems[i].ItemNum[1]-pAwardItem.RatioItems[i].ItemNum[0]+1)
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

func GetItemsAwardIDTimes(awardid int, times int) []ST_ItemData {
	pAwardItem, ok := GT_AwardList[awardid]
	if pAwardItem == nil || !ok {
		gamelog.Error("GetItemsFromAwardID Error: Invalid awardid :%d", awardid)
		return nil
	}

	var ret []ST_ItemData
	var item ST_ItemData
	if pAwardItem.FixItems != nil {
		for i := 0; i < len(pAwardItem.FixItems); i++ {
			item.ItemID = pAwardItem.FixItems[i].ItemID
			if pAwardItem.FixItems[i].ItemNum[0] == pAwardItem.FixItems[i].ItemNum[1] {
				item.ItemNum = pAwardItem.FixItems[i].ItemNum[0]
			} else {
				item.ItemNum = pAwardItem.FixItems[i].ItemNum[0] +
					utility.Rand()%(pAwardItem.FixItems[i].ItemNum[1]-pAwardItem.FixItems[i].ItemNum[0]+1)
			}

			if item.ItemNum > 0 {
				ret = append(ret, item)
			}
		}
	}

	if pAwardItem.RatioItems != nil {
		var tempCount = pAwardItem.RatioCount
		tempCount *= times
		for cycle := 0; cycle < tempCount; cycle++ {
			randvalue := utility.Rand()
			for i := 0; i < (len(pAwardItem.RatioItems) - 1); i++ {
				if (randvalue >= pAwardItem.RatioItems[i].Ratio) && (randvalue < pAwardItem.RatioItems[i+1].Ratio) {
					item.ItemID = pAwardItem.RatioItems[i].ItemID
					if pAwardItem.RatioItems[i].ItemNum[1] == pAwardItem.RatioItems[i].ItemNum[0] {
						item.ItemNum = pAwardItem.RatioItems[i].ItemNum[0]
					} else {
						item.ItemNum = pAwardItem.RatioItems[i].ItemNum[0] + utility.Rand()%(pAwardItem.RatioItems[i].ItemNum[1]-pAwardItem.RatioItems[i].ItemNum[0]+1)
					}

					if item.ItemNum > 0 {
						ret = append(ret, item)
					}

					//if item.ItemNum > 0 {
					//	has := false
					//	for _, t := range ret {
					//		if t.ItemID == item.ItemID {
					//			has = true
					//		}
					//	}
					//	if has == false {
					//		ret = append(ret, item)
					//	}
					//}
				}
			}
		}
	}
	return ret
}

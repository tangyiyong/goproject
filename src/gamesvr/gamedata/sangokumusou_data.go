package gamedata

import (
	"gamelog"
)

//! 三国无双章节
type ST_SangokuMusou_Chapter struct {
	ChapterID      int //! 章节ID
	CopyID         int //! 关卡ID
	Condition      int //! 章节条件
	ConditionValue int //! 章节条件数值
	Diffculty1     int //! 不同难度对应奖励
	Diffculty2     int //! 不同难度对应奖励
	Diffculty3     int //! 不同难度对应奖励
}

//! 三国无双章节奖励
type ST_SangokuMusou_Chapter_Award struct {
	ID      int //! 唯一标识
	Chapter int //! 章节ID
	StarNum int //! 星数
	Award   int //! 奖励
}

//! 三国无双属性加成
type ST_SangokuMusou_Attr_Markup struct {
	AttrID   int //! 属性
	Value    int //! 加成值
	CostStar int //! 消耗星数
}

//! 三国无双精英挑战
type ST_SangokuMusou_Elite_Copy struct {
	EliteID      int //! 唯一标识
	CopyID       int
	NeedPassCopy int //! 需求通关
}

//! 三国无双神装商店
type ST_SangokuMusou_Store struct {
	ID            int //! 唯一标识
	ItemID        int //! 物品ID
	ItemNum       int //! 物品数量
	CostMoneyType int //! 消耗金钱类型
	CostMoneyNum  int //! 消耗金钱数量
	CostItemType  int //! 消耗物品类型
	CostItemNum   int //! 消耗物品数量
	NeedStar      int //! 需要星数
	NeedLevel     int //! 购买需求等级
	BuyTimes      int //! 每日购买次数
	ItemType      int //! 分类(0-商品 1-紫装 2-橙装 3-红装 4-奖励)
}

//! 无双迷藏
type ST_SangokuMusou_Sale struct {
	ID            int //! 唯一标识
	ItemID        int //! 物品ID
	ItemNum       int //! 物品数量
	CostMoneyType int //! 消耗金钱类型
	CostMoneyNum  int //! 消耗金钱数量
	OriginalPrice int //! 原价
	RangeStarMin  int //! 星数取值范围
	RangeStarMax  int //! 星数取值范围
}

var (
	GT_SangokuMusou_Chapter              map[int]*ST_SangokuMusou_Chapter //! 三国无双章节
	GT_SangokuMusou_Chapter_Award        []ST_SangokuMusou_Chapter_Award  //! 三国无双章节奖励
	GT_SangokuMusou_Attr_Markup_Normal   []ST_SangokuMusou_Attr_Markup    //! 三国无双三星属性加成
	GT_SangokuMusou_Attr_Markup_Senior   []ST_SangokuMusou_Attr_Markup    //! 三国无双六星属性加成
	GT_SangokuMusou_Attr_Markup_Ultimate []ST_SangokuMusou_Attr_Markup    //! 三国无双九星属性加成
	GT_SangokuMusou_Elite_Copy           []ST_SangokuMusou_Elite_Copy     //! 三国无双精英挑战
	GT_SangokuMusou_Store                []ST_SangokuMusou_Store          //! 三国无双神装商店
	GT_SangokuMusou_Sale                 []ST_SangokuMusou_Sale           //! 三国无双秘藏
)

func InitSangokuMusouChapter(total int) bool {
	GT_SangokuMusou_Chapter = make(map[int]*ST_SangokuMusou_Chapter)
	return true
}

func InitSangokuMusouChapterAwardParser(total int) bool {
	GT_SangokuMusou_Chapter_Award = make([]ST_SangokuMusou_Chapter_Award, total+1)
	return true
}

func InitSangokuMusouAttrMarkupParser(total int) bool {
	return true
}

func InitSangokuMusouEliteCopyParser(total int) bool {
	GT_SangokuMusou_Elite_Copy = make([]ST_SangokuMusou_Elite_Copy, total+1)
	return true
}

func InitSangokuMusouStoreParser(total int) bool {
	GT_SangokuMusou_Store = make([]ST_SangokuMusou_Store, total+1)
	return true
}

func InitSangokuMusouSaleParser(total int) bool {
	GT_SangokuMusou_Sale = make([]ST_SangokuMusou_Sale, total+1)
	return true
}

func ParseSangokuMusouChapterRecord(rs *RecordSet) {
	copyID := rs.GetFieldInt("copyid")

	copyInfo := new(ST_SangokuMusou_Chapter)
	copyInfo.ChapterID = rs.GetFieldInt("chapterid")
	copyInfo.CopyID = rs.GetFieldInt("copyid")
	copyInfo.Diffculty1 = rs.GetFieldInt("difficulty1")
	copyInfo.Diffculty2 = rs.GetFieldInt("difficulty2")
	copyInfo.Diffculty3 = rs.GetFieldInt("difficulty3")
	copyInfo.Condition = rs.GetFieldInt("condition")
	copyInfo.ConditionValue = rs.GetFieldInt("conditionvalue")

	GT_SangokuMusou_Chapter[copyID] = copyInfo
}

func ParseSangokuMusouChapterAwardRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_SangokuMusou_Chapter_Award[id].ID = id
	GT_SangokuMusou_Chapter_Award[id].Chapter = rs.GetFieldInt("chapter")
	GT_SangokuMusou_Chapter_Award[id].StarNum = rs.GetFieldInt("starnum")
	GT_SangokuMusou_Chapter_Award[id].Award = rs.GetFieldInt("award")
}

func ParseSangokuMusouAttrMarkupRecord(rs *RecordSet) {
	costStar := rs.GetFieldInt("coststar")

	if costStar == 3 { //! 三星属性奖励
		var normal ST_SangokuMusou_Attr_Markup
		normal.AttrID = rs.GetFieldInt("attrid")
		normal.Value = rs.GetFieldInt("value")
		normal.CostStar = costStar
		GT_SangokuMusou_Attr_Markup_Normal = append(GT_SangokuMusou_Attr_Markup_Normal, normal)
	} else if costStar == 6 { //! 六星属性奖励
		var senior ST_SangokuMusou_Attr_Markup
		senior.AttrID = rs.GetFieldInt("attrid")
		senior.Value = rs.GetFieldInt("value")
		senior.CostStar = costStar

		GT_SangokuMusou_Attr_Markup_Senior = append(GT_SangokuMusou_Attr_Markup_Senior, senior)
	} else if costStar == 9 { //! 九星属性奖励
		var ultimate ST_SangokuMusou_Attr_Markup
		ultimate.AttrID = rs.GetFieldInt("attrid")
		ultimate.Value = rs.GetFieldInt("value")
		ultimate.CostStar = costStar
		GT_SangokuMusou_Attr_Markup_Ultimate = append(GT_SangokuMusou_Attr_Markup_Ultimate, ultimate)
	}
}

func ParseSangokuMusouEliteCopyRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_SangokuMusou_Elite_Copy[id].EliteID = id
	GT_SangokuMusou_Elite_Copy[id].CopyID = rs.GetFieldInt("copyid")
	GT_SangokuMusou_Elite_Copy[id].NeedPassCopy = rs.GetFieldInt("needpasscopy")
}

func ParseSangokuMusouStoreRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_SangokuMusou_Store[id].ID = id
	GT_SangokuMusou_Store[id].ItemID = rs.GetFieldInt("itemid")
	GT_SangokuMusou_Store[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_SangokuMusou_Store[id].CostMoneyType = rs.GetFieldInt("costmoneytype")
	GT_SangokuMusou_Store[id].CostMoneyNum = rs.GetFieldInt("costmoneynum")
	GT_SangokuMusou_Store[id].CostItemType = rs.GetFieldInt("costitemid")
	GT_SangokuMusou_Store[id].CostItemNum = rs.GetFieldInt("costitemnum")
	GT_SangokuMusou_Store[id].NeedStar = rs.GetFieldInt("needstar")
	GT_SangokuMusou_Store[id].NeedLevel = rs.GetFieldInt("needlevel")
	GT_SangokuMusou_Store[id].BuyTimes = rs.GetFieldInt("buytimes")
	GT_SangokuMusou_Store[id].ItemType = rs.GetFieldInt("itemtype")
}

func ParseSangokuMusouSaleRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	GT_SangokuMusou_Sale[id].ID = id
	GT_SangokuMusou_Sale[id].ItemID = rs.GetFieldInt("itemid")
	GT_SangokuMusou_Sale[id].ItemNum = rs.GetFieldInt("itemnum")
	GT_SangokuMusou_Sale[id].CostMoneyType = rs.GetFieldInt("costmoneytype")
	GT_SangokuMusou_Sale[id].CostMoneyNum = rs.GetFieldInt("costmoneynum")
	GT_SangokuMusou_Sale[id].OriginalPrice = rs.GetFieldInt("originalprice")
	GT_SangokuMusou_Sale[id].RangeStarMin = rs.GetFieldInt("rangestarmin")
	GT_SangokuMusou_Sale[id].RangeStarMax = rs.GetFieldInt("rangestarmax")
}

//! 获取三国无双章节信息
func GetSangokuMusouChapterInfo(copyID int) *ST_SangokuMusou_Chapter {
	_, ok := GT_SangokuMusou_Chapter[copyID]
	if ok != true {
		gamelog.Error("GetSangokuMusouChapterInfo Error: Invalid copyID %d", copyID)
		return nil
	}

	return GT_SangokuMusou_Chapter[copyID]
}

//! 获取章节关卡
func GetSGWSChapterCopyLst(chapter int) []int {
	copyLst := []int{}
	for _, v := range GT_SangokuMusou_Chapter {
		if v.ChapterID == chapter {
			copyLst = append(copyLst, v.CopyID)
		}
	}
	return copyLst
}

//! 获取三国无双章节信息
// func GetSangokuMusouChapterInfoFromCopy(copyID int) *ST_SangokuMusou_Chapter {
// 	for i, v := range GT_SangokuMusou_Chapter {
// 		for _, n := range v.CopyID {
// 			if n == copyID {
// 				return &GT_SangokuMusou_Chapter[i]
// 			}
// 		}
// 	}
// 	gamelog.Error("GetSangokuMusouChapterInfoFromCopy Error: invalid copyID :%d", copyID)
// 	return nil
// }

//! 获取三国无双商品信息
func GetSangokumusouStoreInfo(id int) *ST_SangokuMusou_Store {
	if id >= len(GT_SangokuMusou_Store) || id <= 0 {
		gamelog.Error("GetSangokumusouStoreInfo Error: invalid id %d", id)
		return nil
	}

	return &GT_SangokuMusou_Store[id]
}

//! 获取三国无双精英挑战信息
func GetSangokuMusouEliteCopyInfo(copyID int) *ST_SangokuMusou_Elite_Copy {
	for i, v := range GT_SangokuMusou_Elite_Copy {
		if v.CopyID == copyID {
			return &GT_SangokuMusou_Elite_Copy[i]
		}
	}

	gamelog.Error("GetSangokuMusouEliteCopyInfo error: not find copyID: %d", copyID)

	return nil
}

//! 获取章节末副本ID
func GetSGWSChapterEndCopyID(chapter int) int {
	endCopy := 0
	for _, v := range GT_SangokuMusou_Chapter {
		if v.ChapterID == chapter && v.CopyID > endCopy {
			endCopy = v.CopyID
		}
	}
	return endCopy
}

//! 获取三国无双章节起始ID
func GetSGWSChapterBeginCopyID(chapter int) int {
	beginCopy := 0
	for _, v := range GT_SangokuMusou_Chapter {
		if v.ChapterID == chapter {

			if beginCopy == 0 {
				beginCopy = v.CopyID
			}

			if v.CopyID < beginCopy {
				beginCopy = v.CopyID
			}
		}
	}
	return beginCopy
}

//! 获取章节数
func GetSGWSChapterCount() int {
	count := 0
	value := 0
	for _, v := range GT_SangokuMusou_Chapter {
		if v.ChapterID != value {
			value = v.ChapterID
			count++
		}
	}
	return count
}

//! 判断是否为章节末
func SangokuMusou_IsChapterEnd(copyID int) (bool, int) {

	count := GetSGWSChapterCount()
	for i := 1; i <= count; i++ {
		endcopy := GetSGWSChapterEndCopyID(i)
		if copyID == endcopy {
			return true, i
		}
	}

	return false, 0
}

//! 获取下一关卡
func SangokuMusou_GetNextChapter(chapter int) int {
	if chapter >= GetSGWSChapterCount() {
		return 0
	}
	return chapter + 1
}

//! 获取三国无双关卡下一关ID
func GetSangokuMusouNextCopy(copyID int) int {
	if copyID == 0 {
		return GetSGWSChapterBeginCopyID(1)
	}

	isEnd, chapter := SangokuMusou_IsChapterEnd(copyID)
	if isEnd == true {
		nextChapterID := SangokuMusou_GetNextChapter(chapter)
		if nextChapterID == 0 {
			return 0
		}

		return GetSGWSChapterBeginCopyID(nextChapterID)
	}

	return copyID + 1
}

//! 获取章节奖励信息
func GetSangokuMusouChapterAwardInfo(chapter int, starnum int) *ST_SangokuMusou_Chapter_Award {
	for i, v := range GT_SangokuMusou_Chapter_Award {
		if v.Chapter == chapter && v.StarNum == starnum {
			return &GT_SangokuMusou_Chapter_Award[i]
		}
	}

	gamelog.Error("Get San Guo Wu Shuang Chapter Award fail. chapter:%d star:%d", chapter, starnum)
	return nil
}

//! 随机属性加成
func RandSangokuMusouAttrMarkup() (attrLst []*ST_SangokuMusou_Attr_Markup) {

	//! 从三星属性加成随机一个
	index := r.Intn(len(GT_SangokuMusou_Attr_Markup_Normal))
	attrLst = append(attrLst, &GT_SangokuMusou_Attr_Markup_Normal[index])

	//! 从六星属性加成随机一个
	index = r.Intn(len(GT_SangokuMusou_Attr_Markup_Senior))
	attrLst = append(attrLst, &GT_SangokuMusou_Attr_Markup_Senior[index])

	//! 从九星属性加成随机一个
	index = r.Intn(len(GT_SangokuMusou_Attr_Markup_Ultimate))
	attrLst = append(attrLst, &GT_SangokuMusou_Attr_Markup_Ultimate[index])
	return attrLst
}

//! 随机无双秘藏
func RandMusouTreasure(starNum int) int {
	itemLst := []ST_SangokuMusou_Sale{}
	for _, v := range GT_SangokuMusou_Sale {
		if starNum >= v.RangeStarMin && starNum < v.RangeStarMax {
			itemLst = append(itemLst, v)
		}
	}

	index := r.Intn(len(itemLst))
	if len(itemLst) == 0 {
		return 0
	}

	return GT_SangokuMusou_Sale[itemLst[index].ID].ID
}

//! 获取无双秘藏
func GetMusouTreasure(id int) *ST_SangokuMusou_Sale {
	if id >= len(GT_SangokuMusou_Sale) || id <= 0 {
		gamelog.Error("GetMusouTreasure Error: invalid id %d", id)
		return nil
	}

	return &GT_SangokuMusou_Sale[id]
}

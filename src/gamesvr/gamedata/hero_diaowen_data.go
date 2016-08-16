package gamedata

import (
	"gamelog"
)

type ST_DiaoWenItem struct {
	Quality      int         //品质
	CostMoneyID  int         //消耗货币ID
	CostMoneyNum int         //消耗货币数量
	NeedLevel    int         //需要等级
	Propertys    [5]ST_Range //五个基本属性的范围
	NeedCulture  [5]int      //需要的培养值
}

type ST_HeroDiaoWen struct {
	ID    int
	Items [10]ST_DiaoWenItem //雕文列表
}

type ST_XiLianCostItem struct {
	LockNum        int //锁定数量
	FirstMoneyID   int //第一货币ID
	FirstMoneyNum  int //第一货币数量
	SecondMoneyID  int //第二货币ID
	SecondMoneyNum int //第二货币数量
}

var (
	GT_DiaoWenList []ST_HeroDiaoWen    //雕文表
	GT_XiLianList  []ST_XiLianCostItem //雕文洗炼表
)

func InitDiaoWenParser(total int) bool {
	GT_DiaoWenList = make([]ST_HeroDiaoWen, 10)
	return true
}

func InitXiLianParser(total int) bool {
	GT_XiLianList = make([]ST_XiLianCostItem, total+1)
	return true
}

func ParseDiaoWenRecord(rs *RecordSet) {
	id := rs.GetFieldInt("id")
	quality := rs.GetFieldInt("quality")
	GT_DiaoWenList[id].Items[quality].CostMoneyID = rs.GetFieldInt("cost_money_id")
	GT_DiaoWenList[id].Items[quality].CostMoneyNum = rs.GetFieldInt("cost_money_num")
	GT_DiaoWenList[id].Items[quality].NeedLevel = rs.GetFieldInt("needlevel")
	GT_DiaoWenList[id].Items[quality].Propertys[0].Value = ParseTo2IntSlice(rs.GetFieldString("p1_range"))
	GT_DiaoWenList[id].Items[quality].Propertys[1].Value = ParseTo2IntSlice(rs.GetFieldString("p2_range"))
	GT_DiaoWenList[id].Items[quality].Propertys[2].Value = ParseTo2IntSlice(rs.GetFieldString("p3_range"))
	GT_DiaoWenList[id].Items[quality].Propertys[3].Value = ParseTo2IntSlice(rs.GetFieldString("p4_range"))
	GT_DiaoWenList[id].Items[quality].Propertys[4].Value = ParseTo2IntSlice(rs.GetFieldString("p5_range"))
	GT_DiaoWenList[id].Items[quality].NeedCulture[0] = rs.GetFieldInt("need_p1")
	GT_DiaoWenList[id].Items[quality].NeedCulture[1] = rs.GetFieldInt("need_p2")
	GT_DiaoWenList[id].Items[quality].NeedCulture[2] = rs.GetFieldInt("need_p3")
	GT_DiaoWenList[id].Items[quality].NeedCulture[3] = rs.GetFieldInt("need_p4")
	GT_DiaoWenList[id].Items[quality].NeedCulture[4] = rs.GetFieldInt("need_p5")

}

func ParseXiLianRecord(rs *RecordSet) {
	locknum := rs.GetFieldInt("lock_num")
	GT_XiLianList[locknum].LockNum = locknum
	GT_XiLianList[locknum].FirstMoneyID = rs.GetFieldInt("money_id_1")
	GT_XiLianList[locknum].FirstMoneyNum = rs.GetFieldInt("money_num_1")
	GT_XiLianList[locknum].SecondMoneyID = rs.GetFieldInt("money_id_2")
	GT_XiLianList[locknum].SecondMoneyNum = rs.GetFieldInt("money_num_2")
}

func GetDiaoWenInfo(id int, quality int) *ST_DiaoWenItem {
	if id <= 0 || quality <= 0 {
		gamelog.Error("GetXiLianInfo Error : Invalid id:%d and Invalid quality:%d", id, quality)
		return nil
	}

	return &GT_DiaoWenList[id].Items[quality]
}

func GetXiLianInfo(locknum int) *ST_XiLianCostItem {
	if locknum >= len(GT_XiLianList) {
		gamelog.Error("GetXiLianInfo Error : Invalid locknum  %d", locknum)
		return nil
	}

	return &GT_XiLianList[locknum]
}

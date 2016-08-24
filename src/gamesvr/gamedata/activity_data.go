package gamedata

import (
	"gamelog"
	"time"
)

// type_activity_type.csv
const (
	Activity_Sign             = 1  //! 签到
	Activity_Vip_Gift         = 2  //! VIP礼包
	Activity_Login            = 3  //! 登录奖励
	Activity_Recv_Action      = 4  //! 领取体力
	Activity_Week_Award       = 5  //! 超级周周盈
	Activity_Money_Gold       = 6  //! 财神
	Activity_Recharge_Gift    = 7  //! 累充礼包
	Activity_Open_Fund        = 8  //! 开服基金
	Activity_Discount_Sale    = 9  //! 折扣
	Activity_First_Recharge   = 10 //! 首充
	Activity_Moon_Card        = 11 //! 月卡
	Activity_Singel_Recharge  = 12 //! 单充返利
	Activity_Limit_Daily_Task = 13 //! 限时日常
	Activity_Hunt_Treasure    = 14 //! 巡回探宝
	Activity_Card_Master      = 15 //! 卡牌大师
	Activity_MoonlightShop    = 16 //! 异域商人
	Activity_Luckly_Wheel     = 17 //! 幸运轮盘
	Activity_Beach_Baby       = 18 //! 沙滩宝贝
	Activity_Group_Purchase   = 19 //! 团购
	Activity_Festival         = 20 //! 欢庆佳节
	Activity_Competition      = 21 //! 七天战力排行
	Activity_Level_Gift       = 22 //! 等级礼包
	Activity_Month_Fund       = 23 //! 月基金
	Activity_LimitSale        = 24 //! 限时优惠
	Activity_Rank_Sale        = 25 //! 巅峰特惠
	Activity_Seven            = 26 //! 七天活动
)

//! 活动配置表
type ST_ActivityData struct {
	ID           int
	Name         string
	ServerType   int //! 1->新服 2->公服 3->全期存在
	BeginTime    int
	EndTime      int
	TimeType     int //! 活动时间 1->月 2->周 3->开服 4->指定日期
	AwardTime    int //! 活动奖励时间
	ActivityType int //! 活动套用模板
	AwardType    int //! 活动奖励套用模板
	Status       int //! 开启状态
	Icon         int //! ICON
	Inside       int //! 1->里面 2->外面 3->同时存在(需判断持续days)
	Days         int //! 临时存在天数
}

var GT_ActivityLst map[int]ST_ActivityData

func InitActivityParser(total int) bool {
	GT_ActivityLst = make(map[int]ST_ActivityData)
	return true
}

func ParseActivityRecord(rs *RecordSet) {
	id := CheckAtoi(rs.Values[0], 0)
	if id <= 0 {
		gamelog.Error("ParseActivityRecord Error: invalid id %d", id)
		return
	}

	data := ST_ActivityData{}
	data.ID = id
	data.Name = rs.GetFieldString("name")
	data.ServerType = rs.GetFieldInt("activity_type")
	data.BeginTime = rs.GetFieldInt("begintime")
	data.EndTime = rs.GetFieldInt("endtime")
	data.AwardTime = rs.GetFieldInt("awardtime")
	data.TimeType = rs.GetFieldInt("timetype")
	data.ActivityType = rs.GetFieldInt("type")
	data.AwardType = rs.GetFieldInt("award_type")
	data.Status = rs.GetFieldInt("status")
	data.Icon = rs.GetFieldInt("icon")
	data.Inside = rs.GetFieldInt("inside")
	data.Days = rs.GetFieldInt("days")
	GT_ActivityLst[id] = data
}

func GetActivityInfo(id int) *ST_ActivityData {
	data, ok := GT_ActivityLst[id]
	if ok == false {
		gamelog.Error("GetActivityInfo Error: Can't Find %d", id)
		return nil
	}

	return &data
}

func GetActivityNextBeginTime(activityID int, openDay int) (beginTime int64, endTime int64) {
	activityInfo := GetActivityInfo(activityID)

	now := time.Now()
	if activityInfo.TimeType == 1 { //! 按照月计算
		beginDate := time.Date(now.Year(), now.Month(), activityInfo.BeginTime, 0, 0, 0, 0, now.Location())
		endDate := time.Date(now.Year(), now.Month(), activityInfo.EndTime, 23, 59, 59, 59, now.Location())
		beginDate = beginDate.AddDate(0, 1, 0)
		endDate = endDate.AddDate(0, 1, 0)
		beginTime = beginDate.Unix()
		endTime = endDate.Unix()
	} else if activityInfo.TimeType == 2 { //! 按照日计算
		weekDay := int(now.Weekday())
		if weekDay == 0 { //! 特殊处理周末
			weekDay = 7
		}

		if activityInfo.EndTime == 7 && activityInfo.BeginTime == 1 {
			//! 永久活动
			beginTime = 0
			endTime = 0
		} else {
			beginDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			beginDate = beginDate.AddDate(0, 0, activityInfo.BeginTime-weekDay)
			beginDate = beginDate.AddDate(0, 0, 7)
			beginTime = beginDate.Unix()

			endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())
			endDate = endDate.AddDate(0, 0, activityInfo.EndTime-weekDay)
			endDate = endDate.AddDate(0, 0, 7)
			endTime = endDate.Unix()
		}
	} else if activityInfo.TimeType == 3 {
		return -1, -1 //! 开服活动无下次开启时间
	} else if activityInfo.TimeType == 4 { //! 按照 月*100+日格式写 比如 310 = 3月10日
		day := activityInfo.BeginTime % 100
		month := (activityInfo.BeginTime - day) / 100
		if day < 1 || day > 31 || month < 1 || month > 12 {
			gamelog.Error("Invalid Activity BeginTime: %d", activityInfo.BeginTime)
			return -1, -1
		}

		beginData := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, now.Location())
		beginData = beginData.AddDate(1, 0, 0)
		beginTime = beginData.Unix()

		day = activityInfo.EndTime % 100
		month = (activityInfo.EndTime - day) / 100
		if day < 1 || day > 31 || month < 1 || month > 12 {
			gamelog.Error("Invalid Activity BeginTime: %d", activityInfo.BeginTime)
			return -1, -1
		}

		endData := time.Date(now.Year(), time.Month(month), day, 23, 59, 59, 59, now.Location())
		endData = endData.AddDate(1, 0, 0)
		endTime = endData.Unix()
	}

	return beginTime, endTime
}

func GetNewServerTypeActivityEndTime(activityInfo *ST_ActivityData) int {
	endDay := 0
	for _, v := range GT_ActivityLst {
		if activityInfo.ActivityType == v.ActivityType &&
			v.Status == 1 &&
			v.ServerType == 1 &&
			v.EndTime > endDay {
			endDay = v.EndTime
		}
	}

	return endDay
}

func GetActivityEndTime(activityID int, openDay int) (beginTime int64, endTime int64) {
	activityInfo := GetActivityInfo(activityID)
	if activityInfo == nil {
		gamelog.Error("GetActivityEndTime Error : Invalid activityid:%d", activityID)
		return
	}

	if activityInfo.BeginTime == 0 && activityInfo.EndTime == 0 {
		return 0, 0
	}

	now := time.Now()
	if activityInfo.TimeType == 1 { //! 按照月计算
		beginDate := time.Date(now.Year(), now.Month(), activityInfo.BeginTime, 0, 0, 0, 0, now.Location())
		endDate := time.Date(now.Year(), now.Month(), activityInfo.EndTime, 23, 59, 59, 59, now.Location())

		if endDate.Unix() <= now.Unix() {
			beginDate = beginDate.AddDate(0, 1, 0)
			endDate = endDate.AddDate(0, 1, 0)
		}

		if activityInfo.ServerType == 2 { //! 公服期活动
			endDay := GetNewServerTypeActivityEndTime(activityInfo)
			if endDay > openDay {
				date := time.Now()
				date = date.AddDate(0, 0, endDay-openDay)
				for {
					beginDate = beginDate.AddDate(0, 1, 0)
					endDate = endDate.AddDate(0, 1, 0)
					if (date.Unix() <= endDate.Unix() && date.Unix() >= beginDate.Unix()) || (date.Unix() > endDate.Unix()) {
						//! 新服期与公服期冲突, 公服期活动顺延
						beginDate = beginDate.AddDate(0, 1, 0)
						endDate = endDate.AddDate(0, 1, 0)
					} else {
						break
					}
				}
			}

		}

		beginTime = beginDate.Unix()
		endTime = endDate.Unix()
	} else if activityInfo.TimeType == 2 { //! 按照日计算
		weekDay := int(now.Weekday())
		if weekDay == 0 { //! 特殊处理周末
			weekDay = 7
		}

		if activityInfo.EndTime == 7 && activityInfo.BeginTime == 1 {
			//! 永久活动
			beginTime = 0
			endTime = 0
		} else {
			beginDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			beginDate = beginDate.AddDate(0, 0, activityInfo.BeginTime-weekDay)

			endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())
			endDate = endDate.AddDate(0, 0, activityInfo.EndTime-weekDay)

			if activityInfo.ServerType == 2 { //! 公服期活动
				endDay := GetNewServerTypeActivityEndTime(activityInfo)
				if endDay > openDay {
					date := time.Now()
					date = date.AddDate(0, 0, endDay-openDay)
					for {
						beginDate = beginDate.AddDate(0, 0, 7)
						endDate = endDate.AddDate(0, 0, 7)
						if (date.Unix() <= endDate.Unix() && date.Unix() >= beginDate.Unix()) || (date.Unix() > endDate.Unix()) {
							//! 新服期与公服期冲突, 公服期活动顺延
							beginDate = beginDate.AddDate(0, 0, 7)
							endDate = endDate.AddDate(0, 0, 7)
						} else {
							break
						}
					}
				}
			}

			endTime = endDate.Unix()
			beginTime = beginDate.Unix()

		}
	} else if activityInfo.TimeType == 3 { //! 按照开服时间计算
		beginDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		beginDate = beginDate.AddDate(0, 0, -1*openDay)
		beginDate = beginDate.AddDate(0, 0, activityInfo.BeginTime)
		beginTime = beginDate.Unix()

		endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 59, now.Location())
		endDate = endDate.AddDate(0, 0, activityInfo.EndTime-openDay)
		endTime = endDate.Unix()
	} else if activityInfo.TimeType == 4 {
		day := activityInfo.BeginTime % 100
		month := (activityInfo.BeginTime - day) / 100
		if day < 1 || day > 31 || month < 1 || month > 12 {
			gamelog.Error("Invalid Activity BeginTime: %d", activityInfo.BeginTime)
			return -1, -1
		}

		beginData := time.Date(now.Year(), time.Month(month), day, 0, 0, 0, 0, now.Location())
		beginTime = beginData.Unix()

		day = activityInfo.EndTime % 100
		month = (activityInfo.EndTime - day) / 100
		if day < 1 || day > 31 || month < 1 || month > 12 {
			gamelog.Error("Invalid Activity BeginTime: %d", activityInfo.BeginTime)
			return -1, -1
		}

		endData := time.Date(now.Year(), time.Month(month), day, 23, 59, 59, 59, now.Location())
		endTime = endData.Unix()

		//! 若今年活动时间已过,则时间变更为明年
		if endTime < now.Unix() {
			beginData = beginData.AddDate(1, 0, 0)
			beginTime = beginData.Unix()

			endData = endData.AddDate(1, 0, 0)
			endTime = endData.Unix()
		}
	}

	return beginTime, endTime
}

func GetTodayOpenActivity(openDay int) []int {
	activityIDLst := []int{}
	now := time.Now()

	for _, v := range GT_ActivityLst {
		if v.TimeType == 1 { //! 按照月计算
			if now.Day() >= v.BeginTime && now.Day() <= v.EndTime {
				activityIDLst = append(activityIDLst, v.ID)
			}
		} else if v.TimeType == 2 { //! 按照周计算
			weekDay := int(now.Weekday())
			if weekDay == 0 { //! 特殊处理周末
				weekDay = 7
			}

			if weekDay >= v.BeginTime && weekDay <= v.EndTime {
				activityIDLst = append(activityIDLst, v.ID)
			}
		} else if v.TimeType == 3 { //! 新服活动
			gamelog.Error("OpenDay: %d  BeginTime: %d", openDay, v.BeginTime)
			if openDay >= v.BeginTime && openDay <= v.EndTime {
				activityIDLst = append(activityIDLst, v.ID)
			}
		} else {
			continue
		}
	}

	return activityIDLst
}

//! 开服竞赛配置表
type ST_CompetitionData struct {
	AwardType int
	Rank_min  int
	Rank_max  int
	Award     int
}

var GT_CompetitionData []ST_CompetitionData

func InitCompetitionParser(total int) bool {
	GT_CompetitionData = make([]ST_CompetitionData, total+1)
	return true
}

func ParseCompetitionRecord(rs *RecordSet) {
	awardType := CheckAtoi(rs.Values[0], 0)
	GT_CompetitionData[awardType].AwardType = awardType
	GT_CompetitionData[awardType].Rank_min = rs.GetFieldInt("rank_min")
	GT_CompetitionData[awardType].Rank_max = rs.GetFieldInt("rank_max")
	GT_CompetitionData[awardType].Award = rs.GetFieldInt("award")
}

func GetCompetitionAward(awardType int, rank int) int {
	for _, v := range GT_CompetitionData {
		if v.AwardType != awardType {
			continue
		}
		if rank >= v.Rank_min && rank <= v.Rank_max {
			return v.Award
		}
	}

	return 0
}

//! 领取体力
type ST_Activity_RecvAction struct {
	Award_Type int //! 奖励类型
	Time_Begin int //! 活动开始时间
	Time_End   int //! 活动结束时间
	AwardPro   int //! 额外奖励概率
	ActionID   int //! 奖励行动力ID
	ActionNum  int //! 奖励行动力数量
	MoneyID    int //! 额外奖励货币ID
	MoneyNum   int //! 额外奖励货币数量
}

var GT_RecvActionLst []ST_Activity_RecvAction

func InitRecvActionParser(total int) bool {
	//GT_RecvActionLst = make([]ST_Activity_RecvAction, total+1)
	return true
}

func ParseRecvActionRecord(rs *RecordSet) {
	var recv ST_Activity_RecvAction
	recv.Award_Type = rs.GetFieldInt("award_type")
	recv.Time_Begin = rs.GetFieldInt("time_begin")
	recv.Time_End = rs.GetFieldInt("time_end")
	recv.AwardPro = rs.GetFieldInt("award_pro")
	recv.ActionID = rs.GetFieldInt("action_id")
	recv.ActionNum = rs.GetFieldInt("action_num")
	recv.MoneyID = rs.GetFieldInt("money_id")
	recv.MoneyNum = rs.GetFieldInt("money_num")
	GT_RecvActionLst = append(GT_RecvActionLst, recv)
}

func GetRecvAction(awardType int, index int) *ST_Activity_RecvAction {
	if index > len(GT_RecvActionLst) {
		gamelog.Error("GetRecvAction Error: Invalid index %d", index)
		return nil
	}

	var awardLst []ST_Activity_RecvAction
	for i, v := range GT_RecvActionLst {
		if v.Award_Type == awardType {
			awardLst = append(awardLst, GT_RecvActionLst[i])
		}
	}

	return &awardLst[index]
}

func IsRecvActionTime(awardType int) int {
	now := time.Now()
	sec := now.Hour()*3600 + now.Minute()*60 + now.Second()

	var awardLst []ST_Activity_RecvAction
	for i, v := range GT_RecvActionLst {
		if v.Award_Type == awardType {
			awardLst = append(awardLst, GT_RecvActionLst[i])
		}
	}

	for i, v := range awardLst {
		if sec >= v.Time_Begin && sec <= v.Time_End {
			return i + 1
		}
	}

	return -1
}

//! 折扣销售
type ST_Activity_DiscoutSale struct {
	AwardType int
	MoneyID   int
	MoneyNum  int
	Award     int
	IsSelect  int
	Times     int
}

var GT_DiscountSaleLst []ST_Activity_DiscoutSale

func InitDiscountSaleParser(total int) bool {
	return true
}

func ParseDiscountSaleRecord(rs *RecordSet) {
	var discount ST_Activity_DiscoutSale
	discount.AwardType = rs.GetFieldInt("award_type")
	discount.MoneyID = rs.GetFieldInt("moneyid")
	discount.MoneyNum = rs.GetFieldInt("moneynum")

	discount.Award = rs.GetFieldInt("award_id")

	discount.IsSelect = rs.GetFieldInt("is_select")
	discount.Times = rs.GetFieldInt("times")
	GT_DiscountSaleLst = append(GT_DiscountSaleLst, discount)
}

func GetDiscountSaleInfo(awardType int) (itemLst []ST_Activity_DiscoutSale) {
	for i := 0; i < len(GT_DiscountSaleLst); i++ {
		if GT_DiscountSaleLst[i].AwardType == awardType {
			itemLst = append(itemLst, GT_DiscountSaleLst[i])
		}
	}

	return itemLst
}

//! 累计登录
type ST_Activity_Login struct {
	AwardType int
	Award     int
	IsSelect  int
}

var GT_ActivityLoginMap map[int][]ST_Activity_Login

func InitActivityLoginParser(total int) bool {
	GT_ActivityLoginMap = make(map[int][]ST_Activity_Login)
	return true
}

func ParseActivityLoginRecord(rs *RecordSet) {
	activity_award := CheckAtoi(rs.Values[0], 0)

	var login ST_Activity_Login
	login.AwardType = activity_award
	login.Award = rs.GetFieldInt("award")
	login.IsSelect = rs.GetFieldInt("is_select")

	GT_ActivityLoginMap[activity_award] = append(GT_ActivityLoginMap[activity_award], login)
}

func GetActivityLoginInfo(awardType int) []ST_Activity_Login {
	array, ok := GT_ActivityLoginMap[awardType]
	if ok == false {
		gamelog.Error("GetActivityLoginInfo Error: invalid awardType %d", awardType)
		return []ST_Activity_Login{}
	}
	return array
}

//! 迎财神
type ST_Activity_Money struct {
	AwardType  int //! 奖励类型
	CDTime     int //! 每次领取间隔时间
	MoneyID    int //! 领取货币ID
	MoneyNum   int //! 领取货币数量
	AwardTimes int //! 奖励次数
	ItemID     int //! 额外奖励物品ID
	ItemNum    int //! 额外奖励物品数量
	LuckPro    int //! 额外奖励概率
}

var GT_MoneyGoldConfig []ST_Activity_Money

func InitActivityMoneyParser(total int) bool {
	return true
}

func ParseActivityMoneyRecord(rs *RecordSet) {
	var config ST_Activity_Money
	config.AwardType = rs.GetFieldInt("award_type")
	config.CDTime = rs.GetFieldInt("cd_time")
	config.MoneyID = rs.GetFieldInt("moneyid")
	config.MoneyNum = rs.GetFieldInt("moneynum")
	config.AwardTimes = rs.GetFieldInt("awardtimes")
	config.ItemID = rs.GetFieldInt("itemid")
	config.ItemNum = rs.GetFieldInt("itemnum")
	config.LuckPro = rs.GetFieldInt("luckpro")
	GT_MoneyGoldConfig = append(GT_MoneyGoldConfig, config)
}

func GetMoneyGoldInfo(awardType int) *ST_Activity_Money {

	length := len(GT_MoneyGoldConfig)
	for i := 0; i < length; i++ {
		if GT_MoneyGoldConfig[i].AwardType == awardType {
			return &GT_MoneyGoldConfig[i]
		}
	}

	return nil
}

//! 充值返利
type ST_Activity_Recharge struct {
	AwardType int
	Recharge  int //! 充值额度
	Award     int //! 奖励
	Times     int
}

var GT_ActivityRechargeMap map[int][]ST_Activity_Recharge

func InitActivityRechargeParser(total int) bool {
	GT_ActivityRechargeMap = make(map[int][]ST_Activity_Recharge)
	return true
}

func ParseActivityRechargeRecord(rs *RecordSet) {
	var recharge ST_Activity_Recharge
	recharge.AwardType = rs.GetFieldInt("award_type")
	recharge.Award = rs.GetFieldInt("award")
	recharge.Recharge = rs.GetFieldInt("recharge")
	recharge.Times = rs.GetFieldInt("times")
	GT_ActivityRechargeMap[recharge.AwardType] = append(GT_ActivityRechargeMap[recharge.AwardType], recharge)
}

func GetRechargeInfo(awardType int) []ST_Activity_Recharge {
	array, ok := GT_ActivityRechargeMap[awardType]
	if ok == false {
		gamelog.Error("GetRechargeInfo Error: invalid awardType:%d", awardType)
		return []ST_Activity_Recharge{}
	}

	var rechargeLst []ST_Activity_Recharge
	for _, v := range array {
		rechargeLst = append(rechargeLst, v)
	}

	return rechargeLst
}

//! 限时日常
type ST_Activity_LimitDaily struct {
	AwardType int
	TaskType  int //! 取值taskType表 任务类型
	Count     int //! 达标数额
	Award     int //! 奖励
	IsSelect  int //! 是否为多选一
}

var GT_ActivityLimitDailyMap map[int][]ST_Activity_LimitDaily

func InitActivityLimitDailyParser(total int) bool {
	GT_ActivityLimitDailyMap = make(map[int][]ST_Activity_LimitDaily)
	return true
}

func ParseActivityLimitDailyRecord(rs *RecordSet) {
	var activity ST_Activity_LimitDaily
	activity.AwardType = rs.GetFieldInt("award_type")
	activity.TaskType = rs.GetFieldInt("task_type")
	activity.Count = rs.GetFieldInt("count")
	activity.Award = rs.GetFieldInt("award")
	activity.IsSelect = rs.GetFieldInt("is_select")
	GT_ActivityLimitDailyMap[activity.AwardType] = append(GT_ActivityLimitDailyMap[activity.AwardType], activity)
}

func GetActivityLimitDaily(awardType int) []ST_Activity_LimitDaily {
	array, ok := GT_ActivityLimitDailyMap[awardType]
	if ok == false {
		gamelog.Error("GetActivityLimitDaily Error: invalid awardType:%d", awardType)
		return []ST_Activity_LimitDaily{}
	}

	return array
}

//! 周周盈
type ST_Activity_WeekAward struct {
	ID          int //! 唯一标识
	AwardType   int //! 奖励模板
	LoginDay    int //! 登录天数
	RechargeNum int //! 累充数额
	AwardID     int //! 奖励
	IsSelect    int //! 是否为多选一
}

var GT_ActivityWeekAwardMap map[int][]ST_Activity_WeekAward

func InitActivityWeekAwardParser(total int) bool {
	GT_ActivityWeekAwardMap = make(map[int][]ST_Activity_WeekAward)
	return true
}

func ParseActivityWeekAwardRecord(rs *RecordSet) {
	var activity ST_Activity_WeekAward
	activity.AwardType = rs.GetFieldInt("award_type")
	activity.LoginDay = rs.GetFieldInt("login_day")
	activity.RechargeNum = rs.GetFieldInt("recharge")
	activity.AwardID = rs.GetFieldInt("award")
	activity.IsSelect = rs.GetFieldInt("is_select")
	GT_ActivityWeekAwardMap[activity.AwardType] = append(GT_ActivityWeekAwardMap[activity.AwardType], activity)
}

func GetWeekAwardInfoLst(awardType int) []ST_Activity_WeekAward {
	data, ok := GT_ActivityWeekAwardMap[awardType]
	if ok == false {
		gamelog.Error("GetWeekAwardInfo Error: Invalid type : %d", awardType)
		return nil
	}

	return data
}

func GetWeekAwardInfo(awardType int, index int) *ST_Activity_WeekAward {
	data, ok := GT_ActivityWeekAwardMap[awardType]
	if ok == false {
		gamelog.Error("GetWeekAwardInfo Error: Invalid type : %d  index: %d", awardType, index)
		return nil
	}

	if index-1 < 0 || index > len(data) {
		gamelog.Error("GetWeekAwardInfo Error: Invalid type : %d  index: %d", awardType, index)
		return nil
	}

	return &data[index-1]
}

//! 等级礼包
type ST_Activity_LevelGift struct {
	AwardType int    //! 奖励模板
	ID        int    //! ID
	Level     string //! 需求等级
	Award     int    //! 奖励
	MoneyID   int    //! 货币ID
	MoneyNum  int    //! 货币数量
	BuyTimes  int    //! 可购买次数
	DeadLine  int    //! 过期时间
}

var GT_ActivityLevelGiftMap map[int][]ST_Activity_LevelGift

func InitActivityLevelGiftParser(total int) bool {
	GT_ActivityLevelGiftMap = make(map[int][]ST_Activity_LevelGift)
	return true
}

func ParseActivityLevelGiftRecord(rs *RecordSet) {
	var activity ST_Activity_LevelGift
	activity.AwardType = rs.GetFieldInt("award_type")
	activity.ID = rs.GetFieldInt("id")
	activity.Level = rs.GetFieldString("level")
	activity.Award = rs.GetFieldInt("award_id")
	activity.MoneyID = rs.GetFieldInt("money_id")
	activity.MoneyNum = rs.GetFieldInt("money_num")
	activity.BuyTimes = rs.GetFieldInt("buy_times")
	activity.DeadLine = rs.GetFieldInt("dead_line")
	GT_ActivityLevelGiftMap[activity.AwardType] = append(GT_ActivityLevelGiftMap[activity.AwardType], activity)
}

func GetLevelGiftInfo(awardType int, id int) *ST_Activity_LevelGift {
	data, ok := GT_ActivityLevelGiftMap[awardType]
	if ok == false {
		gamelog.Error("GetLevelGiftInfo Error: Invalid type : %d  id: %d", awardType, id)
		return nil
	}

	if id > len(data) || id <= 0 {
		gamelog.Error("GetWeekAwardInfo Error: Invalid type : %d  id: %d", awardType, id)
		return nil
	}

	return &data[id-1]
}

func GetLevelGiftLst(awardType int) []ST_Activity_LevelGift {
	data, ok := GT_ActivityLevelGiftMap[awardType]
	if ok == false {
		gamelog.Error("GetLevelGiftLst Error: Invalid type : %d", awardType)
		return nil
	}

	return data
}

//! 月基金
type ST_MonthFund struct {
	AwardType int //! 奖励模板
	Day       int //! 天数
	ItemID    int //! 道具ID
	ItemNum   int //! 道具数量
}

var GT_ActivityMonthFund map[int][]ST_MonthFund

func InitActivityMonthFundParser(total int) bool {
	GT_ActivityMonthFund = make(map[int][]ST_MonthFund)
	return true
}

func ParseActivityMonthFundRecord(rs *RecordSet) {
	var fund ST_MonthFund
	fund.AwardType = rs.GetFieldInt("award_type")
	fund.Day = rs.GetFieldInt("day")
	fund.ItemID = rs.GetFieldInt("item_id")
	fund.ItemNum = rs.GetFieldInt("item_num")
	GT_ActivityMonthFund[fund.AwardType] = append(GT_ActivityMonthFund[fund.AwardType], fund)
}

func GetMonthFundAward(awardType int, day int) *ST_MonthFund {
	monthFundLst, ok := GT_ActivityMonthFund[awardType]
	if ok == false {
		gamelog.Error("GetMonthFundAward Error: Invalid awardType %d", awardType)
		return nil
	}

	if day > len(monthFundLst) || day <= 0 {
		gamelog.Error("GetMonthFundAward Error: Invalid day %d", day)
		return nil
	}

	return &monthFundLst[day-1]
}

func GetMonthFundAwardCount(awardType int) int {
	monthFundLst, ok := GT_ActivityMonthFund[awardType]
	if ok == false {
		gamelog.Error("GetMonthFundAwardCount Error: Invalid awardType %d", awardType)
		return 0
	}

	return len(monthFundLst)
}

type ST_LimitSaleItem struct {
	ID       int
	ItemType int //! 1->普通货币 2->元宝
	ItemID   int
	ItemNum  int
	MoneyID  int
	MoneyNum int
	Discount int //! 折扣
	Score    int //! 购买获得积分
	Original int //! 原价
}

var GT_LimitSaleItemLst [][]ST_LimitSaleItem

func InitLimitSaleItemParser(total int) bool {
	GT_LimitSaleItemLst = make([][]ST_LimitSaleItem, 2)
	return true
}

func ParseLimitSaleItemRecord(rs *RecordSet) {
	var item ST_LimitSaleItem

	item.ID = CheckAtoi(rs.Values[0], 0)
	item.ItemType = rs.GetFieldInt("item_type")
	item.ItemID = rs.GetFieldInt("item_id")
	item.ItemNum = rs.GetFieldInt("item_num")
	item.MoneyID = rs.GetFieldInt("money_id")
	item.MoneyNum = rs.GetFieldInt("money_num")
	item.Discount = rs.GetFieldInt("discount")
	item.Score = rs.GetFieldInt("score")
	item.Original = rs.GetFieldInt("original")

	GT_LimitSaleItemLst[item.ItemType-1] = append(GT_LimitSaleItemLst[item.ItemType-1], item)
}

func GetLimitSaleItemInfo(id int) *ST_LimitSaleItem {
	for i, v := range GT_LimitSaleItemLst {
		for j, n := range v {
			if n.ID == id {
				return &GT_LimitSaleItemLst[i][j]
			}
		}
	}

	gamelog.Error("GetLimitSaleIteminfo nil, Invalid ID: %d", id)
	return nil
}

func RandLimitSaleItem() []int {
	randIDLst := []int{}

	//! 2个元宝商品
	for i := 0; i < 2; i++ {
		randIndex := r.Intn(len(GT_LimitSaleItemLst[1]))
		item := GT_LimitSaleItemLst[1][randIndex]
		randIDLst = append(randIDLst, item.ID)
	}

	//! 1个普通货币
	randIndex := r.Intn(len(GT_LimitSaleItemLst[0]))
	item := GT_LimitSaleItemLst[0][randIndex]
	randIDLst = append(randIDLst, item.ID)

	//! 3个元宝商品
	for i := 0; i < 3; i++ {
		randIndex = r.Intn(len(GT_LimitSaleItemLst[1]))
		item = GT_LimitSaleItemLst[1][randIndex]
		randIDLst = append(randIDLst, item.ID)
	}

	return randIDLst
}

type ST_LimitSaleAllAward struct {
	ID      int
	Award   int
	NeedNum int //!  需求购买人数
}

var GT_LimitSaleAllAwardLst []ST_LimitSaleAllAward

func InitLimitSaleAllAwardParser(total int) bool {
	GT_LimitSaleAllAwardLst = make([]ST_LimitSaleAllAward, total+1)
	return true
}

func ParseLimitSaleAllAwardRecord(rs *RecordSet) {

	id := CheckAtoi(rs.Values[0], 0)

	GT_LimitSaleAllAwardLst[id].ID = id
	GT_LimitSaleAllAwardLst[id].Award = rs.GetFieldInt("award")
	GT_LimitSaleAllAwardLst[id].NeedNum = rs.GetFieldInt("need_num")
}

func GetLimitSaleAllAwardInfo(id int) *ST_LimitSaleAllAward {
	if id > len(GT_LimitSaleAllAwardLst)-1 {
		gamelog.Error("GetLimitSaleAllAwardinfo Error: Invalid id %d", id)
		return nil
	}

	return &GT_LimitSaleAllAwardLst[id]
}

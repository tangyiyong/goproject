package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

type TActivity interface {
	//! 创建初始化
	Init(activityID int32, mPtr *TActivityModule, vercode int32, resetcode int32)

	//! 设置父模块指针
	SetModulePtr(mPtr *TActivityModule)

	//! 刷新数据
	Refresh(versionCode int32)

	//! 活动结束(活动重置)
	End(versionCode int32, resetCode int32)

	//取活动的刷新版本
	GetRefreshV() int32

	//取重置版本
	GetResetV() int32

	//活动是否存在玩家操作机会， 用于客户端的红点提示。
	RedTip() bool

	//重置的数据存储(在End接口内部调用)
	DB_Reset()

	//更新的数据存储(在Refresh接口内部调用)
	DB_Refresh()
}

//! 活动模块
type TActivityModule struct {
	PlayerID int32 `bson:"_id"`

	//! 首充/次充
	FirstCharge TActivityFirstCharge

	//! 月卡
	MonthCard TActivityMonthCard

	//! 迎财神
	MoneyGod TActivityMoneyGod

	//! 折扣贩售
	DiscountSale []TActivityDiscount

	//! 充值反馈
	Recharge []TActivityRecharge

	//! 单冲返利
	SingleRecharge []TActivitySingleRecharge

	//! 领体力
	ReceiveAction TActivityAction

	//! 登录送礼
	Login []TActivityLogin

	//! 签到相关
	Sign TActivitySign

	//! VIP礼包
	VipGift TActivityVipGift

	//! 开服基金
	OpenFund TActivityOpenFund

	//! 限时日常任务
	LimitDaily []TActivityLimitDaily

	//! 巡回探宝
	HuntTreasure TActivityHunt

	//! 幸运轮盘
	LuckyWheel TActivityWheel

	//! 团购
	GroupPurchase TActivityGroupPurchase

	// 卡牌大师
	CardMaster TCardMasterInfo

	// 月光集市
	MoonShop TMoonShop

	// 沙滩宝贝
	BeachBaby TBeachBabyInfo

	//! 欢庆佳节
	Festival TActivityFestival

	//! 七日活动列表
	SevenDay []TActivitySevenDay

	//! 周周盈
	WeekAward TActivityWeekAward

	//! 等级礼包
	LevelGift TActivityLevelGift

	//! 月基金
	MonthFund TActivityMonthFund

	//! 巅峰特惠
	RankGift TActivityRankGift

	//! 限时特惠
	LimitSale TActivityLimitSale

	activityPtrs map[int32]TActivity

	ownplayer *TPlayer
}

func (self *TActivityModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TActivityModule) OnCreate(playerid int32) {
	self.activityPtrs = make(map[int32]TActivity)
	for i := 0; i < len(G_GlobalVariables.ActivityLst); i++ {
		activityID := G_GlobalVariables.ActivityLst[i].ActivityID
		activityType := G_GlobalVariables.ActivityLst[i].activityType
		verionCode := G_GlobalVariables.ActivityLst[i].VersionCode
		resetCode := G_GlobalVariables.ActivityLst[i].ResetCode

		if G_GlobalVariables.IsActivityOpen(activityID) == false {
			continue
		}

		if activityType == gamedata.Activity_Sign {
			//! 签到
			self.Sign.Init(activityID, self, verionCode, resetCode)

		} else if activityType == gamedata.Activity_Login {
			//! 登录奖励
			loginActivity := TActivityLogin{}
			loginActivity.Init(activityID, self, verionCode, resetCode)
			self.Login = append(self.Login, loginActivity)

		} else if activityType == gamedata.Activity_Recv_Action {
			//! 领取体力
			self.ReceiveAction.Init(activityID, self, verionCode, resetCode)

		} else if activityType == gamedata.Activity_Money_Gold {
			//! 迎财神
			self.MoneyGod.Init(activityID, self, verionCode, resetCode)

		} else if activityType == gamedata.Activity_Recharge_Gift {
			//! 累积充值
			rechargeActivity := TActivityRecharge{}
			rechargeActivity.Init(activityID, self, verionCode, resetCode)
			self.Recharge = append(self.Recharge, rechargeActivity)

		} else if activityType == gamedata.Activity_Open_Fund {
			//! 开服基金
			self.OpenFund.Init(activityID, self, verionCode, resetCode)

		} else if activityType == gamedata.Activity_Discount_Sale {
			//! 折扣贩售
			discountActivity := TActivityDiscount{}
			discountActivity.Init(activityID, self, verionCode, resetCode)
			self.DiscountSale = append(self.DiscountSale, discountActivity)

		} else if activityType == gamedata.Activity_First_Recharge {
			//! 首充
			self.FirstCharge.Init(activityID, self, verionCode, resetCode)

		} else if activityType == gamedata.Activity_Singel_Recharge {
			//! 单充返利
			singleRechargeActivity := TActivitySingleRecharge{}
			singleRechargeActivity.Init(activityID, self, verionCode, resetCode)
			self.SingleRecharge = append(self.SingleRecharge, singleRechargeActivity)

		} else if activityType == gamedata.Activity_Limit_Daily_Task {
			//! 限时日常
			limitDailyActivity := TActivityLimitDaily{}
			limitDailyActivity.Init(activityID, self, verionCode, resetCode)
			self.LimitDaily = append(self.LimitDaily, limitDailyActivity)
		} else if activityType == gamedata.Activity_Hunt_Treasure {
			//! 巡回探宝
			self.HuntTreasure.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Luckly_Wheel {
			//! 幸运轮盘
			self.LuckyWheel.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Group_Purchase {
			//! 团购
			self.GroupPurchase.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Moon_Card {
			//! 月卡
			self.MonthCard.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Vip_Gift {
			//! VIP礼包
			self.VipGift.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Card_Master {
			// 卡牌大师
			self.CardMaster.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_MoonlightShop {
			// 月光集市
			self.MoonShop.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Beach_Baby {
			// 沙滩宝贝
			self.BeachBaby.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Festival {
			//! 欢庆佳节
			self.Festival.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Seven {
			//! 七日活动
			SevenDay := TActivitySevenDay{}
			SevenDay.Init(activityID, self, verionCode, resetCode)
			self.SevenDay = append(self.SevenDay, SevenDay)
		} else if activityType == gamedata.Activity_Week_Award {
			//! 周周盈
			self.WeekAward.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Level_Gift {
			//! 等级礼包
			self.LevelGift.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Month_Fund {
			//! 月基金
			self.MonthFund.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_Rank_Sale {
			//! 巅峰特惠
			self.RankGift.Init(activityID, self, verionCode, resetCode)
		} else if activityType == gamedata.Activity_LimitSale {
			//! 限时特惠
			self.LimitSale.Init(activityID, self, verionCode, resetCode)
		}
	}
	mongodb.InsertToDB("PlayerActivity", self)
}

func (self *TActivityModule) OnDestroy(playerid int32) {

}

func (self *TActivityModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *TActivityModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *TActivityModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerActivity").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerActivity Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
	self.activityPtrs = make(map[int32]TActivity)

	//! 签到
	self.Sign.SetModulePtr(self)

	//! 登录奖励
	for i, _ := range self.Login {
		self.Login[i].SetModulePtr(self)
	}

	//! 领取体力
	self.ReceiveAction.SetModulePtr(self)

	//! 迎财神
	self.MoneyGod.SetModulePtr(self)

	//! 累积充值
	for i, _ := range self.Recharge {
		self.Recharge[i].SetModulePtr(self)
	}

	//! 开服基金
	self.OpenFund.SetModulePtr(self)

	//! 折扣贩售
	for i, _ := range self.DiscountSale {
		self.DiscountSale[i].SetModulePtr(self)
	}

	//! 首充
	self.FirstCharge.SetModulePtr(self)

	//! 单充返利
	for i, _ := range self.SingleRecharge {
		self.SingleRecharge[i].SetModulePtr(self)
	}

	//! 限时日常
	for i, _ := range self.LimitDaily {
		self.LimitDaily[i].SetModulePtr(self)
	}

	//! 巡回探宝
	self.HuntTreasure.SetModulePtr(self)

	//! 幸运轮盘
	self.LuckyWheel.SetModulePtr(self)

	//! 团购
	self.GroupPurchase.SetModulePtr(self)

	//! 月卡
	self.MonthCard.SetModulePtr(self)

	//! VIP礼包
	self.VipGift.SetModulePtr(self)

	// 卡牌大师
	self.CardMaster.SetModulePtr(self)

	// 月光集市
	self.MoonShop.SetModulePtr(self)

	// 沙滩宝贝
	self.BeachBaby.SetModulePtr(self)

	//! 欢庆佳节
	self.Festival.SetModulePtr(self)

	//! 七日活动
	for i, _ := range self.SevenDay {
		self.SevenDay[i].SetModulePtr(self)
	}

	//! 周周盈
	self.WeekAward.SetModulePtr(self)

	//! 等级礼包
	self.LevelGift.SetModulePtr(self)

	//! 月基金
	self.MonthFund.SetModulePtr(self)

	//! 巅峰特惠
	self.RankGift.SetModulePtr(self)

	//! 限时特惠
	self.LimitSale.SetModulePtr(self)

}

func (self *TActivityModule) AddActivityPtr(activityID int, activityPtr TActivity) {

}

//! 新增活动初始化
func (self *TActivityModule) CheckNewActivity() {
	//! 检测新增活动初始化
	for i, v := range G_GlobalVariables.ActivityLst {
		activityID := G_GlobalVariables.ActivityLst[i].ActivityID
		versionCode := G_GlobalVariables.ActivityLst[i].VersionCode
		resetCode := G_GlobalVariables.ActivityLst[i].ResetCode

		if G_GlobalVariables.IsActivityOpen(activityID) == false {
			continue
		}

		activity, ok := self.activityPtrs[v.ActivityID]
		if activity != nil && ok {
			continue
		}

		if v.activityType == gamedata.Activity_Login {
			//! 登录奖励
			loginActivity := TActivityLogin{}
			loginActivity.Init(v.ActivityID, self, versionCode, resetCode)
			self.Login = append(self.Login, loginActivity)
			self.DB_AddNewLoginActivity(loginActivity)

		} else if v.activityType == gamedata.Activity_Recharge_Gift {
			//! 累积充值
			rechargeActivity := TActivityRecharge{}
			rechargeActivity.Init(v.ActivityID, self, versionCode, resetCode)
			self.Recharge = append(self.Recharge, rechargeActivity)
			self.DB_AddNewRechargeActivity(rechargeActivity)

		} else if v.activityType == gamedata.Activity_Discount_Sale {
			//! 折扣贩售
			discountActivity := TActivityDiscount{}
			discountActivity.Init(v.ActivityID, self, versionCode, resetCode)
			self.DiscountSale = append(self.DiscountSale, discountActivity)
			self.DB_AddNewDiscountActivity(discountActivity)

		} else if v.activityType == gamedata.Activity_Singel_Recharge {
			//! 单充返利
			singleRechargeActivity := TActivitySingleRecharge{}
			singleRechargeActivity.Init(v.ActivityID, self, versionCode, resetCode)
			self.SingleRecharge = append(self.SingleRecharge, singleRechargeActivity)
			self.DB_AddNewSingleRechargeActivity(singleRechargeActivity)

		} else if v.activityType == gamedata.Activity_Limit_Daily_Task {
			//! 限时日常
			limitDailyActivity := TActivityLimitDaily{}
			limitDailyActivity.Init(v.ActivityID, self, versionCode, resetCode)
			self.LimitDaily = append(self.LimitDaily, limitDailyActivity)
			self.DB_AddNewLimitDailyActivity(limitDailyActivity)
		} else if v.activityType == gamedata.Activity_Hunt_Treasure && v.ActivityID != self.HuntTreasure.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.HuntTreasure.Init(v.ActivityID, self, versionCode, resetCode)
			self.HuntTreasure.DB_Reset()
		} else if v.activityType == gamedata.Activity_Luckly_Wheel && v.ActivityID != self.LuckyWheel.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.LuckyWheel.Init(v.ActivityID, self, versionCode, resetCode)
			self.LuckyWheel.DB_Reset()
		} else if v.activityType == gamedata.Activity_Group_Purchase && v.ActivityID != self.GroupPurchase.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.GroupPurchase.Init(v.ActivityID, self, versionCode, resetCode)
			self.GroupPurchase.DB_Reset()
		} else if v.activityType == gamedata.Activity_Card_Master && v.ActivityID != self.CardMaster.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.CardMaster.Init(v.ActivityID, self, versionCode, resetCode)
			self.CardMaster.DB_Reset()
		} else if v.activityType == gamedata.Activity_Festival && v.ActivityID != self.Festival.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) {
			self.Festival.Init(v.ActivityID, self, versionCode, resetCode)
			self.Festival.DB_Reset()
		} else if v.activityType == gamedata.Activity_MoonlightShop && v.ActivityID != self.MoonShop.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.MoonShop.Init(v.ActivityID, self, versionCode, resetCode)
			self.MoonShop.DB_Reset()
		} else if v.activityType == gamedata.Activity_Beach_Baby && v.ActivityID != self.BeachBaby.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.BeachBaby.Init(v.ActivityID, self, versionCode, resetCode)
			self.BeachBaby.DB_Reset()
		} else if v.activityType == gamedata.Activity_Seven {
			sevenDay := TActivitySevenDay{}
			sevenDay.Init(v.ActivityID, self, versionCode, resetCode)
			self.SevenDay = append(self.SevenDay, sevenDay)
			self.DB_AddNewSevenDay(sevenDay)
		} else if v.activityType == gamedata.Activity_Week_Award && v.ActivityID != self.WeekAward.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.WeekAward.Init(v.ActivityID, self, versionCode, resetCode)
			self.WeekAward.DB_Reset()
		} else if v.activityType == gamedata.Activity_Level_Gift && v.ActivityID != self.LevelGift.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.LevelGift.Init(v.ActivityID, self, versionCode, resetCode)
			self.LevelGift.DB_Reset()
		} else if v.activityType == gamedata.Activity_Month_Fund && v.ActivityID != self.MonthFund.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.MonthFund.Init(v.ActivityID, self, versionCode, resetCode)
			self.MonthFund.DB_Reset()
		} else if v.activityType == gamedata.Activity_Rank_Sale && v.ActivityID != self.RankGift.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.RankGift.Init(v.ActivityID, self, versionCode, resetCode)
			self.RankGift.DB_Reset()
		} else if v.activityType == gamedata.Activity_Sign && v.ActivityID != self.Sign.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.Sign.Init(v.ActivityID, self, versionCode, resetCode)
			self.Sign.DB_Reset()
		}

	}
}

//! 检测重置
func (self *TActivityModule) CheckReset() {
	//! 检测迎财神时间
	self.MoneyGod.CheckMoneyGod()
	for i := 0; i < len(G_GlobalVariables.ActivityLst); i++ {
		pActivity, ok := self.activityPtrs[G_GlobalVariables.ActivityLst[i].ActivityID]
		if !ok || pActivity == nil {
			continue //! 开服竞赛无需存储数据,且玩家创建账号时,若活动已经结束,未初始化活动ID
		}

		//! 如果活动己任经关闭则不进行处理
		if G_GlobalVariables.ActivityLst[i].Status == 0 {
			continue
		}

		if G_GlobalVariables.ActivityLst[i].ResetCode == pActivity.GetResetV() {
			if G_GlobalVariables.ActivityLst[i].VersionCode > pActivity.GetRefreshV() {
				pActivity.Refresh(G_GlobalVariables.ActivityLst[i].VersionCode)
			}
		} else {
			pActivity.End(G_GlobalVariables.ActivityLst[i].VersionCode, G_GlobalVariables.ActivityLst[i].ResetCode)
			pActivity.Refresh(G_GlobalVariables.ActivityLst[i].VersionCode)
		}

	}

	//! 检测新增活动
	self.CheckNewActivity()

}

func (self *TActivityModule) AddLoginDay() {
	self.CheckReset()
	for i, v := range self.Login {
		if G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
			self.Login[i].AddLoginDay(i)
		}
	}
}

func (self *TActivityModule) AddRechargeValue(value int) {
	for _, v := range G_GlobalVariables.ActivityLst {
		if v.activityType == gamedata.Activity_Recharge_Gift {
			for i, n := range self.Recharge {
				if n.ActivityID == v.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
					self.Recharge[i].RechargeValue += value
					self.DB_UpdateTotalRecharge(i, self.Recharge[i].RechargeValue)
				}
			}
		} else if v.activityType == gamedata.Activity_Singel_Recharge {
			for i, n := range self.SingleRecharge {
				if n.ActivityID == v.ActivityID && G_GlobalVariables.IsActivityOpen(v.ActivityID) == true {
					self.SingleRecharge[i].RechargeRecord = append(self.SingleRecharge[i].RechargeRecord, TSingleRechargeRecord{value, 0})
					self.DB_UpdateSingleRecharge(i, TSingleRechargeRecord{value, 0})
				}
			}
		}
	}
}

func (self *TActivityModule) CheckSingleRecharge(activityID int32, money int) (bool, int) {
	for _, v := range self.SingleRecharge {
		if v.ActivityID == activityID {
			for i, n := range v.RechargeRecord {
				if n.Money == money && n.Status == 0 {
					return true, i
				}
			}
		}

	}

	return false, -1
}

//! 获取商品信息
func (self *TActivityModule) GetItemShoppingInfo(activityID int32, index int) *TDiscountSaleGoodsInfo {
	for i, v := range self.DiscountSale {
		if v.ActivityID == activityID {
			for j, n := range v.ShopLst {
				if n.Index == index {
					return &self.DiscountSale[i].ShopLst[j]
				}
			}
		}
	}

	return nil
}

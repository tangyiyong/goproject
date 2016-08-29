package mainlogic

import (
	"gamelog"
	"gamesvr/gamedata"
	"sync"
	"time"
	"utility"
)

//玩家子模块的基类
type TModule interface {
	//创建player
	OnCreate(playerid int32)

	//销毁player
	OnDestroy(playerid int32)

	//player 进入游戏
	OnPlayerOnline(playerid int32)

	//player 离开游戏
	OnPlayerOffline(playerid int32)

	//从数据库加载玩家
	OnPlayerLoad(playerid int32)

	//玩家响应新一天
	OnNewDay(newday int)
}

type TPlayer struct {
	RoleMoudle  TRoleMoudle   //基本角色模块
	HeroMoudle  THeroMoudle   //战斗英雄模块
	TaskMoudle  TTaskMoudle   //任务模块
	VipMoudle   TVipMoudle    //VIP数据模块
	MailMoudle  TMailMoudle   //邮件模块
	CopyMoudle  TCopyMoudle   //副本模块
	BagMoudle   TBagMoudle    //背包模块
	HangMoudle  THangUpMoudle //挂机模块
	ScoreMoudle TScoreMoudle  //积分赛模块

	StoreModule        TStoreModule        //! 商店模块
	SanGuoZhiModule    TSanGuoZhiModule    //! 三国志模块
	MallModule         TMallModule         //! 商城模块
	SummonModule       TSummonModule       //! 召唤模块
	ArenaModule        TArenaModule        //! 竞技场模块
	RobModule          TRobModule          //! 夺宝模块
	SangokuMusouModule TSangokuMusouModule //! 三国无双模块
	AwardCenterModule  TAwardCenterModule  //! 领奖中心模块
	TerritoryModule    TTerritoryModule    //! 领地攻伐模块
	RebelModule        TRebelModule        //! 围剿叛军模块
	FriendMoudle       TFriendMoudle       //! 好友模块
	MiningModule       TMiningModule       //! 挖矿模块
	GuildModule        TGuildModule        //! 公会模块
	BlackMarketModule  TBlackMarketModule  //! 黑市模块
	FameHallModule     TFameHallModule     //! 名人堂
	TitleModule        TTitleModule        //! 称号
	FoodWarModule      TFoodWarModule      //! 夺粮战
	ActivityModule     TActivityModule     //! 活动模块
	WanderMoudle       TWanderModule       //! 云游模块
	HeroSoulsModule    THeroSoulsModule    //! 将灵模块
	CamBattleModule    TCampBattleModule   //! 阵营战模块

	ChargeModule TChargeMoudle //! 充值

	//非存数据库的临时状态数据
	playerid    int32        //角色ID
	pSimpleInfo *TSimpleInfo //角色简信息
	mutex       sync.Mutex   //玩家的一些操作的锁
	isLock      bool
}

//玩家初始化
func (player *TPlayer) InitModules(playerid int32) {
	if playerid <= 0 {
		gamelog.Error("InitModules Error : Invalid PlayerID:%d", playerid)
		return
	}
	player.RoleMoudle.SetPlayerPtr(playerid, player)
	player.HeroMoudle.SetPlayerPtr(playerid, player)
	player.TaskMoudle.SetPlayerPtr(playerid, player)
	player.VipMoudle.SetPlayerPtr(playerid, player)
	player.MailMoudle.SetPlayerPtr(playerid, player)
	player.CopyMoudle.SetPlayerPtr(playerid, player)
	player.BagMoudle.SetPlayerPtr(playerid, player)
	player.StoreModule.SetPlayerPtr(playerid, player)
	player.SanGuoZhiModule.SetPlayerPtr(playerid, player)
	player.MallModule.SetPlayerPtr(playerid, player)
	player.SummonModule.SetPlayerPtr(playerid, player)
	player.ArenaModule.SetPlayerPtr(playerid, player)
	player.RobModule.SetPlayerPtr(playerid, player)
	player.SangokuMusouModule.SetPlayerPtr(playerid, player)
	player.AwardCenterModule.SetPlayerPtr(playerid, player)
	player.TerritoryModule.SetPlayerPtr(playerid, player)
	player.FriendMoudle.SetPlayerPtr(playerid, player)
	player.RebelModule.SetPlayerPtr(playerid, player)
	player.MiningModule.SetPlayerPtr(playerid, player)
	player.HangMoudle.SetPlayerPtr(playerid, player)
	player.GuildModule.SetPlayerPtr(playerid, player)
	player.BlackMarketModule.SetPlayerPtr(playerid, player)
	player.ScoreMoudle.SetPlayerPtr(playerid, player)
	player.FameHallModule.SetPlayerPtr(playerid, player)
	player.TitleModule.SetPlayerPtr(playerid, player)
	player.FoodWarModule.SetPlayerPtr(playerid, player)
	player.ActivityModule.SetPlayerPtr(playerid, player)
	player.WanderMoudle.SetPlayerPtr(playerid, player)
	player.HeroSoulsModule.SetPlayerPtr(playerid, player)
	player.CamBattleModule.SetPlayerPtr(playerid, player)
	player.ChargeModule.SetPlayerPtr(playerid, player)

	player.playerid = playerid
	return
}

//响应玩家创建请求
func (player *TPlayer) OnCreate(playerid int32) {
	player.RoleMoudle.OnCreate(playerid)
	player.HeroMoudle.OnCreate(playerid)
	player.TaskMoudle.OnCreate(playerid)
	player.VipMoudle.OnCreate(playerid)
	player.MailMoudle.OnCreate(playerid)
	player.CopyMoudle.OnCreate(playerid)
	player.BagMoudle.OnCreate(playerid)
	player.StoreModule.OnCreate(playerid)
	player.SanGuoZhiModule.OnCreate(playerid)
	player.MallModule.OnCreate(playerid)
	player.SummonModule.OnCreate(playerid)
	player.ArenaModule.OnCreate(playerid)
	player.RobModule.OnCreate(playerid)
	player.SangokuMusouModule.OnCreate(playerid)
	player.AwardCenterModule.OnCreate(playerid)
	player.TerritoryModule.OnCreate(playerid)
	player.FriendMoudle.OnCreate(playerid)
	player.RebelModule.OnCreate(playerid)
	player.MiningModule.OnCreate(playerid)
	player.HangMoudle.OnCreate(playerid)
	player.GuildModule.OnCreate(playerid)
	player.BlackMarketModule.OnCreate(playerid)
	player.ScoreMoudle.OnCreate(playerid)
	player.FameHallModule.OnCreate(playerid)
	player.TitleModule.OnCreate(playerid)
	player.FoodWarModule.OnCreate(playerid)
	player.ActivityModule.OnCreate(playerid)
	player.WanderMoudle.OnCreate(playerid)
	player.HeroSoulsModule.OnCreate(playerid)
	player.CamBattleModule.OnCreate(playerid)
	player.ChargeModule.OnCreate(playerid)
}

//响应玩家的销毁请求
func (player *TPlayer) OnDestroy(playerid int32) {
	player.RoleMoudle.OnDestroy(playerid)
	player.HeroMoudle.OnDestroy(playerid)
	player.TaskMoudle.OnDestroy(playerid)
	player.VipMoudle.OnDestroy(playerid)
	player.MailMoudle.OnDestroy(playerid)
	player.CopyMoudle.OnDestroy(playerid)
	player.BagMoudle.OnDestroy(playerid)
	player.StoreModule.OnDestroy(playerid)
	player.SanGuoZhiModule.OnDestroy(playerid)
	player.MallModule.OnDestroy(playerid)
	player.SummonModule.OnDestroy(playerid)
	player.ArenaModule.OnDestroy(playerid)
	player.RobModule.OnDestroy(playerid)
	player.SangokuMusouModule.OnDestroy(playerid)
	player.AwardCenterModule.OnDestroy(playerid)
	player.TerritoryModule.OnDestroy(playerid)
	player.FriendMoudle.OnDestroy(playerid)
	player.RebelModule.OnDestroy(playerid)
	player.MiningModule.OnDestroy(playerid)
	player.HangMoudle.OnDestroy(playerid)
	player.GuildModule.OnDestroy(playerid)
	player.BlackMarketModule.OnDestroy(playerid)
	player.ScoreMoudle.OnDestroy(playerid)
	player.FameHallModule.OnDestroy(playerid)
	player.TitleModule.OnDestroy(playerid)
	player.FoodWarModule.OnDestroy(playerid)
	player.ActivityModule.OnDestroy(playerid)
	player.WanderMoudle.OnDestroy(playerid)
	player.HeroSoulsModule.OnDestroy(playerid)
	player.CamBattleModule.OnDestroy(playerid)
	player.ChargeModule.OnDestroy(playerid)

	player = nil
}

//响应玩家的上线请求
func (player *TPlayer) OnPlayerOnline(playerid int32) {
	player.RoleMoudle.OnPlayerOnline(playerid)
	player.HeroMoudle.OnPlayerOnline(playerid)
	player.TaskMoudle.OnPlayerOnline(playerid)
	player.VipMoudle.OnPlayerOnline(playerid)
	player.MailMoudle.OnPlayerOnline(playerid)
	player.CopyMoudle.OnPlayerOnline(playerid)
	player.BagMoudle.OnPlayerOnline(playerid)
	player.StoreModule.OnPlayerOnline(playerid)
	player.SanGuoZhiModule.OnPlayerOnline(playerid)
	player.MallModule.OnPlayerOnline(playerid)
	player.SummonModule.OnPlayerOnline(playerid)
	player.ArenaModule.OnPlayerOnline(playerid)
	player.RobModule.OnPlayerOnline(playerid)
	player.SangokuMusouModule.OnPlayerOnline(playerid)
	player.AwardCenterModule.OnPlayerOnline(playerid)
	player.TerritoryModule.OnPlayerOnline(playerid)
	player.FriendMoudle.OnPlayerOnline(playerid)
	player.RebelModule.OnPlayerOnline(playerid)
	player.MiningModule.OnPlayerOnline(playerid)
	player.HangMoudle.OnPlayerOnline(playerid)
	player.GuildModule.OnPlayerOnline(playerid)
	player.BlackMarketModule.OnPlayerOnline(playerid)
	player.ScoreMoudle.OnPlayerOnline(playerid)
	player.FameHallModule.OnPlayerOnline(playerid)
	player.TitleModule.OnPlayerOnline(playerid)
	player.FoodWarModule.OnPlayerOnline(playerid)
	player.ActivityModule.OnPlayerOnline(playerid)
	player.WanderMoudle.OnPlayerOnline(playerid)
	player.HeroSoulsModule.OnPlayerOnline(playerid)
	player.CamBattleModule.OnPlayerOnline(playerid)
	player.ChargeModule.OnPlayerOnline(playerid)
}

//响应玩家的下线请求
func (player *TPlayer) OnPlayerOffline(playerid int32) {
	player.RoleMoudle.OnPlayerOffline(playerid)
	player.HeroMoudle.OnPlayerOffline(playerid)
	player.TaskMoudle.OnPlayerOffline(playerid)
	player.VipMoudle.OnPlayerOffline(playerid)
	player.MailMoudle.OnPlayerOffline(playerid)
	player.CopyMoudle.OnPlayerOffline(playerid)
	player.BagMoudle.OnPlayerOffline(playerid)
	player.StoreModule.OnPlayerOffline(playerid)
	player.SanGuoZhiModule.OnPlayerOffline(playerid)
	player.MallModule.OnPlayerOffline(playerid)
	player.SummonModule.OnPlayerOffline(playerid)
	player.ArenaModule.OnPlayerOffline(playerid)
	player.RobModule.OnPlayerOffline(playerid)
	player.SangokuMusouModule.OnPlayerOffline(playerid)
	player.AwardCenterModule.OnPlayerOffline(playerid)
	player.TerritoryModule.OnPlayerOffline(playerid)
	player.FriendMoudle.OnPlayerOffline(playerid)
	player.RebelModule.OnPlayerOffline(playerid)
	player.MiningModule.OnPlayerOffline(playerid)
	player.HangMoudle.OnPlayerOffline(playerid)
	player.GuildModule.OnPlayerOffline(playerid)
	player.BlackMarketModule.OnPlayerOffline(playerid)
	player.ScoreMoudle.OnPlayerOffline(playerid)
	player.FameHallModule.OnPlayerOffline(playerid)
	player.TitleModule.OnPlayerOffline(playerid)
	player.FoodWarModule.OnPlayerOffline(playerid)
	player.ActivityModule.OnPlayerOffline(playerid)
	player.WanderMoudle.OnPlayerOffline(playerid)
	player.HeroSoulsModule.OnPlayerOffline(playerid)
	player.CamBattleModule.OnPlayerOffline(playerid)
	player.ChargeModule.OnPlayerOffline(playerid)

	G_SimpleMgr.Set_LogoffTime(playerid, time.Now().Unix())
}

//响应玩家的加载请求
func (player *TPlayer) OnPlayerLoad(playerid int32) {
	var wg sync.WaitGroup
	wg.Add(1)
	go player.RoleMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.HeroMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.TaskMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.VipMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.MailMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.CopyMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.BagMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.StoreModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.SanGuoZhiModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.MallModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.SummonModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.ArenaModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.RobModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.SangokuMusouModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.AwardCenterModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.TerritoryModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.FriendMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.RebelModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.MiningModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.HangMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.GuildModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.BlackMarketModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.ScoreMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.FameHallModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.TitleModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.FoodWarModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.ActivityModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.WanderMoudle.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.HeroSoulsModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.CamBattleModule.OnPlayerLoad(playerid, &wg)
	wg.Add(1)
	go player.ChargeModule.OnPlayerLoad(playerid, &wg)
	wg.Wait()

	player.InitModules(playerid)
}

//响应玩家的加载请求
func (player *TPlayer) OnPlayerLoadSync(playerid int32) {
	player.RoleMoudle.OnPlayerLoad(playerid, nil)
	player.HeroMoudle.OnPlayerLoad(playerid, nil)
	player.TaskMoudle.OnPlayerLoad(playerid, nil)
	player.VipMoudle.OnPlayerLoad(playerid, nil)
	player.MailMoudle.OnPlayerLoad(playerid, nil)
	player.CopyMoudle.OnPlayerLoad(playerid, nil)
	player.BagMoudle.OnPlayerLoad(playerid, nil)
	player.StoreModule.OnPlayerLoad(playerid, nil)
	player.SanGuoZhiModule.OnPlayerLoad(playerid, nil)
	player.MallModule.OnPlayerLoad(playerid, nil)
	player.SummonModule.OnPlayerLoad(playerid, nil)
	player.ArenaModule.OnPlayerLoad(playerid, nil)
	player.RobModule.OnPlayerLoad(playerid, nil)
	player.SangokuMusouModule.OnPlayerLoad(playerid, nil)
	player.AwardCenterModule.OnPlayerLoad(playerid, nil)
	player.TerritoryModule.OnPlayerLoad(playerid, nil)
	player.FriendMoudle.OnPlayerLoad(playerid, nil)
	player.RebelModule.OnPlayerLoad(playerid, nil)
	player.MiningModule.OnPlayerLoad(playerid, nil)
	player.HangMoudle.OnPlayerLoad(playerid, nil)
	player.GuildModule.OnPlayerLoad(playerid, nil)
	player.BlackMarketModule.OnPlayerLoad(playerid, nil)
	player.ScoreMoudle.OnPlayerLoad(playerid, nil)
	player.FameHallModule.OnPlayerLoad(playerid, nil)
	player.TitleModule.OnPlayerLoad(playerid, nil)
	player.FoodWarModule.OnPlayerLoad(playerid, nil)
	player.ActivityModule.OnPlayerLoad(playerid, nil)
	player.WanderMoudle.OnPlayerLoad(playerid, nil)
	player.HeroSoulsModule.OnPlayerLoad(playerid, nil)
	player.CamBattleModule.OnPlayerLoad(playerid, nil)
	player.ChargeModule.OnPlayerLoad(playerid, nil)
	player.InitModules(playerid)
}

var GMap_HandleMsg_Lock map[string]bool

func (player *TPlayer) Lock(url string) {
	if val, ok := GMap_HandleMsg_Lock[url]; ok && val {
		player.mutex.Lock()
		player.isLock = true
	}
}
func (player *TPlayer) Unlock() {
	if player.isLock {
		player.isLock = false
		player.mutex.Unlock()
	}
}

//计算战力
func (player *TPlayer) CalcFightValue() int {
	oldValue := player.pSimpleInfo.FightValue
	value := player.HeroMoudle.CalcFightValue(nil)
	if true == G_SimpleMgr.Set_FightValue(player.playerid, value, player.GetLevel()) {
		G_FightRanker.SetRankItemEx(player.playerid, oldValue, value)
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FIGHT_VALUE, value)
	}
	return value
}

//计算战力
func (player *TPlayer) GetFightValue() int {
	if player.pSimpleInfo == nil {
		gamelog.Error("GetFightValue Error pSimpleInfo is nil :%d", player.playerid)
		return G_SimpleMgr.Get_FightValue(player.playerid)
	}

	return player.pSimpleInfo.FightValue
}

//玩家初始化
func (player *TPlayer) SetPlayerName(name string) {
	player.RoleMoudle.Name = name
	return
}

//玩家初始化
func (player *TPlayer) SetMainHeroID(heroid int) {
	player.HeroMoudle.CurHeros[0].Init(heroid)
	return
}

//获取主英雄等级
func (player *TPlayer) GetLevel() int {
	if player.HeroMoudle.CurHeros[0].ID == 0 {
		gamelog.Error("GetMainHeroLevel Error :HeroID is 0 ")
		return 0
	}

	return player.HeroMoudle.CurHeros[0].Level
}

//获取角色的VIP等级
func (player *TPlayer) GetVipLevel() int8 {
	return player.RoleMoudle.VipLevel
}

//获取角色的VIP经验
func (player *TPlayer) GetVipExp() int {
	return player.RoleMoudle.GetMoney(gamedata.VipExpMoneyID)
}

//判断玩家今天是否己登录
func (player *TPlayer) IsTodayLogin() bool {
	if utility.GetCurDay() == player.pSimpleInfo.LoginDay {
		return true
	}

	return false
}

func (player *TPlayer) IsHasHero(heroid int) bool {
	for _, hero := range player.HeroMoudle.CurHeros {
		if hero.ID == heroid {
			return true
		}
	}

	for _, hero := range player.HeroMoudle.BackHeros {
		if hero.ID == heroid {
			return true
		}
	}

	for _, hero := range player.BagMoudle.HeroBag.Heros {
		if hero.ID == heroid {
			return true
		}
	}

	return false
}

func (player *TPlayer) GetHeroByPos(postype int, pos int) *THeroData {
	if pos < 0 {
		gamelog.Error("GetHeroByPos Error : Invalid pos :%d", pos)
		return nil
	}

	if postype == POSTYPE_BATTLE {
		if pos < len(player.HeroMoudle.CurHeros) {
			return &player.HeroMoudle.CurHeros[pos]
		}
	} else if postype == POSTYPE_BACK {
		if pos < len(player.HeroMoudle.BackHeros) {
			return &player.HeroMoudle.BackHeros[pos]
		}
	} else if postype == POSTYPE_BAG {
		if pos < len(player.BagMoudle.HeroBag.Heros) {
			return &player.BagMoudle.HeroBag.Heros[pos]
		}
	}

	gamelog.Error("GetHeroByPos Error : Invalid pos :%d", pos)
	return nil
}

//响应充值人民币
func (player *TPlayer) HandChargeRenMinBi(RenMinBi int, chargeid int) bool {
	//普通充值
	if chargeid >= 1 && chargeid <= 9 {
		pChargeInfo := gamedata.GetChargeItem(chargeid)
		if pChargeInfo == nil {
			gamelog.Error("OnChargeMoney Error : Invalid chargeid :%d", chargeid)
			return false
		}

		player.ChargeModule.AddChargeTimes(pChargeInfo.ID)
		player.RoleMoudle.AddMoney(gamedata.ChargeMoneyID, pChargeInfo.RenMinBi*gamedata.ChargeMoneyRatio)

		// 给充值奖励
		var awardID int
		if player.ChargeModule.IsFirstCharge(pChargeInfo.ID) {
			awardID = pChargeInfo.FirstAwardID
		} else {
			awardID = pChargeInfo.AwardID
		}
		items := gamedata.GetItemsFromAwardID(awardID)
		player.BagMoudle.AddAwardItems(items)
		//! 发放通知邮件
		SendRechargeMail(player.playerid, RenMinBi)
	} else if chargeid == 10 || chargeid == 11 {
		cardId := chargeid - 9
		player.ActivityModule.CheckReset()
		pMonthCard := gamedata.GetMonthCardInfo(cardId)
		if pMonthCard == nil {
			gamelog.Error("OnChargeMoney Error : Invalid Cardid :%d", cardId)
			return false
		}

		if player.ActivityModule.MonthCard.CardDays[cardId] != 0 {
			gamelog.Error("OnChargeMoney Error : Repeat purchase")
			return false
		}

		player.ActivityModule.MonthCard.CardDays[cardId] += 30
		go player.ActivityModule.MonthCard.DB_UpdateCardDays(cardId, player.ActivityModule.MonthCard.CardDays[cardId])
	}

	player.RoleMoudle.AddVipExp(RenMinBi * gamedata.ChargeMoneyRatio)

	player.OnChargeMoney(RenMinBi)

	return true
}

func (player *TPlayer) OnChargeMoney(rmb int) {
	//! 增加任务/七天/限时完成进度
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_RECHARGE, rmb*gamedata.ChargeMoneyRatio)

	//! 增加单笔充值
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SINGLE_RECHARGE, rmb)

	//! 检测月基金激活
	player.ActivityModule.MonthFund.SetMonthFund(rmb)

	//! 周周盈
	player.ActivityModule.WeekAward.AddRechargeNum(rmb)

	//! 首充/次充
	player.ActivityModule.FirstRecharge.CheckRecharge(rmb)

	//! 单充/累充
	player.ActivityModule.AddRechargeValue(rmb)

	//! 增加豪华充值状态改变
	if rmb >= 6 {
		player.ActivityModule.Sign.SetSignPlusStatus()
	}
}

func OnConfigChange(tbname string) bool {
	switch tbname {
	case "type_activity":
		{
			//! 获取今日开启活动
			openDay := GetOpenServerDay()
			activityLength := len(G_GlobalVariables.ActivityLst)
			for _, v := range gamedata.GT_ActivityLst {
				if v.ID == 0 {
					continue
				}
				//! 遍历当前活动表
				isExist := false
				for j := 0; j < activityLength; j++ {
					if v.ID == G_GlobalVariables.ActivityLst[j].ActivityID {
						//! 活动已存在, 改变状态
						G_GlobalVariables.ActivityLst[j].Status = v.Status
						G_GlobalVariables.ActivityLst[j].award = v.AwardType
						beginTime, endTime := gamedata.GetActivityEndTime(v.ID, openDay)
						G_GlobalVariables.ActivityLst[j].beginTime = beginTime
						G_GlobalVariables.ActivityLst[j].endTime = endTime
						isExist = true
						G_GlobalVariables.DB_UpdateActivityInfo(j)
						break
					}
				}

				if isExist == false {
					//! 新加活动
					if v.ActivityType == gamedata.Activity_Seven {
						seven := TSevenDayBuyInfo{}
						seven.ActivityID = v.ID
						G_GlobalVariables.SevenDayLimit = append(G_GlobalVariables.SevenDayLimit, seven)
						G_GlobalVariables.DB_AddSevenDayBuyInfo(seven)
					}

					var activity TActivityData
					activity.ActivityID = v.ID
					activity.activityType = v.ActivityType
					activity.award = v.AwardType
					activity.beginTime, activity.endTime = gamedata.GetActivityEndTime(v.ID, openDay)
					activity.VersionCode = 0
					activity.Status = v.Status
					activity.ResetCode = 0
					G_GlobalVariables.ActivityLst = append(G_GlobalVariables.ActivityLst, activity)
					G_GlobalVariables.DB_AddNewActivity(activity)
				}
			}
		}
	default:
		{
			gamelog.Error("OnConfigChange Error: Table %s is not processed!", tbname)
		}

	}
	return true
}

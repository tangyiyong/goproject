package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TAction struct {
	Value     int
	StartTime int64
}

//角色基本数据表结构
type TRoleMoudle struct {
	PlayerID    int32     `bson:"_id"` //玩家ID
	Name        string    //玩家角色名
	Actions     []TAction //活力值 ....
	Moneys      []int     //货币集 1: 金币 2: 银币 ...
	VipLevel    int       //Vip等级
	NewWizard   string    //新手向导
	ExpIncLvl   int       //经验加成等级
	TodayCharge int       //今天的充值额度
	TotalCharge int       //总的充值额度
	ownplayer   *TPlayer  //父player指针
}

func (role *TRoleMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	role.PlayerID = playerid
	role.ownplayer = pPlayer
}

func (role *TRoleMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	role.PlayerID = playerid

	role.Actions = make([]TAction, gamedata.GetActionCount())
	for i := 0; i < len(role.Actions); i++ {
		pActionInfo := gamedata.GetActionInfo(i + 1)
		if pActionInfo == nil {
			gamelog.Error("TRoleMoudle:OnCreate Error: invalid actionid %d", i)
			return
		}

		role.Actions[i].Value = pActionInfo.Max
		role.Actions[i].StartTime = 0
	}

	role.Moneys = make([]int, gamedata.GetMoneyCount())

	//创建数据库记录
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerRole", role)
}

//玩家对象销毁
func (role *TRoleMoudle) OnDestroy(playerid int32) {
	role = nil
}

//玩家进入游戏
func (role *TRoleMoudle) OnPlayerOnline(playerid int32) {
}

//OnPlayerOffline 玩家离开游戏
func (role *TRoleMoudle) OnPlayerOffline(playerid int32) {
}

//玩家离开游戏
func (role *TRoleMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	var bRet = true
	err := s.DB(appconfig.GameDbName).C("PlayerRole").Find(bson.M{"_id": playerid}).One(role)
	if err != nil {
		gamelog.Error("PlayerRole Load Error :%s， PlayerID: %d", err.Error(), playerid)
		bRet = false
	}

	for i := 0; i < len(role.Actions); i++ {
		pActionInfo := gamedata.GetActionInfo(i + 1)
		if pActionInfo != nil {
			if role.Actions[i].Value >= pActionInfo.Max {
				role.Actions[i].StartTime = 0
			}
		}
	}

	if wg != nil {
		wg.Done()
	}
	role.PlayerID = playerid
	return bRet
}

//扣除货币， 如果返回成功，就是扣除成功， 如果返回失败，就是货币不足
func (role *TRoleMoudle) CostMoney(moneyID int, moneyNum int) bool {
	if (moneyID <= 0) || (moneyID > len(role.Moneys)) {
		gamelog.Error("CostMoney Error: Inavlid moneyID :%d", moneyID)
		return false
	}

	if moneyNum <= 0 {
		gamelog.Error("CostMoney Error : Invalid moneyNum :%d", moneyNum)
		return false
	}

	if role.Moneys[moneyID-1] < moneyNum {
		gamelog.Error("CostMoney Error : Not Enough Money :%d", moneyNum)
		return false
	}

	role.Moneys[moneyID-1] -= moneyNum
	role.DB_SaveMoneysAt(moneyID)

	//! 增加任务进度
	if moneyID == 1 {
		role.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SPENT_MONEY, moneyNum)
	}

	return true
}

func (role *TRoleMoudle) CheckMoneyEnough(moneyID int, moneyNum int) bool {
	if (moneyID <= 0) || (moneyID > len(role.Moneys)) {
		gamelog.Error("CheckMoneyEnough Error: Inavlid moneyID :%d", moneyID)
		return false
	}

	if moneyNum <= 0 {
		gamelog.Error("CheckMoneyEnough Error : Invalid moneyNum :%d", moneyNum)
		return false
	}

	if role.Moneys[moneyID-1] >= moneyNum {
		return true
	}

	return false
}

func (role *TRoleMoudle) GetMoney(moneyID int) int {
	if (moneyID <= 0) || (moneyID > len(role.Moneys)) {
		gamelog.Error("GetMoney Error: Inavlid moneyID :%d", moneyID)
		return 0
	}

	return role.Moneys[moneyID-1]
}

func (role *TRoleMoudle) AddMoney(moneyID int, moneyNum int) int {
	if moneyNum <= 0 {
		gamelog.Error("AddMoney Error: Inavlid moneyNum :%d", moneyNum)
		return role.Moneys[moneyID-1]
	}

	if (moneyID <= 0) || (moneyID > len(role.Moneys)) {
		gamelog.Error("AddMoney Error: Inavlid moneyID :%d", moneyID)
		return 0
	}

	role.Moneys[moneyID-1] += moneyNum
	if role.Moneys[moneyID-1] > gamedata.GetMoneyMaxValue(moneyID) {
		role.Moneys[moneyID-1] = gamedata.GetMoneyMaxValue(moneyID)
	}

	role.DB_SaveMoneysAt(moneyID)
	return role.Moneys[moneyID-1]
}

//扣除行动力， 如果返回成功，就是扣除成功， 如果返回失败，就是行动力不足
func (role *TRoleMoudle) CostAction(actionID int, actionNum int) bool {
	if (actionID <= 0) || (actionID >= len(role.Actions)) {
		gamelog.Error("CostAction Error: Inavlid actionID :%d", actionID)
		return false
	}

	if actionNum <= 0 {
		gamelog.Error("CostAction Error: Inavlid actionNum :%d", actionNum)
		return false
	}

	if role.Actions[actionID-1].Value < actionNum {
		return false
	}

	pActionInfo := gamedata.GetActionInfo(actionID)
	if pActionInfo == nil {
		gamelog.Error("CostAction Invalid Action id :%d", actionID)
		return false
	}

	role.Actions[actionID-1].Value -= actionNum

	if role.Actions[actionID-1].Value < pActionInfo.Max {
		if role.Actions[actionID-1].StartTime <= 0 {
			role.Actions[actionID-1].StartTime = time.Now().Unix()
		}
	} else {
		role.Actions[actionID-1].StartTime = 0
	}

	role.DB_SaveActionsAt(actionID)

	return true
}

func (role *TRoleMoudle) CheckActionEnough(actionID int, actionNum int) bool {
	if (actionID <= 0) || (actionID > len(role.Actions)) {
		gamelog.Error("CheckActionEnough Error: Inavlid actionID :%d", actionID)
		return false
	}

	if actionNum <= 0 {
		gamelog.Error("CheckActionEnough Error: Inavlid actionNum :%d", actionNum)
		return false
	}

	if role.Actions[actionID-1].Value >= actionNum {
		return true
	}

	if role.UpdateAction(actionID) {
		role.DB_SaveActionsAt(actionID)
	}

	if role.Actions[actionID-1].Value < actionNum {
		return false
	}

	return true
}

func (role *TRoleMoudle) GetActionData(actionID int) (int, int64) {
	if (actionID <= 0) || (actionID > len(role.Actions)) {
		gamelog.Error("GetAction Error: Inavlid actionID :%d", actionID)
		return 0, 0
	}

	return role.Actions[actionID-1].Value, role.Actions[actionID-1].StartTime
}

func (role *TRoleMoudle) GetAction(actionID int) int {
	if (actionID <= 0) || (actionID > len(role.Actions)) {
		gamelog.Error("GetAction Error: Inavlid actionID :%d", actionID)
		return 0
	}

	if role.UpdateAction(actionID) {
		role.DB_SaveActionsAt(actionID)
	}

	return role.Actions[actionID-1].Value
}

func (role *TRoleMoudle) AddAction(actionID int, actionNum int) int {
	if (actionID <= 0) || (actionID > len(role.Actions)) {
		gamelog.Error("AddAction Error: Inavlid actionID :%d", actionID)
		return 0
	}

	role.UpdateAction(actionID)

	role.Actions[actionID-1].Value += actionNum

	pActionInfo := gamedata.GetActionInfo(actionID)
	if pActionInfo == nil {
		gamelog.Error("AddAction Invalid Action id :%d", actionID)
		return 0
	}

	if role.Actions[actionID-1].Value >= pActionInfo.Max {
		role.Actions[actionID-1].StartTime = 0
	}

	role.DB_SaveActionsAt(actionID)

	return role.Actions[actionID-1].Value
}

func (role *TRoleMoudle) UpdateAction(actionID int) bool {
	pActionInfo := gamedata.GetActionInfo(actionID)
	if pActionInfo == nil {
		gamelog.Error("UpdateAction Invalid Action id :%d", actionID)
		return false
	}

	if role.Actions[actionID-1].Value >= pActionInfo.Max {
		if role.Actions[actionID-1].StartTime > 0 {
			gamelog.Error("UpdateAction error  StartTime is not 0")
		}
		role.Actions[actionID-1].StartTime = 0
		return false
	}

	if role.Actions[actionID-1].StartTime <= 0 {
		gamelog.Error("UpdateAction error  action not max, but starttime is 0")
	}

	timeElapse := time.Now().Unix() - role.Actions[actionID-1].StartTime

	if timeElapse < int64(pActionInfo.UnitTime) {
		return false
	}

	ActionNum := int(timeElapse) / pActionInfo.UnitTime
	role.Actions[actionID-1].Value += ActionNum

	if role.Actions[actionID-1].Value >= pActionInfo.Max {
		role.Actions[actionID-1].Value = pActionInfo.Max
		role.Actions[actionID-1].StartTime = 0
	} else {
		role.Actions[actionID-1].StartTime = role.Actions[actionID-1].StartTime + int64(ActionNum*pActionInfo.UnitTime)
	}

	return true
}

func (role *TRoleMoudle) UpdateAllAction() {
	var bUpdate = false
	for i := 0; i < len(role.Actions); i++ {
		if role.UpdateAction(i + 1) {
			bUpdate = true
		}
	}

	if bUpdate {
		role.DB_SaveActions()
	}

	return
}

//增加公会经验技能等级
func (role *TRoleMoudle) AddGuildSkillExpIncLevel() bool {
	role.ExpIncLvl += 1
	role.DB_SaveExpIncLevel()
	return true
}

//清空公会经验技能等级
func (role *TRoleMoudle) ClearGuildSkillExpIncLevel() {
	role.ExpIncLvl = 0
	role.DB_SaveExpIncLevel()
}

//! 增加VIP经验
func (role *TRoleMoudle) AddVipExp(exp int) {
	role.AddMoney(gamedata.VipExpMoneyID, exp)
	newLevel := gamedata.CalcVipLevelByExp(role.GetMoney(gamedata.VipExpMoneyID), role.VipLevel)
	if newLevel != role.VipLevel {
		role.VipLevel = newLevel
		role.ownplayer.ActivityModule.VipGift.IsRecvWelfare = false
		role.DB_SaveVipLevel()
	}
}

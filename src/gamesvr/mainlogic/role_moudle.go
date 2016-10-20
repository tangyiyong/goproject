package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TAction struct {
	Value int
	Time  int32
}

//角色基本数据表结构
type TRoleMoudle struct {
	PlayerID    int32     `bson:"_id"` //玩家ID
	Name        string    //玩家角色名
	Actions     []TAction //活力值 ....
	Moneys      []int     //货币集 1: 金币 2: 银币 ...
	VipLevel    int8      //Vip等级
	NewWizard   string    //新手向导
	TodayCharge int32     //今天的充值额度
	TotalCharge int32     //总的充值额度
	CurStarID   int32     //三国志ID
	ownplayer   *TPlayer  //父player指针
}

func (self *TRoleMoudle) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *TRoleMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	self.PlayerID = playerid

	self.Actions = make([]TAction, gamedata.GetActionCount())
	for i := 0; i < len(self.Actions); i++ {
		pActionInfo := gamedata.GetActionInfo(i + 1)
		if pActionInfo == nil {
			gamelog.Error("TRoleMoudle:OnCreate Error: invalid actionid %d", i)
			return
		}

		self.Actions[i].Value = pActionInfo.Max
		self.Actions[i].Time = 0
	}

	self.Moneys = make([]int, gamedata.GetMoneyCount())

	//创建数据库记录
	mongodb.InsertToDB("PlayerRole", self)
}

//玩家对象销毁
func (self *TRoleMoudle) OnDestroy(playerid int32) {
	self = nil
}

//玩家进入游戏
func (self *TRoleMoudle) OnPlayerOnline(playerid int32) {
}

//OnPlayerOffline 玩家离开游戏
func (self *TRoleMoudle) OnPlayerOffline(playerid int32) {
}

//玩家离开游戏
func (self *TRoleMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) bool {
	s := mongodb.GetDBSession()
	defer s.Close()
	var bRet = true
	err := s.DB(appconfig.GameDbName).C("PlayerRole").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerRole Load Error :%s， PlayerID: %d", err.Error(), playerid)
		bRet = false
	}

	for i := 0; i < len(self.Actions); i++ {
		pActionInfo := gamedata.GetActionInfo(i + 1)
		if pActionInfo != nil {
			if self.Actions[i].Value >= pActionInfo.Max {
				self.Actions[i].Time = 0
			}
		}
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
	return bRet
}

//扣除货币， 如果返回成功，就是扣除成功， 如果返回失败，就是货币不足
func (self *TRoleMoudle) CostMoney(moneyID int, moneyNum int) bool {
	if (moneyID <= 0) || (moneyID > len(self.Moneys)) {
		gamelog.Error("CostMoney Error: Inavlid moneyID :%d", moneyID)
		return false
	}

	if moneyNum <= 0 {
		gamelog.Error("CostMoney Error : Invalid moneyNum :%d", moneyNum)
		return false
	}

	if self.Moneys[moneyID-1] < moneyNum {
		gamelog.Error("CostMoney Error : Not Enough Money :%d", moneyNum)
		return false
	}

	self.Moneys[moneyID-1] -= moneyNum
	self.DB_SaveMoneysAt(moneyID)

	//! 增加任务进度
	if moneyID == 1 {
		self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SPENT_MONEY, moneyNum)
	}

	return true
}

func (self *TRoleMoudle) CheckMoneyEnough(moneyID int, moneyNum int) bool {
	if (moneyID <= 0) || (moneyID > len(self.Moneys)) {
		gamelog.Error("CheckMoneyEnough Error: Inavlid moneyID :%d", moneyID)
		return false
	}

	if moneyNum <= 0 {
		gamelog.Error("CheckMoneyEnough Error : Invalid moneyNum :%d", moneyNum)
		return false
	}

	if self.Moneys[moneyID-1] >= moneyNum {
		return true
	}

	return false
}

func (self *TRoleMoudle) GetMoney(moneyID int) int {
	if (moneyID <= 0) || (moneyID > len(self.Moneys)) {
		gamelog.Error("GetMoney Error: Inavlid moneyID :%d", moneyID)
		return 0
	}

	return self.Moneys[moneyID-1]
}

func (self *TRoleMoudle) AddMoney(moneyID int, moneyNum int) int {
	if moneyNum <= 0 {
		gamelog.Error("AddMoney Error: Inavlid moneyNum :%d", moneyNum)
		return self.Moneys[moneyID-1]
	}

	if (moneyID <= 0) || (moneyID > len(self.Moneys)) {
		gamelog.Error("AddMoney Error: Inavlid moneyID :%d", moneyID)
		return 0
	}

	self.Moneys[moneyID-1] += moneyNum
	if self.Moneys[moneyID-1] > gamedata.GetMoneyMaxValue(moneyID) {
		self.Moneys[moneyID-1] = gamedata.GetMoneyMaxValue(moneyID)
	}

	self.DB_SaveMoneysAt(moneyID)
	return self.Moneys[moneyID-1]
}

//扣除行动力， 如果返回成功，就是扣除成功， 如果返回失败，就是行动力不足
func (self *TRoleMoudle) CostAction(actionID int, actionNum int) bool {
	if (actionID <= 0) || (actionID >= len(self.Actions)) {
		gamelog.Error("CostAction Error: Inavlid actionID :%d", actionID)
		return false
	}

	if actionNum <= 0 {
		gamelog.Error("CostAction Error: Inavlid actionNum :%d", actionNum)
		return false
	}

	if self.Actions[actionID-1].Value < actionNum {
		return false
	}

	pActionInfo := gamedata.GetActionInfo(actionID)
	if pActionInfo == nil {
		gamelog.Error("CostAction Invalid Action id :%d", actionID)
		return false
	}

	self.Actions[actionID-1].Value -= actionNum

	if self.Actions[actionID-1].Value < pActionInfo.Max {
		if self.Actions[actionID-1].Time <= 0 {
			self.Actions[actionID-1].Time = utility.GetCurTime()
		}
	} else {
		self.Actions[actionID-1].Time = 0
	}

	self.DB_SaveActionsAt(actionID)

	return true
}

func (self *TRoleMoudle) CheckActionEnough(actionID int, actionNum int) bool {
	if (actionID <= 0) || (actionID > len(self.Actions)) {
		gamelog.Error("CheckActionEnough Error: Inavlid actionID :%d", actionID)
		return false
	}

	if actionNum <= 0 {
		gamelog.Error("CheckActionEnough Error: Inavlid actionNum :%d", actionNum)
		return false
	}

	if self.Actions[actionID-1].Value >= actionNum {
		return true
	}

	if self.UpdateAction(actionID) {
		self.DB_SaveActionsAt(actionID)
	}

	if self.Actions[actionID-1].Value < actionNum {
		return false
	}

	return true
}

func (self *TRoleMoudle) GetActionData(actionID int) (int, int32) {
	if (actionID <= 0) || (actionID > len(self.Actions)) {
		gamelog.Error("GetAction Error: Inavlid actionID :%d", actionID)
		return 0, 0
	}

	return self.Actions[actionID-1].Value, self.Actions[actionID-1].Time
}

func (self *TRoleMoudle) GetAction(actionID int) int {
	if (actionID <= 0) || (actionID > len(self.Actions)) {
		gamelog.Error("GetAction Error: Inavlid actionID :%d", actionID)
		return 0
	}

	if self.UpdateAction(actionID) {
		self.DB_SaveActionsAt(actionID)
	}

	return self.Actions[actionID-1].Value
}

func (self *TRoleMoudle) AddAction(actionID int, actionNum int) int {
	if (actionID <= 0) || (actionID > len(self.Actions)) {
		gamelog.Error("AddAction Error: Inavlid actionID :%d", actionID)
		return 0
	}

	self.UpdateAction(actionID)

	self.Actions[actionID-1].Value += actionNum

	pActionInfo := gamedata.GetActionInfo(actionID)
	if pActionInfo == nil {
		gamelog.Error("AddAction Invalid Action id :%d", actionID)
		return 0
	}

	if self.Actions[actionID-1].Value >= pActionInfo.Max {
		self.Actions[actionID-1].Time = 0
	}

	self.DB_SaveActionsAt(actionID)

	return self.Actions[actionID-1].Value
}

func (self *TRoleMoudle) UpdateAction(actionID int) bool {
	pActionInfo := gamedata.GetActionInfo(actionID)
	if pActionInfo == nil {
		gamelog.Error("UpdateAction Invalid Action id :%d", actionID)
		return false
	}

	if self.Actions[actionID-1].Value >= pActionInfo.Max {
		if self.Actions[actionID-1].Time > 0 {
			gamelog.Error("UpdateAction error  StartTime is not 0")
		}
		self.Actions[actionID-1].Time = 0
		return false
	}

	if self.Actions[actionID-1].Time <= 0 {
		gamelog.Error("UpdateAction error  action not max, but starttime is 0")
	}

	timeElapse := utility.GetCurTime() - self.Actions[actionID-1].Time

	if timeElapse < int32(pActionInfo.UnitTime) {
		return false
	}

	ActionNum := int(timeElapse) / pActionInfo.UnitTime
	self.Actions[actionID-1].Value += ActionNum

	if self.Actions[actionID-1].Value >= pActionInfo.Max {
		self.Actions[actionID-1].Value = pActionInfo.Max
		self.Actions[actionID-1].Time = 0
	} else {
		self.Actions[actionID-1].Time = self.Actions[actionID-1].Time + int32(ActionNum*pActionInfo.UnitTime)
	}

	return true
}

func (self *TRoleMoudle) UpdateAllAction() {
	var bUpdate = false
	for i := 0; i < len(self.Actions); i++ {
		if self.UpdateAction(i + 1) {
			bUpdate = true
		}
	}

	if bUpdate {
		self.DB_SaveActions()
	}

	return
}

//! 增加VIP经验
func (self *TRoleMoudle) AddVipExp(exp int) {
	self.AddMoney(gamedata.VipExpMoneyID, exp)
	newLevel := gamedata.CalcVipLevelByExp(self.GetMoney(gamedata.VipExpMoneyID), self.VipLevel)
	if newLevel != self.VipLevel {
		self.VipLevel = newLevel
		self.ownplayer.ActivityModule.VipGift.IsRecvWelfare = false
		self.DB_SaveVipLevel()
	}
}

func (self *TRoleMoudle) RedTip() bool {
	//! 判断升星材料是否足够
	info := gamedata.GetSanGuoZhiInfo(int(self.CurStarID))
	if info == nil {
		return false
	}

	bEnough := self.ownplayer.BagMoudle.IsItemEnough(info.CostType, info.CostNum)
	return bEnough
}

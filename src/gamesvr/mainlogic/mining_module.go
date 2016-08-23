package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math"
	"mongodb"
	"msg"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TMiningMonster struct {
	ID    int
	Index int
	Life  int
}

type TMiningPos struct {
	X int
	Y int
}

type TMiningBuff struct {
	BuffType int
	Value    int
	Times    int
}

type TMiningBossAward struct { //! Boss打完后的翻牌奖励
	ID      int
	ItemID  int
	ItemNum int
	Status  bool
}

type IndexLst []int

func (self *IndexLst) Set(index int, value int) {
	for i, v := range *self {
		if (v>>16)&0x0000FFFF == index {
			(*self)[i] = index << 16
			(*self)[i] += value
			return
		}
	}

	var mark int
	mark = index << 16
	mark += value
	*self = append(*self, mark)
}

func (self *IndexLst) Get(index int) int {
	for i, v := range *self {
		if v>>16 == index {
			return (*self)[i] & 0x0000FFFF
		}
	}

	return 0
}

func (self *IndexLst) GetIndex(value int) int {
	for i, v := range *self {
		if v&0x0000FFFF == value {
			return i
		}
	}

	return -1
}

type MapLst [60]uint64

func (self *MapLst) Set(index int) {
	y := uint64(index / 60)
	x := uint64(index % 60)
	(*self)[y] |= (1 << x)
}

func (self *MapLst) Get(index int) bool {
	y := uint64(index / 60)
	x := uint64(index % 60)

	return ((*self)[y] & (1 << x)) > 0
}

func (self *MapLst) Count() int {
	count := 0
	for j := 0; j < gamedata.MiningMapLength; j++ {
		for i := 0; i < 60; i++ {
			if (*self).Get(j*60+i) == true {
				count++
			}
		}
	}
	return count
}

type TMiningModule struct {
	PlayerID int32 `bson:"_id"`

	GuaJiType     int   //! 当前挂机类型
	GuajiCalcTime int64 //! 挂机结算时间
	Point         int   //! 玩家当前积分

	Buff TMiningBuff //! 玩家Buff

	MonsterLst []TMiningMonster //! 怪物信息

	MiningMap MapLst   //! 地图挖掘标记 位运算
	Element   IndexLst //! 前16位为index  后16位为值

	LastPos TMiningPos //! 最后操作位置

	ActionBuyTimes int    //! 购买行动力次数
	ResetDay       uint32 //! 重置购买次数天数

	BossAward []TMiningBossAward //! 翻牌奖励

	BlackMarketBuyMark IntLst //! 黑市购买标记

	StatusCode int //! 状态码

	MiningResetTimes int

	ownplayer *TPlayer
}

//! 设置玩家指针
func (self *TMiningModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TMiningModule) OnCreate(playerid int32) {

	//! 创建地图
	self.CreateNewMap(false)

	//! 初始化状态码
	self.StatusCode = 1 + (1 << 16)

	//! 初始化入口坐标
	self.LastPos.X = gamedata.MiningEnterPointX
	self.LastPos.Y = gamedata.MiningEnterPointY

	//! 初始化地图入口
	self.MiningMap.Set(gamedata.MiningEnterPointY*60 + gamedata.MiningEnterPointX)

	//! 重置矿洞次数
	self.MiningResetTimes = 0

	//! 设置重置购买次数时间
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerMining", self)
}

//! 玩家销毁角色
func (self *TMiningModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TMiningModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TMiningModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TMiningModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerMining").Find(bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerMining Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TMiningModule) RedTip() bool {
	now := time.Now().Unix()
	if self.GuajiCalcTime < now {
		return true
	}

	pActionInfo := gamedata.GetActionInfo(gamedata.MiningCostActionID)
	if pActionInfo == nil {
		gamelog.Error("GetActionInfo Invalid Action id :%d", gamedata.MiningCostActionID)
		return false
	}

	action := self.ownplayer.RoleMoudle.GetAction(gamedata.MiningCostActionID)
	if action >= pActionInfo.Max { //! 挖矿精力已满
		return true
	}

	return false
}

//! 重置挖矿地图
func (self *TMiningModule) ResetMiningMap() {
	self.MiningResetTimes++

	//! 标记重置
	self.BossAward = []TMiningBossAward{} //! 翻牌奖励
	self.Buff = TMiningBuff{}
	self.Point = 0

	self.Element = IndexLst{}
	self.MonsterLst = []TMiningMonster{}
	self.StatusCode = 1 + (1 << 16)

	self.CreateNewMap(true)

	//! 初始化状态码
	self.StatusCode = 1 + (1 << 16)

	//! 初始化入口坐标
	self.LastPos.X = gamedata.MiningEnterPointX
	self.LastPos.Y = gamedata.MiningEnterPointY

	//! 初始化地图入口
	self.MiningMap.Set(gamedata.MiningEnterPointY*60 + gamedata.MiningEnterPointX)

	//! 设置重置购买次数时间
	self.ResetDay = utility.GetCurDay()

	go self.DB_ResetMapInfo()
}

//! 删除元素
func (self *TMiningModule) DeleteElement(index int) {
	pos := 0
	value := 0
	for i, v := range self.Element {
		if ((v >> 16) & 0x0000FFFF) == index {
			pos = i
			value = v
			break
		}
	}

	if pos == 0 {
		self.Element = self.Element[1:]
	} else if (pos + 1) == len(self.Element) {
		self.Element = self.Element[:pos]
	} else {
		self.Element = append(self.Element[:pos], self.Element[pos+1:]...)
	}

	go self.DB_RemoveElement(value)
}

//! 删除怪物
func (self *TMiningModule) DeleteMonster(index int) {
	pos := 0
	monster := TMiningMonster{}
	for i, v := range self.MonsterLst {
		if v.Index == index {
			pos = i
			monster = v
			break
		}
	}

	if pos == 0 {
		self.MonsterLst = self.MonsterLst[1:]
	} else if (pos + 1) == len(self.MonsterLst) {
		self.MonsterLst = self.MonsterLst[:pos]
	} else {
		self.MonsterLst = append(self.MonsterLst[:pos], self.MonsterLst[pos+1:]...)
	}

	go self.DB_RemoveMonster(monster)
}

//! 生成新可视区域部分坐标
func (self *TMiningModule) GetNewVisualPosArena(x int, y int) (posLst []TMiningPos) {
	visualLst := self.GetVisualPosArena(x, y)

	//! 去除已经可以看见的区域
	for _, v := range visualLst {
		index := v.Y*gamedata.MiningMapLength + v.X

		if v.X >= gamedata.MiningMapLength || v.Y >= gamedata.MiningMapLength {
			continue
		}

		//! 去除已挖掘
		if self.MiningMap.Get(index) == true {
			continue
		}

		//! 去除已有事件
		if self.Element.Get(index) != 0 {
			continue
		}

		posLst = append(posLst, TMiningPos{v.X, v.Y})
	}

	return posLst
}

//! 生成可探视区域坐标
func (self *TMiningModule) GetVisualPosArena(x int, y int) (posLst []TMiningPos) {
	for i := 1; i <= 2; i++ {
		if x-i >= 0 && x-i <= gamedata.MiningMapLength {
			posLst = append(posLst, TMiningPos{x - i, y})
		}

		if x+i <= gamedata.MiningMapLength {
			posLst = append(posLst, TMiningPos{x + i, y})
		}

		if y-i >= 0 && y-i <= gamedata.MiningMapLength {
			posLst = append(posLst, TMiningPos{x, y - i})
		}

		if y+i <= gamedata.MiningMapLength {
			posLst = append(posLst, TMiningPos{x, y + i})
		}
	}

	if x-1 >= 0 && y-1 >= 0 {
		posLst = append(posLst, TMiningPos{x - 1, y - 1})
	}

	if x-1 >= 0 && y+1 <= gamedata.MiningMapLength {
		posLst = append(posLst, TMiningPos{x - 1, y + 1})
	}

	if y+1 <= gamedata.MiningMapLength && x+1 <= gamedata.MiningMapLength {
		posLst = append(posLst, TMiningPos{x + 1, y + 1})
	}

	if y-1 >= 0 && x+1 <= gamedata.MiningMapLength {
		posLst = append(posLst, TMiningPos{x + 1, y - 1})
	}

	return posLst
}

//! 生成一张新的挖矿地图
func (self *TMiningModule) CreateNewMap(isSave bool) {
	//! 初始化地图
	self.initMap()

	posLst := self.GetVisualPosArena(gamedata.MiningEnterPointX, gamedata.MiningEnterPointY)

	for _, v := range posLst {
		self.GetMapPosData(v.X, v.Y, false)
	}

	//! 存储地图
	if isSave == true {
		go self.DB_SaveMiningMap()
	}
}

//! 初始化地图 以下函数为该类私有,故不对外提供方法
func (self *TMiningModule) initMap() {
	self.MiningMap = MapLst{}
}

//! 获取当前玩家进度
func (self *TMiningModule) GetSchedule() float64 {
	digCount := 0
	totalCount := gamedata.MiningMapLength * gamedata.MiningMapLength

	digCount = self.MiningMap.Count()
	digCount += len(self.Element)

	schedule := float64(digCount) / float64(totalCount)

	return schedule
}

//! 检测重置时间
func (self *TMiningModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TMiningModule) OnNewDay(newday uint32) {

	//! 重置购买行动力次数
	self.ActionBuyTimes = 0
	self.ResetDay = newday
	go self.DB_SaveActionBuyTimes()
}

//! 当前是否刷新Boss
func (self *TMiningModule) isRefreshBoss() bool {
	bossInfo, _ := self.GetBossInfo()
	if bossInfo != nil {
		return true
	}

	return false
}

//! 随机元素
func (self *TMiningModule) randElement(index int, savedb bool) (element int) {
	//! 随机一个元素
	element = gamedata.RandMiningElement()
	//! 判断完成度
	if self.GetSchedule() >= 0.7 && self.isRefreshBoss() == false {
		element = gamedata.MiningElement_Event
	}

	if element == gamedata.MiningElement_Event {
		//! 随机一个事件
		event := gamedata.RandMimingEvent()
		if self.GetSchedule() >= 0.7 && self.isRefreshBoss() == false {
			event = gamedata.MiningEvent_Boss
		}

		if event == gamedata.MiningEvent_Elite_Monster ||
			event == gamedata.MiningEvent_Normal_Monster ||
			event == gamedata.MiningEvent_Boss {
			//! 若为怪物,则随机怪物信息
			monsterID := gamedata.RandMiningMonster(event, self.ownplayer.GetLevel())

			monster := gamedata.GetMonsterEventInfo(monsterID)

			var info TMiningMonster
			info.ID = monsterID
			info.Index = index

			info.Life = int(float64(monster.MonsterLife) * math.Pow(1.2, float64(self.MiningResetTimes)))

			self.MonsterLst = append(self.MonsterLst, info)

			if event == gamedata.MiningEvent_Boss {
				self.RandBossAward(monsterID)
			}

			if savedb == true {
				go self.DB_AddMonster(info)
				go self.DB_SaveBossAward()
			}

		}

		element = event
		self.Element.Set(index, event)
	} else {
		self.Element.Set(index, element)
	}

	if savedb == true {
		go self.DB_AddElement(index<<16 + element)
	}
	return element
}

//! 地图版本号加一
func (self *TMiningModule) AddMapStatusCode() {
	mapCode := (self.StatusCode >> 16) & 0x0000FFFF
	miningCode := self.StatusCode & 0x0000FFFF
	self.StatusCode = (mapCode << 16) + miningCode

	go self.DB_UpdateMiningStatusCode()
}

//! 矿洞版本号加一
func (self *TMiningModule) AddMiningStatusCode() {
	self.StatusCode += 1
	go self.DB_UpdateMiningStatusCode()
}

//! 获取坐标信息
func (self *TMiningModule) GetMapPosData(x int, y int, savedb bool) (isDig bool, element int, errcode int) {
	index := y*gamedata.MiningMapLength + x
	isDig = self.MiningMap.Get(index)
	if self.Element.Get(index) != 0 {
		element = self.Element.Get(index)
		if element == 0 {
			gamelog.Error("GetMapPosData error: get element fail. Index: %d", index)
			errcode = msg.RE_INVALID_PARAM
		}
	} else {
		//! 未生成状态,随机元素
		element = self.randElement(index, savedb)
	}

	return isDig, element, errcode
}

//! 获取怪物信息
func (self *TMiningModule) GetMonsterInfo(index int) (*TMiningMonster, int) {
	for i, v := range self.MonsterLst {
		if v.Index == index {
			return &self.MonsterLst[i], i
		}
	}

	gamelog.Error("GetMonsterInfo error: index: %d", index)
	return nil, -1
}

//! 获取Boss信息
func (self *TMiningModule) GetBossInfo() (*TMiningMonster, int) {
	for i, v := range self.MonsterLst {
		monsterData := gamedata.GetMonsterEventInfo(v.ID)
		if monsterData.Event == gamedata.MiningEvent_Boss {
			return &self.MonsterLst[i], i
		}
	}

	gamelog.Error("GetBossInfo error")
	return nil, -1
}

//! 增加一个Buff信息
func (self *TMiningModule) AddBuff(buffType int, times int, value int) {

	self.Buff.BuffType = buffType
	self.Buff.Times = times
	self.Buff.Value = value
	go self.DB_SavePlayerBuff()
}

//! 随机Boss奖励
func (self *TMiningModule) RandBossAward(bossID int) {
	bossInfo := gamedata.GetMonsterEventInfo(bossID)
	copyBaseInfo := gamedata.GetCopyBaseInfo(bossInfo.CopyID)

	awardLst := gamedata.GetItemsFromAwardID(copyBaseInfo.AwardID)
	for i, v := range awardLst {
		self.BossAward = append(self.BossAward, TMiningBossAward{i, v.ItemID, v.ItemNum, false})
	}
}

//! 获取Boss奖励
func (self *TMiningModule) GetBossAward(id int) *TMiningBossAward {

	for _, v := range self.BossAward {
		if v.ID == id {
			return &v
		}
	}

	gamelog.Error("GetBossAward Error: invlid id: %d", id)
	return nil
}

//! 设置挂机
func (self *TMiningModule) SetGuaji(guajiType int, hour int) {
	self.GuaJiType = guajiType
	self.GuajiCalcTime = time.Now().Unix() + int64(hour*60*60)
	go self.DB_SaveGuajiInfo()
}

//! 获取挂机奖励
func (self *TMiningModule) GetGuajiAward() []gamedata.ST_ItemData {
	info := gamedata.GetMiningGuajiInfo(self.GuaJiType)
	if info == nil {
		gamelog.Error("GetMiningGuajiInfo fail. guajiType: %d", self.GuaJiType)
		return []gamedata.ST_ItemData{}
	}

	awardItems := gamedata.GetItemsFromAwardID(info.Award)
	self.ownplayer.BagMoudle.AddAwardItems(awardItems)

	//! 重置挂机时间信息
	self.GuaJiType = 0
	self.GuajiCalcTime = 0
	go self.DB_SaveGuajiInfo()

	return awardItems
}

//! 检测行动力自增
func (self *TMiningModule) CheckActionAddTime() int {
	// if self.ownplayer.RoleMoudle.CheckActionEnough(gamedata.MiningCostActionID, gamedata.MiningActionRecoverLimit) == true {
	// 	//! 如果已超过界限,则停止增长
	// 	return 0
	// }

	// duration := time.Now().Unix() - self.ActionAddTime
	// if duration < int64(gamedata.MiningActionRecoverTime) {
	// 	return 0
	// }

	// addAction := duration / int64(gamedata.MiningActionRecoverTime)
	// self.ActionAddTime = time.Now().Unix() - duration%int64(gamedata.MiningActionRecoverTime)

	// action := self.ownplayer.RoleMoudle.GetAction(gamedata.MiningCostActionID)
	// if action+int(addAction) > gamedata.MiningActionRecoverLimit {
	// 	addAction = int64(gamedata.MiningActionRecoverLimit) - int64(action)
	// }

	// self.ownplayer.RoleMoudle.AddAction(gamedata.MiningCostActionID, int(addAction))
	// go self.DB_SaveActionAddTime()

	// //! 返回倒计时
	// nextTime := self.ActionAddTime + int64(gamedata.MiningActionRecoverTime) - time.Now().Unix()
	return 0
}

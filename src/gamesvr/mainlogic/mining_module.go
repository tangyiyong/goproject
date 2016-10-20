package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math"
	"mongodb"
	"msg"
	"sync"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TMiningMonster struct {
	ID    int
	Index int32
	Life  int
}

type TMiningPos struct {
	X int32
	Y int32
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

type ElementLst []int32

func (self *ElementLst) Set(index int32, value int32) {
	for i, v := range *self {
		if (v>>16)&0x0000FFFF == index {
			(*self)[i] = index << 16
			(*self)[i] += value
			return
		}
	}

	var newValue int32
	newValue = index << 16
	newValue += value
	*self = append(*self, newValue)
}

func (self *ElementLst) Get(index int32) int32 {
	for i, v := range *self {
		if v>>16 == index {
			return (*self)[i] & 0x0000FFFF
		}
	}

	return 0
}

//! 删除元素
func (self *TMiningModule) DeleteElement(pos int32) {
	index := 0
	var value int32 = 0
	for i, v := range self.Element {
		if ((v >> 16) & 0x0000FFFF) == pos {
			index = i
			value = v
			break
		}
	}

	if index == 0 {
		self.Element = self.Element[1:]
	} else if (index + 1) == len(self.Element) {
		self.Element = self.Element[:index]
	} else {
		self.Element = append(self.Element[:index], self.Element[index+1:]...)
	}

	self.DB_RemoveElement(value)
}

type TMapData [60]uint64

func (self *TMapData) Set(index int32) {
	y := uint64(index / 60)
	x := uint64(index % 60)
	(*self)[y] |= (1 << x)
}

func (self *TMapData) Get(index int32) bool {
	y := uint64(index / 60)
	x := uint64(index % 60)

	return ((*self)[y] & (1 << x)) > 0
}

type TMiningModule struct {
	PlayerID   int32              `bson:"_id"`
	GuaJiType  int                //! 当前挂机类型
	GuajiTime  int32              //! 挂机结算时间
	Point      int                //! 玩家当前积分
	Buff       TMiningBuff        //! 玩家Buff
	MonsterLst []TMiningMonster   //! 怪物信息
	MapData    TMapData           //! 地图挖掘标记 位运算
	MapCnt     int                //! 己经打开的位置个数
	Element    ElementLst         //! 前16位为index  后16位为值
	LastPos    TMiningPos         //! 最后操作位置
	ResetDay   uint32             //! 重置购买次数天数
	BossAward  []TMiningBossAward //! 翻牌奖励
	BuyRecord  Int32Lst           //! 神秘商店购买标记
	StatusCode int                //! 状态码
	ResetTimes int                //! 地图重置次数
	ownplayer  *TPlayer
}

//! 设置玩家指针
func (self *TMiningModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TMiningModule) OnCreate(playerid int32) {

	//! 创建地图
	self.CreateNewMap()

	//! 初始化状态码
	self.StatusCode = 1 + (1 << 16)

	//! 初始化入口坐标
	self.LastPos.X = gamedata.MiningStartPointX
	self.LastPos.Y = gamedata.MiningStartPointY

	//! 初始化地图入口
	self.MapData.Set(gamedata.MiningStartPointY*60 + gamedata.MiningStartPointX)

	//! 重置矿洞次数
	self.ResetTimes = 0

	//! 设置重置购买次数时间
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	mongodb.InsertToDB("PlayerMining", self)
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

	err := s.DB(appconfig.GameDbName).C("PlayerMining").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerMining Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

func (self *TMiningModule) RedTip() bool {
	now := utility.GetCurTime()
	if self.GuajiTime < now {
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
	self.ResetTimes++

	//! 标记重置
	self.BossAward = []TMiningBossAward{} //! 翻牌奖励
	self.Buff = TMiningBuff{}
	self.Point = 0

	self.Element = ElementLst{}
	self.MonsterLst = []TMiningMonster{}
	self.StatusCode = 1 + (1 << 16)

	self.CreateNewMap()

	//! 初始化状态码
	self.StatusCode = 1 + (1 << 16)

	//! 初始化入口坐标
	self.LastPos.X = gamedata.MiningStartPointX
	self.LastPos.Y = gamedata.MiningStartPointY

	//! 初始化地图入口
	self.MapData.Set(gamedata.MiningStartPointY*60 + gamedata.MiningStartPointX)

	//! 设置重置购买次数时间
	self.ResetDay = utility.GetCurDay()

	self.DB_ResetMapData()
}

//! 删除怪物
func (self *TMiningModule) DeleteMonster(index int32) {
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

	self.DB_RemoveMonster(monster)
}

//! 生成新可视区域部分坐标
func (self *TMiningModule) GetNewVisualPosArena(x int32, y int32) (posLst []TMiningPos) {
	visualLst := self.GetVisualPosArena(x, y)

	//! 去除已经可以看见的区域
	for _, v := range visualLst {
		index := v.Y*gamedata.MiningMapLength + v.X

		if v.X >= gamedata.MiningMapLength || v.Y >= gamedata.MiningMapLength {
			continue
		}

		//! 去除已挖掘
		if self.MapData.Get(index) == true {
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
func (self *TMiningModule) GetVisualPosArena(x int32, y int32) (posLst []TMiningPos) {
	posLst = append(posLst, TMiningPos{x + 1, y})
	posLst = append(posLst, TMiningPos{x + 2, y})
	posLst = append(posLst, TMiningPos{x, y + 1})
	posLst = append(posLst, TMiningPos{x, y + 2})
	posLst = append(posLst, TMiningPos{x + 1, y + 1})

	if y >= 2 {
		posLst = append(posLst, TMiningPos{x, y - 2})
	}

	if y >= 1 {
		posLst = append(posLst, TMiningPos{x + 1, y - 1})
		posLst = append(posLst, TMiningPos{x, y - 1})
	}

	if x >= 2 {
		posLst = append(posLst, TMiningPos{x - 2, y})
	}

	if x >= 1 {
		posLst = append(posLst, TMiningPos{x - 1, y + 1})
		posLst = append(posLst, TMiningPos{x - 1, y})
	}

	if x >= 1 && y >= 1 {
		posLst = append(posLst, TMiningPos{x - 1, y - 1})
	}

	return posLst
}

//! 生成一张新的挖矿地图
func (self *TMiningModule) CreateNewMap() {
	//! 初始化地图
	self.MapData = TMapData{}
	self.Element = ElementLst{}
	self.MapCnt = 1

	posLst := self.GetVisualPosArena(gamedata.MiningStartPointX, gamedata.MiningStartPointY)

	for _, v := range posLst {
		self.GetMapPosData(v.X, v.Y, false)
	}
}

//! 获取当前玩家进度
func (self *TMiningModule) GetSchedule() float64 {
	digCount := 0
	totalCount := gamedata.MiningMapLength * gamedata.MiningMapLength
	digCount = self.MapCnt
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
	self.ResetDay = newday
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
func (self *TMiningModule) randElement(index int32, savedb bool) (element int32) {
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

			info.Life = int(float64(monster.MonsterLife) * math.Pow(1.2, float64(self.ResetTimes)))

			self.MonsterLst = append(self.MonsterLst, info)

			if event == gamedata.MiningEvent_Boss {
				self.RandBossAward(monsterID)
			}

			if savedb == true {
				self.DB_AddMonster(info)
				self.DB_SaveBossAward()
			}

		}

		element = event
		self.Element.Set(index, event)
	} else {
		self.Element.Set(index, element)
	}

	if savedb == true {
		self.DB_AddElement(index<<16 + element)
	}
	return element
}

//! 地图版本号加一
func (self *TMiningModule) AddMapStatusCode() {
	mapCode := (self.StatusCode >> 16) & 0x0000FFFF
	miningCode := self.StatusCode & 0x0000FFFF
	self.StatusCode = (mapCode << 16) + miningCode
	self.DB_UpdateMiningStatusCode()
}

//! 矿洞版本号加一
func (self *TMiningModule) AddMiningStatusCode() {
	self.StatusCode += 1
	self.DB_UpdateMiningStatusCode()
}

//! 获取坐标信息
func (self *TMiningModule) GetMapPosData(x int32, y int32, savedb bool) (isDig bool, element int32, errcode int) {
	index := y*gamedata.MiningMapLength + x
	isDig = self.MapData.Get(index)
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
func (self *TMiningModule) GetMonsterInfo(index int32) (*TMiningMonster, int) {
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
	self.DB_SaveBuff()
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
	self.GuajiTime = utility.GetCurTime() + int32(hour*60*60)
	self.DB_SaveGuajiInfo()
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
	self.GuajiTime = 0
	self.DB_SaveGuajiInfo()

	return awardItems
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"mongodb"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

//! 暴动信息
type TTerritoryRiotData struct {
	IsRoit     bool   //! 是否暴动
	BeginTime  int64  //! 开始时间
	DealTime   int64  //! 处理时间
	HelperName string //! 帮忙处理好友姓名
}

type TTerritoryInfo struct {
	ID            int                    //! 领地ID
	HeroID        int                    //! 巡逻武将ID
	PatrolTime    int                    //! 巡逻时间
	AwardTime     int                    //! 奖励间隔时间
	PatrolEndTime int64                  //! 巡逻结束时间
	AwardItem     []gamedata.ST_ItemData //! 已获奖励
	RiotInfo      []TTerritoryRiotData   //! 暴动信息
	SkillLevel    int                    //! 领地技能等级
}

//! 领地攻伐模块
type TTerritoryModule struct {
	PlayerID int32 `bson:"_id"`

	TerritoryLst      []TTerritoryInfo
	SuppressRiotTimes int    //! 当前镇压暴动次数
	TotalPatrolTime   int    //! 总计巡逻时间(小时)
	ResetDay          uint32 //! 镇压次数刷新时间
	ownplayer         *TPlayer
}

//! 设置指针
func (self *TTerritoryModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//! 玩家创建角色
func (self *TTerritoryModule) OnCreate(playerid int32) {
	//! 初始化信息
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerTerritory", self)
}

//! 玩家销毁角色
func (self *TTerritoryModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (self *TTerritoryModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离线
func (self *TTerritoryModule) OnPlayerOffline(playerid int32) {

}

//! 预取玩家信息
func (self *TTerritoryModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerTerritory").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerTerritory Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
}

//! 红点
func (self *TTerritoryModule) RedTip() bool {
	now := time.Now().Unix()
	for _, v := range self.TerritoryLst {
		if v.HeroID != 0 &&
			v.PatrolEndTime < now {
			return true
		}
	}

	return false
}

//! 镇压次数重置
func (self *TTerritoryModule) CheckReset() {
	//! 对比重置时间
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TTerritoryModule) OnNewDay(newday uint32) {
	self.ResetDay = newday
	//! 开始重置
	self.SuppressRiotTimes = 0
	self.ResetDay = utility.GetCurDay()
	self.DB_UpdateResetTime()
}

//! 检测领地是否已被挑战
func (self *TTerritoryModule) IsChallenged(id int) bool {
	for _, v := range self.TerritoryLst {
		if v.ID == id {
			return true
		}
	}
	return false
}

//! 挑战领地结果
func (self *TTerritoryModule) ChallengeTerritory(id int) {
	//! 获取领地信息
	territoryInfo := gamedata.GetTerritoryData(id)

	//! 添加领地信息
	var territory TTerritoryInfo
	territory.ID = territoryInfo.ID
	self.TerritoryLst = append(self.TerritoryLst, territory)
	self.DB_AddTerritory(territory)

	//! 发放奖励
	copyInfo := gamedata.GetCopyBaseInfo(territoryInfo.CopyID)
	itemLst := gamedata.GetItemsFromAwardID(copyInfo.AwardID)
	self.ownplayer.BagMoudle.AddAwardItems(itemLst)
}

//! 获取领地信息
func (self *TTerritoryModule) GetTerritory(id int) (*TTerritoryInfo, int) {
	for i, v := range self.TerritoryLst {
		if v.ID == id {
			return &self.TerritoryLst[i], i
		}
	}
	return nil, 0
}

//! 巡逻领地
func (self *TTerritoryModule) PatrolTerritory(id int, heroID int, patrol *gamedata.ST_TerritoryPatrolType, awardTime int) {
	//! 获取领地信息
	now := time.Now().Unix()
	territory, index := self.GetTerritory(id)
	if territory == nil {
		gamelog.Error("GetTerritory fail. ID: %d", id)
		return
	}
	territory.HeroID = heroID
	territory.PatrolTime = patrol.Time
	territory.PatrolEndTime = now + int64(patrol.Time)
	territory.AwardTime = awardTime

	//! 随机武将碎片奖励
	pHeroInfo := gamedata.GetHeroInfo(heroID)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pieceNum := 0
	if patrol.Type == 1 {
		pieceNum = 1
	} else if patrol.Type == 2 {
		pieceNum = 1 + r.Intn(2)
	} else if patrol.Type == 3 {
		pieceNum = 1 + r.Intn(3)
	}

	var riotEndTime int64

	//! 随机奖励列表
	timeInterval := patrol.Time / awardTime
	for i := 1; i <= timeInterval+1; i++ {
		award := gamedata.RandTerritoryAward(id)

		//! 获取领地技能信息
		if territory.SkillLevel > 0 {
			skill := gamedata.GetTerritorySkillData(id, territory.SkillLevel)
			randValue := r.Intn(100) + 1
			if skill.DoublePro < randValue {
				award.ItemNum *= 2
			}
		}

		territory.AwardItem = append(territory.AwardItem, gamedata.ST_ItemData{award.ItemID, award.ItemNum})

		//! 随机暴动
		randRiot := r.Intn(10000)
		if randRiot < gamedata.RiotPro && int64(riotEndTime) < now+int64(i*awardTime) {
			//! 发生暴动
			var riot TTerritoryRiotData
			riot.BeginTime = now + int64(i*awardTime)
			riot.IsRoit = true
			territory.RiotInfo = append(territory.RiotInfo, riot)
			self.DB_DB_AddTerritoryRiotInfo(index, riot)

			riotEndTime = now + int64(i*awardTime) + int64(gamedata.RiotTime)
		}

	}

	territory.AwardItem = append(territory.AwardItem, gamedata.ST_ItemData{pHeroInfo.PieceID, pieceNum})
	//! 更新到数据库
	self.DB_UpdateTerritory(index, territory)
}

//! 领地是否暴动
func (self *TTerritoryModule) IsRiot(id int) bool {
	//! 获取领地信息
	territory, _ := self.GetTerritory(id)
	isRiot := false
	for _, n := range territory.RiotInfo {
		//! 判断暴动
		if time.Now().Unix() >= n.BeginTime &&
			time.Now().Unix() < n.BeginTime+int64(gamedata.RiotTime) &&
			n.IsRoit == true && time.Now().Unix() < territory.PatrolEndTime {
			isRiot = true
		}
	}
	return isRiot
}

//! 领取领地巡逻奖励
func (self *TTerritoryModule) GetTerritoryAward(id int) {
	//! 获取领地信息
	territory, index := self.GetTerritory(id)
	for _, v := range territory.AwardItem {
		if v.ItemID != 0 {
			self.ownplayer.BagMoudle.AddAwardItem(v.ItemID, v.ItemNum)
		}
	}

	//! 根据巡逻类型加上累积时间
	self.TotalPatrolTime += territory.PatrolTime

	//! 清空相关数据
	territory.HeroID = 0
	territory.AwardItem = []gamedata.ST_ItemData{}
	territory.PatrolEndTime = 0
	territory.PatrolTime = 0
	territory.AwardTime = 0
	territory.RiotInfo = []TTerritoryRiotData{}
	self.DB_UpdateTerritory(index, territory)
}

//! 升级领地技能
func (self *TTerritoryModule) TerritorySkillLevelUp(id int) {
	territory, index := self.GetTerritory(id)
	territory.SkillLevel += 1
	self.DB_UpdateTerritorySkill(index, territory.SkillLevel)
}

//! 获取领地数量
func (self *TTerritoryModule) GetTerritoryNum() int {
	return len(self.TerritoryLst)
}

//! 获取暴动领地数量
func (self *TTerritoryModule) GetRiotTerritoryNum() []int {
	lst := []int{}
	for _, v := range self.TerritoryLst {
		if true == self.IsRiot(v.ID) {
			lst = append(lst, v.ID)
		}
	}
	return lst
}

//! 获取巡逻领地数量
func (self *TTerritoryModule) GetPatrolNum() []int {
	lst := []int{}
	for _, v := range self.TerritoryLst {
		if time.Now().Unix() < v.PatrolEndTime {
			lst = append(lst, v.ID)
		}
	}
	return lst
}

//! 获取好友领地模块指针
func (self *TTerritoryModule) GetFriendTerritory(playerid int32) *TTerritoryModule {
	var friendTerritory *TTerritoryModule
	friendPlayer := GetPlayerByID(playerid)
	if friendPlayer == nil {
		//! 尝试从数据库读取数据
		friendPlayer = LoadPlayerFromDB(playerid)
	}
	friendTerritory = &friendPlayer.TerritoryModule
	return friendTerritory
}

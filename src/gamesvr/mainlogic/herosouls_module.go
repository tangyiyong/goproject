package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type THeroSouls struct {
	ID      int  //! 唯一ID
	HeroID  int  //! 将灵ID
	IsExist bool //! true 存在 false 不存在
}

type THeroSoulsLink struct {
	ID    int //! 将灵链接ID
	Level int //! 将灵等级
}

type THeroSoulsStore struct {
	ItemID   int  //! 商品ID
	IsBuy    bool //! 是否已经购买
	MoneyID  int  //! 货币ID
	MoneyNum int  //! 货币数量
}

type THeroSoulsInfo struct {
	PlayerID     int64 //! 玩家ID
	SoulMapValue int   //! 阵图值
	SoulMapCount int   //! 阵图个数
}

type THeroSoulsProperty struct {
	PropertyID    int
	PropertyValue int
	Camp          int
	IsPercent     bool
}

//! 将灵模块
type THeroSoulsModule struct {
	PlayerID int32 `bson:"_id"`

	TargetIndex          int    //! 指针指向
	UnLockChapter        int    //! 当前解锁章节
	SoulMapValue         int    //! 阵图值
	ChallengeTimes       int    //! 当前剩余挑战将灵次数
	BuyChallengeTimes    int    //! 当前已购买挑战将灵次数
	ResetDay             uint32 //! 重置天数
	RefreshStoreTimeMark Mark   //! 更新商店时间标记

	Achievement int //! 阵图成就

	HeroSoulsStoreLst []THeroSoulsStore //! 将灵商店
	HeroSoulsLst      []THeroSouls      //! 可挑战将灵
	HeroSoulsLink     []THeroSoulsLink  //! 已激活将灵链接

	//临时数据
	propertyInt            [11]int //! 加成实际数值
	propertyPercent        [11]int //! 加成百分比
	campPropertyKillLst    [4]int  //! 对阵营加伤百分比
	campPropertyDefenceLst [4]int  //! 对阵营减伤百分比

	ownplayer *TPlayer
}

func (self *THeroSoulsModule) CalcHeroSoulProperty() {
	self.CalcSoulLinkProperty()
	self.CalcAchievementProperty()
}

//! 增加临时属性
func (self *THeroSoulsModule) AddTempProperty(pid int, pvalue int, percent bool, camp int) {
	if pid == gamedata.AttackPropertyID {
		if percent == true {
			self.propertyPercent[gamedata.AttackPhysicID-1] += pvalue
			self.propertyPercent[gamedata.AttackMagicID-1] += pvalue
		} else {
			self.propertyInt[gamedata.AttackPhysicID-1] += pvalue
			self.propertyInt[gamedata.AttackMagicID-1] += pvalue
		}
	} else if pid == gamedata.DefencePropertyID {
		if percent == true {
			self.propertyPercent[gamedata.DefencePhysicID-1] += pvalue
			self.propertyPercent[gamedata.DefenceMagicID-1] += pvalue
		} else {
			self.propertyInt[gamedata.DefencePhysicID-1] += pvalue
			self.propertyInt[gamedata.DefenceMagicID-1] += pvalue
		}
	} else {
		if camp <= 0 {
			if percent == true {
				self.propertyPercent[pid-1] += pvalue
			} else {
				self.propertyInt[pid-1] += pvalue
			}
		} else {
			if pid == 6 {
				self.campPropertyDefenceLst[pid-1] += pvalue
			} else if pid == 7 {
				self.campPropertyKillLst[pid-1] += pvalue
			} else {
				gamelog.Error("AddExtraProperty Error : if camp:%d != 0, pid :%d should be 6 or 7!")
			}
		}
	}

}

//! 计算英灵阵图属性加成
func (self *THeroSoulsModule) CalcSoulLinkProperty() {
	for _, v := range self.HeroSoulsLink {
		info := gamedata.GetHeroSoulsInfo(v.ID)
		if info == nil {
			gamelog.Error("CalcSoulLinkProperty Error : Invalid herosoul id:%d", v.ID)
			return
		}

		for _, n := range info.Property {
			if n.PropertyID == 0 {
				continue
			}

			if n.PropertyID == gamedata.AttackPropertyID { //! 攻击->物理+魔法
				if n.Is_Percent == true {
					self.propertyPercent[gamedata.AttackPhysicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
					self.propertyPercent[gamedata.AttackMagicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
				} else {
					self.propertyInt[gamedata.AttackPhysicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
					self.propertyInt[gamedata.AttackMagicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
				}

			} else if n.PropertyID == gamedata.DefencePropertyID { //! 防御->物理+魔法
				if n.Is_Percent == true {
					self.propertyPercent[gamedata.DefencePhysicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
					self.propertyPercent[gamedata.DefenceMagicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
				} else {
					self.propertyInt[gamedata.DefencePhysicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
					self.propertyInt[gamedata.DefenceMagicID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
				}

			} else {
				if n.PropertyID == 6 && n.Camp != 0 {
					self.campPropertyDefenceLst[n.Camp-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
				} else if n.PropertyID == 7 && n.Camp != 0 {
					self.campPropertyKillLst[n.Camp-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
				} else {
					if n.Is_Percent == true {
						self.propertyPercent[n.PropertyID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
					} else {
						self.propertyInt[n.PropertyID-1] += n.PropertyValue + n.LevelUp*(v.Level-1)
					}
				}
			}

		}
	}
	return
}

//! 计算阵图成就属性加成
func (self *THeroSoulsModule) CalcAchievementProperty() (propertyLst []THeroSoulsProperty) {
	for i := 1; i <= self.Achievement; i++ {
		info := gamedata.GetSoulMapInfo(i)
		if info == nil {
			break
		}

		isExist := false
		for j, v := range propertyLst {
			if v.PropertyID == info.PropertyID {
				propertyLst[j].PropertyValue += info.PropertyValue
				isExist = true
				break
			}
		}

		if isExist == true {
			continue
		} else {
			var property THeroSoulsProperty
			property.PropertyID = info.PropertyID
			property.PropertyValue = info.PropertyValue
			property.IsPercent = true
			propertyLst = append(propertyLst, property)
		}
	}

	return propertyLst
}

func (self *THeroSoulsModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

func (self *THeroSoulsModule) OnCreate(playerid int32) {
	//! 初始化各类参数
	self.ResetHeroSoulsLst(false)
	self.RefreshHeroSoulsStore(false)
	self.UnLockChapter = 1
	self.TargetIndex = 0
	self.BuyChallengeTimes = 0
	self.ChallengeTimes = gamedata.HeroSoulsChallengeTimes
	self.ResetDay = utility.GetCurDay()

	//! 插入数据库
	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerHeroSouls", self)
}

func (self *THeroSoulsModule) OnDestroy(playerid int32) {

}

func (self *THeroSoulsModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (self *THeroSoulsModule) OnPlayerOffline(playerid int32) {

}

//! 读取玩家
func (self *THeroSoulsModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerHeroSouls").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerHeroSouls Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}

	self.PlayerID = playerid
	self.CalcHeroSoulProperty()
}

func (self *THeroSoulsModule) CheckStoreRefresh() int {
	now := time.Now()
	sec := now.Hour()*3600 + now.Minute()*60 + now.Second()
	index := len(gamedata.HeroSoulsStoreRefreshTime)
	for i, v := range gamedata.HeroSoulsStoreRefreshTime {
		if sec < v {
			index = i
			break
		}
	}

	for i := 0; i < index; i++ {
		if self.RefreshStoreTimeMark.Get(uint32(i+1)) == false {
			self.RefreshHeroSoulsStore(true)
			self.RefreshStoreTimeMark.Set(uint32(i + 1))
			self.DB_SaveHeroSoulsRefreshMark()
		}
	}

	countDown := 0
	if index != len(gamedata.HeroSoulsStoreRefreshTime) {
		countDown = gamedata.HeroSoulsStoreRefreshTime[index] - sec
	} else {
		countDown = (24*60*60 - sec) + gamedata.HeroSoulsStoreRefreshTime[0]
	}

	return countDown
}

func (self *THeroSoulsModule) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *THeroSoulsModule) OnNewDay(newday uint32) {
	self.BuyChallengeTimes = 0
	self.RefreshStoreTimeMark = 0
	self.ChallengeTimes = gamedata.HeroSoulsChallengeTimes
	self.ResetDay = newday
	self.DB_Reset()
}

func (self *THeroSoulsModule) RedTip() bool {
	if self.ChallengeTimes != 0 { //! 挑战英灵次数
		return true
	}

	//! 商店刷新
	now := time.Now()
	sec := now.Hour()*3600 + now.Minute()*60 + now.Second()
	index := len(gamedata.HeroSoulsStoreRefreshTime)
	for i, v := range gamedata.HeroSoulsStoreRefreshTime {
		if sec < v {
			index = i
			break
		}
	}

	for i := 0; i < index; i++ {
		if self.RefreshStoreTimeMark.Get(uint32(i+1)) == false {
			self.CheckStoreRefresh()
			return true
		}
	}

	return false
}

//! 刷新神将
func (self *THeroSoulsModule) RefreshHeroSoulsStore(isSave bool) {
	if len(self.HeroSoulsStoreLst) > 0 {
		self.HeroSoulsStoreLst = []THeroSoulsStore{}
	}

	//! 加入两个固定商品
	self.HeroSoulsStoreLst = append(self.HeroSoulsStoreLst, THeroSoulsStore{gamedata.HeroSoulsStoreFixedItemID, false,
		gamedata.HeroSoulsStoreFixedItemMoneyID, gamedata.HeroSoulsStoreFixedItemMoneyNum})

	self.HeroSoulsStoreLst = append(self.HeroSoulsStoreLst, THeroSoulsStore{gamedata.HeroSoulsStoreFixedItemID2, false,
		gamedata.HeroSoulsStoreFixedItemMoneyID2, gamedata.HeroSoulsStoreFixedItemMoneyNum2})

	heroIDLst := gamedata.RandHeroSoulsStore(4)

	for _, v := range heroIDLst {
		var goods THeroSoulsStore
		goods.ItemID = v
		goods.IsBuy = false

		heroInfo := gamedata.GetHeroSoulsStoreInfo(v)
		goods.MoneyID = heroInfo.MoneyID
		goods.MoneyNum = heroInfo.MoneyNum

		self.HeroSoulsStoreLst = append(self.HeroSoulsStoreLst, goods)
	}

	if isSave == true {
		self.DB_SaveHeroSoulsStoreLst()
	}
}

//! 重置将灵列表
func (self *THeroSoulsModule) ResetHeroSoulsLst(isSave bool) {
	if len(self.HeroSoulsLst) > 0 {
		self.HeroSoulsLst = []THeroSouls{}
	}

	heroIDLst, IDLst := gamedata.RandHeroSouls()

	for i, v := range heroIDLst {
		var herosoul THeroSouls
		herosoul.ID = IDLst[i]
		herosoul.HeroID = v
		herosoul.IsExist = true
		self.HeroSoulsLst = append(self.HeroSoulsLst, herosoul)
	}

	//! 指针指向0坐标
	self.TargetIndex = 0

	if isSave == true {
		self.DB_SaveHeroSoulsLst()
	}
}

//! 增加将灵
func (self *THeroSoulsModule) AddHeroSouls(index int) {
	//! 获取将灵信息
	heroSoulsInfo := self.HeroSoulsLst[index]
	if heroSoulsInfo.IsExist == false {
		gamelog.Error("AddHeroSouls Error: Hero souls info not exist")
		return
	}

	//! 加入将灵背包
	self.ownplayer.BagMoudle.AddHeroSoul(heroSoulsInfo.HeroID, 1)
	heroSoulsInfo.IsExist = false
	self.DB_UpdateHeroSoulsMark(index)
}

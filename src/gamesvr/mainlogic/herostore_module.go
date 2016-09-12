package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"mongodb"
	"msg"
	"sync"
	"time"
	"utility"

	"gopkg.in/mgo.v2/bson"
)

type TStoreItem struct {
	ID       int //! 唯一标识
	ItemID   int //! 商品ID
	ItemNum  int //! 商品数量
	MoneyID  int //! 需求货币种类: 1->元宝 2->将魂
	MoneyNum int //! 价格
	Status   int //! 0->未购买  1->已购买
}

//! 商店模块
type TStoreModule struct {
	PlayerID int32 `bson:"_id"`

	RefreshCount     int   //! 今天剩余可刷新次数
	FreeRefreshCount int   //! 免费刷新次数
	FreeRefreshTime  int64 //! 免费刷新时间

	AwakeRefreshCount     int //! 觉醒商店可刷新次数
	AwakeFreeRefreshCount int //! 觉醒商店
	AwakeFreeRefreshTime  int64

	PetRefreshCount     int //! 战宠商店可刷新次数
	PetFreeRefreshCount int //! 战宠商店
	PetFreeRefreshTime  int64

	ResetDay uint32

	ShopItemLst      []TStoreItem //! 当前商品
	AwakeShopItemLst []TStoreItem //! 当前商品
	PetShopItemLst   []TStoreItem //! 当前商品

	ownplayer *TPlayer //! 父类指针
}

func (storemodule *TStoreModule) SetPlayerPtr(playerid int32, player *TPlayer) {
	storemodule.PlayerID = playerid
	storemodule.ownplayer = player
}

//! 玩家创建角色
func (storemodule *TStoreModule) OnCreate(playerid int32) {
	//! 创建伊始,给予玩家满次免费刷新
	storemodule.FreeRefreshCount = gamedata.StoreFreeRefreshTimes
	storemodule.RefreshCount = storemodule.GetPlayerRefreshCounts()

	storemodule.AwakeFreeRefreshCount = gamedata.StoreFreeRefreshTimes
	storemodule.AwakeRefreshCount = storemodule.GetPlayerRefreshCounts()

	storemodule.PetFreeRefreshCount = gamedata.StoreFreeRefreshTimes
	storemodule.PetRefreshCount = storemodule.GetPlayerRefreshCounts()
	storemodule.ResetDay = utility.GetCurDay()

	//! 创建商品
	storemodule.RefreshGoods(gamedata.StoreType_Hero)
	storemodule.RefreshGoods(gamedata.StoreType_Awake)
	storemodule.RefreshGoods(gamedata.StoreType_Pet)

	//! 插入数据
	mongodb.InsertToDB( "HeroStore", storemodule)
}

//! 获取用户每日可刷新次数上限
func (storemodule *TStoreModule) GetPlayerRefreshCounts() int {
	vipLevel := storemodule.ownplayer.GetVipLevel()
	refreshTimes := gamedata.GetFuncVipValue(gamedata.FUNC_HERO_STORE_RESET, vipLevel)

	return refreshTimes
}

//! 玩家销毁角色
func (storemodule *TStoreModule) OnDestroy(playerid int32) {

}

//! 玩家进入游戏
func (storemodule *TStoreModule) OnPlayerOnline(playerid int32) {

}

//! 玩家离开游戏
func (storemodule *TStoreModule) OnPlayerOffline(playerid int32) {

}

//! 预读取玩家
func (storemodule *TStoreModule) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("HeroStore").Find(&bson.M{"_id": playerid}).One(storemodule)
	if err != nil {
		gamelog.Error("HeroStore Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}
	if wg != nil {
		wg.Done()
	}
	storemodule.PlayerID = playerid
}

//! 刷新商品货物列表
func (storemodule *TStoreModule) RefreshGoods(needType int) {
	if needType == gamedata.StoreType_Hero {
		if len(storemodule.ShopItemLst) > 0 {
			//! 清空原有商品信息
			storemodule.ShopItemLst = []TStoreItem{}
		}

		//! 随机货物
		goodsLst := gamedata.RandomStoreItem(6, storemodule.ownplayer.GetLevel(), gamedata.StoreType_Hero)
		for _, v := range goodsLst {
			var good TStoreItem
			good.ID = v.ID
			good.ItemID = v.ItemID
			good.ItemNum = v.ItemNum
			good.MoneyID = v.MoneyID
			good.MoneyNum = v.MoneyNum
			good.Status = 0
			storemodule.ShopItemLst = append(storemodule.ShopItemLst, good)
		}
	} else if needType == gamedata.StoreType_Awake {
		if len(storemodule.AwakeShopItemLst) > 0 {
			//! 清空原有商品信息
			storemodule.AwakeShopItemLst = []TStoreItem{}
		}

		//! 随机货物
		goodsLst := gamedata.RandomStoreItem(6, storemodule.ownplayer.GetLevel(), gamedata.StoreType_Awake)
		for _, v := range goodsLst {
			var good TStoreItem
			good.ID = v.ID
			good.ItemID = v.ItemID
			good.ItemNum = v.ItemNum
			good.MoneyID = v.MoneyID
			good.MoneyNum = v.MoneyNum
			good.Status = 0
			storemodule.AwakeShopItemLst = append(storemodule.AwakeShopItemLst, good)
		}

		//! 增加任务进度
		storemodule.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_AWAKE_STORE_REFRESH, 1)
	} else if needType == gamedata.StoreType_Pet {
		if len(storemodule.PetShopItemLst) > 0 {
			//! 清空原有商品信息
			storemodule.PetShopItemLst = []TStoreItem{}
		}

		//! 随机货物
		goodsLst := gamedata.RandomStoreItem(6, storemodule.ownplayer.GetLevel(), gamedata.StoreType_Pet)
		for _, v := range goodsLst {
			var good TStoreItem
			good.ID = v.ID
			good.ItemID = v.ItemID
			good.ItemNum = v.ItemNum
			good.MoneyID = v.MoneyID
			good.MoneyNum = v.MoneyNum
			good.Status = 0
			storemodule.PetShopItemLst = append(storemodule.PetShopItemLst, good)
		}
	}

	//! 刷新完毕
}

//! 检测满足刷新条件
func (storemodule *TStoreModule) CheckRefreshDemand(storeType int) (bool, int) {
	isFree := true

	//! 检测免费刷新次数
	if storeType == gamedata.StoreType_Hero {
		if storemodule.FreeRefreshCount <= 0 {
			isFree = false
		}

		if storemodule.RefreshCount <= 0 {
			gamelog.Error("Hand_RefreshHeroStore error: Not enough refresh times.")
			return false, msg.RE_NOT_HAVE_REFRESH_TIMES
		}

	} else if storeType == gamedata.StoreType_Awake {
		if storemodule.AwakeFreeRefreshCount <= 0 {
			isFree = false
		}

		if storemodule.AwakeRefreshCount <= 0 {
			gamelog.Error("Hand_RefreshHeroStore error: Not enough refresh times.")
			return false, msg.RE_NOT_HAVE_REFRESH_TIMES
		}

	} else if storeType == gamedata.StoreType_Pet {
		if storemodule.PetFreeRefreshCount <= 0 {
			isFree = false
		}

		if storemodule.PetRefreshCount <= 0 {
			gamelog.Error("Hand_RefreshHeroStore error: Not enough refresh times.")
			return false, msg.RE_NOT_HAVE_REFRESH_TIMES
		}
	}

	if isFree == false {
		//! 检测用户是否拥有刷新物品
		bEnough := storemodule.ownplayer.BagMoudle.IsItemEnough(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)
		if !bEnough {
			//! 物品不足,检测用户货币是否足够
			if false == storemodule.ownplayer.RoleMoudle.CheckMoneyEnough(gamedata.HeroStoreRefreshNeedMoneyType, gamedata.HeroStoreRefreshNeedMoneyNum) {
				gamelog.Error("Hand_RefreshHeroStore error: Not enough refresh item.")
				return false, msg.RE_NOT_ENOUGH_REFRESH_ITEM
			}
		}

	}

	return true, msg.RE_SUCCESS
}

//! 扣除刷新条件
func (storemodule *TStoreModule) PaymentTerms(storeType int) (int, int) {

	//! 如果存在免费刷新次数
	if storeType == gamedata.StoreType_Hero {
		if storemodule.FreeRefreshCount > 0 {

			//! 免费刷新次数减一
			storemodule.FreeRefreshCount -= 1
			storemodule.RefreshCount -= 1
			storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Hero)
			if storemodule.FreeRefreshTime == 0 {
				//! 设置刷新次数CD时间
				storemodule.FreeRefreshTime = time.Now().Unix() + int64(gamedata.StoreFreeRefreshAddTime)
			}

			storemodule.DB_UpdateRefreshFreeCount(gamedata.StoreType_Hero)

			return 1, 1
		}

		//! 如果不存在免费刷新次数
		//! 扣除道具优先
		bEnough := storemodule.ownplayer.BagMoudle.IsItemEnough(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)
		if bEnough {
			storemodule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)

			//! 当天可刷新总次数减一
			storemodule.RefreshCount -= 1
			storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Hero)
			return 2, 1
		}

		//! 扣除货币其后
		storemodule.ownplayer.RoleMoudle.CostMoney(gamedata.HeroStoreRefreshNeedMoneyType, gamedata.HeroStoreRefreshNeedMoneyNum)

		//! 当天可刷新总次数减一
		storemodule.RefreshCount -= 1
		storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Hero)
	} else if storeType == gamedata.StoreType_Awake {
		if storemodule.AwakeFreeRefreshCount > 0 {

			//! 免费刷新次数减一
			storemodule.AwakeFreeRefreshCount -= 1
			storemodule.AwakeRefreshCount -= 1
			storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Awake)
			if storemodule.AwakeFreeRefreshTime == 0 {
				//! 设置刷新次数CD时间
				storemodule.AwakeFreeRefreshTime = time.Now().Unix() + int64(gamedata.StoreFreeRefreshAddTime)
			}

			storemodule.DB_UpdateRefreshFreeCount(gamedata.StoreType_Awake)

			return 1, 1
		}

		//! 如果不存在免费刷新次数
		//! 扣除道具优先
		bEnough := storemodule.ownplayer.BagMoudle.IsItemEnough(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)
		if bEnough {
			storemodule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)

			//! 当天可刷新总次数减一
			storemodule.AwakeRefreshCount -= 1
			storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Awake)
			return 2, 1
		}

		//! 扣除货币其后
		storemodule.ownplayer.RoleMoudle.CostMoney(gamedata.AwakeStoreRefreshNeedMoneyType, gamedata.AwakeStoreRefreshNeedMoneyNum)

		//! 当天可刷新总次数减一
		storemodule.AwakeRefreshCount -= 1
		storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Awake)

	} else if storeType == gamedata.StoreType_Pet {
		if storemodule.PetFreeRefreshCount > 0 {

			//! 免费刷新次数减一
			storemodule.PetFreeRefreshCount -= 1
			storemodule.PetRefreshCount -= 1
			storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Pet)
			if storemodule.PetFreeRefreshTime == 0 {
				//! 设置刷新次数CD时间
				storemodule.PetFreeRefreshTime = time.Now().Unix() + int64(gamedata.StoreFreeRefreshAddTime)
			}

			storemodule.DB_UpdateRefreshFreeCount(gamedata.StoreType_Pet)

			return 1, 1
		}

		//! 如果不存在免费刷新次数
		//! 扣除道具优先
		bEnough := storemodule.ownplayer.BagMoudle.IsItemEnough(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)
		if bEnough {
			storemodule.ownplayer.BagMoudle.RemoveNormalItem(gamedata.StoreRefreshNeedItem, gamedata.StoreRefreshItemNum)

			//! 当天可刷新总次数减一
			storemodule.PetRefreshCount -= 1
			storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Pet)
			return 2, 1
		}

		//! 扣除货币其后
		storemodule.ownplayer.RoleMoudle.CostMoney(gamedata.PetStoreRefreshNeedMoneyType, gamedata.PetStoreRefreshNeedMoneyNum)

		//! 当天可刷新总次数减一
		storemodule.PetRefreshCount -= 1
		storemodule.DB_UpdateRefreshCount(gamedata.StoreType_Pet)
	}

	return 3, gamedata.HeroStoreRefreshNeedMoneyNum
}

func (self *TStoreModule) RedTip() (bool, IntLst) {
	isRed := false
	redLst := IntLst{}
	if self.FreeRefreshCount != 0 {
		redLst.Add(gamedata.StoreType_Hero)
		isRed = true
	}

	if self.AwakeFreeRefreshCount != 0 {
		redLst.Add(gamedata.StoreType_Awake)
		isRed = true
	}

	if self.PetFreeRefreshCount != 0 {
		redLst.Add(gamedata.StoreType_Pet)
		isRed = true
	}

	return isRed, redLst
}

//! 检测商品状态
func (storemodule *TStoreModule) CheckGoodsStatus(index int, storeType int) (bool, int) {
	if index > 5 {
		return false, msg.RE_INVALID_PARAM
	}

	//! 获取商品
	var item *TStoreItem
	if storeType == gamedata.StoreType_Hero {
		item = &storemodule.ShopItemLst[index]
	} else if storeType == gamedata.StoreType_Awake {
		item = &storemodule.AwakeShopItemLst[index]
	} else if storeType == gamedata.StoreType_Pet {
		item = &storemodule.PetShopItemLst[index]
	}

	//! 检查商品状态
	if item.Status == 1 {
		gamelog.Error("CheckGoodsStatus error: Item is sold out. index: %d", index)
		return false, msg.RE_ITEM_IS_SOLD_OUT
	}

	//! 检查玩家金钱是否足够
	isEnough := storemodule.ownplayer.RoleMoudle.CheckMoneyEnough(item.MoneyID, item.MoneyNum)
	if isEnough == false {
		gamelog.Error("CheckGoodsStatus error: Not enough money. index: %d", index)
		return false, msg.RE_NOT_ENOUGH_MONEY
	}

	return true, msg.RE_SUCCESS
}

//! 支付货币购买商品
func (storemodule *TStoreModule) PayGoods(id int, storeType int) {
	item := storemodule.ShopItemLst[id]

	//! 修改标记
	if storeType == gamedata.StoreType_Hero {
		storemodule.ShopItemLst[id].Status = 1
		storemodule.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_HERO_STORE_BUY, 1)
	} else if storeType == gamedata.StoreType_Awake {
		storemodule.AwakeShopItemLst[id].Status = 1
		item = storemodule.AwakeShopItemLst[id]
		storemodule.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_AWAKE_STORE_REFRESH, 1)
	} else if storeType == gamedata.StoreType_Pet {
		storemodule.PetShopItemLst[id].Status = 1
		item = storemodule.PetShopItemLst[id]
	}

	//! 扣除金钱
	storemodule.ownplayer.RoleMoudle.CostMoney(item.MoneyID, item.MoneyNum)

	//! 发放奖励
	storemodule.ownplayer.BagMoudle.AddAwardItem(item.ItemID, item.ItemNum)

	storemodule.DB_UpdateShopItemStatusToDatabase(id, storeType)

}

//! 刷新免费次数
func (storemodule *TStoreModule) CheckReset(now int64) {
	if utility.IsSameDay(storemodule.ResetDay) == false {
		storemodule.RefreshCount = storemodule.GetPlayerRefreshCounts()
		storemodule.AwakeRefreshCount = storemodule.RefreshCount
		storemodule.PetRefreshCount = storemodule.RefreshCount

		storemodule.ResetDay = utility.GetCurDay()
		storemodule.DB_UpdateResetTime()
	}

	//! 免费次数不得大于最大免费次数限制
	if storemodule.FreeRefreshCount >= gamedata.StoreFreeRefreshTimes {
		storemodule.FreeRefreshTime = 0
	}

	if storemodule.AwakeFreeRefreshCount >= gamedata.StoreFreeRefreshTimes {
		storemodule.AwakeFreeRefreshTime = 0
	}

	if storemodule.PetFreeRefreshCount >= gamedata.StoreFreeRefreshTimes {
		storemodule.PetFreeRefreshTime = 0
	}

	for {
		if now >= storemodule.FreeRefreshTime && storemodule.FreeRefreshTime != 0 {
			//! 免费次数加一
			storemodule.FreeRefreshCount += 1

			if storemodule.FreeRefreshCount >= gamedata.StoreFreeRefreshTimes {
				//! 如果当前已经是免费次数上限,则刷新时间清零
				storemodule.FreeRefreshTime = 0
				break
			} else {
				storemodule.FreeRefreshTime = storemodule.FreeRefreshTime + int64(gamedata.StoreFreeRefreshAddTime)
			}

		} else {
			break
		}
	}

	for {
		if now >= storemodule.AwakeFreeRefreshTime && storemodule.AwakeFreeRefreshTime != 0 {
			//! 免费次数加一
			storemodule.AwakeFreeRefreshCount += 1

			if storemodule.AwakeFreeRefreshCount >= gamedata.StoreFreeRefreshTimes {
				//! 如果当前已经是免费次数上限,则刷新时间清零
				storemodule.AwakeFreeRefreshTime = 0
				break
			} else {
				storemodule.AwakeFreeRefreshTime = storemodule.AwakeFreeRefreshTime + int64(gamedata.StoreFreeRefreshAddTime)
			}

		} else {
			break
		}
	}

	for {
		if now >= storemodule.PetFreeRefreshTime && storemodule.PetFreeRefreshTime != 0 {
			//! 免费次数加一
			storemodule.PetFreeRefreshCount += 1

			if storemodule.PetFreeRefreshCount >= gamedata.StoreFreeRefreshTimes {
				//! 如果当前已经是免费次数上限,则刷新时间清零
				storemodule.PetFreeRefreshTime = 0
				break
			} else {
				storemodule.PetFreeRefreshTime = storemodule.PetFreeRefreshTime + int64(gamedata.StoreFreeRefreshAddTime)
			}

		} else {
			break
		}
	}

	//! 存储到数据库
	storemodule.DB_UpdateRefreshFreeCount(gamedata.StoreType_Hero)
	storemodule.DB_UpdateRefreshFreeCount(gamedata.StoreType_Awake)
	storemodule.DB_UpdateRefreshFreeCount(gamedata.StoreType_Pet)
}

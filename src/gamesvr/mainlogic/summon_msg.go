package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"msg"
	"net/http"
	"time"
)

//! 玩家请求查询召唤刷新状态
func Hand_GetSummonStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 获取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSummonStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSummonStatus Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSummonStatus_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 更新状态状态
	player.SummonModule.UpdateSummonStatus()

	response.NormalSummon.SummonCounts = gamedata.NormalSummonFreeTimes - player.SummonModule.Normal.SummonCounts

	summonTime := player.SummonModule.Normal.SummonTime - time.Now().Unix()
	if summonTime < 0 {
		summonTime = 0
	}

	response.NormalSummon.SummonTime = summonTime
	response.SeniorSummon.Point = player.SummonModule.Senior.SummonPoint

	summonTime = player.SummonModule.Senior.SummonTime - time.Now().Unix()
	if summonTime < 0 {
		summonTime = 0
	}
	response.SeniorSummon.SummonTime = summonTime

	response.SeniorSummon.OrangeCount = 10 - player.SummonModule.Senior.OrangeCount
	response.Discount = gamedata.TenSummonDiscount
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求召唤
func Hand_GetSummon(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 获取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSummon_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSummon Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSummon_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 更新状态状态
	player.SummonModule.UpdateSummonStatus()

	//! 判断玩家英雄背包状态
	if player.BagMoudle.IsHeroBagFull() == true {
		response.RetCode = msg.RE_HERO_BAG_OVERLOAD
		gamelog.Error("Hand_GetSummon error: Hero bag is full")
		return
	}

	//! 判断种类
	if req.SummonType != gamedata.Summon_Normal && req.SummonType != gamedata.Summon_Senior {
		gamelog.Error("Hand_GetSummon error: invalid summonType: %d  PlayerID: %v", req.SummonType, player.playerid)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断类型
	if req.NumberType != 0 && req.NumberType != 1 {
		gamelog.Error("Hand_GetSummon error: invalid numberType: %d  PlayerID: %v", req.NumberType, player.playerid)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 获取对应配置信息
	summonConfig := gamedata.GetSummonConfig(req.SummonType)
	if summonConfig == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 根据召唤种类分别判断
	if req.SummonType == gamedata.Summon_Normal { //! 普通召唤
		if req.NumberType == 0 { //! 单抽逻辑处理
			player.BagMoudle.IsHeroBagFull()

			//! 检测免费次数
			hasFree := false
			if time.Now().Unix() >= player.SummonModule.Normal.SummonTime &&
				player.SummonModule.Normal.SummonCounts < gamedata.NormalSummonFreeTimes {
				hasFree = true

				response.IsFree = true

				//! 修改标记
				player.SummonModule.Normal.SummonCounts += 1
				player.SummonModule.Normal.SummonTime = time.Now().Unix() + int64(gamedata.NormalSummonFreeCDTime)
				go player.SummonModule.UpdateNormalSummon()
			}

			if hasFree == false {
				//! 检测道具数量
				if !player.BagMoudle.IsItemEnough(summonConfig.CostItemID, summonConfig.CostItemNum) {
					response.RetCode = msg.RE_NOT_ENOUGH_SUMMON_ITEM
					return
				}

				//! 扣除道具
				player.BagMoudle.RemoveNormalItem(summonConfig.CostItemID, summonConfig.CostItemNum)
			}

			//! 随机英雄
			heroLst := gamedata.GetSummonInfoRandom(gamedata.Summon_Normal, 1)
			heroID := heroLst[0]
			if player.SummonModule.IsFirst == true {

				if player.HeroMoudle.CurHeros[0].HeroID == 3 { //! 女主人公
					heroID = 407
					player.SummonModule.IsFirst = false
					go player.SummonModule.UpdateFirstSummon()
				} else {
					heroID = 428
					player.SummonModule.IsFirst = false
					go player.SummonModule.UpdateFirstSummon()
				}

			}
			player.BagMoudle.AddHeroByID(heroID, 1)
			response.HeroID = append(response.HeroID, heroID)
			response.RetCode = msg.RE_SUCCESS

			//! 更新状态状态
			player.SummonModule.UpdateSummonStatus()

			response.NormalSummon.SummonCounts = gamedata.NormalSummonFreeTimes - player.SummonModule.Normal.SummonCounts

			summonTime := player.SummonModule.Normal.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}

			response.NormalSummon.SummonTime = summonTime
			response.SeniorSummon.Point = player.SummonModule.Senior.SummonPoint

			summonTime = player.SummonModule.Senior.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}

			response.SeniorSummon.SummonTime = summonTime

			response.SeniorSummon.OrangeCount = 10 - player.SummonModule.Senior.OrangeCount

			return

		} else if req.NumberType == 1 { //! 十连抽逻辑处理
			//! 检测道具数量是否足够
			bEnough := player.BagMoudle.IsItemEnough(summonConfig.CostItemID, summonConfig.CostItemNum*10)
			if !bEnough {
				response.RetCode = msg.RE_NOT_ENOUGH_SUMMON_ITEM
				return
			}

			//! 扣除道具
			player.BagMoudle.RemoveNormalItem(summonConfig.CostItemID, summonConfig.CostItemNum*10)

			//! 随机英雄
			heroLst := gamedata.GetSummonInfoRandom(gamedata.Summon_Normal, 10)
			for _, v := range heroLst {
				player.BagMoudle.AddHeroByID(v, 1)
				response.HeroID = append(response.HeroID, v)
			}

			response.RetCode = msg.RE_SUCCESS

			//! 更新状态状态
			player.SummonModule.UpdateSummonStatus()

			response.NormalSummon.SummonCounts = gamedata.NormalSummonFreeTimes - player.SummonModule.Normal.SummonCounts

			summonTime := player.SummonModule.Normal.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}

			response.NormalSummon.SummonTime = summonTime
			response.SeniorSummon.Point = player.SummonModule.Senior.SummonPoint

			summonTime = player.SummonModule.Senior.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}
			response.SeniorSummon.SummonTime = summonTime

			response.SeniorSummon.OrangeCount = 10 - player.SummonModule.Senior.OrangeCount
			return
		}

	} else if req.SummonType == gamedata.Summon_Senior { //! 高级召唤
		if req.NumberType == 0 { //! 单抽逻辑处理
			//! 检测免费次数
			hasFree := false
			if time.Now().Unix() >= player.SummonModule.Senior.SummonTime {
				hasFree = true

				response.IsFree = true

				//! 修改标记
				player.SummonModule.Senior.SummonTime = time.Now().Unix() + int64(gamedata.SeniorSummonFreeCDTime)
				go player.SummonModule.UpdateSeniorSummon()
			}

			if hasFree == false {

				//! 检查道具是否足够
				itemEnough := true
				if !player.BagMoudle.IsItemEnough(summonConfig.CostItemID, summonConfig.CostItemNum) {

					//! 检查货币是否足够
					itemEnough = false
					if player.RoleMoudle.CheckMoneyEnough(summonConfig.CostMoneyID, summonConfig.CostMoneyNum) == false {
						response.RetCode = msg.RE_NOT_ENOUGH_MONEY
						return
					}

					//! 扣除金钱
					player.RoleMoudle.CostMoney(summonConfig.CostMoneyID, summonConfig.CostMoneyNum)
				}

				if itemEnough == true {
					//! 扣除道具
					player.BagMoudle.RemoveNormalItem(summonConfig.CostItemID, summonConfig.CostItemNum)
				}
			}

			//! 增加积分
			player.SummonModule.Senior.SummonPoint += gamedata.SeniorSummonPoint
			if player.SummonModule.Senior.SummonPoint > summonConfig.NeedPoint {
				player.SummonModule.Senior.SummonPoint = summonConfig.NeedPoint
			}

			//! 增加抽取次数
			player.SummonModule.Senior.OrangeCount += 1

			//! 随机英雄
			heroID := 0
			if player.SummonModule.Senior.OrangeCount >= 10 {
				heroID = gamedata.GetSummonInfoOrangeRandom()

				//! 重置抽取次数
				player.SummonModule.Senior.OrangeCount = 0
			} else {
				heroLst := gamedata.GetSummonInfoRandom(gamedata.Summon_Senior, 1)
				if len(heroLst) != 1 {
					gamelog.Error("GetSummonInfoRandom Error: too long heroLst")
					response.RetCode = msg.RE_UNKNOWN_ERR
					return
				}

				heroID = heroLst[0]
			}

			player.BagMoudle.AddHeroByID(heroID, 1)
			response.HeroID = append(response.HeroID, heroID)
			response.RetCode = msg.RE_SUCCESS

			player.SummonModule.UpdateSeniorSummon()

			//! 更新状态状态
			go player.SummonModule.UpdateSummonStatus()

			response.NormalSummon.SummonCounts = gamedata.NormalSummonFreeTimes - player.SummonModule.Normal.SummonCounts

			summonTime := player.SummonModule.Normal.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}

			response.NormalSummon.SummonTime = summonTime
			response.SeniorSummon.Point = player.SummonModule.Senior.SummonPoint

			summonTime = player.SummonModule.Senior.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}
			response.SeniorSummon.SummonTime = summonTime

			response.SeniorSummon.OrangeCount = 10 - player.SummonModule.Senior.OrangeCount
			player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SENIOR_SUMMON, 1)
			return

		} else if req.NumberType == 1 { //! 十连抽逻辑处理
			//! 检测道具是否足够
			itemEnougth := true
			if !player.BagMoudle.IsItemEnough(summonConfig.CostItemID, summonConfig.CostItemNum*10) {
				itemEnougth = false

				//! 全部花费元宝
				costMoney := summonConfig.CostMoneyNum * 10 * gamedata.TenSummonDiscount / 100
				if player.RoleMoudle.CheckMoneyEnough(summonConfig.CostMoneyID, costMoney) == false {
					response.RetCode = msg.RE_NOT_ENOUGH_MONEY
					return
				}

				//! 扣除元宝
				player.RoleMoudle.CostMoney(summonConfig.CostMoneyID, costMoney)
			}

			if itemEnougth == true {
				//! 扣除道具
				player.BagMoudle.RemoveNormalItem(summonConfig.CostItemID, summonConfig.CostItemNum*10)
			}

			orangeTimes := 10 - player.SummonModule.Senior.OrangeCount

			//! 随机英雄
			heroLst := gamedata.GetSummonInfoRandom(gamedata.Summon_Senior, orangeTimes-1)
			orange := gamedata.GetSummonInfoOrangeRandom()
			heroLst = append(heroLst, gamedata.GetSummonInfoRandom(gamedata.Summon_Senior, 10-orangeTimes)...)

			random := rand.New(rand.NewSource(time.Now().UnixNano()))
			index := random.Intn(len(heroLst))
			newLst := IntLst{}
			for i := 0; i < len(heroLst); i++ {
				if gamedata.GetHeroInfo(heroLst[i]).Quality == 5 {
					//! 同时存在两个橙将,则去除这个替换成普通
					for {
						heroid := gamedata.GetSummonInfoRandom(gamedata.Summon_Senior, 1)
						if gamedata.GetHeroInfo(heroid[0]).Quality != 5 {
							newLst.Add(heroid[0])
							break
						}
					}

					continue
				}

				newLst.Add(heroLst[i])
				if index == i {
					newLst.Add(orange)
				}
			}
			for _, v := range newLst {
				player.BagMoudle.AddHeroByID(v, 1)
				response.HeroID = append(response.HeroID, v)
			}

			//! 增加积分
			player.SummonModule.Senior.SummonPoint += 10 * gamedata.SeniorSummonPoint
			if player.SummonModule.Senior.SummonPoint > summonConfig.NeedPoint {
				player.SummonModule.Senior.SummonPoint = summonConfig.NeedPoint
			}

			go player.SummonModule.UpdateSeniorSummon()
			response.RetCode = msg.RE_SUCCESS

			//! 更新状态状态
			player.SummonModule.UpdateSummonStatus()

			response.NormalSummon.SummonCounts = gamedata.NormalSummonFreeTimes - player.SummonModule.Normal.SummonCounts

			summonTime := player.SummonModule.Normal.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}

			response.NormalSummon.SummonTime = summonTime
			response.SeniorSummon.Point = player.SummonModule.Senior.SummonPoint

			summonTime = player.SummonModule.Senior.SummonTime - time.Now().Unix()
			if summonTime < 0 {
				summonTime = 0
			}
			response.SeniorSummon.SummonTime = summonTime

			response.SeniorSummon.OrangeCount = 10 - player.SummonModule.Senior.OrangeCount
			player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SENIOR_SUMMON, 10)

			return
		}
	} else {
		//! 异常参数
		gamelog.Error("Hand_GetSummon Error: Invalid SummonType  %d", req.SummonType)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}
}

//! 玩家请求积分兑换英雄
func Hand_ExchangeHero(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 获取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ExchangeHero_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ExchangeHero Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_ExchangeHero_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		gamelog.Info("return: %s", b)
		w.Write(b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 更新状态状态
	player.SummonModule.UpdateSummonStatus()

	//! 获取对应配置信息
	summonConfig := gamedata.GetSummonConfig(gamedata.Summon_Senior)
	if summonConfig == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测积分是否足够
	if player.SummonModule.Senior.SummonPoint < summonConfig.NeedPoint {
		response.RetCode = msg.RE_NOT_ENOUGH_POINT
		return
	}

	//! 检查英雄是否属于兑换队列
	isExist := false
	pAwardItem, ok := gamedata.GT_AwardList[gamedata.OrangeSummonAwardID]
	if pAwardItem == nil || !ok {
		gamelog.Error("GetItemsFromAwardID Error: Invalid awardid :%d", gamedata.OrangeSummonAwardID)
		return
	}

	if pAwardItem.RatioItems == nil {
		gamelog.Error("GetItemsFromAwardID Error: Invalid awardid :%d", gamedata.OrangeSummonAwardID)
		return
	}
	for _, v := range pAwardItem.RatioItems {
		if v.ItemID == req.HeroID {
			isExist = true
			break
		}
	}
	if isExist == false {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 扣除积分
	player.SummonModule.Senior.SummonPoint -= summonConfig.NeedPoint

	//! 赐予英雄
	player.BagMoudle.AddHeroByID(req.HeroID, 1)

	//! 存储数据
	go player.SummonModule.UpdateSeniorSummon()

	response.RetCode = msg.RE_SUCCESS

}

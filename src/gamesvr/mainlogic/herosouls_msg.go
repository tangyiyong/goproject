package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"msg"
	"net/http"
	"time"
	"utility"
)

//! 激活/升级将灵
func Hand_ActivateHeroSouls(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ActivateHeroSouls_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ActivateHeroSouls Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_ActivateHeroSouls_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	//! 获取链接信息
	pHeroSoulInfo := gamedata.GetHeroSoulsInfo(req.ID)
	if pHeroSoulInfo == nil {
		gamelog.Error("Hand_ActivateHeroSouls Error: Invalid id %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 获取所属章节
	chapter := gamedata.GetHeroSoulsBelongChapter(req.ID)
	if chapter == 0 {
		gamelog.Error("Hand_ActivateHeroSouls Error: Invalid id %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	} else if chapter > player.HeroSoulsModule.UnLockChapter {
		gamelog.Error("Hand_ActivateHeroSouls Error: Not Unlock chapter %d", chapter)
		response.RetCode = msg.RE_NOT_UNLOCK
		return
	}

	//! 检测将灵背包是否有足够该将灵
	for _, v := range pHeroSoulInfo.HeroIDs {
		if v == 0 {
			continue
		}

		count := player.BagMoudle.GetHeroSoulCount(v)
		if count < 1 {
			gamelog.Error("Hand_ActivateHeroSouls Error: Hero souls not enough, heroID: %d", v)
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}
	}

	//! 检测当前将灵等级是否能够升级
	for _, v := range player.HeroSoulsModule.HeroSoulsLink {
		if v.ID == req.ID {
			if v.Level >= 10 {
				gamelog.Error("Hand_ActivateHeroSouls Error: Has been to full level")
				return
			}

			needLevel := v.Level*5 + 50
			if player.GetLevel() < needLevel {
				gamelog.Error("Hand_ActivateHeroSouls Error: Level not enough, needLevel: %d", needLevel)
				return
			}
		}
	}

	//! 扣除将灵
	awardsouls := 0
	for _, v := range pHeroSoulInfo.HeroIDs {
		if v == 0 {
			continue
		}

		awardsouls += 1
		player.BagMoudle.RemoveHeroSoul(v, 1)
	}

	//! 激活/升级将灵
	isExist := false
	for i, v := range player.HeroSoulsModule.HeroSoulsLink {
		if v.ID == req.ID {
			isExist = true

			//! 升级将灵
			player.HeroSoulsModule.HeroSoulsLink[i].Level += 1
			player.HeroSoulsModule.DB_UpdateHeroSoulsLinkLevel(i, player.HeroSoulsModule.HeroSoulsLink[i].Level)

			for _, n := range pHeroSoulInfo.Property {
				if n.PropertyID != 0 {
					player.HeroMoudle.AddExtraProperty(n.PropertyID, int32(n.LevelUp), n.Is_Percent, n.Camp)
					player.HeroSoulsModule.AddTempProperty(n.PropertyID, n.LevelUp, n.Is_Percent, n.Camp)
				}
			}

			player.HeroMoudle.DB_SaveExtraProperty()

		}
	}

	if isExist == false {
		//! 奖励阵魂值
		player.HeroSoulsModule.SoulMapValue += awardsouls
		player.HeroSoulsModule.DB_UpdateSoulMapValue()

		//! 激活将灵
		player.HeroSoulsModule.HeroSoulsLink = append(player.HeroSoulsModule.HeroSoulsLink, msg.MSG_HeroSoulsLink{req.ID, 1})
		player.HeroSoulsModule.DB_AddHeroSoulsLink(msg.MSG_HeroSoulsLink{req.ID, 1})

		for _, n := range pHeroSoulInfo.Property {
			if n.PropertyID != 0 {
				player.HeroMoudle.AddExtraProperty(n.PropertyID, int32(n.PropertyValue), n.Is_Percent, n.Camp)
				player.HeroSoulsModule.AddTempProperty(n.PropertyID, n.PropertyValue, n.Is_Percent, n.Camp)
			}
		}

		player.HeroMoudle.DB_SaveExtraProperty()
	}

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.UnLockChapter = player.HeroSoulsModule.UnLockChapter

	//! 计入排行榜
	G_HeroSoulsRanker.SetRankItem(req.PlayerID, player.HeroSoulsModule.SoulMapValue*10000+len(player.HeroSoulsModule.HeroSoulsLink))

	//! 检测解锁
	nextChapter := player.HeroSoulsModule.UnLockChapter + 1
	if nextChapter >= gamedata.GetHeroSoulsChapterCount() { //! 已至最终章,无法解锁下一章节
		return
	}

	nextChapterInfo := gamedata.GetHeroSoulsChapterInfo(nextChapter)

	count := 0
	for _, v := range player.HeroSoulsModule.HeroSoulsLink {
		if gamedata.GetHeroSoulsBelongChapter(v.ID) == nextChapterInfo.UnLockChapter {
			count++
		}
	}

	if count >= nextChapterInfo.UnlockCount {
		//! 解锁下一章节
		player.HeroSoulsModule.UnLockChapter = nextChapter
		player.HeroSoulsModule.DB_UnLockChapter()

		//! 返回解锁章节
		response.UnLockChapter = player.HeroSoulsModule.UnLockChapter
	}

}

//! 详细查询章节将灵激活情况
func Hand_QueryChapterHeroSoulsDetail(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHeroSoulsChapter_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryChapterHeroSoulsDetail Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_QueryHeroSoulsChapter_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	response.HeroSouls = player.HeroSoulsModule.HeroSoulsLink
	response.UnLockChapter = player.HeroSoulsModule.UnLockChapter
	response.RetCode = msg.RE_SUCCESS
}

//! 获取将灵轮盘
func Hand_GetHeroSoulsLst(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetHeroSoulsLst_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetHeroSoulsLst Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_GetHeroSoulsLst_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	for _, v := range player.HeroSoulsModule.HeroSoulsLst {
		var souls msg.THeroSouls
		souls.ID = v.ID
		souls.HeroID = v.HeroID
		souls.IsExist = v.IsExist
		response.HeroSoulsLst = append(response.HeroSoulsLst, souls)
	}

	//! 返回成功
	response.CountDown = player.HeroSoulsModule.CheckStoreRefresh()
	response.TargetIndex = player.HeroSoulsModule.TargetIndex
	response.ChallengeTimes = player.HeroSoulsModule.LeftTimes
	response.BuyChallengeTimes = player.HeroSoulsModule.BuyTimes
	response.RetCode = msg.RE_SUCCESS
}

//! 刷新将灵指针
func Hand_RefreshHeroSoulsLst(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RefreshHeroSoulsLst_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_RefreshHeroSoulsLst Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_RefreshHeroSoulsLst_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	//! 计算花费货币
	costMoneyID := gamedata.HeroSoulsRefreshCostMoneyID

	existCount := 0
	randArray := IntLst{}
	for i, v := range player.HeroSoulsModule.HeroSoulsLst {
		if v.IsExist == true {
			randArray.Add(i)
			existCount++
		}
	}

	if existCount == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_RefreshHeroSoulsLst Error: Not have exist hero soul")
		return
	}

	costMoneyNum := 10 + gamedata.HeroSoulsRefreshCostMoneyValue*(8-existCount)

	//! 检查花费是否足够
	if player.RoleMoudle.CheckMoneyEnough(costMoneyID, costMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_RefreshHeroSoulsLst Error: Not enough money  moneyID: %d   moneyNum: %d ", costMoneyID, costMoneyNum)
		return
	}

	//! 扣除花费
	player.RoleMoudle.CostMoney(costMoneyID, costMoneyNum)

	//! 获取目标将灵品质
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	player.HeroSoulsModule.TargetIndex = randArray[random.Intn(len(randArray))]
	player.HeroSoulsModule.DB_UpdateTargetIndex()

	itemInfo := gamedata.GetItemInfo(player.HeroSoulsModule.HeroSoulsLst[player.HeroSoulsModule.TargetIndex].HeroID)

	//! 给予英魂奖励
	getMoneyNum := gamedata.HeroSoulsRefreshGetMoneyValue * itemInfo.Quality
	player.RoleMoudle.AddMoney(gamedata.HeroSoulsRefreshGetMoneyID, getMoneyNum)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.TargetIndex = player.HeroSoulsModule.TargetIndex
	response.GetMoneyID = gamedata.HeroSoulsRefreshGetMoneyID
	response.GetMoneyNum = getMoneyNum
	response.CostMoneyID = costMoneyID
	response.CostMoneyNum = costMoneyNum
}

//! 挑战将灵
func Hand_ChallengeHeroSouls(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	if false == utility.MsgDataCheck(buffer, G_XorCode) {
		//存在作弊的可能
		gamelog.Error("Hand_ChallengeHeroSouls : Message Data Check Error!!!!")
		return
	}
	var req msg.MSG_ChallengeHeroSouls_Req
	if json.Unmarshal(buffer[:len(buffer)-16], &req) != nil {
		gamelog.Error("Hand_ChallengeHeroSouls : Unmarshal error!!!!")
		return
	}

	//! 创建回复
	var response msg.MSG_ChallengeHeroSouls_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	//检查英雄数据是否一致
	if !player.CheckHeroData(req.HeroCkD) {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChallengeHeroSouls : CheckHeroData Error!!!!")
		return
	}

	player.HeroSoulsModule.CheckReset()

	//! 检测当前指向将灵是否存在
	heroSouls := player.HeroSoulsModule.HeroSoulsLst[player.HeroSoulsModule.TargetIndex]
	if heroSouls.IsExist == false {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ChallengeHeroSouls Error: Hero souls not exist")
		return
	}

	//! 检测挑战次数是否超过
	if player.HeroSoulsModule.LeftTimes <= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		gamelog.Error("Hand_ChallengeHeroSouls Error: ChallengeTimes not enough")
		return
	}

	//! 扣除挑战次数
	player.HeroSoulsModule.LeftTimes -= 1
	player.HeroSoulsModule.DB_UpdateChallengeHeroSoulsTimes()

	//! 设置将灵状态
	player.HeroSoulsModule.HeroSoulsLst[player.HeroSoulsModule.TargetIndex].IsExist = false
	player.HeroSoulsModule.DB_UpdateHeroSoulsMark(player.HeroSoulsModule.TargetIndex)

	//! 获取英灵
	player.BagMoudle.AddHeroSoul(heroSouls.HeroID, 1)
	response.RetCode = msg.RE_SUCCESS

	existCount := 0
	randArray := IntLst{}
	for i, v := range player.HeroSoulsModule.HeroSoulsLst {
		if v.IsExist == true {
			randArray.Add(i)
			existCount++
		}
	}

	if existCount <= 0 {
		return
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	player.HeroSoulsModule.TargetIndex = randArray[random.Intn(len(randArray))]
	player.HeroSoulsModule.DB_UpdateTargetIndex()
}

//! 购买挑战英灵次数
func Hand_BuyChallengeHeroSoulsTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyChallengeHeroSoulsTimes_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BuyChallengeHeroSoulsTimes Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_BuyChallengeHeroSoulsTimes_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	//! 检测购买次数是否超过
	buyTimesLimit := gamedata.GetFuncVipValue(gamedata.FUNC_HEROSOULS_TIMES, player.GetVipLevel())
	if buyTimesLimit < player.HeroSoulsModule.BuyTimes+req.Times {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		gamelog.Error("Hand_BuyChallengeHeroSoulsTimes Error: Buy times not enough")
		return
	}

	costMoneyNum := 0

	//! 计算花费
	for i := 1; i <= req.Times; i++ {
		times := player.HeroSoulsModule.LeftTimes + i
		costMoneyNum += gamedata.GetFuncTimeCost(gamedata.FUNC_HEROSOULS_TIMES, times)
	}

	//! 判断货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.BuyChallengeTimesMoneyID, costMoneyNum) == false {
		gamelog.Error("Hand_BuyChallengeHeroSoulsTimes Error: Money not enough")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除花费
	player.RoleMoudle.CostMoney(gamedata.BuyChallengeTimesMoneyID, costMoneyNum)

	//! 增加次数
	player.HeroSoulsModule.BuyTimes += req.Times
	player.HeroSoulsModule.LeftTimes += req.Times
	player.HeroSoulsModule.DB_BuyChallengeHeroSoulsTimes()

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
	response.CostMoneyID = gamedata.BuyChallengeTimesMoneyID
	response.CostMoneyNum = costMoneyNum
}

//! 重置将灵列表
func Hand_ResetHeroSoulsLst(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ResetHeroSoulsLst_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ResetHeroSoulsLst Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_ResetHeroSoulsLst_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	player.HeroSoulsModule.ResetHeroSoulsLst(true)

	for _, v := range player.HeroSoulsModule.HeroSoulsLst {
		var souls msg.THeroSouls
		souls.HeroID = v.HeroID
		souls.IsExist = v.IsExist
		response.HeroSoulsLst = append(response.HeroSoulsLst, souls)
	}

	response.TargetIndex = player.HeroSoulsModule.TargetIndex
	response.RetCode = msg.RE_SUCCESS
}

//! 查询英灵排行榜
func Hand_QueryHeroSoulsRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHeroSoulsRank_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryHeroSoulsRank Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_QueryHeroSoulsRank_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	response.RankLst = []msg.THeroSoulsRank{}
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	for i := 0; i < 10; i++ {
		playerInfo := G_HeroSoulsRanker.List[i]
		if playerInfo.RankID != 0 {
			var rankInfo msg.THeroSoulsRank

			simpleInfo := G_SimpleMgr.GetSimpleInfoByID(playerInfo.RankID)
			rankInfo.HeroID = simpleInfo.HeroID
			rankInfo.Name = simpleInfo.Name
			rankInfo.Quality = simpleInfo.Quality
			rankInfo.SoulsValue, rankInfo.SoulsCount = playerInfo.RankValue/10000, playerInfo.RankValue%10000
			response.RankLst = append(response.RankLst, rankInfo)
		}
	}

	response.SelfRank = -1
	for i, v := range G_HeroSoulsRanker.List {
		if v.RankID == req.PlayerID {
			response.SelfRank = i + 1
		}
	}
	response.SelfSoulsCount = len(player.HeroSoulsModule.HeroSoulsLink)
	response.SelfSoulsValue = player.HeroSoulsModule.SoulMapValue

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家查询英灵商店信息
func Hand_QueryHeroSoulsStoreInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHeroSoulsStore_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QueryHeroSoulsStoreInfo Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_QueryHeroSoulsStore_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	//! 检测商店刷新
	response.CountDown = player.HeroSoulsModule.CheckStoreRefresh()

	for _, v := range player.HeroSoulsModule.HeroSoulsStoreLst {
		var goods msg.THeroSoulsStore
		goods.ItemID = v.ItemID
		goods.MoneyID = v.MoneyID
		goods.MoneyNum = v.MoneyNum
		goods.IsBuy = v.IsBuy
		response.GoodsLst = append(response.GoodsLst, goods)
	}

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买英灵商店
func Hand_BuyHeroSoulsStoreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuyHeroSouls_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_BuyHeroSoulsStoreItem Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_BuyHeroSouls_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	//! 检测商店刷新
	player.HeroSoulsModule.CheckStoreRefresh()

	//! 检测商品是否存在
	isExist := false
	index := 0
	var goodsInfo *THeroSoulsStore
	for i, v := range player.HeroSoulsModule.HeroSoulsStoreLst {
		if v.ItemID == req.ItemID {
			isExist = true
			goodsInfo = &player.HeroSoulsModule.HeroSoulsStoreLst[i]
			index = i
		}
	}

	if isExist == false {
		gamelog.Error("Hand_BuyHeroSoulsStoreItem Error: Not exist itemID: %d", req.ItemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测是否已经购买
	if goodsInfo.IsBuy == true {
		response.RetCode = msg.RE_ALEADY_BUY
		gamelog.Error("Hand_BuyHeroSoulsStoreItem Error: Aleady buy ItemID: %d", req.ItemID)
		return
	}

	//! 检查货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(goodsInfo.MoneyID, goodsInfo.MoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_BuyHeroSoulsStoreItem Error: Not enough money")
		return
	}

	//! 扣除货币
	player.RoleMoudle.CostMoney(goodsInfo.MoneyID, goodsInfo.MoneyNum)
	response.CostMoneyID, response.CostMoneyNum = goodsInfo.MoneyID, goodsInfo.MoneyNum

	//! 修改商品购买标记
	goodsInfo.IsBuy = true
	player.HeroSoulsModule.DB_UpdateStoreGoodsStatus(index, true)

	//! 给予物品
	player.BagMoudle.AddHeroSoul(goodsInfo.ItemID, 1)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 查询阵图成就信息
func Hand_QuerySoulMapInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHeroSoulsAchievement_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_QuerySoulMapInfo Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_QueryHeroSoulsAchievement_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.HeroSoulsModule.CheckReset()

	response.SoulMapValue = player.HeroSoulsModule.SoulMapValue
	response.Achievement = player.HeroSoulsModule.Achievement

	response.RetCode = msg.RE_SUCCESS
}

//! 请求激活阵图成就
func Hand_ActivateheroSoulsAchievement(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_ActivateHeroSoulsAchievement_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ActivateheroSoulsAchievement Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_ActivateHeroSoulsAchievement_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	nextLevel := gamedata.GetSoulMapInfo(player.HeroSoulsModule.Achievement + 1)
	if nextLevel == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_ActivateheroSoulsAchievement Error: Invalid achiecement %d", player.HeroSoulsModule.Achievement+1)
		return
	}

	//! 检查阵图值是否足够
	if nextLevel.Souls > player.HeroSoulsModule.SoulMapValue {
		response.RetCode = msg.RE_NOT_ENOUGH_VALUE
		gamelog.Error("Hand_ActivateheroSoulsAchievement Error: Soul map value is not enough")
		return
	}

	//! 增加阵图成就
	player.HeroSoulsModule.Achievement += 1
	player.HeroSoulsModule.DB_UpdateSoulMapAchievement()

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

//! 查询属性加成
func Hand_QueryHeroSoulsPerproty(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 读取消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_QueryHeroSoulsProperty_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_ActivateheroSoulsAchievement Error: Unmarshal fail")
		return
	}

	//! 创建回复
	var response msg.MSG_QueryHeroSoulsProperty_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 常规检查
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.PropertyPercentLst = player.HeroSoulsModule.propertyPercent
	response.PropertyIntLst = player.HeroSoulsModule.propertyInt
	response.CampPropertyKillLst = player.HeroSoulsModule.campPropertyKillLst
	response.CampPropertyDefenceLst = player.HeroSoulsModule.campPropertyDefenceLst
	response.RetCode = msg.RE_SUCCESS
}

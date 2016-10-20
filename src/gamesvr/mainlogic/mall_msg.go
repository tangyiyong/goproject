package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 用户查询VIP礼包信息
func Hand_GetVipGiftInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析信息
	var req msg.MSG_GetVipGifts_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetVipGiftInfo Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_GetVipGifts_Ack
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

	//! 获取所有VIP礼包
	vipGiftLst := gamedata.GetVipGiftBagIDs()
	for _, v := range vipGiftLst {
		if player.MallModule.VipBagRecord.IsExist(v) >= 0 {
			continue
		}

		mallItem := gamedata.GetMallItemInfo(v)
		if mallItem == nil {
			gamelog.Error("Hand_GetVipGiftInfo Error: Invalid itemID: %d", v)
			return
		}

		funcID := gamedata.GetFuncID(mallItem.ItemID)
		isOpen := gamedata.IsFuncOpen(funcID, player.GetLevel(), player.GetVipLevel()+1)
		if isOpen == true {
			response.ID = append(response.ID, int32(mallItem.ItemID))
		}
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 用户购买VIP礼包信息
func Hand_BuyVipGift(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析信息
	var req msg.MSG_BuyVipGift_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyVipGift Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_BuyVipGift_Ack
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

	//! 检查是否重复购买
	if player.MallModule.VipBagRecord.IsExist(req.ID) >= 0 {
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		gamelog.Error("Hand_BuyVipGift Error: Not enough buy times. VipGiftID: %d", req.ID)
		return
	}

	//! 检查VIP等级是否足够
	itemData := gamedata.GetMallItemInfo(req.ID)
	if itemData == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	vipLevel := player.GetVipLevel()
	funcID := gamedata.GetFuncID(itemData.ItemID)
	if funcID == 0 {
		gamelog.Error("Hand_BuyVipGift Error:GetFuncID fail. itemID: %d", itemData.ItemID)
		return
	}

	if gamedata.IsFuncOpen(funcID, player.GetLevel(), vipLevel) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_VIP_LVL
		return
	}

	//! 检查金钱是否足够
	needMoney := gamedata.GetFuncTimeCost(funcID, 1)
	if needMoney <= 0 {
		if needMoney == -1 {
			response.RetCode = msg.RE_NOT_ENOUGH_TIMES
			return
		}
		gamelog.Error("Hand_BuyVipGift Error:GetResetCost fail.")
		return
	}

	if player.RoleMoudle.CheckMoneyEnough(gamedata.MallGiftMoneyID, needMoney) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(gamedata.MallGiftMoneyID, needMoney)
	response.MoneyID = gamedata.MallGiftMoneyID
	response.MoneyNum = needMoney

	//! 记录购买

	player.MallModule.VipBagRecord = append(player.MallModule.VipBagRecord, req.ID)
	player.MallModule.UpdateVipBagRecord()

	//! 发送礼包
	player.BagMoudle.AddAwardItem(int(itemData.ItemID), 1)
	response.ItemID, response.ItemNum = itemData.ItemID, 1

	//! 获取所有VIP礼包
	vipGiftLst := gamedata.GetVipGiftBagIDs()
	for _, v := range vipGiftLst {
		if player.MallModule.VipBagRecord.IsExist(v) >= 0 {
			continue
		}
		mallItem := gamedata.GetMallItemInfo(v)
		if mallItem == nil {
			gamelog.Error("Hand_BuyVipGift Error:: Invalid itemID: %d", v)
			return
		}

		funcID := gamedata.GetFuncID(mallItem.ItemID)
		isOpen := gamedata.IsFuncOpen(funcID, player.GetLevel(), player.GetVipLevel()+1)
		if isOpen == true {
			response.ID = append(response.ID, int32(mallItem.ItemID))
		}
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询道具购买次数
func Hand_GetMallBuyTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析信息
	var req msg.MSG_GetMallBuyTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetMallGoodsBuyTimes Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_GetMallBuyTimes_Ack
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

	//! 检查重置
	player.MallModule.CheckReset()
	response.BuyTimesLst = []msg.MSG_BuyData{}
	for i := 0; i < len(player.MallModule.NormalRecord); i++ {
		var info msg.MSG_BuyData
		info.ID = player.MallModule.NormalRecord[i].ID

		funcID := gamedata.GetFuncID(player.MallModule.NormalRecord[i].ID)
		totalTimes := gamedata.GetFuncVipValue(funcID, player.GetVipLevel())
		info.Times = player.MallModule.NormalRecord[i].Times
		if totalTimes == 0 { //! 判断是否限购
			info.Times = -2
		}

		response.BuyTimesLst = append(response.BuyTimesLst, info)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买商城商品
func Hand_BuyMallGoods(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析信息
	var req msg.MSG_BuyGoods_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyMallGoods Error : Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_BuyGoods_Ack
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

	//! 检查重置
	player.MallModule.CheckReset()

	//! 检查参数
	if req.Num <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_BuyMallGoods Error : invalid itemNum: %d  playerID: %v", req.Num, player.playerid)
		return
	}

	//! 获取物品信息
	itemInfo := gamedata.GetMallItemInfo(req.ID)
	if itemInfo == nil {
		gamelog.Error("Hand_BuyMallGoods Error : invalid itemID: %d", req.ID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 该协议用于购买普通商品
	if itemInfo.Type != 0 {
		gamelog.Error("Hand_BuyMallGoods Error :invalid itemType: %d", itemInfo.Type)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测当前用户等级是否能够购买该物品
	vipLevel := player.GetVipLevel()
	funcID := gamedata.GetFuncID(itemInfo.ItemID)
	if funcID == 0 {
		gamelog.Error("Hand_BuyMallGoods Error : GetFunID fail. itemID: %d", itemInfo.ItemID)
		return
	}

	totalTimes := gamedata.GetFuncVipValue(funcID, vipLevel)
	if gamedata.IsFuncOpen(funcID, player.GetLevel(), vipLevel) == false {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	pBuyData := player.MallModule.GetItemBuyData(req.ID)
	//! 检测用户购买次数
	curTimes := 0
	if totalTimes != 0 { //! 该商品有限购次数
		if pBuyData == nil {
			//! 初次购买,初始化商品信息
			if req.Num > totalTimes {
				response.RetCode = msg.RE_NOT_ENOUGH_TIMES
				return
			}
			curTimes = 0
		} else {
			if pBuyData.Times+req.Num > totalTimes {
				response.RetCode = msg.RE_NOT_ENOUGH_TIMES
				return
			}
			curTimes = pBuyData.Times
		}

	}

	//! 检测用户金钱是否足够
	needMoney := 0
	for i := 0; i < req.Num; i++ {
		money := gamedata.GetFuncTimeCost(funcID, curTimes+i+1)
		if money <= 0 {
			if money == -1 {
				response.RetCode = msg.RE_NOT_ENOUGH_MONEY
				return
			}
			break
		}

		needMoney += money
	}

	if player.RoleMoudle.CheckMoneyEnough(gamedata.MallItemMoneyID, needMoney) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_BuyMallGoods Error:Money not enough")
		return
	}

	if pBuyData == nil {
		player.MallModule.NormalRecord = append(player.MallModule.NormalRecord, msg.MSG_BuyData{req.ID, req.Num})
	} else {
		pBuyData.Times += req.Num
	}

	player.MallModule.UpdateNormalRecord()

	//! 扣除金钱
	player.RoleMoudle.CostMoney(gamedata.MallItemMoneyID, needMoney)
	response.MoneyID = gamedata.MallItemMoneyID
	response.MoneyNum = needMoney

	//! 发放物品
	player.BagMoudle.AddAwardItem(int(itemInfo.ItemID), req.Num*itemInfo.ItemNum)

	//! 返回购买信息
	response.BuyTimes.ID = req.ID
	response.BuyTimes.Times = curTimes + req.Num

	//! 完成交易
	response.RetCode = msg.RE_SUCCESS

	//! 增加进度 应策划要求写死ID
	if itemInfo.ItemID == 102 {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_BUY_ZHENGTAOLING, req.Num)
	} else if itemInfo.ItemID == 100 {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_BUY_ACTION_STRENGTH, req.Num)
	} else if itemInfo.ItemID == 101 {
		player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_BUY_ACTION_ENERGY, req.Num)
	}

}

//! 玩家查询单个物品购买次数
func Hand_GetGoodsBuyTimes(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接受信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析信息
	var req msg.MSG_GetBuyTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetGoodsBuyTimes Error : Unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建返回消息
	var response msg.MSG_GetBuyTimes_Ack
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

	//! 检查重置
	player.MallModule.CheckReset()

	isExist := false
	funcID := gamedata.GetFuncID(req.ItemID)
	response.FuncID = funcID
	totalTimes := gamedata.GetFuncVipValue(funcID, player.GetVipLevel())
	for i := 0; i < len(player.MallModule.NormalRecord); i++ {
		if player.MallModule.NormalRecord[i].ID == req.ItemID {

			if totalTimes == 0 { //! 判断是否限购
				response.BuyTimes = -1
			} else {
				response.BuyTimes = totalTimes - player.MallModule.NormalRecord[i].Times
			}

			isExist = true
			break
		}
	}

	if isExist == false {
		response.BuyTimes = totalTimes
	}

	response.ItemID = req.ItemID
	response.RetCode = msg.RE_SUCCESS
}

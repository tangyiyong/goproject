package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求七日活动信息
func Hand_GetSevenActivityInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSevenActivity_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSevenActivityInfo : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetSevenActivity_Ack
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

	player.ActivityModule.CheckReset()
	var activity *TActivitySevenDay
	for i, v := range player.ActivityModule.SevenDay {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.SevenDay[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_GetSevenDay Error: Activity not exist ID: %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("Hand_GetSevenDay Error: Activity not open ID: %d", req.ActivityID)
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	for _, v := range activity.TaskList {
		var taskInfo msg.TTaskInfo
		taskInfo.TaskID = v.TaskID
		taskInfo.TaskStatus = v.TaskStatus
		taskInfo.TaskCount = v.TaskCount
		response.SevenActivityLst = append(response.SevenActivityLst, taskInfo)
	}

	response.BuyLst = activity.BuyLst
	response.LimitInfo = G_GlobalVariables.GetSevenDayLimit(req.ActivityID).LimitBuy
	response.RetCode = msg.RE_SUCCESS
	response.ActivityID = activity.ActivityID
	response.OpenDay = GetOpenServerDay()
}

//! 玩家请求领取七日活动奖励
func Hand_GetSevenActivityAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSevenActivityAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSevenActivityAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	gamelog.Info("Receive: %s", buffer)

	var response msg.MSG_GetSevenActivityAward_Ack
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

	player.ActivityModule.CheckReset()
	var activity *TActivitySevenDay
	for i, v := range player.ActivityModule.SevenDay {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.SevenDay[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_GetSevenDay Error: Activity not exist ID: %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("Hand_GetSevenDay Error: Activity not open ID: %d", req.ActivityID)
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	taskData := gamedata.GetSevenTaskInfo(req.TaskID)
	if taskData == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("GetSevenTaskInfo Fail. taskID: %d", req.TaskID)
		return
	}

	var taskInfo *TTaskInfo
	for i, _ := range activity.TaskList {
		if activity.TaskList[i].TaskID == taskData.TaskID {
			if activity.TaskList[i].TaskCount < taskData.Count {
				response.RetCode = msg.RE_TASK_NOT_COMPLETE
				return
			} else if activity.TaskList[i].TaskStatus == Task_Received {
				response.RetCode = msg.RE_ALREADY_RECEIVED
				return
			} else {
				taskInfo = &activity.TaskList[i]
			}
		}
	}

	if taskInfo == nil {
		gamelog.Error("Hand_GetSevenDay Error: Can't find task id: %d", req.TaskID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	awardItems := gamedata.GetItemsFromAwardID(taskData.AwardItem)

	//!  三选一类型奖励
	if taskData.IsSelectOne == 1 {
		isExist := 0
		itemNum := 0
		for _, v := range awardItems {
			if v.ItemID == req.ItemID {
				isExist = 1
				itemNum = v.ItemNum
				break
			}
		}

		if isExist != 1 {
			response.RetCode = msg.RE_INVALID_PARAM
			return
		}

		//! 发放奖励
		player.BagMoudle.AddAwardItem(req.ItemID, itemNum)
	} else {

		//! 发放奖励
		player.BagMoudle.AddAwardItems(awardItems)
	}

	//! 设置状态
	taskInfo.TaskStatus = Task_Received
	go activity.DB_UpdatePlayerSevenTask(taskData.TaskID, taskInfo.TaskCount, taskInfo.TaskStatus)

	for _, v := range awardItems {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	response.RetCode = msg.RE_SUCCESS
	response.ActivityID = req.ActivityID
}

//! 玩家请求购买限购商品
func Hand_BuySevenActivityLimit(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuySevenActivityLimitItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuySevenActivityLimit : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_BuySevenActivityLimitItem_Ack
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

	player.ActivityModule.CheckReset()
	var activity *TActivitySevenDay
	for i, v := range player.ActivityModule.SevenDay {
		if v.ActivityID == req.ActivityID {
			activity = &player.ActivityModule.SevenDay[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_GetSevenDay Error: Activity not exist ID: %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("Hand_GetSevenDay Error: Activity not open ID: %d", req.ActivityID)
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	awardType := G_GlobalVariables.GetActivityAwardType(activity.ActivityID)

	//! 获取物品信息
	itemInfo := gamedata.GetSevnActivityItemInfo(req.OpenDay, awardType)
	if itemInfo == nil {
		gamelog.Error("GetSevnActivityItemInfo error: can't not get item openday:%d, awardType:%d", req.OpenDay, awardType)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查限购物品是否已经购买
	if activity.BuyLst.IsExist(req.OpenDay) >= 0 {
		response.RetCode = msg.RE_BUY_LIMIT
		gamelog.Error("Hand_BuySevenActivityLimit error: aleady buy")
		return
	}

	limitInfo := G_GlobalVariables.GetSevenDayLimit(activity.ActivityID)

	//! 检查限购物品是否已经达到购买上限
	if itemInfo.Limit < limitInfo.LimitBuy[itemInfo.OpenDay-1] {
		response.RetCode = msg.RE_BUY_LIMIT
		gamelog.Error("Hand_BuySevenActivityLimit error: buy limit")
		return
	}

	//! 检查用户金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.MoneyID, itemInfo.MoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		gamelog.Error("Hand_BuySevenActivityLimit error: money is not enough")
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(itemInfo.MoneyID, itemInfo.MoneyNum)

	//! 给予物品
	player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum)

	//! 记录购买
	activity.BuyLst.Add(req.OpenDay)
	go activity.DB_AddPlayerSevenTaskMark(req.OpenDay)

	//! 购买数目加一
	G_GlobalVariables.AddSevenDayLimit(activity.ActivityID, req.OpenDay-1)
	response.RetCode = msg.RE_SUCCESS

	response.BuyTimes = limitInfo.LimitBuy[req.OpenDay-1]
	response.ActivityID = req.ActivityID
}

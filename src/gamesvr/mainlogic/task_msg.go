package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

//! 获取所有日常任务
func Hand_GetTaskData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetTaskData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetTaskData : Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetTaskData_Ack
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

	player.TaskMoudle.CheckReset()

	//! 获取数据
	for _, v := range player.TaskMoudle.TaskList {
		task := msg.MSG_TaskInfo{}
		task.ID = v.ID
		task.Status = v.Status
		task.Count = v.Count
		response.Tasks = append(response.Tasks, task)
	}
	response.RetCode = msg.RE_SUCCESS

	response.TaskScore = player.TaskMoudle.TaskScore

	for _, v := range player.TaskMoudle.ScoreAwardID {
		var scoreaward msg.MSG_ScoreAward
		scoreaward.ScoreAwardID = v
		if player.TaskMoudle.ScoreAwardStatus.IsExist(v) >= 0 {
			scoreaward.Status = true
		} else {
			scoreaward.Status = false
		}
		response.ScoreAwardLst = append(response.ScoreAwardLst, scoreaward)
	}

	for _, v := range player.TaskMoudle.AchieveList {
		node := msg.MSG_TaskInfo{}
		node.ID = v.ID
		node.Count = v.Count
		node.Status = v.Status
		response.List = append(response.List, node)
	}
}

//! 领取任务奖励
func Hand_RecvTaskAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RecvTaskAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RecvTaskAward : Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_RecvTaskAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.TaskMoudle.CheckReset()

	//! 检查任务完成情况
	ret, errcode := player.TaskMoudle.CheckPlayerTask(req.TaskID)
	if ret == false {
		response.RetCode = errcode
		gamelog.Error("Hand_RecvTaskAward error : Task not complete. playerID: %v  taskID: %v", req.PlayerID, req.TaskID)
		return
	}

	//! 发放奖励
	ret, itemLst := player.TaskMoudle.ReceiveTaskAward(req.TaskID)
	if ret == false {
		gamelog.Error("Hand_RecvTaskAward error : TaskAward receive failed. playerID: %v  taskID: %v", req.PlayerID, req.TaskID)
		return
	}

	//! 记入物品信息
	for _, v := range itemLst {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	//! 改变领取标记
	index := -1
	for i, v := range player.TaskMoudle.TaskList {
		if v.ID == req.TaskID {
			player.TaskMoudle.TaskList[i].Status = Task_Received
			index = i
			break
		}
	}

	player.TaskMoudle.DB_UpdateTask(index)

	//! 返回当前任务积分
	response.TaskScore = player.TaskMoudle.TaskScore
	response.RetCode = msg.RE_SUCCESS
}

//! 领取任务积分奖励
func Hand_RecvTaskScoreAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RecvTaskScoreAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RecvTaskScoreAward : Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_RecvTaskScoreAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.TaskMoudle.CheckReset()

	//! 检测参数
	if req.ScoreAwardID <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_RecvTaskScoreAward error: invalid ScoreAwardID: %d  PlayerID: %v", req.ScoreAwardID, player.playerid)
		return
	}

	//! 判断玩家积分领取资格
	ret, errcode := player.TaskMoudle.CheckTaskScore(req.ScoreAwardID)
	if ret == false {
		gamelog.Error("Hand_RecvTaskScoreAward error : Score not enough player : %d  score: %d  ask: %d",
			req.PlayerID, player.TaskMoudle.TaskScore, req.ScoreAwardID)
		response.RetCode = errcode
		return
	}

	//! 发放积分奖励
	ret, itemLst := player.TaskMoudle.ReceiveTaskScoreAward(req.ScoreAwardID)
	if ret == false {
		gamelog.Error("Hand_RecvTaskScoreAward error : TaskScoreAward receive failed. playerID: %v  ScoreID: %v", req.PlayerID, req.ScoreAwardID)
		return
	}

	//! 记入物品信息
	for _, v := range itemLst {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	//! 改变领取标记
	player.TaskMoudle.ScoreAwardStatus.Add(req.ScoreAwardID)

	//! 当前领取状态更新到数据库
	player.TaskMoudle.DB_UpdateTaskScoreAwardStatus()

	response.RetCode = msg.RE_SUCCESS

	for _, v := range player.TaskMoudle.ScoreAwardID {
		var scoreaward msg.MSG_ScoreAward
		scoreaward.ScoreAwardID = v
		if player.TaskMoudle.ScoreAwardStatus.IsExist(v) >= 0 {
			scoreaward.Status = true
		} else {
			scoreaward.Status = false
		}

		response.ScoreAwardLst = append(response.ScoreAwardLst, scoreaward)
	}
}

//! 请求领取成就奖励
func Hand_RecvAchievementAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_RecvAchievementAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RecvAchievementAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_RecvAchievementAward_Ack
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

	player.TaskMoudle.CheckReset()

	//! 检查参数
	if req.AchievementID <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_RecvAchievementAward error: invalid achievementID: %d  PlayerID: %v", req.AchievementID, player.playerid)
		return
	}

	//! 检查成就是否达成
	ret, errcode := player.TaskMoudle.CheckAchievement(req.AchievementID)
	if ret == false {
		gamelog.Error("Hand_RecvAchievementAward error: Achievement not complete. AchievemengtID: %d", req.AchievementID)
		response.RetCode = errcode
		return
	}

	//! 发放成就奖励
	ret = player.TaskMoudle.ReceiveAchievementAward(req.AchievementID)
	if ret == false {
		gamelog.Error("Hand_RecvAchievementAward error : TaskScoreAward receive failed. playerID: %v  AchievementID: %v", req.PlayerID, req.AchievementID)
		return
	}

	for i, _ := range player.TaskMoudle.AchieveList {
		if player.TaskMoudle.AchieveList[i].ID == req.AchievementID {
			//! 修改成就标记
			player.TaskMoudle.AchieveList[i].Status = Task_Received
			player.TaskMoudle.DB_UpdateAchieve(i)

			//! 增加成就完成列表
			player.TaskMoudle.AchieveIDs = append(player.TaskMoudle.AchieveIDs, player.TaskMoudle.AchieveList[i].ID)
			player.TaskMoudle.DB_AddAchieveID(player.TaskMoudle.AchieveList[i].ID)

		}
	}

	//! 查询替换成就
	newTask := player.TaskMoudle.UpdateNextAchievement(req.AchievementID)
	response.RetCode = msg.RE_SUCCESS
	response.NewAchieve.ID = newTask.ID
	response.NewAchieve.Count = newTask.Count
	response.NewAchieve.Status = newTask.Status
}

//! 玩家请求开服天数
func Hand_GetServerOpenDay(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetOpenServerDay_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetServerOpenDay : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetOpenServerDay_Ack
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

	player.TaskMoudle.CheckReset()

	response.OpenDay = GetOpenServerDay()
	response.RetCode = msg.RE_SUCCESS
}

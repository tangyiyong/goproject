package mainlogic

import (
	"encoding/json"
	"gamelog"
	"msg"
	"net/http"
)

//! 获取所有日常任务
func Hand_GetAllTask(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	var req msg.MSG_GetTasks_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAllTask : Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetTasks_Ack
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
		task := msg.TTaskInfo{}
		task.TaskID = v.TaskID
		task.TaskStatus = v.TaskStatus
		task.TaskCount = v.TaskCount
		response.Tasks = append(response.Tasks, task)
	}
	response.RetCode = msg.RE_SUCCESS
}

//! 领取任务奖励
func Hand_GetTaskAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTaskAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAllTask : Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetTaskAward_Ack
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
		gamelog.Error("Hand_GetTaskAward error : Task not complete. playerID: %v  taskID: %v", req.PlayerID, req.TaskID)
		return
	}

	//! 发放奖励
	ret, itemLst := player.TaskMoudle.ReceiveTaskAward(req.TaskID)
	if ret == false {
		gamelog.Error("Hand_GetTaskAward error : TaskAward receive failed. playerID: %v  taskID: %v", req.PlayerID, req.TaskID)
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
		if v.TaskID == req.TaskID {
			player.TaskMoudle.TaskList[i].TaskStatus = Task_Received
			index = i
			break
		}
	}

	player.TaskMoudle.DB_UpdatePlayerTask(req.TaskID, player.TaskMoudle.TaskList[index].TaskCount, Task_Received)

	//! 返回当前任务积分
	response.TaskScore = player.TaskMoudle.TaskScore
	response.RetCode = msg.RE_SUCCESS
}

//! 领取任务积分奖励
func Hand_GetTaskScoreAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTaskScoreAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetTaskScoreAward : Unmarshal fail. Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetTaskScoreAward_Ack
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
		gamelog.Error("Hand_GetTaskScoreAward error: invalid ScoreAwardID: %d  PlayerID: %v", req.ScoreAwardID, player.playerid)
		return
	}

	//! 判断玩家积分领取资格
	ret, errcode := player.TaskMoudle.CheckTaskScore(req.ScoreAwardID)
	if ret == false {
		gamelog.Error("Hand_GetTaskScoreAward error : Score not enough player : %d  score: %d  ask: %d",
			req.PlayerID, player.TaskMoudle.TaskScore, req.ScoreAwardID)
		response.RetCode = errcode
		return
	}

	//! 发放积分奖励
	ret, itemLst := player.TaskMoudle.ReceiveTaskScoreAward(req.ScoreAwardID)
	if ret == false {
		gamelog.Error("Hand_GetTaskScoreAward error : TaskScoreAward receive failed. playerID: %v  ScoreID: %v", req.PlayerID, req.ScoreAwardID)
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
	go player.TaskMoudle.DB_UpdatePlayerTaskScoreAwardStatus()

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

//! 请求日常任务积分信息
func Hand_GetTaskScoreInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetTaskScores_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetTaskScoreAward : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetTaskScores_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 获取玩家信息
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	player.TaskMoudle.CheckReset()

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
}

//! 请求所有成就
func Hand_GetAllAchievement(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetAchievementAll_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_GetAllAchievement : Unmarshal error!!!!")
		return
	}

	var response msg.MSG_GetAchievementAll_Ack
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

	for _, v := range player.TaskMoudle.AchievementList {
		node := msg.TAchievementInfo{}
		node.ID = v.ID
		node.TaskCount = v.TaskCount
		node.TaskStatus = v.TaskStatus
		response.List = append(response.List, node)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 请求领取成就奖励
func Hand_GetAchievementAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetAchievementAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetAchievementAward : Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetAchievementAward_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
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
		gamelog.Error("Hand_GetAchievementAward error: invalid achievementID: %d  PlayerID: %v", req.AchievementID, player.playerid)
		return
	}

	//! 检查成就是否达成
	ret, errcode := player.TaskMoudle.CheckAchievement(req.AchievementID)
	if ret == false {
		gamelog.Error("Hand_GetAchievementAward error: Achievement not complete. AchievemengtID: %d", req.AchievementID)
		response.RetCode = errcode
		return
	}

	//! 发放成就奖励
	ret = player.TaskMoudle.ReceiveAchievementAward(req.AchievementID)
	if ret == false {
		gamelog.Error("Hand_GetAchievementAward error : TaskScoreAward receive failed. playerID: %v  AchievementID: %v", req.PlayerID, req.AchievementID)
		return
	}

	for i, _ := range player.TaskMoudle.AchievementList {
		if player.TaskMoudle.AchievementList[i].ID == req.AchievementID {
			//! 修改成就标记
			player.TaskMoudle.AchievementList[i].TaskStatus = Task_Received
			player.TaskMoudle.DB_UpdatePlayerAchievement(player.TaskMoudle.AchievementList[i].ID,
				player.TaskMoudle.AchievementList[i].TaskCount,
				player.TaskMoudle.AchievementList[i].TaskStatus)

			//! 增加成就完成列表
			player.TaskMoudle.AchievedList = append(player.TaskMoudle.AchievedList, player.TaskMoudle.AchievementList[i].ID)
			player.TaskMoudle.DB_AddAchievementCompleteLst(player.TaskMoudle.AchievementList[i].ID)

		}
	}

	//! 查询替换成就
	newTask := player.TaskMoudle.UpdateNextAchievement(req.AchievementID)
	response.RetCode = msg.RE_SUCCESS
	response.NewAchieve.ID = newTask.ID
	response.NewAchieve.TaskCount = newTask.TaskCount
	response.NewAchieve.TaskStatus = newTask.TaskStatus
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
		gamelog.Info("Return: %s", b)
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

package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 获取限时任务信息
func Hand_GetLimitDailyTaskInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetLimitDailyTask_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetLimitDailyTaskInfo Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetLimitDailyTask_Ack
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

	//! 获取活动信息
	var activity *TActivityLimitDaily
	for i, v := range player.ActivityModule.LimitDaily {
		if v.ActivityID == activity.ActivityID {
			activity = &player.ActivityModule.LimitDaily[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_GetLimitDailyTaskInfo Error: Activity not exist ID: %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("Hand_GetLimitDailyTaskInfo Error: Activity not open ID: %d", req.ActivityID)
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	for _, v := range activity.TaskLst {
		var task msg.TLimitDailyTask
		task.TaskType = v.TaskType
		task.Count = v.Count
		task.Need = v.Need
		task.Status = v.Status
		response.TaskLst = append(response.TaskLst, task)
	}

	response.RetCode = msg.RE_SUCCESS

}

//! 获取限时日常任务奖励
func Hand_GetLimitDailyTaskAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetLimitDailyAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetLimitDailyTaskAward Error: Unmarshal fail, Error: %s", err.Error())
		return
	}

	var response msg.MSG_GetLimitDailyAward_Ack
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

	//! 获取活动信息
	var activity *TActivityLimitDaily
	activityIndex := 0
	for i, v := range player.ActivityModule.LimitDaily {
		if v.ActivityID == activity.ActivityID {
			activityIndex = i
			activity = &player.ActivityModule.LimitDaily[i]
			break
		}
	}

	if activity == nil {
		gamelog.Error("Hand_GetLimitDailyTaskAward Error: Activity not exist ID: %d", req.ActivityID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if G_GlobalVariables.IsActivityOpen(activity.ActivityID) == false {
		gamelog.Error("Hand_GetLimitDailyTaskInfo Error: Activity not open ID: %d", req.ActivityID)
		response.RetCode = msg.RE_ACTIVITY_NOT_OPEN
		return
	}

	if len(activity.TaskLst) <= req.Index || req.Index < 0 {
		gamelog.Error("Hand_GetLimitDailyTaskAward Error: Invalid param")
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	task := &activity.TaskLst[req.Index]
	if task.Status != 1 {
		gamelog.Error("Hand_GetLimitDailyTaskAward Error: Aleady receive or not done")
		response.RetCode = msg.RE_TASK_NOT_COMPLETE
		return
	}

	awardLst := gamedata.GetItemsFromAwardID(task.Award)
	if task.IsSelect != 0 && req.Select > len(awardLst) {
		gamelog.Error("Hand_GetLimitDailyTaskAward Error: Invalid select %d", req.Select)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if task.IsSelect != 0 {
		award := awardLst[req.Select-1]
		player.BagMoudle.AddAwardItem(award.ItemID, award.ItemNum)
		response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{award.ItemID, award.ItemNum})
	} else {
		player.BagMoudle.AddAwardItems(awardLst)
		for _, v := range awardLst {
			response.AwardItem = append(response.AwardItem, msg.MSG_ItemData{v.ItemID, v.ItemNum})
		}
	}

	//! 更新任务状态
	task.Status = 2
	activity.DB_UpdateTaskStatus(activityIndex, req.Index)

	response.RetCode = msg.RE_SUCCESS
}

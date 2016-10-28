package mainlogic

import (
	"appconfig"
	"bytes"
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"time"
)

func SelectScoreTarget(player *TPlayer, value int) bool {
	if player.ScoreMoudle.Score < value {
		//	return false
	}

	return true
}

//获取积分赛主界面信息
func Hand_GetScoreData(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreData_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreTarget Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreData_Ack
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

	if gamedata.IsFuncOpen(gamedata.FUNC_SCORE_SYSTEM, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.Score = player.ScoreMoudle.Score
	response.FightTime = player.ScoreMoudle.FightTime
	response.WinTime = int(player.ScoreMoudle.SeriesWin & 0x0000FFFF)
	response.IsRecv = int((player.ScoreMoudle.SeriesWin & 0xFFFF0000) >> 16)
	response.BuyTime = player.ScoreMoudle.BuyTime
	response.Targets = player.ScoreMoudle.GetScoreTargets()
	response.Rank = G_ScoreRaceRanker.GetRankIndex(player.playerid, player.ScoreMoudle.Score)
	response.ItemLst = player.ScoreMoudle.BuyRecord
	response.RetCode = msg.RE_SUCCESS
}

func Hand_GetScoreBattleCheck(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreBattleCheck_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreBattleCheck Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreBattleCheck_Ack
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

	if req.TargetIndex < 0 || req.TargetIndex >= 3 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GetScoreBattleCheck Error Invalid TargetIndex:%d", req.TargetIndex)
		return
	}

	var GetFightTargetReq msg.MSG_GetFightTarget_Req
	GetFightTargetReq.PlayerID = player.ScoreMoudle.ScoreEnemy[req.TargetIndex].PlayerID
	GetFightTargetReq.SvrID = player.ScoreMoudle.ScoreEnemy[req.TargetIndex].SvrID

	if GetFightTargetReq.PlayerID == 0 || GetFightTargetReq.SvrID == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GetScoreBattleCheck Error Invalid PlayerID:%d, and SvrID:%d", GetFightTargetReq.PlayerID, GetFightTargetReq.SvrID)
		return
	}

	buffer, _ = json.Marshal(GetFightTargetReq)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(appconfig.CrossGetFightTarget, "text/HTML", bytes.NewReader(buffer))
	if err != nil || httpret == nil {
		gamelog.Error("Hand_GetScoreBattleCheck failed, err : %s !!!!", err.Error())
		return
	}

	buffer = make([]byte, httpret.ContentLength)
	httpret.Body.Read(buffer)
	httpret.Body.Close()
	var GetFightTargetAck msg.MSG_GetFightTarget_Ack
	err = json.Unmarshal(buffer, &GetFightTargetAck)
	if err != nil {
		gamelog.Error("Hand_GetScoreBattleCheck  Unmarshal fail, Error: %s", err.Error())
		return
	}

	response.PlayerData = GetFightTargetAck.PlayerData
	response.RetCode = GetFightTargetAck.RetCode
}

//玩家提交战斗结果
func Hand_SetScoreBattleResult(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_SetScoreBattleResult_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetScoreBattleResult Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SetScoreBattleResult_Ack
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

	//如果打赢了
	if req.WinBattle == 1 {
		player.ScoreMoudle.Score += gamedata.OneTimeFightScore
		player.ScoreMoudle.SeriesWin += 1
	} else {
		player.ScoreMoudle.Score -= gamedata.OneTimeFightScore
		player.ScoreMoudle.SeriesWin &= 0xFFFF0000
	}

	if player.ScoreMoudle.Score < 0 {
		player.ScoreMoudle.Score = 0
	}

	player.ScoreMoudle.FightTime += 1
	player.ScoreMoudle.DB_SaveScoreAndFightTime()

	player.RoleMoudle.CostAction(1, 1)
	response.Targets = player.ScoreMoudle.GetScoreTargets()
	response.RetCode = msg.RE_SUCCESS
	response.Rank = G_ScoreRaceRanker.SetRankItem(player.playerid, player.ScoreMoudle.Score)
	player.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_SCORE_RANK, response.Rank)
	return
}

func (score *TScoreMoudle) GetScoreTargets() []msg.MSG_Target {
	var ScoreTargetReq msg.MSG_CrossQueryScoreTarget_Req
	ScoreTargetReq.PlayerID = score.PlayerID
	b, _ := json.Marshal(ScoreTargetReq)
	http.DefaultClient.Timeout = 3 * time.Second
	httpret, err := http.Post(appconfig.CrossQueryScoreTarget, "text/HTML", bytes.NewReader(b))
	if err != nil || httpret == nil {
		gamelog.Error("GetScoreTargets failed, err : %s !!!!", err.Error())
		return nil
	}

	buffer := make([]byte, httpret.ContentLength)
	httpret.Body.Read(buffer)
	httpret.Body.Close()

	var ScoreTargetAck msg.MSG_CrossQueryScoreTarget_Ack
	err = json.Unmarshal(buffer, &ScoreTargetAck)
	if err != nil {
		gamelog.Error("GetScoreTargets  Unmarshal fail, Error: %s", err.Error())
		return nil
	}

	for i := 0; i < len(ScoreTargetAck.TargetList); i++ {
		score.ScoreEnemy[i].FightValue = ScoreTargetAck.TargetList[i].FightValue
		score.ScoreEnemy[i].HeroID = ScoreTargetAck.TargetList[i].HeroID
		score.ScoreEnemy[i].Level = ScoreTargetAck.TargetList[i].Level
		score.ScoreEnemy[i].Name = ScoreTargetAck.TargetList[i].Name
		score.ScoreEnemy[i].PlayerID = ScoreTargetAck.TargetList[i].PlayerID
		score.ScoreEnemy[i].SvrName = ScoreTargetAck.TargetList[i].SvrName
		score.ScoreEnemy[i].SvrID = ScoreTargetAck.TargetList[i].SvrID
		score.ScoreEnemy[i].Quality = ScoreTargetAck.TargetList[i].Quality
	}

	return ScoreTargetAck.TargetList[0:3]
}

//请求积分赛排行榜信息
func Hand_GetScoreRank(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreRank_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreRank Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreRank_Ack
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

	response.ScoreRankList = []msg.MSG_ScoreRankInfo{}
	response.MyRank = -1
	for i := 0; i < len(G_ScoreRaceRanker.List); i++ {
		if G_ScoreRaceRanker.List[i].RankID <= 0 {
			break
		}

		if len(response.ScoreRankList) >= G_FightRanker.ShowNum {
			break
		}
		pSimpleInfo := G_SimpleMgr.GetSimpleInfoByID(G_ScoreRaceRanker.List[i].RankID)
		if pSimpleInfo != nil {
			var info msg.MSG_ScoreRankInfo
			info.FightValue = pSimpleInfo.FightValue
			info.Name = pSimpleInfo.Name
			info.Quality = pSimpleInfo.Quality
			info.HeroID = pSimpleInfo.HeroID
			info.Score = G_ScoreRaceRanker.List[i].RankValue
			response.ScoreRankList = append(response.ScoreRankList, info)
		}

		if G_ScoreRaceRanker.List[i].RankID == req.PlayerID {
			response.MyRank = i + 1
		}
	}

	if response.MyRank < 0 {
		response.MyRank = G_ScoreRaceRanker.GetRankIndex(player.playerid, player.ScoreMoudle.Score)
	}

	response.RetCode = msg.RE_SUCCESS
	response.MyScore = player.ScoreMoudle.Score

	return
}

//请求积分赛战斗次数奖励
func Hand_RecvScoreTimeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RecvScoreTimeAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RcvScoreTimeAward Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RecvScoreTimeAward_Ack
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

	pTimeAwardInfo := gamedata.GetScoreTimeAward(req.TimeAwardID)
	if pTimeAwardInfo == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_RcvScoreTimeAward Invalid Time Award ID: %d", req.TimeAwardID)
		return
	}
	for _, v := range player.ScoreMoudle.RecvAward {
		if v == req.TimeAwardID {
			response.RetCode = msg.RE_ALREADY_RECEIVED
			gamelog.Error("Hand_RcvScoreTimeAward Already received award: %d", req.TimeAwardID)
			return
		}
	}

	if player.ScoreMoudle.FightTime < pTimeAwardInfo.Times {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_RcvScoreTimeAward Not Enough Time: %d", player.ScoreMoudle.FightTime)
		return
	}

	dropItem := gamedata.GetItemsFromAwardID(pTimeAwardInfo.AwardID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	player.BagMoudle.AddAwardItems(dropItem)
	player.ScoreMoudle.RecvAward = append(player.ScoreMoudle.RecvAward, req.TimeAwardID)
	player.ScoreMoudle.DB_UpdateRecvAward()
	response.RetCode = msg.RE_SUCCESS
	return
}

//请求积分赛连胜奖励
func Hand_RecvContinueWinAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_RecvScoreContinueAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_RcvScoreTimeAward Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_RecvScoreContinueAward_Ack
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

	if int(player.ScoreMoudle.SeriesWin&0x0000FFFF) < gamedata.ScoreSeriesWinTimes {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_RcvScoreTimeAward Not Enough Win Times:%d", player.ScoreMoudle.SeriesWin)
		return
	}

	dropItem := gamedata.GetItemsFromAwardID(gamedata.ScoreSeriesWinAwardID)
	for _, v := range dropItem {
		var item msg.MSG_ItemData
		item.ID = v.ItemID
		item.Num = v.ItemNum
		response.ItemLst = append(response.ItemLst, item)
	}

	player.BagMoudle.AddAwardItems(dropItem)
	response.RetCode = msg.RE_SUCCESS
	return
}

//请求积分赛战斗次数奖励信息
func Hand_GetScoreTimeAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_GetScoreTimeAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetScoreTimeAward Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetScoreTimeAward_Ack
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

	response.Awards = player.ScoreMoudle.RecvAward
	response.FightTime = player.ScoreMoudle.FightTime
	response.RetCode = msg.RE_SUCCESS
	return
}

//! 玩家请求购买积分商店道具
//! 消息: /buy_score_store_item
func Hand_BuyScoreStoreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BuyScoreStoreItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyScoreStoreItem Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyScoreStoreItem_Ack
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

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	//! 获取购买物品信息
	itemInfo := gamedata.GetScoreStoreItem(int(req.StoreItemID))
	if itemInfo == nil {
		gamelog.Error("Hand_BuyScoreStoreItem Error: GetScoreStoreItem nil ID: %d ", req.StoreItemID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 判断道具是否足够
	if itemInfo.CostMoneyID > 20 {
		if player.BagMoudle.IsItemEnough(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum) == false {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough Item")
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}
	} else { //! 判断货币是否足够
		if player.RoleMoudle.CheckMoneyEnough(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum) == false {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough money")
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}
	}

	isExist := false
	for i, v := range player.ScoreMoudle.BuyRecord {
		if int32(v.ID) == req.StoreItemID {
			isExist = true
			if v.Times+req.BuyNum > itemInfo.MaxBuyTime && itemInfo.MaxBuyTime != 0 {
				gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
				response.RetCode = msg.RE_NOT_ENOUGH_TIMES
				return
			}
			player.ScoreMoudle.BuyRecord[i].Times += req.BuyNum
			player.ScoreMoudle.DB_UpdateStoreItemBuyTimes(i, player.ScoreMoudle.BuyRecord[i].Times)
		}
	}

	if isExist == false {
		//! 首次购买
		if req.BuyNum > itemInfo.MaxBuyTime && itemInfo.MaxBuyTime != 0 {
			gamelog.Error("Hand_BuyScoreStoreItem Error: Not enough buy times")
			response.RetCode = msg.RE_NOT_ENOUGH_TIMES
			return
		}

		var itemData msg.MSG_BuyData
		itemData.ID = req.StoreItemID
		itemData.Times = req.BuyNum
		player.ScoreMoudle.BuyRecord = append(player.ScoreMoudle.BuyRecord, itemData)
		player.ScoreMoudle.DB_AddStoreItemBuyInfo(itemData)
	}

	//! 扣除货币
	if itemInfo.CostMoneyID > 20 {
		player.BagMoudle.RemoveNormalItem(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum)
	} else {
		player.RoleMoudle.CostMoney(itemInfo.CostMoneyID, itemInfo.CostMoneyNum*req.BuyNum)
	}

	//! 发放物品
	player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum*req.BuyNum)

	response.RetCode = msg.RE_SUCCESS
}

//购买积分赛战斗次数
func Hand_BuyScoreFightTime(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)
	var req msg.MSG_BuyScoreTime_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_BuyScoreFightTime Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuyScoreTime_Ack
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

	if response.RetCode = player.BeginMsgProcess(); response.RetCode != msg.RE_UNKNOWN_ERR {
		return
	}

	defer player.FinishMsgProcess()

	maxTime := gamedata.GetFuncVipValue(gamedata.FUNC_SCORE_FIGHT_TIME, player.GetVipLevel())
	if player.ScoreMoudle.BuyTime >= maxTime {
		response.RetCode = msg.RE_NOT_ENOUGH_ITEM
		gamelog.Error("Hand_BuyScoreFightTime Not Enough Time")
		return
	}

	cost := gamedata.GetFuncTimeCost(gamedata.FUNC_SCORE_FIGHT_TIME, player.ScoreMoudle.FightTime)
	pCopyInfo := gamedata.GetCopyBaseInfo(gamedata.ScoreCopyID)
	if pCopyInfo == nil {
		response.RetCode = msg.RE_INVALID_COPY_ID
		gamelog.Error("Hand_BuyScoreFightTime Invalid CopyID :%d", gamedata.ScoreCopyID)
		return
	}

	response.RetCode = msg.RE_SUCCESS
	player.RoleMoudle.CostMoney(gamedata.ScoreBuyTimeMoneyID, cost)
	player.RoleMoudle.AddAction(pCopyInfo.ActionType, 1)
	player.ScoreMoudle.BuyTime += 1
	player.ScoreMoudle.DB_SaveBuyFightTime()
	response.BuyTime = player.ScoreMoudle.BuyTime
	response.ActionValue, response.ActionTime = player.RoleMoudle.GetActionData(pCopyInfo.ActionType)

	return
}

//玩家请求积分赛战报信息
//消息:/get_score_report_req
type MSG_GetScoreReport_Req struct {
	PlayerID   int32  //玩家ID
	SessionKey string //Sessionkey
}

type MSG_GetScoreReport_Ack struct {
	RetCode int            //返回码
	Reports []TScoreReport //战报表
}

//! 玩家请求积分赛战报信息
func Hand_GetScoreBattleReport(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req MSG_GetScoreReport_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_ReceiveAllMails unmarshal fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response MSG_GetScoreReport_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
	}()

	//! 通用检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	response.Reports = player.MailMoudle.Reports
	response.RetCode = msg.RE_SUCCESS
	player.MailMoudle.Reports = []TScoreReport{}
	player.MailMoudle.DB_ClearAllReports()

}

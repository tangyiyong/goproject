package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
	"utility"
)

//! 获取三国无双星数信息
func Hand_GetSangokuMusou_StarInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouStarInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusou_StarInfo Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouStarInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	player.SangokuMusouModule.CheckReset()

	//! 获取信息
	response.RetCode = msg.RE_SUCCESS
	if player.SangokuMusouModule.IsEnd == true {
		response.IsEnd = 1
	} else {
		response.IsEnd = 0
	}

	response.CanUseStar = player.SangokuMusouModule.CanUseStar
	response.CurStar = player.SangokuMusouModule.CurStar
	response.HistoryStar = player.SangokuMusouModule.HistoryStar
}

//! 获取三国无双闯关信息
func Hand_GetSangokuMuSou_CopyInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouCopyInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMuSou_CopyInfo Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouCopyInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	player.SangokuMusouModule.CheckReset()

	//! 获取信息
	response.RetCode = msg.RE_SUCCESS
	response.BattleTimes = player.SangokuMusouModule.BattleTimes
	response.PassCopyID = player.SangokuMusouModule.PassCopyID

	if player.SangokuMusouModule.PassCopyID == 0 {
		return
	}

	chapter := gamedata.GetSangokuMusouChapterInfo(player.SangokuMusouModule.PassCopyID + 1)

	if chapter != nil {
		if chapter.ChapterID == 1 { //! 第一章不存在上一章节奖励领取

			response.IsRecvAward = 1
			response.IsSelectBuff = 1

		} else {
			if player.SangokuMusouModule.ChapterAwardMark.IsExist(chapter.ChapterID-1) >= 0 {
				response.IsRecvAward = 1
			} else {
				response.IsRecvAward = 0
			}

			if player.SangokuMusouModule.ChapterBuffMark.IsExist(chapter.ChapterID-1) >= 0 {
				response.IsSelectBuff = 1
			} else {
				response.IsSelectBuff = 0
			}
		}

	} else {
		response.IsRecvAward = 1
		response.IsSelectBuff = 1
	}

	chapter = gamedata.GetSangokuMusouChapterInfo(player.SangokuMusouModule.PassCopyID + 1)
	if chapter != nil {
		chapterInfo := gamedata.GetSGWSChapterCopyLst(chapter.ChapterID)
		for i, v := range chapterInfo {
			for _, n := range player.SangokuMusouModule.CopyInfoLst {
				if v == n.CopyID {
					response.CopyLst[i] = n.StarNum
				}
			}
		}
	}

}

//! 获取三国无双精英挑战闯关信息
func Hand_GetSangokuMusou_EliteCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouEliteCopyInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusou_EliteCopy Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouEliteCopyInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取信息
	response.RetCode = msg.RE_SUCCESS
	response.HistoryCopyID = player.SangokuMusouModule.HistoryCopyID
	response.PassEliteCopyID = player.SangokuMusouModule.PassEliteCopyID
	response.BattleTimes = player.SangokuMusouModule.EliteBattleTimes
}

//! 通关三国无双
func Hand_PassSangokuMusou_Copy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//MD5消息验证
	if false == utility.MsgDataCheck(buffer, G_XorCode) {
		//存在作弊的可能
		gamelog.Error("Hand_PassSangokuMusou_Copy : Message Data Check Error!!!!")
		return
	}
	var req msg.MSG_PassSangokuMusouCopy_Req
	if json.Unmarshal(buffer[:len(buffer)-16], &req) != nil {
		gamelog.Error("Hand_PassSangokuMusou_Copy : Unmarshal error!!!!")
		return
	}

	//! 创建回复
	var response msg.MSG_PassSangokuMusouCopy_Ack
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

	if req.CopyID <= 0 {
		gamelog.Error("Hand_PassSangokuMusou_Copy copyID is invaild. id: %d", req.CopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检查称号过期
	player.TitleModule.CheckTitleDeadLine()

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	if player.SangokuMusouModule.IsEnd == true {
		//! 挑战已经结束
		gamelog.Error("Hand_PassSangokuMusou_Copy Error: Challenge aleady ended.")
		response.RetCode = msg.RE_CHALLENGE_ALEADY_END
		return
	}

	//! 检测关卡有效性
	chapter := gamedata.GetSangokuMusouChapterInfo(req.CopyID)
	if chapter == nil {
		gamelog.Error("Hand_PassSangokuMusou_Copy copyID is invaild. id: %d", req.CopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	if req.CopyID < player.SangokuMusouModule.PassCopyID {
		//! 重置前不得重复挑战
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	nextCopyID := gamedata.GetSangokuMusouNextCopy(player.SangokuMusouModule.PassCopyID)
	if nextCopyID == 0 || req.CopyID < nextCopyID {
		//! 已经全部通关
		gamelog.Error("GetSangokuMusouNextCopy error: invaild copyID %d", nextCopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 通关普通关卡
	isVictory := true
	if req.IsVictory == 0 {
		isVictory = false
	}
	dropItemLst := player.SangokuMusouModule.PassCopy(req.CopyID, req.StarNum, isVictory)

	nextCopyID = gamedata.GetSangokuMusouNextCopy(player.SangokuMusouModule.PassCopyID)
	if nextCopyID == 0 {
		//! 全部通关, 将状态置为结束
		player.SangokuMusouModule.IsEnd = true
		player.SangokuMusouModule.DB_UpdateIsEndMark()
	}

	response.RetCode = msg.RE_SUCCESS
	response.DropItem = dropItemLst
	return

}

//! 通关三国无双精英挑战
func Hand_PassSangokuMusou_EliteCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! MD5消息验证
	if false == utility.MsgDataCheck(buffer, G_XorCode) {
		//存在作弊的可能
		gamelog.Error("Hand_PassSangokuMusou_EliteCopy : Message Data Check Error!!!!")
		return
	}
	var req msg.MSG_PassSangokuMusouEliteCopy_Req
	if json.Unmarshal(buffer[:len(buffer)-16], &req) != nil {
		gamelog.Error("Hand_PassSangokuMusou_EliteCopy : Unmarshal error!!!!")
		return
	}

	//! 创建回复
	var response msg.MSG_PassSangokuMusouEliteCopy_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取关卡信息
	copyInfo := gamedata.GetSangokuMusouEliteCopyInfo(req.CopyID)
	if copyInfo == nil {
		gamelog.Error("Hand_PassSangokuMusou_Copy copyID is invaild. id: %d", req.CopyID)
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 检测挑战次数
	if player.SangokuMusouModule.EliteBattleTimes <= 0 {
		//! 挑战次数不足
		response.RetCode = msg.RE_NOT_ENOUGH_TIMES
		return
	}

	//! 检测历史通关关卡是否满足条件
	if player.SangokuMusouModule.HistoryCopyID < copyInfo.NeedPassCopy {
		//! 需要挑战过前置关卡
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 通关精英挑战
	ret := player.SangokuMusouModule.PassEliteCopy(req.CopyID)
	if ret == true {
		response.IsFirstVictory = 1
	}

	response.RetCode = msg.RE_SUCCESS
	return
}

//! 请求扫荡三国无双章节
func Hand_SangokuMusou_Sweep(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_PassSangokuMusouCopy_sweep_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SangokuMusou_Sweep Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_PassSangokuMusouCopy_sweep_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 开始扫荡
	response.DropItem = player.SangokuMusouModule.SweepChapter(req.Chapter)
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求章节奖励
func Hand_GetSangokuMusou_ChapterAward(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouChapterAward_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusou_ChapterAward Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouChapterAward_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取玩家当前所在章节
	isChapterEnd, chapter := gamedata.SangokuMusou_IsChapterEnd(player.SangokuMusouModule.PassCopyID)
	if isChapterEnd == false {
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 检测领取标记
	if player.SangokuMusouModule.ChapterAwardMark.IsExist(chapter) >= 0 {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 获取章节信息
	beginCopy := gamedata.GetSGWSChapterBeginCopyID(chapter)
	endCopy := gamedata.GetSGWSChapterEndCopyID(chapter)
	starNum := 0
	for _, v := range player.SangokuMusouModule.CopyInfoLst {
		if v.CopyID >= beginCopy && v.CopyID <= endCopy {
			starNum += v.StarNum
		}
	}

	starNum = starNum - (starNum % 3)

	//! 获取章节奖励信息
	chapterAward := gamedata.GetSangokuMusouChapterAwardInfo(chapter, starNum)
	awardLst := gamedata.GetItemsFromAwardID(chapterAward.Award)
	player.BagMoudle.AddAwardItems(awardLst)

	response.AwardLst = []msg.MSG_ItemData{}
	for _, v := range awardLst {
		response.AwardLst = append(response.AwardLst, msg.MSG_ItemData{v.ItemID, v.ItemNum})
	}

	//! 增加奖励领取记录
	player.SangokuMusouModule.ChapterAwardMark.Add(chapter)
	player.SangokuMusouModule.DB_UpdateChapterAwardMark()

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求随机章节属性奖励
func Hand_GetSangokuMusou_ChapterAttr(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouAttrInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusou_ChapterAttr Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouAttrInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检测当前章节是否通关
	isChapterEnd, _ := gamedata.SangokuMusou_IsChapterEnd(player.SangokuMusouModule.PassCopyID)
	if isChapterEnd == false {
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 清空之前的属性加成选项
	if len(player.SangokuMusouModule.AwardAttrLst) > 0 {
		player.SangokuMusouModule.AwardAttrLst = []TSangokuMusouAttrData{}
	}

	//! 随机属性奖励
	attrLst := gamedata.RandSangokuMusouAttrMarkup()
	for i, v := range attrLst {
		var attrInfo msg.MSG_SangokuMusou_Attr1
		attrInfo.ID = i
		attrInfo.CostStar = v.CostStar
		attrInfo.Value = v.Value
		attrInfo.AttrID = v.AttrID
		response.AttrLst = append(response.AttrLst, attrInfo)

		info := TSangokuMusouAttrData{
			ID:       attrInfo.ID,
			CostStar: attrInfo.CostStar,
			AttrID:   attrInfo.AttrID,
			Value:    attrInfo.Value}
		player.SangokuMusouModule.AwardAttrLst = append(player.SangokuMusouModule.AwardAttrLst, info)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求选择属性奖励
func Hand_SetSangokuMusou_ChapterAttr(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SetSangokuMusouAttrInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetSangokuMusou_ChapterAttr Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SetSangokuMusouAttrInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	isChapterEnd, chapter := gamedata.SangokuMusou_IsChapterEnd(player.SangokuMusouModule.PassCopyID)
	if isChapterEnd == false {
		gamelog.Error("Hand_SetSangokuMusou_ChapterAttr error: player not in end copy. CopyID: %d", player.SangokuMusouModule.PassCopyID)
		response.RetCode = msg.RE_NEED_PASS_PRE_COPY
		return
	}

	//! 检测ID是否在奖励属性列表中
	isExist := false
	for _, v := range player.SangokuMusouModule.AwardAttrLst {
		if v.ID == req.ID {
			isExist = true

			//! 检测星数是否足够
			if player.SangokuMusouModule.CanUseStar < v.CostStar {
				response.RetCode = msg.RE_NOT_ENOUGH_STAR
				return
			}

			//! 扣除星数
			player.SangokuMusouModule.CanUseStar -= v.CostStar
			player.SangokuMusouModule.DB_UpdateCanUseStar()

			//! 添加属性
			isFind := false
			for i, _ := range player.SangokuMusouModule.AttrMarkupLst {
				if player.SangokuMusouModule.AttrMarkupLst[i].AttrID == v.AttrID {
					player.SangokuMusouModule.AttrMarkupLst[i].Value += v.Value
					isFind = true
					player.SangokuMusouModule.DB_UpdateAttr(i, player.SangokuMusouModule.AttrMarkupLst[i].Value)
					break
				}
			}

			if isFind == false {
				player.SangokuMusouModule.AttrMarkupLst = append(player.SangokuMusouModule.AttrMarkupLst, TSangokuMusouAttrData2{v.AttrID, v.Value})
				player.SangokuMusouModule.DB_AddAttr(TSangokuMusouAttrData2{v.AttrID, v.Value})
			}
		}
	}

	if isExist == false {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	} else {
		//! 记录玩家领奖
		player.SangokuMusouModule.ChapterBuffMark.Add(chapter)
		player.SangokuMusouModule.DB_UpdateChapterBuffMark()
	}

	//! 清空之前的属性加成选项緩存
	if len(player.SangokuMusouModule.AwardAttrLst) > 0 {
		player.SangokuMusouModule.AwardAttrLst = []TSangokuMusouAttrData{}
	}

	//! 返回所有Buff属性
	for _, v := range player.SangokuMusouModule.AttrMarkupLst {
		var attr msg.MSG_SangokuMusou_Attr2
		attr.AttrID = v.AttrID
		attr.Value = v.Value
		response.AttrLst = append(response.AttrLst, attr)
	}

	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询所有属性加成
func Hand_GetSangokuMusou_Attr(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouAllAttrInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetSangokuMusou_ChapterAttr Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouAllAttrInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.RetCode = msg.RE_SUCCESS

	if len(player.SangokuMusouModule.AttrMarkupLst) <= 0 {
		response.AttrLst = []msg.MSG_SangokuMusou_Attr2{}
		return
	} else {

		for _, v := range player.SangokuMusouModule.AttrMarkupLst {
			attrInfo := msg.MSG_SangokuMusou_Attr2{}
			attrInfo.AttrID = v.AttrID
			attrInfo.Value = v.Value
			response.AttrLst = append(response.AttrLst, attrInfo)
		}
	}
}

//! 玩家请求随机无双秘藏
func Hand_GetSangokuMusou_Treasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouTreasureInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetSangokuMusou_ChapterAttr Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouTreasureInfo_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 必须要挑战结束才允许
	if player.SangokuMusouModule.IsEnd == false {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	response.TreasureID = player.SangokuMusouModule.GetMusouTreasure()
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买无双秘藏
func Hand_BuySangokuMusou_Treasure(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuySangokuMusouTreasure_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetSangokuMusou_ChapterAttr Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuySangokuMusouTreasure_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 检查重复购买
	if player.SangokuMusouModule.IsBuyTreasure == true {
		response.RetCode = msg.RE_ALREADY_RECEIVED
		return
	}

	//! 检查秘藏是否刷新
	if player.SangokuMusouModule.TreasureID == 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	//! 获取秘藏信息
	treasure := gamedata.GetMusouTreasure(player.SangokuMusouModule.TreasureID)

	//! 检查金钱
	if player.RoleMoudle.CheckMoneyEnough(treasure.CostMoneyType, treasure.CostMoneyNum) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(treasure.CostMoneyType, treasure.CostMoneyNum)

	//! 给予奖励
	player.BagMoudle.AddAwardItem(treasure.ItemID, treasure.ItemNum)

	player.SangokuMusouModule.IsBuyTreasure = true
	player.SangokuMusouModule.TreasureID = 0
	player.SangokuMusouModule.DB_UpdateTreasure()
}

//! 玩家请求重置关卡挑战
func Hand_SangokuMusou_ResetCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SangokuMusou_Reset_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SangokuMusou_ResetCopy Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SangokuMusou_Reset_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 刷新次数
	player.SangokuMusouModule.CheckReset()

	//! 检测重制次数是否足够
	resetTimes := gamedata.GetFuncVipValue(gamedata.FUNC_SANGUOWUSHUANG_RESET, player.GetVipLevel())
	if player.SangokuMusouModule.BattleTimes >= resetTimes {
		response.RetCode = msg.RE_NOT_ENOUGH_REFRESH_TIMES
		return
	}

	costMoney := gamedata.GetFuncTimeCost(gamedata.FUNC_SANGUOWUSHUANG_RESET, player.SangokuMusouModule.BattleTimes+1)
	if costMoney < 0 {
		gamelog.Error("GetResetCostFail. FuncID: %d", gamedata.FUNC_SANGUOWUSHUANG_RESET)
		return
	}

	//! 检查金钱是否足够
	if costMoney != 0 {
		if player.RoleMoudle.CheckMoneyEnough(gamedata.SangokuMusouResetCopyMoneyID, costMoney) == false {
			response.RetCode = msg.RE_NOT_ENOUGH_MONEY
			return
		}

		//! 扣除金钱
		player.RoleMoudle.CostMoney(gamedata.SangokuMusouResetCopyMoneyID, costMoney)
		response.MoneyNum = costMoney
	}

	//! 重置挑战
	player.SangokuMusouModule.ResetCopy()

	response.ResetTimes = player.SangokuMusouModule.BattleTimes
	response.MoneyID = gamedata.SangokuMusouResetCopyMoneyID
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求查询可增加精英挑战次数
func Hand_GetSangokuMusou_AddEliteCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusou_Add_BattleTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusou_AddEliteCopy Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusou_Add_BattleTimes_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	player.SangokuMusouModule.CheckReset()

	resetTimes := gamedata.GetFuncVipValue(gamedata.FUNC_BUY_SANGUOWUSHUANG_ELITE_TIMES, player.GetVipLevel())

	response.ResetTimes = resetTimes - player.SangokuMusouModule.AddEliteBattleTimes
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求增加精英挑战次数
func Hand_SangokuMusou_AddEliteCopy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SangokuMusou_Add_BattleTimes_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SangokuMusou_AddEliteCopy Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SangokuMusou_Add_BattleTimes_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	player.SangokuMusouModule.CheckReset()

	//! 检测购买精英挑战次数
	resetTimes := gamedata.GetFuncVipValue(gamedata.FUNC_BUY_SANGUOWUSHUANG_ELITE_TIMES, player.GetVipLevel())

	if player.SangokuMusouModule.AddEliteBattleTimes >= resetTimes {
		response.RetCode = msg.RE_NOT_ENOUGH_REFRESH_TIMES
		return
	}

	//! 计算扣除金钱
	costMoney := player.SangokuMusouModule.AddEliteBattleTimes + 1
	costMoney = costMoney * 30

	//! 检查金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(gamedata.SangokuMusouAddEliteCopyMoneyID, costMoney) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(gamedata.SangokuMusouAddEliteCopyMoneyID, costMoney)

	//! 增加次数
	player.SangokuMusouModule.EliteBattleTimes += 1
	player.SangokuMusouModule.AddEliteBattleTimes += 1

	response.RetCode = msg.RE_SUCCESS
	player.SangokuMusouModule.DB_UpdateEliteBattleTimes()
}

//! 玩家请求获取已购买的商品列表
func Hand_GetSangokuMusouStore_AleadyBuy(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouStoreAleadyBuy_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusouStore_AleadyBuy Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouStoreAleadyBuy_Ack
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

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.ItemLst = player.SangokuMusouModule.BuyRecord
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求购买三国无双商店商品
func Hand_GetSangokuMusou_StoreItem(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_BuySangokuMusouStoreItem_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSangokuMusou_StoreItem Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_BuySangokuMusouStoreItem_Ack
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

	//! 参数检查
	if req.Num <= 0 {
		response.RetCode = msg.RE_INVALID_PARAM
		gamelog.Error("Hand_GetSangokuMusou_StoreItem invalid item num. Num: %d  PlayerID: %v", req.Num, player.playerid)
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOWUSHUANG, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取商品信息
	itemInfo := gamedata.GetSangokumusouStoreInfo(int(req.ID))
	if itemInfo == nil {
		gamelog.Error("Hand_GetSangokuMusou_StoreItem get item info fail. ID: %d", req.ID)
		return
	}

	//! 检查等级
	if player.GetLevel() < itemInfo.NeedLevel {
		response.RetCode = msg.RE_NOT_ENOUGH_HERO_LEVEL
		return
	}

	//! 检查金钱是否足够
	if player.RoleMoudle.CheckMoneyEnough(itemInfo.CostMoneyType, itemInfo.CostMoneyNum*req.Num) == false {
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	if itemInfo.CostItemType != 0 {
		//! 检查道具是否足够
		if player.BagMoudle.IsItemEnough(itemInfo.CostItemType, itemInfo.CostItemNum*req.Num) == false {
			response.RetCode = msg.RE_NOT_ENOUGH_ITEM
			return
		}
	}

	if itemInfo.BuyTimes == 0 {
		//! 可无限次购买商品
		//! 扣除金钱,给予物品
		player.RoleMoudle.CostMoney(itemInfo.CostMoneyType, itemInfo.CostMoneyNum*req.Num)

		if itemInfo.CostItemType != 0 {
			player.BagMoudle.RemoveNormalItem(itemInfo.CostItemType, itemInfo.CostItemNum*req.Num)
		}

		player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum*req.Num)
	} else {
		//! 限次购买商品

		//! 检查当前购买次数
		isExist := player.SangokuMusouModule.IsShoppingInfoExist(req.ID)
		if isExist == false {
			//! 没有购买该物品记录则创建新的记录
			var info msg.MSG_BuyData
			info.ID = int32(itemInfo.ID)
			info.Times = 0
			player.SangokuMusouModule.BuyRecord = append(player.SangokuMusouModule.BuyRecord, info)
			player.SangokuMusouModule.DB_AddStoreItemBuyInfo(info)
		}

		for i, v := range player.SangokuMusouModule.BuyRecord {
			if req.ID == int(v.ID) {
				isExist = true
				if v.Times+req.Num > itemInfo.BuyTimes {
					response.RetCode = msg.RE_NOT_ENOUGH_TIMES
					return
				} else {
					//! 扣除金钱,增加次数
					player.RoleMoudle.CostMoney(itemInfo.CostMoneyType, itemInfo.CostMoneyNum*req.Num)

					if itemInfo.CostItemType != 0 {
						player.BagMoudle.RemoveNormalItem(itemInfo.CostItemType, itemInfo.CostItemNum*req.Num)
					}

					player.BagMoudle.AddAwardItem(itemInfo.ItemID, itemInfo.ItemNum*req.Num)
					player.SangokuMusouModule.BuyRecord[i].Times += req.Num
					response.Times = player.SangokuMusouModule.BuyRecord[i].Times
					player.SangokuMusouModule.DB_UpdateStoreItemBuyTimes(i, player.SangokuMusouModule.BuyRecord[i].Times)
				}
			}
		}

	}

	response.RetCode = msg.RE_SUCCESS

}

func Hand_GetSanguowsStatus(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSangokuMusouStatus_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSanguowsStatus Unmarshal fail, Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSangokuMusouStatus_Ack
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

	//! 星数
	player.SangokuMusouModule.CheckReset()

	//! 获取信息
	response.RetCode = msg.RE_SUCCESS
	if player.SangokuMusouModule.IsEnd == true {
		response.IsEnd = 1
	} else {
		response.IsEnd = 0
	}

	response.CanUseStar = player.SangokuMusouModule.CanUseStar
	response.CurStar = player.SangokuMusouModule.CurStar
	response.HistoryStar = player.SangokuMusouModule.HistoryStar

	//! 获取信息
	response.RetCode = msg.RE_SUCCESS
	response.BattleTimes = player.SangokuMusouModule.BattleTimes
	response.PassCopyID = player.SangokuMusouModule.PassCopyID

	//! 属性
	if len(player.SangokuMusouModule.AttrMarkupLst) <= 0 {
		response.AttrLst = []msg.MSG_SangokuMusou_Attr2{}
	} else {

		for _, v := range player.SangokuMusouModule.AttrMarkupLst {
			attrInfo := msg.MSG_SangokuMusou_Attr2{}
			attrInfo.AttrID = v.AttrID
			attrInfo.Value = v.Value
			response.AttrLst = append(response.AttrLst, attrInfo)
		}
	}

	//! 商店
	response.ItemLst = player.SangokuMusouModule.BuyRecord
	response.IsBuyTreasure = player.SangokuMusouModule.IsBuyTreasure

	if player.SangokuMusouModule.PassCopyID == 0 {
		return
	}

	//! 闯关信息
	chapter := gamedata.GetSangokuMusouChapterInfo(player.SangokuMusouModule.PassCopyID + 1)

	if chapter != nil {
		if chapter.ChapterID == 1 { //! 第一章不存在上一章节奖励领取

			response.IsRecvAward = 1
			response.IsSelectBuff = 1

		} else {
			if player.SangokuMusouModule.ChapterAwardMark.IsExist(chapter.ChapterID-1) >= 0 {
				response.IsRecvAward = 1
			} else {
				response.IsRecvAward = 0
			}

			if player.SangokuMusouModule.ChapterBuffMark.IsExist(chapter.ChapterID-1) >= 0 {
				response.IsSelectBuff = 1
			} else {
				response.IsSelectBuff = 0
			}
		}

	} else {
		response.IsRecvAward = 1
		response.IsSelectBuff = 1
	}

	chapter = gamedata.GetSangokuMusouChapterInfo(player.SangokuMusouModule.PassCopyID)
	if chapter != nil {
		chapterInfo := gamedata.GetSGWSChapterCopyLst(chapter.ChapterID)
		for i, v := range chapterInfo {
			for _, n := range player.SangokuMusouModule.CopyInfoLst {
				if v == n.CopyID {
					response.CopyLst[i] = n.StarNum
				}
			}
		}
	}

}

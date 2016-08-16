package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

func Hand_SanGuoZhiInfo(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSanGuoZhiInfo_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SanGuoZhiInfo Unmarshal is fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSanGuoZhiInfo_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOZHI, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	//! 获取已命星的星
	response.CurOpenID = player.SanGuoZhiModule.CurStarID
	response.RetCode = msg.RE_SUCCESS
}

//! 玩家请求命星
func Hand_SetSanGuoZhi(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_SetSanGuoZhi_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_SetSanGuoZhi Unmarshal is fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_SetSanGuoZhi_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOZHI, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	if gamedata.IsStarEnd(player.SanGuoZhiModule.CurStarID) == true {
		response.RetCode = msg.RE_MAX_STAR
		return
	}

	//! 检测命星材料是否足够
	ok, errcode := player.SanGuoZhiModule.CheckItemEnough(player.SanGuoZhiModule.CurStarID + 1)
	if ok == false {
		response.RetCode = errcode
		return
	}

	//! 扣除材料
	info := gamedata.GetSanGuoZhiInfo(player.SanGuoZhiModule.CurStarID + 1)
	player.BagMoudle.RemoveNormalItem(info.CostType, info.CostNum)

	//! 开始升星
	player.SanGuoZhiModule.CurStarID += 1
	player.SanGuoZhiModule.SaveSanGuoZhiStar()

	if info.Type == gamedata.Sanguo_Add_Attr {
		//! 全队增加指定属性
		player.HeroMoudle.AddExtraProperty(info.AttrID, info.Value, false, 0)
		player.HeroMoudle.DB_SaveExtraProperty()
	} else if info.Type == gamedata.Sanguo_Give_Item {
		//! 给予道具
		player.BagMoudle.AddAwardItem(info.AttrID, info.Value)
		response.AwardItem = msg.MSG_ItemData{info.AttrID, info.Value}
	} else if info.Type == gamedata.Sanguo_Main_Hero_Up {
		//! 提升主角品质
		player.HeroMoudle.ChangeMainQuality(info.Value)
	}

	response.FightValue = player.CalcFightValue()
	response.Quality = player.HeroMoudle.CurHeros[0].Quality
	response.RetCode = msg.RE_SUCCESS
}

//! 查询星宿增加属性
func Hand_GetSanGuoStarAddAttribute(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("message: %s", r.URL.String())

	//! 接收信息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_GetSanGuoZhi_Attribute_Req
	err := json.Unmarshal(buffer, &req)
	if err != nil {
		gamelog.Error("Hand_GetSanGuoStarAddAttribute Unmarshal is fail. Error: %s", err.Error())
		return
	}

	//! 创建回复
	var response msg.MSG_GetSanGuoZhi_Attribute_Ack
	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	//! 常规检测
	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 检测功能是否开启
	if gamedata.IsFuncOpen(gamedata.FUNC_SANGUOZHI, player.GetLevel(), player.GetVipLevel()) == false {
		response.RetCode = msg.RE_FUNC_NOT_OPEN
		return
	}

	response.RetCode = msg.RE_SUCCESS
}

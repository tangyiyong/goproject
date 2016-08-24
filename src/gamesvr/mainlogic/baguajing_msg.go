package mainlogic

import (
	"encoding/json"
	"gamelog"
	"gamesvr/gamedata"
	"msg"
	"net/http"
)

//! 玩家请求使用八卦镜
func Hand_UseBaGuaJing(w http.ResponseWriter, r *http.Request) {
	gamelog.Info("msg: %s", r.URL.String())

	//! 接受消息
	buffer := make([]byte, r.ContentLength)
	r.Body.Read(buffer)

	//! 解析消息
	var req msg.MSG_UseBaguajing_Req
	if json.Unmarshal(buffer, &req) != nil {
		gamelog.Error("Hand_UseBaGuaJing Error: invalid json: %s", buffer)
		return
	}

	//! 定义返回
	var response msg.MSG_UseBaguajing_Ack

	response.RetCode = msg.RE_UNKNOWN_ERR
	defer func() {
		b, _ := json.Marshal(&response)
		w.Write(b)
		gamelog.Info("Return: %s", b)
	}()

	var player *TPlayer = nil
	player, response.RetCode = GetPlayerAndCheck(req.PlayerID, req.SessionKey, r.URL.String())
	if player == nil {
		return
	}

	//! 获取英雄信息
	hero := player.BagMoudle.GetBagHeroByPos(req.BagPos)
	if hero == nil {
		response.RetCode = msg.RE_INVALID_PARAM
		return
	}

	heroInfo := gamedata.GetHeroInfo(hero.ID)

	//! 计算需要金钱
	bgjConfig := gamedata.GetBaGuaJingConfigData()

	//! 获取目标英雄信息
	targetHeroInfo := gamedata.GetHeroInfo(req.HeroID)

	//! 判断英雄等级
	if heroInfo.Camp != targetHeroInfo.Camp &&
		player.GetLevel() < bgjConfig.CrossFactionNeedLevel {
		gamelog.Error("Hand_UseBaGuaJing Error: Not enough level:%v", bgjConfig.CrossFactionNeedLevel)
		response.RetCode = msg.RE_NOT_ENOUGH_LEVEL
		return
	}

	//! 判断英雄品质
	if heroInfo.Quality == bgjConfig.LimitQuality &&
		player.GetLevel() < bgjConfig.QualityHeroExchangeNeedLevel {
		response.RetCode = msg.RE_NOT_ENOUGH_LEVEL
		gamelog.Error("Hand_UseBaGuaJing Error: Not enough level:%v", bgjConfig.QualityHeroExchangeNeedLevel)
		return
	}

	moneyID1, moneyNum1 := bgjConfig.BaseMoneyID1, bgjConfig.BaseMoneyNum1
	moneyID2, moneyNum2 := bgjConfig.BaseMoneyID2, bgjConfig.BaseMoneyNum2

	//! 计算突破觉醒花费
	costValue := 0
	for i := 0; i <= hero.WakeLevel; i++ {
		costValue += gamedata.GetWakeLevelItem(i).NeedHeroNum
	}

	for i := int8(0); i <= hero.BreakLevel; i++ {
		costValue += gamedata.GetHeroBreakInfo(i).HeroNum
	}

	if costValue != 0 {
		moneyNum1 *= costValue
		moneyNum2 *= costValue
	}

	//! 判断玩家货币是否足够
	if player.RoleMoudle.CheckMoneyEnough(moneyID1, moneyNum1) == false ||
		player.RoleMoudle.CheckMoneyEnough(moneyID2, moneyNum2) == false {
		gamelog.Error("Hand_UseBaGuaJing Error: Not enough money")
		response.RetCode = msg.RE_NOT_ENOUGH_MONEY
		return
	}

	//! 扣除金钱
	player.RoleMoudle.CostMoney(moneyID1, moneyNum1)
	player.RoleMoudle.CostMoney(moneyID2, moneyNum2)

	//! 改变英雄
	hero.ID = targetHeroInfo.HeroID
	go player.BagMoudle.DB_UpdateHeroID(req.BagPos, hero.ID)

	//! 返回成功
	response.RetCode = msg.RE_SUCCESS
}

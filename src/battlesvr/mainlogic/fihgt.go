package mainlogic

import (
	"battlesvr/gamedata"
	"gamelog"
	"msg"
	"utility"
)

func (self *TBattleRoom) Hand_SkillState(pdata []byte) {
	gamelog.Info("message: MSG_SKILL_STATE")

	var req msg.MSG_Skill_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_SkillState : Message Reader Error!!!!")
		return
	}

	if req.SkillEvents_Cnt <= 0 && req.AttackEvents_Cnt <= 0 {
		gamelog.Error("Hand_SkillState Error: SkillEvents_Cnt :%d, AttackEvents_Cnt:%d!!!!", req.SkillEvents_Cnt, req.AttackEvents_Cnt)
		return
	}

	SendMessageToRoom(req.PlayerID, self.RoomID, msg.MSG_SKILL_STATE, &req)

	pAttackBatObj := self.GetBattleByPID(req.PlayerID)
	if pAttackBatObj == nil {
		gamelog.Error("Hand_SkillState : Invalid playerid:%d!!!!", req.PlayerID)
		return
	}

	var KillEventReq msg.MSG_KillEvent_Req
	KillEventReq.PlayerID = req.PlayerID

	var bNewSkill bool = false
	var HeroStateNty msg.MSG_HeroState_Nty
	for i := 0; i < len(req.SkillEvents); i++ {
		if req.SkillEvents[i].S_Skill_ID == pAttackBatObj.SkillState[3].ID {
			bNewSkill = true
		}

		if len(req.SkillEvents[i].TargetIDs) <= 0 {
			gamelog.Error("Hand_SkillState Error: len of targetids is 0!!!!")
			continue
		}

		pAttackHeroObj := self.GetHeroObject(req.SkillEvents[i].S_ID)
		if pAttackHeroObj.CurHp <= 0 {
			gamelog.Error("Hand_SkillState Error: pHeroAttacker CurHp is 0!!!!")
			continue
		}
		//检查内容
		// 1.技能是否存在
		// 2.技能CD是否可以施放
		// 3.技能是否可以打中指定的目标
		// 4.

		//是否可以放技能，CD, 技能是否是英雄所有
		//if false == CanSkill(pHeroObject, req.SkillEvents[i].S_Skill_ID) {
		//	pTcpConn.Close()
		//	gamelog.Error("player %d ")
		//	return

		//}

		//pHeroObject.Skill[ID]-GetCurTime > CD
		//

		//是否可以打中目标
		//for _, id := range req.SkillEvents[i].TargetIDs {
		//	pTargetObject := pRoom.GetHeroObject(id)
		//距离是否合适

		//}

		for j := 0; j < len(req.SkillEvents[i].TargetIDs); j++ {
			pDefHeroObj := self.GetHeroObject(req.SkillEvents[i].TargetIDs[j])
			if pDefHeroObj == nil {
				gamelog.Error("Hand_SkillState Error: pDefHeroObj is nil i:%d,j :%d", i, j)
				continue
			}

			bkill := HeroFight(pAttackHeroObj, req.SkillEvents[i].S_Skill_ID, pDefHeroObj)
			HeroStateNty.Heros = append(HeroStateNty.Heros, msg.MSG_HeroItem{ObjectID: pDefHeroObj.ObjectID, CurHp: pDefHeroObj.CurHp})

			if bkill == false {
				continue
			}

			gamelog.Error("Hand_SkillState Error: attackerid:%d, defenderid:%d, defender_hp:%d", pAttackHeroObj.ObjectID, pDefHeroObj.ObjectID, pDefHeroObj.CurHp)
			pAttackBatObj.SeriesKill = pAttackBatObj.SeriesKill + 1

			KillEventReq.Kill += 1
			KillEventReq.SeriesKill = pAttackBatObj.SeriesKill

			pDefBatObj := self.GetBattleByOID(pDefHeroObj.ObjectID)
			if pDefBatObj == nil {
				gamelog.Error("Hand_SkillState Error: cannot get the def batobj!!")
			}

			if pDefBatObj.IsAllDie() {
				KillEventReq.Destroy += 1

				//向客户端发复活通知
				var AllDieNty msg.MSG_HeroAllDie_Nty
				AllDieNty.NtyCode = 0
				if pDefBatObj.ReviveTime[0] > 4 || pDefBatObj.ReviveTime[1] > 5 {
					AllDieNty.NtyCode = 1
				}
				pDieConn := GetConnByID(pDefBatObj.PlayerID)
				if pDieConn != nil {
					var writer msg.PacketWriter
					writer.BeginWrite(msg.MSG_ALL_DIE_NTY, 0)
					AllDieNty.Write(&writer)
					writer.EndWrite()
					pDieConn.WriteMsgData(writer.GetDataPtr())
				}

			}

		}
	}

	//向服务器发送击杀事件
	if KillEventReq.Kill > 0 {
		SendMessageToGameSvr(msg.MSG_KILL_EVENT_REQ, int16(self.RoomID), &KillEventReq)
	}

	HeroStateNty.Heros_Cnt = int32(len(HeroStateNty.Heros))
	if HeroStateNty.Heros_Cnt > 0 {
		SendMessageToRoom(0, self.RoomID, msg.MSG_HERO_STATE, &HeroStateNty)
	}

	pTcpConn := GetConnByID(req.PlayerID)
	if pTcpConn == nil {
		gamelog.Error("Hand_SkillState Error: Invalid PlayerID:%d", req.PlayerID)
		return
	}

	if bNewSkill {
		pAttackBatObj.SkillState[3].ID = gamedata.RandSkill()
		gamelog.Error("Hand_SkillState Error: New Skill ID :%d", pAttackBatObj.SkillState[3].ID)
		var msgNewSkill msg.MSG_NewSkill_Nty
		msgNewSkill.NewSkillID = pAttackBatObj.SkillState[3].ID
		var writer msg.PacketWriter
		writer.BeginWrite(msg.MSG_NEW_SKILL_NTY, 0)
		msgNewSkill.Write(&writer)
		writer.EndWrite()
		pTcpConn.WriteMsgData(writer.GetDataPtr())
	}
	return
}

//1	生命值
//2	物理攻击
//3	物理防御
//4	魔法攻击
//5	魔法防御
//6	伤害减免
//7	伤害加成
//8	闪避率
//9	命中率
//10 暴击率
//11 抗暴率

func HeroFight(pAttacker *THeroObj, skillid int32, pDefender *THeroObj) (bkill bool) {
	if pAttacker == nil || pDefender == nil {
		gamelog.Error("HeroFight Error: pAttacker == nil || pDefender == nil ")
		return
	}

	bkill = false
	value := int32(utility.Rand() % 1000)
	//先判断是否命中
	if value > (800+pAttacker.CurProperty[8]-pDefender.CurProperty[7]) && value > 500 {
		return
	}

	//判断是否爆击
	value = int32(utility.Rand() % 1000)
	bSuperHit := false
	if value < (pAttacker.CurProperty[9]-pAttacker.CurProperty[10]) || value < 10 {
		bSuperHit = true
	} else {
		bSuperHit = false
	}

	var pSkillInfo = new(gamedata.ST_SkillInfo)
	pSkillInfo.Hurts = make([]gamedata.ST_Hurts, 1, 1)
	pSkillInfo.Hurts[0].Percent = 100
	pSkillInfo.Hurts[0].Fixed = 0

	//最终伤害加成
	finaladd := pAttacker.CurProperty[6] - pDefender.CurProperty[5] + 1000

	//伤害随机
	fightrand := int32(900 + utility.Rand()%200)
	//hurt := (pSkillInfo.Hurts[0].Percent*(pAttacker.CurProperty[pAttacker.AttackPID-1]-pDefender.CurProperty[pAttacker.AttackPID]) + pSkillInfo.Hurts[0].Fixed)
	hurt := pAttacker.CurProperty[pAttacker.AttackPID-1] - pDefender.CurProperty[pAttacker.AttackPID]
	if hurt <= 0 {
		hurt = 1
	} else {
		hurt = hurt * fightrand / 1000
		hurt = hurt * finaladd / 1000
		if bSuperHit {
			hurt = hurt * 15 / 10
		}
	}

	gamelog.Error("HeroFight Info hurt:%d", hurt)

	pDefender.CurHp -= hurt
	if pDefender.CurHp <= 0 {
		pDefender.CurHp = 0
		bkill = true
	}

	bkill = false
	return
}

package mainlogic

import (
	"battlesvr/gamedata"
	"gamelog"
	"msg"
	"tcpserver"
	"utility"
)

func Hand_SkillState(pTcpConn *tcpserver.TCPConn, pdata []byte) {
	gamelog.Info("message: MSG_SKILL_STATE")
	playerid := pTcpConn.Data.(*TBattleData).PlayerID
	roomid := pTcpConn.Data.(*TBattleData).RoomID

	var req msg.MSG_Skill_Req
	if req.Read(new(msg.PacketReader).BeginRead(pdata, 0)) == false {
		gamelog.Error("Hand_SkillState : Message Reader Error!!!!")
		return
	}

	pRoom := G_RoomMgr.GetRoomByID(pTcpConn.Data.(*TBattleData).RoomID)
	if pRoom == nil {
		gamelog.Error("Hand_SkillState : Invalid RoomID:%d!!!!", pTcpConn.Data.(*TBattleData).RoomID)
		return
	}

	if req.SkillEvents_Cnt <= 0 && req.AttackEvents_Cnt <= 0 {
		gamelog.Error("Hand_SkillState Error: SkillEvents_Cnt :%d, AttackEvents_Cnt:%d!!!!", req.SkillEvents_Cnt, req.AttackEvents_Cnt)
		return
	}

	SendMessageToRoom(playerid, roomid, msg.MSG_SKILL_STATE, &req)

	pSrcBatObj := pRoom.GetBattleByPID(playerid)
	if pSrcBatObj == nil {
		gamelog.Error("Hand_SkillState : Invalid playerid:%d!!!!", playerid)
		return
	}

	bNewSkill := false
	var ackHeroState msg.MSG_HeroState_Nty
	for i := 0; i < len(req.SkillEvents); i++ {
		if req.SkillEvents[i].S_Skill_ID == int32(pSrcBatObj.SkillState[3].ID) {
			pSrcBatObj.SkillState[3].ID = gamedata.RandSkill()
			gamelog.Error("Hand_SkillState Error: New Skill ID :%d", pSrcBatObj.SkillState[3].ID)
			bNewSkill = true
		}

		if len(req.SkillEvents[i].TargetIDs) <= 0 {
			gamelog.Error("Hand_SkillState Error: len of targetids is 0!!!!")
			continue
		}

		pHeroAttacker := pRoom.GetHeroObject(req.SkillEvents[i].S_ID)
		if pHeroAttacker.CurHp <= 0 {
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

		//pHeroObject.Skill[ID]-time.Now().Unix() > CD
		//

		//是否可以打中目标
		//for _, id := range req.SkillEvents[i].TargetIDs {
		//	pTargetObject := pRoom.GetHeroObject(id)
		//距离是否合适

		//}
		var killreq msg.MSG_KillEvent_Req
		killreq.Killer = req.SkillEvents[i].S_ID
		for j := 0; j < len(req.SkillEvents[i].TargetIDs); j++ {
			pHeroDefender := pRoom.GetHeroObject(req.SkillEvents[i].TargetIDs[j])
			if pHeroDefender == nil {
				gamelog.Error("Hand_SkillState Error: pHeroDefender is nil i:%d,j :%d", i, j)
				continue
			}

			if pHeroDefender.CurHp <= 0 {
				gamelog.Error("Hand_SkillState Error: pHeroDefender CurHp is 0!!")
				ackHeroState.Heros = append(ackHeroState.Heros, msg.MSG_HeroItem{ObjectID: pHeroDefender.ObjectID, CurHp: pHeroDefender.CurHp})
				continue
			}

			bkill := HeroFight(pHeroAttacker, req.SkillEvents[i].S_Skill_ID, pHeroDefender)
			if bkill == true {
				gamelog.Error("Hand_SkillState Error: attackerid:%d, defenderid:%d, defender_hp:%d", pHeroAttacker.ObjectID, pHeroDefender.ObjectID, pHeroDefender.CurHp)
				pSrcBatObj.SeriesKill = pSrcBatObj.SeriesKill + 1

				//如果击杀需要做以下几件事:
				//1. 向游戏服发送击杀事件, 游戏服返回的数据，要发还给玩家

				killreq.Kill += 1
				killreq.SeriesKill = pSrcBatObj.SeriesKill

				pDefBatObj := pRoom.GetBattleByOID(pHeroDefender.ObjectID)
				if pDefBatObj == nil {
					gamelog.Error("Hand_SkillState Error: cannot get the def batobj!!")
				}

				if pDefBatObj.IsAllDie() {
					killreq.Destroy += 1
				}
			}

			//向服务器发送击杀事件
			if killreq.Kill > 0 {
				SendMessageToGameSvr(msg.MSG_KILL_EVENT_REQ, &killreq)
			}

			ackHeroState.Heros = append(ackHeroState.Heros, msg.MSG_HeroItem{ObjectID: pHeroDefender.ObjectID, CurHp: pHeroDefender.CurHp})
		}
	}

	ackHeroState.Heros_Cnt = int32(len(ackHeroState.Heros))
	if ackHeroState.Heros_Cnt > 0 {
		SendMessageToRoom(0, roomid, msg.MSG_HERO_STATE, &ackHeroState)
	}

	if bNewSkill {
		var msgNewSkill msg.MSG_NewSkill_Nty
		msgNewSkill.NewSkillID = pSrcBatObj.SkillState[3].ID
		var writer msg.PacketWriter
		writer.BeginWrite(msg.MSG_NEW_SKILL_NTY)
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

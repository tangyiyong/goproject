#include "stdafx.h>"
class MSG_HeroObj
{
	INT32 HeroID;		//英雄ID
	INT32 ObjectID;		//英雄实例ID
	INT32 CurHp;		//英雄血量
	FLOAT Position[5];		//x,y,z,v,d, x, y,z,速度，方向
	void Read(PacketReader *pReader)
	{
		HeroID = pReader->ReadInt32();
		ObjectID = pReader->ReadInt32();
		CurHp = pReader->ReadInt32();
		for(int i = 0; i < 5; i++)
		{
			Position[i] = pReader->ReadFloat();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(HeroID);
		pWriter->WriteInt32(ObjectID);
		pWriter->WriteInt32(CurHp);
		for(int i = 0; i < 5; i++)
		{
			pWriter->WriteFloat(Position[i]);
		}
		return ;
	}
};

class MSG_BattleObj
{
	INT8 BatCamp;
	MSG_HeroObj Heros[6];
	void Read(PacketReader *pReader)
	{
		BatCamp = pReader->ReadInt8();
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt8(BatCamp);
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

//进入阵营战消息(Client)
class MSG_EnterRoom_Req
{
	INT32 PlayerID;
	INT32 EnterCode;		//进入码
	INT32 MsgNo;
	void Read(PacketReader *pReader)
	{
		PlayerID = pReader->ReadInt32();
		EnterCode = pReader->ReadInt32();
		MsgNo = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(EnterCode);
		pWriter->WriteInt32(MsgNo);
		return ;
	}
};

class MSG_EnterRoom_Ack
{
	INT8 BatCamp;
	INT32 CurRank;		//今日排名
	INT32 KillNum;		//今日击杀
	INT32 KillHonor;		//今日杀人荣誉
	INT32 LeftTimes;		//剩余搬动次数
	INT32 MoveEndTime;		//搬运结束时间
	INT32 BeginMsgNo;		//起始消息编号
	INT32 SkillID[4];		//四个技能ID
	MSG_HeroObj Heros[6];		//六个英雄
	void Read(PacketReader *pReader)
	{
		BatCamp = pReader->ReadInt8();
		CurRank = pReader->ReadInt32();
		KillNum = pReader->ReadInt32();
		KillHonor = pReader->ReadInt32();
		LeftTimes = pReader->ReadInt32();
		MoveEndTime = pReader->ReadInt32();
		BeginMsgNo = pReader->ReadInt32();
		for(int i = 0; i < 4; i++)
		{
			SkillID[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt8(BatCamp);
		pWriter->WriteInt32(CurRank);
		pWriter->WriteInt32(KillNum);
		pWriter->WriteInt32(KillHonor);
		pWriter->WriteInt32(LeftTimes);
		pWriter->WriteInt32(MoveEndTime);
		pWriter->WriteInt32(BeginMsgNo);
		for(int i = 0; i < 4; i++)
		{
			pWriter->WriteInt32(SkillID[i]);
		}
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

class MSG_EnterRoom_Notify
{
	INT32 BatObjs_Cnt;
	MSG_BattleObj BatObjs[1];
	void Read(PacketReader *pReader)
	{
		BatObjs_Cnt = pReader->ReadInt32();
		BatObjs = new MSG_BattleObj[BatObjs_Cnt];
		for(int i = 0; i < BatObjs_Cnt; i++)
		{
			BatObjs[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(BatObjs_Cnt);
		for(int i = 0; i < BatObjs_Cnt; i++)
		{
			BatObjs[i].Write(pWriter);
		}
		return ;
	}
};

//(Client)
class MSG_LeaveRoom_Req
{
	INT32 MsgNo;
	INT32 PlayerID;
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		return ;
	}
};

class MSG_LeaveRoom_Ack
{
	INT32 PlayerID;
	void Read(PacketReader *pReader)
	{
		PlayerID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(PlayerID);
		return ;
	}
};

class MSG_LeaveRoom_Notify
{
	INT32 ObjectIDs[6];
	void Read(PacketReader *pReader)
	{
		for(int i = 0; i < 6; i++)
		{
			ObjectIDs[i] = pReader->ReadInt32();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		for(int i = 0; i < 6; i++)
		{
			pWriter->WriteInt32(ObjectIDs[i]);
		}
		return ;
	}
};

class MSG_Skill_Item
{
	INT32 S_ID;
	INT32 S_Skill_ID;
	INT32 TargetIDs_Cnt;
	INT32 TargetIDs[1];
	void Read(PacketReader *pReader)
	{
		S_ID = pReader->ReadInt32();
		S_Skill_ID = pReader->ReadInt32();
		TargetIDs_Cnt = pReader->ReadInt32();
		TargetIDs = new INT32[TargetIDs_Cnt];
		for(int i = 0; i < TargetIDs_Cnt; i++)
		{
			TargetIDs[i] = pReader->ReadInt32();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(S_ID);
		pWriter->WriteInt32(S_Skill_ID);
		pWriter->WriteInt32(TargetIDs_Cnt);
		for(int i = 0; i < TargetIDs_Cnt; i++)
		{
			pWriter->WriteInt32(TargetIDs[i]);
		}
		return ;
	}
};

//(Client)
class MSG_Skill_Req
{
	INT32 MsgNo;
	INT32 PlayerID;
	INT32 SkillEvents_Cnt;
	MSG_Skill_Item SkillEvents[1];
	INT32 AttackEvents_Cnt;
	MSG_Skill_Item AttackEvents[1];
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		SkillEvents_Cnt = pReader->ReadInt32();
		SkillEvents = new MSG_Skill_Item[SkillEvents_Cnt];
		for(int i = 0; i < SkillEvents_Cnt; i++)
		{
			SkillEvents[i].Read(pReader);
		}
		AttackEvents_Cnt = pReader->ReadInt32();
		AttackEvents = new MSG_Skill_Item[AttackEvents_Cnt];
		for(int i = 0; i < AttackEvents_Cnt; i++)
		{
			AttackEvents[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(SkillEvents_Cnt);
		for(int i = 0; i < SkillEvents_Cnt; i++)
		{
			SkillEvents[i].Write(pWriter);
		}
		pWriter->WriteInt32(AttackEvents_Cnt);
		for(int i = 0; i < AttackEvents_Cnt; i++)
		{
			AttackEvents[i].Write(pWriter);
		}
		return ;
	}
};

class MSG_Move_Item
{
	INT32 S_ID;
	FLOAT Position[5];
	void Read(PacketReader *pReader)
	{
		S_ID = pReader->ReadInt32();
		for(int i = 0; i < 5; i++)
		{
			Position[i] = pReader->ReadFloat();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(S_ID);
		for(int i = 0; i < 5; i++)
		{
			pWriter->WriteFloat(Position[i]);
		}
		return ;
	}
};

//(Client)
class MSG_Move_Req
{
	INT32 MsgNo;		//消息编号
	INT32 PlayerID;
	INT32 MoveEvents_Cnt;
	MSG_Move_Item MoveEvents[1];
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		MoveEvents_Cnt = pReader->ReadInt32();
		MoveEvents = new MSG_Move_Item[MoveEvents_Cnt];
		for(int i = 0; i < MoveEvents_Cnt; i++)
		{
			MoveEvents[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(MoveEvents_Cnt);
		for(int i = 0; i < MoveEvents_Cnt; i++)
		{
			MoveEvents[i].Write(pWriter);
		}
		return ;
	}
};

class MSG_HeroItem
{
	INT32 ObjectID;
	INT32 CurHp;
	void Read(PacketReader *pReader)
	{
		ObjectID = pReader->ReadInt32();
		CurHp = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(ObjectID);
		pWriter->WriteInt32(CurHp);
		return ;
	}
};

class MSG_HeroState_Nty
{
	INT32 Heros_Cnt;
	MSG_HeroItem Heros[1];
	void Read(PacketReader *pReader)
	{
		Heros_Cnt = pReader->ReadInt32();
		Heros = new MSG_HeroItem[Heros_Cnt];
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(Heros_Cnt);
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

//玩家查询当前的水晶品质//(Client)
class MSG_PlayerQuery_Req
{
	INT32 MsgNo;		//消息编号
	INT32 PlayerID;
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		return ;
	}
};

//玩家查询当前的水晶品质回复
class MSG_PlayerQuery_Ack
{
	INT32 RetCode;		//返回码
	INT32 PlayerID;		//玩家的ID
	INT32 Quality;		//水晶品质
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		Quality = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(Quality);
		return ;
	}
};

//(Client)
class MSG_StartCarry_Req
{
	INT32 MsgNo;		//消息编号
	INT32 PlayerID;
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		return ;
	}
};

class MSG_StartCarry_Ack
{
	INT32 RetCode;		//返回码
	INT32 PlayerID;		//玩家的ID
	INT32 EndTime;		//搬运截止时间
	INT32 LeftTimes;		//剩余搬动次数
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		EndTime = pReader->ReadInt32();
		LeftTimes = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(EndTime);
		pWriter->WriteInt32(LeftTimes);
		return ;
	}
};

//(Client)
class MSG_FinishCarry_Req
{
	INT32 MsgNo;		//消息编号
	INT32 PlayerID;
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		return ;
	}
};

class MSG_FinishCarry_Ack
{
	INT32 RetCode;		//返回码
	INT32 PlayerID;		//玩家的ID
	INT32 MoneyID[2];		//货币ID
	INT32 MoneyNum[2];		//货币数量
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		for(int i = 0; i < 2; i++)
		{
			MoneyID[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 2; i++)
		{
			MoneyNum[i] = pReader->ReadInt32();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt32(PlayerID);
		for(int i = 0; i < 2; i++)
		{
			pWriter->WriteInt32(MoneyID[i]);
		}
		for(int i = 0; i < 2; i++)
		{
			pWriter->WriteInt32(MoneyNum[i]);
		}
		return ;
	}
};

//(Client)
class MSG_PlayerChange_Req
{
	INT32 MsgNo;		//消息编号
	INT32 PlayerID;
	INT32 HighQuality;		//直接选择最高品质
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		HighQuality = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(HighQuality);
		return ;
	}
};

class MSG_PlayerChange_Ack
{
	INT32 RetCode;		//返回码
	INT32 PlayerID;		//玩家ID
	INT32 NewQuality;		//新的品质
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		NewQuality = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(NewQuality);
		return ;
	}
};

//(Client)
class MSG_PlayerRevive_Req
{
	INT32 MsgNo;
	INT32 PlayerID;
	INT8 ReviveOpt;		//复活选项
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		ReviveOpt = pReader->ReadInt8();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt8(ReviveOpt);
		return ;
	}
};

class MSG_ServerRevive_Ack
{
	INT32 RetCode;		//返回码 复活结果
	INT8 ReviveOpt;		//复活选项
	INT32 PlayerID;		//玩家ID
	INT32 Stay;		//是否原地复活
	INT32 ProInc;		//属性增加比例
	INT32 BuffTime;		//buff时长
	INT32 MoneyID;		//货币ID
	INT32 MoneyNum;		//货币数
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		ReviveOpt = pReader->ReadInt8();
		PlayerID = pReader->ReadInt32();
		Stay = pReader->ReadInt32();
		ProInc = pReader->ReadInt32();
		BuffTime = pReader->ReadInt32();
		MoneyID = pReader->ReadInt32();
		MoneyNum = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt8(ReviveOpt);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(Stay);
		pWriter->WriteInt32(ProInc);
		pWriter->WriteInt32(BuffTime);
		pWriter->WriteInt32(MoneyID);
		pWriter->WriteInt32(MoneyNum);
		return ;
	}
};

//客户收到和复法返回消息
class MSG_PlayerRevive_Ack
{
	INT32 RetCode;		//返回码 复活结果
	INT32 PlayerID;		//玩家ID
	INT32 MoneyID;		//货币ID
	INT32 MoneyNum;		//货币数
	INT8 BatCamp;		//角色阵营
	INT32 Heros_Cnt;		//英雄数
	MSG_HeroObj Heros[1];		//英雄数据
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		MoneyID = pReader->ReadInt32();
		MoneyNum = pReader->ReadInt32();
		BatCamp = pReader->ReadInt8();
		Heros_Cnt = pReader->ReadInt32();
		Heros = new MSG_HeroObj[Heros_Cnt];
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(MoneyID);
		pWriter->WriteInt32(MoneyNum);
		pWriter->WriteInt8(BatCamp);
		pWriter->WriteInt32(Heros_Cnt);
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

//玩家复活的通知消息
class MSG_Revive_Nty
{
	INT8 BattleCamp;		//角色阵营
	INT32 Heros_Cnt;
	MSG_HeroObj Heros[1];
	void Read(PacketReader *pReader)
	{
		BattleCamp = pReader->ReadInt8();
		Heros_Cnt = pReader->ReadInt32();
		Heros = new MSG_HeroObj[Heros_Cnt];
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt8(BattleCamp);
		pWriter->WriteInt32(Heros_Cnt);
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

class MSG_KillEvent_Req
{
	INT32 PlayerID;		//杀手
	INT32 Kill;		//杀人数
	INT32 Destroy;		//团灭数
	INT32 SeriesKill;		//连杀人数
	void Read(PacketReader *pReader)
	{
		PlayerID = pReader->ReadInt32();
		Kill = pReader->ReadInt32();
		Destroy = pReader->ReadInt32();
		SeriesKill = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(Kill);
		pWriter->WriteInt32(Destroy);
		pWriter->WriteInt32(SeriesKill);
		return ;
	}
};

class MSG_KillEvent_Ack
{
	INT32 PlayerID;		//杀手
	INT32 KillHonor;		//杀人荣誉
	INT32 KillNum;		//杀人数
	INT32 CurRank;		//当前排名
	void Read(PacketReader *pReader)
	{
		PlayerID = pReader->ReadInt32();
		KillHonor = pReader->ReadInt32();
		KillNum = pReader->ReadInt32();
		CurRank = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(KillHonor);
		pWriter->WriteInt32(KillNum);
		pWriter->WriteInt32(CurRank);
		return ;
	}
};

//阵营战服务器向游戏服务器加载数据
class MSG_LoadCampBattle_Req
{
	INT32 PlayerID;
	INT32 EnterCode;		//进入码
	void Read(PacketReader *pReader)
	{
		PlayerID = pReader->ReadInt32();
		EnterCode = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(EnterCode);
		return ;
	}
};

class MSG_LoadObject
{
	INT32 HeroID;		//英雄ID
	INT8 Camp;		//英雄阵营
	INT32 PropertyValue[11];		//数值属性
	INT32 PropertyPercent[11];		//百分比属性
	INT32 CampDef[5];		//抗阵营属性
	INT32 CampKill[5];		//灭阵营属性
	INT32 SkillID;		//英雄技能
	INT32 AttackID;		//攻击属性ID
	void Read(PacketReader *pReader)
	{
		HeroID = pReader->ReadInt32();
		Camp = pReader->ReadInt8();
		for(int i = 0; i < 11; i++)
		{
			PropertyValue[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 11; i++)
		{
			PropertyPercent[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampDef[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampKill[i] = pReader->ReadInt32();
		}
		SkillID = pReader->ReadInt32();
		AttackID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(HeroID);
		pWriter->WriteInt8(Camp);
		for(int i = 0; i < 11; i++)
		{
			pWriter->WriteInt32(PropertyValue[i]);
		}
		for(int i = 0; i < 11; i++)
		{
			pWriter->WriteInt32(PropertyPercent[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			pWriter->WriteInt32(CampDef[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			pWriter->WriteInt32(CampKill[i]);
		}
		pWriter->WriteInt32(SkillID);
		pWriter->WriteInt32(AttackID);
		return ;
	}
};

class MSG_LoadCampBattle_Ack
{
	INT32 RetCode;		//返回码
	INT32 PlayerID;		//角色ID
	INT8 BattleCamp;		//角色阵营
	INT32 RoomType;		//房间类型
	INT32 Level;		//主角等级
	INT32 LeftTimes;		//剩余的移动次数
	INT32 MoveEndTime;		//搬运结束时间
	INT32 CurRank;		//今日排名
	INT32 KillNum;		//今日击杀
	INT32 KillHonor;		//今日杀人荣誉
	MSG_LoadObject Heros[6];		//英雄对象
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		BattleCamp = pReader->ReadInt8();
		RoomType = pReader->ReadInt32();
		Level = pReader->ReadInt32();
		LeftTimes = pReader->ReadInt32();
		MoveEndTime = pReader->ReadInt32();
		CurRank = pReader->ReadInt32();
		KillNum = pReader->ReadInt32();
		KillHonor = pReader->ReadInt32();
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt8(BattleCamp);
		pWriter->WriteInt32(RoomType);
		pWriter->WriteInt32(Level);
		pWriter->WriteInt32(LeftTimes);
		pWriter->WriteInt32(MoveEndTime);
		pWriter->WriteInt32(CurRank);
		pWriter->WriteInt32(KillNum);
		pWriter->WriteInt32(KillHonor);
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

//阵营战服务器向玩家通知新的技能
class MSG_NewSkill_Nty
{
	INT32 NewSkillID;		//新的技能ID
	void Read(PacketReader *pReader)
	{
		NewSkillID = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(NewSkillID);
		return ;
	}
};

//游戏战斗目标数据
class MSG_HeroData
{
	INT32 HeroID;
	INT32 PropertyValue[11];
	INT32 PropertyPercent[11];
	INT32 CampDef[5];
	INT32 CampKill[5];
	void Read(PacketReader *pReader)
	{
		HeroID = pReader->ReadInt32();
		for(int i = 0; i < 11; i++)
		{
			PropertyValue[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 11; i++)
		{
			PropertyPercent[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampDef[i] = pReader->ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampKill[i] = pReader->ReadInt32();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(HeroID);
		for(int i = 0; i < 11; i++)
		{
			pWriter->WriteInt32(PropertyValue[i]);
		}
		for(int i = 0; i < 11; i++)
		{
			pWriter->WriteInt32(PropertyPercent[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			pWriter->WriteInt32(CampDef[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			pWriter->WriteInt32(CampKill[i]);
		}
		return ;
	}
};

//游戏战斗目标数据
class MSG_PlayerData
{
	INT32 PlayerID;
	INT8 Quality;
	INT32 FightValue;
	MSG_HeroData Heros[6];
	void Read(PacketReader *pReader)
	{
		PlayerID = pReader->ReadInt32();
		Quality = pReader->ReadInt8();
		FightValue = pReader->ReadInt32();
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Read(pReader);
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt8(Quality);
		pWriter->WriteInt32(FightValue);
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(pWriter);
		}
		return ;
	}
};

//游戏服务器的运营数据
class MSG_SvrLogData
{
	INT32 SvrID;		//服务器ID
	INT32 ChnlID;		//渠道ID
	INT32 PlayerID;		//玩家角色ID
	INT32 EventID;		//事件ID
	INT32 SrcID;		//来源ID
	INT32 Time;		//事件发生时间
	INT32 Level;		//角色等级
	INT8 VipLvl;		//角色VIP等级
	INT32 Param[2];		//事件的参数
	void Read(PacketReader *pReader)
	{
		SvrID = pReader->ReadInt32();
		ChnlID = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		EventID = pReader->ReadInt32();
		SrcID = pReader->ReadInt32();
		Time = pReader->ReadInt32();
		Level = pReader->ReadInt32();
		VipLvl = pReader->ReadInt8();
		for(int i = 0; i < 2; i++)
		{
			Param[i] = pReader->ReadInt32();
		}
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(SvrID);
		pWriter->WriteInt32(ChnlID);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteInt32(EventID);
		pWriter->WriteInt32(SrcID);
		pWriter->WriteInt32(Time);
		pWriter->WriteInt32(Level);
		pWriter->WriteInt8(VipLvl);
		for(int i = 0; i < 2; i++)
		{
			pWriter->WriteInt32(Param[i]);
		}
		return ;
	}
};

//游戏服务器的心跳消息//(Client)
class MSG_HeartBeat_Req
{
	INT32 MsgNo;
	INT32 SendID;		//客户端是玩家ID, 服务器是服务器ID
	INT32 BeatCode;		//心跳码
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		SendID = pReader->ReadInt32();
		BeatCode = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(SendID);
		pWriter->WriteInt32(BeatCode);
		return ;
	}
};

class MSG_HeroAllDie_Nty
{
	INT32 NtyCode;		//通知类型，暂不用
	void Read(PacketReader *pReader)
	{
		NtyCode = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(NtyCode);
		return ;
	}
};

//(Client)
class MSG_CmapBatChat_Req
{
	INT32 MsgNo;		//消息编号
	INT32 PlayerID;		//角色ID
	string Name;		//角色名
	string Content;		//消息内容
	void Read(PacketReader *pReader)
	{
		MsgNo = pReader->ReadInt32();
		PlayerID = pReader->ReadInt32();
		Name = pReader->ReadString();
		Content = pReader->ReadString();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(MsgNo);
		pWriter->WriteInt32(PlayerID);
		pWriter->WriteString(Name);
		pWriter->WriteString(Content);
		return ;
	}
};

class MSG_CmapBatChat_Ack
{
	INT32 RetCode;		//聊天请求返回码
	void Read(PacketReader *pReader)
	{
		RetCode = pReader->ReadInt32();
		return ;
	}
	void Write(PacketWriter *pWriter)
	{
		pWriter->WriteInt32(RetCode);
		return ;
	}
};


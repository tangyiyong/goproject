using System;
public class MSG_HeroObj
{
	public Int32 HeroID;		//英雄ID
	public Int32 ObjectID;		//英雄实例ID
	public Int32 CurHp;		//英雄血量
	public Single [] Position = new Single[5];		//x,y,z,v,d, x, y,z,速度，方向
	public void Read(PacketReader reader)
	{
		HeroID = reader.ReadInt32();
		ObjectID = reader.ReadInt32();
		CurHp = reader.ReadInt32();
		for(int i = 0; i < 5; i++)
		{
			Position[i] = reader.ReadFloat();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(HeroID);
		writer.WriteInt32(ObjectID);
		writer.WriteInt32(CurHp);
		for(int i = 0; i < 5; i++)
		{
			writer.WriteFloat(Position[i]);
		}
		return ;
	}
};

public class MSG_BattleObj
{
	public SByte BatCamp;
	public MSG_HeroObj [] Heros = new MSG_HeroObj[6];
	public void Read(PacketReader reader)
	{
		BatCamp = reader.ReadInt8();
		for(int i = 0; i < 6; i++)
		{
			Heros[i] = new MSG_HeroObj();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt8(BatCamp);
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

//进入阵营战消息(Client)
public class MSG_EnterRoom_Req
{
	public Int32 PlayerID;
	public Int32 EnterCode;		//进入码
	public Int32 MsgNo;
	public void Read(PacketReader reader)
	{
		PlayerID = reader.ReadInt32();
		EnterCode = reader.ReadInt32();
		MsgNo = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(EnterCode);
		writer.WriteInt32(MsgNo);
		return ;
	}
};

public class MSG_EnterRoom_Ack
{
	public SByte BatCamp;
	public Int32 CurRank;		//今日排名
	public Int32 KillNum;		//今日击杀
	public Int32 KillHonor;		//今日杀人荣誉
	public Int32 LeftTimes;		//剩余搬动次数
	public Int32 MoveEndTime;		//搬运结束时间
	public Int32 BeginMsgNo;		//起始消息编号
	public Int32 [] SkillID = new Int32[4];		//四个技能ID
	public MSG_HeroObj [] Heros = new MSG_HeroObj[6];		//六个英雄
	public void Read(PacketReader reader)
	{
		BatCamp = reader.ReadInt8();
		CurRank = reader.ReadInt32();
		KillNum = reader.ReadInt32();
		KillHonor = reader.ReadInt32();
		LeftTimes = reader.ReadInt32();
		MoveEndTime = reader.ReadInt32();
		BeginMsgNo = reader.ReadInt32();
		for(int i = 0; i < 4; i++)
		{
			SkillID[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 6; i++)
		{
			Heros[i] = new MSG_HeroObj();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt8(BatCamp);
		writer.WriteInt32(CurRank);
		writer.WriteInt32(KillNum);
		writer.WriteInt32(KillHonor);
		writer.WriteInt32(LeftTimes);
		writer.WriteInt32(MoveEndTime);
		writer.WriteInt32(BeginMsgNo);
		for(int i = 0; i < 4; i++)
		{
			writer.WriteInt32(SkillID[i]);
		}
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

public class MSG_EnterRoom_Notify
{
	public Int32 BatObjs_Cnt;
	public MSG_BattleObj [] BatObjs = null;
	public void Read(PacketReader reader)
	{
		BatObjs_Cnt = reader.ReadInt32();
		BatObjs = new MSG_BattleObj[BatObjs_Cnt];
		for(int i = 0; i < BatObjs_Cnt; i++)
		{
			BatObjs[i] = new MSG_BattleObj();
 			BatObjs[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(BatObjs_Cnt);
		for(int i = 0; i < BatObjs_Cnt; i++)
		{
			BatObjs[i].Write(writer);
		}
		return ;
	}
};

//(Client)
public class MSG_LeaveRoom_Req
{
	public Int32 MsgNo;
	public Int32 PlayerID;
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		return ;
	}
};

public class MSG_LeaveRoom_Ack
{
	public Int32 PlayerID;
	public void Read(PacketReader reader)
	{
		PlayerID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(PlayerID);
		return ;
	}
};

public class MSG_LeaveRoom_Notify
{
	public Int32 [] ObjectIDs = new Int32[6];
	public void Read(PacketReader reader)
	{
		for(int i = 0; i < 6; i++)
		{
			ObjectIDs[i] = reader.ReadInt32();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		for(int i = 0; i < 6; i++)
		{
			writer.WriteInt32(ObjectIDs[i]);
		}
		return ;
	}
};

public class MSG_Skill_Item
{
	public Int32 S_ID;
	public Int32 S_Skill_ID;
	public Int32 TargetIDs_Cnt;
	public Int32 [] TargetIDs = null;
	public void Read(PacketReader reader)
	{
		S_ID = reader.ReadInt32();
		S_Skill_ID = reader.ReadInt32();
		TargetIDs_Cnt = reader.ReadInt32();
		TargetIDs = new Int32[TargetIDs_Cnt];
		for(int i = 0; i < TargetIDs_Cnt; i++)
		{
			TargetIDs[i] = reader.ReadInt32();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(S_ID);
		writer.WriteInt32(S_Skill_ID);
		writer.WriteInt32(TargetIDs_Cnt);
		for(int i = 0; i < TargetIDs_Cnt; i++)
		{
			writer.WriteInt32(TargetIDs[i]);
		}
		return ;
	}
};

//(Client)
public class MSG_Skill_Req
{
	public Int32 MsgNo;
	public Int32 PlayerID;
	public Int32 SkillEvents_Cnt;
	public MSG_Skill_Item [] SkillEvents = null;
	public Int32 AttackEvents_Cnt;
	public MSG_Skill_Item [] AttackEvents = null;
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		SkillEvents_Cnt = reader.ReadInt32();
		SkillEvents = new MSG_Skill_Item[SkillEvents_Cnt];
		for(int i = 0; i < SkillEvents_Cnt; i++)
		{
			SkillEvents[i] = new MSG_Skill_Item();
 			SkillEvents[i].Read(reader);
		}
		AttackEvents_Cnt = reader.ReadInt32();
		AttackEvents = new MSG_Skill_Item[AttackEvents_Cnt];
		for(int i = 0; i < AttackEvents_Cnt; i++)
		{
			AttackEvents[i] = new MSG_Skill_Item();
 			AttackEvents[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(SkillEvents_Cnt);
		for(int i = 0; i < SkillEvents_Cnt; i++)
		{
			SkillEvents[i].Write(writer);
		}
		writer.WriteInt32(AttackEvents_Cnt);
		for(int i = 0; i < AttackEvents_Cnt; i++)
		{
			AttackEvents[i].Write(writer);
		}
		return ;
	}
};

public class MSG_Move_Item
{
	public Int32 S_ID;
	public Single [] Position = new Single[5];
	public void Read(PacketReader reader)
	{
		S_ID = reader.ReadInt32();
		for(int i = 0; i < 5; i++)
		{
			Position[i] = reader.ReadFloat();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(S_ID);
		for(int i = 0; i < 5; i++)
		{
			writer.WriteFloat(Position[i]);
		}
		return ;
	}
};

//(Client)
public class MSG_Move_Req
{
	public Int32 MsgNo;		//消息编号
	public Int32 PlayerID;
	public Int32 MoveEvents_Cnt;
	public MSG_Move_Item [] MoveEvents = null;
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		MoveEvents_Cnt = reader.ReadInt32();
		MoveEvents = new MSG_Move_Item[MoveEvents_Cnt];
		for(int i = 0; i < MoveEvents_Cnt; i++)
		{
			MoveEvents[i] = new MSG_Move_Item();
 			MoveEvents[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(MoveEvents_Cnt);
		for(int i = 0; i < MoveEvents_Cnt; i++)
		{
			MoveEvents[i].Write(writer);
		}
		return ;
	}
};

public class MSG_HeroItem
{
	public Int32 ObjectID;
	public Int32 CurHp;
	public void Read(PacketReader reader)
	{
		ObjectID = reader.ReadInt32();
		CurHp = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(ObjectID);
		writer.WriteInt32(CurHp);
		return ;
	}
};

public class MSG_HeroState_Nty
{
	public Int32 Heros_Cnt;
	public MSG_HeroItem [] Heros = null;
	public void Read(PacketReader reader)
	{
		Heros_Cnt = reader.ReadInt32();
		Heros = new MSG_HeroItem[Heros_Cnt];
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i] = new MSG_HeroItem();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(Heros_Cnt);
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

//玩家查询当前的水晶品质//(Client)
public class MSG_PlayerQuery_Req
{
	public Int32 MsgNo;		//消息编号
	public Int32 PlayerID;
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		return ;
	}
};

//玩家查询当前的水晶品质回复
public class MSG_PlayerQuery_Ack
{
	public Int32 RetCode;		//返回码
	public Int32 PlayerID;		//玩家的ID
	public Int32 Quality;		//水晶品质
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		Quality = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(Quality);
		return ;
	}
};

//(Client)
public class MSG_StartCarry_Req
{
	public Int32 MsgNo;		//消息编号
	public Int32 PlayerID;
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		return ;
	}
};

public class MSG_StartCarry_Ack
{
	public Int32 RetCode;		//返回码
	public Int32 PlayerID;		//玩家的ID
	public Int32 EndTime;		//搬运截止时间
	public Int32 LeftTimes;		//剩余搬动次数
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		EndTime = reader.ReadInt32();
		LeftTimes = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(EndTime);
		writer.WriteInt32(LeftTimes);
		return ;
	}
};

//(Client)
public class MSG_FinishCarry_Req
{
	public Int32 MsgNo;		//消息编号
	public Int32 PlayerID;
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		return ;
	}
};

public class MSG_FinishCarry_Ack
{
	public Int32 RetCode;		//返回码
	public Int32 PlayerID;		//玩家的ID
	public Int32 [] MoneyID = new Int32[2];		//货币ID
	public Int32 [] MoneyNum = new Int32[2];		//货币数量
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		for(int i = 0; i < 2; i++)
		{
			MoneyID[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 2; i++)
		{
			MoneyNum[i] = reader.ReadInt32();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt32(PlayerID);
		for(int i = 0; i < 2; i++)
		{
			writer.WriteInt32(MoneyID[i]);
		}
		for(int i = 0; i < 2; i++)
		{
			writer.WriteInt32(MoneyNum[i]);
		}
		return ;
	}
};

//(Client)
public class MSG_PlayerChange_Req
{
	public Int32 MsgNo;		//消息编号
	public Int32 PlayerID;
	public Int32 HighQuality;		//直接选择最高品质
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		HighQuality = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(HighQuality);
		return ;
	}
};

public class MSG_PlayerChange_Ack
{
	public Int32 RetCode;		//返回码
	public Int32 PlayerID;		//玩家ID
	public Int32 NewQuality;		//新的品质
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		NewQuality = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(NewQuality);
		return ;
	}
};

//(Client)
public class MSG_PlayerRevive_Req
{
	public Int32 MsgNo;
	public Int32 PlayerID;
	public SByte ReviveOpt;		//复活选项
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		ReviveOpt = reader.ReadInt8();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		writer.WriteInt8(ReviveOpt);
		return ;
	}
};

public class MSG_ServerRevive_Ack
{
	public Int32 RetCode;		//返回码 复活结果
	public SByte ReviveOpt;		//复活选项
	public Int32 PlayerID;		//玩家ID
	public Int32 Stay;		//是否原地复活
	public Int32 ProInc;		//属性增加比例
	public Int32 BuffTime;		//buff时长
	public Int32 MoneyID;		//货币ID
	public Int32 MoneyNum;		//货币数
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		ReviveOpt = reader.ReadInt8();
		PlayerID = reader.ReadInt32();
		Stay = reader.ReadInt32();
		ProInc = reader.ReadInt32();
		BuffTime = reader.ReadInt32();
		MoneyID = reader.ReadInt32();
		MoneyNum = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt8(ReviveOpt);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(Stay);
		writer.WriteInt32(ProInc);
		writer.WriteInt32(BuffTime);
		writer.WriteInt32(MoneyID);
		writer.WriteInt32(MoneyNum);
		return ;
	}
};

//客户收到和复法返回消息
public class MSG_PlayerRevive_Ack
{
	public Int32 RetCode;		//返回码 复活结果
	public Int32 PlayerID;		//玩家ID
	public Int32 MoneyID;		//货币ID
	public Int32 MoneyNum;		//货币数
	public SByte BatCamp;		//角色阵营
	public Int32 Heros_Cnt;		//英雄数
	public MSG_HeroObj [] Heros = null;		//英雄数据
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		MoneyID = reader.ReadInt32();
		MoneyNum = reader.ReadInt32();
		BatCamp = reader.ReadInt8();
		Heros_Cnt = reader.ReadInt32();
		Heros = new MSG_HeroObj[Heros_Cnt];
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i] = new MSG_HeroObj();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(MoneyID);
		writer.WriteInt32(MoneyNum);
		writer.WriteInt8(BatCamp);
		writer.WriteInt32(Heros_Cnt);
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

//玩家复活的通知消息
public class MSG_Revive_Nty
{
	public SByte BattleCamp;		//角色阵营
	public Int32 Heros_Cnt;
	public MSG_HeroObj [] Heros = null;
	public void Read(PacketReader reader)
	{
		BattleCamp = reader.ReadInt8();
		Heros_Cnt = reader.ReadInt32();
		Heros = new MSG_HeroObj[Heros_Cnt];
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i] = new MSG_HeroObj();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt8(BattleCamp);
		writer.WriteInt32(Heros_Cnt);
		for(int i = 0; i < Heros_Cnt; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

public class MSG_KillEvent_Req
{
	public Int32 PlayerID;		//杀手
	public Int32 Kill;		//杀人数
	public Int32 Destroy;		//团灭数
	public Int32 SeriesKill;		//连杀人数
	public void Read(PacketReader reader)
	{
		PlayerID = reader.ReadInt32();
		Kill = reader.ReadInt32();
		Destroy = reader.ReadInt32();
		SeriesKill = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(Kill);
		writer.WriteInt32(Destroy);
		writer.WriteInt32(SeriesKill);
		return ;
	}
};

public class MSG_KillEvent_Ack
{
	public Int32 PlayerID;		//杀手
	public Int32 KillHonor;		//杀人荣誉
	public Int32 KillNum;		//杀人数
	public Int32 CurRank;		//当前排名
	public void Read(PacketReader reader)
	{
		PlayerID = reader.ReadInt32();
		KillHonor = reader.ReadInt32();
		KillNum = reader.ReadInt32();
		CurRank = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(KillHonor);
		writer.WriteInt32(KillNum);
		writer.WriteInt32(CurRank);
		return ;
	}
};

//阵营战服务器向游戏服务器加载数据
public class MSG_LoadCampBattle_Req
{
	public Int32 PlayerID;
	public Int32 EnterCode;		//进入码
	public void Read(PacketReader reader)
	{
		PlayerID = reader.ReadInt32();
		EnterCode = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(EnterCode);
		return ;
	}
};

public class MSG_LoadObject
{
	public Int32 HeroID;		//英雄ID
	public SByte Camp;		//英雄阵营
	public Int32 [] PropertyValue = new Int32[11];		//数值属性
	public Int32 [] PropertyPercent = new Int32[11];		//百分比属性
	public Int32 [] CampDef = new Int32[5];		//抗阵营属性
	public Int32 [] CampKill = new Int32[5];		//灭阵营属性
	public Int32 SkillID;		//英雄技能
	public Int32 AttackID;		//攻击属性ID
	public void Read(PacketReader reader)
	{
		HeroID = reader.ReadInt32();
		Camp = reader.ReadInt8();
		for(int i = 0; i < 11; i++)
		{
			PropertyValue[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 11; i++)
		{
			PropertyPercent[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampDef[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampKill[i] = reader.ReadInt32();
		}
		SkillID = reader.ReadInt32();
		AttackID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(HeroID);
		writer.WriteInt8(Camp);
		for(int i = 0; i < 11; i++)
		{
			writer.WriteInt32(PropertyValue[i]);
		}
		for(int i = 0; i < 11; i++)
		{
			writer.WriteInt32(PropertyPercent[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			writer.WriteInt32(CampDef[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			writer.WriteInt32(CampKill[i]);
		}
		writer.WriteInt32(SkillID);
		writer.WriteInt32(AttackID);
		return ;
	}
};

public class MSG_LoadCampBattle_Ack
{
	public Int32 RetCode;		//返回码
	public Int32 PlayerID;		//角色ID
	public SByte BattleCamp;		//角色阵营
	public Int32 RoomType;		//房间类型
	public Int32 Level;		//主角等级
	public Int32 LeftTimes;		//剩余的移动次数
	public Int32 MoveEndTime;		//搬运结束时间
	public Int32 CurRank;		//今日排名
	public Int32 KillNum;		//今日击杀
	public Int32 KillHonor;		//今日杀人荣誉
	public MSG_LoadObject [] Heros = new MSG_LoadObject[6];		//英雄对象
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		BattleCamp = reader.ReadInt8();
		RoomType = reader.ReadInt32();
		Level = reader.ReadInt32();
		LeftTimes = reader.ReadInt32();
		MoveEndTime = reader.ReadInt32();
		CurRank = reader.ReadInt32();
		KillNum = reader.ReadInt32();
		KillHonor = reader.ReadInt32();
		for(int i = 0; i < 6; i++)
		{
			Heros[i] = new MSG_LoadObject();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		writer.WriteInt32(PlayerID);
		writer.WriteInt8(BattleCamp);
		writer.WriteInt32(RoomType);
		writer.WriteInt32(Level);
		writer.WriteInt32(LeftTimes);
		writer.WriteInt32(MoveEndTime);
		writer.WriteInt32(CurRank);
		writer.WriteInt32(KillNum);
		writer.WriteInt32(KillHonor);
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

//阵营战服务器向玩家通知新的技能
public class MSG_NewSkill_Nty
{
	public Int32 NewSkillID;		//新的技能ID
	public void Read(PacketReader reader)
	{
		NewSkillID = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(NewSkillID);
		return ;
	}
};

//游戏战斗目标数据
public class MSG_HeroData
{
	public Int32 HeroID;
	public Int32 [] PropertyValue = new Int32[11];
	public Int32 [] PropertyPercent = new Int32[11];
	public Int32 [] CampDef = new Int32[5];
	public Int32 [] CampKill = new Int32[5];
	public void Read(PacketReader reader)
	{
		HeroID = reader.ReadInt32();
		for(int i = 0; i < 11; i++)
		{
			PropertyValue[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 11; i++)
		{
			PropertyPercent[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampDef[i] = reader.ReadInt32();
		}
		for(int i = 0; i < 5; i++)
		{
			CampKill[i] = reader.ReadInt32();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(HeroID);
		for(int i = 0; i < 11; i++)
		{
			writer.WriteInt32(PropertyValue[i]);
		}
		for(int i = 0; i < 11; i++)
		{
			writer.WriteInt32(PropertyPercent[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			writer.WriteInt32(CampDef[i]);
		}
		for(int i = 0; i < 5; i++)
		{
			writer.WriteInt32(CampKill[i]);
		}
		return ;
	}
};

//游戏战斗目标数据
public class MSG_PlayerData
{
	public Int32 PlayerID;
	public SByte Quality;
	public Int32 FightValue;
	public MSG_HeroData [] Heros = new MSG_HeroData[6];
	public void Read(PacketReader reader)
	{
		PlayerID = reader.ReadInt32();
		Quality = reader.ReadInt8();
		FightValue = reader.ReadInt32();
		for(int i = 0; i < 6; i++)
		{
			Heros[i] = new MSG_HeroData();
 			Heros[i].Read(reader);
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(PlayerID);
		writer.WriteInt8(Quality);
		writer.WriteInt32(FightValue);
		for(int i = 0; i < 6; i++)
		{
			Heros[i].Write(writer);
		}
		return ;
	}
};

//游戏服务器的运营数据
public class MSG_SvrLogData
{
	public Int32 SvrID;		//服务器ID
	public Int32 ChnlID;		//渠道ID
	public Int32 PlayerID;		//玩家角色ID
	public Int32 EventID;		//事件ID
	public Int32 SrcID;		//来源ID
	public Int32 Time;		//事件发生时间
	public Int32 Level;		//角色等级
	public SByte VipLvl;		//角色VIP等级
	public Int32 [] Param = new Int32[2];		//事件的参数
	public void Read(PacketReader reader)
	{
		SvrID = reader.ReadInt32();
		ChnlID = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		EventID = reader.ReadInt32();
		SrcID = reader.ReadInt32();
		Time = reader.ReadInt32();
		Level = reader.ReadInt32();
		VipLvl = reader.ReadInt8();
		for(int i = 0; i < 2; i++)
		{
			Param[i] = reader.ReadInt32();
		}
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(SvrID);
		writer.WriteInt32(ChnlID);
		writer.WriteInt32(PlayerID);
		writer.WriteInt32(EventID);
		writer.WriteInt32(SrcID);
		writer.WriteInt32(Time);
		writer.WriteInt32(Level);
		writer.WriteInt8(VipLvl);
		for(int i = 0; i < 2; i++)
		{
			writer.WriteInt32(Param[i]);
		}
		return ;
	}
};

//游戏服务器的心跳消息//(Client)
public class MSG_HeartBeat_Req
{
	public Int32 MsgNo;
	public Int32 SendID;		//客户端是玩家ID, 服务器是服务器ID
	public Int32 BeatCode;		//心跳码
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		SendID = reader.ReadInt32();
		BeatCode = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(SendID);
		writer.WriteInt32(BeatCode);
		return ;
	}
};

public class MSG_HeroAllDie_Nty
{
	public Int32 NtyCode;		//通知类型，暂不用
	public void Read(PacketReader reader)
	{
		NtyCode = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(NtyCode);
		return ;
	}
};

//(Client)
public class MSG_CmapBatChat_Req
{
	public Int32 MsgNo;		//消息编号
	public Int32 PlayerID;		//角色ID
	public string Name;		//角色名
	public string Content;		//消息内容
	public void Read(PacketReader reader)
	{
		MsgNo = reader.ReadInt32();
		PlayerID = reader.ReadInt32();
		Name = reader.ReadString();
		Content = reader.ReadString();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(MsgNo);
		writer.WriteInt32(PlayerID);
		writer.WriteString(Name);
		writer.WriteString(Content);
		return ;
	}
};

public class MSG_CmapBatChat_Ack
{
	public Int32 RetCode;		//聊天请求返回码
	public void Read(PacketReader reader)
	{
		RetCode = reader.ReadInt32();
		return ;
	}
	public void Write(PacketWriter writer)
	{
		writer.WriteInt32(RetCode);
		return ;
	}
};


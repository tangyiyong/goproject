struct MSG_HeroObj 
{
	INT32 HeroID;   //英雄ID
	INT32 ObjectID; //英雄实例ID
	INT32 CurHp;    //英雄血量
	FLOAT Position[5];//x,y,z,v,d, x, y,z,速度，方向
}

struct MSG_BattleObj 
{
	INT8 BatCamp;
	MSG_HeroObj Heros[6];
}

//进入阵营战消息(Client)
struct  MSG_EnterRoom_Req  
{
	INT32 PlayerID;
	INT32 EnterCode;  //进入码
	INT32 MsgNo;
}

struct MSG_EnterRoom_Ack  
{
	INT8  BatCamp;
	INT32 CurRank;  	   //今日排名
	INT32 KillNum;         //今日击杀
	INT32 KillHonor;	   //今日杀人荣誉
	INT32 LeftTimes;       //剩余搬动次数
	INT32 MoveEndTime;     //搬运结束时间
	INT32 BeginMsgNo;	   //起始消息编号
	INT32 SkillID[4];      //四个技能ID
	MSG_HeroObj Heros[6];
}

struct MSG_EnterRoom_Notify  
{
	INT32 BatObjs_Cnt;
	MSG_BattleObj BatObjs[1];
}

//(Client)
struct MSG_LeaveRoom_Req  
{	
	INT32 MsgNo;
	INT32 PlayerID;
}

struct MSG_LeaveRoom_Ack 
{
	INT32 PlayerID;
}

struct MSG_LeaveRoom_Notify  
{
	INT32 ObjectIDs[6];
}


struct MSG_Skill_Item  
{
	INT32 S_ID;
	INT32 S_Skill_ID;
	INT32 TargetIDs_Cnt;
	INT32 TargetIDs[1]
}

//(Client)
struct MSG_Skill_Req  
{	
	INT32 MsgNo;
	INT32 PlayerID;
	INT32 SkillEvents_Cnt;
	MSG_Skill_Item SkillEvents[1];
	INT32 AttackEvents_Cnt;
	MSG_Skill_Item AttackEvents[1];     
}

struct MSG_Move_Item  
{
	INT32 S_ID;
	FLOAT Position[5];
}

//(Client)
struct MSG_Move_Req  
{	
	INT32 MsgNo;
	INT32 PlayerID;
	INT32 MoveEvents_Cnt;
	MSG_Move_Item MoveEvents[1];
}

struct MSG_HeroItem  
{
	INT32 ObjectID;
	INT32 CurHp;
}

struct MSG_HeroState_Nty 
{	
	INT32 Heros_Cnt;
	MSG_HeroItem Heros[1];
}

//玩家查询当前的水晶品质//(Client)
struct MSG_PlayerQuery_Req  
{
	INT32 MsgNo;
	INT32 PlayerID;
}

//玩家查询当前的水晶品质回复
struct MSG_PlayerQuery_Ack  
{
	INT32 RetCode;  //返回码
	INT32 PlayerID; //玩家的ID
	INT32 Quality;  //水晶品质
}

//(Client)
struct MSG_StartCarry_Req  
{	
	INT32 MsgNo;
	INT32 PlayerID;
}

struct MSG_StartCarry_Ack  
{
	INT32 RetCode; 		//返回码
	INT32 PlayerID;		//玩家的ID
	INT32 EndTime; 		//搬运截止时间
	INT32 LeftTimes; 	//剩余搬动次数
}

//(Client)
struct MSG_FinishCarry_Req  
{
	INT32 MsgNo;
	INT32 PlayerID;
	
}

struct MSG_FinishCarry_Ack  
{
	INT32 RetCode; 		//返回码
	INT32 PlayerID;		//玩家的ID
	INT32 MoneyID[2]; 	//货币ID
	INT32 MoneyNum[2]; 	//货币数量
}

//(Client)
struct MSG_PlayerChange_Req  
{
	INT32 MsgNo;
	INT32 PlayerID;
	INT32 HighQuality; //直接选择最高品质
	
}

struct MSG_PlayerChange_Ack  
{
	INT32 RetCode; //返回码
	INT32 PlayerID; //玩家ID
	INT32 NewQuality; //新的品质
}

//(Client)
struct MSG_PlayerRevive_Req  
{
	INT32 MsgNo;
	INT32 PlayerID;
	INT8 ReviveOpt; //复活选项
}

struct MSG_ServerRevive_Ack  
{
	INT32 RetCode;  //返回码 复活结果
	INT8  ReviveOpt; //复活选项
	INT32 PlayerID; //玩家ID
	INT32 Stay;     //是否原地复活
	INT32 ProInc;   //属性增加比例
	INT32 BuffTime; //buff时长
	INT32 MoneyID;  //货币ID
	INT32 MoneyNum; //货币数
}

//客户收到和复法返回消息
struct MSG_PlayerRevive_Ack  
{
	INT32 RetCode; 		//返回码 复活结果
	INT32 PlayerID; 	//玩家ID
	INT32 MoneyID;  	//货币ID
	INT32 MoneyNum; 	//货币数
	INT8 BatCamp;	//角色阵营
	INT32 Heros_Cnt;	//英雄数
	MSG_HeroObj Heros[1]; //英雄数据
}

//MSG_REVIVE_NTY
//玩家复活的通知消息
struct MSG_Revive_Nty 
{	
	INT8  BattleCamp;		//角色阵营
	INT32 Heros_Cnt;
	MSG_HeroObj Heros[1];
}

struct MSG_KillEvent_Req  
{
	INT32 PlayerID; 	//杀手
	INT32 Kill;     	//杀人数
	INT32 Destroy;   	//团灭数
	INT32 SeriesKill;   //连杀人数
}

struct MSG_KillEvent_Ack 
{
	INT32 PlayerID; 	//杀手
	INT32 KillHonor; 	//杀人荣誉
	INT32 KillNum;   	//杀人数
	INT32 CurRank;   	//当前排名
}



//阵营战服务器向游戏服务器加载数据
struct MSG_LoadCampBattle_Req {
	INT32 PlayerID;
	INT32 EnterCode; //进入码
}

struct MSG_LoadObject {
	INT32 HeroID;				//英雄ID
	INT8  Camp;     			//英雄阵营
	INT32 PropertyValue[11]; 	//数值属性
	INT32 PropertyPercent[11]; 	//百分比属性
	INT32 CampDef[5];  			//抗阵营属性
	INT32 CampKill[5];  		//灭阵营属性
	INT32 SkillID;    			//英雄技能
	INT32 AttackID;				//攻击属性ID
}

struct MSG_LoadCampBattle_Ack {
	INT32 RetCode;   		//返回码
	INT32 PlayerID;  		//角色ID
	INT8  BattleCamp;		//角色阵营
	INT32 RoomType;  		//房间类型
	INT32 Level;     		//主角等级
	INT32 LeftTimes; 		//剩余的移动次数
	INT32 MoveEndTime; 		//搬运结束时间
	INT32 CurRank;  	   //今日排名
	INT32 KillNum;         //今日击杀
	INT32 KillHonor;	   //今日杀人荣誉
	MSG_LoadObject Heros[6]//英雄对象
}

//MSG_NEW_SKILL_NTY
//阵营战服务器向玩家通知新的技能
struct MSG_NewSkill_Nty {
	INT32 NewSkillID;  //新的技能ID
}

//游戏战斗目标数据
struct MSG_HeroData  {
	INT32 HeroID;
	INT32 PropertyValue[11];
	INT32 PropertyPercent[11];
	INT32 CampDef[5];
	INT32 CampKill[5];
}

//游戏战斗目标数据
struct MSG_PlayerData  {
	INT32 PlayerID;
	INT8  Quality;
	INT32 FightValue;
	MSG_HeroData Heros[6]
}

//游戏服务器的运营数据
struct MSG_SvrLogData  {
	INT32 EventID; 	//事件ID
	INT32 SrcID;   	//发起人ID
	INT32 TargetID;	//目标人ID
	INT32 Time; 	//事件发生时间
	INT32 Param[4]; //事件的参数
}

//游戏服务器的心跳消息//(Client)
struct MSG_HeartBeat_Req  {
	INT32 MsgNo;
	INT32 SendID;   //客户端是玩家ID, 服务器是服务器ID
	INT32 BeatCode; //心跳码
}

struct MSG_HeroAllDie_Nty 
{	
	INT32 NtyCode; //通知类型，暂不用
}

//(Client)
struct MSG_CmapBatChat_Req 
{	
	INT32 MsgNo;	//消息编号
	INT32 PlayerID; //角色ID
	STRING Name;    //角色名
	STRING Content; //消息内容
}

struct MSG_CmapBatChat_Ack 
{	
	INT32 RetCode; //聊天请求返回码
}



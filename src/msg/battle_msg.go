package msg;
type MSG_HeroObj struct {
	HeroID int32		//英雄ID
	ObjectID int32		//英雄实例ID
	CurHp int32		//英雄血量
	Position[5] float32		//x,y,z,v,d, x, y,z,速度，方向
}

func (self *MSG_HeroObj) Read(reader *PacketReader) bool {
	self.HeroID = reader.ReadInt32()
	self.ObjectID = reader.ReadInt32()
	self.CurHp = reader.ReadInt32()
	for i := 0; i < int(5); i++ {
		self.Position[i] = reader.ReadFloat()
	}
	return true
}

func (self *MSG_HeroObj) Write(writer *PacketWriter) {
	writer.WriteInt32(self.HeroID)
	writer.WriteInt32(self.ObjectID)
	writer.WriteInt32(self.CurHp)
	for i := 0; i < int(5); i++ {
		writer.WriteFloat(self.Position[i]);
	}
	return
}

type MSG_BattleObj struct {
	BatCamp int8
	Heros[6] MSG_HeroObj
}

func (self *MSG_BattleObj) Read(reader *PacketReader) bool {
	self.BatCamp = reader.ReadInt8()
	for i := 0; i < int(6); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_BattleObj) Write(writer *PacketWriter) {
	writer.WriteInt8(self.BatCamp)
	for i := 0; i < int(6); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

//进入阵营战消息(Client)
type MSG_EnterRoom_Req struct {
	PlayerID int32
	EnterCode int32		//进入码
	MsgNo int32
}

func (self *MSG_EnterRoom_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.EnterCode = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_EnterRoom_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.EnterCode)
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_EnterRoom_Ack struct {
	BatCamp int8
	CurRank int32		//今日排名
	KillNum int32		//今日击杀
	KillHonor int32		//今日杀人荣誉
	LeftTimes int32		//剩余搬动次数
	MoveEndTime int32		//搬运结束时间
	BeginMsgNo int32		//起始消息编号
	SkillID[4] int32		//四个技能ID
	Heros[6] MSG_HeroObj		//六个英雄
}

func (self *MSG_EnterRoom_Ack) Read(reader *PacketReader) bool {
	self.BatCamp = reader.ReadInt8()
	self.CurRank = reader.ReadInt32()
	self.KillNum = reader.ReadInt32()
	self.KillHonor = reader.ReadInt32()
	self.LeftTimes = reader.ReadInt32()
	self.MoveEndTime = reader.ReadInt32()
	self.BeginMsgNo = reader.ReadInt32()
	for i := 0; i < int(4); i++ {
		self.SkillID[i] = reader.ReadInt32()
	}
	for i := 0; i < int(6); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_EnterRoom_Ack) Write(writer *PacketWriter) {
	writer.WriteInt8(self.BatCamp)
	writer.WriteInt32(self.CurRank)
	writer.WriteInt32(self.KillNum)
	writer.WriteInt32(self.KillHonor)
	writer.WriteInt32(self.LeftTimes)
	writer.WriteInt32(self.MoveEndTime)
	writer.WriteInt32(self.BeginMsgNo)
	for i := 0; i < int(4); i++ {
		writer.WriteInt32(self.SkillID[i]);
	}
	for i := 0; i < int(6); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

type MSG_EnterRoom_Notify struct {
	BatObjs_Cnt int32
	BatObjs[] MSG_BattleObj
}

func (self *MSG_EnterRoom_Notify) Read(reader *PacketReader) bool {
	self.BatObjs_Cnt = reader.ReadInt32()
	self.BatObjs = make([]MSG_BattleObj,self.BatObjs_Cnt)
	for i := 0; i < int(self.BatObjs_Cnt); i++ {
		self.BatObjs[i].Read(reader)
	}
	return true
}

func (self *MSG_EnterRoom_Notify) Write(writer *PacketWriter) {
	writer.WriteInt32(self.BatObjs_Cnt)
	for i := 0; i < int(self.BatObjs_Cnt); i++ {
		self.BatObjs[i].Write(writer)
	}
	return
}

//(Client)
type MSG_LeaveRoom_Req struct {
	MsgNo int32
	PlayerID int32
}

func (self *MSG_LeaveRoom_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	return true
}

func (self *MSG_LeaveRoom_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	return
}

type MSG_LeaveRoom_Ack struct {
	PlayerID int32
}

func (self *MSG_LeaveRoom_Ack) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	return true
}

func (self *MSG_LeaveRoom_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	return
}

type MSG_LeaveRoom_Notify struct {
	ObjectIDs[6] int32
}

func (self *MSG_LeaveRoom_Notify) Read(reader *PacketReader) bool {
	for i := 0; i < int(6); i++ {
		self.ObjectIDs[i] = reader.ReadInt32()
	}
	return true
}

func (self *MSG_LeaveRoom_Notify) Write(writer *PacketWriter) {
	for i := 0; i < int(6); i++ {
		writer.WriteInt32(self.ObjectIDs[i]);
	}
	return
}

type MSG_Skill_Item struct {
	S_ID int32
	S_Skill_ID int32
	TargetIDs_Cnt int32
	TargetIDs[] int32
}

func (self *MSG_Skill_Item) Read(reader *PacketReader) bool {
	self.S_ID = reader.ReadInt32()
	self.S_Skill_ID = reader.ReadInt32()
	self.TargetIDs_Cnt = reader.ReadInt32()
	self.TargetIDs = make([]int32,self.TargetIDs_Cnt)
	for i := 0; i < int(self.TargetIDs_Cnt); i++ {
		self.TargetIDs[i] = reader.ReadInt32()
	}
	return true
}

func (self *MSG_Skill_Item) Write(writer *PacketWriter) {
	writer.WriteInt32(self.S_ID)
	writer.WriteInt32(self.S_Skill_ID)
	writer.WriteInt32(self.TargetIDs_Cnt)
	for i := 0; i < int(self.TargetIDs_Cnt); i++ {
		writer.WriteInt32(self.TargetIDs[i]);
	}
	return
}

//(Client)
type MSG_Skill_Req struct {
	MsgNo int32
	PlayerID int32
	SkillEvents_Cnt int32
	SkillEvents[] MSG_Skill_Item
	AttackEvents_Cnt int32
	AttackEvents[] MSG_Skill_Item
}

func (self *MSG_Skill_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.SkillEvents_Cnt = reader.ReadInt32()
	self.SkillEvents = make([]MSG_Skill_Item,self.SkillEvents_Cnt)
	for i := 0; i < int(self.SkillEvents_Cnt); i++ {
		self.SkillEvents[i].Read(reader)
	}
	self.AttackEvents_Cnt = reader.ReadInt32()
	self.AttackEvents = make([]MSG_Skill_Item,self.AttackEvents_Cnt)
	for i := 0; i < int(self.AttackEvents_Cnt); i++ {
		self.AttackEvents[i].Read(reader)
	}
	return true
}

func (self *MSG_Skill_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.SkillEvents_Cnt)
	for i := 0; i < int(self.SkillEvents_Cnt); i++ {
		self.SkillEvents[i].Write(writer)
	}
	writer.WriteInt32(self.AttackEvents_Cnt)
	for i := 0; i < int(self.AttackEvents_Cnt); i++ {
		self.AttackEvents[i].Write(writer)
	}
	return
}

type MSG_Move_Item struct {
	S_ID int32
	Position[5] float32
}

func (self *MSG_Move_Item) Read(reader *PacketReader) bool {
	self.S_ID = reader.ReadInt32()
	for i := 0; i < int(5); i++ {
		self.Position[i] = reader.ReadFloat()
	}
	return true
}

func (self *MSG_Move_Item) Write(writer *PacketWriter) {
	writer.WriteInt32(self.S_ID)
	for i := 0; i < int(5); i++ {
		writer.WriteFloat(self.Position[i]);
	}
	return
}

//(Client)
type MSG_Move_Req struct {
	MsgNo int32		//消息编号
	PlayerID int32
	MoveEvents_Cnt int32
	MoveEvents[] MSG_Move_Item
}

func (self *MSG_Move_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.MoveEvents_Cnt = reader.ReadInt32()
	self.MoveEvents = make([]MSG_Move_Item,self.MoveEvents_Cnt)
	for i := 0; i < int(self.MoveEvents_Cnt); i++ {
		self.MoveEvents[i].Read(reader)
	}
	return true
}

func (self *MSG_Move_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MoveEvents_Cnt)
	for i := 0; i < int(self.MoveEvents_Cnt); i++ {
		self.MoveEvents[i].Write(writer)
	}
	return
}

type MSG_HeroItem struct {
	ObjectID int32
	CurHp int32
}

func (self *MSG_HeroItem) Read(reader *PacketReader) bool {
	self.ObjectID = reader.ReadInt32()
	self.CurHp = reader.ReadInt32()
	return true
}

func (self *MSG_HeroItem) Write(writer *PacketWriter) {
	writer.WriteInt32(self.ObjectID)
	writer.WriteInt32(self.CurHp)
	return
}

type MSG_HeroState_Nty struct {
	Heros_Cnt int32
	Heros[] MSG_HeroItem
}

func (self *MSG_HeroState_Nty) Read(reader *PacketReader) bool {
	self.Heros_Cnt = reader.ReadInt32()
	self.Heros = make([]MSG_HeroItem,self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_HeroState_Nty) Write(writer *PacketWriter) {
	writer.WriteInt32(self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

//玩家查询当前的水晶品质//(Client)
type MSG_PlayerQuery_Req struct {
	MsgNo int32		//消息编号
	PlayerID int32
}

func (self *MSG_PlayerQuery_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerQuery_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	return
}

//玩家查询当前的水晶品质回复
type MSG_PlayerQuery_Ack struct {
	RetCode int32		//返回码
	PlayerID int32		//玩家的ID
	Quality int32		//水晶品质
}

func (self *MSG_PlayerQuery_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.Quality = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerQuery_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.Quality)
	return
}

//(Client)
type MSG_StartCarry_Req struct {
	MsgNo int32		//消息编号
	PlayerID int32
}

func (self *MSG_StartCarry_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	return true
}

func (self *MSG_StartCarry_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	return
}

type MSG_StartCarry_Ack struct {
	RetCode int32		//返回码
	PlayerID int32		//玩家的ID
	EndTime int32		//搬运截止时间
	LeftTimes int32		//剩余搬动次数
}

func (self *MSG_StartCarry_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.EndTime = reader.ReadInt32()
	self.LeftTimes = reader.ReadInt32()
	return true
}

func (self *MSG_StartCarry_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.EndTime)
	writer.WriteInt32(self.LeftTimes)
	return
}

//(Client)
type MSG_FinishCarry_Req struct {
	MsgNo int32		//消息编号
	PlayerID int32
}

func (self *MSG_FinishCarry_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	return true
}

func (self *MSG_FinishCarry_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	return
}

type MSG_FinishCarry_Ack struct {
	RetCode int32		//返回码
	PlayerID int32		//玩家的ID
	MoneyID[2] int32		//货币ID
	MoneyNum[2] int32		//货币数量
}

func (self *MSG_FinishCarry_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	for i := 0; i < int(2); i++ {
		self.MoneyID[i] = reader.ReadInt32()
	}
	for i := 0; i < int(2); i++ {
		self.MoneyNum[i] = reader.ReadInt32()
	}
	return true
}

func (self *MSG_FinishCarry_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt32(self.PlayerID)
	for i := 0; i < int(2); i++ {
		writer.WriteInt32(self.MoneyID[i]);
	}
	for i := 0; i < int(2); i++ {
		writer.WriteInt32(self.MoneyNum[i]);
	}
	return
}

//(Client)
type MSG_PlayerChange_Req struct {
	MsgNo int32		//消息编号
	PlayerID int32
	HighQuality int32		//直接选择最高品质
}

func (self *MSG_PlayerChange_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.HighQuality = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerChange_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.HighQuality)
	return
}

type MSG_PlayerChange_Ack struct {
	RetCode int32		//返回码
	PlayerID int32		//玩家ID
	NewQuality int32		//新的品质
}

func (self *MSG_PlayerChange_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.NewQuality = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerChange_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.NewQuality)
	return
}

//(Client)
type MSG_PlayerRevive_Req struct {
	MsgNo int32
	PlayerID int32
	ReviveOpt int8		//复活选项
}

func (self *MSG_PlayerRevive_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.ReviveOpt = reader.ReadInt8()
	return true
}

func (self *MSG_PlayerRevive_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt8(self.ReviveOpt)
	return
}

type MSG_ServerRevive_Ack struct {
	RetCode int32		//返回码 复活结果
	ReviveOpt int8		//复活选项
	PlayerID int32		//玩家ID
	Stay int32		//是否原地复活
	ProInc int32		//属性增加比例
	BuffTime int32		//buff时长
	MoneyID int32		//货币ID
	MoneyNum int32		//货币数
}

func (self *MSG_ServerRevive_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.ReviveOpt = reader.ReadInt8()
	self.PlayerID = reader.ReadInt32()
	self.Stay = reader.ReadInt32()
	self.ProInc = reader.ReadInt32()
	self.BuffTime = reader.ReadInt32()
	self.MoneyID = reader.ReadInt32()
	self.MoneyNum = reader.ReadInt32()
	return true
}

func (self *MSG_ServerRevive_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt8(self.ReviveOpt)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.Stay)
	writer.WriteInt32(self.ProInc)
	writer.WriteInt32(self.BuffTime)
	writer.WriteInt32(self.MoneyID)
	writer.WriteInt32(self.MoneyNum)
	return
}

//客户收到和复法返回消息
type MSG_PlayerRevive_Ack struct {
	RetCode int32		//返回码 复活结果
	PlayerID int32		//玩家ID
	MoneyID int32		//货币ID
	MoneyNum int32		//货币数
	BatCamp int8		//角色阵营
	Heros_Cnt int32		//英雄数
	Heros[] MSG_HeroObj		//英雄数据
}

func (self *MSG_PlayerRevive_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.MoneyID = reader.ReadInt32()
	self.MoneyNum = reader.ReadInt32()
	self.BatCamp = reader.ReadInt8()
	self.Heros_Cnt = reader.ReadInt32()
	self.Heros = make([]MSG_HeroObj,self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_PlayerRevive_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MoneyID)
	writer.WriteInt32(self.MoneyNum)
	writer.WriteInt8(self.BatCamp)
	writer.WriteInt32(self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

//玩家复活的通知消息
type MSG_Revive_Nty struct {
	BattleCamp int8		//角色阵营
	Heros_Cnt int32
	Heros[] MSG_HeroObj
}

func (self *MSG_Revive_Nty) Read(reader *PacketReader) bool {
	self.BattleCamp = reader.ReadInt8()
	self.Heros_Cnt = reader.ReadInt32()
	self.Heros = make([]MSG_HeroObj,self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_Revive_Nty) Write(writer *PacketWriter) {
	writer.WriteInt8(self.BattleCamp)
	writer.WriteInt32(self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

type MSG_KillEvent_Req struct {
	PlayerID int32		//杀手
	Kill int32		//杀人数
	Destroy int32		//团灭数
	SeriesKill int32		//连杀人数
}

func (self *MSG_KillEvent_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.Kill = reader.ReadInt32()
	self.Destroy = reader.ReadInt32()
	self.SeriesKill = reader.ReadInt32()
	return true
}

func (self *MSG_KillEvent_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.Kill)
	writer.WriteInt32(self.Destroy)
	writer.WriteInt32(self.SeriesKill)
	return
}

type MSG_KillEvent_Ack struct {
	PlayerID int32		//杀手
	KillHonor int32		//杀人荣誉
	KillNum int32		//杀人数
	CurRank int32		//当前排名
}

func (self *MSG_KillEvent_Ack) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.KillHonor = reader.ReadInt32()
	self.KillNum = reader.ReadInt32()
	self.CurRank = reader.ReadInt32()
	return true
}

func (self *MSG_KillEvent_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.KillHonor)
	writer.WriteInt32(self.KillNum)
	writer.WriteInt32(self.CurRank)
	return
}

//阵营战服务器向游戏服务器加载数据
type MSG_LoadCampBattle_Req struct {
	PlayerID int32
	EnterCode int32		//进入码
}

func (self *MSG_LoadCampBattle_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.EnterCode = reader.ReadInt32()
	return true
}

func (self *MSG_LoadCampBattle_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.EnterCode)
	return
}

type MSG_LoadObject struct {
	HeroID int32		//英雄ID
	Camp int8		//英雄阵营
	PropertyValue[11] int32		//数值属性
	PropertyPercent[11] int32		//百分比属性
	CampDef[5] int32		//抗阵营属性
	CampKill[5] int32		//灭阵营属性
	SkillID int32		//英雄技能
	AttackID int32		//攻击属性ID
}

func (self *MSG_LoadObject) Read(reader *PacketReader) bool {
	self.HeroID = reader.ReadInt32()
	self.Camp = reader.ReadInt8()
	for i := 0; i < int(11); i++ {
		self.PropertyValue[i] = reader.ReadInt32()
	}
	for i := 0; i < int(11); i++ {
		self.PropertyPercent[i] = reader.ReadInt32()
	}
	for i := 0; i < int(5); i++ {
		self.CampDef[i] = reader.ReadInt32()
	}
	for i := 0; i < int(5); i++ {
		self.CampKill[i] = reader.ReadInt32()
	}
	self.SkillID = reader.ReadInt32()
	self.AttackID = reader.ReadInt32()
	return true
}

func (self *MSG_LoadObject) Write(writer *PacketWriter) {
	writer.WriteInt32(self.HeroID)
	writer.WriteInt8(self.Camp)
	for i := 0; i < int(11); i++ {
		writer.WriteInt32(self.PropertyValue[i]);
	}
	for i := 0; i < int(11); i++ {
		writer.WriteInt32(self.PropertyPercent[i]);
	}
	for i := 0; i < int(5); i++ {
		writer.WriteInt32(self.CampDef[i]);
	}
	for i := 0; i < int(5); i++ {
		writer.WriteInt32(self.CampKill[i]);
	}
	writer.WriteInt32(self.SkillID)
	writer.WriteInt32(self.AttackID)
	return
}

type MSG_LoadCampBattle_Ack struct {
	RetCode int32		//返回码
	PlayerID int32		//角色ID
	BattleCamp int8		//角色阵营
	RoomType int32		//房间类型
	Level int32		//主角等级
	LeftTimes int32		//剩余的移动次数
	MoveEndTime int32		//搬运结束时间
	CurRank int32		//今日排名
	KillNum int32		//今日击杀
	KillHonor int32		//今日杀人荣誉
	Heros[6] MSG_LoadObject		//英雄对象
}

func (self *MSG_LoadCampBattle_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.BattleCamp = reader.ReadInt8()
	self.RoomType = reader.ReadInt32()
	self.Level = reader.ReadInt32()
	self.LeftTimes = reader.ReadInt32()
	self.MoveEndTime = reader.ReadInt32()
	self.CurRank = reader.ReadInt32()
	self.KillNum = reader.ReadInt32()
	self.KillHonor = reader.ReadInt32()
	for i := 0; i < int(6); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_LoadCampBattle_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt8(self.BattleCamp)
	writer.WriteInt32(self.RoomType)
	writer.WriteInt32(self.Level)
	writer.WriteInt32(self.LeftTimes)
	writer.WriteInt32(self.MoveEndTime)
	writer.WriteInt32(self.CurRank)
	writer.WriteInt32(self.KillNum)
	writer.WriteInt32(self.KillHonor)
	for i := 0; i < int(6); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

//阵营战服务器向玩家通知新的技能
type MSG_NewSkill_Nty struct {
	NewSkillID int32		//新的技能ID
}

func (self *MSG_NewSkill_Nty) Read(reader *PacketReader) bool {
	self.NewSkillID = reader.ReadInt32()
	return true
}

func (self *MSG_NewSkill_Nty) Write(writer *PacketWriter) {
	writer.WriteInt32(self.NewSkillID)
	return
}

//游戏战斗目标数据
type MSG_HeroData struct {
	HeroID int32
	PropertyValue[11] int32
	PropertyPercent[11] int32
	CampDef[5] int32
	CampKill[5] int32
}

func (self *MSG_HeroData) Read(reader *PacketReader) bool {
	self.HeroID = reader.ReadInt32()
	for i := 0; i < int(11); i++ {
		self.PropertyValue[i] = reader.ReadInt32()
	}
	for i := 0; i < int(11); i++ {
		self.PropertyPercent[i] = reader.ReadInt32()
	}
	for i := 0; i < int(5); i++ {
		self.CampDef[i] = reader.ReadInt32()
	}
	for i := 0; i < int(5); i++ {
		self.CampKill[i] = reader.ReadInt32()
	}
	return true
}

func (self *MSG_HeroData) Write(writer *PacketWriter) {
	writer.WriteInt32(self.HeroID)
	for i := 0; i < int(11); i++ {
		writer.WriteInt32(self.PropertyValue[i]);
	}
	for i := 0; i < int(11); i++ {
		writer.WriteInt32(self.PropertyPercent[i]);
	}
	for i := 0; i < int(5); i++ {
		writer.WriteInt32(self.CampDef[i]);
	}
	for i := 0; i < int(5); i++ {
		writer.WriteInt32(self.CampKill[i]);
	}
	return
}

//游戏战斗目标数据
type MSG_PlayerData struct {
	PlayerID int32
	Quality int8
	FightValue int32
	Heros[6] MSG_HeroData
}

func (self *MSG_PlayerData) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.Quality = reader.ReadInt8()
	self.FightValue = reader.ReadInt32()
	for i := 0; i < int(6); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_PlayerData) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt8(self.Quality)
	writer.WriteInt32(self.FightValue)
	for i := 0; i < int(6); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

//游戏服务器的运营数据
type MSG_SvrLogData struct {
	SvrID int32		//服务器ID
	PlatID int32		//渠道ID
	PlayerID int32		//玩家角色ID
	EventID int32		//事件ID
	SrcID int32		//来源ID
	Time int32		//事件发生时间
	Level int32		//角色等级
	VipLvl int8		//角色VIP等级
	Param[2] int32		//事件的参数
}

func (self *MSG_SvrLogData) Read(reader *PacketReader) bool {
	self.SvrID = reader.ReadInt32()
	self.PlatID = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.EventID = reader.ReadInt32()
	self.SrcID = reader.ReadInt32()
	self.Time = reader.ReadInt32()
	self.Level = reader.ReadInt32()
	self.VipLvl = reader.ReadInt8()
	for i := 0; i < int(2); i++ {
		self.Param[i] = reader.ReadInt32()
	}
	return true
}

func (self *MSG_SvrLogData) Write(writer *PacketWriter) {
	writer.WriteInt32(self.SvrID)
	writer.WriteInt32(self.PlatID)
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.EventID)
	writer.WriteInt32(self.SrcID)
	writer.WriteInt32(self.Time)
	writer.WriteInt32(self.Level)
	writer.WriteInt8(self.VipLvl)
	for i := 0; i < int(2); i++ {
		writer.WriteInt32(self.Param[i]);
	}
	return
}

//游戏服务器的心跳消息//(Client)
type MSG_HeartBeat_Req struct {
	MsgNo int32
	SendID int32		//客户端是玩家ID, 服务器是服务器ID
	BeatCode int32		//心跳码
}

func (self *MSG_HeartBeat_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.SendID = reader.ReadInt32()
	self.BeatCode = reader.ReadInt32()
	return true
}

func (self *MSG_HeartBeat_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.SendID)
	writer.WriteInt32(self.BeatCode)
	return
}

type MSG_HeroAllDie_Nty struct {
	NtyCode int32		//通知类型，暂不用
}

func (self *MSG_HeroAllDie_Nty) Read(reader *PacketReader) bool {
	self.NtyCode = reader.ReadInt32()
	return true
}

func (self *MSG_HeroAllDie_Nty) Write(writer *PacketWriter) {
	writer.WriteInt32(self.NtyCode)
	return
}

//(Client)
type MSG_CmapBatChat_Req struct {
	MsgNo int32		//消息编号
	PlayerID int32		//角色ID
	Name string		//角色名
	Content string		//消息内容
}

func (self *MSG_CmapBatChat_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.Name = reader.ReadString()
	self.Content = reader.ReadString()
	return true
}

func (self *MSG_CmapBatChat_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.PlayerID)
	writer.WriteString(self.Name)
	writer.WriteString(self.Content)
	return
}

type MSG_CmapBatChat_Ack struct {
	RetCode int32		//聊天请求返回码
}

func (self *MSG_CmapBatChat_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	return true
}

func (self *MSG_CmapBatChat_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	return
}


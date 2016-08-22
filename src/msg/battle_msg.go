package msg;
type MSG_HeroObj struct {
	HeroID int
	ObjectID int
	CurHp int
	Position[5] float32
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
	BatCamp int
	Heros[6] MSG_HeroObj
}

func (self *MSG_BattleObj) Read(reader *PacketReader) bool {
	self.BatCamp = reader.ReadInt32()
	for i := 0; i < int(6); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_BattleObj) Write(writer *PacketWriter) {
	writer.WriteInt32(self.BatCamp)
	for i := 0; i < int(6); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

type MSG_EnterRoom_Req struct {
	PlayerID int
	EnterCode int
	MsgNo int
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
	BatCamp int
	CurRank int
	KillNum int
	KillHonor int
	LeftTimes int
	MoveEndTime int
	BeginMsgNo int
	SkillID[4] int
	Heros[6] MSG_HeroObj
}

func (self *MSG_EnterRoom_Ack) Read(reader *PacketReader) bool {
	self.BatCamp = reader.ReadInt32()
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
	writer.WriteInt32(self.BatCamp)
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
	BatObjs_Cnt int
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

type MSG_LeaveRoom_Req struct {
	MsgNo int
	PlayerID int
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
	PlayerID int
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
	ObjectIDs[6] int
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
	S_ID int
	S_Skill_ID int
	TargetIDs_Cnt int
	TargetIDs[] int
}

func (self *MSG_Skill_Item) Read(reader *PacketReader) bool {
	self.S_ID = reader.ReadInt32()
	self.S_Skill_ID = reader.ReadInt32()
	self.TargetIDs_Cnt = reader.ReadInt32()
	self.TargetIDs = make([]int,self.TargetIDs_Cnt)
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

type MSG_Skill_Req struct {
	MsgNo int
	SkillEvents_Cnt int
	SkillEvents[] MSG_Skill_Item
	AttackEvents_Cnt int
	AttackEvents[] MSG_Skill_Item
}

func (self *MSG_Skill_Req) Read(reader *PacketReader) bool {
	self.MsgNo = reader.ReadInt32()
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
	S_ID int
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

type MSG_Move_Req struct {
	MoveEvents_Cnt int
	MoveEvents[] MSG_Move_Item
	MsgNo int
}

func (self *MSG_Move_Req) Read(reader *PacketReader) bool {
	self.MoveEvents_Cnt = reader.ReadInt32()
	self.MoveEvents = make([]MSG_Move_Item,self.MoveEvents_Cnt)
	for i := 0; i < int(self.MoveEvents_Cnt); i++ {
		self.MoveEvents[i].Read(reader)
	}
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_Move_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.MoveEvents_Cnt)
	for i := 0; i < int(self.MoveEvents_Cnt); i++ {
		self.MoveEvents[i].Write(writer)
	}
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_HeroItem struct {
	ObjectID int
	CurHp int
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
	Heros_Cnt int
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

type MSG_PlayerQuery_Req struct {
	PlayerID int
	MsgNo int
}

func (self *MSG_PlayerQuery_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerQuery_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_PlayerQuery_Ack struct {
	RetCode int
	PlayerID int
	Quality int
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

type MSG_StartCarry_Req struct {
	PlayerID int
	MsgNo int
}

func (self *MSG_StartCarry_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_StartCarry_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_StartCarry_Ack struct {
	RetCode int
	PlayerID int
	EndTime int
	LeftTimes int
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

type MSG_FinishCarry_Req struct {
	PlayerID int
	MsgNo int
}

func (self *MSG_FinishCarry_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_FinishCarry_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_FinishCarry_Ack struct {
	RetCode int
	PlayerID int
	MoneyID[2] int
	MoneyNum[2] int
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

type MSG_PlayerChange_Req struct {
	PlayerID int
	HighQuality int
	MsgNo int
}

func (self *MSG_PlayerChange_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.HighQuality = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerChange_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.HighQuality)
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_PlayerChange_Ack struct {
	RetCode int
	PlayerID int
	NewQuality int
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

type MSG_PlayerRevive_Req struct {
	PlayerID int
	MsgNo int
	ReviveOpt int
}

func (self *MSG_PlayerRevive_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	self.ReviveOpt = reader.ReadInt32()
	return true
}

func (self *MSG_PlayerRevive_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MsgNo)
	writer.WriteInt32(self.ReviveOpt)
	return
}

type MSG_ServerRevive_Ack struct {
	RetCode int
	PlayerID int
	Stay int
	ProInc int
	BuffTime int
	MoneyID int
	MoneyNum int
}

func (self *MSG_ServerRevive_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
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
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.Stay)
	writer.WriteInt32(self.ProInc)
	writer.WriteInt32(self.BuffTime)
	writer.WriteInt32(self.MoneyID)
	writer.WriteInt32(self.MoneyNum)
	return
}

type MSG_PlayerRevive_Ack struct {
	RetCode int
	PlayerID int
	MoneyID int
	MoneyNum int
	BattleCamp int
	Heros_Cnt int
	Heros[] MSG_HeroObj
}

func (self *MSG_PlayerRevive_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.MoneyID = reader.ReadInt32()
	self.MoneyNum = reader.ReadInt32()
	self.BattleCamp = reader.ReadInt32()
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
	writer.WriteInt32(self.BattleCamp)
	writer.WriteInt32(self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

type MSG_Revive_Nty struct {
	BattleCamp int
	Heros_Cnt int
	Heros[] MSG_HeroObj
}

func (self *MSG_Revive_Nty) Read(reader *PacketReader) bool {
	self.BattleCamp = reader.ReadInt32()
	self.Heros_Cnt = reader.ReadInt32()
	self.Heros = make([]MSG_HeroObj,self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_Revive_Nty) Write(writer *PacketWriter) {
	writer.WriteInt32(self.BattleCamp)
	writer.WriteInt32(self.Heros_Cnt)
	for i := 0; i < int(self.Heros_Cnt); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

type MSG_KillEvent_Req struct {
	Killer int
	Kill int
	Destroy int
	SeriesKill int
}

func (self *MSG_KillEvent_Req) Read(reader *PacketReader) bool {
	self.Killer = reader.ReadInt32()
	self.Kill = reader.ReadInt32()
	self.Destroy = reader.ReadInt32()
	self.SeriesKill = reader.ReadInt32()
	return true
}

func (self *MSG_KillEvent_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.Killer)
	writer.WriteInt32(self.Kill)
	writer.WriteInt32(self.Destroy)
	writer.WriteInt32(self.SeriesKill)
	return
}

type MSG_KillEvent_Ack struct {
	PlayerID int
	KillHonor int
	KillNum int
	CurRank int
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

type MSG_LoadCampBattle_Req struct {
	PlayerID int
	EnterCode int
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
	HeroID int
	Camp int
	PropertyValue[11] int
	PropertyPercent[11] int
	CampDef[5] int
	CampKill[5] int
	SkillID int
	AttackID int
}

func (self *MSG_LoadObject) Read(reader *PacketReader) bool {
	self.HeroID = reader.ReadInt32()
	self.Camp = reader.ReadInt32()
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
	writer.WriteInt32(self.Camp)
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
	RetCode int
	PlayerID int
	BattleCamp int
	RoomType int
	Level int
	LeftTimes int
	MoveEndTime int
	CurRank int
	KillNum int
	KillHonor int
	Heros[6] MSG_LoadObject
}

func (self *MSG_LoadCampBattle_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	self.PlayerID = reader.ReadInt32()
	self.BattleCamp = reader.ReadInt32()
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
	writer.WriteInt32(self.BattleCamp)
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

type MSG_NewSkill_Nty struct {
	NewSkillID int
}

func (self *MSG_NewSkill_Nty) Read(reader *PacketReader) bool {
	self.NewSkillID = reader.ReadInt32()
	return true
}

func (self *MSG_NewSkill_Nty) Write(writer *PacketWriter) {
	writer.WriteInt32(self.NewSkillID)
	return
}

type MSG_HeroData struct {
	HeroID int
	PropertyValue[11] int
	PropertyPercent[11] int
	CampDef[5] int
	CampKill[5] int
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

type MSG_PlayerData struct {
	PlayerID int
	Quality int
	FightValue int
	Heros[6] MSG_HeroData
}

func (self *MSG_PlayerData) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.Quality = reader.ReadInt32()
	self.FightValue = reader.ReadInt32()
	for i := 0; i < int(6); i++ {
		self.Heros[i].Read(reader)
	}
	return true
}

func (self *MSG_PlayerData) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.Quality)
	writer.WriteInt32(self.FightValue)
	for i := 0; i < int(6); i++ {
		self.Heros[i].Write(writer)
	}
	return
}

type MSG_SvrLogData struct {
	EventID int
	SrcID int
	TargetID int
	Time int
	Param[4] int
}

func (self *MSG_SvrLogData) Read(reader *PacketReader) bool {
	self.EventID = reader.ReadInt32()
	self.SrcID = reader.ReadInt32()
	self.TargetID = reader.ReadInt32()
	self.Time = reader.ReadInt32()
	for i := 0; i < int(4); i++ {
		self.Param[i] = reader.ReadInt32()
	}
	return true
}

func (self *MSG_SvrLogData) Write(writer *PacketWriter) {
	writer.WriteInt32(self.EventID)
	writer.WriteInt32(self.SrcID)
	writer.WriteInt32(self.TargetID)
	writer.WriteInt32(self.Time)
	for i := 0; i < int(4); i++ {
		writer.WriteInt32(self.Param[i]);
	}
	return
}

type MSG_HeartBeat_Req struct {
	SendID int
	BeatCode int
	MsgNo int
}

func (self *MSG_HeartBeat_Req) Read(reader *PacketReader) bool {
	self.SendID = reader.ReadInt32()
	self.BeatCode = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	return true
}

func (self *MSG_HeartBeat_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.SendID)
	writer.WriteInt32(self.BeatCode)
	writer.WriteInt32(self.MsgNo)
	return
}

type MSG_HeroAllDie_Nty struct {
	NtyCode int
}

func (self *MSG_HeroAllDie_Nty) Read(reader *PacketReader) bool {
	self.NtyCode = reader.ReadInt32()
	return true
}

func (self *MSG_HeroAllDie_Nty) Write(writer *PacketWriter) {
	writer.WriteInt32(self.NtyCode)
	return
}

type MSG_CmapBatChat_Req struct {
	PlayerID int
	MsgNo int
	Name string
	Content string
}

func (self *MSG_CmapBatChat_Req) Read(reader *PacketReader) bool {
	self.PlayerID = reader.ReadInt32()
	self.MsgNo = reader.ReadInt32()
	self.Name = reader.ReadString()
	self.Content = reader.ReadString()
	return true
}

func (self *MSG_CmapBatChat_Req) Write(writer *PacketWriter) {
	writer.WriteInt32(self.PlayerID)
	writer.WriteInt32(self.MsgNo)
	writer.WriteString(self.Name)
	writer.WriteString(self.Content)
	return
}

type MSG_CmapBatChat_Ack struct {
	RetCode int
}

func (self *MSG_CmapBatChat_Ack) Read(reader *PacketReader) bool {
	self.RetCode = reader.ReadInt32()
	return true
}

func (self *MSG_CmapBatChat_Ack) Write(writer *PacketWriter) {
	writer.WriteInt32(self.RetCode)
	return
}


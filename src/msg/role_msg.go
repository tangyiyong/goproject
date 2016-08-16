package msg

//玩家更改角色名
//消息:/change_role_name
type MSG_ChangeRoleName_Req struct {
	PlayerID   int    //玩家ID
	SessionKey string //Sessionkey
	NewName    string //新的角色名
}

type MSG_ChangeRoleName_Ack struct {
	RetCode int //返回码
}

//玩家请求新手向导信息
//消息:/get_new_wizard
type MSG_GetNewWizard_Req struct {
	PlayerID   int
	SessionKey string
}

type MSG_GetNewWizard_Ack struct {
	RetCode   int    //返回码
	NewWizard string //新手向导信息
}

//玩家保存新手向导信息
//消息:/set_new_wizard
type MSG_SetNewWizard_Req struct {
	PlayerID   int
	SessionKey string
	NewWizard  string //新手向导信息
}

type MSG_SetNewWizard_Ack struct {
	RetCode int //返回码
}

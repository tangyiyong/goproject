package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"msg"
	"sync"
	"utility"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const Start_Account_ID = 10000 //10000以下的ID留给服务器使用。

//账号表结构
type TAccount struct {
	ID         int32  `bson:"_id"` //账号ID
	Name       string //账户名
	Pwd        string //密码
	CreateTime int32  //创建时间
	LastTime   int32  //上次登录时间
	Channel    int32  //渠道ID
	Enable     int32  //是否禁用
	LastSvrID  int32  //上次登录的GameSvrID
}

type TAccountMgr struct {
	accmutex       sync.Mutex
	curAccountID   int32
	accountNameMap map[string]int32
	accountMap     map[int32]*TAccount
	loginKeyMap    map[int32]string
}

var (
	G_AccountMgr TAccountMgr
)

func (self *TAccountMgr) GetAccountByName(name string) (*TAccount, bool) {
	self.accmutex.Lock()
	defer self.accmutex.Unlock()
	id, ok := self.accountNameMap[name]
	if ok && (id > 0) {
		return self.accountMap[id], true
	}

	return nil, false
}

func (self *TAccountMgr) GetAccountByID(id int32) (*TAccount, bool) {
	self.accmutex.Lock()
	defer self.accmutex.Unlock()
	pAccount, ok := self.accountMap[id]
	if ok && (pAccount != nil) {
		return pAccount, true
	}

	return nil, false
}

func (self *TAccountMgr) ResetAccount(name string, password string, newname string, newpassword string) bool {
	self.accmutex.Lock()
	accountid, ok := self.accountNameMap[name]
	if ok != true || accountid <= 0 {
		self.accmutex.Unlock()
		return false
	}

	_, ok = self.accountNameMap[newname]
	if ok == true {
		//新的账号己被人使用
		self.accmutex.Unlock()
		return false
	}

	var pAccount *TAccount = self.accountMap[accountid]
	pAccount.Name = newname
	pAccount.Pwd = newpassword
	self.accountNameMap[newname] = pAccount.ID
	delete(self.accountNameMap, name)
	self.accmutex.Unlock()
	mongodb.UpdateToDB("Account", &bson.M{"_id": accountid}, &bson.M{"$set": bson.M{
		"name":     newname,
		"password": newpassword}})
	return true
}

func (self *TAccountMgr) AddNewAccount(name string, password string) (*TAccount, int) {
	self.accmutex.Lock()
	_, ok := self.accountNameMap[name]
	if ok == true {
		self.accmutex.Unlock()
		return nil, msg.RE_ACCOUNT_EXIST
	}

	var account TAccount
	account.CreateTime = utility.GetCurTime()
	account.Enable = 1
	account.Name = name
	account.Pwd = password
	account.LastTime = utility.GetCurTime()
	account.LastSvrID = 0
	account.ID = self.GetNextAccountID()
	self.accountMap[account.ID] = &account
	self.accountNameMap[name] = account.ID
	self.accmutex.Unlock()
	return &account, msg.RE_SUCCESS
}

func (self *TAccountMgr) GetNextAccountID() (ret int32) {
	ret = self.curAccountID
	self.curAccountID += 1
	return
}

func (self *TAccountMgr) AddLoginKey(accountid int32, loginkey string) {
	self.accmutex.Lock()
	defer self.accmutex.Unlock()
	self.loginKeyMap[accountid] = loginkey
	return
}

func (self *TAccountMgr) ResetLastSvrID(accountid int32, svrid int32) {
	self.accmutex.Lock()
	defer self.accmutex.Unlock()
	pAccount, ok := self.accountMap[accountid]
	if pAccount == nil || ok == false {
		gamelog.Error("ResetLastSvrID Error!!!, invalid accountid:%d", accountid)
		return
	}

	pAccount.LastSvrID = svrid

	return
}

func (self *TAccountMgr) ResetLastLoginTime(accountID int32, loginTime int32) {
	self.accmutex.Lock()
	defer self.accmutex.Unlock()
	pAccount, ok := self.accountMap[accountID]
	if pAccount == nil || ok == false {
		gamelog.Error("ResetLastSvrID Error!!!, invalid accountid:%d", accountID)
		return
	}

	pAccount.LastTime = loginTime
}

func (self *TAccountMgr) CheckLoginKey(accountid int32, loginkey string) bool {
	self.accmutex.Lock()
	defer self.accmutex.Unlock()

	key, ok := self.loginKeyMap[accountid]
	if ok {
		if key == loginkey {
			return true
		}
	}

	return false
}

func InitAccountMgr() bool {
	var accountset []TAccount
	s := mongodb.GetDBSession()
	defer s.Close()
	err := s.DB(appconfig.AccountDbName).C("Account").Find(nil).Sort("+_id").All(&accountset)
	if err != nil && err != mgo.ErrNotFound {
		gamelog.Error("InitAccountMgr DB Error!!!")
		return false
	}

	if len(accountset) <= 0 {
		G_AccountMgr.curAccountID = Start_Account_ID
	} else {
		G_AccountMgr.curAccountID = accountset[len(accountset)-1].ID + 1
	}

	G_AccountMgr.accountNameMap = make(map[string]int32, 1024)
	G_AccountMgr.accountMap = make(map[int32]*TAccount, 1024)
	G_AccountMgr.loginKeyMap = make(map[int32]string, 1024)

	for i := 0; i < len(accountset); i++ {
		G_AccountMgr.accountNameMap[accountset[i].Name] = accountset[i].ID
		G_AccountMgr.accountMap[accountset[i].ID] = &accountset[i]
	}
	return true
}

func CheckPassword(password string) bool {
	if len(password) <= 0 {
		return false
	}
	return true
}

func CheckAccountName(account string) bool {
	if len(account) <= 0 {
		return false
	}

	return true
}

package mainlogic

import (
	"appconfig"
	"gamelog"
	"mongodb"
	"msg"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const Start_Account_ID = 10000 //10000以下的ID留给服务器使用。

//账号表结构
type TAccount struct {
	AccountID  int32  `bson:"_id"` //账号ID
	Name       string //账户名
	Password   string //密码
	CreateTime int64  //创建时间
	LastLgTime int64  //上次登录时间
	LoginCount int    //登录次数
	DeviceID   int    //设备ID
	Forbidden  bool   //是否禁用
	LastSvrID  int    //上次登录的GameSvrID
}

type TAccountMgr struct {
	accmutex       sync.Mutex
	curAccountID   int32
	accountNameMap map[string]int32
	accountMap     map[int32]TAccount
	loginKeyMap    map[int32]string
}

var (
	G_AccountMgr TAccountMgr
)

func (accountmgr *TAccountMgr) GetAccountByName(name string) (TAccount, bool) {
	accountmgr.accmutex.Lock()
	defer accountmgr.accmutex.Unlock()
	id, ok := accountmgr.accountNameMap[name]
	if ok && (id > 0) {
		return accountmgr.accountMap[id], true
	}

	return TAccount{}, false
}

func (accountmgr *TAccountMgr) ResetAccount(name string, password string, newname string, newpassword string) bool {
	accountmgr.accmutex.Lock()
	accountid, ok := accountmgr.accountNameMap[name]
	if ok != true || accountid <= 0 {
		accountmgr.accmutex.Unlock()
		return false
	}

	_, ok = accountmgr.accountNameMap[newname]
	if ok == true {
		//新的账号己被人使用
		accountmgr.accmutex.Unlock()
		return false
	}

	var account TAccount = accountmgr.accountMap[accountid]
	account.Name = newname
	account.Password = newpassword
	accountmgr.accountMap[accountid] = account
	accountmgr.accountNameMap[newname] = account.AccountID
	accountmgr.accmutex.Unlock()
	mongodb.UpdateToDB(appconfig.AccountDbName, "Account", bson.M{"_id": accountid}, bson.M{"$set": bson.M{
		"name":     newname,
		"password": newpassword}})
	return true
}

func (accountmgr *TAccountMgr) AddNewAccount(name string, password string) (*TAccount, int) {
	accountmgr.accmutex.Lock()
	_, ok := accountmgr.accountNameMap[name]
	if ok == true {
		accountmgr.accmutex.Unlock()
		return nil, msg.RE_ACCOUNT_EXIST
	}

	var account TAccount
	account.CreateTime = time.Now().Unix()
	account.DeviceID = 1
	account.Forbidden = false
	account.LoginCount = 1
	account.Name = name
	account.Password = password
	account.LastSvrID = 0
	account.AccountID = accountmgr.GetNextAccountID()
	accountmgr.accountMap[account.AccountID] = account
	accountmgr.accountNameMap[name] = account.AccountID
	accountmgr.accmutex.Unlock()
	return &account, msg.RE_SUCCESS
}

func (accountmgr *TAccountMgr) GetNextAccountID() (ret int32) {
	ret = accountmgr.curAccountID
	accountmgr.curAccountID += 1
	return
}

func (accountmgr *TAccountMgr) AddLoginKey(accountid int32, loginkey string) {
	accountmgr.accmutex.Lock()
	defer accountmgr.accmutex.Unlock()

	accountmgr.loginKeyMap[accountid] = loginkey

	return
}

func (accountmgr *TAccountMgr) CheckLoginKey(accountid int32, loginkey string) bool {
	accountmgr.accmutex.Lock()
	defer accountmgr.accmutex.Unlock()

	key, ok := accountmgr.loginKeyMap[accountid]
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
	if err != nil {
		if err != mgo.ErrNotFound {
			gamelog.Error("InitAccountMgr DB Error!!!")
			return false
		}
	}

	if len(accountset) <= 0 {
		G_AccountMgr.curAccountID = Start_Account_ID
	} else {
		G_AccountMgr.curAccountID = accountset[len(accountset)-1].AccountID + 1
	}

	G_AccountMgr.accountNameMap = make(map[string]int32, 1024)
	G_AccountMgr.accountMap = make(map[int32]TAccount, 1024)
	G_AccountMgr.loginKeyMap = make(map[int32]string, 1024)

	var acc TAccount
	for _, acc = range accountset {
		G_AccountMgr.accountNameMap[acc.Name] = acc.AccountID
		G_AccountMgr.accountMap[acc.AccountID] = acc
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

func ChangeLoginCountAndLast(AccountID int32, GameSvrDomainID int) {
	db_Session := mongodb.GetDBSession()
	defer db_Session.Close()
	AcountColn := db_Session.DB(appconfig.AccountDbName).C("Account")
	AcountColn.Update(bson.M{"_id": AccountID}, bson.M{"$set": bson.M{"lastsvrid": GameSvrDomainID, "logincount": 1}})
	AcountColn = nil
}

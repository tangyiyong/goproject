package appconfig

import (
	"strconv"
	"time"
)

//服务器配置数据
var (
	Version string

	//账号服配置
	AccountSvrInnerIp string
	AccountSvrOuterIp string
	AccountSvrPort    int
	AccountDbName     string
	AccountDbAddr     string
	AccountMaxCon     int //最大连接数
	AccountLogLevel   int

	//跨服服配置
	CrossSvrInnerIp  string
	CrossSvrOuterIp  string
	CrossSvrPort     int
	CrossSvrHttpPort int //http端口
	CrossSvrMaxCon   int //最大连接数
	CrossLogLevel    int

	//游戏服配置
	GameSvrInnerIp  string
	GameSvrOuterIp  string
	GameSvrPort     int
	GameDbName      string
	GameDbAddr      string
	GameMaxCon      int
	GameLogLevel    int
	GameOpenSvrTime int32  //服务器开服时间
	GameSvrName     string //服务器域名
	GameSvrID       int    //服务器ID

	//聊天服配置
	ChatSvrInnerIp string
	ChatSvrOuterIp string
	ChatSvrPort    int
	ChatSvrMaxCon  int
	ChatLogLevel   int //最大连接数

	//战斗服务器
	BattleSvrInnerIp string
	BattleSvrOuterIp string
	BattleSvrPort    int
	BattleLogLevel   int
	BattleSvrMaxCon  int //最大连接数

	//日志服务器
	LogSvrInnerIp  string
	LogSvrOuterIp  string
	LogSvrPort     int
	LogSvrLogLevel int
	LogSvrMaxCon   int    //最大连接数
	LogSvrFlushCnt int    //
	LogFileType    int    //日志文件类型
	LogFileName    string //日志文件名

	//SDK
	SdkSvrInnerIp string
	SdkSvrOuterIp string
	SdkSvrPort    int
	SdkLogLevel   int
	SdkDataSource string
)

//一些全局变量
var (
	ChatSvrAddr           string //聊天服的地址
	VerifyUserLoginUrl    string //验证登录的URL
	RegToAccountSvrUrl    string //游戏服注删到账号服URL
	RegToCrossSvrUrl      string //游戏服注删到跨服服URL
	RegToSdkSvrUrl        string //游戏服注册表SDK服URL
	RegToGameSvrUrl       string //阵营战服注册表游戏服URL
	CrossQueryScoreTarget string //请求积分目标的URL
	CrossQueryScoreRank   string //请求积分排行的URL
	CrossGetFightTarget   string //Game服请求跨服战斗目标
	GiftCodeSvrUrl        string //礼包码游戏服务器查询礼包URL
)

func ParseConfigValue(key string, value string) {
	switch key {
	case "version":
		Version = value
	case "domain":
		GameSvrName = value
	case "domainid":
		GameSvrID, _ = strconv.Atoi(value)
	case "account_svr_inner_ip":
		AccountSvrInnerIp = value
	case "account_svr_outer_ip":
		AccountSvrOuterIp = value
	case "account_svr_port":
		AccountSvrPort, _ = strconv.Atoi(value)
	case "account_db_name":
		AccountDbName = value
	case "account_db_addr":
		AccountDbAddr = value
	case "account_max_conn":
		AccountMaxCon, _ = strconv.Atoi(value)
	case "account_log_level":
		AccountLogLevel, _ = strconv.Atoi(value)
	case "game_svr_inner_ip":
		GameSvrInnerIp = value
	case "game_svr_outer_ip":
		GameSvrOuterIp = value
	case "game_svr_port":
		GameSvrPort, _ = strconv.Atoi(value)
	case "game_db_name":
		GameDbName = value
	case "game_db_addr":
		GameDbAddr = value
	case "game_max_conn":
		GameMaxCon, _ = strconv.Atoi(value)
	case "game_log_level":
		GameLogLevel, _ = strconv.Atoi(value)
	case "game_open_time":
		t, _ := time.ParseInLocation("20060102_150405", value, time.Local)
		GameOpenSvrTime = int32(t.Unix())
	case "chat_svr_inner_ip":
		ChatSvrInnerIp = value
	case "chat_svr_outer_ip":
		ChatSvrOuterIp = value
	case "chat_svr_port":
		ChatSvrPort, _ = strconv.Atoi(value)
	case "chat_svr_max_con":
		ChatSvrMaxCon, _ = strconv.Atoi(value)
	case "chat_svr_log_level":
		ChatLogLevel, _ = strconv.Atoi(value)
	case "cross_svr_inner_ip":
		CrossSvrInnerIp = value
	case "cross_svr_outer_ip":
		CrossSvrOuterIp = value
	case "cross_svr_port":
		CrossSvrPort, _ = strconv.Atoi(value)
	case "cross_svr_max_con":
		CrossSvrMaxCon, _ = strconv.Atoi(value)
	case "cross_svr_log_level":
		CrossLogLevel, _ = strconv.Atoi(value)
	case "cross_svr_http_port":
		CrossSvrHttpPort, _ = strconv.Atoi(value)
	case "battle_svr_inner_ip":
		BattleSvrInnerIp = value
	case "battle_svr_outer_ip":
		BattleSvrOuterIp = value
	case "battle_svr_port":
		BattleSvrPort, _ = strconv.Atoi(value)
	case "battle_svr_max_con":
		BattleSvrMaxCon, _ = strconv.Atoi(value)
	case "battle_svr_log_level":
		BattleLogLevel, _ = strconv.Atoi(value)
	case "log_svr_inner_ip":
		LogSvrInnerIp = value
	case "log_svr_outer_ip":
		LogSvrOuterIp = value
	case "log_svr_port":
		LogSvrPort, _ = strconv.Atoi(value)
	case "log_svr_max_con":
		LogSvrMaxCon, _ = strconv.Atoi(value)
	case "log_svr_flush_count":
		LogSvrFlushCnt, _ = strconv.Atoi(value)
	case "log_svr_log_level":
		LogSvrLogLevel, _ = strconv.Atoi(value)
	case "log_svr_file_type":
		LogFileType, _ = strconv.Atoi(value)
	case "log_svr_file_name":
		LogFileName = value
	case "sdk_svr_inner_ip":
		SdkSvrInnerIp = value
	case "sdk_svr_outer_ip":
		SdkSvrOuterIp = value
	case "sdk_svr_port":
		SdkSvrPort, _ = strconv.Atoi(value)
	case "sdk_svr_log_level":
		SdkLogLevel, _ = strconv.Atoi(value)
	case "sdk_svr_datasrc":
		SdkDataSource = value

	case "gmuser":
		ParseGmUser(value)
	default:
		panic("ParseConfigValue key:[" + key + "] need declare a var")
	}
}

func initGlobalVar() bool {
	VerifyUserLoginUrl = "http://" + AccountSvrInnerIp + ":" + strconv.Itoa(AccountSvrPort) + "/verifyuserlogin"
	ChatSvrAddr = ChatSvrOuterIp + ":" + strconv.Itoa(ChatSvrPort)
	RegToAccountSvrUrl = "http://" + AccountSvrInnerIp + ":" + strconv.Itoa(AccountSvrPort) + "/reggameserver"
	RegToCrossSvrUrl = "http://" + CrossSvrInnerIp + ":" + strconv.Itoa(CrossSvrHttpPort) + "/reggameserver"
	RegToSdkSvrUrl = "http://" + SdkSvrOuterIp + ":" + strconv.Itoa(SdkSvrPort) + "/reggameserver"
	RegToGameSvrUrl = "http://" + GameSvrInnerIp + ":" + strconv.Itoa(GameSvrPort) + "/register_battle_svr"
	CrossQueryScoreTarget = "http://" + CrossSvrInnerIp + ":" + strconv.Itoa(CrossSvrHttpPort) + "/cross_query_score_target"
	CrossQueryScoreRank = "http://" + CrossSvrInnerIp + ":" + strconv.Itoa(CrossSvrHttpPort) + "/cross_query_score_rank"
	CrossGetFightTarget = "http://" + CrossSvrInnerIp + ":" + strconv.Itoa(CrossSvrHttpPort) + "/cross_get_fight_target"
	GiftCodeSvrUrl = "http://" + AccountSvrOuterIp + ":" + strconv.Itoa(AccountSvrPort) + "/gamesvr_giftcode"

	return true
}

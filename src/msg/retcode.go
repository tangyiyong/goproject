package msg

const (
	RE_SUCCESS                    = 0   //成功
	RE_FAILED                     = 1   //失败
	RE_UNKNOWN_ERR                = 2   //未知错误
	RE_INVALID_SESSIONKEY         = 3   //无效的SessionKey
	RE_INVALID_LOGINKEY           = 4   //无效的登录key
	RE_INVALID_NAME               = 5   //无效的名字
	RE_INVALID_PASSWORD           = 6   //无效的密码
	RE_ACCOUNT_NOT_EXIST          = 7   //账号不存在
	RE_ACCOUNT_EXIST              = 8   //账号己存在
	RE_FORBIDDED_ACCOUNT          = 9   //账号己被禁用
	RE_ALEADY_HAVE_ROLE           = 10  //角色己创建
	RE_NO_AVALIBLE_SVR            = 11  //没有可用的游戏服务器
	RE_INVALID_PLAYERID           = 12  //无效的playerid
	RE_INVALID_PARAM              = 13  //无效的参数
	RE_ALREADY_RECEIVED           = 14  //奖品己被领取
	RE_TASK_SCORE_NOT_ENOUGH      = 15  //任务积分不足
	RE_TASK_NOT_COMPLETE          = 16  //任务尚未完成
	RE_CAN_NOT_SIGN_PLUS          = 17  //没有充值,无法豪华签名
	RE_NOT_ENOUGH_ITEM            = 18  //道具不足
	RE_NOT_ENOUGH_MONEY           = 19  //货币不足
	RE_NOT_ENOUGH_VIP_LVL         = 20  //VIP等级不足
	RE_NOT_ENOUGH_HERO            = 21  //同名英雄不足
	RE_NOT_ENOUGH_ACTION          = 22  //行动力不足
	RE_NOT_ENOUGH_GEM             = 23  //同名宝物不足
	RE_NOT_ENOUGH_STAR            = 24  //星数不足
	RE_NOT_ENOUGH_LEVEL           = 25  //等级不足
	RE_NOT_ENOUGH_QUALITY         = 26  //品质不足
	RE_NOT_ENOUGH_HERO_LEVEL      = 27  //英雄等级不足
	RE_NOT_ENOUGH_GUILD_EXP       = 28  //公会经验不足
	RE_NOT_ENOUGH_TIMES           = 29  //次数不足
	RE_CANNOT_BE_USE              = 30  //无法被使用
	RE_NOT_HAVE_REFRESH_TIMES     = 31  //已无刷新次数
	RE_NOT_ENOUGH_REFRESH_ITEM    = 32  //道具不足刷新
	RE_ITEM_NOT_EXIST             = 33  //物品不存在
	RE_ITEM_IS_SOLD_OUT           = 34  //物品已售罄
	RE_ALREADY_MAX_LEVEL          = 35  //己达到最大等级
	RE_CHAPTER_NOT_PASS           = 36  //章节未通关
	RE_STRENGTH_NOT_ENOUGH        = 37  //体力不足
	RE_STAR_NOT_ENOUGH            = 38  //星数不足
	RE_NEED_PASS_PRE_COPY         = 39  //需要通关前置关卡
	RE_NEED_THREE_STAR            = 40  //需要三星才能挑战
	RE_COPY_NOT_PASS              = 41  //关卡未通过
	RE_REFRESH_TIMES_NOT_ENOUGH   = 42  //刷新次数不足
	RE_COPY_IS_LOCK               = 43  //关卡未开启
	RE_SANGUOZHI_ITEM_NOT_ENOUGH  = 44  //三国志升星材料不足
	RE_SANGUOZHI_ALEADY_HAVE      = 45  //三国志升星该星已存在
	RE_SANGUOZHI_NEED_PRE_STAR    = 46  //三国志只允许顺序升星
	RE_NOT_RECHARGE               = 47  //玩家没有首充
	RE_NOT_ENOUGH_SUMMON_ITEM     = 48  //没有足够召唤道具
	RE_NOT_ENOUGH_POINT           = 49  //没有足够召唤积分
	RE_NOT_IN_CHALLANGE_LIST      = 50  //竞技场玩家不在可挑战的名单中
	RE_INVALID_COPY_ID            = 51  //无效的副本ID
	RE_NOT_ENOUGH_RANK            = 52  //排名不足购买
	RE_NOT_ENOUGH_PIECE           = 53  //碎片不足
	RE_CHALLENGE_ALEADY_END       = 54  //挑战已经结束
	RE_FUNC_NOT_OPEN              = 55  //功能尚未开放
	RE_NOT_ENOUGH_FIGHT_VALUE     = 56  //战力不足
	RE_CNT_OVER_MAIN_HERO_LEVEL   = 57  //普通英雄等级不能超过主角等级
	RE_NOT_CHALLANGE              = 58  //领地尚未被攻伐
	RE_NOT_HAVE_HERO              = 59  //无此英雄
	RE_ALEADY_HAVE_HERO           = 60  //领地已有武将巡逻
	RE_PATROL_NOT_END             = 61  //领地巡逻尚未结束
	RE_MAX_TERRITORY_SKILL_LEVEL  = 62  //领地技能已至上限
	RE_NOT_ENOUGH_PATROL_TIME     = 63  //累计巡逻时间不足
	RE_NOT_RIOT                   = 64  //领地没有暴动
	RE_SUPPRESS_TIMES_NOT_ENOUGH  = 65  //镇压次数不足
	RE_ROLE_NAME_EXIST            = 66  //角色名己存在
	RE_NOT_FIND_REBEL             = 67  //没有发现叛军
	RE_REBEL_ALEADY_KILL          = 68  //已击杀叛军
	RE_REBEL_ALEADY_ESCAPE        = 69  //叛军已逃跑
	RE_NOT_ENOUGH_EXPLOIT         = 70  //功勋不足
	RE_ALREADY_DIG                = 71  //此处已挖掘过了
	RE_INVALID_EVENT              = 72  //非法事件
	RE_MAP_NOT_COMPLETION         = 73  //地图完成度不足
	RE_REPEATED_GUAJI             = 74  //不得重复挂机
	RE_NOT_REACH_TIME             = 75  //挂机挑战时间未到
	RE_INVADE_ALEADY_ESCAPE       = 76  //入侵已经逃跑
	RE_MAX_STAR                   = 77  //最大星星
	RE_BATTLE_POS_NOT_OPEN        = 78  //上阵位置还未开放
	RE_BACK_POS_NOT_OPEN          = 79  //援军位置还未开放
	RE_HERO_BAG_OVERLOAD          = 80  //英雄背包已满
	RE_AWARD_TIME_END             = 81  //领奖时间已经结束
	RE_BUY_LIMIT                  = 82  //已到购买上限
	RE_ALEADY_BUY                 = 83  //重复购买
	RE_NOT_ENOUGH_NUMBER          = 84  //购买基金人数不足
	RE_ALEADY_HAVE_GUILD          = 85  //已有公会
	RE_HAVE_NOT_GUILD             = 86  //未加入公会
	RE_ALEADY_APPLY               = 87  //重复申请
	RE_NOT_HAVE_PERMISSION        = 88  //无此权限
	RE_NOT_HAVE_APPLY             = 89  //玩家尚未申请
	RE_ALEADY_SACRIFICE           = 90  //已经祭天
	RE_NOT_ENOUGH_SACRIFICE_TIMES = 91  //公会祭天次数不足
	RE_NOT_ENOUGH_GUILD_LEVEL     = 92  //公会等级不足
	RE_CAMP_IS_KILLED             = 93  //攻击阵营已经灭亡
	RE_NOT_HAVE_TREASURE          = 94  //宝箱已被领取完毕
	RE_CANNOT_BE_RECV             = 95  //入帮时间不足以领取此奖励
	RE_BLACK_MARKET_NOT_OPEN      = 96  //黑市尚未开启
	RE_MESSAGE_TOO_LONG           = 97  //消息字数过长
	RE_GUILD_LEADER_CAN_NOT_EXIT  = 98  //公会会长不允许退出公会
	RE_EXIT_GUILD_TIME_NOT_ENOUGH = 99  //离开公会不满24小时
	RE_GUILD_SKILL_LIMIT          = 100 //公会技能研究以致上限
	RE_GUILD_MEMBER_MAX           = 101 //公会成员上限
	RE_NOT_HAVE_TITLE             = 102 //无此称号
	RE_ACTIVITY_NOT_OPEN          = 103 //活动尚未开启
	RE_SELECT_PLAYRE_FAILED       = 104 //选择目标玩家失败
	RE_NOT_ENOUGH_ATTACK_TIMES    = 105 //攻击次数不足
	RE_NOT_ENOUGH_REVENGE_TIMES   = 106 //复仇次数不足
	RE_REACH_FRIEND_NUM_LIMIT     = 107 //己到好友上限
	RE_ALREADY_FRIEND             = 108 //己到好友上限
	RE_AlEADY_SEND                = 110 //已经赠送过鲜花
	RE_NEED_REFRESH               = 111 //需要刷新
	RE_NOT_ENOUGH_FOOD            = 112 //粮草不够
	RE_NOT_ENOUGH_LOGIN_DAY       = 113 //登录天数不足
	RE_PLEASE_GET_JUBAOPENG       = 114 //请先领取聚宝盆奖励
	RE_NOT_UNLOCK                 = 115 //未解锁章节
	RE_NOT_ENOUGH_VALUE           = 116 //阵图值不足
	RE_NOT_ENOUGH_SCORE           = 117 //积分不足
	RE_GUILD_NAME_REPEAT          = 118 //公会名重复
	RE_NOT_FOUND_GUILD            = 119 //未找到公会
	RE_ITEM_NOT_ENOUGH            = 120 //道具不足
	RE_TURNS_NOT_ENOUGH           = 121 //巡回次数不足
	RE_ACTIVITY_IS_OVER           = 122 //活动已结束
	RE_ACTIVITY_NOT_OVER          = 123 //活动尚未结束
	RE_REPEATED_BUY               = 124 //重复购买
	RE_SERVER_LIMIT_NUM           = 125 //服务器人数己满
	RE_SERVER_CANNT_LOGIN         = 126 //服务器繁忙
	RE_ALEADY_REG                 = 127 //重复报名
)

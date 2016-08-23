package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"math/rand"
	"mongodb"

	"gopkg.in/mgo.v2/bson"
	//	"msg"
	"sync"
	"time"
	"utility"
)

type TMainCopy struct {
	CopyID      int //! 副本ID
	BattleTimes int //! 战斗次数
	ResetCount  int //! 当天刷新次数
	StarNum     int //! 星数
}

type TMainChapter struct {
	Chapter    int
	StarAward  [3]bool //! 6 12 15星 额外星级宝箱
	SceneAward [3]bool //! 1 3 5关卡 场景宝箱
}

type TMainCopyData struct {
	CurCopyID  int //! 当前副本ID
	CurChapter int //! 当前章节ID

	CopyInfo []TMainCopy
	Chapter  []TMainChapter
}

type TEliteCopy struct {
	CopyID      int //! 副本ID
	BattleTimes int //! 战斗次数
	ResetCount  int //! 当天刷新次数
	StarNum     int //! 星数
}

type TEliteChapter struct {
	Chapter    int
	StarAward  [3]bool //! 6 12 15星 额外星级宝箱
	SceneAward bool    //! 场景宝箱领取标记
}

type TEliteCopyData struct {
	CurCopyID     int             //! 当前副本ID
	CurChapter    int             //! 当前章节ID
	Chapter       []TEliteChapter //! 章节信息
	CopyInfo      []TEliteCopy    //! 副本信息
	InvadeChapter []int           //! 入侵章节
}

type TDailyCopy struct {
	ResID       int //! 资源类型ID
	IsChallenge bool
}

type TDailyCopyData struct {
	CopyInfo [6]TDailyCopy //! 副本信息
}

type TFamousCopy struct {
	CopyID      int //! 副本ID
	BattleTimes int //! 挑战次数
}

type TFamousChapterData struct {
	PassedCopy   []TFamousCopy //! 已通关的副本
	ChapterAward bool          //! 章节宝箱领取状态
}

type TFamousCopyData struct {
	Chapter     []TFamousChapterData
	CurCopyID   int //当前挂副本ID
	BattleTimes int //! 挑战次数
}

type TCopyMoudle struct {
	PlayerID int32 `bson:"_id"` //玩家ID

	Main   TMainCopyData   //主线副本
	Elite  TEliteCopyData  //精英副本
	Famous TFamousCopyData //名将副本
	Daily  TDailyCopyData  //日常副本

	LastInvadeTime int64  //! 上次产生入侵时间
	ResetDay       uint32 //! 重置天数

	ownplayer *TPlayer //父player指针
}

func (copym *TCopyMoudle) SetPlayerPtr(playerid int32, pPlayer *TPlayer) {
	copym.PlayerID = playerid
	copym.ownplayer = pPlayer
}

//响应玩家创建
func (copym *TCopyMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	copym.PlayerID = playerid
	copym.ResetDay = utility.GetCurDay()

	//创建数据库记录
	//! 主线副本
	copym.Main.CurChapter = 1 //! 初始化为主线关卡第一章

	//! 精英副本
	copym.Elite.CurChapter = 1 //! 初始化精英副本关卡为第一章

	//! 名将副本
	chapters := gamedata.GetFamousChapterCount()
	copym.Famous.Chapter = make([]TFamousChapterData, chapters+1)

	for i, v := range gamedata.GT_DailyResTypeList {
		copym.Daily.CopyInfo[i].ResID = v
		copym.Daily.CopyInfo[i].IsChallenge = false
	}

	go mongodb.InsertToDB(appconfig.GameDbName, "PlayerCopy", copym)
}

//玩家对象销毁
func (copym *TCopyMoudle) OnDestroy(playerid int32) {

}

//玩家进入游戏
func (copym *TCopyMoudle) OnPlayerOnline(playerid int32) {
	//
}

//玩家离开游戏
func (copym *TCopyMoudle) OnPlayerOffline(playerid int32) {
	//
	return
}

//玩家数据从数据库加载
func (copym *TCopyMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	//
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerCopy").Find(bson.M{"_id": playerid}).One(copym)
	if err != nil {
		gamelog.Error("PlayerCopy Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	copym.PlayerID = playerid
	return
}

//! 玩家通过主线关卡
func (main_copy *TCopyMoudle) PlayerPassMainLevels(copyID int, chapter int, star int) {

	//! 设置当前关卡
	if copyID > main_copy.Main.CurCopyID {
		main_copy.Main.CurCopyID = copyID
	}

	//! 设置当前章节
	if chapter > main_copy.Main.CurChapter {
		main_copy.Main.CurChapter = chapter
	}

	isExist := false
	for i := 0; i < len(main_copy.Main.Chapter); i++ {
		if main_copy.Main.Chapter[i].Chapter == chapter {
			isExist = true
		}
	}

	if isExist == false {
		//! 添加章节信息
		var chapterInfo TMainChapter
		chapterInfo.Chapter = chapter
		main_copy.Main.Chapter = append(main_copy.Main.Chapter, chapterInfo)
		go main_copy.AddMainChapterInfo(chapterInfo)
	}

	isExist = false
	for i := 0; i < len(main_copy.Main.CopyInfo); i++ {
		if main_copy.Main.CopyInfo[i].CopyID == copyID {
			main_copy.Main.CopyInfo[i].BattleTimes += 1

			//! 设置挑战星数
			if star > main_copy.Main.CopyInfo[i].StarNum {
				//! 成就任务总星数更新
				main_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_MAINCOPY_STAR, star-main_copy.Main.CopyInfo[i].StarNum)

				main_copy.Main.CopyInfo[i].StarNum = star
			}

			isExist = true
			go main_copy.UpdateMainCopyAt(i)
			break
		}
	}

	//! 如果该关卡不存在,则为新挑战关卡,存储关卡信息
	if isExist == false {
		var mainCopy TMainCopy
		mainCopy.CopyID = copyID
		mainCopy.BattleTimes = 1
		mainCopy.StarNum = star
		main_copy.Main.CopyInfo = append(main_copy.Main.CopyInfo, mainCopy)

		//! 成就任务总星数更新
		main_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_MAINCOPY_STAR, star)

		go main_copy.AddMainCopyInfo(mainCopy)
	}

	//! 日常任务进度加一
	main_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_MAINCOPY_CHALLENGE, 1)
}

//! 玩家通关精英副本
func (elite_copy *TCopyMoudle) PlayerPassEliteLevels(copyID int, chapter int, star int) {

	//! 设置关卡
	if copyID > elite_copy.Elite.CurCopyID {
		elite_copy.Elite.CurCopyID = copyID
	}

	//! 设置当前章节
	if chapter > elite_copy.Elite.CurChapter {
		elite_copy.Elite.CurChapter = chapter
	}

	isExist := false
	for i := 0; i < len(elite_copy.Elite.Chapter); i++ {
		if elite_copy.Elite.Chapter[i].Chapter == chapter {
			isExist = true
		}
	}

	if isExist == false {
		//! 添加章节信息
		var chapterInfo TEliteChapter
		chapterInfo.Chapter = chapter
		elite_copy.Elite.Chapter = append(elite_copy.Elite.Chapter, chapterInfo)
		go elite_copy.AddEliteChapterInfo(chapterInfo)
	}

	isExist = false
	for i := 0; i < len(elite_copy.Elite.CopyInfo); i++ {
		if elite_copy.Elite.CopyInfo[i].CopyID == copyID {
			elite_copy.Elite.CopyInfo[i].BattleTimes += 1

			//! 设置挑战星数
			if star > elite_copy.Elite.CopyInfo[i].StarNum {
				//! 成就任务总星数更新
				elite_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ELITECOPY_STAR, star-elite_copy.Elite.CopyInfo[i].StarNum)

				elite_copy.Elite.CopyInfo[i].StarNum = star
			}

			isExist = true
			go elite_copy.UpdateEliteCopyAt(i)
			break
		}
	}

	//! 如果该关卡不存在,则为新挑战关卡,存储关卡信息
	if isExist == false {
		var eliteCopy TEliteCopy
		eliteCopy.CopyID = copyID
		eliteCopy.BattleTimes = 1
		eliteCopy.StarNum = star
		elite_copy.Elite.CopyInfo = append(elite_copy.Elite.CopyInfo, eliteCopy)

		//! 成就任务总星数更新
		elite_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ELITECOPY_STAR, star)

		go elite_copy.AddEliteCopyInfo(eliteCopy)
	}

	//! 日常任务进度加一
	elite_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ELITECOPY_CHALLENGE, 1)
}

//! 玩家通过日常副本
func (daily_copy *TCopyMoudle) PlayerPassDailyLevels(copyID int) {

	//! 设置通关标记
	dailyCopy := gamedata.GetDailyCopyData(copyID)

	for i := 0; i < len(daily_copy.Daily.CopyInfo); i++ {
		if daily_copy.Daily.CopyInfo[i].ResID == dailyCopy.ResType {
			daily_copy.Daily.CopyInfo[i].IsChallenge = true

			go daily_copy.UpdateDailyCopyMask(i, true)
		}
	}

	//! 增加日常任务完成度
	daily_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_DAILYCOPY_CHALLENGE, 1)
}

//! 玩家通过名将副本
func (famous_copy *TCopyMoudle) PlayerPassFamousLevels(copyID int, curChapter int) bool {
	//! 挑战次数+1
	famous_copy.Famous.BattleTimes += 1
	go famous_copy.UpdateFamousCopyTotalBattleTimes()

	//! 赋值通过关卡ID
	if copyID > famous_copy.Famous.CurCopyID {
		famous_copy.Famous.CurCopyID = copyID
		go famous_copy.UpdateFamousCopyCurCopyID()
	}

	chapterInfo := gamedata.GetFamousChapterInfo(curChapter)
	if chapterInfo.SerialID == copyID {
		famous_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PASS_EPIC_COPY, curChapter)
	}

	//! 不存在则为首胜
	isFirstVictory := false
	isExist := false
	battleTimes := 0
	copyIndex := 0

	for i := 0; i < len(famous_copy.Famous.Chapter[curChapter].PassedCopy); i++ {
		if famous_copy.Famous.Chapter[curChapter].PassedCopy[i].CopyID == copyID {
			isExist = true
			famous_copy.Famous.Chapter[curChapter].PassedCopy[i].BattleTimes += 1
			battleTimes = famous_copy.Famous.Chapter[curChapter].PassedCopy[i].BattleTimes
			copyIndex = i
			break
		}
	}

	if isExist == false {
		isFirstVictory = true
		var famousCopy TFamousCopy
		famousCopy.CopyID = copyID
		famousCopy.BattleTimes = 1
		battleTimes = 1
		famous_copy.Famous.Chapter[curChapter].PassedCopy = append(famous_copy.Famous.Chapter[curChapter].PassedCopy, famousCopy)
		go famous_copy.IncFamousCopy(curChapter, famousCopy)
	} else {
		go famous_copy.UpdateFamousCopyBattleTimes(curChapter, copyIndex, battleTimes)
	}

	famous_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FAMOUSCOPY_CHALLENGE, 1)

	return isFirstVictory
}

//! 获取玩家章节总星数
func (main_copy *TCopyMoudle) GetMainChapterStarNumber(chapter int) int {
	starNum := 0

	chapterInfo := gamedata.GetMainChapterInfo(chapter)
	for n := chapterInfo.StartID; n <= chapterInfo.EndID; n++ {
		isChange := false
		for _, v := range main_copy.Main.CopyInfo {
			if v.CopyID == n {
				//! 有变动的关卡
				starNum += v.StarNum
				isChange = true
				break
			}
		}

		if isChange == false {
			starNum += 3
		}
	}

	return starNum
}

//! 获取玩家精英关卡章节总星数
func (elite_copy *TCopyMoudle) GetEliteChapterStarNumber(chapter int) int {
	starNum := 0

	chapterInfo := gamedata.GetEliteChapterInfo(chapter)
	for n := chapterInfo.StartID; n <= chapterInfo.EndID; n++ {
		isChange := false
		for _, v := range elite_copy.Elite.CopyInfo {
			if v.CopyID == n {
				//! 有变动的关卡
				starNum += v.StarNum
				isChange = true
				break
			}
		}

		if isChange == false {
			starNum += 3
		}
	}

	return starNum
}

//! 查询玩家主线副本是否有可领取的章节奖励
func (main_copy *TCopyMoudle) IsHaveNotReceiveAward(chapter int) bool {

	for _, v := range main_copy.Main.Chapter {
		if v.Chapter == chapter {
			for i := 0; i < 3; i++ {
				if v.SceneAward[i] == false {
					chapterData := gamedata.GetMainChapterInfo(v.Chapter)
					needCopyID := chapterData.SceneAwards[i].Levels

					if needCopyID <= main_copy.Main.CurCopyID {
						return true
					}
				}

				if v.StarAward[i] == false {
					chapterData := gamedata.GetMainChapterInfo(v.Chapter)
					needStarNum := chapterData.StarAwards[i].StarNum

					if needStarNum <= main_copy.GetMainChapterStarNumber(v.Chapter) {
						return true
					}
				}
			}
		}
	}

	return false
}

//! 查询玩家精英副本是否有可领取的章节奖励
func (elite_copy *TCopyMoudle) EliteIsHaveNotReceiveAward(chapter int) bool {

	for _, v := range elite_copy.Elite.Chapter {
		if v.Chapter == chapter {
			for i := 0; i < 3; i++ {
				if v.SceneAward == false {
					chapterData := gamedata.GetEliteChapterInfo(v.Chapter)
					needCopyID := chapterData.SceneAwards.Levels

					if needCopyID <= elite_copy.Elite.CurCopyID {
						return true
					}
				}

				if v.StarAward[i] == false {
					chapterData := gamedata.GetEliteChapterInfo(v.Chapter)
					needStarNum := chapterData.StarAwards[i].StarNum

					if needStarNum <= elite_copy.GetEliteChapterStarNumber(v.Chapter) {
						return true
					}
				}
			}
		}
	}

	return false
}

//! 发放奖励  1->星级奖励 2->场景奖励
const (
	MAIN_AWARD_TYPE_STAR  = 1
	MAIN_AWARD_TYPE_SCENE = 2
)

func (main_copy *TCopyMoudle) PaymentMainAward(chapter int, award int, awardtype int) {
	chapterData := gamedata.GetMainChapterInfo(chapter)
	awardID := 0
	index := 0
	if awardtype == MAIN_AWARD_TYPE_STAR {
		awardID = chapterData.StarAwards[award].AwardID

		for i, v := range main_copy.Main.Chapter {
			if v.Chapter == chapter {
				main_copy.Main.Chapter[i].StarAward[award] = true
				index = i
			}
		}

	} else if awardtype == MAIN_AWARD_TYPE_SCENE {
		awardID = chapterData.SceneAwards[award].AwardID
		for i, v := range main_copy.Main.Chapter {
			if v.Chapter == chapter {
				main_copy.Main.Chapter[i].SceneAward[award] = true
				index = i
			}
		}
	}

	awardItem := gamedata.GetItemsFromAwardID(awardID)
	main_copy.ownplayer.BagMoudle.AddAwardItems(awardItem)
	go main_copy.UpdateMainAward(index)
}

func (elite_copy *TCopyMoudle) PaymentEliteAward(chapter int, award int, awardtype int) {
	chapterData := gamedata.GetEliteChapterInfo(chapter)
	awardID := 0
	index := 0
	if awardtype == MAIN_AWARD_TYPE_STAR {
		awardID = chapterData.StarAwards[award].AwardID

		for i, v := range elite_copy.Elite.Chapter {
			if v.Chapter == chapter {
				elite_copy.Elite.Chapter[i].StarAward[award] = true
				index = i
			}
		}

	} else if awardtype == MAIN_AWARD_TYPE_SCENE {
		awardID = chapterData.SceneAwards.AwardID
		for i, v := range elite_copy.Elite.Chapter {
			if v.Chapter == chapter {
				elite_copy.Elite.Chapter[i].SceneAward = true
				index = i
			}
		}
	}

	awardItem := gamedata.GetItemsFromAwardID(awardID)
	elite_copy.ownplayer.BagMoudle.AddAwardItems(awardItem)
	go elite_copy.UpdateEliteAward(index)
}

//! 获取未有入侵的精英副本章节数
func (elite_copy *TCopyMoudle) GetNoInvadeEliteCount() int {
	//! 获取已通关关卡数目
	chapterCount := elite_copy.GetPassEliteChapter()
	invadeCount := 0

	for i := 1; i <= chapterCount; i++ {
		for _, v := range elite_copy.Elite.InvadeChapter {
			if v == i {
				invadeCount += 1
			}
		}
	}
	return (chapterCount - invadeCount)
}

//! 获取已通过精英副本章节数
func (elite_copy *TCopyMoudle) GetPassEliteChapter() int {
	isEnd := gamedata.IsChapterEnd(elite_copy.Elite.CurCopyID, elite_copy.Elite.CurChapter, gamedata.COPY_TYPE_Elite)

	chapterCount := elite_copy.Elite.CurChapter
	if isEnd == false {
		chapterCount -= 1
	}
	return chapterCount
}

func (elite_copy *TCopyMoudle) IsHaveInvade(chapter int) bool {
	for _, v := range elite_copy.Elite.InvadeChapter {
		if v == chapter {
			return true
		}
	}

	return false
}

func (elite_copy *TCopyMoudle) RemoveInvade(chapter int) bool {
	pos := 0
	for i, v := range elite_copy.Elite.InvadeChapter {
		if v == chapter {
			pos = i
		}
	}

	if pos == 0 {
		elite_copy.Elite.InvadeChapter = elite_copy.Elite.InvadeChapter[1:]
	} else if (pos + 1) == len(elite_copy.Elite.InvadeChapter) {
		elite_copy.Elite.InvadeChapter = elite_copy.Elite.InvadeChapter[:pos]
	} else {
		elite_copy.Elite.InvadeChapter = append(elite_copy.Elite.InvadeChapter[:pos], elite_copy.Elite.InvadeChapter[pos+1:]...)
	}

	go elite_copy.RemoveEliteInvade(chapter)
	return true
}

//! 随机无入侵精英副本章节
func (elite_copy *TCopyMoudle) RandNoInvadeEliteChapter(num int) IntLst {
	if elite_copy.GetPassEliteChapter() == 0 {
		return []int{}
	}

	//! 判断是否还有未产生入侵的章节
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var chapter IntLst
	for {
		randChapter := r.Intn(elite_copy.GetPassEliteChapter()) + 1
		if elite_copy.IsHaveInvade(randChapter) == false {
			elite_copy.Elite.InvadeChapter = append(elite_copy.Elite.InvadeChapter, randChapter)
			chapter.Add(randChapter)
			go elite_copy.AddEliteInvade(randChapter)
		}

		if chapter.Len() == num {
			break
		}

		if elite_copy.GetNoInvadeEliteCount() < num {
			break
		}
	}

	return chapter
}

//! 产生入侵
func (elite_copy *TCopyMoudle) CheckEliteInvade() {
	//! 获取入侵时间
	invadeTime := Int64Lst{int64(gamedata.EliteInvadeTime1), int64(gamedata.EliteInvadeTime2),
		int64(gamedata.EliteInvadeTime3), int64(gamedata.EliteInvadeTime4)}

	//! 获取入侵个数
	invadeNum := []int{gamedata.EliteInvadeNum1, gamedata.EliteInvadeNum2,
		gamedata.EliteInvadeNum3, gamedata.EliteInvadeNum4}

	//! 获取今日凌晨时间
	todayTime := GetTodayTime()

	//! 获取当期按时间
	now := time.Now().Unix()
	for i := 0; i < len(invadeTime); i++ {
		invadeTime[i] = invadeTime[i]*60*60 + todayTime

		if elite_copy.LastInvadeTime > invadeTime[i] {
			//! 去除已刷新个数
			continue
		}

		//! 刷新入侵
		if now >= invadeTime[i] {
			//! 获取刷新个数
			number := invadeNum[i]

			//! 随机两个没有叛军的章节
			elite_copy.RandNoInvadeEliteChapter(number)

			//! 重置上次刷新时间
			elite_copy.LastInvadeTime = invadeTime[i]
		}
	}
	go elite_copy.UpdateEliteInvadeTime()
}

//! 检查扫荡体力是否足够
// func (main_copy *TCopyMoudle) CheckSweepMainAction(times int, copyID int, chapter int) (bool, int) {
// 	baseData := gamedata.GetCopyBaseInfo(copyID)

// 	//! 检查体力
// 	ret := main_copy.ownplayer.RoleMoudle.CheckActionEnough(baseData.ActionType, baseData.ActionValue*times)
// 	if ret == false {
// 		gamelog.Error("Hand_BattleResult error : Not Enough Action")
// 		return false, msg.RE_STRENGTH_NOT_ENOUGH //! 体力不足
// 	}

// 	//! 检查挑战次数
// 	chapterInfo := main_copy.Main.Chapter[chapter]
// 	isExist := false
// 	for _, v := range chapterInfo.CopyInfo {
// 		if v.CopyID == copyID {
// 			isExist = true
// 			if v.BattleTimes+times > 10 {
// 				return false, msg.RE_CHALLENGE_TIMES_NOT_ENOUGH //! 挑战次数不足
// 			}

// 			if v.StarNum != 3 {
// 				return false, msg.RE_NEED_THREE_STAR //! 必须三星才能够扫荡
// 			}
// 		}
// 	}

// 	if isExist == false {
// 		return false, msg.RE_COPY_NOT_PASS //! 关卡未通过
// 	}

// 	return true, msg.RE_SUCCESS
// }

func (self *TCopyMoudle) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TCopyMoudle) OnNewDay(newday uint32) {
	self.UpdateMainReset()
	self.UpdateFamousReset()
	self.UpdateDailyReset()
	self.UpdateEliteReset()

	self.ResetDay = utility.GetCurDay()
	go self.UpdateCopy()
}

func (main_copy *TCopyMoudle) UpdateMainReset() {
	//! 刪除已三星通关的信息
	copyLst := []TMainCopy{}
	for _, v := range main_copy.Main.CopyInfo {
		if v.StarNum != 3 {

			var copyInfo TMainCopy
			copyInfo.CopyID = v.CopyID
			copyInfo.BattleTimes = v.BattleTimes
			copyInfo.ResetCount = v.ResetCount
			copyInfo.StarNum = v.StarNum
			copyLst = append(copyLst, copyInfo)
		}
	}

	main_copy.Main.CopyInfo = copyLst

	//! 删除章节奖励已领取关卡
	chapterLst := []TMainChapter{}
	for _, v := range main_copy.Main.Chapter {
		if v.SceneAward[0] == false || v.SceneAward[1] == false || v.SceneAward[2] == false ||
			v.StarAward[0] == false || v.StarAward[1] == false || v.StarAward[2] == false {
			var chapterInfo TMainChapter
			chapterInfo.Chapter = v.Chapter
			chapterInfo.SceneAward = v.SceneAward
			chapterInfo.StarAward = v.StarAward
			chapterLst = append(chapterLst, chapterInfo)
		}
	}

	main_copy.Main.Chapter = chapterLst
}

//! 精英副本重置
func (elite_copy *TCopyMoudle) UpdateEliteReset() {
	//! 刪除已三星通关的信息
	copyLst := []TEliteCopy{}
	for _, v := range elite_copy.Elite.CopyInfo {
		if v.StarNum != 3 {

			var copyInfo TEliteCopy
			copyInfo.CopyID = v.CopyID
			copyInfo.BattleTimes = v.BattleTimes
			copyInfo.ResetCount = v.ResetCount
			copyInfo.StarNum = v.StarNum
			copyLst = append(copyLst, copyInfo)
		}
	}
	elite_copy.Elite.CopyInfo = copyLst

	//! 删除章节奖励已领取关卡
	chapterLst := []TEliteChapter{}
	for _, v := range elite_copy.Elite.Chapter {
		if v.SceneAward == false ||
			v.StarAward[0] == false || v.StarAward[1] == false || v.StarAward[2] == false {
			var chapterInfo TEliteChapter
			chapterInfo.Chapter = v.Chapter
			chapterInfo.SceneAward = v.SceneAward
			chapterInfo.StarAward = v.StarAward
			chapterLst = append(chapterLst, chapterInfo)
		}
	}

	elite_copy.Elite.Chapter = chapterLst
}

func (daily_copy *TCopyMoudle) UpdateDailyReset() {

	//! 刷新各种数据
	for i, _ := range daily_copy.Daily.CopyInfo {
		daily_copy.Daily.CopyInfo[i].IsChallenge = false
	}
}

func (famous_copy *TCopyMoudle) UpdateFamousReset() {

	//! 刷新各种数据
	famous_copy.Famous.BattleTimes = gamedata.FamousCopyChallengeTimes

	for j, v := range famous_copy.Famous.Chapter {
		for i, _ := range v.PassedCopy {
			famous_copy.Famous.Chapter[j].PassedCopy[i].BattleTimes = 0
		}
	}
}

//! 获取今日开启的日常副本资源类型
func (daily_copy *TCopyMoudle) GetTodayDailyCopy() []int {
	//! 获取今日时间
	today := time.Now()
	openRes := []int{}
	dateType := 0

	day := today.Weekday()
	if day == time.Monday || day == time.Wednesday || day == time.Friday {
		//! 时间类型1副本开启
		dateType = 1
	} else if day == time.Thursday || day == time.Tuesday || day == time.Saturday {
		dateType = 2
	} else if day == time.Sunday {
		dateType = 3
	}

	copyLst := gamedata.GetDailyCopyDataFromType(dateType)

	for _, v := range copyLst {
		isExist := false
		for _, b := range openRes {
			if v.ResType == b {
				isExist = true
				break
			}
		}

		if isExist == false {
			openRes = append(openRes, v.ResType)
		}
	}

	return openRes
}

func (famous_copy *TCopyMoudle) GetFamousCopyInfo(copyID int, chapter int) *TFamousCopy {
	for i, v := range famous_copy.Famous.Chapter[chapter].PassedCopy {
		if v.CopyID == copyID {
			return &famous_copy.Famous.Chapter[chapter].PassedCopy[i]
		}
	}

	gamelog.Error("GetFamousCopyInfo fail. CopyID: %d Chapter: %d", copyID, chapter)
	return nil
}

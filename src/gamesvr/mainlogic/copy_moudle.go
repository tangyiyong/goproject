package mainlogic

import (
	"appconfig"
	"gamelog"
	"gamesvr/gamedata"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"mongodb"
	"sync"
	"time"
	"utility"
)

type TMainCopy struct {
	ID       int //! 副本ID
	Times    int //! 战斗次数
	ResetCnt int //! 当天刷新次数
	StarNum  int //! 星数
}

type TMainChapter struct {
	Chapter    int
	StarAward  BitsType //! 6 12 15星 额外星级宝箱
	SceneAward BitsType //! 1 3 5关卡 场景宝箱
}

type TMainData struct {
	CurID      int //! 当前副本ID
	CurChapter int //! 当前章节ID
	CopyLst    []TMainCopy
	Chapter    []TMainChapter
}

type TEliteCopy struct {
	ID       int //! 副本ID
	Times    int //! 战斗次数
	ResetCnt int //! 当天刷新次数
	StarNum  int //! 星数
}

type TEliteChapter struct {
	Chapter    int
	StarAward  BitsType //! 6 12 15星 额外星级宝箱
	SceneAward bool     //! 场景宝箱领取标记
}

type TEliteData struct {
	CurID         int             //! 当前副本ID
	CurChapter    int             //! 当前章节ID
	Chapter       []TEliteChapter //! 章节信息
	CopyLst       []TEliteCopy    //! 副本信息
	InvadeChapter []int           //! 入侵章节
}

type TDailyCopy struct {
	ResID       int //! 资源类型ID
	IsChallenge bool
}

type TDailyData struct {
	CopyLst [6]TDailyCopy //! 副本信息
}

type TFamousChapter struct {
	PassedCopy IntLst //! 已通关的副本
	BoxAward   bool   //! 章节宝箱领取状态
	Extra      bool   //! 连环计
}

type TFamousData struct {
	Chapter    []TFamousChapter
	CurID      int //! 当前挂副本ID
	CurChapter int //! 当前章节ID
	Times      int //! 挑战次数
}

type TCopyMoudle struct {
	PlayerID       int32       `bson:"_id"` //玩家ID
	Main           TMainData   //! 主线副本
	Elite          TEliteData  //! 精英副本
	Famous         TFamousData //! 名将副本
	Daily          TDailyData  //! 日常副本
	LastInvadeTime int32       //! 上次产生入侵时间
	ResetDay       uint32      //! 重置天数
	ownplayer      *TPlayer    //! 父player指针
}

func (self *TCopyMoudle) SetPlayerPtr(playerid int32, player *TPlayer) {
	self.PlayerID = playerid
	self.ownplayer = player
}

//响应玩家创建
func (self *TCopyMoudle) OnCreate(playerid int32) {
	//初始化各个成员数值
	self.PlayerID = playerid
	self.ResetDay = utility.GetCurDay()

	//创建数据库记录
	//! 主线副本
	self.Main.CurChapter = 1 //! 初始化为主线关卡第一章

	//! 精英副本
	self.Elite.CurChapter = 1 //! 初始化精英副本关卡为第一章

	//! 名将副本
	chapters := gamedata.GetFamousChapterCount()
	self.Famous.Chapter = make([]TFamousChapter, chapters+1)

	for i, v := range gamedata.GT_DailyResTypeList {
		self.Daily.CopyLst[i].ResID = v
		self.Daily.CopyLst[i].IsChallenge = false
	}

	mongodb.InsertToDB("PlayerCopy", self)
}

//玩家对象销毁
func (self *TCopyMoudle) OnDestroy(playerid int32) {

}

//玩家进入游戏
func (self *TCopyMoudle) OnPlayerOnline(playerid int32) {
	//
}

//玩家离开游戏
func (self *TCopyMoudle) OnPlayerOffline(playerid int32) {
	//
	return
}

//玩家数据从数据库加载
func (self *TCopyMoudle) OnPlayerLoad(playerid int32, wg *sync.WaitGroup) {
	s := mongodb.GetDBSession()
	defer s.Close()

	err := s.DB(appconfig.GameDbName).C("PlayerCopy").Find(&bson.M{"_id": playerid}).One(self)
	if err != nil {
		gamelog.Error("PlayerCopy Load Error :%s， PlayerID: %d", err.Error(), playerid)
	}

	if wg != nil {
		wg.Done()
	}
	self.PlayerID = playerid
	return
}

//! 玩家通过主线关卡
func (self *TCopyMoudle) PlayerPassMainLevels(copyID int, chapter int, star int) {
	if copyID > self.Main.CurID {
		self.Main.CurID = copyID
	}

	if chapter > self.Main.CurChapter {
		self.Main.CurChapter = chapter
	}

	isExist := false
	for i := 0; i < len(self.Main.Chapter); i++ {
		if self.Main.Chapter[i].Chapter == chapter {
			isExist = true
		}
	}

	if isExist == false {
		//! 添加章节信息
		var chapterInfo TMainChapter
		chapterInfo.Chapter = chapter
		self.Main.Chapter = append(self.Main.Chapter, chapterInfo)
		self.DB_AddMainChapterInfo(chapterInfo)
	}

	isExist = false
	for i := 0; i < len(self.Main.CopyLst); i++ {
		if self.Main.CopyLst[i].ID == copyID {
			self.Main.CopyLst[i].Times += 1

			//! 设置挑战星数
			if star > self.Main.CopyLst[i].StarNum {
				//! 成就任务总星数更新
				self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_MAINCOPY_STAR, star-self.Main.CopyLst[i].StarNum)

				self.Main.CopyLst[i].StarNum = star
			}

			isExist = true
			self.DB_UpdateMainCopyAt(i)
			break
		}
	}

	//! 如果该关卡不存在,则为新挑战关卡,存储关卡信息
	if isExist == false {
		var mainCopy TMainCopy
		mainCopy.ID = copyID
		mainCopy.Times = 1
		mainCopy.StarNum = star
		self.Main.CopyLst = append(self.Main.CopyLst, mainCopy)

		//! 成就任务总星数更新
		self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_MAINCOPY_STAR, star)

		self.DB_AddMainCopyInfo(mainCopy)
	}

	//! 日常任务进度加一
	self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_MAINCOPY_CHALLENGE, 1)
}

//! 玩家通关精英副本
func (self *TCopyMoudle) PlayerPassEliteLevels(copyID int, chapter int, star int) {

	//! 设置关卡
	if copyID > self.Elite.CurID {
		self.Elite.CurID = copyID
	}

	//! 设置当前章节
	if chapter > self.Elite.CurChapter {
		self.Elite.CurChapter = chapter
	}

	isExist := false
	for i := 0; i < len(self.Elite.Chapter); i++ {
		if self.Elite.Chapter[i].Chapter == chapter {
			isExist = true
		}
	}

	if isExist == false {
		//! 添加章节信息
		var chapterInfo TEliteChapter
		chapterInfo.Chapter = chapter
		self.Elite.Chapter = append(self.Elite.Chapter, chapterInfo)
		self.DB_AddEliteChapterInfo(chapterInfo)
	}

	isExist = false
	for i := 0; i < len(self.Elite.CopyLst); i++ {
		if self.Elite.CopyLst[i].ID == copyID {
			self.Elite.CopyLst[i].Times += 1

			//! 设置挑战星数
			if star > self.Elite.CopyLst[i].StarNum {
				//! 成就任务总星数更新
				self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ELITECOPY_STAR, star-self.Elite.CopyLst[i].StarNum)

				self.Elite.CopyLst[i].StarNum = star
			}

			isExist = true
			self.DB_UpdateEliteCopyAt(i)
			break
		}
	}

	//! 如果该关卡不存在,则为新挑战关卡,存储关卡信息
	if isExist == false {
		var eliteCopy TEliteCopy
		eliteCopy.ID = copyID
		eliteCopy.Times = 1
		eliteCopy.StarNum = star
		self.Elite.CopyLst = append(self.Elite.CopyLst, eliteCopy)

		//! 成就任务总星数更新
		self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ELITECOPY_STAR, star)

		self.DB_AddEliteCopyInfo(eliteCopy)
	}

	//! 日常任务进度加一
	self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_ELITECOPY_CHALLENGE, 1)
}

//! 玩家通过日常副本
func (daily_copy *TCopyMoudle) PlayerPassDailyLevels(copyID int) {

	//! 设置通关标记
	dailyCopy := gamedata.GetDailyCopyData(copyID)

	for i := 0; i < len(daily_copy.Daily.CopyLst); i++ {
		if daily_copy.Daily.CopyLst[i].ResID == dailyCopy.ResType {
			daily_copy.Daily.CopyLst[i].IsChallenge = true

			daily_copy.DB_UpdateDailyCopyMask(i, true)
		}
	}

	//! 增加日常任务完成度
	daily_copy.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_DAILYCOPY_CHALLENGE, 1)
}

//! 玩家通过名将副本
func (self *TCopyMoudle) PlayerPassFamousCopy(curChapter int, copyID int) bool {
	isFirstVictory := false
	if true == gamedata.IsSerialCopy(curChapter, copyID) {
		self.Famous.Chapter[curChapter].Extra = true
		self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_PASS_EPIC_COPY, curChapter)
		self.DB_UpdateFamousExtra(curChapter)
	} else {
		self.Famous.Times += 1
		//! 赋值通过关卡ID
		if copyID > self.Famous.CurID {
			self.Famous.CurID = copyID
			self.Famous.CurChapter = curChapter
			isFirstVictory = true
		}

		if self.Famous.Chapter[curChapter].PassedCopy.IsExist(copyID) < 0 {
			self.Famous.Chapter[curChapter].PassedCopy = append(self.Famous.Chapter[curChapter].PassedCopy, copyID)
			self.DB_AddFamousPassCopy(curChapter, copyID)
		}

		self.DB_UpdateFamousCopyData()
	}

	self.ownplayer.TaskMoudle.AddPlayerTaskSchedule(gamedata.TASK_FAMOUSCOPY_CHALLENGE, 1)
	return isFirstVictory
}

//! 获取玩家章节总星数
func (self *TCopyMoudle) GetMainChapterStarNumber(chapter int) int {
	starNum := 0

	chapterInfo := gamedata.GetMainChapterInfo(chapter)
	for n := chapterInfo.StartID; n <= chapterInfo.EndID; n++ {
		isChange := false
		for _, v := range self.Main.CopyLst {
			if v.ID == n {
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
func (self *TCopyMoudle) GetEliteChapterStarNumber(chapter int) int {
	starNum := 0

	chapterInfo := gamedata.GetEliteChapterInfo(chapter)
	for n := chapterInfo.StartID; n <= chapterInfo.EndID; n++ {
		isChange := false
		for _, v := range self.Elite.CopyLst {
			if v.ID == n {
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
func (self *TCopyMoudle) IsHaveNotReceiveAward(chapter int) bool {

	for _, v := range self.Main.Chapter {
		if v.Chapter == chapter {
			for i := 0; i < 3; i++ {
				if v.SceneAward.Get(i+1) == false {
					chapterData := gamedata.GetMainChapterInfo(v.Chapter)
					needCopyID := chapterData.SceneAwards[i].Levels

					if needCopyID <= self.Main.CurID {
						return true
					}
				}

				if v.StarAward.Get(i+1) == false {
					chapterData := gamedata.GetMainChapterInfo(v.Chapter)
					needStarNum := chapterData.StarAwards[i].StarNum

					if needStarNum <= self.GetMainChapterStarNumber(v.Chapter) {
						return true
					}
				}
			}
		}
	}

	return false
}

//! 查询玩家精英副本是否有可领取的章节奖励
func (self *TCopyMoudle) EliteIsHaveNotReceiveAward(chapter int) bool {

	for _, v := range self.Elite.Chapter {
		if v.Chapter == chapter {
			for i := 0; i < 3; i++ {
				if v.SceneAward == false {
					chapterData := gamedata.GetEliteChapterInfo(v.Chapter)
					needCopyID := chapterData.SceneAwards.Levels

					if needCopyID <= self.Elite.CurID {
						return true
					}
				}

				if v.StarAward.Get(i+1) == false {
					chapterData := gamedata.GetEliteChapterInfo(v.Chapter)
					needStarNum := chapterData.StarAwards[i].StarNum

					if needStarNum <= self.GetEliteChapterStarNumber(v.Chapter) {
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

func (self *TCopyMoudle) PaymentMainAward(chapter int, award int, awardtype int) {
	chapterData := gamedata.GetMainChapterInfo(chapter)
	awardID := 0
	index := 0
	if awardtype == MAIN_AWARD_TYPE_STAR {
		awardID = chapterData.StarAwards[award].AwardID

		for i, v := range self.Main.Chapter {
			if v.Chapter == chapter {
				self.Main.Chapter[i].StarAward.Set(award + 1)
				index = i
			}
		}

	} else if awardtype == MAIN_AWARD_TYPE_SCENE {
		awardID = chapterData.SceneAwards[award].AwardID
		for i, v := range self.Main.Chapter {
			if v.Chapter == chapter {
				self.Main.Chapter[i].SceneAward.Set(award + 1)
				index = i
			}
		}
	}

	awardItem := gamedata.GetItemsFromAwardID(awardID)
	self.ownplayer.BagMoudle.AddAwardItems(awardItem)
	self.DB_UpdateMainAward(index)
}

func (self *TCopyMoudle) PaymentEliteAward(chapter int, award int, awardtype int) {
	chapterData := gamedata.GetEliteChapterInfo(chapter)
	awardID := 0
	index := 0
	if awardtype == MAIN_AWARD_TYPE_STAR {
		awardID = chapterData.StarAwards[award].AwardID

		for i, v := range self.Elite.Chapter {
			if v.Chapter == chapter {
				self.Elite.Chapter[i].StarAward.Set(award + 1)
				index = i
			}
		}

	} else if awardtype == MAIN_AWARD_TYPE_SCENE {
		awardID = chapterData.SceneAwards.AwardID
		for i, v := range self.Elite.Chapter {
			if v.Chapter == chapter {
				self.Elite.Chapter[i].SceneAward = true
				index = i
			}
		}
	}

	awardItem := gamedata.GetItemsFromAwardID(awardID)
	self.ownplayer.BagMoudle.AddAwardItems(awardItem)
	self.DB_UpdateEliteAward(index)
}

//! 获取未有入侵的精英副本章节数
func (self *TCopyMoudle) GetNoInvadeEliteCount() int {
	//! 获取已通关关卡数目
	chapterCount := self.GetPassEliteChapter()
	invadeCount := 0

	for i := 1; i <= chapterCount; i++ {
		for _, v := range self.Elite.InvadeChapter {
			if v == i {
				invadeCount += 1
			}
		}
	}
	return (chapterCount - invadeCount)
}

//! 获取已通过精英副本章节数
func (self *TCopyMoudle) GetPassEliteChapter() int {
	isEnd := gamedata.IsChapterEnd(self.Elite.CurID, self.Elite.CurChapter, gamedata.COPY_TYPE_Elite)

	chapterCount := self.Elite.CurChapter
	if isEnd == false {
		chapterCount -= 1
	}
	return chapterCount
}

func (self *TCopyMoudle) IsHaveInvade(chapter int) bool {
	for _, v := range self.Elite.InvadeChapter {
		if v == chapter {
			return true
		}
	}

	return false
}

func (self *TCopyMoudle) RemoveInvade(chapter int) bool {
	pos := 0
	for i, v := range self.Elite.InvadeChapter {
		if v == chapter {
			pos = i
		}
	}

	if pos == 0 {
		self.Elite.InvadeChapter = self.Elite.InvadeChapter[1:]
	} else if (pos + 1) == len(self.Elite.InvadeChapter) {
		self.Elite.InvadeChapter = self.Elite.InvadeChapter[:pos]
	} else {
		self.Elite.InvadeChapter = append(self.Elite.InvadeChapter[:pos], self.Elite.InvadeChapter[pos+1:]...)
	}

	self.DB_RemoveEliteInvade(chapter)
	return true
}

//! 随机无入侵精英副本章节
func (self *TCopyMoudle) RandNoInvadeEliteChapter(num int) IntLst {
	if self.GetPassEliteChapter() == 0 {
		return []int{}
	}

	//! 判断是否还有未产生入侵的章节
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var chapter IntLst
	for {
		randChapter := r.Intn(self.GetPassEliteChapter()) + 1
		if self.IsHaveInvade(randChapter) == false {
			self.Elite.InvadeChapter = append(self.Elite.InvadeChapter, randChapter)
			chapter.Add(randChapter)
			self.DB_AddEliteInvade(randChapter)
		}

		if chapter.Len() == num {
			break
		}

		if self.GetNoInvadeEliteCount() < num {
			break
		}
	}

	return chapter
}

//! 产生入侵
func (self *TCopyMoudle) CheckEliteInvade() {
	//! 获取今日凌晨时间
	todayTime := utility.GetTodayTime()

	//! 获取当期按时间
	now := utility.GetCurTime()
	for i := 0; i < len(gamedata.EliteInvadeTime); i++ {
		invadeTime := int32(gamedata.EliteInvadeTime[i]*60*60) + todayTime

		if self.LastInvadeTime > invadeTime {
			//! 去除已刷新个数
			continue
		}

		//! 刷新入侵
		if now >= invadeTime {
			//! 获取刷新个数
			number := gamedata.EliteInvadeNum[i]

			//! 随机两个没有叛军的章节
			self.RandNoInvadeEliteChapter(number)

			//! 重置上次刷新时间
			self.LastInvadeTime = invadeTime
		}
	}
	self.DB_UpdateEliteInvadeTime()
}

func (self *TCopyMoudle) CheckReset() {
	if utility.IsSameDay(self.ResetDay) == true {
		return
	}

	self.OnNewDay(utility.GetCurDay())
}

func (self *TCopyMoudle) OnNewDay(newday uint32) {
	self.MainReset()
	self.FamousReset()
	self.DailyReset()
	self.EliteReset()
	self.ResetDay = utility.GetCurDay()
	self.DB_UpdateCopy()
}

func (self *TCopyMoudle) MainReset() {
	//! 刪除已三星通关的信息
	copyLst := []TMainCopy{}
	for _, v := range self.Main.CopyLst {
		if v.StarNum != 3 {

			var copyInfo TMainCopy
			copyInfo.ID = v.ID
			copyInfo.Times = v.Times
			copyInfo.ResetCnt = v.ResetCnt
			copyInfo.StarNum = v.StarNum
			copyLst = append(copyLst, copyInfo)
		}
	}

	self.Main.CopyLst = copyLst

	//! 删除章节奖励已领取关卡
	chapterLst := []TMainChapter{}
	for _, v := range self.Main.Chapter {
		if v.SceneAward.Get(1) == false || v.SceneAward.Get(2) == false || v.SceneAward.Get(3) == false ||
			v.StarAward.Get(1) == false || v.StarAward.Get(2) == false || v.StarAward.Get(3) == false {
			var chapterInfo TMainChapter
			chapterInfo.Chapter = v.Chapter
			chapterInfo.SceneAward = v.SceneAward
			chapterInfo.StarAward = v.StarAward
			chapterLst = append(chapterLst, chapterInfo)
		}
	}

	self.Main.Chapter = chapterLst
}

//! 精英副本重置
func (self *TCopyMoudle) EliteReset() {
	//! 刪除已三星通关的信息
	copyLst := []TEliteCopy{}
	for _, v := range self.Elite.CopyLst {
		if v.StarNum != 3 {

			var copyInfo TEliteCopy
			copyInfo.ID = v.ID
			copyInfo.Times = v.Times
			copyInfo.ResetCnt = v.ResetCnt
			copyInfo.StarNum = v.StarNum
			copyLst = append(copyLst, copyInfo)
		}
	}
	self.Elite.CopyLst = copyLst

	//! 删除章节奖励已领取关卡
	chapterLst := []TEliteChapter{}
	for _, v := range self.Elite.Chapter {
		if v.SceneAward == false ||
			v.StarAward.Get(1) == false || v.StarAward.Get(2) == false || v.StarAward.Get(3) == false {
			var chapterInfo TEliteChapter
			chapterInfo.Chapter = v.Chapter
			chapterInfo.SceneAward = v.SceneAward
			chapterInfo.StarAward = v.StarAward
			chapterLst = append(chapterLst, chapterInfo)
		}
	}

	self.Elite.Chapter = chapterLst
}

func (daily_copy *TCopyMoudle) DailyReset() {

	//! 刷新各种数据
	for i, _ := range daily_copy.Daily.CopyLst {
		daily_copy.Daily.CopyLst[i].IsChallenge = false
	}
}

func (self *TCopyMoudle) FamousReset() {

	//! 刷新各种数据
	self.Famous.Times = 0
	for i := 0; i < len(self.Famous.Chapter); i++ {
		self.Famous.Chapter[i].PassedCopy = IntLst{}
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

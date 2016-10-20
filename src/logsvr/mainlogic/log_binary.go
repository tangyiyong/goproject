package mainlogic

import (
	"bufio"
	"fmt"
	"gamelog"
	"os"
	"time"
	"utility"
)

type TBinaryLog struct {
	file     *os.File
	writer   *bufio.Writer
	writeCnt int
	flushCnt int
	dir      string
	svrid    int32
	curday   uint32
}

func (self *TBinaryLog) Start(dir string, svrid int32) bool {
	self.dir = dir
	self.svrid = svrid
	self.file = nil
	self.writer = nil
	timeStr := time.Now().Format("20060102")
	logFileName := fmt.Sprintf("%s/%d_%s.blog", self.dir, self.svrid, timeStr)
	var err error
	self.file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		gamelog.Error("BinaryLog Open File Error : %s", err.Error())
		return false
	}
	self.writer = bufio.NewWriter(self.file)
	return true
}

func (self *TBinaryLog) WriteLog(pdata []byte) {
	self.writer.Write(pdata)

	self.writeCnt++
	if self.writeCnt >= self.flushCnt {
		self.Flush()
	}
}

func (self *TBinaryLog) Close() {
	self.writer.Flush()
	self.file.Close()
}

func (self *TBinaryLog) Flush() {
	self.writer.Flush()
	if utility.GetCurDayByUnix() == self.curday {
		return
	}

	if self.file != nil {
		self.file.Close()
	}

	timeStr := time.Now().Format("20060102")
	logFileName := fmt.Sprintf("%s/%d_%s.blog", self.dir, self.svrid, timeStr)

	var err error
	self.file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		gamelog.Error("BinaryLog Open File Error : %s", err.Error())
		return
	}

	self.writer = bufio.NewWriter(self.file)
}

func (self *TBinaryLog) SetFlushCnt(cnt int) {
	self.flushCnt = cnt
	if cnt <= 0 {
		self.flushCnt = 100
	}
}

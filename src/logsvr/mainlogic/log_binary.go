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
}

func (self *TBinaryLog) Start() bool {
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
}

func (self *TBinaryLog) SetFlushCnt(cnt int) {
	self.flushCnt = cnt
	if cnt <= 0 {
		self.flushCnt = 100
	}
}

func CreateBinaryFile(name string, svrid int32) *TBinaryLog {
	var err error = nil
	timeStr := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("%slog\\%s_svr%d_%s.blog", utility.GetCurrPath(), name, svrid, timeStr)

	var blog TBinaryLog
	blog.file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		gamelog.Error("CreateBinaryFile Error : %s", err.Error())
		return nil
	}

	blog.writer = bufio.NewWriter(blog.file)
	return &blog
}

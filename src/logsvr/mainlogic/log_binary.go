package mainlogic

import (
	"bufio"
	"gamelog"
	"os"
	"time"
	"utility"
)

type TBinaryLog struct {
	file   *os.File
	writer *bufio.Writer
}

func (self *TBinaryLog) Start() bool {
	return true
}

func (self *TBinaryLog) WriteLog(pdata []byte) {
	self.writer.Write(pdata)
}

func (self *TBinaryLog) Close() {
	self.file.Close()
}

func (self *TBinaryLog) Flush() {
	self.writer.Flush()
}

func CreateBinaryFile(name string) *TBinaryLog {
	var err error = nil
	timeStr := time.Now().Format("20060102_150405")
	logFileName := utility.GetCurrPath() + "log\\" + name + "_" + timeStr + ".blog"

	var blog TBinaryLog
	blog.file, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		gamelog.Error("CreateBinaryFile Error : %s", err.Error())
		return nil
	}

	blog.writer = bufio.NewWriter(blog.file)
	return &blog
}

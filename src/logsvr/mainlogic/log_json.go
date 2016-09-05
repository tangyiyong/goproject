package mainlogic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gamelog"
	"os"
	"time"
	"utility"
)

type M map[string]interface{}
type TJsonLog struct {
	file *os.File
	wr   *bufio.Writer
	json *json.Encoder
}

func NewJsonLog(name string, svrid int32) *TJsonLog {
	var err error = nil
	timeStr := time.Now().Format("20060102_150405")
	fullName := fmt.Sprintf("%slog\\%s_svr%d_%s.blog", utility.GetCurrPath(), name, svrid, timeStr)

	log := new(TJsonLog)
	log.file, err = os.OpenFile(fullName, os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
		gamelog.Error("JsonLog OpenFile:%s", err.Error())
		return nil
	}
	log.wr = bufio.NewWriterSize(log.file, 1024)
	log.json = json.NewEncoder(log.wr)

	return log
}
func (self *TJsonLog) Close() {
	self.wr.Flush()
	self.file.Close()
}
func (self *TJsonLog) Flush() {
	self.wr.Flush()
}

// WriteLog(M{"a":1, "b":Struct{233,"zhoumf"}})
func (self *TJsonLog) WriteLog(data M) {
	self.json.Encode(data)
}

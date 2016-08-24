package utility

import (
	"crypto/md5"
	"encoding/csv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type StrError struct {
	Str string
	Err error
}

func (self *StrError) Error() string {
	if self.Err == nil {
		return self.Str
	} else {
		return self.Str + " " + self.Err.Error()
	}
}

var G_CurPath string
var G_CurCsvPath string

/*获取当前文件执行的路径*/
func GetCurrPath() string {
	if len(G_CurPath) <= 0 {
		file, _ := exec.LookPath(os.Args[0])
		G_CurPath, _ = filepath.Abs(file)
		G_CurPath = string(G_CurPath[0 : 1+strings.LastIndex(G_CurPath, "\\")])
	}

	return G_CurPath
}
func GetCurrCsvPath() string {
	if len(G_CurCsvPath) <= 0 {
		file, _ := exec.LookPath(os.Args[0])
		G_CurCsvPath, _ = filepath.Abs(file)
		G_CurCsvPath = string(G_CurCsvPath[0 : 1+strings.LastIndex(G_CurCsvPath, "\\")])
		G_CurCsvPath += "csv/"
	}
	return G_CurCsvPath
}
func GetCurrPath2() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	return ret
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	return true
}
func LoadCsv(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fstate, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fstate.IsDir() == true {
		return nil, &StrError{"LoadCsv is dir!", nil}
	}

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func GetCurDay() uint32 {
	day := uint32(time.Now().Day())
	year, week := time.Now().ISOWeek()
	var curday = uint32(year)
	curday = curday << 8
	curday += uint32(week)
	curday = curday << 8
	curday += day
	return curday
}

// Activity系统中：
// 今天是奇数则 data[0]昨日积分、data[1]今日积分；偶数反之
func GetCurDayMod() int {
	day := time.Now().Unix() / 86400

	return int(day % 2)
}

func IsSameDay(day uint32) bool {
	if GetCurDay() != day {
		return false
	}

	return true
}

func IsSameWeek(day uint32) bool {
	nowYear, nowWeek := time.Now().ISOWeek()
	dayYear2, dayWeek := day&0xffff0000>>16, day&0x0000ff00>>8

	if uint32(nowYear) == dayYear2 && uint32(nowWeek) == dayWeek {
		return true
	}

	return false
}

func TestBit(value int, nPos uint8) bool {
	nRet := value & (1 << (31 - nPos))
	return nRet != 0
}

func SetBit(value int, nPos uint8) int {
	value |= 1 << (31 - nPos)
	return value
}

func MsgDataCheck(buffer []byte) bool {
	Lenth := len(buffer)
	if Lenth <= 16 {
		return false
	}
	retmd5 := md5.Sum(buffer[:Lenth-16])
	Lenth -= 16
	for i := 0; i < 16; i++ {
		if retmd5[i] != buffer[Lenth+i] {
			return false
		}
	}

	return true
}

func GetCurDayByUnix() uint32 {
	day := time.Now().Unix() / 86400
	return uint32(day)
}

package gamedata

import (
	"encoding/csv"
	"fmt"
	"gamelog"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"utility"
)

type TDataParser struct {
	OnInit      func(total int) bool
	OnParseData func(rs *RecordSet)
	OnFinish    func() bool
}

var DataParserMap = make(map[string]TDataParser)
var G_ReflectParserMap map[string]interface{}
var NullParser = TDataParser{nil, nil, nil}

func LoadGameData() {
	InitDataParser()
	InitReflectParser()
	LoadAllFiles()
	FinishWakeLevelParser()
}

type RecordSet struct {
	Values []string
	colmap map[string]int
}

func (rs *RecordSet) GetFieldInt(name string) int {
	nIndex, ok := rs.colmap[name]
	if !ok {
		panic(fmt.Sprintf("field: %-10s does not exist", name))
	}

	return CheckAtoiName(rs.Values[nIndex], name)
}
func (rs *RecordSet) GetFieldBool(name string) bool {
	nIndex, ok := rs.colmap[name]
	if !ok {
		panic(fmt.Sprintf("field: %-10s does not exist", name))
	}

	return CheckAtoBool(rs.Values[nIndex])
}
func (rs *RecordSet) GetFieldFloat(name string) float32 {
	nIndex, ok := rs.colmap[name]
	if !ok {
		panic(fmt.Sprintf("field: %-10s does not exist", name))
	}

	return CheckAtofloat(rs.Values[nIndex])
}
func (rs *RecordSet) GetFieldString(name string) string {
	nIndex, ok := rs.colmap[name]
	if !ok {
		panic(fmt.Sprintf("field: %-10s does not exist", name))
	}

	return rs.Values[nIndex]
}
func (rs *RecordSet) GetFieldItems(name string) []ST_ItemData {
	nIndex, ok := rs.colmap[name]
	if !ok {
		panic(fmt.Sprintf("field: %-10s does not exist", name))
	}

	return ParseStringToItem(rs.Values[nIndex])
}

// 格式：(id1|num1)(id2|num2)
func ParseStringToItem(str string) []ST_ItemData {
	sFix := strings.Trim(str, "()")
	slice := strings.Split(sFix, ")(")
	items := make([]ST_ItemData, len(slice))
	for i, v := range slice {
		pv := strings.Split(v, "|")
		if len(pv) != 2 {
			gamelog.Error("ParseStringToItem : %s", str)
			return items
		}
		items[i].ItemID = CheckAtoiName(pv[0], pv[0])
		items[i].ItemNum = CheckAtoiName(pv[1], pv[1])
	}
	return items
}
func ParseStringToPair(str string) []IntPair {
	sFix := strings.Trim(str, "()")
	slice := strings.Split(sFix, ")(")
	items := make([]IntPair, len(slice))
	for i, v := range slice {
		pv := strings.Split(v, "|")
		if len(pv) != 2 {
			gamelog.Error("ParseStringToPair : %s", str)
			return items
		}
		items[i].ID = CheckAtoiName(pv[0], pv[0])
		items[i].Cnt = CheckAtoiName(pv[1], pv[1])
	}
	return items
}

// 格式：32400|43200|64800|75600
func ParseStringToIntArray(str string) []int {
	slice := strings.Split(str, "|")
	nums := make([]int, len(slice))
	for i, v := range slice {
		nums[i] = CheckAtoiName(v, v)
	}
	return nums
}
func ParseStringToByteArray(str string) []byte {
	slice := strings.Split(str, "|")
	nums := make([]byte, len(slice))
	for i, v := range slice {
		nums[i] = byte(CheckAtoiName(v, v))
	}
	return nums
}

func CheckAtoi(s string, nindex int) int {
	if len(s) <= 0 {
		panic(fmt.Sprintf("field: %-10d Is Empty", nindex))
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("field: %-10d text can't convert to int", nindex))
	}

	return i
}

func CheckAtoiName(s string, name string) int {
	if len(s) <= 0 {
		panic(fmt.Sprintf("field: %-10s is empty", name))
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("field: %-10s text can't convert to int", name))
	}

	return i
}

func CheckAtofloat(s string) float32 {
	if len(s) <= 0 {
		panic("Error empty field!!")
		return 0
	}

	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err.Error())
	}

	return float32(f)
}

func CheckAtoBool(s string) bool {
	if len(s) <= 0 {
		panic("Error empty field!!")
		return false
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err.Error())
	}
	return b
}

func LoadAllFiles() {
	pattern := utility.GetCurrCsvPath() + "*.csv"
	files, err := filepath.Glob(pattern)
	if err != nil {
		gamelog.Error("LoadAllFiles error : %s", err.Error())
	}

	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			gamelog.Error("LoadAllFiles error : %s", err.Error())
		}

		LoadOneFile(file)
		file.Close()
	}
}

func ReloadOneFile(tbname string) bool {
	fn := utility.GetCurrCsvPath() + tbname + ".csv"
	file, err := os.Open(fn)
	if err != nil {
		gamelog.Error("ReloadOneFile error : %s", err.Error())
		return false
	}
	LoadOneFile(file)
	file.Close()
	return true
}

func LoadOneFile(file *os.File) {
	// 处理表名
	fstate, err := file.Stat()
	if err != nil {
		gamelog.Error("LoadOneFile error : %s", err.Error())
		return
	}

	if fstate.IsDir() == true {
		return
	}

	tblname := strings.TrimSuffix(fstate.Name(), path.Ext(file.Name()))
	var mapInterface interface{}
	DataParser, ok := DataParserMap[tblname]
	if !ok {
		mapInterface, ok = G_ReflectParserMap[tblname]
		if !ok {
			gamelog.Error("table: %-30s need a parser!!", tblname)
		}
		return
	}

	//明确表示不需要解析的表
	if DataParser.OnInit == nil {
		return
	}

	csv_reader := csv.NewReader(file)
	records, err := csv_reader.ReadAll()
	if err != nil {
		gamelog.Error("LoadOneFile %s error : %s", fstate.Name(), err.Error())
		return
	}

	// 是否为空档
	if len(records) == 0 {
		gamelog.Error("LoadOneFile %s empty csv file ", fstate.Name())
		return
	}

	if mapInterface != nil {
		ParseRefCsv(records, mapInterface)
		return
	}

	nCount := len(records) - 1 //实际有效数据的长度
	if !DataParser.OnInit(nCount) {
		gamelog.Error("table: %-30s OnInitParser error!!", tblname)
		return
	}

	// 解析列
	colCount := len(records[0])
	ColMap := make(map[string]int)
	for i := 0; i < colCount; i++ {
		ColMap[records[0][i]] = i
	}

	var line int
	defer func() {
		if err := recover(); err != nil {
			gamelog.Error("table: %-30s line: %-3d %s", tblname, line+1, err)
		}
	}()

	var rs RecordSet
	rs.colmap = ColMap
	// 记录数据, 第一行为表头，因此从第二行开始
	for line = 1; line <= nCount; line++ {
		rs.Values = records[line]
		DataParser.OnParseData(&rs)
	}

	if DataParser.OnFinish != nil {
		DataParser.OnFinish()
	}

	ColMap = nil
}

package gamedata

import (
	"fmt"
	"reflect"
	"strings"
)

func ParseRefCsv(records [][]string, ptr interface{}) {
	switch reflect.TypeOf(ptr).Elem().Kind() {
	case reflect.Map:
		{
			ParseRefCsvByMap(records, ptr)
		}
	case reflect.Slice:
		{
			ParseRefCsvBySlice(records, ptr)
		}
	default:
		{
			panic(fmt.Sprintf("Csv Type Error: TypeName:%s", reflect.TypeOf(ptr).Elem().String()))
		}
	}
}
func ParseRefCsvByMap(records [][]string, pMap interface{}) {
	table := reflect.ValueOf(pMap).Elem()
	typ := table.Type().Elem().Elem() // map内保存的指针，第二次Elem()得到所指对象类型

	total, idx := GetRecordsValidCnt(records), 0
	slice := reflect.MakeSlice(reflect.SliceOf(typ), total, total) // 避免多次new对象，直接new数组，拆开用

	bParsedName, nilFlag := false, int64(0)
	for _, v := range records {
		if strings.Index(v[0], "#") == -1 { // "#"起始的不读
			if !bParsedName {
				nilFlag = parseRefName(v)
				bParsedName = true
			} else {
				// data := reflect.New(typ).Elem()
				data := slice.Index(idx)
				idx++
				parseRefData(v, nilFlag, data)
				table.SetMapIndex(data.Field(0), data.Addr())
			}
		}
	}
}
func ParseRefCsvBySlice(records [][]string, pSlice interface{}) { // slice可减少对象数量，降低gc
	slice := reflect.ValueOf(pSlice).Elem() // 这里slice是nil
	typ := reflect.TypeOf(pSlice).Elem()

	// 表的数组，从1起始
	idx := 1
	total := GetRecordsValidCnt(records) + 1
	slice.Set(reflect.MakeSlice(typ, total, total))

	bParsedName, nilFlag := false, int64(0)
	for _, v := range records {
		if strings.Index(v[0], "#") == -1 { // "#"起始的不读
			if !bParsedName {
				nilFlag = parseRefName(v)
				bParsedName = true
			} else {
				data := slice.Index(idx)
				idx++
				parseRefData(v, nilFlag, data)
			}
		}
	}
}
func parseRefName(record []string) (ret int64) { // 不读的列：没命名/前缀"(c)"
	length := len(record)
	if length > 64 {
		panic(fmt.Sprintf("csv column is over to 64 !!!"))
	}
	for i := 0; i < length; i++ {
		if record[i] == "" || strings.Index(record[i], "(c)") == 0 {
			ret |= (1 << uint(i))
		}
	}
	return ret
}
func parseRefData(record []string, nilFlag int64, data reflect.Value) {
	idx := 0
	for i, s := range record {
		if nilFlag&(1<<uint(i)) > 0 { // 跳过没命名的列
			continue
		}

		field := data.Field(idx)
		idx++

		if s == "" { // 没填的就不必解析了，跳过，idx还是要自增哟
			continue
		}

		switch field.Kind() {
		case reflect.Int:
			{
				field.SetInt(int64(CheckAtoiName(s, s)))
			}
		case reflect.String:
			{
				field.SetString(s)
			}
		case reflect.Slice:
			{
				switch field.Type().Elem().Kind() {
				case reflect.Int:
					{
						vec := ParseStringToIntArray(s)
						field.Set(reflect.ValueOf(vec))
					}
				case reflect.Struct:
					{
						vec := ParseStringToPair(s)
						field.Set(reflect.ValueOf(vec))
					}
				default:
					{
						panic(fmt.Sprintf("Csv Type Error: TypeName:%s", data.Field(i).Type().String()))
					}
				}
			}
		default:
			{
				panic(fmt.Sprintf("Csv Type Error: TypeName:%s", data.Field(i).Type().String()))
			}
		}
	}
}

func GetRecordsValidCnt(records [][]string) (ret int) {
	for _, v := range records {
		if strings.Index(v[0], "#") == -1 { // "#"起始的不读
			ret++
		}
	}
	return ret
}

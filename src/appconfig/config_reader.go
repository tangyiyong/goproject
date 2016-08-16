package appconfig

import (
	"bufio"
	"os"
	"strings"
	"utility"
)

var CONFIG_PATH = utility.GetCurrPath() + "config.ini"

func LoadConfig() bool {
	bRet := parseConfigFile(CONFIG_PATH)
	if bRet {
		bRet = initGlobalVar()
	}

	return bRet
}

func parseConfigFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		panic("parseConfigFile Failed Error :" + err.Error())
		return false
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) <= 2 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		slice := strings.Split(line, "=")
		if len(slice) != 2 {
			continue
		}

		ParseConfigValue(strings.TrimSpace(slice[0]), strings.TrimSpace(slice[1]))
	}

	return true
}

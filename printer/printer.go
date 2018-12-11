package printer

import (
	"bytes"
	"fmt"
	"encoding/json"

	"github.com/zhengcf/goutil/config"
	log "github.com/sirupsen/logrus"
)

// Version information.
var (
	AppName      = "Node"

	AppReleaseVersion = "None"
	AppBuildTS   = "None"
	AppGitHash   = "None"
	AppGitBranch = "None"
	GoVersion     = "None"
)

// PrintTiDBInfo prints the TiDB version information.
func PrintAppInfo() {
	log.Infof("Welcome to %s.", AppName)
	log.Infof("Git Commit Hash: %s", AppGitHash)
	log.Infof("Git Branch: %s", AppGitBranch)
	log.Infof("UTC Build Time:  %s", AppBuildTS)
	log.Infof("GoVersion:  %s", GoVersion)
	configJSON, err := json.Marshal(config.GetGlobalConfig())
	if err != nil {
		panic(err)
	}
	log.Infof("Config: %s", configJSON)
}

// GetTiDBInfo returns the git hash and build time of this tidb-server binary.
func GetAppInfo() string {
	return fmt.Sprintf("Release Version: %s\n"+
		"Git Commit Hash: %s\n"+
		"Git Branch: %s\n"+
		"UTC Build Time: %s\n"+
		"GoVersion: %s\n",
		AppReleaseVersion,
		AppGitHash,
		AppGitBranch,
		AppBuildTS,
		GoVersion)
}

// checkValidity checks whether cols and every data have the same length.
func checkValidity(cols []string, datas [][]string) bool {
	colLen := len(cols)
	if len(datas) == 0 || colLen == 0 {
		return false
	}

	for _, data := range datas {
		if colLen != len(data) {
			return false
		}
	}

	return true
}

func getMaxColLen(cols []string, datas [][]string) []int {
	maxColLen := make([]int, len(cols))
	for i, col := range cols {
		maxColLen[i] = len(col)
	}

	for _, data := range datas {
		for i, v := range data {
			if len(v) > maxColLen[i] {
				maxColLen[i] = len(v)
			}
		}
	}

	return maxColLen
}

func getPrintDivLine(maxColLen []int) []byte {
	var value = make([]byte, 0)
	for _, v := range maxColLen {
		value = append(value, '+')
		value = append(value, bytes.Repeat([]byte{'-'}, v+2)...)
	}
	value = append(value, '+')
	value = append(value, '\n')
	return value
}

func getPrintCol(cols []string, maxColLen []int) []byte {
	var value = make([]byte, 0)
	for i, v := range cols {
		value = append(value, '|')
		value = append(value, ' ')
		value = append(value, []byte(v)...)
		value = append(value, bytes.Repeat([]byte{' '}, maxColLen[i]+1-len(v))...)
	}
	value = append(value, '|')
	value = append(value, '\n')
	return value
}

func getPrintRow(data []string, maxColLen []int) []byte {
	var value = make([]byte, 0)
	for i, v := range data {
		value = append(value, '|')
		value = append(value, ' ')
		value = append(value, []byte(v)...)
		value = append(value, bytes.Repeat([]byte{' '}, maxColLen[i]+1-len(v))...)
	}
	value = append(value, '|')
	value = append(value, '\n')
	return value
}

func getPrintRows(datas [][]string, maxColLen []int) []byte {
	var value = make([]byte, 0)
	for _, data := range datas {
		value = append(value, getPrintRow(data, maxColLen)...)
	}
	return value
}

// GetPrintResult gets a result with a formatted string.
func GetPrintResult(cols []string, datas [][]string) (string, bool) {
	if !checkValidity(cols, datas) {
		return "", false
	}

	var value = make([]byte, 0)
	maxColLen := getMaxColLen(cols, datas)

	value = append(value, getPrintDivLine(maxColLen)...)
	value = append(value, getPrintCol(cols, maxColLen)...)
	value = append(value, getPrintDivLine(maxColLen)...)
	value = append(value, getPrintRows(datas, maxColLen)...)
	value = append(value, getPrintDivLine(maxColLen)...)
	return string(value), true
}

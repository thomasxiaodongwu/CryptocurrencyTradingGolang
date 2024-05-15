/*
 * @Author: xwu
 * @Date: 2021-12-26 18:42:10
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 16:10:35
 */
package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"winter/messages"
)

func String2int(x string) int {
	var ans int = 0
	ans, err := strconv.Atoi(x)
	if err != nil {
		logger.Info("string2int failed: " + err.Error())
	}
	return ans
}

func String2int64(x string) int64 {
	var ans int = 0
	ans, err := strconv.Atoi(x)
	if err != nil {
		logger.Info("string2int failed: " + err.Error())
	}
	return int64(ans)
}

func String2float(x string) float64 {
	var ans float64 = 0
	ans, err := strconv.ParseFloat(x, 64)
	if err != nil {
		logger.Info("string2float64 failed.")
	}
	return ans
}

func String2localtime(x string, format string) int64 {
	// "2006-01-02T15:04:05.000Z"
	stamp, err := time.ParseInLocation(format, x, time.Local)
	if err != nil {
		logger.Info("string2localtime " + err.Error())
	}
	return stamp.Unix()
}

func SymbolFormat(A string, B string, exch string) string {
	switch exch {
	case "Okex":
		return Strcat(strings.ToUpper(A), "-", strings.ToUpper(B))
	default:
		return "Error"
	}
}

func SubMsgFormat(symbols []string, msgkind string, exch string) []byte {
	switch exch {
	// 交易所
	case "Okex":
		ans := messages.NewOkexSubMsg(symbols, msgkind)
		res, err := json.Marshal(ans)
		if err != nil {
			logger.Info("Json Marshal OkexSubMsg Failed.")
		}
		return res
	case "Binance":
		ans := messages.NewBinanceSubMsg(symbols, msgkind)
		res, err := json.Marshal(ans)
		if err != nil {
			logger.Info("Json Marshal OkexSubMsg Failed.")
		}
		fmt.Println(string(res))
		return res
	default:
		return nil
	}
}

func DropStrTail(s string, pattern string) string {
	a := len(pattern)
	if s[len(s)-a:] == pattern {
		return s[:len(s)-a]
	}
	return s
}

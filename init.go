package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"winter/Logger"
	"winter/global"
	"winter/params"
	"winter/utils"

	"go.uber.org/zap"
)

var logger *zap.Logger

func GetAlphaUid() (map[int](map[string]int), int) {
	// 哪个交易所的什么币种就可以确定他在数组中是第几个
	// 每个交易所
	ans := make(map[int](map[string]int))

	count := 0

	instr_alpha := reflect.ValueOf(global.AggParameters.Data)
	for ExchangeName, ExchangeId := range global.ExIDMap {
		ans[ExchangeId] = make(map[string]int) // 判断是否是重复的

		_p := instr_alpha.FieldByName(ExchangeName).Interface().(params.DataParams)

		for i := range _p.Subscribe_symbols {
			symbol := strings.Split(_p.Subscribe_symbols[i], "_")
			if _, ok := ans[ExchangeId][symbol[0]]; !ok {
				ans[ExchangeId][symbol[0]] = count
				count += 1
			}
		}
	}

	fmt.Println("AlphaUid:\n", ans)

	return ans, count
}

func GetHftUid() (map[int](map[string](map[string]int)), int) {
	/*
		哪个交易所的哪个币种的哪个策略
	*/
	ans := make(map[int](map[string](map[string]int)))

	var count int = 0

	instr_trader := reflect.ValueOf(global.AggParameters.Trader)
	for ExchangeName, ExchangeId := range global.ExIDMap {
		ans[ExchangeId] = make(map[string](map[string]int))
		_instr := instr_trader.FieldByName(ExchangeName).Interface().(params.TradeParams)

		for i := range _instr.HftStrategyNames {
			ans[ExchangeId][_instr.HftStrategyNames[i]] = make(map[string]int)
			for j := range _instr.HftStrategyParams[i] {
				symbol := _instr.HftStrategyParams[i][j].Symbol
				if ExchangeName == "Binance" {
					symbol = strings.ToUpper(symbol)
				}
				ans[ExchangeId][_instr.HftStrategyNames[i]][symbol] = count
				count += 1
			}
		}
	}

	fmt.Println("HftUid:\n", ans)

	return ans, count
}

func GetTraderParams() []params.TradeParams {
	ans := make([]params.TradeParams, len(global.ExIDList))

	instr_trader := reflect.ValueOf(global.AggParameters.Trader)
	for i := range global.ExIDList {
		ExchangeName := global.ExIDList[i]
		_instr := instr_trader.FieldByName(ExchangeName).Interface().(params.TradeParams)
		if ExchangeName == "Binance" {
			for i := range _instr.HftStrategyParams[0] {
				_instr.HftStrategyParams[0][i].Symbol = strings.ToUpper(_instr.HftStrategyParams[0][i].Symbol)
			}
		}
		ans[i] = _instr
	}

	return ans
}

func InitInstrument() map[string]params.Instrument {
	ans := make(map[string]params.Instrument)

	file, err := os.Open("Configuration/symbol.csv")
	if err != nil {
		logger.Fatal("open Configuration/symbol.csv failed.")
	}
	defer file.Close()
	rd := bufio.NewReader(file)

	var first_saw bool = true
	var columns []string
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		if first_saw {
			columns = strings.Split(line, ",")
			first_saw = false
		} else {
			values := strings.Split(line, ",")
			if len(values) != len(columns) {
				logger.Fatal("parse symbol.csv failed.")
			}
			var symbol string
			var minsize float64
			var tick int64
			// var tick_str string
			for i, v := range columns {
				if v == "ctVal" {
					minsize, err = strconv.ParseFloat(values[i], 64)
					if err != nil {
						logger.Fatal("parse ctVal in symbol.csv failed.")
					}
				} else if v == "tickSz" {
					// tick_str = values[i]
					_tick, err := strconv.ParseFloat(values[i], 64)
					if err != nil {
						logger.Fatal("parse tickSz in symbol.csv failed.")
					}
					// if symbol == "BTC-USDT-SWAP" {
					// 	fmt.Println(_tick, " ", math.Round(math.Log10(_tick)), " ", int64(-1*math.Log10(_tick)))
					// }
					tick = int64(math.Round((-1 * math.Log10(_tick))))
				} else if v == "instId" {
					symbol = values[i]
				}
			}

			// fmt.Println(symbol, " ", tick, " ", tick_str, " | ", minsize)

			ans[symbol] = params.Instrument{Symbol: symbol, Size: minsize, Tick: tick, II: -1}
		}
	}

	for symbol, instr := range ans {
		if ii, ok := global.AlphaUid[1][symbol]; ok {
			instr.II = ii
			ans[symbol] = instr
		}
	}

	// logger.Fatal("parse symbol.csv failed.")
	return ans
}

func checkGlobal() {
	// 输出global关键配置
	switch global.Mode {
	case "test":
		logger.Info("Version: " + global.Version + ", Mode: " + global.Mode)
	case "backtest":
		logger.Info("Version: " + global.Version + ", Mode: " + global.Mode + ", Trigger: " + global.Trigger)
		if global.Dumper {
			logger.Info("factors dumper: on")
		} else {
			logger.Info("factors dumper: off")
		}
		if global.IsDebugChecksum {
			logger.Info("debug checksum: on")
		} else {
			logger.Info("debug checksum: off")
		}
		if global.IsSimTrade {
			logger.Info("simulate trade: on")
		} else {
			logger.Info("simulate trade: off")
		}
		if global.LGBPredict {
			logger.Info("calc lgb: on")
		} else {
			logger.Info("calc lgb: off")
		}
		if global.IsLogCpu {
			logger.Info("record performace of system: on")
		} else {
			logger.Info("record performace of system: off")
		}
		logger.Info("The window of return: " + strconv.FormatInt(global.ReturnInterval, 10))
	case "hft":
		logger.Info("Version: " + global.Version + ", Mode: " + global.Mode + ", Trigger: " + global.Trigger)
		if global.USING_PROXY {
			logger.Info("proxy: on")
		} else {
			logger.Info("proxy: off")
		}
		if global.IsDebugChecksum {
			logger.Info("debug checksum: on")
		} else {
			logger.Info("debug checksum: off")
		}
		if global.LGBPredict {
			logger.Info("calc lgb: on")
		} else {
			logger.Info("calc lgb: off")
		}
		if global.IsLogCpu {
			logger.Info("record performace of system: on")
		} else {
			logger.Info("record performace of system: off")
		}
		logger.Info("Hft Name: " + global.HftName + ", trade amount: " + strconv.FormatFloat(global.Amount, 'f', 2, 64))
	case "recData":
		logger.Info("Version: " + global.Version + ", Mode: " + global.Mode)
		if global.USING_PROXY {
			logger.Info("proxy: on")
		} else {
			logger.Info("proxy: off")
		}
		if global.IsLogCpu {
			logger.Info("record performace of system: on")
		} else {
			logger.Info("record performace of system: off")
		}
		logger.Info("the path of data: " + global.Data_Dump_Path)
	}
}

func init() {
	// 这里针对不同模式有不同的log方便查看错误
	// 这里只针对global.go的配置问题
	logger = Logger.BaseLogger

	// 输出global关键配置
	checkGlobal()

	switch global.Mode {
	case "test":
		global.AggParameters = utils.LoadConfig(global.ConfigurationTest)
	case "backtest":
		global.AggParameters = utils.LoadConfig(global.ConfigurationBacktest)
	case "hft":
		global.AggParameters = utils.LoadConfig(global.ConfigurationHft)
	case "recData":
		global.AggParameters = utils.LoadConfig(global.ConfigurationRecData)
	}

	// 初始化字典的对应全局变量
	// v0.7版本之后的系统需要有两个map其中一个map是用于process的,
	// 另一个map是用于trader那边的
	// 因为这里的所有操作只需要进行一次,所以可以使用反射

	// 哪个交易所的什么币种就可以确定他在数组中是第几个
	global.AlphaUid, global.FMCount = GetAlphaUid()

	global.HftUid, global.PMCount = GetHftUid()

	global.AggTradeParams = GetTraderParams()
	// for i := range global.AggTradeParams {
	// 	global.AggParameters
	// }

	global.Instruments = InitInstrument()

	global.Okex_Kind_Decimals = make(map[string]([2]int))
	for k, v := range global.Instruments {
		global.Okex_Kind_Decimals[k] = [2]int{int(v.Tick), 0}
	}

	global.HistData.Init(global.AlphaUid, int64(global.LOG_SIZE))

	logger.Info("Initial Global Variances Finished")
	logger.Info("TradeSysGo Starts at " + time.Now().Format("2006-01-02 15:04:05"))

	global.StartTime = time.Now().UnixMilli()
}

// ans["BTC-USDT-SWAP"] = params.Instrument{
// 	Symbol: "BTC-USDT-SWAP",
// 	Size:   0.01,
// 	Tick:   1,
// }
// ans["ETH-USDT-SWAP"] = params.Instrument{
// 	Symbol: "ETH-USDT-SWAP",
// 	Size:   0.1,
// 	Tick:   2,
// }

// ans["CRV-USDT-SWAP"] = params.Instrument{
// 	Symbol: "CRV-USDT-SWAP",
// 	Size:   1,
// 	Tick:   4,
// }

// ans["AXS-USDT-SWAP"] = params.Instrument{
// 	Symbol: "AXS-USDT-SWAP",
// 	Size:   0.1,
// 	Tick:   2,
// }

// ans["DOGE-USDT-SWAP"] = params.Instrument{
// 	Symbol: "DOGE-USDT-SWAP",
// 	Size:   1000,
// 	Tick:   5,
// }

// ans["FTM-USDT-SWAP"] = params.Instrument{
// 	Symbol: "FTM-USDT-SWAP",
// 	Size:   10,
// 	Tick:   4,
// }

// ans["SAND-USDT-SWAP"] = params.Instrument{
// 	Symbol: "SAND-USDT-SWAP",
// 	Size:   10,
// 	Tick:   4,
// }

// ans["LINK-USDT-SWAP"] = params.Instrument{
// 	Symbol: "LINK-USDT-SWAP",
// 	Size:   1,
// 	Tick:   3,
// }

// ans["CELO-USDT-SWAP"] = params.Instrument{
// 	Symbol: "CELO-USDT-SWAP",
// 	Size:   1,
// 	Tick:   4,
// }

// ans["SOL-USDT-SWAP"] = params.Instrument{
// 	Symbol: "SOL-USDT-SWAP",
// 	Size:   1,
// 	Tick:   2,
// }

// ans["MATIC-USDT-SWAP"] = params.Instrument{
// 	Symbol: "MATIC-USDT-SWAP",
// 	Size:   1,
// 	Tick:   2,
// }

// ans["LUNA-USDT-SWAP"] = params.Instrument{
// 	Symbol: "LUNA-USDT-SWAP",
// 	Size:   1,
// 	Tick:   4,
// }

// ans["WAVES-USDT-SWAP"] = params.Instrument{
// 	Symbol: "WAVES-USDT-SWAP",
// 	Size:   1,
// 	Tick:   3,
// }

// ans["OP-USDT-SWAP"] = params.Instrument{
// 	Symbol: "OP-USDT-SWAP",
// 	Size:   1,
// 	Tick:   4,
// }

// ans["GMT-USDT-SWAP"] = params.Instrument{
// 	Symbol: "GMT-USDT-SWAP",
// 	Size:   1,
// 	Tick:   4,
// }

// ans["ADA-USDT-SWAP"] = params.Instrument{
// 	Symbol: "ADA-USDT-SWAP",
// 	Size:   100,
// 	Tick:   5,
// }

// ans["PEOPLE-USDT-SWAP"] = params.Instrument{
// 	Symbol: "PEOPLE-USDT-SWAP",
// 	Size:   100,
// 	Tick:   5,
// }

// ans["TRX-USDT-SWAP"] = params.Instrument{
// 	Symbol: "TRX-USDT-SWAP",
// 	Size:   1000,
// 	Tick:   5,
// }

// ans["DOT-USDT-SWAP"] = params.Instrument{
// 	Symbol: "DOT-USDT-SWAP",
// 	Size:   1,
// 	Tick:   3,
// }

// ans["LTC-USDT-SWAP"] = params.Instrument{
// 	Symbol: "LTC-USDT-SWAP",
// 	Size:   1,
// 	Tick:   2,
// }

// ans["APE-USDT-SWAP"] = params.Instrument{
// 	Symbol: "APE-USDT-SWAP",
// 	Size:   0.1,
// 	Tick:   4,
// }

// ans["ENS-USDT-SWAP"] = params.Instrument{
// 	Symbol: "ENS-USDT-SWAP",
// 	Size:   0.1,
// 	Tick:   3,
// }

// ans["LOOKS-USDT-SWAP"] = params.Instrument{
// 	Symbol: "LOOKS-USDT-SWAP",
// 	Size:   1,
// 	Tick:   4,
// }

// ans["AVAX-USDT-SWAP"] = params.Instrument{
// 	Symbol: "AVAX-USDT-SWAP",
// 	Size:   1,
// 	Tick:   2,
// }

// ans["DYDX-USDT-SWAP"] = params.Instrument{
// 	Symbol: "DYDX-USDT-SWAP",
// 	Size:   1,
// 	Tick:   3,
// }

// ans["GALA-USDT-SWAP"] = params.Instrument{
// 	Symbol: "GALA-USDT-SWAP",
// 	Size:   10,
// 	Tick:   5,
// }

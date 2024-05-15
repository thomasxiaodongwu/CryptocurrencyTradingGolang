/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:43
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-11 11:58:11
 */
package global

import (
	"winter/messages"
	"winter/params"
)

const (
	// ********************************** 系统的版本号等 **********************************
	// system's version
	Version string = "cliff_20220614_v0.8"

	// ********************************** 系统的控制常数 **********************************
	// Mode: test, testDelay, backtest, hft, recData
	Mode string = "backtest"

	// SubMode: follow all这个参数已经被弃用，被Trigger完全取代，但是删除可能会有影响，后面整理
	SubMode string = "all"

	// Trigger: all, follow, time
	Trigger string = "follow"

	// 是否dump因子
	Dumper bool = false

	// 是否开代理，一般来说大陆本地需要开，服务器上不需要
	USING_PROXY bool = false

	// 是否检查okx的checksum，实盘的时候不开，开了会很慢
	IsDebugChecksum bool = false

	// 是否会进行回测交易，还没有和实盘进行对比
	IsSimTrade bool = false

	// ********************************** 需要的各种常数 **********************************
	// 代理的端口
	PROXY_URI string = "http://127.0.0.1:4780"

	// 不同模式下链接的文件  test, backtest, hft, recData
	ConfigurationTest     string = "./Configuration/test.json"
	ConfigurationBacktest string = "./Configuration/backtest.json"
	ConfigurationHft      string = "./Configuration/hft.json"
	ConfigurationRecData  string = "./Configuration/recdata.json"

	// 交易网站相关的常数
	RESTFUL_OKEX_URI           string = `https://www.okx.com/`
	WEBSOCKET_OKEX_URI         string = `wss://ws.okx.com:8443/ws/v5/public`
	WEBSOCKET_OKEX_PRIVATE_URI string = `wss://ws.okx.com:8443/ws/v5/private`

	RESTFUL_BINANCE_URI           string = ``
	WEBSOCKET_BINANCE_URI         string = `wss://fstream.binance.com/ws/ethusdt@aggTrade`
	WEBSOCKET_BINANCE_PRIVATE_URI string = `wss://fstream.binance.com/ws`

	RESTFUL_FTX_URI   string = `https://ftx.com/api`
	WEBSOCKET_FTX_URI string = `wss://ftx.com/ws/`

	Okex_Passphrase string = ``
	Okex_ApiKey     string = ``
	Okex_SecertKey  string = ``

	// 纳秒和秒的兑换比例的常数
	NANO_PER_SEC int64 = 1e9
	NANO_PER_MIN int64 = 60 * 1e9

	LOG_SIZE = 20

	// backtest的时候预测的时间长度，单位是ms
	ReturnInterval int64 = 3000

	// ********************************** 交易模式的常数 **********************************
	// 代理的端口
	Amount  float64 = 125 // 用来计算一次交易多少张
	HftName string  = `hft`

	// 是否运行模型
	LGBPredict bool = false

	// recData模式下，数据放在哪里
	Data_Dump_Path = `/home/ubuntu/program/binance_dydx/TradeSysGo`

	// go tool pprof 测试性能用
	IsLogCpu bool = true

	// 计算因子线程数量，固定必须是1,因为多线程会很慢
	// Num_threads int = 1
)

var (
	// json文件load进来的配置
	AggParameters params.AggParams

	// 交易参数
	AggTradeParams []params.TradeParams

	// 不同交易所不同币种对应的位置
	AlphaUid map[int](map[string]int)

	// 交易所的num,策略的名字,币种的名字,最后对应的key是slice中对应的位置
	HftUid map[int](map[string](map[string]int))

	// 全局变量，可以多线程读取，但是不可以多线程修改
	HistData messages.HistData

	// 特征的数量
	FMCount int

	// 仓位的总数量
	PMCount int

	// 交易所的名字和顺序
	ExIDList = [4]string{"FTX", "Okex", "Binance", "Dydx"}

	// 交易所的map
	ExIDMap = map[string]int{
		"FTX":     0,
		"Okex":    1,
		"Binance": 2,
		"Dydx":    3,
	}

	// 系统开始时间和TickTime
	StartTime int64

	// 从symbol.csv load进来的币种基本信息
	Okex_Kind_Decimals map[string]([2]int)

	// 从symbol.csv load进来的币种基本信息, 功能和上面有点重合，上面的是个过度的产物
	Instruments map[string]params.Instrument

	// backtest模式下，用来看有多少行的，计算效率
	LINE_NUM int64 = 0
)

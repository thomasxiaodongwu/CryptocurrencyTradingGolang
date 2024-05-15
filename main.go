package main

import (
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"
	"winter/clientAlpha"
	"winter/clientData"
	"winter/clientProcess"
	"winter/clientTrader"
	"winter/global"
	"winter/messages"
	"winter/utils"
)

func main() {
	if global.IsLogCpu {
		f, _ := os.OpenFile("cpu.pprof", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	start_time := time.Now()

	switch {
	case global.Mode == "test":
		// go utils.FeiShuMsg("DataClient Test")

		DataClient := clientData.NewBinanceClient()
		var chanMsg chan messages.MsgDataToProcess = make(chan messages.MsgDataToProcess, 10)
		wg := new(sync.WaitGroup)
		wg.Add(2)
		go DataClient.Run(chanMsg, wg)
		go clientData.Fake(chanMsg, wg)
		// go DataClient.RunPing()
		wg.Wait()
	case global.Mode == "testDelay":
		// HftClient := clientTrader.NewTrader()
		// HftClient.RunMonitor_TestNetWorkLag()
		logger.Info("testDelay not implemeted")
	case global.Mode == "Synchronous":
		// go SimOkexDataClient.Run(chanMsg, wg)         // Part Data
		// go Convert.Run(chanMsg, chanStmsg, wg)        // Part Process
		// go AlphaClient.Run(chanStmsg, chanSignal, wg) // Part Alpha
		// go SimHftClient.Run(chanSignal, wg)           // Part Trader

		// SimOkexDataClient := clientData.NewOkexSim()
		// Convert := clientProcess.NewConvert()
		// AlphaClient := clientAlpha.NewAlphaFactory()
		// SimHftClient := clientTrader.NewHftSimClient()

		// chanMsg := SimOkexDataClient.RunSynchronous()       // Part Data
		// chanStmsg := Convert.RunSynchronous(chanMsg)        // Part Process
		// chanSignal := AlphaClient.RunSynchronous(chanStmsg) // Part Alpha
		// SimHftClient.RunSynchronous(chanSignal)             // Part Trader
		logger.Info("Synchronous")
	case global.Mode == "backtest":
		logger.Info("start at " +
			global.AggParameters.Data.Okex.Start_time +
			". end at " +
			global.AggParameters.Data.Okex.End_time +
			".")

		SimOkexDataClient := clientData.NewOkexSim()
		Convert := clientProcess.NewConvert()
		AlphaClient := clientAlpha.NewAlphaFactory()
		SimHftClient := clientTrader.NewHftSimClient()

		chanMsg := make(chan messages.MsgDataToProcess, 30)
		chanStmsg := make(chan messages.AggStMsg, 30)
		chanSignal := make(chan messages.AggSignal, 30)

		// chanMsg := make(chan messages.MsgDataToProcess)
		// chanStmsg := make(chan messages.AggStMsg)
		// chanSignal := make(chan messages.AggSignal)

		wg := new(sync.WaitGroup)
		wg.Add(4)
		go SimOkexDataClient.Run(chanMsg, wg)         // Part Data
		go Convert.Run(chanMsg, chanStmsg, wg)        // Part Process
		go AlphaClient.Run(chanStmsg, chanSignal, wg) // Part Alpha
		go SimHftClient.Run(chanSignal, wg)           // Part Trader
		wg.Wait()

		//utils.FeiShuMsg("Notice: backtest done.")
	case global.Mode == "hft":
		OkexDataClient := clientData.NewOkexClient()
		// FTXDataClient := clientData.NewFTXClient()
		Convert := clientProcess.NewConvert()
		AlphaClient := clientAlpha.NewAlphaFactory()
		HftClient := clientTrader.NewHftClient()

		chanMsg := make(chan messages.MsgDataToProcess, 30)
		chanStmsg := make(chan messages.AggStMsg, 30)
		chanSignal := make(chan messages.AggSignal, 30)
		messageEvent := make(chan string, 102400)
		messageOrders := make(chan string, 102400)
		wg := new(sync.WaitGroup)
		wg.Add(7)
		go OkexDataClient.Run(chanMsg, wg)            // Part Data okex
		go Convert.Run(chanMsg, chanStmsg, wg)        // Part Process
		go AlphaClient.Run(chanStmsg, chanSignal, wg) // Part Alpha

		// Part Trader
		go HftClient.Monitor(messageEvent, messageOrders, wg)
		go HftClient.MessageEvent(messageEvent, wg)
		go HftClient.MessageOrders(messageOrders, wg)
		go HftClient.Run(chanSignal, wg)
		go HftClient.Ping()

		wg.Wait()
	case global.Mode == "recData":
		BinanceDataClient := clientData.NewBinanceClient()
		OkexDataClient := clientData.NewOkexClient()
		var chanMsg chan messages.MsgDataToProcess = make(chan messages.MsgDataToProcess, 10)
		wg := new(sync.WaitGroup)
		wg.Add(3)
		go OkexDataClient.Run(chanMsg, wg)
		go BinanceDataClient.Run(chanMsg, wg)

		go clientData.DumpRawData(chanMsg, wg)
		wg.Wait()
	default:
		logger.Info("Unknown Mode")
	}

	logger.Info(utils.Strcat(
		"TradeSysGo cost time: ",
		time.Since(start_time).String(),
		". all line numbers: ",
		strconv.FormatInt(global.LINE_NUM, 10),
		". cost time of each line: ",
		strconv.FormatFloat(float64(time.Since(start_time).Microseconds())/float64(global.LINE_NUM), 'f', 2, 64),
	))
}

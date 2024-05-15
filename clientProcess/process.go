/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:11
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-02 21:42:40
 */
package clientProcess

import (
	"strings"
	"sync"
	"winter/messages"
	"winter/utils"
)

type Convert struct {
	// 这里的ii是global的AlphaUid
	Is_alive bool

	FTXOrderbook []FTXQteEvt
	FTX_ii       map[string]int

	OkexOrderbook []okexQteEvt
	Okex_ii       map[string]int

	BinanceOrderbook []BinanceQteEvt
	Binance_ii       map[string]int
}

func (con *Convert) Run(chanMsg chan messages.MsgDataToProcess, chanStmsg chan messages.AggStMsg, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var line messages.MsgDataToProcess = <-chanMsg

		var is_end bool = strings.Contains(line.Contents, "all files are done.")
		if is_end {
			chanStmsg <- messages.AggStMsg{Status: 0}
			break
		}

		// 这段代码的目的就是为了给chanStmsg这个channel给出一个需要的值，但是要分不同的交易所,不同的消息种类等等。
		var line_split []string = strings.Split(line.Contents, "|")
		// 这里是旧数据，只有一个|
		// var localtime int64 = int64(utils.String2int(line_split[0]))
		// var msg string = line_split[1]
		// 这里是新数据有两个|
		var localtime int64 = int64(utils.String2int(line_split[0]))
		//fmt.Println(line_split)
		var msg string = line_split[1]

		var stMsg []messages.AggStMsg // 这个就是每个循环过来要输出的结构体

		switch line.ExID {
		case 0: // 0代表FTX交易所
			var is_trades bool = strings.Contains(msg, "trades")
			var is_orderbook bool = strings.Contains(msg, "orderbook")

			switch {
			case is_orderbook: // 初始化，orderbook
				stMsg = append(stMsg, con.FTXQuoteRawToStand(msg, localtime))
			case is_trades: //初始化，trades
				stMsg = con.FTXTradeRawToStand(msg, localtime)
			default:
				logger.Info("FTX msg not found(not trades, orderbook).")
			}
		case 1: // 1代表是okx交易所的数据
			var is_v5_trades bool = strings.Contains(msg, "trades")
			var is_v5_orderbook bool = strings.Contains(msg, "books")
			var is_v5_kline bool = strings.Contains(msg, "candle")

			switch {
			case is_v5_orderbook:
				stMsg = append(stMsg, con.OkexQuoteRawToStand(msg, localtime))
			case is_v5_trades:
				stMsg = append(stMsg, con.OkexTradeRawToStand(msg, localtime))
			case is_v5_kline:
				stMsg = append(stMsg, con.OkexKlineRawToStand(msg, localtime)) // 这里应该是每一条都进行了存储
			default:
				logger.Info("okex msg not found(not init,update,v5_trades,v5_orderbook,kline).")
			}
		case 2: // 2代表是binance交易所的数据
			var is_aggTrades bool = strings.Contains(msg, "aggTrade")
			var is_orderbook bool = strings.Contains(msg, "depthUpdate")
			switch {
			case is_orderbook:
				stMsg = append(stMsg, con.BinanceQuoteRawToStand(msg, localtime))
			case is_aggTrades:
				stMsg = append(stMsg, con.BinanceTradeRawToStand(msg, localtime))
			default:
				logger.Info("binance msg not found(aggtrade,orderbook).")
			}
		default:
			logger.Info("unknown exchange.")
		}

		for i := range stMsg {
			chanStmsg <- stMsg[i]
		}
	}
}

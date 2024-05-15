/*
 * @Author: xwu
 * @Date: 2022-05-23 16:19:04
 * @Last Modified by: xwu
 * @Last Modified time: 2022-07-25 16:53:39
 */
package clientTrader

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
	"winter/Logger"
	"winter/global"
	"winter/messages"
	"winter/utils"

	"github.com/gorilla/websocket"
)

func NewHftClient() HftClient {
	ans := HftClient{okex: NewOkexTraderClient()}
	ans.Count = 0
	ans.PosManager = make([]LogicTrade, global.PMCount)
	ans.Cost = make(map[int64]float64)

	for exid := range global.ExIDList {
		for i, strategyName := range global.AggTradeParams[exid].HftStrategyNames {
			for j := range global.AggTradeParams[exid].HftStrategyParams[i] {
				symbol := global.AggTradeParams[exid].HftStrategyParams[i][j].Symbol
				ii := global.HftUid[exid][strategyName][symbol]
				_p := global.AggTradeParams[exid].HftStrategyParams[i][j]
				ans.PosManager[ii].LDT.Init(_p.ThreEnter, _p.ThreExit)
				fmt.Println(strategyName, " ", symbol, " ", ii, " ", _p.ThreEnter, " ", _p.ThreExit)
			}
		}
	}

	ans.TradeVol = make(map[string]int64)
	ans.Pnl = NewRealPnlManager()
	ans.mutex = new(sync.RWMutex)
	return ans
}
func NewRealPnlManager() RealPnlManager {
	return RealPnlManager{
		Pnl:      make(map[string]float64),
		Fee:      make(map[string]float64),
		Slippage: make(map[string]float64),
		Aggpnl:   make(map[string]float64),
		TrdCount: make(map[string]int64),
	}
}

type RealPnlManager struct {
	Pnl      map[string]float64
	Fee      map[string]float64
	Slippage map[string]float64
	Aggpnl   map[string]float64
	TrdCount map[string]int64
}

func (manager *RealPnlManager) Update(symbol string, side string, fee float64, pnl float64, slippage float64) {
	if _, ok := manager.Pnl[symbol]; ok {
		manager.Pnl[symbol] += pnl
		manager.Fee[symbol] += fee
		manager.Slippage[symbol] += slippage
		manager.Aggpnl[symbol] = manager.Pnl[symbol] + manager.Fee[symbol]
		manager.TrdCount[symbol] += 1
	} else {
		manager.Pnl[symbol] = pnl
		manager.Fee[symbol] = fee
		manager.Slippage[symbol] = slippage
		manager.Aggpnl[symbol] = manager.Pnl[symbol] + manager.Fee[symbol]
		manager.TrdCount[symbol] = 1
	}

	fmt.Println(manager)
}

type HftClient struct {
	PosManager []LogicTrade

	Count int64
	tm    int64
	Cost  map[int64]float64
	okex  OkexTraderClient

	TradeVol map[string]int64
	Pnl      RealPnlManager

	mutex *sync.RWMutex
}

func (hft *HftClient) setCost(Count int64, Mid float64) {
	hft.mutex.Lock()
	hft.Cost[Count] = Mid
	hft.mutex.Unlock()
}

func (hft *HftClient) getCost(Count int64) float64 {
	hft.mutex.RLock()
	temp := hft.Cost[Count]
	hft.mutex.RUnlock()
	return temp
}

func (hft *HftClient) Ping() {
	hft.okex.Ping()
}

func (hft *HftClient) MessageEvent(channel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg_ := <-channel
		var msg_struct messages.OkexChannelResponse
		err := jsonIterator.Unmarshal([]byte(msg_), &msg_struct)
		if err != nil {
			logger.Info("Monitor: event parsed failed : " + msg_)
		}

		switch msg_struct.Event {
		case "login":
			if strings.Contains(msg_, `"code": "0"`) {
				logger.Info("login success")
				go hft.okex.Subcrsibe("orders", "SWAP")
			} else {
				logger.Info("login failed")
				go utils.FeiShuMsg("TradeClient: login failed")
			}

		case "subscribe":
			logger.Info(utils.Strcat(
				"subscribe ",
				msg_struct.Arg.Channel,
				" + ",
				msg_struct.Arg.InstType,
				" success",
			))
		}
	}
}

func (hft *HftClient) MessageOrders(channel chan string, wg *sync.WaitGroup) {
	var exid int = 1 // 交易所的ID，因为这个是okex的监控，所以一定是1
	defer wg.Done()
	for {
		msg_ := <-channel
		var msg_struct messages.OkexTradeChannelMsg
		err := jsonIterator.Unmarshal([]byte(msg_), &msg_struct)
		if err != nil {
			logger.Info("Monitor channel orders: event parsed failed : " + msg_)
		}

		// state:live unfilled filled

		for i := range msg_struct.Data {

			msg_struct_data := msg_struct.Data[i]
			strategy_name := "WINTER"
			symbol := msg_struct_data.InstId
			trader_location := global.HftUid[exid][strategy_name][symbol]

			switch {
			case msg_struct_data.State == "live":
				logger.Info("order in exchage is alive.")
			case msg_struct_data.State == "partially_filled":
				avgprice := utils.String2float(msg_struct_data.AveragePrice)
				side := msg_struct_data.Side
				var slippage float64 = 11111
				countid, err := strconv.ParseInt(msg_struct_data.ClOrdId[3:], 10, 64)
				if err != nil {
					logger.Info("Parse countid failed.")
				} else {
					slippage = (avgprice/hft.getCost(countid) - 1) * 10000
				}
				if side == "sell" {
					slippage *= -1
				}

				fee, err := strconv.ParseFloat(msg_struct_data.FillFee, 64)
				if err != nil {
					logger.Info("Parse fillFee failed.")
				}

				pnl, err := strconv.ParseFloat(msg_struct_data.Pnl, 64)
				if err != nil {
					logger.Info("Parse pnl failed.")
				}

				// 这里只更新pnl，因为订单没有完全更新
				hft.Pnl.Update(msg_struct_data.InstId, side, fee, pnl, slippage)
			case msg_struct_data.State == "filled":
				avgprice := utils.String2float(msg_struct_data.AveragePrice)
				// 这里都是用交易所时间，订单更新时间
				tm := utils.String2int64(msg_struct_data.UTime)
				local_side := hft.PosManager[trader_location].LDT.Update(avgprice, tm) // update LDT StartTm Cost CurrentPosition
				hft.PosManager[trader_location].LDT.TradingStatus(false)               // UPDATE isTrading

				side := msg_struct_data.Side

				var slippage float64 = 11111
				countid, err := strconv.ParseInt(msg_struct_data.ClOrdId[3:], 10, 64)
				if err != nil {
					logger.Info("Parse countid failed.")
				} else {
					slippage = (avgprice/hft.getCost(countid) - 1) * 10000
				}

				if side == "sell" {
					slippage *= -1
				}

				if side != local_side {
					go utils.FeiShuMsg("local side not equal the exchange side.")
					logger.Info("local side not equal the exchange side.")
				}

				fee, err := strconv.ParseFloat(msg_struct_data.FillFee, 64)
				if err != nil {
					logger.Info("Parse fillFee failed.")
				}

				pnl, err := strconv.ParseFloat(msg_struct_data.Pnl, 64)
				if err != nil {
					logger.Info("Parse pnl failed.")
				}

				hft.Pnl.Update(msg_struct_data.InstId, side, fee, pnl, slippage)

				// 输出日志不断进行调整
				logger.Info(utils.Strcat(
					"Okex ",
					symbol,
					" Leverage:",
					msg_struct.Data[i].Lever,
					" Position:",
					strconv.FormatInt(hft.PosManager[trader_location].LDT.CurrentPosition, 10),
					" Create Time:",
					msg_struct.Data[i].CTime,
					" Side:",
					msg_struct.Data[i].Side,
					" Fee:",
					msg_struct.Data[i].FillFee,
					" Pnl:",
					msg_struct.Data[i].Pnl,
					" Slippage:",
					strconv.FormatFloat(slippage, 'f', 2, 64),
				))
			default:
				logger.Info(`msg_struct_data.State: ` + msg_struct_data.State)
			}

			// if msg_struct.Data[i].State == "filled" {
			// 	trader_location := global.HftUid[exid][strategy_name][symbol]
			// 	avgprice := utils.String2float(msg_struct.Data[i].AveragePrice)
			// 	side := hft.PosManager[trader_location].LDT.Update(avgprice, hft.tm) // update LDT StartTm Cost CurrentPosition
			// 	hft.PosManager[trader_location].LDT.TradingStatus(false)             // UPDATE isTrading

			// 	slippage := (utils.String2float(msg_struct.Data[i].AveragePrice)/hft.getCost(symbol) - 1) * 10000 //cost[]
			// 	if side == "SELL" {
			// 		slippage *= -1
			// 	}

			// 	hft.Pnl[0] += utils.String2float(msg_struct.Data[i].FillFee)
			// 	hft.Pnl[1] += utils.String2float(msg_struct.Data[i].Pnl)

			// 	// 输出日志不断进行调整
			// 	logger.Info(utils.Strcat(
			// 		"\nOkex ",
			// 		symbol,
			// 		" Leverage:",
			// 		msg_struct.Data[i].Lever,
			// 		" Position:",
			// 		strconv.FormatInt(hft.PosManager[trader_location].LDT.CurrentPosition, 10),
			// 		" Create Time:",
			// 		msg_struct.Data[i].CTime,
			// 		" Side:",
			// 		msg_struct.Data[i].Side,
			// 		" Fee:",
			// 		msg_struct.Data[i].FillFee,
			// 		" Pnl:",
			// 		msg_struct.Data[i].Pnl,
			// 		" Slippage:",
			// 		strconv.FormatFloat(slippage, 'f', 2, 64),
			// 	))
			// 	fmt.Println(hft.Pnl)
			// }
		}
	}
}

func (hft *HftClient) Monitor(channelEvent chan string, chanMessageOrders chan string, wg *sync.WaitGroup) {
	// 一共有几种消息类型，然后需要确保重连和 重新登录只在这里，这样我觉得逻辑会清晰一点TODO:
	defer wg.Done()
	for {
		hft.mutex.Lock()
		if !hft.okex.getConnFlag() { // 如果没有的话那么进行登录
			hft.okex.Create_Private_Dialogue() // 这里会修改hft.okex.getConnFlag()的状态
		}
		hft.mutex.Unlock()

		for {
			if hft.okex.getConnFlag() {
				_, message_okex, err := hft.okex.Conn.ReadMessage()
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) || err != nil {
					hft.okex.Conn.Close()
					hft.okex.Conn = nil
					hft.okex.setConnFlag(false) // 修改状态
					break
				}

				// 这里解析成为一个空接口?暂时只能if else
				Logger.TerminalSlientLogger.Info(string(message_okex))
				msg_ := string(message_okex)
				switch {
				case strings.Contains(string(message_okex), `pong`):
					logger.Info("pont msg")
				case strings.Contains(msg_, "event"): // 登录消息
					channelEvent <- msg_ // 这里用函数来传递或者类结构体来传递？
				case strings.Contains(msg_, `"channel":"orders"`):
					chanMessageOrders <- msg_
				case strings.Contains(msg_, `"op":"order"`):
					// TODO: 这里应该写一个交易所收到了下单得操作
					logger.Info("exchange receive the order")
				default:
					logger.Info("nothing happened here : " + msg_)
				}
				// if !strings.Contains(string(message_okex), `pong`) {

				// 	// 这里开始对消息进行分类.
				// 	if strings.Contains(msg_, "event") { // 登录消息
				// 		logger.Info("kkkkkkkkkkkk")
				// 	} else if strings.Contains(msg_, `"channel":"orders"`) {
				// 		/*var msg_struct messages.OkexTradeChannelMsg
				// 		err := jsonIterator.Unmarshal([]byte(msg_), &msg_struct)
				// 		if err != nil {
				// 			logger.Info("Monitor channel orders: event parsed failed")
				// 		}

				// 		for i := range msg_struct.Data {
				// 			symbol := msg_struct.Data[i].InstId
				// 			strategy_name := "WINTER"
				// 			if msg_struct.Data[i].State == "filled" {
				// 				trader_location := global.HftUid[1][strategy_name][symbol]
				// 				avgprice := utils.String2float(msg_struct.Data[i].AveragePrice)
				// 				side := hft.PosManager[trader_location].LDT.Update(avgprice, hft.tm) // update LDT StartTm Cost CurrentPosition
				// 				go hft.PosManager[trader_location].LDT.TradingStatus(false) // UPDATE isTrading

				// 				slippage := (utils.String2float(msg_struct.Data[i].AveragePrice)/hft.getCost(symbol) - 1) * 10000 //cost[]
				// 				if side == "SELL" {
				// 					slippage *= -1
				// 				}

				// 				hft.Pnl[0] += utils.String2float(msg_struct.Data[i].FillFee)
				// 				hft.Pnl[1] += utils.String2float(msg_struct.Data[i].Pnl)

				// 				// 输出日志不断进行调整
				// 				logger.Info(utils.Strcat(
				// 					"\nOkex ",
				// 					symbol,
				// 					" Leverage:",
				// 					msg_struct.Data[i].Lever,
				// 					" Position:",
				// 					strconv.FormatInt(hft.PosManager[trader_location].LDT.CurrentPosition, 10),
				// 					" Create Time:",
				// 					msg_struct.Data[i].CTime,
				// 					" Side:",
				// 					msg_struct.Data[i].Side,
				// 					" Fee:",
				// 					msg_struct.Data[i].FillFee,
				// 					" Pnl:",
				// 					msg_struct.Data[i].Pnl,
				// 					" Slippage:",
				// 					strconv.FormatFloat(slippage, 'f', 2, 64),
				// 				))
				// 				fmt.Println(hft.Pnl)
				// 			}
				// 		}*/

				// 	} else {
				// 		logger.Info("nothing happened here : " + msg_)
				// 	}
				// }
			} else {
				hft.okex.Conn.Close()
				hft.okex.Conn = nil
				hft.okex.setConnFlag(false) // 改变状态
				break
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// 需要保证时刻是连接的
func (hft *HftClient) Run(channel chan messages.AggSignal, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		hft.mutex.Lock()
		if !hft.okex.getConnFlag() { // 如果没有的话那么进行登录
			hft.okex.Create_Private_Dialogue() // 这里会修改hft.okex.getConnFlag()的状态
		}
		hft.mutex.Unlock()

		for {
			signal := <-channel

			// 这里是根据Amount自动平衡需要下多少张，默认是5倍杠杆，最少是1张
			// 如果是第一次下单那么mao中找不到，就去赋值
			if _, ok := hft.TradeVol[signal.Symbol]; (!ok) && (signal.Kind == 0) {
				// 向下取整
				var leverage float64 = 5
				var multiplier int64 = int64(math.Floor(global.Amount * leverage / (global.Instruments[signal.Symbol].Size * signal.Mid)))
				if multiplier == 0 {
					multiplier = 1 // 确保最少是一张
				}
				hft.TradeVol[signal.Symbol] = multiplier
				fmt.Println(global.Amount, " ", leverage, " ", global.Instruments[signal.Symbol].Size, " ", signal.Mid)
				fmt.Println(hft.TradeVol)
			}

			// prewarm,一分钟内不交易
			if signal.Localtime-global.StartTime < 60000 {
				continue
			}

			// fmt.Println(signal.FinalSignal, " ", signal.Ask1Price, " ", signal.Mid, " ", signal.Bid1Price)

			hft.tm = signal.Localtime

			// 这里检验是否是连接的状态，这里考虑改成原子操作，会不会比读写锁合适一些？
			// 这里因为目前的交易逻辑是当前币种的message只会触发当前币种的交易，不会触发其他的，所以不需要对币种进行循环
			// 但是不同的策略可能会有叠加的仓位，所以理论上需要对策略名字进行循环，但是暂时限定只有一个策略，所以也不需要进行循环
			if hft.okex.getConnFlag() {
				// 遍历每一个策略
				var PositionChg int64 = 0
				params := global.AggTradeParams[signal.ExID]

				// 这里的0是策略的名字，暂时只有一个策略，明确是什么交易所，什么策略，然后什么品种
				// for j := range params.HftStrategyParams[0] {} // 第一个策略名字下，遍历币种
				ii := global.HftUid[signal.ExID][params.HftStrategyNames[0]][signal.Symbol]
				param := params.HftStrategyParams[0][ii]

				// 检查是否正在交易
				if hft.PosManager[ii].LDT.GetTradingStatus() {
					continue
				}

				// sanity check
				if signal.Symbol != param.Symbol {
					logger.Fatal("error!!!!!")
				}

				currentPosition := hft.PosManager[ii].LDT.GetCurrentPosition()
				targetPosition := hft.PosManager[ii].LDT.GetTargetPosition(
					hft.tm,
					signal.FinalSignal,
					signal.Ask1Price,
					signal.Bid1Price,
				)

				if targetPosition != currentPosition {
					PositionChg = targetPosition - currentPosition
					hft.PosManager[ii].LDT.TradingStatus(true)
				}

				// 如果有仓位变化，那么进行下单的逻辑
				if PositionChg != 0 {
					var side, size string
					if PositionChg > 0 {
						side = "buy"
						size = strconv.FormatInt(PositionChg*hft.TradeVol[signal.Symbol], 10)
					} else if PositionChg < 0 {
						side = "sell"
						size = strconv.FormatInt(-1*PositionChg*hft.TradeVol[signal.Symbol], 10)
					}

					// 这里是为了计算滑点
					hft.setCost(hft.Count, signal.Mid)
					//hft.Cost[signal.Symbol] = signal.Mid

					// 这里一定需要是发起协程，否则会堵塞后面
					go hft.okex.PostMarketOrder(
						signal.ExID,
						signal.Symbol,
						side,
						size,
						int(hft.Count),
						global.HftName,
						signal,
						hft,
					)

					// 这里是计算一共有几个订单，因为这个count会在monitor中返回
					// 所以可以用来校对下单的时候的mid价格和实际成交的avgprice之间的差距
					hft.Count += 1

					logger.Info(
						utils.Strcat(
							"post order starts: ",
							strconv.FormatFloat(signal.FinalSignal, 'f', 4, 64),
							", ",
							side,
							", ",
							size,
						),
					)
				}
			} else {
				break
			}
		}
	}
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:45
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-05 15:33:30
 */
package clientAlpha

import (
	"sync"

	"winter/features"
	"winter/global"
	"winter/lightgbm"
	"winter/mathpro"
	"winter/messages"
)

// ants需要闭包，为了多线程准备的
type mytask func()

func OrderWrapper(alpha features.Alphas, evt *messages.OrderbookBase, wg *sync.WaitGroup) mytask {
	return func() {
		alpha.Quote(evt)
		wg.Done()
	}
}

func TradeWrapper(alpha features.Alphas, evt *messages.TradeBase, wg *sync.WaitGroup) mytask {
	return func() {
		alpha.Trade(evt)
		wg.Done()
	}
}

func KlineWrapper(alpha features.Alphas, evt *messages.KlineBase, wg *sync.WaitGroup) mytask {
	return func() {
		alpha.Kline(evt)
		wg.Done()
	}
}

// 感觉这个框架需要重新改一改
type AlphaSymbol struct {
	Factors []features.Alphas

	FactorCols []string

	Mid          float64
	Ask1Price    float64
	MinAsk1Price float64
	Bid1Price    float64
	MaxBid1Price float64
	TradeVol     float64
}

type AlphaFactory struct {
	AlphasExIDSymbols []AlphaSymbol // 先是哪个交易所,然后是哪个币种
	TickTime          []int64

	Model *lightgbm.Ensemble

	II []map[string]int // 哪个交易所的哪个币种的对应的vector的位置
}

func (client *AlphaFactory) Run(
	chanStmsg chan messages.AggStMsg,
	chanSignal chan messages.AggSignal,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	// // 这里定义ants多线程相关的代码
	// defer ants.Release()
	// var p *ants.Pool
	// if global.Num_threads > 1 {
	// 	p, _ = ants.NewPool(global.Num_threads)
	// 	defer p.Release()
	// }

	for stmsg := range chanStmsg {
		// ExID Status Kinds Symbols OrdEvt TrdEvt
		// var stmsg messages.AggStMsg = <-chanStmsg
		if stmsg.Status == 0 {
			chanSignal <- messages.AggSignal{Status: 0}
			break
		}

		// if stmsg.Kinds == 0 {
		// 	if mathpro.AlmostEqual(stmsg.OrdEvt.Asks[0][0], 0) || mathpro.AlmostEqual(stmsg.OrdEvt.Bids[0][0], 0) {
		// 		continue // 这里加一个过滤，如果是零的暂时过滤掉
		// 	}
		// }

		var signals messages.AggSignal = messages.AggSignal{}

		signals.ExID = stmsg.ExID
		signals.Symbol = stmsg.Symbols
		signals.Kind = stmsg.Kinds
		signals.Status = stmsg.Status

		switch stmsg.ExID {
		case 0: // FTX交易所，这里应该改成单线程，但是FTX因为弃置，所以暂时没改
			localtime := stmsg.OrdEvt.Localtime
			var ii int = client.II[stmsg.ExID][stmsg.Symbols]

			// 这里应该明确是哪个交易所的，方便进行更新
			switch stmsg.Kinds {
			case 0: // Orderbook
				client.AlphasExIDSymbols[ii].Mid = (stmsg.OrdEvt.Bids[0][0] + stmsg.OrdEvt.Asks[0][0]) / 2
				client.AlphasExIDSymbols[ii].Ask1Price = stmsg.OrdEvt.Asks[0][0]
				client.AlphasExIDSymbols[ii].Bid1Price = stmsg.OrdEvt.Bids[0][0]
				localtime = stmsg.OrdEvt.Localtime

				for f_i := 0; f_i < len(global.AggParameters.Alpha.FTX.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Quote(&stmsg.OrdEvt)
				}
			case 1:
				localtime = stmsg.TrdEvt.Localtime

				for f_i := 0; f_i < len(global.AggParameters.Alpha.FTX.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Trade(&stmsg.TrdEvt)
				}
			}

			// 更新signal
			signals.Mid = client.AlphasExIDSymbols[ii].Mid
			signals.Ask1Price = client.AlphasExIDSymbols[ii].Ask1Price
			signals.Bid1Price = client.AlphasExIDSymbols[ii].Bid1Price
			signals.Localtime = localtime
			signals.Signals = make([]float64, len(global.AggParameters.Alpha.FTX.P))
			for i := 0; i < len(global.AggParameters.Alpha.FTX.P); i++ {
				signals.Signals[i] = client.AlphasExIDSymbols[ii].Factors[i].Get()
			}
			signals.Columns = client.AlphasExIDSymbols[ii].FactorCols
		case 1: // Okex交易所
			var localtime, starttime int64
			var ii int = client.II[stmsg.ExID][stmsg.Symbols]

			// 这里应该明确是哪个交易所的，方便进行更新,不开多线程的效果更好，暂时没有对比非多线程和多线程池子
			switch stmsg.Kinds {
			case 0: // OrderBook
				client.AlphasExIDSymbols[ii].Mid = (stmsg.OrdEvt.Asks[0][0] + stmsg.OrdEvt.Bids[0][0]) / 2
				client.AlphasExIDSymbols[ii].Ask1Price = stmsg.OrdEvt.Asks[0][0]
				client.AlphasExIDSymbols[ii].Bid1Price = stmsg.OrdEvt.Bids[0][0]
				localtime = stmsg.OrdEvt.Localtime

				for f_i := 0; f_i < len(global.AggParameters.Alpha.Okex.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Quote(&stmsg.OrdEvt)
				}
			case 1: // Trade
				localtime = stmsg.TrdEvt.Localtime

				// if global.Num_threads == 1 {
				for f_i := 0; f_i < len(global.AggParameters.Alpha.Okex.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Trade(&stmsg.TrdEvt)
				}
				// } else {
				// 	var wg_task sync.WaitGroup
				// 	for f_i := 0; f_i < len(global.AggParameters.Alpha.Okex.P); f_i++ {
				// 		wg_task.Add(1)
				// 		p.Submit(TradeWrapper(
				// 			client.AlphasExIDSymbols[ii].Factors[f_i],
				// 			&stmsg.TrdEvt,
				// 			&wg_task,
				// 		))
				// 	}
				// 	wg_task.Wait()
				// }
			case 2: // Kline
				localtime = stmsg.KlineEvt.Localtime
				starttime = stmsg.KlineEvt.StartTm

				for f_i := 0; f_i < len(global.AggParameters.Alpha.Okex.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Kline(&stmsg.KlineEvt)
				}
			default:
				logger.Fatal("error msg kind")
			}

			// 更新signal
			signals.Mid = client.AlphasExIDSymbols[ii].Mid
			signals.Ask1Price = client.AlphasExIDSymbols[ii].Ask1Price
			signals.Bid1Price = client.AlphasExIDSymbols[ii].Bid1Price
			signals.Localtime = localtime
			signals.StartTm = starttime
			signals.Signals = make([]float64, len(global.AggParameters.Alpha.Okex.P))
			for i := 0; i < len(global.AggParameters.Alpha.Okex.P); i++ {
				signals.Signals[i] = client.AlphasExIDSymbols[ii].Factors[i].Get() // 这个Get非常重要
			}
			signals.Columns = client.AlphasExIDSymbols[ii].FactorCols
		case 2:
			var localtime, starttime int64
			var ii int = client.II[stmsg.ExID][stmsg.Symbols]

			// 这里应该明确是哪个交易所的，方便进行更新,不开多线程的效果更好，暂时没有对比非多线程和多线程池子
			switch stmsg.Kinds {
			case 0: // OrderBook
				client.AlphasExIDSymbols[ii].Mid = (stmsg.OrdEvt.Asks[0][0] + stmsg.OrdEvt.Bids[0][0]) / 2
				client.AlphasExIDSymbols[ii].Ask1Price = stmsg.OrdEvt.Asks[0][0]
				client.AlphasExIDSymbols[ii].Bid1Price = stmsg.OrdEvt.Bids[0][0]
				localtime = stmsg.OrdEvt.Localtime

				for f_i := 0; f_i < len(global.AggParameters.Alpha.Binance.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Quote(&stmsg.OrdEvt)
				}
			case 1: // Trade
				localtime = stmsg.TrdEvt.Localtime
				for f_i := 0; f_i < len(global.AggParameters.Alpha.Binance.P); f_i++ {
					client.AlphasExIDSymbols[ii].Factors[f_i].Trade(&stmsg.TrdEvt)
				}
			default:
				logger.Fatal("error msg kind")
			}

			// 更新signal
			signals.Mid = client.AlphasExIDSymbols[ii].Mid
			signals.Ask1Price = client.AlphasExIDSymbols[ii].Ask1Price
			signals.Bid1Price = client.AlphasExIDSymbols[ii].Bid1Price
			signals.Localtime = localtime
			signals.StartTm = starttime
			signals.Signals = make([]float64, len(global.AggParameters.Alpha.Binance.P))
			for i := 0; i < len(global.AggParameters.Alpha.Binance.P); i++ {
				signals.Signals[i] = client.AlphasExIDSymbols[ii].Factors[i].Get() // 这个Get非常重要
			}
			signals.Columns = client.AlphasExIDSymbols[ii].FactorCols

			// fmt.Println(signals.Signals)
		}

		var coef_ = []float64{5.69786531e-02, 5.61567071e+00, -4.95655158e-01, 6.84224778e-01,
			-2.37378427e-01, 2.68076449e-02, 2.89637700e-03, -6.33137356e-02,
			8.55781988e-02, -4.93063352e-03, 1.02846877e-02, 2.98216444e-01,
			-3.55340851e-01, -3.54826702e-01, 1.27252731e-06, -6.93417681e-02,
			-2.87325909e-02, -1.01763198e-01, -4.23544173e-02, -1.85089466e-02,
			1.98044431e-01, -1.15998522e-01, -3.47876453e-02, 3.43642675e-02,
			-9.33416333e-06, 3.98810310e-06, -1.15186218e-04, 2.22064949e-02,
			1.48873399e-02, 4.59989441e-02, 5.54420364e-02, 7.12449214e-05,
			-1.08104071e-01, 7.08911745e-06, -3.01409386e-06, 6.18531973e-06,
			-2.52683536e-01, -1.82966905e-02, 8.16360304e-02, 2.70504162e-01,
			-1.31321608e-01, -6.16425182e-02, 4.91872267e-02, 5.22253036e-01,
			-4.60746502e-02, -9.55613678e-02, 7.60920899e+00, -6.39374204e+00,
			4.69370371e+00, -7.12102547e-03, -9.38922556e-08, 8.01063509e-08}
		// TODO: 这里使用lightgbm模型或者线性模型,然后假定就是okex交易所
		if global.LGBPredict || global.Mode == "hft" {
			// signals.FinalSignal = client.Model.PredictSingle(signals.Signals, 0)
			signals.FinalSignal = 0
			for i := range coef_ {
				signals.FinalSignal += coef_[i] * signals.Signals[i]
			}
		}

		var ii int = client.II[stmsg.ExID][stmsg.Symbols]
		CurrentTickTime := (signals.Localtime + 10*int64(ii)) / 1000
		if CurrentTickTime > client.TickTime[ii] {
			signals.IsOpen = true
			client.TickTime[ii] = CurrentTickTime
		} else {
			signals.IsOpen = false
		}

		if !mathpro.Isfinite(signals.Ask1Price) {
			continue
		}

		// 这里是所有消息都会触发或者是只有trade消息会触发（跟单）
		if global.Trigger == "all" {
			chanSignal <- signals
		} else if global.Trigger == "follow" {
			if stmsg.Kinds == 1 || stmsg.Kinds == 2 {
				chanSignal <- signals
			}
		} else if global.Trigger == "time" {
			var ii int = client.II[stmsg.ExID][stmsg.Symbols]

			if stmsg.Kinds == 0 {
				CurrentTickTime := (signals.Localtime + 10*int64(ii)) / 1000

				if CurrentTickTime > client.TickTime[ii] {
					chanSignal <- signals
					client.TickTime[ii] = CurrentTickTime
					// fmt.Println(signals.FinalSignal)
				}
			}
		}
	}
}

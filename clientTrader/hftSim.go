/*
 * @Author: xwu
 * @Date: 2022-05-23 16:19:04
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-11 10:31:02
 */
package clientTrader

import (
	"math"
	"os"
	"strconv"
	"sync"

	// "winter/container/orders"
	"winter/global"
	"winter/messages"

	"github.com/olekukonko/tablewriter"
)

func NewHftSimClient() HftSimClient {
	ans := HftSimClient{}
	ans.PosManager = make([]LogicTrade, global.PMCount)
	ans.Cost = make(map[string]float64)
	ans.TradeVol = make(map[string]int64)
	ans.PnlManager = NewPnlCalc()

	for exid := range global.ExIDList {
		for i, strategyName := range global.AggTradeParams[exid].HftStrategyNames {
			for j := range global.AggTradeParams[exid].HftStrategyParams[i] {
				symbol := global.AggTradeParams[exid].HftStrategyParams[i][j].Symbol
				ii := global.HftUid[exid][strategyName][symbol]
				_p := global.AggTradeParams[exid].HftStrategyParams[i][j]
				ans.PosManager[ii].LDT.Init(_p.ThreEnter, _p.ThreExit)
			}
		}
	}

	// ans.FactorsManager = make([]FactorsCache, global.FMCount)
	// for k1 := range global.AlphaUid {
	// 	for k2 := range global.AlphaUid[k1] {
	// 		// fmt.Println(k1, k2)
	// 		ans.FactorsManager[global.AlphaUid[k1][k2]] = NewFactorsCache(global.ExIDList[k1], k2)
	// 	}
	// }

	num_signals := len(global.AggParameters.Data.Okex.Subscribe_symbols)/2 + len(global.AggParameters.Data.Binance.Subscribe_symbols)/2 // 因为每次都订阅了order和trade
	ans.ProcessSignal = NewHandleToSignal(
		global.ReturnInterval,
		num_signals,
		len(global.AggParameters.Alpha.Okex.P)+len(global.AggParameters.Alpha.Binance.P),
		global.Dumper,
	)

	ans.StartTm = make(map[string]int64)
	ans.Pos = make(map[string]int)

	return ans
}

/*
这里主要有三个结构体的同步（非线程同步那种，状态要统一），第一个是PosManager，第二个是PnlManager，第三个是Orders
PosManager计算出仓位变化，然后下单，Orders计算下单的状态，当成交后需要有PnlManager来计算pnl
然后希望这三个结构体独立，只通过数据来实现交换
第一个PosManager计算仓位的变动，这里计算预计的仓位，然后绑定一个oid或者几个oid，当oid都成交，证明仓位应该同步了如果订单cancel了，就取消oid
第二个是Orders，需要给出成交订单成交量和成交价格，
第三个是计算pnl，根据Orders返回的信息来计算pnl
*/
type HftSimClient struct {
	PosManager []LogicTrade

	StartTm map[string]int64
	Pos     map[string]int

	// FactorsManager []FactorsCache
	PnlManager PnlCalc

	FactorNames []string

	Count         int64
	Cost          map[string]float64
	TradeVol      map[string]int64
	ProcessSignal handleToSignal

	global_oid int // 订单所需要的oid
}

func (hft *HftSimClient) Run(channel chan messages.AggSignal, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		signal := <-channel

		if signal.Status == 0 {
			break
		}

		// ii := global.AlphaUid[signal.ExID][signal.Symbol]

		hft.FactorNames = signal.Columns
		// hft.FactorsManager[ii].Update(&signal) // 这里负责计算IC以及计算dump因子

		hft.ProcessSignal.Update(&signal) // 这里负责计算IC以及计算dump因子

		// fmt.Println(signal.FinalSignal)

		// 这里计算的是每次应该交易多少，计算的是每次交易多少张
		if _, ok := hft.TradeVol[signal.Symbol]; !ok {
			var multiplier int64 = int64(math.Floor(global.Amount / (global.Instruments[signal.Symbol].Size * signal.Mid)))
			if multiplier == 0 {
				multiplier = 1
			}
			// hft.TradeVol[signal.Symbol] = multiplier
			hft.TradeVol[signal.Symbol] = 1
		}

	}

	performance := hft.ProcessSignal.End()

	// output performance
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ic_td", "ir_td", "ic_total", "ir_total", "min", "max", "mean", "std", "nan(bp)"})
	for ii, name := range hft.FactorNames {
		table.Append([]string{
			name,
			strconv.FormatFloat(performance[0][ii], 'f', 4, 64),
			strconv.FormatFloat(performance[1][ii], 'f', 2, 64),
			strconv.FormatFloat(performance[2][ii], 'f', 4, 64),
			strconv.FormatFloat(performance[3][ii], 'f', 2, 64),
			strconv.FormatFloat(performance[4][ii], 'f', 4, 64),
			strconv.FormatFloat(performance[5][ii], 'f', 4, 64),
			strconv.FormatFloat(performance[6][ii], 'f', 4, 64),
			strconv.FormatFloat(performance[7][ii], 'f', 4, 64),
			strconv.FormatFloat(performance[8][ii], 'f', 2, 64),
		})
	}
	table.Render() // Send output
}

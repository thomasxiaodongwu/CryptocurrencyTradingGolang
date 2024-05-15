/*
 * @Author: xwu
 * @Date: 2022-01-22 14:45:16
 * @Last Modified by: xwu
 * @Last Modified time: 2022-01-26 20:43:54
 */

package params

type ClientTradeParams struct {
	FTX     TradeParams
	Okex    TradeParams
	Binance TradeParams
	Dydx    TradeParams
}

type TradeParams struct {
	// 策略名字，这里算是一个历史遗留问题，这里策略名字只能写一个，这个是为了monitor方便写
	HftStrategyNames  []string             // 这里只有一个名字
	HftStrategyParams [][]HftStrategyParam // 这里也只有一个相当于是[][1]HftStrategyParam

	ArbitrageStrategyParams []ArbitrageStrategyParam
}

type ArbitrageStrategyParam struct {
	SymbolArb1 string
	SymbolArb2 string

	Interval int64
	Range    int64
}

type HftStrategyParam struct {
	Model_name string // Linear,LGBM,XGB
	Symbol     string
	SubSymbol  string
	Weights    []float64
	SubWeights []float64
	LGBMPath   string
	ThreEnter  float64
	ThreExit   float64
	Size       int64
}

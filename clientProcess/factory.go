/*
 * @Author: xwu
 * @Date: 2022-01-22 17:32:51
 * @Last Modified by: xwu
 * @Last Modified time: 2022-09-26 15:21:18
 */
package clientProcess

import (
	"winter/global"
)

func NewConvert() Convert {
	ans := Convert{}
	p := global.AggParameters.Data

	// initial FTX
	ans.FTXOrderbook = make([]FTXQteEvt, len(p.FTX.Subscribe_symbols))
	ans.FTX_ii = global.AlphaUid[global.ExIDMap["FTX"]]

	// initial okex
	ans.OkexOrderbook = make([]okexQteEvt, len(p.Okex.Subscribe_symbols))
	for i := range ans.OkexOrderbook {
		ans.OkexOrderbook[i].Init()
	}
	ans.Okex_ii = global.AlphaUid[global.ExIDMap["Okex"]]

	ans.BinanceOrderbook = make([]BinanceQteEvt, len(p.Binance.Subscribe_symbols)) // binance暂时不用增量更新，所以不用Init函数
	ans.Binance_ii = global.AlphaUid[global.ExIDMap["Binance"]]

	return ans
}

/*
 * @Author: xwu
 * @Date: 2022-02-24 17:42:06
 * @Last Modified by: xwu
 * @Last Modified time: 2022-03-10 12:43:44
 */
package messages

// import (
// 	"winter/operator"
// )

// type Index struct {
// 	bsize3s float64
// 	asize3s float64
// }
// type HistGlobalData struct {
// 	CurQuote  []float64 // 所有的数据的最新的quote
// 	HistTrd   []float64 // 所有品种一段时间的逐笔成交数据
// 	HistBar   []float64 // 所有的数据历史Bar
// 	HistIndex []float64 // 所有品种的一些指标
// }

type HistData struct {
	Mid       [][]float64
	AdjustMid [][]float64

	Ask1Price [][]float64
	Bid1Price [][]float64

	// AdjustAsk1Price [][]operator.OpAdjustMid
	// AdjustBid1Price [][]operator.OpAdjustMid

	TradeVol []float64
}

func (hist *HistData) Init(AlphaUid map[int](map[string]int), LOG_SIZE int64) {
	hist.Mid = make([][]float64, len(AlphaUid))
	hist.AdjustMid = make([][]float64, len(AlphaUid))
	hist.Ask1Price = make([][]float64, len(AlphaUid))
	hist.Bid1Price = make([][]float64, len(AlphaUid))
	// hist.AdjustAsk1Price = make([][]operator.OpAdjustMid, len(AlphaUid))
	// hist.AdjustBid1Price = make([][]operator.OpAdjustMid, len(AlphaUid))
	for i := range AlphaUid {
		hist.Mid[i] = make([]float64, len(AlphaUid[i]))
		hist.AdjustMid[i] = make([]float64, len(AlphaUid[i]))
		hist.Ask1Price[i] = make([]float64, len(AlphaUid[i]))
		hist.Bid1Price[i] = make([]float64, len(AlphaUid[i]))
		// hist.AdjustAsk1Price[i] = make([]operator.OpAdjustMid, len(AlphaUid[i]))
		// hist.AdjustBid1Price[i] = make([]operator.OpAdjustMid, len(AlphaUid[i]))

		// for _, j := range AlphaUid[i] {
		// 	hist.AdjustAsk1Price[i][j].Init(LOG_SIZE)
		// 	hist.AdjustBid1Price[i][j].Init(LOG_SIZE)
		// }
	}
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:47
 * @Last Modified by:   xwu
 * @Last Modified time: 2021-12-26 18:45:47
 */
package features

import (
	"container/list"
	"sync"
	"winter/messages"
	"winter/params"
)

type Alphas interface {
	Init(p *params.AlphaParamsBase)
	Quote(evt *messages.OrderbookBase) float64
	QuoteSync(evt *messages.OrderbookBase, wg *sync.WaitGroup) float64
	Trade(evt *messages.TradeBase) float64
	TradeSync(evt *messages.TradeBase, wg *sync.WaitGroup) float64
	Kline(evt *messages.KlineBase) float64
	KlineSync(evt *messages.KlineBase, wg *sync.WaitGroup) float64
	Get() float64
}

type HistData struct {
	Ob []*messages.OrderbookBase
	Td []*messages.TradeBase

	Obnew *list.List
	Tdnew *list.List
}

// type aggAlphaParams struct {
// 	Interval int64
// }
// func (client *aggAlphaParams) Get() int64 {
// 	return client.Interval
// }

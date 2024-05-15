/*
 * @Author: xwu
 * @Date: 2022-05-29 10:32:39
 * @Last Modified by: xwu
 * @Last Modified time: 2022-09-29 21:38:16
 */

package features

import (
	"sync"
	"winter/mathpro"
	"winter/messages"
	"winter/operator"
	"winter/params"
)

// ////////////////////////////////////////////////////////////////////
type STA000 struct {
	interval int64
	alpha    float64

	bVol, sVol operator.TimeSum
}

func (cliff *STA000) Init(p *params.AlphaParamsBase) {
	cliff.interval = p.Interval
	cliff.alpha = 0

	cliff.bVol = operator.NewTimeSum(p.Interval)
	cliff.sVol = operator.NewTimeSum(p.Interval)
}

func (cliff *STA000) Quote(evt *messages.OrderbookBase) float64 {
	return cliff.alpha
}

func (cliff *STA000) QuoteSync(evt *messages.OrderbookBase, wg *sync.WaitGroup) float64 {
	defer wg.Done()
	cliff.Quote(evt)
	return cliff.alpha
}

func (cliff *STA000) Trade(evt *messages.TradeBase) float64 {
	if evt.Side > 0 {
		cliff.bVol.Update(evt.Localtime, evt.Size)
	} else if evt.Side < 0 {
		cliff.sVol.Update(evt.Localtime, evt.Size)
	}

	b := cliff.bVol.Value()
	s := cliff.sVol.Value()

	if mathpro.AlmostEqual(b+s, 0) {
		cliff.alpha = 0
	} else {
		cliff.alpha = (b - s) / (b + s) * 10000
	}
	return cliff.alpha
}

func (cliff *STA000) TradeSync(evt *messages.TradeBase, wg *sync.WaitGroup) float64 {
	defer wg.Done()
	cliff.Trade(evt)
	return cliff.alpha
}

func (cliff *STA000) Kline(evt *messages.KlineBase) float64 {
	return cliff.alpha
}

func (cliff *STA000) KlineSync(evt *messages.KlineBase, wg *sync.WaitGroup) float64 {
	defer wg.Done()
	cliff.Kline(evt)
	return cliff.alpha
}

func (cliff *STA000) Get() float64 {
	return cliff.alpha
}

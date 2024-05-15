/*
 * @Author: xwu
 * @Date: 2022-02-23 17:09:17
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-24 11:23:36
 */
package operator

import (
	"fmt"
)

func NewTimeBar(interval int64, window int64) TimeBar {
	op := TimeBar{}

	op.Tail_ = -1
	op.HistData = make([]Bar, window)

	op.Interval = interval
	op.Window = window
	op.CurrentBarNum = -1
	return op
}

type Bar struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64

	Tr  float64
	Atr float64
}

type TimeBar struct {
	Tail_    int64
	HistData []Bar // 这里的最后一个Bar是没有生成完成的Bar

	Interval int64 // 一个Bar的时间长度
	Window   int64 // 一共有多少个Bar，这个就是size

	CurrentBarNum int64 // 用来标记是否是下一个Bar
}

func (op *TimeBar) Update(time int64, midprice float64, vol float64) float64 {
	var num int64 = time / op.Interval

	// 判断是否是当前bar
	if num > op.CurrentBarNum {
		if op.Tail_ > 0 {
			op.HistData[op.Tail_%op.Window].Atr = 0.8*op.HistData[op.Tail_%op.Window].Atr + 0.2*op.HistData[op.Tail_%op.Window].Tr
		} else if op.Tail_ == 0 {
			op.HistData[op.Tail_%op.Window].Atr = op.HistData[op.Tail_%op.Window].Tr
		}

		op.Tail_ += 1

		op.HistData[op.Tail_%op.Window].Open = midprice
		op.HistData[op.Tail_%op.Window].High = midprice
		op.HistData[op.Tail_%op.Window].Low = midprice
		op.HistData[op.Tail_%op.Window].Close = midprice
		op.HistData[op.Tail_%op.Window].Volume = vol

		op.HistData[op.Tail_%op.Window].Tr = 0
		op.CurrentBarNum = num
	} else if num == op.CurrentBarNum {
		if midprice > op.HistData[op.Tail_%op.Window].High {
			op.HistData[op.Tail_%op.Window].High = midprice
		}

		if midprice < op.HistData[op.Tail_%op.Window].Low {
			op.HistData[op.Tail_%op.Window].Low = midprice
		}

		op.HistData[op.Tail_%op.Window].Close = midprice

		op.HistData[op.Tail_%op.Window].Volume += vol

		op.HistData[op.Tail_%op.Window].Tr = op.HistData[op.Tail_%op.Window].High - op.HistData[op.Tail_%op.Window].Low
	} else {
		fmt.Println("time error.")
	}
	return op.HistData[op.Tail_%op.Window].Tr
}

func (op *TimeBar) Get() Bar {
	if op.Tail_ > 0 {
		return op.HistData[(op.Tail_-1)%op.Window]
	} else {
		return op.HistData[op.Tail_%op.Window]
	}
}

func (op *TimeBar) GetHist() []Bar {
	return op.HistData
}

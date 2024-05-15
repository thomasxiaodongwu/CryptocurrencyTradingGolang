/*
 * @Author: xwu 
 * @Date: 2021-12-26 18:43:41 
 * @Last Modified by:   xwu 
 * @Last Modified time: 2021-12-26 18:43:41 
 */
package operator

import "math"

type OpTimer struct {
	interval int64
	data_    []valueTm
	sum_     float64
	count    float64
	last_    float64

	ema_sum_ float64
}

func (op *OpTimer) Init(interval int64) {
	op.interval = interval
	op.count = 0
	op.last_ = 0
	op.sum_ = 0
	op.ema_sum_ = 0
}

func (op *OpTimer) Update(time int64, x float64) {
	var d valueTm = valueTm{time, x}
	op.data_ = append(op.data_, d)
	op.sum_ += x
	op.count += 1
	op.last_ = x

	if len(op.data_) == 1 {
		op.ema_sum_ = x
	} else {
		var tmdiff float64 = float64(time - op.data_[len(op.data_)-2].Tm)
		var time_decay float64 = math.Exp(-0.07 * tmdiff / 1000)
		op.ema_sum_ = op.ema_sum_*time_decay + x
	}

	var i int = 0
	for len(op.data_) > 0 {
		if time-op.data_[i].Tm > op.interval {
			op.sum_ -= op.data_[i].Value
			op.ema_sum_ -= op.data_[i].Value * math.Exp(-0.07*float64(time-op.data_[i].Tm)/1000)
			op.count -= 1
			i += 1
		} else {
			op.data_ = op.data_[i:]
			break
		}
	}
}

func (op *OpTimer) Get_data() []valueTm {
	return op.data_
}
func (op *OpTimer) Get_ema_sum() float64 {
	return op.ema_sum_
}
func (op *OpTimer) Sum() float64 {
	return op.sum_
}
func (op *OpTimer) Mean() float64 {
	if op.count != 0 {
		return op.sum_ / op.count
	} else {
		return 0
	}
}
func (op *OpTimer) DiffMean() float64 {
	return op.last_ - op.Mean()
}
func (op *OpTimer) DiffSum() float64 {
	if len(op.data_) > 1 {
		return op.data_[len(op.data_)-1].Value - op.data_[0].Value
	} else {
		return 0
	}

}

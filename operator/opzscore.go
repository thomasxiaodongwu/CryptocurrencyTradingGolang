/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:27
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 22:15:05
 */
package operator

import (
	"math"
	"winter/container"
	"winter/mathpro"
)

type PointZScore struct {
	interval int64
	data_    *container.RingBufferFloat64

	sum_   float64
	sum2_  float64
	zscore float64
}

func (op *PointZScore) Init(interval int64) {
	op.interval = interval
	op.data_ = container.NewRingBufferFloat64(interval)
	op.sum_ = 0
	op.sum2_ = 0
	op.zscore = 0
}

func (op *PointZScore) Update(time int64, x float64) float64 {
	if op.data_.Full() {
		lastx := op.data_.Front()
		op.sum_ -= lastx
		op.sum2_ -= lastx * lastx
	}

	op.sum_ += x
	op.sum2_ += x * x
	op.data_.Push_back(x)

	var n float64 = float64(op.data_.Len())
	std := math.Sqrt((op.sum2_ - op.sum_*op.sum_/n) / (n - 1))
	if mathpro.AlmostEqual(std, 0) || !mathpro.Isfinite(std) {
		op.zscore = 0
	} else {
		op.zscore = (x - op.sum_/n) / std
	}
	return op.zscore
}

func (op *PointZScore) Value() float64 {
	if mathpro.Isfinite(op.zscore) {
		return op.zscore
	} else {
		return 0
	}
}

type TimeZScore struct {
	interval int64
	data_    []valueTm

	sum_   float64
	sum2_  float64
	zscore float64
}

func (op *TimeZScore) Init(interval int64) {
	op.interval = interval
	op.sum_ = 0
	op.sum2_ = 0
	op.zscore = 0
}

func (op *TimeZScore) Update(time int64, x float64) float64 {
	var d valueTm = valueTm{Tm: time, Value: x}
	op.data_ = append(op.data_, d)
	op.sum_ += x
	op.sum2_ += x * x

	var i int = 0
	for i < len(op.data_) && time-op.data_[i].Tm > op.interval {
		op.sum_ -= op.data_[i].Value
		op.sum2_ -= op.data_[i].Value * op.data_[i].Value
		i += 1
	}
	op.data_ = op.data_[i:]

	var n float64 = float64(len(op.data_))
	std := math.Sqrt((op.sum2_ - op.sum_*op.sum_/n) / (n - 1))
	// op.std = mathpro.FastSqrt(op.sum2_/op.count - mean_*mean_)
	if mathpro.AlmostEqual(std, 0) || !mathpro.Isfinite(std) {
		op.zscore = 0
	} else {
		op.zscore = (x - op.sum_/n) / std
	}

	return op.zscore
}

func (op *TimeZScore) Value() float64 {
	if mathpro.Isfinite(op.zscore) {
		return op.zscore
	} else {
		return 0
	}
}

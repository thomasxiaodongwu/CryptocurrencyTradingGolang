/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:27
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 22:46:57
 */
package operator

import (
	"winter/container"
	"winter/mathpro"
)

type PointDiffMean struct {
	Data     *container.RingBufferFloat64
	Interval int64
	sum_     float64
	result_  float64
}

func (op *PointDiffMean) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval
	op.sum_ = 0
	op.result_ = 0
}

// head_ - tail_
func (op *PointDiffMean) Update(time int64, x float64) float64 {
	// Check
	if !mathpro.Isfinite(x) {
		return op.result_
	}

	// 剔除旧的元素
	if op.Data.Full() {
		op.sum_ -= op.Data.Front()
	}

	// 增加新的元素
	op.Data.Push_back(x)
	op.sum_ += x

	op.result_ = x - op.sum_/float64(op.Data.Len())

	return op.result_
}

func (op *PointDiffMean) Value() float64 {
	return op.result_
}

type TimeDiffMean struct {
	Data     []valueTm
	Interval int64
	sum_     float64
	result_  float64
}

func (op *TimeDiffMean) Init(interval int64) {
	op.Interval = interval
	op.sum_ = 0
	op.result_ = 0
}

// head_ - tail_
func (op *TimeDiffMean) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    time,
		Value: x,
	}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		op.sum_ -= op.Data[i].Value
		i += 1
	}
	op.Data = op.Data[i:]

	// 增加新的元素
	op.Data = append(op.Data, new_value)
	op.sum_ += x

	op.result_ = x - op.sum_/float64(len(op.Data))

	return op.result_
}

func (op *TimeDiffMean) Value() float64 {
	return op.result_
}

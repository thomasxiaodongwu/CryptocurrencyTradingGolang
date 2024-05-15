/*
 * @Author: xwu
 * @Date: 2022-05-27 22:48:44
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-27 23:27:27
 */
package operator

import (
	"winter/container"
	"winter/mathpro"
)

type PointDivMean struct {
	Data     *container.RingBufferFloat64
	Interval int64
	sum_     float64
	result_  float64
}

func (op *PointDivMean) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval
	op.sum_ = 0
	op.result_ = 0
}

// head_ - tail_
func (op *PointDivMean) Update(time int64, x float64) float64 {
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

	op.result_ = x / (op.sum_ / float64(op.Data.Len()))

	return op.result_
}

func (op *PointDivMean) Value() float64 {
	return op.result_
}

type TimeDivMean struct {
	Data     []valueTm
	Interval int64
	sum_     float64
	result_  float64
}

func (op *TimeDivMean) Init(interval int64) {
	op.Interval = interval
	op.sum_ = 0
	op.result_ = 0
}

// head_ - tail_
func (op *TimeDivMean) Update(time int64, x float64) float64 {
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

	op.result_ = x / (op.sum_ / float64(len(op.Data)))

	return op.result_
}

func (op *TimeDivMean) Value() float64 {
	return op.result_
}

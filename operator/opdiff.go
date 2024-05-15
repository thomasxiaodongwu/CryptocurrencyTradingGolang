/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:30
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 22:51:43
 */
package operator

import (
	"winter/container"
	"winter/mathpro"
)

type PointDiff struct {
	Data     *container.RingBufferFloat64
	Interval int64
	result_  float64
}

func (op *PointDiff) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval
	op.result_ = 0
}

// head_ - tail_
func (op *PointDiff) Update(time int64, x float64) float64 {
	// Check
	if !mathpro.Isfinite(x) {
		return op.result_
	}
	// 增加新的元素
	op.Data.Push_back(x)

	op.result_ = op.Data.Back() - op.Data.Front()

	return op.result_
}

func (op *PointDiff) Value() float64 {
	return op.result_
}

type TimeDiff struct {
	Data     []valueTm
	Interval int64
	result_  float64
}

func (op *TimeDiff) Init(interval int64) {
	op.Interval = interval
	op.result_ = 0
}

// head_ - tail_
func (op *TimeDiff) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    time,
		Value: x,
	}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		i += 1
	}
	op.Data = op.Data[i:]

	// 增加新的元素
	op.Data = append(op.Data, new_value)

	op.result_ = x - op.Data[0].Value

	return op.result_
}

func (op *TimeDiff) Value() float64 {
	return op.result_
}

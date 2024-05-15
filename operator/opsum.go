/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:03
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 21:03:25
 */
package operator

import "winter/container"

type PointSum struct {
	Data     *container.RingBufferFloat64
	Interval int64
	result_  float64
}

func (op *PointSum) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval
	op.result_ = 0
}

func (op *PointSum) Update(time int64, x float64) float64 {
	// 剔除头部应该去掉的
	if op.Data.Full() {
		op.result_ -= op.Data.Front()
	}

	// 增加新的元素
	op.Data.Push_back(x)
	op.result_ += x

	// 返回均值
	return op.result_
}

func (op *PointSum) Value() float64 {
	return op.result_
}

type TimeSum struct {
	Data     []valueTm
	Interval int64
	result_  float64
}

func (op *TimeSum) Init(interval int64) {
	op.Interval = interval
	op.result_ = 0
}

func (op *TimeSum) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    time,
		Value: x,
	}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		op.result_ -= op.Data[i].Value
		i += 1
	}

	// 增加新的元素
	op.Data = append(op.Data[i:], new_value)
	op.result_ += x

	return op.result_
}

func (op *TimeSum) Value() float64 {
	return op.result_
}

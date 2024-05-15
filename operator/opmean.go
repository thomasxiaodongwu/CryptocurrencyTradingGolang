/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:03
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 17:51:06
 */
package operator

import (
	"winter/container"
)

type PointMean struct {
	Data     *container.RingBufferFloat64
	Interval int64
	Sum      float64
	result_  float64
}

func (op *PointMean) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval
	op.Sum = 0
	op.result_ = 0
}

func (op *PointMean) Update(time int64, x float64) float64 {
	// 剔除头部应该去掉的
	if op.Data.Full() {
		op.Sum -= op.Data.Front()
	}

	// 增加新的元素
	op.Data.Push_back(x)
	op.Sum += x

	op.result_ = op.Sum / float64(op.Data.Len())
	// 返回均值
	return op.result_
}

func (op *PointMean) Value() float64 {
	return op.result_
}

type TimeMean struct {
	Data     []valueTm
	Interval int64
	Sum      float64
	result_  float64
}

func (op *TimeMean) Init(interval int64) {
	op.Interval = interval
	op.Sum = 0
	op.result_ = 0
}

func (op *TimeMean) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    time,
		Value: x,
	}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		op.Sum -= op.Data[i].Value
		i += 1
	}

	// 增加新的元素
	op.Data = append(op.Data[i:], new_value)
	op.Sum += x

	// 返回均值
	op.result_ = op.Sum / float64(len(op.Data))

	return op.result_
}

func (op *TimeMean) Value() float64 {
	return op.result_
}

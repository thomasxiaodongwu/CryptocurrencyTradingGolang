/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:36
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 18:26:17
 */
package operator

import (
	"winter/container"
)

type PointCov struct {
	Data_X *container.RingBufferFloat64
	Data_Y *container.RingBufferFloat64

	Interval int64

	x_    float64
	y_    float64
	xy_   float64
	count float64

	Result_ float64
}

func (op *PointCov) Init(interval int64) {
	op.Interval = interval

	op.Data_X = container.NewRingBufferFloat64(interval)
	op.Data_Y = container.NewRingBufferFloat64(interval)

	op.x_ = 0
	op.y_ = 0
	op.xy_ = 0
	op.count = 0

	op.Result_ = 0
}

// TODO: 这里的count不知道是不是会因为时间过长导致误差累计，所以可能需要每隔一段时间进行一次清空
func (op *PointCov) Update(time int64, x float64, y float64) float64 {
	// 这里是删除旧的元素
	if op.Data_X.Full() && op.Data_Y.Full() {
		oldx := op.Data_X.Front()
		oldy := op.Data_Y.Front()

		op.x_ -= oldx
		op.y_ -= oldy
		op.xy_ -= oldx * oldy
		op.count -= 1
	}

	// 增加新的元素
	op.Data_X.Push_back(x)
	op.Data_Y.Push_back(y)

	op.x_ += x
	op.y_ += y
	op.xy_ += x * y
	op.count += 1

	op.Result_ = op.xy_/op.count - op.x_/op.count*op.y_/op.count

	return op.Result_
}

func (op *PointCov) Value() float64 {
	return op.Result_
}

type TimeCov struct {
	Data []valuesTm

	Interval int64

	x_    float64
	y_    float64
	xy_   float64
	count float64

	Result_ float64
}

func (op *TimeCov) Init(interval int64) {
	op.Interval = interval

	op.x_ = 0
	op.y_ = 0
	op.xy_ = 0
	op.count = 0

	op.Result_ = 0
}

// TODO: 这里的count不知道是不是会因为时间过长导致误差累计，所以可能需要每隔一段时间进行一次清空
func (op *TimeCov) Update(time int64, x float64, y float64) float64 {
	d := valuesTm{Tm: time, Value1: x, Value2: y}

	// 这里是删除旧的元素
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		op.x_ -= op.Data[i].Value1
		op.y_ -= op.Data[i].Value2
		op.xy_ -= op.Data[i].Value1 * op.Data[i].Value2
		op.count -= 1
		i += 1
	}

	// new elements
	op.Data = append(op.Data[i:], d)

	op.x_ += x
	op.y_ += y
	op.xy_ += x * y
	op.count += 1

	op.Result_ = op.xy_/op.count - op.x_/op.count*op.y_/op.count

	return op.Result_
}

func (op *TimeCov) Value() float64 {
	return op.Result_
}

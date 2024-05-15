/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:33
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 21:37:43
 */
package operator

import (
	"winter/container"
	"winter/mathpro"
)

type PointVariance struct {
	Data     *container.RingBufferFloat64
	Interval int64
	Sum_     float64
	Sum2_    float64
	result_  float64
}

func (op *PointVariance) Init(interval int64) {
	op.Data = container.NewRingBufferFloat64(interval)
	op.Interval = interval
	op.Sum_ = 0
	op.Sum2_ = 0
}

// x, x2
func (op *PointVariance) Update(time int64, x float64) float64 {
	// 剔除头部应该去掉的
	if op.Data.Full() {
		old_v := op.Data.Front()
		op.Sum_ -= old_v
		op.Sum2_ -= old_v * old_v
	}

	// 增加新的元素
	op.Data.Push_back(x)
	op.Sum_ += x
	op.Sum2_ += x * x

	// sqrtf((sum2 - sum * sum / (float) count) / (float) (count - 1));
	// 返回链表的第一个元素
	var n float64 = float64(op.Data.Len())
	op.result_ = (op.Sum2_ - op.Sum_*op.Sum_/n) / (n - 1)
	if !mathpro.Isfinite(op.result_) {
		op.result_ = 0
	}
	return op.result_
}

func (op *PointVariance) Value() float64 {
	return op.result_
}

type TimeVariance struct {
	// Data     *list.List
	Data     []valueTm
	Interval int64
	Sum_     float64
	Sum2_    float64
	result_  float64
}

func (op *TimeVariance) Init(interval int64) {
	op.Interval = interval
	op.Sum_ = 0
	op.Sum2_ = 0
	op.result_ = 0
}

// x, x2
func (op *TimeVariance) Update(time int64, x float64) float64 {
	new_value := valueTm{Tm: time, Value: x}

	// 剔除头部应该去掉的
	var i int = 0
	for i < len(op.Data) && time-op.Data[i].Tm > op.Interval {
		old_v := op.Data[i].Value
		op.Sum_ -= old_v
		op.Sum2_ -= old_v * old_v
		i += 1
	}

	// 增加新的元素
	op.Data = append(op.Data[i:], new_value)
	op.Sum_ += x
	op.Sum2_ += x * x

	n := float64(len(op.Data))

	// sqrtf((sum2 - sum * sum / (float) count) / (float) (count - 1));
	// 返回链表的第一个元素
	op.result_ = (op.Sum2_ - op.Sum_*op.Sum_/n) / (n - 1)
	if !mathpro.Isfinite(op.result_) {
		op.result_ = 0
	}
	return op.result_
}

func (op *TimeVariance) Value() float64 {
	return op.result_
}

// type TimeVariance struct {
// 	Data     *list.List
// 	Interval int64
// 	Sum_     float64
// 	Sum2_    float64
// 	Result_  float64
// }

// // x, x2
// func (op *TimeVariance) Update(time int64, x float64) float64 {
// 	new_value := &valueTm{
// 		Tm:    time,
// 		Value: x,
// 	}

// 	// 剔除头部应该去掉的
// 	for {
// 		if op.Data.Len() > 0 && time-op.Data.Front().Value.(*valueTm).Tm >= op.Interval {
// 			old_v := op.Data.Front().Value.(*valueTm).Value
// 			op.Sum_ -= old_v
// 			op.Sum2_ -= old_v * old_v
// 			op.Data.Remove(op.Data.Front())
// 		} else {
// 			break
// 		}
// 	}

// 	// 增加新的元素
// 	op.Data.PushBack(new_value)
// 	op.Sum_ -= x
// 	op.Sum2_ -= x * x

// 	n := float64(op.Data.Len())

// 	// sqrtf((sum2 - sum * sum / (float) count) / (float) (count - 1));
// 	// 返回链表的第一个元素
// 	op.Result_ = (op.Sum2_ - op.Sum_*op.Sum_/n) / (n - 1)

// 	return op.Result_
// }

// func (op *TimeVariance) Value() float64 {
// 	return op.Result_
// }

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:56
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 18:05:26
 */
package operator

type PointMin struct {
	Data     []valueTm
	Interval int64
	Size     int64
}

func (op *PointMin) Init(interval int64) {
	op.Interval = interval
	op.Size = 0
}

func (op *PointMin) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    op.Size,
		Value: x,
	}

	// 从尾部剔除小于新元素的,非严格单调递减
	var i int = len(op.Data)
	for i > 0 && op.Data[i-1].Value > x {
		i -= 1
	}
	op.Data = append(op.Data[:i], new_value)

	// 队头的元素
	if op.Data[0].Tm+op.Interval == op.Size {
		op.Data = op.Data[1:]
	}

	op.Size += 1

	// 返回链表的第一个元素
	return op.Data[0].Value
}

func (op *PointMin) Value() float64 {
	if len(op.Data) > 0 {
		return op.Data[0].Value
	} else {
		return 0
	}
}

type TimeMin struct {
	Data     []valueTm
	Interval int64
}

func (op *TimeMin) Init(interval int64) {
	op.Interval = interval
}

func (op *TimeMin) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    time,
		Value: x,
	}

	// 从尾部剔除小于新元素的,非严格单调递减
	var i int = len(op.Data)
	for i > 0 && op.Data[i-1].Value > x {
		i -= 1
	}
	op.Data = append(op.Data[:i], new_value)

	// 剔除头部应该去掉的
	i = 0
	for time-op.Data[i].Tm > op.Interval {
		i += 1
	}
	op.Data = op.Data[i:]

	// 返回链表的第一个元素
	if len(op.Data) > 0 {
		return op.Data[0].Value
	} else {
		return 0
	}
}

func (op *TimeMin) Value() float64 {
	if len(op.Data) > 0 {
		return op.Data[0].Value
	} else {
		return 0
	}
}

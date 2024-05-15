/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:07
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 17:20:21
 */
package operator

type PointMax struct {
	Data     []valueTm
	Interval int64
	Size     int64
}

func (op *PointMax) Init(interval int64) {
	op.Interval = interval
	op.Size = 0
}

func (op *PointMax) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    op.Size,
		Value: x,
	}

	// 从尾部剔除小于新元素的,非严格单调递减
	var i int = len(op.Data)
	for i > 0 && op.Data[i-1].Value < x {
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

func (op *PointMax) Value() float64 {
	if len(op.Data) > 0 {
		return op.Data[0].Value
	} else {
		return 0
	}
}

type TimeMax struct {
	Data     []valueTm
	Interval int64
}

func (op *TimeMax) Init(interval int64) {
	op.Interval = interval
}

func (op *TimeMax) Update(time int64, x float64) float64 {
	new_value := valueTm{
		Tm:    time,
		Value: x,
	}

	// 从尾部剔除小于新元素的,非严格单调递减
	var i int = len(op.Data)
	for i > 0 && op.Data[i-1].Value < x {
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

func (op *TimeMax) Value() float64 {
	if len(op.Data) > 0 {
		return op.Data[0].Value
	} else {
		return 0
	}
}

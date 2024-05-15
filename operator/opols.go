/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:52
 * @Last Modified by: xwu
 * @Last Modified time: 2022-01-08 21:16:10
 */
package operator

type TimeOLS struct {
	interval int64
	data_x_  []float64
	data_y_  []float64
	tm       []int64

	count  float64
	sum_x  float64
	sum_y  float64
	sum_xy float64
	sum_x2 float64

	a float64
	b float64
	r float64
}

func (op *TimeOLS) Init(interval int64) {
	op.interval = interval
	op.count = 0
	op.sum_x = 0
	op.sum_y = 0
	op.sum_xy = 0
	op.sum_x2 = 0
}

func (op *TimeOLS) Update(time int64, x float64, y float64) {
	op.tm = append(op.tm, time)
	op.data_x_ = append(op.data_x_, x)
	op.data_y_ = append(op.data_y_, y)

	op.count += 1
	op.sum_x += x
	op.sum_y += y
	op.sum_xy += x * y
	op.sum_x2 += x * x

	var i int = 0
	for len(op.tm) > 0 {
		if time-op.tm[i] > op.interval {
			op.sum_x -= op.data_x_[i]
			op.sum_y -= op.data_y_[i]
			op.sum_xy -= op.data_x_[i] * op.data_y_[i]
			op.sum_x2 -= op.data_x_[i] * op.data_x_[i]
			op.count -= 1
			i += 1
		} else {
			op.tm = op.tm[i:]
			op.data_x_ = op.data_x_[i:]
			op.data_y_ = op.data_y_[i:]
			break
		}
	}

	op.a = (op.sum_xy - op.sum_x*op.sum_y/op.count) / (op.sum_x2 - op.sum_x*op.sum_x/op.count)
	op.b = op.sum_y/op.count - op.a*op.sum_x/op.count
	op.r = y - op.a*x - op.b
}

func (op *TimeOLS) Coef_X() float64 {
	return op.a
}
func (op *TimeOLS) B() float64 {
	return op.b
}
func (op *TimeOLS) Residual() float64 {
	return op.r
}

// TODO: 没有修改结束，是直接复制的timeols
type PointOLS struct {
	interval int64
	data_x_  []float64
	data_y_  []float64
	tm       []int64

	count  float64
	sum_x  float64
	sum_y  float64
	sum_xy float64
	sum_x2 float64

	a float64
	b float64
	r float64
}

func (op *PointOLS) Init(interval int64) {
	op.interval = interval
	op.count = 0
	op.sum_x = 0
	op.sum_y = 0
	op.sum_xy = 0
	op.sum_x2 = 0
}

func (op *PointOLS) Update(time int64, x float64, y float64) {
	op.tm = append(op.tm, time)
	op.data_x_ = append(op.data_x_, x)
	op.data_y_ = append(op.data_y_, y)

	op.count += 1
	op.sum_x += x
	op.sum_y += y
	op.sum_xy += x * y
	op.sum_x2 += x * x

	var i int = 0
	for len(op.tm) > 0 {
		if time-op.tm[i] > op.interval {
			op.sum_x -= op.data_x_[i]
			op.sum_y -= op.data_y_[i]
			op.sum_xy -= op.data_x_[i] * op.data_y_[i]
			op.sum_x2 -= op.data_x_[i] * op.data_x_[i]
			op.count -= 1
			i += 1
		} else {
			op.tm = op.tm[i:]
			op.data_x_ = op.data_x_[i:]
			op.data_y_ = op.data_y_[i:]
			break
		}
	}

	op.a = (op.sum_xy - op.sum_x*op.sum_y/op.count) / (op.sum_x2 - op.sum_x*op.sum_x/op.count)
	op.b = op.sum_y/op.count - op.a*op.sum_x/op.count
	op.r = y - op.a*x - op.b
}

func (op *PointOLS) Coef_X() float64 {
	return op.a
}
func (op *PointOLS) B() float64 {
	return op.b
}
func (op *PointOLS) Residual() float64 {
	return op.r
}

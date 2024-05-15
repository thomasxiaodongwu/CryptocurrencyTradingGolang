/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:20
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-26 22:42:54
 */
package operator

import (
	"math"
)

type PointEMA struct {
	Halflife_Decay float64
	alpha          float64
	first_saw      bool
	Result_        float64
}

func (op *PointEMA) Init(interval int64) {
	op.Set_Decay_Alpha(float64(interval))
	op.first_saw = true
}

func (op *PointEMA) Set_Decay_Alpha(Halflife_Decay float64) {
	op.Halflife_Decay = Halflife_Decay
	op.alpha = math.Exp(-math.Log(2) / op.Halflife_Decay)
	op.Result_ = 0
}

// time 是ms
func (op *PointEMA) Update(time int64, x float64) float64 {
	if op.first_saw {
		op.Result_ = x
		op.first_saw = false
	} else {
		op.Result_ = op.Result_*op.alpha + x*(1-op.alpha)
	}
	return op.Result_
}

func (op *PointEMA) Value() float64 {
	return op.Result_
}

type TimeEMA struct {
	LastTm         int64
	Halflife_Decay float64 // 单位是s
	alpha          float64
	Result_        float64
}

func (op *TimeEMA) Init(interval int64) {
	op.Set_Decay_Alpha(float64(interval))
}

func (op *TimeEMA) Set_Decay_Alpha(Halflife_Decay float64) {
	op.Halflife_Decay = Halflife_Decay
	op.alpha = -math.Log(2) / op.Halflife_Decay
	op.Result_ = 0
	op.LastTm = -1
}

// time 是ms
func (op *TimeEMA) Update(time int64, x float64) float64 {
	if op.LastTm < 0 {
		op.Result_ = x
		op.LastTm = time
	} else {
		alpha := math.Exp(op.alpha * float64(time-op.LastTm))
		op.Result_ = op.Result_*alpha + x*(1-alpha)
		op.LastTm = time
	}
	return op.Result_
}

func (op *TimeEMA) Value() float64 {
	return op.Result_
}

/*
 * @Author: xwu
 * @Date: 2022-03-09 21:44:35
 * @Last Modified by: xwu
 * @Last Modified time: 2022-03-10 13:04:15
 */
package operator

type OpAdjustMid struct {
	Size    int64
	Result_ float64
}

func (op *OpAdjustMid) Init(interval int64) {
	op.Size = interval
	op.Result_ = 0
}

func (op *OpAdjustMid) Update(Data [][]float64) float64 {
	var acc_value float64 = 0
	var acc_vol float64 = 0
	var i int
	for i = range Data {
		acc_value += Data[i][0] * Data[i][1]
		acc_vol += Data[i][1]
		if acc_value > float64(op.Size) {
			break
		}
	}
	op.Result_ = acc_value / acc_vol
	op.Result_ = Data[i][0]
	return op.Result_
}

func (op *OpAdjustMid) Value() float64 {
	return op.Result_
}

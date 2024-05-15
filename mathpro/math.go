/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:33
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-28 22:36:29
 */
package mathpro

import (
	"math"
	"time"
)

func AlmostEqual(x float64, y float64) bool {
	if math.Abs(x-y) < 1e-6 {
		return true
	} else {
		return false
	}
}

func AlmostGreater(x float64, y float64) bool {
	if x > y+1e-6 {
		return true
	} else {
		return false
	}
}

func AlmostLess(x float64, y float64) bool {
	if x < y-1e-6 {
		return true
	} else {
		return false
	}
}

func Isfinite(x float64) bool {
	if !math.IsNaN(x) && !math.IsInf(x, 0) {
		return true
	} else {
		return false
	}
}

func GetMean(x []float64) float64 {
	var ans, count float64 = 0, 0
	for i := range x {
		if Isfinite(x[i]) {
			ans += x[i]
			count += 1
		}
	}
	return ans / count
}

func GetMax(x []float64) float64 {
	var ans float64 = x[0]
	for i := range x {
		if x[i] > ans {
			ans = x[i]
		}
	}
	return ans
}

func GetMin(x []float64) float64 {
	var ans float64 = x[0]
	for i := range x {
		if x[i] < ans {
			ans = x[i]
		}
	}
	return ans
}

func GetSum(x []float64) float64 {
	var ans float64 = 0
	for i := range x {
		if Isfinite(x[i]) {
			ans += x[i]
		}
	}
	return ans
}

func Power(x float64, n int) float64 {
	var ans float64 = 1
	for i := 0; i < n; i += 1 {
		ans *= x
	}
	return ans
}

func HigherOriginMoment(x []float64, n int) float64 {
	var ans, count float64 = 0, 0
	for i := range x {
		if Isfinite(x[i]) {
			ans += Power(x[i], n)
			count += 1
		}
	}
	return ans / count
}

func HigherCentralMoment(x []float64, n int) float64 {
	var ans, count float64 = 0, 0
	var m float64 = GetMean(x)
	for i := range x {
		if Isfinite(x[i]) {
			ans += Power(x[i]-m, n)
			count += 1
		}
	}
	return ans / count
}

// func Variance(x []float64) float64 {
// 	var M, S, count float64 = 0, 0, 0
// 	for i := range x {
// 		if Isfinite(x[i]) {
// 			count += 1
// 			if !AlmostEqual(count, 1) {
// 				M = x[i]
// 				S = 0

// 			} else {
// 				_M := M + (x[i]-M)/count
// 				S = S + (x[i]-_M)*(x[i]-M)
// 				M = _M

// 			}
// 		}
// 	}
// 	return S / (count - 1)
// }

func GetVariance(x []float64) float64 {
	var M1, M2, count float64 = 0, 0, 0
	for i := range x {
		if Isfinite(x[i]) {
			M1 += x[i]
			M2 += x[i] * x[i]
			count += 1
		}
	}
	return (M2 - M1*M1/count) / (count - 1)
}

func GetStd(x []float64) float64 {
	v := GetVariance(x)
	if v > 1e-6 {
		return math.Sqrt(v)
	} else {
		return 0
	}
}

func GetAbsVector(x []float64) []float64 {
	var ans []float64 = make([]float64, len(x))
	for i := range x {
		ans[i] = math.Abs(x[i])
	}
	return ans
}

// TODO:需要double check
// func Corr(x []float64, y []float64) float64 {
// 	// 首先检查是否相等x和y
// 	if len(x) != len(y) {
// 		logger.Error("length of x not equal length of y")
// 	}

// 	// 检查x和y长度是否大于2
// 	if len(x) <= 2 && len(y) <= 2 {
// 		logger.Error("length of x or y less 3")
// 	}

// 	// 首先进行初始化
// 	var Mx float64 = x[0]
// 	var Sx float64 = 0
// 	var My float64 = y[0]
// 	var Sy float64 = 0
// 	var Mxy float64 = Mx * My
// 	var count float64 = 1
// 	var length int = len(x)

// 	for i := 1; i < length; i++ {
// 		if Isfinite(x[i]) && Isfinite(y[i]) {
// 			count += 1

// 			_Mx := Mx + (x[i]-Mx)/count
// 			_My := My + (y[i]-My)/count
// 			Sx = Sx + (x[i]-_Mx)*(x[i]-Mx)
// 			Sy = Sy + (y[i]-_My)*(y[i]-My)
// 			Mx = _Mx
// 			My = _My
// 			Mxy = Mxy + (x[i]*y[i]-Mxy)/count
// 		}
// 	}

//		return (Mxy - Mx*My) * FastInvSqrt64(Sx*Sy) * (count - 1)
//	}
func WinRate(x []float64, y []float64) float64 {
	// 首先检查是否相等x和y
	if len(x) != len(y) {
		panic("length of x not equal length of y")
	}

	var ans float64 = 0

	for i := range x {
		if x[i] > 0 && y[i] > 0 {
			ans += 1
		} else if x[i] < 0 && y[i] < 0 {
			ans += 1
		} else if AlmostEqual(x[i], 0) && AlmostEqual(y[i], 0) {
			ans += 1
		}
	}

	return ans / float64(len(x))
}

func Corr(x []float64, y []float64) float64 {
	// 首先检查是否相等x和y
	if len(x) != len(y) {
		panic("length of x not equal length of y")
	}

	// 检查x和y长度是否大于2
	if len(x) <= 2 && len(y) <= 2 {
		panic("length of x or y less 3")
	}

	// 首先进行初始化
	var x1, x2, y1, y2, xy, count float64 = 0, 0, 0, 0, 0, 0

	for i := range x {
		if Isfinite(x[i]) && Isfinite(y[i]) {
			x1 += x[i]
			x2 += x[i] * x[i]
			y1 += y[i]
			y2 += y[i] * y[i]
			xy += x[i] * y[i]
			count += 1
		}
	}

	varx := (x2 - x1*x1/count) / (count - 1)
	vary := (y2 - y1*y1/count) / (count - 1)

	if AlmostGreater(varx, 0) && AlmostGreater(vary, 0) {
		return (xy/count - x1*y1/count/count) / math.Sqrt(varx) / math.Sqrt(vary)
	} else {
		return 0
	}
}

func Trim(x float64, upper_limit float64, lower_limit float64) float64 {
	return 0
}

func Absmax(x float64, y float64) float64 {
	var xydiff float64 = x - y
	var denom float64 = math.Max(math.Abs(x), math.Abs(y))
	return xydiff / denom
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// /////////
func Get_return(tm []int64, price []float64, windows int64) []float64 {
	// windows单位ms
	var ret []float64
	var start, end int = 0, 0
	for {
		if start == len(tm) {
			break
		}

		if tm[end]-tm[start] < windows && end < len(tm)-1 {
			end += 1
		} else {
			ret = append(ret, (price[end]/price[start]-1)*10000)
			start += 1
		}
	}

	if len(ret) != len(price) {
		logger.Info("length of ret not equal length of price.")
	}

	return ret
}

func Get_tick_return(tm []int64, price []float64, windows int64) []float64 {
	var ret []float64 = make([]float64, len(tm))
	for i := range price {
		var e int = i + int(windows)
		if e > len(price)-1 {
			e = len(price) - 1
		}
		ret[i] = price[e]/price[i] - 1
	}

	return ret
}

func Get_fixvol_return(tm []int64, price []float64, size []float64, windows int64) []float64 {
	var ret []float64 = make([]float64, len(tm))
	for i := range price {
		var e int = i + int(windows)
		if e > len(price)-1 {
			e = len(price) - 1
		}
		ret[i] = price[e]/price[i] - 1
	}

	return ret
}

func Get_corr_price(tm []int64, values [][]float64, price []float64, window int64, delay int64) []float64 {
	// 计算收益率
	var ans []float64 = make([]float64, len(values[0]))

	var ret1 []float64 = Get_return(tm, price, window)
	var ret2 []float64 = Get_return(tm, price, delay)
	var ret []float64 = make([]float64, len(ret1))
	for i := range ret2 {
		ret[i] = ret1[i] - ret2[i]
	}

	for i := range values[0] {
		var _value []float64
		for j := range values {
			_value = append(_value, values[j][i])
		}
		ans[i] = Corr(_value, ret)
	}
	return ans
}

func Get_corr_abs_price(tm []int64, values [][]float64, price []float64, window int64, delay int64) []float64 {
	// 计算收益率
	var ans []float64 = make([]float64, len(values[0]))

	var ret1 []float64 = Get_return(tm, price, window)
	var ret2 []float64 = Get_return(tm, price, delay)
	var ret []float64 = make([]float64, len(ret1))
	for i := range ret2 {
		ret[i] = ret1[i] - ret2[i]
	}
	ret = GetAbsVector(ret)

	for i := range values[0] {
		var _value []float64
		for j := range values {
			_value = append(_value, values[j][i])
		}
		ans[i] = Corr(_value, ret)
	}
	return ans
}

func Get_ir_price(tm []int64, values [][]float64, price []float64, freq string, delay int64) [][]float64 {
	var ans [][]float64

	var init_freq int
	switch freq {
	case "H":
		init_freq = time.Unix(tm[0]/1000, 0).Hour()
	default:
		logger.Info("params freq wrong")
	}

	var _tm []int64
	var _values [][]float64
	var _price []float64

	var end int = 0
	for {
		if end == len(tm) {
			ans = append(ans, Get_corr_price(_tm, _values, _price, 3000+delay, delay))
			break
		}
		if time.Unix(tm[end]/1000, 0).Hour() == init_freq {
			_tm = append(_tm, tm[end])
			_values = append(_values, values[end])
			_price = append(_price, price[end])
			end += 1
		} else {
			ans = append(ans, Get_corr_price(_tm, _values, _price, 3000+delay, delay))
			init_freq = time.Unix(tm[end]/1000, 0).Hour()
			_tm = _tm[:0]
			_values = _values[:0]
			_price = _price[:0]
		}
	}

	var results [][]float64
	for i := range values[0] {
		var _ans []float64
		for j := range ans {
			_ans = append(_ans, ans[j][i])
		}
		var res []float64
		res = append(res, GetMean(_ans))
		res = append(res, GetMean(_ans)/GetStd(_ans))
		results = append(results, res)
	}

	return results
}

func Get_ir_abs_price(tm []int64, values [][]float64, price []float64, freq string, delay int64) [][]float64 {
	var ans [][]float64

	var init_freq int
	switch freq {
	case "H":
		init_freq = time.Unix(tm[0]/1000, 0).Hour()
	default:
		logger.Info("params freq wrong")
	}

	var _tm []int64
	var _values [][]float64
	var _price []float64

	var end int = 0
	for {
		if end == len(tm) {
			ans = append(ans, Get_corr_abs_price(_tm, _values, _price, 3000+delay, delay))
			break
		}
		if time.Unix(tm[end]/1000, 0).Hour() == init_freq {
			_tm = append(_tm, tm[end])
			_values = append(_values, values[end])
			_price = append(_price, price[end])
			end += 1
		} else {
			ans = append(ans, Get_corr_abs_price(_tm, _values, _price, 3000+delay, delay))
			init_freq = time.Unix(tm[end]/1000, 0).Hour()
			_tm = _tm[:0]
			_values = _values[:0]
			_price = _price[:0]
		}
	}

	var results [][]float64
	for i := range values[0] {
		var _ans []float64
		for j := range ans {
			_ans = append(_ans, ans[j][i])
		}
		var res []float64
		res = append(res, GetMean(_ans))
		res = append(res, GetMean(_ans)/GetStd(_ans))
		results = append(results, res)
	}

	return results
}

package mathpro

import "math"

func SafeDivide(x float64, y float64) float64 {
	if AlmostEqual(y, 0) {
		return 0
	} else {
		return x / y
	}
}

func SafeSqrt(x float64) float64 {
	var ans float64 = 0
	if x > 0 {
		ans = math.Sqrt(x)
	} else if AlmostEqual(x, 0) {
		ans = 0
	} else if x < 0 {
		panic("sqrt negative number")
	}
	return ans
}

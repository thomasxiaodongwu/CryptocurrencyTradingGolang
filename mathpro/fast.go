package mathpro

import "math"

const (
	magic64 = 0x5FE6EB50C7B537A9
)

// 还是math包里面的快
func FastInvSqrt64(n float64) float64 {
	if n < 0 {
		return math.NaN()
	}

	n2, th := n*0.5, float64(1.5)
	b := math.Float64bits(n)
	b = magic64 - (b >> 1)
	f := math.Float64frombits(b)

	// Newton step, repeating increases accuracy
	f *= th - (n2 * f * f)
	f *= th - (n2 * f * f)
	f *= th - (n2 * f * f)

	return f
}

func FastSqrt(n float64) float64 {
	return float64(1) / FastInvSqrt64(n)
}

func FastTanh(x float64) float64 {
	var x2 float64 = x * x
	a := x * (135135.0 + x2*(17325.0+x2*(378.0+x2)))
	b := 135135.0 + x2*(62370.0+x2*(3150.0+x2*28.0))
	return a / b
}

func FastExp(x float64) float64 {
	if AlmostLess(x, 0) {
		return 1 / (1 + x*x)
	} else {
		return math.Exp(x)
	}
}

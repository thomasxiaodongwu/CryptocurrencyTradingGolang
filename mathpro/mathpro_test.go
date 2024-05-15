package mathpro

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

/*
fasttanh 快非常多
*/

var (
	x []float64
	y []float64
)

func init() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {
		x = append(x, rand.Float64())
	}
	for i := 0; i < 100; i++ {
		y = append(y, rand.Float64())
	}
}
func Benchmark_AlmostEqual(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AlmostEqual(float64(i), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Isfinite(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Isfinite(float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Mean(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetMean(x)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Sum(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetSum(x)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_HigherOriginMoment(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HigherOriginMoment(x, 3)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_HigherCentralMoment(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HigherCentralMoment(x, 3)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Variance(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVariance(x)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetStd(x)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Corr(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Corr(x, y)
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Absmax(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Absmax(float64(i), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Sigmoid(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sigmoid(float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Sqrt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math.Sqrt(float64(i))
	}
	b.SetBytes(int64(b.N))
}
func Benchmark_FastSqrt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FastSqrt(float64(i))
		// FastInvSqrt64(float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_Tanh(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		math.Tanh(float64(i))
	}
	b.SetBytes(int64(b.N))
}
func Benchmark_FastTanh(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FastTanh(float64(i))
		// FastInvSqrt64(float64(i))
	}
	b.SetBytes(int64(b.N))
}

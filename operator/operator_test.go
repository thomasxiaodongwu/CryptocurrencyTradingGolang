/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:14
 * @Last Modified by: xwu
 * @Last Modified time: 2022-02-27 11:36:25
 */
package operator

import (
	"math"
	"math/rand"
	"testing"
	"time"
	"winter/mathpro"
)

var Localtime []int64
var X []float64
var Y []float64
var N int = 100000

func init() {
	Localtime = make([]int64, N)
	X = make([]float64, N)
	Y = make([]float64, N)
	for i := 0; i < N; i++ {
		Localtime[i] = int64(100 * i)

		rand.Seed(time.Now().UnixNano())
		X[i] = rand.NormFloat64()
		rand.Seed(time.Now().UnixNano())
		Y[i] = rand.NormFloat64()
	}
}

func Benchmark_PointCorr(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointCorr(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i), float64(i))
	}
	b.SetBytes(int64(b.N))
}
func Benchmark_TimeCorr(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeCorr(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointCov(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointCov(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeCov(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeCov(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointDecay(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointDecay(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeDecay(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeDecay(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointDiff(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointDiff(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeDiff(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeDiff(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointDiffMean(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointDiffMean(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeDiffMean(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeDiffMean(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointDiffMeanPlus(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointDiffMeanPlus(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeDiffMeanPlus(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeDiffMeanPlus(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointEMA(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointEMA(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeEMA(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeEMA(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointMax(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointMax(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeMax(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeMax(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointMean(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointMean(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeMean(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeMean(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointMin(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointMin(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeMin(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeMin(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointStd(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointStd(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeStd(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeStd(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_PointSum(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointSum(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeSum(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeSum(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

// Variance
func Benchmark_PointVariance(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointVariance(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeVariance(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeVariance(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

// ZScore
func Benchmark_PointZScore(b *testing.B) {
	b.ReportAllocs()

	ans := NewPointZScore(30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

func Benchmark_TimeZScore(b *testing.B) {
	b.ReportAllocs()

	ans := NewTimeZScore(3000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ans.Update(int64(i*100), float64(i))
	}
	b.SetBytes(int64(b.N))
}

/*
=================================================================
=================================================================
=================================================================
=================================================================
=================================================================
*/
func Test_PointCorr(t *testing.T) {
	var cache []float64
	var cache1 []float64
	var cache2 []float64
	ans := NewPointCorr(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		cache1 = append(cache1, Y[i])
		cache2 = append(cache2, X[i]*Y[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		if len(cache1) == 31 {
			cache1 = cache1[1:]
		}
		if len(cache2) == 31 {
			cache2 = cache2[1:]
		}
		r := ans.Update(int64(i*100), X[i], Y[i])
		s := mathpro.Corr(cache, cache1)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointCorr ", r, " ", s)
		}
	}
}

func Test_TimeCorr(t *testing.T) {
	var cache []float64
	var cache1 []float64
	var cache2 []float64
	var tm []int64
	ans := NewTimeCorr(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		cache1 = append(cache1, Y[i])
		cache2 = append(cache2, X[i]*Y[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			cache1 = cache1[1:]
			cache2 = cache2[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i], Y[i])
		s := mathpro.Corr(cache, cache1)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeCorr ", r, " ", s)
		}
	}
}

func Test_PointCov(t *testing.T) {
	var cache []float64
	var cache1 []float64
	var cache2 []float64
	ans := NewPointCov(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		cache1 = append(cache1, Y[i])
		cache2 = append(cache2, X[i]*Y[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		if len(cache1) == 31 {
			cache1 = cache1[1:]
		}
		if len(cache2) == 31 {
			cache2 = cache2[1:]
		}
		r := ans.Update(int64(i*100), X[i], Y[i])
		s := mathpro.GetMean(cache2) - mathpro.GetMean(cache)*mathpro.GetMean(cache1)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointCov ", r, " ", s)
		}
	}
}

func Test_TimeCov(t *testing.T) {
	var cache []float64
	var cache1 []float64
	var cache2 []float64
	var tm []int64
	ans := NewTimeCov(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		cache1 = append(cache1, Y[i])
		cache2 = append(cache2, X[i]*Y[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			cache1 = cache1[1:]
			cache2 = cache2[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i], Y[i])
		s := mathpro.GetMean(cache2) - mathpro.GetMean(cache)*mathpro.GetMean(cache1)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeCov ", r, " ", s)
		}
	}
}

func Test_PointDecay(t *testing.T) {
	var cache []float64
	ans := NewPointDecay(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		var s, _s float64 = 0, 0
		for i := range cache {
			s += (float64(i) + 1) * cache[i]
			_s += (float64(i) + 1)
		}
		s = s / _s
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointDecay ", r, " ", s)
		}
	}
}

func Test_TimeDecay(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeDecay(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		var s, _s float64 = 0, 0
		for i := range cache {
			s += (float64(i) + 1) * cache[i]
			_s += (float64(i) + 1)
		}
		s = s / _s
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeDecay ", r, " ", s)
		}
	}
}

func Test_PointDiff(t *testing.T) {
	var cache []float64
	ans := NewPointDiff(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := cache[len(cache)-1] - cache[0]
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointDiff ", r, " ", s)
		}
	}
}

func Test_TimeDiff(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeDiff(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := cache[len(cache)-1] - cache[0]
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeDiff ", r, " ", s)
		}
	}
}

func Test_PointDiffMean(t *testing.T) {
	var cache []float64
	ans := NewPointDiffMean(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := X[i] - mathpro.GetMean(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointDiffMean ", r, " ", s)
		}
	}
}

func Test_TimeDiffMean(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeDiffMean(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := X[i] - mathpro.GetMean(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeDiffMean ", r, " ", s)
		}
	}
}

func Test_PointDiffMeanPlus(t *testing.T) {
	var cache []float64
	ans := NewPointDiffMeanPlus(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := 2*X[i] - mathpro.GetMean(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointDiffMeanPlus ", r, " ", s)
		}
	}
}

func Test_TimeDiffMeanPlus(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeDiffMeanPlus(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := 2*X[i] - mathpro.GetMean(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeDiffMeanPlus ", r, " ", s)
		}
	}
}

func Test_PointEMA(t *testing.T) {
	var cache []float64
	ans := NewPointEMA(30)
	var s float64 = 0
	var alpha float64 = math.Exp(-math.Log(2) / 30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		if i == 0 {
			s = X[0]
		} else {
			s = s*alpha + X[i]*(1-alpha)
		}
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointEMA ", r, " ", s)
		}
	}
}

func Test_TimeEMA(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeEMA(3000)
	var s float64 = 0
	var alpha float64 = -math.Log(2) / 3000
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		alpha_ := math.Exp(alpha * 100)
		if i == 0 {
			s = X[i]
		} else {
			s = alpha_*s + (1-alpha_)*X[i]
		}
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeEMA ", r, " ", s)
		}
	}
}

func Test_PointMax(t *testing.T) {
	var cache []float64
	ans := NewPointMax(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetMax(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointMax ", r, " ", s)
		}
	}
}

func Test_TimeMax(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeMax(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetMax(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeMax ", r, " ", s)
		}
	}
}

func Test_PointMean(t *testing.T) {
	var cache []float64
	ans := NewPointMean(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetMean(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointMean ", r, " ", s)
		}
	}
}

func Test_TimeMean(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeMean(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetMean(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeMean ", r, " ", s)
		}
	}
}

func Test_PointMin(t *testing.T) {
	var cache []float64
	ans := NewPointMin(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetMin(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointMax ", r, " ", s)
		}
	}
}

func Test_TimeMin(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeMin(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetMin(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeMax ", r, " ", s)
		}
	}
}

func Test_PointStd(t *testing.T) {
	var cache []float64
	ans := NewPointStd(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetStd(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointStd ", r, " ", s)
		}
	}
}

func Test_TimeStd(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeStd(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetStd(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeStd ", r, " ", s)
		}
	}
}

func Test_PointSum(t *testing.T) {
	var cache []float64
	ans := NewPointSum(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetSum(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointSum ", r, " ", s)
		}
	}
}

func Test_TimeSum(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeSum(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetSum(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeSum ", r, " ", s)
		}
	}
}

func Test_PointVariance(t *testing.T) {
	var cache []float64
	ans := NewPointVariance(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetVariance(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointVariance ", r, " ", s)
		}
	}
}

func Test_TimeVariance(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeVariance(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := mathpro.GetVariance(cache)
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeVariance ", r, " ", s)
		}
	}
}

func Test_PointZScore(t *testing.T) {
	var cache []float64
	ans := NewPointZScore(30)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		if len(cache) == 31 {
			cache = cache[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := (X[i] - mathpro.GetMean(cache)) / mathpro.GetStd(cache)
		if !mathpro.Isfinite(s) {
			s = 0
		}
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator PointZScore ", r, " ", s)
		}
	}
}

func Test_TimeZScore(t *testing.T) {
	var cache []float64
	var tm []int64
	ans := NewTimeZScore(3000)
	for i := 0; i < N; i++ {
		cache = append(cache, X[i])
		tm = append(tm, Localtime[i])
		for len(tm) > 0 && tm[len(tm)-1]-tm[0] > 3000 {
			cache = cache[1:]
			tm = tm[1:]
		}
		r := ans.Update(int64(i*100), X[i])
		s := (X[i] - mathpro.GetMean(cache)) / mathpro.GetStd(cache)
		if !mathpro.Isfinite(s) {
			s = 0
		}
		if !mathpro.AlmostEqual(r, s) {
			t.Log("Wrong operator TimeZScore ", r, " ", s)
		}
	}
}

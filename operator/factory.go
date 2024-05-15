/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:44
 * @Last Modified by: xwu
 * @Last Modified time: 2022-06-01 13:18:05
 */
package operator

// Corr
func NewPointCorr(interval int64) PointCorr {
	op := PointCorr{}
	op.Init(interval)
	return op
}

func NewTimeCorr(interval int64) TimeCorr {
	op := TimeCorr{}
	op.Init(interval)
	return op
}

// Cov
func NewPointCov(interval int64) PointCov {
	op := PointCov{}
	op.Init(interval)
	return op
}

func NewTimeCov(interval int64) TimeCov {
	op := TimeCov{}
	op.Init(interval)
	return op
}

// Decay
func NewPointDecay(interval int64) PointDecay {
	op := PointDecay{}
	op.Init(interval)
	return op
}

func NewTimeDecay(interval int64) TimeDecay {
	op := TimeDecay{}
	op.Init(interval)
	return op
}

// Diff
func NewPointDiff(interval int64) PointDiff {
	op := PointDiff{}
	op.Init(interval)
	return op
}

func NewTimeDiff(interval int64) TimeDiff {
	op := TimeDiff{}
	op.Init(interval)
	return op
}

// DiffMean
func NewPointDiffMean(interval int64) PointDiffMean {
	op := PointDiffMean{}
	op.Init(interval)
	return op
}

func NewTimeDiffMean(interval int64) TimeDiffMean {
	op := TimeDiffMean{}
	op.Init(interval)
	return op
}

// DiffMeanPlus
func NewPointDiffMeanPlus(interval int64) PointDiffMeanPlus {
	op := PointDiffMeanPlus{}
	op.Init(interval)
	return op
}

func NewTimeDiffMeanPlus(interval int64) TimeDiffMeanPlus {
	op := TimeDiffMeanPlus{}
	op.Init(interval)
	return op
}

// Diff
func NewPointDiv(interval int64) PointDiv {
	op := PointDiv{}
	op.Init(interval)
	return op
}

func NewTimeDiv(interval int64) TimeDiv {
	op := TimeDiv{}
	op.Init(interval)
	return op
}

// DiffMean
func NewPointDivMean(interval int64) PointDivMean {
	op := PointDivMean{}
	op.Init(interval)
	return op
}

func NewTimeDivMean(interval int64) TimeDivMean {
	op := TimeDivMean{}
	op.Init(interval)
	return op
}

// EMA
func NewPointEMA(Halflife_Decay int64) PointEMA {
	op := PointEMA{}
	op.Init(Halflife_Decay)
	return op
}

func NewTimeEMA(Halflife_Decay int64) TimeEMA {
	op := TimeEMA{}
	op.Init(Halflife_Decay)
	return op
}

// Max
func NewPointMax(interval int64) PointMax {
	op := PointMax{}
	op.Init(interval)
	return op
}

func NewTimeMax(interval int64) TimeMax {
	op := TimeMax{}
	op.Init(interval)
	return op
}

// Mean
func NewPointMean(interval int64) PointMean {
	op := PointMean{}
	op.Init(interval)
	return op
}

func NewTimeMean(interval int64) TimeMean {
	op := TimeMean{}
	op.Init(interval)
	return op
}

// Min
func NewPointMin(interval int64) PointMin {
	op := PointMin{}
	op.Init(interval)
	return op
}

func NewTimeMin(interval int64) TimeMin {
	op := TimeMin{}
	op.Init(interval)
	return op
}

// Std
func NewPointStd(interval int64) PointStd {
	op := PointStd{}
	op.Init(interval)
	return op
}

func NewTimeStd(interval int64) TimeStd {
	op := TimeStd{}
	op.Init(interval)
	return op
}

// Sum
func NewPointSum(interval int64) PointSum {
	op := PointSum{}
	op.Init(interval)
	return op
}

func NewTimeSum(interval int64) TimeSum {
	op := TimeSum{}
	op.Init(interval)
	return op
}

// Variance
func NewPointVariance(interval int64) PointVariance {
	op := PointVariance{}
	op.Init(interval)
	return op
}

func NewTimeVariance(interval int64) TimeVariance {
	op := TimeVariance{}
	op.Init(interval)
	return op
}

// ZScore
func NewPointZScore(interval int64) PointZScore {
	op := PointZScore{}
	op.Init(interval)
	return op
}

func NewTimeZScore(interval int64) TimeZScore {
	op := TimeZScore{}
	op.Init(interval)
	return op
}

// func NewOpTimer(interval int64) OpTimer {
// 	op := OpTimer{}
// 	op.Init(interval)
// 	return op
// }

func NewTimeOLS(interval int64) TimeOLS {
	op := TimeOLS{}
	op.Init(interval)
	return op
}

// Ret
func NewTimeRet(interval int64) TimeRet {
	op := TimeRet{}
	op.Init(interval)
	return op
}

func NewTimeSkew(interval int64) TimeSkew {
	op := TimeSkew{}
	op.Init(interval)
	return op
}

func NewPointSkesw(interval int64) PointSkew {
	op := PointSkew{}
	op.Init(interval)
	return op
}

package clientTrader

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"winter/global"
	"winter/mathpro"
	"winter/messages"
	"winter/utils"
)

// *********************************************************************************
// ************************************ PnlCalc ************************************
// *********************************************************************************
type PnlCalc struct {
	fee           float64
	money         float64
	maxOccupyFund float64

	// position  container.SkipList
	position1 map[string](SimPosition)
}

func (manager *PnlCalc) BUY(symbol string, size float64, par float64, avgprice float64, fee float64) {
	new_value := SimPosition{Symbol: symbol, Pos: float64(size) * par, LastPx: avgprice}

	if e, ok := manager.position1[symbol]; ok {
		new_value.Pos += e.Pos
		manager.position1[symbol] = new_value
	} else {
		manager.position1[symbol] = new_value
	}
	manager.money -= float64(size) * par * avgprice * (1 + fee)
}

func (manager *PnlCalc) SELL(symbol string, size float64, par float64, avgprice float64, fee float64) {
	new_value := SimPosition{Symbol: symbol, Pos: -float64(size) * par, LastPx: avgprice}

	if e, ok := manager.position1[symbol]; ok {
		new_value.Pos += e.Pos
		manager.position1[symbol] = new_value
	} else {
		manager.position1[symbol] = new_value
	}
	manager.money += float64(size) * par * avgprice * (1 - fee)
}

func (manager *PnlCalc) String() string {
	output := ""
	var pnl float64 = manager.money
	manager.maxOccupyFund = 0

	for _, s := range manager.position1 {
		output = output + s.Symbol + ":" + strconv.FormatFloat(s.Pos, 'f', 8, 64) + " "
		pp := float64(s.Pos) * s.LastPx
		pnl += pp - math.Abs(pp)*manager.fee
		manager.maxOccupyFund += math.Abs(float64(s.Pos)) * s.LastPx
	}
	output = output + "\npnl:" + strconv.FormatFloat(pnl, 'f', 6, 64) + " maxOccupyFund:" + strconv.FormatFloat(manager.maxOccupyFund, 'f', 6, 64)
	return output
}

func NewPnlCalc() PnlCalc {
	return PnlCalc{fee: 0.001, money: 0, maxOccupyFund: 0, position1: make(map[string]SimPosition)}
}

type SimPosition struct {
	Symbol string
	Pos    float64
	LastPx float64
}

// **************************************************************************************
// ************************************ FactorsCache ************************************
// **************************************************************************************
// func NewFactorsCache(exid string, symbol string) FactorsCache {
// 	fileHandle, err := os.OpenFile(
// 		"../Dump/"+exid+"-"+symbol+".csv",
// 		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
// 		0777,
// 	)
// 	if err != nil {
// 		logger.Info("dump factors open failed.")
// 	}

// 	ans := FactorsCache{}
// 	ans.Symbol = symbol
// 	ans.ExID = exid
// 	ans.CurrentDay = 0
// 	ans.IsAlive = false
// 	ans.FileHandle = fileHandle
// 	ans.Buf = bufio.NewWriter(fileHandle)
// 	ans.IC = make(map[string]([]float64))
// 	ans.WinRate = make(map[string]([]float64))
// 	ans.FirstDump = true
// 	ans.Cache = NewTimeCache(global.ReturnInterval)

// 	return ans
// }

// type FactorsCache struct {
// 	Symbol string
// 	ExID   string

// 	IsAlive bool

// 	CurrentDay int

// 	FirstDump bool

// 	FileHandle *os.File
// 	Buf        *bufio.Writer

// 	Dates   []int
// 	IC      map[string]([]float64)
// 	WinRate map[string]([]float64)

// 	Cache TimeCache

// 	FactorColumns []string
// 	Tm            []int64
// 	Factors       [][]float64
// }

// func (FC *FactorsCache) Analysis() bool {
// 	if len(FC.Factors) == 0 {
// 		return false
// 	}

// 	// fmt.Println(len(FC.Factors))

// 	var length, num int = len(FC.Factors), len(FC.FactorColumns) + 3
// 	var ret []float64 = make([]float64, length)
// 	for i := 0; i < length; i += 1 {
// 		ret[i] = FC.Factors[i][num-1]
// 	}

// 	for i := 0; i < len(FC.FactorColumns); i++ {
// 		var name string = FC.FactorColumns[i]
// 		var factors []float64 = make([]float64, length)
// 		for j := range FC.Factors {
// 			factors[j] = FC.Factors[j][i]
// 		}
// 		FC.IC[name] = append(FC.IC[name], mathpro.Corr(factors, ret))
// 		FC.WinRate[name] = append(FC.WinRate[name], mathpro.WinRate(factors, ret))
// 	}

// 	return true
// }

// func (FC *FactorsCache) Clear() {
// 	// 清空存储
// 	FC.Tm = nil
// 	FC.Factors = nil
// }

// func (FC *FactorsCache) OutIC(verbose int) map[string]([]float64) {
// 	table := uitable.New()
// 	table.MaxColWidth = 50

// 	table.AddRow("ExID", "Symbol", "Name", "IC", "IR", "Win%")
// 	for _, name := range FC.FactorColumns {
// 		table.AddRow(
// 			FC.ExID,
// 			FC.Symbol,
// 			name,
// 			strconv.FormatFloat(mathpro.GetMean(FC.IC[name]), 'f', 4, 64),
// 			strconv.FormatFloat(mathpro.GetMean(FC.IC[name])/mathpro.GetStd(FC.IC[name]), 'f', 2, 64),
// 			strconv.FormatFloat(100*mathpro.GetMean(FC.WinRate[name]), 'f', 1, 64),
// 		)
// 	}
// 	if verbose != 0 {
// 		fmt.Println(table)
// 		fmt.Println("*******************************************")
// 	}
// 	return FC.IC
// }

// func (FC *FactorsCache) End() {
// 	if global.Dumper {
// 		FC.Dumper()
// 	}
// 	FC.Analysis()
// 	FC.FileHandle.Close()
// }

// func (FC *FactorsCache) Update(signal *messages.AggSignal) {
// 	FC.IsAlive = true
// 	FC.FactorColumns = signal.Columns
// 	prefactors := FC.Cache.Update(
// 		signal.Localtime,
// 		signal.Signals,
// 		signal.Ask1Price,
// 		signal.Bid1Price,
// 		signal.MinAsk1Price,
// 		signal.MaxBid1Price,
// 		signal.FinalSignal,
// 	)

// 	if prefactors != nil {
// 		var datetime time.Time
// 		for i := range prefactors {
// 			FC.Tm = append(FC.Tm, prefactors[i].Time)
// 			FC.Factors = append(FC.Factors, prefactors[i].Values)

// 			// 计算Cache中出来的最后一个信号的tm是不是应该换天了
// 			datetime = utils.UnixToTime(prefactors[i].Time)
// 		}

// 		if FC.CurrentDay == 0 {
// 			FC.CurrentDay = datetime.Day()
// 		}

// 		if datetime.Day() != FC.CurrentDay {
// 			// 查看是否需要dump相关数据
// 			if global.Dumper {
// 				FC.Dumper()
// 			}
// 			// 分析因子的IC情况
// 			FC.Analysis()

// 			// 清空数组
// 			FC.Clear()

// 			// 换天
// 			FC.Dates = append(FC.Dates, FC.CurrentDay)
// 			FC.CurrentDay = datetime.Day()
// 		}

// 		// logger.Info("HHHHHHHHHHHHHHHHHH")
// 	}

// }

// func (FC *FactorsCache) Dumper() {
// 	for i := range FC.Tm {
// 		if FC.FirstDump {
// 			var columns = []string{"tm", ","}
// 			for j := range FC.FactorColumns {
// 				columns = append(columns, FC.FactorColumns[j], ",")
// 			}
// 			columns = append(columns, "ask1price,bid1price,ret,minask,maxbid,signal\n")
// 			FC.Buf.WriteString(utils.Strcat(columns...))
// 			FC.FirstDump = false
// 		}

// 		var values = []string{strconv.FormatInt(FC.Tm[i], 10), ","}
// 		for j := 0; j < len(FC.FactorColumns)+6; j++ {
// 			values = append(values, strconv.FormatFloat(FC.Factors[i][j], 'f', 6, 64), ",")
// 		}
// 		values = values[:len(values)-1]
// 		FC.Buf.WriteString(utils.Strcat(values...) + "\n")
// 		err := FC.Buf.Flush()
// 		if err != nil {
// 			logger.Info("flush error")
// 		}
// 	}
// }

// // 默认做一分钟的暂时
// func NewTimeCache(interval int64) TimeCache {
// 	ans := TimeCache{}
// 	ans.Init(interval)
// 	return ans
// }

// type timeCacheValue struct {
// 	Time   int64
// 	Values []float64
// }

// type TimeCache struct {
// 	interval int64
// 	Time     []int64
// 	Factors  [][]float64
// 	Ask1     []float64
// 	Bid1     []float64
// 	Ret1m    []float64

// 	FinalSignal []float64

// 	peak   operator.TimeMax
// 	trough operator.TimeMin
// 	std    operator.TimeStd
// }

// func (Cache *TimeCache) Init(interval int64) {
// 	Cache.interval = interval
// 	Cache.peak = operator.NewTimeMax(interval)
// 	Cache.trough = operator.NewTimeMin(interval)
// 	Cache.std = operator.NewTimeStd(interval)
// }

// func (Cache *TimeCache) Update(tm int64, factors []float64, ask1 float64, bid1 float64, minask1 float64, maxbid1 float64, finalSignal float64) []timeCacheValue {
// 	var ans []timeCacheValue = nil

// 	// 先把新的元素加进来
// 	Cache.Time = append(Cache.Time, tm)
// 	Cache.Factors = append(Cache.Factors, factors)
// 	Cache.Ask1 = append(Cache.Ask1, ask1)
// 	Cache.Bid1 = append(Cache.Bid1, bid1)
// 	Cache.FinalSignal = append(Cache.FinalSignal, finalSignal)

// 	// 然后把旧的元素剔除出去并且返回
// 	var i int = 0
// 	for {
// 		if tm-Cache.Time[i] > Cache.interval {
// 			ret := ((ask1+bid1)/(Cache.Ask1[i]+Cache.Bid1[i]) - 1) * 10000
// 			if !mathpro.Isfinite(ret) {
// 				ret = 0
// 			}

// 			var ret1 float64 = minask1
// 			var ret2 float64 = maxbid1

// 			_ans := timeCacheValue{}
// 			_ans.Time = Cache.Time[i]
// 			_ans.Values = Cache.Factors[i]
// 			_ans.Values = append(_ans.Values, Cache.Ask1[i], Cache.Bid1[i], ret, ret1, ret2, Cache.FinalSignal[i])

// 			ans = append(ans, _ans)
// 			i += 1
// 		} else {
// 			Cache.Time = Cache.Time[i:]
// 			Cache.Factors = Cache.Factors[i:]
// 			Cache.Ask1 = Cache.Ask1[i:]
// 			Cache.Bid1 = Cache.Bid1[i:]
// 			Cache.FinalSignal = Cache.FinalSignal[i:]
// 			break
// 		}
// 	}

// 	return ans
// }

// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************
func NewTimeRingBufferGeneric[T []float64 | float64 | item](capacity int, interval int64) TimeRingBufferGeneric[T] {
	if capacity&(capacity-1) != 0 {
		logger.Fatal("circular buffer capacity must be power of 2 ")
	}

	return TimeRingBufferGeneric[T]{
		head_:     0,
		tail_:     0,
		size_:     0,
		interval_: interval,
		capacity_: capacity,
		buffer_:   make([]T, capacity),
		time_:     make([]int64, capacity),
	}
}

type TimeRingBufferGeneric[T []float64 | float64 | item] struct {
	head_, tail_, size_, capacity_ int

	interval_ int64
	buffer_   []T
	time_     []int64
}

func (rb *TimeRingBufferGeneric[T]) IsFull() bool {
	return rb.size_ == rb.capacity_
}

func (rb *TimeRingBufferGeneric[T]) IsEmpty() bool {
	return rb.size_ == 0
}

func (rb *TimeRingBufferGeneric[T]) Size() int {
	return rb.size_
}

func (rb *TimeRingBufferGeneric[T]) Front() T {
	return rb.buffer_[rb.head_]
}

func (rb *TimeRingBufferGeneric[T]) Back() T {
	return rb.buffer_[(rb.tail_-1)&(rb.capacity_-1)]
}

func (rb *TimeRingBufferGeneric[T]) inc_() {
	rb.head_ = (rb.head_ + 1) & (rb.capacity_ - 1)
	rb.size_ -= 1
}

func (rb *TimeRingBufferGeneric[T]) inc() {
	rb.tail_ = (rb.tail_ + 1) & (rb.capacity_ - 1)
	// fmt.Println(rb.tail_)
	rb.size_ += 1
}

func (rb *TimeRingBufferGeneric[T]) PushBack(tm int64, item T) {
	if rb.IsFull() {
		// double the capacity
		rb.buffer_ = append(rb.buffer_, rb.buffer_...)
		rb.time_ = append(rb.time_, rb.time_...)
		if rb.tail_ < rb.head_ {
			rb.head_ += rb.capacity_
		}
		rb.capacity_ = 2 * rb.capacity_
		fmt.Println("double capacity.")
	}

	// fmt.Println(rb.head_, rb.tail_, rb.size_)
	rb.time_[rb.tail_] = tm
	rb.buffer_[rb.tail_] = item
	// fmt.Println(rb.time_)

	rb.inc()
}

func (rb *TimeRingBufferGeneric[T]) PopFront() ([]int64, []T) {
	var ans []T
	var tm []int64
	// fmt.Println(rb.time_[(rb.tail_-1)&rb.capacity_]-rb.time_[rb.head_], rb.interval_)
	for {
		if rb.time_[(rb.tail_-1)&(rb.capacity_-1)]-rb.time_[rb.head_] > rb.interval_ {
			ans = append(ans, rb.buffer_[rb.head_])
			tm = append(tm, rb.time_[rb.head_])
			rb.inc_()

		} else {
			break
		}
	}
	return tm, ans
}

func (rb *TimeRingBufferGeneric[T]) Clear() {
	rb.head_ = 0
	rb.tail_ = 0
	rb.size_ = 0
}

// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************

func NewCalcBase() CalcBase {
	return CalcBase{min: math.MaxFloat64, max: -math.MaxFloat64, mean: 0, welford: 0, n: 0, num_nan: 0}
}

type CalcBase struct {
	min, max, mean, welford float64 // Welford算法
	n, num_nan              int
}

func (calc *CalcBase) Update(x float64) {
	if mathpro.Isfinite(x) {
		calc.n += 1

		if x < calc.min {
			calc.min = x
		}

		if x > calc.max {
			calc.max = x
		}

		premean := calc.mean
		calc.mean = premean + (x-premean)/float64(calc.n)
		calc.welford += (x - premean) * (x - calc.mean)
	} else {
		calc.num_nan += 1
	}
}

func (calc *CalcBase) Value() []float64 {
	std := math.Sqrt(calc.welford / float64(calc.n-1))
	percent_nan := float64(calc.num_nan) / float64(calc.n+calc.num_nan) * 10000
	return []float64{calc.min, calc.max, calc.mean, std, percent_nan}
}

func (calc *CalcBase) Clear() {
	calc.n = 0
	calc.num_nan = 0
	calc.min = 0
	calc.max = 0
	calc.mean = 0
	calc.welford = 0
}

// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************

// 一个币种一个因子
type CalcCorr struct {
	X1, Y1, X2, Y2, XY, Corr float64
	N                        int64
}

func (calc *CalcCorr) Update(tm int64, signal float64, ask1 float64, bid1 float64, ret float64) {
	if mathpro.Isfinite(signal) && mathpro.Isfinite(ret) {
		calc.X1 += signal
		calc.X2 += signal * signal
		calc.Y1 += ret
		calc.Y2 += ret * ret
		calc.XY += signal * ret
		calc.N += 1
	}
}

func (calc *CalcCorr) Value() float64 {
	COV := (calc.XY - calc.X1*calc.Y1/float64(calc.N)) / (float64(calc.N) - 1)
	DX := (calc.X2 - calc.X1*calc.X1/float64(calc.N)) / (float64(calc.N) - 1)
	DY := (calc.Y2 - calc.Y1*calc.Y1/float64(calc.N)) / (float64(calc.N) - 1)

	calc.Corr = COV / math.Sqrt(DX*DY)
	return calc.Corr
}

func (calc *CalcCorr) Clear() {
	calc.X1 = 0
	calc.X2 = 0
	calc.Y1 = 0
	calc.Y2 = 0
	calc.XY = 0
	calc.N = 0
}

// 单个币种，多个因子
func NewAnalysisManager(num_factors int) AnalysisManager {
	ans := AnalysisManager{
		Corr:    make([]CalcCorr, num_factors),
		Base:    make([]CalcBase, num_factors),
		HisCorr: make([][]float64, num_factors), // 单个币种，多个因子，多个时间段
	}

	// X1, Y1, X2, Y2, XY, Corr N
	for i := range ans.Corr {
		ans.Corr[i] = CalcCorr{X1: 0, X2: 0, Y1: 0, Y2: 0, XY: 0, Corr: 0, N: 0}
	}

	for i := range ans.Base {
		ans.Base[i] = NewCalcBase()
	}

	return ans
}

type AnalysisManager struct {
	Corr    []CalcCorr
	Base    []CalcBase // 基础的统计指标
	HisCorr [][]float64
}

func (manager *AnalysisManager) Update(tm int64, signals []float64, ask1 float64, bid1 float64, ret float64) {
	for i := range signals {
		manager.Corr[i].Update(tm, signals[i], ask1, bid1, ret)
		manager.Base[i].Update(signals[i])
	}
}

func (manager *AnalysisManager) Clear() {
	for i := range manager.Corr {
		// fmt.Println(manager.Corr[i])
		manager.HisCorr[i] = append(manager.HisCorr[i], manager.Corr[i].Value())
		manager.Corr[i].Clear()
	}
}

// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************

func NewDumpManager(symbol string) DumpManager {
	fileHandle, err := os.OpenFile(
		"../Dump/"+symbol+".csv",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0777,
	)
	if err != nil {
		logger.Info("dump factors open failed.")
	}

	return DumpManager{
		FileHandle: fileHandle,
		Buf:        bufio.NewWriter(fileHandle),
		Columns:    "",
	}
}

type DumpManager struct {
	FileHandle *os.File
	Buf        *bufio.Writer
	Columns    string
}

func (manager DumpManager) Write(wait2write string) error {
	manager.Buf.WriteString(wait2write)
	err := manager.Buf.Flush()
	return err
}

func (manager DumpManager) Writeline(wait2write string) error {
	manager.Buf.WriteString(wait2write + "/n")
	err := manager.Buf.Flush()
	return err
}

func (manager DumpManager) Convert(
	tm int64,
	ret float64,
	exchange string,
	symbol string,
	ask1 float64,
	bid1 float64,
	pred float64,
	factors []float64,
) string {
	/*
		exchange,symbol,tm,y1,ask1,bid1,yhat,QTE000...
	*/
	var value1 = []string{
		exchange,
		",",
		symbol,
		",",
		strconv.FormatInt(tm, 10),
		",",
		strconv.FormatFloat(ret, 'f', 6, 64),
		",",
		strconv.FormatFloat(ask1, 'f', 6, 64),
		",",
		strconv.FormatFloat(bid1, 'f', 6, 64),
		",",
		strconv.FormatFloat(pred, 'f', 6, 64),
	}

	var value2 []string = make([]string, 2*len(factors)+1)
	for i := range factors {
		value2[2*i] = ","
		value2[2*i+1] = strconv.FormatFloat(factors[i], 'f', 6, 64)
	}
	value2[2*len(factors)] = "\n"

	value := append(value1, value2...)

	return utils.Strcat(value...)
}

// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************
// ********************************************************************

type item struct {
	Exchange, Symbol           string
	Ask1Price, Bid1Price, Yhat float64
	Signals                    []float64
}

func NewHandleToSignal(interval int64, num_symbols int, num_factors int, isdumper bool) handleToSignal {
	ans := handleToSignal{num_factors: num_factors, num_symbols: num_symbols, Isdumper: isdumper}

	ans.Cache = make([]TimeRingBufferGeneric[item], num_symbols)
	for i := range ans.Cache {
		ans.Cache[i] = NewTimeRingBufferGeneric[item](4096, interval)
	}

	ans.Analysis = make([]AnalysisManager, num_symbols)
	for i := range ans.Analysis {
		ans.Analysis[i] = NewAnalysisManager(num_factors)
	}

	ans.AllAnalysis = NewAnalysisManager(num_factors)

	ans.CurrentDay = make([]int, num_factors)
	for i := range ans.CurrentDay {
		ans.CurrentDay[i] = 0
	}
	ans.AllCurrentDay = 0

	if ans.Isdumper {
		ans.Dumper = make([]DumpManager, num_symbols)
		for symbol, ii := range global.AlphaUid[1] {
			ans.Dumper[ii] = NewDumpManager(symbol)
		}

		ans.AllDumper = NewDumpManager("all")
	}

	return ans
}

type handleToSignal struct {
	num_factors, num_symbols int

	Cache []TimeRingBufferGeneric[item]

	Analysis   []AnalysisManager
	Dumper     []DumpManager
	CurrentDay []int

	AllAnalysis   AnalysisManager
	AllDumper     DumpManager
	AllCurrentDay int

	Isdumper bool
}

func (p *handleToSignal) Update(signal *messages.AggSignal) {
	// 第一步骤是计算出来收益率，然后得到收益率对应的十秒前的数据
	// 然后把delete的数据传到一个是分析的类，一个是dump的类
	ii := global.Instruments[signal.Symbol].II
	cache := &p.Cache[ii]
	analysis := &p.Analysis[ii]

	if p.Isdumper {
		// 写入csv文件的columns
		if strings.Compare(p.Dumper[ii].Columns, "") == 0 {
			var cols = []string{"exchange,symbol,tm,y1,ask1,bid1,yhat"}
			for i := range signal.Columns {
				cols = append(cols, ",", signal.Columns[i])
			}
			cols = append(cols, "\n")
			p.Dumper[ii].Columns = utils.Strcat(cols...)
			p.Dumper[ii].Write(p.Dumper[ii].Columns)
		}
		if strings.Compare(p.AllDumper.Columns, "") == 0 {
			var cols = []string{"exchange,symbol,tm,y1,ask1,bid1,yhat"}
			for i := range signal.Columns {
				cols = append(cols, ",", signal.Columns[i])
			}
			cols = append(cols, "\n")
			p.AllDumper.Columns = utils.Strcat(cols...)
			p.AllDumper.Write(p.Dumper[ii].Columns)
		}
	}

	cache.PushBack(signal.Localtime, item{
		Exchange:  global.ExIDList[signal.ExID],
		Symbol:    signal.Symbol,
		Ask1Price: signal.Ask1Price,
		Bid1Price: signal.Bid1Price,
		Signals:   signal.Signals,
		Yhat:      signal.FinalSignal,
	})

	t, items := cache.PopFront()

	for i := range items {
		ret := (2*signal.Mid/(items[i].Ask1Price+items[i].Bid1Price) - 1) * 10000

		// *********************** Analysis ***********************
		// 换天，分析的时候计算每一天的ic，最终的ic是每天ic的均值
		datetime := utils.UnixToTime(t[i])
		if p.CurrentDay[ii] != 0 {
			if datetime.Day() > p.CurrentDay[ii] {
				// 如果不是第一次，而且天数不一样，那么ic清空，重新计算
				analysis.Clear()
				p.CurrentDay[ii] = datetime.Day()
			}
		} else {
			// 初始化
			p.CurrentDay[ii] = datetime.Day()
		}
		if p.AllCurrentDay != 0 {
			if datetime.Day() > p.AllCurrentDay {
				p.AllAnalysis.Clear()
				p.AllCurrentDay = datetime.Day()
			}
		} else {
			p.AllCurrentDay = datetime.Day()
		}

		// 分析ICIR等指标
		analysis.Update(t[i], items[i].Signals, items[i].Ask1Price, items[i].Bid1Price, ret)
		p.AllAnalysis.Update(t[i], items[i].Signals, items[i].Ask1Price, items[i].Bid1Price, ret)

		if p.Isdumper {
			// *********************** Dumper ***********************
			// dump因子  tm, ret, ask1, bid1, yhat, factors
			wait2write := p.Dumper[ii].Convert(
				t[i],
				ret,
				items[i].Exchange,
				items[i].Symbol,
				items[i].Ask1Price,
				items[i].Bid1Price,
				items[i].Yhat,
				items[i].Signals,
			)
			err := p.Dumper[ii].Write(wait2write)
			if err != nil {
				logger.Info("Wrileline failed: " + wait2write)
			}
			err = p.AllDumper.Write(wait2write)
			if err != nil {
				logger.Info("Wrileline failed: " + wait2write)
			}
		}
	}

}

func (p *handleToSignal) End() [][]float64 {
	// 第一步骤是计算出来收益率，然后得到收益率对应的十秒前的数据
	// 然后把delete的数据传到一个是分析的类，一个是dump的类
	for i := range p.Dumper {
		p.Dumper[i].FileHandle.Close()
	}

	for i := range p.Analysis {
		p.Analysis[i].Clear()
	}
	p.AllAnalysis.Clear()

	// 整理输出结果
	var cache [][]float64 = make([][]float64, p.num_factors)
	for i := range cache {
		cache[i] = make([]float64, p.num_symbols)
	}

	for i := range p.Analysis {
		for j := range p.Analysis[i].HisCorr {
			cache[j][i] = mathpro.GetMean(p.Analysis[i].HisCorr[j])
		}
	}

	var ic_total []float64 = make([]float64, p.num_factors)
	var ir_total []float64 = make([]float64, p.num_factors)

	for j := range p.AllAnalysis.HisCorr {
		ic_total[j] = mathpro.GetMean(p.AllAnalysis.HisCorr[j])
		// fmt.Println("HHHHH", p.AllAnalysis.HisCorr[j])
		// ir_total[j] = ic_total[j] / mathpro.GetStd(p.AllAnalysis.HisCorr[j])
		ir_total[j] = ic_total[j] / 1
	}

	var base [][]float64 = make([][]float64, 5)
	for i := range p.AllAnalysis.Base {
		for j, v := range p.AllAnalysis.Base[i].Value() {
			base[j] = append(base[j], v)
		}
	}

	// 写入需要返回的结果
	var results [][]float64 = make([][]float64, 9)

	{
		// ic_td
		var result []float64 = make([]float64, p.num_factors)
		for i := range cache {
			result[i] = mathpro.GetMean(cache[i])
		}
		results[0] = result
	}

	{
		// ir_td
		var result []float64 = make([]float64, p.num_factors)
		for i := range cache {
			result[i] = mathpro.GetMean(cache[i]) / mathpro.GetStd(cache[i])
		}
		results[1] = result
	}

	results[2] = ic_total // ic_total
	results[3] = ir_total // ir_total
	for i := 4; i < 9; i += 1 {
		results[i] = base[i-4]
	}

	return results
}

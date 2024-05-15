/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:25
 * @Last Modified by: xwu
 * @Last Modified time: 2022-06-07 10:01:52
 */
package messages

//
type Instrument struct {
	ExID   int
	Symbol string // 这个是什么品种
	II     int
}

// 这个结构体是从data到process的标准结构体
type MsgDataToProcess struct {
	ExID     int
	Contents string
}

// 这个结构体是从process到alpha使用的标准结构体
type AggStMsg struct {
	ExID    int
	Status  int    // 这个是信息状态主要是用于回测的时候，来给出一个结束的信号
	Kinds   int    // 这个是品种的什么数据，0是ob，1是td
	Symbols string // 这个是什么品种

	OrdEvt   OrderbookBase
	TrdEvt   TradeBase
	KlineEvt KlineBase
}

type OrderbookBase struct {
	Localtime int64
	Eventime  int64

	Asks [50][4]float64 // price, size, 强制平仓单数, 总单数  默认50档
	Bids [50][4]float64
}

type TradeBase struct {
	Localtime int64
	Eventime  int64

	Price float64
	Size  float64
	Side  float64

	IsTaker bool
}

type KlineBase struct {
	Frequency int64
	Localtime int64
	StartTm   int64

	Open   float64
	High   float64
	Low    float64
	Close  float64
	Vol    float64
	VolCcy float64
}

// 从alpha到trade的标准结构体
type AggSignal struct {
	ExID   int    // 交易所
	Status int    // 这个是信息状态主要是用于回测的时候，来给出一个结束的信号,0代表结束
	Symbol string // 币种
	Kind   int    // 消息种类 0:Order 1:Trade 2:Kline

	Mid          float64
	Ask1Price    float64
	MinAsk1Price float64
	Bid1Price    float64
	MaxBid1Price float64
	Localtime    int64
	Eventime     int64
	StartTm      int64
	IsOpen       bool

	Columns     []string
	Signals     []float64
	FinalSignal float64
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:05
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-03 19:50:13
 */
package clientProcess

// import (
// 	"strconv"
// 	"strings"
// 	"winter/mathpro"
// 	"winter/messages"
// 	"winter/utils"
// )

// type okexQteBase struct {
// 	Price      float64
// 	Size       float64
// 	Compulsory float64
// 	NumOrders  float64
// }

// type okexQteEvt struct {
// 	PriceDecimal int
// 	SizeDecimal  int
// 	checksum     int32

// 	// AskSkipList *container.SkipList
// 	// BidSkipList *container.SkipList

// 	AskSkipList *SkiplistOkexQteBase
// 	BidSkipList *SkiplistOkexQteBase
// }

// func (qte *okexQteEvt) Reset(ord *messages.OkexOrderResponse, decimal [2]int) {
// 	// qte.AskSkipList = container.NewSkipList()
// 	// qte.BidSkipList = container.NewSkipList()

// 	qte.AskSkipList = NewSkiplistOkexQteBase()
// 	qte.BidSkipList = NewSkiplistOkexQteBase()

// 	qte.PriceDecimal = decimal[0]
// 	qte.SizeDecimal = decimal[1]
// 	qte.checksum = ord.Data[0].CheckSum

// 	for i := range ord.Data[0].Asks {
// 		var price float64 = utils.String2float(ord.Data[0].Asks[i][0])
// 		var size float64 = utils.String2float(ord.Data[0].Asks[i][1])
// 		var c float64 = utils.String2float(ord.Data[0].Asks[i][2])
// 		var num float64 = utils.String2float(ord.Data[0].Asks[i][3])

// 		var key float64 = price

// 		qte.AskSkipList.Add(
// 			key,
// 			&okexQteBase{
// 				Price:      price,
// 				Size:       size,
// 				Compulsory: c,
// 				NumOrders:  num,
// 			},
// 		)
// 	}

// 	for i := range ord.Data[0].Bids {
// 		var price float64 = utils.String2float(ord.Data[0].Bids[i][0])
// 		var size float64 = utils.String2float(ord.Data[0].Bids[i][1])
// 		var c float64 = utils.String2float(ord.Data[0].Bids[i][2])
// 		var num float64 = utils.String2float(ord.Data[0].Bids[i][3])

// 		// var key int32 = -int32(qte.Decimal * price)
// 		var key float64 = -price

// 		qte.BidSkipList.Add(
// 			key,
// 			&okexQteBase{
// 				Price:      price,
// 				Size:       size,
// 				Compulsory: c,
// 				NumOrders:  num,
// 			},
// 		)
// 	}

// }

// func (qte *okexQteEvt) Update(ord *messages.OkexOrderResponse) {
// 	qte.checksum = ord.Data[0].CheckSum

// 	for i := range ord.Data[0].Asks {
// 		var price float64 = utils.String2float(ord.Data[0].Asks[i][0])
// 		var size float64 = utils.String2float(ord.Data[0].Asks[i][1])
// 		var c float64 = utils.String2float(ord.Data[0].Asks[i][2])
// 		var num float64 = utils.String2float(ord.Data[0].Asks[i][3])

// 		// var key int32 = int32(qte.Decimal * price)
// 		var key float64 = price

// 		if mathpro.AlmostEqual(size, 0) {
// 			qte.AskSkipList.Remove(key)
// 		} else {
// 			qte.AskSkipList.Add(
// 				key,
// 				&okexQteBase{
// 					Price:      price,
// 					Size:       size,
// 					Compulsory: c,
// 					NumOrders:  num,
// 				},
// 			)
// 		}
// 	}

// 	for i := range ord.Data[0].Bids {
// 		var price float64 = utils.String2float(ord.Data[0].Bids[i][0])
// 		var size float64 = utils.String2float(ord.Data[0].Bids[i][1])
// 		var c float64 = utils.String2float(ord.Data[0].Bids[i][2])
// 		var num float64 = utils.String2float(ord.Data[0].Bids[i][3])

// 		var key float64 = -price

// 		if mathpro.AlmostEqual(size, 0) {
// 			qte.BidSkipList.Remove(key)
// 		} else {
// 			qte.BidSkipList.Add(
// 				key,
// 				&okexQteBase{
// 					Price:      price,
// 					Size:       size,
// 					Compulsory: c,
// 					NumOrders:  num,
// 				},
// 			)
// 		}
// 	}

// }

// func (qte *okexQteEvt) GetCurrentOrderbook(ask_or_bid string) [][]float64 {
// 	var ans [][]float64 = make([][]float64, 50)
// 	for i := 0; i < 50; i += 1 {
// 		ans[i] = make([]float64, 4)
// 	}

// 	switch ask_or_bid {
// 	case "ask":
// 		asks := qte.AskSkipList.GetLoop()

// 		for i := range asks {
// 			if i == 50 {
// 				break
// 			}
// 			// askValue := asks[i].(*okexQteBase)
// 			askValue := asks[i]

// 			ans[i][0] = askValue.Price
// 			ans[i][1] = askValue.Size
// 			ans[i][2] = askValue.Compulsory
// 			ans[i][3] = askValue.NumOrders

// 		}
// 	case "bid":
// 		bids := qte.BidSkipList.GetLoop()
// 		// 发生了cross

// 		for i := range bids {
// 			if i == 50 {
// 				break
// 			}
// 			// bidValue := bids[i].(*okexQteBase)
// 			bidValue := bids[i]

// 			ans[i][0] = bidValue.Price
// 			ans[i][1] = bidValue.Size
// 			ans[i][2] = bidValue.Compulsory
// 			ans[i][3] = bidValue.NumOrders

// 		}
// 	}

// 	return ans
// }

// func (qte *okexQteEvt) float2string(v float64, prec int) string {
// 	var ans string
// 	if v < 1e-4 {
// 		ans = strconv.FormatFloat(v, 'e', prec, 64)
// 	} else {
// 		ans = strconv.FormatFloat(v, 'f', prec, 64)
// 	}

// 	var count int = len(ans)

// 	if strings.Contains(ans, ".") {
// 		for i := len(ans) - 1; i >= 0; i -= 1 {
// 			if ans[i] == '0' {
// 				count -= 1
// 			} else if ans[i] == '.' {
// 				count -= 1
// 				ans = ans[:count]
// 				break
// 			} else {
// 				ans = ans[:count]
// 				break
// 			}
// 		}
// 	}
// 	return ans
// }

// func (qte *okexQteEvt) CheckSum() bool {
// 	var crcStr []string
// 	asks := qte.GetCurrentOrderbook("ask")[:25]
// 	bids := qte.GetCurrentOrderbook("bid")[:25]

// 	var iterAsks, iterBids int = 0, 0

// 	for {
// 		if iterAsks < len(asks) && iterBids < len(bids) {
// 			crcStr = append(
// 				crcStr,
// 				qte.float2string(bids[iterBids][0], qte.PriceDecimal),
// 				":",
// 				qte.float2string(bids[iterBids][1], qte.SizeDecimal),
// 				":",
// 				qte.float2string(asks[iterAsks][0], qte.PriceDecimal),
// 				":",
// 				qte.float2string(asks[iterAsks][1], qte.SizeDecimal),
// 				":",
// 			)
// 			iterAsks += 1
// 			iterBids += 1
// 		} else if iterAsks == len(asks) && iterBids < len(bids) {
// 			crcStr = append(
// 				crcStr,
// 				qte.float2string(bids[iterBids][0], qte.PriceDecimal),
// 				":",
// 				qte.float2string(bids[iterBids][1], qte.SizeDecimal),
// 				":",
// 			)
// 			iterBids += 1
// 		} else if iterAsks < len(asks) && iterBids == len(bids) {
// 			crcStr = append(
// 				crcStr,
// 				qte.float2string(asks[iterAsks][0], qte.PriceDecimal),
// 				":",
// 				qte.float2string(asks[iterAsks][1], qte.SizeDecimal),
// 				":",
// 			)
// 			iterAsks += 1
// 		} else {
// 			crcStr = crcStr[:len(crcStr)-1]
// 			break
// 		}
// 	}

// 	var ans bool = true
// 	if int32(utils.CRC32Slice(crcStr))-int32(qte.checksum) != 0 {
// 		ans = false
// 		logger.Info("okex checksum error  " + utils.Strcatslice(crcStr))
// 	}
// 	return ans
// }

// type okexTrdEvt struct {
// 	Price float64
// 	Size  float64
// 	Side  float64
// 	Tm    int64
// }

// type okexKlineEvt struct {
// 	Frequency int64
// 	StartTm   int64
// 	Open      float64
// 	High      float64
// 	Low       float64
// 	Close     float64
// 	Vol       float64
// 	VolCcy    float64
// }

import (
	"fmt"
	"strconv"
	"strings"
	"winter/container"
	"winter/mathpro"
	"winter/messages"
	"winter/utils"
)

// type okexQteBase struct {
// 	Price      float64
// 	Size       float64
// 	Compulsory float64
// 	NumOrders  float64
// }

type okexQteAskBase struct {
	Price      float64
	Size       float64
	Compulsory float64
	NumOrders  float64
}

func (e okexQteAskBase) ExtractKey() float64 {
	return e.Price
}

func (e okexQteAskBase) String() string {
	return strconv.FormatFloat(e.Price, 'f', 6, 64)
}

type okexQteBidBase struct {
	Price      float64
	Size       float64
	Compulsory float64
	NumOrders  float64
}

func (e okexQteBidBase) ExtractKey() float64 {
	return -e.Price
}

func (e okexQteBidBase) String() string {
	return strconv.FormatFloat(e.Price, 'f', 6, 64)
}

type okexQteEvt struct {
	PriceDecimal int
	SizeDecimal  int
	checksum     int32

	// AskSkipList *container.SkipList
	// BidSkipList *container.SkipList
	asks [][]float64
	bids [][]float64

	AskSkipList container.SkipList
	BidSkipList container.SkipList
}

func (qte *okexQteEvt) Init() {
	qte.asks = make([][]float64, 50)
	for i := 0; i < 50; i += 1 {
		qte.asks[i] = make([]float64, 4)
	}

	qte.bids = make([][]float64, 50)
	for i := 0; i < 50; i += 1 {
		qte.bids[i] = make([]float64, 4)
	}
}

func (qte *okexQteEvt) ClearAsks() {
	for i := range qte.asks {
		qte.asks[i][0] = -1
		qte.asks[i][1] = -1
		qte.asks[i][2] = -1
		qte.asks[i][3] = -1
	}
}

func (qte *okexQteEvt) ClearBids() {
	for i := range qte.bids {
		qte.bids[i][0] = -1
		qte.bids[i][1] = -1
		qte.bids[i][2] = -1
		qte.bids[i][3] = -1
	}
}

func (qte *okexQteEvt) Reset(ord *messages.OkexOrderResponse, decimal [2]int) {
	// qte.AskSkipList = container.NewSkipList()
	// qte.BidSkipList = container.NewSkipList()

	qte.AskSkipList = container.New()
	qte.BidSkipList = container.New()

	qte.PriceDecimal = decimal[0]
	qte.SizeDecimal = decimal[1]
	qte.checksum = ord.Data[0].CheckSum

	for i := range ord.Data[0].Asks {
		var price float64 = utils.String2float(ord.Data[0].Asks[i][0])
		var size float64 = utils.String2float(ord.Data[0].Asks[i][1])
		var c float64 = utils.String2float(ord.Data[0].Asks[i][2])
		var num float64 = utils.String2float(ord.Data[0].Asks[i][3])

		qte.AskSkipList.Insert(
			okexQteAskBase{
				Price:      price,
				Size:       size,
				Compulsory: c,
				NumOrders:  num,
			},
		)
	}

	for i := range ord.Data[0].Bids {
		var price float64 = utils.String2float(ord.Data[0].Bids[i][0])
		var size float64 = utils.String2float(ord.Data[0].Bids[i][1])
		var c float64 = utils.String2float(ord.Data[0].Bids[i][2])
		var num float64 = utils.String2float(ord.Data[0].Bids[i][3])

		qte.BidSkipList.Insert(
			okexQteBidBase{
				Price:      price,
				Size:       size,
				Compulsory: c,
				NumOrders:  num,
			},
		)
	}

}

func (qte *okexQteEvt) Update(ord *messages.OkexOrderResponse) {
	qte.checksum = ord.Data[0].CheckSum

	for i := range ord.Data[0].Asks {
		var price float64 = utils.String2float(ord.Data[0].Asks[i][0])
		var size float64 = utils.String2float(ord.Data[0].Asks[i][1])
		var c float64 = utils.String2float(ord.Data[0].Asks[i][2])
		var num float64 = utils.String2float(ord.Data[0].Asks[i][3])

		if mathpro.AlmostEqual(size, 0) {
			qte.AskSkipList.Delete(
				okexQteAskBase{
					Price:      price,
					Size:       size,
					Compulsory: c,
					NumOrders:  num,
				},
			)
		} else {
			event := okexQteAskBase{
				Price:      price,
				Size:       size,
				Compulsory: c,
				NumOrders:  num,
			}

			if x, ok := qte.AskSkipList.Find(event); ok {
				qte.AskSkipList.ChangeValue(x, event)
			} else {
				qte.AskSkipList.Insert(event)
			}
		}
	}

	for i := range ord.Data[0].Bids {
		var price float64 = utils.String2float(ord.Data[0].Bids[i][0])
		var size float64 = utils.String2float(ord.Data[0].Bids[i][1])
		var c float64 = utils.String2float(ord.Data[0].Bids[i][2])
		var num float64 = utils.String2float(ord.Data[0].Bids[i][3])

		if mathpro.AlmostEqual(size, 0) {
			qte.BidSkipList.Delete(
				okexQteBidBase{
					Price:      price,
					Size:       size,
					Compulsory: c,
					NumOrders:  num,
				},
			)
		} else {
			event := okexQteBidBase{
				Price:      price,
				Size:       size,
				Compulsory: c,
				NumOrders:  num,
			}
			if x, ok := qte.BidSkipList.Find(event); ok {
				qte.BidSkipList.ChangeValue(x, event)
			} else {
				qte.BidSkipList.Insert(event)
			}
		}
	}

}

func (qte *okexQteEvt) GetCurrentOrderbook(ask_or_bid string) [50][4]float64 {
	var ans [50][4]float64 // 这里注意不能使用slice，因为slice不扩容的话很可能是共用同一片内存

	switch ask_or_bid {
	case "ask":
		// qte.ClearAsks()
		asklevel := qte.AskSkipList.GetNodeCount()

		evt := qte.AskSkipList.GetSmallestNode()
		for i := 0; i < asklevel; i++ {
			if i == 50 {
				break
			}
			value := evt.GetValue().(okexQteAskBase)
			// qte.asks[i][0] = value.Price
			// qte.asks[i][1] = value.Size
			// qte.asks[i][2] = value.Compulsory
			// qte.asks[i][3] = value.NumOrders
			ans[i][0] = value.Price
			ans[i][1] = value.Size
			ans[i][2] = value.Compulsory
			ans[i][3] = value.NumOrders
			evt = qte.AskSkipList.Next(evt)
		}
		return ans // 这里的行为有些类似于指针，所以需要注意值的改动
	case "bid":
		// qte.ClearBids()
		bidlevel := qte.BidSkipList.GetNodeCount()

		evt := qte.BidSkipList.GetSmallestNode()
		for i := 0; i < bidlevel; i++ {
			if i == 50 {
				break
			}
			value := evt.GetValue().(okexQteBidBase)
			// qte.bids[i][0] = value.Price
			// qte.bids[i][1] = value.Size
			// qte.bids[i][2] = value.Compulsory
			// qte.bids[i][3] = value.NumOrders

			ans[i][0] = value.Price
			ans[i][1] = value.Size
			ans[i][2] = value.Compulsory
			ans[i][3] = value.NumOrders
			evt = qte.BidSkipList.Next(evt)
		}
		return ans
	default:
		logger.Fatal("error type in Orderbook")
	}
	return ans
}

func (qte *okexQteEvt) float2string(v float64, prec int) string {
	var ans string
	if v < 1e-4 {
		ans = strconv.FormatFloat(v, 'e', prec, 64)
	} else {
		ans = strconv.FormatFloat(v, 'f', prec, 64)
	}

	var count int = len(ans)
	if strings.Contains(ans, ".") {
		for i := len(ans) - 1; i >= 0; i -= 1 {
			if ans[i] == '0' {
				count -= 1
			} else if ans[i] == '.' {
				count -= 1
				ans = ans[:count]
				break
			} else {
				ans = ans[:count]
				break
			}
		}
	}

	return ans
}

func (qte *okexQteEvt) CheckSum() bool {
	var crcStr []string
	var asks, bids [25][4]float64
	_asks := qte.GetCurrentOrderbook("ask")
	_bids := qte.GetCurrentOrderbook("bid")
	if len(_asks) >= 25 {
		for i := 0; i < 25; i++ {
			asks[i] = _asks[i]
		}
	}
	if len(_bids) >= 25 {
		for i := 0; i < 25; i++ {
			bids[i] = _bids[i]
		}
	}

	var iterAsks, iterBids int = 0, 0
	var asklevel, bidlevel int = 0, 0

	for i := range asks {
		if mathpro.AlmostEqual(asks[i][0], 0) {
			asklevel += 1
		} else {
			break
		}
	}

	for i := range bids {
		if mathpro.AlmostEqual(bids[i][0], 0) {
			bidlevel += 1
		} else {
			break
		}
	}

	if asklevel > 0 || bidlevel > 0 {
		for {
			if iterAsks < asklevel && iterBids < bidlevel {
				crcStr = append(
					crcStr,
					qte.float2string(bids[iterBids][0], qte.PriceDecimal),
					":",
					qte.float2string(bids[iterBids][1], qte.SizeDecimal),
					":",
					qte.float2string(asks[iterAsks][0], qte.PriceDecimal),
					":",
					qte.float2string(asks[iterAsks][1], qte.SizeDecimal),
					":",
				)
				iterAsks += 1
				iterBids += 1
			} else if iterAsks == asklevel && iterBids < bidlevel {
				crcStr = append(
					crcStr,
					qte.float2string(bids[iterBids][0], qte.PriceDecimal),
					":",
					qte.float2string(bids[iterBids][1], qte.SizeDecimal),
					":",
				)
				iterBids += 1
			} else if iterAsks < asklevel && iterBids == bidlevel {
				crcStr = append(
					crcStr,
					qte.float2string(asks[iterAsks][0], qte.PriceDecimal),
					":",
					qte.float2string(asks[iterAsks][1], qte.SizeDecimal),
					":",
				)
				iterAsks += 1
			} else {
				crcStr = crcStr[:len(crcStr)-1]
				break
			}
		}

		var ans bool = true
		if int32(utils.CRC32Slice(crcStr))-int32(qte.checksum) != 0 {
			ans = false
			logger.Info("okex checksum error  " + utils.Strcatslice(crcStr))
			fmt.Println(asks)
			fmt.Println(bids)
			fmt.Println(int32(qte.checksum))
		}
		return ans
	} else {
		return true
	}
}

type okexTrdEvt struct {
	Price float64
	Size  float64
	Side  float64
	Tm    int64
}

type okexKlineEvt struct {
	Frequency int64
	StartTm   int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Vol       float64
	VolCcy    float64
}

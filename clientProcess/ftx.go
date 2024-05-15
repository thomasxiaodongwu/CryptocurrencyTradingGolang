/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:14
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-31 23:36:40
 */
package clientProcess

import (
	"strconv"
	"strings"
	"winter/container"
	"winter/mathpro"
	"winter/messages"
	"winter/utils"
)

// 这里是对增量的orderbook信息进行处理，生成最终的orderbook，一般是上下二十档
type FTXQteBase struct {
	Price float64
	Size  float64
}

type FTXQteAskBase struct {
	Price float64
	Size  float64
}

func (e FTXQteAskBase) ExtractKey() float64 {
	return e.Price
}

func (e FTXQteAskBase) String() string {
	return strconv.FormatFloat(e.Price, 'f', 6, 64)
}

type FTXQteBidBase struct {
	Price float64
	Size  float64
}

func (e FTXQteBidBase) ExtractKey() float64 {
	return -e.Price
}

func (e FTXQteBidBase) String() string {
	return strconv.FormatFloat(e.Price, 'f', 6, 64)
}

type FTXQteEvt struct {
	PriceDecimal int
	SizeDecimal  int
	checksum     uint32

	AskSkipList container.SkipList
	BidSkipList container.SkipList
}

func (qte *FTXQteEvt) Reset(ord *messages.FTXOrderResponse, decimal [2]int) error {
	var err error = nil

	qte.PriceDecimal = decimal[0]
	qte.SizeDecimal = decimal[1]
	qte.checksum = ord.Data.Checksum

	// 重置asks和bids的map
	qte.AskSkipList = container.New()
	qte.BidSkipList = container.New()

	for i := range ord.Data.Asks {
		var price float64 = ord.Data.Asks[i][0]
		var size float64 = ord.Data.Asks[i][1]

		qte.AskSkipList.Insert(
			FTXQteAskBase{
				Price: price,
				Size:  size,
			},
		)
	}

	for i := range ord.Data.Bids {
		var price float64 = ord.Data.Bids[i][0]
		var size float64 = ord.Data.Bids[i][1]

		qte.BidSkipList.Insert(
			FTXQteBidBase{
				Price: price,
				Size:  size,
			},
		)
	}

	return err
}

func (qte *FTXQteEvt) Update(ord *messages.FTXOrderResponse) {
	qte.checksum = ord.Data.Checksum

	for i := range ord.Data.Asks {
		var price float64 = ord.Data.Asks[i][0]
		var size float64 = ord.Data.Asks[i][1]

		event := FTXQteAskBase{Price: price, Size: size}
		if mathpro.AlmostEqual(size, 0) {
			qte.AskSkipList.Delete(event)
		} else {
			if x, ok := qte.AskSkipList.Find(event); ok {
				qte.AskSkipList.ChangeValue(x, event)
			} else {
				qte.AskSkipList.Insert(event)
			}
		}
	}

	for i := range ord.Data.Bids {
		var price float64 = ord.Data.Bids[i][0]
		var size float64 = ord.Data.Bids[i][1]

		event := FTXQteBidBase{Price: price, Size: size}
		if mathpro.AlmostEqual(size, 0) {
			qte.BidSkipList.Delete(event)
		} else {
			if x, ok := qte.BidSkipList.Find(event); ok {
				qte.BidSkipList.ChangeValue(x, event)
			} else {
				qte.BidSkipList.Insert(event)
			}
		}
	}

}

func (qte *FTXQteEvt) GetCurrentOrderbook(ask_or_bid string) [50][4]float64 {
	var ans [50][4]float64

	switch ask_or_bid {
	case "ask":
		evt := qte.AskSkipList.GetSmallestNode()
		count := 0
		for {
			var _ans [4]float64
			_ans[0] = evt.GetValue().(FTXQteAskBase).Price
			_ans[1] = evt.GetValue().(FTXQteAskBase).Size
			ans[count] = _ans
			count += 1
			evt = qte.AskSkipList.Next(evt)
			// logger.Info("HHHHHHHHHHHHHHHHH")
			if evt == qte.AskSkipList.GetSmallestNode() {
				break
			}
		}
	case "bid":
		evt := qte.BidSkipList.GetSmallestNode()
		count := 0
		for {
			var _ans [4]float64
			_ans[0] = evt.GetValue().(FTXQteBidBase).Price
			_ans[1] = evt.GetValue().(FTXQteBidBase).Size
			ans[count] = _ans
			count += 1
			evt = qte.BidSkipList.Next(evt)
			if evt == qte.BidSkipList.GetSmallestNode() {
				break
			}
		}
	}

	return ans
}

func (qte *FTXQteEvt) float2string(v float64, prec int) string {
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
				ans = ans[:count+1]
				break
			} else {
				ans = ans[:count]
				break
			}
		}
	} else {
		ans = utils.Strcat(ans, `.0`)
	}
	return ans
}

func (qte *FTXQteEvt) CheckSum() bool {
	var crcStr []string
	asks := qte.GetCurrentOrderbook("ask")
	bids := qte.GetCurrentOrderbook("bid")

	var iterAsks, iterBids int = 0, 0

	for {
		if iterAsks < len(asks) && iterBids < len(bids) {
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
		} else if iterAsks == len(asks) && iterBids < len(bids) {
			crcStr = append(
				crcStr,
				qte.float2string(bids[iterBids][0], qte.PriceDecimal),
				":",
				qte.float2string(bids[iterBids][1], qte.SizeDecimal),
				":",
			)
			iterBids += 1
		} else if iterAsks < len(asks) && iterBids == len(bids) {
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
		logger.Info("ftx checksum error  " + utils.Strcatslice(crcStr))
	}
	return ans
}

// 这里是FTX的trade部分，这部分转化没什么好说的，只要类型变化就可以了
type FTXTrdBase struct {
	Price float64
	Size  float64
	Side  float64
	Tm    int64
}

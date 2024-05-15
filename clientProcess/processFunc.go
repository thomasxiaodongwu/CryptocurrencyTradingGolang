package clientProcess

import (
	"fmt"
	"winter/global"
	"winter/messages"
	"winter/utils"
)

func (con *Convert) FTXQuoteRawToStand(msg string, localtime int64) messages.AggStMsg {
	var stMsg messages.AggStMsg = messages.AggStMsg{}

	var ord messages.FTXOrderResponse
	err := jsonIterator.Unmarshal([]byte(msg), &ord)
	if err != nil {
		logger.Info("error: Unmarshal OrderInitResponse Failed.")
	}

	var ii int = con.FTX_ii[ord.Market]

	if ord.Type == "partial" {
		con.FTXOrderbook[ii].Reset(&ord, FTX_Kind_Decimals[ord.Market])
	} else if ord.Type == "update" {
		con.FTXOrderbook[ii].Update(&ord)
	}

	con.FTXOrderbook[ii].CheckSum()

	stMsg.ExID = 0             // FTX交易所
	stMsg.Status = 1           // 状态是活着
	stMsg.Kinds = 0            // Quote 数据
	stMsg.Symbols = ord.Market // 币种
	stMsg.OrdEvt = messages.OrderbookBase{
		Localtime: localtime,
		Eventime:  int64(ord.Data.Time * 1000),
		Asks:      con.FTXOrderbook[ii].GetCurrentOrderbook("ask"),
		Bids:      con.FTXOrderbook[ii].GetCurrentOrderbook("bid"),
	}
	return stMsg
}

func (con *Convert) FTXTradeRawToStand(msg string, localtime int64) []messages.AggStMsg {
	var stMsg []messages.AggStMsg = []messages.AggStMsg{}
	var trd messages.FTXTradeResponse
	err := jsonIterator.Unmarshal([]byte(msg), &trd)
	if err != nil {
		logger.Info("error: Unmarshal TradeInitResponse Failed.")
	}

	for i := range trd.Data {
		_stMsg := messages.AggStMsg{}
		var side float64 = 0
		if trd.Data[i].Side == "buy" {
			side = 1
		} else if trd.Data[i].Side == "sell" {
			side = -1
		}
		_stMsg.ExID = 0 // FTX交易所
		_stMsg.Status = 1
		_stMsg.Kinds = 1
		_stMsg.Symbols = trd.Market
		_stMsg.TrdEvt = messages.TradeBase{
			Localtime: localtime,
			Eventime:  0,
			Price:     trd.Data[i].Price,
			Size:      trd.Data[i].Size,
			Side:      side,
			IsTaker:   trd.Data[i].Liquidation,
		}
		stMsg = append(stMsg, _stMsg)
	}

	return stMsg
}

func (con *Convert) OkexQuoteRawToStand(msg string, localtime int64) messages.AggStMsg {
	var stMsg messages.AggStMsg = messages.AggStMsg{}

	var ord messages.OkexOrderResponse
	err := jsonIterator.Unmarshal([]byte(msg), &ord)
	if err != nil {
		logger.Info("error: Unmarshal OkexOrderResponse Failed. " + err.Error())
	}

	var ii int = con.Okex_ii[ord.Arg.InstId]

	if ord.Action == "snapshot" {
		con.OkexOrderbook[ii].Reset(&ord, global.Okex_Kind_Decimals[ord.Arg.InstId]) //okexQteEvt AskSkipList BidSkipList
	} else if ord.Action == "update" {
		con.OkexOrderbook[ii].Update(&ord) // find change insert
	} else {
		logger.Info("error Action.")
	}

	if global.IsDebugChecksum {
		if !con.OkexOrderbook[ii].CheckSum() {
			logger.Info(ord.Arg.InstId)
			fmt.Println("============================================")
		}
	}

	stMsg.ExID = 1 // 代表Okex交易所
	stMsg.Status = 1
	stMsg.Kinds = 0
	stMsg.Symbols = ord.Arg.InstId

	stMsg.OrdEvt = messages.OrderbookBase{
		Localtime: localtime,
		Eventime:  utils.String2int64(ord.Data[0].ExchangeTm),
		Asks:      con.OkexOrderbook[ii].GetCurrentOrderbook("ask"),
		Bids:      con.OkexOrderbook[ii].GetCurrentOrderbook("bid"),
	}
	return stMsg
}

func (con *Convert) OkexTradeRawToStand(msg string, localtime int64) messages.AggStMsg {
	var stMsg messages.AggStMsg = messages.AggStMsg{}

	var trd messages.OkexTradeResponse
	err := jsonIterator.Unmarshal([]byte(msg), &trd)
	if err != nil {
		logger.Info("error: Unmarshal OkexTradeResponse Failed.")
	}

	var _side float64
	if trd.Data[0].Side == "buy" {
		_side = 1
	} else if trd.Data[0].Side == "sell" {
		_side = -1
	} else {
		_side = 0
	}

	stMsg.ExID = 1
	stMsg.Status = 1
	stMsg.Kinds = 1
	stMsg.Symbols = trd.Arg.InstId
	stMsg.TrdEvt = messages.TradeBase{
		Localtime: localtime,
		Eventime:  utils.String2int64(trd.Data[0].ExchangeTm),
		Price:     utils.String2float(trd.Data[0].Price),
		Size:      utils.String2float(trd.Data[0].Size),
		Side:      _side,
	}

	return stMsg
}

func (con *Convert) OkexKlineRawToStand(msg string, localtime int64) messages.AggStMsg {
	var stMsg messages.AggStMsg = messages.AggStMsg{}

	var kline messages.OkexKlineResponse
	err := jsonIterator.Unmarshal([]byte(msg), &kline)
	if err != nil {
		logger.Info("error: Unmarshal OkexKlineResponse Failed.")
	}

	var freq int64
	switch kline.Arg.Channel {
	case "candle1m":
		freq = 60000
	default:
		logger.Info("the channel not set.")
	}

	stMsg.ExID = 1
	stMsg.Status = 1
	stMsg.Kinds = 2
	stMsg.Symbols = kline.Arg.InstId
	stMsg.KlineEvt = messages.KlineBase{
		Frequency: freq,
		Localtime: localtime,
		StartTm:   utils.String2int64(kline.Data[0][0]),
		Open:      utils.String2float(kline.Data[0][1]),
		High:      utils.String2float(kline.Data[0][2]),
		Low:       utils.String2float(kline.Data[0][3]),
		Close:     utils.String2float(kline.Data[0][4]),
		Vol:       utils.String2float(kline.Data[0][5]),
		VolCcy:    utils.String2float(kline.Data[0][6]),
	}

	return stMsg
}

func (con *Convert) BinanceQuoteRawToStand(msg string, localtime int64) messages.AggStMsg {
	var stMsg messages.AggStMsg = messages.AggStMsg{}
	var qte messages.BinanceOrderResponse
	err := jsonIterator.Unmarshal([]byte(msg), &qte)
	if err != nil {
		logger.Info("error: Unmarshal BinanceQuoteResponse Failed.")
	}

	var asks, bids [50][4]float64
	for i := range asks {
		if i < len(qte.Asks) {
			asks[i] = [4]float64{
				utils.String2float(qte.Asks[i][0]),
				utils.String2float(qte.Asks[i][1]),
				0,
				0,
			}
		} else {
			asks[i] = [4]float64{0, 0, 0, 0}
		}
	}

	for i := range bids {
		if i < len(qte.Bids) {
			bids[i] = [4]float64{
				utils.String2float(qte.Bids[i][0]),
				utils.String2float(qte.Bids[i][1]),
				0,
				0,
			}
		} else {
			bids[i] = [4]float64{0, 0, 0, 0}
		}
	}

	stMsg.ExID = 2
	stMsg.Status = 1
	stMsg.Kinds = 0
	stMsg.Symbols = qte.Symbol
	stMsg.OrdEvt = messages.OrderbookBase{
		Localtime: localtime,
		Eventime:  localtime,
		Asks:      asks,
		Bids:      bids,
	}

	return stMsg
}
func (con *Convert) BinanceTradeRawToStand(msg string, localtime int64) messages.AggStMsg {
	var stMsg messages.AggStMsg = messages.AggStMsg{ExID: 2, Status: 1, Kinds: 1}

	var trd messages.BinanceTradeResponse
	err := jsonIterator.Unmarshal([]byte(msg), &trd)
	if err != nil {
		logger.Info("error: Unmarshal BinanceTradeResponse Failed.")
	}

	var _side float64
	if trd.IsMaker {
		_side = -1
	} else {
		_side = 1
	}

	stMsg.Symbols = trd.Symbol
	stMsg.TrdEvt = messages.TradeBase{
		Localtime: localtime,
		Eventime:  localtime,
		Price:     utils.String2float(trd.Price),
		Size:      utils.String2float(trd.Size),
		Side:      _side,
	}

	return stMsg
}

// Frequency int64
// 	Localtime int64
// 	StartTm   int64

// 	Open   float64
// 	High   float64
// 	Low    float64
// 	Close  float64
// 	Vol    float64
// 	VolCcy float64

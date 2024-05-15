/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:22
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 19:19:48
 */
package messages

import (
	"sync/atomic"
)

// import "winter/logger"

// {"arg":{"channel":"trades","instId":"BTC-USDT-SWAP"},"data":[{"instId":"BTC-USDT-SWAP","tradeId":"119555440","px":"49549.9","sz":"19","side":"sell","ts":"1633421772969"}]}}

// {"arg":{"channel":"books5","instId":"ETH-USDT-SWAP"},"data":[{"asks":[["3388.55","761","0","19"],["3388.56","9","0","1"],["3388.64","53","0","2"],["3388.65","164","0","1"],["3388.79","38","0","2"]],"bids":[["3388.54","74","0","19"],["3388.46","2","0","1"],["3388.38","1","0","1"],["3388.29","9","0","1"],["3388.17","50","0","2"]],"instId":"ETH-USDT-SWAP","ts":"1633421772925"}]}

type OkexSubMsg_ struct {
	Channel string `json:"channel"`
	InstId  string `json:"instId"`
}

type OkexSubMsg struct {
	Op   string        `json:"op"`
	Args []OkexSubMsg_ `json:"args"`
}

func NewOkexSubMsg(symbols []string, msgkind string) OkexSubMsg {
	ans := OkexSubMsg{}
	switch msgkind {
	case "Trade":
		ans.Op = "subscribe"
		for i := range symbols {
			ans.Args = append(ans.Args, OkexSubMsg_{Channel: "trades", InstId: symbols[i]})
		}
	case "Order":
		ans.Op = "subscribe"
		for i := range symbols {
			ans.Args = append(ans.Args, OkexSubMsg_{Channel: "books", InstId: symbols[i]})
		}
	case "Candle1m":
		ans.Op = "subscribe"
		for i := range symbols {
			ans.Args = append(ans.Args, OkexSubMsg_{Channel: "candle1m", InstId: symbols[i]})
		}
	default:
		logger.Info("error msgkind when okex subscribe")
	}

	return ans

}

type BinanceSubMsg struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int64    `json:"id"`
}

func NewBinanceSubMsg(symbols []string, msgkind string) BinanceSubMsg {
	ans := BinanceSubMsg{}
	switch msgkind {
	case "Trade":
		ans.Method = "SUBSCRIBE"
		for i := range symbols {
			// "btcusdt@trade"
			ans.Params = append(ans.Params, symbols[i]+"@aggTrade")
		}
		ans.ID = global_sub_id
		atomic.AddInt64(&global_sub_id, 1)
	case "Order":
		ans.Method = "SUBSCRIBE"
		for i := range symbols {
			// <symbol>@depth<levels>@100ms.
			// ans.Args = append(ans.Args, OkexSubMsg_{Channel: "books", InstId: symbols[i]})
			ans.Params = append(ans.Params, symbols[i]+"@depth20@100ms")
		}
		ans.ID = global_sub_id
		atomic.AddInt64(&global_sub_id, 1)
	case "Candle1m":
		logger.Info("kline not implemented in binance")
	default:
		logger.Info("error msgkind when binance subscribe")
	}
	return ans

}

// type OkexOrderUpdateContents struct {
// 	Offset string      `json:"offset"`
// 	Bids   [][2]string `json:"bids"`
// 	Asks   [][2]string `json:"asks"`
// }

// type OkexOrderUpdateResponse struct {
// 	Type     string                  `json:"type"`
// 	ConID    string                  `json:"connection_id"`
// 	MsgID    int                     `json:"message_id"`
// 	ID       string                  `json:"id"`
// 	Channel  string                  `json:"channel"`
// 	Contents OkexOrderUpdateContents `json:"contents"`
// }

// {"instId":"BTC-USDT-SWAP","tradeId":"119555440","px":"49549.9","sz":"19","side":"sell","ts":"1633421772969"}
type OkexTradeContents struct {
	InstId     string `json:"instId"`
	TradeId    string `json:"tradeId"`
	Price      string `json:"px"`
	Size       string `json:"sz"`
	Side       string `json:"side"`
	ExchangeTm string `json:"ts"`
}

type OkexTradeResponse struct {
	Arg  OkexSubMsg_         `json:"arg"`
	Data []OkexTradeContents `json:"data"`
}

// {
//     "arg": {
//         "channel": "books",
//         "instId": "BTC-USDT"
//     },
//     "action": "snapshot",
//     "data": [{
//         "asks": [
//             ["8476.98", "415", "0", "13"],
//             ["8477", "7", "0", "2"],
//             ["8477.34", "85", "0", "1"],
//             ["8477.56", "1", "0", "1"],
//             ["8505.84", "8", "0", "1"],
//             ["8506.37", "85", "0", "1"],
//             ["8506.49", "2", "0", "1"],
//             ["8506.96", "100", "0", "2"]
//         ],
//         "bids": [
//             ["8476.97", "256", "0", "12"],
//             ["8475.55", "101", "0", "1"],
//             ["8475.54", "100", "0", "1"],
//             ["8475.3", "1", "0", "1"],
//             ["8447.32", "6", "0", "1"],
//             ["8447.02", "246", "0", "1"],
//             ["8446.83", "24", "0", "1"],
//             ["8446", "95", "0", "3"]
//         ],
//         "ts": "1597026383085",
//         "checksum": -855196043
//     }]
// }

type OkexOrderContents struct {
	Asks       [][4]string `json:"asks"`
	Bids       [][4]string `json:"bids"`
	ExchangeTm string      `json:"ts"`
	CheckSum   int32       `json:"checksum"`
}

type OkexOrderResponse struct {
	Arg    OkexSubMsg_         `json:"arg"`
	Action string              `json:"action"`
	Data   []OkexOrderContents `json:"data"`
}

type OkexKlineResponse struct {
	Arg  OkexSubMsg_ `json:"arg"`
	Data [][]string  `json:"data"`
}

// 只有必须的字段，非必须字段都没有写
type OkexTradeArgs struct {
	InstId  string
	TdMode  string
	Side    string
	OrdType string
	Size    string
	Price   string
	ClOrdId string
}

type OkexTradeRequest struct {
	ID   string        `json:"id"`
	Op   string        `json:"op"`
	Args OkexTradeArgs `json:"args"`
}

//	{
//	    "id": "1512",
//	    "op": "order",
//	    "data": [{
//	        "clOrdId": "",
//	        "ordId": "12345689",
//	        "tag": "",
//	        "sCode": "0",
//	        "sMsg": ""
//	    }],
//	    "code": "0",
//	    "msg": ""
//	}
type OkexTradeChannelMsgArg struct {
	Channel  string `json:"channel"`
	InstType string `json:"instType"`
	Uid      string `json:"uid"`
}
type OkexTradeChannelMsgContent struct {
	AccFillSz       string `json:"accFillSz"`   // 累计成交数量
	AmendResult     string `json:"amendResult"` // 修改订单的结果
	AveragePrice    string `json:"avgPx"`
	CTime           string `json:"cTime"`
	Category        string `json:"category"`
	Ccy             string `json:"ccy"`
	ClOrdId         string `json:"clOrdId"`
	Code            string `json:"code"`
	ExecType        string `json:"execType"`
	Fee             string `json:"fee"`
	FeeCcy          string `json:"feeCcy"`
	FillFee         string `json:"fillFee"`
	FillFeeCcy      string `json:"fillFeeCcy"`
	FillNotionalUsd string `json:"fillNotionalUsd"`
	FillPx          string `json:"fillPx"`
	FillSz          string `json:"fillSz"`
	FillTime        string `json:"fillTime"`
	InstId          string `json:"instId"`
	InstType        string `json:"instType"`
	Lever           string `json:"lever"`
	Msg             string `json:"msg"`
	NotionalUsd     string `json:"notionalUsd"`
	OrdId           string `json:"ordId"`
	OrdType         string `json:"ordType"`
	Pnl             string `json:"pnl"`
	PosSide         string `json:"posSide"`
	Px              string `json:"px"`
	Rebate          string `json:"rebate"`
	RebateCcy       string `json:"rebateCcy"`
	ReduceOnly      string `json:"reduceOnly"`
	ReqId           string `json:"reqId"`
	Side            string `json:"side"`
	SlOrdPx         string `json:"slOrdPx"`
	SlTriggerPx     string `json:"slTriggerPx"`
	SlTriggerPxType string `json:"slTriggerPxType"`
	Source          string `json:"source"`
	State           string `json:"state"`
	Size            string `json:"sz"`
	Tag             string `json:"tag"`
	TdMode          string `json:"tdMode"`
	TgtCcy          string `json:"tgtCcy"`
	TpOrdPx         string `json:"tpOrdPx"`
	TpTriggerPx     string `json:"tpTriggerPx"`
	TpTriggerPxType string `json:"tpTriggerPxType"`
	TradeId         string `json:"tradeId"`
	UTime           string `json:"uTime"`
}
type OkexTradeChannelMsg struct {
	Arg  OkexTradeChannelMsgArg       `json:"arg"`
	Data []OkexTradeChannelMsgContent `json:"data"`
}

type OkexWebTradeResponseContent struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
	Tag     string `json:"tag"`
}

type OkexWebTradeResponse struct {
	Code string                        `json:"code"`
	Data []OkexWebTradeResponseContent `json:"data"`
	ID   string                        `json:"id"`
	Msg  string                        `json:"msg"`
	Op   string                        `json:"op"`
}

type OkexChannelResponse struct {
	Event string                 `json:"event"`
	Arg   OkexTradeChannelMsgArg `json:"arg"`
	Code  string                 `json:"code"`
	Msg   string                 `json:"msg"`
}

// {"arg":{"channel":"orders","instType":"SWAP","uid":"103052175277699072"},"data":[{"accFillSz":"0","amendResult":"","avgPx":"0","cTime":"1641572720764","category":"normal","ccy":"","clOrdId":"","code":"0","execType":"","fee":"0","feeCcy":"USDT","fillFee":"0","fillFeeCcy":"","fillNotionalUsd":"","fillPx":"","fillSz":"0","fillTime":"","instId":"ETH-USDT-SWAP","instType":"SWAP","lever":"5","msg":"","notionalUsd":"1248.5040000000001","ordId":"399720194305851397","ordType":"market","pnl":"0","posSide":"net","px":"","rebate":"0","rebateCcy":"USDT","reduceOnly":"false","reqId":"","side":"buy","slOrdPx":"","slTriggerPx":"","slTriggerPxType":"","source":"","state":"live","sz":"4","tag":"","tdMode":"isolated","tgtCcy":"","tpOrdPx":"","tpTriggerPx":"","tpTriggerPxType":"","tradeId":"","uTime":"1641572720764"}]}
// {"arg":{"channel":"orders","instType":"SWAP","uid":"103052175277699072"},"data":[{"accFillSz":"4","amendResult":"","avgPx":"3140.49","cTime":"1641572720764","category":"normal","ccy":"","clOrdId":"","code":"0","execType":"T","fee":"-0.628098","feeCcy":"USDT","fillFee":"-0.628098","fillFeeCcy":"USDT","fillNotionalUsd":"1248.5040000000001","fillPx":"3140.49","fillSz":"4","fillTime":"1641572720770","instId":"ETH-USDT-SWAP","instType":"SWAP","lever":"5","msg":"","notionalUsd":"1248.5040000000001","ordId":"399720194305851397","ordType":"market","pnl":"0.02","posSide":"net","px":"","rebate":"0","rebateCcy":"USDT","reduceOnly":"false","reqId":"","side":"buy","slOrdPx":"","slTriggerPx":"","slTriggerPxType":"","source":"","state":"filled","sz":"4","tag":"","tdMode":"isolated","tgtCcy":"","tpOrdPx":"","tpTriggerPx":"","tpTriggerPxType":"","tradeId":"171510135","uTime":"1641572720777"}]}

/*
// Standard
*/

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:15
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-22 00:52:23
 */
package messages

// FTX相关结构体，用于解析json

// Order更新用的结构体
type FTXOrderBaseContent struct {
	Time     float64      `json:"time"`
	Checksum uint32       `json:"checksum"`
	Bids     [][2]float64 `json:"bids"`
	Asks     [][2]float64 `json:"asks"`
	Action   string       `json:"action"`
}

type FTXOrderResponse struct {
	Channel string              `json:"channel"`
	Market  string              `json:"market"`
	Type    string              `json:"type"`
	Data    FTXOrderBaseContent `json:"data"`
}

// Trade更新用的结构体
// "id":4080139483,
//             "price":1975.4,
//             "size":0.001,
//             "side":"sell",
//             "liquidation":false,
//             "time":"2022-05-21T16:44:54.695451+00:00"
type FTXTradeContent struct {
	ID          int64   `json:"id"`
	Price       float64 `json:"price"`
	Size        float64 `json:"size"`
	Side        string  `json:"side"`
	Liquidation bool    `json:"liquidation"`
	Time        string  `json:"time"`
}

type FTXTradeResponse struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`

	Data []FTXTradeContent `json:"data"`
}

type FTXTradeUpdateContent struct {
	Size  string `json:"size"`
	Side  string `json:"side"`
	Price string `json:"price"`
	Tm    string `json:"createdAt"`
}

type FTXTradeUpdateContents struct {
	Trds []FTXTradeUpdateContent `json:"trades"`
}

type FTXTradeUpdateResponse struct {
	Type     string                 `json:"type"`
	ConID    string                 `json:"connection_id"`
	MsgID    int                    `json:"message_id"`
	ID       string                 `json:"id"`
	Channel  string                 `json:"channel"`
	Contents FTXTradeUpdateContents `json:"contents"`
}

type FTXSubMsg struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

func NewFTXSubMsg(symbol string, msgkind string) FTXSubMsg {
	ans := FTXSubMsg{}
	switch msgkind {
	case "Trade":
		ans.Op = "subscribe"
		ans.Channel = "trades"
		ans.Market = symbol
	case "Order":
		ans.Op = "subscribe"
		ans.Channel = "orderbook"
		ans.Market = symbol
	default:
		logger.Info("error msgkind when ftx subscribe")
	}

	return ans
}

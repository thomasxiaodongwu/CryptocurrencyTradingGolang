/*
 * @Author: xwu
 * @Date: 2021-12-26 11:18:12
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-24 16:20:19
 */
package clientTrader

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
	"winter/global"
	"winter/messages"
	"winter/utils"

	"github.com/gorilla/websocket"
)

/*
  限制 60次/s
  首先订阅订单频道
  websocket下单，收到回复证明交易所收到了但是不代表上线了
  在订单频道收到live
  订单频道收到filled
  V5 API 已支持所有产品类型的改单，让您修改订单的价格（newPx字段）和/或数量（newSz字段）。另外 API 也提供cxlOnFail参数，设置修改失败时自动撤单。
  同样，我们应会收到服务器相应 REST / WebSocket 的成功返回。当我们从 WebSocket 订单频道收到订单状态为canceled的推送更新时，才代表订单撤单成功。
*/

func NewOkexTraderClient() OkexTraderClient {
	ans := OkexTraderClient{}
	ans.Count = 0
	ans.IsConnected = false
	ans.mutex = new(sync.RWMutex)
	return ans
}

type OkexTraderClient struct {
	Conn        *websocket.Conn
	Count       int
	mutex       *sync.RWMutex
	IsConnected bool
}

func (client *OkexTraderClient) setConnFlag(Symbol bool) {
	client.mutex.Lock()
	client.IsConnected = Symbol
	client.mutex.Unlock()
}

func (client *OkexTraderClient) getConnFlag() bool {
	client.mutex.RLock()
	temp := client.IsConnected
	client.mutex.RUnlock()
	return temp
}

func (client *OkexTraderClient) PingOnce() bool {
	if client.IsConnected {
		pingMsg := "ping"
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(pingMsg))
		if err != nil {
			logger.Info("okex ping failed " + err.Error())
			return false
		}
		return true
	} else {
		return false
	}
}

func (client *OkexTraderClient) Ping() {
	for {
		if client.IsConnected {
			pingMsg := "ping"
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(pingMsg))
			if err != nil {
				logger.Info("okex ping failed " + err.Error())
				time.Sleep(time.Duration(1) * time.Second)
			}
			time.Sleep(time.Duration(20) * time.Second)
		} else {
			time.Sleep(time.Duration(3) * time.Second)
		}
	}
}

func (client *OkexTraderClient) Subcrsibe(channel string, instType string) error {
	// ex channel: orders, instType: SWAP
	var err error

	subMsg := utils.Strcat(
		`{"op":"subscribe","args": [{"channel":"`,
		channel,
		`","instType":"`,
		instType,
		`"}]}`,
	)

	err = client.Conn.WriteMessage(websocket.TextMessage, []byte(subMsg))
	if err != nil {
		logger.Info(`Okex subMsg Message Sent Failed.`)
	} else {
		logger.Info(`Okex subMsg Message Sent Successed.`)
	}

	return err
}

func (client *OkexTraderClient) PostMarketOrderSimple(ExID int, Symbol string, side string, size string, count int, strategy_name string) {
	req := new(messages.OkexTradeRequest)
	req.ID = utils.Strcat(
		strategy_name,
		strconv.FormatInt(int64(count), 10),
	) // 客户端id使用策略名字,方便后面logger的输出
	req.Args.OrdType = "market"
	req.Args.Side = side
	req.Args.InstId = Symbol
	req.Args.Size = size
	req.Args.ClOrdId = strategy_name

	_a := utils.Strcat(
		strategy_name,
		strconv.FormatInt(int64(count), 10),
	)
	subMsg := client.market_order(
		_a,
		side,
		Symbol,
		size,
		_a,
	)
	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(subMsg))
	if err != nil {
		logger.Error("client.Conn.WriteMessage fail" + err.Error())
	}
}

func (client *OkexTraderClient) PostMarketOrder(ExID int, Symbol string, side string, size string, count int, strategy_name string, signal messages.AggSignal, hft *HftClient) {
	req := new(messages.OkexTradeRequest)
	req.ID = utils.Strcat(
		strategy_name,
		strconv.FormatInt(int64(count), 10),
	) // 客户端id使用策略名字,方便后面logger的输出
	req.Args.OrdType = "market"
	req.Args.Side = side
	req.Args.InstId = Symbol
	req.Args.Size = size
	req.Args.ClOrdId = strategy_name

	_a := utils.Strcat(
		strategy_name,
		strconv.FormatInt(int64(count), 10),
	)
	subMsg := client.market_order(
		_a,
		side,
		Symbol,
		size,
		_a,
	)

	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(subMsg))
	if err != nil {
		logger.Error("PostMarketOrder failed " + err.Error())
		go utils.FeiShuMsg(utils.Strcat(
			"TradeClient: PostMarketOrder failed. ",
			Symbol,
			"-",
			side,
			"-",
			size,
			"-",
			req.ID,
		))
	}
}

// 下单,
func (client *OkexTraderClient) Trade(params *messages.OkexTradeRequest) {
	var subMsg string = ""

	switch params.Args.OrdType {
	case "market":
		subMsg = client.market_order(
			params.ID,
			params.Args.Side,
			params.Args.InstId,
			params.Args.Size,
			params.Args.ClOrdId,
		)
	case "limit":
	case "post_only":
	case "fok":
	case "ioc":
	case "optimal_limit_ioc":
	default:
		logger.Info("Unknown order type.")
	}

	// logger.Info(subMsg)

	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(subMsg))
	if err != nil {
		logger.Info(`Okex Trade Message Sent Failed.` + err.Error())
	}
}

// 市价单，全仓，usdt保证金
func (client *OkexTraderClient) market_order(
	strate_order_id string,
	side string,
	instId string,
	size string,
	clOrdId string,
) string {
	var msg string = utils.Strcat(
		`{"id":"`,
		strate_order_id,
		`","op":"order","args":[{"side":"`,
		side,
		`","instId":"`,
		instId,
		`","clOrdId":"`,
		clOrdId,
		`","tdMode":"isolated","ordType":"market","sz":"`,
		size,
		`"}]}`,
	)

	return msg
}

// 限价单
func (client *OkexTraderClient) limit_order(strate_order_id string) string {
	var msg string
	return msg
}

// 只做maker单
func (client *OkexTraderClient) post_only_order() string {
	var msg string
	return msg
}

// 全部成交或立即取消
func (client *OkexTraderClient) fok_order() string {
	var msg string
	return msg
}

// 立即成交并取消剩余
func (client *OkexTraderClient) ioc_order() string {
	var msg string
	return msg
}

// 市价委托立即成交并取消剩余（仅适用交割、永续）
func (client *OkexTraderClient) optimal_limit_ioc() string {
	var msg string
	return msg
}

func (client *OkexTraderClient) Create_Private_Dialogue() error {
	var err error
	if global.USING_PROXY {
		uProxy, _ := url.Parse(global.PROXY_URI)
		dialer := websocket.Dialer{
			Proxy: http.ProxyURL(uProxy),
		}
		client.Conn, _, err = dialer.Dial(global.WEBSOCKET_OKEX_PRIVATE_URI, nil)
	} else {
		client.Conn, _, err = websocket.DefaultDialer.Dial(global.WEBSOCKET_OKEX_PRIVATE_URI, nil)
	}

	// 是否连接成功
	if err != nil {
		logger.Info("send message to server failed  " + err.Error())
	} else {
		logger.Info("send message to server success. ")
	}
	client.Count += 1

	timestamp := time.Now().UnixNano()/1e9 - 3
	sign := utils.OkexSign(timestamp)

	loginMessage := utils.Strcat(
		`{"op":"login","args":[{"apiKey":"`,
		global.Okex_ApiKey,
		`","passphrase":"`,
		global.Okex_Passphrase,
		`","timestamp":"`,
		strconv.FormatInt(timestamp, 10),
		`","sign":"`,
		sign,
		`"}]}`,
	)

	err = client.Conn.WriteMessage(websocket.TextMessage, []byte(loginMessage))
	if err != nil {
		logger.Info(`Okex Login Message Sent Failed.`)
		client.setConnFlag(false)
	} else {
		logger.Info(`Okex Login Message Sent Successed.`)
		//client.IsConnected = true
		client.setConnFlag(true)
	}

	return err
}

func (client *OkexTraderClient) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		err := client.Create_Private_Dialogue()
		if err != nil {
			time.Sleep(time.Duration(7))
			continue
		}
		for {
			_, message_okex, err := client.Conn.ReadMessage()

			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) || err != nil {
				break
			}

			fmt.Println(string(message_okex))
		}
	}
}

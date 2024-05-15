package clientData

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"winter/global"
	"winter/messages"
	"winter/utils"

	"github.com/gorilla/websocket"
)

type FTXClient struct {
	Is_connected bool
	Conn         *websocket.Conn
	Count        int
}

// USING_PROXY = true
// PROXY_URI   = "http://127.0.0.1:10809"
func (ftx *FTXClient) connected() bool {
	var is_success bool = false
	var err error
	// 是否使用代理
	if global.USING_PROXY {
		uProxy, _ := url.Parse(global.PROXY_URI)
		dialer := websocket.Dialer{
			Proxy: http.ProxyURL(uProxy),
		}
		ftx.Conn, _, err = dialer.Dial(global.WEBSOCKET_FTX_URI, nil)
	} else {
		ftx.Conn, _, err = websocket.DefaultDialer.Dial(global.WEBSOCKET_FTX_URI, nil)
	}

	// 是否连接成功
	if err != nil {
		logger.Info("Failed connect the server  " + err.Error())
	} else {
		logger.Info(fmt.Sprintf("Successful connection for %dth time", ftx.Count+1))
		is_success = true
	}
	ftx.Count += 1

	// 返回是否成功
	return is_success
}

func (ftx *FTXClient) subscribe(symbols []string) {
	// ex: BTC-USDT-SWAP_Order
	// ex: BTC-USDT-SWAP_Trade
	// ex: BTC-USDT-SWAP_Candle1m

	var _symbols map[string][]string = make(map[string][]string) //哪个频道需要订阅哪些

	for i := range symbols {
		infos := strings.Split(symbols[i], "_")
		_symbols[infos[1]] = append(_symbols[infos[1]], infos[0])
	}

	for channelName, symbolNames := range _symbols {
		switch channelName {
		case "Trade":
			for _, symbol := range symbolNames {
				_sub := messages.NewFTXSubMsg(symbol, "Trade")
				sub, _ := json.Marshal(_sub)
				err := ftx.Conn.WriteMessage(websocket.TextMessage, sub)
				if err != nil {
					logger.Info(utils.Strcat(`ftx Subscribe `, channelName, ` Failed.`))
				} else {
					logger.Info(utils.Strcat(`ftx Subscribe `, channelName, ` Successed.`))
				}
				time.Sleep(100 * time.Millisecond)
			}
		case "Order":
			for _, symbol := range symbolNames {
				_sub := messages.NewFTXSubMsg(symbol, "Order")
				sub, _ := json.Marshal(_sub)
				err := ftx.Conn.WriteMessage(websocket.TextMessage, sub)
				if err != nil {
					logger.Info(utils.Strcat(`ftx Subscribe Order Failed.`))
				} else {
					logger.Info(utils.Strcat(`ftx Subscribe ` + symbol + ` Successed.`))
				}
				time.Sleep(100 * time.Millisecond)
			}
		default:
			logger.Info("ftx error channelName.")
		}

	}
}

func (ftx *FTXClient) Ping() {
	for {
		if ftx.Is_connected {
			pingMsg := `{'op': 'ping'}`
			err := ftx.Conn.WriteMessage(websocket.TextMessage, []byte(pingMsg))
			if err != nil {
				logger.Info("ftx ping failed" + err.Error())
			}
			time.Sleep(time.Duration(20) * time.Second)
		} else {
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}

func (ftx *FTXClient) Run(chanMsg chan messages.MsgDataToProcess, wg *sync.WaitGroup) {
	p := global.AggParameters.Data
	defer wg.Done()

	for {
		is_connected := ftx.connected()
		if !is_connected {
			time.Sleep(time.Duration(7) * time.Second)
			continue
		} else {
			ftx.Is_connected = true
		}

		ftx.subscribe(p.FTX.Subscribe_symbols)

		for {
			var err error
			var message_ftx []byte

			if ftx.Conn != nil {
				_, message_ftx, err = ftx.Conn.ReadMessage()
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) || err != nil {
					ftx.Is_connected = false
					break
				}
			} else {
				ftx.Is_connected = false
				break
			}

			ct := time.Now().UnixNano()

			var msg_data messages.MsgDataToProcess = messages.MsgDataToProcess{
				ExID:     0,
				Contents: utils.Strcat(utils.ConvertMillisecondString(ct/1e6), "|", string(message_ftx)),
			}

			if strings.Contains(msg_data.Contents, "data") {
				chanMsg <- msg_data
			}
		}

		ftx.Conn.Close()
		ftx.Conn = nil
		time.Sleep(7 * time.Second)
	}
}

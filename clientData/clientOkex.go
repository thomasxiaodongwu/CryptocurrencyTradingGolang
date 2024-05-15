/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:30
 * @Last Modified by: xwu
 * @Last Modified time: 2022-06-07 10:00:10
 */
package clientData

import (
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

type OkexClient struct {
	Is_connected bool
	Conn         *websocket.Conn
	Count        int
	CurrentDay   int
}

// USING_PROXY = true
// PROXY_URI   = "http://127.0.0.1:10809"
func (okex *OkexClient) connected() bool {
	// 只要不是第一次，这个函数被调用就会发送一个飞书消息，除了每天零点固定重连，其他时间都是异常，但是有可能是okx的网络波动
	if okex.Count > 0 {
		go utils.FeiShuMsg("Okex DataClient: reConnect")
	}

	if !okex.Is_connected { // 如果没有连接 okex.Is_connected=false
		var err error
		// 是否使用代理
		if global.USING_PROXY {
			uProxy, _ := url.Parse(global.PROXY_URI)
			dialer := websocket.Dialer{
				Proxy: http.ProxyURL(uProxy),
			}
			okex.Conn, _, err = dialer.Dial(global.WEBSOCKET_OKEX_URI, nil)
		} else {
			okex.Conn, _, err = websocket.DefaultDialer.Dial(global.WEBSOCKET_OKEX_URI, nil)
		}

		// 是否连接成功
		if err != nil {
			logger.Info("Failed connect the server  " + err.Error())
		} else {
			logger.Info(fmt.Sprintf("Successful connection for %dth time", okex.Count+1))
			okex.Is_connected = true
		}
		okex.Count += 1
	} else {
		// Ping一下
		okex.Ping()
	}
	// 返回是否成功
	return okex.Is_connected
}

func (okex *OkexClient) subscribe(symbols []string) error {
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
			sub := utils.SubMsgFormat(symbolNames, "Trade", "Okex")
			if len(sub) > 4096 {
				logger.Fatal("length of Okex trade channel msg > 4096.")
			}
			err := okex.Conn.WriteMessage(websocket.TextMessage, sub)
			if err != nil {
				logger.Info(utils.Strcat(`Okex Subscribe `, channelName, ` Failed.`))
				errorMsg := fmt.Errorf("conn")
				return errorMsg
			} else {
				logger.Info(utils.Strcat(`Okex Subscribe `, channelName, ` Successed.`))
			}
		case "Order":
			sub := utils.SubMsgFormat(symbolNames, "Order", "Okex")
			if len(sub) > 4096 {
				logger.Fatal("length of Okex Order channel msg > 4096.")
			}

			err := okex.Conn.WriteMessage(websocket.TextMessage, sub)
			if err != nil {
				logger.Info(utils.Strcat(`Okex Subscribe Order Failed.`))
				errorMsg := fmt.Errorf("conn")
				return errorMsg
			} else {
				logger.Info(utils.Strcat(`Okex Subscribe Order Successed.`))
			}
		case "Candle1m":
			sub := utils.SubMsgFormat(symbolNames, "Candle1m", "Okex")
			if len(sub) > 4096 {
				logger.Fatal("length of Okex Candle1m channel msg > 4096.")
			}

			err := okex.Conn.WriteMessage(websocket.TextMessage, sub)
			if err != nil {
				logger.Info(utils.Strcat(`Okex Subscribe Candle1m Failed.`))
				errorMsg := fmt.Errorf("conn")
				return errorMsg
			} else {
				logger.Info(utils.Strcat(`Okex Subscribe Candle1m Successed.`))
			}
		default:
			logger.Info("Okex error channelName.")
		}
	}
	return nil
}

func (okex *OkexClient) Ping() {
	if okex.Is_connected {
		pingMsg := "ping"
		err := okex.Conn.WriteMessage(websocket.TextMessage, []byte(pingMsg))
		if err != nil {
			logger.Info("okex ping failed" + err.Error())
			okex.Is_connected = false // ping失败所以需要重新连接
		}
	}
}

func (okex *OkexClient) Run(chanMsg chan messages.MsgDataToProcess, wg *sync.WaitGroup) {
	p := global.AggParameters.Data
	defer wg.Done()

	okex.CurrentDay = time.Now().Day() // 初始化天数

	for { // 这里是第一次连接，一次性代码不涉及重连
		is_connected := okex.connected()
		if is_connected { // 连接成功
			// okex.Is_connected = true // 这里在connected函数里面有过赋值
			break
		} else { // 连接失败
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
	}

	for {
		// 这里如果是第一次进来，okex.Conn一定不是nil，如果是break出来的，那么一定是nil，所以只需要针对nil进行处理
		if okex.Conn == nil {
			logger.Info("enter CloseNormalClosure reConnect")

			for {
				is_connected := okex.connected()
				if is_connected { // 连接成功
					// okex.Is_connected = true // 这里在connected函数里面有过赋值
					break
				} else { // 连接失败
					time.Sleep(time.Duration(1) * time.Second)
					continue
				}
			}
		}

		// 这里原来有一个重连的逻辑，我和前面的逻辑合并在一起了
		// 这里是每次都需要进行的订阅行为，如果订阅失败，那么可以重新订阅一次，如果超过三次订阅失败那么终端程序，发送飞书
		var sub_count int = 0
		for {
			err := okex.subscribe(p.Okex.Subscribe_symbols)
			if err != nil {
				logger.Info("enter ReadMessage reConnect")
				time.Sleep(time.Duration(1) * time.Second)
			} else {
				break
			}

			sub_count += 1
			if sub_count == 3 {
				logger.Fatal("subscribe failed.")
				go utils.FeiShuMsg("DataClient: subscribe failed.")
			}
		}

		for {
			// 这个循环一共有三个break语句，break的时候okex.Conn一定是nil，
			var err error
			var message_okex []byte

			if okex.Conn != nil {
				_, message_okex, err = okex.Conn.ReadMessage()
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) || err != nil {
					okex.Is_connected = false
					okex.Conn.Close()
					okex.Conn = nil
					break
				}
			} else {
				okex.Is_connected = false
				break
			}

			_now := time.Now()
			ct := _now.UnixNano()
			cd := _now.Day()

			var msg_data messages.MsgDataToProcess = messages.MsgDataToProcess{
				ExID:     1,
				Contents: utils.Strcat("1|", utils.ConvertMillisecondString(ct/1e6), "|", string(message_okex)),
			}

			if strings.Contains(msg_data.Contents, "data") {
				chanMsg <- msg_data
			}

			// 这里是换天的语句
			if cd != okex.CurrentDay {
				okex.CurrentDay = cd
				logger.Info("currency break")

				okex.Is_connected = false
				okex.Conn.Close()
				okex.Conn = nil
				break
			}
		}
		time.Sleep(time.Duration(7) * time.Second)
	}
}

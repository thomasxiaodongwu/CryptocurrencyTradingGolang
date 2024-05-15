/*
 * @Author: xwu
 * @Date: 2022-05-21 13:00:56
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 14:31:40
 */
package clientData

import (
	"winter/Logger"

	"go.uber.org/zap"
)

// const (
// 	NANO_PER_SEC = 1e9
// 	// Huobi Constant
// 	// U Future
// 	// public information
// 	WEBSOCKET_HUOBI_URI = `wss://api.hbdm.com/linear-swap-ws`
// 	RESTFUL_BASE_URI    = `api.hbdm.com`

// 	// const RESTFUL_BASE_URI = `api.btcgateway.pro`
// 	RESTFUL_U_ORDER       = `/linear-swap-api/v1/swap_order`
// 	RESTFUL_BATCH_U_ORDER = `/linear-swap-api/v1/swap_batchorder`

// 	// account information huobi
// 	UID        = 29260256
// 	ACCESS_KEY = `3d2xc4v5bu-459ce989-cae3219e-01895`
// 	SECRET_KEY = `7cfd35f5-72f53f09-2b8a3459-6fc45`

// 	// logger setting
// 	IS_LOG_PONG bool = false

// 	// Binance Constant
// 	WEBSOCKET_BINANCE_URI = `wss://fstream.binance.com/ws`
// 	RESTFUL_BINANCE_URI   = `https://fapi.binance.com`
// 	RESTFUL_BINANCE_TRADE = `/fapi/v1/order`

// 	// Okex Constant
// 	RESTFUL_OKEX_URI           = `https://www.okex.com/`
// 	WEBSOCKET_OKEX_URI         = `wss://ws.okex.com:8443/ws/v5/public`
// 	WEBSOCKET_OKEX_PRIVATE_URI = `wss://ws.okex.com:8443/ws/v5/private`
// )

var logger *zap.Logger

func init() {
	logger = Logger.BaseLogger
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:46:50
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-22 20:02:00
 */
package clientTrader

import (
	"winter/Logger"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

const (
	okx_swap_maker_fee float64 = 0.0002
	okx_swap_taker_fee float64 = 0.0005
	okx_slippage       float64 = 0.0005
)

var logger *zap.Logger
var jsonIterator jsoniter.API

func init() {
	logger = Logger.BaseLogger
	jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary // 替代原有的json.UnMarshal
}

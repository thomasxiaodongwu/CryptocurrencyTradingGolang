/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:02
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-31 10:22:08
 */
package clientProcess

import (
	"winter/Logger"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

const (
	Version string = "v0.2"
)

var (
	logger *zap.Logger

	jsonIterator jsoniter.API

	// Kind_Decimals = map[string]float64{
	// 	"BTC-USD":  1,
	// 	"ETH-USD":  10,
	// 	"CRV-USD":  1000,
	// 	"DOGE-USD": 10000,
	// }

	// Okex_Kind_Decimals map[string]([2]int)
	// Okex_Kind_Decimals = map[string]([2]int){
	// 	"BTC-USDT-SWAP":    {1, 0},
	// 	"MKR-USDT-SWAP":    {1, 0},
	// 	"ETH-USDT-SWAP":    {2, 0},
	// 	"CRV-USDT-SWAP":    {4, 0},
	// 	"AXS-USDT-SWAP":    {2, 0},
	// 	"LINK-USDT-SWAP":   {3, 0},
	// 	"SAND-USDT-SWAP":   {4, 0},
	// 	"FTM-USDT-SWAP":    {4, 0},
	// 	"DOGE-USDT-SWAP":   {5, 0},
	// 	"SOL-USDT-SWAP":    {2, 0},
	// 	"CELO-USDT-SWAP":   {4, 0},
	// 	"MATIC-USDT-SWAP":  {4, 0},
	// 	"LUNA-USDT-SWAP":   {4, 0},
	// 	"WAVES-USDT-SWAP":  {3, 0},
	// 	"OP-USDT-SWAP":     {4, 0},
	// 	"GMT-USDT-SWAP":    {4, 0},
	// 	"ADA-USDT-SWAP":    {5, 0},
	// 	"PEOPLE-USDT-SWAP": {6, 0},
	// 	"TRX-USDT-SWAP":    {5, 0},
	// 	"DOT-USDT-SWAP":    {3, 0},
	// 	"DYDX-USDT-SWAP":   {3, 0},
	// 	"APE-USDT-SWAP":    {4, 0},
	// 	"ENS-USDT-SWAP":    {3, 0},
	// 	"LOOKS-USDT-SWAP":  {4, 0},
	// 	"GALA-USDT-SWAP":   {5, 0},
	// 	"LTC-USDT-SWAP":    {2, 0},
	// 	"AVAX-USDT-SWAP":   {2, 0},
	// }

	FTX_Kind_Decimals = map[string]([2]int){
		"ETH-PERP": {1, 3},
		"CRV-PERP": {5, 0},
	}
)

func init() {
	logger = Logger.BaseLogger
	jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:05
 * @Last Modified by: xwu
 * @Last Modified time: 2022-09-27 23:17:50
 */
package params

type AggParams struct {
	Data   ClientDataParams
	Alpha  ClientAlphaParams
	Trader ClientTradeParams
}

type Instrument struct {
	Symbol string
	Size   float64
	Tick   int64
	II     int
}

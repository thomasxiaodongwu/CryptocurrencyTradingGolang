/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:11
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 14:13:42
 */
package params

type ClientDataParams struct {
	FTX     DataParams
	Okex    DataParams
	Binance DataParams
	Dydx    DataParams
	Count   int
}

type DataParams struct {
	Is_disable bool

	Path       string
	Start_time string
	End_time   string

	Subscribe_symbols []string
}

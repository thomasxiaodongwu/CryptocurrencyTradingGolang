/*
 * @Author: xwu
 * @Date: 2021-12-26 18:43:08
 * @Last Modified by: xwu
 * @Last Modified time: 2022-01-22 15:04:26
 */
package params

type ClientAlphaParams struct {
	FTX     AlphaParams
	Okex    AlphaParams
	Binance AlphaParams
}

type AlphaParams struct {
	P []AlphaParamsBase
}

type AlphaParamsBase struct {
	Name     string
	Interval int64
}

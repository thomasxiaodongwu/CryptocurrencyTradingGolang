/*
 * @Author: xwu
 * @Date: 2022-01-22 17:24:00
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 14:16:51
 */
package clientData

func NewFTXClient() FTXClient {
	ans := FTXClient{}
	ans.Is_connected = false
	return ans
}

func NewOkexClient() OkexClient {
	ans := OkexClient{}
	ans.Is_connected = false
	return ans
}

func NewBinanceClient() BinanceClient {
	ans := BinanceClient{pong: 0}
	ans.Is_connected = false
	return ans
}

func NewOkexSim() OkexSim {
	ans := OkexSim{}
	return ans
}

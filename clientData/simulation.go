/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:39
 * @Last Modified by: xwu
 * @Last Modified time: 2022-08-12 17:37:59
 */
package clientData

import (
	"bufio"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"winter/global"
	"winter/messages"

	"winter/utils"
)

// /////////////////////////////////////////////////////////////////
type FTXSim struct {
	Is_connected bool
}

func (FTX *FTXSim) connected(msg string) bool {
	var res bool = false
	if strings.Contains(msg, "connected") {
		res = true
	}
	return res
}

func (FTX *FTXSim) subscribe(line string, symbols []string) bool {
	var res bool = false

	for _, symbol := range symbols {
		res = res || strings.Contains(line, symbol)
		if res {
			break
		}
	}

	return res
}

func (FTX *FTXSim) Run(channel_msg chan *messages.MsgDataToProcess, wg *sync.WaitGroup) {
	p := global.AggParameters.Data
	defer wg.Done()

	// 根据start和end获得一个时间的列表
	var date_list []string = utils.GetDataList(p.FTX.Start_time, p.FTX.End_time)

	// 遍历文件
	FTX.Is_connected = false
	for _, sim_data_name := range date_list {
		var file_path string = path.Join(p.FTX.Path, strings.Split(sim_data_name, "_")[0], sim_data_name+".log")
		logger.Info(file_path + " starts.")

		file, err := os.Open(file_path)
		if err != nil {
			logger.Info("open " + file_path + " failed.")
			continue
		}

		rd := bufio.NewReader(file)

		for {
			line, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}

			if FTX.connected(line) {
				FTX.Is_connected = true
			}

			if FTX.Is_connected {
				if FTX.subscribe(line, p.FTX.Subscribe_symbols) {
					channel_msg <- &messages.MsgDataToProcess{ExID: 0, Contents: line}
				}
			}
		}

		file.Close()
	}
	channel_msg <- &messages.MsgDataToProcess{ExID: -1, Contents: "all files are done."}
}

// ////////////////////////////////////////////////////////
type OkexSim struct {
	Is_connected bool

	Symbols_ob_first_saw map[string]bool
}

func (okex *OkexSim) connected(msg string) bool {
	var res bool = false
	if strings.Contains(msg, "snapshot") {
		res = true
	}
	return res
}

func (okex *OkexSim) subscribe(line string, symbols []string) bool {
	var res bool = false

	for _, symbol := range symbols {
		res = res || strings.Contains(line, symbol)
		if res {
			if strings.Contains(line, "-") {
				if strings.Contains(line, "snapshot") {
					okex.Symbols_ob_first_saw[symbol] = true
				}
			} else {
				okex.Symbols_ob_first_saw[symbol] = true
			}
			res = res && okex.Symbols_ob_first_saw[symbol]
			break
		}
	}

	return res
}

func (okex *OkexSim) Run(channel_msg chan messages.MsgDataToProcess, wg *sync.WaitGroup) {
	p := global.AggParameters.Data
	var line_numbers int64 = 0
	defer wg.Done()

	// 根据start和end获得一个时间的列表
	var date_list []string = utils.GetDataList(p.Okex.Start_time, p.Okex.End_time)

	// 这里是okx交易所的初始化代码，后面append binance交易所的元素
	subscribe_symbols := p.Okex.Subscribe_symbols
	for i, v := range p.Okex.Subscribe_symbols {
		subscribe_symbols[i] = strings.Split(v, "_")[0]
	}
	for _, v := range p.Binance.Subscribe_symbols {
		subscribe_symbols = append(subscribe_symbols, strings.ToUpper(strings.Split(v, "_")[0]))
	}

	// 初始化订阅的map
	okex.Symbols_ob_first_saw = make(map[string]bool)
	for _, symbol := range subscribe_symbols {
		okex.Symbols_ob_first_saw[symbol] = false
	}

	// fmt.Println(subscribe_symbols)
	// 遍历文件
	okex.Is_connected = false
	for i, sim_data_name := range date_list {
		utils.Tqdm(i, len(date_list))
		var file_path string = path.Join(p.Okex.Path, strings.Split(sim_data_name, "_")[0], sim_data_name+".log")

		file, err := os.Open(file_path)
		if err != nil {
			logger.Info("open " + file_path + " failed.")
			continue
		}

		rd := bufio.NewReader(file)

		// V1
		// for {
		// 	line, err := rd.ReadString('\n')
		// 	if err != nil || io.EOF == err {
		// 		break
		// 	}

		// 	if okex.connected(line) {
		// 		okex.Is_connected = true
		// 	}

		// 	if okex.Is_connected {
		// 		if okex.subscribe(line, subscribe_symbols) {
		// 			line_numbers += 1
		// 			channel_msg <- messages.MsgDataToProcess{ExID: 1, Contents: line}
		// 		}
		// 	}
		// }

		// V2
		for {
			line, err := rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}

			if okex.subscribe(line, subscribe_symbols) {
				line_numbers += 1
				// // 这里需要处理一下新旧两种数据
				// // 旧数据只有okex，是十三位时间戳+结构体转化的字符串
				// // 新数据有多个交易所，是交易所代码+十三位时间戳+结构体转化的字符串
				// //
				switch strings.Index(line, "|") {
				case 1: // 这里代表是新数据
					exid, err := strconv.ParseInt(line[:1], 10, 64)
					if err != nil {
						logger.Fatal("error exid")
					}

					channel_msg <- messages.MsgDataToProcess{ExID: int(exid), Contents: line}
				case 13: // 这里代表是旧数据
					channel_msg <- messages.MsgDataToProcess{ExID: 1, Contents: line}
				default:
					logger.Fatal("illegal data format.")
				}

				// channel_msg <- messages.MsgDataToProcess{ExID: 1, Contents: line}
			}

		}

		file.Close()
	}
	channel_msg <- messages.MsgDataToProcess{ExID: -1, Contents: "all files are done."}
	global.LINE_NUM = line_numbers
}

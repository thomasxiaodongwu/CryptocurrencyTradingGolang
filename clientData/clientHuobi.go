/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:27
 * @Last Modified by:   xwu
 * @Last Modified time: 2021-12-26 18:47:27
 */
package clientData

// // Huobi 解压火币返回的信息
// func GZipDecompress(input []byte) (string, error) {
// 	buf := bytes.NewBuffer(input)
// 	reader, gzipErr := gzip.NewReader(buf)
// 	if gzipErr != nil {
// 		return "", gzipErr
// 	}
// 	defer reader.Close()

// 	result, readErr := ioutil.ReadAll(reader)
// 	if readErr != nil {
// 		return "", readErr
// 	}

// 	return string(result), nil
// }

// // Huobi 接受数据的客户端
// func client_huobi(subSymbols []string, data_channel chan *message) {
// 	var count int = 0
// 	var subID int = 0

// 	for {
// 		// 新建客户端
// 		conn, _, err := websocket.DefaultDialer.Dial(WEBSOCKET_HUOBI_URI, nil)
// 		if err != nil {
// 			logger.Info("Failed connect the server")
// 			time.Sleep(7 * time.Second)
// 			continue
// 		} else {
// 			logger.Info(fmt.Sprintf("Successful connection for %dth time", count+1))
// 		}

// 		// send subscribe message to server
// 		for _, info := range subSymbols {
// 			var sub, symbol, msgkind string
// 			symbol = symbol_format(strings.Split(info, "|")[0], "huobi")
// 			msgkind = strings.Split(info, "|")[1]
// 			sub = mkSubMsg(symbol, msgkind, "huobi", subID)

// 			subID += 1

// 			err := conn.WriteMessage(websocket.TextMessage, []byte(sub))
// 			if err != nil {
// 				logger.Info(fmt.Sprintf("Huobi Subscribe %s Failed.", info))
// 			} else {
// 				logger.Info(fmt.Sprintf("Huobi Subscribe %s Success.", info))
// 			}
// 			time.Sleep(100 * time.Millisecond)
// 		}

// 		var reg *regexp.Regexp = regexp.MustCompile(`[0-9]{13}`)
// 		for {
// 			var messages []byte
// 			var result string

// 			if conn != nil {
// 				_, messages, err = conn.ReadMessage()
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
// 					logger.Info(fmt.Sprintf("Reconnect Starts. %d", count+1))
// 					break
// 				}
// 			} else {
// 				break
// 			}

// 			result, err = GZipDecompress(messages)
// 			if err != nil {
// 				logger.Info("Decompress Failed. " + err.Error())
// 				if err.Error() == "EOF" {
// 					time.Sleep(7 * time.Second)
// 				}
// 				break
// 			}

// 			if len(result) > 0 {
// 				if result[2:6] != "ping" {
// 					var msg_data message
// 					msg_data.ExID = 2
// 					msg_data.LocalTime = time.Now().UnixNano()
// 					msg_data.Message = result
// 					if isMessageUseful(msg_data.Message, "huobi") {
// 						data_channel <- &msg_data
// 					}

// 				} else if result[2:6] == "ping" {
// 					var ts string = reg.FindString(result)
// 					pongData := fmt.Sprintf("{\"pong\": %s}", ts)

// 					err = conn.WriteMessage(websocket.TextMessage, []byte(pongData))
// 					if err != nil {
// 						logger.Info(fmt.Sprintf("Send Pong Messages Failed. Error Code is %s.", err.Error()))
// 					} else {
// 						if IS_LOG_PONG {
// 							logger.Info("Send Pong Messages Success.")
// 						}
// 					}
// 				}
// 			}
// 		}
// 		conn.Close()
// 		conn = nil
// 		time.Sleep(7 * time.Second)
// 		count += 1
// 	}
// }

// // Huobi 接受数据的客户端
// func client_huobi(subSymbols []string, data_channel chan *message) {
// 	var uri string = WEBSOCKET_BASE_URI
// 	var sub, symbol, msgkind string
// 	var count int = 0
// 	var subID int = 0

// 	for {
// 		conn, _, err := websocket.DefaultDialer.Dial(uri, nil)
// 		if err != nil {
// 			logger.Info("Failed connect the server")
// 			time.Sleep(7 * time.Second)
// 			continue
// 		} else {
// 			logger.Info(fmt.Sprintf("Successful connection for %dth time", count+1))
// 		}

// 		// send subscribe message to server
// 		for _, info := range subSymbols {
// 			symbol = symbol_format(strings.Split(info, "|")[0], "huobi")
// 			msgkind = strings.Split(info, "|")[1]
// 			if msgkind == "td" {
// 				sub = fmt.Sprintf("{\"sub\":\"market.%s.trade.detail\",\"id\":\"%d\"}", symbol, subID)
// 			} else if msgkind == "obmin" {
// 				sub = fmt.Sprintf("{\"sub\":\"market.%s.bbo\",\"id\":\"%d\"}", symbol, subID)
// 			} else if msgkind == "ob" {
// 				sub = fmt.Sprintf("{\"sub\":\"market.%s.depth.step6\",\"id\":\"%d\"}", symbol, subID)
// 			}

// 			subID += 1
// 			err := conn.WriteMessage(websocket.TextMessage, []byte(sub))
// 			if err != nil {
// 				logger.Info(fmt.Sprintf("Huobi Subscribe %s Failed.", info))
// 			} else {
// 				logger.Info(fmt.Sprintf("Huobi Subscribe %s Success.", info))
// 			}
// 			time.Sleep(100 * time.Millisecond)
// 		}

// 		var messages []byte
// 		var reg *regexp.Regexp = regexp.MustCompile(`[0-9]{13}`)
// 		for {
// 			if conn != nil {
// 				_, messages, err = conn.ReadMessage()
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
// 					logger.Info(fmt.Sprintf("Reconnect Starts. %d", count))
// 					break
// 				}
// 			} else {
// 				break
// 			}

// 			// logger.Info(fmt.Sprintf("Huobi %d", msg_data.LocalTime))

// 			var result string
// 			result, err = GZipDecompress(messages)
// 			if err != nil {
// 				logger.Info("Decompress Failed. " + err.Error())
// 				if err.Error() == "EOF" {
// 					time.Sleep(7 * time.Second)
// 				}
// 				continue
// 			}

// 			if len(result) > 0 {
// 				if result[2:6] != "ping" {
// 					var msg_data message
// 					msg_data.ExID = 2
// 					msg_data.LocalTime = time.Now().UnixNano()
// 					msg_data.Message = result

// 					data_channel <- &msg_data
// 				} else if result[2:6] == "ping" {
// 					var ts string = reg.FindString(result)
// 					pongData := fmt.Sprintf("{\"pong\": %s}", ts)

// 					err = conn.WriteMessage(websocket.TextMessage, []byte(pongData))
// 					if err != nil {
// 						logger.Info(fmt.Sprintf("Send Pong Messages Failed. Error Code is %s.", err.Error()))
// 					} else {
// 						if IS_LOG_PONG {
// 							logger.Info("Send Pong Messages Success.")
// 						}
// 					}
// 				}
// 			}
// 		}
// 		conn.Close()
// 		conn = nil
// 		time.Sleep(7 * time.Second)
// 		count += 1
// 	}
// }

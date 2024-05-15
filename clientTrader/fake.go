/*
 * @Author: xwu
 * @Date: 2022-05-22 20:05:35
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-22 21:50:55
 */
package clientTrader

import (
	"sync"
	"winter/messages"
)

func Fake(channel chan messages.AggSignal, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		signal := <-channel
		if signal.Status == 0 {
			break
		}

		// if signal.Localtime == 0 {
		// 	fmt.Println(signal.Localtime)
		// }
	}
}

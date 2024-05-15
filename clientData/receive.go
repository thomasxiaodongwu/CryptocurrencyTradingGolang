/*
 * @Author: xwu
 * @Date: 2021-12-26 18:47:36
 * @Last Modified by: xwu
 * @Last Modified time: 2022-06-02 09:38:29
 */
package clientData

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"winter/global"
	"winter/messages"
	"winter/utils"
)

func Fake(channel chan messages.MsgDataToProcess, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		msg_data := <-channel
		fmt.Println(msg_data.Contents)
	}
}

func DumpRawData(channel chan messages.MsgDataToProcess, wg *sync.WaitGroup) {
	defer wg.Done()

	var t time.Time = time.Now()
	var File_name string = t.Format("20060102_15")
	var Hour int = t.Hour()

	dump_path := global.Data_Dump_Path

	// 建立文件夹
	if !utils.Exists(dump_path + t.Format("20060102")) {
		err := os.Mkdir(dump_path+t.Format("20060102"), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	fileHandle, err := os.OpenFile(dump_path+t.Format("20060102")+"/"+File_name+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	buf := bufio.NewWriter(fileHandle)
	for {
		msg_data := <-channel
		t = time.Now()

		if Hour != t.Hour() {
			Hour = t.Hour()

			fileHandle.Close()
			File_name = t.Format("20060102_15")
			// 新的一天的文件夹，判断是否存在，如果不存在，则建立文件夹
			if !utils.Exists(dump_path + t.Format("20060102")) {
				err := os.Mkdir(dump_path+t.Format("20060102"), os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}

			// 更换filehandler
			fileHandle, err = os.OpenFile(dump_path+t.Format("20060102")+"/"+File_name+".log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			if err != nil {
				logger.Info("open file error :" + err.Error())
			}
			buf = bufio.NewWriter(fileHandle)

			logger.Info(dump_path + t.Format("20060102") + "/" + File_name + ".log" + " starts.")
		}

		// 字符串写入
		buf.WriteString(msg_data.Contents + "\n")
		// 将缓冲中的数据写入
		err = buf.Flush()
		if err != nil {
			log.Println("flush error :", err)
		}
	}
}

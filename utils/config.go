/*
 * @Author: xwu
 * @Date: 2021-12-26 18:41:54
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 13:18:22
 */
package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"winter/params"
)

//读取到file中，再利用ioutil将file直接读取到[]byte中, 这是最优
func ReadFronJson(file_name string) string {
	f, err := os.Open(file_name)
	if err != nil {
		logger.Info("read file fail. " + err.Error())
		return ""
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		logger.Info("read to fd fail. " + err.Error())
		return ""
	}

	return string(fd)
}

func LoadConfig(file_name string) params.AggParams {
	var config params.AggParams
	ans := ReadFronJson(file_name)
	err := json.Unmarshal([]byte(ans), &config)
	if err != nil {
		panic(err.Error())
	}
	return config
}

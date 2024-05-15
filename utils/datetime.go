/*
 * @Author: xwu
 * @Date: 2021-12-26 18:42:07
 * @Last Modified by: xwu
 * @Last Modified time: 2022-03-16 13:39:34
 */
package utils

import (
	"strconv"
	"time"
)

func GetDataList(start_time string, end_time string) []string {
	var date_list []string

	start_t, err := time.ParseInLocation("20060102_15", start_time, time.Local)
	if err != nil {
		logger.Info("get_date_list: start_time string convert to date failed.")
	}

	end_t, err := time.ParseInLocation("20060102_15", end_time, time.Local)
	if err != nil {
		logger.Info("get_date_list: end_time string convert to date failed.")
	}

	for start_t.Before(end_t) || start_t.Equal(end_t) {
		date_list = append(date_list, start_t.Format("20060102_15"))
		start_t = start_t.Add(time.Hour)
	}

	return date_list
}

func GetMillisecondString() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1e6))
}

func ConvertMillisecondString(t int64) string {
	return strconv.Itoa(int(t))
}

func UnixToTime(e int64) time.Time {
	datatime := time.Unix(e/1000, 0)
	return datatime
}

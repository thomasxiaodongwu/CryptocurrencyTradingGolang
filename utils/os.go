/*
 * @Author: xwu
 * @Date: 2021-12-26 18:42:16
 * @Last Modified by:   xwu
 * @Last Modified time: 2021-12-26 18:42:16
 */
package utils

import "os"

// 判断文件夹是否是存在的
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		// if os.IsExist(err) {
		// 	return true
		// }
		// return false
		return os.IsExist(err)
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

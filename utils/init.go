/*
 * @Author: xwu
 * @Date: 2021-12-26 18:42:13
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 14:31:27
 */
package utils

import (
	"winter/Logger"

	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
)

var (
	logger         *zap.Logger
	feishu_webhook string
	jsonIterator   jsoniter.API
)

func init() {
	logger = Logger.BaseLogger
	feishu_webhook = `https://open.feishu.cn/open-apis/bot/v2/hook/83da836d-8ff6-4de0-9806-8ca29db40300`
	jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:46:32
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-14 16:07:41
 */
package container

import (
	"winter/Logger"

	"go.uber.org/zap"
)

const DEFAULT_SKIP_LIST_LEVEL int64 = 6 // default level

var logger *zap.Logger

func init() {
	logger = Logger.BaseLogger
}

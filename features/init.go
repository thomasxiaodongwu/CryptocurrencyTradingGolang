/*
 * @Author: xwu
 * @Date: 2022-10-09 18:38:21
 * @Last Modified by: xwu
 * @Last Modified time: 2022-10-09 18:38:45
 */

package features

import (
	"winter/Logger"

	"go.uber.org/zap"
)

const Version string = "v0.1"

var logger *zap.Logger

func init() {
	logger = Logger.BaseLogger
}

/*
 * @Author: xwu
 * @Date: 2021-12-26 18:45:30
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 14:33:29
 */
package mathpro

import (
	"winter/Logger"

	"go.uber.org/zap"
)

const Version string = "v0.2"

var logger *zap.Logger

func init() {
	logger = Logger.BaseLogger
}

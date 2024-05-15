/*
 * @Author: xwu
 * @Date: 2021-12-26 18:44:41
 * @Last Modified by: xwu
 * @Last Modified time: 2022-05-21 16:02:52
 */
package operator

import (
	"winter/Logger"

	"go.uber.org/zap"
)

const Version string = "v1.0"

var logger *zap.Logger

func init() {
	logger = Logger.BaseLogger
}

type valueTm struct {
	Tm    int64
	Value float64
}

type valuesTm struct {
	Tm     int64
	Value1 float64
	Value2 float64
}

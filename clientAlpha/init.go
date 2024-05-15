package clientAlpha

import (
	"winter/Logger"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger = Logger.BaseLogger
}

package messages

import (
	"winter/Logger"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	global_sub_id int64
)

func init() {
	logger = Logger.BaseLogger
	global_sub_id = 0
}

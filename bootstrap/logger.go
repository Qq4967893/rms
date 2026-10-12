package bootstrap

import (
	"backend-rms/core/config"
	"backend-rms/core/logger"
)

func InitLogger() {
	logger.Init(config.Conf.Log)
}

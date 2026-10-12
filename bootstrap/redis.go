package bootstrap

import (
	"backend-rms/core/config"
	"backend-rms/core/redis"
)

func InitRedis() {
	redis.InitBaseRedis(config.Conf.Redis.Base)
	redis.InitRMSRedis(config.Conf.Redis.RMS)
	redis.InitQuoteRedis(config.Conf.Redis.Quote)
	redis.InitPositionRedis(config.Conf.Redis.Position)
}

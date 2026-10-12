package bootstrap

import (
	"backend-rms/core/config"
	"backend-rms/core/jwt"
)

func InitJWT() {
	jwt.InitJWT(config.Conf.JWT)
}

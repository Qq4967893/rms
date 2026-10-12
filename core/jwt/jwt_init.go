package jwt

import "backend-rms/core/config"

var Instance *JWT

func InitJWT(conf config.JWTConf) {
	Instance = NewJWT(conf)
}

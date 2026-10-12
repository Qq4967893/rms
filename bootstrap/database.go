package bootstrap

import (
	"backend-rms/core/config"
	"backend-rms/core/database"
)

func InitDB() {
	database.InitAPIDB(config.Conf.MySQL.API)
	database.InitRMSDB(config.Conf.MySQL.RMS)
	database.InitCRMDB(config.Conf.MySQL.CRM)
	database.InitET6DB(config.Conf.MySQL.ET6)
	database.InitCHARTDB(config.Conf.MySQL.CHART)
	database.InitETCDB(config.Conf.MySQL.ETC)
	database.InitET6ConfigDB(config.Conf.MySQL.ET6Config)
}

func InitMongoDB() {
	database.InitTickMongoDB(config.Conf.Mongo.Tick)
}

package config

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	Log    LogConf `mapstructure:"log"`
	Sentry SentryConf
	JWT    JWTConf
	MySQL  MySQLList
	Redis  RedisList
	Mongo  MongoList
	GRPC   GRPCList
	Email  EmailConf
	Kafka  KafkaConf
}

type LogConf struct {
	Level        string
	Stdout       bool
	FileName     string
	AccessFolder string
	TimeFormat   string
	Compress     bool
	MaxSize      int
	MaxBackup    int
	MaxAge       int
}

type SentryConf struct {
	Dsn              string
	EnableTracing    bool
	TracesSampleRate float64
}

type MySQLList struct {
	API       MySQLConf
	RMS       MySQLConf
	CRM       MySQLConf
	ET6       MySQLConf
	CHART     MySQLConf
	ETC       MySQLConf
	ET6Config MySQLConf
}

type MySQLConf struct {
	Username     string
	Password     string
	Host         string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
}

type RedisList struct {
	Base     RedisConf
	RMS      RedisConf
	Quote    RedisConf
	Position RedisConf
}

type RedisConf struct {
	Host     string
	Password string
	DB       int
}

type MongoList struct {
	Tick MongoConf
}

type MongoConf struct {
	Host        string
	Username    string
	Password    string
	Database    string
	MaxPoolSize uint64
	MinPoolSize uint64
}

type JWTConf struct {
	SecretKey      string
	ExpireTime     int
	MaxRefreshTime int
}

type GRPCList struct {
	TSP GRPCConf
}

type GRPCConf struct {
	Address string
	Timeout time.Duration
}

type EmailConf struct {
	EmailDriver    string
	ApiKey         string
	ApiUser        string
	ApiUserBatch   string
	From           string
	FromBatch      string
	FromName       string
	DeveloperEmail string
}

type KafkaConf struct {
	Addr     string
	Topic    Topics
	Consumer string
}

type Topics struct {
	Order string
	Group string
}

func Init() {
	v := viper.New()
	v.SetConfigFile("configs/config.yaml")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error occured while reading config file: %s", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file has changed: ", e.Name)
	})

	if err := v.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}

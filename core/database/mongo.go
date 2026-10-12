package database

import (
	"backend-rms/core/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var TickMongoDB *mongo.Database

func InitTickMongoDB(conf config.MongoConf) {
	TickMongoDB = initMongoDB(conf)
}

func initMongoDB(conf config.MongoConf) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s",
		conf.Username, conf.Password, conf.Host, conf.Database)

	clientOpts := options.Client().ApplyURI(uri).
		SetMaxPoolSize(conf.MaxPoolSize).
		SetMinPoolSize(conf.MinPoolSize)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	return client.Database(conf.Database)
}

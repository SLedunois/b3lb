package app

import (
	"b3lb/pkg/config"

	"github.com/go-redis/redis/v8"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func redisClient(conf *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.RDB.Address,
		Password: conf.RDB.Password,
		DB:       conf.RDB.DB,
	})
}

func influxDBClient(conf *config.Config) api.QueryAPI {
	client := influxdb2.NewClient(conf.IDB.Address, conf.IDB.Token)
	return client.QueryAPI(conf.IDB.Organization)
}

package goredis

import (
	"bitbucket.org/roadrunnr/go_common_lib/config"
	"bitbucket.org/roadrunnr/go_common_lib/logger"
	"time"

	"github.com/Scalingo/go-workers"
	"gopkg.in/redis.v4"
)

//redis global client to be declared
var redisClient *redis.Client

type IClient interface {
	GetValue(key string) (string, error)
	SetValue(key string, value string) error
	SetValueEx(key string, value string, seconds int) error
	LPush(key string, value string) error
}

type Client struct {
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func GetRedisClient() *redis.Client{
	return redisClient
}

// FetchRedisConfig - reading redis config from redis.yml
func FetchRedisConfig() *RedisConfig {
	v1 := myconfig.SetupViperAndReadConfig("redis")

	var redisConfig RedisConfig
	config := v1.GetStringMapString(myconfig.GetEnv())
	redisConfig.Host = config["host"]
	redisConfig.Port = config["port"]
	redisConfig.Password = config["password"]
	return &redisConfig
}

// Init - initializes the redisClient
func Init() {
	if redisClient == nil {
		redisConfig := FetchRedisConfig()
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisConfig.Host + ":" + redisConfig.Port,
			Password: redisConfig.Password, // no password set
			DB:       0,                    // use default DB
		})
	}
}

// ConfigureSidekiq - configures the sidekiq queue.
func ConfigureSidekiq() {
	redisConfig := FetchRedisConfig()
	workers.Configure(map[string]string{"process": "pingclient", "password": redisConfig.Password,
		"database": "12", "server": redisConfig.Host + ":" + redisConfig.Port})
}

// EnqueSidekiqJob - enques the sidekiq job
func EnqueSidekiqJob(queue string, worker string, params []string) {
	workers.Enqueue(queue, worker, params)
}

// GetValue - get data from redis
func (client Client) GetValue(key string) (string, error) {
	value, err := redisClient.Get(key).Result()
	if err == redis.Nil {
		//log.Error("Redis key does not exist: " + key)
		return "", err
	} else if err != nil {
		log.Error("Redis read key error: " + err.Error())
		return "", err
	}
	return value, err
}

func (client Client) SetValue(key string, value string) error {
	err := redisClient.Set(key, value, 0).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}

func (client Client) SetValueEx(key string, value string, seconds int) error {
	err := redisClient.Set(key, value, time.Duration(seconds) * time.Second).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}


func (client Client) LPush(key string, value string) error {
	err := redisClient.LPush(key, value).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}

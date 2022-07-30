package goredis

import (
	myconfig "bitbucket.org/zatasales/go_common_lib/config"
	log "bitbucket.org/zatasales/go_common_lib/logger"
	"time"
	"gopkg.in/redis.v4"
	"strings"
)

//redis global client to be declared
//var redisClient *redis.Client

type IClient interface {
	Init()
	GetValue(key string) (string, error)
	SetValue(key string, value string) error
	SetValueEx(key string, value string, seconds int) error
	LPush(key string, value string) error
	ZAdd(key string, values []redis.Z) error
	ZRem(key string, values ...interface{}) error
	ZRange(key string, start int64, end int64) []string
	ZCard(key string) int64
}

type Client struct {
	RedisConfig *RedisConfig
	redisClient *redis.Client
}

type ClusterClient struct {
	RedisConfig *RedisConfig
	redisClient *redis.ClusterClient
}

type RedisConfig struct {
	Host     string
	Hosts     string
	Port     string
	Password string
}

func (client *ClusterClient) GetRedisClient() *redis.ClusterClient{
	return client.redisClient
}

func (client *Client) GetRedisClient() *redis.Client{
	return client.redisClient
}

// FetchRedisConfig - reading redis config from redis.yml
func FetchRedisConfig(configFile string) *RedisConfig {
	v1 := myconfig.SetupViperAndReadConfig(configFile)

	var redisConfig RedisConfig
	config := v1.GetStringMapString(myconfig.GetEnv())
	redisConfig.Host = config["host"]
	redisConfig.Hosts = config["hosts"]
	redisConfig.Port = config["port"]
	redisConfig.Password = config["password"]
	return &redisConfig
}


// Init - initializes the redisClient
func (client *ClusterClient) Init() {
	if client.redisClient == nil {
		redisConfig := client.RedisConfig
		addr := strings.Split(redisConfig.Hosts, ",")
		client.redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     addr,
		})
	}
}

// Init - initializes the redisClient
func (client *Client) Init() {
	if client.redisClient == nil {
		redisConfig := client.RedisConfig
		client.redisClient = redis.NewClient(&redis.Options{
			Addr:     redisConfig.Host + ":" + redisConfig.Port,
			Password: redisConfig.Password, // no password set
			DB:       0,                    // use default DB
		})
	}
}

// GetValue - get data from redis
func (client *Client) GetValue(key string) (string, error) {
	value, err := client.GetRedisClient().Get(key).Result()
	if err == redis.Nil {
		//log.Error("Redis key does not exist: " + key)
		return "", err
	} else if err != nil {
		log.Error("Redis read key error: " + err.Error())
		return "", err
	}
	return value, err
}

func (client *Client) SetValue(key string, value string) error {
	err := client.GetRedisClient().Set(key, value, 0).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}

func (client *Client) SetValueEx(key string, value string, seconds int) error {
	err := client.GetRedisClient().Set(key, value, time.Duration(seconds) * time.Second).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}


func (client *Client) LPush(key string, value string) error {
	err := client.GetRedisClient().LPush(key, value).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}


func (client *Client) ZAdd(key string, values []redis.Z) error {
	err := client.GetRedisClient().ZAdd(key, values...).Err()
	if err != nil {
		log.Error("Redis ZAdd key error: " + err.Error())
	}
	return err
}

func (client *Client) ZRem(key string, values ...interface{}) error {
	err := client.GetRedisClient().ZRem(key, values).Err()
	if err != nil {
		log.Error("Redis ZRem key error: " + err.Error())
	}
	return err
}


func (client *Client) ZRange(key string, start int64, end int64) []string {
	val := client.GetRedisClient().ZRange(key, start, end)
	return val.Val()
}

func (client *Client) ZCard(key string) int64 {
	val := client.GetRedisClient().ZCard(key)
	return val.Val()
}


////////////////// cluster functions

// GetValue - get data from redis
func (client *ClusterClient) GetValue(key string) (string, error) {
	value, err := client.GetRedisClient().Get(key).Result()
	if err == redis.Nil {
		//log.Error("Redis key does not exist: " + key)
		return "", err
	} else if err != nil {
		log.Error("Redis read key error: " + err.Error())
		return "", err
	}
	return value, err
}

func (client *ClusterClient) SetValue(key string, value string) error {
	err := client.GetRedisClient().Set(key, value, 0).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}

func (client *ClusterClient) SetValueEx(key string, value string, seconds int) error {
	err := client.GetRedisClient().Set(key, value, time.Duration(seconds) * time.Second).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}


func (client *ClusterClient) LPush(key string, value string) error {
	err := client.GetRedisClient().LPush(key, value).Err()
	if err != nil {
		log.Error("Redis read key error: " + err.Error())
	}
	return err
}


func (client *ClusterClient) ZAdd(key string, values []redis.Z) error {
	err := client.GetRedisClient().ZAdd(key, values...).Err()
	if err != nil {
		log.Error("Redis ZAdd key error: " + err.Error())
	}
	return err
}

func (client *ClusterClient) ZRem(key string, values ...interface{}) error {
	err := client.GetRedisClient().ZRem(key, values).Err()
	if err != nil {
		log.Error("Redis ZRem key error: " + err.Error())
	}
	return err
}


func (client *ClusterClient) ZRange(key string, start int64, end int64) []string {
	val := client.GetRedisClient().ZRange(key, start, end)
	return val.Val()
}

func (client *ClusterClient) ZCard(key string) int64 {
	val := client.GetRedisClient().ZCard(key)
	return val.Val()
}

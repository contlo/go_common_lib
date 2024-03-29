package goredis

import (
    "context"
    "fmt"
    "strconv"
    "strings"
    "time"

    myconfig "github.com/contlo/go_common_lib/config"
    log "github.com/contlo/go_common_lib/logger"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

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
    ZRangeByScore(key string, min string, max string, offset int64, count int64) []string
    ZCard(key string) int64
    SMembers(key string) []string
    Expire(key string, expire time.Duration) bool
    Lock(key string, expire time.Duration) bool
    GetRedisClient() *redis.Client
    SAdd(key string, values ...string) error
    SMIsMemberValues(key string, values ...string) ([]string, error)
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
    Hosts    string
    Port     string
    Password string
    DB       string
}

func (client *ClusterClient) GetRedisClient() *redis.ClusterClient {
    return client.redisClient
}

func (client *Client) GetRedisClient() *redis.Client {
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
    if dbVal, ok := config["db"]; ok {
        redisConfig.DB = dbVal
    } else {
        redisConfig.DB = "0"
    }

    return &redisConfig
}

// Init - initializes the redisClient
func (client *ClusterClient) Init() {
    if client.redisClient == nil {
        redisConfig := client.RedisConfig
        addr := strings.Split(redisConfig.Hosts, ",")
        client.redisClient = redis.NewClusterClient(&redis.ClusterOptions{
            Addrs: addr,
        })
    }
}

// Init - initializes the redisClient
func (client *Client) Init() {
    if client.redisClient == nil {
        redisConfig := client.RedisConfig
        db, _ := strconv.Atoi(redisConfig.DB)
        client.redisClient = redis.NewClient(&redis.Options{
            Addr:     redisConfig.Host + ":" + redisConfig.Port,
            Password: redisConfig.Password, // no password set
            DB:       db,                   // use default DB
        })
    }
}

// GetValue - get data from redis
func (client *Client) GetValue(key string) (string, error) {
    value, err := client.GetRedisClient().Get(ctx, key).Result()
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
    err := client.GetRedisClient().Set(ctx, key, value, 0).Err()
    if err != nil {
        log.Error("Redis read key error: " + err.Error())
    }
    return err
}

func (client *Client) SetValueEx(key string, value string, seconds int) error {
    err := client.GetRedisClient().Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
    if err != nil {
        log.Error("Redis read key error: " + err.Error())
    }
    return err
}

func (client *Client) LPush(key string, value string) error {
    err := client.GetRedisClient().LPush(ctx, key, value).Err()
    if err != nil {
        log.Error("Redis read key error: " + err.Error())
    }
    return err
}

func (client *Client) ZAdd(key string, values []redis.Z) error {
    err := client.GetRedisClient().ZAdd(ctx, key, values...).Err()
    if err != nil {
        fmt.Println("Redis ZAdd key error: " + err.Error())
        log.Error("Redis ZAdd key error: " + err.Error())
    }
    return err
}

func (client *Client) ZRem(key string, values ...interface{}) error {
    err := client.GetRedisClient().ZRem(ctx, key, values).Err()
    if err != nil {
        log.Error("Redis ZRem key error: " + err.Error())
    }
    return err
}

func (client *Client) ZRange(key string, start int64, end int64) []string {
    val := client.GetRedisClient().ZRange(ctx, key, start, end)
    return val.Val()
}

func (client *Client) ZCard(key string) int64 {
    val := client.GetRedisClient().ZCard(ctx, key)
    return val.Val()
}

func (client *Client) SMembers(key string) []string {
    val := client.GetRedisClient().SMembers(ctx, key)
    return val.Val()
}

func (client *Client) SAdd(key string, values ...string) error {
    err := client.GetRedisClient().SAdd(ctx, key, values).Err()
    if err != nil {
        fmt.Println("Redis ZAdd key error: " + err.Error())
        log.Error("Redis ZAdd key error: " + err.Error())
    }
    return err
}

func (client *Client) SMIsMemberValues(key string, values ...string) ([]string, error) {
    valuesOut := []string{}
    results, err := client.GetRedisClient().SMIsMember(context.TODO(), key, values).Result()
    if err != nil {
        log.Error("Redis SMIsMember key error: " + err.Error())
        return nil, err
    }
    for i, value := range values {
        if results[i] {
            valuesOut = append(valuesOut, value)
        }
    }

    return valuesOut, nil
}

func (client *Client) Expire(key string, expire time.Duration) bool {
    val := client.GetRedisClient().Expire(ctx, key, expire*time.Second)
    return val.Val()
}
func (client *Client) ZRangeByScore(key string, min string, max string, offset int64, count int64) []string {
    val := client.GetRedisClient().ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count})
    return val.Val()
}
func (client *Client) Lock(key string, expire time.Duration) bool {
    val := client.GetRedisClient().Incr(ctx, key)
    if val.Val() > 1 {
        return false
    } else {
        client.GetRedisClient().Expire(ctx, key, expire)
    }
    return true
}

////////////////// cluster functions

// GetValue - get data from redis
func (client *ClusterClient) GetValue(key string) (string, error) {
    value, err := client.GetRedisClient().Get(ctx, key).Result()
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
    err := client.GetRedisClient().Set(ctx, key, value, 0).Err()
    if err != nil {
        log.Error("Redis read key error: " + err.Error())
    }
    return err
}

func (client *ClusterClient) SetValueEx(key string, value string, seconds int) error {
    err := client.GetRedisClient().Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
    if err != nil {
        log.Error("Redis read key error: " + err.Error())
    }
    return err
}

func (client *ClusterClient) LPush(key string, value string) error {
    err := client.GetRedisClient().LPush(ctx, key, value).Err()
    if err != nil {
        log.Error("Redis read key error: " + err.Error())
    }
    return err
}

func (client *ClusterClient) ZAdd(key string, values []redis.Z) error {
    err := client.GetRedisClient().ZAdd(ctx, key, values...).Err()
    if err != nil {
        fmt.Println("Redis ZAdd key error: " + err.Error())
        log.Error("Redis ZAdd key error: " + err.Error())
    }
    return err
}

func (client *ClusterClient) ZRem(key string, values ...interface{}) error {
    err := client.GetRedisClient().ZRem(ctx, key, values).Err()
    if err != nil {
        log.Error("Redis ZRem key error: " + err.Error())
    }
    return err
}

func (client *ClusterClient) ZRange(key string, start int64, end int64) []string {
    val := client.GetRedisClient().ZRange(ctx, key, start, end)
    return val.Val()
}

func (client *ClusterClient) ZRangeByScore(key string, min string, max string, offset int64, count int64) []string {
    val := client.GetRedisClient().ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count})
    return val.Val()
}

func (client *ClusterClient) ZCard(key string) int64 {
    val := client.GetRedisClient().ZCard(ctx, key)
    return val.Val()
}

func (client *ClusterClient) SAdd(key string, values ...string) error {
    err := client.GetRedisClient().SAdd(ctx, key, values).Err()
    if err != nil {
        fmt.Println("Redis ZAdd key error: " + err.Error())
        log.Error("Redis ZAdd key error: " + err.Error())
    }
    return err
}

func (client *ClusterClient) SMIsMemberValues(key string, values ...string) ([]string, error) {
    valuesOut := []string{}
    results, err := client.GetRedisClient().SMIsMember(context.TODO(), key, values).Result()
    if err != nil {
        log.Error("Redis SMIsMember key error: " + err.Error())
        return nil, err
    }
    for i, value := range values {
        if results[i] {
            valuesOut = append(valuesOut, value)
        }
    }

    return valuesOut, nil
}

func (client *ClusterClient) SMembers(key string) []string {
    val := client.GetRedisClient().SMembers(ctx, key)
    return val.Val()
}
func (client *ClusterClient) Expire(key string, expire time.Duration) bool {
    val := client.GetRedisClient().Expire(ctx, key, expire*time.Second)
    return val.Val()
}
func (client *ClusterClient) Lock(key string, expire time.Duration) bool {
    val := client.GetRedisClient().Incr(ctx, key)
    if val.Val() > 1 {
        return false
    } else {
        client.GetRedisClient().Expire(ctx, key, expire)
    }
    return true
}

package goredis

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func InitTest() {
	os.Setenv("GO_ENV", "test")
	Init()
}

func TestGetValueNilValue(t *testing.T) {
	value, _ := Client{}.GetValue("abcd")

	assert.Equal(t, "", value, "Redis get value should be nil.")
}

func TestSetAndGetValue(t *testing.T) {
	Client{}.SetValue("redis_set1", "abcd")
	value, _ := Client{}.GetValue("redis_set1")

	assert.Equal(t, "abcd", value, "Redis get value should be abcd.")
}

func TestFetchRedisConfig(t *testing.T) {
	config := FetchRedisConfig()

	assert.Equal(t, config.Host, "localhost", "TestFetchRedisConfig host config mismatch.")
	assert.Equal(t, config.Port, "6379", "TestFetchRedisConfig port config mismatch.")
}

func TestMain(m *testing.M) {
	InitTest()
	os.Exit(m.Run())
}

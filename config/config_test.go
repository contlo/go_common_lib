package myconfig

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func InitTest() {
	os.Setenv("GO_ENV", "test")
}

func TestMetricEnabled(t *testing.T) {
	assert.False(t, MetricEnabled(), "TestMetricEnabled should be false.")
}

func TestSetupViperConfig(t *testing.T) {
	v1 := SetupViperAndReadConfig("redis")
	config := v1.GetStringMapString(GetEnv())

	assert.Equal(t, config["host"], "localhost", "host must be localhost")
}

func TestMain(m *testing.M) {
	InitTest()
	os.Exit(m.Run())
}

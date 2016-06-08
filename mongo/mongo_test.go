package mongo

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func InitTest() {
	os.Setenv("GO_ENV", "test")
	InitMongo()
}

func TestSetupViperAndReadConfig(t *testing.T) {
	config := FetchMongoConfig()

	assert.Equal(t, config.Database, "cart_hero_test", "TestSetupViperAndReadConfig database config mismatch.")
	assert.Equal(t, config.Hosts[0], "localhost:27017", "TestSetupViperAndReadConfig hosts config mismatch.")
}

func TestCreateMongoSession(t *testing.T) {
	MgoSession = nil
	session := CreateMongoSession()

	assert.NotNil(t, session, "TestCreateMongoSession mongo session cant be nil.")
}

func TestMain(m *testing.M) {
	InitTest()
	os.Exit(m.Run())
}

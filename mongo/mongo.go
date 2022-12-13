package mongo

import (
	"context"
	"fmt"
	myconfig "go_common_lib/config"
	log "go_common_lib/logger"
	"time"
	_ "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfigData struct {
	Database string
	Hosts    []string
	Username string
	Password string
	MongoURI string
}

type MongoDB struct {
	MongoConfig *MongoConfigData
	Client      *mongo.Client
	Context     context.Context
	CancelFunc  context.CancelFunc
}

// FetchMongoConfig - reading mongo config from mongo.yml and setting it up
func FetchMongoConfig(configFile string) *MongoConfigData {
	v1 := myconfig.SetupViperAndReadConfig(configFile)

	var mongoConfig1 MongoConfigData
	config := v1.GetStringMapString(myconfig.GetEnv())
	mongoConfig1.Database = config["database"]
	mongoConfig1.Hosts = v1.GetStringSlice(myconfig.GetEnv() + ".hosts")
	mongoConfig1.Username = config["username"]
	mongoConfig1.Password = config["password"]
	mongoConfig1.MongoURI = config["mongo_uri"]
	return &mongoConfig1
}

func (db *MongoDB) Connect() {
	if db.Client == nil {
		config := db.MongoConfig
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		db.Context = context
		db.CancelFunc = cancel
		fmt.Println("Connecting to mongodb...")
		client, err := mongo.Connect(context, options.Client().ApplyURI(config.MongoURI))
		if err != nil {
			fmt.Println("Unable to connect to mongodb...", err)
			log.Error("Unable to connect to mongodb...", err)
			return
		}

		err = client.Ping(context, readpref.Primary())
		fmt.Println("Done: Connecting to mongodb...")
		db.Client = client
	}
}

func (db *MongoDB) Disconnect() {
	db.CancelFunc()
}

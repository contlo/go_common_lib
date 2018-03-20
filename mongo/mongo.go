package mongo

import (
	"fmt"
	"bitbucket.org/roadrunnr/go_common_lib/config"
	"bitbucket.org/roadrunnr/go_common_lib/logger"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	MgoSession  *mgo.Session
	MongoConfig *MongoConfigData
)

type MongoConfigData struct {
	Database string
	Hosts    []string
	Username string
	Password string
}

func InitMongo() {
	if MongoConfig == nil {
		MongoConfig = FetchMongoConfig()
	}
	log.Info("Mongo hosts", MongoConfig.Hosts)

	if MgoSession == nil {
		CreateMongoSession()
	}
}

// FetchMongoConfig - reading mongo config from mongo.yml and setting it up
func FetchMongoConfig() *MongoConfigData {
	v1 := myconfig.SetupViperAndReadConfig("mongo")

	var mongoConfig1 MongoConfigData
	config := v1.GetStringMapString(myconfig.GetEnv())
	mongoConfig1.Database = config["database"]
	mongoConfig1.Hosts = v1.GetStringSlice(myconfig.GetEnv() + ".hosts")
	mongoConfig1.Username = config["username"]
	mongoConfig1.Password = config["password"]
	return &mongoConfig1
}

// creating mongo session
func CreateMongoSession() *mgo.Session {
	if MgoSession == nil {
		config := MongoConfig
		var err error
		info := &mgo.DialInfo{
			Addrs:    config.Hosts,
			Timeout:  60 * time.Second,
			Database: config.Database,
			Username: config.Username,
			Password: config.Password,
		}

		fmt.Println("Connecting to mongo db...")
		MgoSession, err = mgo.DialWithInfo(info)
		if err != nil {
			panic(err) // no, not really
		}
		fmt.Println("Done:: Connecting to mongo db...")
		return MgoSession
	}

	return MgoSession.Clone()
}

// func GetDb() *mgo.Database {
// 	session := CreateMongoSession()
// 	return session.DB(MongoConfig.Database)
// }

func CloseSession() {
	MgoSession.Close()
}

package mongo

import (
	"crypto/tls"
	"fmt"
	myconfig "bitbucket.org/zatasales/go_common_lib/config"
	log "bitbucket.org/zatasales/go_common_lib/logger"
	"gopkg.in/mgo.v2"
	"net"
	"time"
	_ "time"
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
	MongoURI string
}

type MongoDB struct {
	MgoSession  *mgo.Session
	MongoConfig *MongoConfigData
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

// creating mongo session
func (db *MongoDB) CreateMongoSession() *mgo.Session {
	if db.MgoSession == nil {
		config := db.MongoConfig
		var err error
		var dialInfo *mgo.DialInfo
		log.Info("RAMRAM")
		log.Info(config.MongoURI)
		if len(config.MongoURI) > 0 {
			dialInfo, err = mgo.ParseURL(config.MongoURI)
			tlsConfig := &tls.Config{}
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
				return conn, err
			}
		} else {
			dialInfo = &mgo.DialInfo{
				Addrs:    config.Hosts,
				Timeout:  60 * time.Second,
				Database: config.Database,
				Username: config.Username,
				Password: config.Password,
			}
		}


		log.Info("Connecting to mongo db...", err)
		fmt.Println("Connecting to mongo db...")
		MgoSession, err := mgo.DialWithInfo(dialInfo)

		//MgoSession, err = mgo.DialWithInfo(info)
		if err != nil {
			panic(err) // no, not really
		}
		fmt.Println("Done:: Connecting to mongo db...")
		return MgoSession
	}

	return db.MgoSession.Clone()
}

// func GetDb() *mgo.Database {
// 	session := CreateMongoSession()
// 	return session.DB(MongoConfig.Database)
// }

func (db *MongoDB) CloseSession() {
	db.MgoSession.Close()
}

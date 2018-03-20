package gomysql

import (
	"fmt"
	"bitbucket.org/roadrunnr/go_common_lib/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	MasterDB *MysqlDB
	SlaveDB  *MysqlDB
)

type MysqlConfigData struct {
	Database string
	Username string
	Password string
	Host     string
	Port     string
}

type MysqlDB struct {
	DB          *sqlx.DB
	IsSlave     bool
	MysqlConfig *MysqlConfigData
}

func Init() {
	MasterDB = &MysqlDB{IsSlave: false}
	MasterDB.Init()
	if myconfig.IsProduction() {
		SlaveDB = &MysqlDB{IsSlave: true}
		SlaveDB.Init()
	} else {
		SlaveDB = MasterDB
	}
}
func (db *MysqlDB) Init() {
	dbEnv := myconfig.GetEnv()
	if db.IsSlave {
		dbEnv = "production_slave"
	}
	if db.MysqlConfig == nil {
		db.MysqlConfig = FetchMysqlConfig(dbEnv)
	}

	if db.DB == nil {
		dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", db.MysqlConfig.Username, db.MysqlConfig.Password, db.MysqlConfig.Host, db.MysqlConfig.Port, db.MysqlConfig.Database)
		fmt.Println(dataSourceName)
		var err error
		db.DB, err = sqlx.Connect("mysql", dataSourceName)
		if err != nil {
			panic(err)
		}
	}
}

func FetchMysqlConfig(env string) *MysqlConfigData {
	v1 := myconfig.SetupViperAndReadConfig("mysql")

	config := v1.GetStringMapString(env)
	var mysqlConfig MysqlConfigData

	mysqlConfig.Database = config["database"]
	mysqlConfig.Username = config["username"]
	mysqlConfig.Password = config["password"]
	mysqlConfig.Port = config["port"]
	mysqlConfig.Host = config["host"]

	return &mysqlConfig
}

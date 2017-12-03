package gomysql

import (
	"fmt"
	"go_common_lib/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB          *sqlx.DB
	MysqlConfig *MysqlConfigData
)

type MysqlConfigData struct {
	Database string
	Username string
	Password string
	Host     string
	Port     string
}

func Init() {
	if MysqlConfig == nil {
		MysqlConfig = FetchMysqlConfig()
	}
	if DB == nil {
		dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", MysqlConfig.Username, MysqlConfig.Password, MysqlConfig.Host, MysqlConfig.Port, MysqlConfig.Database)
		fmt.Println(dataSourceName)
		var err error
		DB, err = sqlx.Connect("mysql", dataSourceName)
		if err != nil {
			panic(err)
		}
	}
}

func FetchMysqlConfig() *MysqlConfigData {
	v1 := myconfig.SetupViperAndReadConfig("mysql")

	config := v1.GetStringMapString(myconfig.GetEnv())
	var mysqlConfig MysqlConfigData

	mysqlConfig.Database = config["database"]
	mysqlConfig.Username = config["username"]
	mysqlConfig.Password = config["password"]
	mysqlConfig.Port = config["port"]
	mysqlConfig.Host = config["host"]

	return &mysqlConfig
}

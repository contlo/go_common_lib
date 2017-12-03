package myconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var PingAppServer string = ""
var AppServer string = ""

func GetEnv() string {
	return os.Getenv("GO_ENV")
}

// http://stackoverflow.com/questions/24487943/invoke-golang-struct-function-gives-cannot-refer-to-unexported-field-or-method
// Only functions with cap letter are exported
func SetupViperAndReadConfig(configName string) *viper.Viper {
	v1 := viper.New()
  v1.AddConfigPath("../../config")
	v1.AddConfigPath("../config")
	v1.AddConfigPath("./config")

	v1.SetConfigName(configName)
	err1 := v1.ReadInConfig()

	if err1 != nil {
		fmt.Println("No configuration file loaded - " + configName)
	}
	return v1
}

func MetricEnabled() bool {
	return GetEnv() == "production"
}

package log

import (
	"fmt"
	"github.com/op/go-logging"
	"bitbucket.org/zatasales/go_common_lib/config"
	myconfig "go_common_lib/config"
	"os"
)

var Logger = logging.MustGetLogger("pinglogger")

func SetupLogger(fileName string) *os.File {
	f, err := os.OpenFile("log/"+fileName+"_"+myconfig.GetEnv()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
		panic("Unable to set the log file")
	}

	var format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)

	logBackend := logging.NewLogBackend(f, "", 0)
	backendFormatter := logging.NewBackendFormatter(logBackend, format)
	logging.SetBackend(backendFormatter)
	return f
}

func Info(args ...interface{}) {
	Logger.Info(args)
}

func Debug(args ...interface{}) {
	Logger.Debug(args)
}

func Error(args ...interface{}) {
	Logger.Error(args)
}

func Warning(args ...interface{}) {
	Logger.Warning(args)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args)
}

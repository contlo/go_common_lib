package mymqtt

import (
	"encoding/json"
	"fmt"
	"bitbucket.org/roadrunnr/go_common_lib/config"
	"bitbucket.org/roadrunnr/go_common_lib/rhttp"
	"math/rand"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func CreateMqttClientOptions() *MQTT.ClientOptions {
	v1 := myconfig.SetupViperAndReadConfig("mqtt")

	env := myconfig.GetEnv()
	if env != "production" {
		env = "default"
	}
	config := v1.GetStringMapString(env)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", config["host"]+":"+config["port"]))
	opts.SetUsername(config["username"])
	opts.SetPassword(config["password"])
	opts.SetClientID("mqtt_server_" + string(rand.Intn(1000000)))
	fmt.Println("Random: " + "mqtt_server_" + strconv.Itoa(rand.Intn(1000000)))
	return opts
}

func GetTopicPrefix(httpFetcher rhttp.IHttpFetcher) string {
	resp, err := httpFetcher.Get("/notification/mqtt_config")
	if err != nil {
		return ""
	}
	var mqttConfigResponse struct {
		Status struct {
			Code int `json:"code"`
		}
		TopicPrefix string `json:"topic_prefix"`
	}

	json.Unmarshal(resp, &mqttConfigResponse)
	if mqttConfigResponse.Status.Code != 200 {
		fmt.Println("Failed to read the topic prefix")
		return ""
	}

	return mqttConfigResponse.TopicPrefix
}

package mymqtt

import (
	"encoding/json"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"go_common_lib/config"
	"go_common_lib/rhttp"
	"math/rand"
	"strconv"
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
			Code int
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
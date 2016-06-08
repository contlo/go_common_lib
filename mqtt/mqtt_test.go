package mymqtt

import (
	"bytes"
	// "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockHttpFetcher struct {
	mqttConfigResponse string
}

func (fetcher MockHttpFetcher) Get(url string) ([]byte, error) {
	if url == "/notification/mqtt_config" && fetcher.mqttConfigResponse != "" {
		return []byte(fetcher.mqttConfigResponse), nil
		// `{"status": {"code": 200}, "mqtt_config": {"topic_prefix": 128912981}}`
	}
	return []byte(`""`), nil
}

func (fetcher MockHttpFetcher) GetWithAuth(url string, authKey string) ([]byte, error) {
	return nil, nil
}

func (fetcher MockHttpFetcher) Post(url string, buffer *bytes.Buffer) ([]byte, error) {
	return nil, nil
}

func (fetcher MockHttpFetcher) PostWithAuth(url string, buffer *bytes.Buffer, authKey string) ([]byte, error) {
	return nil, nil
}

func TestGetTopicPrefix(t *testing.T) {
	prefix := GetTopicPrefix(&MockHttpFetcher{mqttConfigResponse: `{"status": {"code": 200}, "topic_prefix": "128912981"}`})
	assert.Equal(t, "128912981", prefix, "MQTT prefix mismatch")
}

func TestGetTopicPrefixNil(t *testing.T) {
	prefix := GetTopicPrefix(&MockHttpFetcher{mqttConfigResponse: ""})
	assert.Equal(t, "", prefix, "MQTT prefix must be nil")
}

func TestCreateMqttClientOptions(t *testing.T) {
	test := CreateMqttClientOptions()
	assert.NotNil(t, test, "Client options can't be nil")
}

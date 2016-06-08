package rhttp

import (
	"bytes"
	"go_common_lib/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type IHttpFetcher interface {
	Get(url string) ([]byte, error)
	GetWithAuth(url string, authKey string) ([]byte, error)
	Post(url string, buffer *bytes.Buffer) ([]byte, error)
	PostWithAuth(url string, buffer *bytes.Buffer, authKey string) ([]byte, error)
}

var (
	httpClient *http.Client
)

type HttpFetcher struct {
	Server string
}

func getHttpClient() *http.Client {
	if httpClient == nil {
		timeout := time.Duration(5 * time.Second)
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}
	return httpClient
}

func (fetcher *HttpFetcher) Get(url string) ([]byte, error) {
	return fetcher.GetWithAuth(url, "")
}

func (fetcher *HttpFetcher) GetWithAuth(url string, authKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", fetcher.Server+url, nil)
	if authKey != "" {
		req.Header.Add("Authorization", authKey)
	}
	res, err := getHttpClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to call http get "+url, err)
		return nil, err
	}
	// bytes.Buffer
	return contents, nil
}

func (fetcher *HttpFetcher) Post(url string, buffer *bytes.Buffer) ([]byte, error) {
	return fetcher.PostWithAuth(url, buffer, "")
}

func (fetcher *HttpFetcher) PostWithAuth(url string, buffer *bytes.Buffer, authKey string) ([]byte, error) {
	req, err := http.NewRequest("POST", fetcher.Server+url, buffer)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if authKey != "" {
		req.Header.Add("Authorization", authKey)
	}
	res, err := getHttpClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to call http post "+url, err)
		return nil, err
	}
	// bytes.Buffer
	return contents, nil
}
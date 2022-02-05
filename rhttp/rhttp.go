package rhttp

import (
	"bytes"
	log "bitbucket.org/zatasales/go_common_lib/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type IHttpFetcher interface {
	Get(url string) (*HttpResponse, error)
	GetWithAuth(url string, authKey string) (*HttpResponse, error)
	Post(url string, buffer *bytes.Buffer) (*HttpResponse, error)
	PostWithAuth(url string, buffer *bytes.Buffer, authKey string) (*HttpResponse, error)
	PostWithHeaderAuthKeys(url string, buffer *bytes.Buffer, authKey string, headerMap map[string]string) (*HttpResponse, error)
}

var (
	httpClient *http.Client
)

type HttpFetcher struct {
	Server string
}

type HttpResponse struct {
	StatusCode int
	Content []byte
}

func getHttpClient() *http.Client {
	if httpClient == nil {
		timeout := time.Duration(15 * time.Second)
		httpClient = &http.Client{
			Timeout: timeout,
		}
	}
	return httpClient
}

func (fetcher *HttpFetcher) Get(url string) (*HttpResponse, error) {
	return fetcher.GetWithAuth(url, "")
}

func (fetcher *HttpFetcher) GetWithAuth(url string, authKey string) (*HttpResponse, error) {
	req, err := http.NewRequest("GET", fetcher.Server+url, nil)
	if authKey != "" {
		req.Header.Add("Authorization", authKey)
	}
	res, err := getHttpClient().Do(req)
	if err != nil {
		return &HttpResponse{StatusCode: 400, Content: nil}, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to call http get "+url, err)
		return &HttpResponse{StatusCode: 400, Content: nil}, err
	}
	// bytes.Buffer
	return &HttpResponse{StatusCode: res.StatusCode, Content: contents}, nil
}

func (fetcher *HttpFetcher) Post(url string, buffer *bytes.Buffer) (*HttpResponse, error) {
	return fetcher.PostWithAuth(url, buffer, "")
}

func (fetcher *HttpFetcher) PostWithAuth(url string, buffer *bytes.Buffer, authKey string) (*HttpResponse, error) {
	httpResponse, err := fetcher.PostWithHeaderAuthKeys(url, buffer, authKey, map[string]string{})
	return httpResponse, err
}

func (fetcher *HttpFetcher) PostWithHeaderAuthKeys(url string, buffer *bytes.Buffer, authKey string, headerMap map[string]string) (*HttpResponse, error) {
	req, err := http.NewRequest("POST", fetcher.Server+url, buffer)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	if authKey != "" {
		req.Header.Add("Authorization", authKey)
	}

	for key, value := range headerMap {
		req.Header.Add(key, value)
	}

	res, err := getHttpClient().Do(req)
	if err != nil {
		return &HttpResponse{StatusCode: 400, Content:nil}, err
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to call http post "+url, err)
		return &HttpResponse{StatusCode: 400, Content:nil}, err
	}

	// bytes.Buffer
	return &HttpResponse{StatusCode: res.StatusCode, Content: contents}, nil
}

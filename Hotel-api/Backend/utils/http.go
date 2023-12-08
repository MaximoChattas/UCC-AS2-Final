package utils

import "net/http"

type HttpClient struct{}

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

func (h *HttpClient) Do(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	return client.Do(req)
}

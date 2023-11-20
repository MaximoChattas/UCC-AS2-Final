package utils

import "net/http"

type HttpClient struct{}

type HttpClientInterface interface {
	Get(url string) (*http.Response, error)
}

func (h *HttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

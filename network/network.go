package network

import (
	"fmt"
	"net/http"
)

type HttpClient struct {
	baseURL string
}

func NewClient(baseURL string) *HttpClient {
	return &HttpClient{
		baseURL,
	}
}

func (c *HttpClient) GET(endpoit string, query map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", c.baseURL, endpoit), nil)
	req.Header.Add("User-Agent", "prayertime cli")
	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	return http.DefaultClient.Do(req)

}

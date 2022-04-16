package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrorNotFound = errors.New("not found")
)

type (
	Tweet struct {
		Id   string `json:"id"`
		Text string `json:"text"`
	}

	GetTweetsResponse struct {
		Data []Tweet `json:"data"`
	}

	Service interface {
		GetTweets(ctx context.Context) (*GetTweetsResponse, error)
	}

	v1 struct {
		baseURL string
		client  *http.Client
		timeout time.Duration
	}
)

func New(baseURL string, client *http.Client, timeout time.Duration) *v1 {
	return &v1{
		baseURL: baseURL,
		client:  client,
		timeout: timeout,
	}
}

func (v *v1) GetTweets(ctx context.Context) (*GetTweetsResponse, error) {
	url := fmt.Sprintf("%s/2/users/1234/tweets", v.baseURL)
	fmt.Printf("url: %s", url)
	ctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w. %s", ErrorNotFound, http.StatusText(resp.StatusCode))
	}

	var d *GetTweetsResponse
	return d, json.NewDecoder(resp.Body).Decode(&d)
}

package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vorhundert.de/vhitter/cmd/web/config"
	errorDto "vorhundert.de/vhitter/cmd/web/errorDtos"
)

type (
	Service interface {
		GetTweets(ctx context.Context) (*GetTweetsResponse, error)
	}

	ServiceImpl struct {
		config  config.TwitterConfig
		baseURL string
		client  *http.Client
		timeout time.Duration
	}
)

func New(config config.TwitterConfig, baseURL string, client *http.Client, timeout time.Duration) *ServiceImpl {
	return &ServiceImpl{
		config:  config,
		baseURL: baseURL,
		client:  client,
		timeout: timeout,
	}
}

func (serviceImpl *ServiceImpl) GetTweets(ctx context.Context) (*GetTweetsResponse, error) {
	url := fmt.Sprintf("%s/2/users/%s/tweets", serviceImpl.baseURL, serviceImpl.config.UserId)
	ctx, cancel := context.WithTimeout(ctx, serviceImpl.timeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	req.Header.Set("Authorization", "Bearer "+serviceImpl.config.Bearer)

	resp, err := serviceImpl.client.Do(req)
	if err != nil {
		return nil, &errorDto.ErrorConnection{Msg: "No connection possible", Cause: err}
	}
	defer resp.Body.Close()

	if parsedError := checkResponseForError(resp); parsedError != nil {
		return nil, parsedError
	}

	var d *GetTweetsResponse
	return d, json.NewDecoder(resp.Body).Decode(&d)
}

func checkResponseForError(resp *http.Response) error {
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("%w. %s", errorDto.ErrorUnauthorized, http.StatusText(resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w. %s", errorDto.ErrorNotFound, http.StatusText(resp.StatusCode))
	}
	return nil
}

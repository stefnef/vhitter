package twitter_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"vorhundert.de/vhitter/cmd/web/twitter"
)

type twitterMockData struct {
	called       bool
	Bearer       string
	responseData string
	statusCode   int
}

var (
	twitterMockServer *httptest.Server
	mockData          *twitterMockData
	service           twitter.Service
)

func fatal(t *testing.T, want, got interface{}) {
	t.Helper()
	t.Fatalf(`want: %v, got: %v`, want, got)
}

func initMockData() {
	mockData = &twitterMockData{
		called:       false,
		Bearer:       "Bearer <ACCESS_TOKEN>",
		statusCode:   http.StatusOK,
		responseData: `{"data":[{"id":"","text":""}]}`,
	}
}

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	mux.HandleFunc("/2/users/1234/tweets", handleGetTweets) //TODO UserId aus Config nehmen

	fmt.Printf("mocking server")
	twitterMockServer = httptest.NewServer(mux)
	defer twitterMockServer.Close()

	fmt.Println("mocking external")
	service = twitter.New(twitterMockServer.URL, http.DefaultClient, time.Second)

	fmt.Println("run tests")
	m.Run()
}

func handleGetTweets(w http.ResponseWriter, r *http.Request) {
	mockData.called = true

	w.Header().Set("Content-Type", "application/json")

	bearer := r.Header.Get("Authorization")
	if bearer != mockData.Bearer {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"title": "Unauthorized", "type": "about:blank", "status": 401, "detail": "Unauthorized" }`))
		return
	}

	w.WriteHeader(mockData.statusCode)
	w.Write([]byte(mockData.responseData))
}

func Test_GetTweets_should_request_server(t *testing.T) {
	initMockData()
	service.GetTweets(context.Background())
	if !mockData.called {
		fatal(t, true, mockData.called)
	}
}

func Test_GetTweets_should_throw_error_if_status_code_is_not_ok(t *testing.T) {
	initMockData()
	mockData.statusCode = http.StatusNotFound

	foundData, foundError := service.GetTweets(context.Background())
	if foundData != nil {
		fatal(t, nil, foundData)
	}
	if foundError == nil { // TODO Needed Error classes
		fatal(t, twitter.ErrorNotFound, foundError)
	}
}

func Test_GetTweets_should_be_parsed(t *testing.T) {
	tt := []struct {
		testCase string
		mockData string
		wantData *twitter.GetTweetsResponse
		wantErr  error
	}{
		{
			"success",
			`{"data":[{"id":"1234","text":"this is text"}]}`,
			&twitter.GetTweetsResponse{
				Data: []twitter.Tweet{
					{Id: "1234", Text: "this is text"},
				},
			},
			nil,
		},
		{
			"with meta data",
			`{"data":[{"id":"12345","text":"this is text"}],"meta":{"result_count":1,"newest_id":"12345","oldest_id":"12345"}}`,
			&twitter.GetTweetsResponse{
				Data: []twitter.Tweet{
					{Id: "12345", Text: "this is text"},
				},
			},
			nil,
		},
		{
			"more than one tweet",
			`{"data":[{"id":"12345","text":"this is text"}, {"id":"444","text":"an other tweet"}],"meta":{"result_count":2,"newest_id":"12345","oldest_id":"444"}}`,
			&twitter.GetTweetsResponse{
				Data: []twitter.Tweet{
					{Id: "12345", Text: "this is text"},
					{Id: "444", Text: "an other tweet"},
				},
			},
			nil,
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.testCase, func(t *testing.T) {
			initMockData()
			mockData.responseData = tc.mockData
			gotData, gotErr := service.GetTweets(context.Background())

			if !errors.Is(gotErr, tc.wantErr) {
				fatal(t, tc.wantErr, gotErr)
			}

			if !reflect.DeepEqual(gotData, tc.wantData) {
				fatal(t, tc.wantData, gotData)
			}
		})
	}
}

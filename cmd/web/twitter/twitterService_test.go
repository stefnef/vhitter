package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"vorhundert.de/vhitter/cmd/web/config"
	errorDto "vorhundert.de/vhitter/cmd/web/errorDtos"
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
	mockConfig        *config.TwitterConfig = &config.TwitterConfig{UserId: "1234", Bearer: "<ACCESS_TOKEN>"}
	service           Service
)

func fatal(t *testing.T, want, got interface{}) {
	t.Helper()
	t.Fatalf(`want: %v, got: %v`, want, got)
}

func initMockData() {
	mockData = &twitterMockData{
		called:       false,
		Bearer:       "Bearer " + mockConfig.Bearer,
		statusCode:   http.StatusOK,
		responseData: `{"data":[{"id":"","text":""}]}`,
	}
}

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	mux.HandleFunc("/2/users/"+mockConfig.UserId+"/tweets", handleGetTweets)

	fmt.Printf("mocking server")
	twitterMockServer = httptest.NewServer(mux)
	defer twitterMockServer.Close()

	fmt.Println("mocking external")
	service = New(*mockConfig, twitterMockServer.URL, http.DefaultClient, time.Second)

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
		fatal(t, errorDto.ErrorNotFound, foundError)
	}
}

func Test_GetTweets_should_handle_Unauthorized_error(t *testing.T) {
	initMockData()
	mockData.Bearer = "Something Unusual"
	gotData, gotErr := service.GetTweets(context.Background())

	if !errors.Is(gotErr, errorDto.ErrorUnauthorized) {
		fatal(t, nil, gotErr)
	}

	if gotData != nil {
		fatal(t, nil, gotData)
	}
}

func Test_GetTweets_should_report_parsing_error(t *testing.T) {
	initMockData()
	mockData.responseData = "d}"
	gotData, gotErr := service.GetTweets(context.Background())

	if syntaxErr, ok := gotErr.(*json.SyntaxError); !ok {
		fatal(t, &json.SyntaxError{Offset: 1}, gotErr)
	} else if syntaxErr.Offset != 1 {
		fatal(t, 1, syntaxErr)

	}

	if gotData != nil {
		fatal(t, nil, gotData)
	}
}

func Test_GetTweets_should_propagade_error_if_server_is_not_reachable(t *testing.T) {
	newClient := New(*mockConfig, "http://some/url/which/doesn/not/fit", http.DefaultClient, time.Second)

	initMockData()
	gotData, gotErr := newClient.GetTweets(context.Background())
	wantErr := &errorDto.ErrorConnection{Msg: "No connection possible"}
	if gotData != nil {
		fatal(t, nil, gotData)
	}
	if !errors.Is(gotErr, wantErr) {
		fatal(t, wantErr, gotErr)
	}
}

func Test_GetTweets_should_be_parsed(t *testing.T) {
	tt := []struct {
		testCase string
		mockData string
		wantData *GetTweetsResponse
	}{
		{
			"success",
			`{"data":[{"id":"1234","text":"this is text"}]}`,
			&GetTweetsResponse{
				Data: []Tweet{
					{Id: "1234", Text: "this is text"},
				},
			},
		},
		{
			"with meta data",
			`{"data":[{"id":"12345","text":"this is text"}],"meta":{"result_count":1,"newest_id":"12345","oldest_id":"12345"}}`,
			&GetTweetsResponse{
				Data: []Tweet{
					{Id: "12345", Text: "this is text"},
				},
			},
		},
		{
			"more than one tweet",
			`{"data":[{"id":"12345","text":"this is text"}, {"id":"444","text":"an other tweet"}],"meta":{"result_count":2,"newest_id":"12345","oldest_id":"444"}}`,
			&GetTweetsResponse{
				Data: []Tweet{
					{Id: "12345", Text: "this is text"},
					{Id: "444", Text: "an other tweet"},
				},
			},
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.testCase, func(t *testing.T) {
			initMockData()
			mockData.responseData = tc.mockData
			gotData, gotErr := service.GetTweets(context.Background())

			if gotErr != nil {
				fatal(t, nil, gotErr.(*json.SyntaxError).Offset)
			}

			if !reflect.DeepEqual(gotData, tc.wantData) {
				fatal(t, tc.wantData, gotData)
			}
		})
	}
}

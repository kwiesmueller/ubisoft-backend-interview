package feedback

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/bborbe/http/requestbuilder"
	"github.com/gorilla/mux"

	"github.com/playnet-public/libs/log"
)

func TestService_Handler(t *testing.T) {
	svc := New(log.NewNop(), nil)
	h := svc.Handler()
	if h == nil {
		t.Error("Service.Handler() is nil")
	}
}

func TestService_getEntries(t *testing.T) {
	svc := New(log.NewNop(), nil)

	tests := []struct {
		name        string
		wantEntries []Entry
		request     requestbuilder.HttpRequestBuilder

		getLatestFunc         func(uint) ([]Entry, error)
		getLatestFilteredFunc func(uint, int) ([]Entry, error)

		wantErr bool
	}{
		{
			"emptyList",
			[]Entry{},
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET"),
			nil,
			func(n uint, f int) ([]Entry, error) {
				return []Entry{}, nil
			},
			false,
		},
		{
			"error",
			nil,
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET"),
			func(n uint) ([]Entry, error) {
				return []Entry{}, errors.New("test error")
			},
			func(n uint, f int) ([]Entry, error) {
				return []Entry{}, errors.New("test error")
			},
			true,
		},
		{
			"filteredList",
			[]Entry{
				{"1", "1", "1", 1, ""},
				{"3", "3", "1", 1, ""},
				{"5", "5", "1", 1, ""},
			},
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET").AddParameter("filter", "1"),
			nil,
			func(n uint, f int) ([]Entry, error) {
				return []Entry{
					{"1", "1", "1", 1, ""},
					{"3", "3", "1", 1, ""},
					{"5", "5", "1", 1, ""},
				}, nil
			},
			false,
		},
		{
			"customLimitList",
			[]Entry{
				{"1", "1", "1", 1, ""},
			},
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET").AddParameter("filter", "1").AddParameter("limit", "1"),
			nil,
			func(n uint, f int) ([]Entry, error) {
				return []Entry{
					{"1", "1", "1", 1, ""},
				}, nil
			},
			false,
		},
		{
			"invalidFilter",
			nil,
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET").AddParameter("filter", "abc"),
			nil,
			func(n uint, f int) ([]Entry, error) {
				return []Entry{
					{"1", "1", "1", 1, ""},
					{"3", "3", "1", 1, ""},
					{"5", "5", "1", 1, ""},
				}, nil
			},
			true,
		},
		{
			"filterError",
			nil,
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET").AddParameter("filter", "1"),
			nil,
			func(n uint, f int) ([]Entry, error) {
				return []Entry{}, errors.New("test error")
			},
			true,
		},
		{
			"invalidLimit",
			nil,
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/list").
				SetMethod("GET").AddParameter("filter", "1").AddParameter("limit", "x"),
			nil,
			func(n uint, f int) ([]Entry, error) {
				return []Entry{
					{"1", "1", "1", 1, ""},
				}, nil
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc.repo = newMockRepository(nil, tt.getLatestFunc, tt.getLatestFilteredFunc)

			req, err := tt.request.Build()
			if err != nil {
				t.Error(err)
			}

			w := httptest.NewRecorder()

			if err := svc.getEntries(w, req); (err != nil) != tt.wantErr {
				t.Errorf("Service.getEntries() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				retErr := errorResponse{}
				err = json.NewDecoder(w.Result().Body).Decode(&retErr)
				if err != nil {
					t.Error("failed to decode response", w.Result().Body)
				}
				if retErr.Err == "" {
					t.Error("Service.getEntries() wanted err", retErr)
				}
			}

			var entries []Entry
			err = json.NewDecoder(w.Result().Body).Decode(&entries)
			if err != nil {

			}

			if !reflect.DeepEqual(entries, tt.wantEntries) {
				t.Errorf("Service.getEntries() got = %v, want %v", entries, tt.wantEntries)
			}
		})
	}
}

func TestService_addEntry(t *testing.T) {
	svc := New(log.NewNop(), nil)

	tests := []struct {
		name    string
		request requestbuilder.HttpRequestBuilder

		addFunc func(Entry) error

		wantErr bool
	}{
		{
			"emptyEntry",
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/0").
				SetMethod("POST").SetBody(strings.NewReader(`{
					"rating": 1,
					"comment": ""
				}`)).AddHeader("Ubi-UserId", "1"),
			func(e Entry) error {
				return nil
			},
			false,
		},
		{
			"jsonError",
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/0").
				SetMethod("POST").SetBody(strings.NewReader(`{
					"rating"
				}`)).AddHeader("Ubi-UserId", "1"),
			func(e Entry) error {
				return nil
			},
			true,
		},
		{
			"noUserID",
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/0").
				SetMethod("POST").SetBody(strings.NewReader(`{
					"rating": 1,
					"comment": ""
				}`)).AddHeader("Ubi-UserId", ""),
			func(e Entry) error {
				return nil
			},
			true,
		},
		{
			"noUserIDHeader",
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/0").
				SetMethod("POST").SetBody(strings.NewReader(`{
					"rating": 1,
					"comment": ""
				}`)),
			func(e Entry) error {
				return nil
			},
			true,
		},
		{
			"addError",
			requestbuilder.NewHTTPRequestBuilder("http://127.0.0.1:8080/0").
				SetMethod("POST").SetBody(strings.NewReader(`{
					"rating": 1,
					"comment": ""
				}`)).AddHeader("Ubi-UserId", "1"),
			func(e Entry) error {
				return errors.New("test error")
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc.repo = newMockRepository(tt.addFunc, func(n uint) ([]Entry, error) {
				return []Entry{}, nil
			}, nil)
			m := mux.NewRouter()
			m.Path("/{sessionID}").Methods("POST").HandlerFunc(svc.MakeHandler(svc.addEntry))

			req, err := tt.request.Build()
			if err != nil {
				t.Error(err)
			}

			w := httptest.NewRecorder()
			m.ServeHTTP(w, req)

			retErr := errorResponse{}
			err = json.NewDecoder(w.Result().Body).Decode(&retErr)
			if err != nil {
				t.Error("failed to decode response", err)
			}
			if (retErr.Err == "") == tt.wantErr {
				t.Error("Service.getEntries() wanted err", retErr)
			}
		})
	}
}

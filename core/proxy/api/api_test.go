package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telepair/telepair/pkg/httpclient"
)

func TestAPI_Parse(t *testing.T) {
	tests := []struct {
		name    string
		api     API
		wantErr bool
	}{
		{
			name: "valid GET request",
			api: API{
				Method: "GET",
				URL:    "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "valid POST request",
			api: API{
				Method: "POST",
				URL:    "http://example.com",
				Body:   "test body",
			},
			wantErr: false,
		},
		{
			name: "invalid method",
			api: API{
				Method: "INVALID",
				URL:    "http://example.com",
			},
			wantErr: true,
		},
		{
			name: "empty URL and URLs",
			api: API{
				Method: "GET",
			},
			wantErr: true,
		},
		{
			name: "single URL in URLs",
			api: API{
				Method: "GET",
				URLs:   []string{"http://example.com"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.api.Parse()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPI_Do(t *testing.T) {
	tests := []struct {
		name           string
		api            API
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantStatusCode int
		wantErr        bool
	}{
		{
			name: "successful GET request",
			api: API{
				Method: "GET",
				Headers: map[string]string{
					"X-Test": "test",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "test", r.Header.Get("X-Test"))
				w.WriteHeader(http.StatusOK)
			},
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name: "failed request with non-success status code",
			api: API{
				Method: "GET",
			},
			serverResponse: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			tt.api.URL = server.URL
			err := tt.api.Parse()
			assert.NoError(t, err)

			resp, err := tt.api.Do()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if resp != nil {
				assert.Equal(t, tt.wantStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestAPI_IsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		api      API
		response *http.Response
		want     bool
	}{
		{
			name: "default success codes",
			api:  API{},
			response: &http.Response{
				StatusCode: http.StatusOK,
			},
			want: true,
		},
		{
			name: "custom success codes",
			api: API{
				Config: Config{
					Checker: Checker{
						SuccessCodes: []int{299},
					},
				},
			},
			response: &http.Response{
				StatusCode: 299,
			},
			want: true,
		},
		{
			name: "header match success",
			api: API{
				Config: Config{
					Checker: Checker{
						HeaderMatch: map[string]string{
							"X-Test": "test",
						},
					},
				},
			},
			response: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Test": []string{"test"},
				},
			},
			want: true,
		},
		{
			name: "header match failure",
			api: API{
				Config: Config{
					Checker: Checker{
						HeaderMatch: map[string]string{
							"X-Test": "test",
						},
					},
				},
			},
			response: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Test": []string{"wrong"},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.api.IsSuccess(tt.response)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAPI_Options(t *testing.T) {
	ctx := context.Background()
	mockClient := httpclient.DefaultClient

	tests := []struct {
		name    string
		opts    []Option
		wantCtx context.Context
	}{
		{
			name: "with context",
			opts: []Option{WithContext(ctx)},
		},
		{
			name: "with client",
			opts: []Option{WithClient(mockClient)},
		},
		{
			name: "with both options",
			opts: []Option{WithContext(ctx), WithClient(mockClient)},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &API{
				Method: "GET",
				URL:    server.URL,
				Config: Config{
					Timeout: time.Second,
				},
			}

			_, err := api.Do(tt.opts...)
			assert.NoError(t, err)
		})
	}
}

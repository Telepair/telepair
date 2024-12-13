package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/telepair/telepair/pkg/httpclient"
)

func ExampleRegisterAPI() {
	data := []byte(`
- name: my-eip
  method: GET
  urls:
    - https://api.ipify.org
    - https://checkip.amazonaws.com
    - https://ifconfig.co/ip
    - https://ifconfig.me/ip
  headers:
    Content-Type: application/json
  config:
    checker:
      success_codes: [200]
    fallback:
      selector: random
    timeout: 10s
`)
	if err := RegisterAPIData("yaml", data); err != nil {
		panic(err)
	}
	resp, err := Do("my-eip")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}

func TestAPI_Parse(t *testing.T) {
	tests := []struct {
		name    string
		api     API
		wantErr bool
	}{
		{
			name: "valid GET request",
			api: API{
				Method: "get",
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
			name: "missing URL and URLs",
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test") == "fail" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := []struct {
		name    string
		api     API
		wantErr bool
	}{
		{
			name: "successful request",
			api: API{
				Method: "GET",
				URL:    server.URL,
				Config: Config{
					Timeout: 5 * time.Second,
				},
			},
			wantErr: false,
		},
		{
			name: "failed request",
			api: API{
				Method: "GET",
				URL:    server.URL,
				Headers: map[string]string{
					"X-Test": "fail",
				},
				Config: Config{
					Timeout: 5 * time.Second,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.api.Do()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
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
							"X-Test": "pass",
						},
					},
				},
			},
			response: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Test": []string{"pass"},
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
							"X-Test": "pass",
						},
					},
				},
			},
			response: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Test": []string{"fail"},
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

func TestAPI_doWithFallback(t *testing.T) {
	fallSvc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer fallSvc.Close()
	okSvc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer okSvc.Close()

	api := API{
		Method: "GET",
		URL:    fallSvc.URL,
		URLs:   []string{fallSvc.URL, okSvc.URL},
	}

	resp, err := api.doWithFallback(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "ok", string(body))
}

func TestWithClient(t *testing.T) {
	client := httpclient.New()

	// Create config and apply the WithClient option
	cfg := &config{}
	opt := WithClient(client)
	opt(cfg)

	// Assert the client was properly set
	assert.Equal(t, client, cfg.client)
}

func TestWithContext(t *testing.T) {
	// Create a test context
	ctx := context.Background()

	// Create config and apply the WithContext option
	cfg := &config{}
	opt := WithContext(ctx)
	opt(cfg)

	// Assert the context was properly set
	assert.Equal(t, ctx, cfg.ctx)
}

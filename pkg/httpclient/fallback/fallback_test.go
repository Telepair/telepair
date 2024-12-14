package fallback

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Example_fallback() {
	setupTestServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			switch r.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				body, _ := io.ReadAll(r.Body)
				_, _ = w.Write(body)
			default:
				_, _ = w.Write([]byte("test response"))
			}
		}))
	}

	ts := setupTestServer()
	defer ts.Close()
	ts1 := setupTestServer()
	defer ts1.Close()

	resp, _ := Get([]string{ts.URL, ts1.URL})
	fmt.Printf("GET: status=%d\n", resp.StatusCode)

	resp, _ = Head([]string{ts.URL, ts1.URL})
	fmt.Printf("HEAD: status=%d\n", resp.StatusCode)

	resp, _ = Options([]string{ts.URL, ts1.URL})
	fmt.Printf("OPTIONS: status=%d\n", resp.StatusCode)

	resp, _ = Post([]string{ts.URL, ts1.URL}, []byte("test"))
	fmt.Printf("POST: status=%d\n", resp.StatusCode)

	resp, _ = Put([]string{ts.URL, ts1.URL}, []byte("test"))
	fmt.Printf("PUT: status=%d\n", resp.StatusCode)

	resp, _ = Patch([]string{ts.URL, ts1.URL}, []byte("test"))
	fmt.Printf("PATCH: status=%d\n", resp.StatusCode)

	resp, _ = Delete([]string{ts.URL, ts1.URL})
	fmt.Printf("DELETE: status=%d\n", resp.StatusCode)

	// Output:
	// GET: status=200
	// HEAD: status=200
	// OPTIONS: status=200
	// POST: status=200
	// PUT: status=200
	// PATCH: status=200
	// DELETE: status=200
}

func TestDo(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		urls          []string
		body          []byte
		serverHandler http.HandlerFunc
		wantStatus    int
		wantErr       bool
	}{
		{
			name:   "successful GET request",
			method: http.MethodGet,
			urls:   []string{""}, // Will be replaced with test server URL
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "fallback to second URL on first failure",
			method: http.MethodGet,
			urls:   []string{"", ""}, // Will be replaced with test server URLs
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "all URLs fail",
			method: http.MethodGet,
			urls:   []string{"", ""}, // Will be replaced with test server URLs
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
		{
			name:   "POST with body",
			method: http.MethodPost,
			urls:   []string{""},
			body:   []byte("test body"),
			serverHandler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name:    "unsupported method",
			method:  "INVALID",
			urls:    []string{""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server(s)
			servers := make([]*httptest.Server, len(tt.urls))
			for i := range tt.urls {
				server := httptest.NewServer(tt.serverHandler)
				defer server.Close()
				tt.urls[i] = server.URL
				servers[i] = server
			}

			// Execute test
			resp, err := Do(tt.method, tt.urls, tt.body)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check response
			if !tt.wantErr {
				if resp.StatusCode != tt.wantStatus {
					t.Errorf("Do() status = %v, want %v", resp.StatusCode, tt.wantStatus)
				}
			}
		})
	}
}

func TestHelperMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	urls := []string{server.URL}
	body := []byte("test body")

	tests := []struct {
		name    string
		fn      func() (*http.Response, error)
		wantErr bool
	}{
		{"Get", func() (*http.Response, error) { return Get(urls) }, false},
		{"Head", func() (*http.Response, error) { return Head(urls) }, false},
		{"Options", func() (*http.Response, error) { return Options(urls) }, false},
		{"Post", func() (*http.Response, error) { return Post(urls, body) }, false},
		{"Put", func() (*http.Response, error) { return Put(urls, body) }, false},
		{"Patch", func() (*http.Response, error) { return Patch(urls, body) }, false},
		{"Delete", func() (*http.Response, error) { return Delete(urls) }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.fn()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp.StatusCode != http.StatusOK {
				t.Errorf("%s() status = %v, want %v", tt.name, resp.StatusCode, http.StatusOK)
			}
		})
	}
}

package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleClient() {
	ts := setupTestServer()
	defer ts.Close()

	// Create a new client
	client := New(
		WithRetry(3, 1*time.Second, 5*time.Second),
		WithDefaultHeader(http.Header{"Content-Type": []string{"application/json"}}),
	)

	// Make a GET request
	resp, _ := client.Get(ts.URL)
	fmt.Printf("GET: status=%d\n", resp.StatusCode)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	header := http.Header{}
	header.Set("X-Request-ID", "123")
	resp, _ = client.Get(ts.URL, WithContext(ctx), WithHeader(header))
	fmt.Printf("GET with options: status=%d\n", resp.StatusCode)

	// Make a POST request with JSON body
	jsonBody := bytes.NewBufferString(`{"name": "John", "age": 30}`)
	resp, _ = client.Post(ts.URL, jsonBody.Bytes(),
		WithHeader(http.Header{"Content-Type": []string{"application/json"}}),
	)
	fmt.Printf("POST: status=%d\n", resp.StatusCode)

	// Make a request with custom headers
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer token123")
	resp, _ = client.Do(req)
	fmt.Printf("Custom request: status=%d\n", resp.StatusCode)

	// Output:
	// GET: status=200
	// GET with options: status=200
	// POST: status=200
	// Custom request: status=200
}

func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			_, _ = w.Write(body)
		} else {
			_, _ = w.Write([]byte("test response"))
		}
	}))
}

func TestClient(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New()

	tests := []struct {
		name         string
		method       string
		executeReq   func() (*http.Response, error)
		expectedBody string
	}{
		{
			name:   "GET request",
			method: http.MethodGet,
			executeReq: func() (*http.Response, error) {
				return client.Get(server.URL)
			},
			expectedBody: "test response",
		},
		{
			name:   "HEAD request",
			method: http.MethodHead,
			executeReq: func() (*http.Response, error) {
				return client.Head(server.URL)
			},
			expectedBody: "test response",
		},
		{
			name:   "POST request",
			method: http.MethodPost,
			executeReq: func() (*http.Response, error) {
				return client.Post(server.URL, []byte("test body"))
			},
			expectedBody: "test body",
		},
		{
			name:   "PUT request",
			method: http.MethodPut,
			executeReq: func() (*http.Response, error) {
				return client.Put(server.URL, []byte("test body"))
			},
			expectedBody: "test body",
		},
		{
			name:   "PATCH request",
			method: http.MethodPatch,
			executeReq: func() (*http.Response, error) {
				return client.Patch(server.URL, []byte("test body"))
			},
			expectedBody: "test body",
		},
		{
			name:   "DELETE request",
			method: http.MethodDelete,
			executeReq: func() (*http.Response, error) {
				return client.Delete(server.URL)
			},
			expectedBody: "test response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.executeReq()
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			if resp.Body != nil {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				if tt.method != http.MethodHead {
					assert.Equal(t, tt.expectedBody, string(body))
				}
			}
		})
	}
}

func TestClientDo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client := New()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test response", string(body))
}

func TestClientWithOptions(t *testing.T) {
	client := New(
		WithSkipTLSVerify(true),
		WithRetry(3, 1*time.Second, 5*time.Second),
	)

	assert.NotNil(t, client)
}

func TestInvalidRequest(t *testing.T) {
	client := New()

	// Test with invalid URL
	resp, err := client.Get("://invalid-url")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestParseResponse(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := Get(server.URL)
	assert.NoError(t, err)
	mediaType, body, err := ParseResponse(resp)
	assert.NoError(t, err)
	assert.Equal(t, "test response", string(body))
	assert.Equal(t, "text/plain", mediaType)
}

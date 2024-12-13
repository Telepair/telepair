package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleDefaultClient() {
	ts := setupTestServer()
	defer ts.Close()

	resp, _ := Get(ts.URL)
	fmt.Printf("GET: status=%d\n", resp.StatusCode)

	resp, _ = Get(ts.URL,
		WithContext(context.Background()),
		WithHeader(http.Header{"X-Request-ID": []string{"123"}}),
	)
	fmt.Printf("GET with options: status=%d\n", resp.StatusCode)

	jsonBody := bytes.NewBufferString(`{"name": "John", "age": 30}`)
	resp, _ = Post(ts.URL, jsonBody.Bytes(),
		WithHeader(http.Header{"Content-Type": []string{"application/json"}}),
	)
	fmt.Printf("POST: status=%d\n", resp.StatusCode)

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer token123")
	resp, _ = Do(req)
	fmt.Printf("Custom request: status=%d\n", resp.StatusCode)

	// Output:
	// GET: status=200
	// GET with options: status=200
	// POST: status=200
	// Custom request: status=200
}

func TestDefaultClientMethods(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	tests := []struct {
		name           string
		operation      func() (*http.Response, error)
		expectedMethod string
		expectedBody   string
	}{
		{
			name: "GET request",
			operation: func() (*http.Response, error) {
				return Get(server.URL)
			},
			expectedMethod: "GET",
			expectedBody:   "test response",
		},
		{
			name: "POST request",
			operation: func() (*http.Response, error) {
				return Post(server.URL, []byte("test-body"))
			},
			expectedMethod: "POST",
			expectedBody:   "test-body",
		},
		{
			name: "PUT request",
			operation: func() (*http.Response, error) {
				return Put(server.URL, []byte("test-body"))
			},
			expectedMethod: "PUT",
			expectedBody:   "test-body",
		},
		{
			name: "PATCH request",
			operation: func() (*http.Response, error) {
				return Patch(server.URL, []byte("test-body"))
			},
			expectedMethod: "PATCH",
			expectedBody:   "test-body",
		},
		{
			name: "DELETE request",
			operation: func() (*http.Response, error) {
				return Delete(server.URL)
			},
			expectedMethod: "DELETE",
			expectedBody:   "",
		},
		{
			name: "HEAD request",
			operation: func() (*http.Response, error) {
				return Head(server.URL)
			},
			expectedMethod: "HEAD",
			expectedBody:   "",
		},
		{
			name: "OPTIONS request",
			operation: func() (*http.Response, error) {
				return Options(server.URL)
			},
			expectedMethod: "OPTIONS",
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := tt.operation()
			assert.NoError(t, err)
			assert.NotNil(t, resp)

			// For methods with body, verify the response body
			if tt.expectedBody != "" {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
				resp.Body.Close()
			}
		})
	}
}

func TestDo(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	assert.NoError(t, err)

	resp, err := Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

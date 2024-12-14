package fallback

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/telepair/telepair/pkg/httpclient"
)

func TestWithSelectStrategy(t *testing.T) {
	tests := []struct {
		name     string
		strategy SelectStrategy
		want     SelectStrategy
	}{
		{
			name:     "round robin strategy",
			strategy: SelectStrategyRoundRobin,
			want:     SelectStrategyRoundRobin,
		},
		{
			name:     "random strategy",
			strategy: SelectStrategyRandom,
			want:     SelectStrategyRandom,
		},
		{
			name:     "invalid strategy",
			strategy: "invalid",
			want:     SelectStrategyRoundRobin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fallback{}
			WithSelectStrategy(tt.strategy)(f)
			assert.Equal(t, tt.want, f.selector)
		})
	}
}

func TestWithClient(t *testing.T) {
	c := httpclient.New()
	tests := []struct {
		name   string
		client httpclient.Client
		want   httpclient.Client
	}{
		{
			name:   "valid client",
			client: c,
			want:   c,
		},
		{
			name:   "nil client",
			client: nil,
			want:   httpclient.NoRetryClient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fallback{}
			WithRetryClient(tt.client)(f)
			assert.Equal(t, tt.want, f.client)
		})
	}
}

func TestWithRetryChecker(t *testing.T) {
	customChecker := func(resp *http.Response) bool {
		if resp != nil {
			return resp.StatusCode >= http.StatusBadRequest
		}
		return false
	}

	tests := []struct {
		name    string
		checker RetryChecker
		want    RetryChecker
	}{
		{
			name:    "valid checker",
			checker: customChecker,
			want:    customChecker,
		},
		{
			name:    "nil checker",
			checker: nil,
			want:    DefaultRetry,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fallback{}
			WithRetryChecker(tt.checker)(f)
			assert.NotNil(t, f.retry)
			if tt.checker != nil {
				assert.Equal(t, tt.want(nil), f.retry(nil))
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	customLogger := slog.Default()

	tests := []struct {
		name   string
		logger *slog.Logger
		want   *slog.Logger
	}{
		{
			name:   "valid logger",
			logger: customLogger,
			want:   customLogger,
		},
		{
			name:   "nil logger",
			logger: nil,
			want:   slog.With("component", "httpclient/fallback"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fallback{}
			WithLogger(tt.logger)(f)
			assert.NotNil(t, f.logger)
		})
	}
}

func TestDefaultRetry(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"internal server error", http.StatusInternalServerError, true},
		{"bad gateway", http.StatusBadGateway, true},
		{"service unavailable", http.StatusServiceUnavailable, true},
		{"gateway timeout", http.StatusGatewayTimeout, true},
		{"too many requests", http.StatusTooManyRequests, true},
		{"ok", http.StatusOK, false},
		{"not found", http.StatusNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.statusCode}
			assert.Equal(t, tt.want, DefaultRetry(resp))
		})
	}
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	f := &fallback{}
	WithContext(ctx)(f)
	assert.Equal(t, ctx, f.ctx)
}

func TestWithHeader(t *testing.T) {
	header := http.Header{
		"X-Test": []string{"test"},
	}
	f := &fallback{}
	WithHeader(header)(f)
	assert.Equal(t, header, f.header)
}

func TestHandleURLs(t *testing.T) {
	urls := []string{"https://example.com", "https://example.org", "https://example.com"}
	urls1 := handleURLs(urls, SelectStrategyRoundRobin)
	assert.Equal(t, []string{"https://example.com", "https://example.org"}, urls1)
	urls2 := handleURLs(urls, SelectStrategyRandom)
	assert.Len(t, urls2, 2)
	assert.Contains(t, urls2, "https://example.com")
	assert.Contains(t, urls2, "https://example.org")
}

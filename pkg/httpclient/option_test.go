package httpclient

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	opt := WithContext(ctx)
	r := &request{}
	opt(r)
	assert.Equal(t, ctx, r.ctx)
}

func TestWithHeaders(t *testing.T) {
	headers := http.Header{
		"Test-Header": []string{"test-value"},
	}
	opt := WithHeader(headers)
	r := &request{}
	opt(r)
	assert.Equal(t, headers, r.header)
}

func TestWithRetry(t *testing.T) {
	tests := []struct {
		name         string
		retryMax     int
		retryWaitMin time.Duration
		retryWaitMax time.Duration
		want         *clientConfig
	}{
		{
			name:         "valid values",
			retryMax:     5,
			retryWaitMin: 2 * time.Second,
			retryWaitMax: 60 * time.Second,
			want: &clientConfig{
				retryMax:     5,
				retryWaitMin: 2 * time.Second,
				retryWaitMax: 60 * time.Second,
			},
		},
		{
			name:         "invalid values should use defaults",
			retryMax:     -1,
			retryWaitMin: 0,
			retryWaitMax: 0,
			want: &clientConfig{
				retryMax:     DefaultRetryMax,
				retryWaitMin: DefaultRetryMinWait,
				retryWaitMax: DefaultRetryMaxWait,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clientConfig{}
			opt := WithRetry(tt.retryMax, tt.retryWaitMin, tt.retryWaitMax)
			opt(c)
			assert.Equal(t, tt.want.retryMax, c.retryMax)
			assert.Equal(t, tt.want.retryWaitMin, c.retryWaitMin)
			assert.Equal(t, tt.want.retryWaitMax, c.retryWaitMax)
		})
	}
}

func TestWithDefaultHeader(t *testing.T) {
	header := http.Header{
		"Test-Header": []string{"test-value"},
	}
	opt := WithDefaultHeader(header)
	c := &clientConfig{}
	opt(c)
	assert.Equal(t, header, c.defaultHeader)
}

func TestWithLogger(t *testing.T) {
	logger := slog.Default()
	opt := WithLogger(logger)
	c := &clientConfig{}
	opt(c)
	assert.Equal(t, logger, c.logger)
}

func TestWithRequestIDKey(t *testing.T) {
	key := "Test-Request-ID"
	opt := WithRequestIDKey(key)
	c := &clientConfig{}
	opt(c)
	assert.Equal(t, key, c.requestIDKey)
}

func TestWithSkipTLSVerify(t *testing.T) {
	skip := true
	opt := WithSkipTLSVerify(skip)
	c := &clientConfig{}
	opt(c)
	assert.Equal(t, skip, c.skipTLSVerify)
}

func TestWithDump(t *testing.T) {
	ch := make(chan []byte)
	opt := WithDump(ch, true, true, true, true)
	c := &clientConfig{}
	opt(c)
	assert.NotNil(t, c.dumpChan)
	assert.True(t, c.dumpRequest)
	assert.True(t, c.dumpResponse)
	assert.True(t, c.dumpRequestBody)
	assert.True(t, c.dumpResponseBody)

	ch = nil
	opt = WithDump(ch, true, true, true, true)
	c = &clientConfig{}
	opt(c)
	assert.Nil(t, c.dumpChan)
	assert.False(t, c.dumpRequest)
	assert.False(t, c.dumpResponse)
	assert.False(t, c.dumpRequestBody)
	assert.False(t, c.dumpResponseBody)
}

func TestGenRequestID(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		key     string
		wantKey string
	}{
		{
			name:    "empty key should use default",
			req:     &http.Request{},
			key:     "",
			wantKey: DefaultRequestIDKey,
		},
		{
			name:    "custom key",
			req:     &http.Request{},
			key:     "Custom-Request-ID",
			wantKey: "Custom-Request-ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rid := GenRequestID(tt.req, tt.key)
			assert.NotEmpty(t, rid)
			assert.Equal(t, rid, tt.req.Header.Get(tt.wantKey))
		})
	}
}

func TestDefaultClientConfig(t *testing.T) {
	cfg := defaultClientConfig()

	assert.Equal(t, DefaultRetryMax, cfg.retryMax)
	assert.Equal(t, DefaultRetryMinWait, cfg.retryWaitMin)
	assert.Equal(t, DefaultRetryMaxWait, cfg.retryWaitMax)
	assert.Equal(t, DefaultRequestIDKey, cfg.requestIDKey)
	assert.NotNil(t, cfg.logger)
}

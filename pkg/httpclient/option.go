package httpclient

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/telepair/telepair/pkg/utils"
)

const (
	DefaultRetryMax     = 3
	DefaultRetryMinWait = 1 * time.Second
	DefaultRetryMaxWait = 30 * time.Second
	DefaultRequestIDKey = "X-Request-ID"
)

// RequestOption is a function that configures the request.
type RequestOption func(*request)

type request struct {
	ctx    context.Context
	header http.Header
}

// WithContext sets the context for the request.
func WithContext(ctx context.Context) RequestOption {
	return func(r *request) {
		if ctx != nil {
			r.ctx = ctx
		}
	}
}

// WithHeader sets the headers for the request.
func WithHeader(header http.Header) RequestOption {
	return func(r *request) {
		if header != nil {
			r.header = header
		}
	}
}

// Option is a function that configures the client.
type Option func(*clientConfig)

type clientConfig struct {
	retryMax         int
	retryWaitMin     time.Duration
	retryWaitMax     time.Duration
	defaultHeader    http.Header
	logger           *slog.Logger
	requestIDKey     string
	skipTLSVerify    bool
	dumpChan         <-chan []byte
	dumpRequest      bool
	dumpRequestBody  bool
	dumpResponse     bool
	dumpResponseBody bool
}

// WithRetry sets the retry options for the client.
func WithRetry(retryMax int, retryWaitMin, retryWaitMax time.Duration) Option {
	return func(c *clientConfig) {
		if retryMax < 0 {
			retryMax = DefaultRetryMax
		}
		if retryWaitMin <= 0 {
			retryWaitMin = DefaultRetryMinWait
		}
		if retryWaitMax <= 0 || retryWaitMax < retryWaitMin {
			retryWaitMax = DefaultRetryMaxWait
		}

		c.retryMax = retryMax
		c.retryWaitMin = retryWaitMin
		c.retryWaitMax = retryWaitMax
	}
}

// WithDefaultHeader sets the default headers for the client.
func WithDefaultHeader(header http.Header) Option {
	return func(c *clientConfig) {
		if header != nil {
			c.defaultHeader = header
		}
	}
}

// WithLogger sets the logger for the client.
func WithLogger(logger *slog.Logger) Option {
	return func(c *clientConfig) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithRequestIDKey sets the request ID key for the client.
func WithRequestIDKey(requestIDKey string) Option {
	return func(c *clientConfig) {
		if requestIDKey != "" {
			c.requestIDKey = requestIDKey
		}
	}
}

// WithSkipTLSVerify sets the skip TLS verify for the client.
func WithSkipTLSVerify(skipTLSVerify bool) Option {
	return func(c *clientConfig) {
		c.skipTLSVerify = skipTLSVerify
	}
}

// WithDump sets the dump options for the client.
func WithDump(ch <-chan []byte, request bool, response bool, requestBody bool, responseBody bool) Option {
	return func(c *clientConfig) {
		if ch == nil {
			slog.Warn("dump channel is nil, it will disable the dump")
			return
		}
		c.dumpChan = ch
		c.dumpRequest = request
		c.dumpResponse = response
		c.dumpRequestBody = requestBody
		c.dumpResponseBody = responseBody
	}
}

// GenRequestID generates the request ID and sets it to the request header.
func GenRequestID(req *http.Request, key string) string {
	if key == "" {
		key = DefaultRequestIDKey
	}
	if req.Header == nil {
		req.Header = http.Header{}
	}
	rid := req.Header.Get(key)
	if rid == "" {
		rid = utils.UUIDv7().String()
		req.Header.Set(key, rid)
	}
	return rid
}

// Add default values method
func defaultClientConfig() *clientConfig {
	return &clientConfig{
		retryMax:     DefaultRetryMax,
		retryWaitMin: DefaultRetryMinWait,
		retryWaitMax: DefaultRetryMaxWait,
		requestIDKey: DefaultRequestIDKey,
		logger:       slog.With("component", "httpclient"),
	}
}

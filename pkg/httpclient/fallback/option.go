package fallback

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"slices"

	"github.com/telepair/telepair/pkg/httpclient"
)

// Option is a option for the fallback
type Option func(*fallback)

type fallback struct {
	ctx      context.Context
	header   http.Header
	client   httpclient.Client
	selector SelectStrategy
	retry    RetryChecker
	logger   *slog.Logger
}

func defaultFallback() *fallback {
	return &fallback{
		client:   httpclient.NoRetryClient,
		selector: SelectStrategyRoundRobin,
		retry:    DefaultRetry,
		logger:   slog.With("component", "httpclient/fallback"),
	}
}

// SelectStrategy is a strategy for selecting the next url
type SelectStrategy string

const (
	SelectStrategyRoundRobin SelectStrategy = "round_robin"
	SelectStrategyRandom     SelectStrategy = "random"
)

// RetryChecker is a function that checks if the response should be retried
type RetryChecker func(resp *http.Response) bool

// DefaultRetry is a default retry checker that retries on internal server error
var DefaultRetry = func(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	case http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}

// WithSelectStrategy sets the select strategy for the fallback
func WithSelectStrategy(strategy SelectStrategy) Option {
	return func(f *fallback) {
		switch strategy {
		case SelectStrategyRoundRobin, SelectStrategyRandom:
			f.selector = strategy
		default:
			slog.Warn("invalid select strategy, use round_robin instead", "strategy", strategy)
			f.selector = SelectStrategyRoundRobin
		}
	}
}

// WithRetryClient sets the retry client for the fallback
func WithRetryClient(client httpclient.Client) Option {
	return func(f *fallback) {
		if client != nil {
			f.client = client
		} else {
			slog.Warn("client is nil, use default client instead")
			f.client = httpclient.NoRetryClient
		}
	}
}

// WithRetryChecker sets the retry checker for the fallback
func WithRetryChecker(checker RetryChecker) Option {
	return func(f *fallback) {
		if checker != nil {
			f.retry = checker
		} else {
			slog.Warn("retry checker is nil, use default checker instead")
			f.retry = DefaultRetry
		}
	}
}

// WithLogger sets the logger for the fallback
func WithLogger(logger *slog.Logger) Option {
	return func(f *fallback) {
		if logger != nil {
			f.logger = logger
		} else {
			slog.Warn("logger is nil, use default logger instead")
			f.logger = slog.With("component", "httpclient/fallback")
		}
	}
}

// WithContext sets the context for the request.
func WithContext(ctx context.Context) Option {
	return func(f *fallback) {
		if ctx != nil {
			f.ctx = ctx
		}
	}
}

// WithHeader sets the headers for the request.
func WithHeader(header http.Header) Option {
	return func(f *fallback) {
		if header != nil {
			f.header = header
		}
	}
}

// handleURLs handles the urls
func handleURLs(urls []string, selector SelectStrategy) []string {
	if len(urls) == 0 {
		return []string{}
	}
	if selector == SelectStrategyRandom {
		slices.Sort(urls)
		urls = slices.Compact(urls)
		rand.Shuffle(len(urls), func(i, j int) {
			urls[i], urls[j] = urls[j], urls[i]
		})
		return urls
	}

	m := make(map[string]struct{}, len(urls))
	result := make([]string, 0, len(m))
	for _, url := range urls {
		if _, ok := m[url]; !ok {
			m[url] = struct{}{}
			result = append(result, url)
		}
	}
	return result
}

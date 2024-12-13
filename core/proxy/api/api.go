package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/telepair/telepair/pkg/httpclient"
	"github.com/telepair/telepair/pkg/httpclient/fallback"
)

var (
	SuccessCodes   = []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}
	AllowedMethods = []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}
	DefaultTimeout = 10 * time.Second
)

// API is a struct for API
type API struct {
	Name    string            `yaml:"name" json:"name"`
	Method  string            `yaml:"method" json:"method"`
	URL     string            `yaml:"url" json:"url"`
	URLs    []string          `yaml:"urls" json:"urls"`
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body    string            `yaml:"body,omitempty" json:"body,omitempty"`
	Config  Config            `yaml:"config,omitempty" json:"config,omitempty"`
}

// Parse parses the API
func (c *API) Parse() error {
	c.Method = strings.ToUpper(strings.TrimSpace(c.Method))
	switch c.Method {
	case http.MethodOptions, http.MethodGet, http.MethodHead:
		c.Body = ""
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
	default:
		return fmt.Errorf("method %s is invalid", c.Method)
	}

	c.URL = strings.TrimSpace(c.URL)
	if c.URL == "" {
		if len(c.URLs) == 0 {
			return errors.New("url or urls is required")
		}
		if len(c.URLs) == 1 {
			c.URL = c.URLs[0]
		}
	} else {
		c.URLs = nil
	}

	if c.Config.Timeout == 0 {
		c.Config.Timeout = DefaultTimeout
	}

	return nil
}

// Do executes the API
func (c *API) Do(opts ...Option) (resp *http.Response, err error) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.ctx == nil {
		ctx, cancel := context.WithTimeout(context.Background(), c.Config.Timeout)
		defer cancel()
		cfg.ctx = ctx
	}

	if c.URL != "" {
		resp, err = c.do(cfg.ctx, cfg.client)
	} else {
		resp, err = c.doWithFallback(cfg.ctx, cfg.client)
	}
	if err != nil {
		return nil, err
	}

	if !c.IsSuccess(resp) {
		return resp, fmt.Errorf("status code %d is not in success codes", resp.StatusCode)
	}

	return resp, nil
}

func (c *API) do(ctx context.Context, client httpclient.Client) (*http.Response, error) {
	var req *http.Request
	req, err := http.NewRequestWithContext(ctx, c.Method, c.URL, bytes.NewBufferString(c.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	if client == nil {
		return httpclient.Do(req)
	}
	return client.Do(req)
}

func (c *API) doWithFallback(ctx context.Context, client httpclient.Client) (*http.Response, error) {
	if len(c.URLs) == 0 {
		return nil, errors.New("urls is required")
	}

	opts := []fallback.Option{fallback.WithContext(ctx)}
	if client != nil {
		opts = append(opts, fallback.WithRetryClient(client))
	}
	if c.Headers != nil {
		headers := http.Header{}
		for k, v := range c.Headers {
			headers.Add(k, v)
		}
		opts = append(opts, fallback.WithHeader(headers))
	}
	if c.Config.Fallback.Selector != "" {
		opts = append(opts, fallback.WithSelectStrategy(c.Config.Fallback.Selector))
	}
	if len(c.Config.Fallback.RetryCodes) > 0 {
		opts = append(opts, fallback.WithRetryChecker(func(resp *http.Response) bool {
			return slices.Contains(c.Config.Fallback.RetryCodes, resp.StatusCode)
		}))
	}

	return fallback.Do(c.Method, c.URLs, []byte(c.Body), opts...)
}

// IsSuccess checks if the response is successful
func (c *API) IsSuccess(resp *http.Response) bool {
	var flag bool
	if len(c.Config.Checker.SuccessCodes) == 0 {
		flag = slices.Contains(SuccessCodes, resp.StatusCode)
	} else {
		flag = slices.Contains(c.Config.Checker.SuccessCodes, resp.StatusCode)
	}
	if !flag {
		return false
	}

	for k, v := range c.Config.Checker.HeaderMatch {
		if resp.Header.Get(k) != v {
			return false
		}
	}

	// TODO: body match
	return true
}

// Config is a struct for API config
type Config struct {
	Checker  Checker       `yaml:"checker" json:"checker"`
	Fallback Fallback      `yaml:"fallback" json:"fallback"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout"`
}

// Checker is a struct for API checker
type Checker struct {
	SuccessCodes []int             `yaml:"success_codes" json:"success_codes"`
	HeaderMatch  map[string]string `yaml:"header_match" json:"header_match"`
}

// Fallback is a struct for API fallback
type Fallback struct {
	Selector   fallback.SelectStrategy `yaml:"selector" json:"selector"`
	RetryCodes []int                   `yaml:"retry_codes" json:"retry_codes"`
}

type Option func(*config)

type config struct {
	ctx    context.Context
	client httpclient.Client
}

func WithClient(client httpclient.Client) Option {
	return func(c *config) {
		c.client = client
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *config) {
		c.ctx = ctx
	}
}

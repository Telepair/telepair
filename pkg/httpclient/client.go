package httpclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/go-retryablehttp"
)

// Client is a http client.
type Client interface {
	Get(url string, opts ...RequestOption) (*http.Response, error)
	Head(url string, opts ...RequestOption) (*http.Response, error)
	Options(url string, opts ...RequestOption) (*http.Response, error)
	Post(url string, body []byte, opts ...RequestOption) (*http.Response, error)
	Put(url string, body []byte, opts ...RequestOption) (*http.Response, error)
	Patch(url string, body []byte, opts ...RequestOption) (*http.Response, error)
	Delete(url string, opts ...RequestOption) (*http.Response, error)
	Do(req *http.Request, opts ...RequestOption) (*http.Response, error)
}

type client struct {
	cfg      *clientConfig
	c        *retryablehttp.Client
	logger   *slog.Logger
	dumpChan chan []byte
}

// New creates a new http client
func New(opts ...Option) Client {
	cfg := defaultClientConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	hc := &client{
		cfg:    cfg,
		c:      retryablehttp.NewClient(),
		logger: cfg.logger,
	}
	hc.setupClient()
	hc.setupDump()

	return hc
}

func (c *client) setupClient() {
	if c.c == nil {
		c.c = retryablehttp.NewClient()
	}

	if c.cfg.retryWaitMin > 0 {
		c.c.RetryWaitMin = c.cfg.retryWaitMin
	}
	if c.cfg.retryWaitMax > 0 {
		c.c.RetryWaitMax = c.cfg.retryWaitMax
	}
	if c.cfg.retryMax >= 0 {
		c.c.RetryMax = c.cfg.retryMax
	}
	if c.cfg.skipTLSVerify {
		c.c.HTTPClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}
	if c.logger != nil {
		c.c.Logger = c.logger
	}
}

func (c *client) setupDump() {
	if c.cfg.dumpRequest && c.dumpChan != nil {
		c.c.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, retryNumber int) {
			rid := GenRequestID(req, c.cfg.requestIDKey)
			url := req.URL.String()
			method := req.Method
			c.logger.Debug("request", "retry", retryNumber, "request_id", rid, "url", url, "method", method)

			data, err := httputil.DumpRequest(req, c.cfg.dumpRequestBody)
			if err != nil {
				slog.Warn("dump request", "error", err)
			} else {
				c.dumpChan <- data
			}
		}
	} else {
		c.c.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, retryNumber int) {
			rid := GenRequestID(req, c.cfg.requestIDKey)
			url := req.URL.String()
			method := req.Method
			c.logger.Debug("request", "retry", retryNumber, "request_id", rid, "url", url, "method", method)
		}
	}

	if c.cfg.dumpResponse && c.dumpChan != nil {
		c.c.ResponseLogHook = func(_ retryablehttp.Logger, resp *http.Response) {
			rid := GenRequestID(resp.Request, c.cfg.requestIDKey)
			url := resp.Request.URL.String()
			method := resp.Request.Method
			if resp.StatusCode >= http.StatusBadRequest {
				c.logger.Warn("response", "request_id", rid, "url", url, "method", method, "status", resp.Status)
			} else {
				c.logger.Debug("response", "request_id", rid, "url", url, "method", method, "status", resp.Status)
			}

			data, err := httputil.DumpResponse(resp, c.cfg.dumpResponseBody)
			if err != nil {
				c.logger.Warn("dump response", "request_id", rid, "url", url, "method", method, "error", err)
			} else {
				c.dumpChan <- data
			}
		}
	} else {
		c.c.ResponseLogHook = func(_ retryablehttp.Logger, resp *http.Response) {
			rid := GenRequestID(resp.Request, c.cfg.requestIDKey)
			url := resp.Request.URL.String()
			method := resp.Request.Method
			if resp.StatusCode >= http.StatusBadRequest {
				c.logger.Warn("response", "request_id", rid, "url", url, "method", method, "status", resp.Status)
			} else {
				c.logger.Debug("response", "request_id", rid, "url", url, "method", method, "status", resp.Status)
			}
		}
	}
}

// Get sends a GET request and returns a response.
func (c *client) Get(url string, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Head sends a HEAD request and returns a response.
func (c *client) Head(url string, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Options sends a OPTIONS request and returns a response.
func (c *client) Options(url string, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodOptions, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Post sends a POST request and returns a response.
func (c *client) Post(url string, body []byte, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Put sends a PUT request and returns a response.
func (c *client) Put(url string, body []byte, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Patch sends a PATCH request and returns a response.
func (c *client) Patch(url string, body []byte, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Delete sends a DELETE request and returns a response.
func (c *client) Delete(url string, opts ...RequestOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, opts...)
}

// Do sends a request and returns a response.
func (c *client) Do(req *http.Request, opts ...RequestOption) (*http.Response, error) {
	r := &request{}
	for _, opt := range opts {
		opt(r)
	}

	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}
	for k, v := range c.cfg.defaultHeader {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	for k, v := range r.header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	rreq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}
	return c.c.Do(rreq)
}

// ParseResponse is a helper function that reads the response body and returns the bytes
func ParseResponse(resp *http.Response) (mediaType string, body []byte, e error) {
	if resp == nil {
		e = errors.New("response is nil")
		return
	}
	mediaType, _, e = mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if e != nil {
		return
	}
	body, e = io.ReadAll(resp.Body)
	if e != nil {
		return
	}
	_ = resp.Body.Close()
	return
}

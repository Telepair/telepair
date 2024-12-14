package fallback

import (
	"errors"
	"net/http"

	"github.com/telepair/telepair/pkg/httpclient"
)

// Get gets the response from the urls, and returns the first successful response
func Get(urls []string, opts ...Option) (*http.Response, error) {
	return Do(http.MethodGet, urls, nil, opts...)
}

// Head is the fallback head method
func Head(urls []string, opts ...Option) (*http.Response, error) {
	return Do(http.MethodHead, urls, nil, opts...)
}

// Options is the fallback options method
func Options(urls []string, opts ...Option) (*http.Response, error) {
	return Do(http.MethodOptions, urls, nil, opts...)
}

// Post is the fallback post method
func Post(urls []string, body []byte, opts ...Option) (*http.Response, error) {
	return Do(http.MethodPost, urls, body, opts...)
}

// Put is the fallback put method
func Put(urls []string, body []byte, opts ...Option) (*http.Response, error) {
	return Do(http.MethodPut, urls, body, opts...)
}

// Patch is the fallback patch method
func Patch(urls []string, body []byte, opts ...Option) (*http.Response, error) {
	return Do(http.MethodPatch, urls, body, opts...)
}

// Delete is the fallback delete method
func Delete(urls []string, opts ...Option) (*http.Response, error) {
	return Do(http.MethodDelete, urls, nil, opts...)
}

// Do is the fallback do method
func Do(method string, urls []string, body []byte, opts ...Option) (resp *http.Response, err error) {
	f := defaultFallback()
	for _, opt := range opts {
		opt(f)
	}

	copts := make([]httpclient.RequestOption, 0, 2)
	if f.ctx != nil {
		copts = append(copts, httpclient.WithContext(f.ctx))
	}
	if f.header != nil {
		copts = append(copts, httpclient.WithHeader(f.header))
	}

	for _, url := range handleURLs(urls, f.selector) {
		switch method {
		case http.MethodGet:
			resp, err = f.client.Get(url, copts...)
		case http.MethodHead:
			resp, err = f.client.Head(url, copts...)
		case http.MethodOptions:
			resp, err = f.client.Options(url, copts...)
		case http.MethodPost:
			resp, err = f.client.Post(url, body, copts...)
		case http.MethodPut:
			resp, err = f.client.Put(url, body, copts...)
		case http.MethodPatch:
			resp, err = f.client.Patch(url, body, copts...)
		case http.MethodDelete:
			resp, err = f.client.Delete(url, copts...)
		default:
			f.logger.Error("unsupported method", "method", method)
			return nil, errors.New("unsupported method: " + method)
		}

		if err != nil {
			f.logger.Error("request url failed", "url", url, "error", err)
			continue
		}

		if f.retry(resp) {
			f.logger.Error("get url failed", "url", url, "status", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, errors.New("all urls failed")
}

package httpclient

import (
	"net/http"
	"time"
)

var (
	// DefaultClient is the default http client
	DefaultClient = New()
	NoRetryClient = New(WithRetry(0, 1*time.Second, 1*time.Second))
)

// Get is the default get method
func Get(url string, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Get(url, opts...)
}

// Head is the default head method
func Head(url string, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Head(url, opts...)
}

// Options is the default options method
func Options(url string, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Options(url, opts...)
}

// Post is the default post method
func Post(url string, body []byte, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Post(url, body, opts...)
}

// Put is the default put method
func Put(url string, body []byte, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Put(url, body, opts...)
}

// Patch is the default patch method
func Patch(url string, body []byte, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Patch(url, body, opts...)
}

// Delete is the default delete method
func Delete(url string, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Delete(url, opts...)
}

// Do is the default do method
func Do(req *http.Request, opts ...RequestOption) (*http.Response, error) {
	return DefaultClient.Do(req, opts...)
}

package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAPI_Do(t *testing.T) {
	api := API{}
	err := RegisterAPI(api)
	assert.Error(t, err)

	api.Name = "test-api"
	err = RegisterAPI(api)
	assert.Error(t, err)

	api.Method = "GET"
	err = RegisterAPI(api)
	assert.Error(t, err)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	}))
	defer srv.Close()

	api.URL = srv.URL
	err = RegisterAPI(api)
	assert.NoError(t, err)

	assert.Error(t, RegisterAPI(api)) // already registered

	resp, err := Do("test-api")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(body))
}

func TestRegisterAPITemplate_Do(t *testing.T) {
	at := Template{
		Name:   "test-api-template",
		Method: true,
		URL:    true,
		Headers: map[string]bool{
			"Content-Type": true,
		},
		API: API{
			Method: "{{ method }}",
			URL:    "{{ url }}",
			Headers: map[string]string{
				"Content-Type": "{{ content-type }}",
			},
		},
		Defaults: map[string]string{
			"method":       "GET",
			"content-type": "application/json",
		},
	}
	err := RegisterTemplate(at)
	assert.NoError(t, err)

	resp, err := DoTemplate("test-api-template", nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test"))
	}))
	defer srv.Close()

	resp, err = DoTemplate("test-api-template", map[string]string{
		"method":       "GET",
		"url":          srv.URL,
		"content-type": "application/json",
	})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test", string(body))
}

func TestRegisterAPIData(t *testing.T) {
	err := RegisterAPIData("yaml", []byte(`
- name: test-api1
  method: GET
  url: https://api.example.com
`))
	assert.NoError(t, err)
}

func TestRegisterAPITemplateData(t *testing.T) {
	err := RegisterAPITemplateData("yaml", []byte(`
- name: test-api-template1
  api:
    method: "{{ method }}"
    url: "{{ url }}"
    headers:
      "Content-Type": "{{ content-type }}"
  method: true
  url: true
  headers:
    Content-Type: true
  defaults:
    method: GET
    content-type: application/json
    url: https://api.example.com
`))
	assert.NoError(t, err)
}

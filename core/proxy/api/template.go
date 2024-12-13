package api

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/telepair/telepair/pkg/utils"
)

// Template is a template for API
type Template struct {
	Name     string            `yaml:"name" json:"name"`
	API      API               `yaml:"api" json:"api"`
	Method   bool              `yaml:"method"`
	URL      bool              `yaml:"url"`
	Headers  map[string]bool   `yaml:"headers"`
	Body     bool              `yaml:"body"`
	Defaults map[string]string `yaml:"defaults"`
}

// Parse parses the API template
func (t *Template) Parse() error {
	if t.Name == "" {
		return errors.New("name is required")
	}

	if !t.Method {
		t.API.Method = strings.ToUpper(strings.TrimSpace(t.API.Method))
		if !slices.Contains(AllowedMethods, t.API.Method) {
			return fmt.Errorf("method %s is invalid", t.API.Method)
		}
	}

	if !t.URL && len(t.API.URLs) == 0 && t.API.URL == "" {
		return errors.New("api url is required")
	}

	if t.API.Headers == nil {
		t.API.Headers = make(map[string]string)
	}
	for key, val := range t.Headers {
		if val && t.API.Headers[key] == "" {
			return fmt.Errorf("header %s is required", key)
		}
	}

	if t.Body && t.API.Body == "" {
		return errors.New("body is required")
	}

	t.API.Name = ""

	return nil
}

// Render renders the API template
func (t *Template) Render(vars map[string]string) (api API, err error) {
	api = t.API
	if t.Method {
		api.Method, err = utils.SimpleRender(t.API.Method, vars, t.Defaults)
		if err != nil {
			return
		}
	}
	if t.URL {
		api.URL, err = utils.SimpleRender(t.API.URL, vars, t.Defaults)
		if err != nil {
			return
		}
	}
	for i, url := range api.URLs {
		api.URLs[i], err = utils.SimpleRender(url, vars, t.Defaults)
		if err != nil {
			return
		}
	}
	for key, val := range api.Headers {
		if t.Headers[key] {
			api.Headers[key], err = utils.SimpleRender(val, vars, t.Defaults)
			if err != nil {
				return
			}
		}
	}
	if t.Body {
		api.Body, err = utils.SimpleRender(t.API.Body, vars, t.Defaults)
		if err != nil {
			return
		}
	}
	api.Name = fmt.Sprintf("%s::%s", t.Name, utils.UUIDv7().String())
	if err := api.Parse(); err != nil {
		return api, err
	}
	return api, nil
}

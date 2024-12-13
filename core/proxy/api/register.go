package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/telepair/telepair/pkg/cache"
)

var (
	apiCache         = cache.NewMemory("api-proxy")
	apiTemplateCache = cache.NewMemory("api-proxy-template")
)

// RegisterAPI registers the API
func RegisterAPI(api API) error {
	if err := api.Parse(); err != nil {
		return err
	}
	ok, err := apiCache.Exists(context.Background(), api.Name)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("api %s already exists", api.Name)
	}
	return apiCache.Set(context.Background(), api.Name, api)
}

// Do returns the response of the API
func Do(name string) (*http.Response, error) {
	val, err := apiCache.Get(context.Background(), name)
	if err != nil {
		return nil, err
	}
	api := val.(API)
	return api.Do()
}

// RegisterTemplate registers the API template
func RegisterTemplate(template Template) error {
	if err := template.Parse(); err != nil {
		return err
	}
	ok, err := apiTemplateCache.Exists(context.Background(), template.Name)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("api template %s already exists", template.Name)
	}
	return apiTemplateCache.Set(context.Background(), template.Name, template)
}

// DoTemplate renders the API template and returns the response
func DoTemplate(name string, vars map[string]string) (*http.Response, error) {
	val, err := apiTemplateCache.Get(context.Background(), name)
	if err != nil {
		return nil, err
	}
	template := val.(Template)
	api, err := template.Render(vars)
	if err != nil {
		return nil, err
	}
	return api.Do()
}

// RegisterAPIData registers the API data
func RegisterAPIData(dataType string, data []byte) error {
	dataType = strings.ToLower(strings.TrimSpace(dataType))
	var apis []API
	switch dataType {
	case "yaml", "yml":
		err := yaml.Unmarshal(data, &apis)
		if err != nil {
			return fmt.Errorf("failed to unmarshal yaml: %w", err)
		}
	case "json":
		err := json.Unmarshal(data, &apis)
		if err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type: %s", dataType)
	}
	for _, api := range apis {
		if err := RegisterAPI(api); err != nil {
			return fmt.Errorf("failed to register api: %w", err)
		}
	}
	return nil
}

// RegisterAPITemplateData registers the API template data
func RegisterAPITemplateData(dataType string, data []byte) error {
	dataType = strings.ToLower(strings.TrimSpace(dataType))
	var templates []Template
	switch dataType {
	case "yaml", "yml":
		err := yaml.Unmarshal(data, &templates)
		if err != nil {
			return fmt.Errorf("failed to unmarshal yaml: %w", err)
		}
	case "json":
		err := json.Unmarshal(data, &templates)
		if err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported data type: %s", dataType)
	}
	for _, template := range templates {
		if err := RegisterTemplate(template); err != nil {
			return fmt.Errorf("failed to register api template: %w", err)
		}
	}
	return nil
}

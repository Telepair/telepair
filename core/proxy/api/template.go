package api

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/telepair/telepair/pkg/utils"
)

// Template represents an API template that can be rendered with variables
// Example:
//
//	template := Template{
//	    Name: "my-api",
//	    API: API{
//	        Method: "GET",
//	        URL: "https://api.example.com/users/{{.userId}}",
//	    },
//	    TemplateField: TemplateField{
//	        URL: true,
//	    },
//	    Vars: []VarRequired{
//	        {Name: "userId", CanEmpty: false},
//	    },
//	}
type Template struct {
	Name          string        `yaml:"name" json:"name"`
	API           API           `yaml:"api" json:"api"`
	TemplateField TemplateField `yaml:"template_field" json:"template_field"`
	Vars          []VarRequired `yaml:"vars" json:"vars"`
}

// TemplateField is a field for API template
type TemplateField struct {
	Method  bool            `yaml:"method"`
	URL     bool            `yaml:"url"`
	Headers map[string]bool `yaml:"headers"`
	Body    bool            `yaml:"body"`
}

// VarRequired is a variable for API template
type VarRequired struct {
	Name     string   `yaml:"name"`
	Options  []string `yaml:"options"`
	Default  string   `yaml:"default"`
	CanEmpty bool     `yaml:"can_empty"`
}

// Validate validates the variable
func (v *VarRequired) Validate() error {
	if v.Name == "" {
		return errors.New("variable name is required")
	}
	if v.Default != "" && len(v.Options) > 0 && !slices.Contains(v.Options, v.Default) {
		return fmt.Errorf("variable %s default value %s is not in options", v.Name, v.Default)
	}
	return nil
}

// Parse parses the API template
func (t *Template) Parse() error {
	if t.Name == "" {
		return errors.New("name is required")
	}

	for _, v := range t.Vars {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	if !t.TemplateField.Method {
		t.API.Method = strings.ToUpper(strings.TrimSpace(t.API.Method))
		if !slices.Contains(AllowedMethods, t.API.Method) {
			return fmt.Errorf("method template %s is invalid", t.API.Method)
		}
	}

	if !t.TemplateField.URL && t.API.URL == "" && len(t.API.URLs) == 0 {
		return errors.New("either URL or URLs template must be specified")
	}

	if t.API.Headers == nil {
		t.API.Headers = make(map[string]string)
	}
	for key, val := range t.TemplateField.Headers {
		if val && t.API.Headers[key] == "" {
			return fmt.Errorf("header %s template is required", key)
		}
	}

	if t.TemplateField.Body && t.API.Body == "" {
		return errors.New("body template is required")
	}

	t.API.Name = ""

	return nil
}

// Render renders the API template
func (t *Template) Render(vars map[string]string) (api API, err error) {
	vars, err = t.MergeVars(vars)
	if err != nil {
		return
	}
	api = t.API
	if t.TemplateField.Method {
		api.Method, err = utils.SimpleRender(t.API.Method, vars)
		if err != nil {
			return
		}
	}

	if err = t.renderURL(&api, vars); err != nil {
		return
	}

	for key, val := range api.Headers {
		if t.TemplateField.Headers[key] {
			api.Headers[key], err = utils.SimpleRender(val, vars)
			if err != nil {
				return
			}
		}
	}
	if t.TemplateField.Body {
		api.Body, err = utils.SimpleRender(t.API.Body, vars)
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

// MergeVars merges the variables
func (t *Template) MergeVars(vars map[string]string) (map[string]string, error) {
	merged := make(map[string]string)
	for _, v := range t.Vars {
		val, ok := vars[v.Name]
		if !ok {
			if v.Default != "" {
				val = v.Default
			} else if v.CanEmpty {
				val = ""
			} else {
				return nil, fmt.Errorf("variable %s is required", v.Name)
			}
		}

		if !v.CanEmpty && val == "" {
			return nil, fmt.Errorf("variable %s is required", v.Name)
		}
		if len(v.Options) > 0 && !slices.Contains(v.Options, val) {
			return nil, fmt.Errorf("variable %s is invalid, not in options", v.Name)
		}

		merged[v.Name] = val
	}
	return merged, nil
}

func (t *Template) renderURL(api *API, vars map[string]string) (err error) {
	if !t.TemplateField.URL {
		return nil
	}
	if t.API.URL != "" {
		api.URL, err = utils.SimpleRender(t.API.URL, vars)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to render URL template %s", t.API.URL))
		}
	}
	for i, url := range api.URLs {
		api.URLs[i], err = utils.SimpleRender(url, vars)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to render URL template %s", url))
		}
	}
	return nil
}

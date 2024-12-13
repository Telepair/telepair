package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleRegisterTemplate() {
	data := []byte(`
- name: city-weather
  api:
    method: GET
    url: https://wttr.in/{{ city }}
    config:
      checker:
        success_codes: [200]
      timeout: 10s
  url: true
  defaults:
    city: "Beijing"
`)
	if err := RegisterAPITemplateData("yaml", data); err != nil {
		panic(err)
	}
	resp, err := DoTemplate("city-weather", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)

	resp, err = DoTemplate("city-weather", map[string]string{"city": "London"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
	// Output:
	// 200
	// 200
}

func TestAPITemplate_Parse(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    Template
		wantErr bool
	}{
		{
			name: "valid template",
			tmpl: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "https://api.example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			tmpl: Template{
				API: API{
					Method: "GET",
					URL:    "https://api.example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			tmpl: Template{
				Name: "test",
				API: API{
					Method: "INVALID",
					URL:    "https://api.example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "empty url",
			tmpl: Template{
				Name: "test",
				API: API{
					Method: "GET",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tmpl.Parse()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAPITemplate_Render(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    Template
		vars    map[string]string
		want    API
		wantErr bool
	}{
		{
			name: "render template without variables",
			tmpl: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "https://api.example.com",
				},
			},
			vars: map[string]string{},
			want: API{
				Method: "GET",
				URL:    "https://api.example.com",
			},
			wantErr: false,
		},
		{
			name: "render template with variables",
			tmpl: Template{
				Name:   "test",
				Method: true,
				URL:    true,
				Body:   true,
				API: API{
					Method: "{{method}}",
					URL:    "https://api.example.com/{{path}}",
					Body:   "{{body}}",
				},
				Defaults: map[string]string{
					"method": "POST",
					"path":   "test",
					"body":   "default-body",
				},
			},
			vars: map[string]string{
				"path": "custom",
			},
			want: API{
				Method: "POST",
				URL:    "https://api.example.com/custom",
				Body:   "default-body",
			},
			wantErr: false,
		},
		{
			name: "render template with headers",
			tmpl: Template{
				Name: "test",
				Headers: map[string]bool{
					"Authorization": true,
				},
				API: API{
					Method: "GET",
					URL:    "https://api.example.com",
					Headers: map[string]string{
						"Authorization": "Bearer {{token}}",
					},
				},
				Defaults: map[string]string{
					"token": "default-token",
				},
			},
			vars: map[string]string{},
			want: API{
				Method: "GET",
				URL:    "https://api.example.com",
				Headers: map[string]string{
					"Authorization": "Bearer default-token",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tmpl.Render(tt.vars)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Method, got.Method)
			assert.Equal(t, tt.want.URL, got.URL)
			assert.Equal(t, tt.want.Headers, got.Headers)
			assert.Equal(t, tt.want.Body, got.Body)
		})
	}
}

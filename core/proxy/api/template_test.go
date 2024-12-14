package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVarRequired_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       VarRequired
		wantErr bool
	}{
		{
			name: "valid var",
			v: VarRequired{
				Name:     "test",
				Options:  []string{"a", "b"},
				Default:  "a",
				CanEmpty: false,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			v: VarRequired{
				Name:     "",
				Options:  []string{"a", "b"},
				Default:  "a",
				CanEmpty: false,
			},
			wantErr: true,
		},
		{
			name: "default not in options",
			v: VarRequired{
				Name:     "test",
				Options:  []string{"a", "b"},
				Default:  "c",
				CanEmpty: false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplate_Parse(t *testing.T) {
	tests := []struct {
		name    string
		t       Template
		wantErr bool
	}{
		{
			name: "valid template",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "http://example.com",
				},
				TemplateField: TemplateField{
					Method: false,
					URL:    false,
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			t: Template{
				Name: "",
				API: API{
					Method: "GET",
					URL:    "http://example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			t: Template{
				Name: "test",
				API: API{
					Method: "INVALID",
					URL:    "http://example.com",
				},
				TemplateField: TemplateField{
					Method: false,
					URL:    false,
				},
			},
			wantErr: true,
		},
		{
			name: "missing URL and URLs",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
				},
				TemplateField: TemplateField{
					Method: false,
					URL:    false,
				},
			},
			wantErr: true,
		},
		{
			name: "template header required but empty",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "http://example.com",
					Headers: map[string]string{
						"Authorization": "",
					},
				},
				TemplateField: TemplateField{
					Headers: map[string]bool{
						"Authorization": true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "template body required but empty",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "http://example.com",
					Body:   "",
				},
				TemplateField: TemplateField{
					Body: true,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.t.Parse()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTemplate_Render(t *testing.T) {
	tests := []struct {
		name    string
		t       Template
		vars    map[string]string
		wantErr bool
	}{
		{
			name: "valid render",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "http://example.com/{{path}}",
				},
				TemplateField: TemplateField{
					Method: false,
					URL:    true,
				},
				Vars: []VarRequired{
					{
						Name:     "path",
						Default:  "test",
						CanEmpty: false,
					},
				},
			},
			vars: map[string]string{
				"path": "api",
			},
			wantErr: false,
		},
		{
			name: "missing required var",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "http://example.com/{{path}}",
				},
				TemplateField: TemplateField{
					Method: false,
					URL:    true,
				},
				Vars: []VarRequired{
					{
						Name:     "path",
						CanEmpty: false,
					},
				},
			},
			vars:    map[string]string{},
			wantErr: true,
		},
		{
			name: "render with template method",
			t: Template{
				Name: "test",
				API: API{
					Method: "{{method}}",
					URL:    "http://example.com",
				},
				TemplateField: TemplateField{
					Method: true,
				},
				Vars: []VarRequired{
					{
						Name:     "method",
						Default:  "GET",
						CanEmpty: false,
					},
				},
			},
			vars:    map[string]string{},
			wantErr: false,
		},
		{
			name: "render with template headers",
			t: Template{
				Name: "test",
				API: API{
					Method: "GET",
					URL:    "http://example.com",
					Headers: map[string]string{
						"Authorization": "Bearer {{token}}",
					},
				},
				TemplateField: TemplateField{
					Headers: map[string]bool{
						"Authorization": true,
					},
				},
				Vars: []VarRequired{
					{
						Name:     "token",
						Default:  "test-token",
						CanEmpty: false,
					},
				},
			},
			vars:    map[string]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api, err := tt.t.Render(tt.vars)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, api.Name)
			}
		})
	}
}

func TestTemplate_MergeVars(t *testing.T) {
	tests := []struct {
		name    string
		t       Template
		vars    map[string]string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "merge with defaults",
			t: Template{
				Vars: []VarRequired{
					{
						Name:     "var1",
						Default:  "default1",
						CanEmpty: false,
					},
					{
						Name:     "var2",
						Default:  "",
						CanEmpty: true,
					},
				},
			},
			vars: map[string]string{
				"var1": "value1",
			},
			want: map[string]string{
				"var1": "value1",
				"var2": "",
			},
			wantErr: false,
		},
		{
			name: "invalid option",
			t: Template{
				Vars: []VarRequired{
					{
						Name:    "var1",
						Options: []string{"opt1", "opt2"},
					},
				},
			},
			vars: map[string]string{
				"var1": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged, err := tt.t.MergeVars(tt.vars)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, merged)
			}
		})
	}
}

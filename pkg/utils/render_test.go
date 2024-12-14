package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleSimpleRender() {
	jsonStr := `{"name": "{{ name }}", "age": {{ age }} }`
	vars := map[string]string{"name": "John", "age": "20"}
	got, err := SimpleRender(jsonStr, vars)
	fmt.Println(got, err)
	// Output:
	// {"name": "John", "age": 20 } <nil>
}

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		vars    map[string]string
		want    string
		wantErr string
	}{
		{
			name: "basic substitution",
			tmpl: "Hello {{ name }}!",
			vars: map[string]string{"name": "World"},
			want: "Hello World!",
		},
		{
			name: "multiple variables",
			tmpl: "{{ greeting }} {{ name }}!",
			vars: map[string]string{
				"greeting": "Hello",
				"name":     "World",
			},
			want: "Hello World!",
		},
		{
			name:    "missing variable",
			tmpl:    "Hello {{ name }}!",
			vars:    map[string]string{},
			want:    "Hello !",
			wantErr: "not found: [name]",
		},
		{
			name:    "multiple missing variables",
			tmpl:    "{{ greeting }} {{ name }}!",
			vars:    map[string]string{},
			want:    " !",
			wantErr: "not found: [greeting name]",
		},
		{
			name: "whitespace in variable",
			tmpl: "Hello {{   name   }}!",
			vars: map[string]string{"name": "World"},
			want: "Hello World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SimpleRender(tt.tmpl, tt.vars)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRender_JSON(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		vars    map[string]string
		want    string
		wantErr string
	}{
		{
			name: "basic json template",
			tmpl: `{"name": "{{ name }}", "age": {{ age }} }`,
			vars: map[string]string{
				"name": "John",
				"age":  "20",
			},
			want: `{"name": "John", "age": 20 }`,
		},
		{
			name: "missing json values",
			tmpl: `{"name": "{{ name }}", "age": {{ age }} }`,
			vars: map[string]string{
				"name": "John",
			},
			want:    `{"name": "John", "age":  }`,
			wantErr: "not found: [age]",
		},
		{
			name: "bool value",
			tmpl: `{"name": "{{ name }}", "status": {{ status }} }`,
			vars: map[string]string{
				"name":   "John",
				"status": "true",
			},
			want: `{"name": "John", "status": true }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SimpleRender(tt.tmpl, tt.vars)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

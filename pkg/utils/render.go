package utils

import (
	"fmt"
	"regexp"
	"slices"
)

// VariablePattern is the regex pattern for matching variables in the template
var VariablePattern = regexp.MustCompile(`{{\s*([a-zA-Z0-9_-]+)\s*}}`)

// SimpleRender renders the template with the given variables and defaults
func SimpleRender(tmpl string, vars map[string]string, defaults map[string]string) (string, error) {
	result := tmpl
	notFound := []string{}
	result = VariablePattern.ReplaceAllStringFunc(result, func(match string) string {
		varName := VariablePattern.FindStringSubmatch(match)[1]
		if value, exists := vars[varName]; exists {
			return value
		}
		if _, exists := defaults[varName]; !exists {
			notFound = append(notFound, varName)
		}
		return defaults[varName]
	})
	if len(notFound) > 0 {
		slices.Sort(notFound)
		notFound = slices.Compact(notFound)
		return result, fmt.Errorf("not found: %v", notFound)
	}
	return result, nil
}

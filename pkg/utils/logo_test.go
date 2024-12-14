package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleStdoutBanner() {
	StdoutBanner("Test")
	// Output is not deterministic, so we don't specify expected output
}

func TestStdoutBanner(t *testing.T) {
	assert.NotPanics(t, func() {
		StdoutBanner("Test")
	})
}

func TestRenderBanner(t *testing.T) {
	// Test case
	text := "Test"
	result := RenderBanner(text)

	// Verify the result is not empty
	assert.NotEmpty(t, result, "RenderBanner returned empty string")

	// Verify the result contains ASCII art characters
	assert.Contains(t, result, "_", "RenderBanner output doesn't appear to contain ASCII art")
	assert.Contains(t, result, "|", "RenderBanner output doesn't appear to contain ASCII art")
}

package logger

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleLogger() {
	Init(Config{
		Level:  "debug",
		Format: "json",
	})
	slog.Debug("test", "key", "value")
	slog.Info("test", "key", "value")
	slog.Warn("test", "key", "value")
	slog.Error("test", "key", "value", "error", errors.New("test"))
	// Output is not deterministic, so we don't specify expected output
}

func TestConfig_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    Config
		expected Config
	}{
		{
			name:  "empty config",
			input: Config{},
			expected: Config{
				Level:  "debug",
				Format: "text",
				Rotate: RotateConfig{
					MaxSize:    defaultMaxSize,
					MaxAge:     defaultMaxAge,
					MaxBackups: defaultMaxBackups,
				},
			},
		},
		{
			name: "custom config",
			input: Config{
				Level:     "info",
				Format:    "json",
				AddSource: true,
				File:      "test.log",
			},
			expected: Config{
				Level:     "info",
				Format:    "json",
				AddSource: true,
				File:      "test.log",
				Rotate: RotateConfig{
					MaxSize:    defaultMaxSize,
					MaxAge:     defaultMaxAge,
					MaxBackups: defaultMaxBackups,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.parse()
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "stdout logger",
			config: Config{
				Level:  "debug",
				Format: "text",
			},
		},
		{
			name: "file logger",
			config: Config{
				Level:  "info",
				Format: "json",
				File:   filepath.Join(t.TempDir(), "test.log"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				Init(tt.config)
			})
			// Verify logger is working by writing a test message
			slog.Info("test message")
		})
	}
}

func TestRotateConfig_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    RotateConfig
		expected RotateConfig
	}{
		{
			name:  "empty config",
			input: RotateConfig{},
			expected: RotateConfig{
				MaxSize:    defaultMaxSize,
				MaxAge:     defaultMaxAge,
				MaxBackups: defaultMaxBackups,
			},
		},
		{
			name: "custom config",
			input: RotateConfig{
				MaxSize:    200,
				MaxAge:     14,
				MaxBackups: 10,
				LocalTime:  true,
				Compress:   true,
			},
			expected: RotateConfig{
				MaxSize:    200,
				MaxAge:     14,
				MaxBackups: 10,
				LocalTime:  true,
				Compress:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Parse()
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func TestNewRotate(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	config := RotateConfig{
		MaxSize:    1,
		MaxAge:     1,
		MaxBackups: 1,
		LocalTime:  true,
		Compress:   true,
	}

	writer := NewRotate(filename, config)
	defer writer.Close()

	// Test that we can write to the log file
	_, err := writer.Write([]byte("test log message\n"))
	assert.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(filename)
	assert.NoError(t, err)
}

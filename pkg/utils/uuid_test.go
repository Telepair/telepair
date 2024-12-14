package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleUUIDv7() {
	uuid := UUIDv7()
	fmt.Println(uuid.String())
}

func ExampleUUIDv4() {
	uuid := UUIDv4()
	fmt.Println(uuid.String())
}

func TestNewUUID(t *testing.T) {
	uuid := NewUUID()
	assert.NotEmpty(t, uuid.String())
	assert.NotEmpty(t, uuid.B58())
}

func TestUUID_Time(t *testing.T) {
	uuid := NewUUID()
	now := time.Now()
	uuidTime := uuid.Time()

	// UUID timestamp should be within 1 second of current time
	assert.WithinDuration(t, now, uuidTime, time.Second)
}

func TestUUID_String(t *testing.T) {
	uuid := NewUUID()
	str := uuid.String()

	// Test parsing back
	var parsed UUID
	err := parsed.FromString(str)
	assert.NoError(t, err)
	assert.Equal(t, uuid, parsed)
}

func TestUUID_B58(t *testing.T) {
	uuid := NewUUID()
	b58 := uuid.B58()

	// Test parsing back
	var parsed UUID
	err := parsed.FromB58(b58)
	assert.NoError(t, err)
	assert.Equal(t, uuid, parsed)
}

func TestUUID_FromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			input:   "123e4567-e89b-12d3-a456-426614174000",
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			input:   "invalid-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uuid UUID
			err := uuid.FromString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidUUID, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUUID_FromB58(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid B58",
			input:   NewUUID().B58(),
			wantErr: false,
		},
		{
			name:    "invalid B58",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uuid UUID
			err := uuid.FromB58(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidUUID, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUUIDv4(t *testing.T) {
	uuid := UUIDv4()
	assert.NotEmpty(t, uuid.String())
	assert.NotEmpty(t, uuid.B58())
}

func TestUUIDv7(t *testing.T) {
	uuid := UUIDv7()
	assert.NotEmpty(t, uuid.String())
	assert.NotEmpty(t, uuid.B58())
}

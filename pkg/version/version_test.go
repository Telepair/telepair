package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()
	assert.Equal(t, info.Version, "dev")
}

func TestGetInfoString(t *testing.T) {
	assert.NotEmpty(t, GetInfoString())
}

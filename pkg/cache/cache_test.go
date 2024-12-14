package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheOptions(t *testing.T) {
	cfg := getConfig{}
	WithTTL(time.Second)(&cfg)
	assert.Equal(t, time.Second, cfg.ttl)
	WithGetter(func(_ context.Context, _ string) (any, error) {
		return nil, nil
	})(&cfg)
	assert.NotNil(t, cfg.getter)
}

package cache

import (
	"context"
	"errors"
	"time"
)

// Type is the type of the cache
type Type string

const (
	TypeMemory Type = "memory"
	TypeDisk   Type = "disk" // use badger
	TypeRedis  Type = "redis"
	TypeNats   Type = "nats"
)

// Cache is a generic cache interface.
type Cache interface {
	Get(ctx context.Context, key string, opts ...GetOption) (any, error)
	Exists(ctx context.Context, key string) (bool, error)
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Keys(ctx context.Context) ([]string, error)
	Clear(ctx context.Context) error
}

var ErrNotFound = errors.New("not found")

// Getter is a function to get the value from the cache
type Getter func(ctx context.Context, key string) (any, error)

// GetOption is the option for Get
type GetOption func(*getConfig)

type getConfig struct {
	ttl    time.Duration
	getter Getter
}

// WithTTL sets the ttl for the cache
func WithTTL(ttl time.Duration) GetOption {
	return func(cfg *getConfig) {
		cfg.ttl = ttl
	}
}

// WithGetter sets the getter for the cache
func WithGetter(getter Getter) GetOption {
	return func(cfg *getConfig) {
		cfg.getter = getter
	}
}

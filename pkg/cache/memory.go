package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

var _ Cache = (*memory)(nil)

// Item is a cache item.
type Item struct {
	Value     any
	ExpiresAt time.Time
}

type memory struct {
	data map[string]Item
	lock sync.Mutex
	log  *slog.Logger
}

// NewMemory creates a new memory cache.
func NewMemory(name string) Cache {
	return &memory{
		data: make(map[string]Item, 256),
		log:  slog.With("component", "cache/"+name),
	}
}

// Get gets the value from the cache.
func (m *memory) Get(ctx context.Context, key string, opts ...GetOption) (any, error) {
	m.lock.Lock()
	item, ok := m.data[key]
	m.lock.Unlock()
	if ok && (item.ExpiresAt.IsZero() || time.Now().Before(item.ExpiresAt)) {
		return item.Value, nil
	}

	cfg := getConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.getter == nil {
		return nil, ErrNotFound
	}

	m.log.Debug("get value from getter", "key", key, "ttl", cfg.ttl)
	value, err := cfg.getter(ctx, key)
	if err != nil {
		m.log.Error("get value from getter", "key", key, "error", err)
		return nil, err
	}

	if cfg.ttl > 0 {
		_ = m.SetWithTTL(ctx, key, value, cfg.ttl)
	} else {
		_ = m.Set(ctx, key, value)
	}

	return value, nil
}

// Exists checks if the key exists in the cache.
func (m *memory) Exists(_ context.Context, key string) (bool, error) {
	m.lock.Lock()
	_, ok := m.data[key]
	m.lock.Unlock()
	return ok, nil
}

// Set sets the value in the cache.
func (m *memory) Set(_ context.Context, key string, value any) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = Item{
		Value: value,
	}
	m.log.Debug("set", "key", key)
	return nil
}

// Delete removes the value from the cache.
func (m *memory) Delete(_ context.Context, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.data, key)
	m.log.Debug("deleted", "key", key)
	return nil
}

// SetWithTTL sets the value in the cache with a given TTL.
func (m *memory) SetWithTTL(_ context.Context, key string, value any, ttl time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = Item{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	m.log.Debug("set with ttl", "key", key, "ttl", ttl)
	return nil
}

// Keys returns all keys in the cache.
func (m *memory) Keys(_ context.Context) ([]string, error) {
	keys := make([]string, 0, len(m.data))
	m.lock.Lock()
	defer m.lock.Unlock()
	for key := range m.data {
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *memory) Clear(_ context.Context) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data = make(map[string]Item, 256)
	return nil
}

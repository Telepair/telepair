package cache

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleCache_memory() {
	cache := NewMemory("test")

	_ = cache.Set(context.Background(), "key1", "value1")
	_ = cache.Set(context.Background(), "key2", "value2")

	ok, _ := cache.Exists(context.Background(), "key1")
	fmt.Println(ok)

	val, _ := cache.Get(context.Background(), "key1")
	fmt.Println(val)

	keys, _ := cache.Keys(context.Background())
	slices.Sort(keys)
	fmt.Println(keys)

	_ = cache.Delete(context.Background(), "key1")

	_ = cache.SetWithTTL(context.Background(), "key3", "value3", time.Second)

	// Output:
	// true
	// value1
	// [key1 key2]
}

func TestMemory(t *testing.T) {
	key := "key"
	key2 := "key2"
	val := "value"
	val2 := "value2"

	cache := NewMemory("test")

	value, err := cache.Get(context.Background(), key)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, value)

	ok, err := cache.Exists(context.Background(), key)
	assert.NoError(t, err)
	assert.False(t, ok)
	ok, err = cache.Exists(context.Background(), "not_exists")
	assert.NoError(t, err)
	assert.False(t, ok)

	err = cache.Set(context.Background(), key, val)
	assert.NoError(t, err)

	value, err = cache.Get(context.Background(), key)
	assert.NoError(t, err)
	assert.Equal(t, val, value)

	err = cache.Delete(context.Background(), key)
	assert.NoError(t, err)

	value, err = cache.Get(context.Background(), key)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, value)

	err = cache.SetWithTTL(context.Background(), key, val, time.Second)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	value, err = cache.Get(context.Background(), key)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, value)

	cache = NewMemory("test2")
	keys, err := cache.Keys(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, keys)

	err = cache.Set(context.Background(), key, val)
	assert.NoError(t, err)
	err = cache.Set(context.Background(), key2, val2)
	assert.NoError(t, err)
	keys, err = cache.Keys(context.Background())
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{key, key2}, keys)

	err = cache.Delete(context.Background(), key)
	assert.NoError(t, err)
	err = cache.Delete(context.Background(), key2)
	assert.NoError(t, err)
	keys, err = cache.Keys(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, keys)

	ct := 0
	cache = NewMemory("test3")
	v, err := cache.Get(context.Background(), key,
		WithTTL(time.Second),
		WithGetter(func(_ context.Context, key string) (any, error) {
			ct++
			return key, nil
		}))
	assert.NoError(t, err)
	assert.Equal(t, key, v)
	assert.Equal(t, 1, ct)
	v, err = cache.Get(context.Background(), key)
	assert.NoError(t, err)
	assert.Equal(t, key, v)
	assert.Equal(t, 1, ct)
	time.Sleep(2 * time.Second)
	v, err = cache.Get(context.Background(), key)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, v)
}

package cache

import (
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrInvalidMaxSize = errors.New("must specify a positive maxSize")
	ErrInvalidTTL     = errors.New("must specify a positive ttl")
	ErrEmptyKey       = errors.New("key is empty")
)

// item is a cache element.
type item[V any] struct {
	value    V
	lastUsed time.Time
}

// Cacher is a non-thread-safe or thread-safe fixed size LRU-like cache with invalidate items by ttl.
type Cacher[K comparable, V any] struct {
	maxSize int
	ttl     time.Duration

	logger         *zap.Logger
	lock           sync.RWMutex
	items          map[K]item[V]
	updateLastUsed bool
	threadSafe     bool

	now func() time.Time
}

// NewCacher creates an cache of the given maxSize.
func NewCacher[K comparable, V any](maxSize int, ttl time.Duration, opts ...CacherOption[K, V]) (*Cacher[K, V], error) {
	if maxSize <= 0 {
		return nil, ErrInvalidMaxSize
	}

	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}

	c := &Cacher[K, V]{
		maxSize: maxSize,
		ttl:     ttl,
		items:   make(map[K]item[V], maxSize),
		now:     time.Now,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// CacherOption is functional option.
type CacherOption[K comparable, V any] func(o *Cacher[K, V])

// WithDebug enables debug logging.
func WithDebug[K comparable, V any](logger *zap.Logger) CacherOption[K, V] {
	return func(c *Cacher[K, V]) {
		c.logger = logger.Named("cacher")
	}
}

// WithUpdateLastUsed enables update last used.
func WithUpdateLastUsed[K comparable, V any]() CacherOption[K, V] {
	return func(c *Cacher[K, V]) {
		c.updateLastUsed = true
	}
}

// WithThreadSafe enables thread safe.
func WithThreadSafe[K comparable, V any]() CacherOption[K, V] {
	return func(c *Cacher[K, V]) {
		c.threadSafe = true
	}
}

// Get looks up a key's value from the cache.
func (c *Cacher[K, V]) Get(key K) (value V, ok bool) {
	if c.threadSafe {
		c.lock.Lock()
		defer c.lock.Unlock()
	}

	elem, ok := c.items[key]
	if !ok {
		return value, false
	}

	if elem.lastUsed.Add(c.ttl).Before(c.now()) {
		delete(c.items, key)

		c.log("Get: deleted by ttl",
			zap.Any("key", key),
			zap.Any("elem.value", elem.value),
			zap.Time("elem.lastUsed", elem.lastUsed),
		)

		return value, false
	}

	if c.updateLastUsed {
		elem.lastUsed = c.now()
		c.items[key] = elem
	}

	c.log("Get",
		zap.Any("key", key),
		zap.Any("elem.value", elem.value),
		zap.Time("elem.lastUsed", elem.lastUsed),
	)

	return elem.value, true
}

// Set adds a value to the cache.
func (c *Cacher[K, V]) Set(key K, value V) error {
	if isEmpty(key) {
		return ErrEmptyKey
	}

	if c.threadSafe {
		c.lock.Lock()
	}

	if len(c.items) == c.maxSize {
		c.clearSpace()
	}

	elem := item[V]{
		value:    value,
		lastUsed: c.now(),
	}

	c.items[key] = elem

	if c.threadSafe {
		c.lock.Unlock()
	}

	c.log("Set",
		zap.Any("key", key),
		zap.Any("elem.value", elem.value),
		zap.Time("elem.lastUsed", elem.lastUsed),
	)

	return nil
}

// clearSpace removes old items from the cache.
func (c *Cacher[K, V]) clearSpace() {
	var keyForDelete K
	var elemForDelete item[V]

	n := c.now()
	t := c.now()
	for k, v := range c.items {
		if v.lastUsed.Add(c.ttl).Before(n) {
			delete(c.items, k)

			c.log("clearSpace: delete invalidate elem",
				zap.Any("key", k),
				zap.Any("elem.value", v.value),
				zap.Time("elem.lastUsed", v.lastUsed),
				zap.Time("now", n),
			)

			continue
		}

		if v.lastUsed.Before(t) {
			t = v.lastUsed
			keyForDelete = k
			elemForDelete = v
		}
	}

	if len(c.items) < c.maxSize {
		return
	}

	delete(c.items, keyForDelete)

	c.log("clearSpace: delete last used",
		zap.Any("key", keyForDelete),
		zap.Any("elem.value", elemForDelete.value),
		zap.Time("elem.lastUsed", elemForDelete.lastUsed),
	)
}

// log debug message.
func (c *Cacher[K, V]) log(msg string, fields ...zapcore.Field) {
	if c.logger != nil {
		c.logger.Debug(msg, fields...)
	}
}

// isEmpty checking for the default value.
func isEmpty[T comparable](v T) bool {
	var empty T

	return v == empty
}

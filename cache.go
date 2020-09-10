package freeGeoIP

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

// ICache is the custom cache interface used by the library.
// If Info is not found in the cache in Get, then the system
// default error should be returned i.e. `ErrCacheMissed`
// See it's default implementation: _Cache
type ICache interface {
	Set(ctx context.Context, info *Info)
	Get(ctx context.Context, ip IP) (*Info, error)
}

// NoopCache empty cache implementation
type NoopCache struct{}

func (n NoopCache) Set(context.Context, *Info) {}

func (n NoopCache) Get(context.Context, IP) (*Info, error) {
	return nil, ErrCacheMissed
}

const (
	DefaultCacheExpiry = 24 * time.Hour
	NoCacheExpiration  = cache.NoExpiration
)

type CacheExpiryFunction func(ctx context.Context, ip IP) time.Duration

// _Cache implements a default ICache implementation
type _Cache struct {
	cache *cache.Cache
	expFn CacheExpiryFunction
}

func DefaultCache() ICache {
	return NewCache(DefaultCacheExpiry, nil)
}

func NonExpiryCache() ICache {
	return NewCache(NoCacheExpiration, nil)
}

func NewCache(expiry time.Duration, expFn CacheExpiryFunction) ICache {
	if expiry <= 0 {
		expiry = NoCacheExpiration
	}
	if expFn == nil {
		expFn = func(context.Context, IP) time.Duration {
			return expiry
		}
	}
	return &_Cache{
		cache: cache.New(expiry, expiry),
		expFn: expFn,
	}
}

func (c *_Cache) Set(ctx context.Context, info *Info) {
	if info == nil {
		return
	}
	c.cache.Set(info.IP.String(), info, c.expFn(ctx, info.IP))
}

func (c *_Cache) Get(_ context.Context, ip IP) (*Info, error) {
	if got, ok := c.cache.Get(ip.String()); ok {
		info, ok := got.(*Info)
		if ok && info != nil {
			return info, nil
		}
	}
	return nil, ErrCacheMissed
}

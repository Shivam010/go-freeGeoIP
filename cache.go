// Copyright 2020 Shivam Rathore
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Set does nothing
func (n NoopCache) Set(context.Context, *Info) {}

// Get does nothing
func (n NoopCache) Get(context.Context, IP) (*Info, error) {
	return nil, ErrCacheMissed
}

// CacheExpiryFunction is the functional parameter for the Cache implementation
// to change expiry for certain IP set or context conditions.
type CacheExpiryFunction func(ctx context.Context, ip IP) time.Duration

const (
	// DefaultCacheExpiry is the default value for _Cache expiry
	DefaultCacheExpiry = 24 * time.Hour
	// DefaultCacheExpiry is constant for specifying no expiry
	NoCacheExpiration = cache.NoExpiration
	// SkipCache is constant for explicitly denying cache hit. Mainly
	// to be used in CacheExpiryFunction for filtering some ip
	SkipCache = -1 << 63
)

// _Cache implements a default ICache implementation
type _Cache struct {
	cache *cache.Cache
	expFn CacheExpiryFunction
}

// DefaultCache is the default cache implementation with 24 Hours expiry
func DefaultCache() ICache {
	return NewCache(DefaultCacheExpiry, nil)
}

// NonExpiryCache is the default cache implementation without any expiry
func NonExpiryCache() ICache {
	return NewCache(NoCacheExpiration, nil)
}

// NewCache is the constructor that returns the ICache implementation with
// the custom defined expiry duration or the `CacheExpiryFunction` to filter
// any selected ip set for expiry.
//
// If expiry is set to SkipCache and expFn is also nil, the the NoopCache{}
// is returned. And if expiry is negative, then NoCacheExpiration will be used.
// Otherwise, provided functionality will continue.
// Note: If expFn is defined then the expiry duration will not be used, define
// carefully.
func NewCache(expiry time.Duration, expFn CacheExpiryFunction) ICache {
	if expiry == SkipCache && expFn == nil {
		return NoopCache{}
	}
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

// Set will use the provided expiry duration and save info in cache
func (c *_Cache) Set(ctx context.Context, info *Info) {
	if info == nil {
		return
	}
	dur := c.expFn(ctx, info.IP)
	if dur == SkipCache {
		return
	}
	c.cache.Set(info.IP.String(), info, dur)
}

// Get will retrieve the saved/cached ip info and if not found then a cache
// missed error, `ErrCacheMissed` will be returned
// The error will also be returned when explicit cache miss is requested
func (c *_Cache) Get(ctx context.Context, ip IP) (*Info, error) {
	// check for explicit cache miss
	if dur := c.expFn(ctx, ip); dur != SkipCache {
		if got, ok := c.cache.Get(ip.String()); ok {
			info, ok := got.(*Info)
			if ok && info != nil {
				return info, nil
			}
		}
	}
	return nil, ErrCacheMissed
}

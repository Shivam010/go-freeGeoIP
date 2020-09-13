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

package freeGeoIP_test

import (
	"context"
	"testing"
	"time"

	"github.com/Shivam010/go-freeGeoIP"
)

func TestNonExpiryCacheCache(t *testing.T) {
	cache := freeGeoIP.NonExpiryCache()
	ctx, st := context.Background(), time.Now()
	times, info := 0, response()
	cache.Set(ctx, info)
	for {
		got, err := cache.Get(ctx, info.IP)
		if err != nil {
			t.Errorf("cache.Get() error = %v, want no error", err)
			return
		}
		if compare(t, got, info) {
			return
		}
		times++
		time.Sleep(time.Millisecond)
		if time.Since(st) >= 100*time.Millisecond {
			break
		}
	}
	t.Logf("Cache read %v times", times)
}

func TestNormalCache(t *testing.T) {
	cache := freeGeoIP.NewCache(10*time.Millisecond, nil)
	ctx, st := context.Background(), time.Now()
	times, info := 0, response()
	cache.Set(ctx, info)
	for {
		diff := time.Since(st)
		got, err := cache.Get(ctx, info.IP)
		if diff <= 10*time.Millisecond && err != nil {
			t.Errorf("cache.Get() error = %v, want no error", err)
			return
		}
		if diff > 10*time.Millisecond && err != freeGeoIP.ErrCacheMissed {
			t.Errorf("cache.Get() error = %v, want %v", err, freeGeoIP.ErrCacheMissed)
			return
		}
		if err == nil {
			if compare(t, got, info) {
				return
			}
		}
		times++
		time.Sleep(time.Millisecond)
		if diff >= 100*time.Millisecond {
			break
		}
	}
	t.Logf("Cache read %v times", times)
}

func TestNilCacheSet(t *testing.T) {
	cache := freeGeoIP.DefaultCache()
	cache.Set(context.TODO(), nil)
}

func TestNoopCache(t *testing.T) {
	cache := freeGeoIP.NoopCache{}
	cache.Set(context.TODO(), nil)
	if _, err := cache.Get(context.TODO(), nil); err != freeGeoIP.ErrCacheMissed {
		t.Errorf("cache.Get() error = %v, want %v", err, freeGeoIP.ErrCacheMissed)
		return
	}
}

func TestSkipCache(t *testing.T) {
	cache := freeGeoIP.NewCache(freeGeoIP.SkipCache, nil)
	cache.Set(context.TODO(), nil)
	if _, err := cache.Get(context.TODO(), nil); err != freeGeoIP.ErrCacheMissed {
		t.Errorf("cache.Get() error = %v, want %v", err, freeGeoIP.ErrCacheMissed)
		return
	}
}

func TestCacheWithOverrideFunction(t *testing.T) {
	// cache with custom function
	cache := freeGeoIP.NewCache(freeGeoIP.NoCacheExpiration,
		func(ctx context.Context, ip freeGeoIP.IP) time.Duration {
			// one can use context values for conditions
			if ip.String() == broadcastIP {
				return freeGeoIP.SkipCache
			}
			return freeGeoIP.NoCacheExpiration
		},
	)
	ctx, st := context.Background(), time.Now()

	// normal response() info cache will not expire
	times, info := 0, response()
	cache.Set(ctx, info)
	for {
		got, err := cache.Get(ctx, info.IP)
		if err != nil {
			t.Errorf("cache.Get() error = %v, want no error", err)
			return
		}
		if compare(t, got, info) {
			return
		}
		times++
		time.Sleep(time.Millisecond)
		if time.Since(st) >= 100*time.Millisecond {
			break
		}
	}

	// broadcast ip info cache will never be set and will always be skipped
	times, info = 0, broadcastResponse()
	cache.Set(ctx, info)
	for {
		got, err := cache.Get(ctx, info.IP)
		if err != freeGeoIP.ErrCacheMissed {
			t.Errorf("cache.Get() error = %v, want %v", err, freeGeoIP.ErrCacheMissed)
			return
		}
		if got != nil {
			t.Errorf("Info got = %v, want %v", got, nil)
			return
		}
		times++
		time.Sleep(time.Millisecond)
		if time.Since(st) >= 100*time.Millisecond {
			break
		}
	}
}

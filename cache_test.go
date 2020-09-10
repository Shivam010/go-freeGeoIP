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

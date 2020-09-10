package freeGeoIP_test

import (
	"context"
	"testing"

	"github.com/Shivam010/go-freeGeoIP"
)

func TestGetGeoInfo_DefaultClient(t *testing.T) {
	// default client with cache
	cli := freeGeoIP.DefaultClient()
	ctx := context.Background()

	// first call for normal response IP
	res := cli.GetGeoInfoFromString(ctx, responseIP)
	if err := res.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if res.Cached {
		t.Fatalf("GetGeoInfoFromString() for new call output must not be cached")
	}
	if compare(t, res.Info, response()) {
		return
	}
	curLimit := res.Meta.Remaining

	// second call, for same ip; response must be cached
	sec := cli.GetGeoInfoFromString(ctx, responseIP)
	if err := sec.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if !sec.Cached {
		t.Fatalf("GetGeoInfoFromString() for second call output must be cached")
	}
	if compare(t, sec.Info, response()) {
		return
	}
	if remains := sec.Meta.Remaining; remains != curLimit {
		t.Fatalf("GetGeoInfoFromString() remaining limit got = %v, and want %v", remains, curLimit)
	}

	// third call for new ip; remaining limit must be one less
	thr := cli.GetGeoInfoFromString(ctx, broadcastIP)
	if err := thr.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if thr.Cached {
		t.Fatalf("GetGeoInfoFromString() for new call output must not be cached")
	}
	if compare(t, thr.Info, broadcastResponse()) {
		return
	}
	if remains := thr.Meta.Remaining; remains != curLimit-1 {
		t.Fatalf("GetGeoInfoFromString() remaining limit got = %v, and want %v", remains, curLimit-1)
	}
}

func TestGetGeoInfo_EmptyClient(t *testing.T) {
	// default client with cache
	cli := &freeGeoIP.Client{}
	ctx := context.Background()

	// first call for normal response IP
	res := cli.GetGeoInfoFromString(ctx, responseIP)
	if err := res.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if res.Cached {
		t.Fatalf("GetGeoInfoFromString() for new call output must not be cached")
	}
	if compare(t, res.Info, response()) {
		return
	}
	curLimit := res.Meta.Remaining

	// second call, for same ip; response must be not be cached as no cache is used
	sec := cli.GetGeoInfoFromString(ctx, responseIP)
	if err := sec.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if sec.Cached {
		t.Fatalf("GetGeoInfoFromString() for second call output must not be cached")
	}
	if compare(t, sec.Info, response()) {
		return
	}
	if remains := sec.Meta.Remaining; remains != curLimit-1 {
		t.Fatalf("GetGeoInfoFromString() remaining limit got = %v, and want %v", remains, curLimit)
	}
	curLimit = sec.Meta.Remaining

	// third call for new ip; remaining limit must be one less
	thr := cli.GetGeoInfoFromString(ctx, broadcastIP)
	if err := thr.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if thr.Cached {
		t.Fatalf("GetGeoInfoFromString() for new call output must not be cached")
	}
	if compare(t, thr.Info, broadcastResponse()) {
		return
	}
	if remains := thr.Meta.Remaining; remains != curLimit-1 {
		t.Fatalf("GetGeoInfoFromString() remaining limit got = %v, and want %v", remains, curLimit-1)
	}
}

func TestInvalidInput(t *testing.T) {
	// default client with cache
	cli := freeGeoIP.DefaultClient()
	ctx := context.Background()

	res := cli.GetGeoInfoFromString(ctx, "responseIP")
	if err := res.Error; err != freeGeoIP.ErrNoResponse {
		t.Fatalf("GetGeoInfoFromString() error = %v, want %v", err, freeGeoIP.ErrNoResponse)
	}

	res = cli.GetGeoInfo(ctx, freeGeoIP.IP{0})
	if err := res.Error; err != freeGeoIP.ErrNoResponse {
		t.Fatalf("GetGeoInfo() error = %v, want %v", err, freeGeoIP.ErrNoResponse)
	}
}

func TestNoInput(t *testing.T) {
	// default client with cache
	cli := freeGeoIP.DefaultClient()
	ctx := context.Background()

	res := cli.GetGeoInfoFromString(ctx, "")
	if err := res.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if res.Cached {
		t.Fatalf("GetGeoInfoFromString() for new call output must not be cached")
	}
	if res.Info == nil {
		t.Fatalf("GetGeoInfoFromString() Info must not be nil")
	}
	curLimit := res.Meta.Remaining

	sec := cli.GetGeoInfo(ctx, freeGeoIP.IP{})
	if err := sec.Error; err != nil {
		t.Fatalf("GetGeoInfo() error = %v, want no error", err)
	}
	if !sec.Cached {
		t.Fatalf("GetGeoInfo() for second call output must be cached")
	}
	if compare(t, sec.Info, res.Info) {
		t.Fatalf("GetGeoInfo() Info must not be nil")
	}
	if remains := sec.Meta.Remaining; remains != curLimit {
		t.Fatalf("GetGeoInfo() remaining limit got = %v, and want %v", remains, curLimit)
	}
}

func TestManipulatingGotResponse(t *testing.T) {
	// default client with cache
	cli := freeGeoIP.DefaultClient()
	ctx := context.Background()

	res := cli.GetGeoInfoFromString(ctx, broadcastIP)
	if err := res.Error; err != nil {
		t.Fatalf("GetGeoInfoFromString() error = %v, want no error", err)
	}
	if res.Cached {
		t.Fatalf("GetGeoInfoFromString() for new call output must not be cached")
	}
	if compare(t, res.Info, broadcastResponse()) {
		return
	}

	// manipulate response info -> cache should not be manipulated
	res.Info.CountryCode = "IN"

	sec := cli.GetGeoInfo(ctx, freeGeoIP.IP{0, 0, 0, 0}) // broadcast ip
	if err := sec.Error; err != nil {
		t.Fatalf("GetGeoInfo() error = %v, want no error", err)
	}
	if !sec.Cached {
		t.Fatalf("GetGeoInfo() for second call output must be cached")
	}
	if compare(t, sec.Info, broadcastResponse()) {
		t.Fatalf("GetGeoInfo() Info must not be nil")
	}
	if sec.Info.CountryCode == res.Info.CountryCode {
		t.Fatalf("GetGeoInfo() cache must not be manipulative")
	}
}

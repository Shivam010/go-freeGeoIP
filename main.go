package freeGeoIP

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	Endpoint = "https://freegeoip.app/json/"

	_HeaderResetIn   = "x-ratelimit-reset"
	_HeaderLimit     = "x-ratelimit-limit"
	_HeaderRemaining = "x-ratelimit-remaining"
)

var (
	// globalMeta maintains the API meta information for the cached information
	// and mu protects globalMeta against updates
	mu         *sync.Mutex
	globalMeta = &MetaInfo{}
	// noopLogger is the quiet logger to handle nil logger in Client
	noopLogger = log.New(ioutil.Discard, "", 0)
	// defaultLogger is the default logger implementation for DefaultClient
	defaultLogger = log.New(os.Stderr, "freeGeoIP ", log.LstdFlags)
)

func init() {
	mu = &sync.Mutex{}
	mu.Lock()
	globalMeta.Limit = 15000
	globalMeta.Remaining = globalMeta.Limit
	globalMeta.ResetIn = time.Hour
	mu.Unlock()
}

// Client is our Free GeoLocation information client.
// If Cache is not provided then will always make the http request to https://freegeoip.app/json/
// If HttpCli is not provided then will always use http.DefaultClient
// And if Logger is not provided, then a noopLogger will be used
type Client struct {
	Cache   ICache
	HttpCli *http.Client
	Logger  *log.Logger
}

// DefaultClient is the library default geo location client with an in-memory
// cache with a default expiry of 24 Hours and will set a total of 2 seconds
// timeout in http.Client, with the default inbuilt library logger.
func DefaultClient() *Client {
	return &Client{
		Cache:   DefaultCache(),
		HttpCli: &http.Client{Timeout: time.Second * 2},
		Logger:  defaultLogger,
	}
}

// GetGeoInfoFromString will return the API response for provided `ip` string
// It call the GetGeoInfo.
func (c *Client) GetGeoInfoFromString(ctx context.Context, ip string) Response {
	_ip := ParseIP(ip)
	if len(_ip) == 0 && ip != "" {
		return fillResponse(nil, ErrNoResponse, nil, struct{}{})
	}
	return c.GetGeoInfo(ctx, _ip)
}

// GetGeoInfo will return the free geolocation api response for the provided IP
// and uses the Cached response, if cache is used. For default empty Client
// behaviour see Client object description
func (c *Client) GetGeoInfo(ctx context.Context, ip IP) Response {
	if c.Logger == nil {
		c.Logger = noopLogger
	}
	if c.HttpCli == nil {
		c.HttpCli = http.DefaultClient
	}
	if c.Cache == nil {
		c.Cache = NoopCache{}
	}
	// check cache
	info, err := c.Cache.Get(ctx, ip)
	if err == nil {
		c.Logger.Println("cache is hit for '" + ip.String())
		return fillResponse(info, nil, nil, struct{}{})
	}
	c.Logger.Println("cache for '"+ip.String()+"' is missed with error:", err)
	// call api
	return c.do(ctx, ip)
}

// do is the internal method used to make the http request to API
func (c *Client) do(ctx context.Context, ip IP) Response {
	// http request
	u, _ := url.Parse(Endpoint)
	u.Path += ip.String()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		c.Logger.Println("http.NewRequest error:", err)
		return fillResponse(nil, wrapError("http", err), nil)
	}

	// request's response
	resp, err := c.HttpCli.Do(req.WithContext(ctx))
	if err != nil {
		c.Logger.Println("http response error:", err)
		return fillResponse(nil, wrapError("http", err), nil)
	}
	defer resp.Body.Close()

	// meta information
	meta := extractMetaInfo(resp.Header)

	// rate limit check
	if resp.StatusCode == http.StatusForbidden {
		c.Logger.Println(ErrLimitReached)
		return fillResponse(nil, ErrLimitReached, meta)
	}

	// invalid ip
	if resp.StatusCode == http.StatusNotFound {
		c.Logger.Println(ErrNoResponse)
		return fillResponse(nil, ErrNoResponse, meta)
	}

	// non ok status code
	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		c.Logger.Println(ErrInternal, "status:", resp.Status, "response:", string(data))
		return fillResponse(nil, ErrInternal, meta)
	}

	// finally response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger.Println("unreadable response body:", err)
		return fillResponse(nil, wrapError("response", err), meta)
	}

	// decode
	info, err := Decoder(data)
	if info != nil {
		c.Cache.Set(ctx, info)
		if len(ip) == 0 { // hack: to cache empty ip for next call
			tmp := info.IP
			info.IP = nil
			c.Cache.Set(ctx, info)
			info.IP = tmp
		}
	}
	return fillResponse(info, err, meta)
}

// extractMetaInfo extract the meta details regarding the limit and reset timer
// from the API response headers
func extractMetaInfo(header http.Header) *MetaInfo {
	atoi := func(key string) int64 {
		v, _ := strconv.Atoi(header.Get(key))
		return int64(v)
	}
	meta := &MetaInfo{
		ResetIn:   time.Second * time.Duration(atoi(_HeaderResetIn)),
		Limit:     atoi(_HeaderLimit),
		Remaining: atoi(_HeaderRemaining),
	}
	mu.Lock()
	globalMeta.Limit = meta.Limit
	globalMeta.Remaining = meta.Remaining
	globalMeta.ResetIn = meta.ResetIn
	mu.Unlock()
	return meta
}

// fillResponse returns a combined response for any client method call
func fillResponse(info *Info, err error, meta *MetaInfo, cached ...struct{}) Response {
	if info != nil { // create a copy of info
		tmp := *info
		info = &tmp
	}
	res := Response{
		Info:   info,
		Error:  err,
		Cached: len(cached) > 0,
		Meta:   meta,
	}
	if res.Meta == nil {
		res.Meta = &MetaInfo{
			ResetIn:   globalMeta.ResetIn,
			Limit:     globalMeta.Limit,
			Remaining: globalMeta.Remaining,
		}
	}
	return res
}

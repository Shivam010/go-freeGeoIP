go-freeGeoIP client with inbuilt cache
======================================

[![Build](https://github.com/Shivam010/go-freeGeoIP/workflows/Build/badge.svg)](https://github.com/Shivam010/go-freeGeoIP/actions?query=workflow%3ABuild)
[![Tests & Check](https://github.com/Shivam010/go-freeGeoIP/workflows/Tests%20&%20Check/badge.svg)](https://github.com/Shivam010/go-freeGeoIP/actions?query=workflow%3A%22Tests+%26+Check%22)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shivam010/go-freeGeoIP?dropcache)](https://goreportcard.com/report/github.com/Shivam010/go-freeGeoIP)
[![GoDoc](https://godoc.org/github.com/Shivam010/go-freeGeoIP?status.svg)](https://godoc.org/github.com/Shivam010/go-freeGeoIP)
[![License](https://img.shields.io/badge/license-apache2-mildgreen.svg)](https://github.com/Shivam010/go-freeGeoIP/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/Shivam010/go-freeGeoIP.svg)](https://github.com/Shivam010/go-freeGeoIP/releases)
[![Coverage Status](https://coveralls.io/repos/github/Shivam010/go-freeGeoIP/badge.svg)](https://coveralls.io/github/Shivam010/go-freeGeoIP)

_go-freeGeoIP_ is a Golang client for Free IP Geolocation information API with inbuilt cache support to 
increase the **_15k per hour rate limit_** of the application [https://freegeoip.app/](https://freegeoip.app/)

By default, the client will cache the IP Geolocation information for 24 hours, but the expiry can be set manually.
If you want set the information cache with no expiration time set the expiry function to nil.

> A 24-hour cache expiry will be sufficient overcome the 15k per hour limit.

Installation
------------
`go get github.com/Shivam010/go-freeGeoIP`

_FreeGeoIP.app_ description
---------------------------
_freegeoip.app provides a free IP geolocation API for software developers. It uses a database of IP addresses that 
are associated to cities along with other relevant information like time zone, latitude and longitude._

_You're allowed up to 15,000 queries per hour by default. Once this limit is reached, all of your requests will 
result in HTTP 403, forbidden, until your quota is cleared._

_The HTTP API takes GET requests in the following schema:_

_https://freegeoip.app/{format}/{IP_or_hostname}_

_Supported formats are: csv, xml, json and jsonp. If no IP or hostname is provided, then your own IP is looked up._

Usage
-----
```go
package main

import (
	"context"
	"github.com/Shivam010/go-freeGeoIP"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()

	// Using default client which comes with an in-memory cache implementation
	// with 24 Hour expiry and a http.Client timeout of 2 seconds and a default
	// `log.Logger`
	cli := freeGeoIP.DefaultClient()
	res := cli.GetGeoInfoFromString(ctx, "8.8.8.8")
	if err := res.Error; err != nil {
		log.Println(err)
		return
	}
	// first time retrieval and hence, not a cached output
	cli.Logger.Println(res.Cached) // false

	// Trying again
	res = cli.GetGeoInfoFromString(ctx, "8.8.8.8")
	if err := res.Error; err != nil {
		log.Println(err)
		return
	}
	cli.Logger.Println(res.Cached) // true

	// Using an empty client, which comes with default http client and no cache
	// and no logs
	cli = &freeGeoIP.Client{}
	res = cli.GetGeoInfo(ctx, freeGeoIP.IP{8, 8, 8, 8})
	if err := res.Error; err != nil {
		log.Println(err)
		return
	}

	// You can use the `ICache` interface and provide you any of you cache
	// implementation or can use the library's in-memory (thread safe) with
	// or without expiry.
	cache := freeGeoIP.NewCache(freeGeoIP.NoCacheExpiration,
		func(ctx context.Context, ip freeGeoIP.IP) time.Duration {
			// check ip pattern
			if value := ctx.Value("IP_Skip_Pattern"); value != nil {
				if pat, ok := value.(string); ok {
					if strings.Contains(ip.String(), pat) {
						// always skip caching such ip patterns
						return freeGeoIP.SkipCache
					}
				}
			}
			return freeGeoIP.NoCacheExpiration
		},
	)

	// And you can even provide your own combination of arguments in client
	// by providing a self cache implementation for `freeGeoIP.ICache` or the
	// the http.Client or the log.Logger
	// The below call to NewCache will create a non expiry cache implementation
	cache = freeGeoIP.NewCache(freeGeoIP.NoCacheExpiration, nil)
	cli = &freeGeoIP.Client{
		Cache:   cache,
		HttpCli: &http.Client{Timeout: time.Second},
		Logger:  log.New(ioutil.Discard, "", 0),
	}
	res = cli.GetGeoInfo(ctx, freeGeoIP.IP{8, 8, 8, 8})
	if err := res.Error; err != nil {
		log.Println(err)
		return
	}
}
```

Request for Contribution
------------------------
Contributors are more than welcome and much appreciated. Please feel free to open a PR to improve anything you 
don't like, or would like to add.

Please make your changes in a specific branch and create a pull request into master! If you can, please make sure all 
the changes work properly and does not affect the existing functioning.

No PR is too small! Even the smallest effort is countable.

License
-------
This project is licensed under the [Apache License 2.0](https://github.com/Shivam010/go-freeGeoIP/blob/master/LICENSE)

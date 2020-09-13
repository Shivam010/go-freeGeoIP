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
	"github.com/Shivam010/go-freeGeoIP"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var ctx = context.TODO()

func Example() {

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

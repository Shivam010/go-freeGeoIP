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

// Package freeGeoIP or go-freeGeoIP is a Golang client for Free IP Geolocation information API with inbuilt cache
// support to increase the 15k per hour rate limit of the application https://freegeoip.app/
//
// By default, the client will cache the IP Geolocation information for 24 hours, but the expiry can be set manually.
// If you want set the information cache with no expiration time set the expiry function to nil.
//
//	A 24-hour cache expiry will be sufficient overcome the 15k per hour limit.
//
// Installation
//
// You can use the package using the following command:
//	go get github.com/Shivam010/go-freeGeoIP
//
// FreeGeoIP.app description
//
// freegeoip.app provides a free IP geolocation API for software developers. It uses a database of IP addresses that
// are associated to cities along with other relevant information like time zone, latitude and longitude.
//
// You're allowed up to 15,000 queries per hour by default. Once this limit is reached, all of your requests will
// result in HTTP 403, forbidden, until your quota is cleared.
//
// The HTTP API takes GET requests in the following schema:
// 	https://freegeoip.app/{format}/{IP_or_hostname}
//
// Supported formats are: csv, xml, json and jsonp. If no IP or hostname is provided, then your own IP is looked up.
//
// Request for Contribution
//
// Contributors are more than welcome and much appreciated. Please feel free to open a PR to improve anything you
// don't like, or would like to add.
//
// Please make your changes in a specific branch and request to pull into master! If you can please make sure all
// the changes work properly and does not affect the existing functioning.
//
// No PR is too small! Even the smallest effort is countable.
//
// License
//
// This project is licensed under the MIT license.(https://github.com/Shivam010/go-freeGeoIP/blob/master/LICENSE)
package freeGeoIP

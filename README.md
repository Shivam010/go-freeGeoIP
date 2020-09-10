go-freeGeoIP client with inbuilt cache
======================================
_go-freeGeoIP_ is a Golang client for Free IP Geolocation information API with inbuilt cache support to 
increase the **_15k per hour rate limit_** of the application [https://freegeoip.app/](https://freegeoip.app/)

By default, the client will cache the IP Geolocation information for 24 hours, but the expiry can be set manually.
If you want set the information cache with no expiration time set the expiry function to nil.

> A 24-hour cache expiry will be sufficient overcome the 15k per hour limit.

_FreeGeoIP.app_ description
---------------------------
_freegeoip.app provides a free IP geolocation API for software developers. It uses a database of IP addresses that 
are associated to cities along with other relevant information like time zone, latitude and longitude._

_You're allowed up to 15,000 queries per hour by default. Once this limit is reached, all of your requests will 
result in HTTP 403, forbidden, until your quota is cleared._

_The HTTP API takes GET requests in the following schema:_

_https://freegeoip.app/{format}/{IP_or_hostname}_

_Supported formats are: csv, xml, json and jsonp. If no IP or hostname is provided, then your own IP is looked up._

Request for Contribution
------------------------
Contributors are more than welcome and much appreciated. Please feel free to open a PR to improve anything you 
don't like, or would like to add.

Please make your changes in a specific branch and request to pull into master! If you can please make sure all 
the changes work properly and does not affect the existing functioning.

No PR is too small! Even the smallest effort is countable.

License
-------
This project is licensed under the [MIT license.](https://github.com/Shivam010/go-freeGeoIP/blob/master/LICENSE)

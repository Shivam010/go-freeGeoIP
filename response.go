package freeGeoIP

import (
	"encoding/json"
	"time"
)

type DD []string

// Info is the object specifying geo-location information obtained from
// the application https://freegeoip.app/
type Info struct {
	IP IP `json:"ip"`

	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	MetroCode   float64 `json:"metro_code"`

	TimeZone *Location `json:"time_zone"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func Decoder(data []byte) (*Info, error) {
	if len(data) == 0 {
		return nil, ErrNoResponse
	}
	info := &Info{}
	if err := json.Unmarshal(data, info); err != nil {
		return nil, wrapError("decode", err)
	}
	if len(info.IP) == 0 {
		return nil, ErrNoResponse
	}
	return info, nil
}

type Response struct {
	// Info will contain the geo-location information about the IP provided in
	// request, and will be nil, if Error is not nil
	Info *Info
	// Error will contains any error if occurred during the retrieval, including
	// the limit reached error, `ErrLimitReached`. Otherwise will be nil
	Error error
	// Cached will be true if the Info response is retrieved from cached
	Cached bool
	// The MetaInfo may not have correct value if the geo info is retrieved from
	// cache, i.e. if Cached is true
	Meta *MetaInfo
}

type MetaInfo struct {
	// duration in which limit will reset
	ResetIn time.Duration
	// total rate limit per hour
	Limit int64
	// remaining limit for the resetsIn duration
	Remaining int64
}

package freeGeoIP_test

import (
	"net"
	"testing"
	"time"

	"github.com/Shivam010/go-freeGeoIP"
)

const (
	dnsIP       = "8.8.8.8"
	broadcastIP = "0.0.0.0"
	responseIP  = "2401:4900:16ff:f1ef:fff5:f63e:8a25:a38a"
)

// response uses the following IP: 2401:4900:16ff:f1ef:fff5:f63e:8a25:a38a
// (it is randomly selected and has nothing to target any one)
// Data url https://freegeoip.app/json/2401:4900:16ff:f1ef:fff5:f63e:8a25:a38a
func response() *freeGeoIP.Info {
	return &freeGeoIP.Info{
		IP:          freeGeoIP.IP(net.ParseIP("2401:4900:16ff:f1ef:fff5:f63e:8a25:a38a")),
		CountryCode: "IN",
		CountryName: "India",
		RegionCode:  "KA",
		RegionName:  "Karnataka",
		City:        "Belgaum",
		ZipCode:     "590006",
		MetroCode:   0,
		TimeZone:    freeGeoIP.LocationF(time.FixedZone("Asia/Kolkata", int(5*time.Hour+30*time.Minute))),
		Latitude:    15.8521,
		Longitude:   74.5045,
	}
}

// broadcastResponse returns IP info for broadcast ip 0.0.0.0
// Data url https://freegeoip.app/json/0.0.0.0
func broadcastResponse() *freeGeoIP.Info {
	return &freeGeoIP.Info{
		IP:       freeGeoIP.IP(net.ParseIP("0.0.0.0")),
		TimeZone: freeGeoIP.LocationF(time.UTC),
	}
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *freeGeoIP.Info
		err     error
		wantErr bool
	}{
		{
			name:    "Nil Response",
			data:    []byte(nil),
			want:    nil,
			err:     freeGeoIP.ErrNoResponse,
			wantErr: true,
		},
		{
			name:    "Empty object Response",
			data:    []byte(`{"":""}`),
			want:    nil,
			err:     freeGeoIP.ErrNoResponse,
			wantErr: true,
		},
		{
			name:    "Minified Response",
			data:    []byte(`{"ip":"2401:4900:16ff:f1ef:fff5:f63e:8a25:a38a","country_code":"IN","country_name":"India","region_code":"KA","region_name":"Karnataka","city":"Belgaum","zip_code":"590006","time_zone":"Asia/Kolkata","latitude":15.8521,"longitude":74.5045,"metro_code":0}`),
			want:    response(),
			err:     nil,
			wantErr: false,
		},
		{
			name: "Beautified Response",
			data: []byte(`
{
  "ip": "2401:4900:16ff:f1ef:fff5:f63e:8a25:a38a",
  "country_code": "IN",
  "country_name": "India",
  "region_code": "KA",
  "region_name": "Karnataka",
  "city": "Belgaum",
  "zip_code": "590006",
  "time_zone": "Asia/Kolkata",
  "latitude": 15.8521,
  "longitude": 74.5045,
  "metro_code": 0
}
`),
			want:    response(),
			err:     nil,
			wantErr: false,
		},
		{
			name:    "Broadcast Response",
			data:    []byte(`{"ip":"0.0.0.0","country_code":"","country_name":"","region_code":"","region_name":"","city":"","zip_code":"","time_zone":"","latitude":0,"longitude":0,"metro_code":0}`),
			want:    broadcastResponse(),
			err:     nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := freeGeoIP.Decoder(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.err {
				t.Errorf("Decoder() error does not match got = %v, wantErr %v", err, tt.err)
				return
			}
			_ = compare(t, got, tt.want)
		})
	}
}

// compare will compare the provided Info objects and fails the test
// if got is different from want, and returns true if test fails
func compare(t *testing.T, got, want *freeGeoIP.Info) bool {
	if (got == nil) != (want == nil) {
		t.Errorf("Info got = %v, want %v", got, want)
		return true
	}
	if got != nil {
		if got.IP.String() != want.IP.String() {
			t.Errorf("Info.IP got = %v, want %v", got.IP, want.IP)
			return true
		}
		if got.CountryCode != want.CountryCode {
			t.Errorf("Info.CountryCode got = %v, want %v", got.CountryCode, want.CountryCode)
			return true
		}
		if got.CountryName != want.CountryName {
			t.Errorf("Info.CountryName got = %v, want %v", got.CountryName, want.CountryName)
			return true
		}
		if got.RegionCode != want.RegionCode {
			t.Errorf("Info.RegionCode got = %v, want %v", got.RegionCode, want.RegionCode)
			return true
		}
		if got.RegionName != want.RegionName {
			t.Errorf("Info.RegionName got = %v, want %v", got.RegionName, want.RegionName)
			return true
		}
		if got.City != want.City {
			t.Errorf("Info.City got = %v, want %v", got.City, want.City)
			return true
		}
		if got.ZipCode != want.ZipCode {
			t.Errorf("Info.ZipCode got = %v, want %v", got.ZipCode, want.ZipCode)
			return true
		}
		if got.MetroCode != want.MetroCode {
			t.Errorf("Info.MetroCode got = %v, want %v", got.MetroCode, want.MetroCode)
			return true
		}
		if got.TimeZone.String() != want.TimeZone.String() {
			t.Errorf("Info.TimeZone got = %v, want %v", got.TimeZone, want.TimeZone)
			return true
		}
		if got.Latitude != want.Latitude {
			t.Errorf("Info.Latitude got = %v, want %v", got.Latitude, want.Latitude)
			return true
		}
		if got.Longitude != want.Longitude {
			t.Errorf("Info.Longitude got = %v, want %v", got.Longitude, want.Longitude)
			return true
		}
	}
	return false
}

package freeGeoIP

import (
	"errors"
	"strings"
)

const (
	ErrInternal     = _Error("freeGeoIp: something went wrong")
	ErrLimitReached = _Error("freeGeoIp: api limit reached")
	ErrNoResponse   = _Error("freeGeoIp: no information found")
	ErrCacheMissed  = _Error("cache: info not found")
)

type _Error string

func (e _Error) Error() string {
	return string(e)
}

func (e _Error) Unwrap() error {
	list := strings.SplitN(string(e), ": ", 2)
	wrap := list[0]
	if len(list) == 2 {
		wrap = list[1]
	}
	return errors.New(wrap)
}

func wrapError(pre string, err error) _Error {
	if e, ok := err.(_Error); ok {
		return e
	}
	s := err.Error()
	if pre != "" {
		s = pre + ": " + s
	}
	return _Error(s)
}

package filter

import (
	"net"
	"strings"

	"github.com/go-apibox/types"
)

type IPSetFilter struct {
	delimiter  string
	minCount   int
	maxCount   int
	validators []IPSetValidator
	allowVals  []string
	toString   bool
}

type IPSetValidator func(paramName string, paramValue []net.IP) *Error

// IPSet return a IPSet filter.
func IPSet() *IPSetFilter {
	f := new(IPSetFilter)
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// ItemToString return value as string
func (f *IPSetFilter) ItemToString() *IPSetFilter {
	f.toString = true
	return f
}

// Allow allow value is a string in the specified list
func (f *IPSetFilter) Allow(vals ...string) *IPSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Delimiter set the delimiter of set string.
func (f *IPSetFilter) Delimiter(delimiter string) *IPSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *IPSetFilter) MinCount(count int) *IPSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *IPSetFilter) MaxCount(count int) *IPSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *IPSetFilter) AddValidator(validator IPSetValidator) *IPSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemIsIPv4 valid whether ip address in set is ipv4 address.
func (f *IPSetFilter) ItemIsIPv4() *IPSetFilter {
	f.AddValidator(func(paramName string, paramValue []net.IP) *Error {
		// DefaultMask returns nil if ip is not a valid IPv4 address.
		for _, v := range paramValue {
			if v.DefaultMask() == nil {
				return NewError(ErrorInvalidParam, paramName, "ItemNotIPv4")
			}
		}
		return nil
	})
	return f
}

// ItemIsIPv6 valid whether ip address in set is ipv6 address.
func (f *IPSetFilter) ItemIsIPv6() *IPSetFilter {
	f.AddValidator(func(paramName string, paramValue []net.IP) *Error {
		// DefaultMask returns nil if ip is not a valid IPv6 address.
		for _, v := range paramValue {
			if v.DefaultMask() != nil {
				return NewError(ErrorInvalidParam, paramName, "ItemNotIPv6")
			}
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *IPSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var ipVals []net.IP
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		if val != "" {
			fields := strings.Split(val, f.delimiter)
			for _, field := range fields {
				field = strings.Trim(field, " \t\r\n")
				v := net.ParseIP(field)
				if v == nil {
					goto parse_error
				}
				ipVals = append(ipVals, v)
			}
		} else {
			ipVals = []net.IP{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			v := net.ParseIP(field)
			if v == nil {
				goto parse_error
			}
			ipVals = append(ipVals, v)
		}
	case []net.IP:
		ipVals = val
	default:
		goto parse_error
	}

	if len(ipVals) < f.minCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooFew")
	}
	if len(ipVals) > f.maxCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooMany")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, ipVals); err != nil {
			return nil, err
		}
	}

	if f.toString {
		ipStrs := make([]string, 0, len(ipVals))
		for _, ipVal := range ipVals {
			if ipVal != nil {
				ipStrs = append(ipStrs, ipVal.String())
			}
		}
		return ipStrs, nil
	}
	return ipVals, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotIPSet")
}

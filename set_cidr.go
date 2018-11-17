package filter

import (
	"net"
	"strings"

	"github.com/go-apibox/types"
)

type CIDRSetFilter struct {
	delimiter  string
	minCount   int
	maxCount   int
	validators []CIDRSetValidator
	allowVals  []string
}

type CIDRSetValidator func(paramName string, paramValue []*CIDRAddr) *Error

// CIDRSet return a CIDRSet filter.
func CIDRSet() *CIDRSetFilter {
	f := new(CIDRSetFilter)
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *CIDRSetFilter) Allow(vals ...string) *CIDRSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Delimiter set the delimiter of set string.
func (f *CIDRSetFilter) Delimiter(delimiter string) *CIDRSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *CIDRSetFilter) MinCount(count int) *CIDRSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *CIDRSetFilter) MaxCount(count int) *CIDRSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *CIDRSetFilter) AddValidator(validator CIDRSetValidator) *CIDRSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemIsIPv4 valid whether cidr is ipv4 cidr.
func (f *CIDRSetFilter) ItemIsIPv4() *CIDRSetFilter {
	f.AddValidator(func(paramName string, paramValue []*CIDRAddr) *Error {
		// DefaultMask returns nil if ip is not a valid IPv4 address.
		for _, v := range paramValue {
			if v.IP.DefaultMask() == nil {
				return NewError(ErrorInvalidParam, paramName, "ItemNotIPv4")
			}
		}
		return nil
	})
	return f
}

// ItemIsIPv6 valid whether cidr is ipv6 cidr.
func (f *CIDRSetFilter) ItemIsIPv6() *CIDRSetFilter {
	f.AddValidator(func(paramName string, paramValue []*CIDRAddr) *Error {
		// DefaultMask returns nil if ip is not a valid IPv6 address.
		for _, v := range paramValue {
			if v.IP.DefaultMask() != nil {
				return NewError(ErrorInvalidParam, paramName, "ItemNotIPv6")
			}
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *CIDRSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var cidrVals []*CIDRAddr
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
				ip, net, err := net.ParseCIDR(field)
				if err != nil {
					goto parse_error
				}
				cidrVals = append(cidrVals, &CIDRAddr{ip, net, field})
			}
		} else {
			cidrVals = []*CIDRAddr{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			ip, net, err := net.ParseCIDR(field)
			if err != nil {
				goto parse_error
			}
			cidrVals = append(cidrVals, &CIDRAddr{ip, net, field})
		}
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, cidrVals); err != nil {
			return nil, err
		}
	}

	return cidrVals, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotCIDRSet")
}

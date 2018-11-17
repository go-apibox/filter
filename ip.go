package filter

import (
	"net"
	"strings"
)

type IPFilter struct {
	validators []IPValidator
	allowVals  []string
	toString   bool
}

type IPValidator func(paramName string, paramValue *net.IP) *Error

// IP return a IP filter.
func IP() *IPFilter {
	f := new(IPFilter)
	return f
}

// Allow allow value is a string in the specified list
func (f *IPFilter) Allow(vals ...string) *IPFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// ToString return value as string
func (f *IPFilter) ToString() *IPFilter {
	f.toString = true
	return f
}

// AddValidator add a custom validator to filter
func (f *IPFilter) AddValidator(validator IPValidator) *IPFilter {
	f.validators = append(f.validators, validator)
	return f
}

// IsIPv4 valid whether ip address is ipv4 address.
func (f *IPFilter) IsIPv4() *IPFilter {
	f.AddValidator(func(paramName string, paramValue *net.IP) *Error {
		// DefaultMask returns nil if ip is not a valid IPv4 address.
		if paramValue.DefaultMask() == nil {
			return NewError(ErrorInvalidParam, paramName, "NotIPv4")
		}
		return nil
	})
	return f
}

// IsIPv6 valid whether ip address is ipv6 address.
func (f *IPFilter) IsIPv6() *IPFilter {
	f.AddValidator(func(paramName string, paramValue *net.IP) *Error {
		// DefaultMask returns nil if ip is not a valid IPv6 address.
		if paramValue.DefaultMask() != nil {
			return NewError(ErrorInvalidParam, paramName, "NotIPv6")
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *IPFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var ipVal *net.IP
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v := net.ParseIP(val)
		if v == nil {
			goto parse_error
		}
		ipVal = &v
	case net.IP:
		ipVal = &val
	case *net.IP:
		ipVal = val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, ipVal); err != nil {
			return nil, err
		}
	}

	if f.toString {
		if ipVal != nil {
			return ipVal.String(), nil
		} else {
			return "", nil
		}
	}
	return ipVal, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotIP")
}

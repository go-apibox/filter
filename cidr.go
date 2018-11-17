package filter

import (
	"net"
	"strings"
)

type CIDRFilter struct {
	validators []CIDRValidator
	allowVals  []string
	toString   bool
}

type CIDRValidator func(paramName string, paramValue *CIDRAddr) *Error

type CIDRAddr struct {
	IP        net.IP
	IPNet     *net.IPNet
	RawString string
}

func (cidr *CIDRAddr) String() string {
	return cidr.RawString
}

// CIDR return a CIDR filter.
func CIDR() *CIDRFilter {
	f := new(CIDRFilter)
	return f
}

// AddValidator add a custom validator to filter
func (f *CIDRFilter) AddValidator(validator CIDRValidator) *CIDRFilter {
	f.validators = append(f.validators, validator)
	return f
}

// Allow allow value is a string in the specified list
func (f *CIDRFilter) Allow(vals ...string) *CIDRFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// ToString return value as string
func (f *CIDRFilter) ToString() *CIDRFilter {
	f.toString = true
	return f
}

// IsIPv4 valid whether cidr is ipv4 cidr.
func (f *CIDRFilter) IsIPv4() *CIDRFilter {
	f.AddValidator(func(paramName string, paramValue *CIDRAddr) *Error {
		// DefaultMask returns nil if ip is not a valid IPv4 address.
		if paramValue.IP.DefaultMask() == nil {
			return NewError(ErrorInvalidParam, paramName, "NotIPv4")
		}
		return nil
	})
	return f
}

// IsIPv6 valid whether cidr is ipv6 cidr.
func (f *CIDRFilter) IsIPv6() *CIDRFilter {
	f.AddValidator(func(paramName string, paramValue *CIDRAddr) *Error {
		// DefaultMask returns nil if ip is not a valid IPv6 address.
		if paramValue.IP.DefaultMask() != nil {
			return NewError(ErrorInvalidParam, paramName, "NotIPv6")
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *CIDRFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var cidrVal *CIDRAddr
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		ip, net, err := net.ParseCIDR(val)
		if err != nil {
			goto parse_error
		}
		cidrVal = &CIDRAddr{ip, net, val}
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, cidrVal); err != nil {
			return nil, err
		}
	}

	if f.toString {
		if cidrVal != nil {
			return cidrVal.RawString, nil
		} else {
			return "", nil
		}
	}
	return cidrVal, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotCIDR")
}

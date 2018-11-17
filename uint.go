package filter

import (
	"strconv"
	"strings"
)

type UintFilter struct {
	base       int
	validators []UintValidator
	allowVals  []string
}

type UintValidator func(paramName string, paramValue uint) *Error

// Uint return a uint filter.
func Uint() *UintFilter {
	f := new(UintFilter)
	f.base = 10
	return f
}

// Allow allow value is a string in the specified list
func (f *UintFilter) Allow(vals ...string) *UintFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of uint.
// UintFilter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *UintFilter) Base(base int) *UintFilter {
	f.base = base
	return f
}

// AddValidator add a custom validator to filter
func (f *UintFilter) AddValidator(validator UintValidator) *UintFilter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *UintFilter) Min(val uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *UintFilter) Max(val uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *UintFilter) LargerThan(val uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *UintFilter) SmallerThan(val uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *UintFilter) Equal(val uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}

		return nil
	})
	return f
}

// Between valid param value should in the specified range.
func (f *UintFilter) Between(min, max uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		if paramValue < min {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		if paramValue > max {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// In valid param value should in the specified set.
func (f *UintFilter) In(set []uint) *UintFilter {
	f.AddValidator(func(paramName string, paramValue uint) *Error {
		for _, v := range set {
			if v == paramValue {
				return nil
			}
		}
		return NewError(ErrorInvalidParam, paramName, "NotInSet")
	})
	return f
}

// Run make the filter running.
func (f *UintFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVal uint
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v, err := strconv.ParseUint(val, f.base, 0)
		if err != nil {
			goto parse_error
		}
		intVal = uint(v)
	case uint:
		intVal = val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, intVal); err != nil {
			return nil, err
		}
	}

	return intVal, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotUint")
}

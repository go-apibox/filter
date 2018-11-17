package filter

import (
	"strconv"
	"strings"
)

type Uint64Filter struct {
	base       int
	validators []Uint64Validator
	allowVals  []string
}

type Uint64Validator func(paramName string, paramValue uint64) *Error

// Uint64 return a uint64 filter.
func Uint64() *Uint64Filter {
	f := new(Uint64Filter)
	f.base = 10
	return f
}

// Allow allow value is a string in the specified list
func (f *Uint64Filter) Allow(vals ...string) *Uint64Filter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of uint64.
// Uint64Filter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *Uint64Filter) Base(base int) *Uint64Filter {
	f.base = base
	return f
}

// AddValidator add a custom validator to filter
func (f *Uint64Filter) AddValidator(validator Uint64Validator) *Uint64Filter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *Uint64Filter) Min(val uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *Uint64Filter) Max(val uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *Uint64Filter) LargerThan(val uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *Uint64Filter) SmallerThan(val uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *Uint64Filter) Equal(val uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
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
func (f *Uint64Filter) Between(min, max uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
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
func (f *Uint64Filter) In(set []uint64) *Uint64Filter {
	f.AddValidator(func(paramName string, paramValue uint64) *Error {
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
func (f *Uint64Filter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVal uint64
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v, err := strconv.ParseUint(val, f.base, 64)
		if err != nil {
			goto parse_error
		}
		intVal = v
	case uint64:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotUint64")
}

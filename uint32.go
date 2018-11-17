package filter

import (
	"math"
	"strconv"
	"strings"
)

type Uint32Filter struct {
	base       int
	validators []Uint32Validator
	allowVals  []string
}

type Uint32Validator func(paramName string, paramValue uint32) *Error

// Uint32 return a uint32 filter.
func Uint32() *Uint32Filter {
	f := new(Uint32Filter)
	f.base = 10
	return f
}

// Allow allow value is a string in the specified list
func (f *Uint32Filter) Allow(vals ...string) *Uint32Filter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of uint.
// Uint32Filter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *Uint32Filter) Base(base int) *Uint32Filter {
	f.base = base
	return f
}

// AddValidator add a custom validator to filter
func (f *Uint32Filter) AddValidator(validator Uint32Validator) *Uint32Filter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *Uint32Filter) Min(val uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *Uint32Filter) Max(val uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *Uint32Filter) LargerThan(val uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *Uint32Filter) SmallerThan(val uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *Uint32Filter) Equal(val uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
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
func (f *Uint32Filter) Between(min, max uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
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
func (f *Uint32Filter) In(set []uint32) *Uint32Filter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
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
func (f *Uint32Filter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVal uint32
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
		if v > math.MaxUint32 {
			goto parse_error
		}
		intVal = uint32(v)
	case uint32:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotUint32")
}

package filter

import (
	"strconv"
	"strings"
)

type Int64Filter struct {
	base       int
	validators []Int64Validator
	allowVals  []string
}

type Int64Validator func(paramName string, paramValue int64) *Error

// Int64 return a int64 filter.
func Int64() *Int64Filter {
	f := new(Int64Filter)
	f.base = 10
	return f
}

// Allow allow value is a string in the specified list
func (f *Int64Filter) Allow(vals ...string) *Int64Filter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of int64.
// Int64Filter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *Int64Filter) Base(base int) *Int64Filter {
	f.base = base
	return f
}

// AddValidator add a custom validator to filter
func (f *Int64Filter) AddValidator(validator Int64Validator) *Int64Filter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *Int64Filter) Min(val int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *Int64Filter) Max(val int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *Int64Filter) LargerThan(val int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *Int64Filter) SmallerThan(val int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *Int64Filter) Equal(val int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
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
func (f *Int64Filter) Between(min, max int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
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
func (f *Int64Filter) In(set []int64) *Int64Filter {
	f.AddValidator(func(paramName string, paramValue int64) *Error {
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
func (f *Int64Filter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVal int64
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v, err := strconv.ParseInt(val, f.base, 64)
		if err != nil {
			goto parse_error
		}
		intVal = v
	case int64:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt64")
}

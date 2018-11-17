package filter

import (
	"math"
	"strconv"
	"strings"
)

type Int32Filter struct {
	base       int
	validators []Int32Validator
	allowVals  []string
}

type Int32Validator func(paramName string, paramValue int32) *Error

// Int32 return a int32 filter.
func Int32() *Int32Filter {
	f := new(Int32Filter)
	f.base = 10
	return f
}

// Allow allow value is a string in the specified list
func (f *Int32Filter) Allow(vals ...string) *Int32Filter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of int32.
// Int32Filter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *Int32Filter) Base(base int) *Int32Filter {
	f.base = base
	return f
}

// AddValidator add a custom validator to filter
func (f *Int32Filter) AddValidator(validator Int32Validator) *Int32Filter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *Int32Filter) Min(val int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *Int32Filter) Max(val int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *Int32Filter) LargerThan(val int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *Int32Filter) SmallerThan(val int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *Int32Filter) Equal(val int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
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
func (f *Int32Filter) Between(min, max int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
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
func (f *Int32Filter) In(set []int32) *Int32Filter {
	f.AddValidator(func(paramName string, paramValue int32) *Error {
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
func (f *Int32Filter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVal int32
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v, err := strconv.ParseInt(val, f.base, 0)
		if err != nil {
			goto parse_error
		}
		if v > math.MaxInt32 {
			goto parse_error
		}
		intVal = int32(v)
	case int32:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt32")
}

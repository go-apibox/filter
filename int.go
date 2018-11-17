package filter

import (
	"strconv"
	"strings"
)

type IntFilter struct {
	base       int
	validators []IntValidator
	allowVals  []string
}

type IntValidator func(paramName string, paramValue int) *Error

// Int return a int filter.
func Int() *IntFilter {
	f := new(IntFilter)
	f.base = 10
	return f
}

// Allow allow value is a string in the specified list
func (f *IntFilter) Allow(vals ...string) *IntFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of int.
// IntFilter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *IntFilter) Base(base int) *IntFilter {
	f.base = base
	return f
}

// AddValidator add a custom validator to filter
func (f *IntFilter) AddValidator(validator IntValidator) *IntFilter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *IntFilter) Min(val int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *IntFilter) Max(val int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *IntFilter) LargerThan(val int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *IntFilter) SmallerThan(val int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *IntFilter) Equal(val int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
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
func (f *IntFilter) Between(min, max int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
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
func (f *IntFilter) In(set []int) *IntFilter {
	f.AddValidator(func(paramName string, paramValue int) *Error {
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
func (f *IntFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVal int
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
		intVal = int(v)
	case int:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt")
}

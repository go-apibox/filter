package filter

import (
	"math"
	"strconv"
	"strings"
)

type Float64Filter struct {
	validators []Float64Validator
	allowVals  []string
}

type Float64Validator func(paramName string, paramValue float64) *Error

// Float64 return a float64 filter.
func Float64() *Float64Filter {
	f := new(Float64Filter)
	return f
}

// Allow allow value is a string in the specified list
func (f *Float64Filter) Allow(vals ...string) *Float64Filter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// AddValidator add a custom validator to filter
func (f *Float64Filter) AddValidator(validator Float64Validator) *Float64Filter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *Float64Filter) Min(val float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *Float64Filter) Max(val float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *Float64Filter) LargerThan(val float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *Float64Filter) SmallerThan(val float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *Float64Filter) Equal(val float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
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
func (f *Float64Filter) Between(min, max float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
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
func (f *Float64Filter) In(set []float64) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
		for _, v := range set {
			if v == paramValue {
				return nil
			}
		}
		return NewError(ErrorInvalidParam, paramName, "NotInSet")
	})
	return f
}

// DecimalPlace valid whether decimal place is equal to the specified length.
func (f *Float64Filter) DecimalPlace(length int) *Float64Filter {
	f.AddValidator(func(paramName string, paramValue float64) *Error {
		valuef := paramValue * float64(math.Pow(10.0, float64(length)))
		extra := valuef - float64(int(valuef))

		if extra != 0 {
			return NewError(ErrorInvalidParam, paramName, "DecimalPlaceNotMatch")
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *Float64Filter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var floatVal float64
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			goto parse_error
		}
		if v > math.MaxFloat64 {
			goto parse_error
		}
		floatVal = float64(v)
	case float64:
		floatVal = val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, floatVal); err != nil {
			return nil, err
		}
	}

	return floatVal, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotFloat64")
}

package filter

import (
	"math"
	"strconv"
	"strings"
)

type Float32Filter struct {
	validators []Float32Validator
	allowVals  []string
}

type Float32Validator func(paramName string, paramValue float32) *Error

// Float32 return a float32 filter.
func Float32() *Float32Filter {
	f := new(Float32Filter)
	return f
}

// Allow allow value is a string in the specified list
func (f *Float32Filter) Allow(vals ...string) *Float32Filter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// AddValidator add a custom validator to filter
func (f *Float32Filter) AddValidator(validator Float32Validator) *Float32Filter {
	f.validators = append(f.validators, validator)
	return f
}

// Min valid param value should not be smaller than the specified value.
func (f *Float32Filter) Min(val float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// Max valid param value should not be larger than the specified value.
func (f *Float32Filter) Max(val float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// LargerThan valid param value should be larger than the specified value.
func (f *Float32Filter) LargerThan(val float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooSmall")
		}
		return nil
	})
	return f
}

// SmallerThan valid param value should be smaller than the specified value.
func (f *Float32Filter) SmallerThan(val float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLarge")
		}
		return nil
	})
	return f
}

// Equal valid param value should be equal to the specified value.
func (f *Float32Filter) Equal(val float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
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
func (f *Float32Filter) Between(min, max float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
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
func (f *Float32Filter) In(set []float32) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
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
func (f *Float32Filter) DecimalPlace(length int) *Float32Filter {
	f.AddValidator(func(paramName string, paramValue float32) *Error {
		valuef := paramValue * float32(math.Pow(10.0, float64(length)))
		extra := valuef - float32(int(valuef))

		if extra != 0 {
			return NewError(ErrorInvalidParam, paramName, "DecimalPlaceNotMatch")
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *Float32Filter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var floatVal float32
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			goto parse_error
		}
		if v > math.MaxFloat32 {
			goto parse_error
		}
		floatVal = float32(v)
	case float32:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotFloat32")
}

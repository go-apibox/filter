package filter

import (
	"strconv"
	"strings"
)

type TimestampFilter struct {
	validators []TimestampValidator
	allowVals  []string
}

type TimestampValidator func(paramName string, paramValue uint32) *Error

// Timestamp return a timestamp filter.
func Timestamp() *TimestampFilter {
	f := new(TimestampFilter)
	return f
}

// Allow allow value is a string in the specified list
func (f *TimestampFilter) Allow(vals ...string) *TimestampFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// AddValidator add a custom validator to filter
func (f *TimestampFilter) AddValidator(validator TimestampValidator) *TimestampFilter {
	f.validators = append(f.validators, validator)
	return f
}

// StartFrom valid whether start from specified time.
func (f *TimestampFilter) StartFrom(val uint32) *TimestampFilter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}
		return nil
	})
	return f
}

// StartFrom valid whether end to specified time.
func (f *TimestampFilter) EndTo(val uint32) *TimestampFilter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}
		return nil
	})
	return f
}

// After valid whether after specified time.
func (f *TimestampFilter) After(val uint32) *TimestampFilter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue <= val {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}
		return nil
	})
	return f
}

// Before valid whether before specified time.
func (f *TimestampFilter) Before(val uint32) *TimestampFilter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue >= val {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}
		return nil
	})
	return f
}

// Equal valid whether equal to specified time.
func (f *TimestampFilter) Equal(val uint32) *TimestampFilter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue < val {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}
		if paramValue > val {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}

		return nil
	})
	return f
}

// Between valid whether between two times.
func (f *TimestampFilter) Between(start, end uint32) *TimestampFilter {
	f.AddValidator(func(paramName string, paramValue uint32) *Error {
		if paramValue < start {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}
		if paramValue > end {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *TimestampFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
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
		v, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotTimestamp")
}

package filter

import (
	"strings"
	"time"
)

type TimeFilter struct {
	layout     string
	validators []TimeValidator
	allowVals  []string
}

type TimeValidator func(paramName string, paramValue *time.Time) *Error

// Time return a time filter.
func Time() *TimeFilter {
	f := new(TimeFilter)
	f.layout = "2006-01-02"
	return f
}

// Allow allow value is a string in the specified list
func (f *TimeFilter) Allow(vals ...string) *TimeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// HasTime set the layout to include time.
func (f *TimeFilter) HasTime() *TimeFilter {
	f.layout = "2006-01-02 15:04:05"
	return f
}

// Layout set the time layout.
func (f *TimeFilter) Layout(layout string) *TimeFilter {
	f.layout = layout
	return f
}

// AddValidator add a custom validator to filter
func (f *TimeFilter) AddValidator(validator TimeValidator) *TimeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// StartFrom valid whether start from specified time.
func (f *TimeFilter) StartFrom(tm string) *TimeFilter {
	f.AddValidator(func(paramName string, paramValue *time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if paramValue.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}

		return nil
	})
	return f
}

// EndTo valid whether end to specified time.
func (f *TimeFilter) EndTo(tm string) *TimeFilter {
	f.AddValidator(func(paramName string, paramValue *time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if paramValue.After(t) {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}

		return nil
	})
	return f
}

// After valid whether after specified time.
func (f *TimeFilter) After(tm string) *TimeFilter {
	f.AddValidator(func(paramName string, paramValue *time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.After(t) {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}

		return nil
	})
	return f
}

// Before valid whether before specified time.
func (f *TimeFilter) Before(tm string) *TimeFilter {
	f.AddValidator(func(paramName string, paramValue *time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}

		return nil
	})
	return f
}

// Equal valid whether equal to specified time.
func (f *TimeFilter) Equal(tm string) *TimeFilter {
	f.AddValidator(func(paramName string, paramValue *time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if paramValue.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}
		if paramValue.After(t) {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}

		return nil
	})
	return f
}

// Between valid whether between two times.
func (f *TimeFilter) Between(startTime, endTime string) *TimeFilter {
	f.AddValidator(func(paramName string, paramValue *time.Time) *Error {
		startTime, err := time.ParseInLocation(f.layout, startTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}
		endTime, err := time.ParseInLocation(f.layout, endTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if paramValue.Before(startTime) {
			return NewError(ErrorInvalidParam, paramName, "TooEarly")
		}
		if paramValue.After(endTime) {
			return NewError(ErrorInvalidParam, paramName, "TooLate")
		}

		return nil
	})
	return f
}

// Run make the filter running.
func (f *TimeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var timeVal *time.Time
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		t, err := time.ParseInLocation(f.layout, val, timeLoc)
		if err != nil {
			goto parse_error
		}
		timeVal = &t
	case time.Time:
		timeVal = &val
	case *time.Time:
		timeVal = val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, timeVal); err != nil {
			return nil, err
		}
	}

	return timeVal, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotTime")
}

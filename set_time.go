package filter

import (
	"strings"
	"time"

	"github.com/go-apibox/types"
)

type TimeSetFilter struct {
	layout     string
	delimiter  string
	minCount   int
	maxCount   int
	validators []TimeSetValidator
	allowVals  []string
}

type TimeSetValidator func(paramName string, paramValue []*time.Time) *Error

// TimeSet return a timestamp range filter.
func TimeSet() *TimeSetFilter {
	f := new(TimeSetFilter)
	f.delimiter = ","
	f.layout = "2006-01-02"
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *TimeSetFilter) Allow(vals ...string) *TimeSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Delimiter set the delimiter in set string.
func (f *TimeSetFilter) Delimiter(delimiter string) *TimeSetFilter {
	f.delimiter = delimiter
	return f
}

// HasTime set the layout to include time.
func (f *TimeSetFilter) HasTime() *TimeSetFilter {
	f.layout = "2006-01-02 15:04:05"
	return f
}

// Layout set the time layout.
func (f *TimeSetFilter) Layout(layout string) *TimeSetFilter {
	f.layout = layout
	return f
}

// MinCount set the max item count of set.
func (f *TimeSetFilter) MinCount(count int) *TimeSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *TimeSetFilter) MaxCount(count int) *TimeSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *TimeSetFilter) AddValidator(validator TimeSetValidator) *TimeSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemStartFrom valid whether left value in set is start from specified time.
func (f *TimeSetFilter) ItemStartFrom(tm string) *TimeSetFilter {
	f.AddValidator(func(paramName string, paramValue []*time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		for _, v := range paramValue {
			if v.Before(t) {
				return NewError(ErrorInvalidParam, paramName, "ItemTooEarly")
			}
		}

		return nil
	})
	return f
}

// ItemEndTo valid whether left value in set is end to specified time.
func (f *TimeSetFilter) ItemEndTo(tm string) *TimeSetFilter {
	f.AddValidator(func(paramName string, paramValue []*time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		for _, v := range paramValue {
			if v.After(t) {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLate")
			}
		}

		return nil
	})
	return f
}

// ItemAfter valid whether left value in set is after specified time.
func (f *TimeSetFilter) ItemAfter(tm string) *TimeSetFilter {
	f.AddValidator(func(paramName string, paramValue []*time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		for _, v := range paramValue {
			if !v.After(t) {
				return NewError(ErrorInvalidParam, paramName, "ItemTooEarly")
			}
		}

		return nil
	})
	return f
}

// ItemBefore valid whether left value in set is before specified time.
func (f *TimeSetFilter) ItemBefore(tm string) *TimeSetFilter {
	f.AddValidator(func(paramName string, paramValue []*time.Time) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		for _, v := range paramValue {
			if !v.Before(t) {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLate")
			}
		}

		return nil
	})
	return f
}

// ItemBetween valid whether left value in set is in the specified range.
func (f *TimeSetFilter) ItemBetween(startTime, endTime string) *TimeSetFilter {
	f.AddValidator(func(paramName string, paramValue []*time.Time) *Error {
		startTime, err := time.ParseInLocation(f.layout, startTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}
		endTime, err := time.ParseInLocation(f.layout, endTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		for _, v := range paramValue {
			if v.Before(startTime) {
				return NewError(ErrorInvalidParam, paramName, "ItemTooEarly")
			}
			if v.After(endTime) {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLate")
			}
		}

		return nil
	})
	return f
}

// Run make the filter running.
func (f *TimeSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var timeVals []*time.Time
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		if val != "" {
			fields := strings.Split(val, f.delimiter)
			for _, field := range fields {
				field = strings.Trim(field, " \t\r\n")
				t, err := time.ParseInLocation(f.layout, field, timeLoc)
				if err != nil {
					goto parse_error
				}
				timeVals = append(timeVals, &t)
			}
		} else {
			timeVals = []*time.Time{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			t, err := time.ParseInLocation(f.layout, field, timeLoc)
			if err != nil {
				goto parse_error
			}
			timeVals = append(timeVals, &t)
		}
	case []*time.Time:
		timeVals = val
	case []time.Time:
		for _, v := range val {
			timeVals = append(timeVals, &v)
		}
	default:
		goto parse_error
	}

	if len(timeVals) < f.minCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooFew")
	}
	if len(timeVals) > f.maxCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooMany")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, timeVals); err != nil {
			return nil, err
		}
	}

	return timeVals, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotTimeSet")
}

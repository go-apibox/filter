package filter

import (
	"strings"
	"time"

	"github.com/go-apibox/types"
)

type TimeRangeFilter struct {
	layout          string
	delimiter       string
	defaultLeftVal  string
	defaultRightVal string
	validators      []TimeRangeValidator
	allowVals       []string
}

type TimeRangeValidator func(paramName string, paramValue *types.TimeRange) *Error

// TimeRange return a timestamp range filter.
func TimeRange() *TimeRangeFilter {
	f := new(TimeRangeFilter)
	f.layout = "2006-01-02"

	f.defaultLeftVal = time.Unix(0, 0).Format(f.layout)
	f.defaultRightVal = time.Now().Format(f.layout)

	return f
}

// Allow allow value is a string in the specified list
func (f *TimeRangeFilter) Allow(vals ...string) *TimeRangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// HasTime set the layout to include time.
func (f *TimeRangeFilter) HasTime() *TimeRangeFilter {
	f.layout = "2006-01-02 15:04:05"
	return f
}

// Layout set the time layout.
func (f *TimeRangeFilter) Layout(layout string) *TimeRangeFilter {
	f.layout = layout
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *TimeRangeFilter) LeftDefault(tm string) *TimeRangeFilter {
	f.defaultLeftVal = tm
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *TimeRangeFilter) RightDefault(tm string) *TimeRangeFilter {
	f.defaultRightVal = tm
	return f
}

// AddValidator add a custom validator to filter
func (f *TimeRangeFilter) AddValidator(validator TimeRangeValidator) *TimeRangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftStartFrom valid whether left value of range is start from specified time.
func (f *TimeRangeFilter) LeftStartFrom(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.LeftClosed {
			t = t.Add(time.Duration(-1) * time.Second)
		}
		if paramValue.Left.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}

		return nil
	})
	return f
}

// LeftEndTo valid whether left value of range is end to specified time.
func (f *TimeRangeFilter) LeftEndTo(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.LeftClosed {
			t = t.Add(time.Duration(-1) * time.Second)
		}
		if paramValue.Left.After(t) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}

		return nil
	})
	return f
}

// LeftAfter valid whether left value of range is after specified time.
func (f *TimeRangeFilter) LeftAfter(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.LeftClosed {
			t = t.Add(time.Duration(-1) * time.Second)
		}
		if !paramValue.Left.After(t) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}

		return nil
	})
	return f
}

// LeftBefore valid whether left value of range is before specified time.
func (f *TimeRangeFilter) LeftBefore(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.LeftClosed {
			t = t.Add(time.Duration(-1) * time.Second)
		}
		if !paramValue.Left.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}

		return nil
	})
	return f
}

// LeftEqual valid whether left value of range is equal to specified time.
func (f *TimeRangeFilter) LeftEqual(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.LeftClosed {
			t = t.Add(time.Duration(-1) * time.Second)
		}
		if paramValue.Left.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}
		if paramValue.Left.After(t) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}

		return nil
	})
	return f
}

// LeftBetween valid whether left value of range is in the specified range.
func (f *TimeRangeFilter) LeftBetween(startTime, endTime string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		startTime, err := time.ParseInLocation(f.layout, startTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}
		endTime, err := time.ParseInLocation(f.layout, endTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.LeftClosed {
			startTime = startTime.Add(time.Duration(-1) * time.Second)
			endTime = endTime.Add(time.Duration(-1) * time.Second)
		}

		if paramValue.Left.Before(startTime) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}
		if paramValue.Left.After(endTime) {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}

		return nil
	})
	return f
}

// RightStartFrom valid whether right value of range is start from specified time.
func (f *TimeRangeFilter) RightStartFrom(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.RightClosed {
			t = t.Add(time.Duration(1) * time.Second)
		}
		if paramValue.Right.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}

		return nil
	})
	return f
}

// RightEndTo valid whether right value of range is end to specified time.
func (f *TimeRangeFilter) RightEndTo(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.RightClosed {
			t = t.Add(time.Duration(1) * time.Second)
		}
		if paramValue.Right.After(t) {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}

		return nil
	})
	return f
}

// RightAfter valid whether right value of range is after specified time.
func (f *TimeRangeFilter) RightAfter(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.RightClosed {
			t = t.Add(time.Duration(1) * time.Second)
		}
		if !paramValue.Right.After(t) {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}

		return nil
	})
	return f
}

// RightBefore valid whether right value of range is before specified time.
func (f *TimeRangeFilter) RightBefore(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.RightClosed {
			t = t.Add(time.Duration(1) * time.Second)
		}
		if !paramValue.Right.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}

		return nil
	})
	return f
}

// RightEqual valid whether right value of range is equal to specified time.
func (f *TimeRangeFilter) RightEqual(tm string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		t, err := time.ParseInLocation(f.layout, tm, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.RightClosed {
			t = t.Add(time.Duration(1) * time.Second)
		}
		if paramValue.Right.Before(t) {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}
		if paramValue.Right.After(t) {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}

		return nil
	})
	return f
}

// RightBetween valid whether right value of range is in the specified range.
func (f *TimeRangeFilter) RightBetween(startTime, endTime string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		startTime, err := time.ParseInLocation(f.layout, startTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}
		endTime, err := time.ParseInLocation(f.layout, endTime, timeLoc)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		if !paramValue.RightClosed {
			startTime = startTime.Add(time.Duration(1) * time.Second)
			endTime = endTime.Add(time.Duration(1) * time.Second)
		}
		if paramValue.Right.Before(startTime) {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}
		if paramValue.Right.After(endTime) {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}

		return nil
	})
	return f
}

// MinDistance valid whether the distance of range not smaller than the specified value.
func (f *TimeRangeFilter) MinDistance(duration string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		if strings.Contains(duration, "ns") || strings.Contains(duration, "us") || strings.Contains(duration, "ms") {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		valDuration, err := time.ParseDuration(duration)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		dist := paramValue.Right.Sub(*paramValue.Left)
		d := dist.Seconds()
		if !paramValue.LeftClosed {
			d = d - 1
		}
		if !paramValue.RightClosed {
			d = d - 1
		}

		if d < valDuration.Seconds() {
			return NewError(ErrorInvalidParam, paramName, "TooNear")
		}

		return nil
	})
	return f
}

// MaxDistance valid whether the distance of range not larger than the specified value.
func (f *TimeRangeFilter) MaxDistance(duration string) *TimeRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimeRange) *Error {
		if strings.Contains(duration, "ns") || strings.Contains(duration, "us") || strings.Contains(duration, "ms") {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		valDuration, err := time.ParseDuration(duration)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}

		dist := paramValue.Right.Sub(*paramValue.Left)
		d := dist.Seconds()
		if !paramValue.LeftClosed {
			d = d - 1
		}
		if !paramValue.RightClosed {
			d = d - 1
		}

		if d > valDuration.Seconds() {
			return NewError(ErrorInvalidParam, paramName, "TooFar")
		}

		return nil
	})
	return f
}

// Run make the filter running.
func (f *TimeRangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var timeRange *types.TimeRange
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		timeRange, err = types.ParseTimeRange(val, f.layout, timeLoc, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.TimeRange:
		timeRange = val
	case types.TimeRange:
		timeRange = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, timeRange); err != nil {
			return nil, err
		}
	}

	return timeRange, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotTimeRange")
}

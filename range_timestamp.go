package filter

import (
	"math"
	"strings"

	"github.com/go-apibox/types"
)

type TimestampRangeFilter struct {
	defaultLeftVal  uint32
	defaultRightVal uint32
	validators      []TimestampRangeValidator
	allowVals       []string
}

type TimestampRangeValidator func(paramName string, paramValue *types.TimestampRange) *Error

// TimestampRange return a timestamp range filter.
func TimestampRange() *TimestampRangeFilter {
	f := new(TimestampRangeFilter)

	f.defaultLeftVal = 0
	f.defaultRightVal = math.MaxUint32

	return f
}

// Allow allow value is a string in the specified list
func (f *TimestampRangeFilter) Allow(vals ...string) *TimestampRangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *TimestampRangeFilter) LeftDefault(val uint32) *TimestampRangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *TimestampRangeFilter) RightDefault(val uint32) *TimestampRangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *TimestampRangeFilter) AddValidator(validator TimestampRangeValidator) *TimestampRangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftStartFrom valid whether left value of range is start from specified time.
func (f *TimestampRangeFilter) LeftStartFrom(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
				val = val - 1
			}
		}
		if paramValue.Left < val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}
		return nil
	})
	return f
}

// LeftEndTo valid whether left value of range is end to specified time.
func (f *TimestampRangeFilter) LeftEndTo(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
				val = val - 1
			}
		}
		if paramValue.Left > val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}
		return nil
	})
	return f
}

// LeftAfter valid whether left value of range is after specified time.
func (f *TimestampRangeFilter) LeftAfter(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
				val = val - 1
			}
		}
		if paramValue.Left <= val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}
		return nil
	})
	return f
}

// LeftBefore valid whether left value of range is before specified time.
func (f *TimestampRangeFilter) LeftBefore(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
				val = val - 1
			}
		}
		if paramValue.Left >= val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}
		return nil
	})
	return f
}

// LeftEqual valid whether left value of range is equal to specified time.
func (f *TimestampRangeFilter) LeftEqual(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
				val = val - 1
			}
		}
		if paramValue.Left < val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}
		if paramValue.Left > val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}

		return nil
	})
	return f
}

// LeftBetween valid whether left value of range is in the specified range.
func (f *TimestampRangeFilter) LeftBetween(start, end uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.LeftClosed {
			if start > 0 {
				start = start - 1
			}
			if end > 0 {
				end = end - 1
			}
		}
		if paramValue.Left < start {
			return NewError(ErrorInvalidParam, paramName, "LeftTooEarly")
		}
		if paramValue.Left > end {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLate")
		}
		return nil
	})
	return f
}

// RightStartFrom valid whether right value of range is start from specified time.
func (f *TimestampRangeFilter) RightStartFrom(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
				val = val + 1
			}
		}
		if paramValue.Right < val {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}
		return nil
	})
	return f
}

// RightEndTo valid whether right value of range is end to specified time.
func (f *TimestampRangeFilter) RightEndTo(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
				val = val + 1
			}
		}
		if paramValue.Right > val {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}
		return nil
	})
	return f
}

// RightAfter valid whether right value of range is after specified time.
func (f *TimestampRangeFilter) RightAfter(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
				val = val + 1
			}
		}
		if paramValue.Right <= val {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}
		return nil
	})
	return f
}

// RightBefore valid whether right value of range is before specified time.
func (f *TimestampRangeFilter) RightBefore(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
				val = val + 1
			}
		}
		if paramValue.Right >= val {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}
		return nil
	})
	return f
}

// RightEqual valid whether right value of range is equal to specified time.
func (f *TimestampRangeFilter) RightEqual(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
				val = val + 1
			}
		}
		if paramValue.Right < val {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}
		if paramValue.Right > val {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}

		return nil
	})
	return f
}

// RightBetween valid whether right value of range is in the specified range.
func (f *TimestampRangeFilter) RightBetween(start, end uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		if !paramValue.RightClosed {
			if start < math.MaxUint32 {
				start = start + 1
			}
			if start < math.MaxUint32 {
				start = start + 1
			}
		}
		if paramValue.Right < start {
			return NewError(ErrorInvalidParam, paramName, "RightTooEarly")
		}
		if paramValue.Right > end {
			return NewError(ErrorInvalidParam, paramName, "RightTooLate")
		}
		return nil
	})
	return f
}

// MinDistance valid whether the distance of range not smaller than the specified value.
func (f *TimestampRangeFilter) MinDistance(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		dist := uint64(paramValue.Right - paramValue.Left)
		if !paramValue.LeftClosed {
			if dist > 0 {
				dist = dist - 1
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}
		if !paramValue.RightClosed {
			if dist > 0 {
				dist = dist - 1
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}
		if dist < uint64(val) {
			return NewError(ErrorInvalidParam, paramName, "TooNear")
		}

		return nil
	})
	return f
}

// MaxDistance valid whether the distance of range not larger than the specified value.
func (f *TimestampRangeFilter) MaxDistance(val uint32) *TimestampRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.TimestampRange) *Error {
		dist := uint64(paramValue.Right - paramValue.Left)
		if !paramValue.LeftClosed {
			if dist > 0 {
				dist = dist - 1
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}
		if !paramValue.RightClosed {
			if dist > 0 {
				dist = dist - 1
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}
		if dist > uint64(val) {
			return NewError(ErrorInvalidParam, paramName, "TooFar")
		}

		return nil
	})
	return f
}

// Run make the filter running.
func (f *TimestampRangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var tsRange *types.TimestampRange
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		tsRange, err = types.ParseTimestampRange(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.TimestampRange:
		tsRange = val
	case types.TimestampRange:
		tsRange = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, tsRange); err != nil {
			return nil, err
		}
	}

	return tsRange, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotTimestampRange")
}

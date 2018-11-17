package filter

import (
	"math"
	"strings"

	"github.com/go-apibox/types"
)

type Int32RangeFilter struct {
	defaultLeftVal  int32
	defaultRightVal int32
	validators      []Int32RangeValidator
	allowVals       []string
}

type Int32RangeValidator func(paramName string, paramValue *types.Int32Range) *Error

// Int32Range return a int32 range filter.
func Int32Range() *Int32RangeFilter {
	f := new(Int32RangeFilter)

	f.defaultLeftVal = math.MinInt32
	f.defaultRightVal = math.MaxInt32

	return f
}

// Allow allow value is a string in the specified list
func (f *Int32RangeFilter) Allow(vals ...string) *Int32RangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *Int32RangeFilter) LeftDefault(val int32) *Int32RangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *Int32RangeFilter) RightDefault(val int32) *Int32RangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *Int32RangeFilter) AddValidator(validator Int32RangeValidator) *Int32RangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftMin valid whether left value of range is not smaller than specified value.
func (f *Int32RangeFilter) LeftMin(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt32 {
				val = val - 1
			}
		}
		if paramValue.Left < val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooSmall")
		}
		return nil
	})
	return f
}

// LeftMax valid whether left value of range is not larger than specified value.
func (f *Int32RangeFilter) LeftMax(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt32 {
				val = val - 1
			}
		}
		if paramValue.Left > val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLarge")
		}
		return nil
	})
	return f
}

// LeftLargerThan valid whether left value of range is larger than the specified value.
func (f *Int32RangeFilter) LeftLargerThan(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt32 {
				val = val - 1
			}
		}
		if paramValue.Left <= val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooSmall")
		}
		return nil
	})
	return f
}

// LeftSmallerThan valid whether left value of range is smaller than the specified value.
func (f *Int32RangeFilter) LeftSmallerThan(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt32 {
				val = val - 1
			}
		}
		if paramValue.Left >= val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLarge")
		}
		return nil
	})
	return f
}

// LeftEqual valid whether left value of range is equal to the specified value.
func (f *Int32RangeFilter) LeftEqual(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt32 {
				val = val - 1
			}
		}
		if paramValue.Left < val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooSmall")
		}
		if paramValue.Left > val {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLarge")
		}

		return nil
	})
	return f
}

// LeftBetween valid whether left value of range is in the specified range.
func (f *Int32RangeFilter) LeftBetween(min, max int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.LeftClosed {
			if min > math.MinInt32 {
				min = min - 1
			}
			if max > math.MinInt32 {
				max = max - 1
			}
		}
		if paramValue.Left < min {
			return NewError(ErrorInvalidParam, paramName, "LeftTooSmall")
		}
		if paramValue.Left > max {
			return NewError(ErrorInvalidParam, paramName, "LeftTooLarge")
		}
		return nil
	})
	return f
}

// RightMin valid whether right value of range is not smaller than specified value.
func (f *Int32RangeFilter) RightMin(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt32 {
				val = val + 1
			}
		}
		if paramValue.Right < val {
			return NewError(ErrorInvalidParam, paramName, "RightTooSmall")
		}
		return nil
	})
	return f
}

// RightMax valid whether right value of range is not larger than specified value.
func (f *Int32RangeFilter) RightMax(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt32 {
				val = val + 1
			}
		}
		if paramValue.Right > val {
			return NewError(ErrorInvalidParam, paramName, "RightTooLarge")
		}
		return nil
	})
	return f
}

// RightLargerThan valid whether right value of range is larger than the specified value.
func (f *Int32RangeFilter) RightLargerThan(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt32 {
				val = val + 1
			}
		}
		if paramValue.Right <= val {
			return NewError(ErrorInvalidParam, paramName, "RightTooSmall")
		}
		return nil
	})
	return f
}

// RightSmallerThan valid whether right value of range is smaller than the specified value.
func (f *Int32RangeFilter) RightSmallerThan(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt32 {
				val = val + 1
			}
		}
		if paramValue.Right >= val {
			return NewError(ErrorInvalidParam, paramName, "RightTooLarge")
		}
		return nil
	})
	return f
}

// RightEqual valid whether right value of range is equal to the specified value.
func (f *Int32RangeFilter) RightEqual(val int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt32 {
				val = val + 1
			}
		}
		if paramValue.Right < val {
			return NewError(ErrorInvalidParam, paramName, "RightTooSmall")
		}
		if paramValue.Right > val {
			return NewError(ErrorInvalidParam, paramName, "RightTooLarge")
		}

		return nil
	})
	return f
}

// RightBetween valid whether right value of range is in the specified range.
func (f *Int32RangeFilter) RightBetween(min, max int32) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		if !paramValue.RightClosed {
			if min < math.MaxInt32 {
				min = min + 1
			}
			if max < math.MaxInt32 {
				max = max + 1
			}
		}
		if paramValue.Right < min {
			return NewError(ErrorInvalidParam, paramName, "RightTooSmall")
		}
		if paramValue.Right > max {
			return NewError(ErrorInvalidParam, paramName, "RightTooLarge")
		}
		return nil
	})
	return f
}

// MinDistance valid whether the distance of range not smaller than the specified value.
func (f *Int32RangeFilter) MinDistance(val uint) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		dist := uint64(int64(paramValue.Right) - int64(paramValue.Left))
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
func (f *Int32RangeFilter) MaxDistance(val uint) *Int32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int32Range) *Error {
		dist := uint64(int64(paramValue.Right) - int64(paramValue.Left))
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
func (f *Int32RangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var int32Range *types.Int32Range
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		int32Range, err = types.ParseInt32Range(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.Int32Range:
		int32Range = val
	case types.Int32Range:
		int32Range = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, int32Range); err != nil {
			return nil, err
		}
	}

	return int32Range, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt32Range")
}

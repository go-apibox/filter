package filter

import (
	"strings"

	"github.com/go-apibox/types"
)

type IntRangeFilter struct {
	defaultLeftVal  int
	defaultRightVal int
	validators      []IntRangeValidator
	allowVals       []string
}

type IntRangeValidator func(paramName string, paramValue *types.IntRange) *Error

// IntRange return a int range filter.
func IntRange() *IntRangeFilter {
	f := new(IntRangeFilter)

	f.defaultLeftVal = types.MinInt
	f.defaultRightVal = types.MaxInt

	return f
}

// Allow allow value is a string in the specified list
func (f *IntRangeFilter) Allow(vals ...string) *IntRangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *IntRangeFilter) LeftDefault(val int) *IntRangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *IntRangeFilter) RightDefault(val int) *IntRangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *IntRangeFilter) AddValidator(validator IntRangeValidator) *IntRangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftMin valid whether left value of range is not smaller than specified value.
func (f *IntRangeFilter) LeftMin(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.LeftClosed {
			if val > types.MinInt {
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
func (f *IntRangeFilter) LeftMax(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.LeftClosed {
			if val > types.MinInt {
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
func (f *IntRangeFilter) LeftLargerThan(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.LeftClosed {
			if val > types.MinInt {
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
func (f *IntRangeFilter) LeftSmallerThan(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.LeftClosed {
			if val > types.MinInt {
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
func (f *IntRangeFilter) LeftEqual(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.LeftClosed {
			if val > types.MinInt {
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
func (f *IntRangeFilter) LeftBetween(min, max int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.LeftClosed {
			if min > types.MinInt {
				min = min - 1
			}
			if max > types.MinInt {
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
func (f *IntRangeFilter) RightMin(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxInt {
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
func (f *IntRangeFilter) RightMax(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxInt {
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
func (f *IntRangeFilter) RightLargerThan(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxInt {
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
func (f *IntRangeFilter) RightSmallerThan(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxInt {
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
func (f *IntRangeFilter) RightEqual(val int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxInt {
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
func (f *IntRangeFilter) RightBetween(min, max int) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
		if !paramValue.RightClosed {
			if min < types.MaxInt {
				min = min + 1
			}
			if max < types.MaxInt {
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
func (f *IntRangeFilter) MinDistance(val uint) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
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
func (f *IntRangeFilter) MaxDistance(val uint) *IntRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.IntRange) *Error {
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
func (f *IntRangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intRange *types.IntRange
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		intRange, err = types.ParseIntRange(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.IntRange:
		intRange = val
	case types.IntRange:
		intRange = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, intRange); err != nil {
			return nil, err
		}
	}

	return intRange, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotIntRange")
}

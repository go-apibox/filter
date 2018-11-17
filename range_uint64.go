package filter

import (
	"math"
	"strings"

	"github.com/go-apibox/types"
)

type Uint64RangeFilter struct {
	defaultLeftVal  uint64
	defaultRightVal uint64
	validators      []Uint64RangeValidator
	allowVals       []string
}

type Uint64RangeValidator func(paramName string, paramValue *types.Uint64Range) *Error

// Uint64Range return a int64 filter.
func Uint64Range() *Uint64RangeFilter {
	f := new(Uint64RangeFilter)

	f.defaultLeftVal = 0
	f.defaultRightVal = math.MaxUint64

	return f
}

// Allow allow value is a string in the specified list
func (f *Uint64RangeFilter) Allow(vals ...string) *Uint64RangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *Uint64RangeFilter) LeftDefault(val uint64) *Uint64RangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *Uint64RangeFilter) RightDefault(val uint64) *Uint64RangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *Uint64RangeFilter) AddValidator(validator Uint64RangeValidator) *Uint64RangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftMin valid whether left value of range is not smaller than specified value.
func (f *Uint64RangeFilter) LeftMin(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
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
func (f *Uint64RangeFilter) LeftMax(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
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
func (f *Uint64RangeFilter) LeftLargerThan(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
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
func (f *Uint64RangeFilter) LeftSmallerThan(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
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
func (f *Uint64RangeFilter) LeftEqual(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.LeftClosed {
			if val > 0 {
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
func (f *Uint64RangeFilter) LeftBetween(min, max uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.LeftClosed {
			if min > 0 {
				min = min - 1
			}
			if max > 0 {
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
func (f *Uint64RangeFilter) RightMin(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint64 {
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
func (f *Uint64RangeFilter) RightMax(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint64 {
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
func (f *Uint64RangeFilter) RightLargerThan(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint64 {
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
func (f *Uint64RangeFilter) RightSmallerThan(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint64 {
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
func (f *Uint64RangeFilter) RightEqual(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint64 {
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
func (f *Uint64RangeFilter) RightBetween(min, max uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
		if !paramValue.RightClosed {
			if min < math.MaxUint64 {
				min = min + 1
			}
			if max < math.MaxUint64 {
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
func (f *Uint64RangeFilter) MinDistance(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
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
		if dist < val {
			return NewError(ErrorInvalidParam, paramName, "TooNear")
		}

		return nil
	})
	return f
}

// MaxDistance valid whether the distance of range not larger than the specified value.
func (f *Uint64RangeFilter) MaxDistance(val uint64) *Uint64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint64Range) *Error {
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
		if dist > val {
			return NewError(ErrorInvalidParam, paramName, "TooFar")
		}

		return nil
	})
	return f
}

// Run make the filter running.
func (f *Uint64RangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var uint64Range *types.Uint64Range
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		uint64Range, err = types.ParseUint64Range(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.Uint64Range:
		uint64Range = val
	case types.Uint64Range:
		uint64Range = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, uint64Range); err != nil {
			return nil, err
		}
	}

	return uint64Range, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotUint64Range")
}

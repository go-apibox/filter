package filter

import (
	"math"
	"strings"

	"github.com/go-apibox/types"
)

type Uint32RangeFilter struct {
	defaultLeftVal  uint32
	defaultRightVal uint32
	validators      []Uint32RangeValidator
	allowVals       []string
}

type Uint32RangeValidator func(paramName string, paramValue *types.Uint32Range) *Error

// Uint32Range return a uint32 range filter.
func Uint32Range() *Uint32RangeFilter {
	f := new(Uint32RangeFilter)

	f.defaultLeftVal = 0
	f.defaultRightVal = math.MaxUint32

	return f
}

// Allow allow value is a string in the specified list
func (f *Uint32RangeFilter) Allow(vals ...string) *Uint32RangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *Uint32RangeFilter) LeftDefault(val uint32) *Uint32RangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *Uint32RangeFilter) RightDefault(val uint32) *Uint32RangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *Uint32RangeFilter) AddValidator(validator Uint32RangeValidator) *Uint32RangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftMin valid whether left value of range is not smaller than specified value.
func (f *Uint32RangeFilter) LeftMin(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) LeftMax(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) LeftLargerThan(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) LeftSmallerThan(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) LeftEqual(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) LeftBetween(min, max uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) RightMin(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
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
func (f *Uint32RangeFilter) RightMax(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
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
func (f *Uint32RangeFilter) RightLargerThan(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
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
func (f *Uint32RangeFilter) RightSmallerThan(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
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
func (f *Uint32RangeFilter) RightEqual(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxUint32 {
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
func (f *Uint32RangeFilter) RightBetween(min, max uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
		if !paramValue.RightClosed {
			if min < math.MaxUint32 {
				min = min + 1
			}
			if max < math.MaxUint32 {
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
func (f *Uint32RangeFilter) MinDistance(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) MaxDistance(val uint32) *Uint32RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Uint32Range) *Error {
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
func (f *Uint32RangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var uint32Range *types.Uint32Range
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		uint32Range, err = types.ParseUint32Range(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.Uint32Range:
		uint32Range = val
	case types.Uint32Range:
		uint32Range = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, uint32Range); err != nil {
			return nil, err
		}
	}

	return uint32Range, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotUint32Range")
}

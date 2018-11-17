package filter

import (
	"strings"

	"github.com/go-apibox/types"
)

type UintRangeFilter struct {
	defaultLeftVal  uint
	defaultRightVal uint
	validators      []UintRangeValidator
	allowVals       []string
}

type UintRangeValidator func(paramName string, paramValue *types.UintRange) *Error

// UintRange return a uint range filter.
func UintRange() *UintRangeFilter {
	f := new(UintRangeFilter)

	f.defaultLeftVal = 0
	f.defaultRightVal = types.MaxUint

	return f
}

// Allow allow value is a string in the specified list
func (f *UintRangeFilter) Allow(vals ...string) *UintRangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *UintRangeFilter) LeftDefault(val uint) *UintRangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *UintRangeFilter) RightDefault(val uint) *UintRangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *UintRangeFilter) AddValidator(validator UintRangeValidator) *UintRangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftMin valid whether left value of range is not smaller than specified value.
func (f *UintRangeFilter) LeftMin(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) LeftMax(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) LeftLargerThan(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) LeftSmallerThan(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) LeftEqual(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) LeftBetween(min, max uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) RightMin(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxUint {
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
func (f *UintRangeFilter) RightMax(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxUint {
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
func (f *UintRangeFilter) RightLargerThan(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxUint {
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
func (f *UintRangeFilter) RightSmallerThan(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxUint {
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
func (f *UintRangeFilter) RightEqual(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
		if !paramValue.RightClosed {
			if val < types.MaxUint {
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
func (f *UintRangeFilter) RightBetween(min, max uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
		if !paramValue.RightClosed {
			if min < types.MaxUint {
				min = min + 1
			}
			if max < types.MaxUint {
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
func (f *UintRangeFilter) MinDistance(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) MaxDistance(val uint) *UintRangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.UintRange) *Error {
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
func (f *UintRangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var uintRange *types.UintRange
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		uintRange, err = types.ParseUintRange(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.UintRange:
		uintRange = val
	case types.UintRange:
		uintRange = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, uintRange); err != nil {
			return nil, err
		}
	}

	return uintRange, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotUintRange")
}

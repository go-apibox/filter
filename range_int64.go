package filter

import (
	"math"
	"math/big"
	"strings"

	"github.com/go-apibox/types"
)

type Int64RangeFilter struct {
	defaultLeftVal  int64
	defaultRightVal int64
	validators      []Int64RangeValidator
	allowVals       []string
}

type Int64RangeValidator func(paramName string, paramValue *types.Int64Range) *Error

// Int64Range return a int64 range filter.
func Int64Range() *Int64RangeFilter {
	f := new(Int64RangeFilter)

	f.defaultLeftVal = math.MinInt64
	f.defaultRightVal = math.MaxInt64

	return f
}

// Allow allow value is a string in the specified list
func (f *Int64RangeFilter) Allow(vals ...string) *Int64RangeFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// LeftDefault set the default left value of range if not specified.
func (f *Int64RangeFilter) LeftDefault(val int64) *Int64RangeFilter {
	f.defaultLeftVal = val
	return f
}

// RightDefault set the default right value of range if not specified.
func (f *Int64RangeFilter) RightDefault(val int64) *Int64RangeFilter {
	f.defaultRightVal = val
	return f
}

// AddValidator add a custom validator to filter
func (f *Int64RangeFilter) AddValidator(validator Int64RangeValidator) *Int64RangeFilter {
	f.validators = append(f.validators, validator)
	return f
}

// LeftMin valid whether left value of range is not smaller than specified value.
func (f *Int64RangeFilter) LeftMin(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt64 {
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
func (f *Int64RangeFilter) LeftMax(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt64 {
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
func (f *Int64RangeFilter) LeftLargerThan(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt64 {
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
func (f *Int64RangeFilter) LeftSmallerThan(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt64 {
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
func (f *Int64RangeFilter) LeftEqual(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.LeftClosed {
			if val > math.MinInt64 {
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
func (f *Int64RangeFilter) LeftBetween(min, max int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.LeftClosed {
			if min > math.MinInt64 {
				min = min - 1
			}
			if max > math.MinInt64 {
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
func (f *Int64RangeFilter) RightMin(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt64 {
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
func (f *Int64RangeFilter) RightMax(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt64 {
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
func (f *Int64RangeFilter) RightLargerThan(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt64 {
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
func (f *Int64RangeFilter) RightSmallerThan(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt64 {
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
func (f *Int64RangeFilter) RightEqual(val int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.RightClosed {
			if val < math.MaxInt64 {
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
func (f *Int64RangeFilter) RightBetween(min, max int64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		if !paramValue.RightClosed {
			if min < math.MaxInt64 {
				min = min + 1
			}
			if max < math.MaxInt64 {
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
func (f *Int64RangeFilter) MinDistance(val uint64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		left := big.NewInt(paramValue.Left)
		right := big.NewInt(paramValue.Right)
		dist := big.NewInt(0)
		dist.Sub(right, left)

		if !paramValue.LeftClosed {
			if dist.Uint64() > 0 {
				dist.Sub(dist, big.NewInt(1))
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}
		if !paramValue.RightClosed {
			if dist.Uint64() > 0 {
				dist.Sub(dist, big.NewInt(1))
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}

		if dist.Uint64() < val {
			return NewError(ErrorInvalidParam, paramName, "TooNear")
		}

		return nil
	})
	return f
}

// MaxDistance valid whether the distance of range not larger than the specified value.
func (f *Int64RangeFilter) MaxDistance(val uint64) *Int64RangeFilter {
	f.AddValidator(func(paramName string, paramValue *types.Int64Range) *Error {
		left := big.NewInt(paramValue.Left)
		right := big.NewInt(paramValue.Right)
		dist := big.NewInt(0)
		dist.Sub(right, left)

		if !paramValue.LeftClosed {
			if dist.Uint64() > 0 {
				dist.Sub(dist, big.NewInt(1))
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}
		if !paramValue.RightClosed {
			if dist.Uint64() > 0 {
				dist.Sub(dist, big.NewInt(1))
			} else {
				return NewError(ErrorInvalidParam, paramName, "WrongRange")
			}
		}

		if dist.Uint64() > val {
			return NewError(ErrorInvalidParam, paramName, "TooFar")
		}

		return nil
	})
	return f
}

// Run make the filter running.
func (f *Int64RangeFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var int64Range *types.Int64Range
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		var err error
		int64Range, err = types.ParseInt64Range(val, f.defaultLeftVal, f.defaultRightVal)
		if err != nil {
			goto parse_error
		}
	case *types.Int64Range:
		int64Range = val
	case types.Int64Range:
		int64Range = &val
	default:
		goto parse_error
	}

	for _, validator := range f.validators {
		if err := validator(paramName, int64Range); err != nil {
			return nil, err
		}
	}

	return int64Range, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt64Range")
}

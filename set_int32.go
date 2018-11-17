package filter

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-apibox/types"
)

type Int32SetFilter struct {
	base       int
	delimiter  string
	minCount   int
	maxCount   int
	validators []Int32SetValidator
	allowVals  []string
}

type Int32SetValidator func(paramName string, paramValue []int32) *Error

// Int32Set return a int32 set filter.
func Int32Set() *Int32SetFilter {
	f := new(Int32SetFilter)
	f.base = 10
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *Int32SetFilter) Allow(vals ...string) *Int32SetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of int32.
// Int32SetFilter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *Int32SetFilter) Base(base int) *Int32SetFilter {
	f.base = base
	return f
}

// Delimiter set the delimiter of set string.
func (f *Int32SetFilter) Delimiter(delimiter string) *Int32SetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *Int32SetFilter) MinCount(count int) *Int32SetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *Int32SetFilter) MaxCount(count int) *Int32SetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *Int32SetFilter) AddValidator(validator Int32SetValidator) *Int32SetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemMin valid whether item value of set is not smaller than specified value.
func (f *Int32SetFilter) ItemMin(val int32) *Int32SetFilter {
	f.AddValidator(func(paramName string, paramValue []int32) *Error {
		for _, v := range paramValue {
			if v < val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooSmall")
			}
		}
		return nil
	})
	return f
}

// ItemMax valid whether item value of set is not larger than specified value.
func (f *Int32SetFilter) ItemMax(val int32) *Int32SetFilter {
	f.AddValidator(func(paramName string, paramValue []int32) *Error {
		for _, v := range paramValue {
			if v > val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLarge")
			}
		}
		return nil
	})
	return f
}

// ItemLargerThan valid whether item value of set is larger than the specified value.
func (f *Int32SetFilter) ItemLargerThan(val int32) *Int32SetFilter {
	f.AddValidator(func(paramName string, paramValue []int32) *Error {
		for _, v := range paramValue {
			if v <= val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooSmall")
			}
		}
		return nil
	})
	return f
}

// ItemSmallerThan valid whether item value of set is smaller than the specified value.
func (f *Int32SetFilter) ItemSmallerThan(val int32) *Int32SetFilter {
	f.AddValidator(func(paramName string, paramValue []int32) *Error {
		for _, v := range paramValue {
			if v >= val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLarge")
			}
		}
		return nil
	})
	return f
}

// ItemBetween valid whether item value of set is in the specified set.
func (f *Int32SetFilter) ItemBetween(min, max int32) *Int32SetFilter {
	f.AddValidator(func(paramName string, paramValue []int32) *Error {
		for _, v := range paramValue {
			if v < min {
				return NewError(ErrorInvalidParam, paramName, "ItemTooSmall")
			}
			if v > max {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLarge")
			}
		}
		return nil
	})
	return f
}

// ItemIn valid item value of set should in the specified set.
func (f *Int32SetFilter) ItemIn(set []int32) *Int32SetFilter {
	f.AddValidator(func(paramName string, paramValue []int32) *Error {
		for _, item := range paramValue {
			itemFound := false
			for _, v := range set {
				if v == item {
					itemFound = true
					break
				}
			}
			if !itemFound {
				return NewError(ErrorInvalidParam, paramName, "ItemNotInSet")
			}
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *Int32SetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVals []int32
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
				v, err := strconv.ParseInt(field, f.base, 64)
				if err != nil {
					goto parse_error
				}
				if v > int64(math.MaxInt32) || v < int64(math.MinInt32) {
					goto parse_error
				}
				intVals = append(intVals, int32(v))
			}
		} else {
			intVals = []int32{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			v, err := strconv.ParseInt(field, f.base, 64)
			if err != nil {
				goto parse_error
			}
			if v > int64(math.MaxInt32) || v < int64(math.MinInt32) {
				goto parse_error
			}
			intVals = append(intVals, int32(v))
		}
	case []int32:
		intVals = val
	default:
		goto parse_error
	}

	if len(intVals) < f.minCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooFew")
	}
	if len(intVals) > f.maxCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooMany")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, intVals); err != nil {
			return nil, err
		}
	}

	return intVals, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt32Set")
}

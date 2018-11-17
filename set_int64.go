package filter

import (
	"strconv"
	"strings"

	"github.com/go-apibox/types"
)

type Int64SetFilter struct {
	base       int
	delimiter  string
	minCount   int
	maxCount   int
	validators []Int64SetValidator
	allowVals  []string
}

type Int64SetValidator func(paramName string, paramValue []int64) *Error

// Int64Set return a int64 set filter.
func Int64Set() *Int64SetFilter {
	f := new(Int64SetFilter)
	f.base = 10
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *Int64SetFilter) Allow(vals ...string) *Int64SetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of int64.
// Int64SetFilter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *Int64SetFilter) Base(base int) *Int64SetFilter {
	f.base = base
	return f
}

// Delimiter set the delimiter of set string.
func (f *Int64SetFilter) Delimiter(delimiter string) *Int64SetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *Int64SetFilter) MinCount(count int) *Int64SetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *Int64SetFilter) MaxCount(count int) *Int64SetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *Int64SetFilter) AddValidator(validator Int64SetValidator) *Int64SetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemMin valid whether item value of set is not smaller than specified value.
func (f *Int64SetFilter) ItemMin(val int64) *Int64SetFilter {
	f.AddValidator(func(paramName string, paramValue []int64) *Error {
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
func (f *Int64SetFilter) ItemMax(val int64) *Int64SetFilter {
	f.AddValidator(func(paramName string, paramValue []int64) *Error {
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
func (f *Int64SetFilter) ItemLargerThan(val int64) *Int64SetFilter {
	f.AddValidator(func(paramName string, paramValue []int64) *Error {
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
func (f *Int64SetFilter) ItemSmallerThan(val int64) *Int64SetFilter {
	f.AddValidator(func(paramName string, paramValue []int64) *Error {
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
func (f *Int64SetFilter) ItemBetween(min, max int64) *Int64SetFilter {
	f.AddValidator(func(paramName string, paramValue []int64) *Error {
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
func (f *Int64SetFilter) ItemIn(set []int64) *Int64SetFilter {
	f.AddValidator(func(paramName string, paramValue []int64) *Error {
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
func (f *Int64SetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVals []int64
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
				intVals = append(intVals, v)
			}
		} else {
			intVals = []int64{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			v, err := strconv.ParseInt(field, f.base, 64)
			if err != nil {
				goto parse_error
			}
			intVals = append(intVals, v)
		}
	case []int64:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotInt64Set")
}

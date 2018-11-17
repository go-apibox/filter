package filter

import (
	"strconv"
	"strings"

	"github.com/go-apibox/types"
)

type IntSetFilter struct {
	base       int
	delimiter  string
	minCount   int
	maxCount   int
	validators []IntSetValidator
	allowVals  []string
}

type IntSetValidator func(paramName string, paramValue []int) *Error

// IntSet return a int set filter.
func IntSet() *IntSetFilter {
	f := new(IntSetFilter)
	f.base = 10
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *IntSetFilter) Allow(vals ...string) *IntSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of int.
// IntSetFilter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *IntSetFilter) Base(base int) *IntSetFilter {
	f.base = base
	return f
}

// Delimiter set the delimiter of set string.
func (f *IntSetFilter) Delimiter(delimiter string) *IntSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *IntSetFilter) MinCount(count int) *IntSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *IntSetFilter) MaxCount(count int) *IntSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *IntSetFilter) AddValidator(validator IntSetValidator) *IntSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemMin valid whether item value of set is not smaller than specified value.
func (f *IntSetFilter) ItemMin(val int) *IntSetFilter {
	f.AddValidator(func(paramName string, paramValue []int) *Error {
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
func (f *IntSetFilter) ItemMax(val int) *IntSetFilter {
	f.AddValidator(func(paramName string, paramValue []int) *Error {
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
func (f *IntSetFilter) ItemLargerThan(val int) *IntSetFilter {
	f.AddValidator(func(paramName string, paramValue []int) *Error {
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
func (f *IntSetFilter) ItemSmallerThan(val int) *IntSetFilter {
	f.AddValidator(func(paramName string, paramValue []int) *Error {
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
func (f *IntSetFilter) ItemBetween(min, max int) *IntSetFilter {
	f.AddValidator(func(paramName string, paramValue []int) *Error {
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
func (f *IntSetFilter) ItemIn(set []int) *IntSetFilter {
	f.AddValidator(func(paramName string, paramValue []int) *Error {
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
func (f *IntSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVals []int
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
				v, err := strconv.ParseInt(field, f.base, 0)
				if err != nil {
					goto parse_error
				}
				if v > int64(types.MaxInt) || v < int64(types.MinInt) {
					goto parse_error
				}
				intVals = append(intVals, int(v))
			}
		} else {
			intVals = []int{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			v, err := strconv.ParseInt(field, f.base, 0)
			if err != nil {
				goto parse_error
			}
			if v > int64(types.MaxInt) || v < int64(types.MinInt) {
				goto parse_error
			}
			intVals = append(intVals, int(v))
		}
	case []int:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotIntSet")
}

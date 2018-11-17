package filter

import (
	"strconv"
	"strings"

	"github.com/go-apibox/types"
)

type UintSetFilter struct {
	base       int
	delimiter  string
	minCount   int
	maxCount   int
	validators []UintSetValidator
	allowVals  []string
}

type UintSetValidator func(paramName string, paramValue []uint) *Error

// UintSet return a uint set filter.
func UintSet() *UintSetFilter {
	f := new(UintSetFilter)
	f.base = 10
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *UintSetFilter) Allow(vals ...string) *UintSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Base set the base of uint.
// UintSetFilter interprets a string s in the given base (2 to 36) and returns
// the corresponding value i. If base == 0, the base is implied by the
// string's prefix: base 16 for "0x", base 8 for "0", and base 10 otherwise.
func (f *UintSetFilter) Base(base int) *UintSetFilter {
	f.base = base
	return f
}

// Delimiter set the delimiter of set string.
func (f *UintSetFilter) Delimiter(delimiter string) *UintSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *UintSetFilter) MinCount(count int) *UintSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *UintSetFilter) MaxCount(count int) *UintSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *UintSetFilter) AddValidator(validator UintSetValidator) *UintSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemMin valid whether item value of set is not smaller than specified value.
func (f *UintSetFilter) ItemMin(val uint) *UintSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint) *Error {
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
func (f *UintSetFilter) ItemMax(val uint) *UintSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint) *Error {
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
func (f *UintSetFilter) ItemLargerThan(val uint) *UintSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint) *Error {
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
func (f *UintSetFilter) ItemSmallerThan(val uint) *UintSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint) *Error {
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
func (f *UintSetFilter) ItemBetween(min, max uint) *UintSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint) *Error {
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
func (f *UintSetFilter) ItemIn(set []uint) *UintSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint) *Error {
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
func (f *UintSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var intVals []uint
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
				v, err := strconv.ParseUint(field, f.base, 0)
				if err != nil {
					goto parse_error
				}
				if v > uint64(types.MaxUint) {
					goto parse_error
				}
				intVals = append(intVals, uint(v))
			}
		} else {
			intVals = []uint{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			v, err := strconv.ParseUint(field, f.base, 0)
			if err != nil {
				goto parse_error
			}
			if v > uint64(types.MaxUint) {
				goto parse_error
			}
			intVals = append(intVals, uint(v))
		}
	case []uint:
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotUintSet")
}

package filter

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-apibox/types"
)

type TimestampSetFilter struct {
	delimiter  string
	minCount   int
	maxCount   int
	validators []TimestampSetValidator
	allowVals  []string
}

type TimestampSetValidator func(paramName string, paramValue []uint32) *Error

// TimestampSet return a timestamp range filter.
func TimestampSet() *TimestampSetFilter {
	f := new(TimestampSetFilter)
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *TimestampSetFilter) Allow(vals ...string) *TimestampSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// Delimiter set the delimiter of set string.
func (f *TimestampSetFilter) Delimiter(delimiter string) *TimestampSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *TimestampSetFilter) MinCount(count int) *TimestampSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *TimestampSetFilter) MaxCount(count int) *TimestampSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *TimestampSetFilter) AddValidator(validator TimestampSetValidator) *TimestampSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// ItemStartFrom valid whether item value in set is start from specified time.
func (f *TimestampSetFilter) ItemStartFrom(val uint32) *TimestampSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint32) *Error {
		for _, v := range paramValue {
			if v < val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooEarly")
			}
		}
		return nil
	})
	return f
}

// ItemEndTo valid whether item value in set is end to specified time.
func (f *TimestampSetFilter) ItemEndTo(val uint32) *TimestampSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint32) *Error {
		for _, v := range paramValue {
			if v > val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLate")
			}
		}
		return nil
	})
	return f
}

// ItemAfter valid whether item value in set is after specified time.
func (f *TimestampSetFilter) ItemAfter(val uint32) *TimestampSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint32) *Error {
		for _, v := range paramValue {
			if v <= val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooEarly")
			}
		}
		return nil
	})
	return f
}

// ItemBefore valid whether item value in set is before specified time.
func (f *TimestampSetFilter) ItemBefore(val uint32) *TimestampSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint32) *Error {
		for _, v := range paramValue {
			if v >= val {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLate")
			}
		}
		return nil
	})
	return f
}

// ItemBetween valid whether item value in set is in the specified range.
func (f *TimestampSetFilter) ItemBetween(start, end uint32) *TimestampSetFilter {
	f.AddValidator(func(paramName string, paramValue []uint32) *Error {
		for _, v := range paramValue {
			if v < start {
				return NewError(ErrorInvalidParam, paramName, "ItemTooEarly")
			}
			if v > end {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLate")
			}
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *TimestampSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var tsVals []uint32
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
				v, err := strconv.ParseUint(field, 10, 0)
				if err != nil {
					goto parse_error
				}
				if v > uint64(math.MaxUint32) {
					goto parse_error
				}
				tsVals = append(tsVals, uint32(v))
			}
		} else {
			tsVals = []uint32{}
		}
	case []string:
		fields := val
		for _, field := range fields {
			field = strings.Trim(field, " \t\r\n")
			v, err := strconv.ParseUint(field, 10, 0)
			if err != nil {
				goto parse_error
			}
			if v > uint64(math.MaxUint32) {
				goto parse_error
			}
			tsVals = append(tsVals, uint32(v))
		}
	case []uint32:
		tsVals = val
	default:
		goto parse_error
	}

	if len(tsVals) < f.minCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooFew")
	}
	if len(tsVals) > f.maxCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooMany")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, tsVals); err != nil {
			return nil, err
		}
	}

	return tsVals, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotTimestampSet")
}

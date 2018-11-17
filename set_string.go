package filter

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-apibox/types"
)

type StringSetFilter struct {
	strcase    int
	delimiter  string
	minCount   int
	maxCount   int
	validators []StringSetValidator
	allowVals  []string
}

type StringSetValidator func(paramName string, paramValue []string) *Error

// StringSet return a string set filter.
func StringSet() *StringSetFilter {
	f := new(StringSetFilter)
	f.strcase = STRING_RAWCASE
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *StringSetFilter) Allow(vals ...string) *StringSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// KeepCase do no case transform before validation.
func (f *StringSetFilter) KeepCase() *StringSetFilter {
	f.strcase = STRING_RAWCASE
	return f
}

// ToLower lower case string before validation.
func (f *StringSetFilter) ToLower() *StringSetFilter {
	f.strcase = STRING_LOWERCASE
	return f
}

// ToUpper lower case string before validation.
func (f *StringSetFilter) ToUpper() *StringSetFilter {
	f.strcase = STRING_UPPERCASE
	return f
}

// Delimiter set the delimiter of set string.
func (f *StringSetFilter) Delimiter(delimiter string) *StringSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *StringSetFilter) MinCount(count int) *StringSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *StringSetFilter) MaxCount(count int) *StringSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *StringSetFilter) AddValidator(validator StringSetValidator) *StringSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// Length valid whether string's length in set is equal with the specified length.
func (f *StringSetFilter) ItemLength(length int) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if len(v) < length {
				return NewError(ErrorInvalidParam, paramName, "ItemTooShort")
			}
			if len(v) > length {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLong")
			}
		}
		return nil
	})
	return f
}

// MinLen valid whether string in set is not longer than the specified length.
func (f *StringSetFilter) ItemMinLen(length int) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if len(v) < length {
				return NewError(ErrorInvalidParam, paramName, "ItemTooShort")
			}
		}
		return nil
	})
	return f
}

// MaxLen valid whether string in set is not shorter than the specified length.
func (f *StringSetFilter) ItemMaxLen(length int) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if len(v) > length {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLong")
			}
		}
		return nil
	})
	return f
}

// ShorterThan valid whether string in set is shorter than the specified length.
func (f *StringSetFilter) ItemShorterThan(length int) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if len(v) >= length {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLong")
			}
		}
		return nil
	})
	return f
}

// LongerThan valid whether string in set is longer than the specified length.
func (f *StringSetFilter) ItemLongerThan(length int) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if len(v) <= length {
				return NewError(ErrorInvalidParam, paramName, "ItemTooShort")
			}
		}
		return nil
	})
	return f
}

// Between valid whether length of string in set is in the range.
func (f *StringSetFilter) ItemBetween(minLength, maxLength int) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if len(v) < minLength {
				return NewError(ErrorInvalidParam, paramName, "ItemTooShort")
			}
			if len(v) > maxLength {
				return NewError(ErrorInvalidParam, paramName, "ItemTooLong")
			}
		}
		return nil
	})
	return f
}

// Match valid whether string in set is match the specified regular expression.
func (f *StringSetFilter) ItemMatch(pattern string) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		}
		for _, v := range paramValue {
			if !re.MatchString(v) {
				return NewError(ErrorInvalidParam, paramName, "ItemWrongFormat")
			}
		}
		return nil
	})
	return f
}

// IsNumeric valid whether string in set is numeric.
func (f *StringSetFilter) ItemIsNumeric() *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if _, err := strconv.ParseUint(v, 0, 64); err != nil {
				return NewError(ErrorInvalidParam, paramName, "ItemNotNumeric")
			}
		}
		return nil
	})
	return f
}

// IsDigit valid whether string in set is consist of digit numbers.
func (f *StringSetFilter) ItemIsDigit() *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			for _, c := range []byte(v) {
				if !(c >= '0' && c <= '9') {
					return NewError(ErrorInvalidParam, paramName, "ItemNotDigit")
				}
			}
		}
		return nil
	})
	return f
}

// IsAlpha valid whether string in set is consist of alpha letters.
func (f *StringSetFilter) ItemIsAlpha() *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			for _, c := range []byte(v) {
				if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
					return NewError(ErrorInvalidParam, paramName, "ItemNotAlpha")
				}
			}
		}
		return nil
	})
	return f
}

// IsAlphaNumeric valid whether string in set is consist of alpha letters or digit numbers.
func (f *StringSetFilter) ItemIsAlphaNumeric() *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			for _, c := range []byte(v) {
				if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
					return NewError(ErrorInvalidParam, paramName, "ItemNotAlphaNumeric")
				}
			}
		}
		return nil
	})
	return f
}

// In valid param value should in the specified set.
func (f *StringSetFilter) ItemIn(set []string) *StringSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
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
func (f *StringSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var strVals []string
	switch val := paramValue.(type) {
	case string:
		for _, allowVal := range f.allowVals {
			if allowVal == val {
				return val, nil
			}
		}
		if val != "" {
			strVals = strings.Split(val, f.delimiter)
		} else {
			strVals = []string{}
		}
	case []string:
		strVals = val
	default:
		goto parse_error
	}

	switch f.strcase {
	case STRING_LOWERCASE:
		for i, v := range strVals {
			strVals[i] = strings.ToLower(v)
		}
	case STRING_UPPERCASE:
		for i, v := range strVals {
			strVals[i] = strings.ToUpper(v)
		}
	}

	if len(strVals) < f.minCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooFew")
	}
	if len(strVals) > f.maxCount {
		return nil, NewError(ErrorInvalidParam, paramName, "TooMany")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, strVals); err != nil {
			return nil, err
		}
	}

	return strVals, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotStringSet")
}

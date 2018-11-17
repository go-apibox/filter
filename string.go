package filter

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	STRING_RAWCASE = iota
	STRING_LOWERCASE
	STRING_UPPERCASE
)

type StringFilter struct {
	strcase    int
	trim       bool
	validators []StringValidator
	allowVals  []string
}

type StringValidator func(paramName string, paramValue string) *Error

// String return a string filter.
func String() *StringFilter {
	f := new(StringFilter)
	f.strcase = STRING_RAWCASE
	return f
}

// Allow allow value is a string in the specified list
func (f *StringFilter) Allow(vals ...string) *StringFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// KeepCase do no case transform before validation.
func (f *StringFilter) KeepCase() *StringFilter {
	f.strcase = STRING_RAWCASE
	return f
}

// ToLower lower case string before validation.
func (f *StringFilter) ToLower() *StringFilter {
	f.strcase = STRING_LOWERCASE
	return f
}

// ToUpper lower case string before validation.
func (f *StringFilter) ToUpper() *StringFilter {
	f.strcase = STRING_UPPERCASE
	return f
}

// Trim trim empty char like space, \t, \n, \r before validation.
func (f *StringFilter) Trim() *StringFilter {
	f.trim = true
	return f
}

// AddValidator add a custom validator to filter
func (f *StringFilter) AddValidator(validator StringValidator) *StringFilter {
	f.validators = append(f.validators, validator)
	return f
}

// Length valid whether string's length is equal with the specified length.
func (f *StringFilter) Length(length int) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) < length {
			return NewError(ErrorInvalidParam, paramName, "TooShort")
		}
		if len(paramValue) > length {
			return NewError(ErrorInvalidParam, paramName, "TooLong")
		}
		return nil
	})
	return f
}

// MinLen valid whether string is not longer than the specified length.
func (f *StringFilter) MinLen(length int) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) < length {
			return NewError(ErrorInvalidParam, paramName, "TooShort")
		}
		return nil
	})
	return f
}

// MaxLen valid whether string is not shorter than the specified length.
func (f *StringFilter) MaxLen(length int) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) > length {
			return NewError(ErrorInvalidParam, paramName, "TooLong")
		}
		return nil
	})
	return f
}

// ShorterThan valid whether string is shorter than the specified length.
func (f *StringFilter) ShorterThan(length int) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) >= length {
			return NewError(ErrorInvalidParam, paramName, "TooLong")
		}
		return nil
	})
	return f
}

// LongerThan valid whether string is longer than the specified length.
func (f *StringFilter) LongerThan(length int) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) <= length {
			return NewError(ErrorInvalidParam, paramName, "TooShort")
		}
		return nil
	})
	return f
}

// Between valid whether string's length is in the range.
func (f *StringFilter) Between(minLength, maxLength int) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) < minLength {
			return NewError(ErrorInvalidParam, paramName, "TooShort")
		}
		if len(paramValue) > maxLength {
			return NewError(ErrorInvalidParam, paramName, "TooLong")
		}
		return nil
	})
	return f
}

// Match valid whether string is match the specified regular expression.
func (f *StringFilter) Match(pattern string) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if matched, err := regexp.MatchString(pattern, paramValue); err != nil {
			return NewError(ErrorInternalError, paramName, "InvalidValidator")
		} else if !matched {
			return NewError(ErrorInvalidParam, paramName, "WrongFormat")
		}
		return nil
	})
	return f
}

// IsNumeric valid whether string is numeric.
func (f *StringFilter) IsNumeric() *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if _, err := strconv.ParseUint(paramValue, 0, 64); err != nil {
			return NewError(ErrorInvalidParam, paramName, "NotNumeric")
		}
		return nil
	})
	return f
}

// IsDigit valid whether string is consist of digit numbers.
func (f *StringFilter) IsDigit() *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		for _, c := range []byte(paramValue) {
			if !(c >= '0' && c <= '9') {
				return NewError(ErrorInvalidParam, paramName, "NotDigit")
			}
		}
		return nil
	})
	return f
}

// IsAlpha valid whether string is consist of alpha letters.
func (f *StringFilter) IsAlpha() *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		for _, c := range []byte(paramValue) {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
				return NewError(ErrorInvalidParam, paramName, "NotAlpha")
			}
		}
		return nil
	})
	return f
}

// IsAlphaNumeric valid whether string is consist of alpha letters or digit numbers.
func (f *StringFilter) IsAlphaNumeric() *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		for _, c := range []byte(paramValue) {
			if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
				return NewError(ErrorInvalidParam, paramName, "NotAlphaNumeric")
			}
		}
		return nil
	})
	return f
}

// In valid param value should in the specified set.
func (f *StringFilter) In(set []string) *StringFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		for _, v := range set {
			if v == paramValue {
				return nil
			}
		}
		return NewError(ErrorInvalidParam, paramName, "NotInSet")
	})
	return f
}

// Run make the filter running.
func (f *StringFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	strVal, ok := paramValue.(string)
	if !ok {
		return nil, NewError(ErrorInvalidParam, paramName, "NotString")
	}
	for _, allowVal := range f.allowVals {
		if allowVal == strVal {
			return strVal, nil
		}
	}

	switch f.strcase {
	case STRING_LOWERCASE:
		strVal = strings.ToLower(strVal)
	case STRING_UPPERCASE:
		strVal = strings.ToUpper(strVal)
	}

	if f.trim {
		strVal = strings.Trim(strVal, " \t\r\n")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, strVal); err != nil {
			return nil, err
		}
	}

	return strVal, nil
}

package filter

import (
	"regexp"
	"strings"

	"github.com/go-apibox/types"
)

type EmailSetFilter struct {
	strcase    int
	delimiter  string
	minCount   int
	maxCount   int
	validators []EmailSetValidator
	allowVals  []string
}

type EmailSetValidator func(paramName string, paramValue []string) *Error

// EmailSet return a email set filter.
func EmailSet() *EmailSetFilter {
	f := new(EmailSetFilter)
	f.strcase = EMAIL_LOWERCASE
	f.delimiter = ","
	f.minCount = 0
	f.maxCount = types.MaxInt
	return f
}

// Allow allow value is a string in the specified list
func (f *EmailSetFilter) Allow(vals ...string) *EmailSetFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// KeepCase do no case transform before validation.
func (f *EmailSetFilter) KeepCase() *EmailSetFilter {
	f.strcase = STRING_RAWCASE
	return f
}

// ToLower lower case string before validation.
func (f *EmailSetFilter) ToLower() *EmailSetFilter {
	f.strcase = STRING_LOWERCASE
	return f
}

// ToUpper lower case string before validation.
func (f *EmailSetFilter) ToUpper() *EmailSetFilter {
	f.strcase = STRING_UPPERCASE
	return f
}

// Delimiter set the delimiter of set string.
func (f *EmailSetFilter) Delimiter(delimiter string) *EmailSetFilter {
	f.delimiter = delimiter
	return f
}

// MinCount set the max item count of set.
func (f *EmailSetFilter) MinCount(count int) *EmailSetFilter {
	f.minCount = count
	return f
}

// MaxCount set the max item count of set.
func (f *EmailSetFilter) MaxCount(count int) *EmailSetFilter {
	f.maxCount = count
	return f
}

// AddValidator add a custom validator to filter
func (f *EmailSetFilter) AddValidator(validator EmailSetValidator) *EmailSetFilter {
	f.validators = append(f.validators, validator)
	return f
}

// MinLen valid whether string in set not longer than the specified length.
func (f *EmailSetFilter) ItemMinLen(length int) *EmailSetFilter {
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

// MaxLen valid whether string in set not shorter than the specified length.
func (f *EmailSetFilter) ItemMaxLen(length int) *EmailSetFilter {
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

// ShorterThan valid whether string in set shorter than the specified length.
func (f *EmailSetFilter) ItemShorterThan(length int) *EmailSetFilter {
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

// LongerThan valid whether string in set longer than the specified length.
func (f *EmailSetFilter) ItemLongerThan(length int) *EmailSetFilter {
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
func (f *EmailSetFilter) ItemBetween(minLength, maxLength int) *EmailSetFilter {
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

// ItemDomain valid whether email is end with specified domain.
func (f *EmailSetFilter) ItemDomain(domain string) *EmailSetFilter {
	f.AddValidator(func(paramName string, paramValue []string) *Error {
		for _, v := range paramValue {
			if !strings.HasSuffix(v, "@"+domain) {
				return NewError(ErrorInvalidParam, paramName, "ItemWrongFormat")
			}
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *EmailSetFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var strVals []string
	switch val := paramValue.(type) {
	case string:
		val = strings.Trim(val, " \t\r\n")
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
		for i, strVal := range strVals {
			strVals[i] = strings.Trim(strVal, " \t\r\n")
		}
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

	// REF: https://fightingforalostcause.net/content/misc/2006/compare-email-regex.php
	// 以下是改进过的表达式，只支持正常域名为后缀，个性后缀最多支持10位
	for _, v := range strVals {
		pattern := `^(?i:[a-z0-9._-]+@(?:(?:(?:[a-z0-9]{1}[a-z0-9\-]{0,62}[a-z0-9]{1})|[a-z0-9])\.)+[a-z]{2,10})$`
		re := regexp.MustCompile(pattern)
		if !re.MatchString(v) {
			goto parse_error
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
	return nil, NewError(ErrorInvalidParam, paramName, "NotEmailSet")
}

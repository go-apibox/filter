package filter

import (
	"regexp"
	"strings"
)

const (
	EMAIL_RAWCASE = iota
	EMAIL_LOWERCASE
	EMAIL_UPPERCASE
)

type EmailFilter struct {
	strcase    int
	validators []EmailValidator
	allowVals  []string
}

type EmailValidator func(paramName string, paramValue string) *Error

// Email return a email filter.
func Email() *EmailFilter {
	f := new(EmailFilter)
	f.strcase = EMAIL_LOWERCASE
	return f
}

// Allow allow value is a string in the specified list
func (f *EmailFilter) Allow(vals ...string) *EmailFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// KeepCase do no case transform before validation.
func (f *EmailFilter) KeepCase() *EmailFilter {
	f.strcase = EMAIL_RAWCASE
	return f
}

// ToLower lower case string before validation.
func (f *EmailFilter) ToLower() *EmailFilter {
	f.strcase = EMAIL_LOWERCASE
	return f
}

// ToUpper lower case string before validation.
func (f *EmailFilter) ToUpper() *EmailFilter {
	f.strcase = EMAIL_LOWERCASE
	return f
}

// AddValidator add a custom validator to filter
func (f *EmailFilter) AddValidator(validator EmailValidator) *EmailFilter {
	f.validators = append(f.validators, validator)
	return f
}

// MinLen valid whether string is not longer than the specified length.
func (f *EmailFilter) MinLen(length int) *EmailFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) < length {
			return NewError(ErrorInvalidParam, paramName, "TooShort")
		}
		return nil
	})
	return f
}

// MaxLen valid whether string is not shorter than the specified length.
func (f *EmailFilter) MaxLen(length int) *EmailFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) > length {
			return NewError(ErrorInvalidParam, paramName, "TooLong")
		}
		return nil
	})
	return f
}

// ShorterThan valid whether string is shorter than the specified length.
func (f *EmailFilter) ShorterThan(length int) *EmailFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) >= length {
			return NewError(ErrorInvalidParam, paramName, "TooLong")
		}
		return nil
	})
	return f
}

// LongerThan valid whether string is longer than the specified length.
func (f *EmailFilter) LongerThan(length int) *EmailFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if len(paramValue) <= length {
			return NewError(ErrorInvalidParam, paramName, "TooShort")
		}
		return nil
	})
	return f
}

// Between valid whether string's length is in the range.
func (f *EmailFilter) Between(minLength, maxLength int) *EmailFilter {
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

// Domain valid whether email is end with specified domain.
func (f *EmailFilter) Domain(domain string) *EmailFilter {
	f.AddValidator(func(paramName string, paramValue string) *Error {
		if !strings.HasSuffix(paramValue, "@"+domain) {
			return NewError(ErrorInvalidParam, paramName, "WrongFormat")
		}
		return nil
	})
	return f
}

// Run make the filter running.
func (f *EmailFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	strVal, ok := paramValue.(string)
	if !ok {
		return nil, NewError(ErrorInvalidParam, paramName, "NotEmail")
	}
	strVal = strings.Trim(strVal, " \t\r\n")
	for _, allowVal := range f.allowVals {
		if allowVal == strVal {
			return strVal, nil
		}
	}

	switch f.strcase {
	case EMAIL_LOWERCASE:
		strVal = strings.ToLower(strVal)
	case EMAIL_UPPERCASE:
		strVal = strings.ToUpper(strVal)
	}

	// REF: https://fightingforalostcause.net/content/misc/2006/compare-email-regex.php
	// 以下是改进过的表达式，只支持正常域名为后缀，个性后缀最多支持10位
	pattern := `^(?i:[a-z0-9._-]+@(?:(?:(?:[a-z0-9]{1}[a-z0-9\-]{0,62}[a-z0-9]{1})|[a-z0-9])\.)+[a-z]{2,10})$`
	if matched, _ := regexp.MatchString(pattern, strVal); !matched {
		return nil, NewError(ErrorInvalidParam, paramName, "NotEmail")
	}

	for _, validator := range f.validators {
		if err := validator(paramName, strVal); err != nil {
			return nil, err
		}
	}

	return strVal, nil
}

package filter

import (
	"encoding/json"
	"reflect"
	"strings"
)

type JsonFilter struct {
	outVar     interface{}
	validators []JsonValidator
	allowVals  []string
	toString   bool
}

// Json return a json filter.
func Json() *JsonFilter {
	f := new(JsonFilter)
	return f
}

type JsonValidator func(paramName string, paramValue interface{}) *Error

// Output decode json value to specified variable.
func (f *JsonFilter) Output(outVar interface{}) *JsonFilter {
	f.outVar = outVar
	return f
}

// Allow allow value is a string in the specified list
func (f *JsonFilter) Allow(vals ...string) *JsonFilter {
	f.allowVals = append(f.allowVals, vals...)
	return f
}

// ToString return value as string
func (f *JsonFilter) ToString() *JsonFilter {
	f.toString = true
	return f
}

// AddValidator add a custom validator to filter
func (f *JsonFilter) AddValidator(validator JsonValidator) *JsonFilter {
	f.validators = append(f.validators, validator)
	return f
}

// Run make the filter running.
func (f *JsonFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, nil
	}

	var jsonVal interface{}

	strVal, ok := paramValue.(string)
	if !ok {
		goto parse_error
	}
	strVal = strings.Trim(strVal, " \t\r\n")
	for _, allowVal := range f.allowVals {
		if allowVal == strVal {
			return strVal, nil
		}
	}

	if f.outVar != nil {
		if err := json.Unmarshal([]byte(strVal), f.outVar); err != nil {
			goto parse_error
		}
		jsonVal = reflect.Indirect(reflect.ValueOf(f.outVar)).Interface()
	} else {
		if err := json.Unmarshal([]byte(strVal), &jsonVal); err != nil {
			goto parse_error
		}
	}

	for _, validator := range f.validators {
		if err := validator(paramName, jsonVal); err != nil {
			return nil, err
		}
	}

	if f.toString {
		return strVal, nil
	}
	return jsonVal, nil

parse_error:
	return nil, NewError(ErrorInvalidParam, paramName, "NotJson")
}

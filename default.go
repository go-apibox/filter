package filter

type DefaultFilter struct {
	defaultVal interface{}
}

// Default return a default filter.
func Default(defaultVal interface{}) *DefaultFilter {
	return &DefaultFilter{defaultVal}
}

// Run make the filter running.
func (f *DefaultFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return f.defaultVal, nil
	}
	return paramValue, nil
}

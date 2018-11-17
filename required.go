package filter

type RequiredFilter struct {
}

// Required return a require filter.
func Required() *RequiredFilter {
	return &RequiredFilter{}
}

// Run make the filter running.
func (f *RequiredFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == nil {
		return nil, NewError(ErrorMissingParam, paramName)
	}
	return paramValue, nil
}

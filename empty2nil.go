package filter

type EmptyToNilFilter struct {
}

// EmptyToNil return a empty to nil filter.
func EmptyToNil() *EmptyToNilFilter {
	return &EmptyToNilFilter{}
}

// Run make the filter running.
func (f *EmptyToNilFilter) Run(paramName string, paramValue interface{}) (interface{}, *Error) {
	if paramValue == "" {
		return nil, nil
	}
	return paramValue, nil
}

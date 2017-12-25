package errors

func NewError(source error, errorType ErrorType, meta map[string]interface{}) Error {
	return &defaultError{
		source:    source,
		errorType: errorType,
		meta:      meta,
	}
}

type defaultError struct {
	source    error
	errorType ErrorType
	meta      map[string]interface{}
}

func (e *defaultError) Error() string {
	return e.source.Error()
}

func (e *defaultError) GetSource() error {
	return e.source
}

func (e *defaultError) GetType() ErrorType {
	return e.errorType
}

func (e *defaultError) GetMeta() map[string]interface{} {
	return e.meta
}

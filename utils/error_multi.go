package utils

type MultiError struct {
	errors []error
}

func (e *MultiError) Append(err error) {
	e.errors = append(e.errors, err)
}

func (e MultiError) Error() string {
	var str string = "error[s] detected: \n"
	for _, err := range e.errors {
		str += (err.Error() + "\n")
	}
	return str
}

func (e MultiError) Has() bool {
	return len(e.errors) > 0
}

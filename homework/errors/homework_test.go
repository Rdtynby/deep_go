package main

import (
	"errors"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	result := strconv.Itoa(len(e.errors)) + " errors occured:\n"

	for _, err := range e.errors {
		result += "\t* " + err.Error()
	}

	result += "\n"

	return result
}

func Append(err error, errs ...error) *MultiError {
	switch e := any(err).(type) {
	case *MultiError:
		e.errors = append(e.errors, errs...)

		return e
	case nil:
		return &MultiError{errs}
	default:
		return &MultiError{append([]error{err}, errs...)}
	}
}

func (e *MultiError) Unwrap() error {
	if len(e.errors) == 1 {
		return nil
	}

	return &MultiError{e.errors[0 : len(e.errors)-1]}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	err2 := errors.Unwrap(err)
	expectedMessage2 := "1 errors occured:\n\t* error 1\n"
	assert.EqualError(t, err2, expectedMessage2)
	err3 := errors.Unwrap(err2)
	assert.Equal(t, err3, nil)

	err4 := Append(nil, &os.SyscallError{"error", errors.New("error 4")})
	expectedMessage4 := "1 errors occured:\n\t* error: error 4\n"
	assert.EqualError(t, err4, expectedMessage4)
}

package errors

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Kind is the error kind
type Kind string

// Error kinds
const (
	KUndefined  Kind = "Undefined"
)

// Error struct
type Error struct {
	Kind       Kind
	Message    string
	StatusCode int
	Err        error
}

func (e *Error) HTTPStatusCode() int {
	return e.StatusCode
}

func (e *Error) ErrorKind() string {
	return string(e.Kind)
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	fmtErr := fmt.Sprintf("%s: %s", e.Kind, e.Message)
	if e.Err == nil {
		return fmtErr
	}
	return fmt.Sprintf("%s :: caused by: %s", fmtErr, e.Err.Error())
}

func E(message string, args ...interface{}) error {
	e := &Error{
		Kind:    KUndefined,
		Message: message,
	}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			e.Kind = arg
		case int:
			e.StatusCode = arg
		case *Error:
			// Make a copy
			cp := *arg
			e.Err = &cp
		case error:
			e.Err = arg
		case Error:
			e.Err = arg.Err
			e.Kind = arg.Kind
		default:
			log.Panicf("unknown type %T, value %v in errors.E call", arg, arg)
		}
	}
	return e
}

func Ef(format string, args ...interface{}) error {
	numArgs := len(args)
	numPercentageChars := strings.Count(format, "%")
	numEscapedPercentageChars := strings.Count(format, "%%")
	numFormatArgs := numPercentageChars - (2 * numEscapedPercentageChars)
	if numArgs < numFormatArgs {
		log.Panicf("not enough arguments given for format '%s' in errors.Ef call. args: %v", format, args)
	}
	msg := fmt.Sprintf(format, args[0:numFormatArgs]...)
	return E(msg, args[numFormatArgs:numArgs]...)
}


// ErrKind returns the kind of the given error. If err is nil/empty, returns the an empty string. If the kind is
// KUndefined, nested errors are traversed recursively until the first non-undefined kind is found and returned.
//
// If err is a multierror returns the common kind of all wrapped errors or KUndefined if they are diverse.
//
// If no defined kind is found or the given error is not of type *Error / *multierror.Error, returns KUndefined.
func ErrKind(err error) Kind {
	if err == nil {
		return ""
	}
	switch err := err.(type) {
	case *Error:
		if err.Kind != KUndefined {
			return err.Kind
		} else if err.Err != nil {
			return ErrKind(err.Err)
		}
		return KUndefined
	case *multierror.Error:
		return multierrKind(err)
	default:
		return KUndefined
	}
}

// multierrKind returns the common kind of the given muliterror
func multierrKind(multierr *multierror.Error) Kind {
	if multierr == nil || multierr.Len() == 0 {
		return ""
	}
	kind := ErrKind(multierr.Errors[0])
	for _, err := range multierr.Errors {
		if ErrKind(err) != kind {
			return KUndefined // wrapped errors have different kinds, fallback to undefined
		}
	}
	return kind
}

// ErrStatusCode returns the status code of the given error. If the status code is undefined (0), nested errors are traversed
// recursively until the defined status code is found and returned. If no defined status code is found or the given
// error is not of type *errors.Error, returns 0.
func ErrStatusCode(err error) int {
	if err == nil {
		return 0
	} else if e, ok := err.(*Error); ok && e.StatusCode != 0 {
		return e.StatusCode
	} else if ok && e.Err != nil {
		return ErrStatusCode(e.Err)
	}
	return 0
}

// NewMultierr returns a new *multierror.Error with a custom error format function wrapping the given errors
func NewMultierr(errors ...error) *multierror.Error {
	multierr := &multierror.Error{
		ErrorFormat: PlainMultiErrFormatFunc,
	}
	for _, err := range errors {
		multierr = multierror.Append(multierr, err)
	}
	return multierr
}

// PlainMultiErrFormatFunc is a basic formatter that satisfied the ErrorFormatFunc type of multierror. It concatenates
// all errors to a one-line string.
func PlainMultiErrFormatFunc(es []error) string {
	if len(es) == 1 {
		return es[0].Error()
	}

	var errors []string
	for _, e := range es {
		errors = append(errors, e.Error())
	}
	joined := strings.Join(errors, ", ")

	return fmt.Sprintf("multi error: [%s]", joined)
}

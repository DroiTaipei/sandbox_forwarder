package droipkg

import (
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// DroiError - Droi Error interface
type DroiError interface {
	ErrorCode() int
	Error() string
}

// ConstDroiError - 7 digits error code following a space and error message
type ConstDroiError string

// ErrorCode - return error code
func (cde ConstDroiError) ErrorCode() int {
	i, _ := strconv.Atoi(string(cde)[0:7])
	return i
}

// Error - return predefined error message
func (cde ConstDroiError) Error() string {
	return string(cde)[8:]
}

// Wrap - [Not Implemented] Wrap error message
func (cde ConstDroiError) Wrap(message string) error {
	return Wrap(cde, message)
}

// CarrierDroiError - 7 digits error code following a space and wrappable error message
type CarrierDroiError string

// ErrorCode - return error code
func (cde CarrierDroiError) ErrorCode() int {
	i, _ := strconv.Atoi(string(cde)[0:7])
	return i
}

// SetErrorCode - setup the error code and return new CDE
func (cde CarrierDroiError) SetErrorCode(code int) (CarrierDroiError, error) {
	if code < 1000000 || code > 9999999 {
		return cde, fmt.Errorf("Error code should be 7 digits")
	}
	return CarrierDroiError(fmt.Sprintf("%d %s", code, cde.Error())), nil
}

// Error - return predefined error message
func (cde CarrierDroiError) Error() string {
	return string(cde)[8:]
}

// Wrap - append message to original one and return new CDE
func (cde CarrierDroiError) Wrap(message string) CarrierDroiError {
	return CarrierDroiError(fmt.Sprintf("%d %s", cde.ErrorCode(), cde.Error()+message))
}

// NewCarrierDroiError returns CarrierDroiError by given error code and message
func NewCarrierDroiError(code int, message string) CarrierDroiError {
	return CarrierDroiError(fmt.Sprintf("%d %s", code, message))
}

// AsDroiError - The Error "As" Droi Error but more enhanced
type AsDroiError interface {
	DroiError
	CastOff() DroiError
	AsEqual(DroiError) bool
	Trace(string)
	Stack() error
	Format(s fmt.State, verb rune)
	Cause() error
}

// TraceDroiError -  Compatiable error stack trace & DroiError
// Contain a DroiError for "as" DroiError
// Using "error" & Wrap as a stack for tracing error history
type TraceDroiError struct {
	DroiError
	stack error
}

// NewTraceDroiError - Get a new *TraceDroiError from DroiError & external error
func NewTraceDroiError(der DroiError, err error) *TraceDroiError {
	return &TraceDroiError{
		DroiError: der,
		stack:     err,
	}
}

// NewTraceWithMsg  - One sugar function for NewTraceDroiError
func NewTraceWithMsg(der DroiError, msg string) *TraceDroiError {
	return &TraceDroiError{
		DroiError: der,
		stack:     NewError(msg),
	}
}

// Error - return error message
func (fde *TraceDroiError) Error() string {
	msg := fde.DroiError.Error()
	if fde.stack != nil {
		return fde.stack.Error() + ": " + msg
	}
	return msg
}

// CastOff - Unwrap TraceDroiError to return DroiError
func (fde *TraceDroiError) CastOff() DroiError {
	return fde.DroiError
}

// AsEqual - Return the DroiError is equal or not
func (fde *TraceDroiError) AsEqual(target DroiError) bool {
	return fde.DroiError == target
}

// Trace - Add message for trace
func (fde *TraceDroiError) Trace(message string) {
	if fde.stack == nil {
		fde.stack = NewError(message)
	} else {
		fde.stack = Wrap(fde.stack, message)
	}
}

// Stack - Return the stack
func (fde *TraceDroiError) Stack() error {
	return fde.stack
}

// Format - Implement the fmt.Formatter interface
// type Formatter interface {
//         Format(f State, c rune)
// }
// For fmt.Printf("%+v") to print stack trace
func (fde *TraceDroiError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n%+v\n", fde.stack, fde.DroiError)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, fde.Error())
	default:
		io.WriteString(s, fde.Error())
	}
}

// Cause - Implement Cause, Support to found the origin root Cause
func (fde *TraceDroiError) Cause() error {
	type causer interface {
		Cause() error
	}
	err := fde.stack
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func NewError(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}

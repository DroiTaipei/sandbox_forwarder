package droipkg

import (
	"github.com/pkg/errors"
	"strconv"
)

type DroiError interface {
	ErrorCode() int
	Error() string
	Wrap(string)
	SetErrorCode(int)
}

type ConstDroiError string

func (cde ConstDroiError) ErrorCode() int {
	i, _ := strconv.Atoi(string(cde)[0:7])
	return i
}
func (cde ConstDroiError) SetErrorCode(code int) {
	// Not Implemented
	return
}
func (cde ConstDroiError) Error() string {
	return string(cde)[8:]
}
func (cde ConstDroiError) Wrap(msg string) {
	// Not Implemented
	return
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

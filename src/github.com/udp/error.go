package udp

import (
	"fmt"

	"github.com/DroiTaipei/droipkg"
)

const (
	// ErrorCodeUnknown is base error code of udp service
	ErrorCodeUnknown = 1190000
	// ErrorCodeTimeout is base error code for timeout case
	ErrorCodeTimeout = ErrorCodeUnknown + 1
	// ErrorCodePanic is the error code for process panic
	ErrorCodePanic = ErrorCodeUnknown + 2
	// ErrorLoadDBN is the error code for loading DBN fail
	ErrorLoadDBN = ErrorCodeUnknown + 3

	// ErrorCodeInvalidInput is the error code for bad request
	ErrorCodeInvalidInput = ErrorCodeUnknown + 101

	// ErrorCodeSendFull is the error code for send channel full
	ErrorCodeSendFull = ErrorCodeUnknown + 201
)

// ErrTimeout is the http timeout error for HTTP PUSH API
var ErrTimeout = droipkg.CarrierDroiError(fmt.Sprintf("%d push httpapi timeout", ErrorCodeTimeout))

// ErrInvalidInput is the error for invalid input
var ErrInvalidInput = droipkg.CarrierDroiError(fmt.Sprintf("%d invalid input", ErrorCodeInvalidInput))

// ErrPanic is the error for process panic
var ErrPanic = droipkg.CarrierDroiError(fmt.Sprintf("%d udp service panic", ErrorCodePanic))

// ErrFull is the error for send channel full
var ErrFull = droipkg.CarrierDroiError(fmt.Sprintf("%d send push notification too fast", ErrorCodeSendFull))

// ErrDBN is the error for loading DBN fail
var ErrDBN = droipkg.CarrierDroiError(fmt.Sprintf("%d load DBN from mongodb fail", ErrorLoadDBN))

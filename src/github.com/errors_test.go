package droipkg

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockErrCode = 1234567
	mockErrMsg  = "Mock error message. "
)

var (
	mockConstErr   = ConstDroiError(fmt.Sprintf("%d %s", mockErrCode, mockErrMsg))
	mockCarrierErr = CarrierDroiError(fmt.Sprintf("%d %s", mockErrCode, mockErrMsg))
	mockTraceErr   = NewTraceDroiError(mockConstErr, nil)
)

func BeforeTest() {

}

func Test_CarrierErrorCode(t *testing.T) {
	assert.Equal(t, mockErrCode, mockCarrierErr.ErrorCode())
}

func Test_CarrierSetErrorCode(t *testing.T) {
	newCode := 1111111
	newError, err := mockCarrierErr.SetErrorCode(newCode)
	assert.Nil(t, err)
	assert.Equal(t, newCode, newError.ErrorCode())

	// While updating invalid one, it should return error and original code is not updated
	newError, err = mockCarrierErr.SetErrorCode(1234)
	assert.NotNil(t, err)
	assert.Equal(t, mockErrCode, newError.ErrorCode())
}

func Test_CarrierError(t *testing.T) {
	assert.Equal(t, mockCarrierErr.Error(), mockErrMsg)
}

func Test_CarrierWrap(t *testing.T) {
	newMessage := "wrap new message"
	newError := mockCarrierErr.Wrap(newMessage)
	assert.Equal(t, mockCarrierErr.Error()+newMessage, newError.Error())
}

func TestConstErrorCode(t *testing.T) {
	assert.Equal(t, mockErrCode, mockConstErr.ErrorCode())
}

func TestConstError(t *testing.T) {
	assert.Equal(t, mockErrMsg, mockConstErr.Error())
}

func TestConstWrap(t *testing.T) {
	newMessage := "wrap new message"
	newError := mockConstErr.Wrap(newMessage)
	assert.Equal(t, fmt.Sprintf("%s: %s", newMessage, mockConstErr.Error()), newError.Error())
}

func TestTraceErrorCode(t *testing.T) {
	assert.Equal(t, mockErrCode, mockTraceErr.ErrorCode())
}

func TestTraceErr(t *testing.T) {
	assert.Equal(t, mockErrMsg, mockTraceErr.Error())
}

func TestTraceWrap(t *testing.T) {
	initMsg := mockTraceErr.Error()
	newMsg := "wrap new message"
	mockTraceErr.Trace(newMsg)
	assert.Equal(t, fmt.Sprintf("%s: %s", newMsg, initMsg), mockTraceErr.Error())
	anoError := NewTraceDroiError(mockCarrierErr, mockTraceErr)
	assert.Equal(t, fmt.Sprintf("%s: %s", mockTraceErr.Error(), mockCarrierErr.Error()), anoError.Error())
}

func TestAsDroiError(t *testing.T) {
	var p AsDroiError
	p = NewTraceDroiError(mockConstErr, nil)
	assert.Equal(t, mockErrCode, p.ErrorCode())
	assert.Equal(t, mockErrMsg, p.Error())
	assert.Equal(t, true, p.AsEqual(mockConstErr))
	assert.Equal(t, mockConstErr, p.CastOff())
}

func TestTraceDroiError(t *testing.T) {
	var p *TraceDroiError
	p = NewTraceDroiError(mockConstErr, nil)
	assert.Equal(t, mockErrCode, p.ErrorCode())
	assert.Equal(t, mockErrMsg, p.Error())
	assert.Equal(t, true, p.AsEqual(mockConstErr))
	assert.Equal(t, mockConstErr, p.CastOff())
	assert.Equal(t, mockConstErr, p.DroiError)
}

func TestTraceAsDroiError(t *testing.T) {
	var p DroiError
	p = NewTraceDroiError(mockConstErr, nil)
	assert.Equal(t, mockErrCode, p.ErrorCode())
	assert.Equal(t, mockErrMsg, p.Error())
}

func TestTraceCause(t *testing.T) {
	msgRoot := "Test Error Cause"
	msgTrace := "Another Error Message"
	c := NewError(msgRoot)
	p := NewTraceDroiError(mockConstErr, c)
	p.Trace(msgTrace)
	assert.Equal(t, c, p.Cause())
	assert.Equal(t, c, Cause(p))
	assert.Equal(t, fmt.Sprintf("%s: %s: %s", msgTrace, msgRoot, mockErrMsg), p.Error())
}

func TestTraceFormat(t *testing.T) {
	sampleMsg := "error sample"
	tests := []struct {
		err    *TraceDroiError
		format string
		want   string
	}{{
		NewTraceDroiError(mockConstErr, nil),
		"%s",
		mockErrMsg,
	}, {
		NewTraceDroiError(mockConstErr, nil),
		"%v",
		mockErrMsg,
	}, {
		NewTraceDroiError(mockConstErr, nil),
		"%+v",
		fmt.Sprintf("<nil>\n%s\n", mockErrMsg),
	}, {
		NewTraceDroiError(mockConstErr, NewError(sampleMsg)),
		"%s",
		fmt.Sprintf("%s: %s", sampleMsg, mockErrMsg),
	}, {
		NewTraceDroiError(mockConstErr, NewError(sampleMsg)),
		"%v",
		fmt.Sprintf("%s: %s", sampleMsg, mockErrMsg),
	}, {
		NewTraceDroiError(mockConstErr, NewError(sampleMsg)),
		"%+v",
		sampleMsg + "\n" +
			".+\n\t(.+?)\n(.*?)" +
			"github.com/DroiTaipei/droipkg.TestTraceFormat\n" +
			"(.+?)/github.com/DroiTaipei/droipkg/errors_test.go:.*",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.err, tt.format, tt.want)
	}
}

// Copy From https://github.com/pkg/errors/blob/master/format_test.go#L120
func TestFormatWrapf(t *testing.T) {
	tests := []struct {
		error
		format string
		want   string
	}{{
		Wrapf(io.EOF, "error%d", 2),
		"%s",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%v",
		"error2: EOF",
	}, {
		Wrapf(io.EOF, "error%d", 2),
		"%+v",
		"EOF\n" +
			"error2\n" +
			"(.+?)\n\t(.+?)\n(.*?)" +
			"github.com/DroiTaipei/droipkg.TestFormatWrapf\n" +
			"\t(.+?)/github.com/DroiTaipei/droipkg/errors_test.go:.+",
	}, {
		Wrapf(NewError("error"), "error%d", 2),
		"%s",
		"error2: error",
	}, {
		Wrapf(NewError("error"), "error%d", 2),
		"%v",
		"error2: error",
	}, {
		Wrapf(NewError("error"), "error%d", 2),
		"%+v",
		"error\n" +
			"(.+?)\n\t(.+?)\n(.*?)" +
			"github.com/DroiTaipei/droipkg.TestFormatWrapf\n" +
			"\t(.+?)/github.com/DroiTaipei/droipkg/errors_test.go:.+",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.error, tt.format, tt.want)
	}
}

func TestTraceStack(t *testing.T) {
	msgRoot := "Test Error Cause"
	msgTrace := "Another Error Message"
	c := NewError(msgRoot)
	p := NewTraceDroiError(mockConstErr, c)
	p.Trace(msgTrace)
	s := p.Stack()
	assert.Equal(t, p.Cause(), Cause(s))
	assert.Equal(t, Cause(p), Cause(s))
	assert.Equal(t, fmt.Sprintf("%s: %s", msgTrace, msgRoot), s.Error())

}

func TestNewCarrierDroiError(t *testing.T) {
	p := NewCarrierDroiError(mockErrCode, mockErrMsg)
	assert.Equal(t, mockErrCode, p.ErrorCode())
	assert.Equal(t, mockErrMsg, p.Error())
}

// Copy from https://github.com/pkg/errors/blob/master/format_test.go#L363
func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}

// Do somethings after all test cases
func AfterTest() {

}

func TestMain(m *testing.M) {
	BeforeTest()
	retCode := m.Run()
	AfterTest()
	os.Exit(retCode)
}

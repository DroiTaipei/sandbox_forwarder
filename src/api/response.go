package api

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

const (
	CODE_OK   = 0
	HTTP_CODE = 200
	STR_OK    = "OK"
)

type Response struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message,omitempty"`
	Result  interface{} `json:"Result,omitempty"`
	Count   *int        `json:"Count,omitempty"`
}

func NewResponse() *Response {
	return &Response{Code: CODE_OK, Message: STR_OK}
}

func NewErrorwResponse(err error) (resp *Response) {
	resp = NewResponse()
	resp.Code = getDroiErrorCode(err)
	errorLog(err)
	resp.Message = err.Error()
	return
}

func (r *Response) Write(c *fasthttp.RequestCtx) {
	buf, _ := json.Marshal(r)
	c.SetContentType("application/json")
	c.Write(buf)
}

// WriteError parse error message and gen corresponding error code
func WriteError(c *fasthttp.RequestCtx, err error) {
	resp := NewResponse()
	resp.Message = err.Error()
	resp.Code = getDroiErrorCode(err)
	errorLog(err)
	resp.Write(c)
}

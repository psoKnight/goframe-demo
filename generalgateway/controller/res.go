package controller

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	Success       = 0
	ParamError    = 1300
	InternalError = 1301
)

var errorText = map[int]string{
	Success:       "成功",
	ParamError:    "参数错误",
	InternalError: "内部错误",
}

// HelloData is the successful greeting payload.
type HelloData struct {
	Message string `json:"message" dc:"greeting message"`
}

// CalculateData is the successful calculator payload.
type CalculateData struct {
	Result int64 `json:"result" dc:"calculation result"`
}

// CodeResponse keeps the API response compatible with grpc-demo.
type CodeResponse struct {
	ErrorNo  int    `json:"error_no" dc:"business error code"`
	ErrorMsg string `json:"error_msg" dc:"business error message"`
	Data     any    `json:"data,omitempty" dc:"business data"`
}

func response(code int, message string, data any) *CodeResponse {
	if message == "" {
		message = errorText[code]
	}
	return &CodeResponse{ErrorNo: code, ErrorMsg: message, Data: data}
}

// Response serializes standard GoFrame controller results without adding
// another response envelope, keeping the API contract compatible with grpc-demo.
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.BufferLength() > 0 || r.Response.BytesWritten() > 0 {
		return
	}
	if err := r.GetError(); err != nil {
		code := InternalError
		if gerror.Code(err) == gcode.CodeInvalidParameter {
			code = ParamError
		}
		r.Response.WriteJson(response(code, err.Error(), nil))
		return
	}
	if result := r.GetHandlerResponse(); result != nil {
		r.Response.WriteJson(result)
	}
}

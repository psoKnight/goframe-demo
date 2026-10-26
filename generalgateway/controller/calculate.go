package controller

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/codes"

	"goframe-demo/generalgateway/service"
)

// CalculateReq is the POST /v1/calculate request.
type CalculateReq struct {
	g.Meta    `path:"/calculate" method:"post" tags:"Calculator" summary:"调用 Calculator 微服务"`
	Operation string          `json:"operation" dc:"add or multiply"`
	Left      json.RawMessage `json:"left" dc:"int64 value or string"`
	Right     json.RawMessage `json:"right" dc:"int64 value or string"`
}

// Calculate handles POST /v1/calculate.
func (c *Controller) Calculate(ctx context.Context, req *CalculateReq) (*CodeResponse, error) {
	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	if operation == "" {
		return response(ParamError, "operation 不能为空", nil), nil
	}
	if operation != service.OperationAdd && operation != service.OperationMultiply {
		return response(ParamError, "operation 仅支持 add 或 multiply", nil), nil
	}

	left, err := parseInt64Field(req.Left)
	if err != nil {
		return response(ParamError, "left 必须是 int64 整数", nil), nil
	}
	right, err := parseInt64Field(req.Right)
	if err != nil {
		return response(ParamError, "right 必须是 int64 整数", nil), nil
	}

	reply, err := service.Calculate(ctx, operation, left, right)
	if err != nil {
		code := gerror.Code(err).Code()
		if code == int(codes.InvalidArgument) || code == int(codes.OutOfRange) {
			return response(ParamError, gerror.Current(err).Error(), nil), nil
		}
		return response(InternalError, err.Error(), nil), nil
	}

	return response(Success, "", &CalculateData{Result: reply.GetResult()}), nil
}

func parseInt64Field(raw json.RawMessage) (int64, error) {
	value := strings.TrimSpace(string(raw))
	if value == "" {
		return 0, strconv.ErrSyntax
	}
	if value[0] == '"' {
		if err := json.Unmarshal(raw, &value); err != nil {
			return 0, err
		}
	}
	return strconv.ParseInt(strings.TrimSpace(value), 10, 64)
}

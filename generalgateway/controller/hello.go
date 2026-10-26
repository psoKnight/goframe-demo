package controller

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc/codes"

	"goframe-demo/generalgateway/service"
)

// Controller exposes the HTTP gateway APIs.
type Controller struct{}

// New creates an HTTP controller.
func New() *Controller {
	return &Controller{}
}

// HelloReq is the POST /v1/hello request.
type HelloReq struct {
	g.Meta `path:"/hello" method:"post" tags:"Greeter" summary:"调用 Greeter 微服务"`
	Name   string `json:"name" dc:"greeted name"`
}

// Hello handles POST /v1/hello.
func (c *Controller) Hello(ctx context.Context, req *HelloReq) (*CodeResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return response(ParamError, "name 不能为空", nil), nil
	}

	reply, err := service.SayHello(ctx, name)
	if err != nil {
		if gerror.Code(err).Code() == int(codes.InvalidArgument) {
			return response(ParamError, gerror.Current(err).Error(), nil), nil
		}
		return response(InternalError, err.Error(), nil), nil
	}

	return response(Success, "", &HelloData{Message: reply.GetMessage()}), nil
}

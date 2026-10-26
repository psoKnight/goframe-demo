package healthing

import "github.com/gogf/gf/v2/net/ghttp"

// Status is the shared HTTP readiness endpoint for GoFrame services.
func Status(r *ghttp.Request) {
	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.Response.Write("ok\n")
}

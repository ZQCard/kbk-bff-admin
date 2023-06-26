package requestInfo

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2/middleware"
)

const DomainKey = "x-md-global-domain"
const ipKey = "x-md-global-ip"

// setRequestInfo 设置Request信息
func SetRequestInfo() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {

			if tr, ok := transport.FromServerContext(ctx); ok {
				// 将请求信息放入ctx中
				if ht, ok := tr.(*http.Transport); ok {
					ctx = context.WithValue(ctx, ipKey, ht.Request().RemoteAddr)
				}
				// 获取请求域
				domain := tr.RequestHeader().Get(DomainKey)
				if domain == "" {
					domain = "default"
				}
				ctx = metadata.AppendToClientContext(ctx, DomainKey, domain)
				ctx = context.WithValue(ctx, DomainKey, domain)
				tr.ReplyHeader().Set(DomainKey, domain)

			}
			return handler(ctx, req)
		}
	}
}

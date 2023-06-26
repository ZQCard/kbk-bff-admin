package auth

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v4"

	authorizationV1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
)

const CabinObj = "role"

func CasbinMiddleware(authorizationClient authorizationV1.AuthorizationServiceClient) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			claim, _ := jwt.FromContext(ctx)
			if claim == nil {
				return nil, errors.Unauthorized("Token Missing", "token解析失败")
			}
			claimInfo := claim.(jwt2.MapClaims)
			if claimInfo[CabinObj] == nil {
				return nil, errors.Unauthorized("Unauthorized", "权限不足")
			}

			role := claimInfo[CabinObj].(string)
			// 获取当前服务operation，验证策略为 operation + method + role  暂时关闭
			if tr, ok := transport.FromServerContext(ctx); ok {
				// 发起grpc请求，确认权限
				// 获取请求方法
				act := ""
				if ht, ok := tr.(*http.Transport); ok {
					act = ht.Request().Method
				}
				// 获取请求的PATH
				obj := tr.Operation()

				reply, err := authorizationClient.CheckAuthorization(ctx, &authorizationV1.CheckAuthorizationReq{
					Sub: role,
					Obj: obj,
					Act: act,
				})
				if err != nil {
					return nil, v1.ErrorSystemError("权限服务异常").WithCause(err)
				}
				if reply.Success != true {
					return nil, v1.ErrorForbidden("权限不足")
				}
			}
			return handler(ctx, req)
		}
	}
}

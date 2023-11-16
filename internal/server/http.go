package server

import (
	"context"

	"github.com/go-redis/redis"
	jwt2 "github.com/golang-jwt/jwt/v4"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	adminV1 "github.com/ZQCard/kbk-bff-admin/api/bff-admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	"github.com/ZQCard/kbk-bff-admin/internal/service"
	"github.com/ZQCard/kbk-bff-admin/pkg/middleware/apiLog"
	"github.com/ZQCard/kbk-bff-admin/pkg/middleware/auth"
	"github.com/ZQCard/kbk-bff-admin/pkg/middleware/requestInfo"
	"github.com/ZQCard/kbk-bff-admin/pkg/middleware/userInfo"

	authorizationV1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	logV1 "github.com/ZQCard/kbk-log/api/log/v1"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/admin.v1.Admin/Login"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	cfg *conf.Bootstrap,
	service *service.AdminInterface,
	authorizationClient authorizationV1.AuthorizationServiceClient,
	logger log.Logger,
	logClient logV1.LogServiceClient,
	redisCli *redis.Client,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			tracing.Server(),
			// 访问日志
			logging.Server(logger),
			// 设置domain
			requestInfo.SetRequestInfo(),
			// 对于需要登录的路由进行jwt中间件验证
			selector.Server(
				// 解析jwt
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte(cfg.Jwt.Key), nil
				},
					jwt.WithSigningMethod(jwt2.SigningMethodHS256),
					jwt.WithClaims(func() jwt2.Claims {
						return jwt2.MapClaims{}
					})),
				userInfo.SetUserInfo(),
				auth.CasbinMiddleware(authorizationClient),

				apiLog.Server(logger, logClient, authorizationClient, redisCli),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openAPIhandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIhandler)
	adminV1.RegisterAdminHTTPServer(srv, service)
	return srv
}

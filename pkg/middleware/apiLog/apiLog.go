package apiLog

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	stdHttp "net/http"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis"
	"google.golang.org/protobuf/types/known/emptypb"

	authorizationV1 "github.com/ZQCard/kratos-base-kit/kbk-authorization/api/authorization/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-bff-admin/pkg/utils/redisHelper"
	logV1 "github.com/ZQCard/kratos-base-kit/kbk-log/api/log/v1"
)

// Redacter defines how to log an object
type Redacter interface {
	Redact() string
}

// Server is an server logging middleware.
func Server(logger log.Logger, logClient logV1.LogServiceClient, authClient authorizationV1.AuthorizationServiceClient, redisCli *redis.Client) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			info, ok := transport.FromServerContext(ctx)
			if !ok {
				err := fmt.Errorf("transport From context failed, ctx not kratos server context, please check frameWork")
				return "", errors.InternalServer("SYSTEM ERROR ", err.Error())
			}
			var (
				traceId   string
				domain    string
				component string
				userId    string
				username  string
				role      string
				apiName   string
				method    string
				path      string
				request   string
				code      int32
				reason    string
				ip        string
				operation string
				latency   string
			)
			operation = info.Operation()
			userIdCtx := ctx.Value("x-md-global-userId")
			if userIdCtx != nil {
				// 从上游ctx中获取userId
				userId = strconv.FormatInt(userIdCtx.(int64), 10)
			}

			roleCtx := ctx.Value("x-md-global-role")
			if roleCtx != nil {
				// 从上游ctx中获取userId
				role = roleCtx.(string)
			}
			domainCtx := ctx.Value("x-md-global-domain")
			if domainCtx != nil {
				// 从上游ctx中获取domain
				domain = domainCtx.(string)
			}
			usernameCtx := ctx.Value("x-md-global-username")
			if usernameCtx != nil {
				// 从上游ctx中获取domain
				username = usernameCtx.(string)
			}
			ipCtx := ctx.Value("x-md-global-ip")
			if usernameCtx != nil {
				// 从上游ctx中获取domain
				ip = ipCtx.(string)
			}
			switch tpKind := info.Kind(); tpKind {
			case transport.KindHTTP:
				tpHttp := info.(*http.Transport)
				method = tpHttp.Request().Method
				// 如果是get请求，则不使用日志
				if method == stdHttp.MethodGet {
					return handler(ctx, req)
				}
				path = tpHttp.Request().URL.Path
				component = tpHttp.Operation()
				path = tpHttp.PathTemplate()
				args, _ := ioutil.ReadAll(tpHttp.Request().Body)
				request = string(args)
			default:
				err := fmt.Errorf("transport type err : %v", tpKind)
				return "", errors.InternalServer("SYSTEM ERROR ", err.Error())
			}

			// 单词读取api管理中的api名称 存放于redis中, 如果不存在则多读取一次，如果还不存在，则进行报错 api不存在
			apiName, err = getApiName(ctx, domain, method, operation, authClient, redisCli)
			if err != nil {
				return "", err
			}

			startTime := time.Now()
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}

			latencyTime := time.Since(startTime).Seconds()
			latency = strconv.FormatFloat(latencyTime, 'f', -1, 64) + "s"
			logClient.CreateLog(ctx, &logV1.CreateLogReq{
				UserId:    userId,
				Username:  username,
				Role:      role,
				Operation: operation,
				Name:      apiName,
				Code:      strconv.Itoa(int(code)),
				Component: component,
				Reason:    reason,
				Path:      path,
				Ip:        ip,
				TraceId:   traceId,
				Method:    method,
				Request:   request,
				Latency:   latency,
			})
			return
		}
	}
}

func getApiName(ctx context.Context, domain string, method string, operation string, authClient authorizationV1.AuthorizationServiceClient, redisCli *redis.Client) (string, error) {
	// 先从redis获取数据,如果没有数据则请求权限服务
	apiNameMap := make(map[string]string)
	key := domain + "_" + "api_all"
	apiNameMapStr := redisHelper.GetRedisCache(redisCli, key)
	if apiNameMapStr == "" {
		apiNameMap, err := getApiAll(ctx, authClient)
		if err != nil {
			return "", errors.InternalServer("SYSTEM ERROR ", err.Error())
		}
		apiNameMapStrBytes, _ := json.Marshal(apiNameMap)
		json.Unmarshal(apiNameMapStrBytes, &apiNameMap)
		apiNameMapStr = string(apiNameMapStrBytes)
		redisHelper.SaveRedisCache(redisCli, key, apiNameMapStr, 120*time.Second)
	}
	json.Unmarshal([]byte(apiNameMapStr), &apiNameMap)
	apiName := apiNameMap[domain+"-"+method+"-"+operation]
	if apiName == "" {
		// 不存在且是redis数据,刷新
		if apiNameMapStr != "" {
			apiNameMap, err := getApiAll(ctx, authClient)
			if err != nil {
				return "", errors.InternalServer("SYSTEM ERROR ", err.Error())
			}
			apiNameMapStrBytes, _ := json.Marshal(apiNameMap)
			json.Unmarshal(apiNameMapStrBytes, &apiNameMap)
			redisHelper.SaveRedisCache(redisCli, key, string(apiNameMapStrBytes), 120*time.Second)
			apiName = apiNameMap[domain+"-"+method+"-"+operation]
			if apiName == "" {
				return "", errors.BadRequest("BAD REQUEST", domain+"-"+method+"-"+operation+" api not found")
			}
		} else {
			return "", errors.BadRequest("BAD REQUEST", domain+"-"+method+"-"+operation+" api not found")
		}
	}
	return apiName, nil
}

func getApiAll(ctx context.Context, authClient authorizationV1.AuthorizationServiceClient) (map[string]string, error) {
	apiNameMap := make(map[string]string)
	apiAll, err := authClient.GetApiListAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, errors.InternalServer("SYSTEM ERROR ", err.Error())
	}
	for _, v := range apiAll.List {
		apiNameMap[v.Domain+"-"+v.Method+"-"+v.Path] = v.Group + "-" + v.Name
	}
	return apiNameMap, nil
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if redacter, ok := req.(Redacter); ok {
		return redacter.Redact()
	}
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

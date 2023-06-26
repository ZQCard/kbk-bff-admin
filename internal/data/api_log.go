package data

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	logV1 "github.com/ZQCard/kbk-log/api/log/v1"
	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewApiLogClient(sr *conf.Endpoint, r registry.Discovery, tp *tracesdk.TracerProvider) logV1.LogServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Log),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(tracing.WithTracerProvider(tp)),
			// 元信息
			metadata.Client(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := logV1.NewLogServiceClient(conn)
	return c
}

type ApiLogRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewApiLogRepo(data *Data, logger log.Logger) *ApiLogRepo {
	return &ApiLogRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/apiLog")),
		sg:   &singleflight.Group{},
	}
}

func (rp ApiLogRepo) ListApiLog(ctx context.Context, req *v1.GetApiLogListReq) (*v1.GetApiLogListRes, error) {
	reply, err := rp.data.logClient.GetLogList(ctx, &logV1.GetLogListReq{
		Page:      req.Page,
		PageSize:  req.PageSize,
		TraceId:   req.TraceId,
		UserId:    req.UserId,
		Username:  req.Username,
		Role:      req.Role,
		Operation: req.Operation,
		Ip:        req.Ip,
	})
	if err != nil {
		return nil, err
	}
	list := []*v1.ApiLog{}

	for _, v := range reply.List {
		list = append(list, apiLogServiceToApi(v))
	}

	response := &v1.GetApiLogListRes{}
	response.Total = reply.Total
	response.List = list
	return response, nil
}

func apiLogServiceToApi(info *logV1.Log) *v1.ApiLog {
	return &v1.ApiLog{
		Id:        info.Id,
		Name:      info.Name,
		TraceId:   info.TraceId,
		Component: info.Component,
		UserId:    info.UserId,
		Username:  info.Username,
		Role:      info.Role,
		Method:    info.Method,
		Path:      info.Path,
		Request:   info.Request,
		Code:      info.Code,
		Ip:        info.Ip,
		Latency:   info.Latency,
		CreatedAt: info.CreatedAt,
		Operation: info.Operation,
	}
}

package data

import (
	"context"

	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/jinzhu/copier"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/singleflight"

	aiV1 "github.com/ZQCard/kbk-ai/api/ai/v1"
	v1 "github.com/ZQCard/kbk-bff-admin/api/admin/v1"
)

func NewAIClient(sr *conf.Endpoint, r registry.Discovery, tp *tracesdk.TracerProvider) aiV1.AIServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Ai),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(tracing.WithTracerProvider(tp)),
			// 元信息
			metadata.Client(),
		),
		grpc.WithTimeout(sr.TimeOut.AsDuration()),
	)
	if err != nil {
		panic(err)
	}
	c := aiV1.NewAIServiceClient(conn)
	return c
}

type AIRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewAIRepo(data *Data, logger log.Logger) *AIRepo {
	return &AIRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/apiLog")),
		sg:   &singleflight.Group{},
	}
}

func (r AIRepo) CreateGptMessage(ctx context.Context, req *v1.CreateGptMessageReq) (*v1.CreateGptMessageRes, error) {
	data := &aiV1.CreateGptMessageReq{}
	copier.Copy(data, req)
	reply, err := r.data.aiClient.CreateGptMessage(ctx, data)
	if err != nil {
		return nil, err
	}
	response := &v1.CreateGptMessageRes{}
	copier.Copy(response, reply)
	return response, nil
}

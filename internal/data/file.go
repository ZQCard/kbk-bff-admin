package data

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	fileV1 "github.com/ZQCard/kbk-file/api/file/v1"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/jinzhu/copier"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func NewFileClient(sr *conf.Endpoint, r registry.Discovery, tp *tracesdk.TracerProvider) fileV1.FileServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.File),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
			// 元信息
			metadata.Client(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := fileV1.NewFileServiceClient(conn)
	return c
}

type FileRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewFileRepo(data *Data, logger log.Logger) *FileRepo {
	return &FileRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/file")),
		sg:   &singleflight.Group{},
	}
}

func (rp FileRepo) GetOssStsToken(ctx context.Context) (*v1.OssStsTokenResponse, error) {
	res, err := rp.data.fileClient.GetOssStsToken(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	pb := &v1.OssStsTokenResponse{}
	if err := copier.Copy(pb, res); err != nil {
		return nil, v1.ErrorSystemError("获取OSS STS TOKEN失败").WithCause(err)
	}
	return pb, nil
}

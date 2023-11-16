package data

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/bff-admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	ossV1 "github.com/ZQCard/kbk-oss/api/oss/v1"

	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/jinzhu/copier"
)

func NewOssClient(sr *conf.Endpoint) ossV1.OSSClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Oss),
		grpc.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := ossV1.NewOSSClient(conn)
	return c
}

type OssRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewOssRepo(data *Data, logger log.Logger) *OssRepo {
	return &OssRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/file")),
		sg:   &singleflight.Group{},
	}
}

func (rp OssRepo) GetOssStsToken(ctx context.Context) (*v1.OssStsTokenResponse, error) {
	res, err := rp.data.ossClient.GetOssStsToken(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	pb := &v1.OssStsTokenResponse{}
	if err := copier.Copy(pb, res); err != nil {
		return nil, v1.ErrorSystemError("获取OSS STS TOKEN失败").WithCause(err)
	}
	return pb, nil
}

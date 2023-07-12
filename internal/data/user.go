package data

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	userV1 "github.com/ZQCard/kbk-user/api/user/v1"
	"github.com/jinzhu/copier"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/singleflight"
)

func NewUserServiceClient(sr *conf.Endpoint, r registry.Discovery, tp *tracesdk.TracerProvider) userV1.UserServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.User),
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
	c := userV1.NewUserServiceClient(conn)
	return c
}

type UserRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewUserRepo(data *Data, logger log.Logger) *UserRepo {
	return &UserRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/user")),
		sg:   &singleflight.Group{},
	}
}

func (rp UserRepo) ListUser(ctx context.Context, req *v1.GetUserListReq) (*v1.GetUserListPageRes, error) {
	list := []*v1.User{}
	reply, err := rp.data.userClient.GetUserList(ctx, &userV1.GetUserListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Mobile:   req.Mobile,
		Username: req.Username,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range reply.List {
		tmp := &v1.User{}
		copier.Copy(tmp, v)
		list = append(list, tmp)
	}
	response := &v1.GetUserListPageRes{}
	response.Total = reply.Total
	response.List = list
	return response, nil
}

func (rp UserRepo) GetUser(ctx context.Context, id int64) (*v1.User, error) {
	reply, err := rp.data.userClient.GetUser(ctx, &userV1.UserIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.User{}
	copier.Copy(res, reply)
	return res, nil
}

func (rp UserRepo) CreateUser(ctx context.Context, req *userV1.CreateUserReq) (*v1.User, error) {
	// 多个角色,创建管理员的时候只输入一个角色，其他角色放入权限服务
	reply, err := rp.data.userClient.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &v1.User{}
	copier.Copy(res, reply)
	return res, nil
}

func (rp UserRepo) UpdateUser(ctx context.Context, req *userV1.UpdateUserReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.userClient.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{
		Success: reply.Success,
	}, nil
}

func (rp UserRepo) DeleteUser(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	reply, err := rp.data.userClient.DeleteUser(ctx, &userV1.UserIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}

	return res, nil
}

func (rp UserRepo) RecoverUser(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	reply, err := rp.data.userClient.RecoverUser(ctx, &userV1.UserIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp UserRepo) ForbidUser(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	reply, err := rp.data.userClient.UserStatusChange(ctx, &userV1.UserStatusChangeReq{
		Id:     id,
		Status: false,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp UserRepo) ApproveUser(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	reply, err := rp.data.userClient.UserStatusChange(ctx, &userV1.UserStatusChangeReq{
		Id:     id,
		Status: true,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

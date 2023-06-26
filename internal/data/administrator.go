package data

import (
	"context"

	administratorV1 "github.com/ZQCard/kratos-base-kit/kbk-administrator/api/administrator/v1"
	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-bff-admin/internal/conf"
	"github.com/golang-jwt/jwt/v4"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/singleflight"
)

func NewAdministratorServiceClient(sr *conf.Endpoint, r registry.Discovery, tp *tracesdk.TracerProvider) administratorV1.AdministratorServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Administrator),
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
	c := administratorV1.NewAdministratorServiceClient(conn)
	return c
}

type AdministratorRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewAdministratorRepo(data *Data, logger log.Logger) *AdministratorRepo {
	return &AdministratorRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/administrator")),
		sg:   &singleflight.Group{},
	}
}

const (
	firstSuperAdminId = 1
)

func administratorServiceToApi(info *administratorV1.Administrator) *v1.Administrator {
	return &v1.Administrator{
		Id:            info.Id,
		Username:      info.Username,
		Nickname:      info.Nickname,
		Mobile:        info.Mobile,
		Status:        info.Status,
		Avatar:        info.Avatar,
		Role:          info.Role,
		LastLoginTime: info.LastLoginTime,
		LastLoginIp:   info.LastLoginIp,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.CreatedAt,
		DeletedAt:     info.DeletedAt,
	}
}

func (rp AdministratorRepo) FindLoginAdministratorByUsername(ctx context.Context, username string) (*administratorV1.Administrator, error) {
	reply, err := rp.data.administratorClient.GetAdministrator(ctx, &administratorV1.GetAdministratorReq{
		Username: username,
	})
	return reply, err
}

func (rp AdministratorRepo) VerifyPassword(ctx context.Context, id int64, password string) error {
	reply, err := rp.data.administratorClient.VerifyAdministratorPassword(ctx, &administratorV1.VerifyAdministratorPasswordReq{
		Id:       id,
		Password: password,
	})
	if err != nil {
		return err
	}

	if !reply.Success {
		return errors.BadRequest("PASSWORD ERR", "密码错误")
	}
	return nil
}

func (rp AdministratorRepo) AdministratorLoginSuccess(ctx context.Context, id int64, ip string, time string) error {
	reply, err := rp.data.administratorClient.AdministratorLoginSuccess(ctx, &administratorV1.AdministratorLoginSuccessReq{
		Id:            id,
		LastLoginIp:   ip,
		LastLoginTime: time,
	})
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.InternalServer("SYSTEM ERR", "更新登录信息失败")
	}
	return nil
}

func (rp AdministratorRepo) GenerateAdministratorToken(ctx context.Context, administrator *administratorV1.Administrator) (string, error) {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   administrator.Id,
			"username": administrator.Username,
			"role":     administrator.Role,
		})
	signedString, _ := claims.SignedString([]byte(GetAuthApiKey()))
	return signedString, nil
}

func (rp AdministratorRepo) ListAdministrator(ctx context.Context, req *v1.GetAdministratorListReq) (*v1.GetAdministratorListPageRes, error) {
	list := []*v1.Administrator{}
	reply, err := rp.data.administratorClient.GetAdministratorList(ctx, &administratorV1.GetAdministratorListReq{
		Page:           req.Page,
		PageSize:       req.PageSize,
		Mobile:         req.Mobile,
		Username:       req.Username,
		Nickname:       req.Nickname,
		Status:         req.Status,
		CreatedAtStart: req.CreatedAtStart,
		CreatedAtEnd:   req.CreatedAtEnd,
	})
	if err != nil {
		return nil, err
	}

	for _, v := range reply.List {
		tmp := administratorServiceToApi(v)
		list = append(list, tmp)
	}

	response := &v1.GetAdministratorListPageRes{}
	response.Total = reply.Total
	response.List = list
	return response, nil
}

func (rp AdministratorRepo) GetAdministrator(ctx context.Context, id int64) (*v1.Administrator, error) {
	reply, err := rp.data.administratorClient.GetAdministrator(ctx, &administratorV1.GetAdministratorReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Administrator{
		Id:            reply.Id,
		Username:      reply.Username,
		Mobile:        reply.Mobile,
		Nickname:      reply.Nickname,
		Avatar:        reply.Avatar,
		Status:        reply.Status,
		Role:          reply.Role,
		LastLoginTime: reply.LastLoginTime,
		LastLoginIp:   reply.LastLoginIp,
		CreatedAt:     reply.CreatedAt,
		UpdatedAt:     reply.UpdatedAt,
		DeletedAt:     reply.DeletedAt,
	}, nil
}

func (rp AdministratorRepo) CreateAdministrator(ctx context.Context, req *administratorV1.CreateAdministratorReq) (*v1.Administrator, error) {
	// 多个角色,创建管理员的时候只输入一个角色，其他角色放入权限服务
	reply, err := rp.data.administratorClient.CreateAdministrator(ctx, req)
	if err != nil {
		return nil, err
	}
	return administratorServiceToApi(reply), err
}

func (rp AdministratorRepo) UpdateAdministrator(ctx context.Context, req *administratorV1.UpdateAdministratorReq) (*v1.CheckResponse, error) {

	reply, err := rp.data.administratorClient.UpdateAdministrator(ctx, req)

	if err != nil {
		return nil, err
	}

	return &v1.CheckResponse{
		Success: reply.Success,
	}, nil
}

func (rp AdministratorRepo) DeleteAdministrator(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	reply, err := rp.data.administratorClient.DeleteAdministrator(ctx, &administratorV1.DeleteAdministratorReq{
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

func (rp AdministratorRepo) RecoverAdministrator(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	reply, err := rp.data.administratorClient.RecoverAdministrator(ctx, &administratorV1.RecoverAdministratorReq{
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

func (rp AdministratorRepo) ForbidAdministrator(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	reply, err := rp.data.administratorClient.AdministratorStatusChange(ctx, &administratorV1.AdministratorStatusChangeReq{
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

func (rp AdministratorRepo) ApproveAdministrator(ctx context.Context, id int64) (*v1.CheckResponse, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	reply, err := rp.data.administratorClient.AdministratorStatusChange(ctx, &administratorV1.AdministratorStatusChangeReq{
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

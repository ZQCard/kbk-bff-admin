package data

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/bff-admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"

	administratorV1 "github.com/ZQCard/kbk-administrator/api/administrator/v1"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"golang.org/x/sync/singleflight"
)

func NewAdministratorServiceClient(sr *conf.Endpoint) administratorV1.AdministratorServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Administrator),
		grpc.WithMiddleware(
			recovery.Recovery(),
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

func (rp AdministratorRepo) FindLoginAdministratorByUsername(ctx context.Context, username string) (*administratorV1.Administrator, error) {
	reply, err := rp.data.administratorClient.GetAdministrator(ctx, &administratorV1.GetAdministratorReq{
		Username: username,
	})
	return reply, err
}

func (rp AdministratorRepo) VerifyPassword(ctx context.Context, id int64, password string) error {
	_, err := rp.data.administratorClient.VerifyAdministratorPassword(ctx, &administratorV1.VerifyAdministratorPasswordReq{
		Id:       id,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (rp AdministratorRepo) AdministratorLoginSuccess(ctx context.Context, id int64, ip string, time string) error {
	_, err := rp.data.administratorClient.AdministratorLoginSuccess(ctx, &administratorV1.AdministratorLoginSuccessReq{
		Id:            id,
		LastLoginIp:   ip,
		LastLoginTime: time,
	})
	return err
}

func (rp AdministratorRepo) GenerateAdministratorToken(ctx context.Context, administrator *administratorV1.Administrator) (string, error) {
	claims := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   administrator.Id,
			"username": administrator.Username,
			"role":     administrator.Role,
		})
	signedString, _ := claims.SignedString([]byte(rp.data.cfg.Jwt.Key))
	return signedString, nil
}

func (rp AdministratorRepo) ListAdministrator(ctx context.Context, req *v1.GetAdministratorListReq) (*v1.GetAdministratorListPageRes, error) {
	list := []*v1.Administrator{}
	reqData := &administratorV1.GetAdministratorListReq{}
	copier.Copy(reqData, req)
	reply, err := rp.data.administratorClient.GetAdministratorList(ctx, reqData)
	if err != nil {
		return nil, err
	}

	for _, v := range reply.List {
		tmp := &v1.Administrator{}
		copier.Copy(tmp, v)
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
	res := &v1.Administrator{}
	copier.Copy(res, reply)
	return res, nil
}

func (rp AdministratorRepo) CreateAdministrator(ctx context.Context, req *administratorV1.CreateAdministratorReq) (*v1.Administrator, error) {
	// 多个角色,创建管理员的时候只输入一个角色，其他角色放入权限服务
	reply, err := rp.data.administratorClient.CreateAdministrator(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &v1.Administrator{}
	copier.Copy(res, reply)
	return res, err
}

func (rp AdministratorRepo) UpdateAdministrator(ctx context.Context, req *administratorV1.UpdateAdministratorReq) (*emptypb.Empty, error) {
	return rp.data.administratorClient.UpdateAdministrator(ctx, req)
}

func (rp AdministratorRepo) DeleteAdministrator(ctx context.Context, id int64) (*emptypb.Empty, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	return rp.data.administratorClient.DeleteAdministrator(ctx, &administratorV1.DeleteAdministratorReq{
		Id: id,
	})
}

func (rp AdministratorRepo) RecoverAdministrator(ctx context.Context, id int64) (*emptypb.Empty, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	return rp.data.administratorClient.RecoverAdministrator(ctx, &administratorV1.RecoverAdministratorReq{
		Id: id,
	})
}

func (rp AdministratorRepo) ForbidAdministrator(ctx context.Context, id int64) (*emptypb.Empty, error) {
	if id == firstSuperAdminId {
		return nil, v1.ErrorSystemError("初始超管无法变更")
	}
	return rp.data.administratorClient.AdministratorStatusChange(ctx, &administratorV1.AdministratorStatusChangeReq{
		Id:     id,
		Status: false,
	})
}

func (rp AdministratorRepo) ApproveAdministrator(ctx context.Context, id int64) (*emptypb.Empty, error) {
	return rp.data.administratorClient.AdministratorStatusChange(ctx, &administratorV1.AdministratorStatusChangeReq{
		Id:     id,
		Status: true,
	})
}

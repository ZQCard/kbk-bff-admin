package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	administratorV1 "github.com/ZQCard/kratos-base-kit/kbk-administrator/api/administrator/v1"
	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-bff-admin/pkg/utils/timeHelper"
)

func (s *AdminInterface) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginRes, error) {
	administrator, err := s.administratorRepo.FindLoginAdministratorByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	// 如果是冻结状态则无法登录
	if !administrator.Status {
		return nil, errors.BadRequest("BAD REQUEST", "用户已被冻结")
	}

	// 验证密码
	err = s.administratorRepo.VerifyPassword(ctx, administrator.Id, req.Password)
	if err != nil {
		return nil, err
	}

	// 更新登录信息
	if tr, ok := transport.FromServerContext(ctx); ok {
		// 将请求信息放入ctx中
		if ht, ok := tr.(*http.Transport); ok {
			now := timeHelper.CurrentTimeYMDHIS()
			ip := ht.Request().RemoteAddr
			err = s.administratorRepo.AdministratorLoginSuccess(ctx, administrator.Id, ip, now)
			if err != nil {
				return nil, err
			}
			administrator.LastLoginIp = ip
			administrator.LastLoginTime = timeHelper.CurrentTimeYMDHIS()
		}

	}

	// 生成token
	token, err := s.administratorRepo.GenerateAdministratorToken(ctx, administrator)
	if err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		Token: token,
	}, nil
}

func (s *AdminInterface) LoginOut(ctx context.Context, req *emptypb.Empty) (*v1.CheckResponse, error) {
	return &v1.CheckResponse{
		Success: true,
	}, nil
}

func (s *AdminInterface) GetAdministratorList(ctx context.Context, req *v1.GetAdministratorListReq) (*v1.GetAdministratorListPageRes, error) {
	return s.administratorRepo.ListAdministrator(ctx, req)
}

func (s *AdminInterface) GetAdministrator(ctx context.Context, req *v1.IdReq) (*v1.Administrator, error) {
	return s.administratorRepo.GetAdministrator(ctx, req.Id)
}

func (s *AdminInterface) GetAdministratorInfo(ctx context.Context, req *emptypb.Empty) (*v1.Administrator, error) {
	// 获取当前登录userId
	userId := ctx.Value("x-md-global-userId").(int64)
	return s.administratorRepo.GetAdministrator(ctx, userId)
}

func (s *AdminInterface) CreateAdministrator(ctx context.Context, req *v1.CreateAdministratorReq) (*v1.Administrator, error) {
	administratorReq := &administratorV1.CreateAdministratorReq{
		Username: req.Username,
		Password: req.Password,
		Mobile:   req.Mobile,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Role:     req.Role[0],
	}
	administrator, err := s.administratorRepo.CreateAdministrator(ctx, administratorReq)
	if err != nil {
		return nil, err
	}
	// 多个角色放入权限服务中
	res, err := s.authorizationRepo.SetRolesForUser(ctx, &v1.SetRolesForUserReq{
		Username: administrator.Username,
		Roles:    req.Role,
	})
	// 授权失败， 删除用户
	if err != nil || !res.Success {
		res, err = s.administratorRepo.DeleteAdministrator(ctx, administrator.Id)
		if err != nil || !res.Success {
			return nil, errors.InternalServer("SYSTEM ERROR", err.Error())
		}
	}
	return administrator, nil
}

func (s *AdminInterface) UpdateAdministrator(ctx context.Context, req *v1.UpdateAdministratorReq) (*v1.CheckResponse, error) {
	administratorReq := &administratorV1.UpdateAdministratorReq{
		Id:       req.Id,
		Username: req.Username,
		Password: req.Password,
		Mobile:   req.Mobile,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Role:     req.Role[0],
	}
	reply, err := s.administratorRepo.UpdateAdministrator(ctx, administratorReq)
	if err != nil {
		return nil, err
	}
	// 多个角色放入权限服务中
	s.authorizationRepo.SetRolesForUser(ctx, &v1.SetRolesForUserReq{
		Username: req.Username,
		Roles:    req.Role,
	})
	return reply, nil
}

func (s *AdminInterface) DeleteAdministrator(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.administratorRepo.DeleteAdministrator(ctx, req.Id)
}

func (s *AdminInterface) RecoverAdministrator(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.administratorRepo.RecoverAdministrator(ctx, req.Id)
}

func (s *AdminInterface) ForbidAdministrator(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.administratorRepo.ForbidAdministrator(ctx, req.Id)
}

func (s *AdminInterface) ApproveAdministrator(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.administratorRepo.ApproveAdministrator(ctx, req.Id)
}

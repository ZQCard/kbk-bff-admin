package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	v1 "github.com/ZQCard/kbk-bff-admin/api/admin/v1"
	userV1 "github.com/ZQCard/kbk-user/api/user/v1"
	"github.com/jinzhu/copier"
)

func (s *AdminInterface) GetUserList(ctx context.Context, req *v1.GetUserListReq) (*v1.GetUserListPageRes, error) {
	return s.userRepo.ListUser(ctx, req)
}

func (s *AdminInterface) GetUser(ctx context.Context, req *v1.IdReq) (*v1.User, error) {
	return s.userRepo.GetUser(ctx, req.Id)
}

func (s *AdminInterface) GetUserInfo(ctx context.Context, req *emptypb.Empty) (*v1.User, error) {
	// 获取当前登录userId
	userId := ctx.Value("x-md-global-userId").(int64)
	return s.userRepo.GetUser(ctx, userId)
}

func (s *AdminInterface) CreateUser(ctx context.Context, req *v1.CreateUserReq) (*v1.User, error) {
	userReq := &userV1.CreateUserReq{}
	copier.Copy(userReq, req)
	user, err := s.userRepo.CreateUser(ctx, userReq)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminInterface) UpdateUser(ctx context.Context, req *v1.UpdateUserReq) (*v1.CheckResponse, error) {
	userReq := &userV1.UpdateUserReq{}
	copier.Copy(userReq, req)
	reply, err := s.userRepo.UpdateUser(ctx, userReq)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (s *AdminInterface) DeleteUser(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.userRepo.DeleteUser(ctx, req.Id)
}

func (s *AdminInterface) RecoverUser(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.userRepo.RecoverUser(ctx, req.Id)
}

func (s *AdminInterface) ForbidUser(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.userRepo.ForbidUser(ctx, req.Id)
}

func (s *AdminInterface) ApproveUser(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.userRepo.ApproveUser(ctx, req.Id)
}

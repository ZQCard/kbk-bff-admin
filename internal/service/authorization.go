package service

import (
	"context"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AdminInterface) GetRoleAll(ctx context.Context, req *emptypb.Empty) (*v1.GetRoleAllRes, error) {
	return s.authorizationRepo.GetRoleAll(ctx)
}

func (s *AdminInterface) CreateRole(ctx context.Context, req *v1.CreateRoleReq) (*v1.Role, error) {
	return s.authorizationRepo.CreateRole(ctx, req)
}
func (s *AdminInterface) UpdateRole(ctx context.Context, req *v1.UpdateRoleReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.UpdateRole(ctx, req)
}
func (s *AdminInterface) DeleteRole(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.DeleteRole(ctx, req)
}
func (s *AdminInterface) SetRolesForUser(ctx context.Context, req *v1.SetRolesForUserReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.SetRolesForUser(ctx, req)
}
func (s *AdminInterface) GetRolesForUser(ctx context.Context, req *v1.GetRolesForUserReq) (*v1.GetRolesForUserRes, error) {
	return s.authorizationRepo.GetRolesForUser(ctx, req)
}
func (s *AdminInterface) GetUsersForRole(ctx context.Context, req *v1.RoleNameReq) (*v1.GetUsersForRoleRes, error) {
	return s.authorizationRepo.GetUsersForRole(ctx, req)
}

func (s *AdminInterface) DeleteRoleForUser(ctx context.Context, req *v1.DeleteRoleForUserReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.DeleteRoleForUser(ctx, req)
}

func (s *AdminInterface) DeleteRolesForUser(ctx context.Context, req *v1.DeleteRolesForUserReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.DeleteRolesForUser(ctx, req)
}

func (s *AdminInterface) GetPolicies(ctx context.Context, req *v1.RoleNameReq) (*v1.GetPoliciesRes, error) {
	return s.authorizationRepo.GetPolicies(ctx, req)
}
func (s *AdminInterface) UpdatePolicies(ctx context.Context, req *v1.UpdatePoliciesReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.UpdatePolicies(ctx, req)
}

func (s *AdminInterface) GetApiAll(ctx context.Context, req *emptypb.Empty) (*v1.GetApiAllRes, error) {
	return s.authorizationRepo.GetApiAll(ctx)
}
func (s *AdminInterface) GetApiList(ctx context.Context, req *v1.GetApiListReq) (*v1.GetApiListRes, error) {
	return s.authorizationRepo.GetApiList(ctx, req)
}
func (s *AdminInterface) CreateApi(ctx context.Context, req *v1.CreateApiReq) (*v1.Api, error) {
	return s.authorizationRepo.CreateApi(ctx, req)
}
func (s *AdminInterface) UpdateApi(ctx context.Context, req *v1.UpdateApiReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.UpdateApi(ctx, req)
}
func (s *AdminInterface) DeleteApi(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.DeleteApi(ctx, req)
}
func (s *AdminInterface) GetMenuAll(ctx context.Context, req *emptypb.Empty) (*v1.GetMenuTreeRes, error) {
	return s.authorizationRepo.GetMenuAll(ctx)
}
func (s *AdminInterface) GetMenuTree(ctx context.Context, req *emptypb.Empty) (*v1.GetMenuTreeRes, error) {
	return s.authorizationRepo.GetMenuTree(ctx)
}
func (s *AdminInterface) CreateMenu(ctx context.Context, req *v1.CreateMenuReq) (*v1.Menu, error) {
	return s.authorizationRepo.CreateMenu(ctx, req)
}
func (s *AdminInterface) UpdateMenu(ctx context.Context, req *v1.UpdateMenuReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.UpdateMenu(ctx, req)
}
func (s *AdminInterface) DeleteMenu(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.DeleteMenu(ctx, req)
}
func (s *AdminInterface) GetRoleMenuTree(ctx context.Context, req *v1.RoleNameReq) (*v1.GetMenuTreeRes, error) {
	return s.authorizationRepo.GetRoleMenuTree(ctx, req)
}
func (s *AdminInterface) GetRoleMenu(ctx context.Context, req *v1.RoleNameReq) (*v1.GetMenuTreeRes, error) {
	return s.authorizationRepo.GetRoleMenu(ctx, req)
}
func (s *AdminInterface) SetRoleMenu(ctx context.Context, req *v1.SetRoleMenuReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.SetRoleMenu(ctx, req)
}

func (s *AdminInterface) GetRoleMenuBtn(ctx context.Context, req *v1.GetRoleMenuBtnReq) (*v1.GetRoleMenuBtnRes, error) {
	return s.authorizationRepo.GetRoleMenuBtn(ctx, req)
}

func (s *AdminInterface) SetRoleMenuBtn(ctx context.Context, req *v1.SetRoleMenuBtnReq) (*v1.CheckResponse, error) {
	return s.authorizationRepo.SetRoleMenuBtn(ctx, req)
}

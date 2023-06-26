package data

import (
	"context"

	authorizationV1 "github.com/ZQCard/kbk-authorization/api/authorization/v1"
	v1 "github.com/ZQCard/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kbk-bff-admin/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type AuthorizationRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewAuthorizationRepo(data *Data, logger log.Logger) *AuthorizationRepo {
	return &AuthorizationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/authorization")),
		sg:   &singleflight.Group{},
	}
}

func NewAuthorizationServiceClient(ac *conf.Auth, sr *conf.Endpoint, r registry.Discovery, tp *tracesdk.TracerProvider) authorizationV1.AuthorizationServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(sr.Authorization),
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
	return authorizationV1.NewAuthorizationServiceClient(conn)
}

func (rp AuthorizationRepo) GetRoleAll(ctx context.Context) (*v1.GetRoleAllRes, error) {
	reply, err := rp.data.authorizationClient.GetRoleAll(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, err
	}
	roles := []*v1.Role{}
	for _, v := range reply.List {
		roles = append(roles, &v1.Role{
			Id:        v.Id,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	res := &v1.GetRoleAllRes{
		List: roles,
	}
	return res, err
}

func (rp AuthorizationRepo) CreateRole(ctx context.Context, req *v1.CreateRoleReq) (*v1.Role, error) {
	reply, err := rp.data.authorizationClient.CreateRole(ctx, &authorizationV1.CreateRoleReq{
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.Role{
		Id:        reply.Id,
		Name:      reply.Name,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
	return res, nil
}

func (rp AuthorizationRepo) UpdateRole(ctx context.Context, req *v1.UpdateRoleReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.UpdateRole(ctx, &authorizationV1.UpdateRoleReq{
		Id:   req.Id,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) DeleteRole(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.DeleteRole(ctx, &authorizationV1.DeleteRoleReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) SetRolesForUser(ctx context.Context, req *v1.SetRolesForUserReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.SetRolesForUser(ctx, &authorizationV1.SetRolesForUserReq{
		Username: req.Username,
		Roles:    req.Roles,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetRolesForUser(ctx context.Context, req *v1.GetRolesForUserReq) (*v1.GetRolesForUserRes, error) {
	reply, err := rp.data.authorizationClient.GetRolesForUser(ctx, &authorizationV1.GetRolesForUserReq{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.GetRolesForUserRes{
		Roles: reply.Roles,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetUsersForRole(ctx context.Context, req *v1.RoleNameReq) (*v1.GetUsersForRoleRes, error) {
	reply, err := rp.data.authorizationClient.GetUsersForRole(ctx, &authorizationV1.RoleNameReq{
		Role: req.Role,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.GetUsersForRoleRes{
		Users: reply.Users,
	}
	return res, nil
}

func (rp AuthorizationRepo) DeleteRoleForUser(ctx context.Context, req *v1.DeleteRoleForUserReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.DeleteRoleForUser(ctx, &authorizationV1.DeleteRoleForUserReq{
		Username: req.Username,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) DeleteRolesForUser(ctx context.Context, req *v1.DeleteRolesForUserReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.DeleteRolesForUser(ctx, &authorizationV1.DeleteRolesForUserReq{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetPolicies(ctx context.Context, req *v1.RoleNameReq) (*v1.GetPoliciesRes, error) {
	reply, err := rp.data.authorizationClient.GetPolicies(ctx, &authorizationV1.RoleNameReq{
		Role: req.Role,
	})
	if err != nil {
		return nil, err
	}
	var policyRules []*v1.PolicyRules
	for _, v := range reply.PolicyRules {
		policyRules = append(policyRules, &v1.PolicyRules{
			Path:   v.Path,
			Method: v.Method,
		})
	}
	res := &v1.GetPoliciesRes{
		PolicyRules: policyRules,
	}
	return res, nil
}

func (rp AuthorizationRepo) UpdatePolicies(ctx context.Context, req *v1.UpdatePoliciesReq) (*v1.CheckResponse, error) {
	var policyRules []*authorizationV1.PolicyRules
	for _, v := range req.PolicyRules {
		policyRules = append(policyRules, &authorizationV1.PolicyRules{
			Path:   v.Path,
			Method: v.Method,
		})
	}

	reply, err := rp.data.authorizationClient.UpdatePolicies(ctx, &authorizationV1.UpdatePoliciesReq{
		Role:        req.Role,
		PolicyRules: policyRules,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetApiAll(ctx context.Context) (*v1.GetApiAllRes, error) {
	reply, err := rp.data.authorizationClient.GetApiListAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var list []*v1.Api
	for _, v := range reply.List {
		list = append(list, &v1.Api{
			Id:        v.Id,
			Group:     v.Group,
			Name:      v.Name,
			Path:      v.Path,
			Method:    v.Method,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	res := &v1.GetApiAllRes{
		List: list,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetApiList(ctx context.Context, req *v1.GetApiListReq) (*v1.GetApiListRes, error) {
	reply, err := rp.data.authorizationClient.GetApiList(ctx, &authorizationV1.GetApiListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Group:    req.Group,
		Method:   req.Method,
		Path:     req.Path,
	})
	if err != nil {
		return nil, err
	}
	var list []*v1.Api
	for _, v := range reply.List {
		list = append(list, &v1.Api{
			Id:        v.Id,
			Group:     v.Group,
			Name:      v.Name,
			Path:      v.Path,
			Method:    v.Method,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	res := &v1.GetApiListRes{
		List:  list,
		Total: reply.Total,
	}
	return res, nil
}

func (rp AuthorizationRepo) CreateApi(ctx context.Context, req *v1.CreateApiReq) (*v1.Api, error) {
	reply, err := rp.data.authorizationClient.CreateApi(ctx, &authorizationV1.CreateApiReq{
		Group:  req.Group,
		Name:   req.Name,
		Path:   req.Path,
		Method: req.Method,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.Api{
		Id:        reply.Id,
		Group:     reply.Group,
		Name:      reply.Name,
		Path:      reply.Path,
		Method:    reply.Method,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}
	return res, nil
}

func (rp AuthorizationRepo) UpdateApi(ctx context.Context, req *v1.UpdateApiReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.UpdateApi(ctx, &authorizationV1.UpdateApiReq{
		Id:     req.Id,
		Group:  req.Group,
		Name:   req.Name,
		Path:   req.Path,
		Method: req.Method,
	})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	response := &v1.CheckResponse{
		Success: reply.Success,
	}
	return response, nil
}

func (rp AuthorizationRepo) DeleteApi(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.DeleteApi(ctx, &authorizationV1.DeleteApiReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetMenuAll(ctx context.Context) (*v1.GetMenuTreeRes, error) {
	reply, err := rp.data.authorizationClient.GetMenuAll(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var list []*v1.Menu
	for _, v := range reply.List {
		var btns []*v1.MenuBtn
		for _, btn := range v.MenuBtns {
			btns = append(btns, &v1.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   btn.CreatedAt,
				UpdatedAt:   btn.UpdatedAt,
			})
		}

		list = append(list, &v1.Menu{
			Id:        v.Id,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Path:      v.Path,
			Name:      v.Name,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			MenuBtns:  btns,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	res := &v1.GetMenuTreeRes{
		List: list,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetMenuTree(ctx context.Context) (*v1.GetMenuTreeRes, error) {
	reply, err := rp.data.authorizationClient.GetMenuTree(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	var list []*v1.Menu
	menu := reply.List
	for k, v := range menu {
		children := findChildrenMenu(menu[k])

		var btns []*v1.MenuBtn
		for _, btn := range v.MenuBtns {
			btns = append(btns, &v1.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   btn.CreatedAt,
				UpdatedAt:   btn.UpdatedAt,
			})
		}

		list = append(list, &v1.Menu{
			Id:        v.Id,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Path:      v.Path,
			Name:      v.Name,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			MenuBtns:  btns,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Children:  children,
		})

	}

	res := &v1.GetMenuTreeRes{
		List: list,
	}
	return res, nil
}

func findChildrenMenu(menu *authorizationV1.Menu) []*v1.Menu {
	children := []*v1.Menu{}
	if len(menu.Children) != 0 {
		for k := range menu.Children {
			var btns []*v1.MenuBtn
			for _, btn := range menu.Children[k].MenuBtns {
				btns = append(btns, &v1.MenuBtn{
					Id:          btn.Id,
					MenuId:      btn.MenuId,
					Name:        btn.Name,
					Description: btn.Description,
					Identifier:  btn.Identifier,
					CreatedAt:   btn.CreatedAt,
					UpdatedAt:   btn.UpdatedAt,
				})
			}

			children = append(children, &v1.Menu{
				Id:        menu.Children[k].Id,
				Name:      menu.Children[k].Name,
				Path:      menu.Children[k].Path,
				ParentId:  menu.Children[k].ParentId,
				ParentIds: menu.Children[k].ParentIds,
				Hidden:    menu.Children[k].Hidden,
				Component: menu.Children[k].Component,
				Sort:      menu.Children[k].Sort,
				Title:     menu.Children[k].Title,
				Icon:      menu.Children[k].Icon,
				CreatedAt: menu.Children[k].CreatedAt,
				UpdatedAt: menu.Children[k].UpdatedAt,
				MenuBtns:  btns,
				Children:  findChildrenMenu(menu.Children[k]),
			})
		}
	}
	return children
}

func (rp AuthorizationRepo) CreateMenu(ctx context.Context, req *v1.CreateMenuReq) (*v1.Menu, error) {
	var btns []*authorizationV1.MenuBtn
	for _, btn := range req.MenuBtns {
		btns = append(btns, &authorizationV1.MenuBtn{
			Id:          btn.Id,
			MenuId:      btn.MenuId,
			Name:        btn.Name,
			Description: btn.Description,
			Identifier:  btn.Identifier,
		})
	}

	reply, err := rp.data.authorizationClient.CreateMenu(ctx, &authorizationV1.CreateMenuReq{
		Name:      req.Name,
		Path:      req.Path,
		ParentId:  req.ParentId,
		ParentIds: req.ParentIds,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Title:     req.Title,
		Icon:      req.Icon,
		MenuBtns:  btns,
	})
	if err != nil {
		return nil, err
	}
	var btns2 []*v1.MenuBtn
	for _, btn := range reply.MenuBtns {
		btns2 = append(btns2, &v1.MenuBtn{
			Id:          btn.Id,
			MenuId:      btn.MenuId,
			Name:        btn.Name,
			Description: btn.Description,
			Identifier:  btn.Identifier,
			CreatedAt:   btn.CreatedAt,
			UpdatedAt:   btn.UpdatedAt,
		})
	}

	res := &v1.Menu{
		Id:        reply.Id,
		Name:      req.Name,
		Path:      req.Path,
		ParentId:  req.ParentId,
		ParentIds: req.ParentIds,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Title:     req.Title,
		Icon:      req.Icon,
		MenuBtns:  btns2,
	}
	return res, nil
}

func (rp AuthorizationRepo) UpdateMenu(ctx context.Context, req *v1.UpdateMenuReq) (*v1.CheckResponse, error) {
	var btns []*authorizationV1.MenuBtn
	for _, btn := range req.MenuBtns {
		btns = append(btns, &authorizationV1.MenuBtn{
			Id:          btn.Id,
			MenuId:      btn.MenuId,
			Name:        btn.Name,
			Description: btn.Description,
			Identifier:  btn.Identifier,
		})
	}
	reply, err := rp.data.authorizationClient.UpdateMenu(ctx, &authorizationV1.UpdateMenuReq{
		Id:        req.Id,
		Name:      req.Name,
		Path:      req.Path,
		ParentId:  req.ParentId,
		ParentIds: req.ParentIds,
		Hidden:    req.Hidden,
		Component: req.Component,
		Sort:      req.Sort,
		Title:     req.Title,
		Icon:      req.Icon,
		MenuBtns:  btns,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) DeleteMenu(ctx context.Context, req *v1.IdReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.DeleteMenu(ctx, &authorizationV1.IdReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetRoleMenuTree(ctx context.Context, req *v1.RoleNameReq) (*v1.GetMenuTreeRes, error) {
	reply, err := rp.data.authorizationClient.GetRoleMenuTree(ctx, &authorizationV1.RoleNameReq{
		Role: req.Role,
	})
	if err != nil {
		return nil, err
	}
	var list []*v1.Menu
	menu := reply.List
	for k, v := range menu {
		children := findChildrenMenu(menu[k])

		var btns []*v1.MenuBtn
		for _, btn := range v.MenuBtns {
			btns = append(btns, &v1.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   btn.CreatedAt,
				UpdatedAt:   btn.UpdatedAt,
			})
		}

		list = append(list, &v1.Menu{
			Id:        v.Id,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Path:      v.Path,
			Name:      v.Name,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			MenuBtns:  btns,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Children:  children,
		})

	}

	res := &v1.GetMenuTreeRes{
		List: list,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetRoleMenu(ctx context.Context, req *v1.RoleNameReq) (*v1.GetMenuTreeRes, error) {
	reply, err := rp.data.authorizationClient.GetRoleMenu(ctx, &authorizationV1.RoleNameReq{
		Role: req.Role,
	})
	if err != nil {
		return nil, err
	}
	var list []*v1.Menu
	menu := reply.List
	for _, v := range menu {
		var btns []*v1.MenuBtn
		for _, btn := range v.MenuBtns {
			btns = append(btns, &v1.MenuBtn{
				Id:          btn.Id,
				MenuId:      btn.MenuId,
				Name:        btn.Name,
				Description: btn.Description,
				Identifier:  btn.Identifier,
				CreatedAt:   btn.CreatedAt,
				UpdatedAt:   btn.UpdatedAt,
			})
		}

		list = append(list, &v1.Menu{
			Id:        v.Id,
			ParentId:  v.ParentId,
			ParentIds: v.ParentIds,
			Path:      v.Path,
			Name:      v.Name,
			Hidden:    v.Hidden,
			Component: v.Component,
			Sort:      v.Sort,
			Title:     v.Title,
			Icon:      v.Icon,
			MenuBtns:  btns,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})

	}

	res := &v1.GetMenuTreeRes{
		List: list,
	}
	return res, nil
}

func (rp AuthorizationRepo) SetRoleMenu(ctx context.Context, req *v1.SetRoleMenuReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.SaveRoleMenu(ctx, &authorizationV1.SaveRoleMenuReq{
		RoleId:  req.RoleId,
		MenuIds: req.MenuIds,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

func (rp AuthorizationRepo) GetRoleMenuBtn(ctx context.Context, req *v1.GetRoleMenuBtnReq) (*v1.GetRoleMenuBtnRes, error) {
	reply, err := rp.data.authorizationClient.GetRoleMenuBtn(ctx, &authorizationV1.GetRoleMenuBtnReq{
		RoleId:   req.RoleId,
		RoleName: req.RoleName,
		MenuId:   req.MenuId,
	})
	if err != nil {
		return nil, err
	}
	list := []*v1.MenuBtn{}
	for _, v := range reply.List {
		list = append(list, &v1.MenuBtn{
			Id:          v.Id,
			MenuId:      v.MenuId,
			Name:        v.Name,
			Description: v.Description,
			Identifier:  v.Identifier,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}
	res := &v1.GetRoleMenuBtnRes{
		List: list,
	}
	return res, nil
}

func (rp AuthorizationRepo) SetRoleMenuBtn(ctx context.Context, req *v1.SetRoleMenuBtnReq) (*v1.CheckResponse, error) {
	reply, err := rp.data.authorizationClient.SaveRoleMenuBtn(ctx, &authorizationV1.SaveRoleMenuBtnReq{
		RoleId:     req.RoleId,
		MenuId:     req.MenuId,
		MenuBtnIds: req.MenuBtnIds,
	})
	if err != nil {
		return nil, err
	}
	res := &v1.CheckResponse{
		Success: reply.Success,
	}
	return res, nil
}

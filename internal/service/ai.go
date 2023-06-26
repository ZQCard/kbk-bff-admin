package service

import (
	"context"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
)

func (s *AdminInterface) CreateGptMessage(ctx context.Context, req *v1.CreateGptMessageReq) (*v1.CreateGptMessageRes, error) {
	r, e := s.aiRepo.CreateGptMessage(ctx, req)
	return r, e
}

package service

import (
	"context"

	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
)

func (s *AdminInterface) GetApiLogList(ctx context.Context, req *v1.GetApiLogListReq) (*v1.GetApiLogListRes, error) {
	return s.apiLogRepo.ListApiLog(ctx, req)
}

package service

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/bff-admin/v1"
)

func (s *AdminInterface) GetApiLogList(ctx context.Context, req *v1.GetApiLogListReq) (*v1.GetApiLogListRes, error) {
	return s.apiLogRepo.ListApiLog(ctx, req)
}

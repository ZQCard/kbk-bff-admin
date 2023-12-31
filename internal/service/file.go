package service

import (
	"context"

	v1 "github.com/ZQCard/kbk-bff-admin/api/bff-admin/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AdminInterface) GetOssStsToken(ctx context.Context, req *emptypb.Empty) (*v1.OssStsTokenResponse, error) {
	return s.ossRepo.GetOssStsToken(ctx)
}

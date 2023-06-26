package service

import (
	v1 "github.com/ZQCard/kratos-base-kit/kbk-bff-admin/api/admin/v1"
	"github.com/ZQCard/kratos-base-kit/kbk-bff-admin/internal/data"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAdminInterface)

type AdminInterface struct {
	v1.UnimplementedAdminServer
	administratorRepo *data.AdministratorRepo
	authorizationRepo *data.AuthorizationRepo
	apiLogRepo        *data.ApiLogRepo
	fileRepo          *data.FileRepo
	aiRepo            *data.AIRepo
	log               *log.Helper
}

func NewAdminInterface(
	administratorRepo *data.AdministratorRepo,
	authorizationRepo *data.AuthorizationRepo,
	apiLogRepo *data.ApiLogRepo,
	fileRepo *data.FileRepo,
	logger log.Logger,
	aiRepo *data.AIRepo,
) *AdminInterface {
	return &AdminInterface{
		log:               log.NewHelper(log.With(logger, "module", "service/interface")),
		administratorRepo: administratorRepo,
		apiLogRepo:        apiLogRepo,
		fileRepo:          fileRepo,
		authorizationRepo: authorizationRepo,
		aiRepo:            aiRepo,
	}
}
